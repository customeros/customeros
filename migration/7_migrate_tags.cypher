# Below script to be executed on env with synced data twice, before and after release of https://github.com/openline-ai/openline-customer-os/issues/824

# Add PROSPECT tag for each tenant

MATCH (t:Tenant {name:"openline"})
  MERGE (t)<-[:TAG_BELONGS_TO_TENANT]-(tag:Tag {name:"PROSPECT"})
  ON CREATE SET tag.id=randomUUID(),
                tag.createdAt=datetime({timezone: 'UTC'}),
                tag.updatedAt=datetime({timezone: 'UTC'}),
                tag.source="openline",
                tag.appSource="manual";
MATCH (t:Tenant {name:"openline"})<-[:TAG_BELONGS_TO_TENANT]-(tag:Tag) SET tag:Tag_openline;