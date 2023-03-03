package entity

import "github.com/YFatMR/go_messenger/protocol/pkg/proto"

func PongProtobuf() *proto.Pong {
	return &proto.Pong{
		Message: "pong",
	}
}

func VoidProtobuf() *proto.Void {
	return &proto.Void{}
}
