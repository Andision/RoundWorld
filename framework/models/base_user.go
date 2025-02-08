package models

type BaseUser struct {
	userID   string
	userName string
}

func (u *BaseUser) GetUserID() string {
	return u.userID
}
