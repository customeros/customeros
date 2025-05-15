package private

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	nats_common "github.com/customeros/customeros/packages/server/customer-os-common-module/nats"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/proto/pb"
	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/clients"
	commonconfig "github.com/customeros/customeros/packages/server/customer-os-common-module/config"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/coserrors"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/data_fields"
	commonenum "github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/postmark"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	common_utils "github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	neoEntity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"
	"github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/mapper"
	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	"github.com/customeros/mailsherpa/mailvalidate"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	tokenOauth "golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	googleOauth "google.golang.org/api/oauth2/v2"

	"github.com/customeros/customeros/packages/server/customer-os-api/config"
	"github.com/customeros/customeros/packages/server/customer-os-api/graphql/model"
	cosapi_services "github.com/customeros/customeros/packages/server/customer-os-api/services"
	"github.com/customeros/customeros/packages/server/customer-os-api/utils"
)

func RML(s *cosapi_services.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		contextWithTimeout, cancel := common_utils.GetContextWithTimeout(c.Request.Context(), 30*time.Second)
		defer cancel()

		spans, ctx := telemetry.StartRestSpan(contextWithTimeout, "RML")
		defer spans.Finish()

		var request RequestMagicLinkRequest
		if err := c.BindJSON(&request); err != nil {
			spans.TraceError(err)
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
			spans.TraceError(errors.Wrap(saveErr, "unable to save IP address"))
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
				spans.TraceError(err)
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
			spans.TraceError(err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"result": fmt.Sprintf("INTERNAL_SERVER_ERROR"),
			})
			return
		}

		if byEmail != nil {
			err := s.Repositories.PostgresRepositories.MagicLinkRepository.Delete(ctx, byEmail.ID)
			if err != nil {
				spans.TraceError(err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"result": fmt.Sprintf("INTERNAL_SERVER_ERROR"),
				})
				return
			}
		}

		link := "https://app.customeros.ai/auth/success?mg=" + code
		err = s.Repositories.PostgresRepositories.MagicLinkRepository.Create(ctx, &postgres_entity.MagicLink{
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
			spans.TraceError(err)
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
		contextWithTimeout, cancel := common_utils.GetContextWithTimeout(c.Request.Context(), 30*time.Second)
		defer cancel()

		spans, ctx := telemetry.StartRestSpan(contextWithTimeout, "PML")
		defer spans.Finish()

		var magicLink *postgres_entity.MagicLink
		var signInRequest SignInRequest
		if err := c.BindJSON(&signInRequest); err != nil {
			spans.TraceError(err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"result": fmt.Sprintf("INVALID_REQUEST"),
			})
			return
		}

		var err error

		if signInRequest.Code != "" {
			magicLink, err = s.Repositories.PostgresRepositories.MagicLinkRepository.GetByCode(ctx, signInRequest.Code)
			if err != nil {
				spans.TraceError(err)
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

			signInRequest.Provider = commonenum.WorkspaceProviderMagicLink.String()
			signInRequest.LoggedInEmail = magicLink.Email
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"result": "MAGIC_LINK_NOT_FOUND",
			})
		}

		signIn(ctx, s, c, signInRequest, s.Cfg)

		err = s.Repositories.PostgresRepositories.MagicLinkRepository.Delete(ctx, magicLink.ID)
		if err != nil {
			spans.TraceError(err)
		}

		return
	}
}

func Signin(s *cosapi_services.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		contextWithTimeout, cancel := common_utils.GetContextWithTimeout(c.Request.Context(), 30*time.Second)
		defer cancel()

		spans, ctx := telemetry.StartRestSpan(contextWithTimeout, "Signin")
		defer spans.Finish()

		var signInRequest SignInRequest
		if err := c.BindJSON(&signInRequest); err != nil {
			spans.TraceError(err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"result": fmt.Sprintf("unable to parse json: %v", err.Error()),
			})
			return
		}

		signIn(ctx, s, c, signInRequest, s.Cfg)
	}
}

