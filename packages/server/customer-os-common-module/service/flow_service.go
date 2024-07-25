package service

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	postgresEntity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type FlowService interface {
	GetFlows(ctx context.Context, tenant string) ([]*postgresEntity.Flow, error)
	GetFlowById(ctx context.Context, tenant, id string) (*postgresEntity.Flow, error)
	StoreFlow(ctx context.Context, entity *postgresEntity.Flow) (*postgresEntity.Flow, error)
	ActivateFlow(ctx context.Context, tenant, id string) error
	DeactivateFlow(ctx context.Context, tenant, id string) error
	DeleteFlow(ctx context.Context, tenant, id string) error

	GetFlowSequences(ctx context.Context, tenant, flowId string) ([]*postgresEntity.FlowSequence, error)
	GetFlowSequenceById(ctx context.Context, tenant, id string) (*postgresEntity.FlowSequence, error)
	StoreFlowSequence(ctx context.Context, tenant string, entity *postgresEntity.FlowSequence) (*postgresEntity.FlowSequence, error)
	ActivateFlowSequence(ctx context.Context, tenant, id string) error
	DeactivateFlowSequence(ctx context.Context, tenant, id string) error
	DeleteFlowSequence(ctx context.Context, tenant, id string) error

	GetFlowSequenceSteps(ctx context.Context, tenant, sequenceId string) ([]*postgresEntity.FlowSequenceStep, error)
	GetFlowSequenceStepById(ctx context.Context, tenant, id string) (*postgresEntity.FlowSequenceStep, error)
	StoreFlowSequenceStep(ctx context.Context, tenant string, entity *postgresEntity.FlowSequenceStep) (*postgresEntity.FlowSequenceStep, error)
	ActivateFlowSequenceStep(ctx context.Context, tenant, id string) error
	DeactivateFlowSequenceStep(ctx context.Context, tenant, id string) error
	DeleteFlowSequenceStep(ctx context.Context, tenant, id string) error

	GetFlowSequenceContacts(ctx context.Context, tenant, sequenceId string) ([]*postgresEntity.FlowSequenceContact, error)
	GetFlowSequenceContactById(ctx context.Context, tenant, id string) (*postgresEntity.FlowSequenceContact, error)
	StoreFlowSequenceContact(ctx context.Context, tenant string, entity *postgresEntity.FlowSequenceContact) (*postgresEntity.FlowSequenceContact, error)
	DeleteFlowSequenceContact(ctx context.Context, tenant, id string) error

	GetFlowSequenceSenders(ctx context.Context, tenant, sequenceId string) ([]*postgresEntity.FlowSequenceSender, error)
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

func (s *flowService) GetFlows(ctx context.Context, tenant string) ([]*postgresEntity.Flow, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowService.GetFlows")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	span.SetTag(tracing.SpanTagComponent, constants.ComponentService)

	entities, err := s.services.PostgresRepositories.FlowRepository.Get(ctx, tenant)
	if err != nil {
		return nil, err
	}

	return entities, nil
}

func (s *flowService) GetFlowById(ctx context.Context, tenant, id string) (*postgresEntity.Flow, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowService.GetFlowById")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	span.SetTag(tracing.SpanTagComponent, constants.ComponentService)

	span.LogFields(log.String("id", id))

	entity, err := s.services.PostgresRepositories.FlowRepository.GetById(ctx, tenant, id)
	if err != nil {
		return nil, err
	}

	return entity, nil
}

func (s *flowService) StoreFlow(ctx context.Context, entity *postgresEntity.Flow) (*postgresEntity.Flow, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowService.StoreFlow")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, entity.Tenant)
	span.SetTag(tracing.SpanTagComponent, constants.ComponentService)

	entity, err := s.services.PostgresRepositories.FlowRepository.Store(ctx, entity)
	if err != nil {
		return nil, err
	}

	return entity, nil
}

func (s *flowService) ActivateFlow(ctx context.Context, tenant, id string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowService.ActivateFlow")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	span.SetTag(tracing.SpanTagComponent, constants.ComponentService)

	entity, err := s.services.PostgresRepositories.FlowRepository.GetById(ctx, tenant, id)
	if err != nil {
		return err
	}

	if entity == nil {
		return nil
	}

	entity.Active = true

	_, err = s.services.PostgresRepositories.FlowRepository.Store(ctx, entity)
	if err != nil {
		return err
	}

	return nil
}

func (s *flowService) DeactivateFlow(ctx context.Context, tenant, id string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowService.DeactivateFlow")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	span.SetTag(tracing.SpanTagComponent, constants.ComponentService)

	entity, err := s.services.PostgresRepositories.FlowRepository.GetById(ctx, tenant, id)
	if err != nil {
		return err
	}

	if entity == nil {
		return nil
	}

	entity.Active = false

	_, err = s.services.PostgresRepositories.FlowRepository.Store(ctx, entity)
	if err != nil {
		return err
	}

	return nil
}

func (s *flowService) DeleteFlow(ctx context.Context, tenant, id string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowService.DeleteFlow")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	span.SetTag(tracing.SpanTagComponent, constants.ComponentService)
	span.LogFields(log.String("id", id))

	err := s.services.PostgresRepositories.FlowRepository.Delete(ctx, tenant, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *flowService) GetFlowSequences(ctx context.Context, tenant, flowId string) ([]*postgresEntity.FlowSequence, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowService.GetFlowSequences")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	span.SetTag(tracing.SpanTagComponent, constants.ComponentService)

	entities, err := s.services.PostgresRepositories.FlowSequenceRepository.Get(ctx, tenant, flowId)
	if err != nil {
		return nil, err
	}

	return entities, nil
}

func (s *flowService) GetFlowSequenceById(ctx context.Context, tenant, id string) (*postgresEntity.FlowSequence, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowService.GetFlowSequenceById")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	span.SetTag(tracing.SpanTagComponent, constants.ComponentService)

	span.LogFields(log.String("id", id))

	entity, err := s.services.PostgresRepositories.FlowSequenceRepository.GetById(ctx, tenant, id)
	if err != nil {
		return nil, err
	}

	return entity, nil
}

func (s *flowService) StoreFlowSequence(ctx context.Context, tenant string, entity *postgresEntity.FlowSequence) (*postgresEntity.FlowSequence, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowService.StoreFlowSequence")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	span.SetTag(tracing.SpanTagComponent, constants.ComponentService)

	entity, err := s.services.PostgresRepositories.FlowSequenceRepository.Store(ctx, tenant, entity)
	if err != nil {
		return nil, err
	}

	return entity, nil
}

func (s *flowService) ActivateFlowSequence(ctx context.Context, tenant, id string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowService.ActivateFlowSequence")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	span.SetTag(tracing.SpanTagComponent, constants.ComponentService)

	entity, err := s.services.PostgresRepositories.FlowSequenceRepository.GetById(ctx, tenant, id)
	if err != nil {
		return err
	}

	if entity == nil {
		return nil
	}

	entity.Active = true

	_, err = s.services.PostgresRepositories.FlowSequenceRepository.Store(ctx, tenant, entity)
	if err != nil {
		return err
	}

	return nil
}

func (s *flowService) DeactivateFlowSequence(ctx context.Context, tenant, id string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowService.DeactivateFlowSequence")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	span.SetTag(tracing.SpanTagComponent, constants.ComponentService)

	entity, err := s.services.PostgresRepositories.FlowSequenceRepository.GetById(ctx, tenant, id)
	if err != nil {
		return err
	}

	if entity == nil {
		return nil
	}

	entity.Active = false

	_, err = s.services.PostgresRepositories.FlowSequenceRepository.Store(ctx, tenant, entity)
	if err != nil {
		return err
	}

	return nil
}

func (s *flowService) DeleteFlowSequence(ctx context.Context, tenant, id string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowService.DeleteFlowSequence")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	span.SetTag(tracing.SpanTagComponent, constants.ComponentService)
	span.LogFields(log.String("id", id))

	err := s.services.PostgresRepositories.FlowSequenceRepository.Delete(ctx, tenant, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *flowService) GetFlowSequenceSteps(ctx context.Context, tenant, sequenceId string) ([]*postgresEntity.FlowSequenceStep, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowService.GetFlowSequenceSteps")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	span.SetTag(tracing.SpanTagComponent, constants.ComponentService)

	entities, err := s.services.PostgresRepositories.FlowSequenceStepRepository.Get(ctx, tenant, sequenceId)
	if err != nil {
		return nil, err
	}

	return entities, nil
}

func (s *flowService) GetFlowSequenceStepById(ctx context.Context, tenant, id string) (*postgresEntity.FlowSequenceStep, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowService.GetFlowSequenceStepById")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	span.SetTag(tracing.SpanTagComponent, constants.ComponentService)

	span.LogFields(log.String("id", id))

	entity, err := s.services.PostgresRepositories.FlowSequenceStepRepository.GetById(ctx, tenant, id)
	if err != nil {
		return nil, err
	}

	return entity, nil
}

func (s *flowService) StoreFlowSequenceStep(ctx context.Context, tenant string, entity *postgresEntity.FlowSequenceStep) (*postgresEntity.FlowSequenceStep, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowService.StoreFlowSequenceStep")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	span.SetTag(tracing.SpanTagComponent, constants.ComponentService)

	entity, err := s.services.PostgresRepositories.FlowSequenceStepRepository.Store(ctx, tenant, entity)
	if err != nil {
		return nil, err
	}

	return entity, nil
}

func (s *flowService) ActivateFlowSequenceStep(ctx context.Context, tenant, id string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowService.ActivateFlowSequenceStep")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	span.SetTag(tracing.SpanTagComponent, constants.ComponentService)

	entity, err := s.services.PostgresRepositories.FlowSequenceStepRepository.GetById(ctx, tenant, id)
	if err != nil {
		return err
	}

	if entity == nil {
		return nil
	}

	entity.Active = true

	_, err = s.services.PostgresRepositories.FlowSequenceStepRepository.Store(ctx, tenant, entity)
	if err != nil {
		return err
	}

	return nil
}

func (s *flowService) DeactivateFlowSequenceStep(ctx context.Context, tenant, id string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowService.DeactivateFlowSequenceStep")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	span.SetTag(tracing.SpanTagComponent, constants.ComponentService)

	entity, err := s.services.PostgresRepositories.FlowSequenceStepRepository.GetById(ctx, tenant, id)
	if err != nil {
		return err
	}

	if entity == nil {
		return nil
	}

	entity.Active = false

	_, err = s.services.PostgresRepositories.FlowSequenceStepRepository.Store(ctx, tenant, entity)
	if err != nil {
		return err
	}

	return nil
}

func (s *flowService) DeleteFlowSequenceStep(ctx context.Context, tenant, id string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowService.DeleteFlowSequenceStep")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	span.SetTag(tracing.SpanTagComponent, constants.ComponentService)
	span.LogFields(log.String("id", id))

	err := s.services.PostgresRepositories.FlowSequenceStepRepository.Delete(ctx, tenant, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *flowService) GetFlowSequenceContacts(ctx context.Context, tenant, sequenceId string) ([]*postgresEntity.FlowSequenceContact, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowService.GetFlowSequenceContacts")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	span.SetTag(tracing.SpanTagComponent, constants.ComponentService)

	entities, err := s.services.PostgresRepositories.FlowSequenceContactRepository.Get(ctx, tenant, sequenceId)
	if err != nil {
		return nil, err
	}

	return entities, nil
}

func (s *flowService) GetFlowSequenceContactById(ctx context.Context, tenant, id string) (*postgresEntity.FlowSequenceContact, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowService.GetFlowSequenceContactById")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	span.SetTag(tracing.SpanTagComponent, constants.ComponentService)

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
	span.SetTag(tracing.SpanTagTenant, tenant)
	span.SetTag(tracing.SpanTagComponent, constants.ComponentService)

	entity, err := s.services.PostgresRepositories.FlowSequenceContactRepository.Store(ctx, tenant, entity)
	if err != nil {
		return nil, err
	}

	return entity, nil
}

func (s *flowService) DeleteFlowSequenceContact(ctx context.Context, tenant, id string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowService.DeleteFlowSequenceContact")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	span.SetTag(tracing.SpanTagComponent, constants.ComponentService)
	span.LogFields(log.String("id", id))

	err := s.services.PostgresRepositories.FlowSequenceContactRepository.Delete(ctx, tenant, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *flowService) GetFlowSequenceSenders(ctx context.Context, tenant, sequenceId string) ([]*postgresEntity.FlowSequenceSender, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowService.GetFlowSequenceSenders")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	span.SetTag(tracing.SpanTagComponent, constants.ComponentService)

	entities, err := s.services.PostgresRepositories.FlowSequenceSenderRepository.Get(ctx, tenant, sequenceId)
	if err != nil {
		return nil, err
	}

	return entities, nil
}

func (s *flowService) GetFlowSequenceSenderById(ctx context.Context, tenant, id string) (*postgresEntity.FlowSequenceSender, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowService.GetFlowSequenceSenderById")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	span.SetTag(tracing.SpanTagComponent, constants.ComponentService)

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
	span.SetTag(tracing.SpanTagTenant, tenant)
	span.SetTag(tracing.SpanTagComponent, constants.ComponentService)

	entity, err := s.services.PostgresRepositories.FlowSequenceSenderRepository.Store(ctx, tenant, entity)
	if err != nil {
		return nil, err
	}

	return entity, nil
}

func (s *flowService) DeleteFlowSequenceSender(ctx context.Context, tenant, id string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowService.DeleteFlowSequenceSender")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	span.SetTag(tracing.SpanTagComponent, constants.ComponentService)
	span.LogFields(log.String("id", id))

	err := s.services.PostgresRepositories.FlowSequenceSenderRepository.Delete(ctx, tenant, id)
	if err != nil {
		return err
	}

	return nil
}
