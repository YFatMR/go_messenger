FROM golang:1.19-alpine

RUN mkdir /messenger

COPY . /messenger

WORKDIR /messenger

RUN go build -o ./front_service_bin ./front_service/cmd && chmod +x ./front_service_bin

CMD ["/messenger/front_service_bin"]
