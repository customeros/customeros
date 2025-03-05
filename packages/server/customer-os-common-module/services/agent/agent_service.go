package agent

import (
	"context"
	"github.com/dustin/go-humanize"
	"strconv"
	"strings"

	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	postgresentity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	postgresrepository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/data_fields"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/dto"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/model"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/agent_capability"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/agent_listeners"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/events"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
)

type agentService struct {
	postgresRepositories            *postgresrepository.Repositories
	events                          *events.EventsService
	agentCapabilities               *agent_capability.AgentCapabilities
	agentListeners                  *agent_listeners.AgentListeners
	agentCapabilityExecutionService interfaces.AgentCapabilityExecutionService
	tenantSettingsService           interfaces.TenantSettingsService
	invoiceService                  interfaces.InvoiceService
	currencyService                 interfaces.CurrencyService
}

func (a *agentService) SetListeners(listeners any) {
	a.agentListeners = listeners.(*agent_listeners.AgentListeners)
}

func NewAgentService(
	postgresRepositories *postgresrepository.Repositories,
	events *events.EventsService,
	agentCapabilities *agent_capability.AgentCapabilities,
	tenantSettingsService interfaces.TenantSettingsService,
	invoiceService interfaces.InvoiceService,
	currencyService interfaces.CurrencyService) interfaces.AgentService {
	return &agentService{
		postgresRepositories:            postgresRepositories,
		events:                          events,
		agentCapabilities:               agentCapabilities,
		agentCapabilityExecutionService: agent_capability.NewAgentCapabilityExecutionService(),
		tenantSettingsService:           tenantSettingsService,
		invoiceService:                  invoiceService,
		currencyService:                 currencyService,
	}
}

func (a *agentService) GetAgentById(ctx context.Context, agentID string) (*postgresentity.Agent, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentService.GetAgentById")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	agent, err := a.postgresRepositories.AgentRepository.GetById(ctx, agentID)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	if agent == nil {
		err := errors.New("agent not found")
		tracing.TraceErr(span, err)
		return nil, err
	}

	return agent, nil
}

func (a *agentService) GetAllAgents(ctx context.Context) ([]*postgresentity.Agent, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentService.GetAllAgentsByTenant")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	agents, err := a.postgresRepositories.AgentRepository.GetAll(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return agents, nil
}

func (a *agentService) CreateAgent(ctx context.Context, agentType enum.AgentType) (*postgresentity.Agent, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentService.CreateAgent")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	tenant := common.GetTenantFromContext(ctx)
	if tenant == "" {
		err := errors.New("tenant not found in context")
		tracing.TraceErr(span, err)
		return nil, err
	}

	// get config from registry
	agentRegistry, err := a.postgresRepositories.AgentRegistryRepository.FindByType(ctx, agentType)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	if agentRegistry == nil {
		err := errors.New("no agent found in agent repository")
		tracing.TraceErr(span, err)
		return nil, err
	}

	// if agent registry is unique, check no other instances of agent exist before creation
	if agentRegistry.IsUnique {
		agents, err := a.postgresRepositories.AgentRepository.GetAllAgentsByTypes(ctx, []enum.AgentType{agentType})
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}
		if len(agents) > 0 {
			err := errors.New("agent already exists")
			tracing.TraceErr(span, err)
			return nil, err
		}
	}

	// build new agent
	agent := postgresentity.Agent{
		Type:        agentType,
		Scope:       agentRegistry.Scope,
		Tenant:      tenant,
		Name:        agentRegistry.AgentName,
		IsActive:    false,
		VisibleInUI: true,
		Icon:        agentRegistry.Icon,
		Color:       utils.GetRandomColor(),
		RegistryID:  agentRegistry.ID,
	}

	if agentRegistry.Scope == enum.AgentScopePersonal {
		agent.Owner = common.GetUserIdFromContext(ctx)
	}

	// build capabilities from registry
	var agentCapabilities []postgresentity.Capability
	for i, registryCapability := range agentRegistry.Capabilities {
		// get default config for each capability
		capabilityType, err := enum.GetAgentCapability(registryCapability)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}
		executor, err := a.agentCapabilities.GetExecutor(capabilityType)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}
		defaultCapability, err := a.createDefaultCapability(ctx, capabilityType, executor.Name(), i+1)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}

		agentCapabilities = append(agentCapabilities, *defaultCapability)
	}
	agent.Capabilities = agentCapabilities

	// build listeners from registry
	var agentListeners []postgresentity.Listener
	position := 1
	for _, registryListener := range agentRegistry.ListenerEvents {
		listenerEvent, err := enum.GetAgentListener(registryListener)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}
		listenerHandler, err := a.agentListeners.GetListener(listenerEvent)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}
		defaultListener, err := a.createDefaultListener(ctx, listenerEvent, listenerHandler.Name(), position)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}
		position++

		agentListeners = append(agentListeners, *defaultListener)
	}
	agent.Listeners = agentListeners

	// create agent instance in database
	newAgent, err := a.postgresRepositories.AgentRepository.Create(ctx, agent)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	tracing.TagEntity(span, newAgent.ID)

	// publish agent created event
	err = a.events.Publisher.PublishFanoutEvent(ctx, newAgent.ID, model.AGENT, dto.CreateAgent{
		Active:       newAgent.IsActive,
		Name:         newAgent.Name,
		Type:         newAgent.Type.String(),
		Icon:         newAgent.Icon,
		Color:        newAgent.Color,
		Capabilities: newAgent.Capabilities,
		Listeners:    newAgent.Listeners,
		Status:       newAgent.Status,
		FlowID:       newAgent.FlowID,
		VisibleInUI:  newAgent.VisibleInUI,
		RegistryID:   newAgent.RegistryID,
	})
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "unable to publish message CreateAgent"))
	}

	return newAgent, nil
}

