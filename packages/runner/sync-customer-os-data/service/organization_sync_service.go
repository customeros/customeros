package service

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/common"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/entity"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/logger"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/repository"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/source"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"strings"
	"time"
)

type organizationSyncService struct {
	repositories *repository.Repositories
	services     *Services
	log          logger.Logger
}

func NewDefaultOrganizationSyncService(repositories *repository.Repositories, services *Services, log logger.Logger) SyncService {
	return &organizationSyncService{
		repositories: repositories,
		services:     services,
		log:          log,
	}
}

func (s *organizationSyncService) Sync(ctx context.Context, dataService source.SourceDataService, syncDate time.Time, tenant, runId string, batchSize int) (int, int, int) {
	span, ctx := tracing.StartTracerSpan(ctx, "OrganizationSyncService.Sync")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	completed, failed, skipped := 0, 0, 0
	for {
		organizations := dataService.GetDataForSync(ctx, common.ORGANIZATIONS, batchSize, runId)
		if len(organizations) == 0 {
			s.log.Debugf("no organizations found for sync from %s for tenant %s", dataService.SourceId(), tenant)
			break
		}
		s.log.Infof("syncing %d organizations from %s for tenant %s", len(organizations), dataService.SourceId(), tenant)

		for _, v := range organizations {
			s.syncOrganization(ctx, v.(entity.OrganizationData), dataService, syncDate, tenant, runId, &completed, &failed, &skipped)
		}
		if len(organizations) < batchSize {
			break
		}
	}
	return completed, failed, skipped
}

