package agent

import (
	"context"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/dto"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/opentracing/opentracing-go"
)

type SupportAgent struct {
}

func NewSupportAgent() *SupportAgent {
	return &SupportAgent{}
}

func (a *SupportAgent) Run(ctx context.Context, agentID string, event *dto.SupportEvent) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SupportAgent.Run")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogKV("agentID", agentID)
	tracing.LogObjectAsJson(span, "event", event)

}