func signIn(ctx context.Context, services *cosapi_services.Services, ginContext *gin.Context, signInRequest SignInRequest, config *config.Config) {
	spans, ctx := telemetry.StartRestSpan(ctx, "signIn")
	defer spans.Finish()

	var err error

	saveErr := saveIP(ctx, ginContext, services, signInRequest.LoggedInEmail)
	if saveErr != nil {
		spans.TraceError(errors.Wrap(saveErr, "unable to save IP address"))
	}

	spans.LogObjectAsJson("request", signInRequest)

	firstName, lastName, err := validateRequestAtProvider(ctx, services, config, signInRequest)
	if err != nil {
		spans.TraceError(err)
		ginContext.JSON(http.StatusInternalServerError, gin.H{
			"result": fmt.Sprintf("unable to validate request at provider: %v", err.Error()),
		})
		return
	}

	isNewTenant := false
	isNewUser := false
	isPersonalEmail := false

	var currentTenant string
	var defaultTenant string
	var userId string
	var authId string
	var authUserId string

	signInRequest.LoggedInEmail = strings.ToLower(signInRequest.LoggedInEmail)

	if signInRequest.Tenant == "" {

		_, err = common_utils.ExecuteWriteInTransactionWithPostCommitActions(ctx, services.CommonServices.Neo4jRepositories.Neo4jDriver, services.CommonServices.Neo4jRepositories.Database, nil, func(txWithPostCommit *common_utils.TxWithPostCommit) (any, error) {
			// authentication
			authIdAndProviderNode, err := services.CommonServices.Neo4jRepositories.AuthenticationReadRepository.GetByAuthIdAndProvider(ctx, signInRequest.LoggedInEmail, signInRequest.Provider)
			if err != nil {
				return nil, err
			}

			// auth with the same provider found
			if authIdAndProviderNode != nil {
				authProps := common_utils.GetPropsFromNode(*authIdAndProviderNode)
				authId = common_utils.GetStringPropOrEmpty(authProps, "id")

				authUserNode, err := services.Repositories.Neo4jRepositories.AuthenticationReadRepository.GetAuthUser(ctx, authId)
				if err != nil {
					return nil, err
				}

				authUserProps := common_utils.GetPropsFromNode(*authUserNode)
				authUserId = common_utils.GetStringPropOrEmpty(authUserProps, "id")
				defaultTenant = common_utils.GetStringPropOrEmpty(authUserProps, "defaultTenant")
				currentTenant = common_utils.GetStringPropOrEmpty(authUserProps, "currentTenant")

				ctx = common.WithCustomContext(ctx, &common.CustomContext{
					Tenant: currentTenant,
				})
			} else {
				// auth with different provider found. link back to the same user and add new auth
				authNodes, err := services.CommonServices.Neo4jRepositories.AuthenticationReadRepository.GetByAuthId(ctx, signInRequest.LoggedInEmail)
				if err != nil {
					return nil, err
				}

				if authNodes != nil && len(authNodes) > 0 {
					authProps := common_utils.GetPropsFromNode(*authNodes[0])
					authId = common_utils.GetStringPropOrEmpty(authProps, "id")

					authUserNode, err := services.Repositories.Neo4jRepositories.AuthenticationReadRepository.GetAuthUser(ctx, authId)
					if err != nil {
						return nil, err
					}

					authUserProps := common_utils.GetPropsFromNode(*authUserNode)
					authUserId = common_utils.GetStringPropOrEmpty(authUserProps, "id")
					defaultTenant = common_utils.GetStringPropOrEmpty(authUserProps, "defaultTenant")
					currentTenant = common_utils.GetStringPropOrEmpty(authUserProps, "currentTenant")

					ctx = common.WithCustomContext(ctx, &common.CustomContext{
						Tenant: currentTenant,
					})

					authId, err = services.Repositories.Neo4jRepositories.AuthenticationWriteRepository.CreateAuthentication(ctx, *txWithPostCommit.Tx, neoEntity.AuthenticationEntity{
						AuthId:     signInRequest.LoggedInEmail,
						Provider:   signInRequest.Provider,
						IdentityId: signInRequest.OAuthToken.ProviderAccountId,
					})
					if err != nil {
						return nil, err
					}

					err = services.Repositories.Neo4jRepositories.AuthenticationWriteRepository.LinkAuthenticationWithAuthenticationUser(ctx, *txWithPostCommit.Tx, authId, authUserId)
					if err != nil {
						return nil, err
					}
				}
			}

			if authId != "" && authUserId == "" {
				spans.LogKV("authId", authId)
				spans.LogKV("authUserId", authUserId)
				return nil, fmt.Errorf("authId found but authUserId not found")
			}

			// auth doesn't exist at all
			if authId == "" {
				isNewUser = true
				// create auth + user
				authId, err = services.Repositories.Neo4jRepositories.AuthenticationWriteRepository.CreateAuthentication(ctx, *txWithPostCommit.Tx, neoEntity.AuthenticationEntity{
					AuthId:     signInRequest.LoggedInEmail,
					Provider:   signInRequest.Provider,
					IdentityId: signInRequest.OAuthToken.ProviderAccountId,
				})
				if err != nil {
					return nil, err
				}

				authUserId, err = services.Repositories.Neo4jRepositories.AuthenticationWriteRepository.CreateAuthenticationUser(ctx, *txWithPostCommit.Tx, authId, neoEntity.AuthenticationUserEntity{
					FirstName: firstName,
					LastName:  lastName,
				})
				if err != nil {
					return nil, err
				}
			}

			spans.LogKV("authId", authId)
			spans.LogKV("authUserId", authUserId)

			// tenant
			tenants, err := services.Repositories.Neo4jRepositories.AuthenticationReadRepository.GetTenants(ctx, authUserId)
			if err != nil {
				return nil, err
			}

			if tenants == nil || len(tenants) == 0 {
				isPersonalEmail = services.CommonServices.EmailService.IsPersonalEmailProvider(ctx, signInRequest.LoggedInEmail)
				spans.LogKV("isPersonalEmail", isPersonalEmail)

				domain := common_utils.ExtractDomain(signInRequest.LoggedInEmail)
				spans.LogKV("domainExtractedFromEmail", domain)

				if !isPersonalEmail {
					tenantWithWorkspace, err := services.Repositories.Neo4jRepositories.TenantReadRepository.GetTenantForWorkspace(ctx, domain)
					if err != nil {
						spans.TraceError(err)
						return nil, err
					}

					if tenantWithWorkspace != nil {
						tenantEntity := mapper.MapDbNodeToTenantEntity(tenantWithWorkspace)

						isNewTenant = false
						currentTenant = tenantEntity.Name
						defaultTenant = tenantEntity.Name
					} else {
						isNewTenant = true
					}
				} else {
					isNewTenant = true
				}

				spans.LogKV("isNewTenant", isNewTenant)

				if isNewTenant {
					tenantStr := ""
					if isPersonalEmail {
						// DO NOT ALLOW REGISTERING TENANT WITH PERSONAL EMAIL
						err = errors.New("personal email domain detected")
						spans.TraceError(err)
						return nil, &coserrors.CustomError{
							Code:      "PERSONAL_EMAIL_DOMAIN",
							Message:   "Please use a business email address to register a tenant",
							ErrorText: "Personal email domains are not allowed for tenant registration",
						}
						// tenantStr = utils.GenerateName() // uncomment this when we want to generate a random tenant name for personal email domains
					} else {
						tenantStr = utils.Sanitize(domain)
					}

					spans.LogKV("newTenantCreationWith", tenantStr)

					primaryDomain := domain
					if isPersonalEmail {
						primaryDomain = ""
					}

					tenantEntity, err := services.CommonServices.TenantService.Merge(ctx, *txWithPostCommit.Tx, neoEntity.TenantEntity{
						Name:      tenantStr,
						CreatedBy: signInRequest.LoggedInEmail,
					},
						primaryDomain)
					if err != nil {
						return nil, err
					}

					currentTenant = tenantEntity.Name
					defaultTenant = tenantEntity.Name

					ctx = common.WithCustomContext(ctx, &common.CustomContext{
						Tenant: currentTenant,
					})

					if !isPersonalEmail {
						_, err = services.CommonServices.WorkspaceService.MergeToTenant(ctx, txWithPostCommit.Tx, neoEntity.WorkspaceEntity{
							Name:     domain,
							Provider: signInRequest.Provider,
						}, tenantEntity.Name)
						if err != nil {
							return nil, err
						}
					}
				}

				ctx = common.WithCustomContext(ctx, &common.CustomContext{
					Tenant: currentTenant,
				})

				err = services.CommonServices.Neo4jRepositories.AuthenticationWriteRepository.LinkAuthenticationUserWithTenant(ctx, txWithPostCommit.Tx, authUserId, defaultTenant)
				if err != nil {
					return nil, err
				}

				err = services.CommonServices.Neo4jRepositories.AuthenticationWriteRepository.SetDefaultTenant(ctx, *txWithPostCommit.Tx, authUserId, defaultTenant)
				if err != nil {
					return nil, err
				}

				err = services.CommonServices.Neo4jRepositories.AuthenticationWriteRepository.SetCurrentTenant(ctx, txWithPostCommit.Tx, authUserId, defaultTenant)
				if err != nil {
					return nil, err
				}
			}

			// user in tenant
			userInTenantNode, err := services.Repositories.Neo4jRepositories.UserReadRepository.GetAuthenticatedUserInTenant(ctx, authUserId, signInRequest.LoggedInEmail)
			if err != nil {
				spans.TraceError(err)
				return nil, err
			}
			if userInTenantNode != nil {
				userId = mapper.MapDbNodeToUserEntity(userInTenantNode).Id
				spans.LogKV("user.found", true)
			} else {
				spans.LogKV("user.found", false)
				userId, err = services.CommonServices.AuthenticationService.CreateUserInTenant(ctx, txWithPostCommit, defaultTenant, false, authUserId, signInRequest.LoggedInEmail, firstName, lastName)
				if err != nil {
					return nil, err
				}
				isNewUser = true
			}

			return nil, nil
		})
		if err != nil {
			spans.TraceError(err)
			if customErr, ok := err.(*coserrors.CustomError); ok {
				ginContext.JSON(http.StatusBadRequest, gin.H{
					"code":    customErr.Code,
					"error":   customErr.ErrorText,
					"message": customErr.Message,
				})
			} else {
				ginContext.JSON(http.StatusInternalServerError, gin.H{
					"result": fmt.Sprintf("unable to create auth: %v", err.Error()),
				})
			}
			return
		}

		ctx = common.WithCustomContext(ctx, &common.CustomContext{
			Tenant:    currentTenant,
			UserEmail: signInRequest.LoggedInEmail,
		})

		_, err = initializeUserInTenant(ctx, services, userId)
		if err != nil {
			spans.TraceError(err)
			ginContext.JSON(http.StatusInternalServerError, gin.H{
				"result": fmt.Sprintf("unable to initialize user: %v", err.Error()),
			})
			return
		}

		if isNewTenant {
			// Post tenant creation actions
			err = registerNewTenantAsLeadInProviderTenant(ctx, config, services, signInRequest.LoggedInEmail)
			if err != nil {
				spans.TraceError(err)
			}

			if !isPersonalEmail {
				err = services.CommonServices.RegistrationService.ProvideAccessToPlatformOwners(ctx, defaultTenant)
				if err != nil {
					spans.TraceError(err)
				}

				err = services.CommonServices.RegistrationService.InitialTenantSetup(ctx, signInRequest.LoggedInEmail)
				if err != nil {
					spans.TraceError(err)
				}

				err = sendNewTenantRegistrationSlackNotification(ctx, config.Common.External.SlackConfig, defaultTenant, signInRequest.LoggedInEmail)
				if err != nil {
					spans.TraceError(err)
				}
			}
		} else if isNewUser {
			err = sendNewUserRegistrationSlackNotification(ctx, config.Common.External.SlackConfig, currentTenant, signInRequest.LoggedInEmail)
			if err != nil {
				spans.TraceError(err)
			}
		}

	} else {
		currentTenant = signInRequest.Tenant
		defaultTenant = signInRequest.Tenant

		ctx = common.WithCustomContext(ctx, &common.CustomContext{
			Tenant: currentTenant,
		})

		userInTenantNode, err := services.Repositories.Neo4jRepositories.UserReadRepository.GetFirstUserByEmail(ctx, currentTenant, signInRequest.LoggedInEmail)
		if err != nil {
			spans.TraceError(err)
			ginContext.JSON(http.StatusInternalServerError, gin.H{
				"result": fmt.Sprintf("unable to find user: %v", err.Error()),
			})
			return
		}

		if userInTenantNode != nil {
			userId = mapper.MapDbNodeToUserEntity(userInTenantNode).Id
		}

		ctx = common.WithCustomContext(ctx, &common.CustomContext{
			Tenant:    currentTenant,
			UserEmail: signInRequest.LoggedInEmail,
			UserId:    userId,
		})
	}

	// handle email token
	if isRequestEnablingOAuthSync(signInRequest) {
		err = handleOAuthTokenSync(ctx, services, signInRequest, defaultTenant, userId)
		if err != nil {
			spans.TraceError(err)
		}
	}

	spans.LogKV(
		"result.email", signInRequest.LoggedInEmail,
		"result.authUserId", authUserId,
		"result.currentTenant", currentTenant,
		"result.defaultTenant", defaultTenant,
		"result.userId", userId)

	apiKey, err := GetApiKeyForTenant(ctx, services, currentTenant)
	if err != nil {
		spans.TraceError(errors.Wrap(err, "GetApiKeyForTenant"))
	}

	ginContext.JSON(http.StatusOK, gin.H{
		"email":         signInRequest.LoggedInEmail,
		"authUserId":    authUserId,
		"currentTenant": currentTenant,
		"defaultTenant": defaultTenant,
		"userId":        userId,
		"apiKey":        apiKey,
	})
}

