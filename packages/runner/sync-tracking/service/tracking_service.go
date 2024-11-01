package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-tracking/config"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-tracking/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service/security"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	neo4jmodel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/model"
	neo4jrepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	validationmodel "github.com/openline-ai/openline-customer-os/packages/server/validation-api/model"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"io"
	"net"
	"net/http"
	"strings"
)

type TrackingService interface {
	ProcessNewRecords(ctx context.Context) error
	ProcessIPDataRequests(ctx context.Context) error
	ProcessIPDataResponses(ctx context.Context) error
	IdentifyTrackingRecords(ctx context.Context) error
	CreateOrganizationsFromTrackedData(ctx context.Context) error
	NotifyOnSlack(ctx context.Context)
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

func (s *trackingService) ProcessNewRecords(ctx context.Context) error {
	span, ctx := tracing.StartTracerSpan(ctx, "TrackingService.ProcessNewRecords")
	defer span.Finish()

	newRecords, err := s.services.CommonServices.PostgresRepositories.TrackingRepository.GetNewRecords(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	for _, record := range newRecords {
		err := s.processNewRecord(ctx, record)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}
	}

	return nil
}

func (s *trackingService) ProcessIPDataRequests(ctx context.Context) error {
	span, ctx := tracing.StartTracerSpan(ctx, "TrackingService.ProcessIPDataRequests")
	defer span.Finish()

	sendRequestsRecords, err := s.services.CommonServices.PostgresRepositories.EnrichDetailsPrefilterTrackingRepository.GetForSendingRequests(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	for _, record := range sendRequestsRecords {
		err := s.askAndStoreIPData(ctx, record)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}
	}

	return nil
}

func (s *trackingService) ProcessIPDataResponses(ctx context.Context) error {
	span, ctx := tracing.StartTracerSpan(ctx, "TrackingService.ProcessIPDataResponses")
	defer span.Finish()

	trackingRecordsWithIPData, err := s.services.CommonServices.PostgresRepositories.TrackingRepository.GetForPrefilter(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	for _, record := range trackingRecordsWithIPData {
		err := s.processTrackingRecordWithIPData(ctx, record)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}
	}

	return nil
}

func (s *trackingService) IdentifyTrackingRecords(ctx context.Context) error {
	span, ctx := tracing.StartTracerSpan(ctx, "TrackingService.IdentifyTrackingRecords")
	defer span.Finish()

	notIdentifiedTrackingRecords, err := s.services.CommonServices.PostgresRepositories.TrackingRepository.GetReadyForIdentification(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	for _, record := range notIdentifiedTrackingRecords {
		err := s.processRecordIdentification(ctx, record)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}

	}

	return nil
}

func (s *trackingService) CreateOrganizationsFromTrackedData(ctx context.Context) error {
	span, ctx := tracing.StartTracerSpan(ctx, "TrackingService.CreateOrganizationsFromTrackedData")
	defer span.Finish()

	identifiedRecords, err := s.services.CommonServices.PostgresRepositories.TrackingRepository.GetIdentifiedWithDistinctIP(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	for _, r := range identifiedRecords {

		innerCtx := common.WithCustomContext(ctx, &common.CustomContext{
			Tenant:    r.Tenant,
			AppSource: constants.AppTracking,
		})

		record, err := s.services.CommonServices.PostgresRepositories.TrackingRepository.GetById(innerCtx, r.ID)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}

		if record.State != entity.TrackingIdentificationStateIdentified {
			span.LogFields(log.String("skip", "bad state"))
			continue
		}

		snitcherData, err := s.services.CommonServices.PostgresRepositories.EnrichDetailsTrackingRepository.GetByIP(innerCtx, record.IP)
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
		span.LogFields(log.String("company_domain", *snitcherData.CompanyDomain))
		span.LogFields(log.String("company_website", utils.StringOrEmpty(snitcherData.CompanyWebsite)))

		organizationByDomainNode, err := s.services.CommonServices.Neo4jRepositories.OrganizationReadRepository.GetOrganizationByDomain(innerCtx, record.Tenant, *snitcherData.CompanyDomain)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}

		if organizationByDomainNode == nil {

			// Save organization
			organizationFields := neo4jrepository.OrganizationSaveFields{
				Name:               utils.StringOrEmpty(snitcherData.CompanyName),
				Website:            utils.StringOrEmpty(snitcherData.CompanyWebsite),
				LeadSource:         "Reveal AI",
				Relationship:       neo4jenum.Prospect,
				Stage:              neo4jenum.Lead,
				Domains:            []string{*snitcherData.CompanyDomain},
				UpdateName:         true,
				UpdateWebsite:      true,
				UpdateLeadSource:   true,
				UpdateRelationship: true,
				UpdateStage:        true,
				SourceFields: neo4jmodel.SourceFields{
					Source:    constants.SourceOpenline,
					AppSource: constants.AppTracking,
				},
			}
			orgIdPtr, err := s.services.CommonServices.OrganizationService.Save(innerCtx, nil, record.Tenant, nil, &organizationFields)
			if err != nil {
				tracing.TraceErr(span, errors.Wrap(err, "failed to save organization"))
				return err
			}
			if orgIdPtr == nil {
				tracing.TraceErr(span, errors.New("organization id is nil"))
				return nil
			}

			err = s.services.CommonServices.PostgresRepositories.TrackingRepository.MarkAsOrganizationCreated(innerCtx, record.ID, *orgIdPtr, snitcherData.CompanyName, snitcherData.CompanyDomain, snitcherData.CompanyWebsite)
			if err != nil {
				tracing.TraceErr(span, err)
				return err
			}

			err = s.services.CommonServices.PostgresRepositories.TrackingRepository.MarkAllExcludeIdWithState(innerCtx, record.ID, record.IP, entity.TrackingIdentificationStateOrganizationExists)
			if err != nil {
				tracing.TraceErr(span, err)
				return err
			}
		} else {
			organizationId := utils.GetStringPropOrEmpty(organizationByDomainNode.Props, "id")

			err = s.services.CommonServices.PostgresRepositories.TrackingRepository.MarkAsOrganizationCreated(innerCtx, record.ID, organizationId, snitcherData.CompanyName, snitcherData.CompanyDomain, snitcherData.CompanyWebsite)
			if err != nil {
				tracing.TraceErr(span, err)
				return err
			}

			err = s.services.CommonServices.PostgresRepositories.TrackingRepository.MarkAllWithState(innerCtx, record.IP, entity.TrackingIdentificationStateOrganizationExists)
			if err != nil {
				tracing.TraceErr(span, err)
				return err
			}
		}
	}

	return nil
}

