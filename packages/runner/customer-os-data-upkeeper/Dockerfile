FROM golang:1.23.1-alpine3.20 AS builder

WORKDIR /app
RUN mkdir /app/packages
RUN mkdir /app/packages/runner
RUN mkdir /app/packages/server

COPY server/customer-os-common-module/. /app/packages/server/customer-os-common-module/
COPY server/customer-os-neo4j-repository/. /app/packages/server/customer-os-neo4j-repository/
COPY server/customer-os-postgres-repository/. /app/packages/server/customer-os-postgres-repository/
COPY server/events-processing-proto/. /app/packages/server/events-processing-proto/
COPY server/events/. /app/packages/server/events/
COPY server/enrichment-api/. /app/packages/server/enrichment-api/
COPY server/validation-api/. /app/packages/server/validation-api/
COPY runner/customer-os-data-upkeeper/. /app/packages/runner/customer-os-data-upkeeper

WORKDIR /app/packages/runner/customer-os-data-upkeeper
RUN go build -v -o /go/bin/app

FROM alpine:3.20@sha256:0a4eaa0eecf5f8c050e5bba433f58c052be7587ee8af3e8b3910ef9ab5fbe9f5

COPY --chown=65534:65534 --from=builder /go/bin/app .
COPY --chown=65534:65534 --from=builder /app/packages/runner/customer-os-data-upkeeper/.env .env
USER 65534

ENTRYPOINT [ "./app" ]