func (s *organizationSyncService) syncOrganization(ctx context.Context, orgInput entity.OrganizationData, dataService source.SourceDataService, syncDate time.Time, tenant, runId string, completed, failed, skipped *int) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationSyncService.syncOrganization")
	defer span.Finish()
	tracing.SetDefaultSyncServiceSpanTags(ctx, span)

	var failedSync = false
	var reason string
	orgInput.Normalize()

	if orgInput.Skip {
		if err := dataService.MarkProcessed(ctx, orgInput.SyncId, runId, true, true, orgInput.SkipReason); err != nil {
			*failed++
			span.LogFields(log.Bool("failedSync", true))
			return
		}
		*skipped++
		span.LogFields(log.Bool("skippedSync", true))
		return
	}

	organizationId, err := s.repositories.OrganizationRepository.GetMatchedOrganizationId(ctx, tenant, orgInput)
	if err != nil {
		failedSync = true
		tracing.TraceErr(span, err)
		reason = fmt.Sprintf("failed finding existing matched organization with external reference %v for tenant %v :%v", orgInput.ExternalId, tenant, err)
		s.log.Errorf(reason)
	}

	newOrganization := len(organizationId) == 0
	// Create new organization id if not found
	if newOrganization {
		orgUuid, _ := uuid.NewRandom()
		organizationId = orgUuid.String()
	}
	orgInput.Id = organizationId

	if !failedSync {
		err = s.repositories.OrganizationRepository.MergeOrganization(ctx, tenant, syncDate, orgInput)
		if err != nil {
			failedSync = true
			tracing.TraceErr(span, err)
			reason = fmt.Sprintf("failed merge organization with external reference %v for tenant %v :%v", orgInput.ExternalId, tenant, err)
			s.log.Errorf(reason)
		}
	}

	if newOrganization && !failedSync {
		err := s.repositories.ActionRepository.OrganizationCreatedAction(ctx, tenant, orgInput.Id, orgInput.ExternalSystem, orgInput.ExternalSystem)
		if err != nil {
			failedSync = true
			tracing.TraceErr(span, err)
			reason = fmt.Sprintf("failed create organization created action for organization %v, tenant %v :%v", organizationId, tenant, err)
			s.log.Errorf(reason)
		}
	}

	if orgInput.HasDomains() && !failedSync {
		for _, domain := range orgInput.Domains {
			err = s.repositories.OrganizationRepository.MergeOrganizationDomain(ctx, tenant, organizationId, domain, orgInput.ExternalSystem)
			if err != nil {
				failedSync = true
				tracing.TraceErr(span, err)
				reason = fmt.Sprintf("failed merge organization domain for organization %v, tenant %v :%v", organizationId, tenant, err)
				s.log.Errorf(reason)
				break
			}
		}
	}

	if orgInput.HasPhoneNumber() && !failedSync {
		if err = s.repositories.OrganizationRepository.MergePhoneNumber(ctx, tenant, organizationId, orgInput.PhoneNumber, orgInput.ExternalSystem, *orgInput.CreatedAt); err != nil {
			failedSync = true
			tracing.TraceErr(span, err)
			reason = fmt.Sprintf("failed merge phone number for organization with external reference %v , tenant %v :%v", orgInput.ExternalId, tenant, err)
			s.log.Errorf(reason)
		}
	}

	if orgInput.HasEmail() && !failedSync {
		orgInput.Email = strings.ToLower(orgInput.Email)
		if err = s.repositories.OrganizationRepository.MergeEmail(ctx, tenant, organizationId, orgInput.Email, orgInput.ExternalSystem, *orgInput.CreatedAt); err != nil {
			failedSync = true
			tracing.TraceErr(span, err)
			reason = fmt.Sprintf("failed merge email for organization with external reference %v , tenant %v :%v", orgInput.ExternalId, tenant, err)
			s.log.Errorf(reason)
		}
	}

	if orgInput.HasLocation() && !failedSync {
		err = s.repositories.OrganizationRepository.MergeOrganizationLocation(ctx, tenant, organizationId, orgInput)
		if err != nil {
			failedSync = true
			tracing.TraceErr(span, err)
			reason = fmt.Sprintf("failed merge organization' location with external reference %v for tenant %v :%v", orgInput.ExternalId, tenant, err)
			s.log.Errorf(reason)
		}
	}

	if orgInput.HasOwner() && !failedSync {
		if err = s.repositories.OrganizationRepository.SetOwner(ctx, tenant, organizationId, orgInput.UserExternalOwnerId, dataService.SourceId()); err != nil {
			// Do not mark sync as failed in case owner relationship is not set
			s.log.Errorf("failed set owner user for organization %v, tenant %v :%v", organizationId, tenant, err)
		}
	}

	if orgInput.HasNotes() && !failedSync {
		for _, note := range orgInput.Notes {
			localNote := entity.NoteData{
				BaseData: entity.BaseData{
					CreatedAt:      orgInput.CreatedAt,
					ExternalId:     string(note.FieldSource) + "-" + orgInput.ExternalId,
					ExternalSystem: orgInput.ExternalSystem,
				},
				Html: note.Note,
			}
			noteId, err := s.repositories.NoteRepository.GetMatchedNoteId(ctx, tenant, localNote)
			if err != nil {
				failedSync = true
				tracing.TraceErr(span, err)
				reason = fmt.Sprintf("failed finding existing matched note with external reference id %v for tenant %v :%v", localNote.ExternalId, tenant, err)
				s.log.Errorf(reason)
				break
			}
			// Create new note id if not found
			if len(noteId) == 0 {
				noteUuid, _ := uuid.NewRandom()
				noteId = noteUuid.String()
			}
			localNote.Id = noteId
			err = s.repositories.NoteRepository.MergeNote(ctx, tenant, syncDate, localNote)
			if err != nil {
				failedSync = true
				tracing.TraceErr(span, err)
				reason = fmt.Sprintf("failed merge organization note for organization %v, tenant %v :%v", organizationId, tenant, err)
				s.log.Errorf(reason)
				break
			}
			err = s.repositories.NoteRepository.NoteLinkWithOrganizationByExternalId(ctx, tenant, noteId, orgInput.ExternalId, orgInput.ExternalSystem)
			if err != nil {
				failedSync = true
				tracing.TraceErr(span, err)
				reason = fmt.Sprintf("failed link note with organization %v, tenant %v :%v", organizationId, tenant, err)
				s.log.Errorf(reason)
				break
			}
		}
	}

	if orgInput.HasRelationship() && !failedSync {
		err = s.repositories.OrganizationRepository.MergeOrganizationRelationshipAndStage(ctx, tenant, organizationId, orgInput.RelationshipName, orgInput.RelationshipStage, orgInput.ExternalSystem)
		if err != nil {
			failedSync = true
			tracing.TraceErr(span, err)
			reason = fmt.Sprintf("failed merge organization relationship for organization %v, tenant %v :%v", organizationId, tenant, err)
			s.log.Errorf(reason)
		}
	}

	if orgInput.IsSubsidiary() && !failedSync {
		if err = s.repositories.OrganizationRepository.LinkToParentOrganizationAsSubsidiary(ctx, tenant, organizationId, orgInput.ExternalSystem, orgInput.ParentOrganization); err != nil {
			failedSync = true
			tracing.TraceErr(span, err)
			reason = fmt.Sprintf("failed link current organization as subsidiary %v to parent organization by external id %v, tenant %v :%v", orgInput.Id, orgInput.ParentOrganization.ExternalId, tenant, err)
			s.log.Errorf(reason)
		}
	}

	s.services.OrganizationService.UpdateLastTouchpointByOrganizationId(ctx, tenant, organizationId)

	s.log.Debugf("successfully merged organization with id %v for tenant %v from %v", organizationId, tenant, dataService.SourceId())
	if err := dataService.MarkProcessed(ctx, orgInput.SyncId, runId, failedSync == false, false, reason); err != nil {
		*failed++
		span.LogFields(log.Bool("failedSync", true))
		return
	}
	if failedSync == true {
		*failed++
	} else {
		*completed++
	}
	span.LogFields(log.Bool("failedSync", failedSync))
}
