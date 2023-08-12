package service

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/grpc_client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	organization_grpc_service "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/organization"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"reflect"
)

type OrganizationService interface {
	Create(ctx context.Context, input *OrganizationCreateData) (*entity.OrganizationEntity, error)
	Update(ctx context.Context, input *OrganizationUpdateData) (*entity.OrganizationEntity, error)
	UpdateRenewalLikelihood(ctx context.Context, orgId string, data *entity.RenewalLikelihood) error
	UpdateRenewalForecast(ctx context.Context, orgId string, data *entity.RenewalForecast) error
	UpdateBillingDetails(ctx context.Context, orgId string, data *entity.BillingDetails) error
	GetOrganizationsForJobRoles(ctx context.Context, jobRoleIds []string) (*entity.OrganizationEntities, error)
	GetOrganizationById(ctx context.Context, organizationId string) (*entity.OrganizationEntity, error)
	FindAll(ctx context.Context, page, limit int, filter *model.Filter, sortBy []*model.SortBy) (*utils.Pagination, error)
	GetOrganizationsForContact(ctx context.Context, contactId string, page, limit int, filter *model.Filter, sortBy []*model.SortBy) (*utils.Pagination, error)
	Archive(ctx context.Context, organizationId string) error
	Merge(ctx context.Context, primaryOrganizationId, mergedOrganizationId string) error
	GetOrganizationsForEmails(ctx context.Context, emailIds []string) (*entity.OrganizationEntities, error)
	GetOrganizationsForPhoneNumbers(ctx context.Context, phoneNumberIds []string) (*entity.OrganizationEntities, error)
	GetSubsidiariesForOrganizations(ctx context.Context, parentOrganizationIds []string) (*entity.OrganizationEntities, error)
	GetSubsidiariesOfForOrganizations(ctx context.Context, organizationIds []string) (*entity.OrganizationEntities, error)
	AddSubsidiary(ctx context.Context, organizationId, subsidiaryId, subsidiaryType string) error
	RemoveSubsidiary(ctx context.Context, organizationId, subsidiaryId string) error
	ReplaceOwner(ctx context.Context, organizationId, userId string) (*entity.OrganizationEntity, error)
	RemoveOwner(ctx context.Context, organizationId string) (*entity.OrganizationEntity, error)
	AddRelationship(ctx context.Context, organizationId string, relationship entity.OrganizationRelationship) (*entity.OrganizationEntity, error)
	RemoveRelationship(ctx context.Context, organizationId string, relationship entity.OrganizationRelationship) (*entity.OrganizationEntity, error)
	SetRelationshipStage(ctx context.Context, organizationId string, relationship entity.OrganizationRelationship, stage string) (*entity.OrganizationEntity, error)
	RemoveRelationshipStage(ctx context.Context, organizationId string, relationship entity.OrganizationRelationship) (*entity.OrganizationEntity, error)
	UpdateLastTouchpointSync(ctx context.Context, organizationId string)
	UpdateLastTouchpointSyncByContactId(ctx context.Context, contactId string)
	UpdateLastTouchpointSyncByEmailId(ctx context.Context, emailId string)
	UpdateLastTouchpointSyncByPhoneNumberId(ctx context.Context, phoneNumberId string)
	UpdateLastTouchpointSyncByEmail(ctx context.Context, email string)
	UpdateLastTouchpointSyncByPhoneNumber(ctx context.Context, phoneNumber string)
	ReplaceHealthIndicator(ctx context.Context, organizationId, healthIndicatorId string) (*entity.OrganizationEntity, error)
	RemoveHealthIndicator(ctx context.Context, organizationId string) (*entity.OrganizationEntity, error)
	GetSuggestedMergeToForOrganizations(ctx context.Context, organizationIds []string) (*entity.OrganizationEntities, error)

	mapDbNodeToOrganizationEntity(node dbtype.Node) *entity.OrganizationEntity

	// Deprecated
	UpsertPhoneNumberRelationInEventStore(ctx context.Context, size int) (int, int, error)
	// Deprecated
	UpsertEmailRelationInEventStore(ctx context.Context, size int) (int, int, error)
}

type OrganizationCreateData struct {
	OrganizationEntity *entity.OrganizationEntity
	CustomFields       *entity.CustomFieldEntities
	FieldSets          *entity.FieldSetEntities
	TemplateId         *string
	Domains            []string
}

type OrganizationUpdateData struct {
	Organization *entity.OrganizationEntity
	Domains      []string
}

type organizationService struct {
	log          logger.Logger
	repositories *repository.Repositories
	grpcClients  *grpc_client.Clients
}

func NewOrganizationService(log logger.Logger, repositories *repository.Repositories, grpcClients *grpc_client.Clients) OrganizationService {
	return &organizationService{
		log:          log,
		repositories: repositories,
		grpcClients:  grpcClients,
	}
}

