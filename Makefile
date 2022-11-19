CURRENT_DIR=$(shell pwd)
BINARY_DIR=${CURRENT_DIR}/_bin

create_binary_dir:
	mkdir -p ${BINARY_DIR}

build: create_binary_dir
	go build -o ${BINARY_DIR}/front_service ${CURRENT_DIR}/front_service/cmd
	go build -o ${BINARY_DIR}/user_service ${CURRENT_DIR}/user_service/cmd

run: build
	sudo docker-compose --env-file ${CURRENT_DIR}/.env up --force-recreate
