# syntax=docker/dockerfile:1

FROM golang:1.19-alpine as builder

WORKDIR /app

COPY message-store ./message-store

WORKDIR /app/message-store

RUN go build -o /go/bin/app github.com/openline-ai/openline-customer-os/packages/server/message-store

FROM alpine:3.14
COPY --chown=65534:65534 --from=builder /go/bin/app .
USER 65534

ENTRYPOINT [ "./app" ]
