package private

import (
	"encoding/json"
	"fmt"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/security"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go/log"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	cosapi_services "github.com/customeros/customeros/packages/server/customer-os-api/services"
)

type OauthSlackResponse struct {
	Ok         bool   `json:"ok"`
	Error      string `json:"error"`
	AppId      string `json:"app_id"`
	AuthedUser struct {
		Id string `json:"id"`
	} `json:"authed_user"`
	Scope       string `json:"scope"`
	TokenType   string `json:"token_type"`
	AccessToken string `json:"access_token"`
	BotUserId   string `json:"bot_user_id"`
	Team        struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	} `json:"team"`
	Enterprise          interface{} `json:"enterprise"`
	IsEnterpriseInstall bool        `json:"is_enterprise_install"`
}

type OauthSlackRevokeResponse struct {
	Ok      bool    `json:"ok"`
	Revoked *bool   `json:"revoked"`
	Error   *string `json:"error"`
}

func RequestAccessSlack(s *cosapi_services.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		_, span := tracing.StartHttpServerTracerSpanWithHeader(c, "/internal/v1/settings/slack/requestAccess", c.Request.Header)
		defer span.Finish()

		scopes := []string{
			"channels:history",
			"channels:join",
			"channels:manage",
			"channels:read",
			"chat:write",
			"files:read",
			"groups:history",
			"groups:read",
			"groups:write",
			"incoming-webhook",
			"links:read",
			"reactions:read",
			"team:read",
			"usergroups:read",
			"users.profile:read",
			"users:read",
			"users:read.email",
		}
		slackRequestAccessUrl := "https://slack.com/oauth/v2/authorize?client_id=" + s.Cfg.Common.External.SlackConfig.ClientID + "&scope=" + strings.Join(scopes, ",") + "&user_scope="

		span.LogFields(log.Object("slackRequestAccessUrl", slackRequestAccessUrl))

		c.JSON(http.StatusOK, gin.H{"url": slackRequestAccessUrl})
	}
}

func CallbackSlack(s *cosapi_services.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c, "/internal/v1/settings/slack/oauth/callback", c.Request.Header)
		defer span.Finish()

		tenant, _ := c.Get(security.KEY_TENANT_NAME)

		slackSettingsEntity, err := s.Repositories.PostgresRepositories.SlackSettingsRepository.Get(ctx, tenant.(string))
		if err != nil {
			tracing.TraceErr(span, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		code := c.Request.URL.Query().Get("code")

		requestData := url.Values{}
		requestData.Set("code", code)
		requestData.Set("client_id", s.Cfg.Common.External.SlackConfig.ClientID)
		requestData.Set("client_secret", s.Cfg.Common.External.SlackConfig.ClientSecret)

		// Encode the form data
		requestBody := requestData.Encode()

		request, err := http.NewRequest("POST", "https://slack.com/api/oauth.v2.access", nil)
		if err != nil {
			tracing.TraceErr(span, err)
			return
		}
		request.Body = ioutil.NopCloser(strings.NewReader(requestBody))

		request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		request.Header.Set("Content-Length", fmt.Sprint(len(requestBody)))

		// Perform the HTTP request
		client := &http.Client{}
		resp, err := client.Do(request)
		if err != nil {
			tracing.TraceErr(span, err)
			return
		}
		defer resp.Body.Close()

		// Read and print the response
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			tracing.TraceErr(span, err)
			return
		}

		// convert body to OauthSlackResponse
		var slackResponse OauthSlackResponse
		err = json.Unmarshal(body, &slackResponse)
		if err != nil {
			tracing.TraceErr(span, err)
			return
		}

		if slackResponse.Ok {
			entity := postgres_entity.SlackSettingsEntity{
				TenantName:   tenant.(string),
				AppId:        slackResponse.AppId,
				AuthedUserId: slackResponse.AuthedUser.Id,
				Scope:        slackResponse.Scope,
				TokenType:    slackResponse.TokenType,
				AccessToken:  slackResponse.AccessToken,
				BotUserId:    slackResponse.BotUserId,
				TeamId:       slackResponse.Team.Id,
			}

			if slackSettingsEntity != nil {
				entity.Id = slackSettingsEntity.Id
			}

			_, err := s.Repositories.PostgresRepositories.SlackSettingsRepository.Save(ctx, entity)
			if err != nil {
				tracing.TraceErr(span, err)
				return
			}
		}

		c.JSON(http.StatusOK, gin.H{})
	}
}

func RevokeSlack(s *cosapi_services.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c, "/internal/v1/settings/slack/revoke", c.Request.Header)
		defer span.Finish()

		tenant, _ := c.Get(security.KEY_TENANT_NAME)

		slackSettingsEntity, err := s.Repositories.PostgresRepositories.SlackSettingsRepository.Get(ctx, tenant.(string))
		if err != nil {
			tracing.TraceErr(span, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if slackSettingsEntity != nil {
			request, err := http.NewRequest("GET", "https://slack.com/api/auth.revoke", nil)
			if err != nil {
				tracing.TraceErr(span, err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			request.Header.Set("Authorization", "Bearer "+slackSettingsEntity.AccessToken)

			client := &http.Client{}
			resp, err := client.Do(request)
			if err != nil {
				tracing.TraceErr(span, err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				tracing.TraceErr(span, err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			// convert body to OauthSlackResponse
			var slackResponse OauthSlackRevokeResponse
			err = json.Unmarshal(body, &slackResponse)
			if err != nil {
				tracing.TraceErr(span, err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			if slackResponse.Ok && slackResponse.Revoked != nil && *slackResponse.Revoked {
				err := s.Repositories.PostgresRepositories.SlackSettingsRepository.Delete(c, tenant.(string))
				if err != nil {
					tracing.TraceErr(span, err)
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				c.JSON(http.StatusOK, gin.H{})
			} else {
				tracing.TraceErr(span, err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": slackResponse.Error})
			}
		}

		c.JSON(http.StatusOK, gin.H{})
	}
}
