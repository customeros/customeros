package flows

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	commonEnum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/handlers"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/service"
)

type FlowTransitionResponse struct {
	enum.BaseResponse
	Transitions []FlowTransitionRecord `json:"transitions"`
}

type FlowTransitionSingleResponse struct {
	enum.BaseResponse
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

		tenant := handlers.ValidateTenant(c, ctx, span)
		if tenant == "" {
			handlers.SendError(c, span, http.StatusUnauthorized, enum.ErrInvalidAPIKey)
			return
		}

		// Get query params
		from := c.Query("fromNode")
		query := entity.FlowTransitionsRegistry{
			FromNode: from,
		}

		var transitions *[]entity.FlowTransitionsRegistry
		var err error

		if from == "" {
			transitions, err = s.Repositories.PostgresRepositories.FlowTransitionsRegistryRepository.GetAllActive(ctx, nil)
		} else {
			transitions, err = s.Repositories.PostgresRepositories.FlowTransitionsRegistryRepository.GetAllActive(ctx, &query)
		}
		if err != nil {
			tracing.TraceErr(span, err)
			handlers.SendError(c, span, http.StatusInternalServerError, enum.ErrInternalServer)
			return
		}

		results, err := buildTransitionRecords(ctx, span, s, transitions)
		if err != nil {
			tracing.TraceErr(span, err)
			handlers.SendError(c, span, http.StatusInternalServerError, enum.ErrInternalServer)
			return
		}

		if len(results) == 1 {
			c.JSON(http.StatusOK, FlowTransitionSingleResponse{
				BaseResponse: enum.BuildBaseResponse(enum.StatusSuccess),
				Transition:   results[0],
			})
			return
		}

		c.JSON(http.StatusOK, FlowTransitionResponse{
			BaseResponse: enum.BuildBaseResponse(enum.StatusSuccess),
			Transitions:  results,
		})
	}
}

func buildTransitionRecords(ctx context.Context, span opentracing.Span, s *service.Services,
	allFromNodes *[]entity.FlowTransitionsRegistry,
) ([]FlowTransitionRecord, error) {
	results := make([]FlowTransitionRecord, len(*allFromNodes))

	for i, from := range *allFromNodes {
		record := FlowTransitionRecord{
			Step: from.FromNode,
			Type: from.FromNodeType,
		}

		enumType, err := commonEnum.GetFlowNodeType(from.FromNodeType)
		if err != nil {
			tracing.TraceErr(span, err)
			return results, err
		}

		allNextSteps, err := s.Repositories.PostgresRepositories.FlowTransitionsRegistryRepository.GetAllActive(ctx, &entity.FlowTransitionsRegistry{
			FromNodeType: enumType.String(),
		})
		if err != nil {
			tracing.TraceErr(span, err)
			return results, err
		}

		nextSteps := make([]NextStepRecord, len(*allNextSteps))
		for i, next := range *allNextSteps {
			nextSteps[i] = NextStepRecord{
				Step: next.ToNode,
				Type: next.ToNodeType,
			}
		}

		record.AvailableNextSteps = nextSteps
		results[i] = record
	}

	return results, nil
}
