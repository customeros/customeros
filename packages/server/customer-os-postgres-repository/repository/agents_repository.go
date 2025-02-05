package postgres_repository

import (
	"context"
	"errors"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"gorm.io/gorm"

	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
)

type AgentRepository interface {
	Create(ctx context.Context, agent postgres_entity.Agent) (*postgres_entity.Agent, error)
	Find(ctx context.Context, agent postgres_entity.Agent) (*postgres_entity.Agent, error)
	GetById(ctx context.Context, id string) (*postgres_entity.Agent, error)
	GetAll(ctx context.Context) ([]*postgres_entity.Agent, error)
	GetAllAgentsByTypes(ctx context.Context, agents []enum.AgentType) ([]postgres_entity.Agent, error)
	GetActiveConfiguredAgentsByTypes(ctx context.Context, agents []enum.AgentType) ([]postgres_entity.Agent, error)
	GetActiveConfiguredAgentsByTypesCrossTenant(ctx context.Context, agents []enum.AgentType) ([]postgres_entity.Agent, error)
	Update(ctx context.Context, agent postgres_entity.Agent) (*postgres_entity.Agent, error)
	UpsertCapabilities(ctx context.Context, agentID string, capabilities []postgres_entity.Capability) error
	FindCapability(ctx context.Context, agentID string, capabilityType enum.AgentCapability) (*postgres_entity.Capability, error)
}

type agentsRepository struct {
	gormDb *gorm.DB
}

func NewAgentRepository(gormDb *gorm.DB) AgentRepository {
	return &agentsRepository{gormDb: gormDb}
}

func (f *agentsRepository) GetById(ctx context.Context, id string) (*postgres_entity.Agent, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentRepository.GetById")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	span.LogFields(log.String("id", id))

	var agent postgres_entity.Agent
	err := f.gormDb.
		Preload("Capabilities").
		Where("id = ?", id).
		First(&agent).
		Error
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return &agent, nil
}

func (f *agentsRepository) Create(ctx context.Context, agent postgres_entity.Agent) (*postgres_entity.Agent, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentRepository.Create")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	err := f.gormDb.Transaction(func(tx *gorm.DB) error {
		// Create agent first
		if err := tx.Create(&agent).Error; err != nil {
			return err
		}

		// Create capabilities if any exist
		if len(agent.Capabilities) > 0 {
			for i := range agent.Capabilities {
				agent.Capabilities[i].AgentID = agent.ID
			}
			if err := tx.Create(&agent.Capabilities).Error; err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	// Fetch the created agent with capabilities
	var created postgres_entity.Agent
	if err := f.gormDb.Preload("Capabilities").First(&created, "id = ?", agent.ID).Error; err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return &created, nil
}

func (f *agentsRepository) GetAll(ctx context.Context) ([]*postgres_entity.Agent, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentRepository.GetAll")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	tenant := common.GetTenantFromContext(ctx)
	if tenant == "" {
		err := errors.New("Tenant not set")
		tracing.TraceErr(span, err)
		return nil, err
	}

	var agents []*postgres_entity.Agent
	query := f.gormDb.
		Preload("Capabilities").
		Where("tenant = ?", tenant)
	err := query.Find(&agents).Error
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	return agents, nil
}

func (f *agentsRepository) GetAllAgentsByTypes(ctx context.Context, agentTypes []enum.AgentType) ([]postgres_entity.Agent, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentRepository.GetAllAgentsByTypes")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	tenant := common.GetTenantFromContext(ctx)
	if tenant == "" {
		err := errors.New("Tenant not set")
		tracing.TraceErr(span, err)
		return nil, err
	}

	var records []postgres_entity.Agent
	types := make([]string, len(agentTypes))
	for i, agentType := range agentTypes {
		types[i] = agentType.String()
	}

	query := f.gormDb.
		Preload("Capabilities").
		Where("tenant = ?", tenant)
	if len(types) > 0 {
		query = query.Where("type IN (?)", types)
	}

	err := query.Find(&records).Error
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	return records, nil
}

func (f *agentsRepository) GetActiveConfiguredAgentsByTypes(ctx context.Context, agentTypes []enum.AgentType) ([]postgres_entity.Agent, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentRepository.GetActiveConfiguredAgentsByTypes")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	tenant := common.GetTenantFromContext(ctx)
	if tenant == "" {
		err := errors.New("tenant not set")
		tracing.TraceErr(span, err)
		return nil, err
	}

	var records []postgres_entity.Agent
	types := make([]string, len(agentTypes))
	for i, agentType := range agentTypes {
		types[i] = agentType.String()
	}

	query := f.gormDb.
		Preload("Capabilities").
		Where("tenant = ? AND is_active = ? AND configured = ?", tenant, true, true)
	if len(types) > 0 {
		query = query.Where("type IN (?)", types)
	}

	err := query.Find(&records).Error
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	return records, nil
}

func (f *agentsRepository) GetActiveConfiguredAgentsByTypesCrossTenant(ctx context.Context, agentTypes []enum.AgentType) ([]postgres_entity.Agent, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentRepository.GetActiveConfiguredAgentsByTypesCrossTenant")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	var records []postgres_entity.Agent
	types := make([]string, len(agentTypes))
	for i, agentType := range agentTypes {
		types[i] = agentType.String()
	}

	query := f.gormDb.
		Preload("Capabilities").
		Where("is_active = ? AND configured = ?", true, true)
	if len(types) > 0 {
		query = query.Where("type IN (?)", types)
	}

	err := query.Find(&records).Error
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	span.LogFields(log.Int("count", len(records)))
	return records, nil
}

func (f *agentsRepository) Find(ctx context.Context, agent postgres_entity.Agent) (*postgres_entity.Agent, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentRepository.Find")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	var result postgres_entity.Agent
	err := f.gormDb.
		Preload("Capabilities").
		Where(&agent).
		Where("is_active = true").
		First(&result).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		tracing.TraceErr(span, err)
		return nil, err
	}
	return &result, nil
}

