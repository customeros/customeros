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
	Update(ctx context.Context, contactType *entity.ContactTypeEntity) (*entity.ContactTypeEntity, error)
	Delete(ctx context.Context, id string) (bool, error)
	GetAll(ctx context.Context) (*entity.ContactTypeEntities, error)
	FindContactTypeForContact(ctx context.Context, contactId string) (*entity.ContactTypeEntity, error)
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
	return s.mapDbNodeToContactTypeEntity(contactTypeNode), nil
}

func (s *contactTypeService) Update(ctx context.Context, contactType *entity.ContactTypeEntity) (*entity.ContactTypeEntity, error) {
	contactTypeNode, err := s.repository.ContactTypeRepository.Update(common.GetContext(ctx).Tenant, contactType)
	if err != nil {
		return nil, err
	}
	return s.mapDbNodeToContactTypeEntity(contactTypeNode), nil
}

func (s *contactTypeService) Delete(ctx context.Context, id string) (bool, error) {
	err := s.repository.ContactTypeRepository.Delete(common.GetContext(ctx).Tenant, id)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s *contactTypeService) GetAll(ctx context.Context) (*entity.ContactTypeEntities, error) {
	contactTypeDbNodes, err := s.repository.ContactTypeRepository.FindAll(common.GetContext(ctx).Tenant)
	if err != nil {
		return nil, err
	}
	contactTypeEntities := entity.ContactTypeEntities{}
	for _, dbNode := range contactTypeDbNodes {
		contactTypeEntity := s.mapDbNodeToContactTypeEntity(dbNode)
		contactTypeEntities = append(contactTypeEntities, *contactTypeEntity)
	}
	return &contactTypeEntities, nil
}

func (s *contactTypeService) FindContactTypeForContact(ctx context.Context, contactId string) (*entity.ContactTypeEntity, error) {
	contactTypeDbNode, err := s.repository.ContactTypeRepository.FindForContact(common.GetContext(ctx).Tenant, contactId)
	if err != nil {
		return nil, err
	} else if contactTypeDbNode == nil {
		return nil, nil
	} else {
		return s.mapDbNodeToContactTypeEntity(contactTypeDbNode), nil
	}
}

func (s *contactTypeService) mapDbNodeToContactTypeEntity(dbNode *dbtype.Node) *entity.ContactTypeEntity {
	props := utils.GetPropsFromNode(*dbNode)
	contactType := entity.ContactTypeEntity{
		Id:   utils.GetStringPropOrEmpty(props, "id"),
		Name: utils.GetStringPropOrEmpty(props, "name"),
	}
	return &contactType
}
