package rest

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/service"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service/security"
	commontracing "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	postgresentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	enrichmentmodel "github.com/openline-ai/openline-customer-os/packages/server/enrichment-api/model"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"net/http"
)

const (
	emailTypePersonal       = "personal"
	emailTypeProfessional   = "professional"
	enrichPersonAcceptedUrl = "/enrich/v1/person/results"
)

type EnrichPersonResponse struct {
	Status        string           `json:"status"`
	Message       string           `json:"message,omitempty"`
	IsComplete    bool             `json:"isComplete"`
	PendingFields []string         `json:"pendingFields,omitempty"`
	ResultURL     string           `json:"resultUrl,omitempty"`
	Data          EnrichPersonData `json:"data"`
}

type EnrichPersonData struct {
	Emails       []EnrichPersonEmail       `json:"emails"`
	Jobs         []EnrichPersonJob         `json:"jobs"`
	Location     EnrichPersonLocation      `json:"location"`
	Name         EnrichPersonName          `json:"name"`
	PhoneNumbers []EnrichPersonPhoneNumber `json:"phoneNumbers"`
	ProfilePic   string                    `json:"profilePic"`
	Social       EnrichPersonSocial        `json:"social"`
}

type EnrichPersonEmail struct {
	Address       string  `json:"address"`
	IsDeliverable *bool   `json:"isDeliverable,omitempty"`
	IsRisky       *bool   `json:"isRisky,omitempty"`
	Type          *string `json:"type,omitempty"`
}

type EnrichPersonJob struct {
	Title           string                  `json:"title"`
	Seniority       string                  `json:"seniority"`
	Duration        EnrichPersonJobDuration `json:"duration"`
	Company         string                  `json:"company"`
	CompanyLinkedin string                  `json:"companyLinkedin"`
	CompanyWebsite  string                  `json:"companyWebsite"`
	IsCurrent       bool                    `json:"isCurrent"`
}

type EnrichPersonJobDuration struct {
	StartMonth int  `json:"startMonth"`
	StartYear  int  `json:"startYear"`
	EndMonth   *int `json:"endMonth"`
	EndYear    *int `json:"endYear"`
}

type EnrichPersonLocation struct {
	City     string `json:"city"`
	Region   string `json:"region"`
	Country  string `json:"country"`
	Timezone string `json:"timezone"`
}

type EnrichPersonName struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	FullName  string `json:"fullName"`
}

type EnrichPersonPhoneNumber struct {
	Number string `json:"number"`
	Type   string `json:"type"`
}

type EnrichPersonSocial struct {
	Linkedin EnrichPersonLinkedIn `json:"linkedin"`
	X        EnrichPersonX        `json:"x"`
	Github   EnrichPersonGithub   `json:"github"`
	Discord  EnrichPersonDiscord  `json:"discord"`
}

type EnrichPersonLinkedIn struct {
	ID            string `json:"id"`
	PublicID      string `json:"publicId"`
	URL           string `json:"url"`
	FollowerCount int    `json:"followerCount"`
}

type EnrichPersonX struct {
	Handle string `json:"handle"`
	URL    string `json:"url"`
}

type EnrichPersonGithub struct {
	Username string `json:"username"`
	URL      string `json:"url"`
}

type EnrichPersonDiscord struct {
	Username string `json:"username"`
}

