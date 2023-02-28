package common

import "github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/entity"

type SourceDataService interface {
	Refresh()
	Close()
	SourceId() string
	GetUsersForSync(batchSize int, runId string) []entity.UserData
	GetOrganizationsForSync(batchSize int, runId string) []entity.OrganizationData
	GetContactsForSync(batchSize int, runId string) []entity.ContactData
	GetNotesForSync(batchSize int, runId string) []entity.NoteData
	GetEmailMessagesForSync(batchSize int, runId string) []entity.EmailMessageData
	GetTicketsForSync(batchSize int, runId string) []entity.TicketData
	MarkUserProcessed(externalSyncId, runId string, synced bool) error
	MarkOrganizationProcessed(externalSyncId, runId string, synced bool) error
	MarkContactProcessed(externalSyncId, runId string, synced bool) error
	MarkNoteProcessed(externalSyncId, runId string, synced bool) error
	MarkEmailMessageProcessed(externalSyncId, runId string, synced bool) error
	MarkTicketProcessed(externalSyncId, runId string, synced bool) error
}
