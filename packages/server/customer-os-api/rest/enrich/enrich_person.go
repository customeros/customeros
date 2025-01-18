// @openapi 3.0.0
package restenrich

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	commontracing "github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	postgresentity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	postgresrepository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"

	"github.com/customeros/customeros/packages/server/customer-os-api/enum"
	"github.com/customeros/customeros/packages/server/customer-os-api/rest"
	cosapi_services "github.com/customeros/customeros/packages/server/customer-os-api/services"
	"github.com/customeros/customeros/packages/server/customer-os-api/tracing"
)

const (
	emailTypePersonal       = "personal"
	emailTypeWork           = "work"
	enrichPersonAcceptedUrl = "/enrich/v1/person/results"
)

// EnrichPersonResponse represents the response for person enrichment
// @Description Response structure for person enrichment operations
type EnrichPersonResponse struct { // Inherits standard response fields
	enum.BaseResponse

	// Optional message providing additional information
	// required: false
	Message string `json:"message,omitempty" example:"Enrichment completed"`

	// Indicates if all enrichment operations are complete
	// required: true
	IsComplete bool `json:"isComplete" example:"true"`

	// List of fields still being processed
	// required: false
	PendingFields []string `json:"pendingFields,omitempty" example:"email,phone number"`

	// URL to check the final result when processing is incomplete
	// required: false
	// format: uri
	ResultURL string `json:"resultUrl,omitempty" example:"https://api.customeros.ai/enrich/v1/person/results/550e8400-e29b-41d4-a716-446655440000"`

	// Enriched person data
	// required: true
	Data EnrichPersonData `json:"data"`
}

// EnrichPersonData represents enriched person information
// @Description Comprehensive enriched information about a person
type EnrichPersonData struct {
	// List of email addresses associated with the person
	// required: false
	Emails []EnrichPersonEmail `json:"emails"`

	// Employment history
	// required: false
	Jobs []EnrichPersonJob `json:"jobs"`

	// Geographic location information
	// required: false
	Location EnrichPersonLocation `json:"location"`

	// Person's name information
	// required: true
	Name EnrichPersonName `json:"name"`

	// List of phone numbers
	// required: false
	PhoneNumbers []EnrichPersonPhoneNumber `json:"phoneNumbers"`

	// URL to person's profile picture
	// required: false
	// format: uri
	ProfilePic string `json:"profilePic" example:"https://example.com/profile.jpg"`

	// Social media presence
	// required: false
	Social EnrichPersonSocial `json:"social"`
}

// EnrichPersonEmail represents email information
// @Description Email address with validation details
type EnrichPersonEmail struct {
	// Email address
	// required: true
	// format: email
	Address string `json:"address" example:"john.doe@example.com"`

	// Indicates if the email is deliverable
	// required: false
	Deliverable *string `json:"deliverable,omitempty" example:"true"`

	// Indicates if the email is considered risky
	// required: false
	IsRisky *bool `json:"isRisky,omitempty" example:"false"`

	// Type of email address
	// required: false
	// enum: personal,work
	Type *string `json:"type,omitempty" example:"work"`
}

// EnrichPersonJob represents employment information
// @Description Details about a person's job position
type EnrichPersonJob struct {
	// Job title
	// required: true
	Title string `json:"title" example:"Software Engineer"`

	// Seniority level
	// required: false
	// enum: Junior,Mid-Level,Senior,Lead,Manager,Director,VP,C-Level
	Seniority string `json:"seniority" example:"Senior"`

	// Employment duration
	// required: true
	Duration EnrichPersonJobDuration `json:"duration"`

	// Company name
	// required: true
	Company string `json:"company" example:"Tech Corp"`

	// Company's LinkedIn URL
	// required: false
	// format: uri
	CompanyLinkedin string `json:"companyLinkedin" example:"https://linkedin.com/company/techcorp"`

	// Company's website
	// required: false
	// format: uri
	CompanyWebsite string `json:"companyWebsite" example:"https://techcorp.com"`

	// Indicates if this is the current position
	// required: true
	IsCurrent bool `json:"isCurrent" example:"true"`
}

