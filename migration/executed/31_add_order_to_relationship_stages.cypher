
MATCH (ors:OrganizationRelationshipStage {name: 'Target'}) SET ors.order=10;
MATCH (ors:OrganizationRelationshipStage {name: 'Lead'}) SET ors.order=20;
MATCH (ors:OrganizationRelationshipStage {name: 'Prospect'}) SET ors.order=30;
MATCH (ors:OrganizationRelationshipStage {name: 'Trial'}) SET ors.order=40;
MATCH (ors:OrganizationRelationshipStage {name: 'Lost'}) SET ors.order=50;
MATCH (ors:OrganizationRelationshipStage {name: 'Live'}) SET ors.order=60;
MATCH (ors:OrganizationRelationshipStage {name: 'Former'}) SET ors.order=70;
