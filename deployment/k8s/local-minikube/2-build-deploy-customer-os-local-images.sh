#! /bin/sh
NAMESPACE_NAME="openline"
CUSTOMER_OS_HOME="$(dirname $(readlink -f $0))/../../../"

## Build Images
cd  $CUSTOMER_OS_HOME/customer-os-api/server/

## clean up previous install
kubectl delete service customer-os-api-service --namespace $NAMESPACE_NAME
kubectl delete deployments customer-os-api --namespace $NAMESPACE_NAME
docker image rm customer-os-api

minikube image unload customer-os-api

# build image

docker build --no-cache -t customer-os-api .
minikube image load customer-os-api:latest

# Deploy Images
cd $CUSTOMER_OS_HOME/deployment/k8s/local-minikube
kubectl apply -f apps-config/customer-os-api.yaml --namespace $NAMESPACE_NAME
kubectl apply -f apps-config/customer-os-api-k8s-service.yaml --namespace $NAMESPACE_NAME
