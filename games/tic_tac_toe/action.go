package tic_tac_toe

import (
	"errors"
	"github.com/Andision/RoundWorld/framework/models"
)

type rawTicTacToeAction struct {
	Executor string      `json:"executor"`
	Type     string      `json:"type"`
	Data     interface{} `json:"data"`
}

const (
	TicTacToeActionTypeSet   = "set"
	TicTacToeActionTypePut   = "put"
	TicTacToeActionTypeReset = "reset"
	TicTacToeActionTypeStart = "start"
)

type TicTacToeAction struct {
	models.BaseAction
	Type string
	Data interface{}
}

func (action *TicTacToeAction) Validator() (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (action *TicTacToeAction) GetPosition() (*TicTacToePosition, error) {
	rawPosition, ok := action.Data.(int)
	if !ok {
		return nil, errors.New("raw position data is not a int")
	}

	position := TicTacToePosition{
		X: rawPosition / 3,
		Y: rawPosition % 3,
	}

	return &position, nil
}
