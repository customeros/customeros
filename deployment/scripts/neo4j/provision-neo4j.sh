#!/bin/bash

###### Script Variables #################################
NAMESPACE_NAME="openline"
NEO4J_CYPHER="./openline-setup/customer-os.cypher"
#########################################################

VERSION=$(kubectl describe pod customer-db-neo4j-0 -n openline | grep HELM_NEO4J_VERSION | tr -s ' ' |cut -d ' ' -f 3)

neo_output="not empty"
while  [ ! -z "$neo_output" ]; do
	neo_output=$(cat $NEO4J_CYPHER |kubectl run --rm -i --namespace $NAMESPACE_NAME --image "neo4j:$VERSION" cypher-shell  -- bash -c 'NEO4J_PASSWORD=StrongLocalPa\$\$ cypher-shell -a neo4j://customer-db-neo4j.openline.svc.cluster.local:7687 -u neo4j --non-interactive' 2>&1 |grep -v "see a command prompt" |grep -v "deleted")
	if [ ! -z "$neo_output" ]; then
		kubectl delete pod cypher-shell -n $NAMESPACE_NAME
		sleep 2
	fi
done
