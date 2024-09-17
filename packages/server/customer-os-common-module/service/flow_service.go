package service

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
)

type FlowService interface {
	FlowGetList(ctx context.Context) (*neo4jentity.FlowEntities, error)
	FlowGetById(ctx context.Context, id string) (*neo4jentity.FlowEntity, error)
	FlowGetByActionId(ctx context.Context, flowActionId string) (*neo4jentity.FlowEntity, error)
	FlowGetByContactId(ctx context.Context, flowContactId string) (*neo4jentity.FlowEntity, error)
	FlowsGetListWithContact(ctx context.Context, contactIds []string) (*neo4jentity.FlowEntities, error)
	FlowMerge(ctx context.Context, entity *neo4jentity.FlowEntity) (*neo4jentity.FlowEntity, error)
	FlowChangeStatus(ctx context.Context, id string, status neo4jentity.FlowStatus) (*neo4jentity.FlowEntity, error)

	FlowActionGetList(ctx context.Context, sequenceIds []string) (*neo4jentity.FlowActionEntities, error)
	FlowActionGetById(ctx context.Context, id string) (*neo4jentity.FlowActionEntity, error)
	FlowActionMerge(ctx context.Context, sequenceId string, entity *neo4jentity.FlowActionEntity) (*neo4jentity.FlowActionEntity, error)
	FlowActionChangeIndex(ctx context.Context, id string, index int64) error
	FlowActionChangeStatus(ctx context.Context, id string, status neo4jentity.FlowActionStatus) (*neo4jentity.FlowActionEntity, error)
	FlowActionDelete(ctx context.Context, id string) error

	FlowContactGetList(ctx context.Context, flowIds []string) (*neo4jentity.FlowContactEntities, error)
	FlowContactGetById(ctx context.Context, id string) (*neo4jentity.FlowContactEntity, error)
	FlowContactAdd(ctx context.Context, flowId, contactId string) (*neo4jentity.FlowContactEntity, error)
	FlowContactDelete(ctx context.Context, flowContactId string) error

	FlowActionSenderGetList(ctx context.Context, actionIds []string) (*neo4jentity.FlowActionSenderEntities, error)
	FlowActionSenderGetById(ctx context.Context, id string) (*neo4jentity.FlowActionSenderEntity, error)
	FlowActionSenderMerge(ctx context.Context, flowActionId string, input *neo4jentity.FlowActionSenderEntity) (*neo4jentity.FlowActionSenderEntity, error)
	FlowActionSenderDelete(ctx context.Context, flowActionSenderId string) error
}

type flowService struct {
	services *Services
}

func NewFlowService(services *Services) FlowService {
	return &flowService{
		services: services,
	}
}

func (s *flowService) FlowGetList(ctx context.Context) (*neo4jentity.FlowEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowService.FlowGetList")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	nodes, err := s.services.Neo4jRepositories.FlowReadRepository.GetList(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	entities := make(neo4jentity.FlowEntities, 0)
	for _, v := range nodes {
		e := mapper.MapDbNodeToFlowEntity(v)
		entities = append(entities, *e)
	}

	return &entities, nil
}

func (s *flowService) FlowGetById(ctx context.Context, id string) (*neo4jentity.FlowEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowService.FlowGetById")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	span.LogFields(log.String("id", id))

	node, err := s.services.Neo4jRepositories.FlowReadRepository.GetById(ctx, id)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return mapper.MapDbNodeToFlowEntity(node), nil
}

func (s *flowService) FlowGetByActionId(ctx context.Context, flowActionId string) (*neo4jentity.FlowEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowService.FlowGetByActionId")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	span.LogFields(log.String("flowActionId", flowActionId))

	node, err := s.services.Neo4jRepositories.FlowActionReadRepository.GetFlowByActionId(ctx, flowActionId)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	if node == nil {
		tracing.TraceErr(span, errors.New("flow not found"))
		return nil, errors.New("flow not found")
	}

	return mapper.MapDbNodeToFlowEntity(node), nil
}

