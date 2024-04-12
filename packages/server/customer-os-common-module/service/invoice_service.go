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
	SimulateInvoice(ctx context.Context, invoiceData *SimulateInvoiceData) (*neo4jentity.InvoiceEntities, error)
	NextInvoiceDryRun(ctx context.Context, contractId, appSource string) (string, error)
	PayInvoice(ctx context.Context, invoiceId, appSource string) error
	VoidInvoice(ctx context.Context, invoiceId, appSource string) error
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

type SimulateInvoiceData struct {
	ContractId   string
	InvoiceLines []SimulateInvoiceLineData
}
type SimulateInvoiceLineData struct {
	ServiceLineItemID *string
	ParentID          *string
	Description       string
	Comments          string
	BillingCycle      enum.BilledType
	Price             float64
	Quantity          int
	ServiceStarted    time.Time
	TaxRate           *float64
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

func (s *invoiceService) SimulateInvoice(ctx context.Context, invoiceData *SimulateInvoiceData) (*neo4jentity.InvoiceEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceService.SimulateInvoice")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("invoiceData", invoiceData))

	if invoiceData.InvoiceLines == nil {
		err := fmt.Errorf("no invoice lines to simulate")
		tracing.TraceErr(span, err)
		return nil, err
	}

	return nil, nil
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
