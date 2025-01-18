package postgres_repository

import (
	"context"
	"errors"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/opentracing/opentracing-go"
	"gorm.io/gorm"

	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
)

type AgentCapabilityRegistryRepository interface {
	Initialize(ctx context.Context) error
	Create(ctx context.Context, agent postgres_entity.AgentCapabilityRegistry) (*postgres_entity.AgentCapabilityRegistry, error)
	Find(ctx context.Context, agentID enum.AgentCapabilityType) (*postgres_entity.AgentCapabilityRegistry, error)
	FindAll(ctx context.Context) ([]postgres_entity.AgentCapabilityRegistry, error)
}

type agentCapabilityRegistryRepository struct {
	gormDb *gorm.DB
}

func NewAgentCapabilityRegistryRepository(gormDb *gorm.DB) AgentCapabilityRegistryRepository {
	return &agentCapabilityRegistryRepository{gormDb: gormDb}
}

func (r *agentCapabilityRegistryRepository) FindAll(ctx context.Context) ([]postgres_entity.AgentCapabilityRegistry, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentCapabilityRegistryRepository.FindAll")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	var capabilities []postgres_entity.AgentCapabilityRegistry
	err := r.gormDb.WithContext(ctx).
		Where("is_active = true").
		Find(&capabilities).Error
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return capabilities, nil
}

func (r *agentCapabilityRegistryRepository) Find(ctx context.Context, capability enum.AgentCapabilityType) (*postgres_entity.AgentCapabilityRegistry, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentCapabilityRegistryRepository.Find")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	var capabilities postgres_entity.AgentCapabilityRegistry
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

func (r *agentCapabilityRegistryRepository) Create(ctx context.Context, capability postgres_entity.AgentCapabilityRegistry) (*postgres_entity.AgentCapabilityRegistry, error) {
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

func (a *agentCapabilityRegistryRepository) Update(ctx context.Context, capability postgres_entity.AgentCapabilityRegistry) (*postgres_entity.AgentCapabilityRegistry, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentCapabilityRegistryRepository.Update")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	// Validate input
	if capability.ID == "" && capability.Type == "" {
		err := errors.New("agent identifier is missing")
		tracing.TraceErr(span, err)
		return nil, err
	}

	var updatedCapability postgres_entity.AgentCapabilityRegistry

	err := a.gormDb.
		Model(&postgres_entity.AgentCapabilityRegistry{}).
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

	requiredCapabilities := []postgres_entity.AgentCapabilityRegistry{
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

func buildCapabilities(capType enum.AgentCapabilityType, desc string) postgres_entity.AgentCapabilityRegistry {
	return postgres_entity.AgentCapabilityRegistry{
		Type:        capType.String(),
		Description: desc,
		IsActive:    true,
	}
}
