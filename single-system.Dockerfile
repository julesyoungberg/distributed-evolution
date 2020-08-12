FROM julesyoungberg/gocv

WORKDIR /app

COPY ./go /app/go
COPY ./go.mod ./go.sum ./
RUN go mod download

EXPOSE 9000

ENTRYPOINT go run go/main/single_system/single_system.go
