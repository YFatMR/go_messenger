#!/bin/bash

# output
GO_OUT="./pkg/proto"
GO_GRPC_OUT="./pkg/proto"
GO_GRPC_GATEWAY_OUT="./pkg/proto"
OPEN_API_V2_OUT="./open_api_v2"

# deps
GOOGLE_API_PATH=$GOPATH/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.16.0/third_party/googleapis

# generate go files
protoc -I "./internal" \
       -I"${GOOGLE_API_PATH}" \
    --go_out "${GO_OUT}" --go_opt paths=source_relative \
    --go-grpc_out "${GO_GRPC_OUT}" --go-grpc_opt paths=source_relative \
    --grpc-gateway_out "${GO_GRPC_GATEWAY_OUT}" \
    --grpc-gateway_opt logtostderr=true \
    --grpc-gateway_opt paths=source_relative \
    internal/front.proto \
    internal/user.proto

# generate openapiv2 for REST endpoint
protoc -I "./internal" \
       -I"${GOOGLE_API_PATH}" \
    --openapiv2_out "${OPEN_API_V2_OUT}" \
    --openapiv2_opt logtostderr=true \
    internal/front.proto
