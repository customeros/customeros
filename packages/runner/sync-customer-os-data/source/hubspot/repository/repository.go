package repository

import (
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/source/hubspot/entity"
	"gorm.io/gorm"
	"time"
)

const (
	ContactEntity = "contact"
	CompanyEntity = "company"
	OwnerEntity   = "owner"
)

func GetContacts(db *gorm.DB, limit int, runId string) (entity.Contacts, error) {
	var contacts entity.Contacts

	cte := `
		WITH UpToDateData AS (
    		SELECT row_number() OVER (PARTITION BY id ORDER BY updatedat DESC) AS row_num, *
    		FROM contacts
		)`
	err := db.
		Raw(cte+" SELECT u.* FROM UpToDateData u left join openline_sync_status_contacts s "+
			" on u.id = s.id and u._airbyte_ab_id = s._airbyte_ab_id and u._airbyte_contacts_hashid = s._airbyte_contacts_hashid "+
			" WHERE u.row_num = ? "+
			" and (s.synced_to_customer_os is null or s.synced_to_customer_os = ?) "+
			" and (s.synced_to_customer_os_attempt is null or s.synced_to_customer_os_attempt < ?) "+
			" and (s.run_id is null or s.run_id <> ?) "+
			" limit ?", 1, false, 10, runId, limit).
		Find(&contacts).Error

	if err != nil {
		return nil, err
	}
	return contacts, nil
}

func GetContactProperties(db *gorm.DB, airbyteAbId, airbyteContactsHashId string) (entity.ContactProperties, error) {
	contactProperties := entity.ContactProperties{}
	err := db.Table(entity.ContactProperties{}.TableName()).
		Where(&entity.ContactProperties{AirbyteAbId: airbyteAbId, AirbyteContactsHashid: airbyteContactsHashId}).
		First(&contactProperties).Error
	return contactProperties, err
}

func MarkContactProcessed(db *gorm.DB, contact entity.Contact, synced bool, runId string) error {
	syncStatusContact := entity.SyncStatusContact{
		Id:                    contact.Id,
		AirbyteAbId:           contact.AirbyteAbId,
		AirbyteContactsHashid: contact.AirbyteContactsHashid,
	}
	db.FirstOrCreate(&syncStatusContact, syncStatusContact)

	return db.Model(&syncStatusContact).
		Where(&entity.SyncStatusContact{Id: contact.Id, AirbyteAbId: contact.AirbyteAbId, AirbyteContactsHashid: contact.AirbyteContactsHashid}).
		Updates(entity.SyncStatusContact{
			SyncedToCustomerOs: synced,
			SyncedAt:           time.Now(),
			SyncAttempt:        syncStatusContact.SyncAttempt + 1,
			RunId:              runId,
		}).Error
}

func GetCompanies(db *gorm.DB, limit int, runId string) (entity.CompaniesRaw, error) {
	var rawCompanies entity.CompaniesRaw

	err := db.
		Raw(`SELECT a.*
FROM _airbyte_raw_companies a
LEFT JOIN openline_sync_status s ON a._airbyte_ab_id = s._airbyte_ab_id and s.entity = ?
WHERE (s.synced_to_customer_os IS NULL OR s.synced_to_customer_os = FALSE)
  AND (s.synced_to_customer_os_attempt IS NULL OR s.synced_to_customer_os_attempt < ?)
  AND (s.run_id IS NULL OR s.run_id <> ?)
ORDER BY a._airbyte_emitted_at ASC
LIMIT ?`, CompanyEntity, 10, runId, limit).
		Find(&rawCompanies).Error

	if err != nil {
		return nil, err
	}
	return rawCompanies, nil
}

func MarkProcessed(db *gorm.DB, syncedEntity, airbyteAbId string, synced bool, runId, externalSyncId string) error {
	syncStatus := entity.SyncStatus{
		Entity:         syncedEntity,
		AirbyteAbId:    airbyteAbId,
		ExternalSyncId: externalSyncId,
	}
	db.FirstOrCreate(&syncStatus, syncStatus)

	return db.Model(&syncStatus).
		Where(&entity.SyncStatus{AirbyteAbId: airbyteAbId, Entity: syncedEntity}).
		Updates(entity.SyncStatus{
			SyncedToCustomerOs: synced,
			SyncedAt:           time.Now(),
			SyncAttempt:        syncStatus.SyncAttempt + 1,
			RunId:              runId,
		}).Error
}

