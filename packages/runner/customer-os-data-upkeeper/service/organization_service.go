package service

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/runner/customer-os-data-upkeeper/config"
	"github.com/openline-ai/openline-customer-os/packages/runner/customer-os-data-upkeeper/constants"
	"github.com/openline-ai/openline-customer-os/packages/runner/customer-os-data-upkeeper/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/clients/grpc_client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/dto"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	commonService "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	organizationpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/organization"
	"github.com/pkg/errors"
	"time"
)

type OrganizationService interface {
	RefreshLastTouchpoint()
	UpkeepOrganizations()
	SendReminders()
}

type organizationService struct {
	cfg                    *config.Config
	log                    logger.Logger
	commonServices         *commonService.Services
	eventsProcessingClient *grpc_client.Clients
}

func NewOrganizationService(cfg *config.Config, log logger.Logger, commonServices *commonService.Services, client *grpc_client.Clients) OrganizationService {
	return &organizationService{
		cfg:                    cfg,
		log:                    log,
		commonServices:         commonServices,
		eventsProcessingClient: client,
	}
}

func (s *organizationService) RefreshLastTouchpoint() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Cancel context on exit

	span, ctx := tracing.StartTracerSpan(ctx, "OrganizationService.RefreshLastTouchpoint")
	defer span.Finish()
	tracing.TagComponentCronJob(span)

	if s.eventsProcessingClient == nil {
		s.log.Warn("eventsProcessingClient is nil.")
		return
	}

	limit := 50
	delayFromPreviousCheckInMinutes := 60 // 60 minutes

	for {
		select {
		case <-ctx.Done():
			s.log.Infof("Context cancelled, stopping")
			return
		default:
			// continue as normal
		}

		records, err := s.commonServices.Neo4jRepositories.OrganizationReadRepository.GetOrganizationsForUpdateLastTouchpoint(ctx, limit, delayFromPreviousCheckInMinutes)
		if err != nil {
			tracing.TraceErr(span, err)
			s.log.Errorf("Error getting organizations for renewals: %v", err)
			return
		}

		// no record
		if len(records) == 0 {
			return
		}

		//process organizations
		for _, record := range records {
			innerCtx := common.WithCustomContext(ctx, &common.CustomContext{
				Tenant:    record.Tenant,
				AppSource: constants.AppSourceDataUpkeeper,
			})

			err = s.commonServices.OrganizationService.RequestRefreshLastTouchpoint(innerCtx, record.OrganizationId)
			if err != nil {
				tracing.TraceErr(span, errors.Wrap(err, "error refreshing last touchpoint"))
				s.log.Errorf("Error refreshing last touchpoint for organization {%s}: %s", record.OrganizationId, err.Error())
			}

			err = s.commonServices.Neo4jRepositories.CommonWriteRepository.UpdateTimeProperty(innerCtx, record.Tenant, model.NodeLabelOrganization, record.OrganizationId, string(neo4jentity.OrganizationPropertyLastTouchpointRequestedAt), utils.NowPtr())
			if err != nil {
				tracing.TraceErr(span, errors.Wrap(err, "error updating last touchpoint requested at"))
				s.log.Errorf("Error updating refresh last touchpoint requested at: %s", err.Error())
			}
		}

		// if less than limit records are returned, we are done
		if len(records) < limit {
			return
		}

		// force exit after single iteration
		return
	}
}

func (s *organizationService) UpkeepOrganizations() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Cancel context on exit

	if s.eventsProcessingClient == nil {
		s.log.Warn("eventsProcessingClient is nil.")
		return
	}

	s.updateDerivedNextRenewalDates(ctx)
	s.linkWithDomain(ctx)
	s.enrichOrganization(ctx)
	s.removeEmptySocials(ctx)
	s.removeDuplicatedSocials(ctx)
}

func (s *organizationService) updateDerivedNextRenewalDates(ctx context.Context) {
	span, ctx := tracing.StartTracerSpan(ctx, "OrganizationService.updateDerivedNextRenewalDates")
	defer span.Finish()
	tracing.TagComponentCronJob(span)

	limit := 1000

	for {
		select {
		case <-ctx.Done():
			s.log.Infof("Context cancelled, stopping")
			return
		default:
			// continue as normal
		}

		records, err := s.commonServices.Neo4jRepositories.OrganizationReadRepository.GetOrganizationsForUpdateNextRenewalDate(ctx, limit)
		if err != nil {
			tracing.TraceErr(span, err)
			s.log.Errorf("Error getting organizations for renewals: %v", err)
			return
		}

		// no record
		if len(records) == 0 {
			return
		}

		//process organizations
		for _, record := range records {
			_, err = utils.CallEventsPlatformGRPCWithRetry[*organizationpb.OrganizationIdGrpcResponse](func() (*organizationpb.OrganizationIdGrpcResponse, error) {
				return s.eventsProcessingClient.OrganizationClient.RefreshRenewalSummary(ctx, &organizationpb.RefreshRenewalSummaryGrpcRequest{
					Tenant:         record.Tenant,
					OrganizationId: record.OrganizationId,
					AppSource:      constants.AppSourceDataUpkeeper,
				})
			})
			if err != nil {
				tracing.TraceErr(span, err)
				s.log.Errorf("Error refreshing renewal summary for organization {%s}: %s", record.OrganizationId, err.Error())
			}
		}

		// if less than limit records are returned, we are done
		if len(records) < limit {
			return
		}

		//sleep for async processing, then check again
		time.Sleep(5 * time.Second)

		// force exit after single iteration
		return
	}
}