// TODO: to be deleted
func storeTenantCreatedEvent(ctx context.Context, services *cosapi_services.Services, tenant, domain string) error {
	span, ctx := telemetry.StartRestSpan(ctx, "Registration.storeTenantCreatedEvent")
	defer span.Finish()

	tenantCreatedEvent := &pb.TenantCreated{
		Timestamp: timestamppb.Now(),
		Tenant:    tenant,
		Domain:    domain,
	}

	payload, err := proto.Marshal(tenantCreatedEvent)
	if err != nil {
		span.TraceError(err)
		return fmt.Errorf("failed to marshal tenant created: %w", err)
	}

	// Create nats message with headers
	// TODO alexb move sending event to outbox processing
	msg := nats.NewMsg(commonenum.EventTenantCreated.String())
	msg.Data = payload
	msg.Header.Set(string(nats_common.NATS_HEADER_TENANT), tenant)

	_, err = services.CommonServices.NATSConnections.JS.PublishMsg(msg)
	if err != nil {
		span.TraceError(err)
		return fmt.Errorf("failed to publish tenant created event: %w", err)
	}

	return nil
}

func initializeUserInTenant(ctx context.Context, services *cosapi_services.Services, userId string) (*string, error) {
	spans, ctx := telemetry.StartRestSpan(ctx, "Registration.initializeUserInTenant")
	defer spans.Finish()

	tenant := common.GetTenantFromContext(ctx)

	err := addDefaultMissingRoles(ctx, services, tenant, userId)
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}

	defaultWorkSchedule := postgres_entity.UserWorkingSchedule{
		UserId:    userId,
		DayRange:  "Mon-Fri",
		StartHour: "09:00",
		EndHour:   "18:00",
	}

	workingSchedule, err := services.Repositories.PostgresRepositories.UserWorkingScheduleRepository.GetForUser(ctx, tenant, userId)
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}

	if len(workingSchedule) == 0 {
		err = services.Repositories.PostgresRepositories.UserWorkingScheduleRepository.Store(ctx, tenant, &defaultWorkSchedule)
		if err != nil {
			spans.TraceError(err)
			return nil, err
		}
	}

	err = services.Repositories.Neo4jRepositories.UserWriteRepository.RegisterLogin(ctx, tenant, userId)
	if err != nil {
		spans.TraceError(err)
	}

	// TODO why is this needed?
	//err = services.CommonServices.Events.Publisher.PublishFanoutEvent(innerCtx, userId, model.USER, dto.UserLogin{LoginEmail: email, Provider: provider, IdentityId: providerAccountId})
	//if err != nil {
	//	tracing.TraceErr(span, err)
	//}

	return &userId, nil
}

