package service

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type FinderService interface {
	FindReferencedEntityId(ctx context.Context, externalSystemId string, referencedEntity model.ReferencedEntity) (id string, label string, err error)
}
type finderService struct {
	log      logger.Logger
	services *Services
}

func NewFinderService(log logger.Logger, services *Services) FinderService {
	return &finderService{
		log:      log,
		services: services,
	}
}

func (s *finderService) FindReferencedEntityId(ctx context.Context, externalSystemId string, referencedEntity model.ReferencedEntity) (id string, label string, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FinderService.FindReferencedEntityId")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("externalSystem", externalSystemId))

	id = ""
	label = ""
	err = nil
	tenant := common.GetTenantFromContext(ctx)

	if referencedEntity.Available() {
		switch referencedEntity.(type) {
		case *model.ReferencedParticipant:
			if id == "" {
				id, err = s.services.UserService.GetIdForReferencedUser(ctx, tenant, externalSystemId, model.ReferencedUser{
					ExternalId: referencedEntity.(*model.ReferencedParticipant).ExternalId,
				})
				if id != "" {
					label = entity.NodeLabel_User
				}
			}
			// TODO implement search for contact
			if id == "" {
				id, err = s.services.OrganizationService.GetIdForReferencedOrganization(ctx, tenant, externalSystemId, model.ReferencedOrganization{
					ExternalId: referencedEntity.(*model.ReferencedParticipant).ExternalId,
				})
				if id != "" {
					label = entity.NodeLabel_Organization
				}
			}
		}
	}
	if id == "" {
		label = ""
	}
	return
}
