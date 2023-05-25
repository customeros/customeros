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
  MERGE (t)<-[:TAG_BELONGS_TO_TENANT]-(tag:Tag {name:"CUSTOMER"})
  ON CREATE SET tag.id=randomUUID(),
                tag.createdAt=datetime({timezone: 'UTC'}),
                tag.updatedAt=datetime({timezone: 'UTC'}),
                tag.source="openline",
                tag.appSource="manual";

MATCH (t:Tenant {name:"openline"})
  MERGE (t)<-[:TAG_BELONGS_TO_TENANT]-(tag:Tag {name:"SUPPLIER"})
  ON CREATE SET tag.id=randomUUID(),
                tag.createdAt=datetime({timezone: 'UTC'}),
                tag.updatedAt=datetime({timezone: 'UTC'}),
                tag.source="openline",
                tag.appSource="manual";

MATCH (t:Tenant {name:"openline"})
  MERGE (t)<-[:TAG_BELONGS_TO_TENANT]-(tag:Tag {name:"INVESTOR"})
  ON CREATE SET tag.id=randomUUID(),
                tag.createdAt=datetime({timezone: 'UTC'}),
                tag.updatedAt=datetime({timezone: 'UTC'}),
                tag.source="openline",
                tag.appSource="manual";

MATCH (t:Tenant {name:"openline"})
  MERGE (t)<-[:TAG_BELONGS_TO_TENANT]-(tag:Tag {name:"PROSPECT"})
  ON CREATE SET tag.id=randomUUID(),
                tag.createdAt=datetime({timezone: 'UTC'}),
                tag.updatedAt=datetime({timezone: 'UTC'}),
                tag.source="openline",
                tag.appSource="manual";

MATCH (t:Tenant {name:"openline"})
  MERGE (t)<-[:ORGANIZATION_TYPE_BELONGS_TO_TENANT]-(ot:OrganizationType {name:"COMPANY"})
  ON CREATE SET ot.id=randomUUID(), ot.createdAt=datetime({timezone: 'UTC'});

MATCH (t:Tenant {name:"openline"})
    MERGE (u:User {id:"development@openline.ai"})-[:USER_BELONGS_TO_TENANT]->(t)
    ON CREATE SET
    		u.firstName="Development",
            u.lastName="User",
    		u.createdAt=datetime({timezone: 'UTC'}),
    		u.source="openline",
    		u.sourceOfTruth="openline",
    		u.appSource="manual";

MATCH (t:Tenant {name:"openline"})
MERGE (t)<-[:EMAIL_ADDRESS_BELONGS_TO_TENANT]-(e:Email {rawEmail:"development@openline.ai"})
ON CREATE SET
            e.email='development@openline.ai',
            e.id=randomUUID(),
            e.source="openline",
            e.sourceOfTruth="openline",
            e.appSource="manual",
            e.createdAt=datetime({timezone: 'UTC'}),
            e.updatedAt=datetime({timezone: 'UTC'})
WITH t, e
MATCH (u:User {id:"development@openline.ai"})-[:USER_BELONGS_TO_TENANT]->(t)
MERGE (u)-[rel:HAS]->(e)
ON CREATE SET
            rel.primary=true,
            rel.label="WORK";

MATCH (t:Tenant {name:"openline"})
MERGE (t)<-[:EMAIL_ADDRESS_BELONGS_TO_TENANT]-(e:Email {rawEmail:"edi@openline.ai"})
ON CREATE SET
            e.email="edi@openline.ai",
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
            rel.label="WORK"
MERGE (p:Person {email: "edi@openline.ai", provider: "google"})-[:IDENTIFIES {default: true}]->(u)
ON CREATE SET
        p.id=randomUUID(),
        p.createdAt=datetime({timezone: 'UTC'}),
        p.updatedAt=datetime({timezone: 'UTC'}),
        p.source="openline",
        p.sourceOfTruth="openline",
        p.appSource="manual";

MATCH (t:Tenant {name:"openline"})
MERGE (t)<-[:EMAIL_ADDRESS_BELONGS_TO_TENANT]-(e:Email {rawEmail:"alex@openline.ai"})
ON CREATE SET
            e.email="alex@openline.ai",
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
            u.appSource="manual"
