package routes

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	cosModel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api-sdk/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-auth/repository/postgres/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	commonUtils "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/neo4jutil"
	neo4jrepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"
	postgresEntity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/user-admin-api/config"
	"github.com/openline-ai/openline-customer-os/packages/server/user-admin-api/model"
	"github.com/openline-ai/openline-customer-os/packages/server/user-admin-api/service"
	"github.com/openline-ai/openline-customer-os/packages/server/user-admin-api/utils"
	"github.com/opentracing/opentracing-go"
	tracingLog "github.com/opentracing/opentracing-go/log"
	tokenOauth "golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	googleOauth "google.golang.org/api/oauth2/v2"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

const APP_SOURCE = "user-admin-api"

func addRegistrationRoutes(rg *gin.RouterGroup, config *config.Config, services *service.Services) {
	personalEmailProviders, err := services.CommonServices.PostgresRepositories.PersonalEmailProviderRepository.GetPersonalEmailProviders()
	if err != nil {
		panic(err)
	}

	rg.POST("/signin",
		func(ginContext *gin.Context) {
			c, cancel := commonUtils.GetContextWithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c, "/signin", ginContext.Request.Header)
			defer span.Finish()

			apiKey := ginContext.GetHeader("X-Openline-Api-Key")
			if apiKey != config.Service.ApiKey {
				span.LogFields(tracingLog.String("error", "invalid api key"))
				ginContext.JSON(http.StatusUnauthorized, gin.H{
					"result": fmt.Sprintf("invalid api key"),
				})
				return
			}

			var signInRequest model.SignInRequest
			if err := ginContext.BindJSON(&signInRequest); err != nil {
				tracing.TraceErr(span, err)
				ginContext.JSON(http.StatusInternalServerError, gin.H{
					"result": fmt.Sprintf("unable to parse json: %v", err.Error()),
				})
				return
			}

			span.LogFields(tracingLog.Object("request", signInRequest))

			firstName, lastName, err := validateRequestAtProvider(config, signInRequest, ginContext)
			if err != nil {
				tracing.TraceErr(span, err)
				ginContext.JSON(http.StatusInternalServerError, gin.H{
					"result": fmt.Sprintf("unable to validate request at provider: %v", err.Error()),
				})
				return
			}

			tenantName, isNewTenant, err := getTenant(ctx, services.CustomerOsClient, personalEmailProviders, signInRequest, ginContext, config)
			if err != nil {
				tracing.TraceErr(span, err)
				ginContext.JSON(http.StatusInternalServerError, gin.H{
					"result": fmt.Sprintf("unable to get tenant: %v", err.Error()),
				})
				return
			}

			ctx = common.WithCustomContext(ctx, &common.CustomContext{
				Tenant: *tenantName,
			})

			_, err = initializeUser(services, signInRequest.Provider, signInRequest.OAuthToken.ProviderAccountId, *tenantName, signInRequest.Email, firstName, lastName, ginContext)
			if err != nil {
				tracing.TraceErr(span, err)
				ginContext.JSON(http.StatusInternalServerError, gin.H{
					"result": fmt.Sprintf("unable to initialize user: %v", err.Error()),
				})
				return
			}

			if isNewTenant {
				go func() {
					c, cancelFunc := context.WithTimeout(context.Background(), 300*time.Second)
					defer cancelFunc()

					ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c, "/signin - register new tenant", ginContext.Request.Header)
					defer span.Finish()

					err = registerNewTenantAsLeadInProviderTenant(ctx, config, services, signInRequest.Email)
					if err != nil {
						tracing.TraceErr(span, err)
						ginContext.JSON(http.StatusInternalServerError, gin.H{
							"result": fmt.Sprintf("unable to register new tenant as lead in provider tenant: %v", err.Error()),
						})
						return
					}

					span.LogFields(tracingLog.String("result", "ok"))
				}()
			}

			// Handle Google provider
			if signInRequest.Provider == "google" {
				if isRequestEnablingOAuthSync(signInRequest) {
					var oauthToken, _ = services.AuthServices.OAuthTokenService.GetByPlayerIdAndProvider(ctx, signInRequest.OAuthToken.ProviderAccountId, signInRequest.Provider)
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
					_, err := services.AuthServices.OAuthTokenService.Save(ctx, *oauthToken)
					if err != nil {
						log.Printf("unable to save oauth token: %v", err.Error())
						ginContext.JSON(http.StatusInternalServerError, gin.H{
							"result": fmt.Sprintf("unable to save oauth token: %v", err.Error()),
						})
						return
					}
				}
			} else if signInRequest.Provider == "azure-ad" {
				var oauthToken, _ = services.AuthServices.OAuthTokenService.GetByPlayerIdAndProvider(ctx, signInRequest.OAuthToken.ProviderAccountId, signInRequest.Provider)
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
				_, err := services.AuthServices.OAuthTokenService.Save(ctx, *oauthToken)
				if err != nil {
					log.Printf("unable to save oauth token: %v", err.Error())
					ginContext.JSON(http.StatusInternalServerError, gin.H{
						"result": fmt.Sprintf("unable to save oauth token: %v", err.Error()),
					})
					return
				}
			} else {
				log.Printf("Unsupported provider: %s", signInRequest.Provider)
				ginContext.JSON(http.StatusBadRequest, gin.H{
					"result": fmt.Sprintf("Unsupported provider: %s", signInRequest.Provider),
				})
				return
			}

			ginContext.JSON(http.StatusOK, gin.H{"status": "ok"})
		})

	rg.POST("/revoke",
		tracing.TracingEnhancer(context.Background(), "POST /revoke"),
		func(ginContext *gin.Context) {
			ctx, cancel := commonUtils.GetLongLivedContext(context.Background())
			defer cancel()

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

			var oauthToken, _ = services.AuthServices.OAuthTokenService.GetByPlayerIdAndProvider(ctx, revokeRequest.ProviderAccountId, revokeRequest.Provider)

			if oauthToken != nil && oauthToken.RefreshToken != "" {
				// Handle revocation based on provider
				var revocationURL string
				switch revokeRequest.Provider {
				case "google":
					revocationURL = fmt.Sprintf("https://accounts.google.com/o/oauth2/revoke?token=%s", oauthToken.RefreshToken)
				case "azure-ad":
					revocationURL = fmt.Sprintf("https://graph.microsoft.com/v1.0/me/revokeSignInSessions")
				}

				if revocationURL != "" {
					resp, err := http.Get(revocationURL)
					if err != nil {
						ginContext.JSON(http.StatusInternalServerError, gin.H{})
						return
					}

					if resp.StatusCode == http.StatusOK {
						// Successfully revoked, delete the token
						err := services.AuthServices.OAuthTokenService.DeleteByPlayerIdAndProvider(ctx, revokeRequest.ProviderAccountId, revokeRequest.Provider)
						if err != nil {
							ginContext.JSON(http.StatusInternalServerError, gin.H{})
							return
						}
					} else {
						// Revocation failed
						ginContext.JSON(http.StatusInternalServerError, gin.H{})
						return
					}
				}
			}

			ginContext.JSON(http.StatusOK, gin.H{})
		})
}

