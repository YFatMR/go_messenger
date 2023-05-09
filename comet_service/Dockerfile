FROM golang:1.19-alpine as build
WORKDIR /app

COPY ./comet_service/go.mod ./comet_service/go.mod
COPY ./comet_service/go.sum ./comet_service/go.sum

COPY ./core/go.mod ./core/go.mod
COPY ./core/go.sum ./core/go.sum

RUN go work init && \
    go work use ./comet_service && \
    go work use ./core && \
    go work sync && \
    cd /app/core && go mod download && \
    cd /app/comet_service && go mod download

COPY ./core ./core
COPY ./comet_service ./comet_service

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./comet_service_bin ./comet_service/cmd

FROM scratch
WORKDIR /app
COPY --from=build /app/comet_service_bin comet_service_bin
CMD ["/app/comet_service_bin"]
