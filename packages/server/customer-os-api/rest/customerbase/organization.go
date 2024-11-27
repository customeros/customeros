// @openapi 3.0.0
package customerbase

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/data_fields"
	commonModel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	commonservice "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	"github.com/pkg/errors"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	enummapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/mapper/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/rest"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/service"
)

// @Summary Create a new organization
// @Description Creates an organization if it doesn't exist based on website, custom ID, or LinkedIn URL. Returns existing organization if found.
// @Tags CustomerBASE API
// @Accept json
// @Produce json
// @Param body body CreateOrganizationRequest true "Organization creation request"
// @Success 201 {object} OrganizationResponse "Organization created successfully"
// @Success 206 {object} OrganizationResponse "Organization created with partial data"
// @Failure 400 {object} rest.BaseResponse "Invalid request - Missing required fields"
// @Failure 401 {object} rest.BaseResponse "Unauthorized - Invalid or missing API key"
// @Failure 409 {object} rest.BaseResponse "Conflict - Organization already exists with provided identifiers"
// @Failure 500 {object} rest.BaseResponse "Internal server error"
// @Router /customerbase/v1/organizations [post]
// @Security ApiKeyAuth
func CreateOrganization(services *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "CreateOrganization", c.Request.Header)
		defer span.Finish()
		tracing.TagComponentRest(span)

		tenant := rest.ValidateTenant(c, ctx, span)
		if tenant == "" {
			return
		}

		var request CreateOrganizationRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "Invalid request body"))
			services.Log.Error(ctx, "Invalid request body", err)
			rest.SendError(c, span, http.StatusBadRequest, rest.ErrBadRequest)
			return
		}

		httpContext := rest.HTTPContext{
			GinContext:     c,
			ServiceContext: &ctx,
			Span:           span,
			Services:       services,
			Tenant:         tenant,
		}

		if err := validateOrganizationRequest(httpContext, &request); err != nil {
			return
		}

		orgFields := buildOrganizationFields(request)
		organizationId, err := services.CommonServices.OrganizationService.Save(ctx, nil, nil, orgFields)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "Failed to create organization"))
			services.Log.Error(ctx, "Failed to create organization", err)
			rest.SendError(c, span, http.StatusInternalServerError, rest.ErrInternalServer.WithMessage("Failed to create organization"))
			return
		}

		c.JSON(http.StatusCreated, OrganizationResponse{
			BaseResponse: rest.BuildBaseResponse(rest.StatusSuccess),
			Organization: OrganizationRecord{
				ID: organizationId,
			},
		})
	}
}

// @Summary Get organization details
// @Description Retrieves detailed organization information by ID or COS ID
// @Tags CustomerBASE API
// @Accept json
// @Produce json
// @Param id path string true "Organization ID or COS ID" example(org_123)
// @Success 200 {object} OrganizationResponse "Organization found"
// @Success 206 {object} OrganizationResponse "Organization found with partial data"
// @Failure 400 {object} rest.BaseResponse "Invalid organization ID format"
// @Failure 401 {object} rest.BaseResponse "Unauthorized - Invalid or missing API key"
// @Failure 404 {object} rest.BaseResponse "Organization not found"
// @Failure 500 {object} rest.BaseResponse "Internal server error"
// @Router /customerbase/v1/organizations/{id} [get]
// @Security ApiKeyAuth
func GetOrganization(services *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "GetOrganization", c.Request.Header)
		defer span.Finish()
		tracing.TagComponentRest(span)

		tenant := rest.ValidateTenant(c, ctx, span)
		if tenant == "" {
			return
		}

		orgID := c.Param("id")
		if orgID == "" {
			rest.SendError(c, span, http.StatusBadRequest, rest.ErrBadRequest.WithMessage("Invalid organization ID"))
			return
		}

		httpContext := rest.HTTPContext{
			GinContext:     c,
			ServiceContext: &ctx,
			Span:           span,
			Services:       services,
			Tenant:         tenant,
		}

		result, status := retrieveOrganization(httpContext, orgID)

		switch {
		case status == rest.StatusError:
			rest.SendError(c, span, http.StatusNotFound, rest.ErrNotFound.WithMessage("Organization does not exist"))
			return
		case status == rest.StatusPartialSuccess:
			c.JSON(http.StatusPartialContent, "Unable to retrieve full organization data")
			return
		default:
			c.JSON(http.StatusOK, OrganizationResponse{
				BaseResponse: rest.BuildBaseResponse(rest.StatusSuccess),
				Organization: result,
			})
			return
		}
	}
}

