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

type FlowAgentRegistryRepository interface {
	Initialize(ctx context.Context) error
	Create(ctx context.Context, flowAgent entity.FlowAgentRegistry) (*entity.FlowAgentRegistry, error)
	Find(ctx context.Context, flowAgent entity.FlowAgentRegistry) (*entity.FlowAgentRegistry, error)
	FindAll(ctx context.Context) (*[]entity.FlowAgentRegistry, error)
}

type flowAgentRegistryRepository struct {
	gormDb *gorm.DB
}

func NewFlowAgentRegistryRepository(gormDb *gorm.DB) FlowAgentRegistryRepository {
	return &flowAgentRegistryRepository{gormDb: gormDb}
}

func (r *flowAgentRegistryRepository) FindAll(ctx context.Context) (*[]entity.FlowAgentRegistry, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowAgentRegistryRepository.FindAll")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	var actions []entity.FlowAgentRegistry
	err := r.gormDb.WithContext(ctx).
		Where("status = ?", enum.FlowNodeEdgeStatusActive.String()).
		Order("action DESC").
		Find(&actions).Error
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return &actions, nil
}

func (r *flowAgentRegistryRepository) Find(ctx context.Context, flowAgent entity.FlowAgentRegistry) (*entity.FlowAgentRegistry, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowAgentRegistryRepository.Find")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	var action entity.FlowAgentRegistry
	query := r.gormDb.WithContext(ctx).
		Where("status = ?", enum.FlowNodeEdgeStatusActive.String())

	// Add additional filters based on non-zero fields in flowAgent
	if flowAgent != (entity.FlowAgentRegistry{}) {
		query = query.Where(&flowAgent)
	}

	err := query.First(&action).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		tracing.TraceErr(span, err)
		return nil, err
	}

	return &action, nil
}

func (r *flowAgentRegistryRepository) Create(ctx context.Context, flowAgent entity.FlowAgentRegistry) (*entity.FlowAgentRegistry, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowAgentRegistryRepository.Create")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	err := r.gormDb.WithContext(ctx).Create(&flowAgent).Error
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return &flowAgent, nil
}

func (r *flowAgentRegistryRepository) Initialize(ctx context.Context) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowAgentRegistryRepository.Initialize")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	requiredAgents := []entity.FlowAgentRegistry{
		{
			Agent:        enum.AgentContactCreate.String(),
			FriendlyName: "Create a Contact",
			Description:  "Creates a new Contact",
			Status:       enum.FlowNodeEdgeStatusActive.String(),
		},
		{
			Agent:        enum.AgentEmailSendNew.String(),
			FriendlyName: "Send an email",
			Description:  "Begins a new email thread",
			Status:       enum.FlowNodeEdgeStatusActive.String(),
		},
		{
			Agent:        enum.AgentEmailSendReply.String(),
			FriendlyName: "Reply to an email thread",
			Description:  "Sends a reply to an existing email thread",
			Status:       enum.FlowNodeEdgeStatusActive.String(),
		},
		{
			Agent:        enum.AgentLinkedinConnect.String(),
			FriendlyName: "Send Linkedin Connection Request",
			Description:  "Sends a connection request on LinkedIn",
			Status:       enum.FlowNodeEdgeStatusActive.String(),
		},
		{
			Agent:        enum.AgentLinkedinMessage.String(),
			FriendlyName: "Send Linkedin Message",
			Description:  "Sends a direct message to a LinkedIn connection",
			Status:       enum.FlowNodeEdgeStatusActive.String(),
		},
		{
			Agent:        enum.AgentOrganizationCreate.String(),
			FriendlyName: "Create an Organization",
			Description:  "Creates a new Organization",
			Status:       enum.FlowNodeEdgeStatusActive.String(),
		},
		{
			Agent:        enum.AgentSlackNotify.String(),
			FriendlyName: "Notify a Slack Channel",
			Description:  "Sends a notification to the specified slack channel",
			Status:       enum.FlowNodeEdgeStatusActive.String(),
		},
		{
			Agent:        enum.AgentTimelineEventCreate.String(),
			FriendlyName: "Add Event to Timeline",
			Description:  "Adds a new event to the Organization Timeline",
			Status:       enum.FlowNodeEdgeStatusActive.String(),
		},
		// ... add more here
	}

	for _, action := range requiredAgents {
		existingAgent, err := r.Find(ctx, action)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}

		if existingAgent != nil && existingAgent.Agent != "" {
			continue
		}

		_, createErr := r.Create(ctx, action)
		if err != nil {
			tracing.TraceErr(span, createErr)
			return createErr
		}
	}

	return nil
}
