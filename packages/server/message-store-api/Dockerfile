# syntax=docker/dockerfile:1

## Build
FROM golang:1.19.2-alpine3.16 AS builder

WORKDIR /workspace

RUN mkdir /workspace/build
RUN mkdir /workspace/build/message-store-api

#copy and build common module
RUN mkdir /workspace/build/customer-os-common-module

COPY customer-os-common-module/go.mod /workspace/build/customer-os-common-module/go.mod
COPY customer-os-common-module/go.sum /workspace/build/customer-os-common-module/go.sum

RUN cd /workspace/build/customer-os-common-module && go mod download && go mod verify

COPY customer-os-common-module/. /workspace/build/customer-os-common-module/
RUN cd /workspace/build/customer-os-common-module && go build

#copy and build common module
RUN mkdir /workspace/build/customer-os-api

COPY customer-os-api/go.mod /workspace/build/customer-os-api/go.mod
COPY customer-os-api/go.sum /workspace/build/customer-os-api/go.sum

RUN cd /workspace/build/customer-os-api && go mod download && go mod verify

COPY customer-os-api/. /workspace/build/customer-os-api/
RUN cd /workspace/build/customer-os-api && go build

#copy and build message store api
COPY message-store-api/go.mod /workspace/build/message-store-api/go.mod
COPY message-store-api/go.sum /workspace/build/message-store-api/go.sum

RUN cd /workspace/build/message-store-api && go mod download && go mod verify

COPY message-store-api/. /workspace/build/message-store-api/

RUN cd /workspace/build/message-store-api && go build -v -o /go/bin/app

FROM alpine:3.14
COPY --chown=65534:65534 --from=builder /go/bin/app .
COPY --chown=65534:65534 --from=builder /workspace/build/message-store-api/.env .env
USER 65534

ENTRYPOINT [ "./app" ]