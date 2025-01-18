package interfaces

import (
	"context"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	neo4jentity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"
)

type InteractionEventService interface {
	SetEmailService(EmailService)
	IsInitialized() bool

	GetById(ctx context.Context, id string) (*neo4jentity.InteractionEventEntity, error)
	GetInteractionEventsForInteractionSessions(ctx context.Context, ids []string, loadContent bool) (*neo4jentity.InteractionEventEntities, error)
	GetInteractionEventsForMeetings(ctx context.Context, ids []string, loadContent bool) (*neo4jentity.InteractionEventEntities, error)
	GetInteractionEventsForIssues(ctx context.Context, issueIds []string, loadContent bool) (*neo4jentity.InteractionEventEntities, error)
	GetSentByParticipantsForInteractionEvents(ctx context.Context, ids []string) (*neo4jentity.InteractionEventParticipants, error)
	GetSentToParticipantsForInteractionEvents(ctx context.Context, ids []string) (*neo4jentity.InteractionEventParticipants, error)
	GetReplyToInteractionsEventForInteractionEvents(ctx context.Context, ids []string, loadContent bool) (*neo4jentity.InteractionEventEntities, error)

	Create(ctx context.Context, data *InteractionEventCreateData) (*string, error)
	CreateInTx(ctx context.Context, tx neo4j.ManagedTransaction, data *InteractionEventCreateData) (*string, error)
}

type InteractionEventCreateData struct {
	InteractionEventEntity *neo4jentity.InteractionEventEntity
	SessionIdentifier      *string
	MeetingIdentifier      *string
	RepliesTo              *string
	SentBy                 []InteractionEventParticipantData
	SentTo                 []InteractionEventParticipantData
	SentCc                 []InteractionEventParticipantData
	SentBcc                []InteractionEventParticipantData
	ExternalSystem         *neo4jentity.ExternalSystemEntity
	Source                 neo4jentity.DataSource
	SourceOfTruth          neo4jentity.DataSource
}

type InteractionEventParticipantData struct {
	Email       *string
	PhoneNumber *string
	ContactId   *string
	UserId      *string
}
