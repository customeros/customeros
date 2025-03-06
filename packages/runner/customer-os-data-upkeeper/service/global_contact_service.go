package service

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/customeros/customeros/packages/runner/customer-os-data-upkeeper/config"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
	commonService "github.com/customeros/customeros/packages/server/customer-os-common-module/services"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	postgresentity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
)

type GlobalContactService interface {
	SyncDataIntoGlobalContacts()
}

type globalContactService struct {
	cfg             *config.Config
	log             logger.Logger
	commonServices  *commonService.CommonServices
	scrapinDataChan chan *postgresentity.ScrapInResponseBody
}

func NewGlobalContactService(cfg *config.Config, log logger.Logger, commonServices *commonService.CommonServices) GlobalContactService {
	return &globalContactService{
		cfg:             cfg,
		log:             log,
		commonServices:  commonServices,
		scrapinDataChan: make(chan *postgresentity.ScrapInResponseBody, 100), // Buffer size of 100
	}
}

func (s *globalContactService) findPrimaryDomainFromGlobalOrg(ctx context.Context, linkedInUrl string, linkedInId string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "GlobalContactService.findPrimaryDomainFromGlobalOrg")
	defer span.Finish()
	tracing.TagComponentService(span)

	var org *postgresentity.GlobalOrganization
	var err error

	// Try by LinkedIn URL first
	if linkedInUrl != "" {
		// Remove trailing slash if present
		linkedInUrl = strings.TrimSuffix(linkedInUrl, "/")
		org, err = s.commonServices.PostgresRepositories.GlobalOrganizationRepository.GetByLinkedInUrl(ctx, linkedInUrl)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "error getting organization by LinkedIn URL"))
			return "", err
		}
		if org != nil {
			return org.PrimaryDomain, nil
		}
	}

	// Try by LinkedIn ID
	if linkedInId != "" {
		linkedInUrl = "https://www.linkedin.com/company/" + linkedInId
		org, err = s.commonServices.PostgresRepositories.GlobalOrganizationRepository.GetByLinkedInUrl(ctx, linkedInUrl)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "error getting organization by LinkedIn ID"))
			return "", err
		}
		if org != nil {
			return org.PrimaryDomain, nil
		}
	}

	return "", nil
}

func (s *globalContactService) syncScrapinInRecordIntoGlobalContact(ctx context.Context, record *postgresentity.EnrichDetailsScrapIn) error {
	span, ctx := tracing.StartTracerSpan(ctx, "GlobalContactService.syncScrapinInRecordIntoGlobalContact")
	defer span.Finish()
	tracing.TagComponentService(span)

	// Mark as synced to avoid double processing
	err := s.commonServices.PostgresRepositories.EnrichDetailsScrapInRepository.MarkSyncedToGlobalContacts(ctx, record.ID)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "error marking scrapin data as synced to global contacts"))
		return err
	}

	if record.Data == "" {
		return nil
	}

	// unmarshal cached data
	data := postgresentity.ScrapInResponseBody{}
	if err = json.Unmarshal([]byte(record.Data), &data); err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to unmarshal scrapin data"))
		return err
	}

	if data.Person == nil {
		return nil
	}

	person := data.Person

	// Find all positions without end date
	var currentPositions []struct {
		Title       string
		StartedOn   *time.Time
		EndedOn     *time.Time
		LinkedInUrl string
		LinkedInId  string
	}

	for _, position := range person.Positions.PositionHistory {
		// Check if position has no end date
		if position.StartEndDate.End == nil {
			var currentPosition struct {
				Title       string
				StartedOn   *time.Time
				EndedOn     *time.Time
				LinkedInUrl string
				LinkedInId  string
			}

			// Convert start date
			if position.StartEndDate.Start != nil {
				startDate := utils.FirstTimeOfMonth(position.StartEndDate.Start.Year, position.StartEndDate.Start.Month)
				currentPosition.StartedOn = &startDate
			}
			currentPosition.Title = position.Title
			currentPosition.LinkedInUrl = position.LinkedInUrl
			currentPosition.LinkedInId = position.LinkedInId

			currentPositions = append(currentPositions, currentPosition)
		}
	}

	// Only proceed if we have current positions
	if len(currentPositions) == 0 {
		return nil
	}

	// Process each current position
	for _, currentPosition := range currentPositions {
		// Find primary domain from global organizations
		primaryDomain, err := s.findPrimaryDomainFromGlobalOrg(ctx, currentPosition.LinkedInUrl, currentPosition.LinkedInId)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "error finding primary domain from global organizations"))
			return err
		}

		contact := &postgresentity.GlobalContact{
			FirstName:          strings.TrimSpace(person.FirstName),
			LastName:           strings.TrimSpace(person.LastName),
			LinkedInIdentifier: person.LinkedInIdentifier,
			JobTitle:           currentPosition.Title,
			JobStartedAt:       currentPosition.StartedOn,
			JobEndedAt:         currentPosition.EndedOn,
			PrimaryDomain:      primaryDomain,
		}

		// Try to find existing contact by LinkedIn identifier and primary domain
		existingContact, err := s.commonServices.PostgresRepositories.GlobalContactRepository.GetByLinkedInIdentifier(ctx, contact.LinkedInIdentifier)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "error getting contact by LinkedIn identifier"))
			return err
		}

		if existingContact != nil {
			// Update existing contact
			existingContact.FirstName = contact.FirstName
			existingContact.LastName = contact.LastName
			existingContact.JobTitle = contact.JobTitle
			existingContact.JobStartedAt = contact.JobStartedAt
			existingContact.JobEndedAt = contact.JobEndedAt
			existingContact.PrimaryDomain = contact.PrimaryDomain

			_, err = s.commonServices.PostgresRepositories.GlobalContactRepository.Update(ctx, existingContact)
			if err != nil {
				tracing.TraceErr(span, errors.Wrap(err, "error updating global contact"))
				return err
			}
		} else {
			// Create new contact
			_, err = s.commonServices.PostgresRepositories.GlobalContactRepository.Create(ctx, contact)
			if err != nil {
				tracing.TraceErr(span, errors.Wrap(err, "error creating global contact"))
				return err
			}
		}
	}

	return nil
}

func (s *globalContactService) SyncDataIntoGlobalContacts() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	span, ctx := opentracing.StartSpanFromContext(ctx, "GlobalContactService.SyncDataIntoGlobalContacts")
	defer span.Finish()
	tracing.TagComponentService(span)

	limit := 50

	records, err := s.commonServices.PostgresRepositories.EnrichDetailsScrapInRepository.GetToSyncIntoGlobalContacts(ctx, limit)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "error getting records to sync"))
		s.log.Errorf("Error getting records to sync: %s", err.Error())
		return
	}

	// no records to process
	if len(records) == 0 {
		return
	}

	// process records
	for _, record := range records {
		if err := s.syncScrapinInRecordIntoGlobalContact(ctx, record); err != nil {
			s.log.Errorf("Error processing record: %s", err.Error())
			continue
		}
	}
}
