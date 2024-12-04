package model

import (
	"context"

	commonenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	"github.com/opentracing/opentracing-go"
)

type EventContext struct {
	Context      context.Context
	Span         opentracing.Span
	Services     *service.Services
	Tenant       string
	SourceSystem enum.ExternalSystemId
	SourceEvent  commonenum.FlowEvent
}
