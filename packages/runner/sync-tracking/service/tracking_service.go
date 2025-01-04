package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/data_fields"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service/security"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	validationmodel "github.com/openline-ai/openline-customer-os/packages/server/validation-api/model"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"

	"github.com/openline-ai/openline-customer-os/packages/runner/sync-tracking/config"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-tracking/constants"
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
	span, ctx := opentracing.StartSpanFromContext(ctx, "TrackingService.ProcessNewRecords")
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
	span, ctx := opentracing.StartSpanFromContext(ctx, "TrackingService.ProcessIPDataRequests")
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
	span, ctx := opentracing.StartSpanFromContext(ctx, "TrackingService.ProcessIPDataResponses")
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
	span, ctx := opentracing.StartSpanFromContext(ctx, "TrackingService.IdentifyTrackingRecords")
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
	span, ctx := opentracing.StartSpanFromContext(ctx, "TrackingService.CreateOrganizationsFromTrackedData")
	defer span.Finish()

	identifiedRecords, err := s.services.CommonServices.PostgresRepositories.TrackingRepository.GetIdentifiedWithDistinctIP(ctx, 100)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	for _, r := range identifiedRecords {
		innerCtx := common.WithCustomContext(ctx, &common.CustomContext{
			Tenant:    r.Tenant,
			AppSource: constants.AppTracking,
		})
		innerErr, done := s.createOrganizationFromTrackingRecord(innerCtx, r)
		if done {
			tracing.TraceErr(span, innerErr)
			return innerErr
		}
	}

	return nil
}

func (s *trackingService) createOrganizationFromTrackingRecord(ctx context.Context, r *entity.Tracking) (error, bool) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TrackingService.createOrganizationFromTrackingRecord")
	defer span.Finish()

	record, err := s.services.CommonServices.PostgresRepositories.TrackingRepository.GetById(ctx, r.ID)
	if err != nil {
		tracing.TraceErr(span, err)
		return err, true
	}

	if record.State != entity.TrackingIdentificationStateIdentified {
		span.LogFields(log.String("skip", "bad state"))
		return nil, false
	}

	snitcherData, err := s.services.CommonServices.PostgresRepositories.EnrichDetailsTrackingRepository.GetByIP(ctx, record.IP)
	if err != nil {
		tracing.TraceErr(span, err)
		return err, true
	}

	if snitcherData == nil {
		tracing.TraceErr(span, errors.New("snitcher record is nil"))
		return nil, false
	}

	if snitcherData.CompanyDomain == nil || *snitcherData.CompanyDomain == "" || utils.IsValidDomain(*snitcherData.CompanyDomain) == false {
		tracing.TraceErr(span, errors.New("company domain is empty or not valid"))
		return nil, false
	}
	span.LogFields(log.String("snitcher.company_domain", *snitcherData.CompanyDomain))
	span.LogFields(log.String("snitcher.company_website", utils.StringOrEmpty(snitcherData.CompanyWebsite)))

	organizationByDomainNode, err := s.services.CommonServices.Neo4jRepositories.OrganizationReadRepository.GetOrganizationByDomain(ctx, nil, record.Tenant, *snitcherData.CompanyDomain)
	if err != nil {
		tracing.TraceErr(span, err)
		return err, true
	}

	if organizationByDomainNode == nil {

		// Save organization
		organizationFields := data_fields.OrganizationFields{
			Name:         snitcherData.CompanyName,
			Website:      snitcherData.CompanyWebsite,
			LeadSource:   utils.StringPtr("Reveal AI"),
			Relationship: utils.ToPtr(neo4jenum.OrganizationRelationshipProspect),
			Stage:        utils.ToPtr(neo4jenum.Lead),
			Domains:      []string{*snitcherData.CompanyDomain},
			Source:       utils.StringPtr(constants.SourceOpenline),
		}
		orgId, err := s.services.CommonServices.OrganizationService.Save(ctx, nil, nil, organizationFields)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to save organization"))
			return err, true
		}
		if orgId == "" {
			tracing.TraceErr(span, errors.New("organization id is nil"))
			return nil, true
		}

		err = s.services.CommonServices.PostgresRepositories.TrackingRepository.MarkAsOrganizationCreated(ctx, record.ID, orgId, snitcherData.CompanyName, snitcherData.CompanyDomain, snitcherData.CompanyWebsite)
		if err != nil {
			tracing.TraceErr(span, err)
			return err, true
		}

		err = s.services.CommonServices.PostgresRepositories.TrackingRepository.MarkAllExcludeIdWithState(ctx, record.ID, record.IP, entity.TrackingIdentificationStateOrganizationExists)
		if err != nil {
			tracing.TraceErr(span, err)
			return err, true
		}
	} else {
		organizationId := utils.GetStringPropOrEmpty(organizationByDomainNode.Props, "id")

		err = s.services.CommonServices.PostgresRepositories.TrackingRepository.MarkAsOrganizationCreated(ctx, record.ID, organizationId, snitcherData.CompanyName, snitcherData.CompanyDomain, snitcherData.CompanyWebsite)
		if err != nil {
			tracing.TraceErr(span, err)
			return err, true
		}

		err = s.services.CommonServices.PostgresRepositories.TrackingRepository.MarkAllWithState(ctx, record.IP, entity.TrackingIdentificationStateOrganizationExists)
		if err != nil {
			tracing.TraceErr(span, err)
			return err, true
		}
	}
	return nil, false
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
	tracing.LogObjectAsJson(span, "request", request)

	// Check ip address is valid
	if net.ParseIP(request.IP) == nil {
		// not a valid IP
		innerErr := s.services.CommonServices.PostgresRepositories.EnrichDetailsPrefilterTrackingRepository.RegisterResponse(ctx, request.IP, false, "invalid IP address", "")
		if innerErr != nil {
			tracing.TraceErr(span, innerErr)
			return fmt.Errorf("failed to store response: %s", innerErr.Error())
		}
		return nil
	}

	ipLookupResponse, err := s.callVerifyAPIForIpData(ctx, request.IP)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to call verify API"))
		return fmt.Errorf("failed to call verify API: %s", err.Error())
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

	// if data by ip is not found, or older than 90 days, ask snitcher
	if snitcherByIp == nil || snitcherByIp.UpdatedAt.Before(utils.Now().AddDate(0, 0, -90)) {
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

	// Perform the request
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
	err = s.services.CommonServices.PostgresRepositories.EnrichDetailsTrackingRepository.Save(ctx, entity.EnrichDetailsTracking{
		CreatedAt:      utils.Now(),
		IP:             ip,
		CompanyName:    companyName,
		CompanyDomain:  companyDomain,
		CompanyWebsite: companyWebsite,
		Response:       string(responseBody),
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, fmt.Errorf("failed to store response: %w", err)
	}

	byIP, err := s.services.CommonServices.PostgresRepositories.EnrichDetailsTrackingRepository.GetByIP(ctx, ip)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, fmt.Errorf("failed to get stored response: %w", err)
	}

	return byIP, nil
}

func (s *trackingService) callVerifyAPIForIpData(ctx context.Context, ipAddress string) (*validationmodel.IpLookupResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TrackingService.callVerifyAPIForIpData")
	defer span.Finish()
	span.LogKV("ipAddress", ipAddress)

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
