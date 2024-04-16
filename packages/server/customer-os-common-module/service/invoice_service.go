package service

import (
	"context"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	commonpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/common"
	invoicepb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/invoice"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"time"
)

type InvoiceService interface {
	GetById(ctx context.Context, tenant, invoiceId string) (*neo4jentity.InvoiceEntity, error)
	GetInvoiceLinesForInvoices(ctx context.Context, invoiceIds []string) (*neo4jentity.InvoiceLineEntities, error)
	GetInvoicesForContracts(ctx context.Context, contractIds []string) (*neo4jentity.InvoiceEntities, error)
	SimulateInvoice(ctx context.Context, invoiceData *SimulateInvoiceRequestData) ([]*SimulateInvoiceResponseData, error)
	NextInvoiceDryRun(ctx context.Context, contractId, appSource string) (string, error)
	PayInvoice(ctx context.Context, invoiceId, appSource string) error
	VoidInvoice(ctx context.Context, invoiceId, appSource string) error

	FillCycleInvoice(ctx context.Context, tenant, contractId string, invoiceEntity *neo4jentity.InvoiceEntity, sliEntities neo4jentity.ServiceLineItemEntities) (*neo4jentity.InvoiceEntity, []*invoicepb.InvoiceLine, error)
}
type invoiceService struct {
	log      logger.Logger
	services *Services
}

func NewInvoiceService(log logger.Logger, services *Services) InvoiceService {
	return &invoiceService{
		log:      log,
		services: services,
	}
}

type SimulateInvoiceRequestData struct {
	ContractId   string
	ServiceLines []SimulateInvoiceRequestServiceLineData
}
type SimulateInvoiceRequestServiceLineData struct {
	ServiceLineItemID *string
	ParentID          *string
	Description       string
	Comments          string
	BillingCycle      enum.BilledType
	Price             float64
	Quantity          int64
	ServiceStarted    time.Time
	TaxRate           *float64
}

type SimulateInvoiceResponseData struct {
	Invoice *neo4jentity.InvoiceEntity
	Lines   []*neo4jentity.InvoiceLineEntity
}

func (s *invoiceService) GetById(ctx context.Context, tenant, invoiceId string) (*neo4jentity.InvoiceEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceService.GetById")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("invoiceId", invoiceId))

	if invoiceDbNode, err := s.services.Neo4jRepositories.InvoiceReadRepository.GetInvoiceById(ctx, tenant, invoiceId); err != nil {
		tracing.TraceErr(span, err)
		wrappedErr := errors.Wrap(err, fmt.Sprintf("Invoice with id {%s} not found", invoiceId))
		return nil, wrappedErr
	} else {
		return mapper.MapDbNodeToInvoiceEntity(invoiceDbNode), nil
	}
}

func (s *invoiceService) GetInvoiceLinesForInvoices(ctx context.Context, invoiceIds []string) (*neo4jentity.InvoiceLineEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceService.GetInvoiceLinesForInvoices")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("invoiceIds", invoiceIds))

	invoiceLines, err := s.services.Neo4jRepositories.InvoiceLineReadRepository.GetAllForInvoices(ctx, common.GetTenantFromContext(ctx), invoiceIds)
	if err != nil {
		return nil, err
	}
	invoiceLineEntities := make(neo4jentity.InvoiceLineEntities, 0, len(invoiceLines))
	for _, v := range invoiceLines {
		invoiceLineEntity := mapper.MapDbNodeToInvoiceLineEntity(v.Node)
		invoiceLineEntity.DataloaderKey = v.LinkedNodeId
		invoiceLineEntities = append(invoiceLineEntities, *invoiceLineEntity)
	}
	return &invoiceLineEntities, nil
}

func (s *invoiceService) GetInvoicesForContracts(ctx context.Context, contractIds []string) (*neo4jentity.InvoiceEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceService.GetInvoicesForContracts")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("contractIds", contractIds))

	invoices, err := s.services.Neo4jRepositories.InvoiceReadRepository.GetAllForContracts(ctx, common.GetTenantFromContext(ctx), contractIds)
	if err != nil {
		return nil, err
	}
	invoiceEntities := make(neo4jentity.InvoiceEntities, 0, len(invoices))
	for _, v := range invoices {
		invoiceEntity := mapper.MapDbNodeToInvoiceEntity(v.Node)
		invoiceEntity.DataloaderKey = v.LinkedNodeId
		invoiceEntities = append(invoiceEntities, *invoiceEntity)
	}
	return &invoiceEntities, nil
}

