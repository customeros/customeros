#!/bin/bash

###### Script Variables #################################
NAMESPACE_NAME="openline"
FILES="./openline-setup/setup.sql"
#########################################################

echo "🦦 Provisioning Neo4j DB..."
# Provision neo4j
while [ -z "$pod" ]; do
    pod=$(kubectl get pods -n $NAMESPACE_NAME|grep neo4j-customer-os|grep Running| cut -f1 -d ' ')
    if [ -z "$pod" ]; then
      echo "  ⏳ Neo4j not ready yet, please wait..."
      sleep 2
    fi
done

started=""
while [ -z "$started" ]; do
    started=$(kubectl logs -n $NAMESPACE_NAME $pod|grep password)
    if [ -z "$started" ]; then
      echo "⏳ Neo4j waiting for app to start..."
      sleep 2
    fi
done
sleep 2

neo_output="not empty"
while  [ ! -z "$neo_output" ]; do
	echo "🦦 Provisioning Neo4j -- may take a bit, will prompt when done"
	neo_output=$(cat $NEO4J_CYPHER |kubectl run --rm -i --namespace $NAMESPACE_NAME --image "neo4j:4.4.11" cypher-shell  -- bash -c 'NEO4J_PASSWORD=StrongLocalPa\$\$ cypher-shell -a neo4j://neo4j-customer-os.openline.svc.cluster.local:7687 -u neo4j --non-interactive' 2>&1 |grep -v "see a command prompt" |grep -v "deleted")
	if [ ! -z "$neo_output" ]; then
		echo "❌ Neo4j provisioning failed, trying again"
		echo "  output: $neo_output"
		kubectl delete pod cypher-shell -n $NAMESPACE_NAME
		sleep 2
  else
    echo "✅ Neo4j provisioned"
	fi
done

echo "🦦 Provisioning PostgreSQL DB..."

SQL_USER=openline SQL_DATABABASE=openline SQL_PASSWORD=password

while [ -z "$pod" ]; do
  pod=$(kubectl get pods -n $NAMESPACE_NAME|grep message-store|grep Running| cut -f1 -d ' ')
  if [ -z "$pod" ]; then
    echo "⏳ message-store not ready yet, please wait..."
    sleep 2
  fi
done

pod=$(kubectl get pods -n $NAMESPACE_NAME|grep postgresql-customer-os-dev|grep Running| cut -f1 -d ' ')

while [ -z "$pod" ]; do
  pod=$(kubectl get pods -n $NAMESPACE_NAME|grep postgresql-customer-os-dev|grep Running| cut -f1 -d ' ')
  if [ -z "$pod" ]; then
    echo "⏳ message-store not ready yet, please wait..."
    sleep 2
  fi
done

echo "🦦 connecting to pod $pod"
echo $FILES |xargs cat|kubectl exec -n $NAMESPACE_NAME -i $pod -- bash -c "PGPASSWORD=$SQL_PASSWORD psql -U $SQL_USER $SQL_DATABASE"
if [ $? -eq 0 ]; then
  echo "✅ postgreSQL provisioned"
else
  echo "❌ postgreSQL not provisioned"
fi

rm -r openline-setup