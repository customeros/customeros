package cosapi_interfaces

import (
	"context"

	neo4j_entity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"
)

type LogEntryService interface {
	GetById(ctx context.Context, logEntryId string) (*neo4j_entity.LogEntryEntity, error)
}
