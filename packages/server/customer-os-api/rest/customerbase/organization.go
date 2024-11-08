package customerbase

import (
	"github.com/gin-gonic/gin"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	enummapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/mapper/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/service"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	commonservice "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	commonpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/common"
	organizationpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/organization"
	socialpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/social"
	tracingLog "github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"net/http"
	"strings"
)

// @Summary Create a new organization
// @Description Creates an organization in the system if it doesn't already exist based on website, custom ID, or LinkedIn URL
// @Tags CustomerBASE API
// @Accept  json
// @Produce  json
// @Param   body   body    CreateOrganizationRequest  true  "Organization creation payload"
// @Success 201 {object} CreateOrganizationResponse "Organization created successfully"
// @Success 206 {object} CreateOrganizationResponse "Partial success - failed to add linkedin url"
// @Failure 400  "Invalid request body or missing input fields"
// @Failure 401  "Unauthorized access"
// @Failure 409  "Conflict - organization already exists"
// @Failure 500  "Failed to create organization"
// @Router /customerbase/v1/organizations [post]
// @Security ApiKeyAuth
func CreateOrganization(services *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "CreateOrganization", c.Request.Header)
		defer span.Finish()
		tracing.TagComponentRest(span)
		tracing.TagTenant(span, common.GetTenantFromContext(ctx))

		tenant := common.GetTenantFromContext(ctx)
		// if tenant missing return auth error
		if tenant == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "API key invalid or expired"})
			span.LogFields(tracingLog.String("result", "Missing tenant in context"))
			return
		}

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
		websiteDomain, _ := services.CommonServices.DomainService.GetPrimaryDomainForOrganizationWebsite(ctx, request.Website)
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
			Id:             newOrgId,
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

		ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
		_, err = utils.CallEventsPlatformGRPCWithRetry[*organizationpb.OrganizationIdGrpcResponse](func() (*organizationpb.OrganizationIdGrpcResponse, error) {
			return services.CommonServices.GrpcClients.OrganizationClient.UpsertOrganization(ctx, &upsertOrganizationRequest)
		})
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "Failed to create organization"))
			services.Log.Error(ctx, "Failed to create organization", err)
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to create organization"})
			return
		}

		if request.LinkedinUrl != "" {
			_, err = utils.CallEventsPlatformGRPCWithRetry[*socialpb.SocialIdGrpcResponse](func() (*socialpb.SocialIdGrpcResponse, error) {
				return services.CommonServices.GrpcClients.OrganizationClient.AddSocial(ctx, &organizationpb.AddSocialGrpcRequest{
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
				c.JSON(http.StatusPartialContent,
					CreateOrganizationResponse{
						Status:         "partial_success",
						Message:        "Failed to add linkedin url",
						ID:             newOrgId,
						PartialSuccess: true,
					})
			}
		}

		// Prepare and send the response
		span.LogFields(tracingLog.String("result", "Organization created successfully"))
		c.JSON(http.StatusCreated,
			CreateOrganizationResponse{
				Status:  "success",
				Message: "Organization created successfully",
				ID:      newOrgId,
			})
	}
}

// @Summary Set primary link ID for an organization
// @Description Sets or replaces the primary ID for an organization linked to a specific external system (e.g., Stripe, HubSpot)
// @Tags CustomerBASE API
// @Accept  json
// @Produce  json
// @Param   id             path     string  true  "Organization ID or Organization COS ID"
// @Param   externalSystem path     string  true  "External system name, supported values (stripe, hubspot, close)"
// @Param   body           body     SetPrimaryExternalSystemIdRequest true  "Request payload to set primary external system ID"
// @Success 200 {object} SetPrimaryExternalSystemIdResponse "Primary ID set successfully"
// @Failure 400  "Invalid request body or input"
// @Failure 401  "Unauthorized access"
// @Failure 404  "Organization or external system not found"
// @Failure 500  "Internal server error"
// @Router /customerbase/v1/organizations/{id}/links/{externalSystem}/primary [put]
// @Security ApiKeyAuth
func SetPrimaryExternalSystemId(services *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "SetPrimaryExternalSystemId", c.Request.Header)
		defer span.Finish()
		tracing.TagComponentRest(span)
		tracing.TagTenant(span, common.GetTenantFromContext(ctx))

		tenant := common.GetTenantFromContext(ctx)
		if tenant == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "API key invalid or expired"})
			span.LogFields(tracingLog.String("result", "Missing tenant in context"))
			return
		}

		orgId := c.Param("id")
		externalSystem := strings.ToLower(c.Param("externalSystem"))
		request := SetPrimaryExternalSystemIdRequest{}

		if err := c.ShouldBindJSON(&request); err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "Invalid request body"))
			services.Log.Error(ctx, "Invalid request body", err)
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid request body"})
			return
		}
		tracing.LogObjectAsJson(span, "request.body", request)

		// Validate if the organization exists
		organizationDbNode, err := services.Repositories.Neo4jRepositories.OrganizationReadRepository.GetOrganizationByIdOrCustomerOsId(ctx, tenant, orgId)
		if err != nil {
			tracing.TraceErr(span, err)
			c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "Organization not found"})
			return
		}
		if organizationDbNode == nil {
			c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "Organization not found"})
			return
		}

		// Validate if the external system exists
		externalSystemValid := neo4jentity.IsValidDataSource(externalSystem)
		if !externalSystemValid {
			c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "External system not found"})
			span.LogFields(tracingLog.String("result", "External system not found"))
			return
		}

		// Set or replace the primary ID
		err = services.CommonServices.ExternalSystemService.SetPrimaryExternalId(ctx, externalSystem, request.ExternalId,
			commonservice.LinkWith{
				Type: commonmodel.ORGANIZATION,
				Id:   orgId,
			})
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "Failed to set primary external ID"))
			services.Log.Error(ctx, "Failed to set primary external ID", err)
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to set primary external ID"})
			return
		}

		span.LogFields(tracingLog.String("result", "Primary external ID set successfully"))
		c.JSON(http.StatusOK, SetPrimaryExternalSystemIdResponse{
			Status:            "success",
			Message:           "Primary external ID set successfully",
			OrganizationId:    orgId,
			ExternalSystem:    externalSystem,
			PrimaryExternalId: request.ExternalId,
		})
	}
}

