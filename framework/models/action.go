package models

type Action interface {
	Validator() (bool, error)
}
