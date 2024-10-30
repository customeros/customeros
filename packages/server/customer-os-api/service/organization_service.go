package service

import (
	"context"
	"fmt"
	model2 "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	"reflect"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/grpc_client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	organizationpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/organization"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
)

type OrganizationService interface {
	CountOrganizations(ctx context.Context, tenant string) (int64, error)
	GetOrganizationsForJobRoles(ctx context.Context, jobRoleIds []string) (*neo4jentity.OrganizationEntities, error)
	GetOrganizationsForInvoices(ctx context.Context, invoiceIds []string) (*neo4jentity.OrganizationEntities, error)
	GetOrganizationsForSlackChannels(ctx context.Context, slackChannelIds []string) (*neo4jentity.OrganizationEntities, error)
	GetOrganizationsForOpportunities(ctx context.Context, opportunityIds []string) (*neo4jentity.OrganizationEntities, error)
	GetByCustomerOsId(ctx context.Context, customerOsId string) (*neo4jentity.OrganizationEntity, error)
	GetByReferenceId(ctx context.Context, referenceId string) (*neo4jentity.OrganizationEntity, error)
	ExistsById(ctx context.Context, organizationId string) (bool, error)
	FindAll(ctx context.Context, page, limit int, filter *model.Filter, sortBy []*model.SortBy) (*utils.Pagination, error)
	GetOrganizationsForContact(ctx context.Context, contactId string, page, limit int, filter *model.Filter, sortBy []*model.SortBy) (*utils.Pagination, error)
	Merge(ctx context.Context, primaryOrganizationId, mergedOrganizationId string) error
	GetOrganizationsForEmails(ctx context.Context, emailIds []string) (*neo4jentity.OrganizationEntities, error)
	GetOrganizationsForPhoneNumbers(ctx context.Context, phoneNumberIds []string) (*neo4jentity.OrganizationEntities, error)
	GetSubsidiariesForOrganizations(ctx context.Context, parentOrganizationIds []string) (*neo4jentity.OrganizationEntities, error)
	GetSubsidiariesOfForOrganizations(ctx context.Context, organizationIds []string) (*neo4jentity.OrganizationEntities, error)
	AddSubsidiary(ctx context.Context, parentOrganizationId, subsidiaryOrganizationId, subsidiaryType string, removeExisting bool) error
	RemoveSubsidiary(ctx context.Context, parentOrganizationId, subsidiaryOrganizationId string) error
	ReplaceOwner(ctx context.Context, organizationId, userId string) (*neo4jentity.OrganizationEntity, error)
	RemoveOwner(ctx context.Context, organizationId string) (*neo4jentity.OrganizationEntity, error)
	UpdateLastTouchpoint(ctx context.Context, organizationId string)
	UpdateLastTouchpointByContactId(ctx context.Context, contactId string)
	UpdateLastTouchpointByEmailId(ctx context.Context, emailId string)
	UpdateLastTouchpointByPhoneNumberId(ctx context.Context, phoneNumberId string)
	UpdateLastTouchpointByEmail(ctx context.Context, email string)
	UpdateLastTouchpointByPhoneNumber(ctx context.Context, phoneNumber string)
	GetSuggestedMergeToForOrganizations(ctx context.Context, organizationIds []string) (*neo4jentity.OrganizationEntities, error)
	GetMinMaxRenewalForecastArr(ctx context.Context) (float64, float64, error)
	GetOrganizations(ctx context.Context, organizationIds []string) (*neo4jentity.OrganizationEntities, error)
}

type OrganizationUpdateData struct {
	Organization *neo4jentity.OrganizationEntity
	Domains      []string
}

type organizationService struct {
	log          logger.Logger
	repositories *repository.Repositories
	grpcClients  *grpc_client.Clients
	services     *Services
}

func NewOrganizationService(log logger.Logger, repositories *repository.Repositories, grpcClients *grpc_client.Clients, services *Services) OrganizationService {
	return &organizationService{
		log:          log,
		repositories: repositories,
		grpcClients:  grpcClients,
		services:     services,
	}
}

