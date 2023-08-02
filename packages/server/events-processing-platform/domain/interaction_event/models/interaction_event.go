package models

import (
	"fmt"
	common_models "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/models"
	"time"
)

type InteractionEvent struct {
	ID        string               `json:"id"`
	Source    common_models.Source `json:"source"`
	CreatedAt time.Time            `json:"createdAt"`
	UpdatedAt time.Time            `json:"updatedAt"`
	Summary   string               `json:"summary"`
}

func (i *InteractionEvent) String() string {
	return fmt.Sprintf("ID: %s, Source: %s, CreatedAt: %s, UpdatedAt: %s", i.ID, i.Source, i.CreatedAt, i.UpdatedAt)
}
