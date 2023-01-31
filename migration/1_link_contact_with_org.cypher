# Below script to be executed on env with synced data when https://github.com/openline-ai/openline-customer-os/issues/695 is released

match (c:Contact)-[:HAS_ROLE]->(r:Role)-[:WORKS]->(org:Organization) merge (c)-[:CONTACT_OF]->(org);