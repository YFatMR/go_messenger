package httpapi

import (
	"net/http"

	"github.com/YFatMR/go_messenger/core/pkg/czap"
	"github.com/YFatMR/go_messenger/front_server/websocketapi"
	"github.com/YFatMR/go_messenger/protocol/pkg/proto"
)

type FrontServer struct {
	userServiceClient    proto.UserClient
	sandboxServiceClient proto.SandboxClient
	dialogServiceClient  proto.DialogServiceClient
	botsServiceClient    proto.BotsServiceClient
	websocketClient      *websocketapi.Client
	logger               *czap.Logger
}

func NewFrontServer(userServiceClient proto.UserClient, sandboxServiceClient proto.SandboxClient,
	dialogServiceClient proto.DialogServiceClient, botsServiceClient proto.BotsServiceClient,
	websocketClient *websocketapi.Client, logger *czap.Logger,
) FrontServer {
	return FrontServer{
		userServiceClient:    userServiceClient,
		sandboxServiceClient: sandboxServiceClient,
		dialogServiceClient:  dialogServiceClient,
		botsServiceClient:    botsServiceClient,
		websocketClient:      websocketClient,
		logger:               logger,
	}
}

func (s *FrontServer) Ping(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("pong"))
}
