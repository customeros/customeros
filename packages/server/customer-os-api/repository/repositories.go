package repository

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	neo4jrepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"
)

type Repositories struct {
	Drivers                       Drivers
	Neo4jRepositories             *neo4jrepository.Repositories
	TimelineEventRepository       TimelineEventRepository
	OrganizationRepository        OrganizationRepository
	ContactRepository             ContactRepository
	CustomFieldTemplateRepository CustomFieldTemplateRepository
	CustomFieldRepository         CustomFieldRepository
	EntityTemplateRepository      EntityTemplateRepository
	FieldSetTemplateRepository    FieldSetTemplateRepository
	FieldSetRepository            FieldSetRepository
	UserRepository                UserRepository
	ExternalSystemRepository      ExternalSystemRepository
	NoteRepository                NoteRepository
	JobRoleRepository             JobRoleRepository
	CalendarRepository            CalendarRepository
	LocationRepository            LocationRepository
	EmailRepository               EmailRepository
	PhoneNumberRepository         PhoneNumberRepository
	TagRepository                 TagRepository
	SearchRepository              SearchRepository
	DashboardRepository           DashboardRepository
	DomainRepository              DomainRepository
	IssueRepository               IssueRepository
	InteractionEventRepository    InteractionEventRepository
	InteractionSessionRepository  InteractionSessionRepository
	AnalysisRepository            AnalysisRepository
	AttachmentRepository          AttachmentRepository
	MeetingRepository             MeetingRepository
	TenantRepository              TenantRepository
	WorkspaceRepository           WorkspaceRepository
	SocialRepository              SocialRepository
	PlayerRepository              PlayerRepository
	ActionRepository              ActionRepository
	CountryRepository             CountryRepository
	ActionItemRepository          ActionItemRepository
	LogEntryRepository            LogEntryRepository
	CommentRepository             CommentRepository
	ServiceLineItemRepository     ServiceLineItemRepository
	OpportunityRepository         OpportunityRepository
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
	repositories.Neo4jRepositories = neo4jrepository.InitNeo4jRepositories(driver, database)
	repositories.TimelineEventRepository = NewTimelineEventRepository(driver, database)
	repositories.OrganizationRepository = NewOrganizationRepository(driver, database)
	repositories.ContactRepository = NewContactRepository(driver, database)
	repositories.CustomFieldTemplateRepository = NewCustomFieldTemplateRepository(driver, database)
	repositories.CustomFieldRepository = NewCustomFieldRepository(driver, database)
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
	repositories.DashboardRepository = NewDashboardRepository(driver)
	repositories.DomainRepository = NewDomainRepository(driver, database)
	repositories.IssueRepository = NewIssueRepository(driver, database)
	repositories.InteractionEventRepository = NewInteractionEventRepository(driver, database)
	repositories.InteractionSessionRepository = NewInteractionSessionRepository(driver)
	repositories.AnalysisRepository = NewAnalysisRepository(driver)
	repositories.AttachmentRepository = NewAttachmentRepository(driver)
	repositories.MeetingRepository = NewMeetingRepository(driver)
	repositories.TenantRepository = NewTenantRepository(driver)
	repositories.WorkspaceRepository = NewWorkspaceRepository(driver)
	repositories.SocialRepository = NewSocialRepository(driver)
	repositories.PlayerRepository = NewPlayerRepository(driver)
	repositories.ActionRepository = NewActionRepository(driver)
	repositories.CountryRepository = NewCountryRepository(driver)
	repositories.ActionItemRepository = NewActionItemRepository(driver)
	repositories.LogEntryRepository = NewLogEntryRepository(driver)
	repositories.CommentRepository = NewCommentRepository(driver, database)
	repositories.ServiceLineItemRepository = NewServiceLineItemRepository(driver, database)
	repositories.OpportunityRepository = NewOpportunityRepository(driver, database)
	return &repositories
}
