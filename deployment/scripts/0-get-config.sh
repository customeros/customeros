#! /bin/sh

### Locations for remote file downloads ###

# Minikube config
OPENLINE_NAMESPACE="https://raw.githubusercontent.com/openline-ai/openline-customer-os/otter/deployment/infra/minikube-config/openline-namespace.json"
CUSTOMER_OS_API_CONFIG="https://raw.githubusercontent.com/openline-ai/openline-customer-os/otter/deployment/infra/minikube-config/customer-os-api.yaml"
CUSTOMER_OS_API_K8S_CONFIG="https://raw.githubusercontent.com/openline-ai/openline-customer-os/otter/deployment/infra/minikube-config/customer-os-api-k8s-service.yaml"
CUSTOMER_OS_API_K8S_LOADBALANCER_CONFIG="https://raw.githubusercontent.com/openline-ai/openline-customer-os/otter/deployment/infra/minikube-config/customer-os-api-k8s-loadbalancer-service.yaml"
MESSAGE_STORE_API_CONFIG="https://raw.githubusercontent.com/openline-ai/openline-customer-os/otter/deployment/infra/minikube-config/message-store.yaml"
MESSAGE_STORE_K8S_CONFIG="https://raw.githubusercontent.com/openline-ai/openline-customer-os/otter/deployment/infra/minikube-config/message-store-k8s-service.yaml"
POSTGRESQL_PERSISTENT_VOLUME_CONFIG="https://raw.githubusercontent.com/openline-ai/openline-customer-os/otter/deployment/infra/minikube-config/postgresql-persistent-volume.yaml"
POSTGRESQL_PERSISTENT_VOLUME_CLAIM_CONFIG="https://raw.githubusercontent.com/openline-ai/openline-customer-os/otter/deployment/infra/minikube-config/postgresql-persistent-volume-claim.yaml"

# Helm config
FUSIONAUTH_HELM_CONFIG="https://raw.githubusercontent.com/openline-ai/openline-customer-os/otter/deployment/infra/helm-charts/fusionauth/fusionauth.yaml"
NEO4J_HELM_CONFIG="https://raw.githubusercontent.com/openline-ai/openline-customer-os/otter/deployment/infra/helm-charts/neo4j/neo4j.yaml"
POSTGRESQL_HELM_CONFIG="https://raw.githubusercontent.com/openline-ai/openline-customer-os/otter/deployment/infra/helm-charts/postgresql/postgresql.yaml"

# Neo4j
NEO4J_CYPHER="https://raw.githubusercontent.com/openline-ai/openline-customer-os/otter/packages/server/customer-os-api/customer-os.cypher"

# PostgreSQL
EXAMPLE_PROVISIONING='https://raw.githubusercontent.com/openline-ai/openline-customer-os/otter/deployment/scripts/postgresql/example_provisioning.sql'
###########################################

mkdir openline-setup

echo "  🦦 getting Openline system config files..."

curl -sS $CUSTOMER_OS_API_CONFIG -o openline-setup/customer-os-api.yaml
if [ $? -eq 0 ]; then
    echo "  ✅ customer-os-api.yaml"
else
    echo "  ❌ customer-os-api.yaml"
fi

curl -sS $CUSTOMER_OS_API_K8S_CONFIG -o openline-setup/customer-os-api-k8s-service.yaml
if [ $? -eq 0 ]; then
    echo "  ✅ customer-os-api-k8s-service.yaml"
else
    echo "  ❌ customer-os-api-k8s-service.yaml"
fi

curl -sS $CUSTOMER_OS_API_K8S_LOADBALANCER_CONFIG -o openline-setup/customer-os-api-k8s-loadbalancer-service.yaml
if [ $? -eq 0 ]; then
    echo "  ✅ customer-os-api-k8s-loadbalancer-service.yaml"
else
    echo "  ❌ customer-os-api-k8s-loadbalancer-service.yaml"
fi

curl -sS $FUSIONAUTH_HELM_CONFIG -o openline-setup/fusionauth-values.yaml
if [ $? -eq 0 ]; then
    echo "  ✅ fusionauth-values.yaml"
else
    echo "  ❌ fusionauth-values.yaml"
fi

curl -sS $NEO4J_CYPHER -o openline-setup/customer-os.cypher
if [ $? -eq 0 ]; then
    echo "  ✅ customer-os.cypher"
else
    echo "  ❌ customer-os.cypher"
fi

curl -sS $NEO4J_HELM_CONFIG -o openline-setup/neo4j-helm-values.yaml
if [ $? -eq 0 ]; then
    echo "  ✅ neo4j-helm-values.yaml"
else
    echo "  ❌ neo4j-helm-values.yaml"
fi

curl -sS $MESSAGE_STORE_API_CONFIG -o openline-setup/message-store.yaml
if [ $? -eq 0 ]; then
    echo "  ✅ message-store.yaml"
else
    echo "  ❌ message-store.yaml"
fi

curl -sS $MESSAGE_STORE_K8S_CONFIG -o openline-setup/message-store-k8s-service.yaml
if [ $? -eq 0 ]; then
    echo "  ✅ message-store-k8s-service.yaml"
else
    echo "  ❌ message-store-k8s-service.yaml"
fi

curl -sS $POSTGRESQL_PERSISTENT_VOLUME_CONFIG -o openline-setup/postgresql-persistent-volume.yaml
if [ $? -eq 0 ]; then
    echo "  ✅ postgresql-persistent-volume.yaml"
else
    echo "  ❌ postgresql-persistent-volume.yaml"
fi

curl -sS $POSTGRESQL_PERSISTENT_VOLUME_CLAIM_CONFIG -o openline-setup/postgresql-persistent-volume-claim.yaml
if [ $? -eq 0 ]; then
    echo "  ✅ postgresql-persistent-volume-claim.yaml"
else
    echo "  ❌ postgresql-persistent-volume-claim.yaml"
fi

curl -sS $POSTGRESQL_HELM_CONFIG -o openline-setup/postgresql-values.yaml
if [ $? -eq 0 ]; then
    echo "  ✅ postgresql-values.yaml"
else
    echo "  ❌ postgresql-values.yaml"
fi

curl -sS $OPENLINE_NAMESPACE -o openline-setup/openline-namespace.json
if [ $? -eq 0 ]; then
    echo "  ✅ openline-namespace.json"
else
    echo "  ❌ openline-namespace.json"
fi

curl -sS $POSTGRES_BUILD_CONFIG -o openline-setup/build-postgres.sh
if [ $? -eq 0 ]; then
    echo "  ✅ build-postgres.sh"
else
    echo "  ❌ build-postgres.sh"
fi

curl -sS $EXAMPLE_PROVISIONING -o openline-setup/example_provisioning.sql
if [ $? -eq 0 ]; then
    echo "  ✅ build-postgres.sh"
else
    echo "  ❌ build-postgres.sh"
fi