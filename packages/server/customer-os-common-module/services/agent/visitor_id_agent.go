package agent

import (
	"context"
	"errors"
	"fmt"

	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	postgres_repository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"
	"github.com/opentracing/opentracing-go"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/data_fields"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/dto"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/agent_capability"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
)

type AgentVisitorIDService struct {
	postgresRepositories   *postgres_repository.Repositories
	agentService           interfaces.AgentService
	agentCapabilityService interfaces.AgentCapabilityService
}

func NewAgentVisitorIDService(
	postgresRepositories *postgres_repository.Repositories,
	agentService interfaces.AgentService,
	agentCapabilityService interfaces.AgentCapabilityService,
) *AgentVisitorIDService {
	return &AgentVisitorIDService{
		postgresRepositories:   postgresRepositories,
		agentService:           agentService,
		agentCapabilityService: agentCapabilityService,
	}
}

const DefaultNotificationCooldownInHours = 12

func (a *AgentVisitorIDService) Create(ctx context.Context) (*postgres_entity.Agents, error) {
	return a.agentService.CreateAgent(ctx, enum.AgentVisitorID)
}

func (a *AgentVisitorIDService) Run(ctx context.Context, agentID string, event *data_fields.WebsiteVisitEvent) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentVisitorIDService.Run")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	// create execution record
	executionID, err := a.createAgentExecutionRecord(ctx, agentID, event.Type())
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	// lookup capabilities

	// execute identify visitor capability
	visitorIDResults, err := a.executeVisitorIDCapability(ctx, agentID, executionID, event)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	// no result, return early
	if visitorIDResults.Domain == "" {
		return nil
	}

	// execute org creation capability
	orgCreationResults, err := a.executeOrgCreationCapability(ctx, agentID, executionID, visitorIDResults.Domain)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	// execute web session analysis capability
	_, err = a.executeWebSessionAnalysisCapability(
		ctx,
		agentID,
		executionID,
		visitorIDResults.Domain,
		orgCreationResults.OrganizationID,
		event,
	)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	// execute send notification capability
	return nil
}

func (a *AgentVisitorIDService) executeWebSessionAnalysisCapability(
	ctx context.Context,
	agentID, executionID, domain, organizationID string,
	event *data_fields.WebsiteVisitEvent,
) (agent_capability.AnalyzeWebSessionResult, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentVisitorIDService.executeWebSessionAnalysisCapability")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	executionContainer := dto.CapabilityExecutionContainer{
		AgentID:          agentID,
		AgentExecutionID: executionID,
		Capability:       enum.CapabilityAnalyzeWebSessionIntent,
		InputData: agent_capability.AnalyzeWebSessionInput{
			SessionID:      event.SessionID,
			VisitorID:      event.VisitorID,
			OrganizationID: organizationID,
			Domain:         domain,
		},
	}

	err := a.agentCapabilityService.ExecuteCapability(ctx, &executionContainer)
	if err != nil {
		tracing.TraceErr(span, err)
		return agent_capability.AnalyzeWebSessionResult{}, err
	}

	output, ok := executionContainer.OutputData.(agent_capability.AnalyzeWebSessionResult)
	if !ok {
		err := fmt.Errorf("expected agent_capability.AnalyzeWebSessionResult, got %T", executionContainer.OutputData)
		tracing.TraceErr(span, err)
		return agent_capability.AnalyzeWebSessionResult{}, err
	}

	return output, nil
}

func (a *AgentVisitorIDService) executeOrgCreationCapability(ctx context.Context, agentID, executionID, domain string) (agent_capability.CreateOrganizationResult, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentVisitorIDService.executeOrgCreationCapability")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	executionContainer := dto.CapabilityExecutionContainer{
		AgentID:          agentID,
		AgentExecutionID: executionID,
		Capability:       enum.CapabilityCreateOrganization,
		InputData: data_fields.OrganizationFields{
			Domains: []string{domain},
		},
	}

	err := a.agentCapabilityService.ExecuteCapability(ctx, &executionContainer)
	if err != nil {
		tracing.TraceErr(span, err)
		return agent_capability.CreateOrganizationResult{}, err
	}

	output, ok := executionContainer.OutputData.(agent_capability.CreateOrganizationResult)
	if !ok {
		err := fmt.Errorf("expected agent_capability.CreateOrganizationResult, got %T", executionContainer.OutputData)
		tracing.TraceErr(span, err)
		return agent_capability.CreateOrganizationResult{}, err
	}

	return output, nil

}

func (a *AgentVisitorIDService) executeVisitorIDCapability(ctx context.Context, agentID, executionID string, event *data_fields.WebsiteVisitEvent) (agent_capability.IdentifyWebsiteVisitorResult, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentVisitorIDService.executeVisitorIDCapability")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	executionContainer := dto.CapabilityExecutionContainer{
		AgentID:          agentID,
		AgentExecutionID: executionID,
		Capability:       enum.CapabilityIdentifyWebVisitor,
		InputData: agent_capability.IdentifyWebsiteVisitorInput{
			IPAddress: event.IPAddress,
		},
	}

	err := a.agentCapabilityService.ExecuteCapability(ctx, &executionContainer)
	if err != nil {
		tracing.TraceErr(span, err)
		return agent_capability.IdentifyWebsiteVisitorResult{}, err
	}

	output, ok := executionContainer.OutputData.(agent_capability.IdentifyWebsiteVisitorResult)
	if !ok {
		err := fmt.Errorf("expected agent_capability.IdentifyWebsiteVisitorResult, got %T", executionContainer.OutputData)
		tracing.TraceErr(span, err)
		return agent_capability.IdentifyWebsiteVisitorResult{}, err
	}

	return output, nil
}

func (a *AgentVisitorIDService) createAgentExecutionRecord(ctx context.Context, agentID, triggerEvent string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentVisitorIDService.createAgentExecutionRecord")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	agent, err := a.postgresRepositories.AgentsRepository.Find(ctx, postgres_entity.Agents{
		ID: agentID,
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return "", err
	}

	if agent == nil {
		err := errors.New("agent not found")
		tracing.TraceErr(span, err)
		return "", err
	}

	return a.agentService.CreateAgentExecutionRecord(ctx, *agent, triggerEvent)
}
