# Deployment for customerOS on development env using K8s

## Prerequisites
- [Homebrew](https://brew.sh/)
- [Docker](https://www.docker.com/)
- [Minikube](https://minikube.sigs.k8s.io/docs/start/)
- [Helm](https://helm.sh/)

# Install Prerequisites 

Set Customer OS home env variable by running :
```
export CUSTOMER_OS_HOME=~/[path to the checked out folder]/openline-customer-os
```
or added to .zprofile for Mac or .bashrc for linux


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
helm repo add hashicorp https://helm.releases.hashicorp.com
helm repo add bitnami https://charts.bitnami.com/bitnami
```


# Start on local dev env
Go to $CUSTOMER_OS_HOME/deployment/k8s/local-minikube
and execute :
0-mac-install-prerequisites.sh
1-deploy-customer-os-base-infrastructure-local.sh
2-build-deploy-customer-os-local-images.sh

# Port FWD
```
#Consul UI
kubectl port-forward consul-server-0 --namespace openline-development 8500:8500

#PostgreSQL DB
kubectl port-forward --namespace default svc/postgresql-dev 5432:5432
```