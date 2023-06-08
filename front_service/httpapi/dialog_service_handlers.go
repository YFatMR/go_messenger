package httpapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
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

func (s *FrontServer) GetDialogByID(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	ctx := r.Context()

	// Json data to protobuf.

	vars := mux.Vars(r)
	dialogID, err := strconv.ParseUint(vars["ID"], 10, 64)
	if err != nil {
		http.Error(w, "dualogID params error:"+err.Error(), http.StatusBadRequest)
		return
	}
	request := proto.DialogID{
		ID: dialogID,
	}

	ctx = context.WithValue(ctx, grpcapi.RequireAuthorizationField{}, true)
	ctx = context.WithValue(ctx, grpcapi.AuthorizationField{}, r.Header.Get("Authorization"))
	response, err := s.dialogServiceClient.GetDialogByID(ctx, &request)
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
	Title string `json:"title,omitempty"`
	Text  string `json:"text,omitempty"`
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

	ctx = context.WithValue(ctx, grpcapi.RequireAuthorizationField{}, true)
	ctx = context.WithValue(ctx, grpcapi.AuthorizationField{}, r.Header.Get("Authorization"))

	if message.Title == "" {
		request := proto.CreateDialogMessageRequest{
			Text: message.Text,
			DialogID: &proto.DialogID{
				ID: dialogID,
			},
		}
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
		return
	}
	request := proto.CreateDialogMessageWithCodeRequest{
		Title: message.Title,
		Text:  message.Text,
		DialogID: &proto.DialogID{
			ID: dialogID,
		},
	}

	response, err := s.dialogServiceClient.CreateDialogMessageWithCode(ctx, &request)
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

	offsetTypeFromURL := func(url *url.URL) (proto.GetDialogMessagesRequest_OffsetType, error) {
		switch url.Query().Get("offset_type") {
		case "before":
			return proto.GetDialogMessagesRequest_BEFORE, nil
		case "before_include":
			return proto.GetDialogMessagesRequest_BEFORE_INCLUDE, nil
		case "after":
			return proto.GetDialogMessagesRequest_AFTER, nil
		case "after_include":
			return proto.GetDialogMessagesRequest_AFTER_INCLUDE, nil
		}
		return proto.GetDialogMessagesRequest_BEFORE, fmt.Errorf("unexpected offset_type params value")
	}

	// Json data to protobuf.
	limit, err := strconv.ParseUint(r.URL.Query().Get("limit"), 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	offsetType, err := offsetTypeFromURL(r.URL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)
	dialogID, err := strconv.ParseUint(vars["dialogID"], 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	messageID, err := strconv.ParseUint(vars["messageID"], 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	request := proto.GetDialogMessagesRequest{
		DialogID: &proto.DialogID{
			ID: dialogID,
		},
		MessageID: &proto.MessageID{
			ID: messageID,
		},
		Limit:      limit,
		OffsetType: offsetType,
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

func (s *FrontServer) ReadMessage(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	ctx := r.Context()

	// Json data to protobuf.
	vars := mux.Vars(r)
	dialogID, err := strconv.ParseUint(vars["dialogID"], 10, 64)
	if err != nil {
		http.Error(w, "dialogID params error:"+err.Error(), http.StatusBadRequest)
		return
	}

	messageID, err := strconv.ParseUint(vars["messageID"], 10, 64)
	if err != nil {
		http.Error(w, "dialogID params error:"+err.Error(), http.StatusBadRequest)
		return
	}

	request := proto.ReadMessageRequest{
		DialogID: &proto.DialogID{
			ID: dialogID,
		},
		MessageID: &proto.MessageID{
			ID: messageID,
		},
	}

	ctx = context.WithValue(ctx, grpcapi.RequireAuthorizationField{}, true)
	ctx = context.WithValue(ctx, grpcapi.AuthorizationField{}, r.Header.Get("Authorization"))
	response, err := s.dialogServiceClient.ReadMessage(ctx, &request)
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

func (s *FrontServer) CreateInstruction(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	ctx := r.Context()
	decoder := json.NewDecoder(r.Body)
	// title + messages id

	// Json data to protobuf.
	vars := mux.Vars(r)
	dialogID, err := strconv.ParseUint(vars["dialogID"], 10, 64)
	if err != nil {
		http.Error(w, "dialogID params error:"+err.Error(), http.StatusBadRequest)
		return
	}

	request := new(proto.CreateInstructionRequest)
	err = decoder.Decode(&request)
	if err != nil {
		http.Error(w, "message parse error:"+err.Error(), http.StatusBadRequest)
		return
	}
	request.DialogID = new(proto.DialogID)
	request.DialogID.ID = dialogID

	ctx = context.WithValue(ctx, grpcapi.RequireAuthorizationField{}, true)
	ctx = context.WithValue(ctx, grpcapi.AuthorizationField{}, r.Header.Get("Authorization"))
	response, err := s.dialogServiceClient.CreateInstruction(ctx, request)
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

func (s *FrontServer) GetInstructions(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	ctx := r.Context()

	offsetTypeFromURL := func(url *url.URL) (proto.GetInstructionsByIDRequest_OffsetType, error) {
		switch url.Query().Get("offset_type") {
		case "after":
			return proto.GetInstructionsByIDRequest_AFTER, nil
		}
		return proto.GetInstructionsByIDRequest_AFTER, fmt.Errorf("unexpected offset_type params value")
	}

	// Json data to protobuf.
	vars := mux.Vars(r)
	dialogID, err := strconv.ParseUint(vars["dialogID"], 10, 64)
	if err != nil {
		http.Error(w, "dialogID params error:"+err.Error(), http.StatusBadRequest)
		return
	}

	limit, err := strconv.ParseUint(r.URL.Query().Get("limit"), 10, 64)
	if err != nil || limit == 0 {
		http.Error(w, "limit params error:"+err.Error(), http.StatusBadRequest)
		return
	}

	instructionID := r.URL.Query().Get("instructionID")
	ctx = context.WithValue(ctx, grpcapi.RequireAuthorizationField{}, true)
	ctx = context.WithValue(ctx, grpcapi.AuthorizationField{}, r.Header.Get("Authorization"))

	if instructionID != "" {
		// use offset type by default after, not parse
		offsetType, err := offsetTypeFromURL(r.URL)
		if err != nil {
			http.Error(w, "offset_type params error:"+err.Error(), http.StatusBadRequest)
			return
		}

		parsedInstructionID, err := strconv.ParseUint(instructionID, 10, 64)
		if err != nil {
			http.Error(w, "parsedInstructionID params error:"+err.Error(), http.StatusBadRequest)
			return
		}
		request := &proto.GetInstructionsByIDRequest{
			DialogID: &proto.DialogID{
				ID: dialogID,
			},
			InstructionID: &proto.InstructionID{
				ID: parsedInstructionID,
			},
			Limit:      limit,
			OffsetType: offsetType,
		}
		response, err := s.dialogServiceClient.GetInstructionsByID(ctx, request)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		bytes, err := protojson.Marshal(response)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(bytes)
	} else {
		request := &proto.GetInstructionsRequest{
			DialogID: &proto.DialogID{
				ID: dialogID,
			},
			Limit: limit,
		}
		response, err := s.dialogServiceClient.GetInstructions(ctx, request)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		bytes, err := protojson.Marshal(response)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(bytes)
	}
}

func (s *FrontServer) GetLinks(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	ctx := r.Context()

	offsetTypeFromURL := func(url *url.URL) (proto.GetDialogLinksByIDRequest_OffsetType, error) {
		switch url.Query().Get("offset_type") {
		case "after":
			return proto.GetDialogLinksByIDRequest_AFTER, nil
		}
		return proto.GetDialogLinksByIDRequest_AFTER, fmt.Errorf("unexpected offset_type params value")
	}

	// Json data to protobuf.
	vars := mux.Vars(r)
	dialogID, err := strconv.ParseUint(vars["dialogID"], 10, 64)
	if err != nil {
		http.Error(w, "dialogID params error:"+err.Error(), http.StatusBadRequest)
		return
	}

	limit, err := strconv.ParseUint(r.URL.Query().Get("limit"), 10, 64)
	if err != nil || limit == 0 {
		http.Error(w, "limit params error:"+err.Error(), http.StatusBadRequest)
		return
	}
	ctx = context.WithValue(ctx, grpcapi.RequireAuthorizationField{}, true)
	ctx = context.WithValue(ctx, grpcapi.AuthorizationField{}, r.Header.Get("Authorization"))

	linkID := r.URL.Query().Get("linkID")
	if linkID != "" {
		// use offset type by default after, not parse
		offsetType, err := offsetTypeFromURL(r.URL)
		if err != nil {
			http.Error(w, "offset_type params error:"+err.Error(), http.StatusBadRequest)
			return
		}

		parsedLinkID, err := strconv.ParseUint(linkID, 10, 64)
		if err != nil {
			http.Error(w, "parsedInstructionID params error:"+err.Error(), http.StatusBadRequest)
			return
		}
		request := &proto.GetDialogLinksByIDRequest{
			DialogID: &proto.DialogID{
				ID: dialogID,
			},
			LinkID: &proto.LinkID{
				ID: parsedLinkID,
			},
			Limit:      limit,
			OffsetType: offsetType,
		}
		response, err := s.dialogServiceClient.GetDialogLinksByID(ctx, request)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		bytes, err := protojson.Marshal(response)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(bytes)
	} else {
		request := &proto.GetDialogLinksRequest{
			DialogID: &proto.DialogID{
				ID: dialogID,
			},
			Limit: limit,
		}
		response, err := s.dialogServiceClient.GetDialogLinks(ctx, request)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		bytes, err := protojson.Marshal(response)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(bytes)
	}
}

func (s *FrontServer) GetDialogMembers(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	ctx := r.Context()

	vars := mux.Vars(r)
	dialogID, err := strconv.ParseUint(vars["dialogID"], 10, 64)
	if err != nil {
		http.Error(w, "dialogID params error:"+err.Error(), http.StatusBadRequest)
		return
	}

	request := &proto.DialogID{
		ID: dialogID,
	}
	ctx = context.WithValue(ctx, grpcapi.RequireAuthorizationField{}, true)
	ctx = context.WithValue(ctx, grpcapi.AuthorizationField{}, r.Header.Get("Authorization"))
	response, err := s.dialogServiceClient.GetDialogMembers(ctx, request)
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

func (s *FrontServer) GetUnreadDialogMessagesCount(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	ctx := r.Context()

	vars := mux.Vars(r)
	dialogID, err := strconv.ParseUint(vars["dialogID"], 10, 64)
	if err != nil {
		http.Error(w, "dialogID params error:"+err.Error(), http.StatusBadRequest)
		return
	}

	request := &proto.DialogID{
		ID: dialogID,
	}
	ctx = context.WithValue(ctx, grpcapi.RequireAuthorizationField{}, true)
	ctx = context.WithValue(ctx, grpcapi.AuthorizationField{}, r.Header.Get("Authorization"))
	response, err := s.dialogServiceClient.GetUnreadDialogMessagesCount(ctx, request)
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
