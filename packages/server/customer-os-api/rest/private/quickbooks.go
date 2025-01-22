package private

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	cosapi_services "github.com/customeros/customeros/packages/server/customer-os-api/services"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/security"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go/log"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type OauthQuickbooksResponse struct {
	ExpiresIn              int     `json:"expires_in"`
	TokenType              string  `json:"token_type"`
	XRefreshTokenExpiresIn int     `json:"x_refresh_token_expires_in"`
	RefreshToken           string  `json:"refresh_token"`
	AccessToken            string  `json:"access_token"`
	Error                  *string `json:"error"`
}

func RequestAccessQuickbooks(s *cosapi_services.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		_, span := tracing.StartHttpServerTracerSpanWithHeader(c, "/internal/v1/settings/quickbooks/requestAccess", c.Request.Header)
		defer span.Finish()

		quickbooksRequestAccessUrl := "https://appcenter.intuit.com/connect/oauth2?client_id=" + s.Cfg.Common.External.QuickbooksConfig.ClientId + "&redirect_uri=" + s.Cfg.Common.External.QuickbooksConfig.RedirectUrl + "&response_type=code&scope=com.intuit.quickbooks.accounting&state=12345"

		span.LogFields(log.Object("quickbooksRequestAccessUrl", quickbooksRequestAccessUrl))

		c.JSON(http.StatusOK, gin.H{"url": quickbooksRequestAccessUrl})
	}
}

func CallbackQuickbooks(s *cosapi_services.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c, "/internal/v1/settings/slack/oauth/callback", c.Request.Header)
		defer span.Finish()

		tenant, _ := c.Get(security.KEY_TENANT_NAME)

		quickbooksSettingsEntity, err := s.Repositories.PostgresRepositories.QuickbooksSettingsRepository.Get(ctx, tenant.(string))
		if err != nil {
			tracing.TraceErr(span, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		code := c.Request.URL.Query().Get("code")
		realmId := c.Request.URL.Query().Get("realmId")

		requestData := url.Values{}
		requestData.Set("grant_type", "authorization_code")
		requestData.Set("code", code)
		requestData.Set("redirect_uri", s.Cfg.Common.External.QuickbooksConfig.RedirectUrl)

		// Encode the form data
		requestBody := requestData.Encode()

		request, err := http.NewRequest("POST", "https://oauth.platform.intuit.com/oauth2/v1/tokens/bearer", nil)
		if err != nil {
			fmt.Println("Error creating request:", err)
			return
		}
		request.Body = ioutil.NopCloser(strings.NewReader(requestBody))

		toString := base64.StdEncoding.EncodeToString([]byte(s.Cfg.Common.External.QuickbooksConfig.ClientId + ":" + s.Cfg.Common.External.QuickbooksConfig.ClientSecret))
		request.Header.Set("Authorization", "Basic "+toString)
		request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		// Perform the HTTP request
		client := &http.Client{}
		resp, err := client.Do(request)
		if err != nil {
			fmt.Println("Error making request:", err)
			return
		}
		defer resp.Body.Close()

		// Read and print the response
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading response:", err)
			return
		}

		// convert body to OauthSlackResponse
		var quickbooksResponse OauthQuickbooksResponse
		err = json.Unmarshal(body, &quickbooksResponse)
		if err != nil {
			fmt.Println("Error unmarshalling response:", err)
			return
		}

		if quickbooksResponse.Error == nil {
			now := utils.Now()

			entity := postgres_entity.QuickbooksSettingsEntity{
				Tenant:                tenant.(string),
				RealmId:               realmId,
				AccessToken:           quickbooksResponse.AccessToken,
				AccessTokenExpiresIn:  quickbooksResponse.ExpiresIn,
				AccessTokenExpiresAt:  now.Add(time.Duration(quickbooksResponse.ExpiresIn) * time.Second),
				RefreshToken:          quickbooksResponse.RefreshToken,
				RefreshTokenExpiresIn: quickbooksResponse.XRefreshTokenExpiresIn,
				RefreshTokenExpiresAt: now.Add(time.Duration(quickbooksResponse.XRefreshTokenExpiresIn) * time.Second),
				RefreshTokenExpired:   false,
			}

			if quickbooksSettingsEntity != nil {
				entity.Id = quickbooksSettingsEntity.Id
			}

			_, err := s.Repositories.PostgresRepositories.QuickbooksSettingsRepository.Save(c, entity)
			if err != nil {
				fmt.Println("Error saving slack settings:", err)
				return
			}
		}

		c.JSON(http.StatusOK, gin.H{})
	}
}
