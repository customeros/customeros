# https://github.com/openline-ai/openline-customer-os/issues/857

#========== Link PhoneNumber with Tenant
#========== Execute before release, and after release if exists Phone Numbers with no Tenant

MATCH (t:Tenant)<-[USER_BELONGS_TO_TENANT]-(:User)--(p:PhoneNumber)
MERGE (p)-[:PHONE_NUMBER_BELONGS_TO_TENANT]->(t);

MATCH (t:Tenant)<-[CONTACT_BELONGS_TO_TENANT]-(:Contact)--(p:PhoneNumber)
MERGE (p)-[:PHONE_NUMBER_BELONGS_TO_TENANT]->(t);

#========== Update relation, execute before and after release

MATCH (c:Contact)-[rel:PHONE_ASSOCIATED_WITH]->(p:PhoneNumber)
MERGE (c)-[newRel:HAS]->(p)
ON CREATE SET newRel.primary=rel.primary
DELETE rel;

#========== Move label property from PhoneNumber to relationship, execute before and after release

MATCH (:User)-[rel:HAS]->(p:PhoneNumber)
WHERE p.label is not null
SET rel.label=p.label
REMOVE p.label;

MATCH (:Contact)-[rel:HAS]->(p:PhoneNumber)
WHERE p.label is not null
SET rel.label=p.label
REMOVE p.label;

#========== Set createdAt / updatedAt where missing, execute before and after release

MATCH (p:PhoneNumber)
WHERE p.createdAt is null
SET p.createdAt=datetime({timezone: 'UTC'});

MATCH (p:PhoneNumber)
WHERE p.updatedAt is null
SET p.updatedAt=p.createdAt;

#========== QUERY PER TENANT Move relationships to the first node in duplicates. execute after release

MATCH (p:PhoneNumber_test)
WHERE p.e164 is not null AND p.e164 <> ''
WITH p order by p.createdAt
WITH p.e164 AS e164, collect(p) AS numbers WHERE size(numbers) > 1
WITH head(numbers) as firstNode, tail(numbers) as otherNodes
UNWIND otherNodes as otherNode
WITH firstNode, otherNode
MATCH (otherNode)<-[rel:HAS]-(n)
MERGE (firstNode)<-[newRel:HAS]-(n)
ON CREATE SET newRel.primary=rel.primary, newRel.label=rel.label
DELETE rel;

#==========  QUERY PER TENANT. Delete duplicate phone number nodes. execute after release

MATCH (p:PhoneNumber_test)
WHERE p.e164 is not null AND p.e164 <> ''
WITH p order by p.createdAt
WITH p.e164 AS e164, collect(p) AS numbers WHERE size(numbers) > 1
WITH head(numbers) as firstNode, tail(numbers) as otherNodes
UNWIND otherNodes as otherNode
WITH firstNode, otherNode
WHERE NOT (otherNode)<-[:HAS]-()
detach delete otherNode;

#========== Verification query after release. Should return 0.

MATCH (t:Tenant)--(p:PhoneNumber)
WHERE 'PhoneNumber_'+t.name in labels(p) AND p.e164 is not null AND p.e164 <> ''
WITH t, p.e164 AS e164, collect(p) AS numbers WHERE size(numbers) > 1
return count(e164) as duplicates;