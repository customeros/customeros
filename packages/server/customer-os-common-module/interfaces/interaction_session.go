package interfaces

import (
	"context"

	neo4j_entity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type InteractionSessionService interface {
	GetById(ctx context.Context, id string) (*neo4j_entity.InteractionSessionEntity, error)
	GetAttendedByParticipantsForInteractionSessions(ctx context.Context, ids []string) (*neo4j_entity.InteractionSessionParticipants, error)
	GetInteractionSessionsForInteractionEvents(ctx context.Context, ids []string) (*neo4j_entity.InteractionSessionEntities, error)

	Create(ctx context.Context, data *neo4j_entity.InteractionSessionEntity) (*string, error)
	CreateInTx(ctx context.Context, tx neo4j.ManagedTransaction, data *neo4j_entity.InteractionSessionEntity) (*string, error)
}