func (s *organizationService) CountOrganizations(ctx context.Context, tenant string) (int64, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationService.CountOrganizations")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.SetTag(tracing.SpanTagTenant, tenant)

	return s.repositories.Neo4jRepositories.OrganizationReadRepository.CountByTenant(ctx, tenant)
}

func (s *organizationService) ExistsById(ctx context.Context, organizationId string) (bool, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationService.ExistsById")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("organizationId", organizationId))

	return s.repositories.Neo4jRepositories.CommonReadRepository.ExistsById(ctx, common.GetTenantFromContext(ctx), organizationId, model2.NodeLabelOrganization)
}

func (s *organizationService) FindAll(ctx context.Context, page, limit int, filter *model.Filter, sortBy []*model.SortBy) (*utils.Pagination, error) {
	var paginatedResult = utils.Pagination{
		Limit: limit,
		Page:  page,
	}
	cypherSort, err := buildSort(sortBy, reflect.TypeOf(neo4jentity.OrganizationEntity{}))
	if err != nil {
		return nil, err
	}
	cypherFilter, err := buildFilter(filter, reflect.TypeOf(neo4jentity.OrganizationEntity{}))
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

	organizationEntities := make(neo4jentity.OrganizationEntities, 0, len(dbNodesWithTotalCount.Nodes))
	for _, v := range dbNodesWithTotalCount.Nodes {
		organizationEntities = append(organizationEntities, *neo4jmapper.MapDbNodeToOrganizationEntity(v))
	}
	paginatedResult.SetRows(&organizationEntities)
	return &paginatedResult, nil
}

func (s *organizationService) GetOrganizationsForContact(ctx context.Context, contactId string, page, limit int, filter *model.Filter, sortBy []*model.SortBy) (*utils.Pagination, error) {
	var paginatedResult = utils.Pagination{
		Limit: limit,
		Page:  page,
	}
	cypherSort, err := buildSort(sortBy, reflect.TypeOf(neo4jentity.OrganizationEntity{}))
	if err != nil {
		return nil, err
	}
	cypherFilter, err := buildFilter(filter, reflect.TypeOf(neo4jentity.OrganizationEntity{}))
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

	organizationEntities := make(neo4jentity.OrganizationEntities, 0, len(dbNodesWithTotalCount.Nodes))
	for _, v := range dbNodesWithTotalCount.Nodes {
		organizationEntities = append(organizationEntities, *neo4jmapper.MapDbNodeToOrganizationEntity(v))
	}
	paginatedResult.SetRows(&organizationEntities)
	return &paginatedResult, nil
}

func (s *organizationService) GetOrganizationsForJobRoles(ctx context.Context, jobRoleIds []string) (*neo4jentity.OrganizationEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationService.GetOrganizationsForJobRoles")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("jobRoleIds", jobRoleIds))

	organizations, err := s.repositories.OrganizationRepository.GetAllForJobRoles(ctx, common.GetTenantFromContext(ctx), jobRoleIds)
	if err != nil {
		return nil, err
	}
	organizationEntities := make(neo4jentity.OrganizationEntities, 0, len(organizations))
	for _, v := range organizations {
		organizationEntity := neo4jmapper.MapDbNodeToOrganizationEntity(v.Node)
		organizationEntity.DataloaderKey = v.LinkedNodeId
		organizationEntities = append(organizationEntities, *organizationEntity)
	}
	return &organizationEntities, nil
}

func (s *organizationService) GetByCustomerOsId(ctx context.Context, customerOsId string) (*neo4jentity.OrganizationEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationService.GetByCustomerOsId")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("customerOsId", customerOsId))

	dbNode, err := s.repositories.Neo4jRepositories.OrganizationReadRepository.GetOrganizationByCustomerOsId(ctx, common.GetTenantFromContext(ctx), customerOsId)
	if err != nil {
		wrappedErr := errors.Wrap(err, fmt.Sprintf("Organization with customerOsId {%s} not found", customerOsId))
		return nil, wrappedErr
	}
	if dbNode == nil {
		return nil, nil
	}
	return neo4jmapper.MapDbNodeToOrganizationEntity(dbNode), nil
}

