package service

import (
	"context"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
	model2 "github.com/customeros/customeros/packages/server/customer-os-common-module/model"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"github.com/customeros/customeros/packages/server/customer-os-webhooks/model"
	"github.com/customeros/customeros/packages/server/customer-os-webhooks/repository"
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
	spans, ctx := telemetry.StartServiceSpan(ctx, "FinderService.FindReferencedEntityId")
	defer spans.Finish()
	spans.TagString(telemetry.SpanTagExternalSystem, externalSystemId)
	spans.LogObjectAsJson("referencedEntity", referencedEntity)

	id = ""
	label = ""
	err = nil
	tenant := common.GetTenantFromContext(ctx)

	if referencedEntity.Available() {
		switch r := referencedEntity.(type) {
		case *model.ReferencedInteractionSession:
			id, err = s.services.InteractionSessionService.GetIdForReferencedInteractionSession(ctx, tenant, externalSystemId, *r)
			if id != "" {
				label = model2.NodeLabelInteractionSession
			}
		case *model.ReferencedIssue:
			id, err = s.services.IssueService.GetIdForReferencedIssue(ctx, tenant, externalSystemId, *r)
			if id != "" {
				label = model2.NodeLabelIssue
			}
		case *model.ReferencedUser:
			id, err = s.services.UserService.GetIdForReferencedUser(ctx, tenant, externalSystemId, *r)
			if id != "" {
				label = model2.NodeLabelUser
			}
		case *model.ReferencedOrganization:
			id, err = s.services.OrganizationService.GetIdForReferencedOrganization(ctx, tenant, externalSystemId, *r)
			if id != "" {
				label = model2.NodeLabelOrganization
			}
		case *model.ReferencedJobRole:
			contactId, _ := s.services.ContactService.GetIdForReferencedContact(ctx, tenant, externalSystemId, r.ReferencedContact)
			orgId, _ := s.services.OrganizationService.GetIdForReferencedOrganization(ctx, tenant, externalSystemId, r.ReferencedOrganization)
			id, err = s.repositories.ContactRepository.GetJobRoleId(ctx, tenant, contactId, orgId)
			if id != "" {
				label = "JobRole"
			}
			if id != "" {
				label = model2.NodeLabelJobRole
			}
		case *model.ReferencedParticipant:
			id, err = s.services.UserService.GetIdForReferencedUser(ctx, tenant, externalSystemId, model.ReferencedUser{
				ExternalId: referencedEntity.(*model.ReferencedParticipant).ExternalId,
			})
			if id != "" {
				label = model2.NodeLabelUser
			}
			if id == "" {
				id, err = s.services.ContactService.GetIdForReferencedContact(ctx, tenant, externalSystemId, model.ReferencedContact{
					ExternalId: referencedEntity.(*model.ReferencedParticipant).ExternalId,
				})
				if id != "" {
					label = model2.NodeLabelContact
				}
			}
			if id == "" {
				id, err = s.services.OrganizationService.GetIdForReferencedOrganization(ctx, tenant, externalSystemId, model.ReferencedOrganization{
					ExternalId: referencedEntity.(*model.ReferencedParticipant).ExternalId,
				})
				if id != "" {
					label = model2.NodeLabelOrganization
				}
			}
		}
	}
	if id == "" {
		label = ""
	}
	return
}
