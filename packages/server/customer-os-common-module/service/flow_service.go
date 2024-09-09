package service

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
)

type FlowService interface {
	FlowGetList(ctx context.Context) ([]*neo4jentity.FlowEntity, error)
	FlowGetById(ctx context.Context, id string) (*neo4jentity.FlowEntity, error)
	FlowGetBySequenceId(ctx context.Context, sequenceId string) (*neo4jentity.FlowEntity, error)
	FlowMerge(ctx context.Context, entity *neo4jentity.FlowEntity) (*neo4jentity.FlowEntity, error)
	FlowChangeStatus(ctx context.Context, id string, status neo4jentity.FlowStatus) (*neo4jentity.FlowEntity, error)

	FlowSequenceGetList(ctx context.Context, flowId *string) (*neo4jentity.FlowSequenceEntities, error)
	FlowSequenceGetById(ctx context.Context, id string) (*neo4jentity.FlowSequenceEntity, error)
	FlowSequenceCreate(ctx context.Context, flow *neo4jentity.FlowEntity, entity *neo4jentity.FlowSequenceEntity) (*neo4jentity.FlowSequenceEntity, error)
	FlowSequenceUpdate(ctx context.Context, entity *neo4jentity.FlowSequenceEntity) (*neo4jentity.FlowSequenceEntity, error)
	FlowSequenceChangeStatus(ctx context.Context, id string, status neo4jentity.FlowSequenceStatus) (*neo4jentity.FlowSequenceEntity, error)

	FlowSequenceStepGetList(ctx context.Context, sequenceIds []string) (*neo4jentity.FlowSequenceStepEntities, error)
	FlowSequenceStepGetById(ctx context.Context, id string) (*neo4jentity.FlowSequenceStepEntity, error)
	FlowSequenceStepCreate(ctx context.Context, sequenceId string, entity *neo4jentity.FlowSequenceStepEntity) (*neo4jentity.FlowSequenceStepEntity, error)
	FlowSequenceStepUpdate(ctx context.Context, entity *neo4jentity.FlowSequenceStepEntity) (*neo4jentity.FlowSequenceStepEntity, error)
	FlowSequenceStepChangeStatus(ctx context.Context, id string, status neo4jentity.FlowSequenceStepStatus) (*neo4jentity.FlowSequenceStepEntity, error)

	FlowSequenceContactGetList(ctx context.Context, sequenceIds []string) (*neo4jentity.FlowSequenceContactEntities, error)
	FlowSequenceContactGetById(ctx context.Context, id string) (*neo4jentity.FlowSequenceContactEntity, error)
	FlowSequenceContactLink(ctx context.Context, sequenceId, contactId, emailId string) (*neo4jentity.FlowSequenceContactEntity, error)
	FlowSequenceContactUnlink(ctx context.Context, sequenceId, contactId, emailId string) error

	FlowSequenceSenderGetList(ctx context.Context, sequenceIds []string) (*neo4jentity.FlowSequenceSenderEntities, error)
	FlowSequenceSenderGetById(ctx context.Context, id string) (*neo4jentity.FlowSequenceSenderEntity, error)
	FlowSequenceSenderLink(ctx context.Context, sequenceId, mailbox string) (*neo4jentity.FlowSequenceSenderEntity, error)
	FlowSequenceSenderUnlink(ctx context.Context, sequenceId, mailbox string) error
}

type flowService struct {
	services *Services
}

func NewFlowService(services *Services) FlowService {
	return &flowService{
		services: services,
	}
}

