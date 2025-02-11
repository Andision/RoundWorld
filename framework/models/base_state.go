package models

import "sync"

type BaseState struct {
	StateMutex   sync.RWMutex
	roundTimeout int
}