func getTenant(ctx context.Context, cosClient service.CustomerOsClient, personalEmailProvider []postgresEntity.PersonalEmailProvider, signInRequest model.SignInRequest, ginContext *gin.Context, config *config.Config) (*string, bool, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "getTenant")
	defer span.Finish()

	domain := commonUtils.ExtractDomain(signInRequest.Email)
	span.LogFields(tracingLog.String("domain", domain))

	var tenantName *string
	var err error

	isPersonalEmail := false
	//check if the user is using a personal email provider
	for _, personalEmailProvider := range personalEmailProvider {
		if strings.Contains(domain, personalEmailProvider.ProviderDomain) {
			isPersonalEmail = true
			break
		}
	}

	span.LogFields(tracingLog.Bool("isPersonalEmail", isPersonalEmail))

	if isPersonalEmail {
		player, err := cosClient.GetPlayer(signInRequest.Email, signInRequest.Provider)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, false, err
		}
		if player != nil && player.PlayerByAuthIdProvider.Users != nil && len(*player.PlayerByAuthIdProvider.Users) > 0 {
			span.LogFields(tracingLog.Object("playerIdentified", true))
			for _, user := range *player.PlayerByAuthIdProvider.Users {
				if user.Tenant != "" {
					tenantName = &user.Tenant
					break
				}
			}
		} else {
			span.LogFields(tracingLog.Object("playerIdentified", false))
		}
	} else {
		tenantName, err = cosClient.GetTenantByWorkspace(&model.WorkspaceInput{
			Name:     domain,
			Provider: signInRequest.Provider,
		})
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, false, err
		}
		if tenantName != nil {
			span.LogFields(tracingLog.String("tenantIdentifiedSameWorkspace", *tenantName))
			return tenantName, false, nil
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
			tracing.TraceErr(span, err)
			return nil, false, err
		}

		if tenantName != nil {
			span.LogFields(tracingLog.String("tenantIdentifiedDifferentWorkspace", *tenantName))

			err = createWorkspaceInTenant(ginContext, cosClient, *tenantName, signInRequest.Provider, domain, APP_SOURCE)
			if err != nil {
				tracing.TraceErr(span, err)
				return nil, false, err
			}

			return tenantName, false, nil
		}
	}

	if tenantName == nil {
		var tenantStr string
		if isPersonalEmail {
			tenantStr = utils.GenerateName()
		} else {
			tenantStr = utils.Sanitize(domain)
		}

		span.LogFields(tracingLog.String("newTenantCreationWith", tenantStr))

		tenantName, err = createTenant(cosClient, tenantStr, APP_SOURCE)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, false, err
		}

		if !isPersonalEmail {
			err = createWorkspaceInTenant(ginContext, cosClient, *tenantName, signInRequest.Provider, domain, APP_SOURCE)
			if err != nil {
				tracing.TraceErr(span, err)
				return nil, false, err
			}
		}

		if config.Slack.NotifyNewTenantRegisteredWebhook != "" {
			notifyNewTenantCreation(config.Slack.NotifyNewTenantRegisteredWebhook, tenantStr, signInRequest.Email)
		}

		return tenantName, true, nil
	} else {
		return tenantName, false, nil
	}
}

