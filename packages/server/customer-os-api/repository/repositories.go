package repository

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type Repositories struct {
	Drivers                       Drivers
	TimelineEventRepository       TimelineEventRepository
	OrganizationRepository        OrganizationRepository
	OrganizationTypeRepository    OrganizationTypeRepository
	ContactGroupRepository        ContactGroupRepository
	ContactRepository             ContactRepository
	ConversationRepository        ConversationRepository
	CustomFieldTemplateRepository CustomFieldTemplateRepository
	CustomFieldRepository         CustomFieldRepository
	EntityTemplateRepository      EntityTemplateRepository
	FieldSetTemplateRepository    FieldSetTemplateRepository
	FieldSetRepository            FieldSetRepository
	UserRepository                UserRepository
	ExternalSystemRepository      ExternalSystemRepository
	NoteRepository                NoteRepository
	JobRoleRepository             JobRoleRepository
	LocationRepository            LocationRepository
	EmailRepository               EmailRepository
	PhoneNumberRepository         PhoneNumberRepository
	TagRepository                 TagRepository
	SearchRepository              SearchRepository
	QueryRepository               QueryRepository
	DomainRepository              DomainRepository
	IssueRepository               IssueRepository
	InteractionEventRepository    InteractionEventRepository
	InteractionSessionRepository  InteractionSessionRepository
	AnalysisRepository            AnalysisRepository
	AttachmentRepository          AttachmentRepository
	MeetingRepository             MeetingRepository
}

type Drivers struct {
	Neo4jDriver *neo4j.DriverWithContext
}

func InitRepos(driver *neo4j.DriverWithContext) *Repositories {
	repositories := Repositories{
		Drivers: Drivers{
			Neo4jDriver: driver,
		},
	}
	repositories.TimelineEventRepository = NewTimelineEventRepository(driver)
	repositories.OrganizationRepository = NewOrganizationRepository(driver)
	repositories.ContactGroupRepository = NewContactGroupRepository(driver)
	repositories.ContactRepository = NewContactRepository(driver)
	repositories.ConversationRepository = NewConversationRepository(driver)
	repositories.CustomFieldTemplateRepository = NewCustomFieldTemplateRepository(driver)
	repositories.CustomFieldRepository = NewCustomFieldRepository(driver)
	repositories.EntityTemplateRepository = NewEntityTemplateRepository(driver, &repositories)
	repositories.FieldSetTemplateRepository = NewFieldSetTemplateRepository(driver, &repositories)
	repositories.FieldSetRepository = NewFieldSetRepository(driver)
	repositories.UserRepository = NewUserRepository(driver)
	repositories.ExternalSystemRepository = NewExternalSystemRepository(driver)
	repositories.NoteRepository = NewNoteRepository(driver)
	repositories.JobRoleRepository = NewJobRoleRepository(driver)
	repositories.LocationRepository = NewLocationRepository(driver)
	repositories.EmailRepository = NewEmailRepository(driver)
	repositories.PhoneNumberRepository = NewPhoneNumberRepository(driver)
	repositories.OrganizationTypeRepository = NewOrganizationTypeRepository(driver)
	repositories.TagRepository = NewTagRepository(driver)
	repositories.SearchRepository = NewSearchRepository(driver)
	repositories.QueryRepository = NewQueryRepository(driver)
	repositories.DomainRepository = NewDomainRepository(driver)
	repositories.IssueRepository = NewIssueRepository(driver)
	repositories.InteractionEventRepository = NewInteractionEventRepository(driver)
	repositories.InteractionSessionRepository = NewInteractionSessionRepository(driver)
	repositories.AnalysisRepository = NewAnalysisRepository(driver)
	repositories.AttachmentRepository = NewAttachmentRepository(driver)
	repositories.MeetingRepository = NewMeetingRepository(driver)
	return &repositories
}
