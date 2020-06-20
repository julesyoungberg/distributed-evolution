FROM golang:latest

WORKDIR /app

COPY ./go /app/go

COPY ./go.mod ./go.sum ./

RUN go mod download

# expose http port
EXPOSE 9000 

# expose rpc port
EXPOSE 8080

# RUN go get github.com/githubnemo/CompileDaemon

# ENTRYPOINT CompileDaemon --build="go build go/main/master.go" --command=./master

ENTRYPOINT go run go/main/master.go
