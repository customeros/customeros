package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
)

var orgRelationshipByModel = map[model.OrganizationRelationship]entity.OrganizationRelationship{
	model.OrganizationRelationshipCustomer:    entity.Customer,
	model.OrganizationRelationshipDistributor: entity.Distributor,
	model.OrganizationRelationshipInvestor:    entity.Investor,
	model.OrganizationRelationshipPartner:     entity.Partner,
	model.OrganizationRelationshipSupplier:    entity.Supplier,
}

var orgRelationshipByValue = utils.ReverseMap(orgRelationshipByModel)

func MapOrgRelationshipFromModel(input model.OrganizationRelationship) entity.OrganizationRelationship {
	if v, exists := orgRelationshipByModel[input]; exists {
		return v
	} else {
		return ""
	}
}

func MapOrgRelationshipToModel(input entity.OrganizationRelationship) model.OrganizationRelationship {
	if v, exists := orgRelationshipByValue[input]; exists {
		return v
	} else {
		return ""
	}
}
