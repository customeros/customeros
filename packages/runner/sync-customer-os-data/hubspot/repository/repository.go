package repository

import (
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/hubspot/entity"
	"gorm.io/gorm"
)

const RecordsLimit = 100

func GetContacts(db *gorm.DB) (entity.Contacts, error) {
	var contacts entity.Contacts

	cte := `
		WITH UpToDateData AS (
    		SELECT row_number() OVER (PARTITION BY id ORDER BY updatedat DESC) AS row_num, *
    		FROM contacts
		)`
	err := db.
		Raw(cte+"SELECT * FROM UpToDateData WHERE row_num = ? and (synced_to_customer_os is null or synced_to_customer_os = ?) "+
			" limit ?", 1, false, RecordsLimit).
		Find(&contacts).Error

	if err != nil {
		return nil, err
	}
	return contacts, nil
}

func MarkSynced(db *gorm.DB, contact entity.Contact) error {
	return nil
}

// Query the "Latest" CTE and select only rows where the row_number is 1
// and the died column is either NULL or FALSE.
