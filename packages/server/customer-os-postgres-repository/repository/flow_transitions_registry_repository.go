package repository

import (
	"context"
	"errors"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/opentracing/opentracing-go"
	"gorm.io/gorm"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
)

type FlowTransitionsRegistryRepository interface {
	InitializeFlowTransitions(ctx context.Context) error
	GetFlowTransitionsByNode(ctx context.Context, fromNode string) ([]entity.FlowTransitionsRegistry, error)
	GetAllFlowTransitions(ctx context.Context) ([]entity.FlowTransitionsRegistry, error)
	CreateFlowTransition(ctx context.Context, transition *entity.FlowTransitionsRegistry) error
}

type flowTransitionsRegistryRepository struct {
	gormDb *gorm.DB
}

func NewFlowTransitionsRegistryRepository(gormDb *gorm.DB) FlowTransitionsRegistryRepository {
	return &flowTransitionsRegistryRepository{gormDb: gormDb}
}

func (r *flowTransitionsRegistryRepository) FindFlowAction(ctx context.Context, actionName string) (entity.FlowActionRegistry, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowActionRegistryRepository.FindFlowAction")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	var action entity.FlowActionRegistry
	err := r.gormDb.WithContext(ctx).
		Where("action = ? AND enabled = true", actionName).
		First(&action).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
	}

	return action, err
}

func (r *flowTransitionsRegistryRepository) CreateFlowAction(ctx context.Context, action *entity.FlowActionRegistry) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowActionRegistryRepository.CreateFlowAction")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	return r.gormDb.WithContext(ctx).Create(action).Error
}

func (r *flowTransitionsRegistryRepository) InitializeFlowTransitions(ctx context.Context) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowActionRegistryRepository.InitializeActions")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	requiredActions := []entity.FlowActionRegistry{
		{
			Action:       enum.ActionContactCreate.String(),
			FriendlyName: "Create a Contact",
			Description:  "Creates a new Contact",
			Enabled:      true,
		},
		{
			Action:       enum.ActionOrganizationCreate.String(),
			FriendlyName: "Create an Organization",
			Description:  "Creates a new Organization",
			Enabled:      true,
		},
		{
			Action:       enum.ActionTimelineEventCreate.String(),
			FriendlyName: "Add Event to Timeline",
			Description:  "Adds a new event to the Organization Timeline",
			Enabled:      true,
		},
		// ... add more here
	}

	for _, action := range requiredActions {
		existingAction, err := r.FindFlowAction(ctx, action.Action)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}

		if existingAction.Action != "" {
			continue
		}

		createErr := r.CreateFlowAction(ctx, &action)
		if err != nil {
			tracing.TraceErr(span, createErr)
			return createErr
		}
	}

	return nil

}
