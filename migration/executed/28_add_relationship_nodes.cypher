MATCH (r:OrganizationRelationship) DETACH DELETE r;

CREATE CONSTRAINT organization_relationship_name_unique IF NOT EXISTS FOR (or:OrganizationRelationship) REQUIRE or.name IS UNIQUE;

MERGE (r:OrganizationRelationship {name:"Customer", group:"Sales and marketing"}) ON CREATE SET r.id=randomUUID(), r.createdAt=datetime({timezone: 'UTC'});
MERGE (r:OrganizationRelationship {name:"Distributor", group:"Sales and marketing"}) ON CREATE SET r.id=randomUUID(), r.createdAt=datetime({timezone: 'UTC'});
MERGE (r:OrganizationRelationship {name:"Partner", group:"Sales and marketing"}) ON CREATE SET r.id=randomUUID(), r.createdAt=datetime({timezone: 'UTC'});
MERGE (r:OrganizationRelationship {name:"Licensing partner", group:"Sales and marketing"}) ON CREATE SET r.id=randomUUID(), r.createdAt=datetime({timezone: 'UTC'});
MERGE (r:OrganizationRelationship {name:"Franchisee", group:"Sales and marketing"}) ON CREATE SET r.id=randomUUID(), r.createdAt=datetime({timezone: 'UTC'});
MERGE (r:OrganizationRelationship {name:"Franchisor", group:"Sales and marketing"}) ON CREATE SET r.id=randomUUID(), r.createdAt=datetime({timezone: 'UTC'});
MERGE (r:OrganizationRelationship {name:"Affiliate", group:"Sales and marketing"}) ON CREATE SET r.id=randomUUID(), r.createdAt=datetime({timezone: 'UTC'});
MERGE (r:OrganizationRelationship {name:"Reseller", group:"Sales and marketing"}) ON CREATE SET r.id=randomUUID(), r.createdAt=datetime({timezone: 'UTC'});
MERGE (r:OrganizationRelationship {name:"Influencer or content creator", group:"Sales and marketing"}) ON CREATE SET r.id=randomUUID(), r.createdAt=datetime({timezone: 'UTC'});
MERGE (r:OrganizationRelationship {name:"Media partner", group:"Sales and marketing"}) ON CREATE SET r.id=randomUUID(), r.createdAt=datetime({timezone: 'UTC'});

MERGE (r:OrganizationRelationship {name:"Investor", group:"Financial"}) ON CREATE SET r.id=randomUUID(), r.createdAt=datetime({timezone: 'UTC'});
MERGE (r:OrganizationRelationship {name:"Merger or acquisition target", group:"Financial"}) ON CREATE SET r.id=randomUUID(), r.createdAt=datetime({timezone: 'UTC'});
MERGE (r:OrganizationRelationship {name:"Parent company", group:"Financial"}) ON CREATE SET r.id=randomUUID(), r.createdAt=datetime({timezone: 'UTC'});
MERGE (r:OrganizationRelationship {name:"Subsidiary", group:"Financial"}) ON CREATE SET r.id=randomUUID(), r.createdAt=datetime({timezone: 'UTC'});
MERGE (r:OrganizationRelationship {name:"Joint venture", group:"Financial"}) ON CREATE SET r.id=randomUUID(), r.createdAt=datetime({timezone: 'UTC'});
MERGE (r:OrganizationRelationship {name:"Sponsor", group:"Financial"}) ON CREATE SET r.id=randomUUID(), r.createdAt=datetime({timezone: 'UTC'});

