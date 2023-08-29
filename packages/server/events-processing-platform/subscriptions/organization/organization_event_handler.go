package organization

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/data"
	commonEntity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/repository/postgres/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/ai"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/caches"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/aggregate"
	cmd "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/command"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/command_handler"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/mapper"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/models"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/graph_db"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/graph_db/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"math"
	"strings"
	"time"
)

const (
	Unknown = "Unknown"
)

type WebscrapeRequest struct {
	Domain string `json:"scrape"`
}
type Socials struct {
	Github    string `json:"github,omitempty"`
	Linkedin  string `json:"linkedin,omitempty"`
	Twitter   string `json:"twitter,omitempty"`
	Youtube   string `json:"youtube,omitempty"`
	Instagram string `json:"instagram,omitempty"`
	Facebook  string `json:"facebook,omitempty"`
}

type WebscrapeResponseV1 struct {
	CompanyName      string `json:"companyName,omitempty"`
	Website          string `json:"website,omitempty"`
	Market           string `json:"market,omitempty"`
	Industry         string `json:"industry,omitempty"`
	IndustryGroup    string `json:"industryGroup,omitempty"`
	SubIndustry      string `json:"subIndustry,omitempty"`
	TargetAudience   string `json:"targetAudience,omitempty"`
	ValueProposition string `json:"valueProposition,omitempty"`
	Github           string `json:"github,omitempty"`
	Linkedin         string `json:"linkedin,omitempty"`
	Twitter          string `json:"twitter,omitempty"`
	Youtube          string `json:"youtube,omitempty"`
	Instagram        string `json:"instagram,omitempty"`
	Facebook         string `json:"facebook,omitempty"`
}

type organizationEventHandler struct {
	repositories         *repository.Repositories
	organizationCommands *command_handler.OrganizationCommands
	log                  logger.Logger
	cfg                  *config.Config
	caches               caches.Cache
	domainScraper        *DomainScraper
}

func (h *organizationEventHandler) WebscrapeOrganization(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationEventHandler.WebscrapeOrganization")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", evt.GetAggregateID()))

	var eventData events.OrganizationLinkDomainEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}
	organizationId := aggregate.GetOrganizationObjectID(evt.AggregateID, eventData.Tenant)
	span.LogFields(log.String("organizationId", organizationId))
	if eventData.Domain == "" {
		tracing.TraceErr(span, errors.New("domain is empty"))
		h.log.Errorf("Missing domain in event data: %v", eventData)
		return nil
	}

	alreadyWebscraped, err := h.repositories.OrganizationRepository.OrganizationWebscrapedForDomain(ctx, eventData.Tenant, organizationId, eventData.Domain)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error checking if organization %s already webscraped: %v", organizationId, err)
		return nil
	}
	if alreadyWebscraped {
		h.log.Infof("Organization %s already webscraped for domain %s", organizationId, eventData.Domain)
		return nil
	}
	result, err := h.domainScraper.Scrape(eventData.Domain, eventData.Tenant, organizationId)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error scraping domain %s: %v", eventData.Domain, err)
		return nil
	}

	org, err := h.repositories.OrganizationRepository.GetOrganization(ctx, eventData.Tenant, organizationId)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error getting organization with id %s: %v", organizationId, err)
		return nil
	}
	currentOrgName := utils.GetStringPropOrEmpty(org.Props, "name")

	err = h.organizationCommands.UpdateOrganization.Handle(ctx,
		cmd.NewUpdateOrganizationCommand(organizationId, eventData.Tenant, constants.SourceWebscrape,
			models.OrganizationDataFields{
				Name:             h.prepareOrgName(result.CompanyName, currentOrgName),
				Market:           result.Market,
				Industry:         result.Industry,
				IndustryGroup:    result.IndustryGroup,
				SubIndustry:      result.SubIndustry,
				TargetAudience:   result.TargetAudience,
				ValueProposition: result.ValueProposition,
			},
			utils.TimePtr(utils.Now()), true))
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error updating organization: %v", err)
		return nil
	}
	if result.Youtube != "" {
		h.addSocial(ctx, organizationId, eventData.Tenant, "youtube", result.Youtube)
	}
	if result.Twitter != "" {
		h.addSocial(ctx, organizationId, eventData.Tenant, "twitter", result.Twitter)
	}
	if result.Linkedin != "" {
		h.addSocial(ctx, organizationId, eventData.Tenant, "linkedin", result.Linkedin)
	}
	if result.Github != "" {
		h.addSocial(ctx, organizationId, eventData.Tenant, "github", result.Github)
	}
	if result.Instagram != "" {
		h.addSocial(ctx, organizationId, eventData.Tenant, "instagram", result.Instagram)
	}
	if result.Facebook != "" {
		h.addSocial(ctx, organizationId, eventData.Tenant, "facebook", result.Facebook)
	}

	return nil
}