func (s *invoiceService) SimulateInvoice(ctx context.Context, invoiceData *SimulateInvoiceRequestData) ([]*SimulateInvoiceResponseData, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceService.SimulateInvoice")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("invoiceData", invoiceData))

	if invoiceData.ServiceLines == nil {
		err := fmt.Errorf("no invoice lines to simulate")
		tracing.TraceErr(span, err)
		return nil, err
	}

	invoiceEntity := &neo4jentity.InvoiceEntity{}
	invoiceLines := []*invoicepb.InvoiceLine{}

	tenantSettings, err := s.services.TenantService.GetTenantSettings(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	contractEntity, err := s.services.ContractService.GetById(ctx, invoiceData.ContractId)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	var nextPreviewInvoiceEntity *neo4jentity.InvoiceEntity
	nextPreviewInvoiceNode, err := s.services.Neo4jRepositories.InvoiceReadRepository.GetFirstPreviewFilledInvoice(ctx, common.GetTenantFromContext(ctx), invoiceData.ContractId)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	if nextPreviewInvoiceNode != nil {
		nextPreviewInvoiceEntity = mapper.MapDbNodeToInvoiceEntity(nextPreviewInvoiceNode)
		invoiceEntity.Number = nextPreviewInvoiceEntity.Number
	} else {
		invoiceEntity.Number = "" // todo
	}

	var invoicePeriodStart, invoicePeriodEnd time.Time
	if contractEntity.NextInvoiceDate != nil {
		invoicePeriodStart = *contractEntity.NextInvoiceDate
	} else {
		invoicePeriodStart = *contractEntity.InvoicingStartDate
	}
	invoicePeriodEnd = calculateInvoiceCycleEnd(invoicePeriodStart, contractEntity.BillingCycle)

	invoiceEntity.OffCycle = false
	invoiceEntity.Postpaid = tenantSettings.InvoicingPostpaid
	invoiceEntity.BillingCycle = contractEntity.BillingCycle
	invoiceEntity.PeriodStartDate = invoicePeriodStart
	invoiceEntity.PeriodEndDate = invoicePeriodEnd
	invoiceEntity.Note = contractEntity.InvoiceNote

	if nextPreviewInvoiceNode != nil {
		invoiceEntity.IssuedDate = nextPreviewInvoiceEntity.IssuedDate
		invoiceEntity.DueDate = nextPreviewInvoiceEntity.DueDate
	} else {
		invoiceEntity.IssuedDate = invoicePeriodStart
		invoiceEntity.DueDate = invoicePeriodEnd.AddDate(0, 0, int(contractEntity.DueDays))
	}

	if contractEntity.Currency != "" {
		invoiceEntity.Currency = contractEntity.Currency
	} else {
		invoiceEntity.Currency = tenantSettings.BaseCurrency
	}

	sliEntities := neo4jentity.ServiceLineItemEntities{}
	for _, sliData := range invoiceData.ServiceLines {
		sliEntity := neo4jentity.ServiceLineItemEntity{
			ID:        utils.StringPtrFirstNonEmpty(sliData.ServiceLineItemID),
			ParentID:  utils.StringPtrFirstNonEmpty(sliData.ParentID),
			Name:      sliData.Description,
			Comments:  sliData.Comments,
			Billed:    sliData.BillingCycle,
			Price:     sliData.Price,
			Quantity:  sliData.Quantity,
			StartedAt: sliData.ServiceStarted,
			EndedAt:   nil,
			VatRate:   utils.IfNotNilFloat64(sliData.TaxRate),
		}
		sliEntities = append(sliEntities, sliEntity)
	}

	invoiceEntity, invoiceLines, err = s.FillCycleInvoice(ctx, common.GetTenantFromContext(ctx), invoiceData.ContractId, invoiceEntity, sliEntities)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	response := []*SimulateInvoiceResponseData{}
	response = append(response, &SimulateInvoiceResponseData{
		Invoice: invoiceEntity,
		Lines:   []*neo4jentity.InvoiceLineEntity{},
	})
	for _, line := range invoiceLines {
		invoiceLineEntity := &neo4jentity.InvoiceLineEntity{
			Id:                      line.ServiceLineItemId,
			ServiceLineItemParentId: line.ServiceLineItemParentId,
			Name:                    line.Name,
			Price:                   line.Price,
			Quantity:                line.Quantity,
			Amount:                  line.Amount,
			TotalAmount:             line.Total,
			Vat:                     line.Vat,
		}
		response[0].Lines = append(response[0].Lines, invoiceLineEntity)
	}

	return response, nil
}

func (s *invoiceService) NextInvoiceDryRun(ctx context.Context, contractId, appSource string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceService.NextInvoiceDryRun")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("contractId", contractId))

	tenant := common.GetTenantFromContext(ctx)
	now := time.Now()

	contract, err := s.services.ContractService.GetById(ctx, contractId)
	if err != nil {
		tracing.TraceErr(span, err)
		return "", err
	}

	var invoicePeriodStart, invoicePeriodEnd time.Time
	if contract.NextInvoiceDate != nil {
		invoicePeriodStart = *contract.NextInvoiceDate
	} else if contract.InvoicingStartDate != nil {
		invoicePeriodStart = *contract.InvoicingStartDate
	} else {
		err = fmt.Errorf("contract has no next invoice date or invoicing start date")
		tracing.TraceErr(span, err)
		return "", err
	}
	invoicePeriodEnd = calculateInvoiceCycleEnd(invoicePeriodStart, contract.BillingCycle)

	tenantSettings, err := s.services.TenantService.GetTenantSettings(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return "", err
	}

	currency := contract.Currency.String()
	if currency == "" {
		currency = tenantSettings.BaseCurrency.String()
	}

	dryRunInvoiceRequest := invoicepb.NewInvoiceForContractRequest{
		Tenant:             tenant,
		LoggedInUserId:     common.GetUserIdFromContext(ctx),
		ContractId:         contractId,
		DryRun:             true,
		CreatedAt:          utils.ConvertTimeToTimestampPtr(&now),
		InvoicePeriodStart: utils.ConvertTimeToTimestampPtr(&invoicePeriodStart),
		InvoicePeriodEnd:   utils.ConvertTimeToTimestampPtr(&invoicePeriodEnd),
		Currency:           currency,
		Note:               contract.InvoiceNote,
		Postpaid:           tenantSettings.InvoicingPostpaid,
		SourceFields: &commonpb.SourceFields{
			Source:    neo4jentity.DataSourceOpenline.String(),
			AppSource: appSource,
		},
	}

	switch contract.BillingCycle {
	case enum.BillingCycleMonthlyBilling:
		dryRunInvoiceRequest.BillingCycle = commonpb.BillingCycle_MONTHLY_BILLING
	case enum.BillingCycleQuarterlyBilling:
		dryRunInvoiceRequest.BillingCycle = commonpb.BillingCycle_QUARTERLY_BILLING
	case enum.BillingCycleAnnuallyBilling:
		dryRunInvoiceRequest.BillingCycle = commonpb.BillingCycle_ANNUALLY_BILLING
	}

	ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
	response, err := utils.CallEventsPlatformGRPCWithRetry[*invoicepb.InvoiceIdResponse](func() (*invoicepb.InvoiceIdResponse, error) {
		return s.services.GrpcClients.InvoiceClient.NewInvoiceForContract(ctx, &dryRunInvoiceRequest)
	})
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error from events processing: %s", err.Error())
		return "", err
	}

	span.LogFields(log.String("output - createdInvoiceId", response.Id))
	return response.Id, nil
}