func Revoke(s *cosapi_services.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		spans, ctx := telemetry.StartRestSpan(c.Request.Context(), "Registration.Revoke")
		defer spans.Finish()

		var revokeRequest RevokeRequest
		if err := c.BindJSON(&revokeRequest); err != nil {
			log.Printf("unable to parse json: %v", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"result": fmt.Sprintf("unable to parse json: %v", err.Error()),
			})
			return
		}
		log.Printf("parsed json: %v", revokeRequest)
		spans.LogObjectAsJson("revokeRequest", revokeRequest)

		workspaceProvider := ""
		if revokeRequest.MailboxProvider == model.MailboxProviderGoogleWorkspace.String() {
			workspaceProvider = commonenum.WorkspaceProviderGoogle.String()
		} else if revokeRequest.MailboxProvider == model.MailboxProviderOutlook.String() {
			workspaceProvider = commonenum.WorkspaceProviderAzure.String()
		}
		spans.LogKV("workspaceProvider", workspaceProvider)

		if workspaceProvider == "" {
			c.JSON(http.StatusBadRequest, gin.H{})
		}

		oauthToken, err := s.Repositories.PostgresRepositories.OAuthTokenRepository.GetByEmailAndProvider(ctx, revokeRequest.Tenant, workspaceProvider, revokeRequest.Email)
		if err != nil {
			spans.TraceError(err)
			c.JSON(http.StatusInternalServerError, gin.H{})
			return
		}

		if oauthToken != nil && oauthToken.RefreshToken != "" {
			// Handle revocation based on provider
			var revocationURL string
			var headers map[string]string
			vendor := commonenum.VendorNotSet
			switch workspaceProvider {
			case commonenum.WorkspaceProviderGoogle.String():
				// Decrypt the access token before revoking
				decryptedAccessToken, err := postgres_entity.DecryptToken(s.Cfg.Common.Infrastructure.GoogleOAuthConfig.EncryptionKey, oauthToken.AccessToken)
				if err != nil {
					spans.TraceError(err)
					c.JSON(http.StatusInternalServerError, gin.H{})
					return
				}
				revocationURL = fmt.Sprintf("https://accounts.google.com/o/oauth2/revoke?token=%s", decryptedAccessToken)
				headers = map[string]string{
					"Content-Type": "application/x-www-form-urlencoded",
				}
				vendor = commonenum.VendorGoogle
			case commonenum.WorkspaceProviderAzure.String():
				revocationURL = "https://graph.microsoft.com/v1.0/me/revokeSignInSessions"
				headers = map[string]string{
					"Authorization": fmt.Sprintf("Bearer %s", oauthToken.AccessToken),
				}
				vendor = commonenum.VendorMicrosoft
			}
			spans.LogKV("revocationURL", revocationURL)

			if revocationURL != "" {
				req, err := http.NewRequest("POST", revocationURL, nil)
				if err != nil {
					spans.TraceError(err)
					c.JSON(http.StatusInternalServerError, gin.H{})
					return
				}

				for key, value := range headers {
					req.Header.Add(key, value)
				}

				clientTimeout := 30 * time.Second
				httpClient := clients.NewLoggingClient(s.CommonServices.WarehouseRepositories.APICallLogRepository, vendor, &clientTimeout)
				resp, err := httpClient.Do(req)
				if err != nil {
					spans.TraceError(err)
					c.JSON(http.StatusInternalServerError, gin.H{})
					return
				}
				defer resp.Body.Close()

				body, _ := io.ReadAll(resp.Body)
				spans.LogKV("response.body", string(body))

				// For Google, if token is already revoked (invalid_token error), we can proceed
				if workspaceProvider == commonenum.WorkspaceProviderGoogle.String() {
					var errorResponse struct {
						Error            string `json:"error"`
						ErrorDescription string `json:"error_description"`
					}
					if err := json.Unmarshal(body, &errorResponse); err == nil {
						if errorResponse.Error == "invalid_token" {
							// Token is already revoked, we can proceed
							spans.LogKV("token_status", "already_revoked")
						} else if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
							// Other errors should be treated as failures
							spans.TraceError(fmt.Errorf("revocation failed, status code: %d", resp.StatusCode))
							c.JSON(http.StatusInternalServerError, gin.H{})
							return
						}
					}
				}

				// For Google, also revoke the refresh token
				if workspaceProvider == commonenum.WorkspaceProviderGoogle.String() {
					// Decrypt the refresh token before revoking
					decryptedRefreshToken, err := postgres_entity.DecryptToken(s.Cfg.Common.Infrastructure.GoogleOAuthConfig.EncryptionKey, oauthToken.RefreshToken)
					if err != nil {
						spans.TraceError(err)
						c.JSON(http.StatusInternalServerError, gin.H{})
						return
					}
					revocationURL = fmt.Sprintf("https://accounts.google.com/o/oauth2/revoke?token=%s", decryptedRefreshToken)
					req, err = http.NewRequest("POST", revocationURL, nil)
					if err != nil {
						spans.TraceError(err)
						c.JSON(http.StatusInternalServerError, gin.H{})
						return
					}

					for key, value := range headers {
						req.Header.Add(key, value)
					}

					httpClient = clients.NewLoggingClient(s.CommonServices.WarehouseRepositories.APICallLogRepository, commonenum.VendorGoogle, &clientTimeout)
					resp, err = httpClient.Do(req)
					if err != nil {
						spans.TraceError(err)
						c.JSON(http.StatusInternalServerError, gin.H{})
						return
					}
					defer resp.Body.Close()

					body, _ = io.ReadAll(resp.Body)
					spans.LogKV("refresh_token_response.body", string(body))

					// Handle already revoked refresh token case
					var errorResponse struct {
						Error            string `json:"error"`
						ErrorDescription string `json:"error_description"`
					}
					if err := json.Unmarshal(body, &errorResponse); err == nil {
						if errorResponse.Error == "invalid_token" {
							// Refresh token is already revoked, we can proceed
							spans.LogKV("refresh_token_status", "already_revoked")
						} else if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
							// Other errors should be treated as failures
							spans.TraceError(fmt.Errorf("refresh token revocation failed, status code: %d", resp.StatusCode))
							c.JSON(http.StatusInternalServerError, gin.H{})
							return
						}
					}
				}
			}
		}

		err = s.Repositories.PostgresRepositories.OAuthTokenRepository.DeleteByEmail(ctx, revokeRequest.Tenant, workspaceProvider, revokeRequest.Email)
		if err != nil {
			spans.TraceError(err)
			c.JSON(http.StatusInternalServerError, gin.H{})
			return
		}

		c.JSON(http.StatusOK, gin.H{})
	}
}

