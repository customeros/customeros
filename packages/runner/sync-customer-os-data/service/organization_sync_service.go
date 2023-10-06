package service

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/common"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/constants"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/entity"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/logger"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/repository"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/source"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/tracing"
	comentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/repository/postgres/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"strings"
	"sync"
	"time"
)

type organizationSyncService struct {
	repositories *repository.Repositories
	services     *Services
	log          logger.Logger
}

type domains struct {
	whitelistDomains       []comentity.WhitelistDomain
	personalEmailProviders []comentity.PersonalEmailProvider
}

func NewDefaultOrganizationSyncService(repositories *repository.Repositories, services *Services, log logger.Logger) SyncService {
	return &organizationSyncService{
		repositories: repositories,
		services:     services,
		log:          log,
	}
}

func (s *organizationSyncService) Sync(ctx context.Context, dataService source.SourceDataService, syncDate time.Time, tenant, runId string, batchSize int) (int, int, int) {
	completed, failed, skipped := 0, 0, 0
	organizationSyncMutex := &sync.Mutex{}

	var controlDomains *domains

	for {
		organizations := dataService.GetDataForSync(ctx, common.ORGANIZATIONS, batchSize, runId)

		if len(organizations) == 0 {
			break
		}

		if controlDomains == nil {
			whitelistDomains, err := s.repositories.CommonRepositories.WhitelistDomainRepository.GetWhitelistDomains(tenant)
			if err != nil {
				s.log.Errorf("error while getting whitelist domains: %v", err)
				whitelistDomains = make([]comentity.WhitelistDomain, 0)
			}
			personalEmailProviders, err := s.repositories.CommonRepositories.PersonalEmailProviderRepository.GetPersonalEmailProviders()
			if err != nil {
				s.log.Errorf("error while getting personal email providers: %v", err)
				personalEmailProviders = make([]comentity.PersonalEmailProvider, 0)
			}
			controlDomains = &domains{
				whitelistDomains:       whitelistDomains,
				personalEmailProviders: personalEmailProviders,
			}
		}

		s.log.Infof("syncing %d organizations from %s for tenant %s", len(organizations), dataService.SourceId(), tenant)

		var wg sync.WaitGroup
		wg.Add(len(organizations))

		results := make(chan result, len(organizations))
		done := make(chan struct{})

		for _, v := range organizations {
			v := v

			go func(org entity.OrganizationData) {
				defer wg.Done()

				var comp, fail, skip int
				s.syncOrganization(ctx, organizationSyncMutex, org, dataService, controlDomains, syncDate, tenant, runId, &comp, &fail, &skip)

				results <- result{comp, fail, skip}
			}(v.(entity.OrganizationData))
		}
		// Wait for goroutines to finish
		go func() {
			wg.Wait()
			close(done)
		}()
		go func() {
			<-done
			close(results)
		}()

		for r := range results {
			completed += r.completed
			failed += r.failed
			skipped += r.skipped
		}

		if len(organizations) < batchSize {
			break
		}

	}

	return completed, failed, skipped
}

