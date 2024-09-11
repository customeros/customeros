FROM golang:1.21.0-alpine3.17 AS builder

WORKDIR /app
COPY customer-os-common-module/. ./customer-os-common-module/
COPY customer-os-common-ai/. ./customer-os-common-ai/
COPY customer-os-neo4j-repository/. ./customer-os-neo4j-repository/
COPY customer-os-postgres-repository/. ./customer-os-postgres-repository/
COPY enrichment-api/. ./enrichment-api/
COPY validation-api/. ./validation-api/
COPY events-processing-proto/. ./events-processing-proto/
COPY events/. ./events/
COPY events-processing-platform/. ./events-processing-platform
COPY events-processing-platform-subscribers/. ./events-processing-platform-subscribers

WORKDIR /app/events-processing-platform-subscribers
RUN go build -v -o /go/bin/app

FROM alpine:3.14

COPY --chown=65534:65534 --from=builder /go/bin/app .
COPY --chown=65534:65534 --from=builder /app/events-processing-platform-subscribers/.env .env
COPY --chown=65534:65534 --from=builder /app/events-processing-platform-subscribers/subscriptions/invoice/pdf_template/* /subscriptions/invoice/pdf_template/
COPY --chown=65534:65534 --from=builder /app/events-processing-platform-subscribers/static/* /static/
USER 65534

ENTRYPOINT [ "./app" ]