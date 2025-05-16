package service

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
	model2 "github.com/customeros/customeros/packages/server/customer-os-common-module/model"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"
	neo4jenum "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/enum"
	neo4jrepository "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/repository"

	"github.com/customeros/customeros/packages/server/customer-os-webhooks/errors"
	"github.com/customeros/customeros/packages/server/customer-os-webhooks/model"
	"github.com/customeros/customeros/packages/server/customer-os-webhooks/repository"
)

type InvoiceService interface {
	SyncInvoices(ctx context.Context, contacts []model.InvoiceData) (SyncResult, error)
}

type invoiceService struct {
	log          logger.Logger
	repositories *repository.Repositories
	services     *Services
	maxWorkers   int
}

func NewInvoiceService(log logger.Logger, repositories *repository.Repositories, services *Services) InvoiceService {
	return &invoiceService{
		log:          log,
		repositories: repositories,
		services:     services,
		maxWorkers:   services.Cfg.App.ConcurrencyConfig.InvoiceSyncConcurrency,
	}
}

func (s *invoiceService) SyncInvoices(ctx context.Context, invoices []model.InvoiceData) (SyncResult, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "InvoiceService.SyncInvoices")
	defer spans.Finish()
	spans.LogKV("num of invoices", len(invoices))

	if !s.services.TenantService.Exists(ctx, common.GetTenantFromContext(ctx)) {
		s.log.Errorf("tenant {%s} does not exist", common.GetTenantFromContext(ctx))
		spans.TraceError(errors.ErrTenantNotValid)
		return SyncResult{}, errors.ErrTenantNotValid
	}

	// pre-validate invoices input before syncing
	for _, invoice := range invoices {
		// sync by id or external system is required
		if invoice.Id == "" && invoice.ExternalSystem == "" {
			spans.TraceError(errors.ErrMissingExternalSystem)
			return SyncResult{}, errors.ErrMissingExternalSystem
		}
		if invoice.ExternalSystem != "" {
			if !neo4jentity.IsValidDataSource(strings.ToLower(invoice.ExternalSystem)) {
				spans.TraceError(errors.ErrExternalSystemNotAccepted)
				spans.LogKV("externalSystem", invoice.ExternalSystem)
				return SyncResult{}, errors.ErrExternalSystemNotAccepted
			}
		}
	}

	// Create a wait group to wait for all workers to finish
	var wg sync.WaitGroup
	// Create a channel to control the number of concurrent workers
	workerLimit := make(chan struct{}, s.maxWorkers)

	syncMutex := &sync.Mutex{}
	statusesMutex := &sync.Mutex{}
	syncDate := utils.Now()
	var statuses []SyncStatus

	// Sync all invoices concurrently
	for _, invoiceData := range invoices {
		// Check if context is cancelled
		select {
		case <-ctx.Done():
			return SyncResult{}, ctx.Err()
		default:
		}

		// Acquire a worker slot
		workerLimit <- struct{}{}
		wg.Add(1)

		go func(invoiceData model.InvoiceData) {
			defer wg.Done()
			defer func() {
				// Release the worker slot when done
				<-workerLimit
			}()

			result := s.syncInvoice(ctx, syncMutex, invoiceData, syncDate)
			statusesMutex.Lock()
			statuses = append(statuses, result)
			statusesMutex.Unlock()
		}(invoiceData)
	}
	// Wait for all workers to finish
	wg.Wait()

	s.services.SyncStatusService.SaveSyncResults(ctx, common.GetTenantFromContext(ctx), invoices[0].ExternalSystem,
		invoices[0].AppSource, "invoice", syncDate, statuses)

	return s.services.SyncStatusService.PrepareSyncResult(statuses), nil
}