func (s *flowService) FlowGetList(ctx context.Context) ([]*neo4jentity.FlowEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowService.FlowGetList")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	nodes, err := s.services.Neo4jRepositories.FlowReadRepository.GetList(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	entities := make([]*neo4jentity.FlowEntity, len(nodes))
	for _, node := range nodes {
		entities = append(entities, mapper.MapDbNodeToFlowEntity(node))
	}

	return entities, nil
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

func (s *flowService) FlowGetBySequenceId(ctx context.Context, sequenceId string) (*neo4jentity.FlowEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowService.FlowGetBySequenceId")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	span.LogFields(log.String("sequenceId", sequenceId))

	node, err := s.services.Neo4jRepositories.FlowSequenceReadRepository.GetFlowBySequenceId(ctx, sequenceId)
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

func (s *flowService) FlowMerge(ctx context.Context, input *neo4jentity.FlowEntity) (*neo4jentity.FlowEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowService.FlowMerge")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	var err error

	if input.Id == "" {
		input.Id, err = s.services.Neo4jRepositories.CommonReadRepository.GenerateId(ctx, common.GetTenantFromContext(ctx), model.NodeLabelFlow)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}
	}

	node, err := s.services.Neo4jRepositories.FlowWriteRepository.Merge(ctx, input)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return mapper.MapDbNodeToFlowEntity(node), nil
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

func (s *flowService) FlowSequenceGetList(ctx context.Context, flowId *string) (*neo4jentity.FlowSequenceEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowService.FlowSequenceGetList")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	if flowId != nil {
		span.LogFields(log.String("flowId", *flowId))
	}

	data, err := s.services.Neo4jRepositories.FlowSequenceReadRepository.GetList(ctx, flowId)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	entities := make(neo4jentity.FlowSequenceEntities, 0, len(data))
	for _, v := range data {
		e := mapper.MapDbNodeToFlowSequenceEntity(v.Node)
		e.DataloaderKey = v.LinkedNodeId
		entities = append(entities, *e)
	}

	return &entities, nil
}

func (s *flowService) FlowSequenceGetById(ctx context.Context, id string) (*neo4jentity.FlowSequenceEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowService.FlowSequenceGetById")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	span.LogFields(log.String("id", id))

	node, err := s.services.Neo4jRepositories.FlowSequenceReadRepository.GetById(ctx, id)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return mapper.MapDbNodeToFlowSequenceEntity(node), nil
}

func (s *flowService) FlowSequenceCreate(ctx context.Context, flow *neo4jentity.FlowEntity, input *neo4jentity.FlowSequenceEntity) (*neo4jentity.FlowSequenceEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowService.FlowSequenceCreate")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	var err error
	tenant := common.GetTenantFromContext(ctx)
	flowId := ""

	if flow != nil && flow.Id != "" {
		f, err := s.services.Neo4jRepositories.FlowReadRepository.GetById(ctx, flow.Id)
		if err != nil {
			return nil, err
		}

		if f == nil {
			return nil, errors.New("flow not found")
		}

		flowId = flow.Id
	} else {
		flow, err := s.FlowMerge(ctx, &neo4jentity.FlowEntity{
			Name:   utils.FirstNotEmptyString(flow.Name, input.Name),
			Status: neo4jentity.FlowStatusInactive,
		})
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}
		flowId = flow.Id
	}

	input.Status = neo4jentity.FlowSequenceStatusInactive

	input.Id, err = s.services.Neo4jRepositories.CommonReadRepository.GenerateId(ctx, tenant, model.NodeLabelFlowSequence)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	node, err := s.services.Neo4jRepositories.FlowSequenceWriteRepository.Merge(ctx, input)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	flowSequence := mapper.MapDbNodeToFlowSequenceEntity(node)

	err = s.services.Neo4jRepositories.CommonWriteRepository.Link(ctx, nil, tenant, repository.LinkDetails{
		FromEntityId:   flowId,
		FromEntityType: model.FLOW,
		Relationship:   model.HAS,
		ToEntityId:     flowSequence.Id,
		ToEntityType:   model.FLOW_SEQUENCE,
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return flowSequence, nil
}

func (s *flowService) FlowSequenceUpdate(ctx context.Context, input *neo4jentity.FlowSequenceEntity) (*neo4jentity.FlowSequenceEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowService.FlowSequenceUpdate")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	if input.Id == "" {
		tracing.TraceErr(span, errors.New("id is required"))
		return nil, errors.New("id is required")
	}

	flowSequence, err := s.FlowSequenceGetById(ctx, input.Id)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	if flowSequence == nil {
		tracing.TraceErr(span, errors.New("flow sequence not found"))
		return nil, errors.New("flow sequence not found")
	}

	flowSequence.Name = input.Name
	flowSequence.Description = input.Description

	node, err := s.services.Neo4jRepositories.FlowSequenceWriteRepository.Merge(ctx, flowSequence)
	if err != nil {
		return nil, err
	}

	return mapper.MapDbNodeToFlowSequenceEntity(node), nil
}

func (s *flowService) FlowSequenceChangeStatus(ctx context.Context, id string, status neo4jentity.FlowSequenceStatus) (*neo4jentity.FlowSequenceEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowService.FlowSequenceChangeStatus")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	flowSequence, err := s.FlowSequenceGetById(ctx, id)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	if flowSequence == nil {
		tracing.TraceErr(span, errors.New("flow sequence not found"))
		return nil, errors.New("flow sequence not found")
	}

	if flowSequence.Status == status {
		return flowSequence, nil
	}

	flowSequence.Status = status

	node, err := s.services.Neo4jRepositories.FlowSequenceWriteRepository.Merge(ctx, flowSequence)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	flowSequence = mapper.MapDbNodeToFlowSequenceEntity(node)

	return flowSequence, nil
}

