package usercontroller

import (
	"context"

	"github.com/YFatMR/go_messenger/core/pkg/errors/cerrors"
	"github.com/YFatMR/go_messenger/protocol/pkg/proto"
	"github.com/YFatMR/go_messenger/user_service/internal/entities"
	accountid "github.com/YFatMR/go_messenger/user_service/internal/entities/account_id"
	"github.com/YFatMR/go_messenger/user_service/internal/entities/user"
	userid "github.com/YFatMR/go_messenger/user_service/internal/entities/user_id"
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

func (c *UserController) Create(ctx context.Context, request *proto.CreateUserDataRequest) (*proto.UserID, cerrors.Error) {
	user, err := user.FromProtobuf(request.GetUserData())
	if err != nil {
		return nil, cerrors.New("Wrong format for user data", err, entities.ErrWrongRequestFormat)
	}

	accountID, err := accountid.FromProtobuf(request.GetAccountID())
	if err != nil {
		return nil, cerrors.New("Wrong format for account id", err, entities.ErrWrongRequestFormat)
	}

	insertedID, cerr := c.service.Create(ctx, user, accountID)
	if cerr != nil {
		return nil, cerr
	}

	return &proto.UserID{
		ID: insertedID.GetID(),
	}, nil
}

func (c *UserController) GetByID(ctx context.Context, request *proto.UserID) (*proto.UserData, cerrors.Error) {
	userID, err := userid.FromProtobuf(request)
	if err != nil {
		return nil, cerrors.New("Wrong format for userID data", err, entities.ErrWrongRequestFormat)
	}

	userData, cerr := c.service.GetByID(ctx, userID)
	if err != nil {
		return nil, cerr
	}

	return &proto.UserData{
		Name:    userData.GetName(),
		Surname: userData.GetSurname(),
	}, nil
}

func (c *UserController) DeleteByID(ctx context.Context, request *proto.UserID) (*proto.Void, cerrors.Error) {
	userID, err := userid.FromProtobuf(request)
	if err != nil {
		return nil, cerrors.New("Wrong format for userID data", err, entities.ErrWrongRequestFormat)
	}

	cerr := c.service.DeleteByID(ctx, userID)
	return &proto.Void{}, cerr
}

func (c *UserController) Ping(ctx context.Context, request *proto.Void) (*proto.Pong, cerrors.Error) {
	return &proto.Pong{
		Message: "pong",
	}, nil
}
