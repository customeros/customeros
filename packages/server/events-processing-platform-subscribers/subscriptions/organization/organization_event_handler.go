package organization

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/clients/grpc_client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/data"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/data_fields"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/enum"
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	postgresEntity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"

	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/caches"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/service"
)

const (
	Unknown = "Unknown"
)

type Socials struct {
	Github    string `json:"github,omitempty"`
	Linkedin  string `json:"linkedin,omitempty"`
	Twitter   string `json:"twitter,omitempty"`
	Youtube   string `json:"youtube,omitempty"`
	Instagram string `json:"instagram,omitempty"`
	Facebook  string `json:"facebook,omitempty"`
}

type organizationEventHandler struct {
	log         logger.Logger
	cfg         *config.Config
	caches      caches.Cache
	grpcClients *grpc_client.Clients
	services    *service.Services
}

func NewOrganizationEventHandler(services *service.Services, log logger.Logger, cfg *config.Config, caches caches.Cache, grpcClients *grpc_client.Clients) *organizationEventHandler {
	return &organizationEventHandler{
		log:         log,
		cfg:         cfg,
		caches:      caches,
		grpcClients: grpcClients,
		services:    services,
	}
}

func (h *organizationEventHandler) saveOrganizationIndustryAndMarket(ctx context.Context, tenant, organizationId, market, industry string, updateMarket, updateIndustry bool) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationEventHandler.saveOrganizationIndustryAndMarket")
	defer span.Finish()

	if !updateMarket && !updateIndustry {
		h.log.Infof("No need to update organization %s", organizationId)
		return nil
	}

	// delay to avoid updating organization before main event
	time.Sleep(250 * time.Millisecond)

	_, err := h.services.CommonServices.OrganizationService.Save(ctx, nil, &organizationId, data_fields.OrganizationFields{
		Market:   utils.StringPtr(market),
		Industry: utils.StringPtr(industry),
	})
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error updating organization %s: %s", organizationId, err.Error())
		return err
	}
	return nil
}

func (h *organizationEventHandler) mapMarketValue(inputMarket string) string {
	return data.AdjustOrganizationMarket(inputMarket)
}

func (h *organizationEventHandler) mapIndustryToGICS(ctx context.Context, tenant, orgId, inputIndustry string) string {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationEventHandler.mapIndustryToGICS")
	defer span.Finish()
	span.LogFields(log.String("inputIndustry", inputIndustry))
	span.SetTag(tracing.SpanTagTenant, tenant)
	span.SetTag(tracing.SpanTagEntityId, orgId)

	trimmedInputIndustry := strings.TrimSpace(inputIndustry)

	if trimmedInputIndustry == "" {
		return ""
	}

	industry := trimmedInputIndustry
	if industryValue, ok := h.caches.GetIndustry(trimmedInputIndustry); ok {
		span.LogFields(log.Bool("result.industryFoundInCache", true))
		span.LogFields(log.String("result.cacheMapping", industryValue))
		industry = industryValue
	} else {
		span.LogFields(log.Bool("result.industryFoundInCache", false))
		h.log.Infof("Industry %s not found in cache, asking AI", trimmedInputIndustry)
		industry = h.mapIndustryToGICSWithAI(ctx, tenant, orgId, trimmedInputIndustry)
		if utils.Contains(data.GICSIndustryValues, industry) {
			h.caches.SetIndustry(trimmedInputIndustry, industry)
			span.LogFields(log.String("result.newMapping", industry))
		} else {
			industry = trimmedInputIndustry
		}
	}

	return strings.TrimSpace(industry)
}

