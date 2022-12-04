CURRENT_DIR=$(shell pwd)
BINARY_DIR=${CURRENT_DIR}/_bin

create_binary_dir:
	mkdir -p ${BINARY_DIR}

build: create_binary_dir
	go build -o ${BINARY_DIR}/front_service ${CURRENT_DIR}/front_service/cmd
	go build -o ${BINARY_DIR}/user_service ${CURRENT_DIR}/user_service/cmd

build_docker_compose: build
	sudo docker-compose build --no-cache

run: build_docker_compose
	sudo docker-compose --env-file ${CURRENT_DIR}/.env --verbose up  --force-recreate --remove-orphans

run-tests:
	go test ${CURRENT_DIR}/user_service/internal/repositories/mongo/ -v