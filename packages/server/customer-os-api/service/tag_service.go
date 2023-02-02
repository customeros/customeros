package service

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils"
)

type TagService interface {
	Merge(ctx context.Context, tag *entity.TagEntity) (*entity.TagEntity, error)
	// FIXME alexb refactor
	Update(ctx context.Context, tag *entity.TagEntity) (*entity.TagEntity, error)
	// FIXME alexb refactor
	Delete(ctx context.Context, id string) (bool, error)
	// FIXME alexb refactor
	GetAll(ctx context.Context) (*entity.TagEntities, error)
	// FIXME alexb refactor
	FindTagForContact(ctx context.Context, contactId string) (*entity.TagEntity, error)
}

type tagService struct {
	repository *repository.Repositories
}

func NewTagService(repository *repository.Repositories) TagService {
	return &tagService{
		repository: repository,
	}
}

func (s *tagService) Merge(ctx context.Context, tag *entity.TagEntity) (*entity.TagEntity, error) {
	tagNodePtr, err := s.repository.TagRepository.Merge(common.GetTenantFromContext(ctx), *tag)
	if err != nil {
		return nil, err
	}
	return s.mapDbNodeToTagEntity(*tagNodePtr), nil
}

func (s *tagService) Update(ctx context.Context, tag *entity.TagEntity) (*entity.TagEntity, error) {
	tagNode, err := s.repository.TagRepository.Update(common.GetContext(ctx).Tenant, tag)
	if err != nil {
		return nil, err
	}
	return s.mapDbNodeToTagEntity(*tagNode), nil
}

func (s *tagService) Delete(ctx context.Context, id string) (bool, error) {
	err := s.repository.TagRepository.Delete(common.GetContext(ctx).Tenant, id)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s *tagService) GetAll(ctx context.Context) (*entity.TagEntities, error) {
	tagDbNodes, err := s.repository.TagRepository.FindAll(common.GetContext(ctx).Tenant)
	if err != nil {
		return nil, err
	}
	tagEntities := entity.TagEntities{}
	for _, dbNode := range tagDbNodes {
		tagEntity := s.mapDbNodeToTagEntity(*dbNode)
		tagEntities = append(tagEntities, *tagEntity)
	}
	return &tagEntities, nil
}

func (s *tagService) FindTagForContact(ctx context.Context, contactId string) (*entity.TagEntity, error) {
	tagDbNode, err := s.repository.TagRepository.FindForContact(common.GetContext(ctx).Tenant, contactId)
	if err != nil {
		return nil, err
	} else if tagDbNode == nil {
		return nil, nil
	} else {
		return s.mapDbNodeToTagEntity(*tagDbNode), nil
	}
}

func (s *tagService) mapDbNodeToTagEntity(dbNode dbtype.Node) *entity.TagEntity {
	props := utils.GetPropsFromNode(dbNode)
	tag := entity.TagEntity{
		Id:        utils.GetStringPropOrEmpty(props, "id"),
		Name:      utils.GetStringPropOrEmpty(props, "name"),
		CreatedAt: utils.GetTimePropOrNow(props, "createdAt"),
		UpdatedAt: utils.GetTimePropOrNow(props, "updatedAt"),
		AppSource: utils.GetStringPropOrEmpty(props, "appSource"),
		Source:    entity.GetDataSource(utils.GetStringPropOrEmpty(props, "source")),
	}
	return &tag
}