func (s *flowService) FlowSequenceStepGetList(ctx context.Context, sequenceIds []string) (*neo4jentity.FlowSequenceStepEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowService.FlowSequenceStepGetList")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	span.LogFields(log.Object("sequenceIds", sequenceIds))

	nodes, err := s.services.Neo4jRepositories.FlowSequenceStepReadRepository.GetList(ctx, sequenceIds)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	entities := make(neo4jentity.FlowSequenceStepEntities, 0, len(nodes))
	for _, v := range nodes {
		e := mapper.MapDbNodeToFlowSequenceStepEntity(v.Node)
		e.DataloaderKey = v.LinkedNodeId
		entities = append(entities, *e)
	}

	return &entities, nil
}

func (s *flowService) FlowSequenceStepGetById(ctx context.Context, id string) (*neo4jentity.FlowSequenceStepEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowService.FlowSequenceStepGetById")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	span.LogFields(log.String("id", id))

	node, err := s.services.Neo4jRepositories.FlowSequenceStepReadRepository.GetById(ctx, id)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return mapper.MapDbNodeToFlowSequenceStepEntity(node), nil
}

func (s *flowService) FlowSequenceStepCreate(ctx context.Context, sequenceId string, input *neo4jentity.FlowSequenceStepEntity) (*neo4jentity.FlowSequenceStepEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowService.FlowSequenceStepCreate")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	var err error

	tenant := common.GetTenantFromContext(ctx)

	input.Status = neo4jentity.FlowSequenceStepStatusInactive

	input.Id, err = s.services.Neo4jRepositories.CommonReadRepository.GenerateId(ctx, tenant, model.NodeLabelFlowSequence)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	node, err := s.services.Neo4jRepositories.FlowSequenceStepWriteRepository.Merge(ctx, input)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	entity := mapper.MapDbNodeToFlowSequenceStepEntity(node)

	err = s.services.Neo4jRepositories.CommonWriteRepository.Link(ctx, nil, tenant, repository.LinkDetails{
		FromEntityId:   sequenceId,
		FromEntityType: model.FLOW_SEQUENCE,
		Relationship:   model.HAS,
		ToEntityId:     entity.Id,
		ToEntityType:   model.FLOW_SEQUENCE_STEP,
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return entity, nil
}

func (s *flowService) FlowSequenceStepUpdate(ctx context.Context, input *neo4jentity.FlowSequenceStepEntity) (*neo4jentity.FlowSequenceStepEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowService.FlowSequenceStepUpdate")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	if input.Id == "" {
		tracing.TraceErr(span, errors.New("id is required"))
		return nil, errors.New("id is required")
	}

	flowSequenceStep, err := s.FlowSequenceStepGetById(ctx, input.Id)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	if flowSequenceStep == nil {
		tracing.TraceErr(span, errors.New("flow sequence step not found"))
		return nil, errors.New("flow sequence step not found")
	}

	flowSequenceStep.Name = input.Name

	node, err := s.services.Neo4jRepositories.FlowSequenceStepWriteRepository.Merge(ctx, flowSequenceStep)
	if err != nil {
		return nil, err
	}

	return mapper.MapDbNodeToFlowSequenceStepEntity(node), nil
}

func (s *flowService) FlowSequenceStepChangeStatus(ctx context.Context, id string, status neo4jentity.FlowSequenceStepStatus) (*neo4jentity.FlowSequenceStepEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowService.FlowSequenceStepChangeStatus")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	node, err := s.services.Neo4jRepositories.FlowSequenceStepReadRepository.GetById(ctx, id)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	if node == nil {
		tracing.TraceErr(span, errors.New("flow sequence step not found"))
		return nil, errors.New("flow not found")
	}

	entity := mapper.MapDbNodeToFlowSequenceStepEntity(node)

	if entity.Status == status {
		return entity, nil
	}

	entity.Status = status

	node, err = s.services.Neo4jRepositories.FlowSequenceStepWriteRepository.Merge(ctx, entity)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return mapper.MapDbNodeToFlowSequenceStepEntity(node), nil
}

