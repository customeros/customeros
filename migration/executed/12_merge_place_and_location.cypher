https://github.com/openline-ai/openline-customer-os/issues/870

# STEP 1 QUERY PER TENANT. Execute before and after release, link Location with tenant

:param {tenant:"test"};

MATCH (t:Tenant {name:$tenant})--(:Contact)--(loc:Location)
MERGE (loc)-[:LOCATION_BELONGS_TO_TENANT]->(t);

MATCH (t:Tenant {name:$tenant})--(:Organization)--(loc:Location)
MERGE (loc)-[:LOCATION_BELONGS_TO_TENANT]->(t);

# STEP 2 Execute before and after release, link Location with tenant

MATCH (loc:Location)-[:LOCATED_AT]->(p:Place)
WHERE loc.country is null
WITH loc, p
SET
loc.country = p.country,
loc.region = p.state,
loc.locality = p.city,
loc.address = p.address,
loc.address2 = p.address2,
loc.zip = p.zip,
loc.phone = p.phone;

# STEP 3 Cleanup query, Execute after release

MATCH (alt:AlternatePlace) DETACH DELETE alt;
MATCH (p:Place) DETACH DELETE p;