package commands

import (
	common_models "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/models"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"time"
)

type CreateJobRoleCommand struct {
	eventstore.BaseCommand
	StartedAt   *time.Time
	EndedAt     *time.Time
	JobTitle    string
	Description *string
	Primary     bool
	Source      common_models.Source
	CreatedAt   *time.Time
	UpdatedAt   *time.Time
}
