package postgres_repository

import (
	"context"
	"errors"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/coserrors"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"gorm.io/gorm"

	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
)

type AgentRepository interface {
	Create(ctx context.Context, agent postgres_entity.Agent) (*postgres_entity.Agent, error)
	GetById(ctx context.Context, id string) (*postgres_entity.Agent, error)
	GetAll(ctx context.Context) ([]*postgres_entity.Agent, error)
	GetAllAgentsByTypes(ctx context.Context, agents []enum.AgentType) ([]postgres_entity.Agent, error)
	GetAllAgentsByTypesCrossTenant(ctx context.Context, agents []enum.AgentType) ([]postgres_entity.Agent, error)
	GetActiveConfiguredAgentsByTypes(ctx context.Context, agents []enum.AgentType) ([]postgres_entity.Agent, error)
	GetActiveConfiguredAgentsByUserAndType(ctx context.Context, agents []enum.AgentType) ([]postgres_entity.Agent, error)
	GetActiveConfiguredAgentsByTypesCrossTenant(ctx context.Context, agents []enum.AgentType) ([]postgres_entity.Agent, error)
	Update(ctx context.Context, agent postgres_entity.Agent) (*postgres_entity.Agent, error)
	UpdateCapabilities(ctx context.Context, capabilities []postgres_entity.Capability) error
	UpdateListeners(ctx context.Context, capabilities []postgres_entity.Listener) error
	FindCapability(ctx context.Context, agentID string, capabilityType enum.AgentCapability) (*postgres_entity.Capability, error)
	Delete(ctx context.Context, id string) error
}

type agentsRepository struct {
	gormDb *gorm.DB
}

func NewAgentRepository(gormDb *gorm.DB) AgentRepository {
	return &agentsRepository{gormDb: gormDb}
}

func (f *agentsRepository) GetById(ctx context.Context, id string) (*postgres_entity.Agent, error) {
	spans, ctx := telemetry.StartPostgresSpan(ctx, "AgentRepository.GetById")
	defer spans.Finish()
	spans.LogKV("id", id)

	var agent postgres_entity.Agent
	err := f.gormDb.
		Preload("Capabilities", func(db *gorm.DB) *gorm.DB {
			return db.Order("position ASC")
		}).
		Preload("Listeners", func(db *gorm.DB) *gorm.DB {
			return db.Order("position ASC")
		}).
		Where("id = ?", id).
		First(&agent).
		Error
	if err != nil {
		spans.LogKV("result.found", false)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		spans.TraceError(err)
		return nil, err
	}

	spans.LogKV("result.found", true)
	return &agent, nil
}

func (f *agentsRepository) Create(ctx context.Context, agent postgres_entity.Agent) (*postgres_entity.Agent, error) {
	spans, ctx := telemetry.StartPostgresSpan(ctx, "AgentRepository.Create")
	defer spans.Finish()

	if agent.Scope == enum.AgentScopePersonal {
		agent.Owner = common.GetUserIdFromContext(ctx)
		if agent.Owner == "" {
			err := errors.New("UserID not set on context")
			spans.TraceError(err)
			return nil, err
		}
	}

	agent.Tenant = common.GetTenantFromContext(ctx)
	if agent.Tenant == "" {
		err := errors.New("tenant not set on context")
		spans.TraceError(err)
		return nil, err
	}

	err := f.gormDb.Transaction(func(tx *gorm.DB) error {
		// Create agent first
		if err := tx.Create(&agent).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}

	// Fetch the created agent with capabilities
	var created postgres_entity.Agent
	if err := f.gormDb.Preload("Capabilities", func(db *gorm.DB) *gorm.DB {
		return db.Order("position ASC")
	}).Preload("Listeners", func(db *gorm.DB) *gorm.DB {
		return db.Order("position ASC")
	}).First(&created, "id = ?", agent.ID).Error; err != nil {
		spans.TraceError(err)
		return nil, err
	}

	return &created, nil
}

func (f *agentsRepository) GetAll(ctx context.Context) ([]*postgres_entity.Agent, error) {
	spans, ctx := telemetry.StartPostgresSpan(ctx, "AgentRepository.GetAll")
	defer spans.Finish()

	tenant := common.GetTenantFromContext(ctx)
	if tenant == "" {
		err := errors.New("tenant not set")
		spans.TraceError(err)
		return nil, err
	}

	userId := common.GetUserIdFromContext(ctx)

	var agents []*postgres_entity.Agent
	query := f.gormDb.
		Preload("Capabilities", func(db *gorm.DB) *gorm.DB {
			return db.Order("position ASC")
		}).
		Preload("Listeners", func(db *gorm.DB) *gorm.DB {
			return db.Order("position ASC")
		}).
		Where(
			"(tenant = ? AND agent_scope = ? ) OR (tenant = ? AND owner = ? AND owner <> '' AND agent_scope = ?)",
			tenant, enum.AgentScopeWorkspace, tenant, userId, enum.AgentScopePersonal,
		)

	err := query.Find(&agents).Error
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}

	spans.LogKV("result.count", len(agents))
	return agents, nil
}

