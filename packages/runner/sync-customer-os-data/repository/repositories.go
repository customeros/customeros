package repository

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/config"
	"gorm.io/gorm"
)

type Dbs struct {
	ControlDb      *gorm.DB
	Neo4jDriver    *neo4j.DriverWithContext
	RawDataStoreDB *config.RawDataStoreDB
}

type Repositories struct {
	Dbs                          Dbs
	TenantSyncSettingsRepository TenantSyncSettingsRepository
	SyncRunRepository            SyncRunRepository

	ContactRepository          ContactRepository
	EmailRepository            EmailRepository
	PhoneNumberRepository      PhoneNumberRepository
	LocationRepository         LocationRepository
	ExternalSystemRepository   ExternalSystemRepository
	OrganizationRepository     OrganizationRepository
	RoleRepository             JobRoleRepository
	UserRepository             UserRepository
	NoteRepository             NoteRepository
	InteractionEventRepository InteractionEventRepository
	IssueRepository            IssueRepository
	MeetingRepository          MeetingRepository
	ActionRepository           ActionRepository
}

func InitRepos(driver *neo4j.DriverWithContext, controlDb *gorm.DB, airbyteStoreDb *config.RawDataStoreDB) *Repositories {
	repositories := Repositories{
		Dbs: Dbs{
			Neo4jDriver:    driver,
			ControlDb:      controlDb,
			RawDataStoreDB: airbyteStoreDb,
		},
		TenantSyncSettingsRepository: NewTenantSyncSettingsRepository(controlDb),
		SyncRunRepository:            NewSyncRunRepository(controlDb),
		ContactRepository:            NewContactRepository(driver),
		EmailRepository:              NewEmailRepository(driver),
		PhoneNumberRepository:        NewPhoneNumberRepository(driver),
		LocationRepository:           NewLocationRepository(driver),
		ExternalSystemRepository:     NewExternalSystemRepository(driver),
		OrganizationRepository:       NewOrganizationRepository(driver),
		RoleRepository:               NewJobRoleRepository(driver),
		UserRepository:               NewUserRepository(driver),
		NoteRepository:               NewNoteRepository(driver),
		InteractionEventRepository:   NewInteractionEventRepository(driver),
		IssueRepository:              NewIssueRepository(driver),
		MeetingRepository:            NewMeetingRepository(driver),
		ActionRepository:             NewActionRepository(driver),
	}
	return &repositories
}
