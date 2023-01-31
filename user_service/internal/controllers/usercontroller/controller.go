package usercontroller

import (
	"context"

	"github.com/YFatMR/go_messenger/core/pkg/errors/logerr"
	"github.com/YFatMR/go_messenger/protocol/pkg/proto"
	"github.com/YFatMR/go_messenger/user_service/internal/entities"
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

func (c *UserController) Create(ctx context.Context, request *proto.CreateUserDataRequest) (*proto.UserID, logerr.Error) {
	user, err := user.FromProtobuf(request.GetUserData())
	if err != nil {
		return nil, logerr.NewError(entities.ErrWrongRequestFormat, "Wrong format for user data", logerr.Err(err))
	}

	accountID, err := accountid.FromProtobuf(request.GetAccountID())
	if err != nil {
		return nil, logerr.NewError(entities.ErrWrongRequestFormat, "Wrong format for account id", logerr.Err(err))
	}

	insertedID, lerr := c.service.Create(ctx, user, accountID)
	if lerr != nil {
		return nil, lerr
	}

	return &proto.UserID{
		ID: insertedID.GetID(),
	}, nil
}

func (c *UserController) GetByID(ctx context.Context, request *proto.UserID) (*proto.UserData, logerr.Error) {
	userID, err := userid.FromProtobuf(request)
	if err != nil {
		return nil, logerr.NewError(entities.ErrWrongRequestFormat, "Wrong format for userID data", logerr.Err(err))
	}

	userData, lerr := c.service.GetByID(ctx, userID)
	if lerr != nil {
		return nil, lerr
	}

	return &proto.UserData{
		Name:    userData.GetName(),
		Surname: userData.GetSurname(),
	}, nil
}

func (c *UserController) DeleteByID(ctx context.Context, request *proto.UserID) (*proto.Void, logerr.Error) {
	userID, err := userid.FromProtobuf(request)
	if err != nil {
		return nil, logerr.NewError(entities.ErrWrongRequestFormat, "Wrong format for userID data", logerr.Err(err))
	}

	lerr := c.service.DeleteByID(ctx, userID)
	return &proto.Void{}, lerr
}

func (c *UserController) Ping(ctx context.Context, request *proto.Void) (*proto.Pong, logerr.Error) {
	return &proto.Pong{
		Message: "pong",
	}, nil
}
