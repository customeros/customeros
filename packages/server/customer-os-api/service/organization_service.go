package service

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils"
	"reflect"
)

type OrganizationService interface {
	Create(ctx context.Context, input *entity.OrganizationEntity) (*entity.OrganizationEntity, error)
	Update(ctx context.Context, input *entity.OrganizationEntity) (*entity.OrganizationEntity, error)
	GetOrganizationForRole(ctx context.Context, roleId string) (*entity.OrganizationEntity, error)
	GetOrganizationById(ctx context.Context, organizationId string) (*entity.OrganizationEntity, error)
	FindAll(ctx context.Context, page, limit int, filter *model.Filter, sortBy []*model.SortBy) (*utils.Pagination, error)
	PermanentDelete(ctx context.Context, organizationId string) (bool, error)
}

type organizationService struct {
	repositories *repository.Repositories
}

func NewOrganizationService(repositories *repository.Repositories) OrganizationService {
	return &organizationService{
		repositories: repositories,
	}
}

func (s *organizationService) Create(ctx context.Context, input *entity.OrganizationEntity) (*entity.OrganizationEntity, error) {
	session := utils.NewNeo4jWriteSession(*s.repositories.Drivers.Neo4jDriver)
	defer session.Close()

	dbNodePtr, err := s.repositories.OrganizationRepository.Create(session, common.GetContext(ctx).Tenant, *input)
	if err != nil {
		return nil, err
	}
	return s.mapDbNodeToOrganizationEntity(*dbNodePtr), nil
}

func (s *organizationService) Update(ctx context.Context, input *entity.OrganizationEntity) (*entity.OrganizationEntity, error) {
	session := utils.NewNeo4jWriteSession(*s.repositories.Drivers.Neo4jDriver)
	defer session.Close()

	dbNodePtr, err := s.repositories.OrganizationRepository.Update(session, common.GetContext(ctx).Tenant, *input)
	if err != nil {
		return nil, err
	}
	return s.mapDbNodeToOrganizationEntity(*dbNodePtr), nil
}

func (s *organizationService) FindAll(ctx context.Context, page, limit int, filter *model.Filter, sortBy []*model.SortBy) (*utils.Pagination, error) {
	session := utils.NewNeo4jReadSession(*s.repositories.Drivers.Neo4jDriver)
	defer session.Close()

	var paginatedResult = utils.Pagination{
		Limit: limit,
		Page:  page,
	}
	cypherSort, err := buildSort(sortBy, reflect.TypeOf(entity.ContactGroupEntity{}))
	if err != nil {
		return nil, err
	}
	cypherFilter, err := buildFilter(filter, reflect.TypeOf(entity.ContactGroupEntity{}))
	if err != nil {
		return nil, err
	}

	dbNodesWithTotalCount, err := s.repositories.OrganizationRepository.GetPaginatedOrganizations(
		session,
		common.GetContext(ctx).Tenant,
		paginatedResult.GetSkip(),
		paginatedResult.GetLimit(),
		cypherFilter,
		cypherSort)
	if err != nil {
		return nil, err
	}
	paginatedResult.SetTotalRows(dbNodesWithTotalCount.Count)

	organizationEntities := entity.OrganizationEntities{}

	for _, v := range dbNodesWithTotalCount.Nodes {
		organizationEntities = append(organizationEntities, *s.mapDbNodeToOrganizationEntity(*v))
	}
	paginatedResult.SetRows(&organizationEntities)
	return &paginatedResult, nil
}

func (s *organizationService) GetOrganizationForRole(ctx context.Context, roleId string) (*entity.OrganizationEntity, error) {
	session := utils.NewNeo4jReadSession(*s.repositories.Drivers.Neo4jDriver)
	defer session.Close()

	dbNode, err := s.repositories.OrganizationRepository.GetOrganizationForRole(session, common.GetContext(ctx).Tenant, roleId)
	if dbNode == nil || err != nil {
		return nil, err
	}
	return s.mapDbNodeToOrganizationEntity(*dbNode), nil
}

func (s *organizationService) GetOrganizationById(ctx context.Context, organizationId string) (*entity.OrganizationEntity, error) {
	session := utils.NewNeo4jReadSession(*s.repositories.Drivers.Neo4jDriver)
	defer session.Close()

	dbNode, err := s.repositories.OrganizationRepository.GetOrganizationById(session, common.GetContext(ctx).Tenant, organizationId)
	if err != nil {
		return nil, err
	}
	return s.mapDbNodeToOrganizationEntity(*dbNode), nil
}

func (s *organizationService) PermanentDelete(ctx context.Context, organizationId string) (bool, error) {
	session := utils.NewNeo4jWriteSession(*s.repositories.Drivers.Neo4jDriver)
	defer session.Close()

	err := s.repositories.OrganizationRepository.Delete(session, common.GetContext(ctx).Tenant, organizationId)

	if err != nil {
		return false, err
	}

	return true, nil
}

func (s *organizationService) mapDbNodeToOrganizationEntity(node dbtype.Node) *entity.OrganizationEntity {
	props := utils.GetPropsFromNode(node)
	organizationEntityPtr := new(entity.OrganizationEntity)
	organizationEntityPtr.Id = utils.GetStringPropOrEmpty(props, "id")
	organizationEntityPtr.Name = utils.GetStringPropOrEmpty(props, "name")
	organizationEntityPtr.Description = utils.GetStringPropOrEmpty(props, "description")
	organizationEntityPtr.Domain = utils.GetStringPropOrEmpty(props, "domain")
	organizationEntityPtr.Website = utils.GetStringPropOrEmpty(props, "website")
	organizationEntityPtr.Industry = utils.GetStringPropOrEmpty(props, "industry")
	organizationEntityPtr.IsPublic = utils.GetBoolPropOrFalse(props, "isPublic")
	organizationEntityPtr.CreatedAt = utils.GetTimePropOrNow(props, "createdAt")
	organizationEntityPtr.Readonly = utils.GetBoolPropOrFalse(props, "readonly")
	return organizationEntityPtr
}