// @Summary Set primary external system ID
// @Description Sets or updates the primary external system identifier for an organization
// @Tags CustomerBASE API
// @Accept json
// @Produce json
// @Param id path string true "Organization ID or COS ID" example(org_123)
// @Param externalSystem path string true "External system name" example(salesforce)
// @Param body body SetPrimaryExternalSystemIdRequest true "External system ID details"
// @Success 200 {object} ExternalSystemResponse "Primary ID set successfully"
// @Failure 400 {object} rest.BaseResponse "Invalid request parameters"
// @Failure 401 {object} rest.BaseResponse "Unauthorized - Invalid or missing API key"
// @Failure 404 {object} rest.BaseResponse "Organization or external system not found"
// @Failure 500 {object} rest.BaseResponse "Internal server error"
// @Router /customerbase/v1/organizations/{id}/links/{externalSystem}/primary [put]
// @Security ApiKeyAuth
func SetPrimaryExternalSystemId(services *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "SetPrimaryExternalSystemId", c.Request.Header)
		defer span.Finish()
		tracing.TagComponentRest(span)

		tenant := rest.ValidateTenant(c, ctx, span)
		if tenant == "" {
			return
		}

		httpContext := rest.HTTPContext{
			GinContext:     c,
			ServiceContext: &ctx,
			Span:           span,
			Services:       services,
			Tenant:         tenant,
		}

		result, statusCode := handleExternalSystemUpdate(httpContext)

		switch {
		case statusCode == http.StatusBadRequest:
			rest.SendError(c, span, statusCode, rest.ErrBadRequest)
		case statusCode == http.StatusInternalServerError:
			rest.SendError(c, span, statusCode, rest.ErrInternalServer)
		case statusCode == http.StatusNotFound:
			rest.SendError(c, span, statusCode, rest.ErrNotFound)
		default:
			c.JSON(statusCode, ExternalSystemResponse{
				BaseResponse: rest.BuildBaseResponse(rest.StatusSuccess),
				Organization: result,
			})
		}
	}
}

func validateOrganizationRequest(ctx rest.HTTPContext, request *CreateOrganizationRequest) error {
	if request.Name == "" && request.CustomId == "" && request.Website == "" && request.LinkedinUrl == "" {
		rest.SendError(ctx.GinContext, ctx.Span, http.StatusBadRequest, rest.ErrBadRequest.WithMessage("Missing organization input fields"))
		return errors.New("missing required fields")
	}

	// Validate website domain
	websiteDomain, _ := ctx.Services.CommonServices.DomainService.GetPrimaryDomainForOrganizationWebsite(*ctx.ServiceContext, request.Website)
	if websiteDomain != "" {
		if exists, err := checkOrganizationExistsByDomain(*ctx.ServiceContext, ctx.Services, websiteDomain); err != nil {
			rest.SendError(ctx.GinContext, ctx.Span, http.StatusInternalServerError, rest.ErrInternalServer.WithMessage("Failed to check organization domain"))
			return err
		} else if exists {
			rest.SendError(ctx.GinContext, ctx.Span, http.StatusConflict, rest.ErrConflict.WithMessage("Organization already exists with given domain"))
			return errors.New("organization exists")
		}
	}

	// Validate custom ID
	if request.CustomId != "" {
		if exists, err := checkOrganizationExistsByCustomId(*ctx.ServiceContext, ctx.Services, request.CustomId); err != nil {
			rest.SendError(ctx.GinContext, ctx.Span, http.StatusInternalServerError, rest.ErrInternalServer.WithMessage("Failed to check organization custom id"))
			return err
		} else if exists {
			rest.SendError(ctx.GinContext, ctx.Span, http.StatusConflict, rest.ErrConflict.WithMessage("Organization already exists with given custom id"))
			return errors.New("organization exists")
		}
	}

	// Validate LinkedIn URL
	if request.LinkedinUrl != "" {
		if exists, err := checkOrganizationExistsBySocialUrl(*ctx.ServiceContext, ctx.Services, request.LinkedinUrl); err != nil {
			rest.SendError(ctx.GinContext, ctx.Span, http.StatusInternalServerError, rest.ErrInternalServer.WithMessage("Failed to check organization linkedin url"))
			return err
		} else if exists {
			rest.SendError(ctx.GinContext, ctx.Span, http.StatusConflict, rest.ErrConflict.WithMessage("Organization already exists with given linkedin url"))
			return errors.New("organization exists")
		}
	}

	return nil
}

