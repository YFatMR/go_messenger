FROM golang:1.19-alpine as build
WORKDIR /app

COPY ./auth_service/go.mod ./auth_service/go.mod
COPY ./auth_service/go.sum ./auth_service/go.sum

COPY ./core/go.mod ./core/go.mod
COPY ./core/go.sum ./core/go.sum

COPY ./protocol/go.mod ./protocol/go.mod
COPY ./protocol/go.sum ./protocol/go.sum

RUN go work init && \
    go work use ./auth_service && \
    go work use ./core && \
    go work use ./protocol && \
    go work sync && \
    cd /app/protocol && go mod download && \
    cd /app/core && go mod download && \
    cd /app/auth_service && go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./auth_service_bin ./auth_service/cmd

FROM scratch
WORKDIR /app
COPY --from=build /app/auth_service_bin auth_service_bin
CMD ["/app/auth_service_bin"]