func (s *flowService) FlowGetByContactId(ctx context.Context, flowContactId string) (*neo4jentity.FlowEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowService.FlowGetByContactId")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	span.LogFields(log.String("flowContactId", flowContactId))

	node, err := s.services.Neo4jRepositories.FlowActionReadRepository.GetFlowByContactId(ctx, flowContactId)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	if node == nil {
		tracing.TraceErr(span, errors.New("flow not found"))
		return nil, errors.New("flow not found")
	}

	return mapper.MapDbNodeToFlowEntity(node), nil
}

func (s *flowService) FlowsGetListWithContact(ctx context.Context, contactIds []string) (*neo4jentity.FlowEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowService.FlowsGetListWithContact")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	span.LogFields(log.Object("contactIds", contactIds))

	data, err := s.services.Neo4jRepositories.FlowReadRepository.GetListWithContact(ctx, contactIds)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	entities := make(neo4jentity.FlowEntities, 0)
	for _, v := range data {
		e := mapper.MapDbNodeToFlowEntity(v.Node)
		e.DataloaderKey = v.LinkedNodeId
		entities = append(entities, *e)
	}

	return &entities, nil
}

func (s *flowService) FlowMerge(ctx context.Context, input *neo4jentity.FlowEntity) (*neo4jentity.FlowEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowService.FlowMerge")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	toStore := &neo4jentity.FlowEntity{}
	var err error

	tenant := common.GetTenantFromContext(ctx)

	if input.Id == "" {
		toStore.Id, err = s.services.Neo4jRepositories.CommonReadRepository.GenerateId(ctx, tenant, model.NodeLabelFlow)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}

		toStore.Status = neo4jentity.FlowStatusInactive
	} else {
		toStore, err = s.FlowGetById(ctx, input.Id)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}

		if toStore == nil {
			tracing.TraceErr(span, errors.New("flow not found"))
			return nil, errors.New("flow not found")
		}
	}

	toStore.Name = input.Name
	toStore.Description = input.Description

	node, err := s.services.Neo4jRepositories.FlowWriteRepository.Merge(ctx, toStore)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	entity := mapper.MapDbNodeToFlowEntity(node)

	return entity, nil
}

func (s *flowService) FlowChangeStatus(ctx context.Context, id string, status neo4jentity.FlowStatus) (*neo4jentity.FlowEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowService.FlowChangeStatus")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	node, err := s.services.Neo4jRepositories.FlowReadRepository.GetById(ctx, id)
	if err != nil {
		return nil, err
	}

	if node == nil {
		tracing.TraceErr(span, errors.New("flow not found"))
		return nil, errors.New("flow not found")
	}

	flow := mapper.MapDbNodeToFlowEntity(node)

	if flow.Status == status {
		return flow, nil
	}

	flow.Status = status

	node, err = s.services.Neo4jRepositories.FlowWriteRepository.Merge(ctx, flow)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return mapper.MapDbNodeToFlowEntity(node), nil
}

func (s *flowService) FlowActionGetList(ctx context.Context, flowIds []string) (*neo4jentity.FlowActionEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowService.FlowActionGetList")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	span.LogFields(log.Object("flowIds", flowIds))

	nodes, err := s.services.Neo4jRepositories.FlowActionReadRepository.GetList(ctx, flowIds)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	entities := make(neo4jentity.FlowActionEntities, 0)
	for _, v := range nodes {
		e := mapper.MapDbNodeToFlowActionEntity(v.Node)
		e.DataloaderKey = v.LinkedNodeId
		entities = append(entities, *e)
	}

	return &entities, nil
}

func (s *flowService) FlowActionGetById(ctx context.Context, id string) (*neo4jentity.FlowActionEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowService.FlowActionGetById")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	span.LogFields(log.String("id", id))

	node, err := s.services.Neo4jRepositories.FlowActionReadRepository.GetById(ctx, id)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return mapper.MapDbNodeToFlowActionEntity(node), nil
}