func (s *organizationService) GetByReferenceId(ctx context.Context, referenceId string) (*neo4jentity.OrganizationEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationService.GetByReferenceId")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("referenceId", referenceId))

	dbNode, err := s.repositories.Neo4jRepositories.OrganizationReadRepository.GetOrganizationByReferenceId(ctx, common.GetTenantFromContext(ctx), referenceId)
	if err != nil {
		wrappedErr := errors.Wrap(err, fmt.Sprintf("Organization with customerOsId {%s} not found", referenceId))
		return nil, wrappedErr
	}
	if dbNode == nil {
		return nil, nil
	}
	return neo4jmapper.MapDbNodeToOrganizationEntity(dbNode), nil
}

func (s *organizationService) Merge(ctx context.Context, primaryOrganizationId, mergedOrganizationId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationService.Merge")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("primaryOrganizationId", primaryOrganizationId), log.String("mergedOrganizationId", mergedOrganizationId))

	tenant := common.GetTenantFromContext(ctx)

	_, err := s.services.CommonServices.OrganizationService.GetById(ctx, tenant, primaryOrganizationId)
	if err != nil {
		s.log.Errorf("(organizationService.Merge) Primary organization with id {%s} not found: {%v}", primaryOrganizationId, err.Error())
		return err
	}
	_, err = s.services.CommonServices.OrganizationService.GetById(ctx, tenant, mergedOrganizationId)
	if err != nil {
		s.log.Errorf("(organizationService.Merge) Organization to merge with id {%s} not found: {%v}", mergedOrganizationId, err.Error())
		return err
	}

	session := utils.NewNeo4jWriteSession(ctx, *s.repositories.Drivers.Neo4jDriver)
	defer session.Close(ctx)

	_, err = session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		err = s.repositories.OrganizationRepository.MergeOrganizationPropertiesInTx(ctx, tx, tenant, primaryOrganizationId, mergedOrganizationId, neo4jentity.DataSourceOpenline)
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

	// Update last touchpoint
	s.UpdateLastTouchpoint(ctx, primaryOrganizationId)

	// Refresh forecast ARR
	ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
	_, err = utils.CallEventsPlatformGRPCWithRetry[*organizationpb.OrganizationIdGrpcResponse](func() (*organizationpb.OrganizationIdGrpcResponse, error) {
		return s.grpcClients.OrganizationClient.RefreshArr(ctx, &organizationpb.OrganizationIdGrpcRequest{
			Tenant:         common.GetTenantFromContext(ctx),
			OrganizationId: primaryOrganizationId,
			LoggedInUserId: common.GetUserIdFromContext(ctx),
			AppSource:      constants.AppSourceCustomerOsApi,
		})
	})
	if err != nil {
		s.log.Errorf("error sending event to events-platform: {%s}", err.Error())
		tracing.TraceErr(span, err, log.String("grpcMethod", "RefreshArr"))
	}

	// Refresh renewal likelihood
	ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
	_, err = utils.CallEventsPlatformGRPCWithRetry[*organizationpb.OrganizationIdGrpcResponse](func() (*organizationpb.OrganizationIdGrpcResponse, error) {
		return s.grpcClients.OrganizationClient.RefreshRenewalSummary(ctx, &organizationpb.RefreshRenewalSummaryGrpcRequest{
			Tenant:         common.GetTenantFromContext(ctx),
			OrganizationId: primaryOrganizationId,
			LoggedInUserId: common.GetUserIdFromContext(ctx),
			AppSource:      constants.AppSourceCustomerOsApi,
		})
	})
	if err != nil {
		s.log.Errorf("error sending event to events-platform: {%s}", err.Error())
		tracing.TraceErr(span, err, log.String("grpcMethod", "RefreshRenewalSummary"))
	}

	return err
}