func (s *invoiceService) syncInvoice(ctx context.Context, syncMutex *sync.Mutex, invoiceInput model.InvoiceData, syncDate time.Time) SyncStatus {
	spans, ctx := telemetry.StartServiceSpan(ctx, "InvoiceService.syncInvoice")
	defer spans.Finish()
	spans.TagString(telemetry.SpanTagExternalSystem, invoiceInput.ExternalSystem)
	spans.TagString(telemetry.SpanTagExternalId, invoiceInput.ExternalId)
	spans.LogObjectAsJson("syncDate", syncDate)
	spans.LogObjectAsJson("invoiceInput", invoiceInput)

	tenant := common.GetTenantFromContext(ctx)
	failedSync := false
	reason := ""

	invoiceInput.Normalize()

	err := s.services.ExternalSystemService.MergeExternalSystem(ctx, tenant, invoiceInput.ExternalSystem)
	if err != nil {
		spans.TraceError(err)
		reason = fmt.Sprintf("failed merging external system %s for tenant %s :%s", invoiceInput.ExternalSystem, tenant, err.Error())
		s.log.Error(reason)
		spans.LogKV("result.status", "failed")
		return NewFailedSyncStatus(reason)
	}

	// Check if invoice sync should be skipped
	if invoiceInput.Skip {
		spans.LogKV("result.status", "skipped")
		return NewSkippedSyncStatus(invoiceInput.SkipReason)
	}
	if invoiceInput.ExternalId == "" && invoiceInput.Id == "" {
		reason = fmt.Sprintf("id and external id are empty for invoice, tenant %s", tenant)
		s.log.Warnf("Skip issue sync: %v", reason)
		spans.LogKV("result.status", "skipped")
		return NewSkippedSyncStatus(reason)
	}

	syncMutex.Lock()
	defer syncMutex.Unlock()

	var invoiceId string
	if invoiceInput.Id != "" {
		exists, err := s.repositories.Neo4jRepositories.CommonReadRepository.ExistsById(ctx, tenant, invoiceInput.Id, model2.NodeLabelInvoice)
		if err != nil {
			failedSync = true
			spans.TraceError(err)
			reason = fmt.Sprintf("failed checking if invoice with id %s exists for tenant %s :%s", invoiceInput.Id, tenant, err.Error())
			s.log.Error(reason)
			return NewFailedSyncStatus(reason)
		}
		if exists {
			invoiceId = invoiceInput.Id
		}
	}

	matchingInvoiceExists := invoiceId != ""
	spans.LogKV("found matching invoice", matchingInvoiceExists)
	if invoiceInput.UpdateOnly && !matchingInvoiceExists {
		reason = fmt.Sprintf("update only is true and matching invoice does not exist for tenant %s", tenant)
		s.log.Warnf("Skip invoice sync: %v", reason)
		spans.LogKV("result.status", "skipped")
		return NewSkippedSyncStatus(reason)
	}

	if invoiceInput.UpdateOnly || matchingInvoiceExists {
		// Update invoice
		invoiceUpdateFields := neo4jrepository.InvoiceUpdateFields{}

		if invoiceInput.Status != "" {
			switch strings.ToLower(invoiceInput.Status) {
			case "draft":
				invoiceUpdateFields.Status = neo4jenum.InvoiceStatusInitialized
			case "paid":
				invoiceUpdateFields.Status = neo4jenum.InvoiceStatusPaid
			case "due":
				invoiceUpdateFields.Status = neo4jenum.InvoiceStatusDue
			case "payment_processing":
				invoiceUpdateFields.Status = neo4jenum.InvoiceStatusPaymentProcessing
			}
			invoiceUpdateFields.UpdateStatus = true
		}
		if invoiceInput.PaymentLink != "" {
			invoiceUpdateFields.PaymentLink = invoiceInput.PaymentLink
			invoiceUpdateFields.UpdatePaymentLink = true

			if invoiceInput.PaymentLinkValidHours != "" {
				paymentLinkValidHours, err := strconv.Atoi(invoiceInput.PaymentLinkValidHours)
				if err != nil {
					// default to 24 hours
					paymentLinkValidHours = 24
				}
				if paymentLinkValidHours > 1 {
					// reduce 30 minutes from the valid until time to make sure the payment link is still valid
					validUntil := utils.Now().Add(time.Minute * time.Duration(paymentLinkValidHours*60-30))
					invoiceUpdateFields.PaymentLinkValidUntil = &validUntil
				} else if paymentLinkValidHours == 1 {
					// reduce 15 minutes from the valid until time to make sure the payment link is still valid
					validUntil := utils.Now().Add(time.Minute * 45)
					invoiceUpdateFields.PaymentLinkValidUntil = &validUntil
				}
			}
		}

		err = s.services.CommonServices.InvoiceService.UpdateInvoice(ctx, nil, invoiceId, invoiceUpdateFields)
		if err != nil {
			failedSync = true
			spans.TraceError(err)
			reason = fmt.Sprintf("failed updating invoice with id %s for tenant %s :%s", invoiceId, tenant, err)
			s.log.Error(reason)
		}
	}

	spans.LogKV("failedSync", failedSync)
	if failedSync {
		spans.LogKV("result.status", "failed")
		return NewFailedSyncStatus(reason)
	}
	spans.LogKV("result.status", "success")
	return NewSuccessfulSyncStatus()
}
