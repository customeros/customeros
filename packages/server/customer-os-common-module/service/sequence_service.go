package service

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	postgresEntity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type SequenceService interface {
	GetSequences(ctx context.Context, tenant string) ([]*postgresEntity.Sequence, error)
	GetSequenceById(ctx context.Context, tenant, id string) (*postgresEntity.Sequence, error)

	StoreSequence(ctx context.Context, entity *postgresEntity.Sequence) (*postgresEntity.Sequence, error)
	EnableSequence(ctx context.Context, tenant, id string) error
	DisableSequence(ctx context.Context, tenant, id string) error

	GetSequenceSteps(ctx context.Context, tenant, sequenceId string) ([]*postgresEntity.SequenceStep, error)
	GetSequenceStepById(ctx context.Context, tenant, id string) (*postgresEntity.SequenceStep, error)
	StoreSequenceStep(ctx context.Context, tenant string, entity *postgresEntity.SequenceStep) (*postgresEntity.SequenceStep, error)
	DeleteSequenceStep(ctx context.Context, tenant, id string) error

	GetSequenceContacts(ctx context.Context, tenant, sequenceId string) ([]*postgresEntity.SequenceContact, error)
	GetSequenceContactById(ctx context.Context, tenant, id string) (*postgresEntity.SequenceContact, error)
	StoreSequenceContact(ctx context.Context, tenant string, entity *postgresEntity.SequenceContact) (*postgresEntity.SequenceContact, error)
	DeleteSequenceContact(ctx context.Context, tenant, id string) error
}

type sequenceService struct {
	services *Services
}

func NewSequenceService(services *Services) SequenceService {
	return &sequenceService{
		services: services,
	}
}

func (s *sequenceService) GetSequences(ctx context.Context, tenant string) ([]*postgresEntity.Sequence, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SequenceService.GetSequences")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	span.SetTag(tracing.SpanTagComponent, constants.ComponentService)

	entities, err := s.services.PostgresRepositories.SequenceRepository.Get(ctx, tenant)
	if err != nil {
		return nil, err
	}

	return entities, nil
}

func (s *sequenceService) GetSequenceById(ctx context.Context, tenant, id string) (*postgresEntity.Sequence, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SequenceService.GetSequenceById")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	span.SetTag(tracing.SpanTagComponent, constants.ComponentService)

	span.LogFields(log.String("id", id))

	entity, err := s.services.PostgresRepositories.SequenceRepository.GetById(ctx, tenant, id)
	if err != nil {
		return nil, err
	}

	return entity, nil
}

func (s *sequenceService) StoreSequence(ctx context.Context, entity *postgresEntity.Sequence) (*postgresEntity.Sequence, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SequenceService.GetSequenceById")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, entity.Tenant)
	span.SetTag(tracing.SpanTagComponent, constants.ComponentService)

	entity, err := s.services.PostgresRepositories.SequenceRepository.Store(ctx, entity)
	if err != nil {
		return nil, err
	}

	return entity, nil
}

func (s *sequenceService) EnableSequence(ctx context.Context, tenant, id string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SequenceService.EnableSequence")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	span.SetTag(tracing.SpanTagComponent, constants.ComponentService)

	entity, err := s.services.PostgresRepositories.SequenceRepository.GetById(ctx, tenant, id)
	if err != nil {
		return err
	}

	if entity == nil {
		return nil
	}

	entity.Enabled = true

	_, err = s.services.PostgresRepositories.SequenceRepository.Store(ctx, entity)
	if err != nil {
		return err
	}

	return nil
}

func (s *sequenceService) DisableSequence(ctx context.Context, tenant, id string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SequenceService.DisableSequence")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	span.SetTag(tracing.SpanTagComponent, constants.ComponentService)

	entity, err := s.services.PostgresRepositories.SequenceRepository.GetById(ctx, tenant, id)
	if err != nil {
		return err
	}

	if entity == nil {
		return nil
	}

	entity.Enabled = false

	_, err = s.services.PostgresRepositories.SequenceRepository.Store(ctx, entity)
	if err != nil {
		return err
	}

	return nil
}

func (s *sequenceService) GetSequenceSteps(ctx context.Context, tenant, sequenceId string) ([]*postgresEntity.SequenceStep, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SequenceService.GetSequenceSteps")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	span.SetTag(tracing.SpanTagComponent, constants.ComponentService)

	entities, err := s.services.PostgresRepositories.SequenceStepRepository.Get(ctx, tenant, sequenceId)
	if err != nil {
		return nil, err
	}

	return entities, nil
}

func (s *sequenceService) GetSequenceStepById(ctx context.Context, tenant, id string) (*postgresEntity.SequenceStep, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SequenceService.GetSequenceStepById")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	span.SetTag(tracing.SpanTagComponent, constants.ComponentService)

	span.LogFields(log.String("id", id))

	entity, err := s.services.PostgresRepositories.SequenceStepRepository.GetById(ctx, tenant, id)
	if err != nil {
		return nil, err
	}

	return entity, nil
}

func (s *sequenceService) StoreSequenceStep(ctx context.Context, tenant string, entity *postgresEntity.SequenceStep) (*postgresEntity.SequenceStep, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SequenceService.StoreSequenceStep")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	span.SetTag(tracing.SpanTagComponent, constants.ComponentService)

	entity, err := s.services.PostgresRepositories.SequenceStepRepository.Store(ctx, tenant, entity)
	if err != nil {
		return nil, err
	}

	return entity, nil
}

func (s *sequenceService) DeleteSequenceStep(ctx context.Context, tenant, id string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SequenceService.DeleteSequenceStep")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	span.SetTag(tracing.SpanTagComponent, constants.ComponentService)
	span.LogFields(log.String("id", id))

	err := s.services.PostgresRepositories.SequenceStepRepository.Delete(ctx, tenant, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *sequenceService) GetSequenceContacts(ctx context.Context, tenant, sequenceId string) ([]*postgresEntity.SequenceContact, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SequenceService.GetSequenceContacts")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	span.SetTag(tracing.SpanTagComponent, constants.ComponentService)

	entities, err := s.services.PostgresRepositories.SequenceContactRepository.Get(ctx, tenant, sequenceId)
	if err != nil {
		return nil, err
	}

	return entities, nil
}

func (s *sequenceService) GetSequenceContactById(ctx context.Context, tenant, id string) (*postgresEntity.SequenceContact, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SequenceService.GetSequenceContactById")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	span.SetTag(tracing.SpanTagComponent, constants.ComponentService)

	span.LogFields(log.String("id", id))

	entity, err := s.services.PostgresRepositories.SequenceContactRepository.GetById(ctx, tenant, id)
	if err != nil {
		return nil, err
	}

	return entity, nil
}

func (s *sequenceService) StoreSequenceContact(ctx context.Context, tenant string, entity *postgresEntity.SequenceContact) (*postgresEntity.SequenceContact, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SequenceService.StoreSequenceContact")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	span.SetTag(tracing.SpanTagComponent, constants.ComponentService)

	entity, err := s.services.PostgresRepositories.SequenceContactRepository.Store(ctx, tenant, entity)
	if err != nil {
		return nil, err
	}

	return entity, nil
}

func (s *sequenceService) DeleteSequenceContact(ctx context.Context, tenant, id string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SequenceService.DeleteSequenceContact")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	span.SetTag(tracing.SpanTagComponent, constants.ComponentService)
	span.LogFields(log.String("id", id))

	err := s.services.PostgresRepositories.SequenceContactRepository.Delete(ctx, tenant, id)
	if err != nil {
		return err
	}

	return nil
}
