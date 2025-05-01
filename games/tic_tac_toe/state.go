package tic_tac_toe

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/Andision/RoundWorld/framework/dal"
	"github.com/Andision/RoundWorld/framework/models"
	"log"
)

type TicTacToeStateType = int

const (
	TicTacToeStateTypeCreated = iota
	TicTacToeStateTypeStarted
	TicTacToeStateTypeEnded
)

type TicTacToeState struct {
	models.BaseState
	cells             [3][3]TicTacToeCellType
	hostPlayer        models.User
	guestPlayer       models.User
	firstTurnPlayer   models.User
	currentTurnPlayer models.User
	winnerPlayer      models.User
	stateType         TicTacToeStateType
}

func (s *TicTacToeState) isEnd() bool {
	if s.stateType == TicTacToeStateTypeEnded {
		return true
	} else if s.stateType == TicTacToeStateTypeCreated {
		return false
	} else if s.stateType == TicTacToeStateTypeStarted {
		log.Println("[isEnd] TicTacToeState is started")
		for i := 0; i < 3; i++ {
			for j := 0; j < 3; j++ {
				log.Printf("cells[%d][%d]: %d", i, j, s.cells[i][j])
			}
		}
		for i := 0; i < 3; i++ {
			if s.cells[i][0] == s.cells[i][1] && s.cells[i][1] == s.cells[i][2] && s.cells[i][0] != TicTacToeCellTypeNone {
				return true
			}
			if s.cells[0][i] == s.cells[1][i] && s.cells[1][i] == s.cells[2][i] && s.cells[0][i] != TicTacToeCellTypeNone {
				return true
			}
		}
		if s.cells[0][0] == s.cells[1][1] && s.cells[1][1] == s.cells[2][2] && s.cells[0][0] != TicTacToeCellTypeNone {
			return true
		}
		if s.cells[0][2] == s.cells[1][1] && s.cells[1][1] == s.cells[2][0] && s.cells[0][2] != TicTacToeCellTypeNone {
			return true
		}
	}
	return false
}

func (s *TicTacToeState) Lock() {
	s.BaseState.StateMutex.Lock()
}

func (s *TicTacToeState) Unlock() {
	s.BaseState.StateMutex.Unlock()
}

type TicTacToeStateJson struct {
	Cells             []int          `json:"cells"`
	HostPlayer        dal.UserIdType `json:"host_player"`
	GuestPlayer       dal.UserIdType `json:"guest_player"`
	FirstTurnPlayer   dal.UserIdType `json:"first_turn_player"`
	CurrentTurnPlayer dal.UserIdType `json:"current_turn_player"`
	WinnerPlayer      dal.UserIdType `json:"winner_player"`
}

func NewTicTacToeState(hostName string) *TicTacToeState {
	hostPlayer, err := dal.NewUserImplByName(hostName)
	if err != nil {
		log.Printf("[NewTicTacToeState] NewUserImplByName failed: %v", err)
		return nil
	}
	return &TicTacToeState{
		hostPlayer:   hostPlayer,
		winnerPlayer: nil,
	}
}

func (s *TicTacToeState) Encode(ctx context.Context) ([]byte, error) {
	log.Print("encode called")

	cells := make([]int, 9)

	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			cells[i*3+j] = int(s.cells[i][j])
		}
	}

	var winnerPlayerId dal.UserIdType
	if s.winnerPlayer != nil {
		winnerPlayerId = s.winnerPlayer.GetUserId()
	}

	return json.Marshal(TicTacToeStateJson{
		Cells:             cells,
		HostPlayer:        s.hostPlayer.GetUserId(),
		GuestPlayer:       s.guestPlayer.GetUserId(),
		FirstTurnPlayer:   s.firstTurnPlayer.GetUserId(),
		CurrentTurnPlayer: s.currentTurnPlayer.GetUserId(),
		WinnerPlayer:      winnerPlayerId,
	})
}

func (s *TicTacToeState) Parse(ctx context.Context, bytes []byte) (models.Action, error) {
	var rawAction rawTicTacToeAction
	err := json.Unmarshal(bytes, &rawAction)
	if err != nil {
		return nil, err
	}
	log.Printf("[TicTacToe.Parse] rawAction: %+v", rawAction)

	executorName, ok := ctx.Value("username").(string)
	if !ok {
		return nil, errors.New("username not found in context")
	}
	executor, err := dal.NewUserImplByName(executorName)
	if err != nil {
		log.Printf("[TicTacToe.Parse] NewUserImplByName failed: %v", err)
		return nil, err
	}

	return &TicTacToeAction{
		BaseAction: &models.BaseAction{
			Executor: executor,
			TableId:  ctx.Value("tableId").(string),
		},
		Type: &rawAction.Type,
		Data: rawAction.Data,
	}, nil

}

