package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-tracking/config"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-tracking/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	commonpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/common"
	organizationpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/organization"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"io"
	"net/http"
	"strings"
)

type TrackingService interface {
	ShouldIdentifyTrackingRecords(ctx context.Context) error
	IdentifyTrackingRecords(ctx context.Context) error
	CreateOrganizationsFromTrackedData(ctx context.Context) error
	NotifyOnSlack(ctx context.Context) error
}

type trackingService struct {
	cfg      *config.Config
	services *Services
}

func NewTrackingService(cfg *config.Config, services *Services) TrackingService {
	return &trackingService{
		cfg:      cfg,
		services: services,
	}
}

func (s *trackingService) ShouldIdentifyTrackingRecords(c context.Context) error {
	span, ctx := tracing.StartTracerSpan(c, "TrackingService.ShouldIdentifyTrackingRecords")
	defer span.Finish()

	needsIdentificationRecords, err := s.services.CommonServices.PostgresRepositories.TrackingRepository.GetForPrefilterBeforeIdentification(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	for _, record := range needsIdentificationRecords {
		shouldIdentify, err := s.processShouldIdentifyRecord(ctx, record.IP)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}

		if shouldIdentify != nil {

			state := entity.TrackingIdentificationStatePrefilteredFail
			if *shouldIdentify {
				state = entity.TrackingIdentificationStatePrefilteredPass
			}

			err = s.services.CommonServices.PostgresRepositories.TrackingRepository.SetStateById(ctx, record.ID, state)
			if err != nil {
				tracing.TraceErr(span, err)
				return err
			}
		}
	}

	return nil
}

func (s *trackingService) IdentifyTrackingRecords(c context.Context) error {
	span, ctx := tracing.StartTracerSpan(c, "TrackingService.IdentifyTrackingRecords")
	defer span.Finish()

	notIdentifiedTrackingRecords, err := s.services.CommonServices.PostgresRepositories.TrackingRepository.GetReadyForIdentification(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	for _, record := range notIdentifiedTrackingRecords {
		isIdentified, err := s.processRecordIdentification(ctx, record.IP)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}

		state := entity.TrackingIdentificationStateNotIdentified
		if isIdentified {
			state = entity.TrackingIdentificationStateIdentified
		}

		err = s.services.CommonServices.PostgresRepositories.TrackingRepository.SetStateById(ctx, record.ID, state)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}
	}

	return nil
}

