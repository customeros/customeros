# syntax=docker/dockerfile:1

FROM golang:1.19-alpine as builder
WORKDIR /app

COPY comms-api ./comms-api

WORKDIR /app/comms-api
RUN make clean all

FROM alpine:3.14
COPY --chown=65534:65534 --from=builder /app/comms-api/bin/comms-api .
USER 65534

ENTRYPOINT [ "./comms-api" ]
