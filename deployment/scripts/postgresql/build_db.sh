#!/bin/bash

NAMESPACE_NAME="openline"
FILES="example_provisioning.sql"
if [ "x$1" == "xlocal-kube" ]; then
  while [ -z "$pod" ]; do
    pod=$(kubectl get pods -n $NAMESPACE_NAME|grep message-store|grep Running| cut -f1 -d ' ')
    if [ -z "$pod" ]; then
      echo "message-store not ready waiting"
      sleep 1
    fi
    sleep 1
  done
  pod=$(kubectl get pods -n $NAMESPACE_NAME|grep postgresql-customer-os-dev|grep Running| cut -f1 -d ' ')
  while [ -z "$pod" ]; do
    pod=$(kubectl get pods -n $NAMESPACE_NAME|grep postgresql-customer-os-dev|grep Running| cut -f1 -d ' ')
    if [ -z "$pod" ]; then
      echo "message-store not ready waiting"
      sleep 1
    fi
    sleep 1
  done

  echo "connecting to pod $pod"
  echo $FILES |xargs cat|kubectl exec -n $NAMESPACE_NAME -i $pod -- bash -c "PGPASSWORD=$SQL_PASSWORD psql -U $SQL_USER $SQL_DATABASE"
else
  echo $FILES |xargs cat| PGPASSWORD=$SQL_PASSWORD  psql -h $SQL_HOST $SQL_USER $SQL_DATABASE
fi
