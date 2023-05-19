MATCH (tmporg:TmpOrg)-[:HAS_LINK]->(tmplink:TmpLink)
MATCH (o:Organization {name:tmporg.child})-[:ORGANIZATION_BELONGS_TO_TENANT]->(:Tenant {name:'openline'})<-[:ENTITY_TEMPLATE_BELONGS_TO_TENANT]-(e:EntityTemplate {id:'59fe701a-14b5-4a61-bd11-3c5fcf0b13e7'})-[:CONTAINS]->(d:CustomFieldTemplate {id:'58f49b51-aeb6-4574-be39-3a80f321c73e'})
  WHERE e.extends='ORGANIZATION'
MERGE (o)-[:IS_DEFINED_BY]->(e)
MERGE (f:TextField:CustomField {name: 'Contract', datatype:'TEXT'})<-[:HAS_PROPERTY]-(o)
  ON CREATE SET f.textValue=tmplink.link, f.id=randomUUID(), f.appSource='CSV', f.source='openline', f.sourceOfTruth='openline', f.createdAt=datetime({timezone: 'UTC'}), f.updatedAt=datetime({timezone: 'UTC'}), f:CustomField_openline
  ON MATCH SET f.textValue=tmplink.link, f.sourceOfTruth='openline', f.updatedAt=datetime({timezone: 'UTC'})
MERGE (f)-[:IS_DEFINED_BY]->(d)
RETURN o.id