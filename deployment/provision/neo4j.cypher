MERGE(t:Tenant {name: "customerosai"}) ON CREATE SET t.createdAt=datetime({timezone: 'UTC'}), t.id=randomUUID();
MATCH (t:Tenant {name:"customerosai"})
MERGE (t)-[:HAS_WORKSPACE]->(w:Workspace {name:"customeros.ai", provider: "google", appSource: "manual"});

MATCH (t:Tenant {name:"customerosai"})
MERGE (t)-[:HAS_SETTINGS]->(s:TenantSettings {id:randomUUID()})
				ON CREATE SET
					s.createdAt=datetime({timezone: 'UTC'}),
					s.invoicingEnabled=false,
					s.tenant="customerosai",
					s.defaultCurrency="EUR";

MATCH (t:Tenant {name:"customerosai"})
MERGE (t)<-[:EMAIL_ADDRESS_BELONGS_TO_TENANT]-(e:Email {rawEmail:"edi@customeros.ai"})
ON CREATE SET
            e.email="edi@customeros.ai",
            e.id=randomUUID(),
            e.source="openline",
            e.sourceOfTruth="openline",
            e.appSource="manual",
            e.createdAt=datetime({timezone: 'UTC'}),
            e.updatedAt=datetime({timezone: 'UTC'})
WITH t, e
MERGE (t)<-[:USER_BELONGS_TO_TENANT]-(u:User)-[rel:HAS]->(e)
ON CREATE SET
            u.id=randomUUID(),
            u.firstName="Eduard",
            u.lastName="Firut",
            u.roles=["USER", "OWNER"],
    		u.createdAt=datetime({timezone: 'UTC'}),
    		u.updatedAt=datetime({timezone: 'UTC'}),
    		u.source="openline",
    		u.sourceOfTruth="openline",
    		u.appSource="manual",
            rel.primary=true,
            rel.label="WORK";

MATCH (t:Tenant {name:"customerosai"})
MERGE (t)<-[:EMAIL_ADDRESS_BELONGS_TO_TENANT]-(e:Email {rawEmail:"alex@customeros.ai"})
ON CREATE SET
            e.email="alex@customeros.ai",
            e.id=randomUUID(),
            e.source="openline",
            e.sourceOfTruth="openline",
            e.appSource="manual",
            e.createdAt=datetime({timezone: 'UTC'}),
            e.updatedAt=datetime({timezone: 'UTC'})
WITH t, e
MERGE (t)<-[:USER_BELONGS_TO_TENANT]-(u:User)-[rel:HAS]->(e)
ON CREATE SET
            rel.primary=true,
            rel.label="WORK",
            u.id=randomUUID(),
            u.firstName="Alex",
            u.lastName="Basarab",
            u.roles=["USER", "OWNER"],
            u.createdAt=datetime({timezone: 'UTC'}),
            u.updatedAt=datetime({timezone: 'UTC'}),
            u.source="openline",
            u.sourceOfTruth="openline",
            u.appSource="manual";

MATCH (t:Tenant {name:"customerosai"})
MERGE (t)<-[:EMAIL_ADDRESS_BELONGS_TO_TENANT]-(e:Email {rawEmail:"acalinica@customeros.ai"})
ON CREATE SET
            e.email="acalinica@customeros.ai",
            e.id=randomUUID(),
            e.source="openline",
            e.sourceOfTruth="openline",
            e.appSource="manual",
            e.createdAt=datetime({timezone: 'UTC'}),
            e.updatedAt=datetime({timezone: 'UTC'})
WITH t, e
MERGE (t)<-[:USER_BELONGS_TO_TENANT]-(u:User)-[rel:HAS]->(e)
ON CREATE SET
            rel.primary=true,
            rel.label="WORK",
            u.id=randomUUID(),
            u.firstName="Alex",
            u.lastName="Calinica",
            u.roles=["USER","OWNER"],
    		u.createdAt=datetime({timezone: 'UTC'}),
    		u.updatedAt=datetime({timezone: 'UTC'}),
    		u.source="openline",
    		u.sourceOfTruth="openline",
    		u.appSource="manual";

