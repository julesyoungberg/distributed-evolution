FROM golang:latest

WORKDIR /app

COPY ./ /app

COPY ./go.mod ./go.sum ./

RUN go mod download

# RUN mkdir /shared
# VOLUME /shared

WORKDIR /app/worker

CMD ["go" "run" "worker/worker.go"]