func (a *agentService) DeleteAgent(ctx context.Context, agentId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentService.DeleteAgent")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	return a.postgresRepositories.AgentRepository.Delete(ctx, agentId)
}

func (r *agentService) GetAgentInfo(ctx context.Context) (*map[enum.AgentType]interfaces.AgentInfo, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentRegistryService.GetAgentInfo")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	agentRegistry, err := r.postgresRepositories.AgentRegistryRepository.FindAll(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	infoMap := make(map[enum.AgentType]interfaces.AgentInfo)
	for _, agent := range agentRegistry {
		infoMap[agent.Type] = interfaces.AgentInfo{
			Goal:        agent.Goal,
			Description: agent.Description,
			Metric:      agent.Metric,
		}
	}
	return &infoMap, nil
}

func (a *agentService) GetNorthStarMetricById(ctx context.Context, agentID string, agentType enum.AgentType) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentService.GetNorthStarMetricById")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	registryAgent, err := a.postgresRepositories.AgentRegistryRepository.FindByType(ctx, agentType)
	if err != nil || registryAgent == nil {
		tracing.TraceErr(span, err)
		return "", err
	}

	switch agentType {
	case enum.AgentCashflowGuardian:
		tenantSettings, err := a.tenantSettingsService.GetTenantSettings(ctx)
		if err != nil {
			tracing.TraceErr(span, err)
			return "", nil
		}
		output := strings.Replace(registryAgent.Metric, "{currencySymbol}", tenantSettings.BaseCurrency.Symbol(), 1)

		// get invoice ids
		invoiceIds, err := a.postgresRepositories.AgentExecutionRepository.GetGoalAchievedImpactedIdsLast30Days(ctx, agentID)
		if err != nil {
			tracing.TraceErr(span, err)
			return "", nil
		}

		invoices, err := a.invoiceService.GetInvoicesByIds(ctx, invoiceIds)
		if err != nil {
			tracing.TraceErr(span, err)
			return "", nil
		}
		amount := float64(0)
		if invoices != nil && len(*invoices) >= 0 {
			for _, invoice := range *invoices {
				if invoice.Currency == tenantSettings.BaseCurrency {
					amount += invoice.TotalAmount
				} else {
					amountInBaseCurrency, err := a.currencyService.GetAmountInCurrency(ctx, invoice.TotalAmount, invoice.Currency.String(), tenantSettings.BaseCurrency.String())
					if err != nil {
						tracing.TraceErr(span, err)
						return "", nil
					}
					amount += amountInBaseCurrency
				}
			}
		}
		output = strings.Replace(output, "{amount}", humanize.CommafWithDigits(amount, 2), 1)
		span.LogKV("output", output)
		return output, nil
	default:
		goalAchievedCount, err := a.postgresRepositories.AgentExecutionRepository.GetGoalAchievedCountLast30Days(ctx, agentID)
		if err != nil {
			tracing.TraceErr(span, err)
			return "", err
		}
		output := strings.Replace(registryAgent.Metric, "{count}", strconv.FormatInt(goalAchievedCount, 10), 1)
		span.LogKV("output", output)
		return output, nil
	}
}

func (a *agentService) createDefaultCapability(ctx context.Context, capabilityType enum.AgentCapability, capabilityName string, position int) (*postgres_entity.Capability, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentService.createDefaultCapability")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	executor, err := a.agentCapabilities.GetExecutor(capabilityType)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	config := executor.DefaultConfig()
	agentCapability := postgresentity.Capability{
		ID:       utils.GenerateNanoIdWithPrefix("cap", 16),
		Name:     capabilityName,
		Type:     capabilityType,
		Active:   executor.DefaultActive(),
		Tenant:   common.GetTenantFromContext(ctx),
		Position: position,
	}
	err = agentCapability.SetConfig(config)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	tracing.LogObjectAsJson(span, "result", agentCapability)
	return &agentCapability, nil
}

func (a *agentService) createDefaultListener(ctx context.Context, listenerEvent enum.AgentListenerEvent, listenerName string, position int) (*postgres_entity.Listener, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentService.createDefaultListener")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	handler, err := a.agentListeners.GetListener(listenerEvent)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	config := handler.DefaultConfig()
	agentListener := postgresentity.Listener{
		ID:       utils.GenerateNanoIdWithPrefix("lst", 16),
		Name:     listenerName,
		Type:     listenerEvent,
		Active:   true,
		Tenant:   common.GetTenantFromContext(ctx),
		Position: position,
	}
	err = agentListener.SetConfig(config)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	tracing.LogObjectAsJson(span, "result", agentListener)
	return &agentListener, nil
}

