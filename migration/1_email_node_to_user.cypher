# Below script to be executed on env when https://github.com/openline-ai/openline-customer-os/issues/690 is released

#Create email node associated with a user
match (u:User) MERGE (u)-[r:EMAIL_ASSOCIATED_WITH]->(e:Email)  ON CREATE SET e.email=u.email, e.label="WORK", r.primary=true, e.id=randomUUID(), e.source=u.source, e.sourceOfTruth=u.sourceOfTruth, e.appSource=u.appSource,  e.createdAt=u.createdAt, e.updatedAt=u.updatedAt  return u, r, e;

#Check the new relation
match (u:User)-[r:EMAIL_ASSOCIATED_WITH]->(e:Email) return u, e, r

#Remove relation if wrong
match (u:User)-[r:EMAIL_ASSOCIATED_WITH]->(e:Email) detach delete e,r;

#Remove email property from user
#match (u:User)-[r:EMAIL_ASSOCIATED_WITH]->(e:Email) REMOVE u.email RETURN u

#Rename existing relations
MATCH (n)-[r:EMAILED_AT]->(m)
MERGE (n)-[:EMAIL_ASSOCIATED_WITH{primary:r.primary}]->(m)
DELETE r;