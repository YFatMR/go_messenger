package httpapi

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/YFatMR/go_messenger/front_server/grpcapi"
	"github.com/YFatMR/go_messenger/protocol/pkg/proto"
	"google.golang.org/protobuf/encoding/protojson"
)

type BotMessage struct {
	Text string `json:"text,omitempty"`
}

func (s *FrontServer) GetBotMessageCompletion(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	ctx := r.Context()
	decoder := json.NewDecoder(r.Body)

	// Json data to protobuf.
	var message BotMessage
	err := decoder.Decode(&message)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	request := proto.BotMessage{
		Text: message.Text,
	}

	ctx = context.WithValue(ctx, grpcapi.RequireAuthorizationField{}, true)
	ctx = context.WithValue(ctx, grpcapi.AuthorizationField{}, r.Header.Get("Authorization"))
	response, err := s.botsServiceClient.GetBotMessageCompletion(ctx, &request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Response protobuf to json.
	bytes, err := protojson.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(bytes)
}