func EnrichPerson(services *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "EnrichPerson", c.Request.Header)
		defer span.Finish()

		tenant := common.GetTenantFromContext(ctx)
		if tenant == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Missing tenant context"})
			return
		}
		span.SetTag(tracing.SpanTagTenant, common.GetTenantFromContext(ctx))

		linkedinUrl := c.Query("linkedinUrl")
		email := c.Query("email")
		searchPhoneNumberParam := c.Query("phoneNumber")
		// convert to boolean
		searchPhoneNumber := searchPhoneNumberParam == "true"

		// TODO add enrichment based on email + first and last names + other scrapin data
		if linkedinUrl == "" {
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Missing linkedin url parameter"})
			return
		}

		span.LogFields(log.String("request.email", email))
		span.LogFields(log.String("request.linkedinUrl", linkedinUrl))
		span.LogFields(log.Bool("request.searchPhoneNumber", searchPhoneNumber))

		// Call enrichPerson API
		enrichPersonApiResponse, err := callApiEnrichPerson(ctx, services, span, email, linkedinUrl)
		if err != nil || enrichPersonApiResponse == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Internal error"})
			return
		}

		// map the response to the customer-os model
		scrapInPersonResponse := postgresentity.ScrapInPersonResponse{}
		if enrichPersonApiResponse.Data != nil && enrichPersonApiResponse.Data.PersonProfile != nil {
			scrapInPersonResponse = *enrichPersonApiResponse.Data.PersonProfile
		}
		enrichedPersonData := mapScrapInData(&scrapInPersonResponse)

		// Call findWorkEmail API
		findWorkEmailApiResponse, err := callApiFindWorkEmail(ctx, services, span, *enrichPersonApiResponse)
		if err != nil || findWorkEmailApiResponse == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Internal error"})
			return
		}

		// Compose response
		response := EnrichPersonResponse{
			Status: "success",
		}
		if enrichedPersonData != nil {
			response.Data = *enrichedPersonData
		}

		betterContactResponseBody := findWorkEmailApiResponse.Data
		if betterContactResponseBody == nil {
			dbRecord, err := services.CommonServices.PostgresRepositories.CosApiEnrichPersonTempResultRepository.Create(ctx, postgresentity.CosApiEnrichPersonTempResult{
				Tenant:                tenant,
				ScrapinRecordId:       enrichPersonApiResponse.RecordId,
				BettercontactRecordId: findWorkEmailApiResponse.RecordId,
			})
			if err != nil {
				tracing.TraceErr(span, errors.Wrap(err, "failed to create temp result"))
				c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Internal error"})
				return
			}
			response.IsComplete = false
			response.PendingFields = []string{"email"}
			response.ResultURL = services.Cfg.Services.CustomerOsApiUrl + enrichPersonAcceptedUrl + "/" + dbRecord.ID.String()
		} else {
			response.IsComplete = true
			for _, datas := range betterContactResponseBody.Data {
				if datas.ContactEmailAddress != "" {
					response.Data.Emails = append(response.Data.Emails, EnrichPersonEmail{
						Address: datas.ContactEmailAddress,
					})
				}
			}
		}

		// call email verify
		for i := range response.Data.Emails {
			emailRecord := &response.Data.Emails[i]
			if emailRecord.Address != "" {
				emailValidationResult, err := callApiValidateEmail(ctx, services, span, emailRecord.Address)
				if err != nil {
					tracing.TraceErr(span, errors.Wrap(err, "failed to validate email"))
					continue
				}
				if emailValidationResult.Status != "success" {
					tracing.TraceErr(span, errors.New("failed to validate email"))
					continue
				}
				emailRecord.IsDeliverable = utils.BoolPtr(emailValidationResult.Data.EmailData.IsDeliverable)
				emailRecord.IsRisky = utils.BoolPtr(emailValidationResult.Data.DomainData.IsFirewalled ||
					emailValidationResult.Data.EmailData.IsRoleAccount ||
					emailValidationResult.Data.EmailData.IsFreeAccount ||
					emailValidationResult.Data.EmailData.IsMailboxFull ||
					emailValidationResult.Data.DomainData.IsCatchAll)
				if emailValidationResult.Data.EmailData.IsFreeAccount {
					emailRecord.Type = utils.StringPtr(emailTypePersonal)
				} else {
					emailRecord.Type = utils.StringPtr(emailTypeProfessional)
				}
			}
		}

		responseStatusCode := http.StatusOK
		if !response.IsComplete {
			responseStatusCode = http.StatusAccepted
		}
		c.JSON(responseStatusCode, response)
	}
}

