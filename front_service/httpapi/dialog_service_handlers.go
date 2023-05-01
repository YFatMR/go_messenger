package httpapi

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/YFatMR/go_messenger/front_server/grpcapi"
	"github.com/YFatMR/go_messenger/protocol/pkg/proto"
	"github.com/gorilla/mux"
	"google.golang.org/protobuf/encoding/protojson"
)

func (s *FrontServer) CreateDialogWith(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	ctx := r.Context()

	// Json data to protobuf.
	userID, err := strconv.ParseUint(r.URL.Query().Get("userID"), 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	request := proto.UserID{
		ID: userID,
	}

	ctx = context.WithValue(ctx, grpcapi.RequireAuthorizationField{}, true)
	ctx = context.WithValue(ctx, grpcapi.AuthorizationField{}, r.Header.Get("Authorization"))
	response, err := s.dialogServiceClient.CreateDialogWith(ctx, &request)
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

func (s *FrontServer) GetDialogs(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	ctx := r.Context()

	// Json data to protobuf.
	limit, err := strconv.ParseUint(r.URL.Query().Get("limit"), 10, 64)
	if err != nil {
		http.Error(w, "limit params error:"+err.Error(), http.StatusBadRequest)
		return
	}

	offset, err := strconv.ParseUint(r.URL.Query().Get("offset"), 10, 64)
	if err != nil {
		http.Error(w, "offset params error:"+err.Error(), http.StatusBadRequest)
		return
	}
	request := proto.GetDialogsRequest{
		Offset: offset,
		Limit:  limit,
	}

	ctx = context.WithValue(ctx, grpcapi.RequireAuthorizationField{}, true)
	ctx = context.WithValue(ctx, grpcapi.AuthorizationField{}, r.Header.Get("Authorization"))
	response, err := s.dialogServiceClient.GetDialogs(ctx, &request)
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

type TextMessage struct {
	Text string `json:"text,omitempty"`
}

func (s *FrontServer) CreateDialogMessage(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	ctx := r.Context()
	decoder := json.NewDecoder(r.Body)

	// Json data to protobuf.
	var message TextMessage
	err := decoder.Decode(&message)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	vars := mux.Vars(r)
	dialogID, err := strconv.ParseUint(vars["ID"], 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	request := proto.CreateDialogMessageRequest{
		Text: message.Text,
		DialogID: &proto.DialogID{
			ID: dialogID,
		},
	}

	ctx = context.WithValue(ctx, grpcapi.RequireAuthorizationField{}, true)
	ctx = context.WithValue(ctx, grpcapi.AuthorizationField{}, r.Header.Get("Authorization"))
	response, err := s.dialogServiceClient.CreateDialogMessage(ctx, &request)
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

func (s *FrontServer) GetDialogMessages(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	ctx := r.Context()

	// Json data to protobuf.
	limit, err := strconv.ParseUint(r.URL.Query().Get("limit"), 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	offset, err := strconv.ParseUint(r.URL.Query().Get("offset"), 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	vars := mux.Vars(r)
	dialogID, err := strconv.ParseUint(vars["ID"], 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	request := proto.GetDialogMessagesRequest{
		DialogID: &proto.DialogID{
			ID: dialogID,
		},
		Limit:  limit,
		Offset: offset,
	}

	ctx = context.WithValue(ctx, grpcapi.RequireAuthorizationField{}, true)
	ctx = context.WithValue(ctx, grpcapi.AuthorizationField{}, r.Header.Get("Authorization"))
	response, err := s.dialogServiceClient.GetDialogMessages(ctx, &request)
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
