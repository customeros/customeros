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
	GetUniqueFromNodes(ctx context.Context) ([]FromNodeRecord, error)
	GetFlowTransitionsByNode(ctx context.Context, fromNode enum.FlowNodeType) ([]entity.FlowTransitionsRegistry, error)
	GetAllFlowTransitions(ctx context.Context) ([]entity.FlowTransitionsRegistry, error)
	CreateFlowTransition(ctx context.Context, transition *entity.FlowTransitionsRegistry) error
}

type flowTransitionsRegistryRepository struct {
	gormDb *gorm.DB
}

type FromNodeRecord struct {
	FromNode     string `gorm:"column:from_node"`
	FromNodeType string `gorm:"column:from_node_type"`
}

func NewFlowTransitionsRegistryRepository(gormDb *gorm.DB) FlowTransitionsRegistryRepository {
	return &flowTransitionsRegistryRepository{gormDb: gormDb}
}

func (r *flowTransitionsRegistryRepository) GetUniqueFromNodes(ctx context.Context) ([]FromNodeRecord, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowTransitionRegistryRepository.GetUniqueFromNodes")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	var uniqueFromNodes []FromNodeRecord
	err := r.gormDb.WithContext(ctx).
		Distinct("from_node").
		Where("enabled = true").
		Order("from_node DESC").
		Find(&uniqueFromNodes).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
	}

	return uniqueFromNodes, err
}

func (r *flowTransitionsRegistryRepository) GetFlowTransitionsByNode(ctx context.Context, fromNode enum.FlowNodeType) ([]entity.FlowTransitionsRegistry, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowTransitionRegistryRepository.GetFlowTransitionsByNode")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	var transition []entity.FlowTransitionsRegistry
	err := r.gormDb.WithContext(ctx).
		Where("from_node = ? AND enabled = true", fromNode).
		Order("to_node_type DESC").
		Find(&transition).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
	}

	return transition, err
}

func (r *flowTransitionsRegistryRepository) GetAllFlowTransitions(ctx context.Context) ([]entity.FlowTransitionsRegistry, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowTransactionRegistryRepository.GetAllFlowTransactions")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	var transitions []entity.FlowTransitionsRegistry
	err := r.gormDb.WithContext(ctx).
		Where("enabled = true").
		Order("from_node DESC").
		Find(&transitions).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
	}

	if err != nil {
		tracing.TraceErr(span, err)
	}

	return transitions, err
}

func (r *flowTransitionsRegistryRepository) CreateFlowTransition(ctx context.Context, transition *entity.FlowTransitionsRegistry) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowTransactionRegistryRepository.CreateFlowTransition")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	return r.gormDb.WithContext(ctx).Create(transition).Error
}

func (r *flowTransitionsRegistryRepository) InitializeFlowTransitions(ctx context.Context) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowTransitionsRegistry.InitializeFlowTransitions")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	required := []entity.FlowTransitionsRegistry{
		{
			FromNodeType: enum.NodeFlowListenerEvent.String(),
			FromNode:     enum.EventFathomMeetingSummaryCreated.String(),
			ToNodeType:   enum.NodeFlowAction.String(),
			ToNode:       enum.ActionTimelineEventCreate.String(),
			Enabled:      true,
		},
		// ... add more here
	}

	dbTransitions, err := r.GetAllFlowTransitions(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	for _, transition := range required {
		exists := false

		for _, t := range dbTransitions {
			if t.FromNode == transition.FromNode && t.ToNode == transition.ToNode {
				exists = true
			}
		}

		if !exists {
			createErr := r.CreateFlowTransition(ctx, &transition)
			if createErr != nil {
				tracing.TraceErr(span, err)
				return err
			}
		}
	}

	return nil

}
