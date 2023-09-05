package routes

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-auth/repository/postgres/entity"
	authCommonService "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-auth/service"
	"github.com/openline-ai/openline-customer-os/packages/server/user-admin-api/config"
	"github.com/openline-ai/openline-customer-os/packages/server/user-admin-api/model"
	"github.com/openline-ai/openline-customer-os/packages/server/user-admin-api/service"
	"github.com/openline-ai/openline-customer-os/packages/server/user-admin-api/utils"
	tokenOauth "golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/oauth2/v2"
	"log"
	"net/http"
	"strings"
)

const APP_SOURCE = "user-admin-api"

func addRegistrationRoutes(rg *gin.RouterGroup, config *config.Config, cosClient service.CustomerOsClient, authServices *authCommonService.Services) {
	rg.POST("/signin", func(ginContext *gin.Context) {
		log.Printf("Sign in User")
		apiKey := ginContext.GetHeader("X-Openline-Api-Key")
		if apiKey != config.Service.ApiKey {
			ginContext.JSON(http.StatusUnauthorized, gin.H{
				"result": fmt.Sprintf("invalid api key"),
			})
			return
		}
		log.Printf("api key is valid")
		var signInRequest model.SignInRequest
		if err := ginContext.BindJSON(&signInRequest); err != nil {
			log.Printf("unable to parse json: %v", err.Error())
			ginContext.JSON(http.StatusInternalServerError, gin.H{
				"result": fmt.Sprintf("unable to parse json: %v", err.Error()),
			})
			return
		}
		log.Printf("parsed json: %v", signInRequest)

		//get wokrspace for user

		var userInfo *oauth2.Userinfo

		if signInRequest.Provider == "google" {

			conf := &tokenOauth.Config{
				ClientID:     config.GoogleOAuth.ClientId,
				ClientSecret: config.GoogleOAuth.ClientSecret,
				Endpoint:     google.Endpoint,
			}

			token := tokenOauth.Token{
				AccessToken:  signInRequest.OAuthToken.AccessToken,
				RefreshToken: signInRequest.OAuthToken.RefreshToken,
				Expiry:       signInRequest.OAuthToken.ExpiresAt,
				TokenType:    "Bearer",
			}

			client := conf.Client(ginContext, &token)

			oauth2Service, err := oauth2.New(client)

			if err != nil {
				log.Printf("unable to create oauth2 service: %v", err.Error())
				ginContext.JSON(http.StatusInternalServerError, gin.H{
					"result": fmt.Sprintf("unable to create oauth2 service: %v", err.Error()),
				})
				return
			}
			userInfoService := oauth2.NewUserinfoV2MeService(oauth2Service)

			userInfo, err = userInfoService.Get().Do()

			if err != nil {
				log.Printf("unable to get user info: %v", err.Error())
				ginContext.JSON(http.StatusInternalServerError, gin.H{
					"result": fmt.Sprintf("unable to get user info: %v", err.Error()),
				})
				return
			}

		}

		var tenantName *string
		if userInfo.Hd != "" {
			tenant, err := cosClient.GetTenantByWorkspace(&model.WorkspaceInput{
				Name:     userInfo.Hd,
				Provider: signInRequest.Provider,
			})
			if err != nil {
				log.Printf("unable to get workspace: %v", err.Error())
				ginContext.JSON(http.StatusInternalServerError, gin.H{
					"result": fmt.Sprintf("unable to get workspace: %v", err.Error()),
				})
				return
			}
			if tenant != nil {
				log.Printf("tenant found %s", *tenant)
				var appSource = APP_SOURCE
				playerId, errorIsPlayer := cosClient.IsPlayer(signInRequest.Email, signInRequest.Provider)
				if errorIsPlayer != nil {
					playerId, err = cosClient.CreateUser(&model.UserInput{
						FirstName: userInfo.GivenName,
						LastName:  userInfo.FamilyName,
						Email: model.EmailInput{
							Email:     signInRequest.Email,
							Primary:   true,
							AppSource: &appSource,
						},
						Player: model.PlayerInput{
							IdentityId: signInRequest.OAuthToken.ProviderAccountId,
							AuthId:     signInRequest.Email,
							Provider:   signInRequest.Provider,
							AppSource:  &appSource,
						},
						AppSource: &appSource,
					}, *tenant, []service.Role{service.ROLE_USER})
					if err != nil {
						log.Printf("unable to create user: %v", err.Error())
						ginContext.JSON(http.StatusInternalServerError, gin.H{
							"result": fmt.Sprintf("unable to create user: %v", err.Error()),
						})
						return
					}
				}
				log.Printf("user created: %s", playerId)
				tenantName = tenant
			} else {
				var appSource = APP_SOURCE
				tenantStr := utils.Sanitize(userInfo.Hd)
				log.Printf("tenant not found for workspace, creating new tenant %s", tenantStr)
				// Workspace is not mapped to a tenant create a new tenant and map it to the workspace
				id, failed := makeTenantAndUser(ginContext, cosClient, tenantStr, appSource, signInRequest, userInfo)
				if failed {
					return
				}
				log.Printf("user created: %s", id)
				tenantName = &tenantStr
			}
		} else {
			// no workspace for this e-mail
			// check tenant exists for this e-mail
			if userInfo.Email != "" {
				var err error
				tenantName, err = cosClient.GetTenantByUserEmail(userInfo.Email)
				if err != nil {
					log.Printf("unable to get tenant: %v", err.Error())
					ginContext.JSON(http.StatusInternalServerError, gin.H{
						"result": fmt.Sprintf("unable to get tenant: %v", err.Error()),
					})
					return
				}
			}
			// no tenant for this e-mail, invent a tenant name
			if tenantName == nil {
				var appSource = APP_SOURCE
				tenantStr := utils.GenerateName()
				log.Printf("user has no workspace, inventing tenant %s", tenantStr)

				id, failed := makeTenantAndUser(ginContext, cosClient, tenantStr, appSource, signInRequest, userInfo)
				if failed {
					return
				}
				log.Printf("user created: %s", id)
				tenantName = &tenantStr
			}
		}

		if isRequestEnablingOAuthSync(signInRequest) {
			//TODO Move this logic to a service
			var oauthToken, _ = authServices.OAuthTokenService.GetByPlayerIdAndProvider(signInRequest.OAuthToken.ProviderAccountId, signInRequest.Provider)
			if oauthToken == nil {
				oauthToken = &entity.OAuthTokenEntity{}
			}
			oauthToken.Provider = signInRequest.Provider
			oauthToken.TenantName = *tenantName
			oauthToken.PlayerIdentityId = signInRequest.OAuthToken.ProviderAccountId
			oauthToken.EmailAddress = signInRequest.Email
			oauthToken.AccessToken = signInRequest.OAuthToken.AccessToken
			oauthToken.RefreshToken = signInRequest.OAuthToken.RefreshToken
			oauthToken.IdToken = signInRequest.OAuthToken.IdToken
			oauthToken.ExpiresAt = signInRequest.OAuthToken.ExpiresAt
			oauthToken.Scope = signInRequest.OAuthToken.Scope
			if isRequestEnablingGmailSync(signInRequest) {
				oauthToken.GmailSyncEnabled = true
			}
			if isRequestEnablingGoogleCalendarSync(signInRequest) {
				oauthToken.GoogleCalendarSyncEnabled = true
			}
			authServices.OAuthTokenService.Save(*oauthToken)
		}
		ginContext.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
}

func isRequestEnablingGmailSync(signInRequest model.SignInRequest) bool {
	if strings.Contains(signInRequest.OAuthToken.Scope, "gmail") {
		return true
	}
	return false
}

func isRequestEnablingGoogleCalendarSync(signInRequest model.SignInRequest) bool {
	if strings.Contains(signInRequest.OAuthToken.Scope, "calendar.events") {
		return true
	}
	return false
}

func isRequestEnablingOAuthSync(signInRequest model.SignInRequest) bool {
	if isRequestEnablingGmailSync(signInRequest) || isRequestEnablingGoogleCalendarSync(signInRequest) {
		return true
	}
	return false
}

func makeTenantAndUser(c *gin.Context, cosClient service.CustomerOsClient, tenantStr string, appSource string, req model.SignInRequest, userInfo *oauth2.Userinfo) (string, bool) {
	newTenantStr, err := cosClient.MergeTenant(&model.TenantInput{
		Name:      tenantStr,
		AppSource: &appSource})
	if err != nil {
		log.Printf("unable to create tenant: %v", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"result": fmt.Sprintf("unable to create tenant: %v", err.Error()),
		})
		return "", true
	}

	if userInfo.Hd != "" {
		mergeWorkspaceRes, err := cosClient.MergeTenantToWorkspace(&model.WorkspaceInput{
			Name:      userInfo.Hd,
			Provider:  req.Provider,
			AppSource: &appSource,
		}, newTenantStr)

		if err != nil {
			log.Printf("unable to merge workspace: %v", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"result": fmt.Sprintf("unable to merge workspace: %v", err.Error()),
			})
			return "", true
		}
		if !mergeWorkspaceRes {
			log.Printf("unable to merge workspace")
			c.JSON(http.StatusInternalServerError, gin.H{
				"result": fmt.Sprintf("unable to merge workspace"),
			})
			return "", true
		}
	}

	id, err := cosClient.CreateUser(&model.UserInput{
		FirstName: userInfo.GivenName,
		LastName:  userInfo.FamilyName,
		Email: model.EmailInput{
			Email:     req.Email,
			Primary:   true,
			AppSource: &appSource,
		},
		Player: model.PlayerInput{
			IdentityId: req.OAuthToken.ProviderAccountId,
			AuthId:     req.Email,
			Provider:   req.Provider,
			AppSource:  &appSource,
		},
		AppSource: &appSource,
	}, newTenantStr, []service.Role{service.ROLE_USER, service.ROLE_OWNER})
	if err != nil {
		log.Printf("unable to create user: %v", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"result": fmt.Sprintf("unable to create user: %v", err.Error()),
		})
		return "", true
	}
	return id, false
}
