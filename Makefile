ROOT_PROJECT_DIRECTORY=$(shell pwd)
BINARY_DIRECTORY=${ROOT_PROJECT_DIRECTORY}/_bin
GO_WORK_SUBDIRECTORIES=$(shell go work edit -json | jq -c -r '[.Use[].DiskPath] | map_values(. + "/...")[]')
GO_LINTER_BINARY=$(shell go env GOPATH)/bin/golangci-lint
GOFUMP_BINATY=$(shell go env GOPATH)/bin/gofumpt
GO_SOURCES=$(shell find . -type f \( -iname "*.go" ! -iname "*.pb.go" ! -iname "*.gw.go" \))

create_binary_directory:
	mkdir -p ${BINARY_DIRECTORY}

gen:
	${ROOT_PROJECT_DIRECTORY}/protocol/generate.sh

raw_build:
	go build -o ${BINARY_DIRECTORY}/auth_service ${ROOT_PROJECT_DIRECTORY}/auth_service/cmd
	go build -o ${BINARY_DIRECTORY}/front_service ${ROOT_PROJECT_DIRECTORY}/front_service/cmd
	go build -o ${BINARY_DIRECTORY}/user_service ${ROOT_PROJECT_DIRECTORY}/user_service/cmd

build: gen
	make raw_build

build_docker_compose: build
	sudo docker-compose build

run: build_docker_compose
	sudo docker-compose --env-file ${ROOT_PROJECT_DIRECTORY}/production.env --verbose up

run-tests:
	go test -v -race -o ${BINARY_DIRECTORY}/user_service_repository_tests ${ROOT_PROJECT_DIRECTORY}/user_service/internal/repositories/mongorepository/ -args -mongo_config_path="${ROOT_PROJECT_DIRECTORY}/core/pkg/recipes/go/mongo/.env"

run-huge-tests:
	go test -v -race -o ${BINARY_DIRECTORY}/user_service_repository_tests ${ROOT_PROJECT_DIRECTORY}/user_service/internal/repositories/mongorepository/ -args -mongo_config_path="${ROOT_PROJECT_DIRECTORY}/core/pkg/recipes/go/mongo/.env"

lint:
	${GO_LINTER_BINARY} run -v --config ${ROOT_PROJECT_DIRECTORY}/.golangci.yml -- ${GO_WORK_SUBDIRECTORIES}

fump-diff:
	${GOFUMP_BINATY} -d ${GO_SOURCES}

fump-write:
	${GOFUMP_BINATY} -w ${GO_SOURCES}

# cover:
# 	go test -short -count=1 -race -coverprofile=coverage.out ./...
# 	go tool cover -html=coverage.out
# 	rm coverage.out