package usercontroller

import (
	"context"

	"github.com/YFatMR/go_messenger/core/pkg/ulo"
	"github.com/YFatMR/go_messenger/protocol/pkg/proto"
	"github.com/YFatMR/go_messenger/user_service/internal/entities/accountid"
	"github.com/YFatMR/go_messenger/user_service/internal/entities/user"
	"github.com/YFatMR/go_messenger/user_service/internal/entities/userid"
	"github.com/YFatMR/go_messenger/user_service/internal/services"
)

type UserController struct {
	service services.UserService
}

func New(service services.UserService) *UserController {
	return &UserController{
		service: service,
	}
}

func (c *UserController) Create(ctx context.Context, request *proto.CreateUserDataRequest) (
	*proto.UserID, ulo.LogStash, error,
) {
	user, err := user.FromProtobuf(request.GetUserData())
	if err != nil {
		return nil, ulo.FromErrorWithMsg("Wrong format for user data", err), ErrWrongRequestFormat
	}

	accountID, err := accountid.FromProtobuf(request.GetAccountID())
	if err != nil {
		return nil, ulo.FromErrorWithMsg("Wrong format for account id", err), ErrWrongRequestFormat
	}

	insertedID, _, err := c.service.Create(ctx, user, accountID)
	if err != nil {
		return nil, nil, err
	}

	return &proto.UserID{
		ID: insertedID.GetID(),
	}, nil, nil
}

func (c *UserController) GetByID(ctx context.Context, request *proto.UserID) (
	*proto.UserData, ulo.LogStash, error,
) {
	userID, err := userid.FromProtobuf(request)
	if err != nil {
		return nil, ulo.FromErrorWithMsg("Wrong format for userID data", err), ErrWrongRequestFormat
	}

	userData, _, err := c.service.GetByID(ctx, userID)
	if err != nil {
		return nil, nil, err
	}

	return &proto.UserData{
		Name:    userData.GetName(),
		Surname: userData.GetSurname(),
	}, nil, nil
}

func (c *UserController) DeleteByID(ctx context.Context, request *proto.UserID) (
	*proto.Void, ulo.LogStash, error,
) {
	userID, err := userid.FromProtobuf(request)
	if err != nil {
		return nil, ulo.FromErrorWithMsg("Wrong format for userID data", err), ErrWrongRequestFormat
	}

	_, err = c.service.DeleteByID(ctx, userID)
	return &proto.Void{}, nil, err
}

func (c *UserController) Ping(ctx context.Context, request *proto.Void) (
	*proto.Pong, ulo.LogStash, error,
) {
	return &proto.Pong{
		Message: "pong",
	}, nil, nil
}
