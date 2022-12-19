package common

import "github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/entity"

type DataService interface {
	Refresh()
	Close()
	SourceId() string
	GetContactsForSync(batchSize int) []entity.ContactData
	GetCompaniesForSync(batchSize int) []entity.CompanyData
	MarkContactProcessed(externalId string, synced bool) error
	MarkCompanyProcessed(externalId string, synced bool) error
}
