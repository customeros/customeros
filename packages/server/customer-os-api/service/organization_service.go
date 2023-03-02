package service

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/sirupsen/logrus"
	"reflect"
)

type OrganizationService interface {
	Create(ctx context.Context, input *OrganizationCreateData) (*entity.OrganizationEntity, error)
	Update(ctx context.Context, input *OrganizationUpdateData) (*entity.OrganizationEntity, error)
	GetOrganizationForJobRole(ctx context.Context, roleId string) (*entity.OrganizationEntity, error)
	GetOrganizationById(ctx context.Context, organizationId string) (*entity.OrganizationEntity, error)
	FindAll(ctx context.Context, page, limit int, filter *model.Filter, sortBy []*model.SortBy) (*utils.Pagination, error)
	GetOrganizationsForContact(ctx context.Context, contactId string, page, limit int, filter *model.Filter, sortBy []*model.SortBy) (*utils.Pagination, error)
	PermanentDelete(ctx context.Context, organizationId string) (bool, error)
	Merge(ctx context.Context, primaryOrganizationId, mergedOrganizationId string) error

	mapDbNodeToOrganizationEntity(node dbtype.Node) *entity.OrganizationEntity
}

type OrganizationCreateData struct {
	OrganizationEntity *entity.OrganizationEntity
	OrganizationTypeID *string
	Domains            []string
}

type OrganizationUpdateData struct {
	OrganizationEntity *entity.OrganizationEntity
	OrganizationTypeID *string
	Domains            []string
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

		err = s.repositories.OrganizationRepository.LinkWithDomainsInTx(ctx, tx, tenant, organizationId, input.Domains)
		if err != nil {
			return nil, err
		}

		if input.OrganizationTypeID != nil {
			err = s.repositories.OrganizationRepository.LinkWithOrganizationTypeInTx(ctx, tx, tenant, organizationId, *input.OrganizationTypeID)
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

		organizationDbNodePtr, err := s.repositories.OrganizationRepository.Update(ctx, tx, tenant, *input.OrganizationEntity)
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

		err = s.repositories.OrganizationRepository.UnlinkFromOrganizationTypesInTx(ctx, tx, tenant, organizationId)
		if err != nil {
			return nil, err
		}
		if input.OrganizationTypeID != nil {
			err := s.repositories.OrganizationRepository.LinkWithOrganizationTypeInTx(ctx, tx, tenant, organizationId, *input.OrganizationTypeID)
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
	session := utils.NewNeo4jReadSession(ctx, *s.repositories.Drivers.Neo4jDriver)
	defer session.Close(ctx)

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

func (s *organizationService) GetOrganizationsForContact(ctx context.Context, contactId string, page, limit int, filter *model.Filter, sortBy []*model.SortBy) (*utils.Pagination, error) {
	session := utils.NewNeo4jReadSession(ctx, *s.repositories.Drivers.Neo4jDriver)
	defer session.Close(ctx)

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
		session,
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

	organizationEntities := entity.OrganizationEntities{}

	for _, v := range dbNodesWithTotalCount.Nodes {
		organizationEntities = append(organizationEntities, *s.mapDbNodeToOrganizationEntity(*v))
	}
	paginatedResult.SetRows(&organizationEntities)
	return &paginatedResult, nil
}

func (s *organizationService) GetOrganizationForJobRole(ctx context.Context, roleId string) (*entity.OrganizationEntity, error) {
	session := utils.NewNeo4jReadSession(ctx, *s.repositories.Drivers.Neo4jDriver)
	defer session.Close(ctx)

	dbNode, err := s.repositories.OrganizationRepository.GetOrganizationForJobRole(ctx, session, common.GetContext(ctx).Tenant, roleId)
	if dbNode == nil || err != nil {
		return nil, err
	}
	return s.mapDbNodeToOrganizationEntity(*dbNode), nil
}

func (s *organizationService) GetOrganizationById(ctx context.Context, organizationId string) (*entity.OrganizationEntity, error) {
	dbNode, err := s.repositories.OrganizationRepository.GetOrganizationById(ctx, common.GetTenantFromContext(ctx), organizationId)
	if err != nil {
		return nil, err
	}
	return s.mapDbNodeToOrganizationEntity(*dbNode), nil
}

func (s *organizationService) PermanentDelete(ctx context.Context, organizationId string) (bool, error) {
	session := utils.NewNeo4jWriteSession(ctx, *s.repositories.Drivers.Neo4jDriver)
	defer session.Close(ctx)

	err := s.repositories.OrganizationRepository.Delete(ctx, session, common.GetContext(ctx).Tenant, organizationId)

	if err != nil {
		return false, err
	}

	return true, nil
}

func (s *organizationService) Merge(ctx context.Context, primaryOrganizationId, mergedOrganizationId string) error {
	session := utils.NewNeo4jWriteSession(ctx, *s.repositories.Drivers.Neo4jDriver)
	defer session.Close(ctx)

	_, err := s.GetOrganizationById(ctx, primaryOrganizationId)
	if err != nil {
		logrus.Errorf("Primary organization with id %s not found: %v", primaryOrganizationId, err)
		return err
	}
	_, err = s.GetOrganizationById(ctx, mergedOrganizationId)
	if err != nil {
		logrus.Errorf("Organization to merge with id %s not found: %v", mergedOrganizationId, err)
		return err
	}

	tenant := common.GetContext(ctx).Tenant
	_, err = session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		err = s.repositories.OrganizationRepository.MergeOrganizationPropertiesInTx(ctx, tx, tenant, primaryOrganizationId, mergedOrganizationId, entity.DataSourceOpenline)
		if err != nil {
			return nil, err
		}

		err = s.repositories.OrganizationRepository.MergeOrganizationRelationsInTx(ctx, tx, tenant, primaryOrganizationId, mergedOrganizationId)
		if err != nil {
			return nil, err
		}

		err = s.repositories.OrganizationRepository.AdaptMergedOrganizationLabelsInTx(ctx, tx, tenant, mergedOrganizationId)
		if err != nil {
			return nil, err
		}

		return nil, nil
	})
	return err
}

func (s *organizationService) mapDbNodeToOrganizationEntity(node dbtype.Node) *entity.OrganizationEntity {
	organizationEntityPtr := new(entity.OrganizationEntity)
	props := utils.GetPropsFromNode(node)
	organizationEntityPtr.ID = utils.GetStringPropOrEmpty(props, "id")
	organizationEntityPtr.Name = utils.GetStringPropOrEmpty(props, "name")
	organizationEntityPtr.Description = utils.GetStringPropOrEmpty(props, "description")
	organizationEntityPtr.Website = utils.GetStringPropOrEmpty(props, "website")
	organizationEntityPtr.Industry = utils.GetStringPropOrEmpty(props, "industry")
	organizationEntityPtr.IsPublic = utils.GetBoolPropOrFalse(props, "isPublic")
	organizationEntityPtr.CreatedAt = utils.GetTimePropOrEpochStart(props, "createdAt")
	organizationEntityPtr.UpdatedAt = utils.GetTimePropOrEpochStart(props, "updatedAt")
	organizationEntityPtr.Source = entity.GetDataSource(utils.GetStringPropOrEmpty(props, "source"))
	organizationEntityPtr.SourceOfTruth = entity.GetDataSource(utils.GetStringPropOrEmpty(props, "sourceOfTruth"))
	organizationEntityPtr.AppSource = utils.GetStringPropOrEmpty(props, "appSource")
	return organizationEntityPtr
}
