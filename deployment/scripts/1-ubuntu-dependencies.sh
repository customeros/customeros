#!/bin/bash

echo "  ⏳ Openline dependency check..."

# Docker
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
    if [ $? -eq 0 ]; then
        echo "  ✅ Docker"
    else
        echo "  ❌ Docker"
    fi
    sudo usermod -aG docker $(whoami)
else
    echo "  ✅ Docker"
fi

# Minikube
if [ -z $(which minikube) ]; then
    wget -O /tmp/minikube-linux-amd64 https://storage.googleapis.com/minikube/releases/latest/minikube-linux-amd64
    sudo install -o root -g root -m 0755 /tmp/minikube-linux-amd64 /usr/local/bin/minikube
    minikube config set driver docker
    if [ $? -eq 0 ]; then
        echo "  ✅ Minikube"
    else
        echo "  ❌ Minikube"
    fi
else
    echo "  ✅ Minikube"
fi

# Helm
if [ -z $(which helm) ]; then
    mkdir -p /tmp/helm/
    wget -O /tmp/helm/helm.tar.gz "https://get.helm.sh/helm-v3.10.2-linux-amd64.tar.gz"
    cd /tmp/helm/
    tar -zxvf helm.tar.gz
    sudo install -o root -g root -m 0755 /tmp/helm/linux-amd64/helm /usr/local/bin/helm
    if [ $? -eq 0 ]; then
        echo "  ✅ Helm"
    else
        echo "  ❌ Helm"
    fi
    cd -
else
    echo "  ✅ Helm"
fi

# Kubernetes
if [ -z $(which kubectl) ]; then
    wget -O /tmp/kubectl "https://dl.k8s.io/release/v1.25.3/bin/linux/amd64/kubectl"
    sudo install -o root -g root -m 0755 /tmp/kubectl /usr/local/bin/kubectl
    if [ $? -eq 0 ]; then
        echo "  ✅ Kubernetes"
    else
        echo "  ❌ Kubernetes"
    fi
else
    echo "  ✅ Kubernetes"
fi