func EnrichPersonCallback(services *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "EnrichPerson", c.Request.Header)
		defer span.Finish()

		tenant := common.GetTenantFromContext(ctx)
		if tenant == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Missing tenant context"})
			return
		}
		span.SetTag(tracing.SpanTagTenant, common.GetTenantFromContext(ctx))

		tempId := c.Param("id")
		span.LogFields(log.String("request.tempId", tempId))

		getTempRecord, err := services.CommonServices.PostgresRepositories.CosApiEnrichPersonTempResultRepository.GetById(ctx, tempId, tenant)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to get temp record"))
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Internal error"})
			return
		}
		if getTempRecord == nil {
			c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "Record not found"})
			return
		}

		// Enrich person data
		scrapInDbRecord, err := services.CommonServices.PostgresRepositories.EnrichDetailsScrapInRepository.GetById(ctx, getTempRecord.ScrapinRecordId)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to get scrapin record"))
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Internal error"})
			return
		}
		if scrapInDbRecord == nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to get scrapin record"))
			c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "Record not found"})
			return
		}

		betterContactDbRecord, err := services.CommonServices.PostgresRepositories.EnrichDetailsBetterContactRepository.GetById(ctx, getTempRecord.BettercontactRecordId)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to get bettercontact record"))
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Internal error"})
			return
		}
		if betterContactDbRecord == nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to get bettercontact record"))
			c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "Record not found"})
			return
		}

		// Compose response
		response := EnrichPersonResponse{
			Status: "success",
		}

		// extract scrapin data
		var scrapInPersonResponse postgresentity.ScrapInPersonResponse
		err = json.Unmarshal([]byte(scrapInDbRecord.Data), &scrapInPersonResponse)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to unmarshal scrapin record"))
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Internal error"})
			return
		}
		enrichedPersonData := mapScrapInData(&scrapInPersonResponse)

		if enrichedPersonData != nil {
			response.Data = *enrichedPersonData
		}

		// extract better contact data
		var betterContactResponseBody *postgresentity.BetterContactResponseBody
		if betterContactDbRecord.Response != "" {
			err = json.Unmarshal([]byte(betterContactDbRecord.Response), &betterContactResponseBody)
			if err != nil {
				tracing.TraceErr(span, errors.Wrap(err, "failed to unmarshal bettercontact record"))
				c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Internal error"})
				return
			}
		}

		if betterContactResponseBody == nil {
			response.IsComplete = false
			response.PendingFields = []string{"email"}
			response.ResultURL = services.Cfg.Services.CustomerOsApiUrl + enrichPersonAcceptedUrl + "/" + tempId
		} else {
			response.IsComplete = true
			for _, datas := range betterContactResponseBody.Data {
				if datas.ContactEmailAddress != "" {
					response.Data.Emails = append(response.Data.Emails, EnrichPersonEmail{
						Address: datas.ContactEmailAddress,
					})
				}
			}
		}

		// call email verify
		for _, email := range response.Data.Emails {
			if email.Address != "" {
				emailValidationResult, err := callApiValidateEmail(ctx, services, span, email.Address)
				if err != nil {
					tracing.TraceErr(span, errors.Wrap(err, "failed to validate email"))
					continue
				}
				if emailValidationResult.Status != "success" {
					tracing.TraceErr(span, errors.New("failed to validate email"))
					continue
				}
				email.IsDeliverable = utils.BoolPtr(emailValidationResult.Data.EmailData.IsDeliverable)
				email.IsRisky = utils.BoolPtr(emailValidationResult.Data.DomainData.IsFirewalled ||
					emailValidationResult.Data.EmailData.IsRoleAccount ||
					emailValidationResult.Data.EmailData.IsFreeAccount ||
					emailValidationResult.Data.EmailData.IsMailboxFull ||
					emailValidationResult.Data.DomainData.IsCatchAll)
				if emailValidationResult.Data.EmailData.IsFreeAccount {
					email.Type = utils.StringPtr(emailTypePersonal)
				} else {
					email.Type = utils.StringPtr(emailTypeProfessional)
				}
			}
		}

		responseStatusCode := http.StatusOK
		if !response.IsComplete {
			responseStatusCode = http.StatusAccepted
		}
		c.JSON(responseStatusCode, response)
	}
}

func callApiEnrichPerson(ctx context.Context, services *service.Services, span opentracing.Span, email, linkedinUrl string) (*enrichmentmodel.EnrichPersonResponse, error) {
	requestJSON, err := json.Marshal(enrichmentmodel.EnrichPersonRequest{
		Email:       email,
		LinkedinUrl: linkedinUrl,
	})
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to marshal request"))
		return nil, err
	}
	requestBody := []byte(string(requestJSON))
	req, err := http.NewRequest("GET", services.Cfg.Services.EnrichmentApiUrl+"/enrichPerson", bytes.NewBuffer(requestBody))
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to create request"))
		return nil, err
	}
	// Inject span context into the HTTP request
	req = commontracing.InjectSpanContextIntoHTTPRequest(req, span)

	// Set the request headers
	req.Header.Set(security.ApiKeyHeader, services.Cfg.Services.EnrichmentApiKey)
	req.Header.Set(security.TenantHeader, common.GetTenantFromContext(ctx))

	// Make the HTTP request
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to perform request"))
		return nil, err
	}
	defer response.Body.Close()
	span.LogFields(log.Int("response.status.enrichPerson", response.StatusCode))

	var enrichPersonApiResponse enrichmentmodel.EnrichPersonResponse
	err = json.NewDecoder(response.Body).Decode(&enrichPersonApiResponse)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to decode enrich person response"))
		return nil, err
	}
	return &enrichPersonApiResponse, nil
}

