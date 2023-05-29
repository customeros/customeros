MERGE (r:OrganizationRelationship {name:"INVESTOR"}) ON CREATE SET r.id=randomUUID();
MERGE (r:OrganizationRelationship {name:"SUPPLIER"}) ON CREATE SET r.id=randomUUID();
MERGE (r:OrganizationRelationship {name:"PARTNER"}) ON CREATE SET r.id=randomUUID();
MERGE (r:OrganizationRelationship {name:"CUSTOMER"}) ON CREATE SET r.id=randomUUID();
MERGE (r:OrganizationRelationship {name:"DISTRIBUTOR"}) ON CREATE SET r.id=randomUUID();