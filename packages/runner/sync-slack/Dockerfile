FROM golang:1.21.0-alpine3.17 AS builder

WORKDIR /app
RUN mkdir /app/packages
RUN mkdir /app/packages/runner
RUN mkdir /app/packages/server

COPY server/customer-os-common-module/. /app/packages/server/customer-os-common-module/
COPY server/customer-os-neo4j-repository/. /app/packages/server/customer-os-neo4j-repository/
COPY server/customer-os-postgres-repository/. /app/packages/server/customer-os-postgres-repository/
COPY server/events-processing-proto/. /app/packages/server/events-processing-proto/
COPY runner/sync-slack/. /app/packages/runner/sync-slack

WORKDIR /app/packages/runner/sync-slack
RUN go build -v -o /go/bin/app

FROM alpine:3.14

COPY --chown=65534:65534 --from=builder /go/bin/app .
COPY --chown=65534:65534 --from=builder /app/packages/runner/sync-slack/.env .env
USER 65534

ENTRYPOINT [ "./app" ]