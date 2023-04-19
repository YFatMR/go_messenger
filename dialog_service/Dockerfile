FROM golang:1.19-alpine as build
WORKDIR /app

COPY ./dialog_service/go.mod ./dialog_service/go.mod
COPY ./dialog_service/go.sum ./dialog_service/go.sum

COPY ./core/go.mod ./core/go.mod
COPY ./core/go.sum ./core/go.sum

COPY ./protocol/go.mod ./protocol/go.mod
COPY ./protocol/go.sum ./protocol/go.sum

RUN go work init && \
    go work use ./dialog_service && \
    go work use ./core && \
    go work use ./protocol && \
    go work sync && \
    cd /app/protocol && go mod download && \
    cd /app/core && go mod download && \
    cd /app/dialog_service && go mod download

COPY ./core ./core
COPY ./protocol ./protocol
COPY ./dialog_service ./dialog_service

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./dialog_service_bin ./dialog_service/cmd

FROM scratch
WORKDIR /app
COPY --from=build /app/dialog_service_bin dialog_service_bin
CMD ["/app/dialog_service_bin"]
