# Below script to be executed on env when https://github.com/openline-ai/openline-customer-os/issues/690 is released

MATCH (t:Tenant {name:"openline"})
    MERGE (u:User {id:"dev@openline.ai"})-[:USER_BELONGS_TO_TENANT]->(t)
    ON CREATE SET
    		u.firstName ="Dev",
            u.lastName="User",
            u.email="dev@openline.ai",
    		u.createdAt=datetime({timezone: 'UTC'}),
    		u.source="openline",
    		u.sourceOfTruth="openline",
    		u.appSource="manual";

#Create email node associated with a user (for all tenants)

match (u:User)-[:USER_BELONGS_TO_TENANT]-(t:Tenant{name:"openline"}) MERGE (u)-[r:HAS]->(e:Email)
ON CREATE SET e.email=u.email, e.label="WORK", r.primary=true, e.id=randomUUID(), e.source=u.source, e.sourceOfTruth=u.sourceOfTruth,
e.appSource=u.appSource,  e.createdAt=u.createdAt, e.updatedAt=u.updatedAt, e:Email_openline return u, r, e;

#Check the new relation
match (u:User)-[r:HAS]->(e:Email) return u, e, r

#Remove relation if wrong
#match (u:User)-[r:HAS]->(e:Email) detach delete e,r;

#Remove email property from user
match (u:User)-[r:HAS]->(e:Email) REMOVE u.email RETURN u

#Rename existing relations
MATCH (n)-[r:EMAILED_AT]->(m)
MERGE (n)-[:HAS{primary:r.primary}]->(m)
DELETE r;