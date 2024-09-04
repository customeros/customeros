package service

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	postgresEntity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
)

type FlowService interface {
	FlowGetList(ctx context.Context) ([]*postgresEntity.Flow, error)
	FlowGetById(ctx context.Context, id string) (*postgresEntity.Flow, error)
	FlowGetBySequenceId(ctx context.Context, sequenceId string) (*postgresEntity.Flow, error)
	FlowStore(ctx context.Context, entity *postgresEntity.Flow) (*postgresEntity.Flow, error)
	FlowChangeStatus(ctx context.Context, id string, status postgresEntity.FlowStatus) (*postgresEntity.Flow, error)

	FlowSequenceGetList(ctx context.Context, flowId *string) ([]*postgresEntity.FlowSequence, error)
	FlowSequenceGetById(ctx context.Context, id string) (*postgresEntity.FlowSequence, error)
	FlowSequenceStore(ctx context.Context, entity *postgresEntity.FlowSequence) (*postgresEntity.FlowSequence, error)
	FlowSequenceChangeStatus(ctx context.Context, id string, status postgresEntity.FlowSequenceStatus) (*postgresEntity.FlowSequence, error)

	FlowSequenceStepGetList(ctx context.Context, sequenceId string) ([]*postgresEntity.FlowSequenceStep, error)
	FlowSequenceStepGetById(ctx context.Context, id string) (*postgresEntity.FlowSequenceStep, error)
	FlowSequenceStepStore(ctx context.Context, entity *postgresEntity.FlowSequenceStep) (*postgresEntity.FlowSequenceStep, error)
	FlowSequenceStepChangeStatus(ctx context.Context, id string, status postgresEntity.FlowSequenceStepStatus) (*postgresEntity.FlowSequenceStep, error)

	GetFlowSequenceContacts(ctx context.Context, tenant, sequenceId string, page, limit int) (*utils.Pagination, error)
	GetFlowSequenceContactById(ctx context.Context, tenant, id string) (*postgresEntity.FlowSequenceContact, error)
	StoreFlowSequenceContact(ctx context.Context, tenant string, entity *postgresEntity.FlowSequenceContact) (*postgresEntity.FlowSequenceContact, error)
	DeleteFlowSequenceContact(ctx context.Context, tenant, id string) error

	GetFlowSequenceSenders(ctx context.Context, tenant, sequenceId string, page, limit int) (*utils.Pagination, error)
	GetFlowSequenceSenderById(ctx context.Context, tenant, id string) (*postgresEntity.FlowSequenceSender, error)
	StoreFlowSequenceSender(ctx context.Context, tenant string, entity *postgresEntity.FlowSequenceSender) (*postgresEntity.FlowSequenceSender, error)
	DeleteFlowSequenceSender(ctx context.Context, tenant, id string) error
}

type flowService struct {
	services *Services
}

func NewFlowService(services *Services) FlowService {
	return &flowService{
		services: services,
	}
}

func (s *flowService) FlowGetList(ctx context.Context) ([]*postgresEntity.Flow, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowService.FlowGetList")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	entities, err := s.services.PostgresRepositories.FlowRepository.GetList(ctx)
	if err != nil {
		return nil, err
	}

	return entities, nil
}

func (s *flowService) FlowGetById(ctx context.Context, id string) (*postgresEntity.Flow, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowService.FlowGetById")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	span.LogFields(log.String("id", id))

	entity, err := s.services.PostgresRepositories.FlowRepository.GetById(ctx, id)
	if err != nil {
		return nil, err
	}

	return entity, nil
}

func (s *flowService) FlowGetBySequenceId(ctx context.Context, sequenceId string) (*postgresEntity.Flow, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowService.FlowGetBySequenceId")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	span.LogFields(log.String("sequenceId", sequenceId))

	entity, err := s.services.PostgresRepositories.FlowSequenceRepository.GetFlowBySequenceId(ctx, sequenceId)
	if err != nil {
		return nil, err
	}

	return entity, nil
}

func (s *flowService) FlowStore(ctx context.Context, entity *postgresEntity.Flow) (*postgresEntity.Flow, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowService.FlowStore")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	entity, err := s.services.PostgresRepositories.FlowRepository.Store(ctx, entity)
	if err != nil {
		return nil, err
	}

	return entity, nil
}

