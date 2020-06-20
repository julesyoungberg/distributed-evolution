FROM golang:latest

WORKDIR /app

COPY ./go /app/go

COPY ./go.mod ./go.sum ./

RUN go mod download

# RUN go get github.com/githubnemo/CompileDaemon

# ENTRYPOINT CompileDaemon --build="go build go/main/worker.go" --command=./worker

ENTRYPOINT go run go/main/worker.go
