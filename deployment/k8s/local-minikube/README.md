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

#### 5.For customerOs API, prepare neo4j constraints and onboard first tenant
* Connect to neo4j, when aksed provide password `StrongLocalPa$$`
```
kubectl run --rm -it --namespace "openline" --image "neo4j:4.4.11" cypher-shell \
     -- cypher-shell -a "neo4j://neo4j-customer-os.openline.svc.cluster.local:7687" -u "neo4j"
```
* Execute cypher commands located in the file
```
customer-os-api/init-customer-os.cypher
```

#### 6. Port FWD for customerOS API
```
kubectl port-forward --namespace openline svc/customer-os-api-service 10010:10010 &
```

#### 7. Open in browser
```
http://127.0.0.1:10010/
```

#### 8. Create first contact in your tenant workspace
```
mutation CreateContact {
  createContact(input: {firstName: "Mr", lastName:"Otter"}) {
    firstName
    lastName
    createdAt
  }
}
```

# Port FWD
```
#PostgreSQL DB
kubectl port-forward --namespace openline svc/postgresql-customer-os-dev 5432:5432
```
```
#Neo4j DB
kubectl port-forward --namespace openline svc/neo4j-customer-os 7687:7687
```

or run the script for both DBs
```
./port-forwarding.sh
```