func validateRequestAtProvider(ctx context.Context, services *cosapi_services.Services, config *config.Config, signInRequest SignInRequest) (string, string, error) {
	spans, ctx := telemetry.StartRestSpan(ctx, "Registration.validateRequestAtProvider")
	defer spans.Finish()

	if signInRequest.Provider == commonenum.WorkspaceProviderMagicLink.String() {
		return "", "", nil
	} else if signInRequest.Provider == commonenum.WorkspaceProviderGoogle.String() {
		userInfo, err := getUserInfoFromGoogle(ctx, config, signInRequest)
		if err != nil {
			spans.TraceError(err)
			return "", "", err
		}

		return userInfo.GivenName, userInfo.FamilyName, nil
	} else if signInRequest.Provider == commonenum.WorkspaceProviderAzure.String() {
		// Create a GET request with the Authorization header.
		req, err := http.NewRequest("GET", "https://graph.microsoft.com/oidc/userinfo", nil)
		if err != nil {
			spans.TraceError(err)
			return "", "", err
		}

		req.Header.Set("Authorization", "Bearer "+signInRequest.OAuthToken.AccessToken)

		clientTimeout := 30 * time.Second
		httpClient := clients.NewLoggingClient(services.CommonServices.WarehouseRepositories.APICallLogRepository, commonenum.VendorMicrosoft, &clientTimeout)
		resp, err := httpClient.Do(req)
		if err != nil {
			spans.TraceError(err)
			return "", "", err
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			var data map[string]string
			json.NewDecoder(resp.Body).Decode(&data)

			firstName := data["given_name"]
			lastName := data["family_name"]
			return firstName, lastName, nil
		} else {
			spans.TraceError(err)
			return "", "", err
		}
	} else {
		spans.TraceError(fmt.Errorf("provider not supported"))
		return "", "", fmt.Errorf("provider not supported")
	}
}

