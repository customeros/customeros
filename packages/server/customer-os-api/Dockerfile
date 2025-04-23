FROM golang:1.21.0-alpine3.17 AS builder

WORKDIR /app
COPY customer-os-common-module/. ./customer-os-common-module/
COPY customer-os-neo4j-repository/. ./customer-os-neo4j-repository/
COPY customer-os-postgres-repository/. ./customer-os-postgres-repository/
COPY mailsherpa-api/. ./mailsherpa-api/
COPY core-crm/. ./core-crm/
COPY customer-os-api/. ./customer-os-api

WORKDIR /app/customer-os-api
RUN go build -v -o /go/bin/app

# Use a minimal base image for the runtime stage
FROM alpine:3.14

COPY --chown=65534:65534 --from=builder /go/bin/app .
COPY --chown=65534:65534 --from=builder /app/customer-os-api/.env .env
COPY --chown=65534:65534 --from=builder /app/customer-os-api/static/* /static/
USER 65534

ENTRYPOINT [ "./app" ]