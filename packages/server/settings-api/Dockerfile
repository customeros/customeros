# syntax=docker/dockerfile:1

## Build
FROM golang:1.19.2-alpine3.16 AS builder

WORKDIR /workspace

RUN mkdir /workspace/build
RUN mkdir /workspace/build/settings-api

#copy and build common module
RUN mkdir /workspace/build/customer-os-common-module

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

#copy and build file storage api
COPY settings-api/go.mod /workspace/build/settings-api/go.mod
COPY settings-api/go.sum /workspace/build/settings-api/go.sum

RUN cd /workspace/build/settings-api && go mod download && go mod verify

COPY settings-api/. /workspace/build/settings-api/

RUN cd /workspace/build/settings-api && go build -v -o /go/bin/app

FROM alpine:3.14
COPY --chown=65534:65534 --from=builder /go/bin/app .
COPY --chown=65534:65534 --from=builder /workspace/build/settings-api/.env .env
USER 65534

EXPOSE 10001
ENTRYPOINT [ "./app" ]