MERGE (p:Person {email: "alex@openline.ai", provider: "google"})-[:IDENTIFIES {default: true}]->(u)
ON CREATE SET
        p.id=randomUUID(),
        p.createdAt=datetime({timezone: 'UTC'}),
        p.updatedAt=datetime({timezone: 'UTC'}),
        p.source="openline",
        p.sourceOfTruth="openline",
        p.appSource="manual";


MATCH (t:Tenant {name:"openline"})
MERGE (t)<-[:EMAIL_ADDRESS_BELONGS_TO_TENANT]-(e:Email {rawEmail:"kasia@openline.ai"})
ON CREATE SET
            e.email="kasia@openline.ai",
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
    		u.appSource="manual"
MERGE (p:Person {email: "kasia@openline.ai", provider: "google"})-[:IDENTIFIES {default: true}]->(u)
ON CREATE SET
        p.id=randomUUID(),
        p.identityId="b7aeff67-ca86-4f68-8344-37748ae792fe",
        p.createdAt=datetime({timezone: 'UTC'}),
        p.updatedAt=datetime({timezone: 'UTC'}),
        p.source="openline",
        p.sourceOfTruth="openline",
        p.appSource="manual";

MATCH (t:Tenant {name:"openline"})
MERGE (t)<-[:EMAIL_ADDRESS_BELONGS_TO_TENANT]-(e:Email {rawEmail:"gabi@openline.ai"})
ON CREATE SET
            e.email="gabi@openline.ai",
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
            u.firstName="Gabriel",
            u.lastName="Gontariu",
            u.roles=["USER","OWNER"],
    		u.createdAt=datetime({timezone: 'UTC'}),
    		u.updatedAt=datetime({timezone: 'UTC'}),
    		u.source="openline",
    		u.sourceOfTruth="openline",
    		u.appSource="manual"
MERGE (p:Person {email: "gabi@openline.ai", provider: "google"})-[:IDENTIFIES {default: true}]->(u)
ON CREATE SET
        p.id=randomUUID(),
        p.identityId="8327e04a-877b-4b05-8aaa-ef6a582f7836",
        p.createdAt=datetime({timezone: 'UTC'}),
        p.updatedAt=datetime({timezone: 'UTC'}),
        p.source="openline",
        p.sourceOfTruth="openline",
        p.appSource="manual";

MATCH (t:Tenant {name:"openline"})
MERGE (t)<-[:EMAIL_ADDRESS_BELONGS_TO_TENANT]-(e:Email {rawEmail:"torrey@openline.ai"})
ON CREATE SET
            e.email="torrey@openline.ai",
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
            u.firstName="Torrey",
            u.lastName="Searle",
            u.roles=["OWNER","USER"],
    		u.createdAt=datetime({timezone: 'UTC'}),
    		u.updatedAt=datetime({timezone: 'UTC'}),
    		u.source="openline",
    		u.sourceOfTruth="openline",
    		u.appSource="manual"
MERGE (p:Person {email: "torrey@openline.ai", provider: "google"})-[:IDENTIFIES {default: true}]->(u)
ON CREATE SET
        p.id=randomUUID(),
        p.identityId="46a69d24-e15a-4a04-ae44-067186ab1c87",
        p.createdAt=datetime({timezone: 'UTC'}),
        p.updatedAt=datetime({timezone: 'UTC'}),
        p.source="openline",
        p.sourceOfTruth="openline",
        p.appSource="manual";


MATCH (t:Tenant {name:"openline"})
MERGE (t)<-[:EMAIL_ADDRESS_BELONGS_TO_TENANT]-(e:Email {rawEmail:"vasi@openline.ai"})
ON CREATE SET
            e.email="vasi@openline.ai",
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
            u.firstName="Vasi",
            u.lastName="Coscotin",
        u.roles=["OWNER","USER"],
    		u.createdAt=datetime({timezone: 'UTC'}),
    		u.updatedAt=datetime({timezone: 'UTC'}),
    		u.source="openline",
    		u.sourceOfTruth="openline",
    		u.appSource="manual"
