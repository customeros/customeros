FROM golang:1.20.1-alpine3.17@sha256:87d0a3309b34e2ca732efd69fb899d3c420d3382370fd6e7e6d2cb5c930f27f9 AS builder

WORKDIR /workspace

RUN mkdir /workspace/packages
RUN mkdir /workspace/packages/runner
RUN mkdir /workspace/packages/runner/customer-os-data-upkeeper
RUN mkdir /workspace/packages/server
RUN mkdir /workspace/packages/server/customer-os-common-module
RUN mkdir /workspace/packages/server/events-processing-proto
RUN mkdir /workspace/packages/server/customer-os-neo4j-repository

#copy and build common module
COPY customer-os-common-module/go.mod /workspace/packages/server/customer-os-common-module/go.mod
COPY customer-os-common-module/go.sum /workspace/packages/server/customer-os-common-module/go.sum
RUN cd /workspace/packages/server/customer-os-common-module && go mod download && go mod verify
COPY customer-os-common-module/. /workspace/packages/server/customer-os-common-module/
RUN cd /workspace/packages/server/customer-os-common-module && go build

#copy and build events-processing-proto
COPY events-processing-proto/go.mod /workspace/packages/server/events-processing-proto/go.mod
COPY events-processing-proto/go.sum /workspace/packages/server/events-processing-proto/go.sum
COPY events-processing-proto/. /workspace/packages/server/events-processing-proto/

#copy and build customer-os-neo4j-repository
COPY customer-os-neo4j-repository/go.mod /workspace/packages/server/customer-os-neo4j-repository/go.mod
COPY customer-os-neo4j-repository/go.sum /workspace/packages/server/customer-os-neo4j-repository/go.sum
RUN cd /workspace/packages/server/customer-os-neo4j-repository && go mod download && go mod verify
COPY customer-os-neo4j-repository/. /workspace/packages/server/customer-os-neo4j-repository/
RUN cd /workspace/packages/server/customer-os-neo4j-repository && go build

#copy and build customer-os-data-upkeeper
COPY customer-os-data-upkeeper/go.mod /workspace/packages/runner/customer-os-data-upkeeper/go.mod
COPY customer-os-data-upkeeper/go.sum /workspace/packages/runner/customer-os-data-upkeeper/go.sum
RUN cd /workspace/packages/runner/customer-os-data-upkeeper && go mod download && go mod verify

COPY customer-os-data-upkeeper/. /workspace/packages/runner/customer-os-data-upkeeper/
RUN cd /workspace/packages/runner/customer-os-data-upkeeper && go build -v -o /go/bin/app

FROM alpine:3.14@sha256:0f2d5c38dd7a4f4f733e688e3a6733cb5ab1ac6e3cb4603a5dd564e5bfb80eed

COPY --chown=65534:65534 --from=builder /go/bin/app .
COPY --chown=65534:65534 --from=builder /workspace/packages/runner/customer-os-data-upkeeper/.env .env
USER 65534

ENTRYPOINT [ "./app" ]