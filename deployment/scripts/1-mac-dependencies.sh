#!/bin/bash

echo "⏳ Openline dependency check..."

# Xcode
xcode-select -p
if [ $? -eq 0 ]; then
    echo "✅ Xcode"
else
    echo "  🦦 Installing Xcode.  This may take awhile, please let the script do it's thing.  It will prompt when completed."
    xcode-select --install
    wait
    if [ $? -eq 0 ]; then
        echo "✅ Xcode"
    else
        echo "❌ Xcode installation failed"
        exit 1
    fi
fi

# Homebrew
if [[ $(brew --version) == *"Homebrew"* ]]; then
    echo "✅ Homebrew"
else
    echo "🦦 Installing Homebrew..."
    /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
    wait
    if [ $? -eq 0 ]; then
        echo "✅ Homebrew"
    else
        echo "❌ Homebrew installation failed"
        exit 1
    fi
fi

# Docker
if [[ $(docker --version) == *"Docker version"* ]]; then
    echo "✅ Docker"
else
    echo "🦦 Installing Docker..."
    brew install docker
    wait
    if [ $? -eq 0 ]; then
        echo "✅ Docker"
    else
        echo "❌ Docker installation failed"
        exit 1
    fi
fi

# Colima
if [[ $(colima version) == *"colima version"* ]]; then
    echo "✅ Colima"
else
    echo "🦦 Installing Colima..."
    brew install colima
    wait
    if [ $? -eq 0 ]; then
        echo "✅ Colima"
    else
        echo "❌ Colima installation failed"
        exit 1
    fi
fi

# Kubectl
if [[ $(kubectl version) == *"Client Version"* ]]; then
    echo "✅ kubectl"
else
    echo "🦦 Installing kubectl..."
    brew install kubectl
    wait
    if [ $? -eq 0 ]; then
        echo "✅ kubectl"
    else
        echo "❌ kubectl installation failed"
        exit 1
    fi
fi

# Helm
if [[ $(helm version) == *"version.BuildInfo"* ]]; then
    echo "✅ Helm"
else
    echo "🦦 Installing Helm..."
    brew install helm
    wait
    if [ $? -eq 0 ]; then
        echo "✅ Helm"
    else
        echo "❌ Helm installation failed."
        exit 1
    fi
fi