func (s *flowService) FlowSequenceContactGetList(ctx context.Context, sequenceIds []string) (*neo4jentity.FlowSequenceContactEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowService.FlowSequenceContactGetList")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	nodes, err := s.services.Neo4jRepositories.FlowSequenceContactReadRepository.GetList(ctx, sequenceIds)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	entities := make(neo4jentity.FlowSequenceContactEntities, 0, len(nodes))
	for _, v := range nodes {
		e := mapper.MapDbNodeToFlowSequenceContactEntity(v.Node)
		e.DataloaderKey = v.LinkedNodeId
		entities = append(entities, *e)
	}

	return &entities, nil
}

func (s *flowService) FlowSequenceContactGetById(ctx context.Context, id string) (*neo4jentity.FlowSequenceContactEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowService.FlowSequenceContactGetById")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	span.LogFields(log.String("id", id))

	node, err := s.services.Neo4jRepositories.FlowSequenceContactReadRepository.GetById(ctx, id)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return mapper.MapDbNodeToFlowSequenceContactEntity(node), nil
}

// (fs:FlowSequence {id: $id})-[:HAS]->(:FlowSequenceContact)
// (fs:FlowSequenceContact {id: $id})-[:HAS]->(:Contact)
// (fs:FlowSequenceContact {id: $id})-[:HAS]->(:Email)
func (s *flowService) FlowSequenceContactLink(ctx context.Context, sequenceId, contactId, emailId string) (*neo4jentity.FlowSequenceContactEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowService.FlowSequenceContactLink")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	identified, err := s.services.Neo4jRepositories.FlowSequenceContactReadRepository.Identify(ctx, sequenceId, contactId, emailId)
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

		emailNode, err := s.services.Neo4jRepositories.EmailReadRepository.GetById(ctx, common.GetTenantFromContext(ctx), emailId)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}
		if emailNode == nil {
			tracing.TraceErr(span, errors.New("email not found"))
			return nil, errors.New("email not found")
		}

		toStore := neo4jentity.FlowSequenceContactEntity{
			ContactId: contactId,
			EmailId:   emailId,
		}
		toStore.Id, err = s.services.Neo4jRepositories.CommonReadRepository.GenerateId(ctx, common.GetTenantFromContext(ctx), model.NodeLabelFlowSequenceContact)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}

		identified, err := s.services.Neo4jRepositories.FlowSequenceContactWriteRepository.Merge(ctx, &toStore)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}
		entity := mapper.MapDbNodeToFlowSequenceContactEntity(identified)

		err = s.services.Neo4jRepositories.CommonWriteRepository.Link(ctx, nil, common.GetTenantFromContext(ctx), repository.LinkDetails{
			FromEntityId:   sequenceId,
			FromEntityType: model.FLOW_SEQUENCE,
			Relationship:   model.HAS,
			ToEntityId:     entity.Id,
			ToEntityType:   model.FLOW_SEQUENCE_CONTACT,
		})
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}

		err = s.services.Neo4jRepositories.CommonWriteRepository.Link(ctx, nil, common.GetTenantFromContext(ctx), repository.LinkDetails{
			FromEntityId:   entity.Id,
			FromEntityType: model.FLOW_SEQUENCE_CONTACT,
			Relationship:   model.HAS,
			ToEntityId:     contactId,
			ToEntityType:   model.CONTACT,
		})
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}

		err = s.services.Neo4jRepositories.CommonWriteRepository.Link(ctx, nil, common.GetTenantFromContext(ctx), repository.LinkDetails{
			FromEntityId:   entity.Id,
			FromEntityType: model.FLOW_SEQUENCE_CONTACT,
			Relationship:   model.HAS,
			ToEntityId:     emailId,
			ToEntityType:   model.EMAIL,
		})
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}
	}

	return mapper.MapDbNodeToFlowSequenceContactEntity(identified), nil
}

