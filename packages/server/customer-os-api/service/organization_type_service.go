package service

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils"
)

type OrganizationTypeService interface {
	Create(ctx context.Context, organizationType *entity.OrganizationTypeEntity) (*entity.OrganizationTypeEntity, error)
	Update(ctx context.Context, organizationType *entity.OrganizationTypeEntity) (*entity.OrganizationTypeEntity, error)
	Delete(ctx context.Context, id string) (bool, error)
	GetAll(ctx context.Context) (*entity.OrganizationTypeEntities, error)
	FindOrganizationTypeForOrganization(ctx context.Context, organizationId string) (*entity.OrganizationTypeEntity, error)
}

type organizationTypeService struct {
	repository *repository.Repositories
}

func NewOrganizationTypeService(repository *repository.Repositories) OrganizationTypeService {
	return &organizationTypeService{
		repository: repository,
	}
}

func (s *organizationTypeService) Create(ctx context.Context, organizationType *entity.OrganizationTypeEntity) (*entity.OrganizationTypeEntity, error) {
	organizationTypeNode, err := s.repository.OrganizationTypeRepository.Create(ctx, common.GetContext(ctx).Tenant, organizationType)
	if err != nil {
		return nil, err
	}
	return s.mapDbNodeToOrganizationTypeEntity(organizationTypeNode), nil
}

func (s *organizationTypeService) Update(ctx context.Context, organizationType *entity.OrganizationTypeEntity) (*entity.OrganizationTypeEntity, error) {
	organizationTypeNode, err := s.repository.OrganizationTypeRepository.Update(ctx, common.GetContext(ctx).Tenant, organizationType)
	if err != nil {
		return nil, err
	}
	return s.mapDbNodeToOrganizationTypeEntity(organizationTypeNode), nil
}

func (s *organizationTypeService) Delete(ctx context.Context, id string) (bool, error) {
	err := s.repository.OrganizationTypeRepository.Delete(ctx, common.GetContext(ctx).Tenant, id)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s *organizationTypeService) GetAll(ctx context.Context) (*entity.OrganizationTypeEntities, error) {
	organizationTypeDbNodes, err := s.repository.OrganizationTypeRepository.FindAll(ctx, common.GetContext(ctx).Tenant)
	if err != nil {
		return nil, err
	}
	organizationTypeEntities := entity.OrganizationTypeEntities{}
	for _, dbNode := range organizationTypeDbNodes {
		organizationTypeEntity := s.mapDbNodeToOrganizationTypeEntity(dbNode)
		organizationTypeEntities = append(organizationTypeEntities, *organizationTypeEntity)
	}
	return &organizationTypeEntities, nil
}

func (s *organizationTypeService) FindOrganizationTypeForOrganization(ctx context.Context, organizationId string) (*entity.OrganizationTypeEntity, error) {
	organizationTypeDbNode, err := s.repository.OrganizationTypeRepository.FindForOrganization(ctx, common.GetContext(ctx).Tenant, organizationId)
	if err != nil {
		return nil, err
	} else if organizationTypeDbNode == nil {
		return nil, nil
	} else {
		return s.mapDbNodeToOrganizationTypeEntity(organizationTypeDbNode), nil
	}
}

func (s *organizationTypeService) mapDbNodeToOrganizationTypeEntity(dbNode *dbtype.Node) *entity.OrganizationTypeEntity {
	props := utils.GetPropsFromNode(*dbNode)
	organizationType := entity.OrganizationTypeEntity{
		Id:        utils.GetStringPropOrEmpty(props, "id"),
		Name:      utils.GetStringPropOrEmpty(props, "name"),
		CreatedAt: utils.GetTimePropOrEpochStart(props, "createdAt"),
		UpdatedAt: utils.GetTimePropOrEpochStart(props, "updatedAt"),
	}
	return &organizationType
}
