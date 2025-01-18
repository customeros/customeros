package private

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/customeros/mailsherpa/mailvalidate"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/data_fields"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/dto"
	common_enum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/interfaces"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	common_srv "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services/postmark"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	common_utils "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neoEntity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	tracingLog "github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	tokenOauth "golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	googleOauth "google.golang.org/api/oauth2/v2"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/config"
	cosapi_services "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/services"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils"
)

func RML(s *cosapi_services.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		contextWithTimeout, cancel := common_utils.GetContextWithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(contextWithTimeout, "/rml", c.Request.Header)
		defer span.Finish()

		var request RequestMagicLinkRequest
		if err := c.BindJSON(&request); err != nil {
			tracing.TraceErr(span, err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"result": fmt.Sprintf("INVALID_REQUEST"),
			})
			return
		}

		if request.Email == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"result": fmt.Sprintf("EMAIL_EMPTY"),
			})
			return
		}

		saveErr := saveIP(ctx, c, s, request.Email)
		if saveErr != nil {
			tracing.TraceErr(span, errors.Wrap(saveErr, "unable to save IP address"))
		}

		emailValidation := mailvalidate.ValidateEmailSyntax(request.Email)

		if emailValidation.IsValid == false {
			c.JSON(http.StatusBadRequest, gin.H{
				"result": fmt.Sprintf("EMAIL_INVALID"),
			})
			return
		}

		var code string

		for {
			code = uuid.New().String() + uuid.New().String() + uuid.New().String()
			code = strings.ReplaceAll(code, "-", "")

			// Check if the code already exists
			magicLink, err := s.Repositories.PostgresRepositories.MagicLinkRepository.GetByCode(ctx, code)
			if err != nil {
				tracing.TraceErr(span, err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"result": fmt.Sprintf("INTERNAL_SERVER_ERROR"),
				})
				return
			}

			if magicLink == nil {
				break
			}
		}

		byEmail, err := s.Repositories.PostgresRepositories.MagicLinkRepository.GetByEmail(ctx, request.Email)
		if err != nil {
			tracing.TraceErr(span, err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"result": fmt.Sprintf("INTERNAL_SERVER_ERROR"),
			})
			return
		}

		if byEmail != nil {
			err := s.Repositories.PostgresRepositories.MagicLinkRepository.Delete(ctx, byEmail.ID)
			if err != nil {
				tracing.TraceErr(span, err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"result": fmt.Sprintf("INTERNAL_SERVER_ERROR"),
				})
				return
			}
		}

		link := "https://app.customeros.ai/auth/success?mg=" + code
		err = s.Repositories.PostgresRepositories.MagicLinkRepository.Create(ctx, &entity.MagicLink{
			Email: request.Email,
			Code:  code,
			Url:   link,
		})

		err = s.CommonServices.PostmarkService.SendNotification(ctx, interfaces.PostmarkEmail{
			MessageStream: postmark.PostmarkMessageStreamMagicLink,
			WorkflowId:    postmark.WorkflowMagicLink,
			Subject:       postmark.WorkflowMagicLinkSubject,
			From:          "notification@app.customeros.ai",
			To:            request.Email,
			TemplateData: map[string]string{
				"{{magicLink}}": link,
			},
		}, "openlineai")
		if err != nil {
			tracing.TraceErr(span, err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"result": fmt.Sprintf("INTERNAL_SERVER_ERROR"),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"result": "OK",
		})
	}
}

func PML(s *cosapi_services.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		// todo move to server init?
		personalEmailProviders, err := s.Repositories.PostgresRepositories.PersonalEmailProviderRepository.GetPersonalEmailProviders()
		if err != nil {
			panic(err)
		}
		contextWithTimeout, cancel := common_utils.GetContextWithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(contextWithTimeout, "/pml", c.Request.Header)
		defer span.Finish()

		var magicLink *entity.MagicLink
		var signInRequest SignInRequest
		if err := c.BindJSON(&signInRequest); err != nil {
			tracing.TraceErr(span, err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"result": fmt.Sprintf("INVALID_REQUEST"),
			})
			return
		}

		if signInRequest.Code != "" {
			magicLink, err := s.Repositories.PostgresRepositories.MagicLinkRepository.GetByCode(ctx, signInRequest.Code)
			if err != nil {
				tracing.TraceErr(span, err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"result": fmt.Sprintf("INTERNAL_SERVER_ERROR"),
				})
				return
			}

			if magicLink == nil {
				c.JSON(http.StatusUnauthorized, gin.H{
					"result": "MAGIC_LINK_NOT_FOUND",
				})
				return
			}

			signInRequest.Provider = common_enum.WorkspaceProviderMagicLink.String()
			signInRequest.LoggedInEmail = magicLink.Email
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"result": "MAGIC_LINK_NOT_FOUND",
			})
		}

		signIn(ctx, s, c, signInRequest, personalEmailProviders, s.Cfg)

		err = s.Repositories.PostgresRepositories.MagicLinkRepository.Delete(ctx, magicLink.ID)
		if err != nil {
			tracing.TraceErr(span, err)
		}

		return
	}
}

