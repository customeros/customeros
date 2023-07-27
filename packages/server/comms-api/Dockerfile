# syntax=docker/dockerfile:1

FROM golang:1.19-alpine as builder
RUN apk add --update make bash

WORKDIR /app
COPY customer-os-common-module/. ./customer-os-common-module/
COPY comms-api/. ./comms-api

WORKDIR /app/comms-api
RUN go build -o bin/comms-api

FROM alpine:3.14
COPY --chown=65534:65534 --from=builder /app/comms-api/bin/comms-api .
USER 65534

ENTRYPOINT [ "./comms-api" ]
