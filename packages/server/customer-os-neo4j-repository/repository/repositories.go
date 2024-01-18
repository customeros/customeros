package repository

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type Repositories struct {
	ActionReadRepository              ActionReadRepository
	ActionWriteRepository             ActionWriteRepository
	BillingProfileWriteRepository     BillingProfileWriteRepository
	CommentWriteRepository            CommentWriteRepository
	CommonReadRepository              CommonReadRepository
	ContactWriteRepository            ContactWriteRepository
	ContractReadRepository            ContractReadRepository
	ContractWriteRepository           ContractWriteRepository
	CountryReadRepository             CountryReadRepository
	CountryWriteRepository            CountryWriteRepository
	CustomFieldWriteRepository        CustomFieldWriteRepository
	EmailReadRepository               EmailReadRepository
	EmailWriteRepository              EmailWriteRepository
	ExternalSystemReadRepository      ExternalSystemReadRepository
	ExternalSystemWriteRepository     ExternalSystemWriteRepository
	InteractionEventReadRepository    InteractionEventReadRepository
	InteractionEventWriteRepository   InteractionEventWriteRepository
	InteractionSessionWriteRepository InteractionSessionWriteRepository
	InvoiceReadRepository             InvoiceReadRepository
	InvoiceWriteRepository            InvoiceWriteRepository
	InvoiceLineReadRepository         InvoiceLineReadRepository
	InvoiceLineWriteRepository        InvoiceLineWriteRepository
	InvoicingCycleReadRepository      InvoicingCycleReadRepository
	InvoicingCycleWriteRepository     InvoicingCycleWriteRepository
	IssueWriteRepository              IssueWriteRepository
	JobRoleWriteRepository            JobRoleWriteRepository
	LocationWriteRepository           LocationWriteRepository
	LogEntryWriteRepository           LogEntryWriteRepository
	MasterPlanReadRepository          MasterPlanReadRepository
	MasterPlanWriteRepository         MasterPlanWriteRepository
	OpportunityReadRepository         OpportunityReadRepository
	OpportunityWriteRepository        OpportunityWriteRepository
	OrganizationReadRepository        OrganizationReadRepository
	OrganizationWriteRepository       OrganizationWriteRepository
	OrganizationPlanReadRepository    OrganizationPlanReadRepository
	OrganizationPlanWriteRepository   OrganizationPlanWriteRepository
	PhoneNumberReadRepository         PhoneNumberReadRepository
	PhoneNumberWriteRepository        PhoneNumberWriteRepository
	PlayerWriteRepository             PlayerWriteRepository
	ServiceLineItemReadRepository     ServiceLineItemReadRepository
	ServiceLineItemWriteRepository    ServiceLineItemWriteRepository
	SocialWriteRepository             SocialWriteRepository
	TagWriteRepository                TagWriteRepository
	TenantReadRepository              TenantReadRepository
	TenantWriteRepository             TenantWriteRepository
	TimelineEventReadRepository       TimelineEventReadRepository
	UserReadRepository                UserReadRepository
	UserWriteRepository               UserWriteRepository
}

