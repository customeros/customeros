# Execute before release

WITH [
{name: 'Green', order:10},
{name: 'Yellow', order:20},
{name: 'Orange', order:30},
{name: 'Red', order:40}
] AS indicators
UNWIND indicators AS indicator
MATCH (t:Tenant)
MERGE (t)<-[:HEALTH_INDICATOR_BELONGS_TO_TENANT]-(h:HealthIndicator {name: indicator.name})
ON CREATE SET h.id=randomUUID(), h.createdAt=datetime({timezone: 'UTC'}), h.order=indicator.order;