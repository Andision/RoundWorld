package models

type Action interface {
	GetPerformers() *[]User
	GetFirstPerformer() (*User, error)
}