func InitNeo4jRepositories(driver *neo4j.DriverWithContext, neo4jDatabase string) *Repositories {
	repositories := Repositories{
		ActionReadRepository:              NewActionReadRepository(driver, neo4jDatabase),
		ActionWriteRepository:             NewActionWriteRepository(driver, neo4jDatabase),
		BillingProfileWriteRepository:     NewBillingProfileWriteRepository(driver, neo4jDatabase),
		CommentWriteRepository:            NewCommentWriteRepository(driver, neo4jDatabase),
		CommonReadRepository:              NewCommonReadRepository(driver, neo4jDatabase),
		ContactWriteRepository:            NewContactWriteRepository(driver, neo4jDatabase),
		ContractReadRepository:            NewContractReadRepository(driver, neo4jDatabase),
		ContractWriteRepository:           NewContractWriteRepository(driver, neo4jDatabase),
		CountryReadRepository:             NewCountryReadRepository(driver, neo4jDatabase),
		CountryWriteRepository:            NewCountryWriteRepository(driver, neo4jDatabase),
		CustomFieldWriteRepository:        NewCustomFieldWriteRepository(driver, neo4jDatabase),
		EmailReadRepository:               NewEmailReadRepository(driver, neo4jDatabase),
		EmailWriteRepository:              NewEmailWriteRepository(driver, neo4jDatabase),
		ExternalSystemReadRepository:      NewExternalSystemReadRepository(driver, neo4jDatabase),
		ExternalSystemWriteRepository:     NewExternalSystemWriteRepository(driver, neo4jDatabase),
		InteractionEventReadRepository:    NewInteractionEventReadRepository(driver, neo4jDatabase),
		InteractionEventWriteRepository:   NewInteractionEventWriteRepository(driver, neo4jDatabase),
		InteractionSessionWriteRepository: NewInteractionSessionWriteRepository(driver, neo4jDatabase),
		InvoiceReadRepository:             NewInvoiceReadRepository(driver, neo4jDatabase),
		InvoiceWriteRepository:            NewInvoiceWriteRepository(driver, neo4jDatabase),
		InvoiceLineReadRepository:         NewInvoiceLineReadRepository(driver, neo4jDatabase),
		InvoiceLineWriteRepository:        NewInvoiceLineWriteRepository(driver, neo4jDatabase),
		InvoicingCycleReadRepository:      NewInvoicingCycleReadRepository(driver, neo4jDatabase),
		InvoicingCycleWriteRepository:     NewInvoicingCycleWriteRepository(driver, neo4jDatabase),
		IssueWriteRepository:              NewIssueWriteRepository(driver, neo4jDatabase),
		JobRoleWriteRepository:            NewJobRoleWriteRepository(driver, neo4jDatabase),
		LocationWriteRepository:           NewLocationWriteRepository(driver, neo4jDatabase),
		LogEntryWriteRepository:           NewLogEntryWriteRepository(driver, neo4jDatabase),
		MasterPlanReadRepository:          NewMasterPlanReadRepository(driver, neo4jDatabase),
		MasterPlanWriteRepository:         NewMasterPlanWriteRepository(driver, neo4jDatabase),
		OpportunityReadRepository:         NewOpportunityReadRepository(driver, neo4jDatabase),
		OpportunityWriteRepository:        NewOpportunityWriteRepository(driver, neo4jDatabase),
		OrganizationReadRepository:        NewOrganizationReadRepository(driver, neo4jDatabase),
		OrganizationWriteRepository:       NewOrganizationWriteRepository(driver, neo4jDatabase),
		OrganizationPlanReadRepository:    NewOrganizationPlanReadRepository(driver, neo4jDatabase),
		OrganizationPlanWriteRepository:   NewOrganizationPlanWriteRepository(driver, neo4jDatabase),
		PhoneNumberReadRepository:         NewPhoneNumberReadRepository(driver, neo4jDatabase),
		PhoneNumberWriteRepository:        NewPhoneNumberWriteRepository(driver, neo4jDatabase),
		PlayerWriteRepository:             NewPlayerWriteRepository(driver, neo4jDatabase),
		ServiceLineItemReadRepository:     NewServiceLineItemReadRepository(driver, neo4jDatabase),
		ServiceLineItemWriteRepository:    NewServiceLineItemWriteRepository(driver, neo4jDatabase),
		SocialWriteRepository:             NewSocialWriteRepository(driver, neo4jDatabase),
		TagWriteRepository:                NewTagWriteRepository(driver, neo4jDatabase),
		TenantWriteRepository:             NewTenantWriteRepository(driver, neo4jDatabase),
		TenantReadRepository:              NewTenantReadRepository(driver, neo4jDatabase),
		TimelineEventReadRepository:       NewTimelineEventReadRepository(driver, neo4jDatabase),
		UserReadRepository:                NewUserReadRepository(driver, neo4jDatabase),
		UserWriteRepository:               NewUserWriteRepository(driver, neo4jDatabase),
	}
	return &repositories
}
