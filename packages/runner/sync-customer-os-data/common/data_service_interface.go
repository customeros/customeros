package common

import "github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/entity"

type SourceDataService interface {
	Refresh()
	Close()
	SourceId() string
	GetContactsForSync(batchSize int, runId string) []entity.ContactData
	GetOrganizationsForSync(batchSize int, runId string) []entity.OrganizationData
	GetUsersForSync(batchSize int, runId string) []entity.UserData
	GetNotesForSync(batchSize int, runId string) []entity.NoteData
	GetEmailMessagesForSync(batchSize int, runId string) []entity.EmailMessageData
	MarkContactProcessed(externalSyncId, runId string, synced bool) error
	MarkOrganizationProcessed(externalSyncId, runId string, synced bool) error
	MarkUserProcessed(externalSyncId, runId string, synced bool) error
	MarkNoteProcessed(externalSyncId, runId string, synced bool) error
	MarkEmailMessageProcessed(externalSyncId, runId string, synced bool) error
}
