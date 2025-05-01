package models

import (
	"context"
	"github.com/Andision/RoundWorld/framework/dal"
)

type State interface {
	UpdatePlayers(ctx context.Context, players []dal.UserIdType) error
	Encode(ctx context.Context) ([]byte, error)
	Parse(ctx context.Context, raw []byte) (Action, error)
	Validate(ctx context.Context, action Action) (bool, error)
	Execute(ctx context.Context, action Action) error
	Lock()
	Unlock()
}