func (s *flowService) FlowActionMerge(ctx context.Context, flowId string, input *neo4jentity.FlowActionEntity) (*neo4jentity.FlowActionEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowService.FlowActionMerge")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	toStore := &neo4jentity.FlowActionEntity{}
	var err error

	tenant := common.GetTenantFromContext(ctx)

	flow, err := s.FlowGetById(ctx, flowId)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	if flow == nil {
		tracing.TraceErr(span, errors.New("flow action not found"))
		return nil, errors.New("flow action not found")
	}

	if input.Id == "" {
		toStore.Id, err = s.services.Neo4jRepositories.CommonReadRepository.GenerateId(ctx, tenant, model.NodeLabelFlowAction)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}

		toStore.Status = neo4jentity.FlowActionStatusInactive
	} else {
		toStore, err = s.FlowActionGetById(ctx, input.Id)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}

		if toStore == nil {
			tracing.TraceErr(span, errors.New("flow sequence step not found"))
			return nil, errors.New("flow sequence step not found")
		}
	}

	exitingActions, err := s.FlowActionGetList(ctx, []string{flowId})
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	toStore.Index = int64(len(*exitingActions))
	toStore.Name = input.Name
	toStore.ActionType = input.ActionType
	toStore.ActionData = input.ActionData

	node, err := s.services.Neo4jRepositories.FlowActionWriteRepository.Merge(ctx, toStore)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	entity := mapper.MapDbNodeToFlowActionEntity(node)

	err = s.services.Neo4jRepositories.CommonWriteRepository.Link(ctx, nil, tenant, repository.LinkDetails{
		FromEntityId:   flowId,
		FromEntityType: model.FLOW,
		Relationship:   model.HAS,
		ToEntityId:     entity.Id,
		ToEntityType:   model.FLOW_ACTION,
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return entity, nil
}

// TODO this is not working. not ordering correctly
func (s *flowService) FlowActionChangeIndex(ctx context.Context, id string, index int64) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowService.FlowActionChangeIndex")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	action, err := s.FlowActionGetById(ctx, id)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if action == nil {
		tracing.TraceErr(span, errors.New("flow action not found"))
		return errors.New("flow action not found")
	}

	flow, err := s.FlowGetByActionId(ctx, id)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	if flow == nil {
		tracing.TraceErr(span, errors.New("flow not found"))
		return errors.New("flow not found")
	}

	existingActions, err := s.FlowActionGetList(ctx, []string{flow.Id})
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	for _, v := range *existingActions {
		if v.Index >= index {
			v.Index = v.Index + 1

			_, err = s.services.Neo4jRepositories.FlowActionWriteRepository.Merge(ctx, &v)
			if err != nil {
				tracing.TraceErr(span, err)
				return err
			}
		}
	}

	action.Index = index

	_, err = s.services.Neo4jRepositories.FlowActionWriteRepository.Merge(ctx, action)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}

func (s *flowService) FlowActionChangeStatus(ctx context.Context, id string, status neo4jentity.FlowActionStatus) (*neo4jentity.FlowActionEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowService.FlowActionChangeStatus")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	node, err := s.services.Neo4jRepositories.FlowActionReadRepository.GetById(ctx, id)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	if node == nil {
		tracing.TraceErr(span, errors.New("flow action not found"))
		return nil, errors.New("flow action")
	}

	entity := mapper.MapDbNodeToFlowActionEntity(node)

	if entity.Status == status {
		return entity, nil
	}

	entity.Status = status

	node, err = s.services.Neo4jRepositories.FlowActionWriteRepository.Merge(ctx, entity)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return mapper.MapDbNodeToFlowActionEntity(node), nil
}

