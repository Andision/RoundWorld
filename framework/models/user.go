package models

import "github.com/Andision/RoundWorld/framework/dal"

type User interface {
	GetUserName() string
	GetUserId() dal.UserIdType
}