func Signin(s *cosapi_services.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		contextWithTimeout, cancel := common_utils.GetContextWithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		personalEmailProviders, err := s.Repositories.PostgresRepositories.PersonalEmailProviderRepository.GetPersonalEmailProviders()
		if err != nil {
			panic(err)
		}

		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(contextWithTimeout, "/signin", c.Request.Header)
		defer span.Finish()

		var signInRequest SignInRequest
		if err := c.BindJSON(&signInRequest); err != nil {
			tracing.TraceErr(span, err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"result": fmt.Sprintf("unable to parse json: %v", err.Error()),
			})
			return
		}

		signIn(ctx, s, c, signInRequest, personalEmailProviders, s.Cfg)
	}
}

func Revoke(s *cosapi_services.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := common_utils.GetLongLivedContext(context.Background())
		defer cancel()

		var revokeRequest RevokeRequest
		if err := c.BindJSON(&revokeRequest); err != nil {
			log.Printf("unable to parse json: %v", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"result": fmt.Sprintf("unable to parse json: %v", err.Error()),
			})
			return
		}
		log.Printf("parsed json: %v", revokeRequest)

		oauthToken, _ := s.Repositories.PostgresRepositories.OAuthTokenRepository.GetByEmail(ctx, revokeRequest.Tenant, revokeRequest.Provider, revokeRequest.Email)

		if oauthToken != nil && oauthToken.RefreshToken != "" {
			// Handle revocation based on provider
			var revocationURL string
			switch revokeRequest.Provider {
			case common_enum.WorkspaceProviderGoogle.String():
				revocationURL = fmt.Sprintf("https://accounts.google.com/o/oauth2/revoke?token=%s", oauthToken.RefreshToken)
			case common_enum.WorkspaceProviderAzure.String():
				revocationURL = fmt.Sprintf("https://graph.microsoft.com/v1.0/me/revokeSignInSessions")
			}

			if revocationURL != "" {
				resp, err := http.Get(revocationURL)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{})
					return
				}

				if resp.StatusCode != http.StatusOK {
					// Revocation failed
					c.JSON(http.StatusInternalServerError, gin.H{})
					return
				}
			}
		}

		err := s.Repositories.PostgresRepositories.OAuthTokenRepository.DeleteByEmail(ctx, revokeRequest.Tenant, revokeRequest.Provider, revokeRequest.Email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{})
			return
		}

		c.JSON(http.StatusOK, gin.H{})
	}
}

