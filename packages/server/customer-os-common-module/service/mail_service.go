package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	postgresentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
)

type mailService struct {
	services *Services
}

type MailService interface {
	ExtractEmails(s string) []string
	GetEmailsForProcessingForUser(ctx context.Context, tenant, userEmailAddress string)
	LoadEmail(ctx context.Context, rawEmail *postgresentity.RawEmail) (EmailMessageData, error)
	ProcessEmailCheck(ctx context.Context, tenant string, email *EmailMessageData) HeaderAnalysis
	ProcessEmail(ctx context.Context, tenant string, emailId uuid.UUID) entity.UpdateRawEmailTable
	ProcessEmailByMessageId(ctx context.Context, tenant, usernameSource, messageId string) entity.UpdateRawEmailTable
	SendMail(ctx context.Context, emailMessage *entity.EmailMessage) error
	ProcessSentEmail(ctx context.Context, tx *neo4j.ManagedTransaction, emailMessage *entity.EmailMessage) (*string, error)
	GetEmailIdForEmail(ctx context.Context, txWithPostCommit *utils.TxWithPostCommit, tenant, email string, source string) (string, error)
}

func NewMailService(services *Services) MailService {
	return &mailService{
		services: services,
	}
}

func (p *mailService) initializeTracing(ctx context.Context, operationName string) (opentracing.Span, context.Context) {
	span, ctx := opentracing.StartSpanFromContext(ctx, operationName)
	tracing.SetDefaultServiceSpanTags(ctx, span)
	return span, ctx
}
