FROM julesyoungberg/gocv AS gocv

WORKDIR /app

COPY ./go /app/go
COPY ./go.mod ./go.sum ./
RUN go mod download

EXPOSE 9000 
EXPOSE 8080

ENTRYPOINT go run go/main/master/master.go

# RUN go build go/main/master/master.go

# FROM julesyoungberg/gocv
# COPY --from=gocv ./app .
# CMD ["./master"]