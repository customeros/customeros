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
	docker pull ghcr.io/openline-ai/openline-customer-os/customer-os-api:otter 
	minikube image load ghcr.io/openline-ai/openline-customer-os/customer-os-api:otter

	docker pull ghcr.io/openline-ai/openline-customer-os/message-store:otter 
	minikube image load ghcr.io/openline-ai/openline-customer-os/message-store:otter
fi

# Deploy Images
cd $CUSTOMER_OS_HOME/deployment/k8s/local-minikube
kubectl apply -f apps-config/customer-os-api.yaml --namespace $NAMESPACE_NAME
kubectl apply -f apps-config/customer-os-api-k8s-service.yaml --namespace $NAMESPACE_NAME
kubectl apply -f apps-config/message-store.yaml --namespace $NAMESPACE_NAME
kubectl apply -f apps-config/message-store-k8s-service.yaml --namespace $NAMESPACE_NAME

#provision neo4j

while [ -z "$pod" ]; do
    pod=$(kubectl get pods -n $NAMESPACE_NAME|grep neo4j-customer-os|grep Running| cut -f1 -d ' ')
    if [ -z "$pod" ]; then
      echo "neo4j not ready waiting"
      sleep 1
    fi
done

started=""
while [ -z "$started" ]; do
    started=$(kubectl logs -n $NAMESPACE_NAME $pod|grep password)
    if [ -z "$started" ]; then
      echo "neo4j waiting for app to start"
      sleep 1
    fi
done
sleep 1

neo_output="not empty"
while  [ ! -z "$neo_output" ]; do
	echo "provisioning neo4j"
	neo_output=$(cat $CUSTOMER_OS_HOME/packages/server/customer-os-api/customer-os.cypher |kubectl run --rm -i --namespace $NAMESPACE_NAME --image "neo4j:4.4.11" cypher-shell  -- bash -c 'NEO4J_PASSWORD=StrongLocalPa\$\$ cypher-shell -a neo4j://neo4j-customer-os.openline.svc.cluster.local:7687 -u neo4j --non-interactive' 2>&1 |grep -v "see a command prompt" |grep -v "deleted")
	if [ ! -z "$neo_output" ]; then
		echo "neo4j provisioning failed, trying again"
		echo "output: $neo_output"
		sleep 1
	fi
done

echo "provisioning postrgess"
cd $CUSTOMER_OS_HOME/packages/server/message-store/sql
SQL_USER=openline SQL_DATABABASE=openline SQL_PASSWORD=password ./build_db.sh local-kube
