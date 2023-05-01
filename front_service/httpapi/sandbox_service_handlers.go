package httpapi

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/YFatMR/go_messenger/front_server/grpcapi"
	"github.com/YFatMR/go_messenger/protocol/pkg/proto"
	"github.com/gorilla/mux"
	"google.golang.org/protobuf/encoding/protojson"
)

func (s *FrontServer) GetProgramByID(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	ctx := r.Context()

	// Json data to protobuf.
	vars := mux.Vars(r)
	request := proto.ProgramID{
		ID: vars["ID"],
	}

	ctx = context.WithValue(ctx, grpcapi.RequireAuthorizationField{}, true)
	ctx = context.WithValue(ctx, grpcapi.AuthorizationField{}, r.Header.Get("Authorization"))
	response, err := s.sandboxServiceClient.GetProgramByID(ctx, &request)
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

func (s *FrontServer) CreateProgram(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	ctx := r.Context()
	decoder := json.NewDecoder(r.Body)

	// Json data to protobuf.
	var request proto.ProgramSource
	err := decoder.Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx = context.WithValue(ctx, grpcapi.RequireAuthorizationField{}, true)
	ctx = context.WithValue(ctx, grpcapi.AuthorizationField{}, r.Header.Get("Authorization"))
	response, err := s.sandboxServiceClient.CreateProgram(ctx, &request)
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

func (s *FrontServer) UpdateProgramSource(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	ctx := r.Context()
	decoder := json.NewDecoder(r.Body)

	// Json data to protobuf.
	var request proto.UpdateProgramSourceRequest
	err := decoder.Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx = context.WithValue(ctx, grpcapi.RequireAuthorizationField{}, true)
	ctx = context.WithValue(ctx, grpcapi.AuthorizationField{}, r.Header.Get("Authorization"))
	response, err := s.sandboxServiceClient.UpdateProgramSource(ctx, &request)
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

func (s *FrontServer) RunProgram(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	ctx := r.Context()

	// Json data to protobuf.
	vars := mux.Vars(r)
	request := proto.ProgramID{
		ID: vars["ID"],
	}

	ctx = context.WithValue(ctx, grpcapi.RequireAuthorizationField{}, true)
	ctx = context.WithValue(ctx, grpcapi.AuthorizationField{}, r.Header.Get("Authorization"))
	response, err := s.sandboxServiceClient.RunProgram(ctx, &request)
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

func (s *FrontServer) LintProgram(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	ctx := r.Context()

	// Json data to protobuf.
	vars := mux.Vars(r)
	request := proto.ProgramID{
		ID: vars["ID"],
	}

	ctx = context.WithValue(ctx, grpcapi.RequireAuthorizationField{}, true)
	ctx = context.WithValue(ctx, grpcapi.AuthorizationField{}, r.Header.Get("Authorization"))
	response, err := s.sandboxServiceClient.LintProgram(ctx, &request)
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
