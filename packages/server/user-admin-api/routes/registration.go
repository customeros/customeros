package routes

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	cosModel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api-sdk/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	model2 "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service/security"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	commonUtils "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
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
		security.ApiKeyCheckerHTTP(services.CommonServices.PostgresRepositories.TenantWebhookApiKeyRepository, services.CommonServices.PostgresRepositories.AppKeyRepository, security.USER_ADMIN_API),
		func(ginContext *gin.Context) {
			c, cancel := commonUtils.GetContextWithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c, "/signin", ginContext.Request.Header)
			defer span.Finish()

			var signInRequest model.SignInRequest
			if err := ginContext.BindJSON(&signInRequest); err != nil {
				tracing.TraceErr(span, err)
				ginContext.JSON(http.StatusInternalServerError, gin.H{
					"result": fmt.Sprintf("unable to parse json: %v", err.Error()),
				})
				return
			}

			span.LogFields(tracingLog.Object("request", signInRequest))

			firstName, lastName, err := validateRequestAtProvider(ctx, config, signInRequest)
			if err != nil {
				tracing.TraceErr(span, err)
				ginContext.JSON(http.StatusInternalServerError, gin.H{
					"result": fmt.Sprintf("unable to validate request at provider: %v", err.Error()),
				})
				return
			}

			if firstName == nil {
				s := ""
				firstName = &s
			}

			if lastName == nil {
				s := ""
				lastName = &s
			}

			var tenantName *string

			if signInRequest.Tenant == "" {
				span.LogFields(tracingLog.String("flow", "authentication"))

				tn, isNewTenant, err := getTenant(ctx, services, personalEmailProviders, signInRequest, ginContext, config)
				if err != nil {
					tracing.TraceErr(span, err)
					ginContext.JSON(http.StatusInternalServerError, gin.H{
						"result": fmt.Sprintf("unable to get tenant: %v", err.Error()),
					})
					return
				}
				tenantName = tn

				ctx = common.WithCustomContext(ctx, &common.CustomContext{
					Tenant: *tenantName,
				})

				err = initializeUser(ctx, services, signInRequest.Provider, signInRequest.OAuthToken.ProviderAccountId, *tenantName, signInRequest.LoggedInEmail, firstName, lastName, ginContext)
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

						err = registerNewTenantAsLeadInProviderTenant(ctx, config, services, signInRequest.LoggedInEmail)
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
			} else {
				span.LogFields(tracingLog.String("flow", "authorization"))

				emailId, err := services.CommonServices.Neo4jRepositories.EmailReadRepository.GetEmailIdIfExists(ctx, signInRequest.Tenant, signInRequest.LoggedInEmail)
				if err != nil {
					tracing.TraceErr(span, err)
					ginContext.JSON(http.StatusInternalServerError, gin.H{
						"result": fmt.Sprintf("unable to get email id: %v", err.Error()),
					})
					return
				}

				if emailId == "" {
					ginContext.JSON(http.StatusUnauthorized, gin.H{
						"result": fmt.Sprintf("email not found"),
					})
					return
				}

				tenantName = &signInRequest.Tenant
			}

			span.SetTag(tracing.SpanTagTenant, *tenantName)

			// Handle Google provider
			if signInRequest.Provider == "google" {
				if isRequestEnablingOAuthSync(signInRequest) {
					var oauthToken, _ = services.CommonServices.PostgresRepositories.OAuthTokenRepository.GetByEmail(ctx, *tenantName, signInRequest.Provider, signInRequest.OAuthTokenForEmail)
					if oauthToken == nil {
						oauthToken = &postgresEntity.OAuthTokenEntity{}
					}
					oauthToken.Provider = signInRequest.Provider
					oauthToken.TenantName = *tenantName
					oauthToken.PlayerIdentityId = signInRequest.OAuthToken.ProviderAccountId
					oauthToken.EmailAddress = signInRequest.OAuthTokenForEmail
					oauthToken.Type = signInRequest.OAuthTokenType
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
					_, err := services.CommonServices.PostgresRepositories.OAuthTokenRepository.Save(ctx, *oauthToken)
					if err != nil {
						log.Printf("unable to save oauth token: %v", err.Error())
						ginContext.JSON(http.StatusInternalServerError, gin.H{
							"result": fmt.Sprintf("unable to save oauth token: %v", err.Error()),
						})
						return
					}
				}
			} else if signInRequest.Provider == "azure-ad" {
				var oauthToken, _ = services.CommonServices.PostgresRepositories.OAuthTokenRepository.GetByEmail(ctx, *tenantName, signInRequest.Provider, signInRequest.OAuthTokenForEmail)
				if oauthToken == nil {
					oauthToken = &postgresEntity.OAuthTokenEntity{}
				}
				oauthToken.Provider = signInRequest.Provider
				oauthToken.TenantName = *tenantName
				oauthToken.PlayerIdentityId = signInRequest.OAuthToken.ProviderAccountId
				oauthToken.EmailAddress = signInRequest.OAuthTokenForEmail
				oauthToken.Type = signInRequest.OAuthTokenType
				oauthToken.AccessToken = signInRequest.OAuthToken.AccessToken
				oauthToken.RefreshToken = signInRequest.OAuthToken.RefreshToken
				oauthToken.IdToken = signInRequest.OAuthToken.IdToken
				oauthToken.ExpiresAt = signInRequest.OAuthToken.ExpiresAt
				oauthToken.Scope = signInRequest.OAuthToken.Scope
				oauthToken.NeedsManualRefresh = false
				_, err := services.CommonServices.PostgresRepositories.OAuthTokenRepository.Save(ctx, *oauthToken)
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
		security.ApiKeyCheckerHTTP(services.CommonServices.PostgresRepositories.TenantWebhookApiKeyRepository, services.CommonServices.PostgresRepositories.AppKeyRepository, security.USER_ADMIN_API),
		tracing.TracingEnhancer(context.Background(), "POST /revoke"),
		func(ginContext *gin.Context) {
			ctx, cancel := commonUtils.GetLongLivedContext(context.Background())
			defer cancel()

			var revokeRequest model.RevokeRequest
			if err := ginContext.BindJSON(&revokeRequest); err != nil {
				log.Printf("unable to parse json: %v", err.Error())
				ginContext.JSON(http.StatusInternalServerError, gin.H{
					"result": fmt.Sprintf("unable to parse json: %v", err.Error()),
				})
				return
			}
			log.Printf("parsed json: %v", revokeRequest)

			var oauthToken, _ = services.CommonServices.PostgresRepositories.OAuthTokenRepository.GetByEmail(ctx, revokeRequest.Tenant, revokeRequest.Provider, revokeRequest.Email)

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

					if resp.StatusCode != http.StatusOK {
						// Revocation failed
						ginContext.JSON(http.StatusInternalServerError, gin.H{})
						return
					}
				}
			}

			err := services.CommonServices.PostgresRepositories.OAuthTokenRepository.DeleteByEmail(ctx, revokeRequest.Tenant, revokeRequest.Provider, revokeRequest.Email)
			if err != nil {
				ginContext.JSON(http.StatusInternalServerError, gin.H{})
				return
			}

			ginContext.JSON(http.StatusOK, gin.H{})
		})
}

