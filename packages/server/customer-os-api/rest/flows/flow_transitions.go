package flows

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/opentracing/opentracing-go"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/rest"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/service"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/repository"
)

type FlowTransitionResponse struct {
	rest.BaseResponse
	Transitions []FlowTransitionRecord `json:"transitions"`
}

type FlowTransitionSingleResponse struct {
	rest.BaseResponse
	Transition FlowTransitionRecord `json:"transition"`
}

type FlowTransitionRecord struct {
	Step               string           `json:"step"`
	Type               string           `json:"type"`
	AvailableNextSteps []NextStepRecord `json:"availableNextSteps"`
}

type NextStepRecord struct {
	Step string `json:"step"`
	Type string `json:"type"`
}

func GetTransitions(s *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "Flows.GetTransitions", c.Request.Header)
		defer span.Finish()
		tracing.TagComponentRest(span)

		// Validate tenant
		tenant := rest.ValidateTenant(c, ctx, span)
		if tenant == "" {
			rest.SendError(c, span, http.StatusUnauthorized, rest.ErrInvalidAPIKey)
			return
		}

		// Get all from nodes
		allFromNodes, err := s.Repositories.PostgresRepositories.FlowTransitionsRegistryRepository.GetUniqueFromNodes(ctx)
		if err != nil {
			tracing.TraceErr(span, err)
			rest.SendError(c, span, http.StatusInternalServerError, rest.ErrInternalServer)
			return
		}

		// Build transition records
		results := buildTransitionRecords(ctx, span, s, allFromNodes)

		// Send response
		if len(results) == 1 {
			c.JSON(http.StatusOK, FlowTransitionSingleResponse{
				BaseResponse: rest.BuildBaseResponse(rest.StatusSuccess),
				Transition:   results[0],
			})
			return
		}

		c.JSON(http.StatusOK, FlowTransitionResponse{
			BaseResponse: rest.BuildBaseResponse(rest.StatusSuccess),
			Transitions:  results,
		})
	}
}

func buildTransitionRecords(ctx context.Context, span opentracing.Span, s *service.Services,
	allFromNodes []repository.FromNodeRecord) []FlowTransitionRecord {

	results := make([]FlowTransitionRecord, len(allFromNodes))

	for i, from := range allFromNodes {
		record := FlowTransitionRecord{
			Step: from.FromNode,
			Type: from.FromNodeType,
		}

		enumType, err := enum.GetFlowNodeType(from.FromNode)
		if err != nil {
			tracing.TraceErr(span, err)
			continue
		}

		allNextSteps, err := s.Repositories.PostgresRepositories.FlowTransitionsRegistryRepository.GetFlowTransitionsByNode(ctx, enumType)
		if err != nil {
			tracing.TraceErr(span, err)
			continue
		}

		nextSteps := make([]NextStepRecord, len(allNextSteps))
		for i, next := range allNextSteps {
			nextSteps[i] = NextStepRecord{
				Step: next.ToNode,
				Type: next.ToNodeType,
			}
		}

		record.AvailableNextSteps = nextSteps
		results[i] = record
	}

	return results
}