func buildOrganizationFields(request CreateOrganizationRequest) data_fields.OrganizationFields {
	relationship := model.OrganizationRelationshipProspect
	if request.Relationship != "" && model.OrganizationRelationship(request.Relationship).IsValid() {
		relationship = model.OrganizationRelationship(request.Relationship)
	}

	fields := data_fields.OrganizationFields{
		Name:         utils.StringPtr(request.Name),
		ReferenceId:  utils.StringPtr(request.CustomId),
		Website:      utils.StringPtr(request.Website),
		Source:       utils.StringPtr(string(neo4jentity.DataSourceOpenline)),
		AppSource:    utils.StringPtr(constants.AppSourceCustomerOsApiRest),
		LeadSource:   utils.StringPtr(request.LeadSource),
		IcpFit:       utils.BoolPtr(request.IcpFit),
		Relationship: utils.ToPtr(enummapper.MapRelationshipFromModel(relationship)),
	}

	if request.LinkedinUrl != "" {
		fields.LinkedInUrl = utils.StringPtr(request.LinkedinUrl)
	}

	// Set stage based on relationship
	fields.Stage = determineOrganizationStage(relationship)

	return fields
}

func determineOrganizationStage(relationship model.OrganizationRelationship) *enum.OrganizationStage {
	var stage model.OrganizationStage
	switch relationship {
	case model.OrganizationRelationshipCustomer:
		stage = model.OrganizationStageOnboarding
	case model.OrganizationRelationshipProspect:
		stage = model.OrganizationStageLead
	case model.OrganizationRelationshipNotAFit:
		stage = model.OrganizationStageUnqualified
	case model.OrganizationRelationshipFormerCustomer:
		stage = model.OrganizationStageTarget
	}
	return utils.ToPtr(enummapper.MapStageFromModel(stage))
}

func retrieveOrganization(ctx rest.HTTPContext, orgID string) (OrganizationRecord, rest.Status) {
	var result OrganizationRecord

	organizationDbNode, err := ctx.Services.Repositories.Neo4jRepositories.OrganizationReadRepository.GetOrganizationByIdOrCustomerOsId(
		*ctx.ServiceContext, common.GetTenantFromContext(*ctx.ServiceContext), orgID)
	if err != nil {
		tracing.TraceErr(ctx.Span, err)
		return result, rest.StatusError
	}
	if organizationDbNode == nil {
		return result, rest.StatusError
	}

	organizationEntity := neo4jmapper.MapDbNodeToOrganizationEntity(organizationDbNode)
	result = mapOrganizationEntityToResult(organizationEntity)

	// Fetch additional data
	partialSuccess := false
	if err := enrichOrganizationWithDomains(*ctx.ServiceContext, ctx.Services, &result, organizationEntity.ID); err != nil {
		partialSuccess = true
	}
	if err := enrichOrganizationWithExternalLinks(*ctx.ServiceContext, ctx.Services, &result, organizationEntity.ID); err != nil {
		partialSuccess = true
	}

	if partialSuccess {
		return result, rest.StatusPartialSuccess
	}

	return result, rest.StatusSuccess
}

