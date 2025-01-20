// @openapi 3.0.0
package customerbase

import (
	"context"
	"net/http"
	"strings"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/data_fields"
	commonModel "github.com/customeros/customeros/packages/server/customer-os-common-module/model"
	common_srv "github.com/customeros/customeros/packages/server/customer-os-common-module/services/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"
	neoEnum "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/enum"
	neo4jmapper "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/mapper"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/customeros/customeros/packages/server/customer-os-api/constants"
	"github.com/customeros/customeros/packages/server/customer-os-api/graphql/model"
	enummapper "github.com/customeros/customeros/packages/server/customer-os-api/mapper/enum"
	"github.com/customeros/customeros/packages/server/customer-os-api/rest/response"
	cosapi_services "github.com/customeros/customeros/packages/server/customer-os-api/services"
)

type OrganizationHandler struct {
	services        *cosapi_services.Services
	responseHandler *response.Response
}

func NewOrganizationHandler(services *cosapi_services.Services, responseHandler *response.Response) *OrganizationHandler {
	return &OrganizationHandler{
		services:        services,
		responseHandler: responseHandler,
	}
}

type APIStatus string

const (
	APIStatusError          APIStatus = "error"
	APIStatusSuccess        APIStatus = "success"
	APIStatusPartialSuccess APIStatus = "partial_success"
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
func (h *OrganizationHandler) CreateOrganization() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "Customerbase.CreateOrganization", c.Request.Header)
		defer span.Finish()
		tracing.TagComponentRest(span)

		tenant := common.GetTenantFromContext(ctx)
		if tenant == "" {
			h.responseHandler.HandleError(c, http.StatusNotFound, nil)
			return
		}

		var request CreateOrganizationRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			message := "Invalid request body"
			tracing.TraceErr(span, errors.Wrap(err, message))
			h.services.Log.Error(ctx, message, err)
			h.responseHandler.HandleError(c, http.StatusBadRequest, &message)
			return
		}
		if err := h.validateOrganizationRequest(c, &request); err != nil {
			return
		}

		orgFields := h.buildOrganizationFields(request)
		organizationId, err := h.services.CommonServices.OrganizationService.Save(ctx, nil, nil, orgFields)
		if err != nil {
			message := "Failed to create organization"
			tracing.TraceErr(span, errors.Wrap(err, message))
			h.services.Log.Error(ctx, message, err)
			h.responseHandler.HandleError(c, http.StatusInternalServerError, &message)
			return
		}

		resp := OrganizationResponse{
			Organization: OrganizationRecord{
				ID: organizationId,
			},
		}
		h.responseHandler.HandleCreated(c, resp)
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
func (h *OrganizationHandler) GetOrganization() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "Customerbase.GetOrganization", c.Request.Header)
		defer span.Finish()
		tracing.TagComponentRest(span)

		tenant := common.GetTenantFromContext(ctx)
		if tenant == "" {
			h.responseHandler.HandleError(c, http.StatusNotFound, nil)
			return
		}

		orgID := c.Param("id")
		if orgID == "" {
			message := "Invalid organization ID"
			h.responseHandler.HandleError(c, http.StatusBadRequest, &message)
			return
		}

		result, status := h.retrieveOrganization(c, orgID)

		switch {
		case status == APIStatusError:
			h.responseHandler.HandleError(c, http.StatusNotFound, nil)
			return
		case status == APIStatusPartialSuccess:
			message := "Unable to retrieve full organization data"
			h.responseHandler.HandleError(c, http.StatusPartialContent, &message)
			return
		default:
			resp := OrganizationResponse{
				Organization: result,
			}
			h.responseHandler.HandleSuccess(c, resp)
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
func (h *OrganizationHandler) SetPrimaryExternalSystemId() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "Customerbase.SetPrimaryExternalSystemId", c.Request.Header)
		defer span.Finish()
		tracing.TagComponentRest(span)

		tenant := common.GetTenantFromContext(ctx)
		if tenant == "" {
			h.responseHandler.HandleError(c, http.StatusNotFound, nil)
			return
		}

		result, statusCode := h.handleExternalSystemUpdate(c)

		switch {
		case statusCode == http.StatusBadRequest:
			h.responseHandler.HandleError(c, statusCode, nil)
		case statusCode == http.StatusInternalServerError:
			h.responseHandler.HandleError(c, statusCode, nil)
		case statusCode == http.StatusNotFound:
			h.responseHandler.HandleError(c, statusCode, nil)
		default:
			resp := ExternalSystemResponse{
				Organization: result,
			}
			h.responseHandler.HandleSuccess(c, resp)
		}
	}
}

