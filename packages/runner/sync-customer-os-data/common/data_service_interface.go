package common

import "github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/entity"

type DataService interface {
	Refresh()
	Close()
	SourceName() string
	GetContactsForSync(batchSize int) []entity.ContactData
	MarkContactProcessed(externalId string, synced bool) error
}
