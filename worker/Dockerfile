FROM golang:latest

WORKDIR /app

COPY ./ /app

# RUN mkdir /shared
# VOLUME /shared

ENTRYPOINT ["go" "run" "worker.go"]
