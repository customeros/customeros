# syntax=docker/dockerfile:1

## Build
FROM golang:1.19.2-alpine3.16 AS builder

WORKDIR /workspace

RUN mkdir /workspace/build
RUN mkdir /workspace/build/message-store

#copy and build common module
RUN mkdir /workspace/build/customer-os-common-module

COPY customer-os-common-module/go.mod /workspace/build/customer-os-common-module/go.mod
COPY customer-os-common-module/go.sum /workspace/build/customer-os-common-module/go.sum

RUN cd /workspace/build/customer-os-common-module && go mod download && go mod verify

COPY customer-os-common-module/. /workspace/build/customer-os-common-module/
RUN cd /workspace/build/customer-os-common-module && go build

#copy and build file storage api
COPY message-store/go.mod /workspace/build/message-store/go.mod
COPY message-store/go.sum /workspace/build/message-store/go.sum

RUN cd /workspace/build/message-store && go mod download && go mod verify

COPY message-store/. /workspace/build/message-store/

RUN cd /workspace/build/message-store && go build -v -o /go/bin/app

FROM alpine:3.14
COPY --chown=65534:65534 --from=builder /go/bin/app .
COPY --chown=65534:65534 --from=builder /workspace/build/message-store/.env .env
USER 65534

ENTRYPOINT [ "./app" ]