func (s *organizationService) Create(ctx context.Context, input *OrganizationCreateData) (*entity.OrganizationEntity, error) {
	session := utils.NewNeo4jWriteSession(ctx, *s.repositories.Drivers.Neo4jDriver)
	defer session.Close(ctx)

	organizationDbNodePtr, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		tenant := common.GetTenantFromContext(ctx)

		for _, domain := range input.Domains {
			_, err := s.repositories.DomainRepository.Merge(ctx, entity.DomainEntity{
				Domain:    domain,
				Source:    input.OrganizationEntity.Source,
				AppSource: input.OrganizationEntity.AppSource,
			})
			if err != nil {
				return nil, err
			}
		}

		organizationDbNodePtr, err := s.repositories.OrganizationRepository.Create(ctx, tx, tenant, *input.OrganizationEntity)
		if err != nil {
			return nil, err
		}
		var organizationId = utils.GetPropsFromNode(*organizationDbNodePtr)["id"].(string)
		var organizationCreatedAt = utils.GetTimePropOrNow(utils.GetPropsFromNode(*organizationDbNodePtr), "createdAt")

		err = s.repositories.OrganizationRepository.LinkWithDomainsInTx(ctx, tx, tenant, organizationId, input.Domains)
		if err != nil {
			return nil, err
		}

		entityType := &model.CustomFieldEntityType{
			ID:         organizationId,
			EntityType: model.EntityTypeOrganization,
		}
		if input.TemplateId != nil {
			err := s.repositories.ContactRepository.LinkWithEntityTemplateInTx(ctx, tx, tenant, entityType, *input.TemplateId)
			if err != nil {
				return nil, err
			}
		}
		if input.CustomFields != nil {
			for _, customField := range *input.CustomFields {
				dbNode, err := s.repositories.CustomFieldRepository.MergeCustomFieldInTx(ctx, tx, tenant, entityType, customField)
				if err != nil {
					return nil, err
				}
				if customField.TemplateId != nil {
					var fieldId = utils.GetPropsFromNode(*dbNode)["id"].(string)
					err := s.repositories.CustomFieldRepository.LinkWithCustomFieldTemplateInTx(ctx, tx, fieldId, entityType, *customField.TemplateId)
					if err != nil {
						return nil, err
					}
				}
			}
		}
		if input.FieldSets != nil {
			for _, fieldSet := range *input.FieldSets {
				setDbNode, err := s.repositories.FieldSetRepository.MergeFieldSetInTx(ctx, tx, tenant, entityType, fieldSet)
				if err != nil {
					return nil, err
				}
				var fieldSetId = utils.GetPropsFromNode(*setDbNode)["id"].(string)
				if fieldSet.TemplateId != nil {
					err := s.repositories.FieldSetRepository.LinkWithFieldSetTemplateInTx(ctx, tx, tenant, fieldSetId, *fieldSet.TemplateId, model.EntityTypeOrganization)
					if err != nil {
						return nil, err
					}
				}
				if fieldSet.CustomFields != nil {
					for _, customField := range *fieldSet.CustomFields {
						fieldDbNode, err := s.repositories.CustomFieldRepository.MergeCustomFieldToFieldSetInTx(ctx, tx, tenant, entityType, fieldSetId, customField)
						if err != nil {
							return nil, err
						}
						if customField.TemplateId != nil {
							var fieldId = utils.GetPropsFromNode(*fieldDbNode)["id"].(string)
							err := s.repositories.CustomFieldRepository.LinkWithCustomFieldTemplateForFieldSetInTx(ctx, tx, fieldId, fieldSetId, *customField.TemplateId)
							if err != nil {
								return nil, err
							}
						}
					}
				}
			}
		}

		createdAction, err := s.repositories.ActionRepository.Create(ctx, tx, tenant, organizationId, entity.ORGANIZATION, entity.ActionCreated, input.OrganizationEntity.Source, input.OrganizationEntity.AppSource)
		if err != nil {
			return nil, err
		}

		var createdActionId = utils.GetPropsFromNode(*createdAction)["id"].(string)
		err = s.repositories.OrganizationRepository.UpdateLastTouchpointInTx(ctx, tx, tenant, organizationId, organizationCreatedAt, createdActionId)
		if err != nil {
			return nil, err
		}

		return organizationDbNodePtr, nil
	})
	if err != nil {
		return nil, err
	}
	return s.mapDbNodeToOrganizationEntity(*organizationDbNodePtr.(*dbtype.Node)), nil
}

func (s *organizationService) Update(ctx context.Context, input *OrganizationUpdateData) (*entity.OrganizationEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationService.Update")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("organizationId", input.Organization.ID), log.Object("organizationUpdateDtls", input.Organization), log.Object("domains", input.Domains))

	session := utils.NewNeo4jWriteSession(ctx, *s.repositories.Drivers.Neo4jDriver)
	defer session.Close(ctx)

	organizationDbNodePtr, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		tenant := common.GetTenantFromContext(ctx)

		for _, domain := range input.Domains {
			_, err := s.repositories.DomainRepository.Merge(ctx, entity.DomainEntity{
				Domain:    domain,
				Source:    input.Organization.Source,
				AppSource: input.Organization.AppSource,
			})
			if err != nil {
				return nil, err
			}
		}

		organizationDbNodePtr, err := s.repositories.OrganizationRepository.Update(ctx, tx, tenant, *input.Organization)
		if err != nil {
			return nil, err
		}
		var organizationId = utils.GetPropsFromNode(*organizationDbNodePtr)["id"].(string)

		err = s.repositories.OrganizationRepository.LinkWithDomainsInTx(ctx, tx, tenant, organizationId, input.Domains)
		if err != nil {
			return nil, err
		}

		err = s.repositories.OrganizationRepository.UnlinkFromDomainsNotInListInTx(ctx, tx, tenant, organizationId, input.Domains)
		if err != nil {
			return nil, err
		}

		return organizationDbNodePtr, nil
	})
	if err != nil {
		return nil, err
	}
	return s.mapDbNodeToOrganizationEntity(*organizationDbNodePtr.(*dbtype.Node)), nil
}