func (s *TicTacToeState) Validate(ctx context.Context, actionInterface models.Action) (bool, error) {
	s.Lock()
	defer s.Unlock()

	log.Printf("[TicTacToe.Validate] validate action: %v", actionInterface)
	action, ok := actionInterface.(*TicTacToeAction)
	if !ok {
		return false, errors.New("action is not a TicTacToeAction")
	}

	switch *action.Type {
	case TicTacToeActionTypeReset:
		if s.stateType != TicTacToeStateTypeEnded {
			return false, errors.New("state is not ended")
		}
		return true, nil
	case TicTacToeActionTypeSet:
		if s.stateType != TicTacToeStateTypeCreated {
			return false, errors.New("state is not created")
		}
		return false, errors.New("not support yet")
	case TicTacToeActionTypePut:
		if s.stateType != TicTacToeStateTypeStarted {
			return false, errors.New("state is not started")
		}

		if s.currentTurnPlayer.GetUserId() != action.Executor.GetUserId() {
			return false, errors.New("not your turn, current turn is " + s.currentTurnPlayer.GetUserName() + " but you are " + action.Executor.GetUserName())
		}

		actionPosition, err := action.GetPosition()
		if err != nil {
			return false, err
		}

		if actionPosition.X < 0 || actionPosition.X > 2 || actionPosition.Y < 0 || actionPosition.Y > 2 {
			return false, errors.New("position is out of range")
		}

		if s.cells[actionPosition.X][actionPosition.Y] != TicTacToeCellTypeNone {
			return false, errors.New("cell is already occupied")
		}

		if s.currentTurnPlayer.GetUserName() == s.hostPlayer.GetUserName() {
			s.currentTurnPlayer = s.guestPlayer
		} else {
			s.currentTurnPlayer = s.hostPlayer
		}

		return true, nil
	case TicTacToeActionTypeStart:
		if s.stateType != TicTacToeStateTypeCreated {
			return false, errors.New("state is not created")
		}

		if s.guestPlayer == nil {
			return false, errors.New("guest player is not set")
		}

		if s.firstTurnPlayer == nil {
			return false, errors.New("first player is not set")
		}

		playerCount := ctx.Value("playerCount")
		playerCount, ok := playerCount.(int)
		if !ok {
			return false, errors.New("can not get playerCount from context")
		}
		if playerCount != 2 {
			return false, errors.New("players number is not 2")
		}

		return true, nil
	}

	return false, errors.New("unknown action type")

}

func (s *TicTacToeState) Execute(ctx context.Context, rawAction models.Action) error {
	s.Lock()
	defer s.Unlock()

	log.Printf("execute action: %v", rawAction)
	action, ok := rawAction.(*TicTacToeAction)
	if !ok {
		return errors.New("action is not a TicTacToeAction")
	}

	switch *action.Type {
	case TicTacToeActionTypeReset:
		s.firstTurnPlayer = s.hostPlayer
		s.winnerPlayer = nil
		s.cells = [3][3]TicTacToeCellType{}
	case TicTacToeActionTypeSet:
		// TODO implement me
		return nil
	case TicTacToeActionTypePut:
		var piece TicTacToeCellType
		if action.BaseAction.Executor.GetUserName() == s.firstTurnPlayer.GetUserName() {
			piece = TicTacToeCellTypeX
		} else {
			piece = TicTacToeCellTypeO
		}

		actionPosition, _ := action.GetPosition()
		s.cells[actionPosition.X][actionPosition.Y] = piece

		if s.isEnd() {
			s.stateType = TicTacToeStateTypeEnded
			s.winnerPlayer = action.BaseAction.Executor
			return nil
		}

		return nil
	case TicTacToeActionTypeStart:
		s.stateType = TicTacToeStateTypeStarted
		s.currentTurnPlayer = s.firstTurnPlayer
		return nil
	}

	return errors.New("unknown action type")
}

func (s *TicTacToeState) UpdatePlayers(ctx context.Context, players []dal.UserIdType) error {
	s.Lock()
	defer s.Unlock()

	if len(players) == 1 {
		s.hostPlayer = dal.MustNewUserImplById(players[0])
		s.firstTurnPlayer = s.hostPlayer
		return nil
	} else if len(players) >= 2 {
		s.hostPlayer = dal.MustNewUserImplById(players[0])
		s.guestPlayer = dal.MustNewUserImplById(players[1])
		s.firstTurnPlayer = s.hostPlayer
		return nil
	}

	return errors.New("invalid players count: " + string(len(players)))
}