func getUserInfoFromGoogle(c context.Context, config *config.Config, signInRequest SignInRequest) (*googleOauth.Userinfo, error) {
	spans, ctx := telemetry.StartRestSpan(c, "Registration.getUserInfoFromGoogle")
	defer spans.Finish()

	conf := &tokenOauth.Config{
		ClientID:     config.Common.Infrastructure.GoogleOAuthConfig.ClientId,
		ClientSecret: config.Common.Infrastructure.GoogleOAuthConfig.ClientSecret,
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
		spans.TraceError(err)
		return nil, err
	}
	userInfoService := googleOauth.NewUserinfoV2MeService(oauth2Service)

	userInfo, err := userInfoService.Get().Do()
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}

	return userInfo, nil
}

func addDefaultMissingRoles(c context.Context, services *cosapi_services.Services, tenant, userId string) error {
	spans, ctx := telemetry.StartRestSpan(c, "Registration.addDefaultMissingRoles")
	defer spans.Finish()

	userRoleFound := false
	ownerRoleFound := false
	impersonatedRoleFound := false

	userNode, err := services.Repositories.Neo4jRepositories.UserReadRepository.GetUserById(ctx, tenant, userId)
	if err != nil {
		spans.TraceError(err)
		return err
	}
	if userNode == nil {
		err = fmt.Errorf("user not found")
		spans.TraceError(err)
		return err
	}

	existingUser := mapper.MapDbNodeToUserEntity(userNode)

	if len(existingUser.Roles) > 0 {
		for _, role := range existingUser.Roles {
			if role == commonenum.UserRoleUser.String() {
				userRoleFound = true
			}
			if role == commonenum.UserRoleOwner.String() {
				ownerRoleFound = true
			}
			if role == commonenum.UserRoleImpersonated.String() {
				impersonatedRoleFound = true
			}
		}
	}

	if !impersonatedRoleFound {
		if !userRoleFound {
			err := services.Repositories.Neo4jRepositories.UserWriteRepository.AddRole(ctx, existingUser.Id, commonenum.UserRoleUser.String())
			if err != nil {
				spans.TraceError(err)
				return err
			}
		}
		if !ownerRoleFound {
			err := services.Repositories.Neo4jRepositories.UserWriteRepository.AddRole(ctx, existingUser.Id, commonenum.UserRoleOwner.String())
			if err != nil {
				spans.TraceError(err)
				return err
			}
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
	spans, ctx := telemetry.StartRestSpan(ctx, "Registration.registerNewTenantAsLeadInProviderTenant")
	defer spans.Finish()

	providerTenantCtx := common.WithCustomContext(ctx, &common.CustomContext{
		Tenant: config.App.AuthConfig.ProviderTenantName,
	})

	organizationId, contactId, err := createOrganizationAndContact(providerTenantCtx, services, config.App.AuthConfig.ProviderTenantName, registeredEmail, true, "Tenant Registration")
	if err != nil {
		spans.TraceError(err)
		return err
	}

	spans.LogKV("providerOrganizationId", *organizationId)
	spans.LogKV("providerContactId", *contactId)

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
	//        <p>To be honest, our self-service onboarding kinda sucks right now as we're still building it out. I'd love to get you setup and using the tool, would you be open to spending 10 mins with me to help you get things configured?</p>
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
	spans, ctx := telemetry.StartRestSpan(ctx, "RegistrationService.CreateOrganizationAndContact")
	defer spans.Finish()

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
			spans.TraceError(err)
			return nil, nil, err
		}

		if organizationByDomain == nil {
			organizationId, err = services.CommonServices.OrganizationService.Save(ctx, nil, nil, data_fields.OrganizationFields{
				Domains:    []string{domain},
				Name:       common_utils.StringPtr(domain),
				Stage:      common_utils.ToPtr(commonenum.Target),
				LeadSource: common_utils.StringPtr(leadSource),
			})
			if err != nil {
				spans.TraceError(err)
				return nil, nil, err
			}
			if organizationId == "" {
				err := errors.New("organization id empty")
				spans.TraceError(err)
				return nil, nil, err
			}
		} else {
			organizationId = mapper.MapDbNodeToOrganizationEntity(organizationByDomain).ID
		}

		if organizationId == "" {
			spans.TraceError(err)
			return nil, nil, err
		}
		spans.LogKV("result.organizationId", organizationId)

		contactNode, err := services.Repositories.Neo4jRepositories.ContactReadRepository.GetContactInOrganizationByEmail(ctx, tenant, organizationId, email)
		if err != nil {
			spans.TraceError(err)
			return nil, nil, err
		}

		if contactNode == nil {
			contactId, err = services.CommonServices.ContactService.CreateContactByEmail(ctx, nil, email)
			if err != nil {
				spans.TraceError(err)
				return nil, nil, err
			}

			err = services.CommonServices.ContactService.LinkContactWithOrganization(ctx, nil, contactId, organizationId, "", "",
				neoEntity.DataSourceOpenline.String(), false, nil, nil)
			if err != nil {
				spans.TraceError(err)
				return nil, nil, err
			}
		} else {
			contactId = mapper.MapDbNodeToContactEntity(contactNode).Id
		}

		if contactId == "" {
			spans.TraceError(errors.New("contact id empty"))
			return nil, nil, errors.New("contact id empty")
		}
		spans.LogKV("result.contactId", contactId)
	}

	return &organizationId, &contactId, nil
}