func (s *organizationSyncService) syncOrganization(ctx context.Context, organizationSyncMutex *sync.Mutex, orgInput entity.OrganizationData, dataService source.SourceDataService, controlDomains *domains, syncDate time.Time, tenant, runId string, completed, failed, skipped *int) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationSyncService.syncOrganization")
	defer span.Finish()
	tracing.SetDefaultSyncServiceSpanTags(ctx, span)

	var failedSync = false
	var reason string
	orgInput.Normalize()

	if orgInput.ExternalSystem == "" {
		_ = dataService.MarkProcessed(ctx, orgInput.SyncId, runId, false, false, "External system is empty. Error during reading data from source")
		*failed++
		return
	}

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

	// populate domain if missing and website available
	if !orgInput.HasDomains() && orgInput.Website != "" {
		domainFromWebsite := utils.ExtractDomainFromUrl(orgInput.Website)
		if domainFromWebsite != "" {
			orgInput.Domains = []string{domainFromWebsite}
		}
	}

	nonPersonalEmailProviderDomains := make([]string, 0)
	for _, domain := range orgInput.Domains {
		if !controlDomains.isPersonalEmailProvider(domain) {
			nonPersonalEmailProviderDomains = append(nonPersonalEmailProviderDomains, domain)
		}
	}
	orgInput.Domains = nonPersonalEmailProviderDomains

	if orgInput.DomainRequired {
		if !orgInput.HasDomains() {
			if err := dataService.MarkProcessed(ctx, orgInput.SyncId, runId, true, true, "Missing non-personal email provider domain"); err != nil {
				*failed++
				span.LogFields(log.Bool("failedSync", true))
				return
			}
			*skipped++
			span.LogFields(log.Bool("skippedSync", true))
			return
		}
	}
	orgHasWhitelistedDomain := false
	for _, domain := range orgInput.Domains {
		if controlDomains.isWhitelistedDomain(domain) {
			orgHasWhitelistedDomain = true
		}
	}

	organizationSyncMutex.Lock()
	organizationId, err := s.repositories.OrganizationRepository.GetMatchedOrganizationId(ctx, tenant, orgInput.ExternalSystem, orgInput.ExternalId, orgInput.Domains)
	if err != nil {
		failedSync = true
		tracing.TraceErr(span, err)
		reason = fmt.Sprintf("failed finding existing matched organization with external reference %v for tenant %v :%v", orgInput.ExternalId, tenant, err)
		s.log.Errorf(reason)
	}

	if !failedSync && orgInput.UpdateOnly && organizationId == "" {
		organizationSyncMutex.Unlock()
		if err := dataService.MarkProcessed(ctx, orgInput.SyncId, runId, true, true, "This record is for update only, organization not available yet."); err != nil {
			*failed++
			span.LogFields(log.Bool("failedSync", true))
			return
		}
		*skipped++
		span.LogFields(log.Bool("skippedSync", true))
		return
	}

	newOrganization := len(organizationId) == 0
	// Create new organization id if not found
	if newOrganization {
		orgUuid, _ := uuid.NewRandom()
		organizationId = orgUuid.String()
	}
	orgInput.Id = organizationId
	span.LogFields(log.String("organizationId", organizationId))

	if !failedSync {
		err = s.repositories.OrganizationRepository.MergeOrganization(ctx, tenant, syncDate, orgInput, orgHasWhitelistedDomain || orgInput.Whitelisted)
		if err != nil {
			failedSync = true
			tracing.TraceErr(span, err)
			reason = fmt.Sprintf("failed merge organization with external reference %v for tenant %v :%v", orgInput.ExternalId, tenant, err)
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
	organizationSyncMutex.Unlock()

	if newOrganization && !failedSync {
		err := s.repositories.ActionRepository.OrganizationCreatedAction(ctx, tenant, orgInput.Id, orgInput.ExternalSystem, constants.AppSourceSyncCustomerOsData)
		if err != nil {
			failedSync = true
			tracing.TraceErr(span, err)
			reason = fmt.Sprintf("failed create organization created action for organization %v, tenant %v :%v", organizationId, tenant, err)
			s.log.Errorf(reason)
		}
	}

	if orgInput.HasPhoneNumbers() && !failedSync {
		for _, phoneNumber := range orgInput.PhoneNumbers {
			if err = s.repositories.OrganizationRepository.MergePhoneNumber(ctx, tenant, organizationId, orgInput.ExternalSystem, *orgInput.CreatedAt, phoneNumber); err != nil {
				failedSync = true
				tracing.TraceErr(span, err)
				reason = fmt.Sprintf("failed merge phone number for organization with external reference %v , tenant %v :%v", orgInput.ExternalId, tenant, err)
				s.log.Errorf(reason)
			}
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

	if orgInput.HasNotes() && !failedSync {
		for _, note := range orgInput.Notes {
			localNote := entity.NoteData{
				BaseData: entity.BaseData{
					CreatedAt:      orgInput.CreatedAt,
					ExternalId:     string(note.FieldSource) + "-" + orgInput.ExternalId,
					ExternalSystem: orgInput.ExternalSystem,
				},
				Content: note.Note,
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
			reason = fmt.Sprintf("failed link current organization as subsidiary %v to parent organization by external id %v, tenant %v :%v", orgInput.Id, orgInput.ParentOrganization.Organization.ExternalId, tenant, err)
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

func (d *domains) isPersonalEmailProvider(domain string) bool {
	for _, v := range d.personalEmailProviders {
		if strings.ToLower(domain) == strings.ToLower(v.ProviderDomain) {
			return true
		}
	}
	return false
}

func (d *domains) isWhitelistedDomain(domain string) bool {
	if domain == "" {
		return false
	}
	for _, v := range d.whitelistDomains {
		if v.Domain != "*" && strings.ToLower(domain) == strings.ToLower(v.Domain) && v.Allowed {
			return true
		}
	}
	return false
}
