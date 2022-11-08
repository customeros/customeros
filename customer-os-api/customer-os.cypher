CREATE CONSTRAINT tenant_name_unique IF NOT EXISTS ON (t:Tenant) ASSERT t.name IS UNIQUE;
CREATE(t:Tenant {name: "openline"});

CREATE INDEX ON :Contact(id);
CREATE INDEX ON :ContactGroup(id);
CREATE INDEX ON :TextCustomField(id);
CREATE INDEX ON :FieldSet(id);
CREATE INDEX ON :Email(id);
CREATE INDEX ON :Email(email);
CREATE INDEX ON :PhoneNumber(id);
CREATE INDEX ON :PhoneNumber(number);