func (s *invoiceService) PayInvoice(ctx context.Context, invoiceId, appSource string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceService.PayInvoice")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("invoiceId", invoiceId))

	tenant := common.GetTenantFromContext(ctx)

	ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
	response, err := utils.CallEventsPlatformGRPCWithRetry[*invoicepb.InvoiceIdResponse](func() (*invoicepb.InvoiceIdResponse, error) {
		return s.services.GrpcClients.InvoiceClient.UpdateInvoice(ctx, &invoicepb.UpdateInvoiceRequest{
			Tenant:         tenant,
			InvoiceId:      invoiceId,
			LoggedInUserId: common.GetUserIdFromContext(ctx),
			AppSource:      appSource,
			Status:         invoicepb.InvoiceStatus_INVOICE_STATUS_PAID,
			FieldsMask:     []invoicepb.InvoiceFieldMask{invoicepb.InvoiceFieldMask_INVOICE_FIELD_STATUS},
		})
	})

	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error from events processing: %s", err.Error())
		return err
	}

	span.LogFields(log.String("output - payInvoiceId", response.Id))
	return nil
}

func (s *invoiceService) VoidInvoice(ctx context.Context, invoiceId, appSource string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceService.VoidInvoice")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("invoiceId", invoiceId))

	tenant := common.GetTenantFromContext(ctx)

	ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
	response, err := utils.CallEventsPlatformGRPCWithRetry[*invoicepb.InvoiceIdResponse](func() (*invoicepb.InvoiceIdResponse, error) {
		return s.services.GrpcClients.InvoiceClient.VoidInvoice(ctx, &invoicepb.VoidInvoiceRequest{
			Tenant:         tenant,
			InvoiceId:      invoiceId,
			LoggedInUserId: common.GetUserIdFromContext(ctx),
			AppSource:      appSource,
		})
	})

	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error from events processing: %s", err.Error())
		return err
	}

	span.LogFields(log.String("output - voidInvoiceId", response.Id))
	return nil
}