func getTenant(c context.Context, services *service.Services, personalEmailProvider []postgresEntity.PersonalEmailProvider, signInRequest model.SignInRequest, ginContext *gin.Context, config *config.Config) (*string, bool, error) {
	span, ctx := opentracing.StartSpanFromContext(c, "getTenant")
	defer span.Finish()

	domain := commonUtils.ExtractDomain(signInRequest.LoggedInEmail)
	span.LogFields(tracingLog.String("domain", domain))

	isPersonalEmail := false
	//check if the user is using a personal email provider
	for _, personalEmailProviderItem := range personalEmailProvider {
		domainLowercase := strings.ToLower(strings.TrimSpace(domain))
		personalEmailProviderDomainLowercase := strings.ToLower(strings.TrimSpace(personalEmailProviderItem.ProviderDomain))
		if domainLowercase == personalEmailProviderDomainLowercase {
			isPersonalEmail = true
			break
		}
	}

	span.LogFields(tracingLog.Bool("isPersonalEmail", isPersonalEmail))

	if isPersonalEmail {
		playerNode, err := services.CommonServices.Neo4jRepositories.PlayerReadRepository.GetPlayerByAuthIdProvider(ctx, signInRequest.LoggedInEmail, signInRequest.Provider)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, false, err
		}
		if playerNode != nil {
			span.LogFields(tracingLog.Object("playerIdentified", true))

			playerId := commonUtils.GetStringPropOrNil(playerNode.Props, "id")

			if playerId != nil {
				usersDb, err := services.CommonServices.Neo4jRepositories.PlayerReadRepository.GetUsersForPlayer(ctx, []string{*playerId})
				if err != nil {
					tracing.TraceErr(span, err)
					return nil, false, err
				}

				if usersDb == nil || len(usersDb) == 0 {
					tracing.TraceErr(span, fmt.Errorf("users not found"))
					return nil, false, fmt.Errorf("users not found")
				}

				tenantFromLabel := model2.GetTenantFromLabels(usersDb[0].Node.Labels, model2.NodeLabelUser)

				if tenantFromLabel == "" {
					tracing.TraceErr(span, fmt.Errorf("tenant not found"))
					return nil, false, fmt.Errorf("tenant not found")
				}

				span.LogFields(tracingLog.String("tenantIdentifiedFromPlayer", tenantFromLabel))
				return &tenantFromLabel, false, nil
			}
		} else {
			span.LogFields(tracingLog.Object("playerIdentified", false))
		}
	} else {
		tenantNode, err := services.CommonServices.Neo4jRepositories.TenantReadRepository.GetTenantForWorkspaceProvider(ctx, domain, signInRequest.Provider)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, false, err
		}
		if tenantNode != nil {
			tenant := neo4jmapper.MapDbNodeToTenantEntity(tenantNode)
			span.LogFields(tracingLog.String("tenantIdentifiedSameWorkspace", tenant.Name))
			return &tenant.Name, false, nil
		}

		//tenant not found by the requested login info, try to find it by another workspace with the same domain
		var provider string
		if signInRequest.Provider == "google" {
			provider = "azure-ad"
		} else if signInRequest.Provider == "azure-ad" {
			provider = "google"
		}
		tenantNode, err = services.CommonServices.Neo4jRepositories.TenantReadRepository.GetTenantForWorkspaceProvider(ctx, domain, provider)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, false, err
		}

		if tenantNode != nil {
			tenant := neo4jmapper.MapDbNodeToTenantEntity(tenantNode)
			span.LogFields(tracingLog.String("tenantIdentifiedDifferentWorkspace", tenant.Name))

			_, err := services.CommonServices.WorkspaceService.MergeToTenant(ctx, neo4jentity.WorkspaceEntity{
				Name:     domain,
				Provider: signInRequest.Provider,
			}, tenant.Name)
			if err != nil {
				tracing.TraceErr(span, err)
				return nil, false, err
			}

			return &tenant.Name, false, nil
		}
	}

	var tenantStr string
	if isPersonalEmail {
		tenantStr = utils.GenerateName()
	} else {
		tenantStr = utils.Sanitize(domain)
	}

	span.LogFields(tracingLog.String("newTenantCreationWith", tenantStr))

	tenantEntity, err := services.CommonServices.TenantService.Merge(ctx, neo4jentity.TenantEntity{
		Name: tenantStr,
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, false, err
	}

	if !isPersonalEmail {
		_, err := services.CommonServices.WorkspaceService.MergeToTenant(ctx, neo4jentity.WorkspaceEntity{
			Name:      domain,
			Provider:  signInRequest.Provider,
			AppSource: APP_SOURCE,
		}, tenantEntity.Name)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, false, err
		}
	}

	if config.Slack.NotifyNewTenantRegisteredWebhook != "" {
		commonUtils.SendSlackMessage(ctx, config.Slack.NotifyNewTenantRegisteredWebhook, tenantStr+" tenant registered by "+signInRequest.LoggedInEmail)
	}

	return &tenantEntity.Name, true, nil
}

