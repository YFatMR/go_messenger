package entity

import (
	"github.com/YFatMR/go_messenger/protocol/pkg/proto"
)

type User struct {
	Nickname    string
	Name        string
	Surname     string
	Github      string
	Linkedin    string
	PublicEmail string
	Status      string
}

func UserFromProtobuf(user *proto.UserData) (*User, error) {
	if user == nil || user.Name == "" || user.Surname == "" || user.Nickname == "" {
		return nil, ErrWrongRequestFormat
	}
	return &User{
		Name:        user.Name,
		Surname:     user.Surname,
		Nickname:    user.Nickname,
		Github:      user.Github,
		Linkedin:    user.Linkedin,
		PublicEmail: user.PublicEmail,
		Status:      user.StatusText,
	}, nil
}

func UserToProtobuf(user *User) *proto.UserData {
	return &proto.UserData{
		Name:        user.Name,
		Surname:     user.Surname,
		Nickname:    user.Nickname,
		Github:      user.Github,
		Linkedin:    user.Linkedin,
		PublicEmail: user.PublicEmail,
		StatusText:  user.Status,
	}
}

type UserWithID struct {
	UserID UserID
	User
}

func UserWithIDToProtobuf(user *UserWithID) *proto.UserDataWithID {
	return &proto.UserDataWithID{
		UserID:      UserIDToProtobuf(&user.UserID),
		Name:        user.Name,
		Surname:     user.Surname,
		Nickname:    user.Nickname,
		Github:      user.Github,
		Linkedin:    user.Linkedin,
		PublicEmail: user.PublicEmail,
		StatusText:  user.Status,
	}
}

func UsersWithIDToProtobuf(users []*UserWithID) []*proto.UserDataWithID {
	result := make([]*proto.UserDataWithID, 0, len(users))
	for _, user := range users {
		result = append(result, UserWithIDToProtobuf(user))
	}
	return result
}

type UpdateUserRequest struct {
	Name        string
	Surname     string
	Github      string
	Linkedin    string
	PublicEmail string
	Status      string
}

func UpdateUserRequestFromProtobuf(user *proto.UpdateUserDataRequest) (*UpdateUserRequest, error) {
	if user == nil || user.Name == "" || user.Surname == "" {
		return nil, ErrWrongRequestFormat
	}
	return &UpdateUserRequest{
		Name:        user.Name,
		Surname:     user.Surname,
		Github:      user.Github,
		Linkedin:    user.Linkedin,
		PublicEmail: user.PublicEmail,
		Status:      user.StatusText,
	}, nil
}
