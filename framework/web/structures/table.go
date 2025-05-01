package structures

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Andision/RoundWorld/framework/dal"
	"github.com/Andision/RoundWorld/framework/models"
	"github.com/golang/glog"
	"github.com/google/uuid"
	"log"
	"net/http"
	"sync"
)

const (
	tableResponseCodeUpdate = iota
	tableResponseCodeError
	tableResponseCodeEnd
	tableResponseCodeSuccess
)

type tableResponse struct {
	TableId string      `json:"table_id"`
	Code    int         `json:"code"`
	Message string      `json:"Message"`
	Data    interface{} `json:"data"`
}

type TableMessage struct {
	Sender  Player
	Message []byte
}

type Table interface {
	GetTableId() string

	Run()
	Join(ctx context.Context, w http.ResponseWriter, r *http.Request)

	SendUnregister(player Player) bool
	SendMessage(msg TableMessage) bool
}

func NewTable(config models.Config, executor string) Table {
	return &TableImpl{
		gameConfig: config,
		gameState:  config.GetNewGameState(executor),

		tableId: uuid.New().String(),
		players: make(map[Player]bool),

		broadcastQueue:  make(chan []byte),
		messageQueue:    make(chan TableMessage),
		unregisterQueue: make(chan Player),
	}
}

// TableImpl maintains the set of active clients and broadcasts messages to the clients.
type TableImpl struct {
	// meta data
	tableId    string
	gameConfig models.Config
	gameState  models.State

	players      map[Player]bool
	playersMutex sync.RWMutex

	broadcastQueue  chan []byte
	messageQueue    chan TableMessage
	unregisterQueue chan Player
}

func (t *TableImpl) GetTableId() string {
	return t.tableId
}

func (t *TableImpl) Run() {
	for {
		select {
		case player := <-t.unregisterQueue:
			t.playersMutex.Lock()
			if _, ok := t.players[player]; ok {
				delete(t.players, player)
				player.Close()
			}
			t.playersMutex.Unlock()
		case message := <-t.broadcastQueue:
			glog.Infof("broadcastQueue: %v", string(message))
			t.playersMutex.Lock()
			for player := range t.players {
				if !player.Send(message) {
					glog.Infof("Closing send channel for player: %v", player)
					player.Close()
					delete(t.players, player)
				}
			}
			t.playersMutex.Unlock()
		case message := <-t.messageQueue:
			log.Printf("Received Message: %s", message.Message)

			ctx := context.Background()
			ctx = context.WithValue(ctx, "tableId", t.tableId)
			ctx = context.WithValue(ctx, "playerCount", len(t.players))
			ctx = context.WithValue(ctx, "players", t.players)

			t.tableMessageHandler(ctx, &message)
		default:
			//log.Printf("No Message received")
		}
	}
}

func (t *TableImpl) Join(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	t.playersMutex.Lock()
	defer t.playersMutex.Unlock()

	currentPlayerCount := len(t.players)
	if currentPlayerCount >= t.gameConfig.GetMaxPlayer() {
		http.Error(w, "TableImpl is full", http.StatusBadRequest)
		return
	}

	username := ctx.Value("username")
	if username == nil {
		http.Error(w, "Username not found", http.StatusBadRequest)
		return
	}

	user := dal.NewUserImpl(username.(string))

	player := NewPlayer(user, t, w, r)
	if player == nil {
		http.Error(w, "Failed to create player", http.StatusInternalServerError)
	}

	t.players[player] = true
	log.Printf("player %s joined", username)
	return
}

func (t *TableImpl) SendUnregister(player Player) bool {
	select {
	case t.unregisterQueue <- player:
		return true
	default:
		return false
	}
}

func (t *TableImpl) SendMessage(msg TableMessage) bool {
	select {
	case t.messageQueue <- msg:
		return true
	default:
		return false
	}
}

func (t *TableImpl) tableMessageHandler(ctx context.Context, message *TableMessage) {
	sender := message.Sender
	rawMessage := message.Message

	parse, err := t.gameState.Parse(ctx, rawMessage)
	if err != nil {
		rawResponse, err := json.Marshal(
			tableResponse{
				TableId: t.tableId,
				Code:    tableResponseCodeError,
				Message: "parse error",
				Data:    nil,
			})
		if err != nil {
			log.Printf("Failed to marshal tableResponse: %v", err)
			return
		}
		sender.Send(rawResponse)
		return
	}

	t.gameState.Lock()
	defer t.gameState.Unlock()

	valid, err := t.gameState.Validate(ctx, parse)
	if err != nil {
		rawResponse, err := json.Marshal(
			tableResponse{
				TableId: t.tableId,
				Code:    tableResponseCodeError,
				Message: fmt.Sprintf("validate error: %v", err),
				Data:    nil,
			})
		if err != nil {
			log.Printf("Failed to marshal tableResponse: %v", err)
			return
		}
		sender.Send(rawResponse)
		return
	}
	if !valid {
		rawResponse, err := json.Marshal(
			tableResponse{
				TableId: t.tableId,
				Code:    tableResponseCodeError,
				Message: "invalid action",
				Data:    nil,
			})
		if err != nil {
			log.Printf("Failed to marshal tableResponse: %v", err)
			return
		}
		sender.Send(rawResponse)
		return
	}

	err = t.gameState.Execute(ctx, parse)
	if err != nil {
		rawResponse, err := json.Marshal(
			tableResponse{
				TableId: t.tableId,
				Code:    tableResponseCodeError,
				Message: fmt.Sprintf("execute error: %v", err),
				Data:    nil,
			})
		if err != nil {
			log.Printf("Failed to marshal tableResponse: %v", err)
			return
		}
		sender.Send(rawResponse)
		return
	}

	log.Print("execute finished")

	rawState, err := t.gameState.Encode(ctx)
	if err != nil {
		rawResponse, err := json.Marshal(
			tableResponse{
				TableId: t.tableId,
				Code:    tableResponseCodeError,
				Message: "encode error",
				Data:    nil,
			})
		if err != nil {
			log.Printf("Failed to marshal tableResponse: %v", err)
			return
		}
		sender.Send(rawResponse)
		return
	}

	rawResponseExecutor, err := json.Marshal(
		tableResponse{
			TableId: t.tableId,
			Code:    tableResponseCodeSuccess,
			Message: "success",
			Data:    nil,
		},
	)
	if err != nil {
		log.Printf("Failed to marshal tableResponse: %v", err)
		return
	}
	sender.Send(rawResponseExecutor)

	log.Print("broadcasting...")

	rawResponseBroadcast, err := json.Marshal(
		tableResponse{
			TableId: t.tableId,
			Code:    tableResponseCodeUpdate,
			Message: "update",
			Data:    string(rawState),
		})
	if err != nil {
		log.Printf("Failed to marshal tableResponse: %v", err)
		return
	}
	t.broadcastQueue <- rawResponseBroadcast
}
