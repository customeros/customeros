#!/bin/bash

if [ -z $(which docker) ]; then
    sudo apt update
    sudo apt install ca-certificates curl gnupg lsb-release
    if [ ! -f "/etc/apt/sources.list.d/docker.list" ]; then
        sudo mkdir -p /etc/apt/keyrings
        curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg
        echo "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu  $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
        sudo apt-get update
    fi
    sudo apt-get install -y docker-ce docker-ce-cli containerd.io docker-compose-plugin
    sudo usermod -aG docker $(whoami)
fi

if [ -z $(which minikube) ]; then
    wget -O /tmp/minikube-linux-amd64 https://storage.googleapis.com/minikube/releases/latest/minikube-linux-amd64
    sudo install -o root -g root -m 0755 /tmp/minikube-linux-amd64 /usr/local/bin/kubectl
    minikube config set driver docker
fi

if [ -z $(which kubectl) ]; then
    wget -O /tmp/kubectl "https://dl.k8s.io/release/v1.25.3/bin/linux/amd64/kubectl"
    sudo install -o root -g root -m 0755 /tmp/kubectl /usr/local/bin/kubectl
fi

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
