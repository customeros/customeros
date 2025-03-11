package service

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/opentracing/opentracing-go/log"

	"github.com/customeros/customeros/packages/runner/customer-os-data-upkeeper/config"
	"github.com/customeros/customeros/packages/runner/customer-os-data-upkeeper/constants"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	commonconstants "github.com/customeros/customeros/packages/server/customer-os-common-module/constants"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/data_fields"
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
	EnrichGlobalOrganizationWithBettercontact()
	SyncGlobalContactsToTenantContacts()
}

type globalContactService struct {
	cfg            *config.Config
	log            logger.Logger
	commonServices *commonService.CommonServices
}

func NewGlobalContactService(cfg *config.Config, log logger.Logger, commonServices *commonService.CommonServices) GlobalContactService {
	return &globalContactService{
		cfg:            cfg,
		log:            log,
		commonServices: commonServices,
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

	// Process all positions
	var positions []struct {
		Title       string
		StartedOn   *time.Time
		EndedOn     *time.Time
		LinkedInUrl string
		LinkedInId  string
	}

	for _, position := range person.Positions.PositionHistory {
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

		// Convert end date
		if position.StartEndDate.End != nil {
			endDate := utils.FirstTimeOfMonth(position.StartEndDate.End.Year, position.StartEndDate.End.Month)
			currentPosition.EndedOn = &endDate
		}

		currentPosition.Title = position.Title
		currentPosition.LinkedInUrl = position.LinkedInUrl
		currentPosition.LinkedInId = position.LinkedInId

		positions = append(positions, currentPosition)
	}

	// Only proceed if we have positions
	if len(positions) == 0 {
		return nil
	}

	// Process each position
	for _, currentPosition := range positions {
		// Find primary domain from global organizations
		primaryDomain, err := s.findPrimaryDomainFromGlobalOrg(ctx, currentPosition.LinkedInUrl, currentPosition.LinkedInId)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "error finding primary domain from global organizations"))
			return err
		}

		// Skip positions where we can't find a primary domain
		if primaryDomain == "" {
			s.log.Debugf("Skipping position with title '%s' - no primary domain found for LinkedIn URL: %s, LinkedIn ID: %s",
				currentPosition.Title, currentPosition.LinkedInUrl, currentPosition.LinkedInId)
			continue
		}

		contact := &postgresentity.GlobalContact{
			FirstName:               utils.CleanName(person.FirstName),
			LastName:                utils.CleanName(person.LastName),
			LinkedInIdentifier:      person.LinkedInIdentifier,
			LinkedInAlias:           person.PublicIdentifier,
			JobTitle:                currentPosition.Title,
			JobStartedAt:            currentPosition.StartedOn,
			JobEndedAt:              currentPosition.EndedOn,
			PrimaryDomain:           primaryDomain,
			ProfilePhotoExternalUrl: person.PhotoUrl,
			DataFetchedAt:           &record.CreatedAt,
		}

		err = s.commonServices.GlobalContactService.SaveContact(ctx, contact)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "error saving global contact"))
			return err
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

	limit := 500

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

func (s *globalContactService) EnrichGlobalOrganizationWithBettercontact() {
	s.sendRequestToBetterContact()
	s.processBetterContactResponses()
}

func (s *globalContactService) sendRequestToBetterContact() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Cancel context on exit

	// Better contact is limited to 60 requests per minute
	// https://bettercontact.notion.site/Documentation-API-e8e1b352a0d647ee9ff898609bf1a168
	limit := 40

	span, ctx := tracing.StartTracerSpan(ctx, "GlobalContactService.sendRequestToBetterContact")
	defer span.Finish()
	tracing.TagComponentCronJob(span)

	records, err := s.commonServices.PostgresRepositories.GlobalContactRepository.GetContactsToFindWorkEmailWithBetterContact(ctx, limit)
	if err != nil {
		tracing.TraceErr(span, err)
		return
	}

	// no record
	if len(records) == 0 {
		return
	}

	for _, record := range records {
		func(globalContact *postgresentity.GlobalContact) {
			innerSpan, innerCtx := tracing.StartTracerSpan(ctx, "GlobalContactService.sendRequestToBetterContact.Record")
			defer innerSpan.Finish()

			linkedInUrl := "https://linkedin.com/in/" + globalContact.LinkedInIdentifier
			// get global organization by primary domain
			globalOrg, err := s.commonServices.PostgresRepositories.GlobalOrganizationRepository.GetByPrimaryDomain(innerCtx, globalContact.PrimaryDomain)
			if err != nil {
				tracing.TraceErr(innerSpan, err)
			}
			companyName := ""
			if globalOrg != nil {
				companyName = globalOrg.Name
			}

			_, betterContactRequestId, _, err := s.commonServices.EnrichmentService.FindWorkEmail(innerCtx, linkedInUrl, globalContact.FirstName, globalContact.LastName, companyName, globalContact.PrimaryDomain, false)
			if err != nil {
				tracing.TraceErr(span, err)
				span.LogFields(log.Object("record", record))
			} else {
				// mark contact with enrich requested
				err = s.commonServices.PostgresRepositories.GlobalContactRepository.MarkBetterContactRequested(innerCtx, globalContact.ID, betterContactRequestId)
				if err != nil {
					tracing.TraceErr(innerSpan, err)
				}
			}
		}(record)
	}
}

