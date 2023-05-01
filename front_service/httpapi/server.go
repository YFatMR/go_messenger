package httpapi

import (
	"net/http"

	"github.com/YFatMR/go_messenger/protocol/pkg/proto"
)

type FrontServer struct {
	userServiceClient    proto.UserClient
	sandboxServiceClient proto.SandboxClient
	dialogServiceClient  proto.DialogServiceClient
}

func NewFrontServer(userServiceClient proto.UserClient, sandboxServiceClient proto.SandboxClient,
	dialogServiceClient proto.DialogServiceClient,
) FrontServer {
	return FrontServer{
		userServiceClient:    userServiceClient,
		sandboxServiceClient: sandboxServiceClient,
		dialogServiceClient:  dialogServiceClient,
	}
}

func (s *FrontServer) Ping(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("pong"))
}
