#!/bin/bash

read -p "Which ORY project to tunnel to? [kind-babbage-plrgzvk56q]" project
project=${project:-kind-babbage-plrgzvk56q}
read -p "Which app to redirect to? [http://localhost.openline.svc.cluster.local:3001]" redirect 
redirect=${redirect:-http://localhost:3001}

ory tunnel --dev --project $project $redirect