func validateRequestAtProvider(c context.Context, config *config.Config, signInRequest model.SignInRequest) (*string, *string, error) {
	span, ctx := opentracing.StartSpanFromContext(c, "Registration.getUserInfoFromGoogle")
	defer span.Finish()

	if signInRequest.Provider == "google" {
		userInfo, err := getUserInfoFromGoogle(ctx, config, signInRequest)
		if err != nil {
			tracing.TraceErr(nil, err)
			return nil, nil, err
		}

		return &userInfo.GivenName, &userInfo.FamilyName, nil
	} else if signInRequest.Provider == "azure-ad" {
		client := &http.Client{}
		// Create a GET request with the Authorization header.
		req, err := http.NewRequest("GET", "https://graph.microsoft.com/oidc/userinfo", nil)
		if err != nil {
			tracing.TraceErr(nil, err)
			return nil, nil, err
		}

		req.Header.Set("Authorization", "Bearer "+signInRequest.OAuthToken.AccessToken)

		resp, err := client.Do(req)
		if err != nil {
			tracing.TraceErr(nil, err)
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
			tracing.TraceErr(nil, err)
			return nil, nil, err
		}
	} else {
		tracing.TraceErr(nil, fmt.Errorf("provider not supported"))
		return nil, nil, fmt.Errorf("provider not supported")
	}
}

func getUserInfoFromGoogle(c context.Context, config *config.Config, signInRequest model.SignInRequest) (*googleOauth.Userinfo, error) {
	span, ctx := opentracing.StartSpanFromContext(c, "Registration.getUserInfoFromGoogle")
	defer span.Finish()

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

	client := conf.Client(ctx, &token)

	oauth2Service, err := googleOauth.New(client)

	if err != nil {
		tracing.TraceErr(nil, err)
		return nil, err
	}
	userInfoService := googleOauth.NewUserinfoV2MeService(oauth2Service)

	userInfo, err := userInfoService.Get().Do()
	if err != nil {
		tracing.TraceErr(nil, err)
		return nil, err
	}

	return userInfo, nil
}

