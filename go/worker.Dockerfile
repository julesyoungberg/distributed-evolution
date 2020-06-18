FROM golang:latest

WORKDIR /app

COPY ./ /app

COPY ./go.mod ./go.sum ./

RUN go mod download

CMD ["go" "run" "main/worker/worker.go"]
