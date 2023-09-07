# Step1 Remove non customer stages
match (or:OrganizationRelationship)--(os:OrganizationRelationshipStage) where or.name <> 'Customer' detach delete os;

# Step2 Create new non customer stages
WITH [{name:"Active",order:110},
      {name:"Inactive",order:120}] AS stages
UNWIND stages AS stage
MATCH (t:Tenant), (or:OrganizationRelationship)
WHERE or.name <> 'Customer'
MERGE (t)<-[:STAGE_BELONGS_TO_TENANT]-(s:OrganizationRelationshipStage {name:stage.name})<-[:HAS_STAGE]-(or)
				ON CREATE SET 	s.id=randomUUID(),
								s.order=stage.order,
								s.createdAt=datetime({timezone: 'UTC'});


# Step 3 delete deprecated stages
MATCH (t:Tenant)<-[:STAGE_BELONGS_TO_TENANT]-(s:OrganizationRelationshipStage)
WHERE s.name IN ['Target','Prospect','Unqualified'] detach delete s;

# Step 4 create new stages
WITH [{name:"MQL",order:20},
      {name:"SQL",order:30},
      {name:"Proposal",order:50},
      {name:"Not a fit",order:50}
      ] AS stages
UNWIND stages AS stage
MATCH (t:Tenant), (or:OrganizationRelationship)
WHERE or.name = 'Customer'
MERGE (t)<-[:STAGE_BELONGS_TO_TENANT]-(s:OrganizationRelationshipStage {name:stage.name})<-[:HAS_STAGE]-(or)
				ON CREATE SET 	s.id=randomUUID(),
								s.order=stage.order,
								s.createdAt=datetime({timezone: 'UTC'});

# Step 5 update existing orders
MATCH (s:OrganizationRelationshipStage) WHERE s.name='Lead' SET s.order=10;
MATCH (s:OrganizationRelationshipStage) WHERE s.name='MQL' SET s.order=20;
MATCH (s:OrganizationRelationshipStage) WHERE s.name='SQL' SET s.order=30;
MATCH (s:OrganizationRelationshipStage) WHERE s.name='Trial' SET s.order=40;
MATCH (s:OrganizationRelationshipStage) WHERE s.name='Proposal' SET s.order=50;
MATCH (s:OrganizationRelationshipStage) WHERE s.name='Live' SET s.order=60;
MATCH (s:OrganizationRelationshipStage) WHERE s.name='Lost' SET s.order=70;
MATCH (s:OrganizationRelationshipStage) WHERE s.name='Former' SET s.order=80;
MATCH (s:OrganizationRelationshipStage) WHERE s.name='Not a fit' SET s.order=90;

# Step 6 correct labels for new stages
MATCH (t:Tenant)<-[:STAGE_BELONGS_TO_TENANT]-(s:OrganizationRelationshipStage)
call  apoc.create.setLabels( s, [ "OrganizationRelationshipStage","OrganizationRelationshipStage_"+t.name ] )
YIELD node
return count(node);