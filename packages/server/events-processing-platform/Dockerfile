# syntax=docker/dockerfile:1

## Build
FROM golang:1.20.1-alpine3.17 AS builder

WORKDIR /workspace

RUN mkdir /workspace/build
RUN mkdir /workspace/build/customer-os-api
RUN mkdir /workspace/build/customer-os-common-module
RUN mkdir /workspace/build/events-processing-common

#copy and build common module
COPY customer-os-common-module/go.mod /workspace/build/customer-os-common-module/go.mod
COPY customer-os-common-module/go.sum /workspace/build/customer-os-common-module/go.sum
RUN cd /workspace/build/customer-os-common-module && go mod download && go mod verify

COPY customer-os-common-module/. /workspace/build/customer-os-common-module/
RUN cd /workspace/build/customer-os-common-module && go build

#copy and events-processing-common
COPY events-processing-common/go.mod /workspace/build/events-processing-common/go.mod
COPY events-processing-common/go.sum /workspace/build/events-processing-common/go.sum
COPY events-processing-common/. /workspace/build/events-processing-common/

#copy and build events-processing-platform
COPY events-processing-platform/go.mod /workspace/build/events-processing-platform/go.mod
COPY events-processing-platform/go.sum /workspace/build/events-processing-platform/go.sum
RUN cd /workspace/build/events-processing-platform && go mod download && go mod verify

COPY events-processing-platform/. /workspace/build/events-processing-platform/
RUN cd /workspace/build/events-processing-platform && go build -v -o /go/bin/app

FROM alpine:3.14
COPY --chown=65534:65534 --from=builder /go/bin/app .
COPY --chown=65534:65534 --from=builder /workspace/build/events-processing-platform/.env .env
USER 65534

ENTRYPOINT [ "./app" ]
