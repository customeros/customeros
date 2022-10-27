CREATE INDEX ON :Contact(id);

CREATE INDEX ON :ContactGroup(id);
CREATE CONSTRAINT ON (g:ContactGroup) ASSERT g.name IS UNIQUE;