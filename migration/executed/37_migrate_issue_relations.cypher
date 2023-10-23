MATCH (i:Issue)<-[r:REPORTED_BY]-(u:User)
MERGE (i)-[:REPORTED_BY]->(u)
DELETE r;

MATCH (i:Issue)<-[r:FOLLOWS]-(u:User)
MERGE (i)-[:FOLLOWED_BY]->(u)
DELETE r;

MATCH (i:Issue)<-[r:IS_ASSIGNED_TO]-(u:User)
MERGE (i)-[:ASSIGNED_TO]->(u)
DELETE r;
