# syntax=docker/dockerfile:1

FROM golang:1.20.1-alpine3.17 as builder
RUN apk add --update make bash

WORKDIR /app

COPY user-admin-api ./user-admin-api
COPY customer-os-common-auth ./customer-os-common-auth

WORKDIR /app/user-admin-api
RUN go build -o bin/user-admin-api

FROM alpine:3.14
COPY --chown=65534:65534 --from=builder /app/user-admin-api/bin/user-admin-api .
USER 65534

ENTRYPOINT [ "./user-admin-api" ]
