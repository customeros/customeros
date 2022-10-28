# Tenant
CREATE CONSTRAINT tenant_name_unique IF NOT EXISTS ON (t:Tenant) ASSERT t.name IS UNIQUE;
CREATE(t:Tenant {name: "openline"});
CREATE(t:Tenant {name: "gaia"});

# Contact
CREATE INDEX ON :Contact(id);

# Contact group
CREATE INDEX ON :ContactGroup(id);
CREATE CONSTRAINT contact_group_name_unique IF NOT EXISTS ON (g:ContactGroup) ASSERT g.name IS UNIQUE;
#Property existence constraint requires Neo4j Enterprise Edition
#CREATE CONSTRAINT contact_group_name_not_null IF NOT EXISTS FOR (n:ContactGroup) REQUIRE n.name IS NOT NULL;