MATCH (t:Tenant {name:"customerosai"})
MERGE (t)<-[:EMAIL_ADDRESS_BELONGS_TO_TENANT]-(e:Email {rawEmail:"silviu@customeros.ai"})
ON CREATE SET
            e.email="silviu@customeros.ai",
            e.id=randomUUID(),
            e.source="openline",
            e.sourceOfTruth="openline",
            e.appSource="manual",
            e.createdAt=datetime({timezone: 'UTC'}),
            e.updatedAt=datetime({timezone: 'UTC'})
WITH t, e
MERGE (t)<-[:USER_BELONGS_TO_TENANT]-(u:User)-[rel:HAS]->(e)
ON CREATE SET
            rel.primary=true,
            rel.label="WORK",
            u.id=randomUUID(),
            u.firstName="Silviu",
            u.lastName="Basu",
            u.roles=["USER","OWNER"],
    		u.createdAt=datetime({timezone: 'UTC'}),
    		u.updatedAt=datetime({timezone: 'UTC'}),
    		u.source="openline",
    		u.sourceOfTruth="openline",
    		u.appSource="manual";

MATCH (t:Tenant {name:"customerosai"})
MERGE (t)<-[:EMAIL_ADDRESS_BELONGS_TO_TENANT]-(e:Email {rawEmail:"kasia@customeros.ai"})
ON CREATE SET
            e.email="kasia@customeros.ai",
            e.id=randomUUID(),
            e.source="openline",
            e.sourceOfTruth="openline",
            e.appSource="manual",
            e.createdAt=datetime({timezone: 'UTC'}),
            e.updatedAt=datetime({timezone: 'UTC'})
WITH t, e
MERGE (t)<-[:USER_BELONGS_TO_TENANT]-(u:User)-[rel:HAS]->(e)
ON CREATE SET
            rel.primary=true,
            rel.label="WORK",
            u.id=randomUUID(),
            u.firstName="Kasia",
            u.lastName="Marciniszyn",
            u.roles=["USER","OWNER"],
    		u.createdAt=datetime({timezone: 'UTC'}),
    		u.updatedAt=datetime({timezone: 'UTC'}),
    		u.source="openline",
    		u.sourceOfTruth="openline",
    		u.appSource="manual";

MATCH (t:Tenant {name:"customerosai"})
MERGE (t)<-[:EMAIL_ADDRESS_BELONGS_TO_TENANT]-(e:Email {rawEmail:"mihai@customeros.ai"})
ON CREATE SET
            e.email="mihai@customeros.ai",
            e.id=randomUUID(),
            e.source="openline",
            e.sourceOfTruth="openline",
            e.appSource="manual",
            e.createdAt=datetime({timezone: 'UTC'}),
            e.updatedAt=datetime({timezone: 'UTC'})
WITH t, e
MERGE (t)<-[:USER_BELONGS_TO_TENANT]-(u:User)-[rel:HAS]->(e)
ON CREATE SET
            rel.primary=true,
            rel.label="WORK",
            u.id=randomUUID(),
            u.firstName="Mihai",
            u.lastName="Mihai",
            u.roles=["USER","OWNER"],
    		u.createdAt=datetime({timezone: 'UTC'}),
    		u.updatedAt=datetime({timezone: 'UTC'}),
    		u.source="openline",
    		u.sourceOfTruth="openline",
    		u.appSource="manual";

MATCH (t:Tenant {name:"customerosai"})<-[:TAG_BELONGS_TO_TENANT]-(tag:Tag) SET tag:Tag_customerosai;
MATCH (t:Tenant {name:"customerosai"})<-[:ORGANIZATION_TYPE_BELONGS_TO_TENANT]-(ot:OrganizationType) SET ot:OrganizationType_customerosai;
MATCH (t:Tenant {name:"customerosai"})<-[:USER_BELONGS_TO_TENANT]-(u:User) SET u:User_customerosai;
MATCH (t:Tenant {name:"customerosai"})<-[:CONTACT_BELONGS_TO_TENANT]-(c:Contact) SET c:Contact_customerosai;
MATCH (t:Tenant {name:"customerosai"})<-[:EMAIL_ADDRESS_BELONGS_TO_TENANT]-(e:Email) SET e:Email_customerosai;

