# Below script to be executed on env with synced data after release of https://github.com/openline-ai/openline-customer-os/issues/824

# Remove NOT_SET tag
MATCH(tag:Tag {name:"NOT_SET"}) DETACH DELETE tag;

# Remove tag for "evangelist" execute each tenant

MATCH (t:Tenant {name:"test"}),
(t)<--(c:Contact)-->(cf:CustomField {name:"Hubspot Lifecycle Stage", textValue: "evangelist"}),
(c)-[rel:TAGGED]->(tag:Tag {name:"CUSTOMER"})
delete rel;

# Verify query, count should be 0

MATCH (t:Tenant {name:"test"}),
(t)<--(c:Contact)-->(cf:CustomField {name:"Hubspot Lifecycle Stage", textValue: "evangelist"}),
(c)-[rel:TAGGED]->(tag:Tag {name:"CUSTOMER"})
return count(c);

# Add prospect tag to non customers execute each tenant

MATCH (t:Tenant {name:"test"}),
(t)<--(c:Contact)-->(cf:CustomField {name:"Hubspot Lifecycle Stage"}),
(c)-[rel:TAGGED]->(tag:Tag {name:"CUSTOMER"})-->(t)
WHERE cf.textValue <> "customer"
WITH c,rel,tag,t
MATCH (t)<--(pt:Tag {name:"PROSPECT"})
MERGE (c)-[newrel:TAGGED]->(pt)
ON CREATE SET newrel.taggedAt = rel.taggedAt
WITH rel
DELETE rel;

# Verify query, count should be 0

MATCH (t:Tenant {name:"test"}),
(t)<--(c:Contact)-->(cf:CustomField {name:"Hubspot Lifecycle Stage"}),
(c)-[rel:TAGGED]->(tag:Tag {name:"CUSTOMER"})-->(t)
WHERE cf.textValue <> "customer"
return count(c);
