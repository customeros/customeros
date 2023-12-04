package service

import (
	"context"
	"fmt"
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
	organizationpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/organization"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"reflect"
)

type OrganizationService interface {
	CountOrganizations(ctx context.Context, tenant string) (int64, error)
	GetOrganizationsForJobRoles(ctx context.Context, jobRoleIds []string) (*entity.OrganizationEntities, error)
	GetById(ctx context.Context, organizationId string) (*entity.OrganizationEntity, error)
	ExistsById(ctx context.Context, organizationId string) (bool, error)
	FindAll(ctx context.Context, page, limit int, filter *model.Filter, sortBy []*model.SortBy) (*utils.Pagination, error)
	GetOrganizationsForContact(ctx context.Context, contactId string, page, limit int, filter *model.Filter, sortBy []*model.SortBy) (*utils.Pagination, error)
	Archive(ctx context.Context, organizationId string) error
	Merge(ctx context.Context, primaryOrganizationId, mergedOrganizationId string) error
	GetOrganizationsForEmails(ctx context.Context, emailIds []string) (*entity.OrganizationEntities, error)
	GetOrganizationsForPhoneNumbers(ctx context.Context, phoneNumberIds []string) (*entity.OrganizationEntities, error)
	GetSubsidiariesForOrganizations(ctx context.Context, parentOrganizationIds []string) (*entity.OrganizationEntities, error)
	GetSubsidiariesOfForOrganizations(ctx context.Context, organizationIds []string) (*entity.OrganizationEntities, error)
	AddSubsidiary(ctx context.Context, parentOrganizationId, subOrganizationId, subsidiaryType string) error
	RemoveSubsidiary(ctx context.Context, parentOrganizationId, subOrganizationId string) error
	ReplaceOwner(ctx context.Context, organizationId, userId string) (*entity.OrganizationEntity, error)
	RemoveOwner(ctx context.Context, organizationId string) (*entity.OrganizationEntity, error)
	UpdateLastTouchpoint(ctx context.Context, organizationId string)
	UpdateLastTouchpointByContactId(ctx context.Context, contactId string)
	UpdateLastTouchpointByEmailId(ctx context.Context, emailId string)
	UpdateLastTouchpointByPhoneNumberId(ctx context.Context, phoneNumberId string)
	UpdateLastTouchpointByEmail(ctx context.Context, email string)
	UpdateLastTouchpointByPhoneNumber(ctx context.Context, phoneNumber string)
	GetSuggestedMergeToForOrganizations(ctx context.Context, organizationIds []string) (*entity.OrganizationEntities, error)
	GetMinMaxRenewalForecastAmount(ctx context.Context) (float64, float64, error)
	GetMinMaxRenewalForecastArr(ctx context.Context) (float64, float64, error)

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

func (s *organizationService) CountOrganizations(ctx context.Context, tenant string) (int64, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationService.CountOrganizations")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("tenant", tenant))

	return s.repositories.OrganizationRepository.CountOrganizations(ctx, tenant)
}

func (s *organizationService) ExistsById(ctx context.Context, organizationId string) (bool, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationService.ExistsById")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("organizationId", organizationId))

	return s.repositories.CommonRepository.ExistsById(ctx, common.GetTenantFromContext(ctx), organizationId, entity.NodeLabel_Organization)
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

func (s *organizationService) GetById(ctx context.Context, organizationId string) (*entity.OrganizationEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationService.GetById")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("organizationId", organizationId))

	dbNode, err := s.repositories.OrganizationRepository.GetOrganizationById(ctx, common.GetTenantFromContext(ctx), organizationId)
	if err != nil {
		wrappedErr := errors.Wrap(err, fmt.Sprintf("Organization with id {%s} not found", organizationId))
		return nil, wrappedErr
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

	_, err := s.GetById(ctx, primaryOrganizationId)
	if err != nil {
		s.log.Errorf("(organizationService.Merge) Primary organization with id {%s} not found: {%v}", primaryOrganizationId, err.Error())
		return err
	}
	_, err = s.GetById(ctx, mergedOrganizationId)
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

	s.UpdateLastTouchpoint(ctx, primaryOrganizationId)

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

		//check if the hide bool is false and only then append the sub org to the parent org
		if !organizationEntity.Hide {
			s.addLinkedOrganizationRelationshipToOrganizationEntity(*v.Relationship, organizationEntity)
			organizationEntity.DataloaderKey = v.LinkedNodeId
			organizationEntities = append(organizationEntities, *organizationEntity)
		}
	}
	return &organizationEntities, nil
}

func (s *organizationService) AddSubsidiary(ctx context.Context, parentOrganizationId, subOrganizationId, subsidiaryType string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationService.AddSubsidiary")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("parentOrganizationId", parentOrganizationId), log.String("subOrganizationId", subOrganizationId), log.String("subsidiaryType", subsidiaryType))

	parentExists, err := s.repositories.CommonRepository.ExistsById(ctx, common.GetTenantFromContext(ctx), parentOrganizationId, entity.NodeLabel_Organization)
	if err != nil {
		s.log.Errorf("error checking if parent organization exists: {%v}", err.Error())
		tracing.TraceErr(span, err)
		return err
	}
	if !parentExists {
		err = fmt.Errorf("Parent organization with id {%s} not found", parentOrganizationId)
		s.log.Errorf("%v", err.Error())
		tracing.TraceErr(span, err)
		return err
	}

	subExists, err := s.repositories.CommonRepository.ExistsById(ctx, common.GetTenantFromContext(ctx), subOrganizationId, entity.NodeLabel_Organization)
	if err != nil {
		s.log.Errorf("error checking if sub organization exists: {%v}", err.Error())
		tracing.TraceErr(span, err)
		return err
	}
	if !subExists {
		err = fmt.Errorf("sub organization with id {%s} not found", subOrganizationId)
		s.log.Errorf("%v", err.Error())
		tracing.TraceErr(span, err)
		return err
	}

	_, err = s.grpcClients.OrganizationClient.AddParentOrganization(ctx, &organizationpb.AddParentOrganizationGrpcRequest{
		Tenant:               common.GetTenantFromContext(ctx),
		OrganizationId:       subOrganizationId,
		ParentOrganizationId: parentOrganizationId,
		Type:                 subsidiaryType,
		LoggedInUserId:       common.GetUserIdFromContext(ctx),
		AppSource:            constants.AppSourceCustomerOsApi,
	})
	if err != nil {
		s.log.Errorf("error sending event to events-platform: {%v}", err.Error())
		tracing.TraceErr(span, err, log.String("grpcMethod", "AddParentOrganization"))
	}
	return err
}

func (s *organizationService) RemoveSubsidiary(ctx context.Context, parentOrganizationId, subOrganizationId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationService.AddSubsidiary")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("parentOrganizationId", parentOrganizationId), log.String("subOrganizationId", subOrganizationId))

	parentExists, err := s.repositories.CommonRepository.ExistsById(ctx, common.GetTenantFromContext(ctx), parentOrganizationId, entity.NodeLabel_Organization)
	if err != nil {
		s.log.Errorf("error checking if parent organization exists: {%v}", err.Error())
		tracing.TraceErr(span, err)
		return err
	}
	if !parentExists {
		err = fmt.Errorf("Parent organization with id {%s} not found", parentOrganizationId)
		s.log.Errorf("%v", err.Error())
		tracing.TraceErr(span, err)
		return err
	}

	subExists, err := s.repositories.CommonRepository.ExistsById(ctx, common.GetTenantFromContext(ctx), subOrganizationId, entity.NodeLabel_Organization)
	if err != nil {
		s.log.Errorf("error checking if sub organization exists: {%v}", err.Error())
		tracing.TraceErr(span, err)
		return err
	}
	if !subExists {
		err = fmt.Errorf("sub organization with id {%s} not found", subOrganizationId)
		s.log.Errorf("%v", err.Error())
		tracing.TraceErr(span, err)
		return err
	}

	_, err = s.grpcClients.OrganizationClient.RemoveParentOrganization(ctx, &organizationpb.RemoveParentOrganizationGrpcRequest{
		Tenant:               common.GetTenantFromContext(ctx),
		OrganizationId:       subOrganizationId,
		ParentOrganizationId: parentOrganizationId,
		LoggedInUserId:       common.GetUserIdFromContext(ctx),
		AppSource:            constants.AppSourceCustomerOsApi,
	})
	if err != nil {
		s.log.Errorf("error sending event to events-platform: {%v}", err.Error())
		tracing.TraceErr(span, err, log.String("grpcMethod", "RemoveParentOrganization"))
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

func (s *organizationService) GetMinMaxRenewalForecastAmount(ctx context.Context) (float64, float64, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationService.GetMinMaxRenewalForecastAmount")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	min, max, err := s.repositories.OrganizationRepository.GetMinMaxRenewalForecastAmount(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error getting min and max renewal forecast amount: %s", err.Error())
		return 0, 0, err
	}
	return min, max, nil
}

func (s *organizationService) GetMinMaxRenewalForecastArr(ctx context.Context) (float64, float64, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationService.GetMinMaxRenewalForecastArr")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	min, max, err := s.repositories.OrganizationRepository.GetMinMaxRenewalForecastArr(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error getting min and max renewal forecast ARR: %s", err.Error())
		return 0, 0, err
	}
	if min == float64(0) && max == float64(0) {
		min, max, err = s.repositories.OrganizationRepository.GetMinMaxRenewalForecastAmount(ctx)
		if err != nil {
			tracing.TraceErr(span, err)
			s.log.Errorf("Error getting min and max renewal forecast amount: %s", err.Error())
			return 0, 0, err
		}
	}
	return min, max, nil
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
			_, err := s.grpcClients.OrganizationClient.LinkPhoneNumberToOrganization(context.Background(), &organizationpb.LinkPhoneNumberToOrganizationGrpcRequest{
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
			_, err := s.grpcClients.OrganizationClient.LinkEmailToOrganization(context.Background(), &organizationpb.LinkEmailToOrganizationGrpcRequest{
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

func (s *organizationService) UpdateLastTouchpoint(ctx context.Context, organizationID string) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationService.UpdateLastTouchpoint")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("organizationID", organizationID))

	if organizationID == "" {
		return
	}
	s.updateLastTouchpoint(ctx, organizationID)
}

func (s *organizationService) UpdateLastTouchpointByContactId(ctx context.Context, contactID string) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationService.UpdateLastTouchpointByContactId")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("contactID", contactID))

	if contactID == "" {
		return
	}
	s.updateLastTouchpointByContactId(ctx, contactID)
}

func (s *organizationService) UpdateLastTouchpointByEmailId(ctx context.Context, emailID string) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationService.UpdateLastTouchpointByContactId")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("emailID", emailID))

	if emailID == "" {
		return
	}
	s.updateLastTouchpointByEmailId(ctx, emailID)
}

func (s *organizationService) UpdateLastTouchpointByEmail(ctx context.Context, email string) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationService.UpdateLastTouchpointByEmail")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("email", email))

	if email == "" {
		return
	}
	s.updateLastTouchpointByEmail(ctx, email)
}

func (s *organizationService) UpdateLastTouchpointByPhoneNumberId(ctx context.Context, phoneNumberID string) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationService.UpdateLastTouchpointByPhoneNumberId")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("phoneNumberID", phoneNumberID))

	if phoneNumberID == "" {
		return
	}
	s.updateLastTouchpointByPhoneNumberId(ctx, phoneNumberID)
}