func (s *organizationService) UpdateRenewalLikelihood(ctx context.Context, orgId string, data *entity.RenewalLikelihood) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationService.UpdateRenewalLikelihood")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("organizationId", orgId), log.Object("data", data))

	organization, err := s.GetOrganizationById(ctx, orgId)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	data.UpdatedAt = utils.TimePtr(utils.Now())
	data.UpdatedBy = utils.StringPtr(common.GetUserIdFromContext(ctx))
	if organization.RenewalLikelihood.RenewalLikelihood != data.RenewalLikelihood {
		data.PreviousRenewalLikelihood = organization.RenewalLikelihood.RenewalLikelihood
	}

	err = s.repositories.OrganizationRepository.UpdateRenewalLikelihood(ctx, orgId, data)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}

func (s *organizationService) UpdateRenewalForecast(ctx context.Context, orgId string, data *entity.RenewalForecast) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationService.UpdateRenewalForecast")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("organizationId", orgId), log.Object("data", data))

	organization, err := s.GetOrganizationById(ctx, orgId)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	data.UpdatedAt = utils.TimePtr(utils.Now())
	data.UpdatedBy = utils.StringPtr(common.GetUserIdFromContext(ctx))
	if organization.RenewalForecast.Amount != data.Amount {
		data.PreviousAmount = organization.RenewalForecast.Amount
	}

	err = s.repositories.OrganizationRepository.UpdateRenewalForecast(ctx, orgId, data)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}

func (s *organizationService) UpdateBillingDetails(ctx context.Context, orgId string, data *entity.BillingDetails) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationService.UpdateBillingDetails")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("organizationId", orgId), log.Object("data", data))

	err := s.repositories.OrganizationRepository.UpdateBillingDetails(ctx, orgId, data)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}

func (s *organizationService) FindAll(ctx context.Context, page, limit int, filter *model.Filter, sortBy []*model.SortBy) (*utils.Pagination, error) {
	var paginatedResult = utils.Pagination{
		Limit: limit,
		Page:  page,
	}
	cypherSort, err := buildSort(sortBy, reflect.TypeOf(entity.OrganizationEntity{}))
	if err != nil {
		return nil, err
	}
	cypherFilter, err := buildFilter(filter, reflect.TypeOf(entity.OrganizationEntity{}))
	if err != nil {
		return nil, err
	}

	dbNodesWithTotalCount, err := s.repositories.OrganizationRepository.GetPaginatedOrganizations(
		ctx,
		common.GetTenantFromContext(ctx),
		paginatedResult.GetSkip(),
		paginatedResult.GetLimit(),
		cypherFilter,
		cypherSort)
	if err != nil {
		return nil, err
	}
	paginatedResult.SetTotalRows(dbNodesWithTotalCount.Count)

	organizationEntities := make(entity.OrganizationEntities, 0, len(dbNodesWithTotalCount.Nodes))
	for _, v := range dbNodesWithTotalCount.Nodes {
		organizationEntities = append(organizationEntities, *s.mapDbNodeToOrganizationEntity(*v))
	}
	paginatedResult.SetRows(&organizationEntities)
	return &paginatedResult, nil
}

func (s *organizationService) GetOrganizationsForContact(ctx context.Context, contactId string, page, limit int, filter *model.Filter, sortBy []*model.SortBy) (*utils.Pagination, error) {
	var paginatedResult = utils.Pagination{
		Limit: limit,
		Page:  page,
	}
	cypherSort, err := buildSort(sortBy, reflect.TypeOf(entity.OrganizationEntity{}))
	if err != nil {
		return nil, err
	}
	cypherFilter, err := buildFilter(filter, reflect.TypeOf(entity.OrganizationEntity{}))
	if err != nil {
		return nil, err
	}

	dbNodesWithTotalCount, err := s.repositories.OrganizationRepository.GetPaginatedOrganizationsForContact(
		ctx,
		common.GetTenantFromContext(ctx),
		contactId,
		paginatedResult.GetSkip(),
		paginatedResult.GetLimit(),
		cypherFilter,
		cypherSort)
	if err != nil {
		return nil, err
	}
	paginatedResult.SetTotalRows(dbNodesWithTotalCount.Count)

	organizationEntities := make(entity.OrganizationEntities, 0, len(dbNodesWithTotalCount.Nodes))
	for _, v := range dbNodesWithTotalCount.Nodes {
		organizationEntities = append(organizationEntities, *s.mapDbNodeToOrganizationEntity(*v))
	}
	paginatedResult.SetRows(&organizationEntities)
	return &paginatedResult, nil
}

func (s *organizationService) GetOrganizationsForJobRoles(ctx context.Context, jobRoleIds []string) (*entity.OrganizationEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationService.GetOrganizationsForJobRoles")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("jobRoleIds", jobRoleIds))

	organizations, err := s.repositories.OrganizationRepository.GetAllForJobRoles(ctx, common.GetTenantFromContext(ctx), jobRoleIds)
	if err != nil {
		return nil, err
	}
	organizationEntities := make(entity.OrganizationEntities, 0, len(organizations))
	for _, v := range organizations {
		organizationEntity := s.mapDbNodeToOrganizationEntity(*v.Node)
		organizationEntity.DataloaderKey = v.LinkedNodeId
		organizationEntities = append(organizationEntities, *organizationEntity)
	}
	return &organizationEntities, nil
}

func (s *organizationService) GetOrganizationById(ctx context.Context, organizationId string) (*entity.OrganizationEntity, error) {
	dbNode, err := s.repositories.OrganizationRepository.GetOrganizationById(ctx, common.GetTenantFromContext(ctx), organizationId)
	if err != nil {
		return nil, err
	}
	return s.mapDbNodeToOrganizationEntity(*dbNode), nil
}

