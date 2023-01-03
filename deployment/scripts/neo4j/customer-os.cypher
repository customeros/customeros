CREATE CONSTRAINT tenant_name_unique IF NOT EXISTS ON (t:Tenant) ASSERT t.name IS UNIQUE;
MERGE(t:Tenant {id:"2086420f-05fd-42c8-a7f3-a9688e65fe53", name: "openline"});

MATCH (t:Tenant {name:"openline"})
 MERGE (t)<-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]-(e:ExternalSystem {id:"hubspot"})
 ON CREATE SET e.name="HubSpot";

MATCH (t:Tenant {name:"openline"})
  MERGE (t)<-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]-(e:ExternalSystem {id:"zendesk"})
  ON CREATE SET e.name="Zendesk";

MATCH (t:Tenant {name:"openline"})
  MERGE (t)<-[:CONTACT_TYPE_BELONGS_TO_TENANT]-(ct:ContactType {name:"CUSTOMER"})
  ON CREATE SET ct.id=randomUUID(), ct:ContactType_openline;

MATCH (t:Tenant {name:"openline"})
  MERGE (t)<-[:CONTACT_TYPE_BELONGS_TO_TENANT]-(ct:ContactType {name:"SUPPLIER"})
  ON CREATE SET ct.id=randomUUID(), ct:ContactType_openline;

MATCH (t:Tenant {name:"openline"})
  MERGE (t)<-[:CONTACT_TYPE_BELONGS_TO_TENANT]-(ct:ContactType {name:"INVESTOR"})
  ON CREATE SET ct.id=randomUUID(), ct:ContactType_openline;

MATCH (t:Tenant {name:"openline"})
  MERGE (t)<-[:CONTACT_TYPE_BELONGS_TO_TENANT]-(ct:ContactType {name:"NOT_SET"})
  ON CREATE SET ct.id=randomUUID(), ct:ContactType_openline;

MATCH (t:Tenant {name:"openline"})
  MERGE (t)<-[:ORGANIZATION_TYPE_BELONGS_TO_TENANT]-(ot:OrganizationType {name:"COMPANY"})
  ON CREATE SET ot.id=randomUUID(), ot:OrganizationType_openline;

MATCH (t:Tenant {name:"openline"})
    MERGE (u:User {id:"AgentSmith"})-[:USER_BELONGS_TO_TENANT]->(t)
    ON CREATE SET
    		u.firstName ="Agent",
            u.lastName="Smith",
            u.email="AgentSmith@oasis.openline.ninja",
    		u.createdAt=datetime({timezone: 'UTC'}),
    		u:User_openline;

MATCH (t:Tenant {name:"openline"})
    MERGE (c:Contact {id:"echotest"})-[:CONTACT_BELONGS_TO_TENANT]->(t)
    ON CREATE SET
    		c.firstName ="Echo",
            c.lastName="Test",
    		c.createdAt=datetime({timezone: 'UTC'});

MATCH (c:Contact {id:"echotest"})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:"openline"})
			MERGE (c)-[r:EMAILED_AT]->(e:Email {email: "echo@oasis.openline.ai"})
            ON CREATE SET e.label="MAIN", r.primary=true, e.id=randomUUID(), e:Email_openline
            ON MATCH SET e.label="MAIN", r.primary=true;

CREATE INDEX contact_id_idx IF NOT EXISTS FOR (n:Contact) ON (n.id);
CREATE INDEX contact_type_id_idx IF NOT EXISTS FOR (n:ContactType) ON (n.id);
CREATE INDEX contact_group_id_idx IF NOT EXISTS FOR (n:ContactGroup) ON (n.id);
CREATE INDEX organization_id_idx IF NOT EXISTS FOR (n:Organization) ON (n.id);
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
