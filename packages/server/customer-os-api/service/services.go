package service

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
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
	TicketService              TicketService
	InteractionSessionService  InteractionSessionService
	InteractionEventService    InteractionEventService
	PageViewService            PageViewService
}

func InitServices(driver *neo4j.DriverWithContext) *Services {
	repositories := repository.InitRepos(driver)

	services := Services{
		ContactService:             NewContactService(repositories),
		OrganizationService:        NewOrganizationService(repositories),
		ContactGroupService:        NewContactGroupService(repositories),
		CustomFieldService:         NewCustomFieldService(repositories),
		PhoneNumberService:         NewPhoneNumberService(repositories),
		EmailService:               NewEmailService(repositories),
		UserService:                NewUserService(repositories),
		FieldSetService:            NewFieldSetService(repositories),
		EntityTemplateService:      NewEntityTemplateService(repositories),
		FieldSetTemplateService:    NewFieldSetTemplateService(repositories),
		CustomFieldTemplateService: NewCustomFieldTemplateService(repositories),
		ConversationService:        NewConversationService(repositories),
		OrganizationTypeService:    NewOrganizationTypeService(repositories),
		NoteService:                NewNoteService(repositories),
		JobRoleService:             NewJobRoleService(repositories),
		LocationService:            NewLocationService(repositories),
		TagService:                 NewTagService(repositories),
		DomainService:              NewDomainService(repositories),
		TicketService:              NewTicketService(repositories),
		InteractionSessionService:  NewInteractionSessionService(repositories),
		PageViewService:            NewPageViewService(repositories),
	}
	services.TimelineEventService = NewTimelineEventService(repositories, &services)
	services.SearchService = NewSearchService(repositories, &services)
	services.QueryService = NewQueryService(repositories, &services)
	services.InteractionEventService = NewInteractionEventService(repositories, &services)

	return &services
}
