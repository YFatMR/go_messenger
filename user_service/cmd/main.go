package main

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"net"
	pb "user_server/pkg/proto"
)

type GRPCUserServer struct {
	pb.UnimplementedUserServiceServer
}

func (server *GRPCUserServer) CreateUser(ctx context.Context, request *pb.CreateUserRequest) (*pb.UserDataResponse, error) {
	return &pb.UserDataResponse{
		Id:      "1",
		Name:    request.Name,
		Surname: request.Surname,
	}, nil
}

//func (c *ExpensesConsumer) ExpensesReport(ctx context.Context, request *pb.ExpensesRequest) (*pb.ExpensesResponse, error) {
//	if err := c.telegramClient.SendMessage(request.UserId, request.Expenses); err != nil {
//		return &pb.ExpensesResponse{
//			StatusCode: pb.ExpensesResponse_ERROR,
//			Message:    err.Error(),
//		}, nil
//	}
//
//	if err := c.cache.SaveExpenses(ctx, request.UserId, request.Expenses, request.Period); err != nil {
//		return nil, err
//	}
//
//	return &pb.ExpensesResponse{
//		StatusCode: pb.ExpensesResponse_OK,
//		Message:    "ok",
//	}, nil
//}

func main() {
	s := grpc.NewServer()
	server := &GRPCUserServer{}
	pb.RegisterUserServiceServer(s, server)

	l, err := net.Listen("tcp", ":8000")
	if err == nil {
		s.Serve(l)
	}

}
