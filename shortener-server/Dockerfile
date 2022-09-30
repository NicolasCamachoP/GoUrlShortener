# syntax=docker/dockerfile:1
FROM golang:1.18-alpine

ADD . /app
WORKDIR /app

ENV SERVER_PORT_NUMBER=3333

EXPOSE 3333

RUN go mod download
RUN go build

CMD ["/app/shortener-server"]