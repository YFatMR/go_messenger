FROM golang:1.19-alpine as build
WORKDIR /app

COPY ./sandbox_service/go.mod ./sandbox_service/go.mod
COPY ./sandbox_service/go.sum ./sandbox_service/go.sum

COPY ./core/go.mod ./core/go.mod
COPY ./core/go.sum ./core/go.sum

COPY ./protocol/go.mod ./protocol/go.mod
COPY ./protocol/go.sum ./protocol/go.sum

RUN go work init && \
    go work use ./sandbox_service && \
    go work use ./core && \
    go work use ./protocol && \
    go work sync && \
    cd /app/protocol && go mod download && \
    cd /app/core && go mod download && \
    cd /app/sandbox_service && go mod download

COPY ./core ./core
COPY ./protocol ./protocol
COPY ./sandbox_service ./sandbox_service

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./sandbox_service_bin ./sandbox_service/cmd

FROM scratch
WORKDIR /app
COPY --from=build /app/sandbox_service_bin sandbox_service_bin
CMD ["/app/sandbox_service_bin"]
