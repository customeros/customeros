package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	postgresentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
)

type mailService struct {
	services *Services
}

type MailService interface {
	FindEmailsForUser(tenant, userId string) ([]*neo4jentity.EmailEntity, error)
	InitializeTracing(ctx context.Context, operationName string) (opentracing.Span, context.Context)
	LoadEmail(rawEmail *postgresentity.RawEmail) (EmailMessageData, error)
	ProcessEmailCheck(email *EmailMessageData) HeaderAnalysis
	ProcessSentEmail(ctx context.Context, tx *neo4j.ManagedTransaction, emailMessage *entity.EmailMessage) (*string, error)
	SendMail(ctx context.Context, emailMessage *entity.EmailMessage) error
	SyncEmail(tenant string, emailId uuid.UUID) (postgresentity.RawState, *string, error)
	SyncEmailByMessageId(tenant, usernameSource, messageId string) (postgresentity.RawState, *string, error)
	SyncEmailsForUser(tenant string, userSource string)
}

func NewMailService(services *Services) MailService {
	return &mailService{
		services: services,
	}
}
