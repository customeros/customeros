# replace <tenant> with the tenant name

MATCH (ie:InteractionEvent)<-[rel:SENT]-(n)
MERGE (ie)-[:SENT_BY]->(n)
DELETE rel