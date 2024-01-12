package entity

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	"time"
)

type TenantSettingsEntity struct {
	LogoUrl          string
	DefaultCurrency  enum.Currency
	InvoicingEnabled bool
	CreatedAt        time.Time
	UpdatedAt        time.Time
}
