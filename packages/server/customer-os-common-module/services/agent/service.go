package agent

import (
	"context"
	"errors"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	postgres_repository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"
	"github.com/opentracing/opentracing-go"
)

type agentService struct {
	postgresRepositories *postgres_repository.Repositories
}

func NewAgentService(postgresRepositories *postgres_repository.Repositories) interfaces.AgentService {
	return &agentService{
		postgresRepositories: postgresRepositories,
	}
}

func (a *agentService) CreateAgent(ctx context.Context, agentType enum.AgentType) (*postgres_entity.Agents, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentService.CreateAgent")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	tenant := common.GetTenantFromContext(ctx)
	if tenant == "" {
		err := errors.New("tenant not found in context")
		tracing.TraceErr(span, err)
		return nil, err
	}

	// get config from registry
	masterAgent, err := a.postgresRepositories.AgentRegistryRepository.Find(ctx, agentType)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	if masterAgent == nil {
		err := errors.New("no agent found in agent repository")
		tracing.TraceErr(span, err)
		return nil, err
	}

	// build new agent
	agent := postgres_entity.Agents{
		Type:         masterAgent.ID,
		Tenant:       tenant,
		Name:         masterAgent.Name,
		Capabilities: masterAgent.Capabilities,
		Goal:         masterAgent.Goal,
		IsActive:     false,
		VisibleInUI:  true,
		Icon:         masterAgent.Icon,
		Color:        utils.GetRandomColor(),
	}

	// create agent instance in database
	newAgent, err := a.postgresRepositories.AgentsRepository.Create(ctx, agent)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return newAgent, nil
}

func (a *agentService) CreateAgentExecutionRecord(ctx context.Context, agent postgres_entity.Agents, triggerEvent string) (string, error) {
	span, ctx := tracing.StartTracerSpan(ctx, "AgentService.CreateAgentExecutionRecord")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	agentExecutionRecord := postgres_entity.AgentExecution{
		AgentID:      &agent.ID,
		TriggerEvent: triggerEvent,
		Status:       enum.AgentExecutionRunning.String(),
		StartedAt:    utils.NowPtr(),
	}

	createdRecord, err := a.postgresRepositories.AgentExecutionRepository.Create(ctx, agentExecutionRecord)
	if err != nil {
		tracing.TraceErr(span, err)
		return "", err
	}

	if createdRecord == nil {
		err := errors.New("unable to create agent execution record")
		tracing.TraceErr(span, err)
		return "", err
	}

	return createdRecord.ID, nil

}

func (a *agentService) SaveAgentExecutionSuccess(ctx context.Context, executionID string) error {
	span, ctx := tracing.StartTracerSpan(ctx, "AgentService.SaveAgentExecutionSuccess")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	executionRecord, err := a.postgresRepositories.AgentExecutionRepository.Find(ctx, postgres_entity.AgentExecution{
		ID: executionID,
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	executionRecord.Status = enum.AgentExecutionSuccess.String()
	executionRecord.CompletedAt = utils.NowPtr()
	// update

	return nil
}

func (a *agentService) SaveAgentExecutionError(ctx context.Context, executionID, errorMessage string) error {
	span, ctx := tracing.StartTracerSpan(ctx, "AgentService.SaveAgentExecutionError")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	executionRecord, err := a.postgresRepositories.AgentExecutionRepository.Find(ctx, postgres_entity.AgentExecution{
		ID: executionID,
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	executionRecord.Status = enum.AgentExecutionFail.String()
	executionRecord.ErrorMessage = &errorMessage
	//update

	return nil
}

// func (a *agentVisitorIDService) RunAgent(ctx context.Context, agent *postgres_entity.Agents, eventData any) error {
// 	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentVisitorIDService.RunAgent")
// 	defer span.Finish()
// 	tracing.SetDefaultServiceSpanTags(ctx, span)
//
// 	// validate type
// 	event, ok := eventData.(data_fields.WebsiteVisitEvent)
// 	if !ok {
// 		return errors.New("invalid trigger event")
// 	}
//
// 	// validate agent
// 	if agent.Type != enum.AgentVisitorID.String() {
// 		return errors.New("agent does not match VisitorID agent")
// 	}
//
// 	tenant := common.GetTenantFromContext(ctx)
// 	if tenant == "" {
// 		return errors.New("tenant not set on context")
// 	}
//
// 	if tenant != agent.Tenant {
// 		return errors.New("agent does not belong to tenant")
// 	}
//
// 	// try to deanonymize the IP address
//
// 	// Create org if doesn't exist
// 	orgID, err := a.organizationService.Save(ctx, nil, nil, data_fields.OrganizationFields{
// 		Domains: []string{*domain},
// 	})
// 	if err != nil {
// 		tracing.TraceErr(span, err)
// 		return err
// 	}
//
// 	// update session table with ID data
// 	session, err := a.postgresRepositories.WebSessionRepository.UpdateSessionWithDomain(ctx, event.SessionID, *domain)
// 	if err != nil {
// 		tracing.TraceErr(span, err)
// 		return err
// 	}
//
// 	// determine if new org
// 	isNewOrg, err := a.isNewCompanyVisit(ctx, &tenant, domain)
// 	if err != nil {
// 		tracing.TraceErr(span, errors.Wrap(err, "failed isNewCompany lookup"))
// 		return err
// 	}
//
// 	// determine if new person
// 	isNewVisitor, err := a.isNewWebsiteVisitor(ctx, tenant, event.VisitorID)
// 	if err != nil {
// 		tracing.TraceErr(span, errors.Wrap(err, "failed isNewVisitor lookup"))
// 		return err
// 	}
//
// 	// log visit on timeline
// 	timelineMessage, err := a.buildTimelineMessage(ctx, &event, isNewOrg, isNewVisitor)
// 	if err != nil {
// 		tracing.TraceErr(span, err)
// 		return err
// 	}
// 	if timelineMessage == nil {
// 		err := errors.New("unable to build timeline message")
// 		tracing.TraceErr(span, err)
// 		return err
// 	}
// 	source := enum.SourceAgent.String()
// 	actionType := enum.ActionGeneric
//
// 	action := data_fields.ActionFields{
// 		Source:     &source,
// 		CreatedAt:  utils.NowPtr(),
// 		ActionType: &actionType,
// 		Content:    timelineMessage,
// 	}
// 	_, err = a.actionService.CreateActionForOrganization(ctx, nil, orgID, action)
// 	if err != nil {
// 		tracing.TraceErr(span, err)
// 		return err
// 	}
//
// 	// determine if slack notification is configured
// 	slackEnabled, agentConfig, err := a.isSlackNotificationEnabled(ctx, agent)
// 	if err != nil {
// 		tracing.TraceErr(span, err)
// 	}
// 	if slackEnabled == false || agentConfig.SlackChannelID == "" {
// 		return nil
// 	}
//
// 	// determine if notification should be skipped
// 	skip, err := a.skipNotification(ctx, agentConfig, session)
// 	if err != nil {
// 		tracing.TraceErr(span, err)
// 	}
// 	if skip {
// 		return nil
// 	}
//
// 	// handle slack notification
// 	message, err := a.buildWebVisitorSlackNotification(ctx, session, orgID)
// 	if err != nil {
// 		tracing.TraceErr(span, err)
// 	}
//
// 	err = a.notificationService.NotifySlackChannel(ctx, tenant, agentConfig.SlackChannelID, message)
// 	if err != nil {
// 		tracing.TraceErr(span, err)
// 	}
//
// 	return nil
// }
