package models

import "context"

type State interface {
	Encode(ctx context.Context) ([]byte, error)
	Parse(ctx context.Context, raw []byte) (Action, error)
	Validate(ctx context.Context, action Action) (bool, error)
	Execute(ctx context.Context, action Action) error
	Lock()
	Unlock()
}