func (f *agentsRepository) Update(ctx context.Context, agent postgres_entity.Agent) (*postgres_entity.Agent, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentRepository.Update")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	tracing.LogObjectAsJson(span, "agentEntity", agent)

	if agent.ID == "" {
		err := errors.New("agent ID is missing")
		tracing.TraceErr(span, err)
		return nil, err
	}

	err := f.gormDb.Transaction(func(tx *gorm.DB) error {
		// Update agent
		if err := tx.Model(&agent).Updates(&agent).Error; err != nil {
			return err
		}

		// Handle capabilities if they exist
		if len(agent.Capabilities) > 0 {
			if err := f.UpsertCapabilities(ctx, agent.ID, agent.Capabilities); err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	// Fetch updated agent with capabilities
	var updatedAgent postgres_entity.Agent
	if err := f.gormDb.Preload("Capabilities").First(&updatedAgent, "id = ?", agent.ID).Error; err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return &updatedAgent, nil
}

func (f *agentsRepository) UpsertCapabilities(ctx context.Context, agentID string, capabilities []postgres_entity.Capability) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentRepository.UpsertCapabilities")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	return f.gormDb.Transaction(func(tx *gorm.DB) error {
		// Delete existing capabilities
		if err := tx.Where("agent_id = ?", agentID).Delete(&postgres_entity.Capability{}).Error; err != nil {
			return err
		}

		// Set agent ID for all capabilities
		for i := range capabilities {
			capabilities[i].AgentID = agentID
		}

		// Create new capabilities
		if len(capabilities) > 0 {
			if err := tx.Create(&capabilities).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (f *agentsRepository) FindCapability(ctx context.Context, agentID string, capabilityType enum.AgentCapability) (*postgres_entity.Capability, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentRepository.FindActiveCapability")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	span.LogFields(
		log.String("agentID", agentID),
		log.String("capabilityType", capabilityType.String()),
	)

	var capability postgres_entity.Capability
	err := f.gormDb.
		Where("agent_id = ? AND type = ? AND active = ?", agentID, capabilityType.String(), true).
		First(&capability).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		tracing.TraceErr(span, err)
		return nil, err
	}

	return &capability, nil
}