MERGE (p:Person {email: "vasi@openline.ai", provider: "google"})-[:IDENTIFIES {default: true}]->(u)
ON CREATE SET
        p.id=randomUUID(),
        p.identityId="c6591b09-4e2a-48ba-bff2-a30c33e26a3a",
        p.createdAt=datetime({timezone: 'UTC'}),
        p.updatedAt=datetime({timezone: 'UTC'}),
        p.source="openline",
        p.sourceOfTruth="openline",
        p.appSource="manual";



MATCH (t:Tenant {name:"openline"})
MERGE (t)<-[:EMAIL_ADDRESS_BELONGS_TO_TENANT]-(e:Email {rawEmail:"antoine@openline.ai"})
ON CREATE SET
            e.email="antoine@openline.ai",
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
            u.firstName="Antoine",
            u.lastName="Valot",
            u.roles=["USER"],
    		u.createdAt=datetime({timezone: 'UTC'}),
    		u.updatedAt=datetime({timezone: 'UTC'}),
    		u.source="openline",
    		u.sourceOfTruth="openline",
    		u.appSource="manual"
MERGE (p:Person {email: "antoine@openline.ai", provider: "google"})-[:IDENTIFIES {default: true}]->(u)
ON CREATE SET
        p.id=randomUUID(),
        p.createdAt=datetime({timezone: 'UTC'}),
        p.updatedAt=datetime({timezone: 'UTC'}),
        p.source="openline",
        p.sourceOfTruth="openline",
        p.appSource="manual";

MATCH (t:Tenant {name:"openline"})
MERGE (t)<-[:EMAIL_ADDRESS_BELONGS_TO_TENANT]-(e:Email {rawEmail:"matt@openline.ai"})
ON CREATE SET
            e.email="matt@openline.ai",
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
            u.firstName="Matt",
            u.lastName="Brown",
            u.roles=["USER"],
    		u.createdAt=datetime({timezone: 'UTC'}),
    		u.updatedAt=datetime({timezone: 'UTC'}),
    		u.source="openline",
    		u.sourceOfTruth="openline",
    		u.appSource="manual"
MERGE (p:Person {email: "matt@openline.ai", provider: "google"})-[:IDENTIFIES {default: true}]->(u)
ON CREATE SET
        p.id=randomUUID(),
        p.createdAt=datetime({timezone: 'UTC'}),
        p.updatedAt=datetime({timezone: 'UTC'}),
        p.source="openline",
        p.sourceOfTruth="openline",
        p.appSource="manual";


MATCH (t:Tenant {name:"openline"})
MERGE (t)<-[:EMAIL_ADDRESS_BELONGS_TO_TENANT]-(e:Email {rawEmail:"jonty@openline.ai"})
ON CREATE SET
            e.email="jonty@openline.ai",
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
            u.firstName="Jonty",
            u.lastName="Knox",
            u.roles=["USER"],
    		u.createdAt=datetime({timezone: 'UTC'}),
    		u.updatedAt=datetime({timezone: 'UTC'}),
    		u.source="openline",
    		u.sourceOfTruth="openline",
    		u.appSource="manual"
MERGE (p:Person {email: "jonty@openline.ai", provider: "google"})-[:IDENTIFIES {default: true}]->(u)
ON CREATE SET
        p.id=randomUUID(),
        p.createdAt=datetime({timezone: 'UTC'}),
        p.updatedAt=datetime({timezone: 'UTC'}),
        p.source="openline",
        p.sourceOfTruth="openline",
        p.appSource="manual";


MATCH (t:Tenant {name:"openline"})
    MERGE (c:Contact {id:"echotest"})-[:CONTACT_BELONGS_TO_TENANT]->(t)
    ON CREATE SET
    		c.firstName ="Echo",
            c.lastName="Test",
    		c.createdAt=datetime({timezone: 'UTC'}),
    		c.updatedAt=datetime({timezone: 'UTC'}),
    		c.source="openline",
            c.sourceOfTruth="openline",
            c.appSource="manual";

