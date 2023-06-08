#!/bin/bash

SCRIPT_DIRECTORY=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
cd $SCRIPT_DIRECTORY

# output
GO_OUT="./pkg/proto"
GO_GRPC_OUT="./pkg/proto"
GO_GRPC_GATEWAY_OUT="./pkg/proto"
OPEN_API_V2_OUT="./open_api_v2"

export GOPATH=/home/am/go
# https://grpc.io/docs/languages/go/quickstart/
export PATH="$PATH:$(go env GOPATH)/bin"

# deps
GOOGLE_API_PATH=$GOPATH/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.16.0/third_party/googleapis

protoc-gen-go --version
# generate go files
protoc \
    -I "./internal" \
    -I"${GOOGLE_API_PATH}" \
    --go_out "${GO_OUT}" --go_opt paths=source_relative \
    --go-grpc_out "${GO_GRPC_OUT}" --go-grpc_opt paths=source_relative \
    --grpc-gateway_out "${GO_GRPC_GATEWAY_OUT}" \
    --grpc-gateway_opt logtostderr=true \
    --grpc-gateway_opt paths=source_relative \
    internal/front.proto \
    internal/user.proto \
    internal/common.proto \
    internal/sandbox.proto \
    internal/dialog.proto \
    internal/bots.proto
    # internal/auth.proto \

# generate openapiv2 for REST endpoint
# protoc \
#     -I "./internal" \
#     -I"${GOOGLE_API_PATH}" \
#     --openapiv2_out "${OPEN_API_V2_OUT}" \
#     --openapiv2_opt logtostderr=true \
#     internal/front.proto
