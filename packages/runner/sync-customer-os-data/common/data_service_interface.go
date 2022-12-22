package common

import "github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/entity"

type DataService interface {
	Refresh()
	Close()
	SourceId() string
	GetContactsForSync(batchSize int, runId string) []entity.ContactData
	GetCompaniesForSync(batchSize int, runId string) []entity.CompanyData
	GetUsersForSync(batchSize int, runId string) []entity.UserData
	GetNotesForSync(batchSize int, runId string) []entity.NoteData
	MarkContactProcessed(externalId, runId string, synced bool) error
	MarkCompanyProcessed(externalId, runId string, synced bool) error
	MarkUserProcessed(externalId, runId string, synced bool) error
	MarkNoteProcessed(externalId, runId string, synced bool) error
}
