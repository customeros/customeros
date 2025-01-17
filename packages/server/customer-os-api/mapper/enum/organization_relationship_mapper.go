package enummapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graphql/model"
)

var relationshipByModel = map[model.OrganizationRelationship]neo4jenum.OrganizationRelationship{
	model.OrganizationRelationshipCustomer:       neo4jenum.OrganizationRelationshipCustomer,
	model.OrganizationRelationshipProspect:       neo4jenum.OrganizationRelationshipProspect,
	model.OrganizationRelationshipNotAFit:        neo4jenum.OrganizationRelationshipNotAFit,
	model.OrganizationRelationshipFormerCustomer: neo4jenum.OrganizationRelationshipFormerCustomer,
}

var relationshipByValue = utils.ReverseMap(relationshipByModel)

func MapRelationshipFromModel(input model.OrganizationRelationship) neo4jenum.OrganizationRelationship {
	return relationshipByModel[input]
}

func MapRelationshipToModel(input neo4jenum.OrganizationRelationship) model.OrganizationRelationship {
	return relationshipByValue[input]
}
