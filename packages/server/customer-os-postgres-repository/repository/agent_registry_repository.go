package postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/opentracing/opentracing-go"
	"gorm.io/gorm"

	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
)

var (
	ErrAgentIDMissing    = errors.New("agent ID is missing")
	ErrAgentNotFound     = errors.New("agent not found")
	ErrAgentCreateFailed = errors.New("failed to create agent")
)

type AgentRegistryRepository interface {
	Create(ctx context.Context, agent postgres_entity.AgentRegistry, plays []postgres_entity.AgentPlay) (*postgres_entity.AgentRegistry, error)
	Update(ctx context.Context, agent postgres_entity.AgentRegistry, plays []postgres_entity.AgentPlay) (*postgres_entity.AgentRegistry, error)
	FindByType(ctx context.Context, agentType enum.AgentType) (*postgres_entity.AgentRegistry, error)
	FindAll(ctx context.Context) ([]postgres_entity.AgentRegistry, error)
	FindPlay(ctx context.Context, agentType enum.AgentType, triggerEvent enum.AgentListenerEvent) (*postgres_entity.AgentPlay, error)
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

func (r *agentRegistryRepository) FindAll(ctx context.Context) ([]postgres_entity.AgentRegistry, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentRegistryRepository.FindAll")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	var agents []postgres_entity.AgentRegistry
	err := r.gormDb.WithContext(ctx).
		Where("is_active = ?", true).
		Find(&agents).Error
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, fmt.Errorf("failed to find agents: %w", err)
	}

	return agents, nil
}

func (r *agentRegistryRepository) FindPlay(ctx context.Context, agentType enum.AgentType, triggerEvent enum.AgentListenerEvent) (*postgres_entity.AgentPlay, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentRegistryRepository.FindPlay")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	var play postgres_entity.AgentPlay
	err := r.gormDb.WithContext(ctx).
		Where("agent_type = ? AND trigger_type = ?", agentType.String(), triggerEvent.String()).
		First(&play).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		tracing.TraceErr(span, err)
		return nil, fmt.Errorf("failed to find play for agent type %s and trigger %s: %w",
			agentType.String(), triggerEvent.String(), err)
	}

	return &play, nil
}

func (r *agentRegistryRepository) FindByType(ctx context.Context, agentType enum.AgentType) (*postgres_entity.AgentRegistry, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentRegistryRepository.FindByType")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	var agent postgres_entity.AgentRegistry
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

func (r *agentRegistryRepository) Create(ctx context.Context, agent postgres_entity.AgentRegistry, plays []postgres_entity.AgentPlay) (*postgres_entity.AgentRegistry, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentRegistryRepository.Create")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	err := r.gormDb.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Check for existing agent with same type
		var count int64
		if err := tx.Model(&postgres_entity.AgentRegistry{}).
			Where("type = ? AND is_active = ?", agent.Type, true).
			Count(&count).Error; err != nil {
			return err
		}

		if count > 0 {
			return fmt.Errorf("agent with type %s already exists", agent.Type)
		}

		// Create the agent
		if err := tx.Create(&agent).Error; err != nil {
			return err
		}

		// Create plays with the agent ID
		for i := range plays {
			plays[i].AgentRegistryID = agent.ID
		}
		if len(plays) > 0 {
			if err := tx.Create(&plays).Error; err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, fmt.Errorf("%w: %v", ErrAgentCreateFailed, err)
	}

	return &agent, nil
}

func (r *agentRegistryRepository) Update(ctx context.Context, agent postgres_entity.AgentRegistry, plays []postgres_entity.AgentPlay) (*postgres_entity.AgentRegistry, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentRegistryRepository.Update")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	if agent.ID == "" {
		return nil, ErrAgentIDMissing
	}

	var updatedAgent postgres_entity.AgentRegistry
	err := r.gormDb.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Check if agent exists
		if err := tx.First(&postgres_entity.AgentRegistry{}, "id = ?", agent.ID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrAgentNotFound
			}
			return err
		}

		// Update agent
		result := tx.Model(&postgres_entity.AgentRegistry{}).
			Where("id = ?", agent.ID).
			Updates(&agent)

		if result.Error != nil {
			return result.Error
		}

		if result.RowsAffected == 0 {
			return fmt.Errorf("no rows affected during update")
		}

		// Delete existing plays
		if err := tx.Delete(&postgres_entity.AgentPlay{}, "agent_registry_id = ?", agent.ID).Error; err != nil {
			return err
		}

		// Create new plays
		if len(plays) > 0 {
			for i := range plays {
				plays[i].AgentRegistryID = agent.ID
			}
			if err := tx.Create(&plays).Error; err != nil {
				return err
			}
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