func (s *organizationService) GetOrganizationsForEmails(ctx context.Context, emailIds []string) (*neo4jentity.OrganizationEntities, error) {
	organizations, err := s.repositories.OrganizationRepository.GetAllForEmails(ctx, common.GetTenantFromContext(ctx), emailIds)
	if err != nil {
		return nil, err
	}
	organizationEntities := make(neo4jentity.OrganizationEntities, 0, len(organizations))
	for _, v := range organizations {
		organizationEntity := neo4jmapper.MapDbNodeToOrganizationEntity(v.Node)
		organizationEntity.DataloaderKey = v.LinkedNodeId
		organizationEntities = append(organizationEntities, *organizationEntity)
	}
	return &organizationEntities, nil
}

func (s *organizationService) GetOrganizationsForPhoneNumbers(ctx context.Context, phoneNumberIds []string) (*neo4jentity.OrganizationEntities, error) {
	organizations, err := s.repositories.OrganizationRepository.GetAllForPhoneNumbers(ctx, common.GetTenantFromContext(ctx), phoneNumberIds)
	if err != nil {
		return nil, err
	}
	organizationEntities := make(neo4jentity.OrganizationEntities, 0, len(organizations))
	for _, v := range organizations {
		organizationEntity := neo4jmapper.MapDbNodeToOrganizationEntity(v.Node)
		organizationEntity.DataloaderKey = v.LinkedNodeId
		organizationEntities = append(organizationEntities, *organizationEntity)
	}
	return &organizationEntities, nil
}

func (s *organizationService) GetSubsidiariesForOrganizations(ctx context.Context, parentOrganizationIds []string) (*neo4jentity.OrganizationEntities, error) {
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
	organizationEntities := make(neo4jentity.OrganizationEntities, 0, len(dbEntries))
	for _, v := range dbEntries {
		organizationEntity := neo4jmapper.MapDbNodeToOrganizationEntity(v.Node)

		//check if the hide bool is false and only then append the sub org to the parent org
		if !organizationEntity.Hide {
			s.addLinkedOrganizationRelationshipToOrganizationEntity(*v.Relationship, organizationEntity)
			organizationEntity.DataloaderKey = v.LinkedNodeId
			organizationEntities = append(organizationEntities, *organizationEntity)
		}
	}
	return &organizationEntities, nil
}

func (s *organizationService) AddSubsidiary(ctx context.Context, parentOrganizationId, subsidiaryOrganizationId, subsidiaryType string, removeExisting bool) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationService.AddSubsidiary")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(
		log.String("parentOrganizationId", parentOrganizationId),
		log.String("subsidiaryOrganizationId", subsidiaryOrganizationId),
		log.String("subsidiaryType", subsidiaryType),
		log.Bool("removeExisting", removeExisting))

	parentExists, err := s.repositories.Neo4jRepositories.CommonReadRepository.ExistsById(ctx, common.GetTenantFromContext(ctx), parentOrganizationId, model2.NodeLabelOrganization)
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

	subExists, err := s.repositories.Neo4jRepositories.CommonReadRepository.ExistsById(ctx, common.GetTenantFromContext(ctx), subsidiaryOrganizationId, model2.NodeLabelOrganization)
	if err != nil {
		s.log.Errorf("error checking if sub organization exists: {%v}", err.Error())
		tracing.TraceErr(span, err)
		return err
	}
	if !subExists {
		err = fmt.Errorf("subsidiary organization with id {%s} not found", subsidiaryOrganizationId)
		s.log.Errorf("%v", err.Error())
		tracing.TraceErr(span, err)
		return err
	}

	if removeExisting {
		// fetch existing subsidiaries of the parent organization
		existingSubsidiaries, err := s.GetSubsidiariesForOrganizations(ctx, []string{parentOrganizationId})
		if err != nil {
			tracing.TraceErr(span, err)
			s.log.Errorf("error fetching existing subsidiaries: {%v}", err.Error())
			return err
		}
		for _, subsidiary := range *existingSubsidiaries {
			if subsidiary.ID != subsidiaryOrganizationId {
				err = s.RemoveSubsidiary(ctx, parentOrganizationId, subsidiary.ID)
				if err != nil {
					tracing.TraceErr(span, err)
					s.log.Errorf("error removing existing subsidiary: {%v}", err.Error())
				}
			}
		}
	}

	ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
	_, err = utils.CallEventsPlatformGRPCWithRetry[*organizationpb.OrganizationIdGrpcResponse](func() (*organizationpb.OrganizationIdGrpcResponse, error) {
		return s.grpcClients.OrganizationClient.AddParentOrganization(ctx, &organizationpb.AddParentOrganizationGrpcRequest{
			Tenant:               common.GetTenantFromContext(ctx),
			OrganizationId:       subsidiaryOrganizationId,
			ParentOrganizationId: parentOrganizationId,
			Type:                 subsidiaryType,
			LoggedInUserId:       common.GetUserIdFromContext(ctx),
			AppSource:            constants.AppSourceCustomerOsApi,
		})
	})
	if err != nil {
		s.log.Errorf("error sending event to events-platform: {%v}", err.Error())
		tracing.TraceErr(span, err, log.String("grpcMethod", "AddParentOrganization"))
	}
	return err
}

