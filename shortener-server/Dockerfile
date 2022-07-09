# syntax=docker/dockerfile:1

FROM golang:1.18-alpine

ADD . /app

WORKDIR /app

RUN go mod download

RUN go build main.go

RUN ls -a

EXPOSE 3333

CMD ["/app/main"]