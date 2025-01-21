package agent_capability

import (
	"context"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/multierr"
)

func (c *agentCapabilityService) InitCapabilityRegistry() error {
	ctx := context.Background()
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentCapabilityService.InitCapabilityRegistry")
	defer span.Finish()
	tracing.TagComponentService(span)

	// register all agent capabilities here
	requiredCapabilities := []postgres_entity.AgentCapabilityRegistry{
		c.buildCapability(
			enum.CapabilityIdentifyWebVisitor,
			"Identify website visitors",
		),
		c.buildCapability(
			enum.CapabilityAnalyzeWebSessionIntent,
			"Analyze web session for intent signals",
		),
		c.buildCapability(
			enum.CapabilitySendSlackNotification,
			"Send Slack notification",
		),
	}

	var errs error
	for _, capability := range requiredCapabilities {
		err := c.RegisterCapability(ctx, capability)
		errs = multierr.Append(errs, err)
	}
	return errs
}

func (c *agentCapabilityService) RegisterCapability(ctx context.Context, capability postgres_entity.AgentCapabilityRegistry) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentCapabilityService.RegisterCapability")
	defer span.Finish()
	tracing.TagComponentService(span)

	validCapability, err := enum.GetAgentCapability(capability.Type)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	existingCapability, err := c.postgresRepositories.AgentCapabilityRegistryRepository.Find(ctx, validCapability)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if existingCapability != nil {
		return nil // capability already exists in db
	}

	_, createErr := c.postgresRepositories.AgentCapabilityRegistryRepository.Create(ctx, capability)
	if createErr != nil {
		tracing.TraceErr(span, createErr)
		return createErr
	}
	return nil
}

func (c *agentCapabilityService) buildCapability(capability enum.AgentCapabilityType, description string) postgres_entity.AgentCapabilityRegistry {
	return postgres_entity.AgentCapabilityRegistry{
		Type:        capability.String(),
		Description: description,
		IsActive:    true,
	}
}
