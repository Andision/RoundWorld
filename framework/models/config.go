package models

type Config interface {
	GetGameType() string
	GetMaxPlayer() int
	GetMinPlayer() int

	// GetNewGameState returns a new game state for the given host and players.
	// The host is the player who created the game, and players are the participants.
	// The players slice should not contain the host!
	GetNewGameState(executor string) State
}
