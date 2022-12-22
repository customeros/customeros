package repository

import (
	hubspotEntity "github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/hubspot/entity"
	"gorm.io/gorm"
	"time"
)

func GetContacts(db *gorm.DB, limit int, runId string) (hubspotEntity.Contacts, error) {
	var contacts hubspotEntity.Contacts

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

func GetContactProperties(db *gorm.DB, airbyteAbId, airbyteContactsHashId string) (hubspotEntity.ContactProperties, error) {
	contactProperties := hubspotEntity.ContactProperties{}
	err := db.Table(hubspotEntity.ContactProperties{}.TableName()).
		Where(&hubspotEntity.ContactProperties{AirbyteAbId: airbyteAbId, AirbyteContactsHashid: airbyteContactsHashId}).
		First(&contactProperties).Error
	return contactProperties, err
}

func MarkContactProcessed(db *gorm.DB, contact hubspotEntity.Contact, synced bool, runId string) error {
	syncStatusContact := hubspotEntity.SyncStatusContact{
		Id:                    contact.Id,
		AirbyteAbId:           contact.AirbyteAbId,
		AirbyteContactsHashid: contact.AirbyteContactsHashid,
	}
	db.FirstOrCreate(&syncStatusContact, syncStatusContact)

	return db.Model(&syncStatusContact).
		Where(&hubspotEntity.SyncStatusContact{Id: contact.Id, AirbyteAbId: contact.AirbyteAbId, AirbyteContactsHashid: contact.AirbyteContactsHashid}).
		Updates(hubspotEntity.SyncStatusContact{
			SyncedToCustomerOs: synced,
			SyncedAt:           time.Now(),
			SyncAttempt:        syncStatusContact.SyncAttempt + 1,
			RunId:              runId,
		}).Error
}

func GetCompanies(db *gorm.DB, limit int, runId string) (hubspotEntity.Companies, error) {
	var companies hubspotEntity.Companies

	cte := `
		WITH UpToDateData AS (
    		SELECT row_number() OVER (PARTITION BY id ORDER BY updatedat DESC) AS row_num, *
    		FROM companies
		)`
	err := db.
		Raw(cte+" SELECT u.* FROM UpToDateData u left join openline_sync_status_companies s "+
			" on u.id = s.id and u._airbyte_ab_id = s._airbyte_ab_id and u._airbyte_companies_hashid = s._airbyte_companies_hashid "+
			" WHERE u.row_num = ? "+
			" and (s.synced_to_customer_os is null or s.synced_to_customer_os = ?) "+
			" and (s.synced_to_customer_os_attempt is null or s.synced_to_customer_os_attempt < ?) "+
			" and (s.run_id is null or s.run_id <> ?) "+
			" limit ?", 1, false, 10, runId, limit).
		Find(&companies).Error

	if err != nil {
		return nil, err
	}
	return companies, nil
}

func GetCompanyProperties(db *gorm.DB, airbyteAbId, airbyteCompaniesHashId string) (hubspotEntity.CompanyProperties, error) {
	companyProperties := hubspotEntity.CompanyProperties{}
	err := db.Table(hubspotEntity.CompanyProperties{}.TableName()).
		Where(&hubspotEntity.CompanyProperties{AirbyteAbId: airbyteAbId, AirbyteCompaniesHashid: airbyteCompaniesHashId}).
		First(&companyProperties).Error
	return companyProperties, err
}

func MarkCompanyProcessed(db *gorm.DB, company hubspotEntity.Company, synced bool, runId string) error {
	syncStatusCompany := hubspotEntity.SyncStatusCompany{
		Id:                     company.Id,
		AirbyteAbId:            company.AirbyteAbId,
		AirbyteCompaniesHashid: company.AirbyteCompaniesHashid,
	}
	db.FirstOrCreate(&syncStatusCompany, syncStatusCompany)

	return db.Model(&syncStatusCompany).
		Where(&hubspotEntity.SyncStatusCompany{Id: company.Id, AirbyteAbId: company.AirbyteAbId, AirbyteCompaniesHashid: company.AirbyteCompaniesHashid}).
		Updates(hubspotEntity.SyncStatusCompany{
			SyncedToCustomerOs: synced,
			SyncedAt:           time.Now(),
			SyncAttempt:        syncStatusCompany.SyncAttempt + 1,
			RunId:              runId,
		}).Error
}

func GetOwners(db *gorm.DB, limit int, runId string) (hubspotEntity.Owners, error) {
	var owners hubspotEntity.Owners

	cte := `
		WITH UpToDateData AS (
    		SELECT row_number() OVER (PARTITION BY id ORDER BY updatedat DESC) AS row_num, *
    		FROM owners
		)`
	err := db.
		Raw(cte+" SELECT u.* FROM UpToDateData u left join openline_sync_status_owners s "+
			" on u.id = s.id and u._airbyte_ab_id = s._airbyte_ab_id and u._airbyte_owners_hashid = s._airbyte_owners_hashid "+
			" WHERE u.row_num = ? "+
			" and (s.synced_to_customer_os is null or s.synced_to_customer_os = ?) "+
			" and (s.synced_to_customer_os_attempt is null or s.synced_to_customer_os_attempt < ?) "+
			" and (s.run_id is null or s.run_id <> ?) "+
			" limit ?", 1, false, 10, runId, limit).
		Find(&owners).Error

	if err != nil {
		return nil, err
	}
	return owners, nil
}

func MarkOwnerProcessed(db *gorm.DB, owner hubspotEntity.Owner, synced bool, runId string) error {
	syncStatusOwner := hubspotEntity.SyncStatusOwner{
		Id:                  owner.Id,
		AirbyteAbId:         owner.AirbyteAbId,
		AirbyteOwnersHashid: owner.AirbyteOwnersHashid,
	}
	db.FirstOrCreate(&syncStatusOwner, syncStatusOwner)

	return db.Model(&syncStatusOwner).
		Where(&hubspotEntity.SyncStatusOwner{Id: owner.Id, AirbyteAbId: owner.AirbyteAbId, AirbyteOwnersHashid: owner.AirbyteOwnersHashid}).
		Updates(hubspotEntity.SyncStatusOwner{
			SyncedToCustomerOs: synced,
			SyncedAt:           time.Now(),
			SyncAttempt:        syncStatusOwner.SyncAttempt + 1,
			RunId:              runId,
		}).Error
}

func GetNotes(db *gorm.DB, limit int, runId string) (hubspotEntity.Notes, error) {
	var notes hubspotEntity.Notes

	cte := `
		WITH UpToDateData AS (
    		SELECT row_number() OVER (PARTITION BY id ORDER BY updatedat DESC) AS row_num, *
    		FROM engagements_notes
		)`
	err := db.
		Raw(cte+" SELECT u.* FROM UpToDateData u left join openline_sync_status_notes s "+
			" on u.id = s.id and u._airbyte_ab_id = s._airbyte_ab_id and u._airbyte_engagements_notes_hashid = s._airbyte_engagements_notes_hashid "+
			" WHERE u.row_num = ? "+
			" and (u.contacts is not null) "+
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

func GetNoteProperties(db *gorm.DB, airbyteAbId, airbyteNotesHashId string) (hubspotEntity.NoteProperties, error) {
	noteProperties := hubspotEntity.NoteProperties{}
	err := db.Table(hubspotEntity.NoteProperties{}.TableName()).
		Where(&hubspotEntity.NoteProperties{AirbyteAbId: airbyteAbId, AirbyteNotesHashid: airbyteNotesHashId}).
		First(&noteProperties).Error
	return noteProperties, err
}

func MarkNoteProcessed(db *gorm.DB, note hubspotEntity.Note, synced bool, runId string) error {
	syncStatusNote := hubspotEntity.SyncStatusNote{
		Id:                 note.Id,
		AirbyteAbId:        note.AirbyteAbId,
		AirbyteNotesHashid: note.AirbyteNotesHashid,
	}
	db.FirstOrCreate(&syncStatusNote, syncStatusNote)

	return db.Model(&syncStatusNote).
		Where(&hubspotEntity.SyncStatusNote{Id: note.Id, AirbyteAbId: note.AirbyteAbId, AirbyteNotesHashid: note.AirbyteNotesHashid}).
		Updates(hubspotEntity.SyncStatusNote{
			SyncedToCustomerOs: synced,
			SyncedAt:           time.Now(),
			SyncAttempt:        syncStatusNote.SyncAttempt + 1,
			RunId:              runId,
		}).Error
}