func (h *organizationEventHandler) mapIndustryToGICSWithAI(ctx context.Context, tenant, orgId, inputIndustry string) string {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationEventHandler.mapIndustryToGICSWithAI")
	defer span.Finish()
	tracing.TagTenant(span, tenant)
	span.LogFields(log.String("inputIndustry", inputIndustry))

	firstPrompt := fmt.Sprintf(h.cfg.Services.Anthropic.IndustryLookupPrompt1, inputIndustry)

	promptLog1 := postgresEntity.AiPromptLog{
		CreatedAt:      utils.Now(),
		AppSource:      constants.AppSourceEventProcessingPlatformSubscribers,
		Provider:       constants.Anthropic,
		Model:          enum.AIModelAnthropicHaiku.String(),
		PromptType:     constants.PromptType_MapIndustry,
		Tenant:         &tenant,
		NodeId:         &orgId,
		NodeLabel:      utils.StringPtr(commonmodel.NodeLabelOrganization),
		PromptTemplate: &h.cfg.Services.Anthropic.IndustryLookupPrompt1,
		Prompt:         firstPrompt,
	}
	promptStoreLogId1, err := h.services.CommonServices.PostgresRepositories.AiPromptLogRepository.Store(promptLog1)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to store prompt log"))
		h.log.Errorf("Error storing prompt log: %v", err)
	} else {
		span.LogFields(log.String("promptStoreLogId1", promptStoreLogId1))
	}

	firstResult, err := h.services.CommonServices.AIService.AskAI(ctx, enum.AIModelAnthropicHaiku, &firstPrompt)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to invoke AI for first prompt"))
		h.log.Errorf("Error invoking AI: %v", err)
		storeErr := h.services.CommonServices.PostgresRepositories.AiPromptLogRepository.UpdateError(promptStoreLogId1, err.Error())
		if storeErr != nil {
			tracing.TraceErr(span, errors.Wrap(storeErr, "failed to update prompt log with error"))
			h.log.Errorf("Error updating prompt log with error: %v", storeErr)
		}
		return ""
	} else {
		storeErr := h.services.CommonServices.PostgresRepositories.AiPromptLogRepository.UpdateResponse(promptStoreLogId1, *firstResult)
		if storeErr != nil {
			tracing.TraceErr(span, errors.Wrap(storeErr, "failed to update prompt log with ai response"))
			h.log.Errorf("Error updating prompt log with ai response: %v", storeErr)
		}
	}
	if firstResult == nil {
		return ""
	}
	secondPrompt := fmt.Sprintf(h.cfg.Services.Anthropic.IndustryLookupPrompt2, firstResult)

	promptLog2 := postgresEntity.AiPromptLog{
		CreatedAt:      utils.Now(),
		AppSource:      constants.AppSourceEventProcessingPlatformSubscribers,
		Provider:       constants.Anthropic,
		Model:          enum.AIModelAnthropicHaiku.String(),
		PromptType:     constants.PromptType_ExtractIndustryValue,
		Tenant:         &tenant,
		NodeId:         &orgId,
		NodeLabel:      utils.StringPtr(commonmodel.NodeLabelOrganization),
		PromptTemplate: &h.cfg.Services.Anthropic.IndustryLookupPrompt2,
		Prompt:         secondPrompt,
	}
	promptStoreLogId2, err := h.services.CommonServices.PostgresRepositories.AiPromptLogRepository.Store(promptLog2)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to store prompt log"))
		h.log.Errorf("Error storing prompt log with error: %v", err)
	}
	secondResult, err := h.services.CommonServices.AIService.AskAI(ctx, enum.AIModelAnthropicHaiku, &secondPrompt)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to invoke AI for second prompt"))
		h.log.Errorf("Error invoking AI: %v", err)
		err = h.services.CommonServices.PostgresRepositories.AiPromptLogRepository.UpdateError(promptStoreLogId2, err.Error())
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to update prompt log with error"))
			h.log.Errorf("Error updating prompt log with error: %v", err)
		}
		return ""
	}
	if secondResult == nil {
		return ""
	}

	err = h.services.CommonServices.PostgresRepositories.AiPromptLogRepository.UpdateResponse(promptStoreLogId2, *secondResult)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to update prompt log with ai response"))
		h.log.Errorf("Error updating prompt log with ai response: %v", err)
	}
	return *secondResult
}

func (h *organizationEventHandler) OnAdjustIndustry(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationEventHandler.OnAdjustIndustry")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", evt.GetAggregateID()))

	var eventData events.OrganizationAdjustIndustryEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "evt.GetJsonData"))
		return errors.Wrap(err, "evt.GetJsonData")
	}
	organizationId := aggregate.GetOrganizationObjectID(evt.AggregateID, eventData.Tenant)
	span.SetTag(tracing.SpanTagEntityId, organizationId)
	span.SetTag(tracing.SpanTagTenant, eventData.Tenant)

	innerCtx := common.WithCustomContext(ctx, &common.CustomContext{
		Tenant:    eventData.Tenant,
		AppSource: constants.AppSourceEventProcessingPlatformSubscribers,
	})

	orgDbNode, err := h.services.CommonServices.Neo4jRepositories.OrganizationReadRepository.GetOrganization(innerCtx, eventData.Tenant, organizationId)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error getting organization with id %s: %v", organizationId, err)
		return err
	}
	organizationEntity := neo4jmapper.MapDbNodeToOrganizationEntity(orgDbNode)

	industry := h.mapIndustryToGICS(innerCtx, eventData.Tenant, organizationId, organizationEntity.Industry)

	if industry != "" && organizationEntity.Industry != industry {
		_, err = h.services.CommonServices.OrganizationService.Save(innerCtx, nil, &organizationId, data_fields.OrganizationFields{
			Industry: utils.StringPtr(industry),
		})
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("Error updating organization %s: %s", organizationId, err.Error())
			return err
		}
	}
	return nil
}
