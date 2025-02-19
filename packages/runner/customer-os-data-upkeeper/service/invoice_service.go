package service

import (
	"context"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/agent_capability"
	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	"time"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/data_fields"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	commonService "github.com/customeros/customeros/packages/server/customer-os-common-module/services"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"
	neo4jenum "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/enum"
	neo4jmapper "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/mapper"
	neo4jrepository "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/repository"
	"github.com/pkg/errors"

	"github.com/customeros/customeros/packages/runner/customer-os-data-upkeeper/config"
	"github.com/customeros/customeros/packages/runner/customer-os-data-upkeeper/constants"
	"github.com/customeros/customeros/packages/runner/customer-os-data-upkeeper/logger"
	"github.com/customeros/customeros/packages/runner/customer-os-data-upkeeper/repository"
)

type GeneratePaymentLinkEventBody struct {
	Tenant                       string `json:"tenant"`
	Currency                     string `json:"currency"`
	AmountInSmallestCurrencyUnit int64  `json:"amountInSmallestCurrencyUnit"`
	InvoiceId                    string `json:"invoiceId"`
	InvoiceDescription           string `json:"invoiceDescription"`
	CustomerEmail                string `json:"customerEmail"`
	PrimaryStripeCustomerId      string `json:"stripeCustomerId"`
}

type InvoiceService interface {
	CleanupInvoices()
	GenerateNextPreviewInvoices()
	AdjustInvoiceStatus()

	// TODO stopped offcycle invoicing
	GenerateOffCycleInvoices()
}

type invoiceService struct {
	cfg            *config.Config
	log            logger.Logger
	commonServices *commonService.CommonServices
	repositories   *repository.Repositories
}

func NewInvoiceService(cfg *config.Config, log logger.Logger, commonServices *commonService.CommonServices, repositories *repository.Repositories) InvoiceService {
	return &invoiceService{
		cfg:            cfg,
		log:            log,
		commonServices: commonServices,
		repositories:   repositories,
	}
}

func (s *invoiceService) GenerateNextPreviewInvoices() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Cancel context on exit

	span, ctx := tracing.StartTracerSpan(ctx, "InvoiceService.GenerateNextPreviewInvoices")
	defer span.Finish()
	tracing.TagComponentCronJob(span)

	referenceTime := utils.Now()

	// Get all agents for cashflow guardian
	agents, err := s.repositories.PostgresRepositories.AgentRepository.GetActiveConfiguredAgentsByTypesCrossTenant(ctx, []enum.AgentType{enum.AgentCashflowGuardian})
	if err != nil {
		tracing.TraceErr(span, err)
		return
	}
	if len(agents) == 0 {
		s.log.Infof("No agents found for invoicing")
		return
	}
	// collect tenants
	var tenants []string
	agentByTenant := map[string]postgres_entity.Agent{}
	for _, agent := range agents {
		tenants = append(tenants, agent.Tenant)
		agentByTenant[agent.Tenant] = agent
	}

	for {
		select {
		case <-ctx.Done():
			s.log.Infof("Context cancelled, stopping")
			return
		default:
			// continue as normal
		}

		records, err := s.repositories.Neo4jRepositories.ContractReadRepository.GetContractsToGenerateNextScheduledInvoices(ctx, tenants, referenceTime, 15)
		if err != nil {
			tracing.TraceErr(span, err)
			s.log.Errorf("Error getting contracts for invoicing: %v", err)
			return
		}

		// no contracts found
		if len(records) == 0 {
			return
		}

		// process records
		for _, record := range records {
			innerCtx := common.WithCustomContext(ctx, &common.CustomContext{
				Tenant:    record.Tenant,
				AppSource: constants.AppSourceDataUpkeeper,
			})
			recordSpan, innerCtx := tracing.StartTracerSpan(innerCtx, "InvoiceService.GenerateNextPreviewInvoices.Record")
			defer recordSpan.Finish()
			tracing.TagTenant(recordSpan, record.Tenant)
			contract := neo4jmapper.MapDbNodeToContractEntity(record.Node)
			tenant := record.Tenant

			// mark next preview invoice requested
			err = s.repositories.Neo4jRepositories.ContractWriteRepository.MarkNextPreviewInvoicingRequested(ctx, tenant, contract.Id, utils.Now())
			if err != nil {
				tracing.TraceErr(span, err)
				s.log.Errorf("Error marking invoicing started for contract %s: %s", contract.Id, err.Error())
				return
			}

			// get agent for tenant
			dataFields := data_fields.InvoiceFields{
				DryRun:  true,
				Preview: true,
			}
			agent := agentByTenant[tenant]
			capabilityConfig := agent_capability.GenerateInvoiceConfig{}
			err = agent.GetCapabilityConfigByType(enum.CapabilityGenerateInvoice, &capabilityConfig)
			if err != nil {
				dataFields.TenantBillingProfile = &data_fields.TenantBillingProfile{
					LogoRepositoryFileId:       capabilityConfig.LogoRepositoryFileId.Value,
					Country:                    capabilityConfig.Country.Value,
					LegalName:                  capabilityConfig.LegalName.Value,
					AddressLine1:               capabilityConfig.AddressLine1.Value,
					AddressLine2:               capabilityConfig.AddressLine2.Value,
					Zip:                        capabilityConfig.ZIP.Value,
					Locality:                   capabilityConfig.Locality.Value,
					Region:                     capabilityConfig.Region.Value,
					IncludeBankTransferDetails: capabilityConfig.IncludeBankTransferDetails.Value,
					BankName:                   capabilityConfig.BankName.Value,
					AccountNumber:              capabilityConfig.AccountNumber.Value,
					IBAN:                       capabilityConfig.IBAN.Value,
					BIC:                        capabilityConfig.BIC.Value,
					SortCode:                   capabilityConfig.SortCode.Value,
					RoutingNumber:              capabilityConfig.RoutingNumber.Value,
					OtherDetails:               capabilityConfig.OtherDetails.Value,
				}
			}
			_, err = s.commonServices.InvoiceService.InvoiceContract(innerCtx, nil, contract.Id, dataFields)
			if err != nil {
				tracing.TraceErr(span, err)
				s.log.Errorf("Error generating invoice for contract %s: %s", contract.Id, err.Error())
			}
		}
		// sleep for async processing, then check again
		time.Sleep(10 * time.Second)
	}
}

