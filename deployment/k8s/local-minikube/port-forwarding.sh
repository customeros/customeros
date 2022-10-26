#!/bin/sh
NAMESPACE_NAME=openline-development
kubectl port-forward --namespace $NAMESPACE_NAME svc/postgresql-customer-os-dev 5432:5432 &
