package routes

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-auth/repository/postgres/entity"
	commonUtils "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/user-admin-api/config"
	"github.com/openline-ai/openline-customer-os/packages/server/user-admin-api/model"
	"github.com/openline-ai/openline-customer-os/packages/server/user-admin-api/service"
	"github.com/openline-ai/openline-customer-os/packages/server/user-admin-api/utils"
	tokenOauth "golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	googleOauth "google.golang.org/api/oauth2/v2"
	"log"
	"net/http"
	"strings"
)

const APP_SOURCE = "user-admin-api"

func addRegistrationRoutes(rg *gin.RouterGroup, config *config.Config, services *service.Services) {
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

		firstName, lastName, err := validateRequestAtProvider(config, signInRequest, ginContext)
		if err != nil {
			log.Printf("unable to validate request at provider: %v", err.Error())
			ginContext.JSON(http.StatusInternalServerError, gin.H{
				"result": fmt.Sprintf("unable to validate request at provider: %v", err.Error()),
			})
			return
		}

		tenantName, err := getTenant(services.CustomerOsClient, signInRequest, ginContext)
		if err != nil {
			log.Printf("unable to get tenant: %v", err.Error())
			ginContext.JSON(http.StatusInternalServerError, gin.H{
				"result": fmt.Sprintf("unable to get tenant: %v", err.Error()),
			})
			return
		}

		_, err = initializeUser(services, signInRequest.Provider, signInRequest.OAuthToken.ProviderAccountId, *tenantName, signInRequest.Email, firstName, lastName, ginContext)
		if err != nil {
			log.Printf("unable to initialize user: %v", err.Error())
			ginContext.JSON(http.StatusInternalServerError, gin.H{
				"result": fmt.Sprintf("unable to initialize user: %v", err.Error()),
			})
			return
		}

		if signInRequest.Provider == "google" {
			if isRequestEnablingOAuthSync(signInRequest) {
				//TODO Move this logic to a service
				var oauthToken, _ = services.AuthServices.OAuthTokenService.GetByPlayerIdAndProvider(signInRequest.OAuthToken.ProviderAccountId, signInRequest.Provider)
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
				oauthToken.NeedsManualRefresh = false
				if isRequestEnablingGmailSync(signInRequest) {
					oauthToken.GmailSyncEnabled = true
				}
				if isRequestEnablingGoogleCalendarSync(signInRequest) {
					oauthToken.GoogleCalendarSyncEnabled = true
				}
				services.AuthServices.OAuthTokenService.Save(*oauthToken)
			}
		}

		ginContext.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	rg.POST("/google/revoke", func(ginContext *gin.Context) {
		log.Printf("revoke oauth token")

		apiKey := ginContext.GetHeader("X-Openline-Api-Key")
		if apiKey != config.Service.ApiKey {
			ginContext.JSON(http.StatusUnauthorized, gin.H{
				"result": fmt.Sprintf("invalid api key"),
			})
			return
		}

		var revokeRequest model.RevokeRequest
		if err := ginContext.BindJSON(&revokeRequest); err != nil {
			log.Printf("unable to parse json: %v", err.Error())
			ginContext.JSON(http.StatusInternalServerError, gin.H{
				"result": fmt.Sprintf("unable to parse json: %v", err.Error()),
			})
			return
		}
		log.Printf("parsed json: %v", revokeRequest)

		var oauthToken, _ = services.AuthServices.OAuthTokenService.GetByPlayerIdAndProvider(revokeRequest.ProviderAccountId, "google")

		var resp *http.Response
		var err error

		if oauthToken.RefreshToken != "" {
			url := fmt.Sprintf("https://accounts.google.com/o/oauth2/revoke?token=%s", oauthToken.RefreshToken)
			resp, err = http.Get(url)
			if err != nil {
				ginContext.JSON(http.StatusInternalServerError, gin.H{})
				return
			}
		}

		if resp == nil || resp.StatusCode == 200 {
			err := services.AuthServices.OAuthTokenService.DeleteByPlayerIdAndProvider(revokeRequest.ProviderAccountId, "google")
			if err != nil {
				ginContext.JSON(http.StatusInternalServerError, gin.H{})
				return
			}
		} else {
			if resp != nil && resp.StatusCode != 200 {
				ginContext.JSON(http.StatusInternalServerError, gin.H{})
				return
			}
		}

		ginContext.JSON(http.StatusOK, gin.H{})
	})
}