func (s *invoiceService) calculateInvoiceCycleEnd(ctx context.Context, start time.Time, tenant string, contractEntity neo4jentity.ContractEntity) time.Time {
	nextStart := start.AddDate(0, int(contractEntity.BillingCycleInMonths), 0)
	if start.Day() == 1 {
		// if previous invoice was generated end of month, we need to substract extra 1 day
		previousCycleInvoiceDbNode, err := s.repositories.Neo4jRepositories.InvoiceReadRepository.GetPreviousCycleInvoice(ctx, tenant, contractEntity.Id)
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
	return nextStart.AddDate(0, 0, -1)
}

func (s *invoiceService) getTenantBaseCurrency(ctx context.Context, tenant string, cachedTenantBaseCurrencies map[string]neo4jenum.Currency) neo4jenum.Currency {
	if currency, ok := cachedTenantBaseCurrencies[tenant]; ok {
		return currency
	}

	dbNode, _ := s.repositories.Neo4jRepositories.TenantReadRepository.GetTenantSettings(ctx, tenant)
	tenantSettings := neo4jmapper.MapDbNodeToTenantSettingsEntity(dbNode)

	currency := tenantSettings.BaseCurrency
	cachedTenantBaseCurrencies[tenant] = currency
	return currency
}

func (s *invoiceService) GenerateOffCycleInvoices() {
	return

	// ctx, cancel := context.WithCancel(context.Background())
	// defer cancel() // Cancel context on exit

	//span, ctx := tracing.StartTracerSpan(ctx, "InvoiceService.GenerateOffCycleInvoices")
	//defer span.Finish()
	//tracing.TagComponentCronJob(span)
	//
	//if s.cfg.App.ProcessConfig.OffCycleInvoicingEnabled == false {
	//	s.log.Infof("Off-cycle invoicing is disabled, stopping")
	//	span.LogFields(log.Bool("off_cycle_invoicing_enabled", s.cfg.App.ProcessConfig.OffCycleInvoicingEnabled))
	//	return
	//}
	//
	//if s.eventsProcessingClient == nil {
	//	err := errors.New("eventsProcessingClient is nil")
	//	tracing.TraceErr(span, err)
	//	s.log.Error(err.Error())
	//	return
	//}
	//
	//referenceTime := utils.Now()
	//dryRun := false
	//cachedTenantBaseCurrencies := make(map[string]neo4jenum.Currency)
	//
	//limit := 100
	//
	//for {
	//	select {
	//	case <-ctx.Done():
	//		s.log.Infof("Context cancelled, stopping")
	//		return
	//	default:
	//		// continue as normal
	//	}
	//
	//	records, err := s.repositories.Neo4jRepositories.ContractReadRepository.GetContractsToGenerateOffCycleInvoices(ctx, referenceTime, s.cfg.App.ProcessConfig.DelayGenerateOffCycleInvoiceInMinutes, limit)
	//	if err != nil {
	//		tracing.TraceErr(span, err)
	//		s.log.Errorf("Error getting contracts for off-cycle invoicing: %v", err)
	//		return
	//	}
	//
	//	// no contracts found
	//	if len(records) == 0 {
	//		return
	//	}
	//
	//	// process records
	//	for _, record := range records {
	//		contract := neo4jmapper.MapDbNodeToContractEntity(record.Node)
	//		tenant := record.Tenant
	//
	//		currency := contract.Currency.String()
	//		if currency == "" {
	//			currency = s.getTenantBaseCurrency(ctx, tenant, cachedTenantBaseCurrencies).String()
	//		}
	//
	//		invoicePeriodStart := utils.ToDate(referenceTime)
	//		invoicePeriodEnd := utils.ToDate(utils.IfNotNilTimeWithDefault(contract.NextInvoiceDate, referenceTime).AddDate(0, 0, -1))
	//
	//		readyToRequestInvoice := !invoicePeriodEnd.Before(invoicePeriodStart)
	//		if readyToRequestInvoice {
	//			newInvoiceRequest := invoicepb.NewInvoiceForContractRequest{
	//				Tenant:             record.Tenant,
	//				ContractId:         contract.Id,
	//				Currency:           currency,
	//				InvoicePeriodStart: utils.ConvertTimeToTimestampPtr(&invoicePeriodStart),
	//				InvoicePeriodEnd:   utils.ConvertTimeToTimestampPtr(&invoicePeriodEnd),
	//				DryRun:             dryRun,
	//				OffCycle:           true,
	//				SourceFields: &commonpb.SourceFields{
	//					AppSource: constants.AppSourceDataUpkeeper,
	//					Source:    neo4jentity.DataSourceOpenline.String(),
	//				},
	//			}
	//			_, err = CallEventsPlatformGRPCWithRetry[*invoicepb.InvoiceIdResponse](func() (*invoicepb.InvoiceIdResponse, error) {
	//				return s.eventsProcessingClient.InvoiceClient.NewInvoiceForContract(ctx, &newInvoiceRequest)
	//			})
	//			if err != nil {
	//				tracing.TraceErr(span, err)
	//				s.log.Errorf("Error generating off-cycle invoice for contract %s: %s", contract.Id, err.Error())
	//			}
	//		}
	//		// mark invoicing started
	//		err = s.repositories.Neo4jRepositories.ContractWriteRepository.MarkOffCycleInvoicingRequested(ctx, tenant, contract.Id, utils.Now())
	//		if err != nil {
	//			tracing.TraceErr(span, err)
	//			s.log.Errorf("Error marking invoicing started for contract %s: %s", contract.Id, err.Error())
	//		}
	//	}
	//	// sleep for async processing, then check again
	//	if len(records) < limit {
	//		return
	//	}
	//	time.Sleep(10 * time.Second)
	//}
}

func (s *invoiceService) CleanupInvoices() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Cancel context on exit

	span, ctx := tracing.StartTracerSpan(ctx, "InvoiceService.CleanupInvoices")
	defer span.Finish()
	tracing.TagComponentCronJob(span)

	for {
		select {
		case <-ctx.Done():
			s.log.Infof("Context cancelled, stopping")
			return
		default:
			// continue as normal
		}

		records, err := s.repositories.Neo4jRepositories.InvoiceReadRepository.GetExpiredDryRunInvoices(ctx)
		if err != nil {
			tracing.TraceErr(span, err)
			s.log.Errorf("Error getting invoices for cleanup: %v", err)
			return
		}

		// no invoices found
		if len(records) == 0 {
			return
		}

		// process records
		for _, record := range records {
			invoice := neo4jmapper.MapDbNodeToInvoiceEntity(record.Node)
			tenant := record.Tenant

			err = s.repositories.Neo4jRepositories.InvoiceWriteRepository.DeleteDryRunInvoice(ctx, tenant, invoice.Id)
			if err != nil {
				tracing.TraceErr(span, err)
				s.log.Errorf("Error deleting dry run invoice %s: %v", invoice.Id, err)
			}
		}
	}
}

