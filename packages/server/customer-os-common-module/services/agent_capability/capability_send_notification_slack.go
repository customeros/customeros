package agent_capability

import (
	"context"
	"fmt"

	"github.com/opentracing/opentracing-go"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/dto"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
)

type SendSlackNotificationInput struct {
	Message                   string
	NotificationCoolDownInHrs int
}

type SendSlackNotificationResult struct{}

func (c *agentCapabilityService) handleSendSlackNotificationExecution(ctx context.Context, executionContainer *dto.CapabilityExecutionContainer) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentCapabilityService.handleSendSlackNotificationExecution")
	defer span.Finish()
	tracing.TagComponentService(span)

	input, ok := executionContainer.InputData.(SendSlackNotificationInput)
	if !ok {
		err := fmt.Errorf("expected SendSlackNotificationInput, got %T", executionContainer.InputData)
		tracing.TraceErr(span, err)
		return err
	}
	result, err := c.executeSendSlackNotification(ctx, input)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	executionContainer.OutputData = result
	span.LogKV(
		"event", "capability_executed",
		"output_type", fmt.Sprintf("%T", result),
	)
	return nil
}

// func (a *agentVisitorIDService) skipNotification(ctx context.Context, agentConfig *AgentConfig, session *postgres_entity.WebSession) (bool, error) {
// 	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentVisitorIDService.skipNotifications")
// 	tracing.SetDefaultServiceSpanTags(ctx, span)
// 	defer span.Finish()
//
// 	if session.Domain == nil {
// 		return true, nil
// 	}
//
// 	// don't send if from workspace domain
// 	isWorkspaceDomain := a.isWorkspaceDomain(ctx, *session.Domain)
// 	if isWorkspaceDomain {
// 		return true, nil
// 	}
//
// 	if agentConfig == nil {
// 		err := errors.New("agent config not set")
// 		tracing.TraceErr(span, err)
// 		return true, err
// 	}
//
// 	// determine last notification from this domain
// 	lastNotification, err := a.postgresRepositories.WebSessionRepository.FindLastNotification(ctx, session.Tenant, *session.Domain)
// 	if err != nil {
// 		tracing.TraceErr(span, err)
// 		return false, nil
// 	}
//
// 	if lastNotification == nil {
// 		return false, nil
// 	}
//
// 	// determine how long since last notification
// 	hoursSinceLastNotification := time.Now().Sub(*lastNotification.SentSlackNotification).Hours()
// 	if hoursSinceLastNotification < float64(agentConfig.NotificationCooldownInHours) {
// 		return true, nil
// 	}
//
// 	return false, nil
// }
//
// func (a *agentVisitorIDService) isWorkspaceDomain(ctx context.Context, domain string) bool {
// 	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentVisitorIDService.isWorkspaceDomain")
// 	tracing.SetDefaultServiceSpanTags(ctx, span)
// 	defer span.Finish()
//
// 	workspaceDomains, err := a.workspaceService.GetWorkspaceDomainsForTenant(ctx)
// 	if err != nil {
// 		tracing.TraceErr(span, err)
// 		return false
// 	}
//
// 	for _, d := range workspaceDomains {
// 		if strings.EqualFold(d, domain) {
// 			return true
// 		}
// 	}
//
// 	return false
// }
