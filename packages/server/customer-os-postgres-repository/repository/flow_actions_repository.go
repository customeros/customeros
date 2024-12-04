package repository

import (
	"context"
	"errors"
	"log"

	commonenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/opentracing/opentracing-go"
	"gorm.io/gorm"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
)

type FlowActionRepository interface {
	InitializeActions(ctx context.Context) error
	FindFlowActionByName(ctx context.Context, actionName string) (entity.FlowAction, error)
	CreateFlowAction(ctx context.Context, action *entity.FlowAction) error
	GetFlowActionsByEvent(ctx context.Context, event commonenum.FlowEvent) ([]entity.FlowAction, error)
}

type flowActionRepository struct {
	gormDb *gorm.DB
}

func NewFlowActionRepository(gormDb *gorm.DB) FlowActionRepository {
	repo := &flowActionRepository{gormDb: gormDb}
	if err := repo.InitializeActions(context.Background()); err != nil {
		log.Printf("Failed to initialize flow actions: %v", err)
	}
	return repo
}

func (r *flowActionRepository) InitializeActions(ctx context.Context) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowActionRepository.InitializeActions")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	requiredActions := []entity.FlowAction{
		{
			EventName:   commonenum.EventFathomMeetingSummaryCreated.String(),
			ActionType:  commonenum.ActionCreateTimelineEvent.String(),
			Description: "Add meeting summary to the timeline of each participant",
			Enabled:     true,
		},
		{
			EventName:   commonenum.EventFathomMeetingSummaryCreated.String(),
			ActionType:  commonenum.ActionCreateContacts.String(),
			Description: "Create contacts for each participant (if they don't exist)",
			Enabled:     true,
		},
		{
			EventName:   commonenum.EventGrainMeetingSummaryCreated.String(),
			ActionType:  commonenum.ActionCreateTimelineEvent.String(),
			Description: "Add meeting summary to the timeline of each participant",
			Enabled:     true,
		},
		{
			EventName:   commonenum.EventGrainMeetingSummaryCreated.String(),
			ActionType:  commonenum.ActionCreateContacts.String(),
			Description: "Create contacts for each participant (if they don't exist)",
			Enabled:     true,
		},
	}

	for _, action := range requiredActions {
		existingAction, err := r.FindFlowActionByName(ctx, action.ActionType)
		if err != nil {
			return err
		}

		if existingAction.ActionType != "" {
			continue
		}

		if err := r.CreateFlowAction(ctx, &action); err != nil {
			return err
		}
	}

	return nil
}

func (r *flowActionRepository) FindFlowActionByName(ctx context.Context, actionName string) (entity.FlowAction, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowActionRepository.FindFlowActionByName")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	var action entity.FlowAction
	err := r.gormDb.WithContext(ctx).
		Where("action_type = ? AND enabled = true", actionName).
		First(&action).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
	}

	return action, err
}

func (r *flowActionRepository) CreateFlowAction(ctx context.Context, action *entity.FlowAction) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowActionRepository.CreateFlowAction")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	return r.gormDb.WithContext(ctx).Create(action).Error
}

func (r *flowActionRepository) GetFlowActionsByEvent(ctx context.Context, event commonenum.FlowEvent) ([]entity.FlowAction, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowActionRepository.GetFlowActionsByEvent")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	var actions []entity.FlowAction
	err := r.gormDb.WithContext(ctx).
		Where("enabled = true AND event_name = ?", event.String()).
		Order("action_type DESC").
		Find(&actions).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
	}
	return actions, err
}
