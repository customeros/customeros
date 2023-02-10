# https://github.com/openline-ai/openline-customer-os/issues/857

# Execute before release for each tenant, link Email and Tenant nodes

:param { tenant: "openline" };

MATCH (t:Tenant {name:$tenant})<-[USER_BELONGS_TO_TENANT]-(:User)-[HAS]->(e:Email)
MERGE (e)-[:EMAIL_ADDRESS_BELONGS_TO_TENANT]->(t);

MATCH (t:Tenant {name:$tenant})<-[CONTACT_BELONGS_TO_TENANT]-(:Contact)-[HAS]->(e:Email)
MERGE (e)-[:EMAIL_ADDRESS_BELONGS_TO_TENANT]->(t);

MATCH (t:Tenant {name:$tenant})<-[ORGANIZATION_BELONGS_TO_TENANT]-(:Organization)-[HAS]->(e:Email)
MERGE (e)-[:EMAIL_ADDRESS_BELONGS_TO_TENANT]->(t);


# Move label property from Email to relationship, not tenant specific

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