func (h *invoiceService) FillCycleInvoice(ctx context.Context, tenant, contractId string, invoiceEntity *neo4jentity.InvoiceEntity, sliEntities neo4jentity.ServiceLineItemEntities) (*neo4jentity.InvoiceEntity, []*invoicepb.InvoiceLine, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceEventHandler.fillCycleInvoice")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	span.SetTag(tracing.SpanTagEntityId, invoiceEntity.Id)
	span.LogFields(log.String("contractId", contractId))

	amount, vat := float64(0), float64(0)
	var invoiceLines []*invoicepb.InvoiceLine

	referenceTime := invoiceEntity.PeriodStartDate
	periodEndTime := utils.EndOfDayInUTC(invoiceEntity.PeriodEndDate)
	if invoiceEntity.Postpaid {
		referenceTime = periodEndTime
	}
	for _, sliEntity := range sliEntities {
		// skip for now usage SLIs
		if sliEntity.Billed == neo4jenum.BilledTypeUsage {
			continue
		}
		// skip SLI if of None type
		if sliEntity.Billed == neo4jenum.BilledTypeNone {
			continue
		}
		// skip SLI if ended on the reference time
		if sliEntity.EndedAt != nil && sliEntity.EndedAt.Before(referenceTime) {
			continue
		}
		// skip SLI if not active on the reference time
		if sliEntity.IsRecurrent() && !sliEntity.IsActiveAt(referenceTime) {
			continue
		}
		// skip ONE TIME SLI if started after the end period
		if sliEntity.IsOneTime() && sliEntity.StartedAt.After(periodEndTime) {
			continue
		}
		// cancelled ONE TIME SLI should not be invoiced
		if sliEntity.IsOneTime() && sliEntity.Canceled {
			continue
		}

		// skip SLI if quantity or price is negative
		if sliEntity.Quantity < 0 || sliEntity.Price < 0 {
			continue
		}

		calculatedSLIAmount, calculatedSLIVat := float64(0), float64(0)
		invoiceLineCalculationsReady := false
		// process one time SLIs
		if sliEntity.Billed == neo4jenum.BilledTypeOnce {
			// Check any version of SLI not invoiced
			result, err := h.services.Neo4jRepositories.InvoiceLineReadRepository.GetLatestInvoiceLineWithInvoiceIdByServiceLineItemParentId(ctx, tenant, sliEntity.ParentID)
			if err != nil {
				tracing.TraceErr(span, err)
				h.log.Errorf("Error getting latest invoice line for sli parent id {%s}: {%s}", sliEntity.ParentID, err.Error())
			}
			if result != nil {
				// SLI already invoiced
				continue
			}
			quantity := sliEntity.Quantity
			calculatedSLIAmount = utils.TruncateFloat64(float64(quantity)*sliEntity.Price, 2)
			calculatedSLIVat = utils.TruncateFloat64(calculatedSLIAmount*sliEntity.VatRate/100, 2)
			invoiceLineCalculationsReady = true
		}
		// process monthly, quarterly and annually SLIs
		if sliEntity.Billed == neo4jenum.BilledTypeMonthly || sliEntity.Billed == neo4jenum.BilledTypeQuarterly || sliEntity.Billed == neo4jenum.BilledTypeAnnually {
			calculatedSLIAmount = calculateSLIAmountForCycleInvoicing(sliEntity.Quantity, sliEntity.Price, sliEntity.Billed, invoiceEntity.BillingCycle)
			calculatedSLIAmount = utils.TruncateFloat64(calculatedSLIAmount, 2)
			calculatedSLIVat = utils.TruncateFloat64(calculatedSLIAmount*sliEntity.VatRate/100, 2)
			invoiceLineCalculationsReady = true
		}
		if invoiceLineCalculationsReady {
			amount += calculatedSLIAmount
			vat += calculatedSLIVat
			invoiceLine := invoicepb.InvoiceLine{
				Name:                    sliEntity.Name,
				Price:                   utils.TruncateFloat64(calculatePriceForBilledType(sliEntity.Price, sliEntity.Billed, invoiceEntity.BillingCycle), 2),
				Quantity:                sliEntity.Quantity,
				Amount:                  calculatedSLIAmount,
				Total:                   calculatedSLIAmount + calculatedSLIVat,
				Vat:                     calculatedSLIVat,
				ServiceLineItemId:       sliEntity.ID,
				ServiceLineItemParentId: sliEntity.ParentID,
			}
			switch sliEntity.Billed {
			case neo4jenum.BilledTypeMonthly:
				invoiceLine.BilledType = commonpb.BilledType_MONTHLY_BILLED
			case neo4jenum.BilledTypeQuarterly:
				invoiceLine.BilledType = commonpb.BilledType_QUARTERLY_BILLED
			case neo4jenum.BilledTypeAnnually:
				invoiceLine.BilledType = commonpb.BilledType_ANNUALLY_BILLED
			case neo4jenum.BilledTypeOnce:
				invoiceLine.BilledType = commonpb.BilledType_ONCE_BILLED
			}
			invoiceLines = append(invoiceLines, &invoiceLine)
			continue
		}
		// if remained any unprocessed SLI log an error
		err := errors.Errorf("Unprocessed SLI %s", sliEntity.ID)
		tracing.TraceErr(span, err)
		h.log.Errorf("Error processing SLI during invoicing %s: %s", sliEntity.ID, err.Error())
	}

	if len(invoiceLines) == 0 {
		return nil, nil, errors.New("No invoice lines to fill")
	}

	invoiceEntity.Amount = amount
	invoiceEntity.Vat = vat
	invoiceEntity.TotalAmount = amount + vat

	return invoiceEntity, invoiceLines, nil
}

