FROM golang:1.20.1-alpine3.17@sha256:87d0a3309b34e2ca732efd69fb899d3c420d3382370fd6e7e6d2cb5c930f27f9 AS builder

WORKDIR /app
COPY customer-os-common-module/. ./customer-os-common-module/
COPY customer-os-api/. ./customer-os-api/
COPY comms-api/. ./comms-api

WORKDIR /app/comms-api
RUN go build -o bin/comms-api

FROM alpine:3.14@sha256:0f2d5c38dd7a4f4f733e688e3a6733cb5ab1ac6e3cb4603a5dd564e5bfb80eed

COPY --chown=65534:65534 --from=builder /app/comms-api/bin/comms-api .
USER 65534

ENTRYPOINT [ "./comms-api" ]