MATCH (t:Tenant {name:"openline"})
MERGE (t)<-[:EMAIL_ADDRESS_BELONGS_TO_TENANT]-(e:Email {rawEmail:"echo@oasis.openline.ai"})
ON CREATE SET
            e.id=randomUUID(),
            e.source="openline",
            e.sourceOfTruth="openline",
            e.appSource="manual",
            e.createdAt=datetime({timezone: 'UTC'}),
            e.updatedAt=datetime({timezone: 'UTC'})
WITH t, e
MATCH (c:Contact {id:"echotest"})-[:CONTACT_BELONGS_TO_TENANT]->(t)
MERGE (c)-[rel:HAS]->(e)
ON CREATE SET
            rel.primary=true,
            rel.label="MAIN";

MATCH (t:Tenant {name:"openline"})<-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]-(e:ExternalSystem) SET e:ExternalSystem_openline;
MATCH (t:Tenant {name:"openline"})<-[:TAG_BELONGS_TO_TENANT]-(tag:Tag) SET tag:Tag_openline;
MATCH (t:Tenant {name:"openline"})<-[:ORGANIZATION_TYPE_BELONGS_TO_TENANT]-(ot:OrganizationType) SET ot:OrganizationType_openline;
MATCH (t:Tenant {name:"openline"})<-[:USER_BELONGS_TO_TENANT]-(u:User) SET u:User_openline;
MATCH (t:Tenant {name:"openline"})<-[:CONTACT_BELONGS_TO_TENANT]-(c:Contact) SET c:Contact_openline;
MATCH (t:Tenant {name:"openline"})<-[:EMAIL_ADDRESS_BELONGS_TO_TENANT]-(e:Email) SET e:Email_openline;

MERGE (c:Country {name:"Romania"}) ON CREATE SET
c.id=randomUUID(),
c.name="Romania",
c.phoneCode="40",
c.codeA2="RO",
c.codeA3="ROU",
c.appSource="csvImport",
c.createdAt=datetime({timezone: 'UTC'}),
c.source="openline",
c.sourceOfTruth= "openline",
c.updatedAt=datetime({timezone: 'UTC'});

DROP INDEX basicSearchStandard_openline IF EXISTS;
CREATE FULLTEXT INDEX basicSearchStandard_openline FOR (n:Contact_openline|Email_openline|Organization_openline) ON EACH [n.firstName, n.lastName, n.name, n.email]
OPTIONS {
  indexConfig: {
    `fulltext.analyzer`: 'standard',
    `fulltext.eventually_consistent`: true
  }
};

DROP INDEX basicSearchSimple_openline IF EXISTS;
CREATE FULLTEXT INDEX basicSearchSimple_openline FOR (n:Contact_openline|Email_openline|Organization_openline) ON EACH [n.firstName, n.lastName, n.email, n.name]
OPTIONS {
  indexConfig: {
    `fulltext.analyzer`: 'simple',
    `fulltext.eventually_consistent`: true
  }
};

CREATE CONSTRAINT tenant_name_unique IF NOT EXISTS FOR (t:Tenant) REQUIRE t.name IS UNIQUE;
CREATE CONSTRAINT domain_domain_unique IF NOT EXISTS FOR (n:Domain) REQUIRE n.domain IS UNIQUE;

CREATE INDEX contact_id_idx IF NOT EXISTS FOR (n:Contact) ON (n.id);
CREATE INDEX tag_id_idx IF NOT EXISTS FOR (n:Tag) ON (n.id);
CREATE INDEX organization_id_idx IF NOT EXISTS FOR (n:Organization) ON (n.id);
CREATE INDEX custom_field_id_idx IF NOT EXISTS FOR (n:CustomField) ON (n.id);
CREATE INDEX field_set_id_idx IF NOT EXISTS FOR (n:FieldSet) ON (n.id);
CREATE INDEX email_id_idx IF NOT EXISTS FOR (n:Email) ON (n.id);
CREATE INDEX email_email_idx IF NOT EXISTS FOR (n:Email) ON (n.email);
CREATE INDEX phone_id_idx IF NOT EXISTS FOR (n:PhoneNumber) ON (n.id);
CREATE INDEX phone_e164_idx IF NOT EXISTS FOR (n:PhoneNumber) ON (n.e164);

:exit;
