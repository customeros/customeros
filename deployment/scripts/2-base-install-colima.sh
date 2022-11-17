#!/bin/bash

# Start minikube
echo "  🦦 starting Colima"
colima start --with-kubernetes
if [ $? -eq 0 ]; then
    echo "  ✅ Colima running"
else
    echo "  ❌ Colima failed to start"
fi


#Get namespace config & setup namespace in kubernetes
OPENLINE_NAMESPACE="./openline-setup/openline-namespace.json"
NAMESPACE_NAME="openline"
if [[ $(kubectl get namespaces) == *"$NAMESPACE_NAME"* ]];
  then
    echo "  🦦 Continue deploy on namespace $NAMESPACE_NAME"
  else
    echo "  🦦 Creating $NAMESPACE_NAME namespace in Kubernetes"
    kubectl create -f $OPENLINE_NAMESPACE
    if [ $? -eq 0 ]; then
        echo "  ✅ $NAMESPACE_NAME namespace created in Kubernetes"
    else
        echo "  ❌ failed to create $NAMESPACE_NAME namespace in Kubernetes"
    fi
fi

#Adding helm repos :
echo "  🦦 adding helm repos"
helm repo add bitnami https://charts.bitnami.com/bitnami
if [ $? -eq 0 ]; then
    echo "  ✅ bitnami"
else
    echo "  ❌ bitnami"
    exit 1
fi
helm repo add neo4j https://helm.neo4j.com/neo4j
if [ $? -eq 0 ]; then
    echo "  ✅ neo4j"
else
    echo "  ❌ neo4j"
    exit 1
fi
helm repo add fusionauth https://fusionauth.github.io/charts
if [ $? -eq 0 ]; then
    echo "  ✅ fusionauth"
else
    echo "  ❌ fusionauth"
    exit 1
fi
helm repo update
echo "  ✅ helm repos updated"

#Get postgresql config and install 
POSTGRESQL_PERSISTENT_VOLUME="./openline-setup/postgresql-persistent-volume.yaml"
kubectl apply -f $POSTGRESQL_PERSISTENT_VOLUME --namespace $NAMESPACE_NAME
if [ $? -eq 0 ]; then
    echo "  ✅ postgresql-persistent-volume.yaml configured"
else
    echo "  ❌ postgresql-persistent-volume.yaml not configured"
fi

POSTGRESQL_PERSISTENT_VOLUME_CLAIM="./openline-setup/postgresql-persistent-volume-claim.yaml"
kubectl apply -f $POSTGRESQL_PERSISTENT_VOLUME_CLAIM --namespace $NAMESPACE_NAME
if [ $? -eq 0 ]; then
    echo "  ✅ postgresql-persistent-volume-claim.yaml configured"
else
    echo "  ❌ postgresql-persistent-volume-claim.yaml not configured"
fi

POSTGRESQL_VALUES="./openline-setup/postgresql-values.yaml"
helm install --values $POSTGRESQL_VALUES postgresql-customer-os-dev bitnami/postgresql --namespace $NAMESPACE_NAME
if [ $? -eq 0 ]; then
    echo "  ✅ postgresql installed"
else
    echo "  ❌ postgresql not installed"
fi

#Get Neo4j config and install
NEO4J_HELM_VALUES="./openline-setup/neo4j-helm-values.yaml"
helm install neo4j-customer-os neo4j/neo4j-standalone --set volumes.data.mode=defaultStorageClass -f $NEO4J_HELM_VALUES --namespace $NAMESPACE_NAME
if [ $? -eq 0 ]; then
    echo "  ✅ Neo4j installed"
else
    echo "  ❌ Neo4j not installed"
fi

#Get FusionAuth config and install
FUSIONAUTH_VALUES="./openline-setup/fusionauth-values.yaml"
helm install fusionauth-customer-os fusionauth/fusionauth -f $FUSIONAUTH_VALUES --namespace $NAMESPACE_NAME
if [ $? -eq 0 ]; then
    echo "  ✅ FusionAuth installed"
else
    echo "  ❌ FusionAuth not installed"
fi