func (s *flowService) FlowSequenceContactUnlink(ctx context.Context, sequenceId, contactId, emailId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowService.FlowSequenceContactUnlink")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	span.LogFields(log.String("sequenceId", sequenceId), log.String("contactId", contactId), log.String("emailId", emailId))

	node, err := s.services.Neo4jRepositories.FlowSequenceContactReadRepository.Identify(ctx, sequenceId, contactId, emailId)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if node == nil {
		tracing.TraceErr(span, errors.New("flow sequence contact not found"))
		return errors.New("flow sequence contact not found")
	}

	entity := mapper.MapDbNodeToFlowSequenceContactEntity(node)

	//todo use TX

	err = s.services.Neo4jRepositories.CommonWriteRepository.Unlink(ctx, nil, common.GetTenantFromContext(ctx), repository.LinkDetails{
		FromEntityId:   sequenceId,
		FromEntityType: model.FLOW_SEQUENCE,
		Relationship:   model.HAS,
		ToEntityId:     entity.Id,
		ToEntityType:   model.FLOW_SEQUENCE_CONTACT,
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	err = s.services.Neo4jRepositories.CommonWriteRepository.Unlink(ctx, nil, common.GetTenantFromContext(ctx), repository.LinkDetails{
		FromEntityId:   entity.Id,
		FromEntityType: model.FLOW_SEQUENCE_CONTACT,
		Relationship:   model.HAS,
		ToEntityId:     contactId,
		ToEntityType:   model.CONTACT,
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	err = s.services.Neo4jRepositories.CommonWriteRepository.Unlink(ctx, nil, common.GetTenantFromContext(ctx), repository.LinkDetails{
		FromEntityId:   entity.Id,
		FromEntityType: model.FLOW_SEQUENCE_CONTACT,
		Relationship:   model.HAS,
		ToEntityId:     emailId,
		ToEntityType:   model.EMAIL,
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	err = s.services.Neo4jRepositories.CommonWriteRepository.Delete(ctx, nil, common.GetTenantFromContext(ctx), entity.Id, model.NodeLabelFlowSequenceContact)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}

func (s *flowService) FlowSequenceSenderGetList(ctx context.Context, sequenceIds []string) (*neo4jentity.FlowSequenceSenderEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowService.FlowSequenceSenderGetList")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	nodes, err := s.services.Neo4jRepositories.FlowSequenceSenderReadRepository.GetList(ctx, sequenceIds)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	entities := make(neo4jentity.FlowSequenceSenderEntities, 0, len(nodes))
	for _, v := range nodes {
		e := mapper.MapDbNodeToFlowSequenceSenderEntity(v.Node)
		e.DataloaderKey = v.LinkedNodeId
		entities = append(entities, *e)
	}

	return &entities, nil
}

func (s *flowService) FlowSequenceSenderGetById(ctx context.Context, id string) (*neo4jentity.FlowSequenceSenderEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowService.FlowSequenceSenderGetById")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	span.LogFields(log.String("id", id))

	node, err := s.services.Neo4jRepositories.FlowSequenceSenderReadRepository.GetById(ctx, id)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return mapper.MapDbNodeToFlowSequenceSenderEntity(node), nil
}

// (fs:FlowSequence {id: $id})-[:HAS]->(:FlowSequenceSender)
func (s *flowService) FlowSequenceSenderLink(ctx context.Context, sequenceId, mailbox string) (*neo4jentity.FlowSequenceSenderEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowService.FlowSequenceSenderLink")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	identified, err := s.services.Neo4jRepositories.FlowSequenceSenderReadRepository.Identify(ctx, sequenceId, mailbox)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	if identified == nil {

		m, err := s.services.PostgresRepositories.TenantSettingsMailboxRepository.GetByMailbox(ctx, common.GetTenantFromContext(ctx), mailbox)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}

		if m == nil {
			tracing.TraceErr(span, errors.New("mailbox not found"))
			return nil, errors.New("mailbox not found")
		}

		toStore := neo4jentity.FlowSequenceSenderEntity{
			Mailbox: mailbox,
		}
		toStore.Id, err = s.services.Neo4jRepositories.CommonReadRepository.GenerateId(ctx, common.GetTenantFromContext(ctx), model.NodeLabelFlowSequenceSender)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}

		identified, err = s.services.Neo4jRepositories.FlowSequenceSenderWriteRepository.Merge(ctx, &toStore)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}
		entity := mapper.MapDbNodeToFlowSequenceSenderEntity(identified)

		err = s.services.Neo4jRepositories.CommonWriteRepository.Link(ctx, nil, common.GetTenantFromContext(ctx), repository.LinkDetails{
			FromEntityId:   sequenceId,
			FromEntityType: model.FLOW_SEQUENCE,
			Relationship:   model.HAS,
			ToEntityId:     entity.Id,
			ToEntityType:   model.FLOW_SEQUENCE_SENDER,
		})
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}
	}

	return mapper.MapDbNodeToFlowSequenceSenderEntity(identified), nil
}

func (s *flowService) FlowSequenceSenderUnlink(ctx context.Context, sequenceId, mailbox string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowService.FlowSequenceSenderUnlink")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	span.LogFields(log.String("sequenceId", sequenceId), log.String("mailbox", mailbox))

	node, err := s.services.Neo4jRepositories.FlowSequenceSenderReadRepository.Identify(ctx, sequenceId, mailbox)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if node == nil {
		tracing.TraceErr(span, errors.New("flow sequence sender not found"))
		return errors.New("flow sequence sender not found")
	}

	entity := mapper.MapDbNodeToFlowSequenceContactEntity(node)

	//todo use TX

	err = s.services.Neo4jRepositories.CommonWriteRepository.Unlink(ctx, nil, common.GetTenantFromContext(ctx), repository.LinkDetails{
		FromEntityId:   sequenceId,
		FromEntityType: model.FLOW_SEQUENCE,
		Relationship:   model.HAS,
		ToEntityId:     entity.Id,
		ToEntityType:   model.FLOW_SEQUENCE_SENDER,
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	err = s.services.Neo4jRepositories.CommonWriteRepository.Delete(ctx, nil, common.GetTenantFromContext(ctx), entity.Id, model.NodeLabelFlowSequenceContact)
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
//func (s *flowService) GetFlowSequenceSenderById(ctx context.Context, tenant, id string) (*neo4jentity.FlowSequenceSenderEntity, error) {
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
//func (s *flowService) StoreFlowSequenceSender(ctx context.Context, tenant string, entity *neo4jentity.FlowSequenceSenderEntity) (*neo4jentity.FlowSequenceSenderEntity, error) {
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
