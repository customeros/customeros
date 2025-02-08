package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/data_fields"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/dto"
	"net/http"
	"time"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/data"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/model"
	commonService "github.com/customeros/customeros/packages/server/customer-os-common-module/services"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"
	neo4jenum "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/enum"
	neo4jmapper "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/mapper"
	neo4jrepository "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/repository"
	"github.com/opentracing/opentracing-go"
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

type InvoiceFinalizedEventBody struct {
	Tenant                       string `json:"tenant"`
	Currency                     string `json:"currency"`
	AmountInSmallestCurrencyUnit int64  `json:"amountInSmallestCurrencyUnit"`
	InvoiceId                    string `json:"invoiceId"`
	InvoiceDescription           string `json:"invoiceDescription"`
	CustomerOsId                 string `json:"customerOsId"`
	Status                       string `json:"status"`
	CustomerEmail                string `json:"customerEmail"`
	CustomerName                 string `json:"customerName"`
	PrimaryStripeCustomerId      string `json:"stripeCustomerId"`
	Pay                          struct {
		PayAutomatically      bool `json:"payAutomatically"`
		CanPayWithCard        bool `json:"canPayWithCard"`
		CanPayWithDirectDebit bool `json:"canPayWithDirectDebit"`
	} `json:"pay"`
}

type InvoiceService interface {
	GenerateCycleInvoices()
	GenerateOffCycleInvoices()
	SendPayNotifications()
	SendRemindNotifications()
	GenerateInvoicePaymentLinks()
	CleanupInvoices()
	GenerateNextPreviewInvoices()
	AdjustInvoiceStatus()
	SendInvoiceFinalizedEvent()
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

func (s *invoiceService) GenerateCycleInvoices() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Cancel context on exit

	if s.cfg.App.ProcessConfig.CycleInvoicingEnabled == false {
		s.log.Infof("Cycle invoicing is disabled, stopping")
		return
	}

	span, ctx := tracing.StartTracerSpan(ctx, "InvoiceService.GenerateCycleInvoices")
	defer span.Finish()
	tracing.TagComponentCronJob(span)

	referenceTime := utils.Now()
	dryRun := false
	preview := false

	limit := 0

	for {
		select {
		case <-ctx.Done():
			s.log.Infof("Context cancelled, stopping")
			return
		default:
			// continue as normal
		}

		records, err := s.repositories.Neo4jRepositories.ContractReadRepository.GetContractsToGenerateCycleInvoices(ctx, referenceTime, s.cfg.App.ProcessConfig.DelayGenerateCycleInvoiceInMinutes, limit)
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
			recordSpan, innerCtx := tracing.StartTracerSpan(innerCtx, "InvoiceService.GenerateCycleInvoices.Record")
			defer recordSpan.Finish()
			tracing.TagTenant(recordSpan, record.Tenant)
			contract := neo4jmapper.MapDbNodeToContractEntity(record.Node)
			tenant := record.Tenant

			// mark invoicing requested to prevent double processing of the record
			err = s.repositories.Neo4jRepositories.ContractWriteRepository.MarkCycleInvoicingRequested(innerCtx, tenant, contract.Id, utils.Now())
			if err != nil {
				tracing.TraceErr(recordSpan, errors.Wrap(err, "Error marking invoicing requested"))
				s.log.Errorf("Error marking invoicing started for contract %s: %s", contract.Id, err.Error())
				return
			}

			// check if tenant has agent for invoicing enabled and configured
			invoicingAgents, err := s.commonServices.PostgresRepositories.AgentRepository.GetActiveConfiguredAgentsByTypesCrossTenant(innerCtx, []enum.AgentType{enum.AgentCashflowGuardian})
			if err != nil {
				tracing.TraceErr(span, err)
				s.log.Errorf("Error getting invoicing agents: %v", err)
				continue
			}
			if len(invoicingAgents) == 0 {
				s.log.Infof("No invoicing agents configured for tenant %s", tenant)
				continue
			}

			event := dto.IntentDetected{
				EventName: enum.EventIntentSignal,
				//IntentType: enum.IntentGenerateCycleInvoice,
				ContractID: contract.Id,
				DryRun:     dryRun,
				Preview:    preview,
			}

			err = s.commonServices.Events.Publisher.PublishFanoutEvent(innerCtx, contract.Id, model.INTENT_SIGNAL, &event)
			if err != nil {
				tracing.TraceErr(span, errors.Wrap(err, "error publishing icp check event"))
				s.log.Errorf("Error publishing icp check event: %v", err)
				continue
			}
		}

		// sleep for async processing, then check again
		time.Sleep(10 * time.Second)
	}
}

