package entity

import "time"

type OrganizationRelationship string

const (
	Customer                         OrganizationRelationship = "Customer"
	Distributor                      OrganizationRelationship = "Distributor"
	Partner                          OrganizationRelationship = "Partner"
	LicensingPartner                 OrganizationRelationship = "Licensing partner"
	Franchisee                       OrganizationRelationship = "Franchisee"
	Franchisor                       OrganizationRelationship = "Franchisor"
	Affiliate                        OrganizationRelationship = "Affiliate"
	Reseller                         OrganizationRelationship = "Reseller"
	Influencer                       OrganizationRelationship = "Influencer or content creator"
	MediaPartner                     OrganizationRelationship = "Media partner"
	Investor                         OrganizationRelationship = "Investor"
	Merger                           OrganizationRelationship = "Merger or acquisition target"
	ParentCompany                    OrganizationRelationship = "Parent company"
	Subsidiary                       OrganizationRelationship = "Subsidiary"
	JointVenture                     OrganizationRelationship = "Joint venture"
	Sponsor                          OrganizationRelationship = "Sponsor"
	Supplier                         OrganizationRelationship = "Supplier"
	Vendor                           OrganizationRelationship = "Vendor"
	ContractManufacturer             OrganizationRelationship = "Contract manufacturer"
	OriginalEquipmentManufacturer    OrganizationRelationship = "Original equipment manufacturer"
	OriginalDesignManufacturer       OrganizationRelationship = "Original design manufacturer"
	PrivateLabelManufacturer         OrganizationRelationship = "Private label manufacturer"
	LogisticsPartner                 OrganizationRelationship = "Logistics partner"
	Consultant                       OrganizationRelationship = "Consultant"
	ServiceProvider                  OrganizationRelationship = "Service provider"
	OutsourcingProvider              OrganizationRelationship = "Outsourcing provider"
	InsourcingPartner                OrganizationRelationship = "Insourcing partner"
	TechnologyProvider               OrganizationRelationship = "Technology provider"
	DataProvider                     OrganizationRelationship = "Data provider"
	CertificationBody                OrganizationRelationship = "Certification body"
	StandardsOrganization            OrganizationRelationship = "Standards organization"
	IndustryAnalyst                  OrganizationRelationship = "Industry analyst"
	RealEstatePartner                OrganizationRelationship = "Real estate partner"
	TalentAcquisitionPartner         OrganizationRelationship = "Talent acquisition partner"
	ProfessionalEmployerOrganization OrganizationRelationship = "Professional employer organization"
	ResearchCollaborator             OrganizationRelationship = "Research collaborator"
	RegulatoryBody                   OrganizationRelationship = "Regulatory body"
	TradeAssociationMember           OrganizationRelationship = "Trade association member"
	Competitor                       OrganizationRelationship = "Competitor"
)

type OrganizationRelationships []OrganizationRelationship

type OrganizationRelationshipWithDataloaderKey struct {
	OrganizationRelationship OrganizationRelationship
	DataloaderKey            string
}

var AllOrganizationRelationship = []OrganizationRelationship{
	Customer, Distributor, Partner, LicensingPartner, Franchisee, Franchisor, Affiliate, Reseller, Influencer,
	MediaPartner, Investor, Merger, ParentCompany, Subsidiary, JointVenture, Sponsor, Supplier, Vendor,
	ContractManufacturer, OriginalEquipmentManufacturer, OriginalDesignManufacturer, PrivateLabelManufacturer,
	LogisticsPartner, Consultant, ServiceProvider, OutsourcingProvider, InsourcingPartner, TechnologyProvider,
	DataProvider, CertificationBody, StandardsOrganization, IndustryAnalyst, RealEstatePartner,
	TalentAcquisitionPartner, ProfessionalEmployerOrganization, ResearchCollaborator, RegulatoryBody,
	TradeAssociationMember, Competitor,
}

func OrganizationRelationshipFromString(input string) OrganizationRelationship {
	for _, relationship := range AllOrganizationRelationship {
		if string(relationship) == input {
			return relationship
		}
	}
	// Return a default value or handle the case when the input string doesn't match any OrganizationRelationship
	return ""
}

func (o OrganizationRelationship) String() string {
	return string(o)
}

func (o OrganizationRelationship) IsValid() bool {
	switch o {
	case Customer, Distributor, Partner, LicensingPartner, Franchisee, Franchisor, Affiliate, Reseller, Influencer,
		MediaPartner, Investor, Merger, ParentCompany, Subsidiary, JointVenture, Sponsor, Supplier, Vendor,
		ContractManufacturer, OriginalEquipmentManufacturer, OriginalDesignManufacturer, PrivateLabelManufacturer,
		LogisticsPartner, Consultant, ServiceProvider, OutsourcingProvider, InsourcingPartner, TechnologyProvider,
		DataProvider, CertificationBody, StandardsOrganization, IndustryAnalyst, RealEstatePartner, TalentAcquisitionPartner,
		ProfessionalEmployerOrganization, ResearchCollaborator, RegulatoryBody, TradeAssociationMember, Competitor:
		return true
	default:
		return false
	}
}

type OrganizationRelationshipEntity struct {
	ID            string
	CreatedAt     time.Time `neo4jDb:"property:createdAt;lookupName:CREATED_AT;supportCaseSensitive:false"`
	Name          string    `neo4jDb:"property:name;lookupName:NAME;supportCaseSensitive:true"`
	Group         string    `neo4jDb:"property:group;lookupName:GROUP;supportCaseSensitive:true"`
	DataloaderKey string
}

type OrganizationRelationshipEntities []OrganizationRelationshipEntity

func (o OrganizationRelationshipEntity) GetOrganizationRelationship() OrganizationRelationship {
	return OrganizationRelationshipFromString(o.Name)
}
