FROM golang:latest

WORKDIR /app

COPY ./ /app

COPY ./go.mod ./go.sum ./

RUN go mod download

# expose http port
EXPOSE 9000 

# expose rpc port
EXPOSE 8080

WORKDIR /app/master

CMD ["go" "run" "."]
