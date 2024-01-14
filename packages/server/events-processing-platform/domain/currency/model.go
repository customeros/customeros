package currency

import (
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"time"
)

type Currency struct {
	ID           string             `json:"id"`
	CreatedAt    time.Time          `json:"createdAt"`
	UpdatedAt    time.Time          `json:"updatedAt"`
	SourceFields commonmodel.Source `json:"source"`

	Name   string `json:"name"`
	Symbol string `json:"symbol"`
}
