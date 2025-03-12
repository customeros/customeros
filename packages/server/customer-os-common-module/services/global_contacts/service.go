package globalcontacts

import (
	"context"
	"strconv"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"

	postgres_repository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	tracingLog "github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
)

type globalContactService struct {
	log      logger.Logger
	postgres *postgres_repository.Repositories
}

func NewGlobalContactService(log logger.Logger, postgres *postgres_repository.Repositories) interfaces.GlobalContactService {
	return &globalContactService{
		log:      log,
		postgres: postgres,
	}
}

func (s *globalContactService) findExistingContact(ctx context.Context, contact *postgres_entity.GlobalContact) (*postgres_entity.GlobalContact, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "GlobalContactService.findExistingContact")
	defer span.Finish()
	tracing.TagComponentService(span)

	var existingContact *postgres_entity.GlobalContact
	var err error

	if contact.PrimaryDomain != "" {
		// If primary domain is present, first try to find by domain combinations
		span.LogFields(tracingLog.String("search_mode", "with_primary_domain"))

		// Try LinkedIn identifier + primary domain
		if contact.LinkedInIdentifier != "" {
			span.LogFields(tracingLog.String("search_by", "linkedin_and_domain"))
			existingContact, err = s.postgres.GlobalContactRepository.GetByLinkedInIdentifierAndDomain(ctx, contact.LinkedInIdentifier, contact.PrimaryDomain)
			if err != nil {
				tracing.TraceErr(span, errors.Wrap(err, "error getting contact by LinkedIn identifier and domain"))
				return nil, err
			}
		} else if contact.WorkEmail != "" {
			// Try work email + primary domain
			span.LogFields(tracingLog.String("search_by", "work_email_and_domain"))
			existingContact, err = s.postgres.GlobalContactRepository.GetByWorkEmailAndDomain(ctx, contact.WorkEmail, contact.PrimaryDomain)
			if err != nil {
				tracing.TraceErr(span, errors.Wrap(err, "error getting contact by work email and domain"))
				return nil, err
			}
		} else if contact.PersonalEmail != "" {
			// Try personal email + primary domain
			span.LogFields(tracingLog.String("search_by", "personal_email_and_domain"))
			existingContact, err = s.postgres.GlobalContactRepository.GetByPersonalEmailAndDomain(ctx, contact.PersonalEmail, contact.PrimaryDomain)
			if err != nil {
				tracing.TraceErr(span, errors.Wrap(err, "error getting contact by personal email and domain"))
				return nil, err
			}
		}
	}

	if existingContact != nil && existingContact.LinkedInIdentifier != "" && contact.LinkedInIdentifier != "" && existingContact.LinkedInIdentifier != contact.LinkedInIdentifier {
		return nil, nil
	}
	if existingContact != nil && existingContact.PrimaryDomain != "" && contact.PrimaryDomain != "" && existingContact.PrimaryDomain != contact.PrimaryDomain {
		return nil, nil
	}

	return existingContact, nil
}

func (s *globalContactService) updateContactFields(existing *postgres_entity.GlobalContact, new *postgres_entity.GlobalContact) {
	// Only update fields if they are non-empty in the new contact
	existing.JobStartedAt = new.JobStartedAt
	existing.JobEndedAt = new.JobEndedAt
	if new.FirstName != "" {
		existing.FirstName = new.FirstName
	}
	if new.LastName != "" {
		existing.LastName = new.LastName
	}
	if new.Description != "" {
		existing.Description = new.Description
	}
	if new.JobTitle != "" {
		existing.JobTitle = new.JobTitle
	}
	if new.PrimaryDomain != "" {
		existing.PrimaryDomain = new.PrimaryDomain
	}
	if new.WorkEmail != "" {
		existing.WorkEmail = new.WorkEmail
	}
	if new.PersonalEmail != "" {
		existing.PersonalEmail = new.PersonalEmail
	}
	if new.LinkedInIdentifier != "" {
		existing.LinkedInIdentifier = new.LinkedInIdentifier
	}
	if new.LinkedInAlias != "" {
		existing.LinkedInAlias = new.LinkedInAlias
	}
	if new.LocationText != "" {
		existing.LocationText = new.LocationText
	}
	if new.ProfilePhotoExternalUrl != "" {
		// if new URL is different then existing one and download status is error, reset download status
		if existing.ProfilePhotoExternalUrl != new.ProfilePhotoExternalUrl {
			existing.DownloadStatus = enum.DownloadNotStarted
		}
		existing.ProfilePhotoExternalUrl = new.ProfilePhotoExternalUrl
	}
	if new.DataFetchedAt != nil {
		existing.DataFetchedAt = new.DataFetchedAt
	}
}

func (s *globalContactService) SaveContact(ctx context.Context, contact *postgres_entity.GlobalContact) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "GlobalContactService.SaveContact")
	defer span.Finish()
	tracing.TagComponentService(span)

	// Try to find existing contact using various identifiers
	existingContact, err := s.findExistingContact(ctx, contact)
	if err != nil {
		return err
	}

	if existingContact != nil {
		// Update existing contact preserving non-empty fields
		s.updateContactFields(existingContact, contact)

		_, err = s.postgres.GlobalContactRepository.Update(ctx, existingContact)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "error updating global contact"))
			return err
		}
	} else {
		// Create new contact
		_, err = s.postgres.GlobalContactRepository.Create(ctx, contact)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "error creating global contact"))
			return err
		}
	}

	return nil
}

func (s *globalContactService) SetWorkEmail(ctx context.Context, id uint64, workEmail string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "GlobalContactService.SetWorkEmail")
	defer span.Finish()
	tracing.TagComponentService(span)
	tracing.TagEntity(span, strconv.FormatUint(id, 10))

	contact, err := s.postgres.GlobalContactRepository.GetById(ctx, id)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "error getting contact by id"))
		return err
	}

	contact.WorkEmail = workEmail

	_, err = s.postgres.GlobalContactRepository.Update(ctx, contact)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "error updating global contact"))
		return err
	}

	return nil
}

func (s *globalContactService) GetGlobalContactsByLinkedIn(ctx context.Context, linkedIn string) ([]*postgres_entity.GlobalContact, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "GlobalContactService.GetGlobalContactsByLinkedIn")
	defer span.Finish()
	tracing.TagComponentService(span)
	span.LogFields(tracingLog.String("linkedIn", linkedIn))

	if linkedIn == "" {
		span.LogFields(tracingLog.Int("result.count", 0))
		return nil, nil
	}

	contacts, err := s.postgres.GlobalContactRepository.GetByLinkedInIdentifier(ctx, linkedIn)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "error getting contacts by linkedin identifier"))
		return nil, err
	}

	if len(contacts) == 0 {
		contacts, err = s.postgres.GlobalContactRepository.GetByLinkedInAlias(ctx, linkedIn)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "error getting contacts by linkedin alias"))
			return nil, err
		}
	}

	span.LogFields(tracingLog.Int("result.count", len(contacts)))
	return contacts, nil
}