func (s *trackingService) NotifyOnSlack(ctx context.Context) {
	span, ctx := tracing.StartTracerSpan(ctx, "TrackingService.NotifyOnSlack")
	defer span.Finish()

	if s.cfg.SlackBotApiKey == "" {
		span.LogFields(log.String("skip", "no slack bot api key"))
		return
	}

	limit := 100
	notifyOnSlackRecords, err := s.services.CommonServices.PostgresRepositories.TrackingRepository.GetForSlackNotifications(ctx, limit)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to get records for slack notifications"))
		s.services.Logger.Errorf("failed to get records for slack notifications: %s", err.Error())
		return
	}

	if len(notifyOnSlackRecords) == 0 {
		span.LogFields(log.String("skip", "no records to notify"))
		return
	}

	for _, r := range notifyOnSlackRecords {
		err := s.notifyOnSlack(ctx, r)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to notify on slack"))
		}
	}

	return
}

func (s *trackingService) processNewRecord(c context.Context, newRecord *entity.Tracking) error {
	span, ctx := opentracing.StartSpanFromContext(c, "TrackingService.processNewRecord")
	defer span.Finish()

	ipDataByIp, err := s.services.CommonServices.PostgresRepositories.EnrichDetailsPrefilterTrackingRepository.GetByIP(ctx, newRecord.IP)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if ipDataByIp == nil {
		span.LogFields(log.String("result", "registering ip data request"))
		err := s.services.CommonServices.PostgresRepositories.EnrichDetailsPrefilterTrackingRepository.RegisterRequest(ctx, newRecord.IP)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}

		err = s.services.CommonServices.PostgresRepositories.TrackingRepository.SetStateById(ctx, newRecord.ID, entity.TrackingIdentificationStatePrefilteredAsked)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}
	} else {
		if ipDataByIp.ShouldIdentify == nil {
			err = s.services.CommonServices.PostgresRepositories.TrackingRepository.SetStateById(ctx, newRecord.ID, entity.TrackingIdentificationStatePrefilteredAsked)
			if err != nil {
				tracing.TraceErr(span, err)
				return err
			}
			return nil
		}

		state := entity.TrackingIdentificationStatePrefilteredFail
		if *ipDataByIp.ShouldIdentify {
			state = entity.TrackingIdentificationStatePrefilteredPass
		}

		err = s.services.CommonServices.PostgresRepositories.TrackingRepository.SetStateById(ctx, newRecord.ID, state)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}

		span.LogFields(log.String("result", "processed"))
	}

	return nil
}

