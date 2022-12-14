CREATE CONSTRAINT tenant_name_unique IF NOT EXISTS ON (t:Tenant) ASSERT t.name IS UNIQUE;
MERGE(t:Tenant {name: "openline"});

MATCH (t:Tenant {name:"openline"})
 MERGE (t)<-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]-(e:ExternalSystem {id:"hubspot"})
 ON CREATE SET e.name="HubSpot";

 MATCH (t:Tenant {name:"openline"})
  MERGE (t)<-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]-(e:ExternalSystem {id:"zendesk"})
  ON CREATE SET e.name="Zendesk";

MATCH (t:Tenant {name:"openline"})
    MERGE (u:User {id:"AgentSmith"})-[:USER_BELONGS_TO_TENANT]->(t)
    ON CREATE SET
    		u.firstName ="Agent",
            u.lastName="Smith",
            u.email="AgentSmith@oasis.openline.ninja",
    		u.createdAt=datetime({timezone: 'UTC'});

MATCH (t:Tenant {name:"openline"})
    MERGE (c:Contact {id:"echotest"})-[:CONTACT_BELONGS_TO_TENANT]->(t)
    ON CREATE SET
    		c.firstName ="Echo",
            c.lastName="Test",
    		c.createdAt=datetime({timezone: 'UTC'});
MATCH (c:Contact {id:"echotest"})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:"openline"})
			MERGE (c)-[r:EMAILED_AT]->(e:Email {email: "echo.test@openline.ninja"})
            ON CREATE SET e.label="MAIN", r.primary=true, e.id=randomUUID()
            ON MATCH SET e.label="MAIN", r.primary=true
			RETURN e, r

CREATE INDEX contact_id_idx IF NOT EXISTS FOR (n:Contact) ON (n.id);
CREATE INDEX contact_group_id_idx IF NOT EXISTS FOR (n:ContactGroup) ON (n.id);
CREATE INDEX company_id_idx IF NOT EXISTS FOR (n:Company) ON (n.id);
CREATE INDEX company_name_idx IF NOT EXISTS FOR (n:Company) ON (n.name);
CREATE INDEX custom_field_id_idx IF NOT EXISTS FOR (n:CustomField) ON (n.id);
CREATE INDEX field_set_id_idx IF NOT EXISTS FOR (n:FieldSet) ON (n.id);
CREATE INDEX email_id_idx IF NOT EXISTS FOR (n:Email) ON (n.id);
CREATE INDEX email_email_idx IF NOT EXISTS FOR (n:Email) ON (n.email);
CREATE INDEX phone_id_idx IF NOT EXISTS FOR (n:PhoneNumber) ON (n.id);
CREATE INDEX phone_e164_idx IF NOT EXISTS FOR (n:PhoneNumber) ON (n.e164);
CREATE INDEX entity_definition_id_idx IF NOT EXISTS FOR (n:EntityDefinition) ON (n.id);
CREATE INDEX field_set_definition_id_idx IF NOT EXISTS FOR (n:FieldSetDefinition) ON (n.id);
CREATE INDEX custom_field_definition_id_idx IF NOT EXISTS FOR (n:CustomFieldDefinition) ON (n.id);
CREATE INDEX conversation_id_idx IF NOT EXISTS FOR (n:Conversation) ON (n.id);
CREATE INDEX message_id_idx IF NOT EXISTS FOR (n:Message) ON (n.id);

:exit;