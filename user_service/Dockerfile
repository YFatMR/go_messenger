FROM golang:1.19-alpine as build
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./user_service_bin ./user_service/cmd

FROM scratch
WORKDIR /app
COPY --from=build /app/user_service_bin user_service_bin
CMD ["/app/user_service_bin"]
