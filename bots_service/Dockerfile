FROM golang:1.19-alpine as build
WORKDIR /app

COPY ./bots_service/go.mod ./bots_service/go.mod
COPY ./bots_service/go.sum ./bots_service/go.sum

COPY ./core/go.mod ./core/go.mod
COPY ./core/go.sum ./core/go.sum

COPY ./protocol/go.mod ./protocol/go.mod
COPY ./protocol/go.sum ./protocol/go.sum

RUN go work init && \
    go work use ./bots_service && \
    go work use ./core && \
    go work use ./protocol && \
    go work sync && \
    cd /app/protocol && go mod download && \
    cd /app/core && go mod download && \
    cd /app/bots_service && go mod download

COPY ./core ./core
COPY ./protocol ./protocol
COPY ./bots_service ./bots_service

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bots_service_bin ./bots_service/cmd

FROM scratch
WORKDIR /app
COPY --from=build /app/bots_service_bin bots_service_bin
CMD ["/app/bots_service_bin"]
