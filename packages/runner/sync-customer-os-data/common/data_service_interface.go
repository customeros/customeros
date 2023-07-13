package common

import "github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/entity"

type SourceDataService interface {
	Start()
	Close()
	SourceId() string
	GetUsersForSync(batchSize int, runId string) []entity.UserData
	GetOrganizationsForSync(batchSize int, runId string) []entity.OrganizationData
	GetContactsForSync(batchSize int, runId string) []entity.ContactData
	GetNotesForSync(batchSize int, runId string) []entity.NoteData
	GetEmailMessagesForSync(batchSize int, runId string) []entity.EmailMessageData
	GetIssuesForSync(batchSize int, runId string) []entity.IssueData
	GetMeetingsForSync(batchSize int, runId string) []entity.MeetingData
	GetInteractionEventsForSync(batchSize int, runId string) []entity.InteractionEventData
	MarkUserProcessed(externalSyncId, runId string, synced bool) error
	MarkOrganizationProcessed(externalSyncId, runId string, synced bool) error
	MarkContactProcessed(externalSyncId, runId string, synced bool) error
	MarkNoteProcessed(externalSyncId, runId string, synced bool) error
	MarkEmailMessageProcessed(externalSyncId, runId string, synced bool) error
	MarkIssueProcessed(externalSyncId, runId string, synced bool) error
	MarkMeetingProcessed(externalSyncId, runId string, synced bool) error
	MarkInteractionEventProcessed(externalSyncId, runId string, synced bool) error
}
