FROM golang:1.20.1-alpine3.17@sha256:87d0a3309b34e2ca732efd69fb899d3c420d3382370fd6e7e6d2cb5c930f27f9 AS builder

WORKDIR /app
COPY customer-os-common-module/. ./customer-os-common-module/
COPY customer-os-common-ai/. ./customer-os-common-ai/
COPY customer-os-neo4j-repository/. ./customer-os-neo4j-repository/
COPY customer-os-postgres-repository/. ./customer-os-postgres-repository/
COPY events-processing-proto/. ./events-processing-proto/
COPY events/. ./events/
COPY events-processing-platform/. ./events-processing-platform
COPY events-processing-platform-subscribers/. ./events-processing-platform-subscribers

WORKDIR /app/events-processing-platform-subscribers
RUN go build -v -o /go/bin/app

FROM alpine:3.14@sha256:0f2d5c38dd7a4f4f733e688e3a6733cb5ab1ac6e3cb4603a5dd564e5bfb80eed

COPY --chown=65534:65534 --from=builder /go/bin/app .
COPY --chown=65534:65534 --from=builder /app/events-processing-platform-subscribers/.env .env
COPY --chown=65534:65534 --from=builder /app/events-processing-platform-subscribers/subscriptions/invoice/pdf_template/* /subscriptions/invoice/pdf_template/
COPY --chown=65534:65534 --from=builder /app/events-processing-platform-subscribers/static/* /static/
USER 65534

ENTRYPOINT [ "./app" ]