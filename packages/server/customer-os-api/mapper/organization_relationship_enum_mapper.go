package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
)

var orgRelationshipByModel = map[model.OrganizationRelationship]entity.OrganizationRelationship{
	model.OrganizationRelationshipCustomer:                         entity.Customer,
	model.OrganizationRelationshipDistributor:                      entity.Distributor,
	model.OrganizationRelationshipPartner:                          entity.Partner,
	model.OrganizationRelationshipLicensingPartner:                 entity.LicensingPartner,
	model.OrganizationRelationshipFranchisee:                       entity.Franchisee,
	model.OrganizationRelationshipFranchisor:                       entity.Franchisor,
	model.OrganizationRelationshipAffiliate:                        entity.Affiliate,
	model.OrganizationRelationshipReseller:                         entity.Reseller,
	model.OrganizationRelationshipInfluencerOrContentCreator:       entity.Influencer,
	model.OrganizationRelationshipMediaPartner:                     entity.MediaPartner,
	model.OrganizationRelationshipInvestor:                         entity.Investor,
	model.OrganizationRelationshipMergerOrAcquisitionTarget:        entity.Merger,
	model.OrganizationRelationshipParentCompany:                    entity.ParentCompany,
	model.OrganizationRelationshipSubsidiary:                       entity.Subsidiary,
	model.OrganizationRelationshipJointVenture:                     entity.JointVenture,
	model.OrganizationRelationshipSponsor:                          entity.Sponsor,
	model.OrganizationRelationshipSupplier:                         entity.Supplier,
	model.OrganizationRelationshipVendor:                           entity.Vendor,
	model.OrganizationRelationshipContractManufacturer:             entity.ContractManufacturer,
	model.OrganizationRelationshipOriginalEquipmentManufacturer:    entity.OriginalEquipmentManufacturer,
	model.OrganizationRelationshipOriginalDesignManufacturer:       entity.OriginalDesignManufacturer,
	model.OrganizationRelationshipPrivateLabelManufacturer:         entity.PrivateLabelManufacturer,
	model.OrganizationRelationshipLogisticsPartner:                 entity.LogisticsPartner,
	model.OrganizationRelationshipConsultant:                       entity.Consultant,
	model.OrganizationRelationshipServiceProvider:                  entity.ServiceProvider,
	model.OrganizationRelationshipOutsourcingProvider:              entity.OutsourcingProvider,
	model.OrganizationRelationshipInsourcingPartner:                entity.InsourcingPartner,
	model.OrganizationRelationshipTechnologyProvider:               entity.TechnologyProvider,
	model.OrganizationRelationshipDataProvider:                     entity.DataProvider,
	model.OrganizationRelationshipCertificationBody:                entity.CertificationBody,
	model.OrganizationRelationshipStandardsOrganization:            entity.StandardsOrganization,
	model.OrganizationRelationshipIndustryAnalyst:                  entity.IndustryAnalyst,
	model.OrganizationRelationshipRealEstatePartner:                entity.RealEstatePartner,
	model.OrganizationRelationshipTalentAcquisitionPartner:         entity.TalentAcquisitionPartner,
	model.OrganizationRelationshipProfessionalEmployerOrganization: entity.ProfessionalEmployerOrganization,
	model.OrganizationRelationshipResearchCollaborator:             entity.ResearchCollaborator,
	model.OrganizationRelationshipRegulatoryBody:                   entity.RegulatoryBody,
	model.OrganizationRelationshipTradeAssociationMember:           entity.TradeAssociationMember,
	model.OrganizationRelationshipCompetitor:                       entity.Competitor,
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

func MapOrgRelationshipsToModel(input []entity.OrganizationRelationship) []model.OrganizationRelationship {
	var result []model.OrganizationRelationship
	for _, v := range input {
		result = append(result, MapOrgRelationshipToModel(v))
	}
	return result
}

func MapOrgRelationshipFromModelString(str string) entity.OrganizationRelationship {
	return MapOrgRelationshipFromModel(getModelOrganizationRelationshipFromString(str))
}

func getModelOrganizationRelationshipFromString(str string) model.OrganizationRelationship {
	for k, _ := range orgRelationshipByModel {
		if string(k) == str {
			return k
		}
	}

	return ""
}
