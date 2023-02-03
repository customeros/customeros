# Below script to be executed on env with synced data when https://github.com/openline-ai/openline-customer-os/issues/692 is released

# Create a tag for each contact type
MATCH (ct:ContactType)-[:CONTACT_TYPE_BELONGS_TO_TENANT]->(t:Tenant)
MERGE (t)<-[:TAG_BELONGS_TO_TENANT]-(tg:Tag {name: ct.name})-[:TEMP_LINK]->(ct)
ON CREATE SET   tg.createdAt=datetime({timezone: 'UTC'}),
                tg.updatedAt=datetime({timezone: 'UTC'}),
                tg.source=ct.source,
                tg.appSource=ct.appSource,
                tg.id=randomUUID(),
                tg.tenant=t.name;

# Link contact of each contact type with Tag
MATCH (t:Tenant)<-[:CONTACT_TYPE_BELONGS_TO_TENANT]-(ct:ContactType)<-[:IS_OF_TYPE]-(c:Contact)--(t),
        (ct)<-[:TEMP_LINK]-(tag:Tag)--(t)
MERGE (c)-[r:TAGGED]->(tag)
ON CREATE SET r.taggedAt=datetime({timezone: 'UTC'});

# Set tenant specific label to tag node. execute per each tenant
MATCH (tag:Tag {tenant:"openline"}) SET tag:Tag_openline;
MATCH (tag:Tag {tenant:"test"}) SET tag:Tag_test;
# MATCH (tag:Tag {tenant:"any other tenant here"}) SET tag:TAG_any other tenant here:

# Clean tenant name on tag node
MATCH (t:Tag) REMOVE t.tenant;

