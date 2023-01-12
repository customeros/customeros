package common

import "github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/entity"

type DataService interface {
	Refresh()
	Close()
	SourceId() string
	GetContactsForSync(batchSize int, runId string) []entity.ContactData
	GetOrganizationsForSync(batchSize int, runId string) []entity.OrganizationData
	GetUsersForSync(batchSize int, runId string) []entity.UserData
	GetNotesForSync(batchSize int, runId string) []entity.NoteData
	GetEmailMessagesForSync(batchSize int, runId string) []entity.EmailMessageData
	MarkContactProcessed(externalId, runId string, synced bool) error
	MarkOrganizationProcessed(externalId, runId string, synced bool) error
	MarkUserProcessed(externalId, runId string, synced bool) error
	MarkNoteProcessed(externalId, runId string, synced bool) error
	MarkEmailMessageProcessed(externalId, runId string, synced bool) error
}
