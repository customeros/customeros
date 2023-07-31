package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/machinebox/graphql"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/runner/customer-os-dedup/config"
	"github.com/openline-ai/openline-customer-os/packages/runner/customer-os-dedup/logger"
	"github.com/openline-ai/openline-customer-os/packages/runner/customer-os-dedup/repository"
	"github.com/openline-ai/openline-customer-os/packages/runner/customer-os-dedup/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"math"
	"net/http"
	"strings"
	"time"
)

const maxFetchOrganizations = 10000

type OrganizationService interface {
	DedupOrganizations()
}

type organizationService struct {
	cfg           *config.Config
	log           logger.Logger
	repositories  *repository.Repositories
	graphqlClient *graphql.Client
}

func NewOrganizationService(cfg *config.Config, log logger.Logger, repositories *repository.Repositories, graphqlClient *graphql.Client) OrganizationService {
	return &organizationService{
		cfg:           cfg,
		log:           log,
		repositories:  repositories,
		graphqlClient: graphqlClient,
	}
}

type OrganizationIdName struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}
type DuplicatePair struct {
	First      OrganizationIdName `json:"first"`
	Second     OrganizationIdName `json:"second"`
	Confidence float64            `json:"confidence"`
}
type DuplicatePairs struct {
	Duplicates []DuplicatePair `json:"duplicates"`
}
type OrgsCompareResponse struct {
	Primary    string  `json:"primary"`
	Secondary  string  `json:"secondary"`
	Confidence float64 `json:"confidence"`
}

func (s *organizationService) DedupOrganizations() {
	ctx, cancel := context.WithCancel(context.Background())

	defer cancel() // Cancel context on exit

	tenants, err := s.getTenantsWithOrganizations(ctx)
	if err != nil {
		s.log.Errorf("Failed to get tenants for organizations dedup: %v", err)
		return
	} else {
		s.log.Infof("Got %d tenants for organizations dedup", len(tenants))
	}

	// Long-running dedup
	for _, tenant := range tenants {
		// Check if context is cancelled
		select {
		case <-ctx.Done():
			return
		default:
			// Continue fetching organizations
		}
		s.dedupTenantOrganizations(ctx, tenant)
	}
}

func (s *organizationService) dedupTenantOrganizations(ctx context.Context, tenant string) {
	span, ctx := tracing.StartTracerSpan(ctx, "OrganizationService.dedupTenantOrganizations")
	defer span.Finish()
	span.LogFields(log.String("tenant", tenant))

	lastDedupAt, err := s.getLastDedupTime(ctx, tenant)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Failed to get last dedup time: %v", err)
		return
	}
	skip, err := s.checkSkipDedupOrgsForTenant(ctx, tenant, lastDedupAt)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Failed to check skip dedup orgs for tenant %s: %v", tenant, err)
		return
	}
	span.LogFields(log.Bool("skipDedup", skip))
	if skip {
		s.log.Infof("Skipping dedup for tenant %s", tenant)
		return
	}

	orgs, err := s.repositories.OrganizationRepository.GetOrganizationsForDedupComparison(ctx, tenant, maxFetchOrganizations)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Failed to get organizations for dedup comparison: %v", err)
		return
	}

	for i := 0; i < len(orgs); i += s.cfg.Organizations.OrganizationsPerPrompt {
		end := int(math.Min(float64(len(orgs)), float64(i+s.cfg.Organizations.OrganizationsPerPrompt)))
		batch := orgs[i:end]
		err = s.dedupBatchOfTenantOrganizations(ctx, tenant, batch)
		if err != nil {
			tracing.TraceErr(span, err)
			s.log.Errorf("Failed to dedup batch of organizations: %v", err)
			return
		}
	}

	err = s.repositories.TenantRepository.UpdateTenantMetadataOrgDedupAt(ctx, tenant, utils.Now())
	if err != nil {
		s.log.Errorf("Failed to update tenant metadata org dedup at: %v", err)
		tracing.TraceErr(span, err)
	}
}

