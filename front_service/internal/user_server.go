package internal

import (
	"context"
	pb "core/pkg/go_proto"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type FrontUserServer struct {
	pb.UnimplementedFrontUserServer
	userServerAddress string
}

func (server *FrontUserServer) CreateUser(ctx context.Context, request *pb.CreateUserRequest) (*pb.UserDataResponse, error) {
	conn, err := grpc.Dial(server.userServerAddress, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	client := pb.NewUserClient(conn)
	return client.CreateUser(ctx, request)
}

func RegisterRestUserServer(mux *runtime.ServeMux, userServerAddress string) {
	// TODO: if pass  context.WithValue, can't extract it via function
	err := pb.RegisterFrontUserHandlerServer(context.Background(), mux, &FrontUserServer{
		userServerAddress: userServerAddress,
	})
	if err != nil {
		panic(err)
	}
}