// EnrichPersonJobDuration represents employment duration
// @Description Time period of employment
type EnrichPersonJobDuration struct {
	// Starting month (1-12)
	// required: false
	// minimum: 1
	// maximum: 12
	StartMonth *int `json:"startMonth,omitempty" example:"1"`

	// Starting year
	// required: false
	// minimum: 1900
	// maximum: 2100
	StartYear *int `json:"startYear,omitempty" example:"2020"`

	// Ending month (1-12)
	// required: false
	// minimum: 1
	// maximum: 12
	EndMonth *int `json:"endMonth,omitempty" example:"12"`

	// Ending year
	// required: false
	// minimum: 1900
	// maximum: 2100
	EndYear *int `json:"endYear,omitempty" example:"2023"`
}

// EnrichPersonLocation represents location information
// @Description Geographic and timezone information about a person
type EnrichPersonLocation struct {
	// City name
	// required: false
	City string `json:"city" example:"San Francisco"`

	// State or region
	// required: false
	Region string `json:"region" example:"California"`

	// Country name
	// required: false
	Country string `json:"country" example:"United States"`

	// Timezone identifier
	// required: false
	// example: America/Los_Angeles
	Timezone string `json:"timezone" example:"PST"`
}

// EnrichPersonName represents name information
// @Description Person's name details
type EnrichPersonName struct {
	// First name
	// required: true
	// minLength: 1
	FirstName string `json:"firstName" example:"John"`

	// Last name
	// required: true
	// minLength: 1
	LastName string `json:"lastName" example:"Doe"`

	// Full name (typically firstName + lastName)
	// required: false
	FullName string `json:"fullName" example:"John Doe"`
}

// EnrichPersonPhoneNumber represents phone information
// @Description Phone number with type classification
type EnrichPersonPhoneNumber struct {
	// Phone number in E.164 format
	// required: true
	// pattern: ^\+[1-9]\d{1,14}$
	Number string `json:"number" example:"+14155552671"`

	// Type of phone number
	// required: true
	// enum: mobile,work,home,other
	Type string `json:"type" example:"mobile"`
}

// EnrichPersonSocial represents social media presence
// @Description Collection of social media profile information
type EnrichPersonSocial struct {
	// LinkedIn profile information
	// required: false
	Linkedin EnrichPersonLinkedIn `json:"linkedin"`

	// X (Twitter) profile information
	// required: false
	X EnrichPersonX `json:"x"`

	// GitHub profile information
	// required: false
	Github EnrichPersonGithub `json:"github"`

	// Discord profile information
	// required: false
	Discord EnrichPersonDiscord `json:"discord"`
}

// EnrichPersonLinkedIn represents LinkedIn profile information
// @Description LinkedIn specific profile details
type EnrichPersonLinkedIn struct {
	// LinkedIn internal ID
	// required: false
	ID string `json:"id" example:"123456789"`

	// LinkedIn public identifier
	// required: false
	PublicID string `json:"publicId" example:"john-doe"`

	// Full LinkedIn profile URL
	// required: false
	// format: uri
	URL string `json:"url" example:"https://linkedin.com/in/john-doe"`

	// Number of LinkedIn followers
	// required: false
	// minimum: 0
	FollowerCount int `json:"followerCount" example:"500"`
}

// EnrichPersonX represents X (Twitter) profile information
// @Description X (formerly Twitter) profile details
type EnrichPersonX struct {
	// X handle (without @)
	// required: true
	Handle string `json:"handle" example:"johndoe"`

	// Full X profile URL
	// required: false
	// format: uri
	URL string `json:"url" example:"https://x.com/johndoe"`
}

