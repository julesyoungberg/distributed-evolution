FROM golang:latest

WORKDIR /app

COPY ./go /app/go

COPY ./go.mod ./go.sum ./

RUN go mod download

ENTRYPOINT go run go/main/worker/worker.go