func (s *organizationService) Archive(ctx context.Context, organizationId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationService.Archive")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("organizationId", organizationId))

	err := s.repositories.OrganizationRepository.Archive(ctx, organizationId)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(organizationService.Archive) Error archiving organization with id {%s}: {%v}", organizationId, err.Error())
	}
	return err
}

func (s *organizationService) Merge(ctx context.Context, primaryOrganizationId, mergedOrganizationId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationService.Merge")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("primaryOrganizationId", primaryOrganizationId), log.String("mergedOrganizationId", mergedOrganizationId))

	_, err := s.GetOrganizationById(ctx, primaryOrganizationId)
	if err != nil {
		s.log.Errorf("(organizationService.Merge) Primary organization with id {%s} not found: {%v}", primaryOrganizationId, err.Error())
		return err
	}
	_, err = s.GetOrganizationById(ctx, mergedOrganizationId)
	if err != nil {
		s.log.Errorf("(organizationService.Merge) Organization to merge with id {%s} not found: {%v}", mergedOrganizationId, err.Error())
		return err
	}

	session := utils.NewNeo4jWriteSession(ctx, *s.repositories.Drivers.Neo4jDriver)
	defer session.Close(ctx)

	tenant := common.GetTenantFromContext(ctx)
	_, err = session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		err = s.repositories.OrganizationRepository.MergeOrganizationPropertiesInTx(ctx, tx, tenant, primaryOrganizationId, mergedOrganizationId, entity.DataSourceOpenline)
		if err != nil {
			return nil, err
		}

		err = s.repositories.OrganizationRepository.MergeOrganizationRelationsInTx(ctx, tx, tenant, primaryOrganizationId, mergedOrganizationId)
		if err != nil {
			return nil, err
		}

		err = s.repositories.OrganizationRepository.UpdateMergedOrganizationLabelsInTx(ctx, tx, tenant, mergedOrganizationId)
		if err != nil {
			return nil, err
		}

		return nil, nil
	})

	s.UpdateLastTouchpointSync(ctx, primaryOrganizationId)

	return err
}

func (s *organizationService) GetOrganizationsForEmails(ctx context.Context, emailIds []string) (*entity.OrganizationEntities, error) {
	organizations, err := s.repositories.OrganizationRepository.GetAllForEmails(ctx, common.GetTenantFromContext(ctx), emailIds)
	if err != nil {
		return nil, err
	}
	organizationEntities := make(entity.OrganizationEntities, 0, len(organizations))
	for _, v := range organizations {
		organizationEntity := s.mapDbNodeToOrganizationEntity(*v.Node)
		organizationEntity.DataloaderKey = v.LinkedNodeId
		organizationEntities = append(organizationEntities, *organizationEntity)
	}
	return &organizationEntities, nil
}

func (s *organizationService) GetOrganizationsForPhoneNumbers(ctx context.Context, phoneNumberIds []string) (*entity.OrganizationEntities, error) {
	organizations, err := s.repositories.OrganizationRepository.GetAllForPhoneNumbers(ctx, common.GetTenantFromContext(ctx), phoneNumberIds)
	if err != nil {
		return nil, err
	}
	organizationEntities := make(entity.OrganizationEntities, 0, len(organizations))
	for _, v := range organizations {
		organizationEntity := s.mapDbNodeToOrganizationEntity(*v.Node)
		organizationEntity.DataloaderKey = v.LinkedNodeId
		organizationEntities = append(organizationEntities, *organizationEntity)
	}
	return &organizationEntities, nil
}

func (s *organizationService) GetSubsidiariesForOrganizations(ctx context.Context, parentOrganizationIds []string) (*entity.OrganizationEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationService.GetSubsidiariesForOrganizations")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("parentOrganizationIds", parentOrganizationIds))

	dbEntries, err := s.repositories.OrganizationRepository.GetLinkedSubOrganizations(ctx, common.GetTenantFromContext(ctx), parentOrganizationIds, repository.Relationship_Subsidiary)
	if err != nil {
		s.log.Errorf("(organizationService.GetSubsidiariesForOrganizations) Error getting linked organizations: {%v}", err.Error())
		tracing.TraceErr(span, err)
		return nil, err
	}
	organizationEntities := make(entity.OrganizationEntities, 0, len(dbEntries))
	for _, v := range dbEntries {
		organizationEntity := s.mapDbNodeToOrganizationEntity(*v.Node)
		s.addLinkedOrganizationRelationshipToOrganizationEntity(*v.Relationship, organizationEntity)
		organizationEntity.DataloaderKey = v.LinkedNodeId
		organizationEntities = append(organizationEntities, *organizationEntity)
	}
	return &organizationEntities, nil
}

func (s *organizationService) AddSubsidiary(ctx context.Context, organizationId, subsidiaryId, subsidiaryType string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationService.AddSubsidiary")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("organizationId", organizationId), log.String("subsidiaryId", subsidiaryId), log.String("subsidiaryType", subsidiaryType))

	err := s.repositories.OrganizationRepository.LinkSubOrganization(ctx, common.GetTenantFromContext(ctx), organizationId, subsidiaryId, subsidiaryType, repository.Relationship_Subsidiary)
	if err != nil {
		s.log.Errorf("(organizationService.AddSubsidiary) Error adding subsidiary: {%v}", err.Error())
		tracing.TraceErr(span, err)
	}
	return err
}

func (s *organizationService) RemoveSubsidiary(ctx context.Context, organizationId, subsidiaryId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationService.AddSubsidiary")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("organizationId", organizationId), log.String("subsidiaryId", subsidiaryId))

	err := s.repositories.OrganizationRepository.UnlinkSubOrganization(ctx, common.GetTenantFromContext(ctx), organizationId, subsidiaryId, repository.Relationship_Subsidiary)
	if err != nil {
		s.log.Errorf("(organizationService.RemoveSubsidiary) Error removing subsidiary: {%v}", err.Error())
		tracing.TraceErr(span, err)
	}
	return err
}

