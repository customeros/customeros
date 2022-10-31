# Deployment for customerOS on development env using K8s

## Prerequisites
- [Homebrew](https://brew.sh/)
- [Docker](https://www.docker.com/)
- [Minikube](https://minikube.sigs.k8s.io/docs/start/)
- [Helm](https://helm.sh/)

# Install Prerequisites 

## Setup Environment for OSX

TODO create script to install Prerequisites

### Brew
```
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
```

### Docker
```
https://docs.docker.com/desktop/install/mac-install/
```

### Minikube 

```
brew install minikube
```

### Helm
```
brew install helm
```

### Helm repositories dependencies  
```
helm repo add bitnami https://charts.bitnami.com/bitnami
```

# Start on local dev env
Go to deployment/k8s/local-minikube
and execute :
0-mac-install-prerequisites.sh
1-deploy-customer-os-base-infrastructure-local.sh
2-build-deploy-customer-os-local-images.sh

## Setup Environment for Ubuntu

Go to eployment/k8s/local-minikube
and execute :
``` 
0-ubuntu-install-prerequisites.sh
1-deploy-customer-os-base-infrastructure-local.sh
2-build-deploy-customer-os-local-images.sh
``` 

# Port FWD
```
#PostgreSQL DB
kubectl port-forward --namespace openline-development svc/postgresql-customer-os-dev 5432:5432
```
```
#Neo4j DB
kubectl port-forward --namespace openline-development svc/neo4j-customer-os 7687:7687
```

or run the script for both DBs
```
./port-forwarding.sh
```