func signIn(ctx context.Context, services *cosapi_services.Services, ginContext *gin.Context, signInRequest SignInRequest, personalEmailProviders []entity.PersonalEmailProvider, config *config.Config) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "getTenant")
	defer span.Finish()

	var err error

	saveErr := saveIP(ctx, ginContext, services, signInRequest.LoggedInEmail)
	if saveErr != nil {
		tracing.TraceErr(span, errors.Wrap(err, "unable to save IP address"))
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
	var userId *string

	if signInRequest.Tenant == "" {
		span.LogFields(tracingLog.String("flow", "authentication"))

		tn, isNewTenant, err := getTenant(ctx, services, personalEmailProviders, signInRequest, config)
		if err != nil {
			tracing.TraceErr(span, err)
			ginContext.JSON(http.StatusInternalServerError, gin.H{
				"result": fmt.Sprintf("unable to get tenant: %v", err.Error()),
			})
			return
		}
		tenantName = tn

		ctx = common.WithCustomContext(ctx, &common.CustomContext{
			Tenant:    *tenantName,
			AppSource: constants.AppSourceUserAdminApi,
		})
		userId, err = initializeUser(ctx, services, signInRequest.Provider, signInRequest.OAuthToken.ProviderAccountId, *tenantName, signInRequest.LoggedInEmail, firstName, lastName)
		if err != nil {
			tracing.TraceErr(span, err)
			ginContext.JSON(http.StatusInternalServerError, gin.H{
				"result": fmt.Sprintf("unable to initialize user: %v", err.Error()),
			})
			return
		}

		if isNewTenant {
			domain := common_utils.ExtractDomain(signInRequest.LoggedInEmail)
			isPersonalEmail := false
			// check if the user is using a personal email provider
			for _, personalEmailProviderItem := range personalEmailProviders {
				domainLowercase := strings.ToLower(strings.TrimSpace(domain))
				personalEmailProviderDomainLowercase := strings.ToLower(strings.TrimSpace(personalEmailProviderItem.ProviderDomain))
				if domainLowercase == personalEmailProviderDomainLowercase {
					isPersonalEmail = true
					break
				}
			}

			if !isPersonalEmail {
				err = services.CommonServices.RegistrationService.PrepareDefaultTenantSetup(ctx, signInRequest.LoggedInEmail)
				if err != nil {
					tracing.TraceErr(span, err)
				}
			}

			go func() {
				c, cancelFunc := context.WithTimeout(context.Background(), 300*time.Second)
				defer cancelFunc()

				ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c, "/signin - register new tenant", ginContext.Request.Header)
				defer span.Finish()

				err = registerNewTenantAsLeadInProviderTenant(ctx, config, services, signInRequest.LoggedInEmail)
				if err != nil {
					tracing.TraceErr(span, err)
					return
				}

				span.LogFields(tracingLog.String("result", "ok"))
			}()
		} else {

			domain := common_utils.ExtractDomain(signInRequest.LoggedInEmail)

			isPersonalEmail := false
			// check if the user is using a personal email provider
			for _, personalEmailProviderItem := range personalEmailProviders {
				domainLowercase := strings.ToLower(strings.TrimSpace(domain))
				personalEmailProviderDomainLowercase := strings.ToLower(strings.TrimSpace(personalEmailProviderItem.ProviderDomain))
				if domainLowercase == personalEmailProviderDomainLowercase {
					isPersonalEmail = true
					break
				}
			}

			if !isPersonalEmail {
				err = services.CommonServices.RegistrationService.PrepareDefaultTenantSetup(ctx, signInRequest.LoggedInEmail)
				if err != nil {
					tracing.TraceErr(span, err)
				}
			}
		}
	} else {
		span.LogFields(tracingLog.String("flow", "authorization"))

		userDbNode, err := services.Repositories.Neo4jRepositories.UserReadRepository.GetFirstUserByEmail(ctx, signInRequest.Tenant, signInRequest.LoggedInEmail)
		if err != nil {
			tracing.TraceErr(span, err)
			ginContext.JSON(http.StatusInternalServerError, gin.H{
				"result": fmt.Sprintf("unable to get email id: %v", err.Error()),
			})
			return
		}

		if userDbNode == nil {
			ginContext.JSON(http.StatusUnauthorized, gin.H{
				"result": fmt.Sprintf("email not found"),
			})
			return
		}

		tenantName = &signInRequest.Tenant
	}

	span.SetTag(tracing.SpanTagTenant, *tenantName)

	// Handle Google provider
	if signInRequest.Provider == common_enum.WorkspaceProviderGoogle.String() {
		if isRequestEnablingOAuthSync(signInRequest) {
			oauthToken, _ := services.Repositories.PostgresRepositories.OAuthTokenRepository.GetByEmail(ctx, *tenantName, signInRequest.Provider, signInRequest.OAuthTokenForEmail)
			if oauthToken == nil {
				oauthToken = &entity.OAuthTokenEntity{}
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
			_, err := services.Repositories.PostgresRepositories.OAuthTokenRepository.Save(ctx, *oauthToken)
			if err != nil {
				log.Printf("unable to save oauth token: %v", err.Error())
				ginContext.JSON(http.StatusInternalServerError, gin.H{
					"result": fmt.Sprintf("unable to save oauth token: %v", err.Error()),
				})
				return
			}
		}
	} else if signInRequest.Provider == common_enum.WorkspaceProviderAzure.String() {
		oauthToken, _ := services.Repositories.PostgresRepositories.OAuthTokenRepository.GetByEmail(ctx, *tenantName, signInRequest.Provider, signInRequest.OAuthTokenForEmail)
		if oauthToken == nil {
			oauthToken = &entity.OAuthTokenEntity{}
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
		_, err := services.Repositories.PostgresRepositories.OAuthTokenRepository.Save(ctx, *oauthToken)
		if err != nil {
			log.Printf("unable to save oauth token: %v", err.Error())
			ginContext.JSON(http.StatusInternalServerError, gin.H{
				"result": fmt.Sprintf("unable to save oauth token: %v", err.Error()),
			})
			return
		}
	} else if signInRequest.Provider == common_enum.WorkspaceProviderMagicLink.String() {
	} else {
		log.Printf("Unsupported provider: %s", signInRequest.Provider)
		ginContext.JSON(http.StatusBadRequest, gin.H{
			"result": fmt.Sprintf("Unsupported provider: %s", signInRequest.Provider),
		})
		return
	}

	ginContext.JSON(http.StatusOK, gin.H{
		"status": "OK",
		"email":  signInRequest.LoggedInEmail,
		"tenant": *tenantName,
		"userId": userId,
	})
}

func getTenant(c context.Context, services *cosapi_services.Services, personalEmailProvider []entity.PersonalEmailProvider, signInRequest SignInRequest, config *config.Config) (*string, bool, error) {
	span, ctx := opentracing.StartSpanFromContext(c, "getTenant")
	defer span.Finish()

	domain := common_utils.ExtractDomain(signInRequest.LoggedInEmail)
	span.LogFields(tracingLog.String("domain", domain))

	isPersonalEmail := false
	// check if the user is using a personal email provider
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
		playerNode, err := services.Repositories.Neo4jRepositories.PlayerReadRepository.GetPlayerByAuthIdProvider(ctx, signInRequest.LoggedInEmail, signInRequest.Provider)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, false, err
		}
		if playerNode != nil {
			span.LogFields(tracingLog.Object("playerIdentified", true))

			playerId := mapper.MapDbNodeToPlayerEntity(playerNode).Id

			usersDb, err := services.Repositories.Neo4jRepositories.PlayerReadRepository.GetUsersForPlayer(ctx, []string{playerId})
			if err != nil {
				tracing.TraceErr(span, err)
				return nil, false, err
			}

			if usersDb == nil || len(usersDb) == 0 {
				tracing.TraceErr(span, fmt.Errorf("users not found"))
				return nil, false, fmt.Errorf("users not found")
			}

			tenantFromLabel := model.GetTenantFromLabels(usersDb[0].Node.Labels, model.NodeLabelUser)

			if tenantFromLabel == "" {
				tracing.TraceErr(span, fmt.Errorf("tenant not found"))
				return nil, false, fmt.Errorf("tenant not found")
			}

			span.LogFields(tracingLog.String("tenantIdentifiedFromPlayer", tenantFromLabel))
			return &tenantFromLabel, false, nil
		} else {
			span.LogFields(tracingLog.Object("playerIdentified", false))
		}
	} else {
		tenantNode, err := services.Repositories.Neo4jRepositories.TenantReadRepository.GetTenantForWorkspaceProvider(ctx, domain, signInRequest.Provider)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, false, err
		}
		if tenantNode != nil {
			tenant := mapper.MapDbNodeToTenantEntity(tenantNode)
			span.LogFields(tracingLog.String("tenantIdentifiedSameWorkspace", tenant.Name))
			return &tenant.Name, false, nil
		}

		// tenant not found by the requested login info, try to find it by another workspace with the same domain
		tenantNode, err = services.Repositories.Neo4jRepositories.TenantReadRepository.GetTenantForWorkspace(ctx, domain)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, false, err
		}

		if tenantNode != nil {
			tenant := mapper.MapDbNodeToTenantEntity(tenantNode)
			span.LogFields(tracingLog.String("tenantIdentifiedDifferentWorkspace", tenant.Name))

			_, err := services.CommonServices.WorkspaceService.MergeToTenant(ctx, neoEntity.WorkspaceEntity{
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

	tenantEntity, err := services.CommonServices.TenantService.Merge(ctx, neoEntity.TenantEntity{
		Name: tenantStr,
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, false, err
	}

	if !isPersonalEmail {
		_, err := services.CommonServices.WorkspaceService.MergeToTenant(ctx, neoEntity.WorkspaceEntity{
			Name:      domain,
			Provider:  signInRequest.Provider,
			AppSource: constants.AppSourceUserAdminApi,
		}, tenantEntity.Name)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, false, err
		}
	}

	if config.CommonServices.External.SlackConfig.NotifyNewTenantRegisteredHook != "" {
		common_utils.SendSlackMessage(ctx, config.CommonServices.External.SlackConfig.NotifyNewTenantRegisteredHook, tenantStr+" tenant registered by "+signInRequest.LoggedInEmail)
	}

	return &tenantEntity.Name, true, nil
}

func validateRequestAtProvider(c context.Context, config *config.Config, signInRequest SignInRequest) (*string, *string, error) {
	span, ctx := opentracing.StartSpanFromContext(c, "Registration.getUserInfoFromGoogle")
	defer span.Finish()

	if signInRequest.Provider == common_enum.WorkspaceProviderMagicLink.String() {
		return nil, nil, nil
	} else if signInRequest.Provider == common_enum.WorkspaceProviderGoogle.String() {
		userInfo, err := getUserInfoFromGoogle(ctx, config, signInRequest)
		if err != nil {
			tracing.TraceErr(nil, err)
			return nil, nil, err
		}

		return &userInfo.GivenName, &userInfo.FamilyName, nil
	} else if signInRequest.Provider == common_enum.WorkspaceProviderAzure.String() {
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

func getUserInfoFromGoogle(c context.Context, config *config.Config, signInRequest SignInRequest) (*googleOauth.Userinfo, error) {
	span, ctx := opentracing.StartSpanFromContext(c, "Registration.getUserInfoFromGoogle")
	defer span.Finish()

	conf := &tokenOauth.Config{
		ClientID:     config.CommonServices.Infrastructure.GoogleOAuthConfig.ClientId,
		ClientSecret: config.CommonServices.Infrastructure.GoogleOAuthConfig.ClientSecret,
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

func initializeUser(c context.Context, services *cosapi_services.Services, provider, providerAccountId, tenant, email string, firstName, lastName *string) (*string, error) {
	span, ctx := opentracing.StartSpanFromContext(c, "Registration.initializeUser")
	defer span.Finish()

	innerCtx := common.WithCustomContext(ctx, &common.CustomContext{
		Tenant:    tenant,
		AppSource: constants.AppSourceUserAdminApi,
	})

	userId := ""
	playerId := ""

	userNode, err := services.Repositories.Neo4jRepositories.UserReadRepository.GetFirstUserByEmail(ctx, tenant, email)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	if userNode != nil {
		userId = mapper.MapDbNodeToUserEntity(userNode).Id
		span.LogFields(tracingLog.Object("user", "found"))
	} else {
		span.LogFields(tracingLog.Object("user", "not found"))
	}

	playerNode, err := services.Repositories.Neo4jRepositories.PlayerReadRepository.GetPlayerByAuthIdProvider(ctx, email, provider)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	if playerNode != nil {
		playerId = mapper.MapDbNodeToPlayerEntity(playerNode).Id
		span.LogFields(tracingLog.Object("player", "found"))
	} else {
		span.LogFields(tracingLog.Object("player", "not found"))
	}

	defaultWorkSchedule := entity.UserWorkingSchedule{
		UserId:    userId,
		DayRange:  "Mon-Fri",
		StartHour: "09:00",
		EndHour:   "18:00",
	}

	if userId == "" {
		userId, err = services.CommonServices.UserService.Save(innerCtx, nil, nil, data_fields.UserFields{
			FirstName: firstName,
			LastName:  lastName,
			Roles:     common_utils.ToPtr([]string{"USER", "OWNER"}),
		})

		_, err = services.CommonServices.EmailService.Merge(innerCtx, nil, tenant, interfaces.EmailFields{
			Primary:   true,
			Email:     email,
			Source:    neoEntity.DataSourceOpenline,
			AppSource: constants.AppSourceUserAdminApi,
		}, &common_srv.LinkWith{
			Type: model.USER,
			Id:   userId,
		})
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "unable to link user to email"))
			return nil, err
		}

		err = services.Repositories.PostgresRepositories.UserWorkingScheduleRepository.Store(innerCtx, tenant, &defaultWorkSchedule)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}
	} else {
		err = addDefaultMissingRoles(ctx, services, tenant, userId)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}

		workingSchedule, err := services.Repositories.PostgresRepositories.UserWorkingScheduleRepository.GetForUser(ctx, tenant, userId)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}

		if len(workingSchedule) == 0 {
			err = services.Repositories.PostgresRepositories.UserWorkingScheduleRepository.Store(innerCtx, tenant, &defaultWorkSchedule)
			if err != nil {
				tracing.TraceErr(span, err)
				return nil, err
			}
		}
	}

	if playerId == "" {
		err := services.Repositories.Neo4jRepositories.PlayerWriteRepository.Merge(ctx, userId, neoEntity.PlayerEntity{
			AuthId:     email,
			Provider:   provider,
			IdentityId: providerAccountId,
			AppSource:  constants.AppSourceUserAdminApi,
		})
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}
	}

	if userId != "" {
		innerCtx := common.WithCustomContext(ctx, &common.CustomContext{
			Tenant:    tenant,
			AppSource: constants.AppSourceUserAdminApi,
			UserEmail: email,
		})
		err = services.Repositories.Neo4jRepositories.UserWriteRepository.RegisterLogin(innerCtx, tenant, userId)
		if err != nil {
			tracing.TraceErr(span, err)
		}
		err = services.CommonServices.Events.Publisher.PublishEvent(innerCtx, userId, model.USER, dto.UserLogin{LoginEmail: email, Provider: provider, IdentityId: providerAccountId})
		if err != nil {
			tracing.TraceErr(span, err)
		}
	}

	return &userId, nil
}