func (s *trackingService) CreateOrganizationsFromTrackedData(c context.Context) error {
	span, ctx := tracing.StartTracerSpan(c, "TrackingService.CreateOrganizationsFromTrackedData")
	defer span.Finish()

	identifiedRecords, err := s.services.CommonServices.PostgresRepositories.TrackingRepository.GetIdentifiedWithDistinctIP(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	for _, r := range identifiedRecords {

		record, err := s.services.CommonServices.PostgresRepositories.TrackingRepository.GetById(ctx, r.ID)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}

		if record.State != entity.TrackingIdentificationStateIdentified {
			span.LogFields(log.String("skip", "bad state"))
			continue
		}

		snitcherData, err := s.services.CommonServices.PostgresRepositories.EnrichDetailsTrackingRepository.GetByIP(ctx, record.IP)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}

		if snitcherData == nil {
			tracing.TraceErr(span, errors.New("snitcher record is nil"))
			continue
		}

		if snitcherData.CompanyDomain == nil || *snitcherData.CompanyDomain == "" {
			tracing.TraceErr(span, errors.New("company domain is empty"))
			continue
		}

		organizationByDomainNode, err := s.services.CommonServices.Neo4jRepositories.OrganizationReadRepository.GetOrganizationWithDomain(ctx, record.Tenant, *snitcherData.CompanyDomain)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}

		if organizationByDomainNode == nil {

			if snitcherData.CompanyWebsite == nil || *snitcherData.CompanyWebsite == "" {
				span.LogFields(log.String("skip", "no website"))
				continue
			}

			companyName := ""
			companyWebsite := ""
			if snitcherData.CompanyName != nil {
				companyName = *snitcherData.CompanyName
			}
			if snitcherData.CompanyWebsite != nil {
				companyWebsite = *snitcherData.CompanyWebsite
			}

			upsertOrganizationRequest := organizationpb.UpsertOrganizationGrpcRequest{
				Tenant:       record.Tenant,
				Name:         companyName,
				Website:      companyWebsite,
				Relationship: neo4jenum.Prospect.String(),
				Stage:        neo4jenum.Lead.String(),
				LeadSource:   "Reveal AI",
				SourceFields: &commonpb.SourceFields{
					Source:    constants.SourceOpenline,
					AppSource: constants.AppTracking,
				},
			}

			organizationResponse, err := utils.CallEventsPlatformGRPCWithRetry[*organizationpb.OrganizationIdGrpcResponse](func() (*organizationpb.OrganizationIdGrpcResponse, error) {
				return s.services.GrpcClient.OrganizationClient.UpsertOrganization(ctx, &upsertOrganizationRequest)
			})
			if err != nil {
				tracing.TraceErr(span, err)
				return err
			}

			err = s.services.CommonServices.PostgresRepositories.TrackingRepository.MarkAsOrganizationCreated(ctx, record.ID, organizationResponse.Id, *snitcherData.CompanyName)
			if err != nil {
				tracing.TraceErr(span, err)
				return err
			}

			err = s.services.CommonServices.PostgresRepositories.TrackingRepository.MarkAllExcludeIdWithState(ctx, record.ID, record.IP, entity.TrackingIdentificationStateOrganizationExists)
			if err != nil {
				tracing.TraceErr(span, err)
				return err
			}
		} else {
			organizationId := utils.GetStringPropOrEmpty(organizationByDomainNode.Props, "id")
			organizationName := utils.GetStringPropOrEmpty(organizationByDomainNode.Props, "name")

			err = s.services.CommonServices.PostgresRepositories.TrackingRepository.MarkAsOrganizationCreated(ctx, record.ID, organizationId, organizationName)
			if err != nil {
				tracing.TraceErr(span, err)
				return err
			}

			err = s.services.CommonServices.PostgresRepositories.TrackingRepository.MarkAllWithState(ctx, record.IP, entity.TrackingIdentificationStateOrganizationExists)
			if err != nil {
				tracing.TraceErr(span, err)
				return err
			}
		}

	}

	return nil
}

func (s *trackingService) NotifyOnSlack(c context.Context) error {
	span, ctx := tracing.StartTracerSpan(c, "TrackingService.NotifyOnSlack")
	defer span.Finish()

	if s.cfg.SlackBotApiKey == "" {
		span.LogFields(log.String("skip", "no slack bot api key"))
		return nil
	}

	notifyOnSlackRecords, err := s.services.CommonServices.PostgresRepositories.TrackingRepository.GetForSlackNotifications(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if len(notifyOnSlackRecords) == 0 {
		span.LogFields(log.String("skip", "no records to notify"))
		return nil
	}

	for _, r := range notifyOnSlackRecords {
		err := s.notifyOnSlack(ctx, r)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}
	}

	return nil
}

func (s *trackingService) processShouldIdentifyRecord(c context.Context, ip string) (*bool, error) {
	span, ctx := opentracing.StartSpanFromContext(c, "TrackingService.processShouldIdentifyRecord")
	defer span.Finish()

	ipDataByIp, err := s.services.CommonServices.PostgresRepositories.EnrichDetailsPrefilterTrackingRepository.GetByIP(ctx, ip)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, fmt.Errorf("failed to get better contact details: %v", err)
	}

	if ipDataByIp == nil {
		ipDataByIp, err = s.askAndStoreIPData(ctx, ip)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}

	}

	if ipDataByIp == nil {
		tracing.TraceErr(span, errors.New("ipdata record is nil"))
	}

	return &ipDataByIp.ShouldIdentify, nil
}