func saveIP(ctx context.Context, c *gin.Context, s *cosapi_services.Services, email string) error {
	spans, ctx := telemetry.StartRestSpan(ctx, "registration.saveIP")
	defer spans.Finish()

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
		spans.TraceError(err)
		spans.LogKV("email", email)
	}

	if validEmail.IsFreeAccount || validEmail.IsRoleAccount || validEmail.IsSystemGenerated {
		return nil
	}

	website := fmt.Sprintf("https://%s", validEmail.Domain)

	details := postgres_entity.EnrichDetailsTracking{
		IP:             clientIP,
		CompanyDomain:  &validEmail.Domain,
		CompanyWebsite: &website,
		SourceEmail:    &validEmail.CleanEmail,
	}
	spans.LogObjectAsJson("ipToEmailDetails", details)

	err := s.Repositories.PostgresRepositories.EnrichDetailsTrackingRepository.Save(ctx, details)
	if err != nil {
		spans.TraceError(err)
		return err
	}

	return nil
}

func handleGoogleOAuthToken(ctx context.Context, services *cosapi_services.Services, signInRequest SignInRequest, defaultTenant string, userId string) error {
	spans, ctx := telemetry.StartRestSpan(ctx, "handleGoogleOAuthToken")
	defer spans.Finish()

	var err error

	oauthToken, _ := services.Repositories.PostgresRepositories.OAuthTokenRepository.GetByEmailAndProvider(ctx, defaultTenant, commonenum.SourceGmail.String(), signInRequest.OAuthTokenForEmail)
	newOAuthToken := false
	if oauthToken == nil {
		oauthToken = &postgres_entity.OAuthTokenEntity{}
		newOAuthToken = true
	}
	oauthToken.Provider = commonenum.SourceGmail.String()
	oauthToken.TenantName = defaultTenant
	oauthToken.PlayerIdentityId = signInRequest.OAuthToken.ProviderAccountId
	oauthToken.EmailAddress = signInRequest.OAuthTokenForEmail
	oauthToken.UserId = userId

	oauthToken.AccessToken, err = postgres_entity.EncryptToken(services.Cfg.Common.Infrastructure.GoogleOAuthConfig.EncryptionKey, signInRequest.OAuthToken.AccessToken)
	if err != nil {
		spans.TraceError(err)
	}

	oauthToken.RefreshToken, err = postgres_entity.EncryptToken(services.Cfg.Common.Infrastructure.GoogleOAuthConfig.EncryptionKey, signInRequest.OAuthToken.RefreshToken)
	if err != nil {
		spans.TraceError(err)
	}

	oauthToken.IdToken = signInRequest.OAuthToken.IdToken
	oauthToken.ExpiresAt = signInRequest.OAuthToken.ExpiresAt
	oauthToken.Scope = signInRequest.OAuthToken.Scope
	oauthToken.NeedsManualRefresh = false
	if !newOAuthToken {
		if isRequestEnablingGmailSync(signInRequest) {
			oauthToken.GmailSyncEnabled = true
		}
		if isRequestEnablingGoogleCalendarSync(signInRequest) {
			oauthToken.GoogleCalendarSyncEnabled = true
		}
	}
	_, err = services.Repositories.PostgresRepositories.OAuthTokenRepository.Save(ctx, *oauthToken)
	if err != nil {
		spans.TraceError(err)
	}

	return nil
}