func handleExternalSystemUpdate(ctx rest.HTTPContext) (ExternalSystemRecord, int) {
	var results ExternalSystemRecord

	orgId := ctx.GinContext.Param("id")
	externalSystem := strings.ToLower(ctx.GinContext.Param("externalSystem"))

	var request SetPrimaryExternalSystemIdRequest
	if err := ctx.GinContext.ShouldBindJSON(&request); err != nil {
		return results, http.StatusBadRequest
	}

	// Validate organization exists
	if exists, err := validateOrganizationExists(*ctx.ServiceContext, ctx.Services, orgId); err != nil {
		return results, http.StatusInternalServerError
	} else if !exists {
		return results, http.StatusNotFound
	}

	// Validate external system
	if !neo4jentity.IsValidDataSource(externalSystem) {
		return results, http.StatusNotFound
	}

	// Set primary ID
	err := ctx.Services.CommonServices.ExternalSystemService.SetPrimaryExternalId(*ctx.ServiceContext, externalSystem, request.ExternalId,
		commonservice.LinkWith{
			Type: commonModel.ORGANIZATION,
			Id:   orgId,
		})
	if err != nil {
		return results, http.StatusInternalServerError
	}

	return ExternalSystemRecord{
		OrganizationId: orgId,
		ExternalSystem: externalSystem,
		ExternalId:     request.ExternalId,
		Primary:        true,
	}, http.StatusOK
}

// Helper functions for checking existence
func checkOrganizationExistsByDomain(ctx context.Context, services *service.Services, domain string) (bool, error) {
	orgDbNode, err := services.Repositories.Neo4jRepositories.OrganizationReadRepository.GetOrganizationByDomain(ctx, common.GetTenantFromContext(ctx), domain)
	return orgDbNode != nil, err
}

func checkOrganizationExistsByCustomId(ctx context.Context, services *service.Services, customId string) (bool, error) {
	orgDbNode, err := services.Repositories.Neo4jRepositories.OrganizationReadRepository.GetOrganizationByReferenceId(ctx, common.GetTenantFromContext(ctx), customId)
	return orgDbNode != nil, err
}

func checkOrganizationExistsBySocialUrl(ctx context.Context, services *service.Services, url string) (bool, error) {
	orgDbNode, err := services.Repositories.Neo4jRepositories.OrganizationReadRepository.GetOrganizationBySocialUrl(ctx, common.GetTenantFromContext(ctx), url)
	return orgDbNode != nil, err
}

func validateOrganizationExists(ctx context.Context, services *service.Services, orgId string) (bool, error) {
	organizationDbNode, err := services.Repositories.Neo4jRepositories.OrganizationReadRepository.GetOrganizationByIdOrCustomerOsId(ctx, common.GetTenantFromContext(ctx), orgId)
	return organizationDbNode != nil, err
}

func enrichOrganizationWithDomains(ctx context.Context, services *service.Services, result *OrganizationRecord, orgId string) error {
	domainEntities, err := services.CommonServices.DomainService.GetAllDomainsForOrganizations(ctx, []string{orgId})
	if err != nil {
		return err
	}

	result.Domains = make([]string, 0, len(*domainEntities))
	for _, domain := range *domainEntities {
		result.Domains = append(result.Domains, domain.Domain)
	}
	return nil
}

func enrichOrganizationWithExternalLinks(ctx context.Context, services *service.Services, result *OrganizationRecord, orgId string) error {
	externalSystemEntities, err := services.CommonServices.ExternalSystemService.GetExternalSystemsForEntities(ctx, []string{orgId}, commonModel.ORGANIZATION)
	if err != nil {
		return err
	}

	result.ExternalLinks = make([]ExternalLink, 0, len(*externalSystemEntities))
	for _, entity := range *externalSystemEntities {
		if entity.Relationship.ExternalId != "" {
			result.ExternalLinks = append(result.ExternalLinks, ExternalLink{
				Name:    entity.ExternalSystemId.String(),
				Id:      entity.Relationship.ExternalId,
				Primary: entity.Relationship.Primary,
			})
		}
	}
	return nil
}

func mapOrganizationEntityToResult(entity *neo4jentity.OrganizationEntity) OrganizationRecord {
	return OrganizationRecord{
		ID:           entity.ID,
		CustomId:     entity.ReferenceId,
		CosId:        entity.CustomerOsId,
		Name:         entity.Name,
		Website:      entity.Website,
		LeadSource:   entity.LeadSource,
		Relationship: entity.Relationship.String(),
		IcpFit:       entity.IcpFit,
		Stage:        entity.Stage.String(),
	}
}
