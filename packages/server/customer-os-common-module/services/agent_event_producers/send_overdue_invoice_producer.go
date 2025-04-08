package agent_producers

import (
	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/constants"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/dto"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/model"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/agent_listeners"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/events"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	neo4j_entity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"
	neo4jmapper "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/mapper"
	neo4j_repository "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/repository"
	postgres_repository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

type SendOverdueInvoiceProducer struct {
	postgresRepository *postgres_repository.Repositories
	neo4jRepository    *neo4j_repository.Repositories
	events             *events.EventsService
	log                logger.Logger
}

func NewSendOverdueInvoiceProducer(
	postgresRepository *postgres_repository.Repositories,
	neo4jRepository *neo4j_repository.Repositories,
	events *events.EventsService,
	log logger.Logger,
) *SendOverdueInvoiceProducer {
	return &SendOverdueInvoiceProducer{
		postgresRepository: postgresRepository,
		neo4jRepository:    neo4jRepository,
		events:             events,
		log:                log,
	}
}

func (p *SendOverdueInvoiceProducer) Execute() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Cancel context on exit

	spans, ctx := telemetry.StartServiceSpan(ctx, "SendOverdueInvoiceProducer.Execute")
	defer spans.Finish()
	tracing.TagComponentCronJob(span)

	referenceTime := utils.Now()

	// Get all agents for cashflow guardian
	agents, err := p.postgresRepository.AgentRepository.GetActiveConfiguredAgentsByTypesCrossTenant(ctx, []enum.AgentType{enum.AgentCashflowGuardian})
	if err != nil {
		spans.TraceError(err)
		return
	}
	if len(agents) == 0 {
		p.log.Infof("No agents found for invoicing")
		return
	}

	limit := 100

	for _, agent := range agents {
		tenant := agent.Tenant
		agentCtx := common.WithCustomContext(ctx, &common.CustomContext{
			Tenant:    tenant,
			AppSource: constants.AppSourceUpkeeper,
		})

		// check if agent listener for overdue invoice is enabled
		overdueDays := 0
		listenerIsActive := false
		for _, listener := range agent.Listeners {
			if listener.Type != enum.EventInvoicePastDue {
				continue
			}
			if !listener.Active {
				break
			}
			var config agent_listeners.PastDueInvoiceConfig
			err = listener.GetConfig(&config)
			if err != nil {
				spans.TraceError(err)
				break
			}
			overdueDays = int(config.OverdueDays.Value)
			listenerIsActive = true
		}
		if !listenerIsActive {
			continue
		}

		for {
			select {
			case <-ctx.Done():
				p.log.Infof("Context cancelled, stopping")
				return
			default:
				// continue as normal
			}

			records, err := p.neo4jRepository.InvoiceReadRepository.GetInvoicesForPastDueNotifications(agentCtx, tenant, referenceTime, overdueDays, limit)
			if err != nil {
				spans.TraceError(err)
				p.log.Errorf("Error getting invoices for past due notifications: %s", err.Error())
				return
			}

			// no invoices found
			if len(records) == 0 {
				break // exit the inner loop, then process the next agent
			}

			// process records
			for _, record := range records {
				innerCtx := common.WithCustomContext(ctx, &common.CustomContext{
					Tenant:    tenant,
					AppSource: constants.AppSourceUpkeeper,
				})
				recordSpan, innerCtx := tracing.StartTracerSpan(innerCtx, "SendOverdueInvoiceProducer.Execute.Record")
				defer recordSpan.Finish()
				tracing.TagTenant(recordSpan, record.Tenant)

				invoice := neo4jmapper.MapDbNodeToInvoiceEntity(record.Node)

				// Mark notification requested, to avoid double notifications
				err = p.neo4jRepository.CommonWriteRepository.UpdateTimeProperty(innerCtx, record.Tenant, model.NodeLabelInvoice, invoice.Id, string(neo4j_entity.InvoicePropertyRemindInvoiceNotificationRequestedAt), utils.NowPtr())
				if err != nil {
					tracing.TraceErr(recordSpan, err)
					continue
				}

				err = p.events.Publisher.PublishFanoutEvent(innerCtx, invoice.Id, model.INVOICE, dto.PastDueInvoice{
					InvoiceId: invoice.Id,
				})
				if err != nil {
					spans.TraceError(errors.Wrap(err, "error publishing past due invoice event"))
					continue
				}
			}
		}
	}
}
