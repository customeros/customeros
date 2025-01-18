package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/opentracing/opentracing-go"
	"gorm.io/gorm"

	"github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
)

var (
	ErrAgentIDMissing    = errors.New("agent ID is missing")
	ErrAgentNotFound     = errors.New("agent not found")
	ErrAgentCreateFailed = errors.New("failed to create agent")
)

type AgentRegistryRepository interface {
	Initialize(ctx context.Context) error
	Create(ctx context.Context, agent entity.AgentRegistry) (*entity.AgentRegistry, error)
	Update(ctx context.Context, agent entity.AgentRegistry) (*entity.AgentRegistry, error)
	Find(ctx context.Context, agentType enum.AgentType) (*entity.AgentRegistry, error)
	FindAll(ctx context.Context) ([]entity.AgentRegistry, error)
}

type agentRegistryRepository struct {
	gormDb *gorm.DB
}

func NewAgentRegistryRepository(gormDb *gorm.DB) AgentRegistryRepository {
	if gormDb == nil {
		panic("gormDb cannot be nil")
	}
	return &agentRegistryRepository{gormDb: gormDb}
}

func (r *agentRegistryRepository) FindAll(ctx context.Context) ([]entity.AgentRegistry, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentRegistryRepository.FindAll")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	var agents []entity.AgentRegistry
	err := r.gormDb.WithContext(ctx).
		Where("is_active = ?", true).
		Find(&agents).Error
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, fmt.Errorf("failed to find agents: %w", err)
	}

	return agents, nil
}

func (r *agentRegistryRepository) Find(ctx context.Context, agentType enum.AgentType) (*entity.AgentRegistry, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentRegistryRepository.Find")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	var agent entity.AgentRegistry
	err := r.gormDb.WithContext(ctx).
		Where("is_active = ? AND type = ?", true, agentType.String()).
		First(&agent).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		tracing.TraceErr(span, err)
		return nil, fmt.Errorf("failed to find agent: %w", err)
	}

	return &agent, nil
}

func (r *agentRegistryRepository) Create(ctx context.Context, agent entity.AgentRegistry) (*entity.AgentRegistry, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentRegistryRepository.Create")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	err := r.gormDb.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Check for existing agent with same type
		var count int64
		if err := tx.Model(&entity.AgentRegistry{}).
			Where("type = ? AND is_active = ?", agent.Type, true).
			Count(&count).Error; err != nil {
			return err
		}

		if count > 0 {
			return fmt.Errorf("agent with type %s already exists", agent.Type)
		}

		return tx.Create(&agent).Error
	})

	if err != nil {
		tracing.TraceErr(span, err)
		return nil, fmt.Errorf("%w: %v", ErrAgentCreateFailed, err)
	}

	return &agent, nil
}

func (r *agentRegistryRepository) Update(ctx context.Context, agent entity.AgentRegistry) (*entity.AgentRegistry, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentRegistryRepository.Update")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	if agent.ID == "" {
		return nil, ErrAgentIDMissing
	}

	var updatedAgent entity.AgentRegistry
	err := r.gormDb.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Check if agent exists
		if err := tx.First(&entity.AgentRegistry{}, "id = ?", agent.ID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrAgentNotFound
			}
			return err
		}

		// Perform update
		result := tx.Model(&entity.AgentRegistry{}).
			Where("id = ?", agent.ID).
			Updates(&agent)

		if result.Error != nil {
			return result.Error
		}

		if result.RowsAffected == 0 {
			return fmt.Errorf("no rows affected during update")
		}

		// Fetch updated record
		return tx.First(&updatedAgent, "id = ?", agent.ID).Error
	})

	if err != nil {
		tracing.TraceErr(span, err)
		return nil, fmt.Errorf("failed to update agent: %w", err)
	}

	return &updatedAgent, nil
}

func (r *agentRegistryRepository) Initialize(ctx context.Context) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentRegistryRepository.Initialize")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	requiredAgents := []entity.AgentRegistry{
		registerVisitorIdentityAgent(),
		// Add more agents here
	}

	for _, agent := range requiredAgents {
		// Look for existing agent by type (not ID since it might not be set yet)
		existingAgent, err := r.Find(ctx, enum.AgentType(agent.Type))
		if err != nil {
			tracing.TraceErr(span, err)
			return fmt.Errorf("failed to check existing agent: %w", err)
		}

		// Skip if agent already exists
		if existingAgent != nil {
			continue
		}

		// Create new agent
		if _, err := r.Create(ctx, agent); err != nil {
			tracing.TraceErr(span, err)
			return fmt.Errorf("failed to initialize agent: %w", err)
		}
	}

	return nil
}

func registerVisitorIdentityAgent() entity.AgentRegistry {
	return entity.AgentRegistry{
		Type:         enum.AgentVisitorID.String(),
		Name:         "Identify website visitors",
		Goal:         enum.AgentGoalIdentifyVisitors.String(),
		Icon:         "",
		Capabilities: "",
		IsActive:     true,
	}
}

