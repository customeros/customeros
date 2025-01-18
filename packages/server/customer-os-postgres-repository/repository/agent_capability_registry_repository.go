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

type AgentCapabilityRegistryRepository interface {
	Initialize(ctx context.Context) error
	Create(ctx context.Context, agent entity.AgentCapabilityRegistry) (*entity.AgentCapabilityRegistry, error)
	Find(ctx context.Context, agentID enum.AgentCapabilityType) (*entity.AgentCapabilityRegistry, error)
	FindAll(ctx context.Context) ([]entity.AgentCapabilityRegistry, error)
}

type agentCapabilityRegistryRepository struct {
	gormDb *gorm.DB
}

func NewAgentCapabilityRegistryRepository(gormDb *gorm.DB) AgentCapabilityRegistryRepository {
	return &agentCapabilityRegistryRepository{gormDb: gormDb}
}

func (r *agentCapabilityRegistryRepository) FindAll(ctx context.Context) ([]entity.AgentCapabilityRegistry, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentCapabilityRegistryRepository.FindAll")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	var capabilities []entity.AgentCapabilityRegistry
	err := r.gormDb.WithContext(ctx).
		Where("is_active = true").
		Find(&capabilities).Error
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return capabilities, nil
}

func (r *agentCapabilityRegistryRepository) Find(ctx context.Context, capability enum.AgentCapabilityType) (*entity.AgentCapabilityRegistry, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentCapabilityRegistryRepository.Find")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	var capabilities entity.AgentCapabilityRegistry
	query := r.gormDb.WithContext(ctx).
		Where("is_active = true").
		Where("type = ?", capability.String())

	err := query.First(&capabilities).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		tracing.TraceErr(span, err)
		return nil, err
	}

	return &capabilities, nil
}

func (r *agentCapabilityRegistryRepository) Create(ctx context.Context, capability entity.AgentCapabilityRegistry) (*entity.AgentCapabilityRegistry, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentCapabilityRegistryRepository.Create")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	err := r.gormDb.WithContext(ctx).Create(&capability).Error
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return &capability, nil
}

func (a *agentCapabilityRegistryRepository) Update(ctx context.Context, capability entity.AgentCapabilityRegistry) (*entity.AgentCapabilityRegistry, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentCapabilityRegistryRepository.Update")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	// Validate input
	if capability.ID == "" && capability.Type == "" {
		err := errors.New("agent identifier is missing")
		tracing.TraceErr(span, err)
		return nil, err
	}

	var updatedCapability entity.AgentCapabilityRegistry

	err := a.gormDb.
		Model(&entity.AgentCapabilityRegistry{}).
		Where("id = ?", capability.ID).
		Updates(&capability).
		First(&updatedCapability, "id = ?", capability.ID).
		Error

	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return &updatedCapability, nil
}

func (r *agentCapabilityRegistryRepository) Initialize(ctx context.Context) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentRegistryRepository.Initialize")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	requiredCapabilities := []entity.AgentCapabilityRegistry{
		buildCapabilities(
			enum.CapabilityTrackWebSession,
			"Track and identify website visitors",
		),
		buildCapabilities(
			enum.CapabilityIdentifyWebVisitor,
			"Identify website visitors",
		),
		buildCapabilities(
			enum.CapabilityCreateOrganization,
			"Create new organizations",
		),
		buildCapabilities(
			enum.CapabilityAnalyzeWebSessionIntent,
			"Analyze web session for intent signals",
		),
		buildCapabilities(
			enum.CapabilitySendSlackNotification,
			"Send Slack notification",
		),
	}

	for _, capability := range requiredCapabilities {
		validCapability, err := enum.GetAgentCapability(capability.Type)
		if err != nil {
			tracing.TraceErr(span, err)
			continue
		}
		existingCapability, err := r.Find(ctx, validCapability)
		if err != nil {
			tracing.TraceErr(span, err)
			continue
		}

		if existingCapability != nil {
			continue
		}

		_, createErr := r.Create(ctx, capability)
		if createErr != nil {
			tracing.TraceErr(span, createErr)
			return createErr
		}
	}
	return nil
}

func buildCapabilities(capType enum.AgentCapabilityType, desc string) entity.AgentCapabilityRegistry {
	return entity.AgentCapabilityRegistry{
		Type:        capType.String(),
		Description: desc,
		IsActive:    true,
	}
}
