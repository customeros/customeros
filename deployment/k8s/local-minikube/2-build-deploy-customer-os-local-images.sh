#! /bin/sh

## Build Images
cd  $CUSTOMER_OS_HOME/customer-os-api/server/
docker build -t customer-os-api .

minikube image load customer-os-api:latest

# Deploy Images
NAMESPACE_NAME="openline-development"

cd $CUSTOMER_OS_HOME/deployment/k8s/local-minikube
kubectl apply -f apps-config/customer-os-api.yaml --namespace $NAMESPACE_NAME
kubectl apply -f apps-config/customer-os-api-k8s-service.yaml --namespace $NAMESPACE_NAME
