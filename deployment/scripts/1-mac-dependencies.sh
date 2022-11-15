#!/bin/bash

echo "  ⏳ Openline dependency check..."

# Xcode
xcode-select -p
if [ $? -eq 0 ]; then
    echo "  ✅ Xcode"
else
    echo "  🦦 Installing Xcode.  This may take awhile, please let the script do it's thing.  It will prompt when completed."
    xcode-select --install
    if [ $? -eq 0 ]; then
        echo "  ✅ Xcode"
    else
        echo "  ❌ Xcode installation failed"
    fi
fi

# Homebrew
if [[ $(brew --version) == *"Homebrew"* ]];
    then
        echo "  ✅ Homebrew"
    else
        echo "  🦦 Installing Homebrew..."
        /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
        if [ $? -eq 0 ]; then
            echo "  ✅ Homebrew"
        else
            echo "  ❌ Homebrew installation failed."
        fi
fi

# Docker
if [[ $(docker --version) == *"Docker version"* ]];
    then
        echo "  ✅ Docker"
    else
        echo "  🦦 Installing Docker..."
        echo "  ❗️ This can take a while, please let the script do it's thing.  It will prompt when completed."
        softwareupdate --install-rosetta
        
        if [[ $(arch) == 'arm64' ]]; 
            then
                echo "  Installing Apple silicon version..."
                curl -L https://desktop.docker.com/mac/main/arm64/Docker.dmg --output openline-setup/Docker.dmg
            else
                echo "  Installing Intel version..."
                curl -L https://desktop.docker.com/mac/main/amd64/Docker.dmg --output openline-setup/Docker.dmg
        fi

        sudo hdiutil attach openline-setup/Docker.dmg
        sudo /Volumes/Docker/Docker.app/Contents/MacOS/install
        sudo hdiutil detach /Volumes/Docker

        echo "  ✅ Docker"
        echo "  ❗️Please open Docker desktop via the GUI to initialize the application before proceeding."
        rm -r openline-setup/Docker.dmg
        echo "  Attempting to open Docker..."
        open -a Docker.app
        read -p "  => Press enter to continue once Docker GUI has opened..."
fi

# Minikube
if [[ $(minikube version) == *"minikube version"* ]];
    then
        echo "  ✅ Minikube"
    else
        echo "  🦦 Installing Minikube..."
        brew install minikube
        if [ $? -eq 0 ]; then
            echo "  ✅ Minikube"
        else
            echo "  ❌ Minikube installation failed."
        fi
fi

# Helm
if [[ $(helm version) == *"version.BuildInfo"* ]];
    then
        echo "  ✅ Helm"
    else
        echo "  🦦 Installing Helm..."
        brew install helm
        if [ $? -eq 0 ]; then
            echo "  ✅ Helm"
        else
            echo "  ❌ Helm installation failed."
        fi
fi