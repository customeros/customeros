# Execute the script to clean Role nodes and it's relations https://github.com/openline-ai/openline-customer-os/issues/705

# Delete all HAS_ROLE relationships
match (r:Role)<-[rel:HAS_ROLE]-(c:Contact), (c)-[:WORKS_AS]->(r) DELETE rel;

# Verify all relations removed, return should be 0
match (r:Role)<-[rel:HAS_ROLE]-(c:Contact) return count(rel);

# Delete all WORKS relationships
match (r:Role)-[rel:WORKS]->(org:Organization), (r)-[:ROLE_IN]->(org) DELETE rel;

# Verify all relations removed, return should be 0
match (r:Role)-[rel:WORKS]->(org:Organization) return count(rel);

# Delete all Role labels, that has JobRole label
MATCH (j:JobRole) REMOVE j:Role

# Verify all Role labels removed, return should be 0
match (r:Role) return count(r);

# Delete all tenant specific Role labels
MATCH (j:JobRole_openline) REMOVE j:Role_openline
MATCH (j:JobRole_openline) REMOVE j:ROLE_openline
MATCH (j:JobRole_test) REMOVE j:Role_test
MATCH (j:JobRole_test) REMOVE j:ROLE_test

# Verify all Role labels removed, return should be 0
MATCH (n:JobRole)
WHERE SIZE(labels(n)) > 2
RETURN count(n);