func GetOwners(db *gorm.DB, limit int, runId string) (entity.OwnersRaw, error) {
	var owners entity.OwnersRaw

	err := db.
		Raw(`SELECT a.*
FROM _airbyte_raw_owners a
LEFT JOIN openline_sync_status s ON a._airbyte_ab_id = s._airbyte_ab_id and s.entity = ?
WHERE (s.synced_to_customer_os IS NULL OR s.synced_to_customer_os = FALSE)
  AND (s.synced_to_customer_os_attempt IS NULL OR s.synced_to_customer_os_attempt < ?)
  AND (s.run_id IS NULL OR s.run_id <> ?)
ORDER BY a._airbyte_emitted_at ASC
LIMIT ?`, OwnerEntity, 10, runId, limit).
		Find(&owners).Error

	if err != nil {
		return nil, err
	}
	return owners, nil
}

func GetNotes(db *gorm.DB, limit int, runId string) (entity.Notes, error) {
	var notes entity.Notes

	cte := `
		WITH UpToDateData AS (
    		SELECT row_number() OVER (PARTITION BY id ORDER BY updatedat DESC) AS row_num, *
    		FROM engagements_notes
		)`
	err := db.
		Raw(cte+" SELECT u.* FROM UpToDateData u left join openline_sync_status_notes s "+
			" on u.id = s.id and u._airbyte_ab_id = s._airbyte_ab_id and u._airbyte_engagements_notes_hashid = s._airbyte_engagements_notes_hashid "+
			" WHERE u.row_num = ? "+
			" and (u.contacts is not null or u.companies is not null) "+
			" and (s.synced_to_customer_os is null or s.synced_to_customer_os = ?) "+
			" and (s.synced_to_customer_os_attempt is null or s.synced_to_customer_os_attempt < ?) "+
			" and (s.run_id is null or s.run_id <> ?) "+
			" limit ?", 1, false, 10, runId, limit).
		Find(&notes).Error

	if err != nil {
		return nil, err
	}
	return notes, nil
}

func GetNoteProperties(db *gorm.DB, airbyteAbId, airbyteNotesHashId string) (entity.NoteProperties, error) {
	noteProperties := entity.NoteProperties{}
	err := db.Table(entity.NoteProperties{}.TableName()).
		Where(&entity.NoteProperties{AirbyteAbId: airbyteAbId, AirbyteNotesHashid: airbyteNotesHashId}).
		First(&noteProperties).Error
	return noteProperties, err
}

func MarkNoteProcessed(db *gorm.DB, note entity.Note, synced bool, runId string) error {
	syncStatusNote := entity.SyncStatusNote{
		Id:                 note.Id,
		AirbyteAbId:        note.AirbyteAbId,
		AirbyteNotesHashid: note.AirbyteNotesHashid,
	}
	db.FirstOrCreate(&syncStatusNote, syncStatusNote)

	return db.Model(&syncStatusNote).
		Where(&entity.SyncStatusNote{Id: note.Id, AirbyteAbId: note.AirbyteAbId, AirbyteNotesHashid: note.AirbyteNotesHashid}).
		Updates(entity.SyncStatusNote{
			SyncedToCustomerOs: synced,
			SyncedAt:           time.Now(),
			SyncAttempt:        syncStatusNote.SyncAttempt + 1,
			RunId:              runId,
		}).Error
}

func GetEmails(db *gorm.DB, limit int, runId string) (entity.Emails, error) {
	var emails entity.Emails

	cte := `
		WITH UpToDateData AS (
    		SELECT row_number() OVER (PARTITION BY id ORDER BY updatedat DESC) AS row_num, *
    		FROM engagements_emails
		)`
	err := db.
		Raw(cte+" SELECT u.* FROM UpToDateData u left join openline_sync_status_emails s "+
			" on u.id = s.id and u._airbyte_ab_id = s._airbyte_ab_id and u._airbyte_engagements_emails_hashid = s._airbyte_engagements_emails_hashid "+
			" left join engagements_emails_properties p "+
			" on u._airbyte_ab_id = p._airbyte_ab_id and u._airbyte_engagements_emails_hashid = p._airbyte_engagements_emails_hashid "+
			" WHERE u.row_num = ? "+
			" and (p.hs_email_status = 'SENT' and p.hs_email_thread_id is not null) "+
			" and (s.synced_to_customer_os is null or s.synced_to_customer_os = ?) "+
			" and (s.synced_to_customer_os_attempt is null or s.synced_to_customer_os_attempt < ?) "+
			" and (s.run_id is null or s.run_id <> ?) "+
			" order by u.createdat "+
			" limit ?", 1, false, 10, runId, limit).
		Find(&emails).Error

	if err != nil {
		return nil, err
	}
	return emails, nil
}

