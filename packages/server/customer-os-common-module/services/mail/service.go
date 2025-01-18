package mail

import (
	"context"

	neo4j_repository "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/repository"
	postgres_repository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"
	"github.com/opentracing/opentracing-go"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/caches"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
)

type mailService struct {
	cache              *caches.Cache
	postgres           *postgres_repository.Repositories
	neo4j              *neo4j_repository.Repositories
	azure              interfaces.AzureService
	google             interfaces.GoogleService
	interactionSession interfaces.InteractionSessionService
	opensrs            interfaces.OpenSrsService
	contact            interfaces.ContactService
	email              interfaces.EmailService
	interactionEvent   interfaces.InteractionEventService
	org                interfaces.OrganizationService
}

func NewMailService(cache *caches.Cache, postgres *postgres_repository.Repositories, neo4j *neo4j_repository.Repositories, azure interfaces.AzureService, contact interfaces.ContactService, email interfaces.EmailService, google interfaces.GoogleService, interactionEvent interfaces.InteractionEventService, interactionSession interfaces.InteractionSessionService, opensrs interfaces.OpenSrsService, org interfaces.OrganizationService) interfaces.MailService {
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
	return utils.IsInitialized(p)
}

func (p *mailService) initializeTracing(ctx context.Context, operationName string) (opentracing.Span, context.Context) {
	span, ctx := opentracing.StartSpanFromContext(ctx, operationName)
	tracing.SetDefaultServiceSpanTags(ctx, span)
	return span, ctx
}