func (s *invoiceService) GenerateNextPreviewInvoices() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Cancel context on exit

	span, ctx := tracing.StartTracerSpan(ctx, "InvoiceService.GenerateNextPreviewInvoices")
	defer span.Finish()
	tracing.TagComponentCronJob(span)

	referenceTime := utils.Now()

	for {
		select {
		case <-ctx.Done():
			s.log.Infof("Context cancelled, stopping")
			return
		default:
			// continue as normal
		}

		records, err := s.repositories.Neo4jRepositories.ContractReadRepository.GetContractsToGenerateNextScheduledInvoices(ctx, referenceTime, 10)
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

			dataFields := data_fields.InvoiceFields{
				DryRun:  true,
				Preview: true,
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

func (s *invoiceService) SendPayNotifications() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Cancel context on exit

	span, ctx := tracing.StartTracerSpan(ctx, "InvoiceService.SendPayNotifications")
	defer span.Finish()
	tracing.TagComponentCronJob(span)

	referenceTime := utils.Now()

	for {
		select {
		case <-ctx.Done():
			s.log.Infof("Context cancelled, stopping")
			return
		default:
			// continue as normal
		}

		records, err := s.repositories.Neo4jRepositories.InvoiceReadRepository.GetInvoicesForPayNotifications(
			ctx, s.cfg.App.ProcessConfig.DelaySendPayInvoiceNotificationInMinutes, s.cfg.App.ProcessConfig.RetrySendPayInvoiceNotificationDays, referenceTime)
		if err != nil {
			tracing.TraceErr(span, err)
			s.log.Errorf("Error getting invoices for pay notifications: %v", err)
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
			innerCtx := common.WithCustomContext(ctx, &common.CustomContext{
				Tenant:    tenant,
				AppSource: constants.AppSourceDataUpkeeper,
			})

			err = s.commonServices.InvoiceService.SendPayInvoiceNotification(innerCtx, invoice.Id)
			if err != nil {
				tracing.TraceErr(span, err)
				s.log.Errorf("Error sending pay notification for invoice %s: %s", invoice.Id, err.Error())
			}
		}
		// sleep for async processing, then check again
		time.Sleep(5 * time.Second)
	}
}

func (s *invoiceService) SendRemindNotifications() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Cancel context on exit

	span, ctx := tracing.StartTracerSpan(ctx, "InvoiceService.SendRemindNotifications")
	defer span.Finish()
	tracing.TagComponentCronJob(span)

	referenceTime := utils.Now()
	limit := 100
	overdueDays := 15

	for {
		select {
		case <-ctx.Done():
			s.log.Infof("Context cancelled, stopping")
			return
		default:
			// continue as normal
		}

		records, err := s.repositories.Neo4jRepositories.InvoiceReadRepository.GetInvoicesForRemindNotifications(ctx, referenceTime, overdueDays, limit)
		if err != nil {
			tracing.TraceErr(span, err)
			s.log.Errorf("Error getting invoices for pay notifications: %v", err)
			return
		}

		// no invoices found
		if len(records) == 0 {
			return
		}

		// process records
		for _, record := range records {
			invoice := neo4jmapper.MapDbNodeToInvoiceEntity(record.Node)
			innerCtx := common.WithCustomContext(ctx, &common.CustomContext{
				Tenant:    record.Tenant,
				AppSource: constants.AppSourceDataUpkeeper,
			})

			err = s.commonServices.InvoiceService.SendPayReminderInvoiceNotification(innerCtx, invoice.Id)
			if err != nil {
				tracing.TraceErr(span, err)
				s.log.Errorf("Error sending pay notification for invoice %s: %s", invoice.Id, err.Error())
			}
		}
		// sleep for async processing, then check again
		time.Sleep(1 * time.Second)
	}
}