func (s *organizationService) dedupBatchOfTenantOrganizations(ctx context.Context, tenant string, orgs []*dbtype.Node) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationService.dedupBatchOfTenantOrganizations")
	defer span.Finish()
	span.LogFields(log.String("tenant", tenant), log.Int("orgsBatchSize", len(orgs)))

	orgIdNames := make([]OrganizationIdName, 0, len(orgs))
	for _, v := range orgs {
		id := utils.GetStringPropOrEmpty(v.Props, "id")
		name := utils.GetStringPropOrEmpty(v.Props, "name")
		if id != "" && name != "" {
			orgIdNames = append(orgIdNames, OrganizationIdName{
				Id:   id,
				Name: name,
			})
		}
	}
	if len(orgIdNames) == 0 {
		return nil
	}
	data, err := json.Marshal(orgIdNames)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Failed to marshal orgs to json: %v", err)
		return err
	}
	jsonStr := string(data)

	var duplicatedOrgsByName DuplicatePairs
	jsonResponse, err := s.invokeAIForNamesCheck(ctx, jsonStr)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Failed to invoke AI for names check: %v", err)
		return err
	}

	err = json.Unmarshal([]byte(jsonResponse), &duplicatedOrgsByName)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Failed to unmarshal Anthropic response: %s . error: %v", jsonResponse, err)
		return err
	}

	for _, v := range duplicatedOrgsByName.Duplicates {
		err = s.dedupTwoOrgsByDetails(ctx, tenant, v.First.Id, v.Second.Id)
		if err != nil {
			tracing.TraceErr(span, err)
			s.log.Errorf("Failed to compare two orgs for duplicate: %v", err)
		}
	}

	return nil
}

func (s *organizationService) dedupTwoOrgsByDetails(ctx context.Context, tenant string, org1Id, org2Id string) error {
	span, ctx := opentracing.StartSpanFromContext(context.Background(), "OrganizationService.dedupTwoOrgsByDetails")
	defer span.Finish()
	span.LogFields(log.String("org1", org1Id), log.String("org2", org2Id))

	orgsAlreadyCompared, err := s.repositories.OrganizationRepository.OrgsAlreadyComparedForDuplicates(ctx, tenant, org1Id, org2Id)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Failed to check if orgs already compared for duplicates: %v", err)
		return err
	}
	if orgsAlreadyCompared {
		s.log.Infof("Orgs already compared for duplicates: %s, %s", org1Id, org2Id)
		return nil
	}
	org1Str := s.getOrganizationDetailsAsString(ctx, tenant, org1Id)
	org2Str := s.getOrganizationDetailsAsString(ctx, tenant, org2Id)
	if org1Str == "" || org2Str == "" {
		err = errors.New(fmt.Sprintf("Failed to get org details for dedup %s, %s", org1Id, org2Id))
		tracing.TraceErr(span, err)
		s.log.Errorf("Failed to get org details for dedup %s, %s", org1Id, org2Id)
		return err
	}
	jsonResponse, err := s.invokeAIForOrgsCompare(ctx, org1Id, org1Str, org2Id, org2Str)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Failed to invoke AI for orgs compare: %v", err)
		return err
	}

	var orgsCompareResponse OrgsCompareResponse
	err = json.Unmarshal([]byte(jsonResponse), &orgsCompareResponse)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Failed to unmarshal Anthropic response: %s . error: %v", jsonResponse, err)
		return err
	}
	if orgsCompareResponse.Primary == "" || orgsCompareResponse.Secondary == "" {
		s.log.Infof("No duplicate found for orgs: %s, %s", org1Id, org2Id)
		return nil
	}
	err = s.repositories.OrganizationRepository.SuggestOrganizationsMerge(ctx, tenant, orgsCompareResponse.Primary, orgsCompareResponse.Secondary,
		"Anthropic", s.cfg.Organizations.Anthropic.PromptCompareOrgs, orgsCompareResponse.Confidence)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Failed to save to DB suggested orgs merge: %v", err)
		return err
	}

	return nil
}

func (s *organizationService) getTenantsWithOrganizations(ctx context.Context) ([]string, error) {
	return s.repositories.TenantRepository.GetTenantsWithOrganizations(ctx, s.cfg.Organizations.AtLeastPerTenant)
}

func (s *organizationService) addHeadersToGraphRequest(req *graphql.Request, tenant string) {
	req.Header.Add("X-Openline-API-KEY", s.cfg.Service.CustomerOsAdminAPIKey)
	if tenant != "" {
		req.Header.Add("X-Openline-TENANT", tenant)
	}
}