func (s *organizationService) GetSubsidiariesOfForOrganizations(ctx context.Context, organizationIds []string) (*entity.OrganizationEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationService.GetSubsidiariesOfForOrganizations")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("organizationIds", organizationIds))

	dbEntries, err := s.repositories.OrganizationRepository.GetLinkedParentOrganizations(ctx, common.GetTenantFromContext(ctx), organizationIds, repository.Relationship_Subsidiary)
	if err != nil {
		s.log.Errorf("(organizationService.GetSubsidiariesOfForOrganizations) Error getting linked parent organizations: {%v}", err.Error())
		tracing.TraceErr(span, err)
		return nil, err
	}
	organizationEntities := make(entity.OrganizationEntities, 0, len(dbEntries))
	for _, v := range dbEntries {
		organizationEntity := s.mapDbNodeToOrganizationEntity(*v.Node)
		s.addLinkedOrganizationRelationshipToOrganizationEntity(*v.Relationship, organizationEntity)
		organizationEntity.DataloaderKey = v.LinkedNodeId
		organizationEntities = append(organizationEntities, *organizationEntity)
	}
	return &organizationEntities, nil
}

func (s *organizationService) GetSuggestedMergeToForOrganizations(ctx context.Context, organizationIds []string) (*entity.OrganizationEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationService.GetSuggestedMergeToForOrganizations")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("organizationIds", organizationIds))

	dbEntries, err := s.repositories.OrganizationRepository.GetSuggestedMergePrimaryOrganizations(ctx, organizationIds)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error getting suggested merge primary organizations: {%v}", err.Error())
		return nil, err
	}
	organizationEntities := make(entity.OrganizationEntities, 0, len(dbEntries))
	for _, v := range dbEntries {
		organizationEntity := s.mapDbNodeToOrganizationEntity(*v.Node)
		s.addSuggestedMergeRelationshipToOrganizationEntity(*v.Relationship, organizationEntity)
		organizationEntity.DataloaderKey = v.LinkedNodeId
		organizationEntities = append(organizationEntities, *organizationEntity)
	}
	return &organizationEntities, nil
}

func (s *organizationService) ReplaceOwner(ctx context.Context, organizationID, userID string) (*entity.OrganizationEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationService.ReplaceOwner")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("organizationID", organizationID), log.String("userID", userID))

	dbNode, err := s.repositories.OrganizationRepository.ReplaceOwner(ctx, common.GetTenantFromContext(ctx), organizationID, userID)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	return s.mapDbNodeToOrganizationEntity(*dbNode), nil
}

func (s *organizationService) RemoveOwner(ctx context.Context, organizationID string) (*entity.OrganizationEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationService.RemoveOwner")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("organizationID", organizationID))

	dbNode, err := s.repositories.OrganizationRepository.RemoveOwner(ctx, common.GetTenantFromContext(ctx), organizationID)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	return s.mapDbNodeToOrganizationEntity(*dbNode), nil
}

func (s *organizationService) AddRelationship(ctx context.Context, organizationID string, relationship entity.OrganizationRelationship) (*entity.OrganizationEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationService.AddRelationship")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("organizationID", organizationID), log.String("relationship", relationship.String()))

	dbNode, err := s.repositories.OrganizationRepository.AddRelationship(ctx, common.GetTenantFromContext(ctx), organizationID, relationship.String())
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	return s.mapDbNodeToOrganizationEntity(*dbNode), nil
}

func (s *organizationService) SetRelationshipStage(ctx context.Context, organizationID string, relationship entity.OrganizationRelationship, stage string) (*entity.OrganizationEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationService.SetRelationshipWithStage")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("organizationID", organizationID), log.String("relationship", relationship.String()), log.String("stage", stage))

	dbNode, err := s.repositories.OrganizationRepository.SetRelationshipWithStage(ctx, common.GetTenantFromContext(ctx), organizationID, relationship.String(), stage)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	return s.mapDbNodeToOrganizationEntity(*dbNode), nil
}

func (s *organizationService) RemoveRelationship(ctx context.Context, organizationID string, relationship entity.OrganizationRelationship) (*entity.OrganizationEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationService.RemoveRelationship")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("organizationID", organizationID), log.String("relationship", relationship.String()))

	dbNode, err := s.repositories.OrganizationRepository.RemoveRelationship(ctx, common.GetTenantFromContext(ctx), organizationID, relationship.String())
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	return s.mapDbNodeToOrganizationEntity(*dbNode), nil
}

func (s *organizationService) RemoveRelationshipStage(ctx context.Context, organizationID string, relationship entity.OrganizationRelationship) (*entity.OrganizationEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationService.RemoveRelationshipStage")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("organizationID", organizationID), log.String("relationship", relationship.String()))

	dbNode, err := s.repositories.OrganizationRepository.RemoveRelationshipStage(ctx, common.GetTenantFromContext(ctx), organizationID, relationship.String())
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	return s.mapDbNodeToOrganizationEntity(*dbNode), nil
}

func (s *organizationService) ReplaceHealthIndicator(ctx context.Context, organizationId, healthIndicatorId string) (*entity.OrganizationEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationService.ReplaceHealthIndicator")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("organizationId", organizationId), log.String("healthIndicatorId", healthIndicatorId))

	dbNode, err := s.repositories.OrganizationRepository.ReplaceHealthIndicator(ctx, organizationId, healthIndicatorId)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	return s.mapDbNodeToOrganizationEntity(*dbNode), nil
}

