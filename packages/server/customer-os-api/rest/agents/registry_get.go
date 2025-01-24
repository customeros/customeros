package agents

import (
	"context"
	"net/http"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/multierr"
)

func (a *AgentHandler) AgentRegistry() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "AgentHandler.AgentRegistry", c.Request.Header)
		defer span.Finish()
		tracing.TagComponentRest(span)

		agents, err := a.services.Repositories.PostgresRepositories.AgentRegistryRepository.FindAll(ctx)
		if err != nil {
			tracing.TraceErr(span, err)
			a.responseHandler.HandleError(c, http.StatusInternalServerError, nil)
			return
		}

		if agents == nil || len(agents) == 0 {
			message := "Agent registry is empty"
			a.responseHandler.HandleError(c, http.StatusNotFound, &message)
			return
		}

		response, err := a.buildGetAgentRegistryResponse(ctx, agents)
		if err != nil {
			tracing.TraceErr(span, err)
			message := "Unable to retrieve agents from registry"
			a.responseHandler.HandleError(c, http.StatusInternalServerError, &message)
			return
		}

		a.responseHandler.HandleSuccess(c, response)
		return
	}
}

func (a *AgentHandler) buildGetAgentRegistryResponse(ctx context.Context, agents []postgres_entity.AgentRegistry) (GetAgentRegistryResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentHandler.buildGetAgentRegistryResponse")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	agentRecords := make([]AgentRegistryRecord, 0, len(agents))
	var errs error

	for _, agent := range agents {
		record, err := a.buildAgentRegistryRecord(ctx, agent)
		if err != nil {
			tracing.TraceErr(span, err)
			errs = multierr.Append(errs, err)
			continue
		}
		agentRecords = append(agentRecords, record)
	}

	return GetAgentRegistryResponse{
		Agents: agentRecords,
	}, errs
}

func (a *AgentHandler) buildAgentRegistryRecord(ctx context.Context, agentRegistry postgres_entity.AgentRegistry) (AgentRegistryRecord, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentHandler.buildAgentRegistryRecord")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	capabilities := make([]AgentCapability, 0, len(agentRegistry.CapabilitiesConfig.Capabilities))
	for _, capability := range agentRegistry.CapabilitiesConfig.Capabilities {
		capabilityRecord := AgentCapability{
			Type:        capability.Type.String(),
			Description: capability.Description,
			Optional:    capability.Optional,
			Name:        capability.Name,
		}
		capabilities = append(capabilities, capabilityRecord)
	}

	record := AgentRegistryRecord{
		ID:                  agentRegistry.ID,
		Type:                agentRegistry.Type.String(),
		Name:                agentRegistry.Name,
		Goal:                agentRegistry.Goal,
		DefaultCapabilities: capabilities,
	}
	return record, nil
}
