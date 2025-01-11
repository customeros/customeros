package service

import (
	"context"

	"github.com/opentracing/opentracing-go"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/data_fields"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
)

type AgentICPQualificationService interface {
}

type agentICPQualificationService struct {
	services *Services
}

func NewAgentICPQualificationService(services *Services) AgentICPQualificationService {
	return &agentICPQualificationService{
		services: services,
	}
}

func (a *agentICPQualificationService) ICPAgent(ctx context.Context, eventData *data_fields.OrganizationQualifyEventFields) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentService.ICPAgent")
	tracing.SetDefaultServiceSpanTags(ctx, span)
	defer span.Finish()

	return a.buildICPQualificationReport(ctx, eventData)

}

func (a *agentICPQualificationService) buildICPQualificationReport(ctx context.Context, eventData *data_fields.OrganizationQualifyEventFields) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentService.buildICPQualificationReport")
	tracing.SetDefaultServiceSpanTags(ctx, span)
	defer span.Finish()

	// get ICP definition

	// get all company context

	// build prompt

	// askAI

	// save to timeline

	// build agent execution record

	// write agent execution to db

	// fire action completion event

	return nil
}