func (s *organizationService) RemoveHealthIndicator(ctx context.Context, organizationId string) (*entity.OrganizationEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationService.RemoveHealthIndicator")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("organizationId", organizationId))

	dbNode, err := s.repositories.OrganizationRepository.RemoveHealthIndicator(ctx, organizationId)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	return s.mapDbNodeToOrganizationEntity(*dbNode), nil
}

func (s *organizationService) UpsertPhoneNumberRelationInEventStore(ctx context.Context, size int) (int, int, error) {
	processedRecords := 0
	failedRecords := 0
	outputErr := error(nil)
	for size > 0 {
		batchSize := constants.Neo4jBatchSize
		if size < constants.Neo4jBatchSize {
			batchSize = size
		}
		records, err := s.repositories.OrganizationRepository.GetAllOrganizationPhoneNumberRelationships(ctx, batchSize)
		if err != nil {
			return 0, 0, err
		}
		for _, v := range records {
			_, err := s.grpcClients.OrganizationClient.LinkPhoneNumberToOrganization(context.Background(), &organization_grpc_service.LinkPhoneNumberToOrganizationGrpcRequest{
				Primary:        utils.GetBoolPropOrFalse(v.Values[0].(neo4j.Relationship).Props, "primary"),
				Label:          utils.GetStringPropOrEmpty(v.Values[0].(neo4j.Relationship).Props, "label"),
				OrganizationId: v.Values[1].(string),
				PhoneNumberId:  v.Values[2].(string),
				Tenant:         v.Values[3].(string),
			})
			if err != nil {
				failedRecords++
				if outputErr != nil {
					outputErr = err
				}
				s.log.Errorf("(organizationService.UpsertPhoneNumberRelationInEventStore) Failed to call method: {%v}", err.Error())
			} else {
				processedRecords++
			}
		}

		size -= batchSize
	}

	return processedRecords, failedRecords, outputErr
}

func (s *organizationService) UpsertEmailRelationInEventStore(ctx context.Context, size int) (int, int, error) {
	processedRecords := 0
	failedRecords := 0
	outputErr := error(nil)
	for size > 0 {
		batchSize := constants.Neo4jBatchSize
		if size < constants.Neo4jBatchSize {
			batchSize = size
		}
		records, err := s.repositories.OrganizationRepository.GetAllOrganizationEmailRelationships(ctx, batchSize)
		if err != nil {
			return 0, 0, err
		}
		for _, v := range records {
			_, err := s.grpcClients.OrganizationClient.LinkEmailToOrganization(context.Background(), &organization_grpc_service.LinkEmailToOrganizationGrpcRequest{
				Primary:        utils.GetBoolPropOrFalse(v.Values[0].(neo4j.Relationship).Props, "primary"),
				Label:          utils.GetStringPropOrEmpty(v.Values[0].(neo4j.Relationship).Props, "label"),
				OrganizationId: v.Values[1].(string),
				EmailId:        v.Values[2].(string),
				Tenant:         v.Values[3].(string),
			})
			if err != nil {
				failedRecords++
				if outputErr != nil {
					outputErr = err
				}
				s.log.Errorf("(organizationService.UpsertEmailRelationInEventStore) Failed to call method: {%v}", err.Error())
			} else {
				processedRecords++
			}
		}

		size -= batchSize
	}

	return processedRecords, failedRecords, outputErr
}

func (s *organizationService) UpdateLastTouchpointSync(ctx context.Context, organizationID string) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationService.UpdateLastTouchpointSync")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("organizationID", organizationID))

	if organizationID == "" {
		return
	}
	s.updateLastTouchpoint(ctx, organizationID)
}

func (s *organizationService) UpdateLastTouchpointSyncByContactId(ctx context.Context, contactID string) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationService.UpdateLastTouchpointSyncByContactId")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("contactID", contactID))

	if contactID == "" {
		return
	}
	s.updateLastTouchpointByContactId(ctx, contactID)
}

func (s *organizationService) UpdateLastTouchpointSyncByEmailId(ctx context.Context, emailID string) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationService.UpdateLastTouchpointSyncByContactId")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("emailID", emailID))

	if emailID == "" {
		return
	}
	s.updateLastTouchpointByEmailId(ctx, emailID)
}

func (s *organizationService) UpdateLastTouchpointSyncByEmail(ctx context.Context, email string) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationService.UpdateLastTouchpointSyncByEmail")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("email", email))

	if email == "" {
		return
	}
	s.updateLastTouchpointByEmail(ctx, email)
}

func (s *organizationService) UpdateLastTouchpointSyncByPhoneNumberId(ctx context.Context, phoneNumberID string) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationService.UpdateLastTouchpointSyncByPhoneNumberId")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("phoneNumberID", phoneNumberID))

	if phoneNumberID == "" {
		return
	}
	s.updateLastTouchpointByPhoneNumberId(ctx, phoneNumberID)
}

func (s *organizationService) UpdateLastTouchpointSyncByPhoneNumber(ctx context.Context, phoneNumber string) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationService.UpdateLastTouchpointSyncByPhoneNumber")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("phoneNumber", phoneNumber))

	if phoneNumber == "" {
		return
	}
	s.updateLastTouchpointByPhoneNumber(ctx, phoneNumber)
}

