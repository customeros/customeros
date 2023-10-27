package service

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type FinderService interface {
	FindReferencedEntityId(ctx context.Context, externalSystemId string, referencedEntity model.ReferencedEntity) (id string, label string, err error)
}
type finderService struct {
	log          logger.Logger
	services     *Services
	repositories *repository.Repositories
}

func NewFinderService(log logger.Logger, repositories *repository.Repositories, services *Services) FinderService {
	return &finderService{
		log:          log,
		services:     services,
		repositories: repositories,
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
		switch r := referencedEntity.(type) {
		case *model.ReferencedInteractionSession:
			id, err = s.services.InteractionSessionService.GetIdForReferencedInteractionSession(ctx, tenant, externalSystemId, *r)
			if id != "" {
				label = entity.NodeLabel_InteractionSession
			}
		case *model.ReferencedIssue:
			id, err = s.services.IssueService.GetIdForReferencedIssue(ctx, tenant, externalSystemId, *r)
			if id != "" {
				label = entity.NodeLabel_Issue
			}
		case *model.ReferencedUser:
			id, err = s.services.UserService.GetIdForReferencedUser(ctx, tenant, externalSystemId, *r)
			if id != "" {
				label = entity.NodeLabel_User
			}
		case *model.ReferencedOrganization:
			id, err = s.services.OrganizationService.GetIdForReferencedOrganization(ctx, tenant, externalSystemId, *r)
			if id != "" {
				label = entity.NodeLabel_Organization
			}
		case *model.ReferencedJobRole:
			contactId, _ := s.services.ContactService.GetIdForReferencedContact(ctx, tenant, externalSystemId, r.ReferencedContact)
			orgId, _ := s.services.OrganizationService.GetIdForReferencedOrganization(ctx, tenant, externalSystemId, r.ReferencedOrganization)
			id, err = s.repositories.ContactRepository.GetJobRoleId(ctx, tenant, contactId, orgId)
			if id != "" {
				label = "JobRole"
			}
			if id != "" {
				label = entity.NodeLabel_JobRole
			}
		case *model.ReferencedParticipant:
			id, err = s.services.UserService.GetIdForReferencedUser(ctx, tenant, externalSystemId, model.ReferencedUser{
				ExternalId: referencedEntity.(*model.ReferencedParticipant).ExternalId,
			})
			if id != "" {
				label = entity.NodeLabel_User
			}
			if id == "" {
				id, err = s.services.ContactService.GetIdForReferencedContact(ctx, tenant, externalSystemId, model.ReferencedContact{
					ExternalId: referencedEntity.(*model.ReferencedParticipant).ExternalId,
				})
				if id != "" {
					label = entity.NodeLabel_Contact
				}
			}
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