func (f *agentsRepository) GetAllAgentsByTypes(ctx context.Context, agentTypes []enum.AgentType) ([]postgres_entity.Agent, error) {
	spans, ctx := telemetry.StartPostgresSpan(ctx, "AgentRepository.GetAllAgentsByTypes")
	defer spans.Finish()

	tenant := common.GetTenantFromContext(ctx)
	if tenant == "" {
		err := errors.New("tenant not set")
		spans.TraceError(err)
		return nil, err
	}

	var records []postgres_entity.Agent
	types := make([]string, len(agentTypes))
	for i, agentType := range agentTypes {
		types[i] = agentType.String()
	}

	query := f.gormDb.
		Preload("Capabilities", func(db *gorm.DB) *gorm.DB {
			return db.Order("position ASC")
		}).
		Preload("Listeners", func(db *gorm.DB) *gorm.DB {
			return db.Order("position ASC")
		}).
		Where("tenant = ?", tenant)
	if len(types) > 0 {
		query = query.Where("type IN (?)", types)
	}

	err := query.Find(&records).Error
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}
	return records, nil
}

func (f *agentsRepository) GetAllAgentsByTypesCrossTenant(ctx context.Context, agentTypes []enum.AgentType) ([]postgres_entity.Agent, error) {
	spans, ctx := telemetry.StartPostgresSpan(ctx, "AgentRepository.GetAllAgentsByTypesCrossTenant")
	defer spans.Finish()
	spans.LogKV("agentTypes", agentTypes)

	var records []postgres_entity.Agent
	types := make([]string, len(agentTypes))
	for i, agentType := range agentTypes {
		types[i] = agentType.String()
	}

	query := f.gormDb.
		Preload("Capabilities", func(db *gorm.DB) *gorm.DB {
			return db.Order("position ASC")
		}).
		Preload("Listeners", func(db *gorm.DB) *gorm.DB {
			return db.Order("position ASC")
		})
	if len(types) > 0 {
		query = query.Where("type IN (?)", types)
	}

	err := query.Find(&records).Error
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}
	spans.LogKV("result.count", len(records))
	return records, nil
}

func (f *agentsRepository) GetActiveConfiguredAgentsByTypes(ctx context.Context, agentTypes []enum.AgentType) ([]postgres_entity.Agent, error) {
	spans, ctx := telemetry.StartPostgresSpan(ctx, "AgentRepository.GetActiveConfiguredAgentsByTypes")
	defer spans.Finish()

	tenant := common.GetTenantFromContext(ctx)
	if tenant == "" {
		err := errors.New("tenant not set")
		spans.TraceError(err)
		return nil, err
	}

	var records []postgres_entity.Agent
	types := make([]string, len(agentTypes))
	for i, agentType := range agentTypes {
		types[i] = agentType.String()
	}

	query := f.gormDb.
		Preload("Capabilities", func(db *gorm.DB) *gorm.DB {
			return db.Order("position ASC")
		}).
		Preload("Listeners", func(db *gorm.DB) *gorm.DB {
			return db.Order("position ASC")
		}).
		Where("tenant = ? AND is_active = ? AND configured = ?", tenant, true, true)
	if len(types) > 0 {
		query = query.Where("type IN (?)", types)
	}

	err := query.Find(&records).Error
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}
	spans.LogKV("result.count", len(records))
	return records, nil
}

