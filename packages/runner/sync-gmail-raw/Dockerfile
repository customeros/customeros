# syntax=docker/dockerfile:1

## Build
FROM golang:1.20.1-alpine3.17 AS builder

WORKDIR /workspace

RUN mkdir /workspace/build
RUN mkdir /workspace/build/sync-gmail-raw
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

#copy and build sync-gmail-raw
COPY sync-gmail-raw/go.mod /workspace/build/sync-gmail-raw/go.mod
COPY sync-gmail-raw/go.sum /workspace/build/sync-gmail-raw/go.sum

RUN cd /workspace/build/sync-gmail-raw && go mod download && go mod verify

COPY sync-gmail-raw/. /workspace/build/sync-gmail-raw/

RUN cd /workspace/build/sync-gmail-raw && go build -v -o /go/bin/app

FROM alpine:3.14
COPY --chown=65534:65534 --from=builder /go/bin/app .
COPY --chown=65534:65534 --from=builder /workspace/build/sync-gmail-raw/.env .env
USER 65534

ENTRYPOINT [ "./app" ]