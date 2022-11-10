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

# Setup local environment
#### 1. Go to 
```
deployment/k8s/local-minikube
```
#### 2. Install prerequisites, execute:

* For Mac:
```
0-mac-install-prerequisites.sh
```
* For Ubuntu:
```
0-ubuntu-install-prerequisites.sh
``` 
#### 3. Execute
```
1-deploy-customer-os-base-infrastructure-local.sh
2-build-deploy-customer-os-local-images.sh
``` 

#### 4. Check all services are running:
```
kubectl get pods -n openline
```
#### 5. Setup Port Forwards
# Port FWD
run the script
```
./port-forwarding.sh

#### 6. Open in browser
```
http://127.0.0.1:10010/
```

#### 7. Create first contact in your tenant workspace
```
mutation CreateContact {
  createContact(input: {firstName: "Mr", lastName:"Otter"}) {
    firstName
    lastName
    createdAt
  }
}
```

