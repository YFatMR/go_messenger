package main

import (
	pb "core/pkg/go_proto"
	. "core/pkg/utils"
	"google.golang.org/grpc"
	"net"
	. "user_server/internal"
)

func main() {
	s := grpc.NewServer()

	pb.RegisterUserServer(s, &GRPCUserServer{})

	userServerAddress := GetFullServiceAddress("USER")
	l, err := net.Listen("tcp", userServerAddress)
	if err == nil {
		s.Serve(l)
	}
}