func (h *organizationEventHandler) addSocial(ctx context.Context, organizationId, tenant, platform, url string) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationEventHandler.addSocial")
	defer span.Finish()
	span.LogFields(log.String("organizationId", organizationId), log.String("tenant", tenant), log.String("platform", platform), log.String("url", url))

	err := h.organizationCommands.AddSocialCommand.Handle(ctx,
		command_handler.NewAddSocialCommand(organizationId, tenant, uuid.New().String(),
			platform, url, constants.SourceWebscrape, constants.SourceWebscrape, constants.AppSourceEventProcessingPlatform, nil, nil))
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error adding %s social: %v", platform, err)
	}
}

func (h *organizationEventHandler) prepareOrgName(webscrabedOrgName, currentOrgName string) string {
	if currentOrgName == "" {
		return webscrabedOrgName
	} else {
		return currentOrgName
	}
}

func (h *organizationEventHandler) AdjustNewOrganizationFields(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationEventHandler.AdjustNewOrganizationFields")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", evt.GetAggregateID()))

	var eventData events.OrganizationCreateEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}
	organizationId := aggregate.GetOrganizationObjectID(evt.AggregateID, eventData.Tenant)

	market := h.mapMarketValue(eventData.Market)
	industry := h.mapIndustryToGICS(ctx, eventData.Tenant, organizationId, eventData.Industry)

	if eventData.Market != market || eventData.Industry != industry {
		err := h.callUpdateOrganizationCommand(ctx, eventData.Tenant, organizationId, eventData.SourceOfTruth, market, industry, span)
		return err
	} else {
		h.log.Infof("No need to update organization %s", organizationId)
	}
	return nil
}

func (h *organizationEventHandler) AdjustUpdatedOrganizationFields(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationEventHandler.AdjustUpdatedOrganizationFields")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", evt.GetAggregateID()))

	var eventData events.OrganizationUpdateEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}
	organizationId := aggregate.GetOrganizationObjectID(evt.AggregateID, eventData.Tenant)

	market := h.mapMarketValue(eventData.Market)
	industry := h.mapIndustryToGICS(ctx, eventData.Tenant, organizationId, eventData.Industry)

	if eventData.Market != market || eventData.Industry != industry {
		err := h.callUpdateOrganizationCommand(ctx, eventData.Tenant, organizationId, eventData.SourceOfTruth, market, industry, span)
		return err
	} else {
		h.log.Infof("No need to update organization %s", organizationId)
	}
	return nil
}

func (h *organizationEventHandler) callUpdateOrganizationCommand(ctx context.Context, tenant, organizationId, sourceOfTruth, market, industry string, span opentracing.Span) error {
	err := h.organizationCommands.UpdateOrganization.Handle(ctx,
		cmd.NewUpdateOrganizationCommand(organizationId, tenant, sourceOfTruth,
			models.OrganizationDataFields{
				Market:   market,
				Industry: industry,
			},
			utils.TimePtr(utils.Now()),
			true))
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error updating organization %s: %v", organizationId, err.Error())
		return err
	}
	return nil
}

func (h *organizationEventHandler) mapMarketValue(inputMarket string) string {
	return data.AdjustOrganizationMarket(inputMarket)
}

