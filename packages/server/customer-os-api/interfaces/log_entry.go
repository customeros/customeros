package cosapi_interfaces

import (
	"context"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
)

type LogEntryService interface {
	GetById(ctx context.Context, logEntryId string) (*entity.LogEntryEntity, error)
}