func (s *invoiceService) AdjustInvoiceStatus() {
	s.updateInvoiceStatusToOverdue()
	s.updateInvoiceStatusToOnHold()
	s.updateInvoiceStatusScheduled()
	s.updateInvoiceStatusFromPaymentProcessingToDue()
}

func (s *invoiceService) updateInvoiceStatusToOverdue() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Cancel context on exit

	span, ctx := tracing.StartTracerSpan(ctx, "InvoiceService.updateInvoiceStatusToOverdue")
	defer span.Finish()
	tracing.TagComponentCronJob(span)

	limit := 500
	records, err := s.repositories.Neo4jRepositories.InvoiceReadRepository.GetInvoicesForOverdue(ctx, limit)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error getting invoices for overdue: %v", err)
		return
	}

	for _, record := range records {
		invoice := neo4jmapper.MapDbNodeToInvoiceEntity(record.Node)
		tenant := record.Tenant

		innerCtx := common.WithCustomContext(ctx, &common.CustomContext{
			Tenant:    tenant,
			AppSource: constants.AppSourceDataUpkeeper,
		})

		err = s.commonServices.InvoiceService.UpdateInvoice(innerCtx, nil, invoice.Id, neo4jrepository.InvoiceUpdateFields{
			Status:       neo4jenum.InvoiceStatusOverdue,
			UpdateStatus: true,
		})
		if err != nil {
			tracing.TraceErr(span, err)
			s.log.Errorf("Error updating invoice %s to overdue: %s", invoice.Id, err.Error())
			return // stop processing
		}
	}
}

