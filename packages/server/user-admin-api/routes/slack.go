package routes

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-auth/repository/postgres/entity"
	commonService "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	"github.com/openline-ai/openline-customer-os/packages/server/user-admin-api/config"
	"github.com/openline-ai/openline-customer-os/packages/server/user-admin-api/service"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

func addSlackRoutes(rg *gin.RouterGroup, config *config.Config, services *service.Services) {
	rg.GET("/slack/requestAccess", commonService.TenantUserContextEnhancer(commonService.USERNAME_OR_TENANT, services.CommonServices.CommonRepositories), func(ctx *gin.Context) {
		slackRequestAccessUrl := "https://slack.com/oauth/v2/authorize?client_id=" + config.Slack.ClientId + "&scope=channels:history,channels:join,channels:read,files:read,groups:history,groups:read,im:history,links:read,reactions:read,team:read,usergroups:read,users.profile:read,users:read,users:read.email&user_scope="

		ctx.JSON(http.StatusOK, gin.H{"url": slackRequestAccessUrl})
	})
	rg.POST("/slack/oauth/callback", commonService.TenantUserContextEnhancer(commonService.USERNAME_OR_TENANT, services.CommonServices.CommonRepositories), func(ctx *gin.Context) {
		tenant, _ := ctx.Get(commonService.KEY_TENANT_NAME)

		slackSettingsEntity, err := services.AuthServices.CommonAuthRepositories.SlackSettingsRepository.Get(tenant.(string))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if slackSettingsEntity != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Slack settings already exists"})
			return
		}

		code := ctx.Request.URL.Query().Get("code")

		requestData := url.Values{}
		requestData.Set("code", code)
		requestData.Set("client_id", config.Slack.ClientId)
		requestData.Set("client_secret", config.Slack.ClientSecret)

		// Encode the form data
		requestBody := requestData.Encode()

		request, err := http.NewRequest("POST", "https://slack.com/api/oauth.v2.access", nil)
		if err != nil {
			fmt.Println("Error creating request:", err)
			return
		}
		request.Body = ioutil.NopCloser(strings.NewReader(requestBody))

		request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		request.Header.Set("Content-Length", fmt.Sprint(len(requestBody)))

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

		//convert body to OauthSlackResponse
		var slackResponse OauthSlackResponse
		err = json.Unmarshal(body, &slackResponse)
		if err != nil {
			fmt.Println("Error unmarshalling response:", err)
			return
		}

		if slackResponse.Ok {
			_, err := services.AuthServices.CommonAuthRepositories.SlackSettingsRepository.Save(entity.SlackSettingsEntity{
				TenantName:   tenant.(string),
				AppId:        slackResponse.AppId,
				AuthedUserId: slackResponse.AuthedUser.Id,
				Scope:        slackResponse.Scope,
				TokenType:    slackResponse.TokenType,
				AccessToken:  slackResponse.AccessToken,
				BotUserId:    slackResponse.BotUserId,
				TeamId:       slackResponse.Team.Id,
			})
			if err != nil {
				fmt.Println("Error saving slack settings:", err)
				return
			}
		}

		ctx.JSON(http.StatusOK, gin.H{})
	})
	rg.POST("/slack/revoke", commonService.TenantUserContextEnhancer(commonService.USERNAME_OR_TENANT, services.CommonServices.CommonRepositories), func(ctx *gin.Context) {
		tenant, _ := ctx.Get(commonService.KEY_TENANT_NAME)

		slackSettingsEntity, err := services.AuthServices.CommonAuthRepositories.SlackSettingsRepository.Get(tenant.(string))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if slackSettingsEntity != nil {

			request, err := http.NewRequest("GET", "https://slack.com/api/auth.revoke", nil)
			if err != nil {
				fmt.Println("Error creating request:", err)
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			request.Header.Set("Authorization", "Bearer "+slackSettingsEntity.AccessToken)

			client := &http.Client{}
			resp, err := client.Do(request)
			if err != nil {
				fmt.Println("Error making request:", err)
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Println("Error reading response:", err)
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			//convert body to OauthSlackResponse
			var slackResponse OauthSlackRevokeResponse
			err = json.Unmarshal(body, &slackResponse)
			if err != nil {
				fmt.Println("Error unmarshalling response:", err)
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			if slackResponse.Ok && slackResponse.Revoked != nil && *slackResponse.Revoked {
				err := services.AuthServices.CommonAuthRepositories.SlackSettingsRepository.Delete(tenant.(string))
				if err != nil {
					ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				ctx.JSON(http.StatusOK, gin.H{})
			} else {
				logrus.Error("Error revoking slack token: ", slackResponse.Error)
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": slackResponse.Error})
			}
		}

		ctx.JSON(http.StatusOK, gin.H{})
	})
}

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
