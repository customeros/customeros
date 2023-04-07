package service

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/grpc_client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository"
)

type Services struct {
	ContactService             ContactService
	OrganizationService        OrganizationService
	ContactGroupService        ContactGroupService
	CustomFieldService         CustomFieldService
	PhoneNumberService         PhoneNumberService
	EmailService               EmailService
	UserService                UserService
	FieldSetService            FieldSetService
	EntityTemplateService      EntityTemplateService
	FieldSetTemplateService    FieldSetTemplateService
	CustomFieldTemplateService CustomFieldTemplateService
	ConversationService        ConversationService
	OrganizationTypeService    OrganizationTypeService
	TimelineEventService       TimelineEventService
	NoteService                NoteService
	JobRoleService             JobRoleService
	LocationService            LocationService
	TagService                 TagService
	SearchService              SearchService
	QueryService               QueryService
	DomainService              DomainService
	IssueService               IssueService
	InteractionSessionService  InteractionSessionService
	InteractionEventService    InteractionEventService
	PageViewService            PageViewService
	AnalysisService            AnalysisService
}

func InitServices(driver *neo4j.DriverWithContext, grpcClients *grpc_client.Clients) *Services {
	repositories := repository.InitRepos(driver)

	services := Services{
		ContactService:             NewContactService(repositories, grpcClients),
		OrganizationService:        NewOrganizationService(repositories),
		ContactGroupService:        NewContactGroupService(repositories),
		CustomFieldService:         NewCustomFieldService(repositories),
		PhoneNumberService:         NewPhoneNumberService(repositories, grpcClients),
		EmailService:               NewEmailService(repositories, grpcClients),
		UserService:                NewUserService(repositories),
		FieldSetService:            NewFieldSetService(repositories),
		EntityTemplateService:      NewEntityTemplateService(repositories),
		FieldSetTemplateService:    NewFieldSetTemplateService(repositories),
		CustomFieldTemplateService: NewCustomFieldTemplateService(repositories),
		ConversationService:        NewConversationService(repositories),
		OrganizationTypeService:    NewOrganizationTypeService(repositories),
		JobRoleService:             NewJobRoleService(repositories),
		LocationService:            NewLocationService(repositories),
		TagService:                 NewTagService(repositories),
		DomainService:              NewDomainService(repositories),
		IssueService:               NewIssueService(repositories),
		PageViewService:            NewPageViewService(repositories),
	}
	services.NoteService = NewNoteService(repositories, &services)
	services.TimelineEventService = NewTimelineEventService(repositories, &services)
	services.SearchService = NewSearchService(repositories, &services)
	services.QueryService = NewQueryService(repositories, &services)
	services.InteractionEventService = NewInteractionEventService(repositories, &services)
	services.InteractionSessionService = NewInteractionSessionService(repositories, &services)
	services.AnalysisService = NewAnalysisService(repositories, &services)

	return &services
}
