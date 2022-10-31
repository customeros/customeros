#!/bin/sh
NAMESPACE_NAME=openline
kubectl port-forward --namespace $NAMESPACE_NAME svc/postgresql-customer-os-dev 5432:5432 &
kubectl port-forward --namespace $NAMESPACE_NAME svc/neo4j-customer-os 7687:7687 &
