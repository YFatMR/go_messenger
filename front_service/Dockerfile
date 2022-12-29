FROM golang:1.19-alpine as build
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./front_service_bin ./front_service/cmd

FROM scratch
WORKDIR /app
COPY --from=build /app/front_service_bin front_service_bin
CMD ["/app/front_service_bin"]
