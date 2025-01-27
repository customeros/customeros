package agents

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"net/http"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
)

func (a *AgentHandler) RegisterAgent() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "AgentHandler.RegisterAgent", c.Request.Header)
		defer span.Finish()
		tracing.TagComponentRest(span)

		tenant := common.GetTenantFromContext(ctx)
		if tenant == "" {
			a.responseHandler.HandleError(c, http.StatusNotFound, nil)
			return
		}

		agentPayload, err := a.parseRegisterAgentPayload(c)
		if err != nil {
			tracing.TraceErr(span, err)
			message := err.Error()
			a.responseHandler.HandleError(c, http.StatusBadRequest, &message)
			return
		}

		newAgentID, err := a.createMasterAgent(ctx, agentPayload)
		if err != nil {
			tracing.TraceErr(span, err)
			message := "Unable to create agent in registry"
			a.responseHandler.HandleError(c, http.StatusInternalServerError, &message)
			return
		}

		a.responseHandler.HandleCreated(c, RegisterMasterAgentResponse{
			ID: newAgentID,
		})
		return
	}
}

func (a *AgentHandler) createMasterAgent(ctx context.Context, agent RegisterMasterAgentRequest) (string, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "AgentHandler.createMasterAgent")
	defer span.Finish()
	tracing.TagComponentRest(span)

	var capabilities []postgres_entity.Capability
	for _, inputCapability := range agent.DefaultCapabilities {
		capability := postgres_entity.Capability{
			Type:        enum.AgentCapabilityType(inputCapability.Type),
			Description: inputCapability.Description,
			Active:      inputCapability.Active,
			Name:        inputCapability.Name,
		}
		capabilities = append(capabilities, capability)
	}
	query := postgres_entity.AgentRegistry{
		Type:               enum.AgentType(agent.Type),
		Name:               agent.Name,
		Goal:               agent.Goal,
		IsActive:           true,
		Icon:               agent.Icon,
		CapabilitiesConfig: postgres_entity.CapabilitiesConfig{Capabilities: capabilities},
	}

	newAgent, err := a.services.Repositories.PostgresRepositories.AgentRegistryRepository.Create(ctx, query)
	if err != nil {
		tracing.TraceErr(span, err)
		return "", err
	}
	if newAgent == nil {
		err := errors.New("Could not create agent in registry")
		tracing.TraceErr(span, err)
		return "", err
	}

	return newAgent.ID, nil
}

func (a *AgentHandler) agentCapabilityIDsToString(capabilities []AgentCapability) string {
	ids := make([]string, len(capabilities))
	for i, cap := range capabilities {
		ids[i] = cap.ID
	}
	jsonBytes, _ := json.Marshal(ids)
	return string(jsonBytes)
}

func (a *AgentHandler) parseRegisterAgentPayload(c *gin.Context) (RegisterMasterAgentRequest, error) {
	span, _ := opentracing.StartSpanFromContext(c.Request.Context(), "AgentHandler.parseRegisterAgentPayload")
	defer span.Finish()
	tracing.TagComponentRest(span)

	var req RegisterMasterAgentRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		return req, err
	}

	if req.Type == "" {
		return req, errors.New("Must provide type")
	}

	if req.Name == "" {
		return req, errors.New("Must provide name")
	}

	if req.Goal == "" {
		return req, errors.New("Must provide goal")
	}

	if req.DefaultCapabilities == nil {
		return req, errors.New("Must provide defaultCapabilities")
	}

	return req, nil
}
