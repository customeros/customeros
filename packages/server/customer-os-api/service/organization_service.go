package service

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils"
	"reflect"
)

type OrganizationService interface {
	Create(ctx context.Context, input *OrganizationCreateData) (*entity.OrganizationEntity, error)
	Update(ctx context.Context, input *OrganizationUpdateData) (*entity.OrganizationEntity, error)
	FindOrganizationForRole(ctx context.Context, roleId string) (*entity.OrganizationEntity, error)
	GetOrganizationById(ctx context.Context, organizationId string) (*entity.OrganizationEntity, error)
	FindAll(ctx context.Context, page, limit int, filter *model.Filter, sortBy []*model.SortBy) (*utils.Pagination, error)
	PermanentDelete(ctx context.Context, organizationId string) (bool, error)
}

type OrganizationCreateData struct {
	OrganizationEntity *entity.OrganizationEntity
	OrganizationTypeID *string
}

type OrganizationUpdateData struct {
	OrganizationEntity *entity.OrganizationEntity
	OrganizationTypeID *string
}

type organizationService struct {
	repositories *repository.Repositories
}

func NewOrganizationService(repositories *repository.Repositories) OrganizationService {
	return &organizationService{
		repositories: repositories,
	}
}

func (s *organizationService) Create(ctx context.Context, input *OrganizationCreateData) (*entity.OrganizationEntity, error) {
	session := utils.NewNeo4jWriteSession(*s.repositories.Drivers.Neo4jDriver)
	defer session.Close()

	organizationDbNodePtr, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		tenant := common.GetContext(ctx).Tenant
		organizationDbNodePtr, err := s.repositories.OrganizationRepository.Create(tx, tenant, *input.OrganizationEntity)
		if err != nil {
			return nil, err
		}
		var organizationId = utils.GetPropsFromNode(*organizationDbNodePtr)["id"].(string)

		if input.OrganizationTypeID != nil {
			err = s.repositories.OrganizationRepository.LinkWithOrganizationTypeInTx(tx, tenant, organizationId, *input.OrganizationTypeID)
			if err != nil {
				return nil, err
			}
		}

		return organizationDbNodePtr, nil
	})
	if err != nil {
		return nil, err
	}
	return s.mapDbNodeToOrganizationEntity(*organizationDbNodePtr.(*dbtype.Node)), nil
}

func (s *organizationService) Update(ctx context.Context, input *OrganizationUpdateData) (*entity.OrganizationEntity, error) {
	session := utils.NewNeo4jWriteSession(*s.repositories.Drivers.Neo4jDriver)
	defer session.Close()

	organizationDbNodePtr, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		tenant := common.GetContext(ctx).Tenant
		organizationDbNodePtr, err := s.repositories.OrganizationRepository.Update(tx, tenant, *input.OrganizationEntity)
		if err != nil {
			return nil, err
		}
		var organizationId = utils.GetPropsFromNode(*organizationDbNodePtr)["id"].(string)

		err = s.repositories.OrganizationRepository.UnlinkFromOrganizationTypesInTx(tx, tenant, organizationId)
		if err != nil {
			return nil, err
		}
		if input.OrganizationTypeID != nil {
			err := s.repositories.OrganizationRepository.LinkWithOrganizationTypeInTx(tx, tenant, organizationId, *input.OrganizationTypeID)
			if err != nil {
				return nil, err
			}
		}

		return organizationDbNodePtr, nil
	})
	if err != nil {
		return nil, err
	}
	return s.mapDbNodeToOrganizationEntity(*organizationDbNodePtr.(*dbtype.Node)), nil
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

func (s *organizationService) FindOrganizationForRole(ctx context.Context, roleId string) (*entity.OrganizationEntity, error) {
	session := utils.NewNeo4jReadSession(*s.repositories.Drivers.Neo4jDriver)
	defer session.Close()

	dbNode, err := s.repositories.OrganizationRepository.FindOrganizationForRole(session, common.GetContext(ctx).Tenant, roleId)
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
	organizationEntityPtr := new(entity.OrganizationEntity)
	props := utils.GetPropsFromNode(node)
	organizationEntityPtr.ID = utils.GetStringPropOrEmpty(props, "id")
	organizationEntityPtr.Name = utils.GetStringPropOrEmpty(props, "name")
	organizationEntityPtr.Description = utils.GetStringPropOrEmpty(props, "description")
	organizationEntityPtr.Domain = utils.GetStringPropOrEmpty(props, "domain")
	organizationEntityPtr.Website = utils.GetStringPropOrEmpty(props, "website")
	organizationEntityPtr.Industry = utils.GetStringPropOrEmpty(props, "industry")
	organizationEntityPtr.IsPublic = utils.GetBoolPropOrFalse(props, "isPublic")
	organizationEntityPtr.CreatedAt = utils.GetTimePropOrNow(props, "createdAt")
	organizationEntityPtr.UpdatedAt = utils.GetTimePropOrNow(props, "updatedAt")
	organizationEntityPtr.Source = entity.GetDataSource(utils.GetStringPropOrEmpty(props, "source"))
	organizationEntityPtr.SourceOfTruth = entity.GetDataSource(utils.GetStringPropOrEmpty(props, "sourceOfTruth"))
	organizationEntityPtr.AppSource = utils.GetStringPropOrEmpty(props, "appSource")
	return organizationEntityPtr
}
