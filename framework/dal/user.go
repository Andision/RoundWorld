package dal

type UserIdType string

type UserImpl struct {
	userName string
	userId   UserIdType
}

func NewUserImplByName(userName string) (*UserImpl, error) {
	return &UserImpl{
		userName: userName,
		userId:   UserIdType(userName),
	}, nil
}

func MustNewUserImplByName(userName string) *UserImpl {
	user, err := NewUserImplByName(userName)
	if err != nil {
		panic(err)
	}
	return user
}

func NewUserImplById(userId UserIdType) (*UserImpl, error) {
	return &UserImpl{
		userName: string(userId),
		userId:   userId,
	}, nil
}

func MustNewUserImplById(userId UserIdType) *UserImpl {
	user, err := NewUserImplById(userId)
	if err != nil {
		panic(err)
	}
	return user
}

func (u *UserImpl) GetUserName() string {
	return u.userName
}

func (u *UserImpl) GetUserId() UserIdType {
	return u.userId
}