func (f *agentsRepository) GetActiveConfiguredAgentsByUserAndType(ctx context.Context, agentTypes []enum.AgentType) ([]postgres_entity.Agent, error) {
	spans, ctx := telemetry.StartPostgresSpan(ctx, "AgentRepository.GetActiveConfiguredAgentsByUserAndType")
	defer spans.Finish()

	tenant := common.GetTenantFromContext(ctx)
	if tenant == "" {
		spans.TraceError(coserrors.ErrTenantNotSet)
		return nil, coserrors.ErrTenantNotSet
	}
	user := common.GetUserIdFromContext(ctx)
	if user == "" {
		spans.TraceError(coserrors.ErrUserIDNotSet)
		return nil, coserrors.ErrUserIDNotSet
	}

	var records []postgres_entity.Agent
	types := make([]string, len(agentTypes))
	for i, agentType := range agentTypes {
		types[i] = agentType.String()
	}

	query := f.gormDb.
		Preload("Capabilities", func(db *gorm.DB) *gorm.DB {
			return db.Order("position ASC")
		}).
		Preload("Listeners", func(db *gorm.DB) *gorm.DB {
			return db.Order("position ASC")
		}).
		Where("tenant = ? AND is_active = ? AND configured = ? AND owner = ?", tenant, true, true, user)
	if len(types) > 0 {
		query = query.Where("type IN (?)", types)
	}

	err := query.Find(&records).Error
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}
	return records, nil
}

func (f *agentsRepository) GetActiveConfiguredAgentsByTypesCrossTenant(ctx context.Context, agentTypes []enum.AgentType) ([]postgres_entity.Agent, error) {
	spans, ctx := telemetry.StartPostgresSpan(ctx, "AgentRepository.GetActiveConfiguredAgentsByTypesCrossTenant")
	defer spans.Finish()
	spans.LogKV("agentTypes", agentTypes)

	var records []postgres_entity.Agent
	types := make([]string, len(agentTypes))
	for i, agentType := range agentTypes {
		types[i] = agentType.String()
	}

	query := f.gormDb.
		Preload("Capabilities", func(db *gorm.DB) *gorm.DB {
			return db.Order("position ASC")
		}).
		Preload("Listeners", func(db *gorm.DB) *gorm.DB {
			return db.Order("position ASC")
		}).
		Where("is_active = ? AND configured = ?", true, true)
	if len(types) > 0 {
		query = query.Where("type IN (?)", types)
	}

	err := query.Find(&records).Error
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}
	spans.LogKV("result.count", len(records))
	return records, nil
}

