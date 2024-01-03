package service

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	neo4jentity "github.com/openline-ai/customer-os-neo4j-repository/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type TagService interface {
	Merge(ctx context.Context, tag *entity.TagEntity) (*entity.TagEntity, error)
	Update(ctx context.Context, tag *entity.TagEntity) (*entity.TagEntity, error)
	UnlinkAndDelete(ctx context.Context, id string) (bool, error)
	GetAll(ctx context.Context) (*entity.TagEntities, error)
	GetById(ctx context.Context, tagId string) (*entity.TagEntity, error)
	GetByName(ctx context.Context, tagName string) (*entity.TagEntity, error)
	GetTagsForContact(ctx context.Context, contactId string) (*entity.TagEntities, error)
	GetTagsForContacts(ctx context.Context, contactIds []string) (*entity.TagEntities, error)
	GetTagsForIssues(ctx context.Context, issueIds []string) (*entity.TagEntities, error)
	GetTagsForOrganizations(ctx context.Context, organizationIds []string) (*entity.TagEntities, error)
	GetTagsForLogEntries(ctx context.Context, logEntryIds []string) (*entity.TagEntities, error)
}

type tagService struct {
	log          logger.Logger
	repositories *repository.Repositories
}

func NewTagService(log logger.Logger, repository *repository.Repositories) TagService {
	return &tagService{
		log:          log,
		repositories: repository,
	}
}

func (s *tagService) Merge(ctx context.Context, tag *entity.TagEntity) (*entity.TagEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TagService.Merge")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("tag", tag))

	tagNodePtr, err := s.repositories.TagRepository.Merge(ctx, common.GetTenantFromContext(ctx), *tag)
	if err != nil {
		return nil, err
	}
	return s.mapDbNodeToTagEntity(*tagNodePtr), nil
}

func (s *tagService) Update(ctx context.Context, tag *entity.TagEntity) (*entity.TagEntity, error) {
	tagNodePtr, err := s.repositories.TagRepository.Update(ctx, common.GetTenantFromContext(ctx), *tag)
	if err != nil {
		return nil, err
	}
	return s.mapDbNodeToTagEntity(*tagNodePtr), nil
}

func (s *tagService) UnlinkAndDelete(ctx context.Context, id string) (bool, error) {
	err := s.repositories.TagRepository.UnlinkAndDelete(ctx, common.GetTenantFromContext(ctx), id)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s *tagService) GetAll(ctx context.Context) (*entity.TagEntities, error) {
	tagDbNodes, err := s.repositories.TagRepository.GetAll(ctx, common.GetTenantFromContext(ctx))
	if err != nil {
		return nil, err
	}
	tagEntities := make(entity.TagEntities, 0, len(tagDbNodes))
	for _, dbNodePtr := range tagDbNodes {
		tagEntities = append(tagEntities, *s.mapDbNodeToTagEntity(*dbNodePtr))
	}
	return &tagEntities, nil
}

func (s *tagService) GetTagsForContact(ctx context.Context, contactId string) (*entity.TagEntities, error) {
	tagDbNodes, err := s.repositories.TagRepository.GetForContact(ctx, common.GetTenantFromContext(ctx), contactId)
	if err != nil {
		return nil, err
	}
	tagEntities := make(entity.TagEntities, 0, len(tagDbNodes))
	for _, dbNodePtr := range tagDbNodes {
		tagEntities = append(tagEntities, *s.mapDbNodeToTagEntity(*dbNodePtr))
	}
	return &tagEntities, nil
}

func (s *tagService) GetTagsForContacts(ctx context.Context, contactIds []string) (*entity.TagEntities, error) {
	tags, err := s.repositories.TagRepository.GetForContacts(ctx, common.GetTenantFromContext(ctx), contactIds)
	if err != nil {
		return nil, err
	}
	tagEntities := make(entity.TagEntities, 0, len(tags))
	for _, v := range tags {
		tagEntity := s.mapDbNodeToTagEntity(*v.Node)
		s.addDbRelationshipToTagEntity(*v.Relationship, tagEntity)
		tagEntity.DataloaderKey = v.LinkedNodeId
		tagEntities = append(tagEntities, *tagEntity)
	}
	return &tagEntities, nil
}

