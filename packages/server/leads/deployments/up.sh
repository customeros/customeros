#!/bin/sh

nerdctl compose -f ./../../../../deployment/docker-compose.yaml up -d postgres timescaledb timescaledb-sidecar nats otel-collector tempo grafana 

