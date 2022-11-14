package service

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/customer-os-api/common"
	"github.com/openline-ai/openline-customer-os/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/customer-os-api/repository"
	"github.com/openline-ai/openline-customer-os/customer-os-api/utils"
)

type ContactTypeService interface {
	Create(ctx context.Context, contactType *entity.ContactTypeEntity) (*entity.ContactTypeEntity, error)
}

type contactTypeService struct {
	repository *repository.RepositoryContainer
}

func NewContactTypeService(repository *repository.RepositoryContainer) ContactTypeService {
	return &contactTypeService{
		repository: repository,
	}
}

func (s *contactTypeService) Create(ctx context.Context, contactType *entity.ContactTypeEntity) (*entity.ContactTypeEntity, error) {
	contactTypeNode, err := s.repository.ContactTypeRepository.Create(common.GetContext(ctx).Tenant, contactType)
	if err != nil {
		return nil, err
	}
	return s.mapDbNodeToContactTypeEntity(*contactTypeNode), nil
}

func (s *contactTypeService) mapDbNodeToContactTypeEntity(dbNode dbtype.Node) *entity.ContactTypeEntity {
	props := utils.GetPropsFromNode(dbNode)
	contactType := entity.ContactTypeEntity{
		Id:   utils.GetStringPropOrEmpty(props, "id"),
		Name: utils.GetStringPropOrEmpty(props, "name"),
	}
	return &contactType
}