func (s *organizationService) RemoveSubsidiary(ctx context.Context, parentOrganizationId, subsidiaryOrganizationId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationService.AddSubsidiary")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("parentOrganizationId", parentOrganizationId), log.String("subsidiaryOrganizationId", subsidiaryOrganizationId))

	parentExists, err := s.repositories.Neo4jRepositories.CommonReadRepository.ExistsById(ctx, common.GetTenantFromContext(ctx), parentOrganizationId, model2.NodeLabelOrganization)
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

	subExists, err := s.repositories.Neo4jRepositories.CommonReadRepository.ExistsById(ctx, common.GetTenantFromContext(ctx), subsidiaryOrganizationId, model2.NodeLabelOrganization)
	if err != nil {
		s.log.Errorf("error checking if sub organization exists: {%v}", err.Error())
		tracing.TraceErr(span, err)
		return err
	}
	if !subExists {
		err = fmt.Errorf("sub organization with id {%s} not found", subsidiaryOrganizationId)
		s.log.Errorf("%v", err.Error())
		tracing.TraceErr(span, err)
		return err
	}

	ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
	_, err = utils.CallEventsPlatformGRPCWithRetry[*organizationpb.OrganizationIdGrpcResponse](func() (*organizationpb.OrganizationIdGrpcResponse, error) {
		return s.grpcClients.OrganizationClient.RemoveParentOrganization(ctx, &organizationpb.RemoveParentOrganizationGrpcRequest{
			Tenant:               common.GetTenantFromContext(ctx),
			OrganizationId:       subsidiaryOrganizationId,
			ParentOrganizationId: parentOrganizationId,
			LoggedInUserId:       common.GetUserIdFromContext(ctx),
			AppSource:            constants.AppSourceCustomerOsApi,
		})
	})
	if err != nil {
		s.log.Errorf("error sending event to events-platform: {%v}", err.Error())
		tracing.TraceErr(span, err, log.String("grpcMethod", "RemoveParentOrganization"))
	}
	return err
}

func (s *organizationService) GetSubsidiariesOfForOrganizations(ctx context.Context, organizationIds []string) (*neo4jentity.OrganizationEntities, error) {
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
	organizationEntities := make(neo4jentity.OrganizationEntities, 0, len(dbEntries))
	for _, v := range dbEntries {
		organizationEntity := neo4jmapper.MapDbNodeToOrganizationEntity(v.Node)
		s.addLinkedOrganizationRelationshipToOrganizationEntity(*v.Relationship, organizationEntity)
		organizationEntity.DataloaderKey = v.LinkedNodeId
		organizationEntities = append(organizationEntities, *organizationEntity)
	}
	return &organizationEntities, nil
}

