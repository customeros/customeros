# syntax=docker/dockerfile:1

## Build
FROM golang:1.19.2-alpine3.16 AS builder

WORKDIR /workspace

RUN mkdir /workspace/build
RUN mkdir /workspace/build/file-storage-api

#copy and build common module
RUN mkdir /workspace/build/customer-os-common-module

COPY customer-os-common-module/go.mod /workspace/build/customer-os-common-module/go.mod
COPY customer-os-common-module/go.sum /workspace/build/customer-os-common-module/go.sum

RUN cd /workspace/build/customer-os-common-module && go mod download && go mod verify

COPY customer-os-common-module/. /workspace/build/customer-os-common-module/
RUN cd /workspace/build/customer-os-common-module && go build

#copy and build file storage api
COPY file-storage-api/go.mod /workspace/build/file-storage-api/go.mod
COPY file-storage-api/go.sum /workspace/build/file-storage-api/go.sum

RUN cd /workspace/build/file-storage-api && go mod download && go mod verify

COPY file-storage-api/. /workspace/build/file-storage-api/

RUN cd /workspace/build/file-storage-api && go build -v -o /go/bin/app

FROM alpine:3.14
COPY --chown=65534:65534 --from=builder /go/bin/app .
COPY --chown=65534:65534 --from=builder /workspace/build/file-storage-api/.env .env
USER 65534

EXPOSE 10001
ENTRYPOINT [ "./app" ]