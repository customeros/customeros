package entity

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	"time"
)

type TenantSettingsEntity struct {
	Id                   string
	LogoUrl              string
	LogoRepositoryFileId string
	DefaultCurrency      enum.Currency //Deprecated
	BaseCurrency         enum.Currency
	InvoicingEnabled     bool
	InvoicingPostpaid    bool
	CreatedAt            time.Time
	UpdatedAt            time.Time
}