func (s *trackingService) processRecordIdentification(c context.Context, ip string) (bool, error) {
	span, ctx := opentracing.StartSpanFromContext(c, "TrackingService.processRecordIdentification")
	defer span.Finish()

	snitcherByIp, err := s.services.CommonServices.PostgresRepositories.EnrichDetailsTrackingRepository.GetByIP(ctx, ip)
	if err != nil {
		tracing.TraceErr(span, err)
		return false, fmt.Errorf("failed to get better contact details: %v", err)
	}

	if snitcherByIp == nil {
		snitcherByIp, err = s.askAndStoreSnitcherData(ctx, ip)
		if err != nil {
			tracing.TraceErr(span, err)
			return false, err
		}
	}

	if snitcherByIp == nil {
		tracing.TraceErr(span, errors.New("snitcher record is nil"))
	}

	if snitcherByIp.CompanyDomain != nil && *snitcherByIp.CompanyDomain != "" {
		return true, nil
	}

	return false, nil
}

func (s *trackingService) askAndStoreIPData(c context.Context, ip string) (*entity.EnrichDetailsPreFilterTracking, error) {
	span, ctx := opentracing.StartSpanFromContext(c, "TrackingService.askAndStoreIPData")
	defer span.Finish()

	// Create HTTP client
	client := &http.Client{}

	// Create POST request
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s?api-key=%s", s.cfg.IPDataApi.Url, ip, s.cfg.IPDataApi.ApiKey), nil)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, fmt.Errorf("failed to create POST request: %v", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")

	//Perform the request
	resp, err := client.Do(req)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, fmt.Errorf("failed to perform POST request: %v", err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	// Parse the JSON request body
	var ipDataResponseBody entity.IPDataResponseBody
	if err = json.Unmarshal(responseBody, &ipDataResponseBody); err != nil {
		tracing.TraceErr(span, err)
		return nil, fmt.Errorf("failed to unmarshal response body: %v", err)
	}

	shouldIdentify := true

	if ipDataResponseBody.Ip == "" {
		shouldIdentify = false
	}

	if ipDataResponseBody.Carrier != nil {
		shouldIdentify = false
	}

	if ipDataResponseBody.Threat.IsTor ||
		ipDataResponseBody.Threat.IsIcloudRelay ||
		ipDataResponseBody.Threat.IsProxy ||
		ipDataResponseBody.Threat.IsDatacenter ||
		ipDataResponseBody.Threat.IsAnonymous ||
		ipDataResponseBody.Threat.IsKnownAttacker ||
		ipDataResponseBody.Threat.IsKnownAbuser ||
		ipDataResponseBody.Threat.IsThreat ||
		ipDataResponseBody.Threat.IsBogon {
		shouldIdentify = false
	}

	marshal, err := json.Marshal(ipDataResponseBody)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, fmt.Errorf("failed to marshal response body: %v", err)
	}

	// Store response
	err = s.services.CommonServices.PostgresRepositories.EnrichDetailsPrefilterTrackingRepository.RegisterRequest(ctx, entity.EnrichDetailsPreFilterTracking{
		CreatedAt:      utils.Now(),
		IP:             ip,
		ShouldIdentify: shouldIdentify,
		Response:       string(marshal),
	})

	if err != nil {
		tracing.TraceErr(span, err)
		return nil, fmt.Errorf("failed to store response: %v", err)
	}

	byIP, err := s.services.CommonServices.PostgresRepositories.EnrichDetailsPrefilterTrackingRepository.GetByIP(ctx, ip)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, fmt.Errorf("failed to get stored response: %v", err)
	}

	return byIP, nil
}

