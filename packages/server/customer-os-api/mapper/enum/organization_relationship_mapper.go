package enummapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
)

var relationshipByModel = map[model.OrganizationRelationship]neo4jenum.OrganizationRelationship{
	model.OrganizationRelationshipCustomer:       neo4jenum.Customer,
	model.OrganizationRelationshipProspect:       neo4jenum.Prospect,
	model.OrganizationRelationshipNotAFit:        neo4jenum.NotAFit,
	model.OrganizationRelationshipFormerCustomer: neo4jenum.FormerCustomer,
}

var relationshipByValue = utils.ReverseMap(relationshipByModel)

func MapRelationshipFromModel(input model.OrganizationRelationship) neo4jenum.OrganizationRelationship {
	return relationshipByModel[input]
}

func MapRelationshipToModel(input neo4jenum.OrganizationRelationship) model.OrganizationRelationship {
	return relationshipByValue[input]
}
