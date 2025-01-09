package handlers

import (
	"context"
	"fmt"
	"strings"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/data_fields"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/enum"
	commonEnum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/multierr"
)

const MinHoursBetweenNotifications int = 12 // on same domain for a tenant

func HandleWebsiteVisitorEvent(c context.Context, s *service.Services, sourceEvent commonEnum.FlowListenerEvent, eventData *data_fields.WebsiteVisitEvent) error {
	span, ctx := opentracing.StartSpanFromContext(c, "EventHandlers.HandleWebsiteVisitorEvent")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)
	tracing.LogObjectAsJson(span, "eventData", eventData)

	ctx = common.WithCustomContext(ctx, &common.CustomContext{
		Tenant: eventData.Tenant,
	})

	// List of agents that listen on event

	// find all live agent implementations
	query := entity.Agents{
		Tenant:     eventData.Tenant,
		IsActive:   true,
		TriggersOn: eventData.Type(),
	}

	agents, err := s.PostgresRepositories.AgentsRepository.FindAll(ctx, query)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	if agents == nil {
		return nil
	}

	var errs error
	for _, agent := range *agents {
        // create automation execution record

        // call agent service

        // update automation execution record with results

		default:
			err = fmt.Errorf("automation agent not handled for %s", automation.ID)
			tracing.TraceErr(span, err)
			errs = multierr.Append(errs, err)
		}
	}

	return errs
}

