package service

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"strings"
)

type TagService interface {
	Merge(ctx context.Context, tag *neo4jentity.TagEntity) (*neo4jentity.TagEntity, error)
	Update(ctx context.Context, tag *neo4jentity.TagEntity) (*neo4jentity.TagEntity, error)
	UnlinkAndDelete(ctx context.Context, id string) (bool, error)
	GetAll(ctx context.Context) (*neo4jentity.TagEntities, error)
	GetById(ctx context.Context, tagId string) (*neo4jentity.TagEntity, error)
	GetByNameOptional(ctx context.Context, tagName string) (*neo4jentity.TagEntity, error)
	GetTagsForContact(ctx context.Context, contactId string) (*neo4jentity.TagEntities, error)
	GetTagsForContacts(ctx context.Context, contactIds []string) (*neo4jentity.TagEntities, error)
	GetTagsForIssues(ctx context.Context, issueIds []string) (*neo4jentity.TagEntities, error)
	GetTagsForOrganizations(ctx context.Context, organizationIds []string) (*neo4jentity.TagEntities, error)
	GetTagsForLogEntries(ctx context.Context, logEntryIds []string) (*neo4jentity.TagEntities, error)
	GetTagId(ctx context.Context, tagId, tagName *string) string
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

func (s *tagService) Merge(ctx context.Context, tag *neo4jentity.TagEntity) (*neo4jentity.TagEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TagService.Merge")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("tag", tag))

	tagNodePtr, err := s.repositories.TagRepository.Merge(ctx, common.GetTenantFromContext(ctx), *tag)
	if err != nil {
		return nil, err
	}
	return neo4jmapper.MapDbNodeToTagEntity(tagNodePtr), nil
}

func (s *tagService) Update(ctx context.Context, tag *neo4jentity.TagEntity) (*neo4jentity.TagEntity, error) {
	tagNodePtr, err := s.repositories.TagRepository.Update(ctx, common.GetTenantFromContext(ctx), *tag)
	if err != nil {
		return nil, err
	}
	return neo4jmapper.MapDbNodeToTagEntity(tagNodePtr), nil
}

func (s *tagService) UnlinkAndDelete(ctx context.Context, id string) (bool, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TagService.UnlinkAndDelete")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	err := s.repositories.Neo4jRepositories.TagWriteRepository.UnlinkAllAndDelete(ctx, common.GetTenantFromContext(ctx), id)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s *tagService) GetAll(ctx context.Context) (*neo4jentity.TagEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TagService.GetAll")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	tagDbNodes, err := s.repositories.Neo4jRepositories.TagReadRepository.GetAll(ctx, common.GetTenantFromContext(ctx))
	if err != nil {
		return nil, err
	}
	tagEntities := make(neo4jentity.TagEntities, 0, len(tagDbNodes))
	for _, dbNodePtr := range tagDbNodes {
		tagEntities = append(tagEntities, *neo4jmapper.MapDbNodeToTagEntity(dbNodePtr))
	}
	return &tagEntities, nil
}

func (s *tagService) GetTagsForContact(ctx context.Context, contactId string) (*neo4jentity.TagEntities, error) {
	tagDbNodes, err := s.repositories.Neo4jRepositories.TagReadRepository.GetForContact(ctx, common.GetTenantFromContext(ctx), contactId)
	if err != nil {
		return nil, err
	}
	tagEntities := make(neo4jentity.TagEntities, 0, len(tagDbNodes))
	for _, dbNodePtr := range tagDbNodes {
		tagEntities = append(tagEntities, *neo4jmapper.MapDbNodeToTagEntity(dbNodePtr))
	}
	return &tagEntities, nil
}

func (s *tagService) GetTagsForContacts(ctx context.Context, contactIds []string) (*neo4jentity.TagEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TagService.GetTagsForContacts")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("contactIds", contactIds))

	tags, err := s.repositories.Neo4jRepositories.TagReadRepository.GetForContacts(ctx, common.GetTenantFromContext(ctx), contactIds)
	if err != nil {
		return nil, err
	}
	tagEntities := make(neo4jentity.TagEntities, 0, len(tags))
	for _, v := range tags {
		tagEntity := neo4jmapper.MapDbNodeToTagEntity(v.Node)
		s.addDbRelationshipToTagEntity(*v.Relationship, tagEntity)
		tagEntity.DataloaderKey = v.LinkedNodeId
		tagEntities = append(tagEntities, *tagEntity)
	}
	return &tagEntities, nil
}