func initializeUser(c context.Context, services *service.Services, provider, providerAccountId, tenant, email string, firstName, lastName *string, ginContext *gin.Context) error {
	span, ctx := opentracing.StartSpanFromContext(c, "Registration.initializeUser")
	defer span.Finish()

	appSource := APP_SOURCE

	playerExists := false
	userExists := false
	var user *neo4jentity.UserEntity

	playerNode, err := services.CommonServices.Neo4jRepositories.PlayerReadRepository.GetPlayerByAuthIdProvider(ctx, email, provider)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	if playerNode != nil {
		playerExists = true
		span.LogFields(tracingLog.Object("player", "found"))
	} else {
		span.LogFields(tracingLog.Object("player", "not found"))
	}

	userNode, err := services.CommonServices.Neo4jRepositories.UserReadRepository.GetFirstUserByEmail(ctx, tenant, email)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	if userNode != nil {
		user = neo4jmapper.MapDbNodeToUserEntity(userNode)

		userExists = true
		span.LogFields(tracingLog.Object("user", "found"))
	} else {
		span.LogFields(tracingLog.Object("user", "not found"))
	}

	if !playerExists && !userExists {
		_, err := services.CustomerOsClient.CreateUser(&model.UserInput{
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
			tracing.TraceErr(span, err)
			return err
		}

		retries := 0

		for {
			userId, _, _, err := services.CommonServices.Neo4jRepositories.UserReadRepository.FindFirstUserWithRolesByEmail(ctx, email)
			if err != nil {
				tracing.TraceErr(span, err)
				return err
			}

			if userId != "" || retries > 25 {
				break
			}

			retries++
			time.Sleep(1 * time.Second)
		}

		userNode, err := services.CommonServices.Neo4jRepositories.UserReadRepository.GetFirstUserByEmail(ctx, tenant, email)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}
		if userNode != nil {
			user = neo4jmapper.MapDbNodeToUserEntity(userNode)
		} else {
			tracing.TraceErr(span, fmt.Errorf("user not created"))
			return fmt.Errorf("user not created")
		}
	} else {
		if !playerExists {
			err = services.CustomerOsClient.CreatePlayer(tenant, user.Id, providerAccountId, email, provider)
			if err != nil {
				tracing.TraceErr(span, err)
				return err
			}
		}
	}

	err = addDefaultMissingRoles(ctx, services, user, &tenant)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}

func addDefaultMissingRoles(ctx context.Context, services *service.Services, user *neo4jentity.UserEntity, tenant *string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Registration.addDefaultMissingRoles")
	defer span.Finish()

	var rolesToAdd []model.Role

	if user.Roles == nil || len(user.Roles) == 0 {
		rolesToAdd = []model.Role{model.RoleUser, model.RoleOwner}
	} else {
		userRoleFound := false
		ownerRoleFound := false
		for _, role := range user.Roles {
			if role == "USER" {
				userRoleFound = true
			}
			if role == "OWNER" {
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

	span.LogFields(tracingLog.Object("rolesToAdd", rolesToAdd))

	if len(rolesToAdd) > 0 {
		_, err := services.CustomerOsClient.AddUserRoles(*tenant, user.Id, rolesToAdd)
		if err != nil {
			tracing.TraceErr(span, err)
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

func isRequestEnablingMicrosoftSync(signInRequest model.SignInRequest) bool {
	if strings.Contains(signInRequest.OAuthToken.Scope, "Mail.ReadWrite") {
		return true
	}
	return false
}

func isRequestEnablingOAuthSync(signInRequest model.SignInRequest) bool {
	if isRequestEnablingGmailSync(signInRequest) || isRequestEnablingGoogleCalendarSync(signInRequest) || isRequestEnablingMicrosoftSync(signInRequest) {
		return true
	}
	return false
}

func registerNewTenantAsLeadInProviderTenant(ctx context.Context, config *config.Config, services *service.Services, registeredEmail string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "registration.registerNewTenantAsLeadInProviderTenant")
	defer span.Finish()

	organizationId, contactId, err := services.RegistrationService.CreateOrganizationAndContact(ctx, config.Service.ProviderTenantName, registeredEmail, true, "Tenant Registration")
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

	neo4jrepository.WaitForNodeCreatedInNeo4jWithConfig(ctx, span, services.CommonServices.Neo4jRepositories, *organizationId, model2.NodeLabelOrganization, 30*time.Second)
	neo4jrepository.WaitForNodeCreatedInNeo4jWithConfig(ctx, span, services.CommonServices.Neo4jRepositories, *contactId, model2.NodeLabelContact, 30*time.Second)
}
