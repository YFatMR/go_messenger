ROOT_PROJECT_DIRECTORY=$(shell pwd)
BINARY_DIRECTORY=${ROOT_PROJECT_DIRECTORY}/_bin
GO_WORK_SUBDIRECTORIES=$(shell go work edit -json | jq -c -r '[.Use[].DiskPath] | map_values(. + "/...")[]')
GO_LINTER_BINARY=$(shell go env GOPATH)/bin/golangci-lint
GO_WRAP_BINARY=$(shell go env GOPATH)/bin/gowrap
GOFUMP_BINATY=$(shell go env GOPATH)/bin/gofumpt
GO_DECORATORS_TEMPLATE_DIRECTORY=${ROOT_PROJECT_DIRECTORY}/core/pkg/decorators/templates/
GENERATED_DECORATORS_EXTENTION=".gen.go"
GO_SOURCES=$(shell find . -type f \( -iname "*.go" ! -iname "*.pb.go" ! -iname "*.gw.go" \))

create_binary_directory:
	mkdir -p ${BINARY_DIRECTORY}

gen:
	${ROOT_PROJECT_DIRECTORY}/protocol/generate.sh

gen_interfaces:
	${GO_WRAP_BINARY} gen \
		-p github.com/YFatMR/go_messenger/user_service/internal/repositories \
		-i UserRepository \
		-t ${GO_DECORATORS_TEMPLATE_DIRECTORY}/common/loggers.go \
		-o ${ROOT_PROJECT_DIRECTORY}/user_service/internal/repositories/decorators/loggers${GENERATED_DECORATORS_EXTENTION} \

	${GO_WRAP_BINARY} gen \
		-p github.com/YFatMR/go_messenger/user_service/internal/repositories \
		-i UserRepository \
		-t ${GO_DECORATORS_TEMPLATE_DIRECTORY}/common/opentelemetry_tracing.go \
		-o ${ROOT_PROJECT_DIRECTORY}/user_service/internal/repositories/decorators/opentelemetry_tracing${GENERATED_DECORATORS_EXTENTION}
	${GO_WRAP_BINARY} gen \
		-p github.com/YFatMR/go_messenger/user_service/internal/repositories \
		-i UserRepository \
		-t ${GO_DECORATORS_TEMPLATE_DIRECTORY}/repositories/prometheus_metrics.go \
		-o ${ROOT_PROJECT_DIRECTORY}/user_service/internal/repositories/decorators/prometheus_metrics${GENERATED_DECORATORS_EXTENTION}

	${GO_WRAP_BINARY} gen \
		-p github.com/YFatMR/go_messenger/user_service/internal/services \
		-i UserService \
		-t ${GO_DECORATORS_TEMPLATE_DIRECTORY}/common/loggers.go \
		-o ${ROOT_PROJECT_DIRECTORY}/user_service/internal/services/decorators/loggers${GENERATED_DECORATORS_EXTENTION}
	${GO_WRAP_BINARY} gen \
		-p github.com/YFatMR/go_messenger/user_service/internal/services \
		-i UserService \
		-t ${GO_DECORATORS_TEMPLATE_DIRECTORY}/common/opentelemetry_tracing.go \
		-o ${ROOT_PROJECT_DIRECTORY}/user_service/internal/services/decorators/opentelemetry_tracing${GENERATED_DECORATORS_EXTENTION}

	${GO_WRAP_BINARY} gen \
		-p github.com/YFatMR/go_messenger/user_service/internal/controllers \
		-i UserController \
		-t ${GO_DECORATORS_TEMPLATE_DIRECTORY}/common/loggers.go \
		-o ${ROOT_PROJECT_DIRECTORY}/user_service/internal/controllers/decorators/loggers${GENERATED_DECORATORS_EXTENTION}

raw_build:
	go build -o ${BINARY_DIRECTORY}/auth_service ${ROOT_PROJECT_DIRECTORY}/auth_service/cmd
	go build -o ${BINARY_DIRECTORY}/front_service ${ROOT_PROJECT_DIRECTORY}/front_service/cmd
	go build -o ${BINARY_DIRECTORY}/user_service ${ROOT_PROJECT_DIRECTORY}/user_service/cmd
	go build -o ${BINARY_DIRECTORY}/user_service ${ROOT_PROJECT_DIRECTORY}/qa/test

build: gen
	make raw_build

build_docker_compose: build
	sudo docker-compose build

run: build_docker_compose
	sudo docker-compose --env-file ${ROOT_PROJECT_DIRECTORY}/production.env --verbose up

run-tests:
	go test -v -race -o ${BINARY_DIRECTORY}/user_service_repository_tests ${ROOT_PROJECT_DIRECTORY}/user_service/internal/repositories/mongorepository/ -args -mongo_config_path="${ROOT_PROJECT_DIRECTORY}/core/pkg/recipes/go/mongo/.env"

run-huge-tests:
	go test -v -race ${ROOT_PROJECT_DIRECTORY}/qa/test/ -args -docker-compose-file="${ROOT_PROJECT_DIRECTORY}/docker-compose-test.yml" -env-file="${ROOT_PROJECT_DIRECTORY}/test.env"

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