func getTenant(cosClient service.CustomerOsClient, signInRequest model.SignInRequest, ginContext *gin.Context) (*string, error) {
	domain := commonUtils.ExtractDomain(signInRequest.Email)
	tenantName, err := cosClient.GetTenantByWorkspace(&model.WorkspaceInput{
		Name:     domain,
		Provider: signInRequest.Provider,
	})
	if err != nil {
		return nil, err
	}
	if tenantName != nil {
		return tenantName, nil
	}

	//tenant not found by the requested login info, try to find it by another workspace with the same domain
	var provider string
	if signInRequest.Provider == "google" {
		provider = "azure-ad"
	} else if signInRequest.Provider == "azure-ad" {
		provider = "google"
	}
	tenantName, err = cosClient.GetTenantByWorkspace(&model.WorkspaceInput{
		Name:     domain,
		Provider: provider,
	})
	if err != nil {
		return nil, err
	}

	//if not other workspace with the same domain, create tenant
	if tenantName == nil {
		tenantStr := utils.GenerateName()
		tenantName, err = createTenant(cosClient, tenantStr, APP_SOURCE)
		if err != nil {
			return nil, err
		}
	}

	err = createWorkspaceInTenant(ginContext, cosClient, *tenantName, signInRequest.Provider, domain, APP_SOURCE)
	if err != nil {
		return nil, err
	}

	return tenantName, nil
}

func validateRequestAtProvider(config *config.Config, signInRequest model.SignInRequest, ginContext *gin.Context) (*string, *string, error) {
	if signInRequest.Provider == "google" {
		userInfo, err := getUserInfoFromGoogle(ginContext, config, signInRequest)
		if err != nil {
			return nil, nil, err
		}

		return &userInfo.GivenName, &userInfo.FamilyName, nil
	} else if signInRequest.Provider == "azure-ad" {
		client := &http.Client{}
		// Create a GET request with the Authorization header.
		req, err := http.NewRequest("GET", "https://graph.microsoft.com/oidc/userinfo", nil)
		if err != nil {
			return nil, nil, err
		}

		req.Header.Set("Authorization", "Bearer "+signInRequest.OAuthToken.AccessToken)

		resp, err := client.Do(req)
		if err != nil {
			return nil, nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			var data map[string]string
			json.NewDecoder(resp.Body).Decode(&data)

			firstName := data["given_name"]
			lastName := data["family_name"]
			return &firstName, &lastName, nil
		} else {
			return nil, nil, err
		}
	} else {
		return nil, nil, fmt.Errorf("provider not supported")
	}
}

func createTenant(cosClient service.CustomerOsClient, tenant string, appSource string) (*string, error) {
	tenant, err := cosClient.MergeTenant(&model.TenantInput{
		Name:      tenant,
		AppSource: &appSource})
	if err != nil {
		return nil, err
	}
	return &tenant, nil
}

func createWorkspaceInTenant(c *gin.Context, cosClient service.CustomerOsClient, tenant, provider, domain string, appSource string) error {
	mergeWorkspaceRes, err := cosClient.MergeTenantToWorkspace(&model.WorkspaceInput{
		Name:      domain,
		Provider:  provider,
		AppSource: &appSource,
	}, tenant)
	if err != nil {
		return err
	}
	if !mergeWorkspaceRes {
		return fmt.Errorf("unable to merge workspace")
	}
	return nil
}