func (s *tagService) GetTagsForIssues(ctx context.Context, issueIds []string) (*entity.TagEntities, error) {
	tags, err := s.repositories.TagRepository.GetForIssues(ctx, common.GetTenantFromContext(ctx), issueIds)
	if err != nil {
		return nil, err
	}
	tagEntities := make(entity.TagEntities, 0, len(tags))
	for _, v := range tags {
		tagEntity := s.mapDbNodeToTagEntity(*v.Node)
		s.addDbRelationshipToTagEntity(*v.Relationship, tagEntity)
		tagEntity.DataloaderKey = v.LinkedNodeId
		tagEntities = append(tagEntities, *tagEntity)
	}
	return &tagEntities, nil
}

func (s *tagService) GetTagsForOrganizations(ctx context.Context, organizationIDs []string) (*entity.TagEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TagService.GetTagsForOrganizations")
	defer span.Finish()
	span.LogFields(log.Object("organizationIDs", organizationIDs))

	tags, err := s.repositories.TagRepository.GetForOrganizations(ctx, common.GetTenantFromContext(ctx), organizationIDs)
	if err != nil {
		return nil, err
	}
	tagEntities := make(entity.TagEntities, 0, len(tags))
	for _, v := range tags {
		tagEntity := s.mapDbNodeToTagEntity(*v.Node)
		s.addDbRelationshipToTagEntity(*v.Relationship, tagEntity)
		tagEntity.DataloaderKey = v.LinkedNodeId
		tagEntities = append(tagEntities, *tagEntity)
	}
	return &tagEntities, nil
}

func (s *tagService) GetTagsForLogEntries(ctx context.Context, logEntryIds []string) (*entity.TagEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TagService.GetTagsForLogEntries")
	defer span.Finish()
	span.LogFields(log.Object("logEntryIds", logEntryIds))

	tags, err := s.repositories.TagRepository.GetForLogEntries(ctx, common.GetTenantFromContext(ctx), logEntryIds)
	if err != nil {
		return nil, err
	}
	tagEntities := make(entity.TagEntities, 0, len(tags))
	for _, v := range tags {
		tagEntity := s.mapDbNodeToTagEntity(*v.Node)
		s.addDbRelationshipToTagEntity(*v.Relationship, tagEntity)
		tagEntity.DataloaderKey = v.LinkedNodeId
		tagEntities = append(tagEntities, *tagEntity)
	}
	return &tagEntities, nil
}

func (s *tagService) GetById(ctx context.Context, tagId string) (*entity.TagEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TagService.GetById")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("tagId", tagId))

	tagDbNode, err := s.repositories.TagRepository.GetById(ctx, tagId)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	return s.mapDbNodeToTagEntity(*tagDbNode), nil
}

func (s *tagService) GetByName(ctx context.Context, tagName string) (*entity.TagEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TagService.GetByName")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("tagName", tagName))

	tagDbNode, err := s.repositories.TagRepository.GetByName(ctx, tagName)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	return s.mapDbNodeToTagEntity(*tagDbNode), nil
}

func (s *tagService) mapDbNodeToTagEntity(dbNode dbtype.Node) *entity.TagEntity {
	props := utils.GetPropsFromNode(dbNode)
	tag := entity.TagEntity{
		Id:        utils.GetStringPropOrEmpty(props, "id"),
		Name:      utils.GetStringPropOrEmpty(props, "name"),
		CreatedAt: utils.GetTimePropOrEpochStart(props, "createdAt"),
		UpdatedAt: utils.GetTimePropOrEpochStart(props, "updatedAt"),
		AppSource: utils.GetStringPropOrEmpty(props, "appSource"),
		Source:    neo4jentity.GetDataSource(utils.GetStringPropOrEmpty(props, "source")),
	}
	return &tag
}

func (s *tagService) addDbRelationshipToTagEntity(relationship dbtype.Relationship, tagEntity *entity.TagEntity) {
	props := utils.GetPropsFromRelationship(relationship)
	tagEntity.TaggedAt = utils.GetTimePropOrEpochStart(props, "taggedAt")
}
