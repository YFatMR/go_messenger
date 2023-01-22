package entities

import proto "github.com/YFatMR/go_messenger/protocol/pkg/proto"

type User struct {
	name    string
	surname string
}

func newUser(name string, surname string) *User {
	return &User{
		name:    name,
		surname: surname,
	}
}

func NewUserFromProtobuf(user *proto.UserData) (*User, error) {
	if user == nil || user.GetName() == "" || user.GetSurname() == "" {
		return nil, ErrWrongRequestFormat
	}
	return newUser(user.GetName(), user.GetSurname()), nil
}

func NewUserFromRawDatabaseDocument(name string, surname string) *User {
	return newUser(name, surname)
}

func NewMockUser(name string, surname string) *User {
	return newUser(name, surname)
}

func (u *User) GetName() string {
	return u.name
}

func (u *User) GetSurname() string {
	return u.surname
}
