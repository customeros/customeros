#! /bin/bash

NAMESPACE_NAME="openline"
CUSTOMER_OS_HOME="$(dirname $(readlink -f $0))/../../../"
echo CUSTOMER_OS_HOME=$CUSTOMER_OS_HOME

if [[ $(kubectl get namespaces) == *"$NAMESPACE_NAME"* ]];
  then
    echo " --- Continue deploy on namespace openline --- "
  else
    echo " --- Creating Openline Development namespace in minikube ---"
    kubectl create -f "$CUSTOMER_OS_HOME/deployment/k8s/configs/openline-namespace.json"
    wait
fi

#Adding helm repos :
helm repo add bitnami https://charts.bitnami.com/bitnami
helm repo add neo4j https://helm.neo4j.com/neo4j
helm repo update

#install postgresql
kubectl apply -f $CUSTOMER_OS_HOME/deployment/k8s/configs/postgresql/postgresql-presistent-volume.yaml --namespace $NAMESPACE_NAME
kubectl apply -f $CUSTOMER_OS_HOME/deployment/k8s/configs/postgresql/postgresql-persistent-volume-claim.yaml --namespace $NAMESPACE_NAME

helm install --values "$CUSTOMER_OS_HOME/deployment/k8s/configs/postgresql/postgresql-values.yaml" postgresql-customer-os-dev bitnami/postgresql --namespace $NAMESPACE_NAME
wait

helm install neo4j-customer-os neo4j/neo4j-standalone --set volumes.data.mode=defaultStorageClass -f $CUSTOMER_OS_HOME/deployment/k8s/configs/neo4j/neo4j-helm-values.yaml --namespace $NAMESPACE_NAME

echo "Completed."