func (s *invoiceService) GenerateOffCycleInvoices() {
	return

	//ctx, cancel := context.WithCancel(context.Background())
	//defer cancel() // Cancel context on exit

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

func (s *invoiceService) GenerateInvoicePaymentLinks() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Cancel context on exit

	span, ctx := tracing.StartTracerSpan(ctx, "InvoiceService.GenerateInvoicePaymentLinks")
	defer span.Finish()
	tracing.TagComponentCronJob(span)

	if s.cfg.App.EventNotifications.IntegrationAppEventWebhookUrls.GeneratePaymentLinkUrl == "" {
		err := errors.New("GeneratePaymentLinkUrl is not configured")
		tracing.TraceErr(span, err)
		s.log.Error(err.Error())
		return
	}

	referenceTime := utils.Now()
	limit := 100

	for {
		select {
		case <-ctx.Done():
			s.log.Infof("Context cancelled, stopping")
			return
		default:
			// continue as normal
		}

		records, err := s.repositories.Neo4jRepositories.InvoiceReadRepository.GetInvoicesForPaymentLinkRequest(
			ctx, s.cfg.App.ProcessConfig.DelayRequestPaymentLinkInMinutes, s.cfg.App.ProcessConfig.RequestPaymentLinkLookBackWindowInDays, referenceTime, limit)
		if err != nil {
			tracing.TraceErr(span, err)
			s.log.Errorf("Error getting invoices for payment links generation: %v", err)
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
				AppSource: constants.AppSourceDataUpkeeper,
			})
			invoice := neo4jmapper.MapDbNodeToInvoiceEntity(record.Node)
			tenant := record.Tenant

			// get contract linked to invoice
			contractDbNode, err := s.repositories.Neo4jRepositories.ContractReadRepository.GetContractForInvoice(innerCtx, tenant, invoice.Id)
			if err != nil {
				tracing.TraceErr(span, err)
				s.log.Errorf("Error getting contract for invoice %s: %s", invoice.Id, err.Error())
			}
			contractEntity := neo4jentity.ContractEntity{}
			if contractDbNode != nil {
				contractEntity = *neo4jmapper.MapDbNodeToContractEntity(contractDbNode)
			}

			// get organization for invoice
			organizationDbNode, err := s.repositories.Neo4jRepositories.OrganizationReadRepository.GetOrganizationByInvoiceId(innerCtx, tenant, invoice.Id)
			if err != nil {
				tracing.TraceErr(span, err)
				s.log.Errorf("Error getting organization for invoice %s: %s", invoice.Id, err.Error())
			}
			organizationEntity := neo4jentity.OrganizationEntity{}
			if organizationDbNode != nil {
				organizationEntity = *neo4jmapper.MapDbNodeToOrganizationEntity(organizationDbNode)
			}

			// convert amount to the smallest currency unit
			amountInSmallestCurrencyUnit, err := data.InSmallestCurrencyUnit(invoice.Currency.String(), invoice.TotalAmount)
			if err != nil {
				tracing.TraceErr(span, err)
			}

			// mark payment link request first, before sending the event
			err = s.repositories.Neo4jRepositories.InvoiceWriteRepository.MarkPaymentLinkRequested(innerCtx, tenant, invoice.Id)
			if err != nil {
				tracing.TraceErr(span, err)
				s.log.Errorf("Error marking payment link requested for invoice %s: %s", invoice.Id, err.Error())
			}

			primaryStripeCustomerId, err := s.commonServices.ExternalSystemService.GetPrimaryExternalId(innerCtx, enum.SourceStripe.String(), organizationEntity.ID, model.ORGANIZATION)
			if err != nil {
				tracing.TraceErr(span, err)
				s.log.Errorf("Error getting primary stripe customer id for contract %s: %s", contractEntity.Id, err.Error())
			}

			requestBody := GeneratePaymentLinkEventBody{
				Tenant:                       tenant,
				Currency:                     invoice.Currency.String(),
				AmountInSmallestCurrencyUnit: amountInSmallestCurrencyUnit,
				InvoiceId:                    invoice.Id,
				InvoiceDescription:           fmt.Sprintf("Invoice %s", invoice.Number),
				CustomerEmail:                contractEntity.InvoiceEmail,
				PrimaryStripeCustomerId:      primaryStripeCustomerId,
			}

			// Convert the request body to JSON
			requestBodyJSON, err := json.Marshal(requestBody)
			if err != nil {
				tracing.TraceErr(span, err)
				s.log.Errorf("error encoding JSON: %v", err)
				continue
			}

			// Create an HTTP client
			client := &http.Client{}

			// Create a POST request with headers and body
			req, err := http.NewRequest("POST", s.cfg.App.EventNotifications.IntegrationAppEventWebhookUrls.GeneratePaymentLinkUrl, bytes.NewBuffer(requestBodyJSON))
			if err != nil {
				tracing.TraceErr(span, err)
				continue
			}

			// Set the content type header
			req.Header.Set("Content-Type", "application/json")

			// Send the POST request
			resp, err := client.Do(req)
			if err != nil {
				tracing.TraceErr(span, err)
				s.log.Errorf("error sending request: %v", err)
				continue
			}
			defer resp.Body.Close()

			// Check the response status code
			if resp.StatusCode != http.StatusOK {
				tracing.TraceErr(span, fmt.Errorf("request failed with status code: %s", resp.Status))
				s.log.Errorf("request failed with status code: %s", resp.Status)
			}
		}
	}
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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Cancel context on exit

	span, ctx := tracing.StartTracerSpan(ctx, "InvoiceService.AdjustInvoiceStatus")
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

		recordsForOverdueInvoices, err := s.repositories.Neo4jRepositories.InvoiceReadRepository.GetInvoicesForOverdue(ctx)
		if err != nil {
			tracing.TraceErr(span, err)
			s.log.Errorf("Error getting invoices for overdue: %v", err)
			return
		}

		for _, record := range recordsForOverdueInvoices {
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

		recordsForOnHoldInvoices, err := s.repositories.Neo4jRepositories.InvoiceReadRepository.GetInvoicesForOnHold(ctx)
		if err != nil {
			tracing.TraceErr(span, err)
			s.log.Errorf("Error getting invoices for on hold: %v", err)
			return
		}

		for _, record := range recordsForOnHoldInvoices {
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

		recordsForOnScheduledInvoices, err := s.repositories.Neo4jRepositories.InvoiceReadRepository.GetInvoicesForScheduled(ctx)
		if err != nil {
			tracing.TraceErr(span, err)
			s.log.Errorf("Error getting invoices for scheduled: %v", err)
			return
		}

		for _, record := range recordsForOnScheduledInvoices {
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

		if len(recordsForOverdueInvoices) == 0 && len(recordsForOnHoldInvoices) == 0 && len(recordsForOnScheduledInvoices) == 0 {
			return
		}

		// sleep for async processing, then check again
		time.Sleep(10 * time.Second)
	}
}

func (s *invoiceService) SendInvoiceFinalizedEvent() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	span, ctx := tracing.StartTracerSpan(ctx, "InvoiceService.SendInvoiceFinalizedEvent")
	defer span.Finish()
	tracing.TagComponentCronJob(span)

	if s.cfg.App.EventNotifications.IntegrationAppEventWebhookUrls.InvoiceFinalizedUrl == "" {
		err := errors.New("InvoiceFinalizedUrl is not configured")
		tracing.TraceErr(span, err)
		s.log.Error(err.Error())
		return
	}

	referenceTime := utils.Now()
	limit := 10

	for {
		select {
		case <-ctx.Done():
			s.log.Infof("Context cancelled, stopping")
			return
		default:
			// continue as normal
		}

		records, err := s.repositories.Neo4jRepositories.InvoiceReadRepository.GetReadyInvoicesForFinalizedEvent(ctx, s.cfg.App.ProcessConfig.DelayAutoPayInvoiceInMinutes, referenceTime, limit)
		if err != nil {
			tracing.TraceErr(span, err)
			s.log.Errorf("Error getting invoices for finalized event: %v", err)
			return
		}

		// no invoices found
		if len(records) == 0 {
			return
		}

		// process records
		for _, record := range records {
			invoiceEntity := neo4jmapper.MapDbNodeToInvoiceEntity(record.Node)
			tenant := record.Tenant

			if invoiceEntity.InvoiceInternalFields.InvoiceFinalizedSentAt == nil {
				err = s.integrationAppInvoiceFinalizedWebhook(ctx, tenant, *invoiceEntity)
				if err != nil {
					tracing.TraceErr(span, errors.Wrap(err, "error invoking invoice finalized webhook"))
					s.log.Errorf("error invoking invoice ready webhook for invoice %s: %s", invoiceEntity.Id, err.Error())
				}

				err = s.repositories.Neo4jRepositories.CommonWriteRepository.UpdateTimeProperty(ctx, record.Tenant, model.NodeLabelInvoice, invoiceEntity.Id, string(neo4jentity.InvoicePropertyInvoiceFinalizedEventSentAt), utils.NowPtr())
				if err != nil {
					tracing.TraceErr(span, errors.Wrap(err, "error updating invoice finalized sent at"))
					s.log.Errorf("Error updating invoice finalized sent at for invoice %s: %s", invoiceEntity.Id, err.Error)
					return
				}
			}
		}
		time.Sleep(2 * time.Second)
	}
}

