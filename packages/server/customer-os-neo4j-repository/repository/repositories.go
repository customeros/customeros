package repository

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type Repositories struct {
	Neo4jDriver *neo4j.DriverWithContext
	Database    string

	ActionReadRepository                     ActionReadRepository
	ActionWriteRepository                    ActionWriteRepository
	AttachmentReadRepository                 AttachmentReadRepository
	AttachmentWriteRepository                AttachmentWriteRepository
	BankAccountReadRepository                BankAccountReadRepository
	BankAccountWriteRepository               BankAccountWriteRepository
	BillingProfileWriteRepository            BillingProfileWriteRepository
	CommentReadRepository                    CommentReadRepository
	CommentWriteRepository                   CommentWriteRepository
	CommonReadRepository                     CommonReadRepository
	CommonWriteRepository                    CommonWriteRepository
	ContactReadRepository                    ContactReadRepository
	ContactWithFiltersReadRepository         ContactWithFiltersReadRepository
	ContactWriteRepository                   ContactWriteRepository
	ContractReadRepository                   ContractReadRepository
	ContractWriteRepository                  ContractWriteRepository
	CountryReadRepository                    CountryReadRepository
	CountryWriteRepository                   CountryWriteRepository
	CustomFieldWriteRepository               CustomFieldWriteRepository
	CustomFieldTemplateReadRepository        CustomFieldTemplateReadRepository
	CustomFieldTemplateWriteRepository       CustomFieldTemplateWriteRepository
	DomainReadRepository                     DomainReadRepository
	DomainWriteRepository                    DomainWriteRepository
	EmailReadRepository                      EmailReadRepository
	EmailWriteRepository                     EmailWriteRepository
	ExternalSystemReadRepository             ExternalSystemReadRepository
	ExternalSystemWriteRepository            ExternalSystemWriteRepository
	FlowReadRepository                       FlowReadRepository
	FlowWriteRepository                      FlowWriteRepository
	FlowActionReadRepository                 FlowActionReadRepository
	FlowActionWriteRepository                FlowActionWriteRepository
	FlowParticipantReadRepository            FlowParticipantReadRepository
	FlowParticipantWriteRepository           FlowParticipantWriteRepository
	FlowSenderReadRepository                 FlowSenderReadRepository
	FlowSenderWriteRepository                FlowSenderWriteRepository
	FlowExecutionSettingsReadRepository      FlowExecutionSettingsReadRepository
	FlowExecutionSettingsWriteRepository     FlowExecutionSettingsWriteRepository
	FlowActionExecutionReadRepository        FlowActionExecutionReadRepository
	FlowActionExecutionWriteRepository       FlowActionExecutionWriteRepository
	InteractionEventReadRepository           InteractionEventReadRepository
	InteractionEventWriteRepository          InteractionEventWriteRepository
	InteractionSessionReadRepository         InteractionSessionReadRepository
	InteractionSessionWriteRepository        InteractionSessionWriteRepository
	InvoiceReadRepository                    InvoiceReadRepository
	InvoiceWriteRepository                   InvoiceWriteRepository
	InvoiceLineReadRepository                InvoiceLineReadRepository
	InvoiceLineWriteRepository               InvoiceLineWriteRepository
	IssueReadRepository                      IssueReadRepository
	IssueWriteRepository                     IssueWriteRepository
	JobRoleReadRepository                    JobRoleReadRepository
	JobRoleWriteRepository                   JobRoleWriteRepository
	LocationWriteRepository                  LocationWriteRepository
	LogEntryReadRepository                   LogEntryReadRepository
	LogEntryWriteRepository                  LogEntryWriteRepository
	LinkedinConnectionRequestReadRepository  LinkedinConnectionRequestReadRepository
	LinkedinConnectionRequestWriteRepository LinkedinConnectionRequestWriteRepository
	OpportunityReadRepository                OpportunityReadRepository
	OpportunityWriteRepository               OpportunityWriteRepository
	OrganizationReadRepository               OrganizationReadRepository
	OrganizationWithFiltersReadRepository    OrganizationWithFiltersReadRepository
	OrganizationWriteRepository              OrganizationWriteRepository
	PhoneNumberReadRepository                PhoneNumberReadRepository
	PhoneNumberWriteRepository               PhoneNumberWriteRepository
	PlayerReadRepository                     PlayerReadRepository
	PlayerWriteRepository                    PlayerWriteRepository
	ReminderReadRepository                   ReminderReadRepository
	ReminderWriteRepository                  ReminderWriteRepository
	ServiceLineItemReadRepository            ServiceLineItemReadRepository
	ServiceLineItemWriteRepository           ServiceLineItemWriteRepository
	StateReadRepository                      StateReadRepository
	SocialReadRepository                     SocialReadRepository
	SocialWriteRepository                    SocialWriteRepository
	TagReadRepository                        TagReadRepository
	TagWriteRepository                       TagWriteRepository
	TenantReadRepository                     TenantReadRepository
	TenantWriteRepository                    TenantWriteRepository
	TimelineEventReadRepository              TimelineEventReadRepository
	UserReadRepository                       UserReadRepository
	UserWriteRepository                      UserWriteRepository
	WorkspaceReadRepository                  WorkspaceReadRepository
	WorkspaceWriteRepository                 WorkspaceWriteRepository
}

