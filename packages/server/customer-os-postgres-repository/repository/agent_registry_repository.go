package repository

import (
	"context"
	"errors"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
	"gorm.io/gorm"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
)

type AgentRegistryRepository interface {
	Initialize(ctx context.Context) error
	Create(ctx context.Context, agent entity.AgentRegistry) (*entity.AgentRegistry, error)
	Find(ctx context.Context, agent entity.AgentRegistry) (*entity.AgentRegistry, error)
	FindAll(ctx context.Context) ([]entity.AgentRegistry, error)
}

type agentRegistryRepository struct {
	gormDb *gorm.DB
}

func NewAgentRegistryRepository(gormDb *gorm.DB) AgentRegistryRepository {
	return &agentRegistryRepository{gormDb: gormDb}
}

func (r *agentRegistryRepository) FindAll(ctx context.Context) ([]entity.AgentRegistry, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentRegistryRepository.FindAll")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	var actions []entity.AgentRegistry
	err := r.gormDb.WithContext(ctx).
		Where("is_active = true").
		Find(&actions).Error
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return actions, nil
}

func (r *agentRegistryRepository) Find(ctx context.Context, agentID string) (*entity.AgentRegistry, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentRegistryRepository.Find")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	if agentID == "" {
		return nil, errors.New("empty agentID")
	}

	var action entity.AgentRegistry
	query := r.gormDb.WithContext(ctx).
		Where("is_active = true").
		Where("id = ?", enum.AgentVisitorID)

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

func (r *agentRegistryRepository) Create(ctx context.Context, agent entity.AgentRegistry) (*entity.AgentRegistry, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentRegistryRepository.Create")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	err := r.gormDb.WithContext(ctx).Create(&agent).Error
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return &agent, nil
}

func (r *agentRegistryRepository) Initialize(ctx context.Context) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentRegistryRepository.Initialize")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	requiredAgents := []entity.AgentRegistry{
		RegisterVisitorIdentityAgent(),
		// ... add more here
	}

	for _, agent := range requiredAgents {
		// Look for existing agent by ID
		existingAgent, err := r.Find(ctx, entity.AgentRegistry{ID: agent.ID})
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}

		if existingAgent != nil && existingAgent.ID != "" {
			// Check if update is needed by comparing fields
			if needsUpdate(existingAgent, &agent) {
				_, updateErr := r.Update(ctx, agent)
				if updateErr != nil {
					tracing.TraceErr(span, updateErr)
					return updateErr
				}
			}
			continue
		}

		// Create new agent if it doesn't exist
		_, createErr := r.Create(ctx, agent)
		if createErr != nil {
			tracing.TraceErr(span, createErr)
			return createErr
		}
	}
	return nil
}

// Helper function to check if agent needs update
func needsUpdate(existing *entity.AgentRegistry, new *entity.AgentRegistry) bool {
	return existing.Name != new.Name ||
		existing.Capabilities != new.Capabilities ||
		existing.ConfigSchema != new.ConfigSchema ||
		existing.IsActive != new.IsActive
}

// agent schema definitions here

func RegisterVisitorIdentityAgent() entity.AgentRegistry {
	return entity.AgentRegistry{
		ID:   "visitor-identity-agent",
		Name: "Identify website visitors",
		Capabilities: `{
            "capabilities": [
                {
                    "name": "visitor_identification",
                    "action": "Track and identify website visitors"
                },
                {
                    "name": "session_tracking",
                    "action": "Log page views and session duration"
                },
                {
                    "name": "lead_creation",
                    "action": "Create identified organizations as leads"
                },
                {
                    "name": "lead_enrishment",
                    "action": "Enrich leads with qualification data"
                },
                {
                    "name": "intent_signals",
                    "action": "Analyze behavior for intent signals"
                },
                {
                    "name": "slack_notifications",
                    "action": "Send Slack notification",
                    "optional": true
                }
            ]
        }`,
		ConfigSchema: utils.StringPtr(`{
            "type": "object",
            "required": ["websites"],
            "properties": {
                "websites": {
                    "type": "array",
                    "description": "Websites where tracker is installed",
                    "items": {
                        "type": "string"
                    }
                },
                "slackEnabled": {
                    "type": "boolean",
                    "description": "Enable Slack notifications",
                    "default": false
                },
                "slackChannelId": {
                    "type": "string",
                    "description": "Slack channel ID to send notifications to",
                },
                "notificationCooldownHours": {
                    "type": "integer",
                    "description": "Minimum hours between notifications for the same company",
                    "minimum": 1,
                    "default": 12
                }
            },
            "dependencies": {
                "slack_enabled": {
                    "if": {
                        "properties": { "slackEnabled": { "const": true } }
                    },
                    "then": {
                        "required": ["slackChannelId"]
                    }
                }
            }
        }`),
		IsActive: "true",
	}
}
