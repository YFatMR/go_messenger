#!/bin/bash

# shellcheck disable=SC2164
script_path="$( cd -- "$(dirname "$0")" >/dev/null 2>&1 ; pwd -P )"

export GO_PATH=/usr/local/go
export PATH=$PATH:/$GO_PATH/bin

#protoc -I "${script_path}" user.proto --go_out=plugins=grpc:"${script_path}/generated"

#protoc -I "${script_path}" "${script_path}/user.proto" \
#  --go_out="${script_path}/generated" --go_opt=paths=source_relative \
#  --go-grpc_out="${script_path}/generated" --go-grpc_opt=paths=source_relative

cd "${script_path}"

protoc --go_out=./ --go_opt=paths=source_relative --go-grpc_out=./ --go-grpc_opt=paths=source_relative user.proto