MATCH (t:Tenant {name:"customerosai"})
MERGE (e:ExternalSystem {id:"calcom"})-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]->(t);
MATCH (t:Tenant {name:"customerosai"})
MERGE (e:ExternalSystem {id:"slack", name: "slack"})-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]->(t);
MATCH (t:Tenant {name:"customerosai"})
MERGE (e:ExternalSystem {id:"intercom", name: "intercom"})-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]->(t);

CREATE CONSTRAINT tenant_name_unique IF NOT EXISTS FOR (t:Tenant) REQUIRE t.name IS UNIQUE;
CREATE CONSTRAINT domain_domain_unique IF NOT EXISTS FOR (n:Domain) REQUIRE n.domain IS UNIQUE;

CREATE INDEX user_id_idx IF NOT EXISTS FOR (n:User) ON (n.id);
CREATE INDEX contact_id_idx IF NOT EXISTS FOR (n:Contact) ON (n.id);
CREATE INDEX tag_id_idx IF NOT EXISTS FOR (n:Tag) ON (n.id);
CREATE INDEX organization_id_idx IF NOT EXISTS FOR (n:Organization) ON (n.id);
CREATE INDEX custom_field_id_idx IF NOT EXISTS FOR (n:CustomField) ON (n.id);
CREATE INDEX field_set_id_idx IF NOT EXISTS FOR (n:FieldSet) ON (n.id);
CREATE INDEX email_id_idx IF NOT EXISTS FOR (n:Email) ON (n.id);
CREATE INDEX email_email_idx IF NOT EXISTS FOR (n:Email) ON (n.email);
CREATE INDEX phone_id_idx IF NOT EXISTS FOR (n:PhoneNumber) ON (n.id);
CREATE INDEX phone_e164_idx IF NOT EXISTS FOR (n:PhoneNumber) ON (n.e164);
CREATE INDEX action_id_idx IF NOT EXISTS FOR (n:Action) ON (n.id);
CREATE INDEX interaction_session_id_idx IF NOT EXISTS FOR (n:InteractionSession) ON (n.id);
CREATE INDEX interaction_event_id_idx IF NOT EXISTS FOR (n:InteractionEvent) ON (n.id);
CREATE INDEX note_id_idx IF NOT EXISTS FOR (n:Note) ON (n.id);
CREATE INDEX job_role_id_idx IF NOT EXISTS FOR (n:JobRole) ON (n.id);
CREATE INDEX location_id_idx IF NOT EXISTS FOR (n:Location) ON (n.id);
CREATE INDEX log_entry_id_idx IF NOT EXISTS FOR (n:LogEntry) ON (n.id);
CREATE INDEX comment_id_idx IF NOT EXISTS FOR (n:Comment) ON (n.id);
CREATE INDEX issue_id_idx IF NOT EXISTS FOR (n:Issue) ON (n.id);
CREATE INDEX meeting_id_idx IF NOT EXISTS FOR (n:Meeting) ON (n.id);
CREATE INDEX timeline_event_id_idx IF NOT EXISTS FOR (n:TimelineEvent) ON (n.id);
CREATE INDEX opportunity_id_idx IF NOT EXISTS FOR (n:Opportunity) ON (n.id);
CREATE INDEX contract_id_idx IF NOT EXISTS FOR (n:Contract) ON (n.id);
CREATE INDEX service_line_item_id_idx IF NOT EXISTS FOR (n:ServiceLineItem) ON (n.id);
CREATE INDEX invoice_id_idx IF NOT EXISTS FOR (n:Invoice) ON (n.id);
CREATE INDEX invoice_line_id_idx IF NOT EXISTS FOR (n:InvoiceLine) ON (n.id);

:exit;
