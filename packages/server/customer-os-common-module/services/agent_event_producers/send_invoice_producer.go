package agent_producers

import (
	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/constants"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/dto"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/model"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/events"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/opentracing/opentracing-go/log"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"
	neo4jmapper "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/mapper"
	neo4j_repository "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/repository"
	postgres_repository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"

	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"time"
)

type SendInvoiceProducer struct {
	postgresRepository *postgres_repository.Repositories
	neo4jRepository    *neo4j_repository.Repositories
	events             *events.EventsService
	log                logger.Logger
}

func NewSendInvoiceProducer(
	postgresRepository *postgres_repository.Repositories,
	neo4jRepository *neo4j_repository.Repositories,
	events *events.EventsService,
	log logger.Logger,
) *SendInvoiceProducer {
	return &SendInvoiceProducer{
		postgresRepository: postgresRepository,
		neo4jRepository:    neo4jRepository,
		events:             events,
		log:                log,
	}
}

func (p *SendInvoiceProducer) Execute() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Cancel context on exit

	spans, ctx := telemetry.StartCronSpan(ctx, "SendInvoiceProducer.Execute")
	defer spans.Finish()

	limit := 100
	minutesFromLastAttempt := 360
	minutesFromCreation := 30
	lookBackWindowDays := 5

	for {
		select {
		case <-ctx.Done():
			p.log.Infof("Context cancelled, stopping")
			return
		default:
			// continue as normal
		}

		records, err := p.neo4jRepository.InvoiceReadRepository.GetInvoicesForPayNotifications(ctx, minutesFromCreation, minutesFromLastAttempt, lookBackWindowDays, limit)
		if err != nil {
			spans.TraceError(err)
			p.log.Errorf("Error getting invoices for pay notifications: %v", err)
			return
		}

		// no invoices found
		if len(records) == 0 {
			return
		}

		// process records
		for _, record := range records {
			innerCtx := common.WithCustomContext(ctx, &common.CustomContext{
				Tenant:    record.Tenant,
				AppSource: constants.AppSourceUpkeeper,
			})
			invoice := neo4jmapper.MapDbNodeToInvoiceEntity(record.Node)

			// mark pay notification requested
			err = p.neo4jRepository.InvoiceWriteRepository.MarkPayNotificationRequested(innerCtx, record.Tenant, invoice.Id, utils.Now())
			if err != nil {
				spans.TraceError(err)
				p.log.Errorf("Error marking invoice %s as pay notification requested: %s", invoice.Id, err.Error())
				return
			}

			err = p.events.Publisher.PublishFanoutEvent(innerCtx, invoice.Id, model.INVOICE, dto.SendInvoice{
				InvoiceId: invoice.Id,
			})
		}
	}
}

func (p *SendInvoiceProducer) isReadyForInvoicing(ctx context.Context, tenant string, contractEntity neo4jentity.ContractEntity, postpaid bool) bool {
	spans, ctx := telemetry.StartServiceSpan(ctx, "SendInvoiceProducer.isReadyForInvoicing")
	defer spans.Finish()
	spans.LogFields(log.Bool("postpaid", postpaid))

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

func (p *SendInvoiceProducer) prepareInvoiceCycleEndDate(ctx context.Context, start time.Time, tenant string, contractEntity neo4jentity.ContractEntity) time.Time {
	spans, ctx := telemetry.StartServiceSpan(ctx, "SendInvoiceProducer.prepareInvoiceCycleEndDate")
	defer spans.Finish()

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
	spans.LogFields(log.Object("result.invoiceCycleEnd", invoiceCycleEnd))
	return invoiceCycleEnd
}
