package enrichment

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/clients"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"github.com/opentracing/opentracing-go/log"
	"io"
	"net/http"
	"strings"
	"time"

	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	"github.com/forPelevin/gomoji"

	"github.com/pkg/errors"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
)

const (
	// better contact cached results set to 29 days,
	// this is to avoid the cache from being used 30 days retry
	BetterContactTTL = 29 * 24 * time.Hour
)

type BetterContactResponseBody struct {
	Success bool   `json:"success"`
	ID      string `json:"id"`
	Message string `json:"message"`
}

type BetterContactRequestBody struct {
	Data              []BetterContactData `json:"data"`
	Webhook           string              `json:"webhook"`
	EnrichPhoneNumber bool                `json:"enrich_phone_number"`
}

type BetterContactData struct {
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	LinkedInUrl   string `json:"linkedin_url"`
	Company       string `json:"company"`
	CompanyDomain string `json:"company_domain"`
}

func (s *enrichmentService) FindWorkEmailWithBetterContact(ctx context.Context, linkedInUrl, firstName, lastName, companyName, companyDomain string, enrichPhoneNumber bool) (string, string, *postgres_entity.BetterContactResponseBody, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "EnrichmentService.FindWorkEmailWithBetterContact")
	defer spans.Finish()
	spans.LogFields(
		log.String("linkedInUrl", linkedInUrl),
		log.String("firstName", firstName),
		log.String("lastName", lastName),
		log.String("companyDomain", companyDomain),
		log.String("companyName", companyName),
		log.Bool("enrichPhoneNumber", enrichPhoneNumber))

	// validate if bettercontact is configured
	if s.config.BetterContactConfig.ApiKey == "" || s.config.BetterContactConfig.Url == "" {
		err := errors.New("bettercontact is not configured")
		spans.TraceError(err)
		s.log.Errorf("bettercontact is not configured")
		return "", "", nil, err
	}

	firstName = strings.TrimSpace(firstName)
	lastName = strings.TrimSpace(lastName)
	companyName = strings.TrimSpace(companyName)
	companyDomain = strings.TrimSpace(companyDomain)
	linkedInUrl = strings.TrimSpace(linkedInUrl)

	// replace special characters
	firstName = utils.NormalizeString(firstName)
	lastName = utils.NormalizeString(lastName)
	companyName = utils.NormalizeString(companyName)
	companyDomain = utils.NormalizeString(companyDomain)

	// strip special characters
	firstName = gomoji.RemoveEmojis(firstName)
	lastName = gomoji.RemoveEmojis(lastName)
	companyName = gomoji.RemoveEmojis(companyName)
	companyDomain = gomoji.RemoveEmojis(companyDomain)

	var existingBetterContactData *postgres_entity.EnrichDetailsBetterContact

	detailsBetterContactList, err := s.postgres.EnrichDetailsBetterContactRepository.GetByRequestParams(ctx, linkedInUrl, firstName, lastName, companyName, companyDomain, enrichPhoneNumber, BetterContactTTL)
	if err != nil {
		spans.TraceError(errors.Wrap(err, "failed to get better contact details"))
		return "", "", nil, err
	}

	if detailsBetterContactList != nil && len(detailsBetterContactList) > 0 {
		existingBetterContactData = detailsBetterContactList[0]
	}

	if existingBetterContactData != nil {
		if existingBetterContactData.Response != "" {
			var responseBody postgres_entity.BetterContactResponseBody
			err := json.Unmarshal([]byte(existingBetterContactData.Response), &responseBody)
			if err != nil {
				spans.TraceError(err)
				return "", "", nil, fmt.Errorf("failed to unmarshal response body: %v", err)
			}
			spans.LogKV("result.bettercontact_request_id", existingBetterContactData.RequestID)
			return existingBetterContactData.ID, existingBetterContactData.RequestID, &responseBody, nil
		} else if existingBetterContactData.RequestID != "" {
			spans.LogKV("result.bettercontact_request_id", existingBetterContactData.RequestID)
			return existingBetterContactData.ID, existingBetterContactData.RequestID, nil, nil
		}
	}

	requestBodyDtls := BetterContactRequestBody{}

	requestBodyDtls.Data = []BetterContactData{
		{
			FirstName:     firstName,
			LastName:      lastName,
			LinkedInUrl:   linkedInUrl,
			Company:       companyName,
			CompanyDomain: companyDomain,
		},
	}

	requestBodyDtls.Webhook = s.config.BetterContactConfig.CallbackUrl
	requestBodyDtls.EnrichPhoneNumber = enrichPhoneNumber

	// Marshal request body to JSON
	requestBody, err := json.Marshal(requestBodyDtls)
	if err != nil {
		spans.TraceError(errors.Wrap(err, "failed to marshal bettercontact request body"))
		return "", "", nil, err
	}

	// Create POST request
	req, err := http.NewRequest("POST", fmt.Sprintf("%s?api_key=%s", s.config.BetterContactConfig.Url, s.config.BetterContactConfig.ApiKey), bytes.NewBuffer(requestBody))
	if err != nil {
		spans.TraceError(errors.Wrap(err, "failed to create bettercontact POST request"))
		return "", "", nil, err
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")

	clientTimeout := 15 * time.Second
	httpClient := clients.NewLoggingClient(s.warehouse.APICallLogRepository, enum.VendorBetterContact, &clientTimeout)

	// Perform the request
	resp, err := httpClient.Do(req)
	if err != nil {
		spans.TraceError(errors.Wrap(err, "failed to perform bettercontact POST request"))
		return "", "", nil, err
	}
	defer resp.Body.Close()

	// Decode response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		spans.LogKV("response_body", string(body))
		spans.TraceError(errors.Wrap(err, "failed to read bettercontact response body"))
		return "", "", nil, err
	}

	var responseBody BetterContactResponseBody
	err = json.Unmarshal(body, &responseBody)
	if err != nil {
		spans.LogKV("response_body", string(body))
		spans.TraceError(errors.Wrap(err, "failed to decode bettercontact response body"))
		return "", "", nil, err
	}

	if responseBody.ID == "" {
		spans.LogKV("response_body", string(body))
		err = errors.New("missing bettercontact response id")
		spans.TraceError(err)
		return "", "", nil, err
	}

	dbRecord, err := s.postgres.EnrichDetailsBetterContactRepository.RegisterRequest(ctx, postgres_entity.EnrichDetailsBetterContact{
		RequestID:          responseBody.ID,
		ContactFirstName:   firstName,
		ContactLastName:    lastName,
		ContactLinkedInUrl: linkedInUrl,
		CompanyName:        companyName,
		CompanyDomain:      companyDomain,
		EnrichPhoneNumber:  enrichPhoneNumber,
		Request:            string(requestBody),
	})
	if err != nil {
		spans.TraceError(err)
		return "", "", nil, err
	}

	spans.LogKV("result.bettercontact_request_id", responseBody.ID)
	return dbRecord.ID, responseBody.ID, nil, nil
}