func calculateInvoiceCycleEnd(start time.Time, cycle enum.BillingCycle) time.Time {
	var end time.Time
	switch cycle {
	case enum.BillingCycleMonthlyBilling:
		end = start.AddDate(0, 1, 0)
	case enum.BillingCycleQuarterlyBilling:
		end = start.AddDate(0, 3, 0)
	case enum.BillingCycleAnnuallyBilling:
		end = start.AddDate(1, 0, 0)
	default:
		return start
	}
	previousDay := end.AddDate(0, 0, -1)
	return previousDay
}

func calculateSLIAmountForCycleInvoicing(quantity int64, price float64, billed neo4jenum.BilledType, cycle neo4jenum.BillingCycle) float64 {
	if quantity == 0 || price == 0 {
		return 0
	}
	unitAmount := calculatePriceForBilledType(price, billed, cycle)
	unitAmount = utils.TruncateFloat64(unitAmount, 2)
	return float64(quantity) * unitAmount
}

func calculatePriceForBilledType(price float64, billed neo4jenum.BilledType, cycle neo4jenum.BillingCycle) float64 {
	if billed == neo4jenum.BilledTypeOnce {
		return price
	}

	switch cycle {
	case neo4jenum.BillingCycleMonthlyBilling:
		switch billed {
		case neo4jenum.BilledTypeMonthly:
			return price
		case neo4jenum.BilledTypeQuarterly:
			return price / 3
		case neo4jenum.BilledTypeAnnually:
			return price / 12
		}
	case neo4jenum.BillingCycleQuarterlyBilling:
		switch billed {
		case neo4jenum.BilledTypeMonthly:
			return price * 3
		case neo4jenum.BilledTypeQuarterly:
			return price
		case neo4jenum.BilledTypeAnnually:
			return price / 4
		}
	case neo4jenum.BillingCycleAnnuallyBilling:
		switch billed {
		case neo4jenum.BilledTypeMonthly:
			return price * 12
		case neo4jenum.BilledTypeQuarterly:
			return price * 4
		case neo4jenum.BilledTypeAnnually:
			return price
		}
	}

	return 0
}
