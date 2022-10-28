#! /bin/sh
echo "Check if CUSTOMER_OS_HOME is set below"
echo "CUSTOMER_OS_HOME = "$CUSTOMER_OS_HOME

chmod 755 1-deploy-customer-os-base-infrastructure-local.sh
chmod 755 2-build-deploy-customer-os-local-images.sh
chmod 755 port-forwarding.sh

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

echo "Completed."