func (s *trackingService) askAndStoreSnitcherData(c context.Context, ip string) (*entity.EnrichDetailsTracking, error) {
	span, ctx := opentracing.StartSpanFromContext(c, "TrackingService.askAndStoreSnitcherData")
	defer span.Finish()

	// Create HTTP client
	client := &http.Client{}

	// Create POST request
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/company/find?ip=%s", s.cfg.SnitcherApi.Url, ip), nil)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, fmt.Errorf("failed to create POST request: %v", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.cfg.SnitcherApi.ApiKey)

	//Perform the request
	resp, err := client.Do(req)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, fmt.Errorf("failed to perform POST request: %v", err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	// Parse the JSON request body
	var snitherResponse entity.SnitcherResponseBody
	if err = json.Unmarshal(responseBody, &snitherResponse); err != nil {
		tracing.TraceErr(span, err)
		return nil, fmt.Errorf("failed to unmarshal response body: %v", err)
	}

	var companyName, companyDomain, companyWebsite *string

	if snitherResponse.Company != nil && snitherResponse.Company.Name != "" {
		companyName = &snitherResponse.Company.Name
	}

	if snitherResponse.Company != nil && snitherResponse.Company.Domain != "" {
		companyDomain = &snitherResponse.Company.Domain
	}

	if snitherResponse.Company != nil && snitherResponse.Company.Website != "" {
		companyWebsite = &snitherResponse.Company.Website
	}

	// Store response
	err = s.services.CommonServices.PostgresRepositories.EnrichDetailsTrackingRepository.RegisterRequest(ctx, entity.EnrichDetailsTracking{
		CreatedAt:      utils.Now(),
		IP:             ip,
		CompanyName:    companyName,
		CompanyDomain:  companyDomain,
		CompanyWebsite: companyWebsite,
		Response:       string(responseBody),
	})

	if err != nil {
		tracing.TraceErr(span, err)
		return nil, fmt.Errorf("failed to store response: %v", err)
	}

	byIP, err := s.services.CommonServices.PostgresRepositories.EnrichDetailsTrackingRepository.GetByIP(ctx, ip)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, fmt.Errorf("failed to get stored response: %v", err)
	}

	return byIP, nil
}

