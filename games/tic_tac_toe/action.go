package tic_tac_toe

import (
	"errors"
	"github.com/Andision/RoundWorld/framework/models"
)

type TicTacToeActionData struct {
	PutActionData int `json:"put_action_data"`
}

type rawTicTacToeAction struct {
	Type string               `json:"type"`
	Data *TicTacToeActionData `json:"data"`
}

const (
	TicTacToeActionTypeSet   = "set"
	TicTacToeActionTypePut   = "put"
	TicTacToeActionTypeReset = "reset"
	TicTacToeActionTypeStart = "start"
)

type TicTacToeAction struct {
	*models.BaseAction
	Type *string
	Data *TicTacToeActionData
}

func (action *TicTacToeAction) Validator() (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (action *TicTacToeAction) GetPosition() (*TicTacToePosition, error) {
	rawPosition := 0
	if action.Data.PutActionData != 0 {
		rawPosition = action.Data.PutActionData
	}

	if rawPosition == 0 {
		return nil, errors.New("invalid action data")
	}

	rawPosition = rawPosition - 1
	position := TicTacToePosition{
		X: rawPosition / 3,
		Y: rawPosition % 3,
	}
	return &position, nil
}