func (s *flowService) FlowChangeStatus(ctx context.Context, id string, status postgresEntity.FlowStatus) (*postgresEntity.Flow, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowService.FlowChangeStatus")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	entity, err := s.services.PostgresRepositories.FlowRepository.GetById(ctx, id)
	if err != nil {
		return nil, err
	}

	if entity == nil {
		tracing.TraceErr(span, errors.New("flow not found"))
		return nil, errors.New("flow not found")
	}

	if entity.Status == status {
		return entity, nil
	}

	entity.Status = status

	entity, err = s.services.PostgresRepositories.FlowRepository.Store(ctx, entity)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return entity, nil
}

func (s *flowService) FlowSequenceGetList(ctx context.Context, flowId *string) ([]*postgresEntity.FlowSequence, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowService.FlowSequenceGetList")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	if flowId != nil {
		span.LogFields(log.String("flowId", *flowId))
	}

	entities, err := s.services.PostgresRepositories.FlowSequenceRepository.GetList(ctx, flowId)
	if err != nil {
		return nil, err
	}

	return entities, nil
}

func (s *flowService) FlowSequenceGetById(ctx context.Context, id string) (*postgresEntity.FlowSequence, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowService.FlowSequenceGetById")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	span.LogFields(log.String("id", id))

	entity, err := s.services.PostgresRepositories.FlowSequenceRepository.GetById(ctx, id)
	if err != nil {
		return nil, err
	}

	return entity, nil
}

func (s *flowService) FlowSequenceStore(ctx context.Context, entity *postgresEntity.FlowSequence) (*postgresEntity.FlowSequence, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowService.FlowSequenceStore")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	entity, err := s.services.PostgresRepositories.FlowSequenceRepository.Store(ctx, entity)
	if err != nil {
		return nil, err
	}

	return entity, nil
}

func (s *flowService) FlowSequenceChangeStatus(ctx context.Context, id string, status postgresEntity.FlowSequenceStatus) (*postgresEntity.FlowSequence, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowService.FlowSequenceChangeStatus")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	entity, err := s.services.PostgresRepositories.FlowSequenceRepository.GetById(ctx, id)
	if err != nil {
		return nil, err
	}

	if entity == nil {
		tracing.TraceErr(span, errors.New("flow not found"))
		return nil, errors.New("flow not found")
	}

	if entity.Status == status {
		return entity, nil
	}

	entity.Status = status

	entity, err = s.services.PostgresRepositories.FlowSequenceRepository.Store(ctx, entity)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return entity, nil
}

func (s *flowService) FlowSequenceStepGetList(ctx context.Context, sequenceId string) ([]*postgresEntity.FlowSequenceStep, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowService.FlowSequenceStepGetList")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	span.LogFields(log.String("sequenceId", sequenceId))

	entities, err := s.services.PostgresRepositories.FlowSequenceStepRepository.GetList(ctx, sequenceId)
	if err != nil {
		return nil, err
	}

	return entities, nil
}

func (s *flowService) FlowSequenceStepGetById(ctx context.Context, id string) (*postgresEntity.FlowSequenceStep, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowService.FlowSequenceStepGetById")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	span.LogFields(log.String("id", id))

	entity, err := s.services.PostgresRepositories.FlowSequenceStepRepository.GetById(ctx, id)
	if err != nil {
		return nil, err
	}

	return entity, nil
}

func (s *flowService) FlowSequenceStepStore(ctx context.Context, entity *postgresEntity.FlowSequenceStep) (*postgresEntity.FlowSequenceStep, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowService.FlowSequenceStepStore")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	entity, err := s.services.PostgresRepositories.FlowSequenceStepRepository.Store(ctx, entity)
	if err != nil {
		return nil, err
	}

	return entity, nil
}