// EnrichPersonGithub represents GitHub profile information
// @Description GitHub profile details
type EnrichPersonGithub struct {
	// GitHub username
	// required: true
	Username string `json:"username" example:"johndoe"`

	// Full GitHub profile URL
	// required: false
	// format: uri
	URL string `json:"url" example:"https://github.com/johndoe"`
}

// EnrichPersonDiscord represents Discord profile information
// @Description Discord profile details
type EnrichPersonDiscord struct {
	// Discord username with discriminator
	// required: true
	// pattern: ^.{3,32}#[0-9]{4}$
	Username string `json:"username" example:"johndoe#1234"`
}

// @Summary Enrich person information
// @Description Enriches person information using LinkedIn URL, email, and other optional details
// @Tags Enrichment API
// @Accept json
// @Produce json
// @Param linkedinUrl query string false "LinkedIn profile URL" example(https://linkedin.com/in/johndoe)
// @Param email query string false "Email address" example(john.doe@example.com) format(email)
// @Param firstName query string false "First name" example(John) minLength(1)
// @Param lastName query string false "Last name" example(Doe) minLength(1)
// @Param includeMobileNumber query bool false "Include mobile number in results" default(false)
// @Success 200 {object} EnrichPersonResponse "Successfully retrieved enriched data"
// @Success 202 {object} EnrichPersonResponse "Processing initiated, check ResultURL for final data"
// @Success 200 {object} rest.ErrorResponse "Person not found (status: warning)"
// @Failure 400 {object} rest.BaseResponse "Missing linkedinUrl or email"
// @Failure 401 {object} rest.BaseResponse "Missing or invalid API key"
// @Failure 500 {object} rest.BaseResponse "Internal server error"
// @Router /enrich/v1/person [get]
// @Security ApiKeyAuth

