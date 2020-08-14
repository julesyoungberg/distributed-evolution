FROM julesyoungberg/gocv AS builder

WORKDIR /app

COPY ./go /app/go
COPY ./go.mod ./go.sum ./
RUN go mod download

# ENTRYPOINT go run go/main/master/master.go

RUN go build go/main/master/master.go

FROM julesyoungberg/gocv

COPY --from=builder ./app/master .

EXPOSE 9001
EXPOSE 8080

CMD ["./master"]
