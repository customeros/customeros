package agents

import (
	"context"
	"encoding/json"
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
		}

		response, err := a.buildGetAgentRegistryResponse(ctx, agents)
		if err != nil {
			tracing.TraceErr(span, err)
			message := "Unable to retrieve agents from registry"
			a.responseHandler.HandleError(c, http.StatusInternalServerError, &message)
		}

		a.responseHandler.HandleSuccess(c, response)
		return
	}
}

func (a *AgentHandler) buildGetAgentRegistryResponse(ctx context.Context, agents []postgres_entity.AgentRegistry) (GetAgentRegistryResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentHandler.buildGetAgentRegistryResponse")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	agentRecords := make([]AgentRegistryRecord, len(agents))
	var errs error
	for _, agent := range agents {
		record, err := a.buildAgentRegistryRecord(ctx, agent)
		if err != nil {
			tracing.TraceErr(span, err)
			errs = multierr.Append(errs, err)
		}
		agentRecords = append(agentRecords, record)
	}

	return GetAgentRegistryResponse{
		Agents: agentRecords,
	}, errs
}

func (a *AgentHandler) buildAgentRegistryRecord(ctx context.Context, agent postgres_entity.AgentRegistry) (AgentRegistryRecord, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentHandler.buildAgentRegistryRecord")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	capabilities, err := a.parseAgentCapabilityIDs(ctx, agent.Capabilities)
	if err != nil {
		tracing.TraceErr(span, err)
		return AgentRegistryRecord{}, err
	}

	record := AgentRegistryRecord{
		ID:                  agent.ID,
		Type:                agent.Type,
		Name:                agent.Name,
		Goal:                agent.Goal,
		DefaultCapabilities: capabilities,
	}
	return record, nil
}

func (a *AgentHandler) parseAgentCapabilityIDs(ctx context.Context, capabilitiesStr string) ([]AgentCapability, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentHandler.parseAgentCapabilityIDs")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	var ids []string
	err := json.Unmarshal([]byte(capabilitiesStr), &ids)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	capabilities := make([]AgentCapability, len(ids))
	for i, id := range ids {
		capabilities[i] = AgentCapability{ID: id}
	}
	return capabilities, nil
}
