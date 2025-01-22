package tags

import (
	"context"

	"github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/constants"
	neo4jentity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"
	neo4jmapper "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/mapper"
	neo4j_repository "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/repository"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/dto"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/model"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/events"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
)

type tagService struct {
	log    logger.Logger
	neo4j  *neo4j_repository.Repositories
	events *events.EventsService
}

func NewTagService(log logger.Logger, neo4j *neo4j_repository.Repositories, events *events.EventsService) interfaces.TagService {
	return &tagService{
		log:    log,
		neo4j:  neo4j,
		events: events,
	}
}

func (s *tagService) GetTagByEntityTypeAndName(ctx context.Context, entityType model.EntityType, name string) (*neo4jentity.TagEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TagService.GetTagByEntityTypeAndName")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("entityType", entityType.String()), log.String("name", name))

	tagDbNode, err := s.neo4j.TagReadRepository.GetByEntityTypeAndName(ctx, common.GetTenantFromContext(ctx), entityType, name)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	return neo4jmapper.MapDbNodeToTagEntity(tagDbNode), nil
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
	if inputTag.ColorCode == "" {
		inputTag.ColorCode = utils.GetRandomColor()
	}

	tagNodePtr, err := s.neo4j.TagWriteRepository.Merge(ctx, tx, tenant, *inputTag)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	tagEntity := neo4jmapper.MapDbNodeToTagEntity(tagNodePtr)
	if tagEntity.Id == "" {
		err = errors.New("tag not saved")
		tracing.TraceErr(span, err)
		return nil, err
	}

	err = s.events.Publisher.PublishEvent(ctx, tagEntity.Id, model.TAG, dto.SaveTag{
		EntityType: utils.StringPtr(tagEntity.EntityType.String()),
		Name:       utils.StringPtr(tagEntity.Name),
		ColorCode:  utils.StringPtr(tagEntity.ColorCode),
	})

	return tagEntity, nil
}

func (s *tagService) AddTagToEntity(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, entityId string, entityType model.EntityType, tagId, tagName string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TagService.AddTagToEntity")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	tracing.TagEntity(span, entityId)
	span.LogFields(log.String("tagId", tagId), log.String("tagName", tagName), log.String("entityType", entityType.String()))

	// check tag exists by id
	tagByIdExists := false
	if tagId != "" {
		exists, err := s.neo4j.CommonReadRepository.ExistsById(ctx, common.GetTenantFromContext(ctx), tagId, model.NodeLabelTag)
		if err != nil {
			tracing.TraceErr(span, err)
			return "", err
		}
		tagByIdExists = exists
	}

	if !tagByIdExists && tagName != "" {
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

	err = s.neo4j.TagWriteRepository.LinkTagByIdToEntity(ctx, tx, tenant, tagId, entityId, entityType)
	if err != nil {
		tracing.TraceErr(span, err)
		return "", err
	}

	// event for tag added
	err = s.events.Publisher.PublishEvent(ctx, entityId, entityType, dto.NewAddTagEvent(tagId, tagName))
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "unable to publish message AddTagEvent"))
	}

	if common.GetAppSourceFromContext(ctx) != constants.AppSourceCustomerOsApi {
		s.events.Publisher.PublishEventCompleted(ctx, tenant, entityId, entityType, utils.NewEventCompletedDetails().WithUpdate())
	}

	return tagId, nil
}

func (s *tagService) RemoveTagFromEntity(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, entityId string, entityType model.EntityType, tagId, tagName string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TagService.RemoveTagFromEntity")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	tracing.TagEntity(span, entityId)
	span.LogFields(log.String("tagId", tagId), log.String("entityType", entityType.String()), log.String("tagName", tagName))

	var tagEntity *neo4jentity.TagEntity
	var err error

	if tagId != "" {
		tagEntity, err = s.GetById(ctx, tagId)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "unable to get tag by id"))
			return err
		}
	} else if tagName != "" {
		tagEntity, err = s.GetTagByEntityTypeAndName(ctx, entityType, tagName)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "unable to get tag by name"))
			return err
		}
	}

	if tagEntity == nil {
		err = errors.New("tag not found")
		tracing.TraceErr(span, err)
		return err
	}

	err = s.neo4j.TagWriteRepository.UnlinkTagByIdFromEntity(ctx, tx, tenant, tagEntity.Id, entityId, entityType)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "unable to unlink tag from entity"))
		return err
	}

	// event for tag removed
	err = s.events.Publisher.PublishEvent(ctx, entityId, entityType, dto.NewRemoveTagEvent(tagEntity.Id, tagEntity.Name))
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "unable to publish message RemoveTagEvent"))
	}

	if common.GetAppSourceFromContext(ctx) != constants.AppSourceCustomerOsApi {
		s.events.Publisher.PublishEventCompleted(ctx, tenant, entityId, entityType, utils.NewEventCompletedDetails().WithUpdate())
	}

	return nil
}

func (s *tagService) Update(ctx context.Context, tagId string, name, colorCode *string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TagService.Update")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	tracing.TagEntity(span, tagId)

	// validate tenant
	err := common.ValidateTenant(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	tenant := common.GetTenantFromContext(ctx)

	// validate tag exists
	exists, err := s.neo4j.CommonReadRepository.ExistsById(ctx, tenant, tagId, model.NodeLabelTag)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	if !exists {
		err = errors.New("tag not found")
		tracing.TraceErr(span, err)
		return err
	}

	err = s.neo4j.TagWriteRepository.Update(ctx, common.GetTenantFromContext(ctx), tagId, name, colorCode)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error updating tag name: %s", err.Error())
		return err
	}

	err = s.events.Publisher.PublishEvent(ctx, tagId, model.TAG, dto.SaveTag{
		Name:      name,
		ColorCode: colorCode,
	})

	return nil
}

func (s *tagService) UnlinkAndDelete(ctx context.Context, id string) (bool, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TagService.UnlinkAndDelete")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	err := s.neo4j.TagWriteRepository.UnlinkAllAndDelete(ctx, common.GetTenantFromContext(ctx), id)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s *tagService) GetAll(ctx context.Context) (*neo4jentity.TagEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TagService.GetAll")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	tagDbNodes, err := s.neo4j.TagReadRepository.GetAll(ctx, common.GetTenantFromContext(ctx))
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

	tags, err := s.neo4j.TagReadRepository.GetForContacts(ctx, common.GetTenantFromContext(ctx), contactIds)
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

	tags, err := s.neo4j.TagReadRepository.GetForIssues(ctx, common.GetTenantFromContext(ctx), issueIds)
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

	tags, err := s.neo4j.TagReadRepository.GetForOrganizations(ctx, common.GetTenantFromContext(ctx), organizationIDs)
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

	tags, err := s.neo4j.TagReadRepository.GetForLogEntries(ctx, common.GetTenantFromContext(ctx), logEntryIds)
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

	tagDbNode, err := s.neo4j.TagReadRepository.GetById(ctx, common.GetTenantFromContext(ctx), tagId)
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

	tagDbNodes, err := s.neo4j.TagReadRepository.GetAllByEntityType(ctx, common.GetTenantFromContext(ctx), entityType)
	if err != nil {
		return nil, err
	}
	tagEntities := make(neo4jentity.TagEntities, 0, len(tagDbNodes))
	for _, dbNodePtr := range tagDbNodes {
		tagEntities = append(tagEntities, *neo4jmapper.MapDbNodeToTagEntity(dbNodePtr))
	}
	return &tagEntities, nil
}
