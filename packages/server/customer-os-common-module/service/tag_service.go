package service

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/dto"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/constants"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
)

type TagService interface {
	Save(ctx context.Context, tx *neo4j.ManagedTransaction, inputTag *neo4jentity.TagEntity) (*neo4jentity.TagEntity, error)
	AddTagToEntity(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, entityId string, entityType model.EntityType, tagId, tagName string) (string, error)
	RemoveTagFromEntity(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, entityId string, entityType model.EntityType, tagId string) error
	Update(ctx context.Context, tagId, name string) error
	UnlinkAndDelete(ctx context.Context, id string) (bool, error)
	GetAll(ctx context.Context) (*neo4jentity.TagEntities, error)
	GetById(ctx context.Context, tagId string) (*neo4jentity.TagEntity, error)
	GetTagByEntityTypeAndName(ctx context.Context, entityType model.EntityType, name string) (*neo4jentity.TagEntity, error)
	GetTagsByEntityType(ctx context.Context, entityType model.EntityType) (*neo4jentity.TagEntities, error)
	GetTagsForContacts(ctx context.Context, contactIds []string) (*neo4jentity.TagEntities, error)
	GetTagsForIssues(ctx context.Context, issueIds []string) (*neo4jentity.TagEntities, error)
	GetTagsForOrganizations(ctx context.Context, organizationIds []string) (*neo4jentity.TagEntities, error)
	GetTagsForLogEntries(ctx context.Context, logEntryIds []string) (*neo4jentity.TagEntities, error)
}

type tagService struct {
	log      logger.Logger
	services *Services
}

func (s *tagService) GetTagByEntityTypeAndName(ctx context.Context, entityType model.EntityType, name string) (*neo4jentity.TagEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TagService.GetTagByEntityTypeAndName")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("entityType", entityType.String()), log.String("name", name))

	tagDbNode, err := s.services.Neo4jRepositories.TagReadRepository.GetByEntityTypeAndName(ctx, common.GetTenantFromContext(ctx), entityType, name)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	return neo4jmapper.MapDbNodeToTagEntity(tagDbNode), nil
}

func NewTagService(log logger.Logger, services *Services) TagService {
	return &tagService{
		log:      log,
		services: services,
	}
}

func (s *tagService) Save(ctx context.Context, tx *neo4j.ManagedTransaction, inputTag *neo4jentity.TagEntity) (*neo4jentity.TagEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TagService.Save")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("inputTag", inputTag))

	// validate tenant
	err := common.ValidateTenant(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	tenant := common.GetTenantFromContext(ctx)

	if inputTag.Source == "" {
		inputTag.Source = constants.SourceOpenline
	}
	if inputTag.AppSource == "" {
		inputTag.AppSource = common.GetAppSourceFromContext(ctx)
	}

	tagNodePtr, err := s.services.Neo4jRepositories.TagWriteRepository.Merge(ctx, tx, tenant, *inputTag)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return neo4jmapper.MapDbNodeToTagEntity(tagNodePtr), nil
}

func (s *tagService) AddTagToEntity(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, entityId string, entityType model.EntityType, tagId, tagName string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TagService.AddTagToEntity")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	tracing.TagEntity(span, entityId)
	span.LogFields(log.String("tagId", tagId), log.String("tagName", tagName), log.String("entityType", entityType.String()))

	if tagId == "" {
		tagEntity, err := s.Save(ctx, tx, &neo4jentity.TagEntity{
			Name:       tagName,
			Source:     constants.SourceOpenline,
			AppSource:  common.GetTenantFromContext(ctx),
			EntityType: entityType,
		})
		if err != nil {
			tracing.TraceErr(span, err)
			return "", err
		}
		if tagEntity != nil {
			tagId = tagEntity.Id
		}
	}

	tagEntity, err := s.GetById(ctx, tagId)
	if err != nil {
		tracing.TraceErr(span, err)
		return "", err
	}
	if tagEntity.EntityType != entityType {
		err = errors.New("tag entity type mismatch")
		tracing.TraceErr(span, err)
		return "", err
	}

	err = s.services.Neo4jRepositories.TagWriteRepository.LinkTagByIdToEntity(ctx, tx, tenant, tagId, entityId, entityType)
	if err != nil {
		tracing.TraceErr(span, err)
		return "", err
	}

	// event for tag added
	err = s.services.RabbitMQService.PublishEvent(ctx, entityId, entityType, dto.NewAddTagEvent(tagId, tagName))
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "unable to publish message AddTagEvent"))
	}

	if common.GetAppSourceFromContext(ctx) != constants.AppSourceCustomerOsApi {
		utils.EventCompleted(ctx, tenant, entityType.String(), entityId, s.services.GrpcClients, utils.NewEventCompletedDetails().WithUpdate())
	}

	return tagId, nil
}