func notifyNewTenantCreation(slackWehbookUrl, tenant, email string) {
	// Create a struct to hold the JSON data
	type SlackMessage struct {
		Text string `json:"text"`
	}
	message := SlackMessage{Text: tenant + " tenant registered by " + email}

	// Convert struct to JSON
	jsonData, err := json.Marshal(message)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return
	}

	// Send POST request
	resp, err := http.Post(slackWehbookUrl, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	// Print response status
	fmt.Println("Response Status:", resp.Status)
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
	log.Printf("Initialize user: %s", email)
	appSource := APP_SOURCE

	playerExists := false
	userExists := false

	player, err := services.CustomerOsClient.GetPlayer(email, provider)
	if err != nil {
		return nil, err
	}
	if player != nil && player.PlayerByAuthIdProvider.Id != "" {
		playerExists = true
		log.Printf("Initialize user - existing player: %v", player.PlayerByAuthIdProvider)
	} else {
		log.Printf("Initialize user - player not found")
	}

	userByEmail, err := services.CustomerOsClient.GetUserByEmail(tenant, email)
	if err != nil {
		return nil, err
	}
	if userByEmail != nil && userByEmail.ID != "" {
		userExists = true
		log.Printf("Initialize user - user by email: %v", userByEmail)
	} else {
		log.Printf("Initialize user - user by email not found")
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
		log.Printf("Initialize user - user created: %v", userByEmail)

		for attempt := 1; attempt <= 5; attempt++ {
			checkUserByEmail, err := services.CustomerOsClient.GetUserByEmail(tenant, email)
			if err == nil && checkUserByEmail.ID != "" {
				break
			}
			time.Sleep(commonUtils.BackOffExponentialDelay(attempt))
		}
	} else {
		if !playerExists {
			err = services.CustomerOsClient.CreatePlayer(tenant, userByEmail.ID, providerAccountId, email, provider)
			if err != nil {
				return nil, err
			}
		}
	}

	err = addDefaultMissingRoles(services, userByEmail, &tenant)
	if err != nil {
		return nil, err
	}

	return userByEmail, nil
}

