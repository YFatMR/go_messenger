package internal

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	proto "protocol/pkg/proto"
)

type frontUserServer struct {
	proto.UnimplementedFrontUserServer
	userServerAddress string
}

func (server *frontUserServer) CreateUser(ctx context.Context, request *proto.UserData) (*proto.UserId, error) {
	fmt.Println(" CreateUser call")
	conn, err := grpc.Dial(server.userServerAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	client := proto.NewUserClient(conn)
	return client.CreateUser(ctx, request)
}

func (server *frontUserServer) GetUserById(ctx context.Context, request *proto.UserId) (*proto.UserData, error) {
	fmt.Println(" GetUserById call")
	conn, err := grpc.Dial(server.userServerAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	client := proto.NewUserClient(conn)
	return client.GetUserById(ctx, request)
}

// registration

func RegisterRestUserServer(ctx context.Context, mux *runtime.ServeMux, grpcFrontServerAddress string) {
	// TODO: if pass  context.WithValue, can't extract it via function
	fmt.Println("CreateUser registration")

	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err := proto.RegisterFrontUserHandlerFromEndpoint(ctx, mux, grpcFrontServerAddress, opts)
	if err != nil {
		panic(err)
	}
}

func RegisterGrpcUserServer(grpcServer grpc.ServiceRegistrar, userServerAddress string) {
	proto.RegisterFrontUserServer(grpcServer, &frontUserServer{
		userServerAddress: userServerAddress})
}