func (s *organizationService) linkWithDomain(ctx context.Context) {
	span, ctx := tracing.StartTracerSpan(ctx, "OrganizationService.linkWithDomain")
	defer span.Finish()
	tracing.TagComponentCronJob(span)

	limit := 100
	delayMinutesFromLastCheck := 60 * 24 * 3 // 3 days

	for {
		select {
		case <-ctx.Done():
			s.log.Infof("Context cancelled, stopping")
			return
		default:
			// continue as normal
		}

		records, err := s.commonServices.Neo4jRepositories.OrganizationReadRepository.GetOrganizationsWithWebsiteAndWithoutDomains(ctx, limit, delayMinutesFromLastCheck)
		if err != nil {
			tracing.TraceErr(span, err)
			s.log.Errorf("Error getting organizations: %v", err)
			return
		}

		// no record
		if len(records) == 0 {
			return
		}

		//process organizations
		for _, record := range records {
			innerCtx := common.WithCustomContext(ctx, &common.CustomContext{
				Tenant:    record.Tenant,
				AppSource: constants.AppSourceDataUpkeeper,
			})

			organizationEntity, err := s.commonServices.OrganizationService.GetById(innerCtx, record.Tenant, record.OrganizationId)
			if err != nil {
				tracing.TraceErr(span, err)
				s.log.Errorf("Error getting organization {%s}: %s", record.OrganizationId, err.Error())
			}

			primaryDomain, _ := s.commonServices.DomainService.GetPrimaryDomainForOrganizationWebsite(innerCtx, organizationEntity.Website)
			if primaryDomain != "" {
				_, err = s.commonServices.OrganizationService.LinkWithDomain(innerCtx, nil, record.OrganizationId, primaryDomain)
				if err != nil {
					tracing.TraceErr(span, err)
					s.log.Errorf("Error linking with domain {%s}: %s", record.OrganizationId, err.Error())
				}
			}
			err = s.commonServices.Neo4jRepositories.CommonWriteRepository.UpdateTimeProperty(innerCtx, record.Tenant, model.NodeLabelOrganization, record.OrganizationId, string(neo4jentity.OrganizationPropertyDomainCheckedAt), utils.NowPtr())
			if err != nil {
				tracing.TraceErr(span, err)
				s.log.Errorf("Error updating domain checked at: %s", err.Error())
			}
		}

		// if less than limit records are returned, we are done
		if len(records) < limit {
			return
		}

		// force exit after single iteration
		return
	}
}

func (s *organizationService) enrichOrganization(ctx context.Context) {
	span, ctx := tracing.StartTracerSpan(ctx, "OrganizationService.enrichOrganization")
	defer span.Finish()
	tracing.TagComponentCronJob(span)

	limit := 250
	delayFromPreviousAttemptInMinutes := 60 * 24 * 10 // 10 days

	for {
		select {
		case <-ctx.Done():
			s.log.Infof("Context cancelled, stopping")
			return
		default:
			// continue as normal
		}

		records, err := s.commonServices.Neo4jRepositories.OrganizationReadRepository.GetOrganizationsForEnrichByDomain(ctx, limit, delayFromPreviousAttemptInMinutes)
		if err != nil {
			tracing.TraceErr(span, err)
			s.log.Errorf("Error getting organizations: %v", err)
			return
		}

		// no record
		if len(records) == 0 {
			return
		}

		//process organizations
		for _, record := range records {
			innerCtx := common.WithCustomContext(ctx, &common.CustomContext{
				Tenant:    record.Tenant,
				AppSource: constants.AppSourceDataUpkeeper,
			})
			err = s.commonServices.RabbitMQService.PublishEvent(innerCtx, record.OrganizationId, model.ORGANIZATION, dto.RequestEnrichOrganization{Url: record.Param1})
			if err != nil {
				tracing.TraceErr(span, err)
				s.log.Errorf("Error enriching organization {%s}: %s", record.OrganizationId, err.Error())
			}
			err = s.commonServices.Neo4jRepositories.OrganizationWriteRepository.UpdateTimeProperty(innerCtx, record.Tenant, record.OrganizationId, string(neo4jentity.OrganizationPropertyEnrichRequestedAt), utils.NowPtr())
			if err != nil {
				tracing.TraceErr(span, err)
				s.log.Errorf("Error updating domain checked at: %s", err.Error())
			}

			// increment enrich attempts
			err = s.commonServices.Neo4jRepositories.CommonWriteRepository.IncrementProperty(innerCtx, record.Tenant, model.NodeLabelOrganization, record.OrganizationId, string(neo4jentity.OrganizationPropertyEnrichAttempts))
			if err != nil {
				tracing.TraceErr(span, err)
				s.log.Errorf("Error incrementing contact' enrich attempts: %s", err.Error())
			}
		}

		// if less than limit records are returned, we are done
		if len(records) < limit {
			return
		}

		// force exit after single iteration
		return
	}
}