func (s *trackingService) askAndStoreIPData(c context.Context, request *entity.EnrichDetailsPreFilterTracking) error {
	span, ctx := opentracing.StartSpanFromContext(c, "TrackingService.askAndStoreIPData")
	defer span.Finish()

	ipLookupResponse, err := s.callVerifyAPIForIpData(ctx, request.IP)
	if err != nil {
		tracing.TraceErr(span, err)
		return fmt.Errorf("failed to call verify API: %v", err)
	}

	shouldIdentify := true
	skipIdenitifyReason := ""

	if ipLookupResponse.IpData == nil {
		shouldIdentify = false
		skipIdenitifyReason = "ip data is nil"
	} else if ipLookupResponse.IpData.Ip == "" {
		skipIdenitifyReason = "ip is empty"
		shouldIdentify = false
	} else if ipLookupResponse.IpData.Carrier != nil {
		skipIdenitifyReason = "carrier is not nil"
		shouldIdentify = false
	} else if ipLookupResponse.IpData.Threat.IsTor {
		skipIdenitifyReason = "tor detected"
		shouldIdentify = false
	} else if ipLookupResponse.IpData.Threat.IsIcloudRelay {
		skipIdenitifyReason = "icloud relay detected"
		shouldIdentify = false
	} else if ipLookupResponse.IpData.Threat.IsProxy {
		skipIdenitifyReason = "proxy detected"
		shouldIdentify = false
	} else if ipLookupResponse.IpData.Threat.IsDatacenter {
		skipIdenitifyReason = "datacenter detected"
		shouldIdentify = false
	} else if ipLookupResponse.IpData.Threat.IsAnonymous {
		skipIdenitifyReason = "anonymous detected"
		shouldIdentify = false
	} else if ipLookupResponse.IpData.Threat.IsKnownAttacker {
		skipIdenitifyReason = "known attacker detected"
		shouldIdentify = false
	} else if ipLookupResponse.IpData.Threat.IsKnownAbuser {
		skipIdenitifyReason = "known abuser detected"
		shouldIdentify = false
	} else if ipLookupResponse.IpData.Threat.IsThreat {
		skipIdenitifyReason = "threat detected"
		shouldIdentify = false
	} else if ipLookupResponse.IpData.Threat.IsBogon {
		skipIdenitifyReason = "bogon detected"
		shouldIdentify = false
	}

	if ipLookupResponse.IpData != nil {
		marshal, err := json.Marshal(ipLookupResponse.IpData)
		if err != nil {
			tracing.TraceErr(span, err)
			return fmt.Errorf("failed to marshal response body: %v", err)
		}

		err = s.services.CommonServices.PostgresRepositories.EnrichDetailsPrefilterTrackingRepository.RegisterResponse(ctx, request.IP, shouldIdentify, skipIdenitifyReason, string(marshal))
		if err != nil {
			tracing.TraceErr(span, err)
			return fmt.Errorf("failed to store response: %v", err)
		}
	}

	return nil
}

func (s *trackingService) processTrackingRecordWithIPData(c context.Context, record *entity.Tracking) error {
	span, ctx := opentracing.StartSpanFromContext(c, "TrackingService.processTrackingRecordWithIPData")
	defer span.Finish()

	ipDataByIp, err := s.services.CommonServices.PostgresRepositories.EnrichDetailsPrefilterTrackingRepository.GetByIP(ctx, record.IP)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if ipDataByIp == nil {
		tracing.TraceErr(span, errors.New("ip data record is nil"))
		return nil
	}

	if ipDataByIp.ShouldIdentify == nil {
		tracing.TraceErr(span, errors.New("should identify is nil"))
		return nil
	}

	state := entity.TrackingIdentificationStatePrefilteredFail
	if *ipDataByIp.ShouldIdentify {
		state = entity.TrackingIdentificationStatePrefilteredPass
	}

	err = s.services.CommonServices.PostgresRepositories.TrackingRepository.SetStateById(ctx, record.ID, state)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}

