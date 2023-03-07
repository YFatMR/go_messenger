FROM golang:1.19-alpine as build
WORKDIR /app

RUN go mod init sandbox && \
    go get github.com/spf13/viper && \
    go get github.com/fatih/color && \
    go mod download

ENTRYPOINT ["go", "run", "."]
