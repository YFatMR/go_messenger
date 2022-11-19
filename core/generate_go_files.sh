#!/_bin/bash

# shellcheck disable=SC2164
script_path="$( cd -- "$(dirname "$0")" >/dev/null 2>&1 ; pwd -P )"

export GO_PATH=/usr/local/go
export PATH=$PATH:/$GO_PATH/_bin




cd "${script_path}"

# buf
buf_path="/home/linuxbrew/.linuxbrew/bin/buf"
# update deps (when edit buf.yaml)
${buf_path} mod update
# generate proto files
#proto_path="${script_path}/pkg/proto"
#cd "${proto_path}"
${buf_path} generate --template "buf.gen.yaml" #--verbose # --output "../.." #--config ./buf.yaml

#${buf_path} build


#proto_path="${script_path}/pkg/proto"
#generated_proto_path="${script_path}/pkg/generated"
#protoc -I "${proto_path}" --go_out="${generated_proto_path}" --go_opt=paths=source_relative --go-grpc_out="${generated_proto_path}" --go-grpc_opt=paths=source_relative \
#  user.proto \
#  front.proto