func (s *trackingService) processRecordIdentification(c context.Context, record *entity.Tracking) error {
	span, ctx := opentracing.StartSpanFromContext(c, "TrackingService.processRecordIdentification")
	defer span.Finish()

	snitcherByIp, err := s.services.CommonServices.PostgresRepositories.EnrichDetailsTrackingRepository.GetByIP(ctx, record.IP)
	if err != nil {
		tracing.TraceErr(span, err)
		return fmt.Errorf("failed to get better contact details: %v", err)
	}

	if snitcherByIp == nil {
		snitcherByIp, err = s.askAndStoreSnitcherData(ctx, record.IP)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}
	}

	if snitcherByIp == nil {
		tracing.TraceErr(span, errors.New("snitcher record is nil"))
		return err
	}

	state := entity.TrackingIdentificationStateNotIdentified

	if snitcherByIp.CompanyDomain != nil && *snitcherByIp.CompanyDomain != "" {
		state = entity.TrackingIdentificationStateIdentified
	}

	err = s.services.CommonServices.PostgresRepositories.TrackingRepository.SetStateById(ctx, record.ID, state)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return nil
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
		span.LogFields(log.String("snitcherResponseBody", string(responseBody)))
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
		tracing.TraceErr(span, errors.Wrap(err, "failed to get tracking record"))
		return err
	}
	if record.Tenant != "" {
		tracing.TagTenant(span, record.Tenant)
	}

	err = s.services.CommonServices.PostgresRepositories.TrackingRepository.IncrementNotificationTry(ctx, record.ID)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to increment notification try"))
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

	//skip notification if the identified company domain is the same as the workspace domain in the tenant ( basically skip employees from triggering notifications)
	if record.OrganizationDomain != nil && *record.OrganizationDomain != "" {

		workspaceNodeList, err := s.services.CommonServices.Neo4jRepositories.WorkspaceReadRepository.Get(ctx, record.Tenant)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to get workspace nodes"))
			return err
		}

		for _, workspaceNode := range workspaceNodeList {
			props := utils.GetPropsFromNode(*workspaceNode)
			domainName := utils.GetStringPropOrEmpty(props, "name")
			if domainName == *record.OrganizationDomain {
				err := s.services.CommonServices.PostgresRepositories.TrackingRepository.MarkAsNotified(ctx, record.ID)
				if err != nil {
					tracing.TraceErr(span, err)
					return err
				}

				span.LogFields(log.String("skip", "workspace is the same as organization domain"))
				return nil
			}
		}
	}

	var organizationName, organizationLocation, organizationWebsiteUrl, organizationLinkedIn, referrer string

	if record.OrganizationName != nil {
		organizationName = *record.OrganizationName
	} else if snitcherDataResponse.Company.Name != "" {
		organizationName = snitcherDataResponse.Company.Name
	} else {
		organizationName = "Unknown"
	}

	organizationLocation = snitcherDataResponse.LocationToString()

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

		err = s.sendSlackMessage(ctx, slackChannel.Tenant, slackChannel.ChannelId, slackBlock)
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

func (s *trackingService) sendSlackMessage(ctx context.Context, tenant, channel, blocks string) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "TrackingService.sendSlackMessage")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	span.LogFields(log.String("channel", channel))

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
		tracing.TraceErr(span, errors.Wrap(err, "failed to marshal request body"))
		return fmt.Errorf("failed to marshal request body: %v", err)
	}

	span.LogFields(log.String("request.body", string(requestBodyBytes)))

	// Create POST request
	req, err := http.NewRequest("POST", "https://slack.com/api/chat.postMessage", bytes.NewBuffer(requestBodyBytes))
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to create POST request"))
		return fmt.Errorf("failed to create POST request: %v", err)
	}

	botApiKey := ""
	// prepare bot key
	slackSettings, err := s.services.CommonServices.PostgresRepositories.SlackSettingsRepository.Get(ctx, tenant)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to get slack settings"))
	}
	if slackSettings == nil {
		span.LogFields(log.String("skip", "slack settings not found"))
		s.services.Logger.Warnf("slack settings not found for tenant %s", tenant)
		return nil
	} else {
		botApiKey = slackSettings.AccessToken
	}

	// display last first 8 and last 3 chars
	maskedBotApiKey := ""
	if len(botApiKey) > 11 {
		maskedBotApiKey = botApiKey[:8] + "..." + botApiKey[len(botApiKey)-3:]
	}
	span.LogFields(log.String("bot.api.key", maskedBotApiKey))

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+botApiKey)

	//Perform the request
	resp, err := client.Do(req)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to perform POST request"))
		return fmt.Errorf("failed to perform POST request: %v", err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to read response body"))
		return fmt.Errorf("failed to read response body: %v", err)
	}

	span.LogFields(log.String("response.body", string(responseBody)))

	return nil
}

func (s *trackingService) callVerifyAPIForIpData(ctx context.Context, ipAddress string) (*validationmodel.IpLookupResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TrackingService.callVerifyAPIForIpData")
	defer span.Finish()

	if net.ParseIP(ipAddress) == nil {
		err := errors.New("invalid IP address")
		tracing.TraceErr(span, err)
		return nil, err
	}

	requestJSON, err := json.Marshal(validationmodel.IpLookupRequest{
		Ip: ipAddress,
	})
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to marshal request"))
		return nil, err
	}
	requestBody := []byte(string(requestJSON))
	req, err := http.NewRequest("POST", s.cfg.ValidationApi.Url+"/ipLookup", bytes.NewBuffer(requestBody))
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to create request"))
		return nil, err
	}
	// Inject span context into the HTTP request
	req = tracing.InjectSpanContextIntoHTTPRequest(req, span)

	// Set the request headers
	req.Header.Set(security.ApiKeyHeader, s.cfg.ValidationApi.ApiKey)
	req.Header.Set(security.TenantHeader, "")

	// Make the HTTP request
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to perform request"))
		return nil, err
	}
	defer response.Body.Close()

	var result validationmodel.IpLookupResponse
	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to decode response"))
		return nil, err
	}

	return &result, nil
}
