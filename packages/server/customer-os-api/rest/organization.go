package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	enummapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/mapper/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/service"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/grpc_client"
	commonTracing "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	commonpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/common"
	organizationpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/organization"
	socialpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/social"
	tracingLog "github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"net/http"
)

type CreateOrganizationRequest struct {
	Name         string `json:"name"`
	CustomId     string `json:"customId"`
	Website      string `json:"website"`
	LinkedinUrl  string `json:"linkedinUrl"`
	LeadSource   string `json:"leadSource"`
	Relationship string `json:"relationship"`
	IcpFit       bool   `json:"icpFit"`
}

func CreateOrganization(services *service.Services, grpcClients *grpc_client.Clients) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "CreateOrganization", c.Request.Header)
		defer span.Finish()

		tenant := common.GetTenantFromContext(ctx)
		// if tenant missing return auth error
		if tenant == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "API key invalid or expired"})
			span.LogFields(tracingLog.String("result", "Missing tenant in context"))
			return
		}

		tracing.SetDefaultRestSpanTags(ctx, span)

		request := CreateOrganizationRequest{}
		// Bind the JSON request body to the struct
		if err := c.ShouldBindJSON(&request); err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "Invalid request body"))
			services.Log.Error(ctx, "Invalid request body", err)
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid request body"})
			return
		}
		tracing.LogObjectAsJson(span, "request", request)

		// at least 1 input field is required
		if request.Name == "" && request.CustomId == "" && request.Website == "" && request.LinkedinUrl == "" {
			span.LogFields(tracingLog.String("result", "Missing organization input fields"))
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Missing organization input fields"})
			return
		}

		// step 1 validate no organization exist with given domain
		websiteDomain := ""
		if request.Website != "" {
			websiteDomain = utils.ExtractDomain(request.Website)
			if websiteDomain != "" {
				orgDbNode, err := services.Repositories.Neo4jRepositories.OrganizationReadRepository.GetOrganizationByDomain(ctx, tenant, websiteDomain)
				if err != nil {
					tracing.TraceErr(span, err)
					c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to check organization domain"})
					return
				}
				if orgDbNode != nil {
					orgId := utils.GetStringPropOrEmpty(utils.GetPropsFromNode(*orgDbNode), "id")
					span.LogFields(tracingLog.String("result", "Organization already exists with given domain"))
					c.JSON(http.StatusConflict, gin.H{"status": "error", "message": "Organization already exists with given domain", "id": orgId})
					return
				}
			}
		}

		// step 2 validate no organization exist with given custom id
		if request.CustomId != "" {
			orgDbNode, err := services.Repositories.Neo4jRepositories.OrganizationReadRepository.GetOrganizationByReferenceId(ctx, tenant, request.CustomId)
			if err != nil {
				tracing.TraceErr(span, err)
				c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to check organization custom id"})
				return
			}
			if orgDbNode != nil {
				orgId := utils.GetStringPropOrEmpty(utils.GetPropsFromNode(*orgDbNode), "id")
				span.LogFields(tracingLog.String("result", "Organization already exists with given custom id"))
				c.JSON(http.StatusConflict, gin.H{"status": "error", "message": "Organization already exists with given custom id", "id": orgId})
				return
			}
		}

		// step 3 validate no organization exist with given linkedin url
		if request.LinkedinUrl != "" {
			orgDbNode, err := services.Repositories.Neo4jRepositories.OrganizationReadRepository.GetOrganizationBySocialUrl(ctx, tenant, request.LinkedinUrl)
			if err != nil {
				tracing.TraceErr(span, err)
				c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to check organization linkedin url"})
				return
			}
			if orgDbNode != nil {
				orgId := utils.GetStringPropOrEmpty(utils.GetPropsFromNode(*orgDbNode), "id")
				span.LogFields(tracingLog.String("result", "Organization already exists with given linkedin url"))
				c.JSON(http.StatusConflict, gin.H{"status": "error", "message": "Organization already exists with given linkedin url", "id": orgId})
				return
			}
		}

		// step 4 reserve org id
		newOrgId, err := services.Repositories.Neo4jRepositories.OrganizationWriteRepository.ReserveOrganizationId(ctx, tenant, "")
		if err != nil {
			tracing.TraceErr(span, err)
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to create organization"})
			return
		}
		span.SetTag(tracing.SpanTagEntityId, newOrgId)

		// step 5 create organization
		orgName := request.Name
		if orgName == "" {
			orgName = websiteDomain
		}
		upsertOrganizationRequest := organizationpb.UpsertOrganizationGrpcRequest{
			Tenant:         common.GetTenantFromContext(ctx),
			LoggedInUserId: common.GetUserIdFromContext(ctx),
			Name:           orgName,
			ReferenceId:    request.CustomId,
			Website:        request.Website,
			SourceFields: &commonpb.SourceFields{
				Source:    string(neo4jentity.DataSourceOpenline),
				AppSource: constants.AppSourceCustomerOsApiRest,
			},
			LeadSource: request.LeadSource,
			IcpFit:     request.IcpFit,
		}
		relationship := model.OrganizationRelationshipProspect
		if request.Relationship != "" && model.OrganizationRelationship(request.Relationship).IsValid() {
			relationship = model.OrganizationRelationship(request.Relationship)
		}
		upsertOrganizationRequest.Relationship = enummapper.MapRelationshipFromModel(relationship).String()

		if upsertOrganizationRequest.Relationship == enummapper.MapRelationshipFromModel(model.OrganizationRelationshipCustomer).String() {
			upsertOrganizationRequest.Stage = enummapper.MapStageFromModel(model.OrganizationStageOnboarding).String()
		} else if upsertOrganizationRequest.Relationship == enummapper.MapRelationshipFromModel(model.OrganizationRelationshipProspect).String() {
			upsertOrganizationRequest.Stage = enummapper.MapStageFromModel(model.OrganizationStageLead).String()
		} else if upsertOrganizationRequest.Relationship == enummapper.MapRelationshipFromModel(model.OrganizationRelationshipNotAFit).String() {
			upsertOrganizationRequest.Stage = enummapper.MapStageFromModel(model.OrganizationStageUnqualified).String()
		} else if upsertOrganizationRequest.Relationship == enummapper.MapRelationshipFromModel(model.OrganizationRelationshipFormerCustomer).String() {
			upsertOrganizationRequest.Stage = enummapper.MapStageFromModel(model.OrganizationStageTarget).String()
		}

		ctx = commonTracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
		_, err = utils.CallEventsPlatformGRPCWithRetry[*organizationpb.OrganizationIdGrpcResponse](func() (*organizationpb.OrganizationIdGrpcResponse, error) {
			return grpcClients.OrganizationClient.UpsertOrganization(ctx, &upsertOrganizationRequest)
		})
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "Failed to create organization"))
			services.Log.Error(ctx, "Failed to create organization", err)
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to create organization"})
			return
		}

		if request.LinkedinUrl != "" {
			_, err = utils.CallEventsPlatformGRPCWithRetry[*socialpb.SocialIdGrpcResponse](func() (*socialpb.SocialIdGrpcResponse, error) {
				return grpcClients.OrganizationClient.AddSocial(ctx, &organizationpb.AddSocialGrpcRequest{
					Tenant:         common.GetTenantFromContext(ctx),
					LoggedInUserId: common.GetUserIdFromContext(ctx),
					OrganizationId: newOrgId,
					Url:            request.LinkedinUrl,
					SourceFields: &commonpb.SourceFields{
						Source:    string(neo4jentity.DataSourceOpenline),
						AppSource: constants.AppSourceCustomerOsApiRest,
					},
				})
			})
			if err != nil {
				tracing.TraceErr(span, errors.Wrap(err, "Failed to add linkedin url"))
				services.Log.Error(ctx, "Failed to add linkedin url", err)
				// partial saving of data
				c.JSON(http.StatusPartialContent, gin.H{"status": "partial_success", "message": "Failed to add linkedin url", "id": newOrgId})
			}
		}

		// Prepare and send the response
		span.LogFields(tracingLog.String("result", "Organization created successfully"))
		c.JSON(http.StatusCreated, gin.H{"status": "success", "message": "Organization created successfully", "id": newOrgId})
	}
}
