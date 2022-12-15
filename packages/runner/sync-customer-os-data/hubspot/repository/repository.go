package repository

import (
	hubspotEntity "github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/hubspot/entity"
	"gorm.io/gorm"
	"time"
)

func GetContacts(db *gorm.DB, limit int) (hubspotEntity.Contacts, error) {
	var contacts hubspotEntity.Contacts

	cte := `
		WITH UpToDateData AS (
    		SELECT row_number() OVER (PARTITION BY id ORDER BY updatedat DESC) AS row_num, *
    		FROM contacts
		)`
	err := db.
		Raw(cte+" SELECT * FROM UpToDateData "+
			" WHERE row_num = ? "+
			" and (synced_to_customer_os is null or synced_to_customer_os = ?) "+
			" and (synced_to_customer_os_attempt is null or synced_to_customer_os_attempt < ?) "+
			" limit ?", 1, false, 10, limit).
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

func MarkContactProcessed(db *gorm.DB, contact hubspotEntity.Contact, synced bool) error {
	return db.Model(&contact).
		Where(&hubspotEntity.ContactProperties{AirbyteAbId: contact.AirbyteAbId, AirbyteContactsHashid: contact.AirbyteContactsHashid}).
		Updates(hubspotEntity.Contact{
			SyncedToCustomerOs: synced,
			SyncedAt:           time.Now(),
			SyncAttempt:        contact.SyncAttempt + 1,
		}).
		Error
}
