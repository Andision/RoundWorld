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
	cells         [3][3]TicTacToeCellType
	firstPlayer   models.User
	currentPlayer models.User
	stateType     TicTacToeStateType
}

func (s *TicTacToeState) Lock() {
	s.BaseState.StateMutex.Lock()
}

func (s *TicTacToeState) Unlock() {
	s.BaseState.StateMutex.Unlock()
}

type TicTacToeStateJson struct {
	Cells         []int  `json:"cells"`
	FirstPlayer   string `json:"firstPlayer"`
	CurrentPlayer string `json:"currentPlayer"`
}

func NewTicTacToeState(executor string) *TicTacToeState {
	return &TicTacToeState{
		firstPlayer:   dal.NewUserImpl(executor),
		currentPlayer: dal.NewUserImpl(executor),
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

	return json.Marshal(TicTacToeStateJson{
		Cells:         cells,
		FirstPlayer:   s.firstPlayer.GetUserName(),
		CurrentPlayer: s.currentPlayer.GetUserName(),
	})
}

func (s *TicTacToeState) Parse(ctx context.Context, bytes []byte) (models.Action, error) {
	var rawAction rawTicTacToeAction
	err := json.Unmarshal(bytes, &rawAction)
	if err != nil {
		return nil, err
	}

	return &TicTacToeAction{
		BaseAction: models.BaseAction{
			Executor: dal.NewUserImpl(rawAction.Executor),
			TableId:  ctx.Value("tableId").(string),
		},
		Type: rawAction.Type,
		Data: rawAction.Data,
	}, nil

}

func (s *TicTacToeState) Validate(ctx context.Context, actionInterface models.Action) (bool, error) {
	log.Printf("validate action: %v", actionInterface)
	action, ok := actionInterface.(*TicTacToeAction)
	if !ok {
		return false, errors.New("action is not a TicTacToeAction")
	}

	switch action.Type {
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
		if s.stateType != TicTacToeStateTypeCreated {
			return false, errors.New("state is not created")
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

		return true, nil
	case TicTacToeActionTypeStart:
		if s.stateType != TicTacToeStateTypeCreated {
			return false, errors.New("state is not created")
		}

		playerCount := ctx.Value("playerCount")
		playerCount, ok := playerCount.(int)
		if !ok {
			return false, errors.New("can not get playerCount from context")
		}

		if playerCount != 2 {
			return false, errors.New("players number is not 2")
		}

		if s.firstPlayer == nil {
			return false, errors.New("first player is not set")
		}

		return true, nil

	}

	return false, errors.New("unknown action type")

}

func (s *TicTacToeState) Execute(ctx context.Context, rawAction models.Action) error {
	log.Printf("execute action: %v", rawAction)
	action, ok := rawAction.(*TicTacToeAction)
	if !ok {
		return errors.New("action is not a TicTacToeAction")
	}

	switch action.Type {
	case TicTacToeActionTypeReset:
		s.firstPlayer = action.BaseAction.Executor
		s.cells = [3][3]TicTacToeCellType{}
	case TicTacToeActionTypeSet:
		// TODO implement me
		return nil
	case TicTacToeActionTypePut:
		var piece TicTacToeCellType
		if action.BaseAction.Executor.GetUserName() == s.firstPlayer.GetUserName() {
			piece = TicTacToeCellTypeX
		} else {
			piece = TicTacToeCellTypeO
		}

		actionPosition, _ := action.GetPosition()
		s.cells[actionPosition.X][actionPosition.Y] = piece

		return nil
	case TicTacToeActionTypeStart:
		s.stateType = TicTacToeStateTypeStarted
		s.currentPlayer = s.firstPlayer
		return nil
	}

	return errors.New("unknown action type")
}
