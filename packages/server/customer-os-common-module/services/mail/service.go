package mail

import (
	"context"

	neoRepo "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/repository"
	"github.com/opentracing/opentracing-go"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/caches"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/interfaces"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
)

type mailService struct {
	cache              *caches.Cache
	postgres           *repository.Repositories
	neo4j              *neoRepo.Repositories
	azure              interfaces.AzureService
	contact            interfaces.ContactService
	email              interfaces.EmailService
	google             interfaces.GoogleService
	interactionEvent   interfaces.InteractionEventService
	interactionSession interfaces.InteractionSessionService
	opensrs            interfaces.OpenSrsService
	org                interfaces.OrganizationService
}

func NewMailService(cache *caches.Cache, postgres *repository.Repositories, neo4j *neoRepo.Repositories, azure interfaces.AzureService, contact interfaces.ContactService, email interfaces.EmailService, google interfaces.GoogleService, interactionEvent interfaces.InteractionEventService, interactionSession interfaces.InteractionSessionService, opensrs interfaces.OpenSrsService, org interfaces.OrganizationService) MailService {
	return &mailService{
		cache:              cache,
		postgres:           postgres,
		neo4j:              neo4j,
		azure:              azure,
		contact:            contact,
		email:              email,
		google:             google,
		interactionEvent:   interactionEvent,
		interactionSession: interactionSession,
		opensrs:            opensrs,
		org:                org,
	}
}

func (p *mailService) initializeTracing(ctx context.Context, operationName string) (opentracing.Span, context.Context) {
	span, ctx := opentracing.StartSpanFromContext(ctx, operationName)
	tracing.SetDefaultServiceSpanTags(ctx, span)
	return span, ctx
}
