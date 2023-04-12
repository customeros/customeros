# Execute before release and after release per tenant
# set tenant name in first and last line

MATCH (c:Contact_test)-[rel:CONTACT_OF]->(o:Organization_test)
WHERE not exists ((c)--(:JobRole)--(o))
MERGE (c)-[:WORKS_AS]->(j:JobRole)-[:ROLE_IN]->(o)
		ON CREATE SET j.id=randomUUID(),
					j.primary=true,
					j.source=c.source,
					j.sourceOfTruth=c.sourceOfTruth,
					j.appSource=c.appSource,
					j.jobTitle="",
					j.createdAt=c.createdAt,
					j.updatedAt=c.updatedAt,
					j:JobRole_test
DELETE rel;