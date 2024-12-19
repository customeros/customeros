package flows

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/service"
)

func validateFlowBelongsToTenant(c *gin.Context, s *service.Services) (string, bool, error) {
	span, ctx := tracing.StartTracerSpan(c.Request.Context(), "Flows.ValidateFlowBelongsToTenant")
	defer span.Finish()
	tracing.TagComponentRest(span)

	flowId := c.Param("flowId")
	if flowId == "" {
		return "", true, nil
	}

	flowPrefix := strings.HasPrefix(strings.ToLower(flowId), "flow_")
	if !flowPrefix {
		flowId = fmt.Sprintf("flow_%s", flowId)
	}

	isValidFlow, err := s.CommonServices.WorkflowService.ValidateFlowBelongsToTenant(ctx, flowId)
	if err != nil {
		tracing.TraceErr(span, err)
		return flowId, false, err
	}

	if !isValidFlow {
		err := errors.New("Flow does not belong to tenant")
		tracing.TraceErr(span, err)
		return flowId, false, nil
	}

	return flowId, true, nil
}
