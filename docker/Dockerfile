FROM golang:1.19

WORKDIR /usr/src/app

RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
RUN go install github.com/go-delve/delve/cmd/dlv@latest
RUN go install github.com/swaggo/swag/cmd/swag@latest

RUN apt-get update
RUN apt-get install -y net-tools lsof gcc libc6 libc-bin
RUN apt-get clean && rm -rf /var/lib/apt/lists/*
