# https://github.com/openline-ai/openline-customer-os/issues/857

# QUERY PER TENANT. Execute before release, link Email and Tenant nodes

:param { tenant: "test" };

MATCH (t:Tenant {name:$tenant})<-[USER_BELONGS_TO_TENANT]-(:User)-[HAS]->(e:Email)
MERGE (e)-[:EMAIL_ADDRESS_BELONGS_TO_TENANT]->(t);

MATCH (t:Tenant {name:$tenant})<-[CONTACT_BELONGS_TO_TENANT]-(:Contact)-[HAS]->(e:Email)
MERGE (e)-[:EMAIL_ADDRESS_BELONGS_TO_TENANT]->(t);

MATCH (t:Tenant {name:$tenant})<-[ORGANIZATION_BELONGS_TO_TENANT]-(:Organization)-[HAS]->(e:Email)
MERGE (e)-[:EMAIL_ADDRESS_BELONGS_TO_TENANT]->(t);


# Move label property from Email to relationship, execute before and after release

MATCH (:User)-[rel:HAS]->(e:Email)
WHERE e.label is not null
WITH e, rel
SET rel.label=e.label
WITH e
REMOVE e.label;

MATCH (:Organization)-[rel:HAS]->(e:Email)
WHERE e.label is not null
WITH e, rel
SET rel.label=e.label
WITH e
REMOVE e.label;

MATCH (:Contact)-[rel:HAS]->(e:Email)
WHERE e.label is not null
WITH e, rel
SET rel.label=e.label
WITH e
REMOVE e.label;

# QUERY PER TENANT Move relationships to the first node in duplicates, execute after release

MATCH (e:Email_test)
WITH e order by e.createdAt
WITH e.email AS email, collect(e) AS emails WHERE size(emails) > 1
WITH head(emails) as firstNode, tail(emails) as otherNodes
UNWIND otherNodes as otherNode
WITH firstNode, otherNode
MATCH (otherNode)<-[rel:HAS]-(n)
MERGE (firstNode)<-[newRel:HAS]-(n)
ON CREATE SET newRel.primary=rel.primary, newRel.label=rel.label
DELETE rel;

# QUERY PER TENANT. Delete duplicate email nodes. execute after release

MATCH (e:Email_test)
WITH e order by e.createdAt
WITH e.email AS email, collect(e) AS emails WHERE size(emails) > 1
WITH head(emails) as firstNode, tail(emails) as otherNodes
UNWIND otherNodes as otherNode
WITH firstNode, otherNode
WHERE NOT (otherNode)<-[:HAS]-()
detach delete otherNode;

# QUERY PER TENANT. Verification query after release. Should return 0.
MATCH (e:Email_test)
WITH e.email AS email, collect(e) AS emails WHERE size(emails) > 1
return count(email) as duplicates;