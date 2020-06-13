FROM golang:latest

WORKDIR /app

COPY ./ /app

COPY ./go.mod ./go.sum ./

RUN go mod download

# RUN mkdir /shared
# VOLUME /shared

EXPOSE 3000

WORKDIR /app/master

CMD ["go" "run" "master.go"]
