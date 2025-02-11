package tic_tac_toe

import "github.com/Andision/RoundWorld/framework/models"

type TicTacToeCellType int

const (
	TicTacToeCellTypeNone TicTacToeCellType = iota
	TicTacToeCellTypeX
	TicTacToeCellTypeO
)

// TicTacToePosition is the position to put the piece, the X refers to row and the Y refers to column.
type TicTacToePosition struct {
	X int
	Y int
}

type TicTacToeConfig struct {
	gameType  string
	maxPlayer int
	minPlayer int
}

func (t *TicTacToeConfig) GetNewGameState(executor string) models.State {
	return NewTicTacToeState(executor)
}

func NewTicTacToeConfig() *TicTacToeConfig {
	return &TicTacToeConfig{
		gameType:  "TicTacToe",
		maxPlayer: 2,
		minPlayer: 2,
	}
}

func (t *TicTacToeConfig) GetGameType() string {
	return t.gameType
}

func (t *TicTacToeConfig) GetMaxPlayer() int {
	return t.maxPlayer
}

func (t *TicTacToeConfig) GetMinPlayer() int {
	return t.minPlayer
}
