package agent_producers

import (
	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/constants"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/dto"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/model"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/events"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	neo4jmapper "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/mapper"
	neo4j_repository "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/repository"
	postgres_repository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"time"
)

type InvoiceProducer struct {
	postgresRepository *postgres_repository.Repositories
	neo4jRepository    *neo4j_repository.Repositories
	events             *events.EventsService
	log                logger.Logger
}

func NewInvoiceProducer(
	postgresRepository *postgres_repository.Repositories,
	neo4jRepository *neo4j_repository.Repositories,
	events *events.EventsService,
	log logger.Logger,
) *InvoiceProducer {
	return &InvoiceProducer{
		postgresRepository: postgresRepository,
		neo4jRepository:    neo4jRepository,
		events:             events,
		log:                log,
	}
}

// Add all Agent types subscribed to this event here
func (p *InvoiceProducer) subscribedAgents() []enum.AgentType {
	return []enum.AgentType{
		enum.AgentCashflowGuardian,
	}
}

func (p *InvoiceProducer) Execute() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Cancel context on exit

	span, ctx := tracing.StartTracerSpan(ctx, "InvoiceProducer.Execute")
	defer span.Finish()
	tracing.TagComponentCronJob(span)

	referenceTime := utils.Now()

	// Get all agents for cashflow guardian
	agents, err := p.postgresRepository.AgentRepository.GetActiveConfiguredAgentsByTypesCrossTenant(ctx, []enum.AgentType{enum.AgentCashflowGuardian})
	if err != nil {
		tracing.TraceErr(span, err)
		return
	}

	limit := 0 // TODO set to non zero to start invoice generation
	delayFromPreviousInvoicingAttemptInMinutes := 240

	for _, cashflowGuardianAgent := range agents {
		tenant := cashflowGuardianAgent.Tenant
		for {
			select {
			case <-ctx.Done():
				p.log.Infof("Context cancelled, stopping")
				return
			default:
				// continue as normal
			}

			records, err := p.neo4jRepository.ContractReadRepository.GetContractsToGenerateCycleInvoices(ctx, tenant, referenceTime, delayFromPreviousInvoicingAttemptInMinutes, limit)
			if err != nil {
				tracing.TraceErr(span, err)
				p.log.Errorf("Error getting contracts for invoicing: %s", err.Error())
				return
			}

			// no contracts found
			if len(records) == 0 {
				break // exit the inner loop, then process the next agent
			}

			// process records
			for _, record := range records {
				innerCtx := common.WithCustomContext(ctx, &common.CustomContext{
					Tenant:    tenant,
					AppSource: constants.AppSourceUpkeeper,
				})
				recordSpan, innerCtx := tracing.StartTracerSpan(innerCtx, "InvoiceProducer.Execute.Record")
				defer recordSpan.Finish()
				tracing.TagTenant(recordSpan, record.Tenant)

				contract := neo4jmapper.MapDbNodeToContractEntity(record.Node)

				// mark invoicing requested to prevent double processing of the record
				err = p.neo4jRepository.ContractWriteRepository.MarkCycleInvoicingRequested(innerCtx, tenant, contract.Id, utils.Now())
				if err != nil {
					tracing.TraceErr(recordSpan, errors.Wrap(err, "Error marking invoicing requested"))
					p.log.Errorf("Error marking invoicing started for contract %s: %s", contract.Id, err.Error())
					return
				}

				if contract.PayAutomatically {
					err = p.events.Publisher.PublishFanoutEvent(innerCtx, contract.Id, model.CONTRACT, dto.InvoiceContractWithAutopayment{
						ContractId: contract.Id,
						DryRun:     false,
						Preview:    false,
					})
				} else {
					err = p.events.Publisher.PublishFanoutEvent(innerCtx, contract.Id, model.CONTRACT, dto.InvoiceContract{
						ContractId: contract.Id,
						DryRun:     false,
						Preview:    false,
					})
				}
				if err != nil {
					tracing.TraceErr(span, errors.Wrap(err, "error publishing new invoice event"))
					return
				}
			}

			time.Sleep(1 * time.Second)
		}
	}
}
