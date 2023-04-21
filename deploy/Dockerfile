FROM golang:1.20.0-buster

ENV GOPATH=/

WORKDIR /usr/src/app

COPY ./ ./

RUN go mod download
RUN GOOS=linux go build -o main cmd/main.go