package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/forPelevin/gomoji"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	postgresentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/enrichment-api/config"
	"github.com/openline-ai/openline-customer-os/packages/server/enrichment-api/logger"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"net/http"
	"strings"
)

type BettercontactService interface {
	FindWorkEmail(ctx context.Context, linkedInUrl, firstName, lastName, companyName, companyDomain string) (string, *postgresentity.BetterContactResponseBody, error)
}

type bettercontactService struct {
	config   *config.Config
	services *Services
	log      logger.Logger
}

func NewBettercontactService(config *config.Config, services *Services, log logger.Logger) BettercontactService {
	return &bettercontactService{
		config:   config,
		services: services,
		log:      log,
	}
}

type BetterContactRequestBody struct {
	Data    []BetterContactData `json:"data"`
	Webhook string              `json:"webhook"`
}

type BetterContactData struct {
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	LinkedInUrl   string `json:"linkedin_url"`
	Company       string `json:"company"`
	CompanyDomain string `json:"company_domain"`
}

type BetterContactResponseBody struct {
	Success bool   `json:"success"`
	ID      string `json:"id"`
	Message string `json:"message"`
}

func (s bettercontactService) FindWorkEmail(ctx context.Context, linkedInUrl, firstName, lastName, companyName, companyDomain string) (string, *postgresentity.BetterContactResponseBody, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "BettercontactService.FindWorkEmail")
	defer span.Finish()
	span.LogFields(log.String("linkedInUrl", linkedInUrl), log.String("firstName", firstName), log.String("lastName", lastName), log.String("companyDomain", companyDomain), log.String("companyName", companyName))

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

	var existingBetterContactData *postgresentity.EnrichDetailsBetterContact

	if linkedInUrl != "" {
		betterContactByLinkedInUrl, err := s.services.CommonServices.PostgresRepositories.EnrichDetailsBetterContactRepository.GetByLinkedInUrl(ctx, linkedInUrl)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to get better contact details"))
			return "", nil, err
		}

		if betterContactByLinkedInUrl != nil {
			existingBetterContactData = betterContactByLinkedInUrl
		}
	} else {
		detailsBetterContactList, err := s.services.CommonServices.PostgresRepositories.EnrichDetailsBetterContactRepository.GetBy(ctx, firstName, lastName, companyName, companyDomain)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to get better contact details"))
			return "", nil, err
		}

		if detailsBetterContactList != nil && len(detailsBetterContactList) > 0 {
			existingBetterContactData = detailsBetterContactList[0]
		}
	}

	if existingBetterContactData != nil {
		if existingBetterContactData.Response != "" {
			var responseBody postgresentity.BetterContactResponseBody
			err := json.Unmarshal([]byte(existingBetterContactData.Response), &responseBody)
			if err != nil {
				tracing.TraceErr(span, err)
				return "", nil, fmt.Errorf("failed to unmarshal response body: %v", err)
			}
			return existingBetterContactData.ID.String(), &responseBody, nil
		}
		return existingBetterContactData.ID.String(), nil, nil
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

	// Marshal request body to JSON
	requestBody, err := json.Marshal(requestBodyDtls)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to marshal bettercontact request body"))
		return "", nil, err
	}

	// Create HTTP client
	client := &http.Client{}

	// Create POST request
	req, err := http.NewRequest("POST", fmt.Sprintf("%s?api_key=%s", s.config.BetterContactConfig.Url, s.config.BetterContactConfig.ApiKey), bytes.NewBuffer(requestBody))
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to create bettercontact POST request"))
		return "", nil, err
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")

	//Perform the request
	resp, err := client.Do(req)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to perform bettercontact POST request"))
		return "", nil, err
	}
	defer resp.Body.Close()

	//Decode response body
	var responseBody BetterContactResponseBody
	err = json.NewDecoder(resp.Body).Decode(&responseBody)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to decode bettercontact response body"))
		return "", nil, err
	}

	dbRecord, err := s.services.CommonServices.PostgresRepositories.EnrichDetailsBetterContactRepository.RegisterRequest(ctx, postgresentity.EnrichDetailsBetterContact{
		RequestID:          responseBody.ID,
		ContactFirstName:   firstName,
		ContactLastName:    lastName,
		ContactLinkedInUrl: linkedInUrl,
		CompanyName:        companyName,
		CompanyDomain:      companyDomain,
		Request:            string(requestBody),
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return "", nil, err
	}

	return dbRecord.ID.String(), nil, nil
}
