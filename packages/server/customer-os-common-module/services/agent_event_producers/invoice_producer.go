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
	neo4jentity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"
	neo4jmapper "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/mapper"
	neo4j_repository "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/repository"
	postgres_repository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
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
	if len(agents) == 0 {
		p.log.Infof("No agents found for invoicing")
		return
	}

	limit := 100
	delayFromPreviousInvoicingAttemptInMinutes := 240

	for _, cashflowGuardianAgent := range agents {
		tenant := cashflowGuardianAgent.Tenant
		agentCtx := common.WithCustomContext(ctx, &common.CustomContext{
			Tenant:    tenant,
			AppSource: constants.AppSourceUpkeeper,
		})
		for {
			select {
			case <-agentCtx.Done():
				p.log.Infof("Context cancelled, stopping")
				return
			default:
				// continue as normal
			}

			records, err := p.neo4jRepository.ContractReadRepository.GetContractsToGenerateCycleInvoices(agentCtx, tenant, referenceTime, delayFromPreviousInvoicingAttemptInMinutes, limit)
			if err != nil {
				tracing.TraceErr(span, err)
				p.log.Errorf("Error getting contracts for invoicing: %s", err.Error())
				return
			}

			// no contracts found
			if len(records) == 0 {
				break // exit the inner loop, then process the next agent
			}

			// prepare postpaid flag
			dbNode, err := p.neo4jRepository.TenantReadRepository.GetTenantSettings(agentCtx, tenant)
			if err != nil {
				tracing.TraceErr(span, err)
				return
			}
			tenantSettingsEntity := neo4jmapper.MapDbNodeToTenantSettingsEntity(dbNode)
			postpaid := tenantSettingsEntity.InvoicingPostpaid

			// process records
			for _, record := range records {
				innerCtx := common.WithCustomContext(agentCtx, &common.CustomContext{
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

				// check if contract is ready for invoicing by dates
				if !p.isReadyForInvoicing(innerCtx, tenant, *contract, postpaid) {
					continue
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

func (p *InvoiceProducer) isReadyForInvoicing(ctx context.Context, tenant string, contractEntity neo4jentity.ContractEntity, postpaid bool) bool {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceProducer.isReadyForInvoicing")
	defer span.Finish()
	tracing.TagTenant(span, tenant)
	span.LogFields(log.Bool("postpaid", postpaid))

	// prepare and validate dates
	var invoicePeriodStart, invoicePeriodEnd time.Time
	if contractEntity.NextInvoiceDate != nil {
		invoicePeriodStart = *contractEntity.NextInvoiceDate
	} else {
		invoicePeriodStart = *contractEntity.InvoicingStartDate
	}

	invoicePeriodEnd = p.prepareInvoiceCycleEndDate(ctx, invoicePeriodStart, tenant, contractEntity)

	contractReadyForInvoicingByDates := true
	if postpaid {
		contractReadyForInvoicingByDates = utils.EndOfDayInUTC(invoicePeriodEnd).Before(utils.Now())
	} else {
		contractReadyForInvoicingByDates = invoicePeriodEnd.After(invoicePeriodStart)
	}
	return contractReadyForInvoicingByDates
}

func (p *InvoiceProducer) prepareInvoiceCycleEndDate(ctx context.Context, start time.Time, tenant string, contractEntity neo4jentity.ContractEntity) time.Time {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceProducer.prepareInvoiceCycleEndDate")
	defer span.Finish()

	nextStart := start.AddDate(0, int(contractEntity.BillingCycleInMonths), 0)
	if start.Day() == 1 {
		// if previous invoice was generated end of month, we need to substract extra 1 day
		previousCycleInvoiceDbNode, err := p.neo4jRepository.InvoiceReadRepository.GetPreviousCycleInvoice(ctx, tenant, contractEntity.Id)
		if err != nil {
			tracing.TraceErr(nil, errors.Wrap(err, "Error getting previous cycle invoice"))
		}
		if previousCycleInvoiceDbNode != nil {
			previousInvoice := neo4jmapper.MapDbNodeToInvoiceEntity(previousCycleInvoiceDbNode)
			if previousInvoice.PeriodStartDate.Day() != 1 {
				nextStart = nextStart.AddDate(0, -1, 0)
				nextStart = time.Date(nextStart.Year(), nextStart.Month(), previousInvoice.PeriodStartDate.Day(), 0, 0, 0, 0, nextStart.Location())
			}
		}
	}
	invoiceCycleEnd := nextStart.AddDate(0, 0, -1)
	span.LogFields(log.Object("result.invoiceCycleEnd", invoiceCycleEnd))
	return invoiceCycleEnd
}
