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
	enrichmentmodel "github.com/openline-ai/openline-customer-os/packages/server/enrichment-api/model"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"net/http"
)

type EnrichPersonResponse struct {
	Status        string           `json:"status"`
	Message       string           `json:"message,omitempty"`
	IsComplete    bool             `json:"isComplete"`
	PendingFields []string         `json:"pendingFields,omitempty"`
	ResultURL     *string          `json:"resultUrl,omitempty"`
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
	Address       string `json:"address"`
	IsDeliverable bool   `json:"isDeliverable"`
	IsRisky       bool   `json:"isRisky"`
	Type          string `json:"type"`
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

		// Check if address is provided
		linkedinUrl := c.Query("linkedinUrl")
		email := c.Query("email")
		searchPhoneNumberParam := c.Query("phoneNumber")
		// convert to boolean
		searchPhoneNumber := searchPhoneNumberParam == "true"

		if email == "" && linkedinUrl == "" {
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Missing email and linkedin url parameters"})
			return
		}

		span.LogFields(log.String("request.email", email))
		span.LogFields(log.String("request.linkedinUrl", linkedinUrl))
		span.LogFields(log.Bool("request.searchPhoneNumber", searchPhoneNumber))

		requestJSON, err := json.Marshal(enrichmentmodel.EnrichPersonRequest{
			Email:       email,
			LinkedinUrl: linkedinUrl,
		})
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to marshal request"))
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Internal error"})
			return
		}
		requestBody := []byte(string(requestJSON))
		req, err := http.NewRequest("GET", services.Cfg.Services.EnrichmentApiUrl+"/enrichPerson", bytes.NewBuffer(requestBody))
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to create request"))
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Internal error"})
			return
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
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Internal error"})
		}
		defer response.Body.Close()
		span.LogFields(log.Int("response.status", response.StatusCode))

		var result enrichmentmodel.EnrichPersonResponse
		err = json.NewDecoder(response.Body).Decode(&result)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to decode response"))
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Internal error"})
			return
		}

		enrichPersonResponse := EnrichPersonResponse{
			Status: "success",
		}

		mapData(result, &enrichPersonResponse)

		c.JSON(http.StatusOK, enrichPersonResponse)
	}
}

func mapData(source enrichmentmodel.EnrichPersonResponse, target *EnrichPersonResponse) {
	if source.Data == nil {
		return
	}
	if source.Data.PersonProfile == nil {
		return
	}
	target.IsComplete = true          // TODO when adding async check better contact results
	target.PendingFields = []string{} // TODO when adding async check better contact results
	target.ResultURL = nil            // TODO when adding async check better contact results
	target.Data = EnrichPersonData{}

	// set emails
	if source.Data.PersonProfile.Email != "" {
		target.Data.Emails = append(target.Data.Emails, EnrichPersonEmail{
			Address: source.Data.PersonProfile.Email,
			//IsDeliverable: true, // TODO implement me with verify API
			//IsRisky:       false, // TODO implement me
			//Type: source.Data.PersonProfile.EmailType, // TODO use verify response, if free -> personal, else -> professional
		})
	}

	// set name
	target.Data.Name = EnrichPersonName{
		FirstName: source.Data.PersonProfile.Person.FirstName,
		LastName:  source.Data.PersonProfile.Person.LastName,
	}
	if source.Data.PersonProfile.Person.FirstName != "" && source.Data.PersonProfile.Person.LastName != "" {
		target.Data.Name.FullName = source.Data.PersonProfile.Person.FirstName + " " + source.Data.PersonProfile.Person.LastName
	}

	// set jobs
	for _, position := range source.Data.PersonProfile.Person.Positions.PositionHistory {
		enrichPersonJob := EnrichPersonJob{
			Title:           position.Title,
			Company:         position.CompanyName,
			CompanyLinkedin: position.LinkedInUrl,
			//CompanyWebsite:  position.CompanyWebsite, // TODO implement with company enrichment
			IsCurrent: position.StartEndDate.End == nil,
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
		target.Data.Jobs = append(target.Data.Jobs, enrichPersonJob)
	}

	// set profile picture
	target.Data.ProfilePic = source.Data.PersonProfile.Person.PhotoUrl

	// set social
	target.Data.Social = EnrichPersonSocial{
		Linkedin: EnrichPersonLinkedIn{
			ID:            source.Data.PersonProfile.Person.LinkedInIdentifier,
			PublicID:      source.Data.PersonProfile.Person.PublicIdentifier,
			URL:           source.Data.PersonProfile.Person.LinkedInUrl,
			FollowerCount: source.Data.PersonProfile.Person.FollowerCount,
		},
		// TODO add X, github, other from Customer OS
	}

	// set location // TODO implement AI lookup to get details
	target.Data.Location = EnrichPersonLocation{
		//City: source.Data.PersonProfile.Person.Location,
		Region: source.Data.PersonProfile.Person.Location,
		//Country:  source.Data.PersonProfile.Person.Country,
		//Timezone: source.Data.PersonProfile.Person.Timezone,
	}

	// set phone numbers
	// TODO implement phone numbers
}
