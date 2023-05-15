FROM golang:1.19-alpine as build
WORKDIR /app

COPY ./sandbox_worker_service/go.mod ./sandbox_worker_service/go.mod
COPY ./sandbox_worker_service/go.sum ./sandbox_worker_service/go.sum

COPY ./core/go.mod ./core/go.mod
COPY ./core/go.sum ./core/go.sum

RUN go work init && \
    go work use ./sandbox_worker_service && \
    go work use ./core && \
    go work sync && \
    cd /app/core && go mod download && \
    cd /app/sandbox_worker_service && go mod download

COPY ./core ./core
COPY ./sandbox_worker_service ./sandbox_worker_service

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./sandbox_worker_service_bin ./sandbox_worker_service/cmd


FROM python:3.9-slim-buster

ENV PYTHONUNBUFFERED=1
ENV PATH=/root/.local/bin:$PATH

# Variables for command execution

ENV EXECUTION_COMMAND=python
ENV EXECUTION_COMMAND_ARGS="-c"

# Run
WORKDIR /app
COPY --from=build /app/sandbox_worker_service_bin sandbox_worker_service_bin
CMD ["/app/sandbox_worker_service_bin"]
