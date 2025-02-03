package postgres_repository

import (
	"context"
	"errors"
	"github.com/opentracing/opentracing-go/log"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
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

	agent.ID = utils.GenerateNanoIdWithPrefix("agent", 16)

	var created postgres_entity.Agent
	err := f.gormDb.Create(&agent).Scan(&created).Error
	if err != nil {
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
	query := f.gormDb.Where("tenant = ? ", tenant)
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
	// Convert enum array to []string
	types := make([]string, len(agentTypes))
	for i, agentType := range agentTypes {
		types[i] = agentType.String()
	}

	query := f.gormDb.Where("tenant = ? ", tenant)
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
	// Convert enum array to []string
	types := make([]string, len(agentTypes))
	for i, agentType := range agentTypes {
		types[i] = agentType.String()
	}

	query := f.gormDb.Where("tenant = ? AND is_active = ? AND configured = ?", tenant, true, true)
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
	// Convert enum array to []string
	types := make([]string, len(agentTypes))
	for i, agentType := range agentTypes {
		types[i] = agentType.String()
	}

	query := f.gormDb.Where("is_active = ? AND configured = ?", true, true)
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

	var automation postgres_entity.Agent
	err := f.gormDb.
		Where(&agent).
		Where("is_active = true").
		First(&automation).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		tracing.TraceErr(span, err)
		return nil, err
	}
	return &automation, nil
}

func (f *agentsRepository) Update(ctx context.Context, agent postgres_entity.Agent) (*postgres_entity.Agent, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentRepository.Update")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	tracing.LogObjectAsJson(span, "agentEntity", agent)

	if agent.ID == "" {
		err := errors.New("agent ID is missing") // Fixed error message to match the context
		tracing.TraceErr(span, err)
		return nil, err
	}

	var updatedAgents postgres_entity.Agent
	err := f.gormDb.
		Model(&postgres_entity.Agent{}).
		Where("id = ?", agent.ID).
		Save(agent).
		First(&updatedAgents, "id = ?", agent.ID).
		Error
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return &updatedAgents, nil
}
