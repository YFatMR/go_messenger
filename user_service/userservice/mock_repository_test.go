package userservice_test

import (
	"context"

	"github.com/YFatMR/go_messenger/user_service/entities/account"
	"github.com/YFatMR/go_messenger/user_service/entities/credential"
	"github.com/YFatMR/go_messenger/user_service/entities/user"
	"github.com/YFatMR/go_messenger/user_service/entities/userid"
)

type CreateResponseData struct {
	UserID *userid.Entity
	Error  error
}

type GetByIDResponseData struct {
	User  *user.Entity
	Error error
}

type DeleteByIDResponseData struct {
	Error error
}

type GetAccountByLoginResponseData struct {
	Account *account.Entity
	Error   error
}

type MockUserRepository struct {
	CreateResponse            CreateResponseData
	GetByIDResponse           GetByIDResponseData
	DeleteByIDResponse        DeleteByIDResponseData
	GetAccountByLoginResponse GetAccountByLoginResponseData
}

func (r *MockUserRepository) Create(ctx context.Context, user *user.Entity, credential *credential.Entity) (
	*userid.Entity, error,
) {
	return r.CreateResponse.UserID, r.CreateResponse.Error
}

func (r *MockUserRepository) GetByID(ctx context.Context, userID *userid.Entity) (
	*user.Entity, error,
) {
	return r.GetByIDResponse.User, r.GetByIDResponse.Error
}

func (r *MockUserRepository) DeleteByID(ctx context.Context, userID *userid.Entity) error {
	return r.DeleteByIDResponse.Error
}

func (r *MockUserRepository) GetAccountByLogin(ctx context.Context, login string) (
	*account.Entity, error,
) {
	return r.GetAccountByLoginResponse.Account, r.GetAccountByLoginResponse.Error
}
