package repository

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type Repositories struct {
	Drivers                            Drivers
	TimelineEventRepository            TimelineEventRepository
	OrganizationRepository             OrganizationRepository
	ContactRepository                  ContactRepository
	CustomFieldTemplateRepository      CustomFieldTemplateRepository
	CustomFieldRepository              CustomFieldRepository
	EntityTemplateRepository           EntityTemplateRepository
	FieldSetTemplateRepository         FieldSetTemplateRepository
	FieldSetRepository                 FieldSetRepository
	UserRepository                     UserRepository
	ExternalSystemRepository           ExternalSystemRepository
	NoteRepository                     NoteRepository
	JobRoleRepository                  JobRoleRepository
	CalendarRepository                 CalendarRepository
	LocationRepository                 LocationRepository
	EmailRepository                    EmailRepository
	PhoneNumberRepository              PhoneNumberRepository
	TagRepository                      TagRepository
	SearchRepository                   SearchRepository
	QueryRepository                    DashboardRepository
	DomainRepository                   DomainRepository
	IssueRepository                    IssueRepository
	InteractionEventRepository         InteractionEventRepository
	InteractionSessionRepository       InteractionSessionRepository
	AnalysisRepository                 AnalysisRepository
	AttachmentRepository               AttachmentRepository
	MeetingRepository                  MeetingRepository
	TenantRepository                   TenantRepository
	WorkspaceRepository                WorkspaceRepository
	SocialRepository                   SocialRepository
	PlayerRepository                   PlayerRepository
	OrganizationRelationshipRepository OrganizationRelationshipRepository
	ActionRepository                   ActionRepository
	CountryRepository                  CountryRepository
	ActionItemRepository               ActionItemRepository
	LogEntryRepository                 LogEntryRepository
}

type Drivers struct {
	Neo4jDriver *neo4j.DriverWithContext
}

func InitRepos(driver *neo4j.DriverWithContext, database string) *Repositories {
	repositories := Repositories{
		Drivers: Drivers{
			Neo4jDriver: driver,
		},
	}
	repositories.TimelineEventRepository = NewTimelineEventRepository(driver)
	repositories.OrganizationRepository = NewOrganizationRepository(driver)
	repositories.ContactRepository = NewContactRepository(driver)
	repositories.CustomFieldTemplateRepository = NewCustomFieldTemplateRepository(driver)
	repositories.CustomFieldRepository = NewCustomFieldRepository(driver)
	repositories.EntityTemplateRepository = NewEntityTemplateRepository(driver, &repositories)
	repositories.FieldSetTemplateRepository = NewFieldSetTemplateRepository(driver, &repositories)
	repositories.FieldSetRepository = NewFieldSetRepository(driver)
	repositories.UserRepository = NewUserRepository(driver, database)
	repositories.ExternalSystemRepository = NewExternalSystemRepository(driver)
	repositories.NoteRepository = NewNoteRepository(driver)
	repositories.JobRoleRepository = NewJobRoleRepository(driver)
	repositories.CalendarRepository = NewCalendarRepository(driver)
	repositories.LocationRepository = NewLocationRepository(driver)
	repositories.EmailRepository = NewEmailRepository(driver, database)
	repositories.PhoneNumberRepository = NewPhoneNumberRepository(driver)
	repositories.TagRepository = NewTagRepository(driver)
	repositories.SearchRepository = NewSearchRepository(driver)
	repositories.QueryRepository = NewDashboardRepository(driver)
	repositories.DomainRepository = NewDomainRepository(driver)
	repositories.IssueRepository = NewIssueRepository(driver)
	repositories.InteractionEventRepository = NewInteractionEventRepository(driver)
	repositories.InteractionSessionRepository = NewInteractionSessionRepository(driver)
	repositories.AnalysisRepository = NewAnalysisRepository(driver)
	repositories.AttachmentRepository = NewAttachmentRepository(driver)
	repositories.MeetingRepository = NewMeetingRepository(driver)
	repositories.TenantRepository = NewTenantRepository(driver)
	repositories.WorkspaceRepository = NewWorkspaceRepository(driver)
	repositories.SocialRepository = NewSocialRepository(driver)
	repositories.PlayerRepository = NewPlayerRepository(driver)
	repositories.OrganizationRelationshipRepository = NewOrganizationRelationshipRepository(driver)
	repositories.ActionRepository = NewActionRepository(driver)
	repositories.CountryRepository = NewCountryRepository(driver)
	repositories.ActionItemRepository = NewActionItemRepository(driver)
	repositories.LogEntryRepository = NewLogEntryRepository(driver)
	return &repositories
}
