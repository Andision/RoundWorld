package web

import (
	"context"
	"github.com/Andision/RoundWorld/framework/dal"
	"github.com/Andision/RoundWorld/framework/models"
	"github.com/google/uuid"
	"log"
	"net/http"
	"sync"
)

// Table maintains the set of active clients and broadcasts messages to the
// clients.
type Table struct {
	gameConfig models.Config
	gameState  models.State

	tableId      string
	tableMutex   sync.RWMutex
	players      map[*Player]bool
	playersMutex sync.RWMutex

	broadcast  chan []byte
	queue      chan tableMessage
	unregister chan *Player
}

func NewTable(config models.Config, executor string) *Table {
	return &Table{
		gameConfig: config,
		gameState:  config.GetNewGameState(executor),

		tableId: uuid.New().String(),
		players: make(map[*Player]bool),

		broadcast:  make(chan []byte),
		queue:      make(chan tableMessage),
		unregister: make(chan *Player),
	}
}

func (t *Table) GetTableId() string {
	return t.tableId
}

func (t *Table) Run() {
	for {
		select {
		case player := <-t.unregister:
			t.playersMutex.Lock()
			if _, ok := t.players[player]; ok {
				delete(t.players, player)
				close(player.send)
			}
			t.playersMutex.Unlock()
		case message := <-t.broadcast:
			log.Printf("Broadcasting message to %d players", len(t.players))
			t.playersMutex.Lock()
			for client := range t.players {
				select {
				case client.send <- message:
				default:
					log.Printf("Closing send channel for client: %v", client)
					close(client.send)
					delete(t.players, client)
				}
			}
			t.playersMutex.Unlock()
		case message := <-t.queue:
			log.Printf("Received message: %s", message.message)

			ctx := context.Background()
			ctx = context.WithValue(ctx, "tableId", t.tableId)
			ctx = context.WithValue(ctx, "playerCount", len(t.players))
			ctx = context.WithValue(ctx, "players", t.players)

			tableMessageHandler(ctx, t, &message)
		default:
			//log.Printf("No message received")
		}
	}
}

func (t *Table) Join(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	t.playersMutex.Lock()
	defer t.playersMutex.Unlock()

	currentPlayerCount := len(t.players)
	if currentPlayerCount >= t.gameConfig.GetMaxPlayer() {
		http.Error(w, "Table is full", http.StatusBadRequest)
		return
	}

	username := ctx.Value("username")
	if username == nil {
		http.Error(w, "Username not found", http.StatusBadRequest)
		return
	}

	user := dal.NewUserImpl(username.(string))

	player := CreatePlayer(user, t, w, r)
	if player == nil {
		http.Error(w, "Failed to create player", http.StatusInternalServerError)
	}

	t.players[player] = true
	log.Printf("player %s joined", username)
	return
}
