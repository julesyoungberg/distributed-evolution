FROM golang:latest

WORKDIR /app

COPY ./ /app

COPY ./go.mod ./go.sum ./

RUN go mod download

WORKDIR /app/worker

CMD ["go" "run" "."]