func EnrichPerson(services *cosapi_services.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "EnrichPerson", c.Request.Header)
		defer span.Finish()
		commontracing.TagComponentRest(span)
		commontracing.TagTenant(span, common.GetTenantFromContext(ctx))

		tenant := common.GetTenantFromContext(ctx)
		if tenant == "" {
			handlers.SendError(c, span, http.StatusUnauthorized, enum.ErrUnauthorized)
			return
		}

		linkedinUrl := c.Query("linkedinUrl")
		email := c.Query("email")
		firstName := c.Query("firstName")
		lastName := c.Query("lastName")
		enrichPhoneNumber := c.Query("includeMobileNumber") == "true"

		if strings.TrimSpace(linkedinUrl) == "" && strings.TrimSpace(email) == "" {
			handlers.SendError(c, span, http.StatusBadRequest, enum.ErrBadRequest.WithMessage("Missing linkedinUrl or email"))
			return
		}

		span.LogFields(
			log.String("request.email", email),
			log.String("request.linkedinUrl", linkedinUrl),
			log.Bool("request.includeMobileNumber", enrichPhoneNumber),
			log.String("request.firstName", firstName),
			log.String("request.lastName", lastName))

		// Call enrichPerson API
		person := interfaces.PersonSearch{
			LinkedinURL: &linkedinUrl,
			FirstName:   &firstName,
			LastName:    &lastName,
			Email:       &email,
		}
		personDbID, enrichPersonResponse, err := services.CommonServices.EnrichmentService.EnrichPerson(ctx, person)
		if err != nil {
			handlers.SendError(c, span, http.StatusInternalServerError, enum.ErrInternalServer)
			return
		}
		if enrichPersonResponse == nil {
			c.JSON(http.StatusNotFound,
				enum.ErrorResponse{
					BaseResponse: enum.BuildBaseResponse(enum.StatusWarning),
					Message:      "Person not found",
				})
			return
		}

		// Map the response
		enrichedPersonData := mapPersonScrapInData(enrichPersonResponse)

		// Call findWorkEmail API
		companyName, companyDomain := "", ""
		if enrichPersonResponse.Company != nil {
			companyName = enrichPersonResponse.Company.Name
			companyDomain, _ = services.CommonServices.DomainService.GetPrimaryDomainForOrganizationWebsite(ctx, enrichPersonResponse.Company.WebsiteUrl)
		}

		dbID, _, findWorkEmailResponse, err := services.CommonServices.EnrichmentService.FindWorkEmail(ctx,
			enrichPersonResponse.Person.LinkedInUrl,
			enrichPersonResponse.Person.FirstName,
			enrichPersonResponse.Person.LastName,
			companyName,
			companyDomain,
			enrichPhoneNumber)
		if err != nil || findWorkEmailResponse == nil {
			handlers.SendError(c, span, http.StatusInternalServerError, enum.ErrInternalServer)
			return
		}

		// Compose response
		response := EnrichPersonResponse{
			BaseResponse: enum.BuildBaseResponse(enum.StatusSuccess),
		}
		if enrichedPersonData != nil {
			response.Data = *enrichedPersonData
		}

		betterContactResponseBody := findWorkEmailResponse.Data
		query := postgresentity.CosApiEnrichPersonTempResult{
			Tenant:                tenant,
			BettercontactRecordId: dbID,
		}
		if personDbID != nil {
			query.ScrapinRecordId = *personDbID
		}
		if betterContactResponseBody == nil {
			dbRecord, err := services.Repositories.PostgresRepositories.CosApiEnrichPersonTempResultRepository.Create(ctx, query)
			if err != nil {
				tracing.TraceErr(span, errors.Wrap(err, "failed to create temp result"))
				handlers.SendError(c, span, http.StatusInternalServerError, enum.ErrInternalServer)
				return
			}
			response.IsComplete = false
			response.PendingFields = []string{"email"}
			if enrichPhoneNumber {
				response.PendingFields = append(response.PendingFields, "phone number")
			}
			response.ResultURL = services.Cfg.CommonServices.Internal.CustomerOsApi.ApiUrl + enrichPersonAcceptedUrl + "/" + dbRecord.ID.String()
		} else {
			response.IsComplete = true
			emailFound, phoneFound := false, false
			for _, item := range betterContactResponseBody {
				if item.ContactEmailAddress != "" {
					emailFound = true
					response.Data.Emails = append(response.Data.Emails, EnrichPersonEmail{
						Address: item.ContactEmailAddress,
					})
				}
				if enrichPhoneNumber {
					if item.ContactPhoneNumber != nil && fmt.Sprintf("%v", item.ContactPhoneNumber) != "" {
						phoneFound = true
						response.Data.PhoneNumbers = append(response.Data.PhoneNumbers, EnrichPersonPhoneNumber{
							Number: fmt.Sprintf("%v", item.ContactPhoneNumber),
							Type:   "mobile",
						})
					}
				}
			}

			// Register billable events
			if emailFound {
				_, err = services.Repositories.PostgresRepositories.ApiBillableEventRepository.RegisterEvent(ctx, tenant,
					postgresentity.BillableEventEnrichPersonEmailFound,
					postgresrepository.BillableEventDetails{
						ExternalID: dbID,
						ReferenceData: fmt.Sprintf("Email: %s, LinkedIn: %s, FirstName: %s, LastName: %s",
							email, linkedinUrl, firstName, lastName),
					})
				if err != nil {
					tracing.TraceErr(span, errors.Wrap(err, "failed to store billable event"))
				}
			}
			if phoneFound {
				_, err = services.Repositories.PostgresRepositories.ApiBillableEventRepository.RegisterEvent(ctx, tenant,
					postgresentity.BillableEventEnrichPersonPhoneFound,
					postgresrepository.BillableEventDetails{
						ExternalID: dbID,
						ReferenceData: fmt.Sprintf("Email: %s, LinkedIn: %s, FirstName: %s, LastName: %s",
							email, linkedinUrl, firstName, lastName),
					})
				if err != nil {
					tracing.TraceErr(span, errors.Wrap(err, "failed to store billable event"))
				}
			}
		}

		// Validate emails
		for i := range response.Data.Emails {
			emailRecord := &response.Data.Emails[i]
			if emailRecord.Address != "" {
				emailValidationResult, err := services.CommonServices.VerifyService.ValidateEmail(ctx, emailRecord.Address)
				if err != nil {
					tracing.TraceErr(span, errors.Wrap(err, "failed to validate email"))
					continue
				}

				emailRecord.Deliverable = utils.StringPtr(emailValidationResult.EmailData.Deliverable)
				emailRecord.IsRisky = utils.BoolPtr(
					emailValidationResult.DomainData.IsFirewalled ||
						emailValidationResult.EmailData.IsRoleAccount ||
						emailValidationResult.EmailData.IsSystemGenerated ||
						emailValidationResult.EmailData.IsFreeAccount ||
						emailValidationResult.EmailData.IsMailboxFull ||
						!emailValidationResult.DomainData.IsPrimaryDomain)

				if emailValidationResult.EmailData.IsFreeAccount {
					emailRecord.Type = utils.StringPtr(emailTypePersonal)
				} else {
					emailRecord.Type = utils.StringPtr(emailTypeWork)
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

// @Summary Retrieve enrichment results
// @Description Retrieves the results of an asynchronous person enrichment operation
// @Tags Enrichment API
// @Accept json
// @Produce json
// @Param id path string true "Result ID" format(uuid)
// @Success 200 {object} EnrichPersonResponse "Successfully retrieved enriched data"
// @Success 202 {object} EnrichPersonResponse "Still processing, check again later"
// @Failure 400 {object} rest.BaseResponse "Invalid result ID"
// @Failure 401 {object} rest.BaseResponse "Missing or invalid API key"
// @Failure 404 {object} rest.BaseResponse "Result not found"
// @Failure 500 {object} rest.BaseResponse "Internal server error"
// @Router /enrich/v1/person/results/{id} [get]
// @Security ApiKeyAuth
func EnrichPersonCallback(services *cosapi_services.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "EnrichPerson", c.Request.Header)
		defer span.Finish()
		commontracing.TagComponentRest(span)
		commontracing.TagTenant(span, common.GetTenantFromContext(ctx))

		tenant := common.GetTenantFromContext(ctx)
		if tenant == "" {
			handlers.SendError(c, span, http.StatusUnauthorized, enum.ErrUnauthorized)
			return
		}

		tempId := c.Param("id")
		span.LogFields(log.String("request.tempId", tempId))

		getTempRecord, err := services.Repositories.PostgresRepositories.CosApiEnrichPersonTempResultRepository.GetById(ctx, tempId, tenant)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to get temp record"))
			handlers.SendError(c, span, http.StatusInternalServerError, enum.ErrInternalServer)
			return
		}
		if getTempRecord == nil {
			handlers.SendError(c, span, http.StatusNotFound, enum.ErrNotFound)
			return
		}

		// Enrich person data
		scrapInDbRecord, err := services.Repositories.PostgresRepositories.EnrichDetailsScrapInRepository.GetById(ctx, getTempRecord.ScrapinRecordId)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to get scrapin record"))
			handlers.SendError(c, span, http.StatusInternalServerError, enum.ErrInternalServer)
			return
		}
		if scrapInDbRecord == nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to get scrapin record"))
			handlers.SendError(c, span, http.StatusNotFound, enum.ErrNotFound)
			return
		}

		betterContactDbRecord, err := services.Repositories.PostgresRepositories.EnrichDetailsBetterContactRepository.GetById(ctx, getTempRecord.BettercontactRecordId)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to get bettercontact record"))
			handlers.SendError(c, span, http.StatusInternalServerError, enum.ErrInternalServer)
			return
		}
		if betterContactDbRecord == nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to get bettercontact record"))
			handlers.SendError(c, span, http.StatusNotFound, enum.ErrNotFound)
			return
		}

		// Compose response
		response := EnrichPersonResponse{
			BaseResponse: enum.BuildBaseResponse(enum.StatusSuccess),
		}

		// extract scrapin data
		var scrapInPersonResponse postgresentity.ScrapInResponseBody
		err = json.Unmarshal([]byte(scrapInDbRecord.Data), &scrapInPersonResponse)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to unmarshal scrapin record"))
			handlers.SendError(c, span, http.StatusInternalServerError, enum.ErrInternalServer)
			return
		}
		enrichedPersonData := mapPersonScrapInData(&scrapInPersonResponse)

		if enrichedPersonData != nil {
			response.Data = *enrichedPersonData
		}

		// extract better contact data
		var betterContactResponseBody *postgresentity.BetterContactResponseBody
		if betterContactDbRecord.Response != "" {
			err = json.Unmarshal([]byte(betterContactDbRecord.Response), &betterContactResponseBody)
			if err != nil {
				tracing.TraceErr(span, errors.Wrap(err, "failed to unmarshal bettercontact record"))
				handlers.SendError(c, span, http.StatusInternalServerError, enum.ErrInternalServer)
				return
			}
		}

		if betterContactResponseBody == nil {
			response.IsComplete = false
			response.PendingFields = []string{"email"}
			if betterContactDbRecord.EnrichPhoneNumber {
				response.PendingFields = append(response.PendingFields, "phone number")
			}
			response.ResultURL = services.Cfg.CommonServices.Internal.CustomerOsApi.ApiUrl + enrichPersonAcceptedUrl + "/" + tempId
		} else {
			response.IsComplete = true
			for _, item := range betterContactResponseBody.Data {
				if item.ContactEmailAddress != "" {
					response.Data.Emails = append(response.Data.Emails, EnrichPersonEmail{
						Address: item.ContactEmailAddress,
					})
				}
				if item.ContactPhoneNumber != nil && fmt.Sprintf("%v", item.ContactPhoneNumber) != "" {
					response.Data.PhoneNumbers = append(response.Data.PhoneNumbers, EnrichPersonPhoneNumber{
						Number: fmt.Sprintf("%v", item.ContactPhoneNumber),
						Type:   "mobile",
					})
				}
			}
		}

		// call email verify
		for i := range response.Data.Emails {
			email := &response.Data.Emails[i] // Get a pointer to the email in the slice
			if email.Address != "" {
				emailValidationResult, err := services.CommonServices.VerifyService.ValidateEmail(ctx, email.Address)
				if err != nil {
					tracing.TraceErr(span, errors.Wrap(err, "failed to validate email"))
					continue
				}
				email.Deliverable = utils.StringPtr(emailValidationResult.EmailData.Deliverable)
				email.IsRisky = utils.BoolPtr(
					emailValidationResult.DomainData.IsFirewalled ||
						emailValidationResult.EmailData.IsRoleAccount ||
						emailValidationResult.EmailData.IsSystemGenerated ||
						emailValidationResult.EmailData.IsFreeAccount ||
						emailValidationResult.EmailData.IsMailboxFull ||
						!emailValidationResult.DomainData.IsPrimaryDomain)
				if emailValidationResult.EmailData.IsFreeAccount {
					email.Type = utils.StringPtr(emailTypePersonal)
				} else {
					email.Type = utils.StringPtr(emailTypeWork)
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

func mapPersonScrapInData(source *postgresentity.ScrapInResponseBody) *EnrichPersonData {
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
			// Seniority:       position.Seniority, // TODO will be implemented later after clarifications
		}
		if position.StartEndDate.Start != nil {
			enrichPersonJob.Duration.StartMonth = &position.StartEndDate.Start.Month
			enrichPersonJob.Duration.StartYear = &position.StartEndDate.Start.Year
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
		// City: source.Data.PersonProfile.Person.Location,
		Region: source.Person.Location,
		// Country:  source.Data.PersonProfile.Person.Country,
		// Timezone: source.Data.PersonProfile.Person.Timezone,
	}

	// set phone numbers
	// TODO implement phone numbers
	return &output
}
