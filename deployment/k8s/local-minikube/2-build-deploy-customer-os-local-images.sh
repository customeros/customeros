#! /bin/sh
NAMESPACE_NAME="openline"
CUSTOMER_OS_HOME="$(dirname $(readlink -f $0))/../../../"

## Build Images

## clean up previous install
kubectl delete service customer-os-api-service --namespace $NAMESPACE_NAME
kubectl delete deployments customer-os-api --namespace $NAMESPACE_NAME
kubectl delete service message-store-service --namespace $NAMESPACE_NAME
kubectl delete deployments message-store --namespace $NAMESPACE_NAME
docker image rm ghcr.io/openline-ai/openline-customer-os/customer-os-api:otter
docker image rm ghcr.io/openline-ai/openline-customer-os/message-store:otter

minikube image unload ghcr.io/openline-ai/openline-customer-os/customer-os-api:otter
minikube image unload ghcr.io/openline-ai/openline-customer-os/message-store:otter

# build or download image

if [ "x$1" == "xbuild" ]; then
	echo "locally building docker images"
	cd  $CUSTOMER_OS_HOME/packages/server/customer-os-api/
	docker build --no-cache -t ghcr.io/openline-ai/openline-customer-os/customer-os-api:otter .
	minikube image load ghcr.io/openline-ai/openline-customer-os/customer-os-api:otter

	cd  $CUSTOMER_OS_HOME
docker build --no-cache -t ghcr.io/openline-ai/openline-customer-os/message-store:otter -f packages/server/message-store/Dockerfile packages/server/
	minikube image load ghcr.io/openline-ai/openline-customer-os/message-store:otter
else
	echo "installing pre-build images"
	docker pull ghcr.io/openline-ai/openline-customer-os/customer-os-api:otter .
	minikube image load ghcr.io/openline-ai/openline-customer-os/customer-os-api:otter

	docker pull ghcr.io/openline-ai/openline-customer-os/message-store:otter .
	minikube image load ghcr.io/openline-ai/openline-customer-os/message-store:otter
fi

# Deploy Images
cd $CUSTOMER_OS_HOME/deployment/k8s/local-minikube
kubectl apply -f apps-config/customer-os-api.yaml --namespace $NAMESPACE_NAME
kubectl apply -f apps-config/customer-os-api-k8s-service.yaml --namespace $NAMESPACE_NAME
kubectl apply -f apps-config/message-store.yaml --namespace $NAMESPACE_NAME
kubectl apply -f apps-config/message-store-k8s-service.yaml --namespace $NAMESPACE_NAME

cd $CUSTOMER_OS_HOME/packages/server/message-store/sql
SQL_USER=openline SQL_DATABABASE=openline SQL_PASSWORD=password ./build_db.sh local-kube