func (s *organizationService) updateLastTouchpointByContactId(ctx context.Context, contactID string) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationService.updateLastTouchpointByContactId")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("contactID", contactID))

	dbNodesWithTotalCount, err := s.repositories.OrganizationRepository.GetPaginatedOrganizationsForContact(ctx, common.GetTenantFromContext(ctx), contactID, 0, 1000, &utils.CypherFilter{}, &utils.CypherSort{})
	if err != nil {
		s.log.Errorf("(organizationService.updateLastTouchpointByContactId) Failed to get organizations for contact: {%v}", err.Error())
		tracing.TraceErr(span, err)
		return
	}
	for _, dbNode := range dbNodesWithTotalCount.Nodes {
		props := utils.GetPropsFromNode(*dbNode)
		orgID := utils.GetStringPropOrEmpty(props, "id")
		if orgID != "" {
			s.updateLastTouchpoint(ctx, orgID)
		}
	}
}

func (s *organizationService) updateLastTouchpointByEmailId(ctx context.Context, emailID string) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationService.updateLastTouchpointByEmailId")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("emailID", emailID))

	contactDbNodes, err := s.repositories.ContactRepository.GetAllForEmails(ctx, common.GetTenantFromContext(ctx), []string{emailID})
	if err != nil {
		s.log.Errorf("(organizationService.updateLastTouchpointByEmailId) Failed to get contacts for email: {%v}", err.Error())
		tracing.TraceErr(span, err)
		return
	}
	for _, dbNode := range contactDbNodes {
		props := utils.GetPropsFromNode(*dbNode.Node)
		contactID := utils.GetStringPropOrEmpty(props, "id")
		if contactID != "" {
			s.updateLastTouchpointByContactId(ctx, contactID)
		}
	}

	orgDbNodes, err := s.repositories.OrganizationRepository.GetAllForEmails(ctx, common.GetTenantFromContext(ctx), []string{emailID})
	if err != nil {
		s.log.Errorf("(organizationService.updateLastTouchpointByEmailId) Failed to get organizations for email: {%v}", err.Error())
		tracing.TraceErr(span, err)
		return
	}
	for _, dbNode := range orgDbNodes {
		props := utils.GetPropsFromNode(*dbNode.Node)
		orgID := utils.GetStringPropOrEmpty(props, "id")
		if orgID != "" {
			s.updateLastTouchpoint(ctx, orgID)
		}
	}
}

func (s *organizationService) updateLastTouchpointByEmail(ctx context.Context, email string) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationService.updateLastTouchpointByEmail")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("email", email))

	dbNode, err := s.repositories.EmailRepository.GetByEmail(ctx, common.GetTenantFromContext(ctx), email)
	if err != nil {
		s.log.Errorf("(organizationService.updateLastTouchpointByEmail) Failed to get email: {%v}", err.Error())
		tracing.TraceErr(span, err)
		return
	}
	props := utils.GetPropsFromNode(*dbNode)
	emailID := utils.GetStringPropOrEmpty(props, "id")
	if emailID != "" {
		s.updateLastTouchpointByEmailId(ctx, emailID)
	}
}

func (s *organizationService) updateLastTouchpointByPhoneNumberId(ctx context.Context, phoneNumberID string) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationService.updateLastTouchpointByPhoneNumberId")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("phoneNumberID", phoneNumberID))

	contactDbNodes, err := s.repositories.ContactRepository.GetAllForPhoneNumbers(ctx, common.GetTenantFromContext(ctx), []string{phoneNumberID})
	if err != nil {
		s.log.Errorf("(organizationService.updateLastTouchpointByPhoneNumberId) Failed to get contacts for phone number: {%v}", err.Error())
		tracing.TraceErr(span, err)
		return
	}
	for _, dbNode := range contactDbNodes {
		props := utils.GetPropsFromNode(*dbNode.Node)
		contactID := utils.GetStringPropOrEmpty(props, "id")
		if contactID != "" {
			s.updateLastTouchpointByContactId(ctx, contactID)
		}
	}

	orgDbNodes, err := s.repositories.OrganizationRepository.GetAllForPhoneNumbers(ctx, common.GetTenantFromContext(ctx), []string{phoneNumberID})
	if err != nil {
		s.log.Errorf("(organizationService.updateLastTouchpointByPhoneNumberId) Failed to get organizations for phone number: {%v}", err.Error())
		tracing.TraceErr(span, err)
		return
	}
	for _, dbNode := range orgDbNodes {
		props := utils.GetPropsFromNode(*dbNode.Node)
		orgID := utils.GetStringPropOrEmpty(props, "id")
		if orgID != "" {
			s.updateLastTouchpoint(ctx, orgID)
		}
	}
}

func (s *organizationService) updateLastTouchpointByPhoneNumber(ctx context.Context, phoneNumber string) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationService.updateLastTouchpointByPhoneNumber")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("phoneNumber", phoneNumber))

	dbNode, err := s.repositories.PhoneNumberRepository.GetByPhoneNumber(ctx, common.GetTenantFromContext(ctx), phoneNumber)
	if err != nil {
		s.log.Errorf("(organizationService.updateLastTouchpointByPhoneNumber) Failed to get phone number: {%v}", err.Error())
		tracing.TraceErr(span, err)
		return
	}
	props := utils.GetPropsFromNode(*dbNode)
	phoneNumberID := utils.GetStringPropOrEmpty(props, "id")
	if phoneNumberID != "" {
		s.updateLastTouchpointByPhoneNumberId(ctx, phoneNumberID)
	}
}

