# Below script to be executed on env with synced data when https://github.com/openline-ai/openline-customer-os/issues/703 is released

MATCH (r:Role)-[rel:WORKS]->(o:Organization) MERGE (r)-[:ROLE_IN]->(o);
MATCH (r:Role)<-[rel:HAS_ROLE]->(c:Contact) MERGE (r)<-[:WORKS_AS]-(c);
MATCH (r:Role) SET r:JobRole;

# Following script execute for all tenants
MATCH (j:JobRole)--(c:Contact)--(t:Tenant {name:"openline"}) SET j:JobRole_openline;