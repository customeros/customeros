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

type AgentRegistryRepository interface {
	Initialize(ctx context.Context) error
	Create(ctx context.Context, flowAgent entity.AgentRegistry) (*entity.AgentRegistry, error)
	Find(ctx context.Context, flowAgent entity.AgentRegistry) (*entity.AgentRegistry, error)
	FindAll(ctx context.Context) ([]entity.AgentRegistry, error)
}

type flowAgentRegistryRepository struct {
	gormDb *gorm.DB
}

func NewAgentRegistryRepository(gormDb *gorm.DB) AgentRegistryRepository {
	return &flowAgentRegistryRepository{gormDb: gormDb}
}

func (r *flowAgentRegistryRepository) FindAll(ctx context.Context) ([]entity.AgentRegistry, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentRegistryRepository.FindAll")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	var actions []entity.AgentRegistry
	err := r.gormDb.WithContext(ctx).
		Where("status = ?", enum.NodeEdgeStatusActive.String()).
		Order("action DESC").
		Find(&actions).Error
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return actions, nil
}

func (r *flowAgentRegistryRepository) Find(ctx context.Context, flowAgent entity.AgentRegistry) (*entity.AgentRegistry, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentRegistryRepository.Find")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	var action entity.AgentRegistry
	query := r.gormDb.WithContext(ctx).
		Where("status = ?", enum.NodeEdgeStatusActive.String())

	// Add additional filters based on non-zero fields in flowAgent
	if flowAgent != (entity.AgentRegistry{}) {
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

func (r *flowAgentRegistryRepository) Create(ctx context.Context, flowAgent entity.AgentRegistry) (*entity.AgentRegistry, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentRegistryRepository.Create")
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
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentRegistryRepository.Initialize")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	requiredAgents := []entity.AgentRegistry{
		{
			Agent:        enum.AgentContactCreate.String(),
			FriendlyName: "Create a Contact",
			Description:  "Creates a new Contact",
			Status:       enum.NodeEdgeStatusActive.String(),
		},
		{
			Agent:        enum.AgentEmailSendNew.String(),
			FriendlyName: "Send an email",
			Description:  "Begins a new email thread",
			Status:       enum.NodeEdgeStatusActive.String(),
		},
		{
			Agent:        enum.AgentEmailSendReply.String(),
			FriendlyName: "Reply to an email thread",
			Description:  "Sends a reply to an existing email thread",
			Status:       enum.NodeEdgeStatusActive.String(),
		},
		{
			Agent:        enum.AgentLinkedinConnect.String(),
			FriendlyName: "Send Linkedin Connection Request",
			Description:  "Sends a connection request on LinkedIn",
			Status:       enum.NodeEdgeStatusActive.String(),
		},
		{
			Agent:        enum.AgentLinkedinMessage.String(),
			FriendlyName: "Send Linkedin Message",
			Description:  "Sends a direct message to a LinkedIn connection",
			Status:       enum.NodeEdgeStatusActive.String(),
		},
		{
			Agent:        enum.AgentOrganizationCreate.String(),
			FriendlyName: "Create an Organization",
			Description:  "Creates a new Organization",
			Status:       enum.NodeEdgeStatusActive.String(),
		},
		{
			Agent:        enum.AgentSlackNotify.String(),
			FriendlyName: "Notify a Slack Channel",
			Description:  "Sends a notification to the specified slack channel",
			Status:       enum.NodeEdgeStatusActive.String(),
		},
		{
			Agent:        enum.AgentTimelineEventCreate.String(),
			FriendlyName: "Add Event to Timeline",
			Description:  "Adds a new event to the Organization Timeline",
			Status:       enum.NodeEdgeStatusActive.String(),
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
