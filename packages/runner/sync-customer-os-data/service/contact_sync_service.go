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
	"sync"
	"time"
)

type contactSyncService struct {
	repositories *repository.Repositories
	services     *Services
	log          logger.Logger
}

func NewDefaultContactSyncService(repositories *repository.Repositories, services *Services, log logger.Logger) SyncService {
	return &contactSyncService{
		repositories: repositories,
		services:     services,
		log:          log,
	}
}

func (s *contactSyncService) Sync(ctx context.Context, dataService source.SourceDataService, syncDate time.Time, tenant, runId string, batchSize int) (int, int, int) {

	completed, failed, skipped := 0, 0, 0

	for {

		contacts := dataService.GetDataForSync(ctx, common.CONTACTS, batchSize, runId)

		if len(contacts) == 0 {
			break
		}

		s.log.Infof("Syncing %d contacts", len(contacts))

		var wg sync.WaitGroup
		wg.Add(len(contacts))

		results := make(chan result, len(contacts))
		done := make(chan struct{})

		for _, v := range contacts {
			v := v

			go func(contact entity.ContactData) {
				defer wg.Done()

				var comp, fail, skip int
				s.syncContact(ctx, contact, dataService, syncDate, tenant, runId, &comp, &fail, &skip)

				results <- result{comp, fail, skip}
			}(v.(entity.ContactData))
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

		if len(contacts) < batchSize {
			break
		}

	}

	return completed, failed, skipped
}

func (s *contactSyncService) syncContact(ctx context.Context, contactInput entity.ContactData, dataService source.SourceDataService, syncDate time.Time, tenant, runId string, completed, failed, skipped *int) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactSyncService.syncContact")
	defer span.Finish()
	tracing.SetDefaultSyncServiceSpanTags(ctx, span)

	var failedSync = false
	var reason string
	contactInput.Normalize()

	if contactInput.Skip {
		if err := dataService.MarkProcessed(ctx, contactInput.SyncId, runId, true, true, contactInput.SkipReason); err != nil {
			*failed++
			span.LogFields(log.Bool("failedSync", true))
			return
		}
		*skipped++
		span.LogFields(log.Bool("skippedSync", true))
		return
	}

	if contactInput.OrganizationRequired {
		found := false
		for _, org := range contactInput.Organizations {
			orgId, _ := s.services.OrganizationService.GetIdForReferencedOrganization(ctx, tenant, contactInput.ExternalSystem, org)
			if orgId != "" {
				found = true
				break
			}
		}
		if !found {
			if err := dataService.MarkProcessed(ctx, contactInput.SyncId, runId, true, true, "Organization not found"); err != nil {
				*failed++
				span.LogFields(log.Bool("failedSync", true))
				return
			}
			*skipped++
			span.LogFields(log.Bool("skippedSync", true))
			return
		}
	}

	if contactInput.Name == "" {
		contactInput.Name = strings.TrimSpace(fmt.Sprintf("%s %s", contactInput.FirstName, contactInput.LastName))
	}

	contactId, err := s.repositories.ContactRepository.GetMatchedContactId(ctx, tenant, contactInput)
	if err != nil {
		failedSync = true
		tracing.TraceErr(span, err)
		reason = fmt.Sprintf("failed finding existing matched contact with external reference %v for tenant %v :%v", contactInput.ExternalId, tenant, err)
		s.log.Errorf(reason)
	}

	// Create new contact id if not found
	if len(contactId) == 0 {
		contactUuid, _ := uuid.NewRandom()
		contactId = contactUuid.String()
	}
	contactInput.Id = contactId
	span.LogFields(log.String("contactId", contactId))

	if !failedSync {
		err = s.repositories.ContactRepository.MergeContact(ctx, tenant, syncDate, contactInput)
		if err != nil {
			failedSync = true
			tracing.TraceErr(span, err)
			reason = fmt.Sprintf("failed merge contact with external reference %v for tenant %v :%v", contactInput.ExternalId, tenant, err)
			s.log.Errorf(reason)
		}
	}

	if len(contactInput.Email) > 0 && !failedSync {
		if err = s.repositories.ContactRepository.MergePrimaryEmail(ctx, tenant, contactId, contactInput.Email, contactInput.ExternalSystem, *contactInput.CreatedAt); err != nil {
			failedSync = true
			tracing.TraceErr(span, err)
			reason = fmt.Sprintf("failed merge primary email for contact with external reference %v , tenant %v :%v", contactInput.ExternalId, tenant, err)
			s.log.Errorf(reason)
		}
	}

	if !failedSync {
		for _, additionalEmail := range contactInput.AdditionalEmails {
			if len(additionalEmail) > 0 {
				if err = s.repositories.ContactRepository.MergeAdditionalEmail(ctx, tenant, contactId, additionalEmail, contactInput.ExternalSystem, *contactInput.CreatedAt); err != nil {
					failedSync = true
					tracing.TraceErr(span, err)
					reason = fmt.Sprintf("failed merge additional email for contact with external reference %v , tenant %v :%v", contactInput.ExternalId, tenant, err)
					s.log.Errorf(reason)
					break
				}
			}
		}
	}

	if contactInput.HasPhoneNumber() && !failedSync {
		if err = s.repositories.ContactRepository.MergePrimaryPhoneNumber(ctx, tenant, contactId, contactInput.PhoneNumber, contactInput.ExternalSystem, *contactInput.CreatedAt); err != nil {
			failedSync = true
			tracing.TraceErr(span, err)
			reason = fmt.Sprintf("failed merge primary phone number for contact with external reference %v , tenant %v :%v", contactInput.ExternalId, tenant, err)
			s.log.Errorf(reason)
		}
	}

	if !failedSync {
		for _, additionalPhoneNumber := range contactInput.AdditionalPhoneNumbers {
			if len(additionalPhoneNumber) > 0 {
				if err = s.repositories.ContactRepository.MergeAdditionalPhoneNumber(ctx, tenant, contactId, additionalPhoneNumber, contactInput.ExternalSystem, *contactInput.CreatedAt); err != nil {
					failedSync = true
					tracing.TraceErr(span, err)
					reason = fmt.Sprintf("failed merge additional phone number for contact with external reference %v , tenant %v :%v", contactInput.ExternalId, tenant, err)
					s.log.Errorf(reason)
					break
				}
			}
		}
	}

	if !failedSync {
		for _, organization := range contactInput.Organizations {
			if !organization.Available() {
				continue
			}
			if organization.ReferencedByExternalId() {
				if err = s.repositories.ContactRepository.LinkContactWithOrganizationByExternalId(ctx, tenant, contactId, dataService.SourceId(), organization); err != nil {
					failedSync = true
					tracing.TraceErr(span, err)
					reason = fmt.Sprintf("failed link contact %v to organization with external id %v, tenant %v :%v", contactId, organization.ExternalId, tenant, err)
					s.log.Errorf(reason)
					break
				}
			} else if organization.ReferencedById() {
				if err = s.repositories.ContactRepository.LinkContactWithOrganizationByInternalId(ctx, tenant, contactId, dataService.SourceId(), organization); err != nil {
					failedSync = true
					tracing.TraceErr(span, err)
					reason = fmt.Sprintf("failed link contact %v to organization with id %v, tenant %v :%v", contactId, organization.Id, tenant, err)
					s.log.Errorf(reason)
				}
			} else if organization.ReferencedByDomain() {
				if err = s.repositories.ContactRepository.LinkContactWithOrganizationByDomain(ctx, tenant, contactId, dataService.SourceId(), organization); err != nil {
					failedSync = true
					tracing.TraceErr(span, err)
					reason = fmt.Sprintf("failed link contact %v to organization with id %v, tenant %v :%v", contactId, organization.Id, tenant, err)
					s.log.Errorf(reason)
				}
			}
		}
	}

	if !failedSync {
		if contactInput.HasOwnerByOwnerId() {
			if err = s.repositories.ContactRepository.SetOwnerByOwnerExternalId(ctx, tenant, contactId, contactInput.UserExternalOwnerId, dataService.SourceId()); err != nil {
				// Do not mark sync as failed in case owner relationship is not set
				tracing.TraceErr(span, err)
				s.log.Errorf("failed set owner user for contact %s, tenant %s :%v", contactId, tenant, err)
			}
		} else if contactInput.HasOwnerByUserId() {
			if err = s.repositories.ContactRepository.SetOwnerByUserExternalId(ctx, tenant, contactId, contactInput.UserExternalUserId, dataService.SourceId()); err != nil {
				// Do not mark sync as failed in case owner relationship is not set
				tracing.TraceErr(span, err)
				s.log.Errorf("failed set owner user for contact %s, tenant %s :%v", contactId, tenant, err)
			}
		}
	}

	if contactInput.HasNotes() && !failedSync {
		for _, note := range contactInput.Notes {
			localNote := entity.NoteData{
				BaseData: entity.BaseData{
					CreatedAt:      contactInput.CreatedAt,
					ExternalId:     string(note.FieldSource) + "-" + contactInput.ExternalId,
					ExternalSystem: contactInput.ExternalSystem,
				},
				Content: note.Note,
			}
			noteId, err := s.repositories.NoteRepository.GetMatchedNoteId(ctx, tenant, localNote)
			if err != nil {
				failedSync = true
				tracing.TraceErr(span, err)
				reason = fmt.Sprintf("failed finding existing matched note with external reference id %s for tenant %s :%s", localNote.ExternalId, tenant, err)
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
				reason = fmt.Sprintf("failed merge note for contact %v, tenant %v :%v", contactId, tenant, err)
				s.log.Errorf(reason)
				break
			}
			err = s.repositories.NoteRepository.NoteLinkWithContactByExternalId(ctx, tenant, noteId, contactInput.ExternalId, contactInput.ExternalSystem)
			if err != nil {
				failedSync = true
				tracing.TraceErr(span, err)
				reason = fmt.Sprintf("failed link note with contact %v, tenant %v :%v", contactId, tenant, err)
				s.log.Errorf(reason)
				break
			}
		}
	}

	if contactInput.HasTextCustomFields() && !failedSync {
		for _, customField := range contactInput.TextCustomFields {
			if err = s.repositories.ContactRepository.MergeTextCustomField(ctx, tenant, contactId, customField); err != nil {
				failedSync = true
				tracing.TraceErr(span, err)
				reason = fmt.Sprintf("failed merge custom field %v for contact %v, tenant %v :%v", customField.Name, contactId, tenant, err)
				s.log.Errorf(reason)
				break
			}
		}
	}

	if contactInput.HasLocation() && !failedSync {
		err = s.repositories.ContactRepository.MergeContactLocation(ctx, tenant, contactId, contactInput)
		if err != nil {
			failedSync = true
			tracing.TraceErr(span, err)
			reason = fmt.Sprintf("failed merge location for contact %v, tenant %v :%v", contactId, tenant, err)
			s.log.Errorf(reason)
		}
	}

	if contactInput.HasTags() && !failedSync {
		for _, tag := range contactInput.Tags {
			err = s.repositories.ContactRepository.MergeTagForContact(ctx, tenant, contactId, tag, contactInput.ExternalSystem)
			if err != nil {
				failedSync = true
				tracing.TraceErr(span, err)
				reason = fmt.Sprintf("failed to merge tag %v for contact %v, tenant %v :%v", tag, contactId, tenant, err)
				s.log.Errorf(reason)
				break
			}
		}
	}

	s.services.OrganizationService.UpdateLastTouchpointByContactId(ctx, tenant, contactId)

	s.log.Debugf("successfully merged contact with id %v for tenant %v from %v", contactId, tenant, dataService.SourceId())
	if err = dataService.MarkProcessed(ctx, contactInput.SyncId, runId, failedSync == false, false, reason); err != nil {
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
