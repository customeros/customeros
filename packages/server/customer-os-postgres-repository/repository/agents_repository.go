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

type AgentsRepository interface {
	Create(ctx context.Context, agent postgres_entity.Agents) (*postgres_entity.Agents, error)
	Find(ctx context.Context, agent postgres_entity.Agents) (*postgres_entity.Agents, error)
	GetById(ctx context.Context, id string) (*postgres_entity.Agents, error)
	GetAll(ctx context.Context) ([]*postgres_entity.Agents, error)
	FindAllFromAgentsList(ctx context.Context, agents []enum.AgentType) ([]postgres_entity.Agents, error)
	Update(ctx context.Context, agent postgres_entity.Agents) (*postgres_entity.Agents, error)
}

type agentsRepository struct {
	gormDb *gorm.DB
}

func NewAgentsRepository(gormDb *gorm.DB) AgentsRepository {
	return &agentsRepository{gormDb: gormDb}
}

func (f *agentsRepository) GetById(ctx context.Context, id string) (*postgres_entity.Agents, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentsRepository.GetById")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	span.LogFields(log.String("id", id))

	var agent postgres_entity.Agents
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

func (f *agentsRepository) Create(ctx context.Context, agent postgres_entity.Agents) (*postgres_entity.Agents, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentsRepository.Create")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	agent.ID = utils.GenerateNanoIdWithPrefix("agent", 16)

	var created postgres_entity.Agents
	err := f.gormDb.Create(&agent).Scan(&created).Error
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return &created, nil
}

func (f *agentsRepository) GetAll(ctx context.Context) ([]*postgres_entity.Agents, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentsRepository.GetAll")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	var agents []*postgres_entity.Agents
	err := f.gormDb.Find(&agents).Error
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	return agents, nil
}

// TODO alexb rename it
func (f *agentsRepository) FindAllFromAgentsList(ctx context.Context, agents []enum.AgentType) ([]postgres_entity.Agents, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentsRepository.FindAllFromAgentsList")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	tenant := common.GetTenantFromContext(ctx)
	if tenant == "" {
		err := errors.New("Tenant not set")
		tracing.TraceErr(span, err)
		return nil, err
	}

	var records []postgres_entity.Agents
	// Convert enum array to []string
	agentIds := make([]string, len(agents))
	for i, agent := range agents {
		agentIds[i] = agent.String()
	}

	query := f.gormDb.Where("tenant = ? AND is_active = ?", tenant, true)
	if len(agentIds) > 0 {
		query = query.Where("type IN (?)", agentIds)
	}

	err := query.Find(&records).Error
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	return records, nil
}

func (f *agentsRepository) Find(ctx context.Context, agent postgres_entity.Agents) (*postgres_entity.Agents, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentsRepository.Find")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	var automation postgres_entity.Agents
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

func (f *agentsRepository) Update(ctx context.Context, agent postgres_entity.Agents) (*postgres_entity.Agents, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentsRepository.Update")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	if agent.ID == "" {
		err := errors.New("agent ID is missing") // Fixed error message to match the context
		tracing.TraceErr(span, err)
		return nil, err
	}

	var updatedAgents postgres_entity.Agents
	err := f.gormDb.
		Model(&postgres_entity.Agents{}).
		Where("id = ?", agent.ID).
		Updates(agent).
		First(&updatedAgents, "id = ?", agent.ID).
		Error
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return &updatedAgents, nil
}
