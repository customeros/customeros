package common

import "github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/entity"

type DataService interface {
	GetContactsForSync(batchSize int) []entity.ContactEntity
	Refresh()
}