func addDefaultMissingRoles(services *service.Services, user *model.UserResponse, tenant *string) error {
	var rolesToAdd []model.Role

	log.Printf("Add default missing roles for user: %v", user)
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

	log.Printf("Roles to add: %v to %s", rolesToAdd, user.ID)
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

func registerNewTenantAsLeadInProviderTenant(ctx context.Context, config *config.Config, services *service.Services, registeredEmail string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "registration.registerNewTenantAsLeadInProviderTenant")
	defer span.Finish()

	organizationId, contactId, err := services.RegistrationService.CreateOrganizationAndContact(ctx, config.Service.ProviderTenantName, registeredEmail, true, "tenant-registration")
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	span.LogFields(tracingLog.String("providerOrganizationId", *organizationId))
	span.LogFields(tracingLog.String("providerContactId", *contactId))

	waitForOrganizationAndContactToBeCreated(ctx, services, organizationId, contactId)

	emailId, err := services.CustomerOSApiClient.MergeEmailToContact(config.Service.ProviderTenantName, *contactId, cosModel.EmailInput{Email: registeredEmail})
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	span.LogFields(tracingLog.String("emailId", emailId))

	url := config.Comms.CommsAPI
	method := "POST"

	type EmailPayload struct {
		Channel   string   `json:"channel"`
		Username  string   `json:"username"`
		Direction string   `json:"direction"`
		To        []string `json:"to"`
		Cc        []string `json:"cc"`
		Bcc       []string `json:"bcc"`
		Content   string   `json:"content"`
		Subject   string   `json:"subject"`
	}

	payload := EmailPayload{
		Channel:   "EMAIL",
		Username:  config.Service.ProviderUsername,
		Direction: "OUTBOUND",
		To:        []string{registeredEmail},
		Cc:        []string{},
		Bcc:       []string{},
		Subject:   "Welcome to CustomerOS",
		Content: `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Welcome to CustomerOS</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            line-height: 1.6;
        }
        .header {
            font-size: 1.4em;
            margin-bottom: 10px;
        }
        .content {
            margin-bottom: 20px;
        }
        .footer {
            font-size: 0.9em;
            color: #888;
        }
        .signature {
            margin-top: 20px;
        }
    </style>
</head>
<body>
<div>
    <div class="header">Hey, welcome to CustomerOS!</div>
    <div class="content">
        <p>Thanks for trying us out.</p>
        <p>To be honest, our self-service onboarding kinda sucks right now as we’re still building it out. I’d love to get you setup and using the tool, would you be open to spending 10 mins with me to help you get things configured?</p>
        <p>Please grab any slot on my <a href="https://app.customeros.ai/organization/cal.com/mbrown/20min" target="_blank">calendar.</a></p>
    </div>
    <div class="signature">
        <p>Thanks again,</p>
        <p>Matt Brown<br>
            CEO @ <a href="https://customeros.ai/?utm_source=signature&utm_medium=email&utm_campaign=signup" target="_blank">CustomerOS</a></p>
        <p class="footer">
            Follow me on <a href="https://www.linkedin.com/in/mateocafe/" target="_blank">LinkedIn</a><br>
            US: +1 650 977 2199<br>
            UK: +44 7700 155 600
        </p>
    </div>
</div>
</body>
</html>`,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	body := bytes.NewBuffer(payloadBytes)

	// Create new request
	req, err := http.NewRequest(method, url+"/mail/send", body)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	// Set headers
	req.Header.Add("X-Openline-Mail-Api-Key", config.Comms.CommsAPIKey)
	req.Header.Add("X-Openline-USERNAME", config.Service.ProviderUsername)
	req.Header.Add("Content-Type", "application/json")

	// Create HTTP client and make request
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	defer res.Body.Close()

	mapBody := make(map[string]interface{})

	// Read response
	responseBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	// Convert response to map
	err = json.Unmarshal(responseBody, &mapBody)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if mapBody["error"] != nil {
		tracing.TraceErr(span, fmt.Errorf("error: %v", mapBody["error"]))
		return fmt.Errorf("error: %v", mapBody["error"])
	}

	span.LogFields(tracingLog.Object("email sent: ", mapBody))

	return nil
}

func waitForOrganizationAndContactToBeCreated(ctx context.Context, services *service.Services, organizationId, contactId *string) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "registration.waitForOrganizationAndContactToBeCreated")
	defer span.Finish()

	neo4jrepository.WaitForNodeCreatedInNeo4jWithConfig(ctx, span, services.CommonServices.Neo4jRepositories, *organizationId, neo4jutil.NodeLabelOrganization, 30*time.Second)
	neo4jrepository.WaitForNodeCreatedInNeo4jWithConfig(ctx, span, services.CommonServices.Neo4jRepositories, *contactId, neo4jutil.NodeLabelContact, 30*time.Second)
}
