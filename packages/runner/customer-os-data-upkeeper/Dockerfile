FROM golang:1.20.1-alpine3.17@sha256:87d0a3309b34e2ca732efd69fb899d3c420d3382370fd6e7e6d2cb5c930f27f9 AS builder

WORKDIR /app
RUN mkdir /app/packages
RUN mkdir /app/packages/runner
RUN mkdir /app/packages/server

COPY server/customer-os-common-module/. /app/packages/server/customer-os-common-module/
COPY server/customer-os-neo4j-repository/. /app/packages/server/customer-os-neo4j-repository/
COPY server/customer-os-postgres-repository/. /app/packages/server/customer-os-postgres-repository/
COPY server/events-processing-proto/. /app/packages/server/events-processing-proto/
COPY server/events/. /app/packages/server/events/
COPY server/customer-os-api-sdk/. /app/packages/server/customer-os-api-sdk/
COPY runner/customer-os-data-upkeeper/. /app/packages/runner/customer-os-data-upkeeper

WORKDIR /app/packages/runner/customer-os-data-upkeeper
RUN go build -v -o /go/bin/app

FROM alpine:3.14@sha256:0f2d5c38dd7a4f4f733e688e3a6733cb5ab1ac6e3cb4603a5dd564e5bfb80eed

COPY --chown=65534:65534 --from=builder /go/bin/app .
COPY --chown=65534:65534 --from=builder /app/packages/runner/customer-os-data-upkeeper/.env .env
USER 65534

ENTRYPOINT [ "./app" ]