func (s *invoiceService) updateInvoiceStatusToOnHold() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Cancel context on exit

	span, ctx := tracing.StartTracerSpan(ctx, "InvoiceService.updateInvoiceStatusToOnHold")
	defer span.Finish()
	tracing.TagComponentCronJob(span)

	limit := 500

	records, err := s.repositories.Neo4jRepositories.InvoiceReadRepository.GetInvoicesForOnHold(ctx, limit)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error getting invoices for on hold: %v", err)
		return
	}

	for _, record := range records {
		invoice := neo4jmapper.MapDbNodeToInvoiceEntity(record.Node)
		tenant := record.Tenant

		innerCtx := common.WithCustomContext(ctx, &common.CustomContext{
			Tenant:    tenant,
			AppSource: constants.AppSourceDataUpkeeper,
		})

		err = s.commonServices.InvoiceService.UpdateInvoice(innerCtx, nil, invoice.Id, neo4jrepository.InvoiceUpdateFields{
			Status:       neo4jenum.InvoiceStatusOnHold,
			UpdateStatus: true,
		})
		if err != nil {
			tracing.TraceErr(span, err)
			s.log.Errorf("Error updating invoice %s to on hold: %s", invoice.Id, err.Error())
			return // stop processing
		}
	}
}

func (s *invoiceService) updateInvoiceStatusScheduled() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Cancel context on exit

	span, ctx := tracing.StartTracerSpan(ctx, "InvoiceService.updateInvoiceStatusScheduled")
	defer span.Finish()
	tracing.TagComponentCronJob(span)

	limit := 500
	records, err := s.repositories.Neo4jRepositories.InvoiceReadRepository.GetInvoicesForScheduled(ctx, limit)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error getting invoices for scheduled: %v", err)
		return
	}

	for _, record := range records {
		invoice := neo4jmapper.MapDbNodeToInvoiceEntity(record.Node)
		tenant := record.Tenant

		innerCtx := common.WithCustomContext(ctx, &common.CustomContext{
			Tenant:    tenant,
			AppSource: constants.AppSourceDataUpkeeper,
		})

		err = s.commonServices.InvoiceService.UpdateInvoice(innerCtx, nil, invoice.Id, neo4jrepository.InvoiceUpdateFields{
			Status:       neo4jenum.InvoiceStatusScheduled,
			UpdateStatus: true,
		})
		if err != nil {
			tracing.TraceErr(span, err)
			s.log.Errorf("Error updating invoice %s to scheduled: %s", invoice.Id, err.Error())
			return // stop processing
		}
	}
}

func (s *invoiceService) updateInvoiceStatusFromPaymentProcessingToDue() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	span, ctx := tracing.StartTracerSpan(ctx, "InvoiceService.updateInvoiceStatusFromPaymentProcessingToDue")
	defer span.Finish()
	tracing.TagComponentCronJob(span)

	limit := 500
	paymentProcessingMaxDays := 8
	records, err := s.repositories.Neo4jRepositories.InvoiceReadRepository.GetExpiredPaymentProcessingInvoices(ctx, paymentProcessingMaxDays, limit)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error getting invoices for overdue: %v", err)
		return
	}

	for _, record := range records {
		invoice := neo4jmapper.MapDbNodeToInvoiceEntity(record.Node)
		tenant := record.Tenant

		innerCtx := common.WithCustomContext(ctx, &common.CustomContext{
			Tenant:    tenant,
			AppSource: constants.AppSourceDataUpkeeper,
		})

		err = s.commonServices.InvoiceService.UpdateInvoice(innerCtx, nil, invoice.Id, neo4jrepository.InvoiceUpdateFields{
			Status:       neo4jenum.InvoiceStatusDue,
			UpdateStatus: true,
		})
		if err != nil {
			tracing.TraceErr(span, err)
			s.log.Errorf("Error updating invoice %s to due: %s", invoice.Id, err.Error())
			return // stop processing
		}
	}

}
