package source

import (
	"encoding/json"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/entity"
)

func MapJsonToUser(jsonData, syncId, source string) (entity.UserData, error) {
	user := entity.UserData{}
	err := json.Unmarshal([]byte(jsonData), &user)
	if err != nil {
		return entity.UserData{}, err
	}
	user.SyncId = syncId
	user.ExternalSystem = source
	user.Id = ""
	return user, nil
}

func MapJsonToOrganization(jsonData, syncId, source string) (entity.OrganizationData, error) {
	org := entity.OrganizationData{}
	err := json.Unmarshal([]byte(jsonData), &org)
	if err != nil {
		return entity.OrganizationData{}, err
	}
	org.SyncId = syncId
	org.ExternalSystem = source
	org.Id = ""
	return org, nil
}

func MapJsonToContact(jsonData, syncId, source string) (entity.ContactData, error) {
	contact := entity.ContactData{}
	err := json.Unmarshal([]byte(jsonData), &contact)
	if err != nil {
		return entity.ContactData{}, err
	}
	contact.SyncId = syncId
	contact.ExternalSystem = source
	for _, textCustomField := range contact.TextCustomFields {
		textCustomField.ExternalSystem = source
	}
	contact.Id = ""
	return contact, nil
}

func MapJsonToLogEntry(jsonData, syncId, source string) (entity.LogEntryData, error) {
	logEntry := entity.LogEntryData{}
	err := json.Unmarshal([]byte(jsonData), &logEntry)
	if err != nil {
		return entity.LogEntryData{}, err
	}
	logEntry.SyncId = syncId
	logEntry.ExternalSystem = source
	logEntry.Id = ""
	return logEntry, nil
}

func MapJsonToNote(jsonData, syncId, source string) (entity.NoteData, error) {
	note := entity.NoteData{}
	err := json.Unmarshal([]byte(jsonData), &note)
	if err != nil {
		return entity.NoteData{}, err
	}
	note.SyncId = syncId
	note.ExternalSystem = source
	note.Id = ""
	return note, nil
}

func MapJsonToEmailMessage(jsonData, syncId, source string) (entity.EmailMessageData, error) {
	note := entity.EmailMessageData{}
	err := json.Unmarshal([]byte(jsonData), &note)
	if err != nil {
		return entity.EmailMessageData{}, err
	}
	note.SyncId = syncId
	note.ExternalSystem = source
	note.Id = ""
	return note, nil
}

func MapJsonToMeeting(jsonData, syncId, source string) (entity.MeetingData, error) {
	note := entity.MeetingData{}
	err := json.Unmarshal([]byte(jsonData), &note)
	if err != nil {
		return entity.MeetingData{}, err
	}
	note.SyncId = syncId
	note.ExternalSystem = source
	note.Id = ""
	return note, nil
}

func MapJsonToIssue(jsonData, syncId, source string) (entity.IssueData, error) {
	issue := entity.IssueData{}
	err := json.Unmarshal([]byte(jsonData), &issue)
	if err != nil {
		return entity.IssueData{}, err
	}
	issue.SyncId = syncId
	issue.ExternalSystem = source
	issue.Id = ""
	return issue, nil
}

func MapJsonToInteractionEvent(jsonData, syncId, source string) (entity.InteractionEventData, error) {
	interactionEvent := entity.InteractionEventData{}
	err := json.Unmarshal([]byte(jsonData), &interactionEvent)
	if err != nil {
		return entity.InteractionEventData{}, err
	}
	interactionEvent.SyncId = syncId
	interactionEvent.ExternalSystem = source
	interactionEvent.Id = ""
	return interactionEvent, nil
}