func (s *organizationService) updateLastTouchpoint(ctx context.Context, organizationID string) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationService.updateLastTouchpoint")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("organizationID", organizationID))

	lastTouchpointAt, lastTouchpointId, err := s.repositories.TimelineEventRepository.CalculateAndGetLastTouchpoint(ctx, common.GetTenantFromContext(ctx), organizationID)

	if err != nil {
		s.log.Errorf("(organizationService.updateLastTouchpoint) Failed to calculate last touchpoint: {%v}", err.Error())
		tracing.TraceErr(span, err)
		return
	}

	if lastTouchpointAt == nil {
		s.log.Infof("(organizationService.updateLastTouchpoint) Last touchpoint not available for organization: {%v}", organizationID)
		return
	}

	if err = s.repositories.OrganizationRepository.UpdateLastTouchpoint(ctx, common.GetTenantFromContext(ctx), organizationID, *lastTouchpointAt, lastTouchpointId); err != nil {
		s.log.Errorf("(organizationService.updateLastTouchpoint) Failed to update last touchpoint: {%v}", err.Error())
		tracing.TraceErr(span, err)
	}
}

func (s *organizationService) mapDbNodeToOrganizationEntity(node dbtype.Node) *entity.OrganizationEntity {
	props := utils.GetPropsFromNode(node)

	output := entity.OrganizationEntity{
		ID:                utils.GetStringPropOrEmpty(props, "id"),
		Name:              utils.GetStringPropOrEmpty(props, "name"),
		Description:       utils.GetStringPropOrEmpty(props, "description"),
		Website:           utils.GetStringPropOrEmpty(props, "website"),
		Industry:          utils.GetStringPropOrEmpty(props, "industry"),
		IndustryGroup:     utils.GetStringPropOrEmpty(props, "industryGroup"),
		SubIndustry:       utils.GetStringPropOrEmpty(props, "subIndustry"),
		TargetAudience:    utils.GetStringPropOrEmpty(props, "targetAudience"),
		ValueProposition:  utils.GetStringPropOrEmpty(props, "valueProposition"),
		LastFundingRound:  utils.GetStringPropOrEmpty(props, "lastFundingRound"),
		LastFundingAmount: utils.GetStringPropOrEmpty(props, "lastFundingAmount"),
		SlackChannelLink:  utils.GetStringPropOrEmpty(props, "slackChannelLink"),
		IsPublic:          utils.GetBoolPropOrFalse(props, "isPublic"),
		Employees:         utils.GetInt64PropOrZero(props, "employees"),
		Market:            utils.GetStringPropOrEmpty(props, "market"),
		CreatedAt:         utils.GetTimePropOrEpochStart(props, "createdAt"),
		UpdatedAt:         utils.GetTimePropOrEpochStart(props, "updatedAt"),
		Source:            entity.GetDataSource(utils.GetStringPropOrEmpty(props, "source")),
		SourceOfTruth:     entity.GetDataSource(utils.GetStringPropOrEmpty(props, "sourceOfTruth")),
		AppSource:         utils.GetStringPropOrEmpty(props, "appSource"),
		LastTouchpointAt:  utils.GetTimePropOrNil(props, "lastTouchpointAt"),
		LastTouchpointId:  utils.GetStringPropOrNil(props, "lastTouchpointId"),
		RenewalLikelihood: entity.RenewalLikelihood{
			RenewalLikelihood:         utils.GetStringPropOrEmpty(props, "renewalLikelihood"),
			PreviousRenewalLikelihood: utils.GetStringPropOrEmpty(props, "renewalLikelihoodPrevious"),
			Comment:                   utils.GetStringPropOrNil(props, "renewalLikelihoodComment"),
			UpdatedBy:                 utils.GetStringPropOrNil(props, "renewalLikelihoodUpdatedBy"),
			UpdatedAt:                 utils.GetTimePropOrNil(props, "renewalLikelihoodUpdatedAt"),
		},
		RenewalForecast: entity.RenewalForecast{
			Amount:         utils.GetFloatPropOrNil(props, "renewalForecast"),
			PreviousAmount: utils.GetFloatPropOrNil(props, "renewalForecastPrevious"),
			Comment:        utils.GetStringPropOrNil(props, "renewalForecastComment"),
			UpdatedBy:      utils.GetStringPropOrNil(props, "renewalForecastUpdatedBy"),
			UpdatedAt:      utils.GetTimePropOrNil(props, "renewalForecastUpdatedAt"),
		},
		BillingDetails: entity.BillingDetails{
			Amount:            utils.GetFloatPropOrNil(props, "billingDetailsAmount"),
			Frequency:         utils.GetStringPropOrEmpty(props, "billingDetailsFrequency"),
			RenewalCycle:      utils.GetStringPropOrEmpty(props, "billingDetailsRenewalCycle"),
			RenewalCycleStart: utils.GetTimePropOrNil(props, "billingDetailsRenewalCycleStart"),
		},
	}
	return &output
}

func (s *organizationService) addLinkedOrganizationRelationshipToOrganizationEntity(relationship dbtype.Relationship, organizationEntity *entity.OrganizationEntity) {
	props := utils.GetPropsFromRelationship(relationship)
	organizationEntity.LinkedOrganizationType = utils.GetStringPropOrNil(props, "type")
}

func (s *organizationService) addSuggestedMergeRelationshipToOrganizationEntity(relationship dbtype.Relationship, organizationEntity *entity.OrganizationEntity) {
	props := utils.GetPropsFromRelationship(relationship)
	organizationEntity.SuggestedMerge.SuggestedBy = utils.GetStringPropOrNil(props, "suggestedBy")
	organizationEntity.SuggestedMerge.SuggestedAt = utils.GetTimePropOrNil(props, "suggestedAt")
	organizationEntity.SuggestedMerge.Confidence = utils.GetFloatPropOrNil(props, "confidence")
}