func (s *globalContactService) processBetterContactResponses() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Cancel context on exit

	span, ctx := tracing.StartTracerSpan(ctx, "GlobalContactService.processBetterContactResponses")
	defer span.Finish()
	tracing.TagComponentCronJob(span)

	limit := 250

	records, err := s.commonServices.PostgresRepositories.GlobalContactRepository.GetContactsToSetWorkEmailFromBetterContactResponse(ctx, limit)
	if err != nil {
		tracing.TraceErr(span, err)
		return
	}

	// no record
	if len(records) == 0 {
		return
	}

	for _, record := range records {
		func(globalContact *postgresentity.GlobalContact) {
			innerSpan, innerCtx := tracing.StartTracerSpan(ctx, "GlobalContactService.processBetterContactResponses.Record")
			defer innerSpan.Finish()

			// get better contact response
			betterContactRecord, err := s.commonServices.PostgresRepositories.EnrichDetailsBetterContactRepository.GetByRequestId(innerCtx, globalContact.BetterContactRequestId)
			if err != nil {
				tracing.TraceErr(span, err)
				return
			}

			if betterContactRecord == nil || betterContactRecord.Response == "" {
				innerSpan.LogFields(log.String("message", "better contact response not ready"))
				return
			}

			var betterContactResponse postgresentity.BetterContactResponseBody
			if err = json.Unmarshal([]byte(betterContactRecord.Response), &betterContactResponse); err != nil {
				tracing.TraceErr(span, err)
				return
			}

			err = s.commonServices.PostgresRepositories.GlobalContactRepository.MarkBetterContactSet(innerCtx, globalContact.ID)
			if err != nil {
				tracing.TraceErr(innerSpan, err)
			}

			if betterContactResponse.Data[0].ContactEmailAddress != "" {
				err = s.commonServices.GlobalContactService.SetWorkEmail(innerCtx, globalContact.ID, betterContactResponse.Data[0].ContactEmailAddress)
				if err != nil {
					tracing.TraceErr(innerSpan, err)
				}
			}

		}(record)
	}
}

func (s *globalContactService) SyncGlobalContactsToTenantContacts() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	span, ctx := tracing.StartTracerSpan(ctx, "GlobalContactService.SyncGlobalContactsToTenantContacts")
	defer span.Finish()
	tracing.TagComponentCronJob(span)

	limit := 250
	daysFromPreviousSync := 1
	forceSyncAfterDays := 30

	records, err := s.commonServices.PostgresRepositories.GlobalContactRepository.GetGlobalContactsToSyncIntoTenantContacts(ctx, daysFromPreviousSync, forceSyncAfterDays, limit)
	if err != nil {
		tracing.TraceErr(span, err)
		return
	}

	if len(records) == 0 {
		return
	}

	for _, record := range records {
		s.syncGlobalContactToTenantContact(ctx, record)
	}
}

func (s *globalContactService) syncGlobalContactToTenantContact(ctx context.Context, globalContact *postgresentity.GlobalContact) {
	span, ctx := tracing.StartTracerSpan(ctx, "GlobalContactService.syncGlobalContactToTenantContact")
	defer span.Finish()
	tracing.TagComponentCronJob(span)
	tracing.TagEntity(span, strconv.FormatUint(globalContact.ID, 10))

	// sync photo
	if globalContact.ProfilePhotoExternalUrl != "" && globalContact.ProfilePhotoPath != "" {
		s.syncGlobalContactPhotoToTenantContact(ctx, globalContact)
	}

	// mark global contact as synced to neo4j
	// add 10 sec to avoid double processing
	globalContact.SyncedToNeoAt = utils.TimePtr(utils.Now().Add(10 * time.Second))
	_, err := s.commonServices.PostgresRepositories.GlobalContactRepository.Update(ctx, globalContact)
	if err != nil {
		tracing.TraceErr(span, err)
	}
}

func (s *globalContactService) syncGlobalContactPhotoToTenantContact(ctx context.Context, globalContact *postgresentity.GlobalContact) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "GlobalContactService.syncGlobalContactPhotoToTenantContact")
	defer span.Finish()
	tracing.TagComponentCronJob(span)

	if globalContact.ProfilePhotoPath == "" {
		return
	}
	if globalContact.ProfilePhotoExternalUrl == "" {
		return
	}

	tenantContact, err := s.commonServices.Neo4jRepositories.ContactReadRepository.GetContactsWithProfilePhotoUrlCrossTenant(ctx, globalContact.ProfilePhotoExternalUrl)
	if err != nil {
		tracing.TraceErr(span, err)
		return
	}

	if len(tenantContact) == 0 {
		return
	}

	for _, tenantContact := range tenantContact {
		innerCtx := common.WithCustomContext(ctx, &common.CustomContext{
			Tenant:    tenantContact.Tenant,
			AppSource: constants.AppSourceDataUpkeeper,
		})
		s.commonServices.ContactService.Save(innerCtx, nil, &tenantContact.ContactId, data_fields.ContactFields{
			ProfilePhotoUrl: utils.StringPtr(commonconstants.S3ImagesCDN + globalContact.ProfilePhotoPath),
		}, false)
	}

}
