#========== migrate org domains to separate nodes and combine with email domains ==========

MATCH (e:EmailDomain) SET e:Domain;

MATCH (e:Domain) REMOVE e:EmailDomain;

MATCH (o:Organization)
WHERE o.domain is not null AND o.domain <> ''
WITH o, o.domain as domain
MERGE (d:Domain {domain:domain})
ON CREATE SET
    d.id = randomUUID(),
    d.source = o.source,
    d.appSource = o.appSource,
    d.createdAt = datetime({timezone: 'UTC'}),
    d.updatedAt = datetime({timezone: 'UTC'})
MERGE (o)-[:HAS_DOMAIN]->(d);

MATCH (e:Domain) REMOVE e.sourceOfTruth;

#========== Add unique constraint to domain nodes, include scripts in common file in openline-cloud ==========
CREATE CONSTRAINT domain_domain_unique IF NOT EXISTS ON (n:Domain) ASSERT n.domain IS UNIQUE;