func (h *organizationEventHandler) mapIndustryToGICS(ctx context.Context, tenant, orgId, inputIndustry string) string {
	trimmedInputIndustry := strings.TrimSpace(inputIndustry)

	if inputIndustry == "" {
		return ""
	}

	var industry string
	if industryValue, ok := h.caches.GetIndustry(trimmedInputIndustry); ok {
		industry = industryValue
	} else {
		h.log.Infof("Industry %s not found in cache, asking AI", trimmedInputIndustry)
		industry = h.mapIndustryToGICSWithAI(ctx, tenant, orgId, trimmedInputIndustry)
		if industry != "" && len(industry) < 45 {
			h.caches.SetIndustry(trimmedInputIndustry, industry)
			h.log.Infof("Industry %s mapped to %s", trimmedInputIndustry, industry)
		} else {
			h.log.Warnf("Industry %s mapped wrongly to (%s) with AI, returning input value", industry, trimmedInputIndustry)
			return trimmedInputIndustry
		}
	}
	if industry == Unknown {
		h.log.Infof("Unknown industry %s, returning as is", trimmedInputIndustry)
		return trimmedInputIndustry
	}

	return strings.TrimSpace(industry)
}

func (h *organizationEventHandler) mapIndustryToGICSWithAI(ctx context.Context, tenant, orgId, inputIndustry string) string {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationEventHandler.mapIndustryToGICSWithAI")
	defer span.Finish()
	span.LogFields(log.String("inputIndustry", inputIndustry))

	firstPrompt := fmt.Sprintf(h.cfg.Services.Anthropic.IndustryLookupPrompt1, inputIndustry)

	promptLog1 := commonEntity.AiPromptLog{
		CreatedAt:      utils.Now(),
		AppSource:      constants.AppSourceEventProcessingPlatform,
		Provider:       constants.Anthropic,
		Model:          "claude-2",
		PromptType:     constants.PromptType_MapIndustry,
		Tenant:         &tenant,
		NodeId:         &orgId,
		NodeLabel:      utils.StringPtr(constants.NodeLabel_Organization),
		PromptTemplate: &h.cfg.Services.Anthropic.IndustryLookupPrompt1,
		Prompt:         firstPrompt,
	}
	promptStoreLogId1, err := h.repositories.CommonRepositories.AiPromptLogRepository.Store(promptLog1)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error storing prompt log: %v", err)
	} else {
		span.LogFields(log.String("promptStoreLogId1", promptStoreLogId1))
	}

	firstResult, err := ai.InvokeAnthropic(ctx, h.cfg, h.log, firstPrompt)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error invoking AI: %v", err)
		storeErr := h.repositories.CommonRepositories.AiPromptLogRepository.UpdateError(promptStoreLogId1, err.Error())
		if storeErr != nil {
			tracing.TraceErr(span, storeErr)
			h.log.Errorf("Error updating prompt log with error: %v", storeErr)
		}
		return ""
	} else {
		storeErr := h.repositories.CommonRepositories.AiPromptLogRepository.UpdateResponse(promptStoreLogId1, firstResult)
		if storeErr != nil {
			tracing.TraceErr(span, storeErr)
			h.log.Errorf("Error updating prompt log with ai response: %v", storeErr)
		}
	}
	if firstResult == "" || firstResult == Unknown {
		return firstResult
	}
	secondPrompt := fmt.Sprintf(h.cfg.Services.Anthropic.IndustryLookupPrompt2, firstResult)

	promptLog2 := commonEntity.AiPromptLog{
		CreatedAt:      utils.Now(),
		AppSource:      constants.AppSourceEventProcessingPlatform,
		Provider:       constants.Anthropic,
		Model:          "claude-2",
		PromptType:     constants.PromptType_ExtractIndustryValue,
		Tenant:         &tenant,
		NodeId:         &orgId,
		NodeLabel:      utils.StringPtr(constants.NodeLabel_Organization),
		PromptTemplate: &h.cfg.Services.Anthropic.IndustryLookupPrompt2,
		Prompt:         secondPrompt,
	}
	promptStoreLogId2, err := h.repositories.CommonRepositories.AiPromptLogRepository.Store(promptLog2)

	secondResult, err := ai.InvokeAnthropic(ctx, h.cfg, h.log, secondPrompt)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error invoking AI: %v", err)
		err = h.repositories.CommonRepositories.AiPromptLogRepository.UpdateError(promptStoreLogId2, err.Error())
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("Error updating prompt log with error: %v", err)
		}
		return ""
	} else {
		err = h.repositories.CommonRepositories.AiPromptLogRepository.UpdateResponse(promptStoreLogId2, secondResult)
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("Error updating prompt log with ai response: %v", err)
		}
	}
	return secondResult
}

