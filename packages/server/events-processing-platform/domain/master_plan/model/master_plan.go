package model

import (
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"time"
)

type MasterPlan struct {
	ID           string             `json:"id"`
	Name         string             `json:"name"`
	CreatedAt    time.Time          `json:"createdAt"`
	SourceFields commonmodel.Source `json:"source"`
}