func (s *flowService) FlowSequenceStepChangeStatus(ctx context.Context, id string, status postgresEntity.FlowSequenceStepStatus) (*postgresEntity.FlowSequenceStep, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowService.FlowSequenceStepChangeStatus")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	entity, err := s.services.PostgresRepositories.FlowSequenceStepRepository.GetById(ctx, id)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	if entity == nil {
		tracing.TraceErr(span, errors.New("flow not found"))
		return nil, errors.New("flow not found")
	}

	if entity.Status == status {
		return entity, nil
	}

	entity.Status = status

	entity, err = s.services.PostgresRepositories.FlowSequenceStepRepository.Store(ctx, entity)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return entity, nil
}

func (s *flowService) GetFlowSequenceContacts(ctx context.Context, tenant, sequenceId string, page, limit int) (*utils.Pagination, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowService.GetFlowSequenceContacts")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	pageResult := utils.Pagination{
		Page:  page,
		Limit: limit,
	}

	count, err := s.services.PostgresRepositories.FlowSequenceContactRepository.Count(ctx, tenant, sequenceId)
	if err != nil {
		return nil, err
	}

	entities, err := s.services.PostgresRepositories.FlowSequenceContactRepository.Get(ctx, tenant, sequenceId, page, limit)
	if err != nil {
		return nil, err
	}

	pageResult.SetTotalRows(count)
	pageResult.SetRows(entities)

	return &pageResult, nil
}

func (s *flowService) GetFlowSequenceContactById(ctx context.Context, tenant, id string) (*postgresEntity.FlowSequenceContact, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowService.GetFlowSequenceContactById")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	span.LogFields(log.String("id", id))

	entity, err := s.services.PostgresRepositories.FlowSequenceContactRepository.GetById(ctx, tenant, id)
	if err != nil {
		return nil, err
	}

	return entity, nil
}

func (s *flowService) StoreFlowSequenceContact(ctx context.Context, tenant string, entity *postgresEntity.FlowSequenceContact) (*postgresEntity.FlowSequenceContact, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowService.StoreFlowSequenceContact")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	entity, err := s.services.PostgresRepositories.FlowSequenceContactRepository.Store(ctx, tenant, entity)
	if err != nil {
		return nil, err
	}

	return entity, nil
}

func (s *flowService) DeleteFlowSequenceContact(ctx context.Context, tenant, id string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowService.DeleteFlowSequenceContact")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	span.LogFields(log.String("id", id))

	err := s.services.PostgresRepositories.FlowSequenceContactRepository.Delete(ctx, tenant, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *flowService) GetFlowSequenceSenders(ctx context.Context, tenant, sequenceId string, page, limit int) (*utils.Pagination, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowService.GetFlowSequenceSenders")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	pageResult := utils.Pagination{
		Page:  page,
		Limit: limit,
	}

	count, err := s.services.PostgresRepositories.FlowSequenceSenderRepository.Count(ctx, tenant, sequenceId)
	if err != nil {
		return nil, err
	}

	entities, err := s.services.PostgresRepositories.FlowSequenceSenderRepository.Get(ctx, tenant, sequenceId, page, limit)
	if err != nil {
		return nil, err
	}

	pageResult.SetTotalRows(count)
	pageResult.SetRows(entities)

	return &pageResult, nil
}

func (s *flowService) GetFlowSequenceSenderById(ctx context.Context, tenant, id string) (*postgresEntity.FlowSequenceSender, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowService.GetFlowSequenceSenderById")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	span.LogFields(log.String("id", id))

	entity, err := s.services.PostgresRepositories.FlowSequenceSenderRepository.GetById(ctx, tenant, id)
	if err != nil {
		return nil, err
	}

	return entity, nil
}

func (s *flowService) StoreFlowSequenceSender(ctx context.Context, tenant string, entity *postgresEntity.FlowSequenceSender) (*postgresEntity.FlowSequenceSender, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowService.StoreFlowSequenceSender")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	entity, err := s.services.PostgresRepositories.FlowSequenceSenderRepository.Store(ctx, tenant, entity)
	if err != nil {
		return nil, err
	}

	return entity, nil
}

func (s *flowService) DeleteFlowSequenceSender(ctx context.Context, tenant, id string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowService.DeleteFlowSequenceSender")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	span.LogFields(log.String("id", id))

	err := s.services.PostgresRepositories.FlowSequenceSenderRepository.Delete(ctx, tenant, id)
	if err != nil {
		return err
	}

	return nil
}