func (h *OrganizationHandler) validateOrganizationRequest(c *gin.Context, request *CreateOrganizationRequest) error {
	span, ctx := tracing.StartTracerSpan(c.Request.Context(), "Customerbase.validateOrganizationRequest")
	defer span.Finish()
	tracing.TagComponentRest(span)

	if request.Name == "" && request.CustomId == "" && request.Website == "" && request.LinkedinUrl == "" {
		message := "Missing organization input fields"
		h.responseHandler.HandleError(c, http.StatusBadRequest, &message)
		return errors.New("missing required fields")
	}

	// Validate website domain
	websiteDomain, _ := h.services.CommonServices.DomainService.GetPrimaryDomainForOrganizationWebsite(ctx, request.Website)
	if websiteDomain != "" {
		if exists, err := h.checkOrganizationExistsByDomain(ctx, websiteDomain); err != nil {
			message := "Failed to check organization domain"
			h.responseHandler.HandleError(c, http.StatusInternalServerError, &message)
			return err
		} else if exists {
			message := "Organization already exists with given domain"
			h.responseHandler.HandleError(c, http.StatusConflict, &message)
			return errors.New("organization exists")
		}
	}

	// Validate custom ID
	if request.CustomId != "" {
		if exists, err := h.checkOrganizationExistsByCustomId(ctx, request.CustomId); err != nil {
			message := "Failed to check organization custom id"
			h.responseHandler.HandleError(c, http.StatusInternalServerError, &message)
			return err
		} else if exists {
			message := "Organization already exists with given custom id"
			h.responseHandler.HandleError(c, http.StatusConflict, &message)
			return errors.New("organization exists")
		}
	}

	// Validate LinkedIn URL
	if request.LinkedinUrl != "" {
		if exists, err := h.checkOrganizationExistsBySocialUrl(ctx, request.LinkedinUrl); err != nil {
			message := "Failed to check organization linkedin url"
			h.responseHandler.HandleError(c, http.StatusInternalServerError, &message)
			return err
		} else if exists {
			message := "Organization already exists with given linkedin url"
			h.responseHandler.HandleError(c, http.StatusConflict, &message)
			return errors.New("organization exists")
		}
	}

	return nil
}

func (h *OrganizationHandler) buildOrganizationFields(request CreateOrganizationRequest) data_fields.OrganizationFields {
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
	fields.Stage = h.determineOrganizationStage(relationship)

	return fields
}