func addDefaultMissingRoles(c context.Context, services *cosapi_services.Services, tenant, userId string) error {
	span, ctx := opentracing.StartSpanFromContext(c, "Registration.addDefaultMissingRoles")
	defer span.Finish()

	userRoleFound := false
	ownerRoleFound := false

	userNode, err := services.Repositories.Neo4jRepositories.UserReadRepository.GetUserById(ctx, tenant, userId)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	if userNode == nil {
		tracing.TraceErr(span, fmt.Errorf("user not found"))
		return fmt.Errorf("user not found")
	}

	existingUser := mapper.MapDbNodeToUserEntity(userNode)

	if existingUser.Roles != nil && len(existingUser.Roles) > 0 {
		for _, role := range existingUser.Roles {
			if role == "USER" {
				userRoleFound = true
			}
			if role == "OWNER" {
				ownerRoleFound = true
			}
		}
	}

	if !userRoleFound {
		err := services.Repositories.Neo4jRepositories.UserWriteRepository.AddRole(ctx, existingUser.Id, "USER")
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}
	}
	if !ownerRoleFound {
		err := services.Repositories.Neo4jRepositories.UserWriteRepository.AddRole(ctx, existingUser.Id, "OWNER")
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}
	}

	return nil
}

func isRequestEnablingGmailSync(signInRequest SignInRequest) bool {
	if strings.Contains(signInRequest.OAuthToken.Scope, "gmail") {
		return true
	}
	return false
}