func (s *flowService) FlowActionDelete(ctx context.Context, id string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowService.FlowActionDelete")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	tenant := common.GetTenantFromContext(ctx)

	node, err := s.services.Neo4jRepositories.FlowActionReadRepository.GetById(ctx, id)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if node == nil {
		tracing.TraceErr(span, errors.New("flow action not found"))
		return errors.New("flow action")
	}

	list, err := s.FlowActionSenderGetList(ctx, []string{id})
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	for _, v := range *list {
		err = s.FlowActionSenderDelete(ctx, v.Id)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}
	}

	flowEntity, err := s.FlowGetByActionId(ctx, id)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	err = s.services.Neo4jRepositories.CommonWriteRepository.Unlink(ctx, nil, tenant, repository.LinkDetails{
		FromEntityId:   flowEntity.Id,
		FromEntityType: model.FLOW,
		Relationship:   model.HAS,
		ToEntityId:     id,
		ToEntityType:   model.FLOW_ACTION,
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	err = s.services.Neo4jRepositories.FlowActionWriteRepository.Delete(ctx, id)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}

func (s *flowService) FlowContactGetList(ctx context.Context, sequenceIds []string) (*neo4jentity.FlowContactEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowService.FlowContactGetList")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	nodes, err := s.services.Neo4jRepositories.FlowContactReadRepository.GetList(ctx, sequenceIds)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	entities := make(neo4jentity.FlowContactEntities, 0)
	for _, v := range nodes {
		e := mapper.MapDbNodeToFlowContactEntity(v.Node)
		e.DataloaderKey = v.LinkedNodeId
		entities = append(entities, *e)
	}

	return &entities, nil
}

func (s *flowService) FlowContactGetById(ctx context.Context, id string) (*neo4jentity.FlowContactEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowService.FlowContactGetById")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	span.LogFields(log.String("id", id))

	node, err := s.services.Neo4jRepositories.FlowContactReadRepository.GetById(ctx, id)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return mapper.MapDbNodeToFlowContactEntity(node), nil
}

// (fs:FlowSequence {id: $id})-[:HAS]->(:FlowSequenceContact)
// (fs:FlowSequenceContact {id: $id})-[:HAS]->(:Contact)
// (fs:FlowSequenceContact {id: $id})-[:HAS]->(:Email)
func (s *flowService) FlowContactAdd(ctx context.Context, flowId, contactId string) (*neo4jentity.FlowContactEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowService.FlowContactLink")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	flow, err := s.FlowGetById(ctx, flowId)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	if flow == nil {
		tracing.TraceErr(span, errors.New("flow not found"))
		return nil, errors.New("flow not found")
	}

	identified, err := s.services.Neo4jRepositories.FlowContactReadRepository.Identify(ctx, flowId, contactId)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	if identified == nil {

		contactNode, err := s.services.Neo4jRepositories.ContactReadRepository.GetContact(ctx, common.GetTenantFromContext(ctx), contactId)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}
		if contactNode == nil {
			tracing.TraceErr(span, errors.New("contact not found"))
			return nil, errors.New("contact not found")
		}

		toStore := neo4jentity.FlowContactEntity{
			ContactId: contactId,
		}
		toStore.Id, err = s.services.Neo4jRepositories.CommonReadRepository.GenerateId(ctx, common.GetTenantFromContext(ctx), model.NodeLabelFlowContact)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}

		toStore.ContactId = contactId

		identified, err = s.services.Neo4jRepositories.FlowContactWriteRepository.Merge(ctx, &toStore)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}
		entity := mapper.MapDbNodeToFlowContactEntity(identified)

		err = s.services.Neo4jRepositories.CommonWriteRepository.Link(ctx, nil, common.GetTenantFromContext(ctx), repository.LinkDetails{
			FromEntityId:   flowId,
			FromEntityType: model.FLOW,
			Relationship:   model.HAS,
			ToEntityId:     entity.Id,
			ToEntityType:   model.FLOW_CONTACT,
		})
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}

		err = s.services.Neo4jRepositories.CommonWriteRepository.Link(ctx, nil, common.GetTenantFromContext(ctx), repository.LinkDetails{
			FromEntityId:   entity.Id,
			FromEntityType: model.FLOW_CONTACT,
			Relationship:   model.HAS,
			ToEntityId:     contactId,
			ToEntityType:   model.CONTACT,
		})
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}
	}

	return mapper.MapDbNodeToFlowContactEntity(identified), nil
}