MERGE (r:OrganizationRelationship {name:"Supplier", group:"Supply chain"}) ON CREATE SET r.id=randomUUID(), r.createdAt=datetime({timezone: 'UTC'});
MERGE (r:OrganizationRelationship {name:"Vendor", group:"Supply chain"}) ON CREATE SET r.id=randomUUID(), r.createdAt=datetime({timezone: 'UTC'});
MERGE (r:OrganizationRelationship {name:"Contract manufacturer", group:"Supply chain"}) ON CREATE SET r.id=randomUUID(), r.createdAt=datetime({timezone: 'UTC'});
MERGE (r:OrganizationRelationship {name:"Original equipment manufacturer", group:"Supply chain"}) ON CREATE SET r.id=randomUUID(), r.createdAt=datetime({timezone: 'UTC'});
MERGE (r:OrganizationRelationship {name:"Original design manufacturer", group:"Supply chain"}) ON CREATE SET r.id=randomUUID(), r.createdAt=datetime({timezone: 'UTC'});
MERGE (r:OrganizationRelationship {name:"Private label manufacturer", group:"Supply chain"}) ON CREATE SET r.id=randomUUID(), r.createdAt=datetime({timezone: 'UTC'});
MERGE (r:OrganizationRelationship {name:"Logistics partner", group:"Supply chain"}) ON CREATE SET r.id=randomUUID(), r.createdAt=datetime({timezone: 'UTC'});

MERGE (r:OrganizationRelationship {name:"Consultant", group:"Service provider"}) ON CREATE SET r.id=randomUUID(), r.createdAt=datetime({timezone: 'UTC'});
MERGE (r:OrganizationRelationship {name:"Service provider", group:"Service provider"}) ON CREATE SET r.id=randomUUID(), r.createdAt=datetime({timezone: 'UTC'});
MERGE (r:OrganizationRelationship {name:"Outsourcing provider", group:"Service provider"}) ON CREATE SET r.id=randomUUID(), r.createdAt=datetime({timezone: 'UTC'});
MERGE (r:OrganizationRelationship {name:"Insourcing partner", group:"Service provider"}) ON CREATE SET r.id=randomUUID(), r.createdAt=datetime({timezone: 'UTC'});
MERGE (r:OrganizationRelationship {name:"Technology provider", group:"Service provider"}) ON CREATE SET r.id=randomUUID(), r.createdAt=datetime({timezone: 'UTC'});
MERGE (r:OrganizationRelationship {name:"Data provider", group:"Service provider"}) ON CREATE SET r.id=randomUUID(), r.createdAt=datetime({timezone: 'UTC'});
MERGE (r:OrganizationRelationship {name:"Certification body", group:"Service provider"}) ON CREATE SET r.id=randomUUID(), r.createdAt=datetime({timezone: 'UTC'});
MERGE (r:OrganizationRelationship {name:"Standards organization", group:"Service provider"}) ON CREATE SET r.id=randomUUID(), r.createdAt=datetime({timezone: 'UTC'});
MERGE (r:OrganizationRelationship {name:"Industry analyst", group:"Service provider"}) ON CREATE SET r.id=randomUUID(), r.createdAt=datetime({timezone: 'UTC'});
MERGE (r:OrganizationRelationship {name:"Real estate partner", group:"Service provider"}) ON CREATE SET r.id=randomUUID(), r.createdAt=datetime({timezone: 'UTC'});
MERGE (r:OrganizationRelationship {name:"Talent acquisition partner", group:"Service provider"}) ON CREATE SET r.id=randomUUID(), r.createdAt=datetime({timezone: 'UTC'});
MERGE (r:OrganizationRelationship {name:"Professional employer organization", group:"Service provider"}) ON CREATE SET r.id=randomUUID(), r.createdAt=datetime({timezone: 'UTC'});

MERGE (r:OrganizationRelationship {name:"Research collaborator", group:"Collaborative"}) ON CREATE SET r.id=randomUUID(), r.createdAt=datetime({timezone: 'UTC'});
MERGE (r:OrganizationRelationship {name:"Regulatory body", group:"Collaborative"}) ON CREATE SET r.id=randomUUID(), r.createdAt=datetime({timezone: 'UTC'});
MERGE (r:OrganizationRelationship {name:"Trade association member", group:"Collaborative"}) ON CREATE SET r.id=randomUUID(), r.createdAt=datetime({timezone: 'UTC'});

MERGE (r:OrganizationRelationship {name:"Competitor", group:"Competitive"}) ON CREATE SET r.id=randomUUID(), r.createdAt=datetime({timezone: 'UTC'});