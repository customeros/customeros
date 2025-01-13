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
	google             interfaces.GoogleService
	interactionSession interfaces.InteractionSessionService
	opensrs            interfaces.OpenSrsService
	contact            interfaces.ContactService
	email              interfaces.EmailService
	interactionEvent   interfaces.InteractionEventService
	org                interfaces.OrganizationService
}

func NewMailService(cache *caches.Cache, postgres *repository.Repositories, neo4j *neoRepo.Repositories, azure interfaces.AzureService, contact interfaces.ContactService, email interfaces.EmailService, google interfaces.GoogleService, interactionEvent interfaces.InteractionEventService, interactionSession interfaces.InteractionSessionService, opensrs interfaces.OpenSrsService, org interfaces.OrganizationService) interfaces.MailService {
	return &mailService{
		cache:              cache,
		postgres:           postgres,
		neo4j:              neo4j,
		azure:              azure,
		google:             google,
		interactionSession: interactionSession,
		opensrs:            opensrs,
		contact:            contact,
		email:              email,
		interactionEvent:   interactionEvent,
		org:                org,
	}
}

func (p *mailService) SetContactService(contact interfaces.ContactService) {
	p.contact = contact
}

func (p *mailService) SetEmailService(email interfaces.EmailService) {
	p.email = email
}

func (p *mailService) SetInteractionEventService(interactionEvent interfaces.InteractionEventService) {
	p.interactionEvent = interactionEvent
}

func (p *mailService) SetOrganizationService(org interfaces.OrganizationService) {
	p.org = org
}

func (p *mailService) IsInitialized() bool {
	if p.cache == nil || p.postgres == nil || p.neo4j == nil || p.azure == nil ||
		p.google == nil || p.interactionSession == nil || p.opensrs == nil || p.contact == nil ||
		p.email == nil || p.interactionEvent == nil || p.org == nil {
		return false
	}
	return true
}

func (p *mailService) initializeTracing(ctx context.Context, operationName string) (opentracing.Span, context.Context) {
	span, ctx := opentracing.StartSpanFromContext(ctx, operationName)
	tracing.SetDefaultServiceSpanTags(ctx, span)
	return span, ctx
}
