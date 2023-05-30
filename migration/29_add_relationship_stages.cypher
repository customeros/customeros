# part 1:

WITH [
{name: 'Target'},
{name: 'Lead'},
{name: 'Prospect'},
{name: 'Trial'},
{name: 'Lost'},
{name: 'Live'},
{name: 'Former'}
] AS stages
UNWIND stages AS stage
MATCH (t:Tenant)
MERGE (t)<-[:STAGE_BELONGS_TO_TENANT]-(s:OrganizationRelationshipStage {name: stage.name})
ON CREATE SET s.id=randomUUID(), s.createdAt=datetime({timezone: 'UTC'})
WITH s
MATCH (or:OrganizationRelationship)
MERGE (s)<-[:HAS_STAGE]-(or)

# part 2: execute per tenant
MATCH (t:Tenant {name:"openline"})<-[:STAGE_BELONGS_TO_TENANT]-(s:OrganizationRelationshipStage)
SET s:OrganizationRelationshipStage_openline;