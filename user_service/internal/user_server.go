package internal

import (
	pb "core/pkg/go_proto"
	"golang.org/x/net/context"
)

// // 1 -> many
// // DBCLient -> cinnection
// // Repository -> create_user, ...
// // Service    -> call create_user 2 times for example (buisness logic)
// // Controller (Endpoints) ->

type GRPCUserServer struct {
	pb.UnimplementedUserServer
}

func (server *GRPCUserServer) CreateUser(ctx context.Context, request *pb.CreateUserRequest) (*pb.UserDataResponse, error) {
	return &pb.UserDataResponse{
		Id:      "1",
		Name:    request.Name,
		Surname: request.Surname,
	}, nil
}
