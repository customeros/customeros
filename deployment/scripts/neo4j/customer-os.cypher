MERGE(t:Tenant {id:"2086420f-05fd-42c8-a7f3-a9688e65fe53", name: "openline"})
 ON CREATE SET t.createdAt=datetime({timezone: 'UTC'});

MATCH (t:Tenant {name:"openline"})
 MERGE (t)<-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]-(e:ExternalSystem {id:"hubspot"})
 ON CREATE SET e.name="HubSpot", e.createdAt=datetime({timezone: 'UTC'});

MATCH (t:Tenant {name:"openline"})
  MERGE (t)<-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]-(e:ExternalSystem {id:"zendesk_support"})
  ON CREATE SET e.name="Zendesk Support", e.createdAt=datetime({timezone: 'UTC'});

MATCH (t:Tenant {name:"openline"})
  MERGE (t)<-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]-(e:ExternalSystem {id:"zendesk_sell"})
  ON CREATE SET e.name="Zendesk Sell", e.createdAt=datetime({timezone: 'UTC'});

MATCH (t:Tenant {name:"openline"})
  MERGE (t)<-[:CONTACT_TYPE_BELONGS_TO_TENANT]-(ct:ContactType {name:"CUSTOMER"})
  ON CREATE SET ct.id=randomUUID(), ct.createdAt=datetime({timezone: 'UTC'});

MATCH (t:Tenant {name:"openline"})
  MERGE (t)<-[:CONTACT_TYPE_BELONGS_TO_TENANT]-(ct:ContactType {name:"SUPPLIER"})
  ON CREATE SET ct.id=randomUUID(), ct.createdAt=datetime({timezone: 'UTC'});

MATCH (t:Tenant {name:"openline"})
  MERGE (t)<-[:CONTACT_TYPE_BELONGS_TO_TENANT]-(ct:ContactType {name:"INVESTOR"})
  ON CREATE SET ct.id=randomUUID(), ct.createdAt=datetime({timezone: 'UTC'});

MATCH (t:Tenant {name:"openline"})
  MERGE (t)<-[:CONTACT_TYPE_BELONGS_TO_TENANT]-(ct:ContactType {name:"NOT_SET"})
  ON CREATE SET ct.id=randomUUID(), ct.createdAt=datetime({timezone: 'UTC'});

MATCH (t:Tenant {name:"openline"})
  MERGE (t)<-[:ORGANIZATION_TYPE_BELONGS_TO_TENANT]-(ot:OrganizationType {name:"COMPANY"})
  ON CREATE SET ot.id=randomUUID(), ot.createdAt=datetime({timezone: 'UTC'});

MATCH (t:Tenant {name:"openline"})
    MERGE (u:User {id:"development@openline.ai"})-[:USER_BELONGS_TO_TENANT]->(t)
    ON CREATE SET
    		u.firstName ="Development",
            u.lastName="User",
    		u.createdAt=datetime({timezone: 'UTC'}),
    		u.source="openline",
    		u.sourceOfTruth="openline",
    		u.appSource="manual";

MATCH (u:User {id:"development@openline.ai"})
    MERGE (:Email {
                      id: randomUUID(),
                      email: "development@openline.ai",
                      label: "WORK",
                      u.source="openline",
                      u.sourceOfTruth="openline",
                      u.appSource="manual",
                      u.createdAt=datetime({timezone: 'UTC'}),
                      u.updatedAt=datetime({timezone: 'UTC'})
                    })<-[:HAS {primary:true}]-(u);

MATCH (t:Tenant {name:"openline"})
    MERGE (c:Contact {id:"echotest"})-[:CONTACT_BELONGS_TO_TENANT]->(t)
    ON CREATE SET
    		c.firstName ="Echo",
            c.lastName="Test",
    		c.createdAt=datetime({timezone: 'UTC'}),
    		c.source="openline",
            c.sourceOfTruth="openline",
            c.appSource="manual";

MATCH (c:Contact {id:"echotest"})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:"openline"})
			MERGE (c)-[r:HAS]->(e:Email {email: "echo@oasis.openline.ai"})
            ON CREATE SET e.label="MAIN", r.primary=true, e.id=randomUUID(), e.createdAt=datetime({timezone: 'UTC'}),
                e.source="openline", e.sourceOfTruth="openline", e.appSource="manual"
            ON MATCH SET e.label="MAIN", r.primary=true;

MATCH (t:Tenant {name:"openline"})
MERGE (o:Conversation{id:"echotest"}) ON CREATE SET  o:Conversation_openline, o.messageCount=1, o.updatedAt=datetime({timezone: 'UTC'}), o.startedAt=datetime({timezone: 'UTC'}), o.initiatorFirstName="", o.initiatorLastName="", o.initiatorUsername="echo@oasis.openline.ai", o.initiatorType="CONTACT", o.lastSenderId="echo@oasis.openline.ai", o.lastSenderType="", o.lastSenderFirstName="", o.lastSenderLastName="", o.lastContentPreview="Hello world!", o.status="ACTIVE", o.channel="WEB_CHAT",  o:Conversation_openline WITH DISTINCT t, o
OPTIONAL MATCH (c:Contact)-[:CONTACT_BELONGS_TO_TENANT]->(t) WHERE c.id ="echotest"  MERGE (c)-[:PARTICIPATES]->(o)  RETURN o;

MATCH (t:Tenant {name:"openline"})<-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]-(e:ExternalSystem) SET e:ExternalSystem_openline;
MATCH (t:Tenant {name:"openline"})<-[:CONTACT_TYPE_BELONGS_TO_TENANT]-(ct:ContactType) SET ct:ContactType_openline;
MATCH (t:Tenant {name:"openline"})<-[:ORGANIZATION_TYPE_BELONGS_TO_TENANT]-(ot:OrganizationType) SET ot:OrganizationType_openline;
MATCH (t:Tenant {name:"openline"})<-[:USER_BELONGS_TO_TENANT]-(u:User) SET u:User_openline;
MATCH (t:Tenant {name:"openline"})<-[:CONTACT_BELONGS_TO_TENANT]-(c:Contact) SET c:Contact_openline;
MATCH (c:Contact_openline)-[:HAS]->(e:Email) SET e:Email_openline;

CREATE CONSTRAINT tenant_name_unique IF NOT EXISTS ON (t:Tenant) ASSERT t.name IS UNIQUE;
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

:exit;
