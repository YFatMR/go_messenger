ROOT_PROJECT_DIRECTORY=$(shell pwd)
BINARY_DIRECTORY=${ROOT_PROJECT_DIRECTORY}/_bin
GO_WORK_SUBDIRECTORIES=$(shell go work edit -json | jq -c -r '[.Use[].DiskPath] | map_values(. + "/...")[]')
GO_LINTER_BINARY=$(shell go env GOPATH)/bin/golangci-lint
GO_WRAP_BINARY=$(shell go env GOPATH)/bin/gowrap
GOFUMP_BINATY=$(shell go env GOPATH)/bin/gofumpt
GENERATED_DECORATORS_EXTENTION=.gen.go
GO_SOURCES=$(shell find . -type f \( -iname "*.go" ! -iname "*.pb.go" ! -iname "*.gw.go" ! -iname "*.gen.go" ! -iname "*.template.go" \))

create_binary_directory:
	mkdir -p ${BINARY_DIRECTORY}

gen:
	${ROOT_PROJECT_DIRECTORY}/protocol/generate.sh

gen_interfaces:
	${ROOT_PROJECT_DIRECTORY}/scripts/generate_interfaces.sh

raw_build:
	go build -o ${BINARY_DIRECTORY}/front_service ${ROOT_PROJECT_DIRECTORY}/front_service/cmd
	go build -o ${BINARY_DIRECTORY}/user_service ${ROOT_PROJECT_DIRECTORY}/user_service/cmd
	go build -o ${BINARY_DIRECTORY}/qa_test ${ROOT_PROJECT_DIRECTORY}/qa/test
	go build -o ${BINARY_DIRECTORY}/sandbox_service ${ROOT_PROJECT_DIRECTORY}/sandbox_service/cmd

build: gen
	make raw_build

build_docker_compose: build
	sudo docker-compose --file ${ROOT_PROJECT_DIRECTORY}/docker-compose.yml --env-file ${ROOT_PROJECT_DIRECTORY}/production.env build

run: build_docker_compose
	sudo docker-compose --file ${ROOT_PROJECT_DIRECTORY}/docker-compose.yml --env-file ${ROOT_PROJECT_DIRECTORY}/production.env --verbose up

run-tests:
	go test -tags=test -v -race -o ${BINARY_DIRECTORY}/user_service_repository_tests ${ROOT_PROJECT_DIRECTORY}/user_service/user/ -args -mongo_config_path="${ROOT_PROJECT_DIRECTORY}/core/pkg/recipes/go/mongo/.env"

# up compose first:
# sudo docker-compose --file ./docker-compose-test.yml --env-file ./test.env build
# sudo docker-compose --file ./docker-compose-test.yml --env-file ./test.env up
# ./scripts/create_kafka_topics.sh
run-huge-tests:
	go test -tags=test -v -race -o ${BINARY_DIRECTORY}/qa_test ${ROOT_PROJECT_DIRECTORY}/qa/test/ -args -env-file="${ROOT_PROJECT_DIRECTORY}/test.env"

run-bench:
	go test -v -race ${ROOT_PROJECT_DIRECTORY}/user_service/userservice -bench=.

run-huge-bench:
	go test -tags=bench -v -race -o ${BINARY_DIRECTORY}/qa_test_bench -bench=. ${ROOT_PROJECT_DIRECTORY}/qa/test/ -args -env-file="${ROOT_PROJECT_DIRECTORY}/test.env"

lint:
	${GO_LINTER_BINARY} run -v --config ${ROOT_PROJECT_DIRECTORY}/.golangci.yml -- ${GO_WORK_SUBDIRECTORIES}

fump-diff:
	${GOFUMP_BINATY} -d ${GO_SOURCES}

fump-write:
	${GOFUMP_BINATY} -w ${GO_SOURCES}

update-sandbox-image:
	docker build --tag sandbox ${ROOT_PROJECT_DIRECTORY}/sandbox_service/sandbox/
	docker save -o ${ROOT_PROJECT_DIRECTORY}/sandbox_service/docker/images/sandbox.tar sandbox:latest


# docker build --tag local_sandbox ${ROOT_PROJECT_DIRECTORY}/sandbox_service/sandbox/
# docker image tag local_sandbox maximfafa/sandbox2:latest
# docker image push maximfafa/sandbox2

# cover:
# 	go test -short -count=1 -race -coverprofile=coverage.out ./...
# 	go tool cover -html=coverage.out
# 	rm coverage.out