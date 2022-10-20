# syntax=docker/dockerfile:1

## Build
FROM golang:1.19.2-alpine3.16 AS builder

WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

RUN go build -o /go/bin/app

FROM alpine:3.14
COPY --chown=65534:65534 --from=builder /go/bin/app .
USER 65534

ENTRYPOINT [ "./app" ]