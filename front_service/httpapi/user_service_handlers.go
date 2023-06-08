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

// Unsafe.
func (s *FrontServer) CreateUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	ctx := r.Context()
	decoder := json.NewDecoder(r.Body)

	// Json data to protobuf.
	var request proto.CreateUserRequest
	err := decoder.Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	request.Credential.Role = "user"

	grpcCtx := context.WithValue(ctx, grpcapi.RequireAuthorizationField{}, false)
	response, err := s.userServiceClient.CreateUser(grpcCtx, &request)
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

// Unsafe.
func (s *FrontServer) GenerateToken(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	ctx := r.Context()
	decoder := json.NewDecoder(r.Body)

	// Json data to protobuf.
	var request proto.Credential
	err := decoder.Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	request.Role = "user"

	grpcCtx := context.WithValue(ctx, grpcapi.RequireAuthorizationField{}, false)
	response, err := s.userServiceClient.GenerateToken(grpcCtx, &request)
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

func (s *FrontServer) GetUserByID(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	ctx := r.Context()

	// Json data to protobuf.
	vars := mux.Vars(r)
	userID, err := strconv.ParseUint(vars["ID"], 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	request := proto.UserID{
		ID: userID,
	}

	ctx = context.WithValue(ctx, grpcapi.RequireAuthorizationField{}, true)
	ctx = context.WithValue(ctx, grpcapi.AuthorizationField{}, r.Header.Get("Authorization"))
	response, err := s.userServiceClient.GetUserByID(ctx, &request)
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

func (s *FrontServer) UpdateUserData(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	ctx := r.Context()
	decoder := json.NewDecoder(r.Body)

	// Json data to protobuf.
	var request proto.UpdateUserDataRequest
	err := decoder.Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx = context.WithValue(ctx, grpcapi.RequireAuthorizationField{}, true)
	ctx = context.WithValue(ctx, grpcapi.AuthorizationField{}, r.Header.Get("Authorization"))
	response, err := s.userServiceClient.UpdateUserData(ctx, &request)
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

func (s *FrontServer) GetUsersByPrefix(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	ctx := r.Context()
	decoder := json.NewDecoder(r.Body)

	// Json data to protobuf.
	var request proto.GetUsersByPrefixRequest
	err := decoder.Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	queryStringLimit := r.URL.Query().Get("limit")
	if queryStringLimit == "" {
		http.Error(w, "no limit argument", http.StatusBadRequest)
		return
	}
	limit, err := strconv.ParseUint(queryStringLimit, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	request.Limit = limit

	ctx = context.WithValue(ctx, grpcapi.RequireAuthorizationField{}, true)
	ctx = context.WithValue(ctx, grpcapi.AuthorizationField{}, r.Header.Get("Authorization"))
	response, err := s.userServiceClient.GetUsersByPrefix(ctx, &request)
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