func (s *flowService) FlowContactDelete(ctx context.Context, flowContactId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowService.FlowContactUnlink")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	span.LogFields(log.String("flowContactId", flowContactId))

	tenant := common.GetTenantFromContext(ctx)

	entity, err := s.FlowContactGetById(ctx, flowContactId)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if entity == nil {
		tracing.TraceErr(span, errors.New("flow sequence contact not found"))
		return errors.New("flow sequence contact not found")
	}

	flow, err := s.FlowGetByContactId(ctx, flowContactId)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	if flow == nil {
		tracing.TraceErr(span, errors.New("flow not found"))
		return errors.New("flow not found")
	}

	//todo use TX

	err = s.services.Neo4jRepositories.CommonWriteRepository.Unlink(ctx, nil, tenant, repository.LinkDetails{
		FromEntityId:   flow.Id,
		FromEntityType: model.FLOW,
		Relationship:   model.HAS,
		ToEntityId:     entity.Id,
		ToEntityType:   model.FLOW_CONTACT,
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	err = s.services.Neo4jRepositories.CommonWriteRepository.Unlink(ctx, nil, tenant, repository.LinkDetails{
		FromEntityId:   entity.Id,
		FromEntityType: model.FLOW_CONTACT,
		Relationship:   model.HAS,
		ToEntityId:     entity.ContactId,
		ToEntityType:   model.CONTACT,
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	err = s.services.Neo4jRepositories.FlowContactWriteRepository.Delete(ctx, entity.Id)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}

func (s *flowService) FlowActionSenderGetList(ctx context.Context, actionIds []string) (*neo4jentity.FlowActionSenderEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowService.FlowActionSenderGetList")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	nodes, err := s.services.Neo4jRepositories.FlowActionSenderReadRepository.GetList(ctx, actionIds)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	entities := make(neo4jentity.FlowActionSenderEntities, 0)
	for _, v := range nodes {
		e := mapper.MapDbNodeToFlowActionSenderEntity(v.Node)
		e.DataloaderKey = v.LinkedNodeId
		entities = append(entities, *e)
	}

	return &entities, nil
}

func (s *flowService) FlowActionSenderGetById(ctx context.Context, id string) (*neo4jentity.FlowActionSenderEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowService.FlowActionSenderGetById")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	span.LogFields(log.String("id", id))

	node, err := s.services.Neo4jRepositories.FlowActionSenderReadRepository.GetById(ctx, id)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return mapper.MapDbNodeToFlowActionSenderEntity(node), nil
}

// (fs:FlowSequence {id: $id})-[:HAS]->(:FlowActionSender)
func (s *flowService) FlowActionSenderMerge(ctx context.Context, flowActionId string, input *neo4jentity.FlowActionSenderEntity) (*neo4jentity.FlowActionSenderEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowService.FlowActionSenderMerge")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	action, err := s.FlowActionGetById(ctx, flowActionId)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	if action == nil {
		tracing.TraceErr(span, errors.New("flow action not found"))
		return nil, errors.New("flow action not found")
	}

	isNew := input.Id == ""
	var toStore *neo4jentity.FlowActionSenderEntity

	if input.Id == "" {
		toStore = &neo4jentity.FlowActionSenderEntity{}
		toStore.Id, err = s.services.Neo4jRepositories.CommonReadRepository.GenerateId(ctx, common.GetTenantFromContext(ctx), model.NodeLabelFlowActionSender)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}
	} else {
		toStore, err = s.FlowActionSenderGetById(ctx, input.Id)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}
		if toStore == nil {
			tracing.TraceErr(span, errors.New("flow action sender not found"))
			return nil, errors.New("flow action sender not found")
		}
	}

	toStore.Mailbox = input.Mailbox
	toStore.UserId = input.UserId

	node, err := s.services.Neo4jRepositories.FlowActionSenderWriteRepository.Merge(ctx, toStore)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	if isNew {
		err = s.services.Neo4jRepositories.CommonWriteRepository.Link(ctx, nil, common.GetTenantFromContext(ctx), repository.LinkDetails{
			FromEntityId:   flowActionId,
			FromEntityType: model.FLOW_ACTION,
			Relationship:   model.HAS,
			ToEntityId:     toStore.Id,
			ToEntityType:   model.FLOW_ACTION_SENDER,
		})
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}
	}

	//TODO LINK WITH USER and UNLINK WITH PREVIOUS IF CHANGED

	return mapper.MapDbNodeToFlowActionSenderEntity(node), nil
}

