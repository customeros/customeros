package service

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/runner/customer-os-data-upkeeper/config"
	"github.com/openline-ai/openline-customer-os/packages/runner/customer-os-data-upkeeper/constants"
	"github.com/openline-ai/openline-customer-os/packages/runner/customer-os-data-upkeeper/events_processing_client"
	"github.com/openline-ai/openline-customer-os/packages/runner/customer-os-data-upkeeper/logger"
	"github.com/openline-ai/openline-customer-os/packages/runner/customer-os-data-upkeeper/repository"
	"github.com/openline-ai/openline-customer-os/packages/runner/customer-os-data-upkeeper/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	commonpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/common"
	invoicepb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/invoice"
	"time"
)

type InvoiceService interface {
	GenerateInvoices()
}

type invoiceService struct {
	cfg                    *config.Config
	log                    logger.Logger
	repositories           *repository.Repositories
	eventsProcessingClient *events_processing_client.Client
}

func NewInvoiceService(cfg *config.Config, log logger.Logger, repositories *repository.Repositories, client *events_processing_client.Client) InvoiceService {
	return &invoiceService{
		cfg:                    cfg,
		log:                    log,
		repositories:           repositories,
		eventsProcessingClient: client,
	}
}

func (s *invoiceService) GenerateInvoices() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Cancel context on exit

	span, ctx := tracing.StartTracerSpan(ctx, "InvoiceService.GenerateInvoices")
	defer span.Finish()

	if s.eventsProcessingClient == nil {
		s.log.Warn("eventsProcessingClient is nil. Will not update next cycle date.")
		return
	}

	referenceTime := utils.Now()
	dryRun := false
	tenantDefaultCurrencies := make(map[string]neo4jenum.Currency)

	for {
		select {
		case <-ctx.Done():
			s.log.Infof("Context cancelled, stopping")
			return
		default:
			// continue as normal
		}

		records, err := s.repositories.Neo4jRepositories.ContractReadRepository.GetContractsToGenerateOnCycleInvoices(ctx, referenceTime)
		if err != nil {
			tracing.TraceErr(span, err)
			s.log.Errorf("Error getting contracts for invoicing: %v", err)
			return
		}

		// no contracts found
		if len(records) == 0 {
			return
		}

		//process records
		for _, record := range records {
			contract := neo4jmapper.MapDbNodeToContractEntity(record.Node)
			tenant := record.Tenant

			currency := contract.Currency.String()
			if currency == "" {
				currency = s.getTenantDefaultCurrency(ctx, tenant, tenantDefaultCurrencies).String()
			}

			var invoicePeriodStart, invoicePeriodEnd time.Time
			if contract.NextInvoiceDate != nil {
				invoicePeriodStart = *contract.NextInvoiceDate
			} else {
				invoicePeriodStart = *contract.InvoicingStartDate
			}
			invoicePeriodEnd = calculateInvoiceCycleEnd(invoicePeriodStart, contract.BillingCycle)

			readyToRequestInvoice := invoicePeriodEnd.After(invoicePeriodStart)
			if readyToRequestInvoice {
				newInvoiceRequest := invoicepb.NewInvoiceForContractRequest{
					Tenant:             record.Tenant,
					ContractId:         contract.Id,
					Currency:           currency,
					InvoicePeriodStart: utils.ConvertTimeToTimestampPtr(&invoicePeriodStart),
					InvoicePeriodEnd:   utils.ConvertTimeToTimestampPtr(&invoicePeriodEnd),
					DryRun:             dryRun,
					SourceFields: &commonpb.SourceFields{
						AppSource: constants.AppSourceDataUpkeeper,
						Source:    neo4jentity.DataSourceOpenline.String(),
					},
				}
				// TODO alexb send billing cycle
				//switch contract.BillingCycle {
				//case neo4jenum.BillingCycleMonthlyBilling:
				//	newInvoiceRequest.BillingCycle = commonpb.BillingCycle_MONTHLY_BILLING
				//
				//}
				_, err = s.eventsProcessingClient.InvoiceClient.NewInvoiceForContract(ctx, &newInvoiceRequest)

				if err != nil {
					tracing.TraceErr(span, err)
					s.log.Errorf("Error generating invoice for contract %s: %s", contract.Id, err.Error())
				}
			}
			// mark invoicing started as long dry run is false
			nextInvoiceDate := contract.NextInvoiceDate
			if !dryRun {
				nextInvoiceDate = utils.ToPtr(invoicePeriodEnd.AddDate(0, 0, 1))
			}
			err = s.repositories.Neo4jRepositories.ContractWriteRepository.MarkInvoicingStarted(ctx, tenant, contract.Id, utils.Now(), nextInvoiceDate)
			if err != nil {
				tracing.TraceErr(span, err)
				s.log.Errorf("Error marking invoicing started for contract %s: %s", contract.Id, err.Error())
			}
		}
		//sleep for async processing, then check again
		time.Sleep(5 * time.Second)
	}
}

func calculateInvoiceCycleEnd(start time.Time, cycle neo4jenum.BillingCycle) time.Time {
	var end time.Time
	switch cycle {
	case neo4jenum.BillingCycleMonthlyBilling:
		end = start.AddDate(0, 1, 0)
	case neo4jenum.BillingCycleQuarterlyBilling:
		end = start.AddDate(0, 3, 0)
	case neo4jenum.BillingCycleAnnualBilling:
		end = start.AddDate(1, 0, 0)
	default:
		return start
	}
	previousDay := end.AddDate(0, 0, -1)
	return previousDay
}

func (s *invoiceService) getTenantDefaultCurrency(ctx context.Context, tenant string, tenantDefaultCurrencies map[string]neo4jenum.Currency) neo4jenum.Currency {
	if currency, ok := tenantDefaultCurrencies[tenant]; ok {
		return currency
	}

	dbNode, _ := s.repositories.Neo4jRepositories.TenantReadRepository.GetTenantSettings(ctx, tenant)
	tenantSettings := neo4jmapper.MapDbNodeToTenantSettingsEntity(dbNode)

	tenantDefaultCurrencies[tenant] = tenantSettings.DefaultCurrency
	return tenantSettings.DefaultCurrency
}