func (s *organizationService) getOrganizationDetailsAsString(ctx context.Context, tenant, organizationId string) string {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationService.getOrganizationDetailsAsString")
	defer span.Finish()
	span.LogFields(log.String("tenant", tenant), log.String("organizationId", organizationId))

	graphqlRequest := graphql.NewRequest(
		`query Organization($id: ID!) {
  				 organization(id: $id) {
    id 
    name
    description
    domains
    website 
    industry
    subIndustry
    industryGroup
    targetAudience
    valueProposition
    lastFundingRound
    lastFundingAmount
    isPublic
    market
    employees
    socials {
      id  
      url
    }
    emails {
      rawEmail
      email
    }
    phoneNumbers {
      rawPhoneNumber
      e164
    }
    contacts(pagination: {page: 1, limit: 20}) {
    	content {
        name
        firstName
        lastName
        emails {
          rawEmail
          email
        }
        phoneNumbers {
          rawPhoneNumber
          e164
        }
      }
      
    }
    tags {
      name
    }
    subsidiaries {
      type
      organization {
        id
        name
      }
    }
    subsidiaryOf {
      organization {
        id
        name
      }
    }
    externalLinks {
      type
      externalId
      externalUrl
      externalSource
    }
    healthIndicator {
      id
      name
    }
    locations {
      rawAddress
      country
      region
      district
      locality
      street
      address
      address2
      zip
      addressType
      houseNumber
      postalCode
      commercial
      predirection
      latitude
      longitude
      timeZone
      utcOffset
    }
    customFields {
      name
      datatype
      value
    }
    source
  }
			}`)
	graphqlRequest.Var("id", organizationId)
	s.addHeadersToGraphRequest(graphqlRequest, tenant)

	var graphqlResponse interface{}
	if err := s.graphqlClient.Run(ctx, graphqlRequest, &graphqlResponse); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Failed query cosApi for organization details :%v", err.Error())
		return ""
	}
	jsonBytes, err := json.Marshal(graphqlResponse)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Failed to marshal graphql response: %v", err)
		return ""
	}
	return string(jsonBytes)
}

func (s *organizationService) getLastDedupTime(ctx context.Context, tenant string) (time.Time, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationService.getLastSyncTime")
	defer span.Finish()
	span.LogFields(log.String("tenant", tenant))

	node, err := s.repositories.TenantRepository.GetTenantMetadata(ctx, tenant)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Failed to get tenant metadata: %v", err)
		return time.Time{}, err
	}
	return utils.GetTimePropOrEpochStart(node.Props, "orgDedupAt"), nil
}

func (s *organizationService) checkSkipDedupOrgsForTenant(ctx context.Context, tenant string, at time.Time) (bool, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationService.checkSkipDedupOrgsForTenant")
	defer span.Finish()
	span.LogFields(log.String("tenant", tenant), log.Object("lastDedupAt", at))

	now := time.Now()
	diffInDays := now.Sub(at).Hours() / 24
	if diffInDays > float64(s.cfg.Organizations.ForceDedupEachDays) {
		return false, nil
	}

	newOrgsFound, err := s.repositories.OrganizationRepository.ExistsNewOrganizationsCreatedAfter(ctx, tenant, at)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Failed to get new organizations created after %v: %v", at, err)
		return false, err
	}
	return !newOrgsFound, nil
}

func (s *organizationService) invokeAIForNamesCheck(ctx context.Context, str string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationService.invokeAIForNamesCheck")
	defer span.Finish()

	if s.cfg.Organizations.Anthropic.Enabled {
		span.LogFields(log.String("usedAI", "anthropic"))
		prompt := fmt.Sprintf(s.cfg.Organizations.Anthropic.PromptSuggestNames, str)
		response, err := s.invokeAnthropic(ctx, prompt)
		if err != nil {
			tracing.TraceErr(span, err)
			s.log.Errorf("Failed to invoke Anthropic: %v", err.Error())
			return "", err
		}
		jsonResponse := s.extractJsonValueFromAIResponse(response)
		s.log.Infof("Got suggested org pairs from Anthropic: %s", jsonResponse)
		return jsonResponse, nil
	} else if s.cfg.Organizations.OpenAI.Enabled {
		span.LogFields(log.String("usedAI", "openai"))
		prompt := fmt.Sprintf(s.cfg.Organizations.OpenAI.PromptSuggestNames, str)
		response, err := s.invokeOpenAI(ctx, prompt)
		if err != nil {
			tracing.TraceErr(span, err)
			s.log.Errorf("Failed to invoke OpenAI: %v", err.Error())
			return "", err
		}
		jsonResponse := s.extractJsonValueFromAIResponse(response)
		s.log.Infof("Got suggested org pairs from OpenAI: %s", jsonResponse)
		return jsonResponse, nil
	}
	return "", nil
}

