package models

type Config interface {
	GetGameType() string
	GetMaxPlayer() int
	GetMinPlayer() int
	GetNewGameState(executor string) State
}