func (s *trackingService) notifyOnSlack(c context.Context, r *entity.Tracking) error {
	span, ctx := opentracing.StartSpanFromContext(c, "TrackingService.notifyOnSlack")
	defer span.Finish()

	record, err := s.services.CommonServices.PostgresRepositories.TrackingRepository.GetById(ctx, r.ID)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if record.Notified || record.OrganizationId == nil {
		return nil
	}

	snitcherData, err := s.services.CommonServices.PostgresRepositories.EnrichDetailsTrackingRepository.GetByIP(ctx, record.IP)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if snitcherData == nil {
		tracing.TraceErr(span, errors.New("snitcher record is nil"))
		return nil
	}

	var snitcherDataResponse entity.SnitcherResponseBody
	err = json.Unmarshal([]byte(snitcherData.Response), &snitcherDataResponse)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	var organizationName, organizationLocation, organizationWebsiteUrl, organizationLinkedIn, referrer string

	if record.OrganizationName != nil {
		organizationName = *record.OrganizationName
	} else if snitcherDataResponse.Company.Name != "" {
		organizationName = snitcherDataResponse.Company.Name
	} else {
		organizationName = "Unknown"
	}

	if snitcherDataResponse.Company != nil && snitcherDataResponse.Company.Location != "" {
		organizationLocation = snitcherDataResponse.Company.Location
	} else {
		organizationLocation = "Unknown"
	}

	if snitcherDataResponse.Company != nil && snitcherDataResponse.Company.Website != "" {
		t := strings.Replace(snitcherDataResponse.Company.Website, "https://", "", -1)
		t = strings.Replace(t, "http://", "", -1)
		organizationWebsiteUrl = fmt.Sprintf(`<%s|%s>`, snitcherDataResponse.Company.Website, t)
	} else {
		organizationWebsiteUrl = "Unknown"
	}

	if snitcherDataResponse.Company != nil && snitcherDataResponse.Company.Profiles != nil && snitcherDataResponse.Company.Profiles.Linkedin != nil && snitcherDataResponse.Company.Profiles.Linkedin.Url != "" {
		t := strings.Replace(snitcherDataResponse.Company.Profiles.Linkedin.Url, "https://linkedin.com/companies", "", -1)
		organizationLinkedIn = fmt.Sprintf(`<%s|%s>`, snitcherDataResponse.Company.Profiles.Linkedin.Url, t)
	} else {
		organizationLinkedIn = "Unknown"
	}

	if record.Referrer != "" {
		referrer = record.Referrer
	} else {
		referrer = "Direct"
	}

	slackBlock := `
						[
							{
								"type": "header",
								"text": {
									"type": "plain_text",
									"text": "A visitor from {placeholder_organization_name} is on your website",
									"emoji": true
								}
							},
							{
								"type": "divider"
							},
							{
								"type": "section",
								"text": {
									"type": "mrkdwn",
									"text": "*Location:* {placeholder_location}\n*Website:* {placeholder_website}\n*LinkedIn:* {placeholder_linkedin}\n*Source:* {placeholder_referrer}"
								}
							},
							{
								"type": "divider"
							},
							{
								"type": "actions",
								"elements": [
									{
										"type": "button",
										"text": {
											"type": "plain_text",
											"text": "Open in CustomerOS"
										},
										"url": "{placeholder_view_organization_url}",
										"value": "click_me_123",
										"action_id": "actionId-0"
									}
								]
							}
						]`

	slackBlock = strings.Replace(slackBlock, "{placeholder_organization_name}", organizationName, -1)
	slackBlock = strings.Replace(slackBlock, "{placeholder_location}", organizationLocation, -1)
	slackBlock = strings.Replace(slackBlock, "{placeholder_website}", organizationWebsiteUrl, -1)
	slackBlock = strings.Replace(slackBlock, "{placeholder_linkedin}", organizationLinkedIn, -1)
	slackBlock = strings.Replace(slackBlock, "{placeholder_referrer}", referrer, -1)
	slackBlock = strings.Replace(slackBlock, "{placeholder_view_organization_url}", "https://app.customeros.ai/organization/"+*record.OrganizationId+"?tab=about", -1)

	slackChannels, err := s.services.CommonServices.PostgresRepositories.SlackChannelNotificationRepository.GetSlackChannels(ctx, record.Tenant, "REVEAL-AI")
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	for _, slackChannel := range slackChannels {

		//do not notify old tracking records
		if slackChannel.CreatedAt.After(record.CreatedAt) {
			err := s.services.CommonServices.PostgresRepositories.TrackingRepository.MarkAsNotified(ctx, record.ID)
			if err != nil {
				tracing.TraceErr(span, err)
				return err
			}

			continue
		}

		err = s.sendSlackMessage(ctx, slackChannel.ChannelId, slackBlock)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}

		err := s.services.CommonServices.PostgresRepositories.TrackingRepository.MarkAsNotified(ctx, record.ID)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}
	}

	return nil
}

func (s *trackingService) sendSlackMessage(c context.Context, channel, blocks string) error {
	span, _ := opentracing.StartSpanFromContext(c, "TrackingService.sendSlackMessage")
	defer span.Finish()

	// Create HTTP client
	client := &http.Client{}

	requestBody := map[string]interface{}{
		"channel":      channel,
		"unfurl_links": false,
		"unfurl_media": false,
		"blocks":       blocks,
	}

	// Marshal the request body
	requestBodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		tracing.TraceErr(span, err)
		return fmt.Errorf("failed to marshal request body: %v", err)
	}

	span.LogFields(log.String("request.body", string(requestBodyBytes)))

	// Create POST request
	req, err := http.NewRequest("POST", "https://slack.com/api/chat.postMessage", bytes.NewBuffer(requestBodyBytes))
	if err != nil {
		tracing.TraceErr(span, err)
		return fmt.Errorf("failed to create POST request: %v", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.cfg.SlackBotApiKey)

	//Perform the request
	resp, err := client.Do(req)
	if err != nil {
		tracing.TraceErr(span, err)
		return fmt.Errorf("failed to perform POST request: %v", err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		tracing.TraceErr(span, err)
		return fmt.Errorf("failed to read response body: %v", err)
	}

	span.LogFields(log.String("response.body", string(responseBody)))

	return nil
}