func isRequestEnablingGoogleCalendarSync(signInRequest SignInRequest) bool {
	if strings.Contains(signInRequest.OAuthToken.Scope, "calendar") {
		return true
	}
	return false
}

func isRequestEnablingMicrosoftSync(signInRequest SignInRequest) bool {
	if strings.Contains(signInRequest.OAuthToken.Scope, "Mail.ReadWrite") {
		return true
	}
	return false
}

func isRequestEnablingOAuthSync(signInRequest SignInRequest) bool {
	if isRequestEnablingGmailSync(signInRequest) || isRequestEnablingGoogleCalendarSync(signInRequest) || isRequestEnablingMicrosoftSync(signInRequest) {
		return true
	}
	return false
}

func registerNewTenantAsLeadInProviderTenant(ctx context.Context, config *config.Config, services *cosapi_services.Services, registeredEmail string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "registration.registerNewTenantAsLeadInProviderTenant")
	defer span.Finish()

	providerTenantCtx := common.WithCustomContext(ctx, &common.CustomContext{
		Tenant:    config.App.AuthConfig.ProviderTenantName,
		AppSource: constants.AppSourceUserAdminApi,
	})

	organizationId, contactId, err := createOrganizationAndContact(providerTenantCtx, services, config.App.AuthConfig.ProviderTenantName, registeredEmail, true, "Tenant Registration")
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	span.LogFields(tracingLog.String("providerOrganizationId", *organizationId))
	span.LogFields(tracingLog.String("providerContactId", *contactId))

	// TODO EDI - send welcome email

	//	type EmailPayload struct {
	//		Channel   string   `json:"channel"`
	//		Username  string   `json:"username"`
	//		Direction string   `json:"direction"`
	//		To        []string `json:"to"`
	//		Cc        []string `json:"cc"`
	//		Bcc       []string `json:"bcc"`
	//		Content   string   `json:"content"`
	//		Subject   string   `json:"subject"`
	//	}
	//
	//	payload := EmailPayload{
	//		Channel:   "EMAIL",
	//		Username:  config.Service.ProviderUsername,
	//		Direction: "OUTBOUND",
	//		To:        []string{registeredEmail},
	//		Cc:        []string{},
	//		Bcc:       []string{},
	//		Subject:   "Welcome to CustomerOS",
	//		Content: `<!DOCTYPE html>
	//<html lang="en">
	//<head>
	//    <meta charset="UTF-8">
	//    <meta name="viewport" content="width=device-width, initial-scale=1.0">
	//    <title>Welcome to CustomerOS</title>
	//    <style>
	//        body {
	//            font-family: Arial, sans-serif;
	//            line-height: 1.6;
	//        }
	//        .header {
	//            font-size: 1.4em;
	//            margin-bottom: 10px;
	//        }
	//        .content {
	//            margin-bottom: 20px;
	//        }
	//        .footer {
	//            font-size: 0.9em;
	//            color: #888;
	//        }
	//        .signature {
	//            margin-top: 20px;
	//        }
	//    </style>
	//</head>
	//<body>
	//<div>
	//    <div class="header">Hey, welcome to CustomerOS!</div>
	//    <div class="content">
	//        <p>Thanks for trying us out.</p>
	//        <p>To be honest, our self-service onboarding kinda sucks right now as we’re still building it out. I’d love to get you setup and using the tool, would you be open to spending 10 mins with me to help you get things configured?</p>
	//        <p>Please grab any slot on my <a href="https://app.customeros.ai/organization/cal.com/mbrown/20min" target="_blank">calendar.</a></p>
	//    </div>
	//    <div class="signature">
	//        <p>Thanks again,</p>
	//        <p>Matt Brown<br>
	//            CEO @ <a href="https://customeros.ai/?utm_source=signature&utm_medium=email&utm_campaign=signup" target="_blank">CustomerOS</a></p>
	//        <p class="footer">
	//            Follow me on <a href="https://www.linkedin.com/in/mateocafe/" target="_blank">LinkedIn</a><br>
	//            US: +1 650 977 2199<br>
	//            UK: +44 7700 155 600
	//        </p>
	//    </div>
	//</div>
	//</body>
	//</html>`,
	//	}

	//payloadBytes, err := json.Marshal(payload)
	//if err != nil {
	//	tracing.TraceErr(span, err)
	//	return err
	//}

	// body := bytes.NewBuffer(payloadBytes)

	// Create new request
	//req, err := http.NewRequest(method, url+"/mail/send", body)
	//if err != nil {
	//	tracing.TraceErr(span, err)
	//	return err
	//}
	//
	//// Set headers
	//req.Header.Add("X-Openline-Mail-Api-Key", config.Comms.CommsAPIKey)
	//req.Header.Add("X-Openline-USERNAME", config.Service.ProviderUsername)
	//req.Header.Add("Content-Type", "application/json")
	//
	//// Create HTTP client and make request
	//client := &http.Client{}
	//res, err := client.Do(req)
	//if err != nil {
	//	tracing.TraceErr(span, err)
	//	return err
	//}
	//defer res.Body.Close()
	//
	//mapBody := make(map[string]interface{})
	//
	//// Read response
	//responseBody, err := ioutil.ReadAll(res.Body)
	//if err != nil {
	//	tracing.TraceErr(span, err)
	//	return err
	//}
	//
	//// Convert response to map
	//err = json.Unmarshal(responseBody, &mapBody)
	//if err != nil {
	//	tracing.TraceErr(span, err)
	//	return err
	//}
	//
	//if mapBody["error"] != nil {
	//	tracing.TraceErr(span, fmt.Errorf("error: %v", mapBody["error"]))
	//	return fmt.Errorf("error: %v", mapBody["error"])
	//}

	// span.LogFields(tracingLog.Object("email sent: ", mapBody))

	return nil
}

