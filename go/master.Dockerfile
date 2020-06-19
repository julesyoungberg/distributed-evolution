FROM golang:latest

WORKDIR /app

COPY ./ /app

COPY ./go.mod ./go.sum ./

RUN go mod download

# expose http port
EXPOSE 9000 

# expose rpc port
EXPOSE 8080

# CMD ["go" "run" "main/master.go"]

RUN go get github.com/githubnemo/CompileDaemon

ENTRYPOINT CompileDaemon --build="go build commands/run_master.go" --command=./run_master