func (s *organizationService) GetSuggestedMergeToForOrganizations(ctx context.Context, organizationIds []string) (*neo4jentity.OrganizationEntities, error) {
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
	organizationEntities := make(neo4jentity.OrganizationEntities, 0, len(dbEntries))
	for _, v := range dbEntries {
		organizationEntity := neo4jmapper.MapDbNodeToOrganizationEntity(v.Node)
		s.addSuggestedMergeRelationshipToOrganizationEntity(*v.Relationship, organizationEntity)
		organizationEntity.DataLoaderKey.DataloaderKey = v.LinkedNodeId
		organizationEntities = append(organizationEntities, *organizationEntity)
	}
	return &organizationEntities, nil
}

func (s *organizationService) GetMinMaxRenewalForecastArr(ctx context.Context) (float64, float64, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationService.GetMinMaxRenewalForecastArr")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	minArr, maxArr, err := s.repositories.OrganizationRepository.GetMinMaxRenewalForecastArr(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error getting min and max renewal forecast ARR: %s", err.Error())
		return 0, 0, err
	}
	return minArr, maxArr, nil
}

func (s *organizationService) ReplaceOwner(ctx context.Context, organizationID, userID string) (*neo4jentity.OrganizationEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationService.ReplaceOwner")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("organizationID", organizationID), log.String("userID", userID))

	ownerUpdateReq := &organizationpb.UpdateOrganizationOwnerGrpcRequest{
		Tenant:         common.GetTenantFromContext(ctx),
		OrganizationId: organizationID,
		OwnerUserId:    userID,
		LoggedInUserId: common.GetUserIdFromContext(ctx),
		AppSource:      constants.AppSourceCustomerOsApi,
	}

	ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
	_, err := utils.CallEventsPlatformGRPCWithRetry[*organizationpb.OrganizationIdGrpcResponse](func() (*organizationpb.OrganizationIdGrpcResponse, error) {
		return s.grpcClients.OrganizationClient.UpdateOrganizationOwner(ctx, ownerUpdateReq)
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	// get org node back from db to keep function signature the same
	dbNode, err := s.repositories.Neo4jRepositories.OrganizationReadRepository.GetOrganization(ctx, common.GetTenantFromContext(ctx), organizationID)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	return neo4jmapper.MapDbNodeToOrganizationEntity(dbNode), nil
}

func (s *organizationService) RemoveOwner(ctx context.Context, organizationID string) (*neo4jentity.OrganizationEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationService.RemoveOwner")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("organizationID", organizationID))

	dbNode, err := s.repositories.OrganizationRepository.RemoveOwner(ctx, common.GetTenantFromContext(ctx), organizationID)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	return neo4jmapper.MapDbNodeToOrganizationEntity(dbNode), nil
}

func (s *organizationService) UpdateLastTouchpoint(ctx context.Context, organizationID string) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationService.UpdateLastTouchpoint")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("organizationID", organizationID))

	if organizationID == "" {
		return
	}
	_ = s.services.CommonServices.OrganizationService.RequestRefreshLastTouchpoint(ctx, organizationID)
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
			_ = s.services.CommonServices.OrganizationService.RequestRefreshLastTouchpoint(ctx, orgID)
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
			_ = s.services.CommonServices.OrganizationService.RequestRefreshLastTouchpoint(ctx, orgID)
		}
	}
}

func (s *organizationService) updateLastTouchpointByEmail(ctx context.Context, email string) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationService.updateLastTouchpointByEmail")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("email", email))

	dbNode, err := s.repositories.Neo4jRepositories.EmailReadRepository.GetFirstByEmail(ctx, common.GetTenantFromContext(ctx), email)
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
			_ = s.services.CommonServices.OrganizationService.RequestRefreshLastTouchpoint(ctx, orgID)
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

