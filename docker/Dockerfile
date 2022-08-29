FROM golang:1.19-alpine

WORKDIR /usr/src/app

RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
RUN go install github.com/swaggo/swag/cmd/swag@latest

RUN apk update
RUN apk add net-tools lsof

RUN rm -rf /var/cache/apk/*
