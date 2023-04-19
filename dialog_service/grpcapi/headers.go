package grpcapi

import "github.com/YFatMR/go_messenger/core/pkg/configs/cviper"

type Headers struct {
	UserID string
}

func NewHeaders(userID string) Headers {
	return Headers{UserID: userID}
}

func HeadersFromConfig(config *cviper.CustomViper) Headers {
	return Headers{
		UserID: config.GetStringRequired("GRPC_USER_ID_HEADER"),
	}
}
