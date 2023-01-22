FROM golang:1.19-alpine as build
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./auth_service_bin ./auth_service/cmd

FROM scratch
WORKDIR /app
COPY --from=build /app/auth_service_bin auth_service_bin
CMD ["/app/auth_service_bin"]