func InitNeo4jRepositories(driver *neo4j.DriverWithContext, neo4jDatabase string) *Repositories {
	repositories := Repositories{
		Neo4jDriver:                              driver,
		Database:                                 neo4jDatabase,
		ActionReadRepository:                     NewActionReadRepository(driver, neo4jDatabase),
		ActionWriteRepository:                    NewActionWriteRepository(driver, neo4jDatabase),
		AttachmentReadRepository:                 NewAttachmentReadRepository(driver, neo4jDatabase),
		AttachmentWriteRepository:                NewAttachmentWriteRepository(driver, neo4jDatabase),
		BankAccountReadRepository:                NewBankAccountReadRepository(driver, neo4jDatabase),
		BankAccountWriteRepository:               NewBankAccountWriteRepository(driver, neo4jDatabase),
		BillingProfileWriteRepository:            NewBillingProfileWriteRepository(driver, neo4jDatabase),
		CommentReadRepository:                    NewCommentReadRepository(driver, neo4jDatabase),
		CommentWriteRepository:                   NewCommentWriteRepository(driver, neo4jDatabase),
		CommonReadRepository:                     NewCommonReadRepository(driver, neo4jDatabase),
		CommonWriteRepository:                    NewCommonWriteRepository(driver, neo4jDatabase),
		ContactReadRepository:                    NewContactReadRepository(driver, neo4jDatabase),
		ContactWithFiltersReadRepository:         NewContactWithFiltersReadRepository(driver, neo4jDatabase),
		ContactWriteRepository:                   NewContactWriteRepository(driver, neo4jDatabase),
		ContractReadRepository:                   NewContractReadRepository(driver, neo4jDatabase),
		ContractWriteRepository:                  NewContractWriteRepository(driver, neo4jDatabase),
		CountryReadRepository:                    NewCountryReadRepository(driver, neo4jDatabase),
		CountryWriteRepository:                   NewCountryWriteRepository(driver, neo4jDatabase),
		CustomFieldWriteRepository:               NewCustomFieldWriteRepository(driver, neo4jDatabase),
		CustomFieldTemplateReadRepository:        NewCustomFieldTemplateReadRepository(driver, neo4jDatabase),
		CustomFieldTemplateWriteRepository:       NewCustomFieldTemplateWriteRepository(driver, neo4jDatabase),
		DomainReadRepository:                     NewDomainReadRepository(driver, neo4jDatabase),
		DomainWriteRepository:                    NewDomainWriteRepository(driver, neo4jDatabase),
		EmailReadRepository:                      NewEmailReadRepository(driver, neo4jDatabase),
		EmailWriteRepository:                     NewEmailWriteRepository(driver, neo4jDatabase),
		ExternalSystemReadRepository:             NewExternalSystemReadRepository(driver, neo4jDatabase),
		ExternalSystemWriteRepository:            NewExternalSystemWriteRepository(driver, neo4jDatabase),
		FlowReadRepository:                       NewFlowReadRepository(driver, neo4jDatabase),
		FlowWriteRepository:                      NewFlowWriteRepository(driver, neo4jDatabase),
		FlowActionReadRepository:                 NewFlowActionReadRepository(driver, neo4jDatabase),
		FlowActionWriteRepository:                NewFlowActionWriteRepository(driver, neo4jDatabase),
		FlowParticipantReadRepository:            NewFlowParticipantReadRepository(driver, neo4jDatabase),
		FlowParticipantWriteRepository:           NewFlowParticipantWriteRepository(driver, neo4jDatabase),
		FlowSenderReadRepository:                 NewFlowSenderReadRepository(driver, neo4jDatabase),
		FlowSenderWriteRepository:                NewFlowSenderWriteRepository(driver, neo4jDatabase),
		FlowActionExecutionReadRepository:        NewFlowActionExecutionReadRepository(driver, neo4jDatabase),
		FlowExecutionSettingsWriteRepository:     NewFlowExecutionSettingsWriteRepository(driver, neo4jDatabase),
		FlowActionExecutionWriteRepository:       NewFlowActionExecutionWriteRepository(driver, neo4jDatabase),
		FlowExecutionSettingsReadRepository:      NewFlowExecutionSettingsReadRepository(driver, neo4jDatabase),
		InteractionEventReadRepository:           NewInteractionEventReadRepository(driver, neo4jDatabase),
		InteractionEventWriteRepository:          NewInteractionEventWriteRepository(driver, neo4jDatabase),
		InteractionSessionReadRepository:         NewInteractionSessionReadRepository(driver, neo4jDatabase),
		InteractionSessionWriteRepository:        NewInteractionSessionWriteRepository(driver, neo4jDatabase),
		InvoiceReadRepository:                    NewInvoiceReadRepository(driver, neo4jDatabase),
		InvoiceWriteRepository:                   NewInvoiceWriteRepository(driver, neo4jDatabase),
		InvoiceLineReadRepository:                NewInvoiceLineReadRepository(driver, neo4jDatabase),
		InvoiceLineWriteRepository:               NewInvoiceLineWriteRepository(driver, neo4jDatabase),
		IssueReadRepository:                      NewIssueReadRepository(driver, neo4jDatabase),
		IssueWriteRepository:                     NewIssueWriteRepository(driver, neo4jDatabase),
		JobRoleReadRepository:                    NewJobRoleReadRepository(driver, neo4jDatabase),
		JobRoleWriteRepository:                   NewJobRoleWriteRepository(driver, neo4jDatabase),
		LocationWriteRepository:                  NewLocationWriteRepository(driver, neo4jDatabase),
		LogEntryReadRepository:                   NewLogEntryReadRepository(driver, neo4jDatabase),
		LogEntryWriteRepository:                  NewLogEntryWriteRepository(driver, neo4jDatabase),
		LinkedinConnectionRequestReadRepository:  NewLinkedinConnectionRequestReadRepository(driver, neo4jDatabase),
		LinkedinConnectionRequestWriteRepository: NewLinkedinConnectionRequestWriteRepository(driver, neo4jDatabase),
		OpportunityReadRepository:                NewOpportunityReadRepository(driver, neo4jDatabase),
		OpportunityWriteRepository:               NewOpportunityWriteRepository(driver, neo4jDatabase),
		OrganizationReadRepository:               NewOrganizationReadRepository(driver, neo4jDatabase),
		OrganizationWithFiltersReadRepository:    NewOrganizationWithFiltersReadRepository(driver, neo4jDatabase),
		OrganizationWriteRepository:              NewOrganizationWriteRepository(driver, neo4jDatabase),
		PhoneNumberReadRepository:                NewPhoneNumberReadRepository(driver, neo4jDatabase),
		PhoneNumberWriteRepository:               NewPhoneNumberWriteRepository(driver, neo4jDatabase),
		PlayerReadRepository:                     NewPlayerReadRepository(driver, neo4jDatabase),
		PlayerWriteRepository:                    NewPlayerWriteRepository(driver, neo4jDatabase),
		ReminderReadRepository:                   NewReminderReadRepository(driver, neo4jDatabase),
		ReminderWriteRepository:                  NewReminderWriteRepository(driver, neo4jDatabase),
		ServiceLineItemReadRepository:            NewServiceLineItemReadRepository(driver, neo4jDatabase),
		ServiceLineItemWriteRepository:           NewServiceLineItemWriteRepository(driver, neo4jDatabase),
		StateReadRepository:                      NewStateReadRepository(driver, neo4jDatabase),
		SocialReadRepository:                     NewSocialReadRepository(driver, neo4jDatabase),
		SocialWriteRepository:                    NewSocialWriteRepository(driver, neo4jDatabase),
		TagReadRepository:                        NewTagReadRepository(driver, neo4jDatabase),
		TagWriteRepository:                       NewTagWriteRepository(driver, neo4jDatabase),
		TenantWriteRepository:                    NewTenantWriteRepository(driver, neo4jDatabase),
		TenantReadRepository:                     NewTenantReadRepository(driver, neo4jDatabase),
		TimelineEventReadRepository:              NewTimelineEventReadRepository(driver, neo4jDatabase),
		UserReadRepository:                       NewUserReadRepository(driver, neo4jDatabase),
		UserWriteRepository:                      NewUserWriteRepository(driver, neo4jDatabase),
		WorkspaceReadRepository:                  NewWorkspaceReadRepository(driver, neo4jDatabase),
		WorkspaceWriteRepository:                 NewWorkspaceWriteRepository(driver, neo4jDatabase),
	}
	return &repositories
}