func (s *organizationService) GetOrganizations(parentCtx context.Context, organizationIds []string) (*neo4jentity.OrganizationEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(parentCtx, "OrganizationService.GetOrganizations")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("organizationIds", organizationIds))

	organizationDbNodes, err := s.repositories.OrganizationRepository.GetOrganizations(ctx, common.GetTenantFromContext(ctx), organizationIds)
	if err != nil {
		return nil, err
	}
	organizationEntities := make(neo4jentity.OrganizationEntities, 0, len(organizationDbNodes))
	for _, dbNode := range organizationDbNodes {
		organizationEntity := neo4jmapper.MapDbNodeToOrganizationEntity(dbNode)
		organizationEntities = append(organizationEntities, *organizationEntity)
	}
	return &organizationEntities, nil
}

func (s *organizationService) addLinkedOrganizationRelationshipToOrganizationEntity(relationship dbtype.Relationship, organizationEntity *neo4jentity.OrganizationEntity) {
	props := utils.GetPropsFromRelationship(relationship)
	organizationEntity.LinkedOrganizationType = utils.GetStringPropOrNil(props, "type")
}

func (s *organizationService) addSuggestedMergeRelationshipToOrganizationEntity(relationship dbtype.Relationship, organizationEntity *neo4jentity.OrganizationEntity) {
	props := utils.GetPropsFromRelationship(relationship)
	organizationEntity.SuggestedMerge.SuggestedBy = utils.GetStringPropOrNil(props, "suggestedBy")
	organizationEntity.SuggestedMerge.SuggestedAt = utils.GetTimePropOrNil(props, "suggestedAt")
	organizationEntity.SuggestedMerge.Confidence = utils.GetFloatPropOrNil(props, "confidence")
}

func (s *organizationService) GetOrganizationsForInvoices(ctx context.Context, invoiceIds []string) (*neo4jentity.OrganizationEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationService.GetOrganizationsForInvoices")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("invoiceIds", invoiceIds))

	organizations, err := s.repositories.Neo4jRepositories.OrganizationReadRepository.GetAllForInvoices(ctx, common.GetTenantFromContext(ctx), invoiceIds)
	if err != nil {
		return nil, err
	}
	organizationEntities := make(neo4jentity.OrganizationEntities, 0, len(organizations))
	for _, v := range organizations {
		organizationEntity := neo4jmapper.MapDbNodeToOrganizationEntity(v.Node)
		organizationEntity.DataloaderKey = v.LinkedNodeId
		organizationEntities = append(organizationEntities, *organizationEntity)
	}
	return &organizationEntities, nil
}

func (s *organizationService) GetOrganizationsForSlackChannels(ctx context.Context, slackChannelIds []string) (*neo4jentity.OrganizationEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationService.GetOrganizationsForSlackChannels")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("slackChannelIds", slackChannelIds))

	organizations, err := s.repositories.Neo4jRepositories.OrganizationReadRepository.GetAllForSlackChannels(ctx, common.GetTenantFromContext(ctx), slackChannelIds)
	if err != nil {
		return nil, err
	}
	organizationEntities := make(neo4jentity.OrganizationEntities, 0, len(organizations))
	for _, v := range organizations {
		organizationEntity := neo4jmapper.MapDbNodeToOrganizationEntity(v.Node)
		organizationEntity.DataloaderKey = v.LinkedNodeId
		organizationEntities = append(organizationEntities, *organizationEntity)
	}
	return &organizationEntities, nil
}

func (s *organizationService) GetOrganizationsForOpportunities(ctx context.Context, opportunityIds []string) (*neo4jentity.OrganizationEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationService.GetOrganizationsForOpportunities")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("opportunityIds", opportunityIds))

	organizations, err := s.repositories.Neo4jRepositories.OrganizationReadRepository.GetAllForOpportunities(ctx, common.GetTenantFromContext(ctx), opportunityIds)
	if err != nil {
		return nil, err
	}
	organizationEntities := make(neo4jentity.OrganizationEntities, 0, len(organizations))
	for _, v := range organizations {
		organizationEntity := neo4jmapper.MapDbNodeToOrganizationEntity(v.Node)
		organizationEntity.DataloaderKey = v.LinkedNodeId
		organizationEntities = append(organizationEntities, *organizationEntity)
	}
	return &organizationEntities, nil
}
