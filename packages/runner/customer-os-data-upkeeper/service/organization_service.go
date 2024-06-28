package service

import (
	"context"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/runner/customer-os-data-upkeeper/config"
	"github.com/openline-ai/openline-customer-os/packages/runner/customer-os-data-upkeeper/constants"
	"github.com/openline-ai/openline-customer-os/packages/runner/customer-os-data-upkeeper/logger"
	"github.com/openline-ai/openline-customer-os/packages/runner/customer-os-data-upkeeper/repository"
	"github.com/openline-ai/openline-customer-os/packages/runner/customer-os-data-upkeeper/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/grpc_client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/neo4jutil"
	organizationpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/organization"
	"net/http"
	"time"
)

type OrganizationService interface {
	RefreshLastTouchpoint()
	UpkeepOrganizations()
}

type organizationService struct {
	cfg                    *config.Config
	log                    logger.Logger
	repositories           *repository.Repositories
	eventsProcessingClient *grpc_client.Clients
}

func NewOrganizationService(cfg *config.Config, log logger.Logger, repositories *repository.Repositories, client *grpc_client.Clients) OrganizationService {
	return &organizationService{
		cfg:                    cfg,
		log:                    log,
		repositories:           repositories,
		eventsProcessingClient: client,
	}
}

func (s *organizationService) RefreshLastTouchpoint() {
	if s.eventsProcessingClient == nil {
		s.log.Warn("eventsProcessingClient is nil.")
		return
	}

	headers := map[string]string{
		"X-Openline-TENANT":  "openlineai",
		"X-Openline-API-KEY": s.cfg.PlatformAdminApi.ApiKey,
	}

	req, err := http.NewRequest("POST", s.cfg.PlatformAdminApi.Url+"/organization/refreshLastTouchpoint", nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("RefreshLastTouchpoint: Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("RefreshLastTouchpoint: Error response:", resp.Status)
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

	now := utils.Now()

	s.updateDerivedNextRenewalDates(ctx, now)
	s.linkWithDomain(ctx)
	s.enrichOrganization(ctx)
	s.removeEmptySocials(ctx)
	s.removeDuplicatedSocials(ctx, now)
}

func (s *organizationService) updateDerivedNextRenewalDates(ctx context.Context, referenceTime time.Time) {
	span, ctx := tracing.StartTracerSpan(ctx, "OrganizationService.updateDerivedNextRenewalDates")
	defer span.Finish()

	limit := 1000

	for {
		select {
		case <-ctx.Done():
			s.log.Infof("Context cancelled, stopping")
			return
		default:
			// continue as normal
		}

		records, err := s.repositories.Neo4jRepositories.OrganizationReadRepository.GetOrganizationsForUpdateNextRenewalDate(ctx, limit)
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

	limit := 100

	for {
		select {
		case <-ctx.Done():
			s.log.Infof("Context cancelled, stopping")
			return
		default:
			// continue as normal
		}

		records, err := s.repositories.Neo4jRepositories.OrganizationReadRepository.GetOrganizationsWithWebsiteAndWithoutDomains(ctx, limit, 360)
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
			organizationDbNode, err := s.repositories.Neo4jRepositories.OrganizationReadRepository.GetOrganization(ctx, record.Tenant, record.OrganizationId)
			if err != nil {
				tracing.TraceErr(span, err)
				s.log.Errorf("Error getting organization {%s}: %s", record.OrganizationId, err.Error())
			}
			organizationEntity := neo4jmapper.MapDbNodeToOrganizationEntity(organizationDbNode)

			domain := utils.ExtractDomain(organizationEntity.Website)

			if domain != "" {
				_, err = utils.CallEventsPlatformGRPCWithRetry[*organizationpb.OrganizationIdGrpcResponse](func() (*organizationpb.OrganizationIdGrpcResponse, error) {
					return s.eventsProcessingClient.OrganizationClient.LinkDomainToOrganization(ctx, &organizationpb.LinkDomainToOrganizationGrpcRequest{
						Tenant:         record.Tenant,
						OrganizationId: record.OrganizationId,
						Domain:         domain,
						AppSource:      constants.AppSourceDataUpkeeper,
					})
				})
				if err != nil {
					tracing.TraceErr(span, err)
					s.log.Errorf("Error linking with domain {%s}: %s", record.OrganizationId, err.Error())
				}
			}
			err = s.repositories.Neo4jRepositories.OrganizationWriteRepository.MarkDomainCheckRequested(ctx, record.Tenant, record.OrganizationId)
			if err != nil {
				tracing.TraceErr(span, err)
				s.log.Errorf("Error marking domain check requested: %s", err.Error())
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

func (s *organizationService) enrichOrganization(ctx context.Context) {
	span, ctx := tracing.StartTracerSpan(ctx, "OrganizationService.enrichOrganization")
	defer span.Finish()

	limit := 100

	for {
		select {
		case <-ctx.Done():
			s.log.Infof("Context cancelled, stopping")
			return
		default:
			// continue as normal
		}

		records, err := s.repositories.Neo4jRepositories.OrganizationReadRepository.GetOrganizationsForEnrich(ctx, limit, 360)
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
			_, err = utils.CallEventsPlatformGRPCWithRetry[*organizationpb.OrganizationIdGrpcResponse](func() (*organizationpb.OrganizationIdGrpcResponse, error) {
				return s.eventsProcessingClient.OrganizationClient.EnrichOrganization(ctx, &organizationpb.EnrichOrganizationGrpcRequest{
					Tenant:         record.Tenant,
					OrganizationId: record.OrganizationId,
					Url:            record.Param1,
					AppSource:      constants.AppSourceDataUpkeeper,
				})
			})
			if err != nil {
				tracing.TraceErr(span, err)
				s.log.Errorf("Error enriching organization {%s}: %s", record.OrganizationId, err.Error())
			}

			err = s.repositories.Neo4jRepositories.OrganizationWriteRepository.MarkDomainCheckRequested(ctx, record.Tenant, record.OrganizationId)
			if err != nil {
				tracing.TraceErr(span, err)
				s.log.Errorf("Error marking domain check requested: %s", err.Error())
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

	limit := 100

	for {
		select {
		case <-ctx.Done():
			s.log.Infof("Context cancelled, stopping")
			return
		default:
			// continue as normal
		}

		minutesSinceLastUpdate := 180
		records, err := s.repositories.Neo4jRepositories.SocialReadRepository.GetEmptySocialsForEntityType(ctx, neo4jutil.NodeLabelOrganization, minutesSinceLastUpdate, limit)
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
			_, err = utils.CallEventsPlatformGRPCWithRetry[*organizationpb.OrganizationIdGrpcResponse](func() (*organizationpb.OrganizationIdGrpcResponse, error) {
				return s.eventsProcessingClient.OrganizationClient.RemoveSocial(ctx, &organizationpb.RemoveSocialGrpcRequest{
					Tenant:         record.Tenant,
					OrganizationId: record.LinkedEntityId,
					SocialId:       record.SocialId,
					AppSource:      constants.AppSourceDataUpkeeper,
				})
			})
			if err != nil {
				tracing.TraceErr(span, err)
				s.log.Errorf("Error removing social {%s}: %s", record.SocialId, err.Error())
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

func (s *organizationService) removeDuplicatedSocials(ctx context.Context, now time.Time) {
	span, ctx := tracing.StartTracerSpan(ctx, "OrganizationService.removeDuplicatedSocials")
	defer span.Finish()

	limit := 100

	for {
		select {
		case <-ctx.Done():
			s.log.Infof("Context cancelled, stopping")
			return
		default:
			// continue as normal
		}

		records, err := s.repositories.Neo4jRepositories.SocialReadRepository.GetDuplicatedSocialsForEntityType(ctx, neo4jutil.NodeLabelOrganization, 180, limit)
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
			_, err = utils.CallEventsPlatformGRPCWithRetry[*organizationpb.OrganizationIdGrpcResponse](func() (*organizationpb.OrganizationIdGrpcResponse, error) {
				return s.eventsProcessingClient.OrganizationClient.RemoveSocial(ctx, &organizationpb.RemoveSocialGrpcRequest{
					Tenant:         record.Tenant,
					OrganizationId: record.LinkedEntityId,
					SocialId:       record.SocialId,
					AppSource:      constants.AppSourceDataUpkeeper,
				})
			})
			if err != nil {
				tracing.TraceErr(span, err)
				s.log.Errorf("Error removing social {%s}: %s", record.SocialId, err.Error())
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