func getUserInfoFromGoogle(ginContext *gin.Context, config *config.Config, signInRequest model.SignInRequest) (*googleOauth.Userinfo, error) {
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

	oauth2Service, err := googleOauth.New(client)

	if err != nil {
		log.Printf("unable to create oauth2 service: %v", err.Error())
		ginContext.JSON(http.StatusInternalServerError, gin.H{
			"result": fmt.Sprintf("unable to create oauth2 service: %v", err.Error()),
		})
		return nil, err
	}
	userInfoService := googleOauth.NewUserinfoV2MeService(oauth2Service)

	userInfo, err := userInfoService.Get().Do()
	if err != nil {
		log.Printf("unable to get user info: %v", err.Error())
		ginContext.JSON(http.StatusInternalServerError, gin.H{
			"result": fmt.Sprintf("unable to get user info: %v", err.Error()),
		})
		return nil, err
	}

	return userInfo, nil
}

func initializeUser(services *service.Services, provider, providerAccountId, tenant, email string, firstName, lastName *string, ginContext *gin.Context) (*model.UserResponse, error) {
	appSource := APP_SOURCE

	playerExists := false
	userExists := false

	player, err := services.CustomerOsClient.GetPlayer(email, provider)
	if err != nil {
		return nil, err
	}
	if player != nil && player.Id != "" {
		playerExists = true
	}

	userByEmail, err := services.CustomerOsClient.GetUserByEmail(tenant, email)
	if err != nil {
		return nil, err
	}
	if userByEmail != nil && userByEmail.ID != "" {
		userExists = true
	}

	if !playerExists && !userExists {
		userByEmail, err = services.CustomerOsClient.CreateUser(&model.UserInput{
			FirstName: *firstName,
			LastName:  *lastName,
			Email: model.EmailInput{
				Email:     email,
				Primary:   true,
				AppSource: &appSource,
			},
			Player: model.PlayerInput{
				IdentityId: providerAccountId,
				AuthId:     email,
				Provider:   provider,
				AppSource:  &appSource,
			},
			AppSource: &appSource,
		}, tenant, []model.Role{model.RoleUser, model.RoleOwner})
		if err != nil {
			return nil, err
		}
	} else {
		if !playerExists {
			err = services.CustomerOsClient.CreatePlayer(tenant, userByEmail.ID, providerAccountId, email, provider)
			if err != nil {
				return nil, err
			}
		}
	}

	err = addDefaultMissingRoles(services, userByEmail, &tenant, ginContext)
	if err != nil {
		return nil, err
	}

	return userByEmail, nil
}

func addDefaultMissingRoles(services *service.Services, user *model.UserResponse, tenant *string, ginContext *gin.Context) error {
	var rolesToAdd []model.Role

	if user.Roles == nil || len(*user.Roles) == 0 {
		rolesToAdd = []model.Role{model.RoleUser, model.RoleOwner}
	} else {
		userRoleFound := false
		ownerRoleFound := false
		for _, role := range *user.Roles {
			if role == model.RoleUser {
				userRoleFound = true
			}
			if role == model.RoleOwner {
				ownerRoleFound = true
			}
		}
		if !userRoleFound {
			rolesToAdd = append(rolesToAdd, model.RoleUser)
		}
		if !ownerRoleFound {
			rolesToAdd = append(rolesToAdd, model.RoleOwner)
		}
	}

	if len(rolesToAdd) > 0 {
		_, err := services.CustomerOsClient.AddUserRoles(*tenant, user.ID, rolesToAdd)
		if err != nil {
			return err
		}
	}

	return nil
}

func isRequestEnablingGmailSync(signInRequest model.SignInRequest) bool {
	if strings.Contains(signInRequest.OAuthToken.Scope, "gmail") {
		return true
	}
	return false
}

func isRequestEnablingGoogleCalendarSync(signInRequest model.SignInRequest) bool {
	if strings.Contains(signInRequest.OAuthToken.Scope, "calendar") {
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
