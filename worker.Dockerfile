FROM golang:latest as builder

WORKDIR /app

COPY ./go /app/go

COPY ./go.mod ./go.sum ./

RUN go mod download

# ENTRYPOINT go run go/main/worker/worker.go

RUN go build go/main/worker/worker.go

FROM golang:latest

COPY --from=builder ./app/worker .

CMD ["./worker"]
