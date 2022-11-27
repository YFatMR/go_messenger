FROM golang:1.19-alpine

RUN mkdir /messenger

COPY . /messenger

WORKDIR /messenger

RUN go build -o ./user_service_bin ./user_service/cmd && chmod +x ./user_service_bin

CMD ["/messenger/user_service_bin"]
