#!/bin/bash

NAMESPACE_NAME="openline"
FILES="./example_provisioning.sql"

echo "  🦦 Provisioning PostgreSQL"

SQL_USER=openline SQL_DATABABASE=openline SQL_PASSWORD=password

while [ -z "$pod" ]; do
  pod=$(kubectl get pods -n $NAMESPACE_NAME|grep message-store|grep Running| cut -f1 -d ' ')
  if [ -z "$pod" ]; then
    echo "  ⏳ message-store not ready, please wait..."
    sleep 2
  fi
  sleep 2
done

pod=$(kubectl get pods -n $NAMESPACE_NAME|grep postgresql-customer-os-dev|grep Running| cut -f1 -d ' ')

while [ -z "$pod" ]; do
  pod=$(kubectl get pods -n $NAMESPACE_NAME|grep postgresql-customer-os-dev|grep Running| cut -f1 -d ' ')
  if [ -z "$pod" ]; then
    echo "  ⏳ message-store not ready, please wait..."
    sleep 2
  fi
  sleep 2
done

echo "  🦦 connecting to pod $pod"
echo $FILES |xargs cat|kubectl exec -n $NAMESPACE_NAME -i $pod -- bash -c "PGPASSWORD=$SQL_PASSWORD psql -U $SQL_USER $SQL_DATABASE"
if [ $? -eq 0 ]; then
  echo "  ✅ postgreSQL provisioned"
else
  echo "  ❌ postgreSQL not provisioned"
fi


rm -r openline-setup