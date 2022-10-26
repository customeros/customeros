#!/bin/sh
NAMESPACE_NAME=openline-development
kubectl port-forward consul-server-0 --namespace $NAMESPACE_NAME 8500:8500 &
kubectl port-forward --namespace $NAMESPACE_NAME svc/postgresql-customer-os-dev 5432:5432 &