func (s *tagService) GetTagsForIssues(ctx context.Context, issueIds []string) (*neo4jentity.TagEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TagService.GetTagsForIssues")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("issueIds", issueIds))

	tags, err := s.repositories.Neo4jRepositories.TagReadRepository.GetForIssues(ctx, common.GetTenantFromContext(ctx), issueIds)
	if err != nil {
		return nil, err
	}
	tagEntities := make(neo4jentity.TagEntities, 0, len(tags))
	for _, v := range tags {
		tagEntity := neo4jmapper.MapDbNodeToTagEntity(v.Node)
		s.addDbRelationshipToTagEntity(*v.Relationship, tagEntity)
		tagEntity.DataloaderKey = v.LinkedNodeId
		tagEntities = append(tagEntities, *tagEntity)
	}
	return &tagEntities, nil
}

func (s *tagService) GetTagsForOrganizations(ctx context.Context, organizationIDs []string) (*neo4jentity.TagEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TagService.GetTagsForOrganizations")
	defer span.Finish()
	span.LogFields(log.Object("organizationIDs", organizationIDs))

	tags, err := s.repositories.Neo4jRepositories.TagReadRepository.GetForOrganizations(ctx, common.GetTenantFromContext(ctx), organizationIDs)
	if err != nil {
		return nil, err
	}
	tagEntities := make(neo4jentity.TagEntities, 0, len(tags))
	for _, v := range tags {
		tagEntity := neo4jmapper.MapDbNodeToTagEntity(v.Node)
		s.addDbRelationshipToTagEntity(*v.Relationship, tagEntity)
		tagEntity.DataloaderKey = v.LinkedNodeId
		tagEntities = append(tagEntities, *tagEntity)
	}
	return &tagEntities, nil
}

func (s *tagService) GetTagsForLogEntries(ctx context.Context, logEntryIds []string) (*neo4jentity.TagEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TagService.GetTagsForLogEntries")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("logEntryIds", logEntryIds))

	tags, err := s.repositories.Neo4jRepositories.TagReadRepository.GetForLogEntries(ctx, common.GetTenantFromContext(ctx), logEntryIds)
	if err != nil {
		return nil, err
	}
	tagEntities := make(neo4jentity.TagEntities, 0, len(tags))
	for _, v := range tags {
		tagEntity := neo4jmapper.MapDbNodeToTagEntity(v.Node)
		s.addDbRelationshipToTagEntity(*v.Relationship, tagEntity)
		tagEntity.DataloaderKey = v.LinkedNodeId
		tagEntities = append(tagEntities, *tagEntity)
	}
	return &tagEntities, nil
}

func (s *tagService) GetById(ctx context.Context, tagId string) (*neo4jentity.TagEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TagService.GetById")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("tagId", tagId))

	tagDbNode, err := s.repositories.Neo4jRepositories.TagReadRepository.GetById(ctx, common.GetTenantFromContext(ctx), tagId)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	return neo4jmapper.MapDbNodeToTagEntity(tagDbNode), nil
}

func (s *tagService) GetByNameOptional(ctx context.Context, tagName string) (*neo4jentity.TagEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TagService.GetByNameOptional")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("tagName", tagName))

	tagDbNode, err := s.repositories.Neo4jRepositories.TagReadRepository.GetByNameOptional(ctx, common.GetTenantFromContext(ctx), tagName)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	if tagDbNode == nil {
		return nil, nil
	}
	return neo4jmapper.MapDbNodeToTagEntity(tagDbNode), nil
}

func (s *tagService) addDbRelationshipToTagEntity(relationship dbtype.Relationship, tagEntity *neo4jentity.TagEntity) {
	props := utils.GetPropsFromRelationship(relationship)
	tagEntity.TaggedAt = utils.GetTimePropOrEpochStart(props, "taggedAt")
}

func (s *tagService) GetTagId(ctx context.Context, tagId, tagName *string) string {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TagService.GetTagId")
	defer span.Finish()

	outputTagId := ""
	if tagId != nil && *tagId != "" {
		exists, err := s.repositories.Neo4jRepositories.CommonReadRepository.ExistsById(ctx, common.GetTenantFromContext(ctx), *tagId, model.NodeLabelTag)
		if err != nil {
			tracing.TraceErr(span, err)
			s.log.Error(ctx, "Error checking if tag exists by id", err)
		}
		if exists {
			outputTagId = *tagId
		}
	}
	if outputTagId == "" && tagName != nil && strings.TrimSpace(*tagName) != "" {
		tagEntity, _ := s.GetByNameOptional(ctx, strings.TrimSpace(*tagName))
		if tagEntity != nil {
			outputTagId = tagEntity.Id
		}
	}
	return outputTagId
}