func (h *OrganizationHandler) determineOrganizationStage(relationship model.OrganizationRelationship) *neoEnum.OrganizationStage {
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

func (h *OrganizationHandler) retrieveOrganization(c *gin.Context, orgID string) (OrganizationRecord, APIStatus) {
	span, ctx := tracing.StartTracerSpan(c.Request.Context(), "Customerbase.retrieveOrganization")
	defer span.Finish()
	tracing.TagComponentRest(span)

	var result OrganizationRecord

	organizationDbNode, err := h.services.Repositories.Neo4jRepositories.OrganizationReadRepository.GetOrganizationByIdOrCustomerOsId(
		ctx, common.GetTenantFromContext(ctx), orgID)
	if err != nil {
		tracing.TraceErr(span, err)
		return result, APIStatusError
	}
	if organizationDbNode == nil {
		return result, APIStatusError
	}

	organizationEntity := neo4jmapper.MapDbNodeToOrganizationEntity(organizationDbNode)
	result = h.mapOrganizationEntityToResult(organizationEntity)

	// Fetch additional data
	partialSuccess := false
	if err := h.enrichOrganizationWithDomains(ctx, &result, organizationEntity.ID); err != nil {
		partialSuccess = true
	}
	if err := h.enrichOrganizationWithExternalLinks(ctx, &result, organizationEntity.ID); err != nil {
		partialSuccess = true
	}

	if partialSuccess {
		return result, APIStatusPartialSuccess
	}

	return result, APIStatusSuccess
}

func (h *OrganizationHandler) handleExternalSystemUpdate(c *gin.Context) (ExternalSystemRecord, int) {
	span, ctx := tracing.StartTracerSpan(c.Request.Context(), "Customerbase.handleExternalSystemUpdate")
	defer span.Finish()
	tracing.TagComponentRest(span)

	var results ExternalSystemRecord

	orgId := c.Param("id")
	externalSystem := strings.ToLower(c.Param("externalSystem"))

	var request SetPrimaryExternalSystemIdRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		return results, http.StatusBadRequest
	}

	// Validate organization exists
	if exists, err := h.validateOrganizationExists(ctx, orgId); err != nil {
		return results, http.StatusInternalServerError
	} else if !exists {
		return results, http.StatusNotFound
	}

	// Validate external system
	if !neo4jentity.IsValidDataSource(externalSystem) {
		return results, http.StatusNotFound
	}

	// Set primary ID
	err := h.services.CommonServices.ExternalSystemService.SetPrimaryExternalId(ctx, externalSystem, request.ExternalId,
		common_srv.LinkWith{
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
func (h *OrganizationHandler) checkOrganizationExistsByDomain(ctx context.Context, domain string) (bool, error) {
	orgDbNode, err := h.services.Repositories.Neo4jRepositories.OrganizationReadRepository.GetOrganizationByDomain(ctx, nil, common.GetTenantFromContext(ctx), domain)
	return orgDbNode != nil, err
}

func (h *OrganizationHandler) checkOrganizationExistsByCustomId(ctx context.Context, customId string) (bool, error) {
	orgDbNode, err := h.services.Repositories.Neo4jRepositories.OrganizationReadRepository.GetOrganizationByReferenceId(ctx, common.GetTenantFromContext(ctx), customId)
	return orgDbNode != nil, err
}

func (h *OrganizationHandler) checkOrganizationExistsBySocialUrl(ctx context.Context, url string) (bool, error) {
	orgDbNode, err := h.services.Repositories.Neo4jRepositories.OrganizationReadRepository.GetOrganizationBySocialUrl(ctx, common.GetTenantFromContext(ctx), url)
	return orgDbNode != nil, err
}

func (h *OrganizationHandler) validateOrganizationExists(ctx context.Context, orgId string) (bool, error) {
	organizationDbNode, err := h.services.Repositories.Neo4jRepositories.OrganizationReadRepository.GetOrganizationByIdOrCustomerOsId(ctx, common.GetTenantFromContext(ctx), orgId)
	return organizationDbNode != nil, err
}

func (h *OrganizationHandler) enrichOrganizationWithDomains(ctx context.Context, result *OrganizationRecord, orgId string) error {
	domainEntities, err := h.services.CommonServices.DomainService.GetAllDomainsForOrganizations(ctx, []string{orgId})
	if err != nil {
		return err
	}

	result.Domains = make([]string, 0, len(*domainEntities))
	for _, domain := range *domainEntities {
		result.Domains = append(result.Domains, domain.Domain)
	}
	return nil
}

func (h *OrganizationHandler) enrichOrganizationWithExternalLinks(ctx context.Context, result *OrganizationRecord, orgId string) error {
	externalSystemEntities, err := h.services.CommonServices.ExternalSystemService.GetExternalSystemsForEntities(ctx, []string{orgId}, commonModel.ORGANIZATION)
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

func (h *OrganizationHandler) mapOrganizationEntityToResult(entity *neo4jentity.OrganizationEntity) OrganizationRecord {
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
