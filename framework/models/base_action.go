package models

import "errors"

type BaseAction struct {
	performers []User
}

func (action *BaseAction) GetPerformers() *[]User {
	return &action.performers
}

func (action *BaseAction) GetFirstPerformer() (*User, error) {
	if len(action.performers) == 0 {
		return nil, errors.New("no performer found")
	}
	return &action.performers[0], nil
}
