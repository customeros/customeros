package service

import (
	"context"
	"encoding/json"
	"github.com/customeros/mailsherpa/domaincheck"
	"github.com/openline-ai/openline-customer-os/packages/runner/customer-os-data-upkeeper/config"
	"github.com/openline-ai/openline-customer-os/packages/runner/customer-os-data-upkeeper/logger"
	commonService "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	postgresentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/pkg/errors"
)

type GlobalOrganizationService interface {
	SyncScrapInToGlobalOrganization()
}

type globalOrganizationService struct {
	cfg            *config.Config
	log            logger.Logger
	commonServices *commonService.Services
}

func NewGlobalOrganizationService(cfg *config.Config, log logger.Logger, commonServices *commonService.Services) GlobalOrganizationService {
	return &globalOrganizationService{
		cfg:            cfg,
		log:            log,
		commonServices: commonServices,
	}
}

func (s *globalOrganizationService) SyncScrapInToGlobalOrganization() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Cancel context on exit

	span, ctx := tracing.StartTracerSpan(ctx, "GlobalOrganizationService.SyncScrapInToGlobalOrganization")
	defer span.Finish()
	tracing.TagComponentCronJob(span)

	limit := 50

	records, err := s.commonServices.PostgresRepositories.EnrichDetailsScrapInRepository.GetToSyncIntoGlobalOrganizations(ctx, limit)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "error getting records to sync"))
		s.log.Errorf("Error getting records to sync: %s", err.Error())
		return
	}

	// no record
	if len(records) == 0 {
		return
	}

	//process records
	for _, record := range records {

		// mark record as synced initially to not process same record again, even if error occurs
		err = s.commonServices.PostgresRepositories.EnrichDetailsScrapInRepository.MarkSyncedToGlobalOrganizations(ctx, record.ID)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "error marking record as synced"))
			s.log.Errorf("Error marking record as synced: %s", err.Error())
			continue
		}

		if record.Data == "" {
			continue
		}

		// unmarshal cached data
		data := postgresentity.ScrapInResponseBody{}
		if err = json.Unmarshal([]byte(record.Data), &data); err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to unmarshal scrapin data"))
			continue
		}

		if data.Company == nil {
			continue
		}

		// if company website is missing skip processing as we cannot validate primary domain
		if data.Company.WebsiteUrl == "" {
			continue
		}

		// identify primary domain
		_, primaryDomain := domaincheck.PrimaryDomainCheck(data.Company.WebsiteUrl)

		// if primary domain is empty, skip processing
		if primaryDomain == "" {
			continue
		}

		// if global organization already exists, update otherwise create
		globalOrganization, err := s.commonServices.PostgresRepositories.GlobalOrganizationRepository.GetByPrimaryDomain(ctx, primaryDomain)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "error getting global organization by primary domain"))
			s.log.Errorf("Error getting global organization by primary domain: %s", err.Error())
			continue
		}

		createGlobalOrg := false
		if globalOrganization == nil {
			createGlobalOrg = true
			now := utils.Now()
			globalOrganization = &postgresentity.GlobalOrganization{
				PrimaryDomain: primaryDomain,
				CreatedAt:     now,
				UpdatedAt:     now,
			}
		}

		// populate global organization entity
		if data.Company.Name != "" {
			globalOrganization.Name = data.Company.Name
		}
		if data.Company.Description != "" {
			globalOrganization.Description = data.Company.Description
		}
		if data.Company.Tagline != nil {
			if tagline, ok := data.Company.Tagline.(string); ok {
				globalOrganization.ValueProposition = tagline
			}
		}
		if data.Company.GetEmployeeCount() > 0 {
			globalOrganization.EmployeeCount = data.Company.GetEmployeeCount()
		}
		if data.Company.WebsiteUrl != "" {
			globalOrganization.Website = data.Company.WebsiteUrl
		}
		if data.Company.FoundedOn.Year > 0 {
			globalOrganization.YearFounded = data.Company.FoundedOn.Year
		}
		if data.Company.LinkedInUrl != "" {
			globalOrganization.LinkedInUrl = data.Company.LinkedInUrl
		}
		if data.Company.UniversalName != "" {
			globalOrganization.LinkedInAlias = data.Company.UniversalName
		}
		if data.Company.Logo != "" {
			globalOrganization.LogoUrl = data.Company.Logo
		}

		if createGlobalOrg {
			// create global organization
			_, err = s.commonServices.PostgresRepositories.GlobalOrganizationRepository.Create(ctx, globalOrganization)
			if err != nil {
				tracing.TraceErr(span, errors.Wrap(err, "error creating global organization"))
				s.log.Errorf("Error creating global organization: %s", err.Error())
				continue
			}
		} else {
			_, err = s.commonServices.PostgresRepositories.GlobalOrganizationRepository.Update(ctx, globalOrganization)
			if err != nil {
				tracing.TraceErr(span, errors.Wrap(err, "error updating global organization"))
				s.log.Errorf("Error updating global organization: %s", err.Error())
				continue
			}
		}
	}
}