func (h *organizationEventHandler) OnRenewalForecastRequested(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationEventHandler.OnRenewalForecastRequested")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", evt.GetAggregateID()))

	var eventData events.OrganizationRequestRenewalForecastEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "GetJsonData")
	}
	organizationId := aggregate.GetOrganizationObjectID(evt.AggregateID, eventData.Tenant)

	orgDbNode, err := h.repositories.OrganizationRepository.GetOrganization(ctx, eventData.Tenant, organizationId)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error getting organization with id %s: %v", organizationId, err)
		return nil
	}
	organizationEntity := graph_db.MapDbNodeToOrganizationEntity(*orgDbNode)

	amount, err := h.calculateForecastAmount(ctx, organizationEntity.BillingDetails, organizationEntity.RenewalLikelihood.RenewalLikelihood)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "calculateForecastAmount")
	}
	potentialAmount, err := h.calculateForecastAmount(ctx, organizationEntity.BillingDetails, string(entity.RenewalLikelihoodHigh))
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "calculateForecastAmount")
	}
	if amount != nil && potentialAmount != nil {
		err := h.organizationCommands.UpdateRenewalForecastCommand.Handle(ctx, cmd.NewUpdateRenewalForecastCommand(
			eventData.Tenant, organizationId, models.RenewalForecastFields{
				Amount:          amount,
				PotentialAmount: potentialAmount,
				UpdatedBy:       "",
				UpdatedAt:       utils.Now(),
				Comment:         utils.StringPtr(""),
			},
			mapper.MapRenewalLikelihoodFromGraphDb(entity.RenewalLikelihoodProbability(organizationEntity.RenewalLikelihood.RenewalLikelihood))))
		if err != nil {
			tracing.TraceErr(span, err)
			return errors.Wrap(err, "UpdateRenewalForecastCommand")
		}
		span.LogFields(log.Float64("output - calculated amount", *amount), log.Float64("output - calculated potentialAmount", *potentialAmount))
	} else {
		span.LogFields(log.String("output", "amount or potentialAmount is nil"))
	}
	return nil
}

func (h *organizationEventHandler) calculateForecastAmount(ctx context.Context, billingDtls entity.BillingDetails, likelihood string) (*float64, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationEventHandler.calculateForecastAmount")
	defer span.Finish()

	if billingDtls.Amount == nil || billingDtls.Frequency == "" || billingDtls.RenewalCycle == "" || likelihood == "" {
		return nil, nil
	}

	billingPeriods := h.getBillingPeriods(billingDtls.Frequency, billingDtls.RenewalCycle)

	var likelihoodFactor float64
	switch entity.RenewalLikelihoodProbability(likelihood) {
	case entity.RenewalLikelihoodHigh:
		likelihoodFactor = 1
	case entity.RenewalLikelihoodMedium:
		likelihoodFactor = 0.5
	case entity.RenewalLikelihoodLow:
		likelihoodFactor = 0.25
	case entity.RenewalLikelihoodZero:
		likelihoodFactor = 0
	default:
		return nil, errors.New("invalid likelihood")
	}

	forecastAmount := *billingDtls.Amount * billingPeriods * likelihoodFactor

	// trim decimal places
	forecastAmount = math.Trunc(forecastAmount*100) / 100

	span.LogFields(log.Float64("return - forecastAmount", forecastAmount))
	return &forecastAmount, nil
}

