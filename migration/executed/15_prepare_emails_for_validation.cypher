https://github.com/openline-ai/openline-customer-os/issues/945

#========== migrate emails for non validated emails ==========

MATCH (e:Email)
WHERE e.acceptsMail is not null AND e.validated is null
SET e.validated=true;

MATCH (e:Email)
WHERE e.validated is null AND e.rawEmail is null AND e.email is not null
SET e.rawEmail=e.email
REMOVE e.email;

MATCH (u:User)-[:HAS]->(e:Email)
WHERE e.email is null and e.rawEmail is not null
SET e.email=e.rawEmail;