func createOrganizationAndContact(ctx context.Context, services *cosapi_services.Services, tenant, email string, allowPersonalEmail bool, leadSource string) (*string, *string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "RegistrationService.CreateOrganizationAndContact")
	defer span.Finish()

	domain := common_utils.ExtractDomain(email)

	isPersonalEmail := false
	// check if the user is using a personal email provider
	for _, personalEmailProvider := range services.CommonServices.Cache.GetPersonalEmailProviders() {
		if strings.Contains(domain, personalEmailProvider) {
			isPersonalEmail = true
			break
		}
	}

	organizationId := ""
	contactId := ""

	if !isPersonalEmail || allowPersonalEmail {
		organizationByDomain, err := services.Repositories.Neo4jRepositories.OrganizationReadRepository.GetOrganizationByDomain(ctx, nil, tenant, domain)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, nil, err
		}

		if organizationByDomain == nil {
			organizationId, err = services.CommonServices.OrganizationService.Save(ctx, nil, nil, data_fields.OrganizationFields{
				Domains:      []string{domain},
				Name:         common_utils.StringPtr(domain),
				Relationship: common_utils.ToPtr(enum.OrganizationRelationshipProspect),
				Stage:        common_utils.ToPtr(enum.Trial),
				LeadSource:   common_utils.StringPtr(leadSource),
			})
			if err != nil {
				tracing.TraceErr(span, err)
				return nil, nil, err
			}
			if organizationId == "" {
				e := errors.New("organization id empty")
				tracing.TraceErr(span, e)
				return nil, nil, e
			}
		} else {
			organizationId = mapper.MapDbNodeToOrganizationEntity(organizationByDomain).ID
		}

		if organizationId == "" {
			tracing.TraceErr(span, err)
			return nil, nil, err
		}
		span.LogFields(tracingLog.String("result.organizationId", organizationId))

		contactNode, err := services.Repositories.Neo4jRepositories.ContactReadRepository.GetContactInOrganizationByEmail(ctx, tenant, organizationId, email)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, nil, err
		}

		if contactNode == nil {
			contactId, err = services.CommonServices.ContactService.CreateContactByEmail(ctx, nil, email)
			if err != nil {
				tracing.TraceErr(span, err)
				return nil, nil, err
			}

			err = services.CommonServices.ContactService.LinkContactWithOrganization(ctx, nil, contactId, organizationId, "", "",
				neoEntity.DataSourceOpenline.String(), false, nil, nil)
			if err != nil {
				tracing.TraceErr(span, err)
				return nil, nil, err
			}
		} else {
			contactId = mapper.MapDbNodeToContactEntity(contactNode).Id
		}

		if contactId == "" {
			tracing.TraceErr(span, errors.New("contact id empty"))
			return nil, nil, errors.New("contact id empty")
		}
		span.LogFields(tracingLog.String("result.contactId", contactId))
	}

	return &organizationId, &contactId, nil
}