func callApiFindWorkEmail(ctx context.Context, services *service.Services, span opentracing.Span, enrichPersonApiResponse enrichmentmodel.EnrichPersonResponse) (*enrichmentmodel.FindWorkEmailResponse, error) {
	companyName, companyDomain := "", ""
	if enrichPersonApiResponse.Data.PersonProfile.Company != nil {
		companyName = enrichPersonApiResponse.Data.PersonProfile.Company.Name
		companyDomain = services.CommonServices.DomainService.ExtractDomainFromOrganizationWebsite(ctx, enrichPersonApiResponse.Data.PersonProfile.Company.WebsiteUrl)
	}
	requestJSON, err := json.Marshal(enrichmentmodel.FindWorkEmailRequest{
		LinkedinUrl:   enrichPersonApiResponse.Data.PersonProfile.Person.LinkedInUrl,
		FirstName:     enrichPersonApiResponse.Data.PersonProfile.Person.FirstName,
		LastName:      enrichPersonApiResponse.Data.PersonProfile.Person.LastName,
		CompanyName:   companyName,
		CompanyDomain: companyDomain,
	})
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to marshal request"))
		return nil, err
	}
	requestBody := []byte(string(requestJSON))
	req, err := http.NewRequest("GET", services.Cfg.Services.EnrichmentApiUrl+"/findWorkEmail", bytes.NewBuffer(requestBody))
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to create request"))
		return nil, err
	}
	// Inject span context into the HTTP request
	req = commontracing.InjectSpanContextIntoHTTPRequest(req, span)

	// Set the request headers
	req.Header.Set(security.ApiKeyHeader, services.Cfg.Services.EnrichmentApiKey)
	req.Header.Set(security.TenantHeader, common.GetTenantFromContext(ctx))

	// Make the HTTP request
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to perform request"))
		return nil, err
	}
	defer response.Body.Close()
	span.LogFields(log.Int("response.status.findWorkEmail", response.StatusCode))

	var findWorkEmailApiResponse enrichmentmodel.FindWorkEmailResponse
	err = json.NewDecoder(response.Body).Decode(&findWorkEmailApiResponse)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to decode find work email response"))
		return nil, err
	}
	return &findWorkEmailApiResponse, nil
}

func mapScrapInData(source *postgresentity.ScrapInPersonResponse) *EnrichPersonData {
	if source == nil {
		return nil
	}
	output := EnrichPersonData{}

	// set emails
	if source.Email != "" {
		output.Emails = append(output.Emails, EnrichPersonEmail{
			Address: source.Email,
		})
	}

	// set name
	output.Name = EnrichPersonName{
		FirstName: source.Person.FirstName,
		LastName:  source.Person.LastName,
	}
	if source.Person.FirstName != "" && source.Person.LastName != "" {
		output.Name.FullName = source.Person.FirstName + " " + source.Person.LastName
	}

	// set jobs
	for _, position := range source.Person.Positions.PositionHistory {
		enrichPersonJob := EnrichPersonJob{
			Title:           position.Title,
			Company:         position.CompanyName,
			CompanyLinkedin: position.LinkedInUrl,
			IsCurrent:       position.StartEndDate.End == nil,
			//Seniority:       position.Seniority, // TODO will be implemented later after clarifications
			Duration: EnrichPersonJobDuration{
				StartMonth: position.StartEndDate.Start.Month,
				StartYear:  position.StartEndDate.Start.Year,
			},
		}
		if position.StartEndDate.End != nil {
			enrichPersonJob.Duration.EndMonth = &position.StartEndDate.End.Month
			enrichPersonJob.Duration.EndYear = &position.StartEndDate.End.Year
		}
		if source.Company != nil && source.Company.LinkedInId == position.LinkedInId {
			enrichPersonJob.CompanyWebsite = source.Company.WebsiteUrl
		}
		output.Jobs = append(output.Jobs, enrichPersonJob)
	}

	// set profile picture
	output.ProfilePic = source.Person.PhotoUrl

	// set social
	output.Social = EnrichPersonSocial{
		Linkedin: EnrichPersonLinkedIn{
			ID:            source.Person.LinkedInIdentifier,
			PublicID:      source.Person.PublicIdentifier,
			URL:           source.Person.LinkedInUrl,
			FollowerCount: source.Person.FollowerCount,
		},
		// TODO add X, github, other from Customer OS
	}

	// set location // TODO implement AI lookup to get details
	output.Location = EnrichPersonLocation{
		//City: source.Data.PersonProfile.Person.Location,
		Region: source.Person.Location,
		//Country:  source.Data.PersonProfile.Person.Country,
		//Timezone: source.Data.PersonProfile.Person.Timezone,
	}

	// set phone numbers
	// TODO implement phone numbers
	return &output
}
