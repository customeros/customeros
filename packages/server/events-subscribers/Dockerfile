FROM golang:1.21.0-alpine3.17 AS builder

WORKDIR /app
COPY customer-os-common-module/. ./customer-os-common-module/
COPY customer-os-neo4j-repository/. ./customer-os-neo4j-repository/
COPY customer-os-postgres-repository/. ./customer-os-postgres-repository/
COPY events-processing-proto/. ./events-processing-proto/
COPY events-subscribers/. ./events-subscribers
COPY enrichment-api/. ./enrichment-api/
COPY validation-api/. ./validation-api/

WORKDIR /app/events-subscribers
RUN go build -v -o /go/bin/app

FROM alpine:3.14

COPY --chown=65534:65534 --from=builder /go/bin/app .
COPY --chown=65534:65534 --from=builder /app/events-subscribers/.env .env
USER 65534

ENTRYPOINT [ "./app" ]