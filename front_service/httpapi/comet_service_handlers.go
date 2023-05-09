package httpapi

import (
	"errors"
	"net/http"

	"github.com/YFatMR/go_messenger/front_server/websocketapi"
)

func (s *FrontServer) WebsocketHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	s.logger.DebugContext(ctx, "WebsocketHandler starting...")
	err := s.websocketClient.ProxifyWithAuth(ctx, w, r)
	if errors.Is(err, websocketapi.ErrAuth) {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	s.logger.DebugContext(ctx, "WebsocketHandler exit")
}