func GetEmailProperties(db *gorm.DB, airbyteAbId, airbyteEmailsHashId string) (entity.EmailProperties, error) {
	emailProperties := entity.EmailProperties{}
	err := db.Table(entity.EmailProperties{}.TableName()).
		Where(&entity.EmailProperties{AirbyteAbId: airbyteAbId, AirbyteEmailsHashid: airbyteEmailsHashId}).
		First(&emailProperties).Error
	return emailProperties, err
}

func MarkEmailProcessed(db *gorm.DB, email entity.Email, synced bool, runId string) error {
	syncStatusEmail := entity.SyncStatusEmail{
		Id:                  email.Id,
		AirbyteAbId:         email.AirbyteAbId,
		AirbyteEmailsHashid: email.AirbyteEmailsHashid,
	}
	db.FirstOrCreate(&syncStatusEmail, syncStatusEmail)

	return db.Model(&syncStatusEmail).
		Where(&entity.SyncStatusEmail{Id: email.Id, AirbyteAbId: email.AirbyteAbId, AirbyteEmailsHashid: email.AirbyteEmailsHashid}).
		Updates(entity.SyncStatusEmail{
			SyncedToCustomerOs: synced,
			SyncedAt:           time.Now(),
			SyncAttempt:        syncStatusEmail.SyncAttempt + 1,
			RunId:              runId,
		}).Error
}

func GetMeetings(db *gorm.DB, limit int, runId string) (entity.Meetings, error) {
	var meetings entity.Meetings

	cte := `
		WITH UpToDateData AS (
    		SELECT row_number() OVER (PARTITION BY id ORDER BY updatedat DESC) AS row_num, *
    		FROM engagements_meetings
		)`
	err := db.
		Raw(cte+" SELECT u.* FROM UpToDateData u left join openline_sync_status_meetings s "+
			" on u.id = s.id and u._airbyte_ab_id = s._airbyte_ab_id and u._airbyte_engagements_meetings_hashid = s._airbyte_engagements_meetings_hashid "+
			" WHERE u.row_num = ? "+
			" and (u.contacts is not null) "+
			" and (s.synced_to_customer_os is null or s.synced_to_customer_os = ?) "+
			" and (s.synced_to_customer_os_attempt is null or s.synced_to_customer_os_attempt < ?) "+
			" and (s.run_id is null or s.run_id <> ?) "+
			" limit ?", 1, false, 10, runId, limit).
		Find(&meetings).Error

	if err != nil {
		return nil, err
	}
	return meetings, nil
}

func GetMeetingProperties(db *gorm.DB, airbyteAbId, airbyteMeetingsHashId string) (entity.MeetingProperties, error) {
	meetingProperties := entity.MeetingProperties{}
	err := db.Table(entity.MeetingProperties{}.TableName()).
		Where(&entity.MeetingProperties{AirbyteAbId: airbyteAbId, AirbyteMeetingsHashid: airbyteMeetingsHashId}).
		First(&meetingProperties).Error
	return meetingProperties, err
}

func MarkMeetingProcessed(db *gorm.DB, meeting entity.Meeting, synced bool, runId string) error {
	syncStatusMeeting := entity.SyncStatusMeeting{
		Id:                    meeting.Id,
		AirbyteAbId:           meeting.AirbyteAbId,
		AirbyteMeetingsHashid: meeting.AirbyteMeetingsHashid,
	}
	db.FirstOrCreate(&syncStatusMeeting, syncStatusMeeting)

	return db.Model(&syncStatusMeeting).
		Where(&entity.SyncStatusMeeting{Id: meeting.Id, AirbyteAbId: meeting.AirbyteAbId, AirbyteMeetingsHashid: meeting.AirbyteMeetingsHashid}).
		Updates(entity.SyncStatusMeeting{
			SyncedToCustomerOs: synced,
			SyncedAt:           time.Now(),
			SyncAttempt:        syncStatusMeeting.SyncAttempt + 1,
			RunId:              runId,
		}).Error
}