func (s *organizationService) removeEmptySocials(ctx context.Context) {
	span, ctx := tracing.StartTracerSpan(ctx, "OrganizationService.removeEmptySocials")
	defer span.Finish()
	tracing.TagComponentCronJob(span)

	limit := 100
	minutesSinceLastUpdate := 180 // 3 hours

	records, err := s.commonServices.Neo4jRepositories.SocialReadRepository.GetEmptySocialsForEntityType(ctx, model.NodeLabelOrganization, minutesSinceLastUpdate, limit)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error getting socials: %v", err)
		return
	}

	// no record
	if len(records) == 0 {
		return
	}

	//remove socials from organization
	for _, record := range records {
		innerCtx := common.WithCustomContext(ctx, &common.CustomContext{
			Tenant:    record.Tenant,
			AppSource: constants.AppSourceDataUpkeeper,
		})
		err = s.commonServices.SocialService.RemoveSocialFromEntity(innerCtx, nil,
			commonService.LinkWith{
				Id:   record.LinkedEntityId,
				Type: model.ORGANIZATION,
			},
			record.SocialId)
		if err != nil {
			tracing.TraceErr(span, err)
			s.log.Errorf("Error removing social {%s}: %s", record.SocialId, err.Error())
		}
	}
}

func (s *organizationService) removeDuplicatedSocials(ctx context.Context) {
	span, ctx := tracing.StartTracerSpan(ctx, "OrganizationService.removeDuplicatedSocials")
	defer span.Finish()
	tracing.TagComponentCronJob(span)

	limit := 100

	records, err := s.commonServices.Neo4jRepositories.SocialReadRepository.GetDuplicatedSocialsForEntityType(ctx, model.NodeLabelOrganization, 180, limit)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error getting socials: %v", err)
		return
	}

	// no record
	if len(records) == 0 {
		return
	}

	//remove socials from organization
	for _, record := range records {
		innerCtx := common.WithCustomContext(ctx, &common.CustomContext{
			Tenant:    record.Tenant,
			AppSource: constants.AppSourceDataUpkeeper,
		})

		// remove social
		err = s.commonServices.SocialService.RemoveSocialFromEntity(innerCtx, nil,
			commonService.LinkWith{
				Id:   record.LinkedEntityId,
				Type: model.ORGANIZATION,
			},
			record.SocialId)
		if err != nil {
			tracing.TraceErr(span, err)
			s.log.Errorf("Error removing social {%s}: %s", record.SocialId, err.Error())
		}

	}

}

func (s *organizationService) SendReminders() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Cancel context on exit

	span, ctx := tracing.StartTracerSpan(ctx, "OrganizationService.SendReminders")
	defer span.Finish()
	tracing.TagComponentCronJob(span)

	readyToSendNodes, err := s.commonServices.Neo4jRepositories.ReminderReadRepository.GetReadyToSend(ctx, utils.Now())
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error getting reminders to send: %v", err)
		return
	}

	for _, reminderNode := range readyToSendNodes {
		tenant := model.GetTenantFromLabels(reminderNode.Labels, model.NodeLabelReminder)
		innerCtx := common.WithCustomContext(ctx, &common.CustomContext{
			Tenant:    tenant,
			AppSource: constants.AppSourceDataUpkeeper,
		})

		reminder := neo4jmapper.MapDbNodeToReminderEntity(reminderNode)
		err := s.commonServices.ReminderService.SendNotification(innerCtx, reminder.Id, s.cfg.NovuConfig.FronteraUrl)
		if err != nil {
			tracing.TraceErr(span, err)
			s.log.Errorf("Error sending reminder {%s}: %s", reminder.Id, err.Error())
			return
		}

		err = s.commonServices.ReminderService.UpdateReminder(ctx, tenant, reminder.Id, nil, nil, nil, utils.BoolPtr(true))
		if err != nil {
			tracing.TraceErr(span, err)
			return
		}
	}
}