func (s *tagService) RemoveTagFromEntity(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, entityId string, entityType model.EntityType, tagId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TagService.RemoveTagFromEntity")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	tracing.TagEntity(span, entityId)
	span.LogFields(log.String("tagId", tagId), log.String("entityType", entityType.String()))

	tagEntity, err := s.GetById(ctx, tagId)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "unable to get tag by id"))
		return err
	}

	err = s.services.Neo4jRepositories.TagWriteRepository.UnlinkTagByIdFromEntity(ctx, tx, tenant, tagId, entityId, entityType)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "unable to unlink tag from entity"))
		return err
	}

	// event for tag removed
	err = s.services.RabbitMQService.PublishEvent(ctx, entityId, entityType, dto.NewRemoveTagEvent(tagId, tagEntity.Name))
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "unable to publish message RemoveTagEvent"))
	}

	if common.GetAppSourceFromContext(ctx) != constants.AppSourceCustomerOsApi {
		utils.EventCompleted(ctx, tenant, entityType.String(), entityId, s.services.GrpcClients, utils.NewEventCompletedDetails().WithUpdate())
	}

	return nil
}

func (s *tagService) Update(ctx context.Context, tagId, name string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TagService.Update")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("tagId", tagId), log.String("name", name))

	if name == "" {
		err := errors.New("name is required")
		tracing.TraceErr(span, err)
	}

	err := s.services.Neo4jRepositories.TagWriteRepository.UpdateName(ctx, common.GetTenantFromContext(ctx), tagId, name)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error updating tag name: %s", err.Error())
		return err
	}
	return nil
}

func (s *tagService) UnlinkAndDelete(ctx context.Context, id string) (bool, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TagService.UnlinkAndDelete")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	err := s.services.Neo4jRepositories.TagWriteRepository.UnlinkAllAndDelete(ctx, common.GetTenantFromContext(ctx), id)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s *tagService) GetAll(ctx context.Context) (*neo4jentity.TagEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TagService.GetAll")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	tagDbNodes, err := s.services.Neo4jRepositories.TagReadRepository.GetAll(ctx, common.GetTenantFromContext(ctx))
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

	tags, err := s.services.Neo4jRepositories.TagReadRepository.GetForContacts(ctx, common.GetTenantFromContext(ctx), contactIds)
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

	tags, err := s.services.Neo4jRepositories.TagReadRepository.GetForIssues(ctx, common.GetTenantFromContext(ctx), issueIds)
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

	tags, err := s.services.Neo4jRepositories.TagReadRepository.GetForOrganizations(ctx, common.GetTenantFromContext(ctx), organizationIDs)
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

	tags, err := s.services.Neo4jRepositories.TagReadRepository.GetForLogEntries(ctx, common.GetTenantFromContext(ctx), logEntryIds)
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

	tagDbNode, err := s.services.Neo4jRepositories.TagReadRepository.GetById(ctx, common.GetTenantFromContext(ctx), tagId)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	return neo4jmapper.MapDbNodeToTagEntity(tagDbNode), nil
}

func (s *tagService) addDbRelationshipToTagEntity(relationship dbtype.Relationship, tagEntity *neo4jentity.TagEntity) {
	props := utils.GetPropsFromRelationship(relationship)
	tagEntity.TaggedAt = utils.GetTimePropOrEpochStart(props, "taggedAt")
}

func (s *tagService) GetTagsByEntityType(ctx context.Context, entityType model.EntityType) (*neo4jentity.TagEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TagService.GetTagsByEntityType")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("entityType", entityType.String()))

	tagDbNodes, err := s.services.Neo4jRepositories.TagReadRepository.GetAllByEntityType(ctx, common.GetTenantFromContext(ctx), entityType)
	if err != nil {
		return nil, err
	}
	tagEntities := make(neo4jentity.TagEntities, 0, len(tagDbNodes))
	for _, dbNodePtr := range tagDbNodes {
		tagEntities = append(tagEntities, *neo4jmapper.MapDbNodeToTagEntity(dbNodePtr))
	}
	return &tagEntities, nil
}