func handleAzureOAuthToken(ctx context.Context, services *cosapi_services.Services, signInRequest SignInRequest, defaultTenant string, userId string) error {
	spans, ctx := telemetry.StartRestSpan(ctx, "handleAzureOAuthToken")
	defer spans.Finish()

	oauthToken, _ := services.Repositories.PostgresRepositories.OAuthTokenRepository.GetByEmailAndProvider(ctx, defaultTenant, commonenum.SourceOutlook.String(), signInRequest.OAuthTokenForEmail)
	if oauthToken == nil {
		oauthToken = &postgres_entity.OAuthTokenEntity{}
	}
	oauthToken.Provider = commonenum.SourceOutlook.String()
	oauthToken.TenantName = defaultTenant
	oauthToken.PlayerIdentityId = signInRequest.OAuthToken.ProviderAccountId
	oauthToken.EmailAddress = signInRequest.OAuthTokenForEmail
	oauthToken.UserId = userId
	oauthToken.AccessToken = signInRequest.OAuthToken.AccessToken
	oauthToken.RefreshToken = signInRequest.OAuthToken.RefreshToken
	oauthToken.IdToken = signInRequest.OAuthToken.IdToken
	oauthToken.ExpiresAt = signInRequest.OAuthToken.ExpiresAt
	oauthToken.Scope = signInRequest.OAuthToken.Scope
	oauthToken.NeedsManualRefresh = false
	_, err := services.Repositories.PostgresRepositories.OAuthTokenRepository.Save(ctx, *oauthToken)
	if err != nil {
		spans.TraceError(err)
	}

	return nil
}

func handleOAuthTokenSync(ctx context.Context, services *cosapi_services.Services, signInRequest SignInRequest, defaultTenant string, userId string) error {
	spans, ctx := telemetry.StartRestSpan(ctx, "handleOAuthTokenSync")
	defer spans.Finish()

	var err error

	switch signInRequest.Provider {
	case commonenum.WorkspaceProviderGoogle.String():
		err = handleGoogleOAuthToken(ctx, services, signInRequest, defaultTenant, userId)
	case commonenum.WorkspaceProviderAzure.String():
		err = handleAzureOAuthToken(ctx, services, signInRequest, defaultTenant, userId)
	case commonenum.WorkspaceProviderMagicLink.String():
		// No action needed for magic link
	default:
		spans.TraceError(fmt.Errorf("unsupported provider: %s", signInRequest.Provider))
	}

	return err
}

func sendNewTenantRegistrationSlackNotification(ctx context.Context, cfg commonconfig.SlackConfig, tenant, userEmail string) error {
	spans, ctx := telemetry.StartServiceSpan(ctx, "sendNewTenantRegistrationSlackNotification")
	defer spans.Finish()

	// Create a struct to hold the JSON data
	type SlackMessage struct {
		Text string `json:"text"`
	}
	message := SlackMessage{Text: fmt.Sprintf("🏢 New tenant registration!\n• Tenant: %s\n• User: %s", tenant, userEmail)}
	// Convert struct to JSON
	jsonData, err := json.Marshal(message)
	if err != nil {
		spans.TraceError(err)
		return err
	}

	// Send POST request
	resp, err := http.Post(cfg.NotifyNewTenantRegisteredHook, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error sending request:", err)
		return err
	}
	defer resp.Body.Close()

	spans.LogKV("result.status", resp.Status)

	return nil
}

func sendNewUserRegistrationSlackNotification(ctx context.Context, cfg commonconfig.SlackConfig, tenant, userEmail string) error {
	spans, ctx := telemetry.StartServiceSpan(ctx, "sendNewUserRegistrationSlackNotification")
	defer spans.Finish()

	// Create a struct to hold the JSON data
	type SlackMessage struct {
		Text string `json:"text"`
	}
	message := SlackMessage{Text: fmt.Sprintf("👤 New user first sign-in!\n• User: %s\n• Tenant: %s", userEmail, tenant)}
	// Convert struct to JSON
	jsonData, err := json.Marshal(message)
	if err != nil {
		spans.TraceError(err)
		return err
	}

	// Send POST request
	resp, err := http.Post(cfg.NotifyNewTenantRegisteredHook, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error sending request:", err)
		return err
	}
	defer resp.Body.Close()

	spans.LogKV("result.status", resp.Status)

	return nil
}
