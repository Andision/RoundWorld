package dal

type UserImpl struct {
	userName string
}

func NewUserImpl(userName string) *UserImpl {
	return &UserImpl{
		userName: userName,
	}
}

func (u *UserImpl) GetUserName() string {
	return u.userName
}
