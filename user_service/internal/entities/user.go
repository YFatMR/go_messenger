package entities

import "github.com/YFatMR/go_messenger/protocol/pkg/proto"

type User struct {
	name    string
	surname string
}

func NewUser(name string, surname string) *User {
	return &User{
		name:    name,
		surname: surname,
	}
}

func NewUserFromProtobuf(user *proto.UserData) (*User, error) {
	if user == nil || user.GetName() == "" || user.GetSurname() == "" {
		return nil, ErrWrongRequestFormat
	}
	return NewUser(user.GetName(), user.GetSurname()), nil
}

func (u *User) GetName() string {
	return u.name
}

func (u *User) GetSurname() string {
	return u.surname
}
