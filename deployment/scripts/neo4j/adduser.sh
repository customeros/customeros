#!/bin/bash

###### Script Variables #################################
NAMESPACE_NAME="openline"
NEO4J_CYPHER="./create-user.cypher.template"
#########################################################

VERSION=$(kubectl describe pod customer-db-neo4j-0 -n openline | grep HELM_NEO4J_VERSION | tr -s ' ' |cut -d ' ' -f 3)

read -p "Email: " USER_EMAIL
read -p "First Name: " USER_FIRST_NAME
read -p "Last Name: " USER_LAST_NAME

SCRIPT=$(cat $NEO4J_CYPHER | sed "s/!USER_EMAIL!/$USER_EMAIL/g" | sed "s/!USER_FIRST_NAME!/$USER_FIRST_NAME/g" | sed "s/!USER_LAST_NAME!/$USER_LAST_NAME/g")
neo_output="not empty"
while  [ ! -z "$neo_output" ]; do
	neo_output=$(echo $SCRIPT |kubectl run --rm -i --namespace $NAMESPACE_NAME --image "neo4j:$VERSION" cypher-shell  -- bash -c 'NEO4J_PASSWORD=StrongLocalPa\$\$ cypher-shell -a neo4j://customer-db-neo4j.openline.svc.cluster.local:7687 -u neo4j --non-interactive' 2>&1 |grep -v "see a command prompt" |grep -v "deleted")
	if [ ! -z "$neo_output" ]; then
		kubectl delete pod cypher-shell -n $NAMESPACE_NAME
		sleep 2
	fi
done