func saveIP(ctx context.Context, c *gin.Context, s *cosapi_services.Services, email string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "registration.saveIP")
	defer span.Finish()

	var clientIP string
	var originalIP string
	var cloudflareIP string

	if c.Request.Header["X-Original-Forwarded-For"] == nil && c.Request.Header["Cf-Connecting-Ip"] == nil {
		return nil
	}

	if c.Request.Header["X-Original-Forwarded-For"] != nil {
		originalIP = c.Request.Header["X-Original-Forwarded-For"][0]
	}

	if c.Request.Header["Cf-Connecting-Ip"] != nil {
		cloudflareIP = c.Request.Header["Cf-Connecting-Ip"][0]
	}

	if cloudflareIP == "" && originalIP == "" {
		return nil
	}

	if cloudflareIP != "" {
		clientIP = cloudflareIP
	} else {
		clientIP = originalIP
	}

	validEmail := mailvalidate.ValidateEmailSyntax(email)
	if !validEmail.IsValid {
		err := errors.New("Email is invalid")
		span.LogKV("email", email)
		tracing.TraceErr(span, err)
	}

	if validEmail.IsFreeAccount || validEmail.IsRoleAccount || validEmail.IsSystemGenerated {
		return nil
	}

	website := fmt.Sprintf("https://%s", validEmail.Domain)

	details := entity.EnrichDetailsTracking{
		IP:             clientIP,
		CompanyDomain:  &validEmail.Domain,
		CompanyWebsite: &website,
		SourceEmail:    &validEmail.CleanEmail,
	}
	tracing.LogObjectAsJson(span, "ipToEmailDetails", details)

	err := s.Repositories.PostgresRepositories.EnrichDetailsTrackingRepository.Save(ctx, details)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}
