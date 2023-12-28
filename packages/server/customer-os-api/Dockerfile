FROM golang:1.20.1-alpine3.17@sha256:87d0a3309b34e2ca732efd69fb899d3c420d3382370fd6e7e6d2cb5c930f27f9 AS builder

WORKDIR /workspace

RUN mkdir /workspace/build
RUN mkdir /workspace/build/customer-os-api
RUN mkdir /workspace/build/customer-os-common-module
RUN mkdir /workspace/build/customer-os-neo4j-repository
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

#copy events-processing-proto
RUN mkdir /workspace/build/events-processing-proto
COPY events-processing-proto/go.mod /workspace/build/events-processing-proto/go.mod
COPY events-processing-proto/go.sum /workspace/build/events-processing-proto/go.sum
COPY events-processing-proto/. /workspace/build/events-processing-proto/

#copy and build customer-os-neo4j-repository
COPY customer-os-neo4j-repository/go.mod /workspace/build/customer-os-neo4j-repository/go.mod
COPY customer-os-neo4j-repository/go.sum /workspace/build/customer-os-neo4j-repository/go.sum
RUN cd /workspace/build/customer-os-neo4j-repository && go mod download && go mod verify
COPY customer-os-neo4j-repository/. /workspace/build/customer-os-neo4j-repository/
RUN cd /workspace/build/customer-os-neo4j-repository && go build

#copy and build customer-os-api
COPY customer-os-api/go.mod /workspace/build/customer-os-api/go.mod
COPY customer-os-api/go.sum /workspace/build/customer-os-api/go.sum

RUN cd /workspace/build/customer-os-api && go mod download && go mod verify

COPY customer-os-api/. /workspace/build/customer-os-api/

RUN cd /workspace/build/customer-os-api && go build -v -o /go/bin/app

FROM alpine:3.14@sha256:0f2d5c38dd7a4f4f733e688e3a6733cb5ab1ac6e3cb4603a5dd564e5bfb80eed

COPY --chown=65534:65534 --from=builder /go/bin/app .
COPY --chown=65534:65534 --from=builder /workspace/build/customer-os-api/.env .env
USER 65534

ENTRYPOINT [ "./app" ]