func (h *organizationEventHandler) getBillingPeriods(billingFreq string, renewalFreq string) float64 {
	switch billingFreq {

	case "WEEKLY":
		switch renewalFreq {
		case "WEEKLY":
			return 1
		case "BIWEEKLY":
			return 2
		case "MONTHLY":
			return 4
		case "QUARTERLY":
			return 13
		case "BIANNUALLY":
			return 26
		case "ANNUALLY":
			return 52
		}

	case "BIWEEKLY":
		switch renewalFreq {
		case "BIWEEKLY":
			return 1
		case "MONTHLY":
			return 2
		case "QUARTERLY":
			return 6
		case "BIANNUALLY":
			return 12
		case "ANNUALLY":
			return 24
		}

	case "MONTHLY":
		switch renewalFreq {
		case "MONTHLY":
			return 1
		case "QUARTERLY":
			return 3
		case "BIANNUALLY":
			return 6
		case "ANNUALLY":
			return 12
		}

	case "QUARTERLY":
		switch renewalFreq {
		case "QUARTERLY":
			return 1
		case "BIANNUALLY":
			return 2
		case "ANNUALLY":
			return 4
		}

	case "BIANNUALLY":
		switch renewalFreq {
		case "BIANNUALLY":
			return 1
		case "ANNUALLY":
			return 2
		}

	case "ANNUALLY":
		switch renewalFreq {
		case "ANNUALLY":
			return 1
		}

	default:
		return 1
	}

	return 1
}

func (h *organizationEventHandler) OnNextCycleDateRequested(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationEventHandler.OnNextCycleDateRequested")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", evt.GetAggregateID()))

	var eventData events.OrganizationRequestNextCycleDateEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "GetJsonData")
	}
	organizationId := aggregate.GetOrganizationObjectID(evt.AggregateID, eventData.Tenant)

	orgDbNode, err := h.repositories.OrganizationRepository.GetOrganization(ctx, eventData.Tenant, organizationId)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error getting organization with id %s: %v", organizationId, err.Error())
		return nil
	}
	organizationEntity := graph_db.MapDbNodeToOrganizationEntity(*orgDbNode)

	nextRenewalDate := h.calculateRenewalCycleNext(ctx, organizationEntity.BillingDetails)

	err = h.organizationCommands.UpdateBillingDetailsCommand.Handle(ctx, cmd.NewUpdateBillingDetailsCommand(
		eventData.Tenant, organizationId, models.BillingDetailsFields{
			Amount:            organizationEntity.BillingDetails.Amount,
			Frequency:         organizationEntity.BillingDetails.Frequency,
			RenewalCycle:      organizationEntity.BillingDetails.RenewalCycle,
			RenewalCycleStart: organizationEntity.BillingDetails.RenewalCycleStart,
			RenewalCycleNext:  nextRenewalDate,
			UpdatedBy:         "",
		}))
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "UpdateBillingDetailsCommand")
	}
	return nil
}

func (h *organizationEventHandler) calculateRenewalCycleNext(ctx context.Context, billingDtls entity.BillingDetails) *time.Time {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationEventHandler.calculateRenewalCycleNext")
	defer span.Finish()

	if billingDtls.RenewalCycleStart == nil || billingDtls.RenewalCycle == "" {
		return nil
	}

	renewalCycleNext := *billingDtls.RenewalCycleStart
	for {
		switch billingDtls.RenewalCycle {
		case "WEEKLY":
			renewalCycleNext = renewalCycleNext.AddDate(0, 0, 7)
		case "BIWEEKLY":
			renewalCycleNext = renewalCycleNext.AddDate(0, 0, 14)
		case "MONTHLY":
			renewalCycleNext = renewalCycleNext.AddDate(0, 1, 0)
		case "QUARTERLY":
			renewalCycleNext = renewalCycleNext.AddDate(0, 3, 0)
		case "BIANNUALLY":
			renewalCycleNext = renewalCycleNext.AddDate(0, 6, 0)
		case "ANNUALLY":
			renewalCycleNext = renewalCycleNext.AddDate(1, 0, 0)
		default:
			return nil // invalid
		}

		if renewalCycleNext.After(time.Now()) {
			break
		}
	}

	span.LogFields(log.Object("return - renewalCycleNext", renewalCycleNext))
	return &renewalCycleNext
}
