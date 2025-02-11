package games

import (
	"github.com/Andision/RoundWorld/framework/models"
	"github.com/Andision/RoundWorld/games/tic_tac_toe"
)

type ConfigGames struct {
	availableGames map[string]models.Config
}

func NewConfigGames() *ConfigGames {
	return &ConfigGames{
		availableGames: map[string]models.Config{
			"TicTacToe": tic_tac_toe.NewTicTacToeConfig(),
		},
	}
}

func (c *ConfigGames) GetAvailableGames() map[string]models.Config {
	return c.availableGames
}