func (s *invoiceService) integrationAppInvoiceFinalizedWebhook(ctx context.Context, tenant string, invoice neo4jentity.InvoiceEntity) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceService.integrationAppInvoiceFinalizedWebhook")
	defer span.Finish()
	tracing.TagTenant(span, tenant)
	tracing.LogObjectAsJson(span, "invoice", invoice)

	innerCtx := common.WithCustomContext(ctx, &common.CustomContext{
		Tenant:    tenant,
		AppSource: constants.AppSourceDataUpkeeper,
	})

	if s.cfg.App.EventNotifications.IntegrationAppEventWebhookUrls.InvoiceFinalizedUrl == "" {
		return nil
	}

	// get organization linked to invoice
	organizationDbNode, err := s.repositories.Neo4jRepositories.OrganizationReadRepository.GetOrganizationByInvoiceId(innerCtx, tenant, invoice.Id)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error getting organization for invoice %s: %s", invoice.Id, err.Error())
		return err
	}
	organizationEntity := neo4jentity.OrganizationEntity{}
	if organizationDbNode != nil {
		organizationEntity = *neo4jmapper.MapDbNodeToOrganizationEntity(organizationDbNode)
	}

	// get contract linked to invoice
	contractDbNode, err := s.repositories.Neo4jRepositories.ContractReadRepository.GetContractForInvoice(innerCtx, tenant, invoice.Id)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error getting contract for invoice %s: %s", invoice.Id, err.Error())
		return err
	}
	contractEntity := neo4jentity.ContractEntity{}
	if contractDbNode != nil {
		contractEntity = *neo4jmapper.MapDbNodeToContractEntity(contractDbNode)
	}

	// convert amount to the smallest currency unit
	amountInSmallestCurrencyUnit, err := data.InSmallestCurrencyUnit(invoice.Currency.String(), invoice.TotalAmount)
	if err != nil {
		return fmt.Errorf("error converting amount to smallest currency unit: %v", err.Error())
	}

	primaryStripeCustomerId, err := s.commonServices.ExternalSystemService.GetPrimaryExternalId(innerCtx, enum.SourceStripe.String(), organizationEntity.ID, model.ORGANIZATION)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error getting primary stripe customer id for contract %s: %s", contractEntity.Id, err.Error())
	}

	requestBody := InvoiceFinalizedEventBody{
		Tenant:                       tenant,
		Currency:                     invoice.Currency.String(),
		AmountInSmallestCurrencyUnit: amountInSmallestCurrencyUnit,
		InvoiceId:                    invoice.Id,
		InvoiceDescription:           fmt.Sprintf("Invoice %s", invoice.Number),
		CustomerOsId:                 organizationEntity.CustomerOsId,
		Status:                       invoice.Status.String(),
		CustomerEmail:                contractEntity.InvoiceEmail,
		CustomerName:                 utils.GetReadableNameFromEmail(contractEntity.InvoiceEmail),
		PrimaryStripeCustomerId:      primaryStripeCustomerId,
		Pay: struct {
			PayAutomatically      bool `json:"payAutomatically"`
			CanPayWithCard        bool `json:"canPayWithCard"`
			CanPayWithDirectDebit bool `json:"canPayWithDirectDebit"`
		}{
			PayAutomatically:      contractEntity.PayAutomatically && (invoice.Status == neo4jenum.InvoiceStatusDue || invoice.Status == neo4jenum.InvoiceStatusOverdue),
			CanPayWithCard:        true,
			CanPayWithDirectDebit: true,
		},
	}

	// Convert the request body to JSON
	requestBodyJSON, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("error encoding JSON: %v", err)
	}

	// Create an HTTP client
	client := &http.Client{}

	// Create a POST request with headers and body
	req, err := http.NewRequest("POST", s.cfg.App.EventNotifications.IntegrationAppEventWebhookUrls.InvoiceFinalizedUrl, bytes.NewBuffer(requestBodyJSON))
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	// Set the content type header
	req.Header.Set("Content-Type", "application/json")

	// Send the POST request
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("request failed with status code: %s", resp.Status)
	}

	return nil
}