func (s *organizationService) invokeAIForOrgsCompare(ctx context.Context, id1 string, dtls1 string, id2 string, dtls2 string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationService.invokeAIForNamesCheck")
	defer span.Finish()

	if s.cfg.Organizations.Anthropic.Enabled {
		span.LogFields(log.String("usedAI", "anthropic"))
		prompt := fmt.Sprintf(s.cfg.Organizations.Anthropic.PromptCompareOrgs, id1, dtls1, id2, dtls2)
		response, err := s.invokeAnthropic(ctx, prompt)
		if err != nil {
			tracing.TraceErr(span, err)
			s.log.Errorf("Failed to invoke Anthropic: %v", err.Error())
			return "", err
		}
		jsonResponse := s.extractJsonValueFromAIResponse(response)
		s.log.Infof("Got response for possible duplicate from Anthropic: %s", jsonResponse)
		return jsonResponse, nil
	} else if s.cfg.Organizations.OpenAI.Enabled {
		span.LogFields(log.String("usedAI", "openai"))
		prompt := fmt.Sprintf(s.cfg.Organizations.OpenAI.PromptCompareOrgs, id1, dtls1, id2, dtls2)
		response, err := s.invokeOpenAI(ctx, prompt)
		if err != nil {
			tracing.TraceErr(span, err)
			s.log.Errorf("Failed to invoke OpenAI: %v", err.Error())
			return "", err
		}
		jsonResponse := s.extractJsonValueFromAIResponse(response)
		s.log.Infof("Got response for possible duplicate from OpenAI: %s", jsonResponse)
		return jsonResponse, nil
	}
	return "", nil
}

func (s *organizationService) invokeAnthropic(ctx context.Context, prompt string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationService.invokeAnthropic")
	defer span.Finish()
	span.LogFields(log.String("prompt", prompt))

	reqBody := map[string]interface{}{
		"prompt": prompt,
		"model":  "claude-2",
	}

	jsonBody, _ := json.Marshal(reqBody)
	reqReader := bytes.NewReader(jsonBody)

	req, err := http.NewRequest("POST", s.cfg.Service.Anthropic.ApiPath+"/ask", reqReader)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error creating request: %v", err.Error())
		return "", err
	}
	req.Header.Set("content-type", "application/json")
	req.Header.Set("X-Openline-API-KEY", s.cfg.Service.Anthropic.ApiKey)

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error sending request to Anthropic API %s: error - %v", s.cfg.Service.Anthropic.ApiPath+"/ask", err.Error())
		return "", err
	}
	defer resp.Body.Close()

	// Print summarized email
	var data map[string]string
	json.NewDecoder(resp.Body).Decode(&data)
	result := strings.TrimSpace(data["completion"])
	span.LogFields(log.String("response", result))

	return result, nil
}

func (s *organizationService) invokeOpenAI(ctx context.Context, prompt string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationService.invokeOpenAI")
	defer span.Finish()
	span.LogFields(log.String("prompt", prompt))

	requestData := map[string]interface{}{}
	requestData["model"] = "gpt-3.5-turbo"
	requestData["prompt"] = prompt

	requestBody, _ := json.Marshal(requestData)
	request, _ := http.NewRequest("POST", s.cfg.Service.OpenAI.ApiPath+"/ask", strings.NewReader(string(requestBody)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-Openline-API-KEY", s.cfg.Service.OpenAI.ApiKey)

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error making the API request: %v", err.Error())
		return "", err
	}
	if response.StatusCode != 200 {
		err = errors.New("Error making the OpenAI API request")
		tracing.TraceErr(span, err)
		s.log.Errorf("Error making the API request: %s", response.Status)
		return "", err
	}
	defer response.Body.Close()

	// Print summarized email
	var data map[string]interface{}
	json.NewDecoder(response.Body).Decode(&data)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error parsing the API response: %v", err.Error())
		return "", err
	}
	choices := data["choices"].([]interface{})
	result := ""
	if len(choices) > 0 {
		result = strings.TrimSpace(choices[0].(map[string]interface{})["message"].(map[string]interface{})["content"].(string))
	}
	span.LogFields(log.String("response", result))

	return result, nil
}

func (s *organizationService) extractJsonValueFromAIResponse(str string) string {
	start := strings.IndexByte(str, '{')
	if start == -1 {
		return ""
	}

	end := strings.LastIndexByte(str, '}')
	if end == -1 {
		return ""
	}

	return str[start : end+1]
}