func (s *organizationService) UpdateLastTouchpointByPhoneNumber(ctx context.Context, phoneNumber string) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationService.UpdateLastTouchpointByPhoneNumber")
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

	_, err := s.grpcClients.OrganizationClient.RefreshLastTouchpoint(ctx, &organizationpb.OrganizationIdGrpcRequest{
		Tenant:         common.GetTenantFromContext(ctx),
		OrganizationId: organizationID,
		LoggedInUserId: common.GetUserIdFromContext(ctx),
		AppSource:      constants.AppSourceCustomerOsApi,
	})
	if err != nil {
		s.log.Errorf("error sending event to events-platform: {%v}", err.Error())
		tracing.TraceErr(span, err, log.String("grpcMethod", "RefreshLastTouchpoint"))
	}
}

func (s *organizationService) mapDbNodeToOrganizationEntity(node dbtype.Node) *entity.OrganizationEntity {
	props := utils.GetPropsFromNode(node)

	output := entity.OrganizationEntity{
		ID:                 utils.GetStringPropOrEmpty(props, "id"),
		CustomerOsId:       utils.GetStringPropOrEmpty(props, "customerOsId"),
		ReferenceId:        utils.GetStringPropOrEmpty(props, "referenceId"),
		Name:               utils.GetStringPropOrEmpty(props, "name"),
		Description:        utils.GetStringPropOrEmpty(props, "description"),
		Website:            utils.GetStringPropOrEmpty(props, "website"),
		Industry:           utils.GetStringPropOrEmpty(props, "industry"),
		IndustryGroup:      utils.GetStringPropOrEmpty(props, "industryGroup"),
		SubIndustry:        utils.GetStringPropOrEmpty(props, "subIndustry"),
		TargetAudience:     utils.GetStringPropOrEmpty(props, "targetAudience"),
		ValueProposition:   utils.GetStringPropOrEmpty(props, "valueProposition"),
		LastFundingRound:   utils.GetStringPropOrEmpty(props, "lastFundingRound"),
		LastFundingAmount:  utils.GetStringPropOrEmpty(props, "lastFundingAmount"),
		Note:               utils.GetStringPropOrEmpty(props, "note"),
		IsPublic:           utils.GetBoolPropOrFalse(props, "isPublic"),
		IsCustomer:         utils.GetBoolPropOrFalse(props, "isCustomer"),
		Hide:               utils.GetBoolPropOrFalse(props, "hide"),
		Employees:          utils.GetInt64PropOrZero(props, "employees"),
		Market:             utils.GetStringPropOrEmpty(props, "market"),
		CreatedAt:          utils.GetTimePropOrEpochStart(props, "createdAt"),
		UpdatedAt:          utils.GetTimePropOrEpochStart(props, "updatedAt"),
		Source:             entity.GetDataSource(utils.GetStringPropOrEmpty(props, "source")),
		SourceOfTruth:      entity.GetDataSource(utils.GetStringPropOrEmpty(props, "sourceOfTruth")),
		AppSource:          utils.GetStringPropOrEmpty(props, "appSource"),
		LastTouchpointId:   utils.GetStringPropOrNil(props, "lastTouchpointId"),
		LastTouchpointAt:   utils.GetTimePropOrNil(props, "lastTouchpointAt"),
		LastTouchpointType: utils.GetStringPropOrNil(props, "lastTouchpointType"),
		RenewalSummary: entity.RenewalSummary{
			ArrForecast:            utils.GetFloatPropOrNil(props, "renewalForecastArr"),
			MaxArrForecast:         utils.GetFloatPropOrNil(props, "renewalForecastMaxArr"),
			NextRenewalAt:          utils.GetTimePropOrNil(props, "derivedNextRenewalAt"),
			RenewalLikelihood:      utils.GetStringPropOrEmpty(props, "derivedRenewalLikelihood"),
			RenewalLikelihoodOrder: utils.GetInt64PropOrNil(props, "derivedRenewalLikelihoodOrder"),
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
