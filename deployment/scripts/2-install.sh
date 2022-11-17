#!/bin/bash

####### Script Variables #####################
OPENLINE_NAMESPACE="./openline-setup/openline-namespace.json"
NAMESPACE_NAME="openline"

POSTGRESQL_PERSISTENT_VOLUME="./openline-setup/postgresql-persistent-volume.yaml"
POSTGRESQL_PERSISTENT_VOLUME_CLAIM="./openline-setup/postgresql-persistent-volume-claim.yaml"

POSTGRESQL_HELM_VALUES="./openline-setup/postgresql-values.yaml"
NEO4J_HELM_VALUES="./openline-setup/neo4j-helm-values.yaml"
FUSIONAUTH_HELM_VALUES="./openline-setup/fusionauth-values.yaml"

CUSTOMER_OS_API_IMAGE="ghcr.io/openline-ai/openline-customer-os/customer-os-api:latest"
CUSTOMER_OS_API="openline-setup/customer-os-api.yaml"
CUSTOMER_OS_API_K8S_SERVICE="openline-setup/customer-os-api-k8s-service.yaml"
CUSTOMER_OS_API_LOADBALANCER="openline-setup/customer-os-api-k8s-loadbalancer-service.yaml"

MESSAGE_STORE_API_IMAGE="ghcr.io/openline-ai/openline-customer-os/message-store:latest"
MESSAGE_STORE_API="openline-setup/message-store.yaml"
MESSAGE_STORE_API_K8S_SERVICE="openline-setup/message-store-k8s-service.yaml"

NEO4J_CYPHER="openline-setup/customer-os.cypher"
##############################################


# Start Colima
echo "  🦦 starting Colima"
colima start --with-kubernetes --cpu 2 --memory 4 --disk 60
if [ $? -eq 0 ]; then
    echo "✅ Colima running"
else
    echo "❌ Colima failed to start"
fi


#Get namespace config & setup namespace in kubernetes
if [[ $(kubectl get namespaces) == *"$NAMESPACE_NAME"* ]];
  then
    echo "🦦 Continue deploy on namespace $NAMESPACE_NAME"
  else
    echo "🦦 Creating $NAMESPACE_NAME namespace in Kubernetes"
    kubectl create -f $OPENLINE_NAMESPACE
    if [ $? -eq 0 ]; then
        echo "✅ $NAMESPACE_NAME namespace created in Kubernetes"
    else
        echo "❌ failed to create $NAMESPACE_NAME namespace in Kubernetes"
    fi
fi

#Adding helm repos :
echo "  🦦 adding helm repos"
helm repo add bitnami https://charts.bitnami.com/bitnami
if [ $? -eq 0 ]; then
    echo "✅ bitnami"
else
    echo "❌ bitnami"
    exit 1
fi
helm repo add neo4j https://helm.neo4j.com/neo4j
if [ $? -eq 0 ]; then
    echo "✅ neo4j"
else
    echo "❌ neo4j"
    exit 1
fi
helm repo add fusionauth https://fusionauth.github.io/charts
if [ $? -eq 0 ]; then
    echo "✅ fusionauth"
else
    echo "❌ fusionauth"
    exit 1
fi
helm repo update
echo "✅ helm repos updated"

# Install Neo4j
helm install neo4j-customer-os neo4j/neo4j-standalone --set volumes.data.mode=defaultStorageClass -f $NEO4J_HELM_VALUES --namespace $NAMESPACE_NAME
if [ $? -eq 0 ]; then
    echo "✅ Neo4j installed"
else
    echo "❌ Neo4j not installed"
fi

# Get PostgreSQL config and setup disc
kubectl apply -f $POSTGRESQL_PERSISTENT_VOLUME --namespace $NAMESPACE_NAME
if [ $? -eq 0 ]; then
    echo "✅ postgresql-persistent-volume.yaml configured"
else
    echo "❌ postgresql-persistent-volume.yaml not configured"
fi

kubectl apply -f $POSTGRESQL_PERSISTENT_VOLUME_CLAIM --namespace $NAMESPACE_NAME
if [ $? -eq 0 ]; then
    echo "✅ postgresql-persistent-volume-claim.yaml configured"
else
    echo "❌ postgresql-persistent-volume-claim.yaml not configured"
fi

# Install PostgreSQL
helm install --values $POSTGRESQL_HELM_VALUES postgresql-customer-os-dev bitnami/postgresql --namespace $NAMESPACE_NAME
if [ $? -eq 0 ]; then
    echo "✅ postgresql installed"
else
    echo "❌ postgresql not installed"
fi

# Download latest container images
echo "  🦦 Getting latest customerOS API docker image..."
docker pull $CUSTOMER_OS_API_IMAGE
  if [ $? -eq 0 ]; then
    echo "✅ grabbed customerOS API image"
  else
    echo "❌ unable to grab customerOS API image"
  fi

echo "  🦦 Getting latest message store API docker image..."
docker pull $MESSAGE_STORE_API_IMAGE
  if [ $? -eq 0 ]; then
    echo "✅ grabbed message store API image"
  else
    echo "❌ unable to grab message store API image"
  fi

# Deploy Images
echo "  🦦 Deploying customerOS API..."
kubectl apply -f $CUSTOMER_OS_API --namespace $NAMESPACE_NAME
  if [ $? -eq 0 ]; then
    echo "✅ customer-os-api.yaml"
  else
    echo "❌ customer-os-api.yaml"
  fi

kubectl apply -f $CUSTOMER_OS_API_K8S_SERVICE --namespace $NAMESPACE_NAME
  if [ $? -eq 0 ]; then
    echo "✅ customer-os-api-k8s-service.yaml"
  else
    echo "❌ customer-os-api-k8s-service.yaml"
  fi

kubectl apply -f $CUSTOMER_OS_API_LOADBALANCER --namespace $NAMESPACE_NAME
  if [ $? -eq 0 ]; then
    echo "✅ customer-os-api-k8s-loadbalancer-service.yaml"
  else
    echo "❌ customer-os-api-k8s-loadbalancer-service.yaml"
  fi

echo "🦦 Deploying message store API..."

kubectl apply -f $MESSAGE_STORE_API --namespace $NAMESPACE_NAME
  if [ $? -eq 0 ]; then
    echo "✅ message-store.yaml"
  else
    echo "❌ message-store.yaml"
  fi

kubectl apply -f $MESSAGE_STORE_API_K8S_SERVICE --namespace $NAMESPACE_NAME
  if [ $? -eq 0 ]; then
    echo "✅ message-store-k8s-service.yaml"
  else
    echo "❌ message-store-k8s-service.yaml"
  fi

# Install FusionAuth
helm install fusionauth-customer-os fusionauth/fusionauth -f $FUSIONAUTH_HELM_VALUES --namespace $NAMESPACE_NAME
if [ $? -eq 0 ]; then
    echo "  ✅ FusionAuth installed"
else
    echo "  ❌ FusionAuth not installed"
fi