func (s *flowService) FlowActionSenderDelete(ctx context.Context, flowActionSenderId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowService.FlowActionSenderDelete")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	span.LogFields(log.String("flowActionSenderId", flowActionSenderId))

	actionSender, err := s.FlowActionSenderGetById(ctx, flowActionSenderId)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if actionSender == nil {
		tracing.TraceErr(span, errors.New("flow action sender not found"))
		return errors.New("flow action sender not found")
	}

	actionNode, err := s.services.Neo4jRepositories.FlowActionSenderReadRepository.GetFlowActionBySenderId(ctx, flowActionSenderId)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if actionNode == nil {
		tracing.TraceErr(span, errors.New("flow action not found"))
		return errors.New("flow action not found")
	}

	action := mapper.MapDbNodeToFlowActionEntity(actionNode)

	//todo use TX

	err = s.services.Neo4jRepositories.CommonWriteRepository.Unlink(ctx, nil, common.GetTenantFromContext(ctx), repository.LinkDetails{
		FromEntityId:   action.Id,
		FromEntityType: model.FLOW_ACTION,
		Relationship:   model.HAS,
		ToEntityId:     actionSender.Id,
		ToEntityType:   model.FLOW_ACTION_SENDER,
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	//TODO unlink from user
	err = s.services.Neo4jRepositories.FlowContactWriteRepository.Delete(ctx, actionSender.Id)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}

//func (s *flowService) GetFlowSequenceSenders(ctx context.Context, tenant, sequenceId string, page, limit int) (*utils.Pagination, error) {
//	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowService.GetFlowSequenceSenders")
//	defer span.Finish()
//	tracing.SetDefaultServiceSpanTags(ctx, span)
//
//	pageResult := utils.Pagination{
//		Page:  page,
//		Limit: limit,
//	}
//
//	count, err := s.services.PostgresRepositories.FlowSequenceSenderRepository.Count(ctx, tenant, sequenceId)
//	if err != nil {
//		return nil, err
//	}
//
//	entities, err := s.services.PostgresRepositories.FlowSequenceSenderRepository.Get(ctx, tenant, sequenceId, page, limit)
//	if err != nil {
//		return nil, err
//	}
//
//	pageResult.SetTotalRows(count)
//	pageResult.SetRows(entities)
//
//	return &pageResult, nil
//}
//
//func (s *flowService) GetFlowSequenceSenderById(ctx context.Context, tenant, id string) (*neo4jentity.FlowActionSenderEntity, error) {
//	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowService.GetFlowSequenceSenderById")
//	defer span.Finish()
//	tracing.SetDefaultServiceSpanTags(ctx, span)
//
//	span.LogFields(log.String("id", id))
//
//	entity, err := s.services.PostgresRepositories.FlowSequenceSenderRepository.GetById(ctx, tenant, id)
//	if err != nil {
//		return nil, err
//	}
//
//	return entity, nil
//}
//
//func (s *flowService) StoreFlowSequenceSender(ctx context.Context, tenant string, entity *neo4jentity.FlowActionSenderEntity) (*neo4jentity.FlowActionSenderEntity, error) {
//	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowService.StoreFlowSequenceSender")
//	defer span.Finish()
//	tracing.SetDefaultServiceSpanTags(ctx, span)
//
//	entity, err := s.services.PostgresRepositories.FlowSequenceSenderRepository.Store(ctx, tenant, entity)
//	if err != nil {
//		return nil, err
//	}
//
//	return entity, nil
//}
//
//func (s *flowService) DeleteFlowSequenceSender(ctx context.Context, tenant, id string) error {
//	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowService.DeleteFlowSequenceSender")
//	defer span.Finish()
//	tracing.SetDefaultServiceSpanTags(ctx, span)
//
//	span.LogFields(log.String("id", id))
//
//	err := s.services.PostgresRepositories.FlowSequenceSenderRepository.Delete(ctx, tenant, id)
//	if err != nil {
//		return err
//	}
//
//	return nil
//}
