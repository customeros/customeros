package interfaces

import (
	"context"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
)

type InteractionSessionService interface {
	GetById(ctx context.Context, id string) (*entity.InteractionSessionEntity, error)
	GetAttendedByParticipantsForInteractionSessions(ctx context.Context, ids []string) (*entity.InteractionSessionParticipants, error)
	GetInteractionSessionsForInteractionEvents(ctx context.Context, ids []string) (*entity.InteractionSessionEntities, error)

	Create(ctx context.Context, data *entity.InteractionSessionEntity) (*string, error)
	CreateInTx(ctx context.Context, tx neo4j.ManagedTransaction, data *entity.InteractionSessionEntity) (*string, error)
}
