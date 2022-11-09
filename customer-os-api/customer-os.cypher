CREATE CONSTRAINT tenant_name_unique IF NOT EXISTS ON (t:Tenant) ASSERT t.name IS UNIQUE;
MERGE(t:Tenant {name: "openline"});

CREATE INDEX contact_id_idx IF NOT EXISTS FOR (n:Contact) ON (n.id);
CREATE INDEX contact_group_id_idx IF NOT EXISTS FOR (n:ContactGroup) ON (n.id);
CREATE INDEX text_custom_field_id_idx IF NOT EXISTS FOR (n:TextCustomField) ON (n.id);
CREATE INDEX field_set_id_idx IF NOT EXISTS FOR (n:FieldSet) ON (n.id);
CREATE INDEX email_id_idx IF NOT EXISTS FOR (n:Email) ON (n.id);
CREATE INDEX email_email_idx IF NOT EXISTS FOR (n:Email) ON (n.email);
CREATE INDEX phone_id_idx IF NOT EXISTS FOR (n:PhoneNumber) ON (n.id);
CREATE INDEX phone_e164_idx IF NOT EXISTS FOR (n:PhoneNumber) ON (n.e164);
CREATE INDEX entity_definition_id_idx IF NOT EXISTS FOR (n:EntityDefinition) ON (n.id);