// @Summary Get an organization
// @Description Retrieves an organization by its ID or COS ID
// @Tags CustomerBASE API
// @Accept  json
// @Produce  json
// @Param   id   path     string  true  "Organization ID or Organization COS ID"
// @Success 200 {object} OrganizationResponse "Organization retrieved successfully"
// @Failure 400  "Invalid organization ID"
// @Failure 401  "Unauthorized access"
// @Failure 404  "Organization not found"
// @Failure 500  "Internal server error"
// @Router /customerbase/v1/organizations/{id} [get]
// @Security ApiKeyAuth
func GetOrganization(services *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "GetOrganization", c.Request.Header)
		defer span.Finish()
		tracing.TagComponentRest(span)
		tracing.TagTenant(span, common.GetTenantFromContext(ctx))

		tenant := common.GetTenantFromContext(ctx)
		// if tenant missing return auth error
		if tenant == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "API key invalid or expired"})
			span.LogFields(tracingLog.String("result", "Missing tenant in context"))
			return
		}

		// Extract organization ID from the path
		orgID := c.Param("id")
		if orgID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid organization ID"})
			span.LogFields(tracingLog.String("result", "Invalid organization ID"))
			return
		}

		// Check organization exists
		organizationDbNode, err := services.Repositories.Neo4jRepositories.OrganizationReadRepository.GetOrganizationByIdOrCustomerOsId(ctx, tenant, orgID)
		if err != nil {
			tracing.TraceErr(span, err)
			c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "Organization not found"})
			return
		}
		if organizationDbNode == nil {
			c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "Organization not found"})
			return
		}
		organizationEntity := neo4jmapper.MapDbNodeToOrganizationEntity(organizationDbNode)

		response := OrganizationResponse{
			Status:        "success",
			Message:       "Organization retrieved successfully",
			ID:            organizationEntity.ID,
			CustomId:      organizationEntity.ReferenceId,
			CosId:         organizationEntity.CustomerOsId,
			Name:          organizationEntity.Name,
			Website:       organizationEntity.Website,
			LeadSource:    organizationEntity.LeadSource,
			Relationship:  organizationEntity.Relationship.String(),
			IcpFit:        organizationEntity.IcpFit,
			Stage:         organizationEntity.Stage.String(),
			Domains:       []string{},
			ExternalLinks: []ExternalLink{},
		}

		// Fetch domains associated with the organization
		partialSuccess := false
		domainEntities, err := services.CommonServices.DomainService.GetAllDomainsForOrganizations(ctx, []string{organizationEntity.ID})
		if err != nil {
			partialSuccess = true
			tracing.TraceErr(span, errors.Wrap(err, "Failed to retrieve domains"))
		} else {
			for _, domain := range *domainEntities {
				response.Domains = append(response.Domains, domain.Domain)
			}
		}

		// Fetch external links associated with the organization
		externalSystemEntities, err := services.ExternalSystemService.GetExternalSystemsForEntities(ctx, []string{organizationEntity.ID}, commonmodel.ORGANIZATION)
		if err != nil {
			partialSuccess = true
			tracing.TraceErr(span, errors.Wrap(err, "Failed to retrieve external links"))
		} else {
			for _, externalSystemEntity := range *externalSystemEntities {
				if externalSystemEntity.Relationship.ExternalId != "" {
					response.ExternalLinks = append(response.ExternalLinks, ExternalLink{
						Name:    externalSystemEntity.ExternalSystemId.String(),
						Id:      externalSystemEntity.Relationship.ExternalId,
						Primary: externalSystemEntity.Relationship.Primary,
					})
				}
			}
		}

		if partialSuccess {
			response.Status = "partial_success"
			response.Message = "Failed to retrieve completed organization data"
			c.JSON(http.StatusPartialContent, response)
		} else {
			c.JSON(http.StatusOK, response)
		}
	}
}
