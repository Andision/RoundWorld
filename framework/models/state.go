package models

type State interface {
	Execute(action *Action)
}