func (f *agentsRepository) Update(ctx context.Context, agent postgres_entity.Agent) (*postgres_entity.Agent, error) {
	spans, ctx := telemetry.StartPostgresSpan(ctx, "AgentRepository.Update")
	defer spans.Finish()
	spans.LogObjectAsJson("agentEntity", agent)

	if agent.ID == "" {
		err := errors.New("agent ID is missing")
		spans.TraceError(err)
		return nil, err
	}

	err := f.gormDb.Transaction(func(tx *gorm.DB) error {
		// Update agent
		if err := tx.Model(&agent).Omit("Capabilities", "Listeners").Save(&agent).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}

	// Fetch updated agent with capabilities
	var updatedAgent postgres_entity.Agent
	if err := f.gormDb.
		Preload("Capabilities", func(db *gorm.DB) *gorm.DB {
			return db.Order("position ASC")
		}).
		Preload("Listeners", func(db *gorm.DB) *gorm.DB {
			return db.Order("position ASC")
		}).
		First(&updatedAgent, "id = ?", agent.ID).Error; err != nil {
		spans.TraceError(err)
		return nil, err
	}

	return &updatedAgent, nil
}

func (f *agentsRepository) UpdateCapabilities(ctx context.Context, capabilities []postgres_entity.Capability) error {
	spans, ctx := telemetry.StartPostgresSpan(ctx, "AgentRepository.UpdateCapabilities")
	defer spans.Finish()
	spans.LogKV("capabilities", capabilities)

	return f.gormDb.Transaction(func(tx *gorm.DB) error {
		for _, capability := range capabilities {
			// Skip records with no ID
			if capability.ID == "" {
				continue
			}
			// Save() will update the record with the primary key cap.ID.
			if err := tx.Save(&capability).Error; err != nil {
				spans.TraceError(err)
				return err
			}
		}
		return nil
	})
}

func (f *agentsRepository) UpdateListeners(ctx context.Context, listeners []postgres_entity.Listener) error {
	spans, ctx := telemetry.StartPostgresSpan(ctx, "AgentRepository.UpdateListeners")
	defer spans.Finish()
	spans.LogKV("listeners", listeners)

	return f.gormDb.Transaction(func(tx *gorm.DB) error {
		for _, listener := range listeners {
			// Skip records with no ID
			if listener.ID == "" {
				continue
			}
			// Save() will update the record with the primary key listener.ID.
			if err := tx.Save(&listener).Error; err != nil {
				spans.TraceError(err)
				return err
			}
		}
		return nil
	})
}

func (f *agentsRepository) FindCapability(ctx context.Context, agentID string, capabilityType enum.AgentCapability) (*postgres_entity.Capability, error) {
	spans, ctx := telemetry.StartPostgresSpan(ctx, "AgentRepository.FindActiveCapability")
	defer spans.Finish()
	spans.LogKV("agentID", agentID)
	spans.LogKV("capabilityType", capabilityType.String())

	var capability postgres_entity.Capability
	err := f.gormDb.
		Where("agent_id = ? AND type = ?", agentID, capabilityType.String()).
		First(&capability).
		Error
	if err != nil {
		spans.LogKV("result.found", false)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		spans.TraceError(err)
		return nil, err
	}

	spans.LogKV("result.found", true)
	spans.LogKV("result.capabilityId", capability.ID)
	return &capability, nil
}

func (f *agentsRepository) Delete(ctx context.Context, id string) error {
	spans, ctx := telemetry.StartPostgresSpan(ctx, "AgentRepository.Delete")
	defer spans.Finish()
	spans.LogKV("id", id)

	tenant := common.GetTenantFromContext(ctx)
	if tenant == "" {
		err := errors.New("tenant not set")
		spans.TraceError(err)
		return err
	}

	return f.gormDb.Transaction(func(tx *gorm.DB) error {
		// First verify the agent exists and belongs to the tenant
		var agent postgres_entity.Agent
		if err := tx.Where("id = ? AND tenant = ?", id, tenant).First(&agent).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("agent not found or access denied")
			}
			spans.TraceError(err)
			return err
		}

		// Delete capabilities
		if err := tx.Where("agent_id = ?", id).Delete(&postgres_entity.Capability{}).Error; err != nil {
			spans.TraceError(err)
			return err
		}

		// Delete listeners
		if err := tx.Where("agent_id = ?", id).Delete(&postgres_entity.Listener{}).Error; err != nil {
			spans.TraceError(err)
			return err
		}

		// Delete the agent
		if err := tx.Delete(&postgres_entity.Agent{}, "id = ?", id).Error; err != nil {
			spans.TraceError(err)
			return err
		}

		return nil
	})
}
