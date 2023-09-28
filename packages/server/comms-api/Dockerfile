# syntax=docker/dockerfile:1

## Build
FROM golang:1.20.1-alpine3.17 AS builder

WORKDIR /workspace

RUN mkdir /workspace/build
RUN mkdir /workspace/build/comms-api
RUN mkdir /workspace/build/customer-os-common-module
RUN mkdir /workspace/build/customer-os-common-auth

#copy and build common module
COPY customer-os-common-module/go.mod /workspace/build/customer-os-common-module/go.mod
COPY customer-os-common-module/go.sum /workspace/build/customer-os-common-module/go.sum
RUN cd /workspace/build/customer-os-common-module && go mod download && go mod verify

COPY customer-os-common-module/. /workspace/build/customer-os-common-module/
RUN cd /workspace/build/customer-os-common-module && go build

#copy and build common auth module
COPY customer-os-common-auth/go.mod /workspace/build/customer-os-common-auth/go.mod
COPY customer-os-common-auth/go.sum /workspace/build/customer-os-common-auth/go.sum
RUN cd /workspace/build/customer-os-common-auth && go mod download && go mod verify

COPY customer-os-common-auth/. /workspace/build/customer-os-common-auth/
RUN cd /workspace/build/customer-os-common-auth && go build

#copy and build comms-api
COPY comms-api/go.mod /workspace/build/comms-api/go.mod
COPY comms-api/go.sum /workspace/build/comms-api/go.sum

RUN cd /workspace/build/comms-api && go mod download && go mod verify

COPY comms-api/. /workspace/build/comms-api/

RUN cd /workspace/build/comms-api && go build -v -o /go/bin/app

FROM alpine:3.14
COPY --chown=65534:65534 --from=builder /go/bin/app .
COPY --chown=65534:65534 --from=builder /workspace/build/comms-api/.env .env
USER 65534

ENTRYPOINT [ "./app" ]