func (a *agentService) UpdateAgent(ctx context.Context, agentId string, agentFields data_fields.AgentFields, capabilities []postgresentity.Capability, listeners []postgresentity.Listener) (*postgresentity.Agent, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentService.UpdateAgent")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	tracing.TagEntity(span, agentId)
	tracing.LogObjectAsJson(span, "agentFields", agentFields)
	tracing.LogObjectAsJson(span, "capabilities", capabilities)

	agentEntity, err := a.postgresRepositories.AgentRepository.GetById(ctx, agentId)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	if agentEntity == nil {
		err := errors.New("agent not found")
		tracing.TraceErr(span, err)
		return nil, err
	}
	if agentFields.Name != nil {
		agentEntity.Name = *agentFields.Name
	}
	if agentFields.Icon != nil {
		agentEntity.Icon = *agentFields.Icon
	}
	if agentFields.Color != nil {
		agentEntity.Color = *agentFields.Color
	}
	if agentFields.Active != nil {
		agentEntity.IsActive = *agentFields.Active
	}
	if agentFields.VisibleInUI != nil {
		agentEntity.VisibleInUI = *agentFields.VisibleInUI
	}
	if agentFields.FlowID != nil {
		agentEntity.FlowID = *agentFields.FlowID
	}

	err = a.updateCapabilities(ctx, agentEntity, capabilities)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	err = a.updateListeners(ctx, agentEntity, listeners)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	err = a.postgresRepositories.AgentRepository.UpdateCapabilities(ctx, agentEntity.Capabilities)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	err = a.postgresRepositories.AgentRepository.UpdateListeners(ctx, agentEntity.Listeners)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	updatedAgent, err := a.postgresRepositories.AgentRepository.Update(ctx, *agentEntity)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	eventFields := dto.UpdateAgent{AgentFields: agentFields, Capabilities: capabilities, Listeners: listeners}
	err = a.events.Publisher.PublishFanoutEvent(ctx, agentId, model.AGENT, eventFields)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "unable to publish message {UpdateAgent}"))
	}

	return updatedAgent, nil
}

func (a *agentService) updateCapabilities(ctx context.Context, agentEntity *postgresentity.Agent, capabilities []postgresentity.Capability) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentService.updateCapabilities")
	defer span.Finish()

	agentEntity.UpdateCapabilities(capabilities)

	allCapabilitiesValid := true
	for i, capability := range agentEntity.Capabilities {
		if !capability.Active {
			// Not active => auto valid
			continue
		}

		executor, err := a.agentCapabilities.GetExecutor(capability.Type)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}
		config := executor.DefaultConfig()
		err = capability.GetConfig(&config)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}

		// see if typedConfig implements ConfigValidator
		if validator, ok := config.(agent_capability.ConfigValidator); ok {
			if !validator.Validate() {
				allCapabilitiesValid = false
			}
			err = capability.SetConfig(config)
			if err != nil {
				tracing.TraceErr(span, err)
				return err
			}
			agentEntity.Capabilities[i] = capability
		}
	}
	agentEntity.Configured = allCapabilitiesValid
	return nil
}

func (a *agentService) updateListeners(ctx context.Context, agentEntity *postgresentity.Agent, listeners []postgresentity.Listener) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentService.updateListeners")
	defer span.Finish()

	agentEntity.UpdateListeners(listeners)

	allListenersValid := true
	for i, listener := range agentEntity.Listeners {
		if !listener.Active {
			// Not active => auto valid
			continue
		}

		handler, err := a.agentListeners.GetListener(listener.Type)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}
		config := handler.DefaultConfig()
		err = listener.GetConfig(&config)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}

		// see if typedConfig implements ConfigValidator
		if validator, ok := config.(agent_capability.ConfigValidator); ok {
			if !validator.Validate() {
				allListenersValid = false
			}
			err = listener.SetConfig(config)
			if err != nil {
				tracing.TraceErr(span, err)
				return err
			}
			agentEntity.Listeners[i] = listener
		}
	}
	agentEntity.Configured = agentEntity.Configured && allListenersValid
	return nil
}

func (a *agentService) CreateAgentExecutionRecord(ctx context.Context, agent postgresentity.Agent, triggerEvent, traceId string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentService.CreateAgentExecutionRecord")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	agentExecutionRecord := postgresentity.AgentExecution{
		Tenant:       agent.Tenant,
		AgentID:      &agent.ID,
		TriggerEvent: triggerEvent,
		Status:       enum.AgentExecutionRunning,
		StartedAt:    utils.NowPtr(),
		TraceId:      traceId,
	}

	createdRecord, err := a.postgresRepositories.AgentExecutionRepository.Create(ctx, agentExecutionRecord)
	if err != nil {
		tracing.TraceErr(span, err)
		return "", err
	}

	if createdRecord == nil {
		err := errors.New("unable to create agent execution record")
		tracing.TraceErr(span, err)
		return "", err
	}

	return createdRecord.ID, nil
}
