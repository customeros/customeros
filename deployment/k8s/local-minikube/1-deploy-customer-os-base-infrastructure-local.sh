#! /bin/sh

MINIKUBE_STATUS=$(minikube status)
MINIKUBE_STARTED_STATUS_TEXT='Running'
if [[ "$MINIKUBE_STATUS" == *"$MINIKUBE_STARTED_STATUS_TEXT"* ]];
  then
     echo " --- Minikube already started --- "
  else
     eval $(minikube docker-env)
     minikube start &
     wait
fi


NAMESPACE_NAME="openline-development"

if [[ $(kubectl get namespaces) == *"$NAMESPACE_NAME"* ]];
  then
    echo " --- Continue deploy on namespace openline-development --- "
  else
    echo " --- Creating Openline Development namespace in minikube ---"
    kubectl create -f "../configs/openline-namespace.json"
    wait
fi

#Adding helm repos :
helm repo add bitnami https://charts.bitnami.com/bitnami
helm repo add neo4j https://helm.neo4j.com/neo4j
helm repo update

#install postgresql
kubectl apply -f ../configs/postgresql/postgresql-presistent-volume.yaml --namespace $NAMESPACE_NAME
kubectl apply -f ../configs/postgresql/postgresql-persistent-volume-claim.yaml --namespace $NAMESPACE_NAME

helm install --values "../configs/postgresql/postgresql-values.yaml" postgresql-customer-os-dev bitnami/postgresql --namespace $NAMESPACE_NAME
wait

helm install neo4j-customer-os-dev neo4j/neo4j-standalone --set volumes.data.mode=defaultStorageClass -f ../configs/neo4j/neo4j-helm-values.yaml --namespace $NAMESPACE_NAME
