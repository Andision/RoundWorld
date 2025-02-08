package models

type BaseState struct {
	gameID       string
	userList     []User
	roundTimeout int
	gameTimeout  int
}
