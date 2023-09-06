# syntax=docker/dockerfile:1

## Build
FROM golang:1.20.1-alpine3.17 AS builder

WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download && go mod verify

COPY . .

RUN go build -v -o /go/bin/app

FROM alpine:3.14
COPY --chown=65534:65534 --from=builder /go/bin/app .
USER 65534

EXPOSE 10001
ENTRYPOINT [ "./app" ]