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
	"sort"
	"time"
)

type InvoiceService interface {
	GenerateNewRandomInvoiceNumber() string

	GetById(ctx context.Context, invoiceId string) (*neo4jentity.InvoiceEntity, error)
	GetInvoiceLinesForInvoices(ctx context.Context, invoiceIds []string) (*neo4jentity.InvoiceLineEntities, error)
	GetInvoicesForContracts(ctx context.Context, contractIds []string) (*neo4jentity.InvoiceEntities, error)
	SimulateInvoice(ctx context.Context, invoiceData *SimulateInvoiceRequestData) ([]*SimulateInvoiceResponseData, error)
	NextInvoiceDryRun(ctx context.Context, contractId, appSource string) (string, error)
	PayInvoice(ctx context.Context, invoiceId, appSource string) error
	VoidInvoice(ctx context.Context, invoiceId, appSource string) error

	FillCycleInvoice(ctx context.Context, invoiceEntity *neo4jentity.InvoiceEntity, sliEntities neo4jentity.ServiceLineItemEntities) (*neo4jentity.InvoiceEntity, []*invoicepb.InvoiceLine, error)
	FillOffCyclePrepaidInvoice(ctx context.Context, invoiceEntity *neo4jentity.InvoiceEntity, sliEntities neo4jentity.ServiceLineItemEntities) (*neo4jentity.InvoiceEntity, []*invoicepb.InvoiceLine, error)
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
	Key               string
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

func (s *invoiceService) GenerateNewRandomInvoiceNumber() string {
	digits := "0123456789"
	consonants := "BCDFGHJKLMNPQRSTVWXYZ"
	invoiceNumber := utils.GenerateRandomStringFromCharset(3, consonants) + "-" + utils.GenerateRandomStringFromCharset(5, digits)
	return invoiceNumber
}

func (s *invoiceService) GetById(ctx context.Context, invoiceId string) (*neo4jentity.InvoiceEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceService.GetById")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("invoiceId", invoiceId))

	if invoiceDbNode, err := s.services.Neo4jRepositories.InvoiceReadRepository.GetInvoiceById(ctx, common.GetTenantFromContext(ctx), invoiceId); err != nil {
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

type SimulateInvoices struct {
	ContractId         string
	IssueDate          time.Time
	DueDate            time.Time
	InvoicePeriodStart time.Time
	InvoicePeriodEnd   time.Time
	InvoiceNumber      string
	InvoiceLines       []*SimulateInvoiceLine
}
type SimulateInvoiceLine struct {
	Key               string
	ServiceLineItemID string
	ParentID          string
	Name              string
	Price             float64
	Quantity          int64
	Amount            float64
	TotalAmount       float64
}

func (s *invoiceService) SimulateInvoice(ctx context.Context, simulateInvoicesWithChanges *SimulateInvoiceRequestData) ([]*SimulateInvoiceResponseData, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceService.SimulateInvoice")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("simulateInvoicesWithChanges", simulateInvoicesWithChanges))

	if simulateInvoicesWithChanges.ServiceLines == nil {
		err := fmt.Errorf("no invoice lines to simulate")
		tracing.TraceErr(span, err)
		return nil, err
	}

	response := []*SimulateInvoiceResponseData{}

	//fetch existing contract and set the next invoice date
	contract, err := s.services.ContractService.GetById(ctx, simulateInvoicesWithChanges.ContractId)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	if contract.InvoicingStartDate == nil && contract.NextInvoiceDate == nil {
		err := fmt.Errorf("contract has no invoicing start date or next invoice date")
		tracing.TraceErr(span, err)
		return nil, err
	}
	if contract.NextInvoiceDate == nil {
		contract.NextInvoiceDate = contract.InvoicingStartDate
	}
	contractDefaultNextInvoiceDate := *contract.NextInvoiceDate

	existingSlis, err := s.services.ServiceLineItemService.GetServiceLineItemsForContract(ctx, contract.Id)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	sliEntities := *existingSlis

	//determine the interval to compute invoices
	//[ current contract invoice date, max service start date ]
	invoicePeriodStartGeneration := contract.NextInvoiceDate
	invoicePeriodEndGeneration := time.Time{}

	for _, sliData := range simulateInvoicesWithChanges.ServiceLines {
		if sliData.ServiceStarted.After(invoicePeriodEndGeneration) {
			invoicePeriodEndGeneration = sliData.ServiceStarted
		}
	}
	invoicePeriodEndGeneration = invoicePeriodEndGeneration.AddDate(0, int(contract.BillingCycleInMonths), 0)

	for true {
		if invoicePeriodStartGeneration.After(invoicePeriodEndGeneration) {
			break
		}

		//identify first SLI change in the simualtion and replace it in the sliEntities
		for _, sliData := range simulateInvoicesWithChanges.ServiceLines {

			prorationNeeded := false

			if sliData.ServiceStarted.After(*invoicePeriodStartGeneration) {

				if sliData.ServiceLineItemID == nil || *sliData.ServiceLineItemID == "" {
					//new sli item - adding it to the sli entities and trigger proration
					sliEntities = append(sliEntities, neo4jentity.ServiceLineItemEntity{
						ID:        sliData.Key,
						Name:      sliData.Description,
						Comments:  sliData.Comments,
						Billed:    sliData.BillingCycle,
						Price:     sliData.Price,
						Quantity:  sliData.Quantity,
						StartedAt: sliData.ServiceStarted,
						EndedAt:   nil,
					})
					prorationNeeded = true
				} else {
					//existing sli item - to check if there is any change in the sli item to decide if proration is needed
					existingSli, err := s.services.ServiceLineItemService.GetById(ctx, *sliData.ServiceLineItemID)
					if err != nil {
						tracing.TraceErr(span, err)
						return nil, err
					}

					if sliData.ServiceStarted.After(utils.Today()) {
						//new version
					} else {
						//update existing version

						if existingSli.Billed != sliData.BillingCycle || existingSli.Price != sliData.Price || existingSli.Quantity != sliData.Quantity || existingSli.StartedAt != sliData.ServiceStarted || existingSli.VatRate != utils.IfNotNilFloat64(sliData.TaxRate) {
							//existing sli item - new version - proration needed
							for i, sliEntity := range sliEntities {
								if sliEntity.ID == *sliData.ServiceLineItemID {
									sliEntities[i].Billed = sliData.BillingCycle
									sliEntities[i].Price = sliData.Price
									sliEntities[i].Quantity = sliData.Quantity
									sliEntities[i].StartedAt = sliData.ServiceStarted
									sliEntities[i].VatRate = utils.IfNotNilFloat64(sliData.TaxRate)
									break
								}
							}
							prorationNeeded = true
						}
					}
				}

			}

			if prorationNeeded {

				//proratedInvoice, err := s.SimulateOffCycleInvoice(ctx, contract, simulateInvoicesWithChanges, span)
				//if err != nil {
				//	tracing.TraceErr(span, err)
				//	return nil, err
				//}
				//
				response = append(response, &SimulateInvoiceResponseData{Invoice: &neo4jentity.InvoiceEntity{
					Number:               "",
					OffCycle:             true,
					Postpaid:             false,
					Amount:               0,
					Vat:                  0,
					TotalAmount:          0,
					BillingCycleInMonths: contract.BillingCycleInMonths,
					PeriodStartDate:      *invoicePeriodStartGeneration,
					PeriodEndDate:        invoicePeriodEndGeneration,
				}, Lines: []*neo4jentity.InvoiceLineEntity{}})

				nextOnCycleInvoiceDate := contract.NextInvoiceDate.AddDate(0, int(contract.BillingCycleInMonths), 0)
				contract.NextInvoiceDate = &nextOnCycleInvoiceDate
				onCycleInvoice, err := s.SimulateOnCycleInvoice(ctx, contract, &sliEntities, span)
				if err != nil {
					tracing.TraceErr(span, err)
					return nil, err
				}

				response = append(response, onCycleInvoice)
			}

			nextInvoiceDate := invoicePeriodStartGeneration.AddDate(0, int(contract.BillingCycleInMonths), 0)
			invoicePeriodStartGeneration = &nextInvoiceDate

			nextOnCycleDate := contract.NextInvoiceDate.AddDate(0, int(contract.BillingCycleInMonths), 0)
			contract.NextInvoiceDate = &nextOnCycleDate
		}
	}

	//no proration needed - only on cycle invoice
	if len(response) == 0 {

		contract.NextInvoiceDate = &contractDefaultNextInvoiceDate

		//build sli entities to reflect the changes for the period
		onCycleSliEntities := neo4jentity.ServiceLineItemEntities{}
		for _, sliData := range simulateInvoicesWithChanges.ServiceLines {

			nextInvoiceDate := contract.NextInvoiceDate.Add(1) // adding 1 nanosecond to the next invoice date to consider the first nanosecond of the month
			if sliData.ServiceStarted.Before(nextInvoiceDate) {
				sliEntity := neo4jentity.ServiceLineItemEntity{
					ID:        sliData.Key,
					Name:      sliData.Description,
					Comments:  sliData.Comments,
					Billed:    sliData.BillingCycle,
					Price:     sliData.Price,
					Quantity:  sliData.Quantity,
					StartedAt: sliData.ServiceStarted,
					EndedAt:   nil,
					VatRate:   utils.IfNotNilFloat64(sliData.TaxRate),
				}

				onCycleSliEntities = append(onCycleSliEntities, sliEntity)
			}
		}

		onCycleInvoice, err := s.SimulateOnCycleInvoice(ctx, contract, &onCycleSliEntities, span)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}

		response = append(response, onCycleInvoice)
	}

	return response, nil
}

func (s *invoiceService) SimulateOnCycleInvoice(ctx context.Context, contract *neo4jentity.ContractEntity, sliEntities *neo4jentity.ServiceLineItemEntities, span opentracing.Span) (*SimulateInvoiceResponseData, error) {
	invoiceEntity := &neo4jentity.InvoiceEntity{}
	invoiceLines := []*invoicepb.InvoiceLine{}

	tenantSettings, err := s.services.TenantService.GetTenantSettings(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	var invoicePeriodStart, invoicePeriodEnd time.Time
	if contract.NextInvoiceDate != nil {
		invoicePeriodStart = *contract.NextInvoiceDate
	} else {
		invoicePeriodStart = *contract.InvoicingStartDate
	}
	invoicePeriodEnd = calculateInvoiceCycleEnd(invoicePeriodStart, contract.BillingCycleInMonths)

	invoiceEntity.Number = s.GenerateNewRandomInvoiceNumber()
	invoiceEntity.OffCycle = false
	invoiceEntity.Postpaid = tenantSettings.InvoicingPostpaid
	invoiceEntity.BillingCycleInMonths = contract.BillingCycleInMonths
	invoiceEntity.PeriodStartDate = invoicePeriodStart
	invoiceEntity.PeriodEndDate = invoicePeriodEnd
	invoiceEntity.Note = contract.InvoiceNote
	invoiceEntity.IssuedDate = invoicePeriodStart
	invoiceEntity.DueDate = invoicePeriodEnd.AddDate(0, 0, int(contract.DueDays))

	if contract.Currency != "" {
		invoiceEntity.Currency = contract.Currency
	} else {
		invoiceEntity.Currency = tenantSettings.BaseCurrency
	}

	invoiceEntity, invoiceLines, err = s.FillCycleInvoice(ctx, invoiceEntity, *sliEntities)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	onCycleInvoice := &SimulateInvoiceResponseData{
		Invoice: invoiceEntity,
		Lines:   []*neo4jentity.InvoiceLineEntity{},
	}
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
		onCycleInvoice.Lines = append(onCycleInvoice.Lines, invoiceLineEntity)
	}

	return onCycleInvoice, nil
}

func (s *invoiceService) SimulateOffCycleInvoice(ctx context.Context, contract *neo4jentity.ContractEntity, invoiceData *SimulateInvoiceRequestData, span opentracing.Span) (*SimulateInvoiceResponseData, error) {
	invoiceEntity := &neo4jentity.InvoiceEntity{}
	invoiceLines := []*invoicepb.InvoiceLine{}

	tenantSettings, err := s.services.TenantService.GetTenantSettings(ctx)
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
	if contract.NextInvoiceDate != nil {
		invoicePeriodStart = *contract.NextInvoiceDate
	} else {
		invoicePeriodStart = *contract.InvoicingStartDate
	}
	invoicePeriodEnd = calculateInvoiceCycleEnd(invoicePeriodStart, contract.BillingCycleInMonths)

	invoiceEntity.OffCycle = false
	invoiceEntity.Postpaid = tenantSettings.InvoicingPostpaid
	invoiceEntity.BillingCycleInMonths = contract.BillingCycleInMonths
	invoiceEntity.PeriodStartDate = invoicePeriodStart
	invoiceEntity.PeriodEndDate = invoicePeriodEnd
	invoiceEntity.Note = contract.InvoiceNote

	if nextPreviewInvoiceNode != nil {
		invoiceEntity.IssuedDate = nextPreviewInvoiceEntity.IssuedDate
		invoiceEntity.DueDate = nextPreviewInvoiceEntity.DueDate
	} else {
		invoiceEntity.IssuedDate = invoicePeriodStart
		invoiceEntity.DueDate = invoicePeriodEnd.AddDate(0, 0, int(contract.DueDays))
	}

	if contract.Currency != "" {
		invoiceEntity.Currency = contract.Currency
	} else {
		invoiceEntity.Currency = tenantSettings.BaseCurrency
	}

	sliEntities := neo4jentity.ServiceLineItemEntities{}
	for _, sliData := range invoiceData.ServiceLines {
		sliEntity := neo4jentity.ServiceLineItemEntity{
			ID:        sliData.Key,
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

	invoiceEntity, invoiceLines, err = s.FillOffCyclePrepaidInvoice(ctx, invoiceEntity, sliEntities)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	onCycleInvoice := &SimulateInvoiceResponseData{
		Invoice: invoiceEntity,
		Lines:   []*neo4jentity.InvoiceLineEntity{},
	}
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
		onCycleInvoice.Lines = append(onCycleInvoice.Lines, invoiceLineEntity)
	}

	return onCycleInvoice, nil
}

func (s *invoiceService) NextInvoiceDryRun(ctx context.Context, contractId, appSource string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceService.NextInvoiceDryRun")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("contractId", contractId))

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
	invoicePeriodEnd = calculateInvoiceCycleEnd(invoicePeriodStart, contract.BillingCycleInMonths)

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
		Tenant:             common.GetTenantFromContext(ctx),
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

	dryRunInvoiceRequest.BillingCycleInMonths = contract.BillingCycleInMonths

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

	ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
	response, err := utils.CallEventsPlatformGRPCWithRetry[*invoicepb.InvoiceIdResponse](func() (*invoicepb.InvoiceIdResponse, error) {
		return s.services.GrpcClients.InvoiceClient.UpdateInvoice(ctx, &invoicepb.UpdateInvoiceRequest{
			Tenant:         common.GetTenantFromContext(ctx),
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

	ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
	response, err := utils.CallEventsPlatformGRPCWithRetry[*invoicepb.InvoiceIdResponse](func() (*invoicepb.InvoiceIdResponse, error) {
		return s.services.GrpcClients.InvoiceClient.VoidInvoice(ctx, &invoicepb.VoidInvoiceRequest{
			Tenant:         common.GetTenantFromContext(ctx),
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

func (h *invoiceService) FillCycleInvoice(ctx context.Context, invoiceEntity *neo4jentity.InvoiceEntity, sliEntities neo4jentity.ServiceLineItemEntities) (*neo4jentity.InvoiceEntity, []*invoicepb.InvoiceLine, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceEventHandler.fillCycleInvoice")
	defer span.Finish()
	span.SetTag(tracing.SpanTagEntityId, invoiceEntity.Id)

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
			result, err := h.services.Neo4jRepositories.InvoiceLineReadRepository.GetLatestInvoiceLineWithInvoiceIdByServiceLineItemParentId(ctx, common.GetTenantFromContext(ctx), sliEntity.ParentID)
			if err != nil {
				tracing.TraceErr(span, err)
				h.log.Errorf("Error getting latest invoice line for sli parent id {%s}: {%s}", sliEntity.ParentID, err.Error())
			}
			if result != nil {
				// SLI already invoiced
				continue
			}
			quantity := sliEntity.Quantity
			calculatedSLIAmount = utils.RoundHalfUpFloat64(float64(quantity)*sliEntity.Price, 2)
			calculatedSLIVat = utils.RoundHalfUpFloat64(calculatedSLIAmount*sliEntity.VatRate/100, 2)
			invoiceLineCalculationsReady = true
		}
		// process monthly, quarterly and annually SLIs
		if sliEntity.Billed == neo4jenum.BilledTypeMonthly || sliEntity.Billed == neo4jenum.BilledTypeQuarterly || sliEntity.Billed == neo4jenum.BilledTypeAnnually {
			calculatedSLIAmount = calculateSLIAmountForCycleInvoicing(sliEntity.Quantity, sliEntity.Price, sliEntity.Billed, invoiceEntity.BillingCycleInMonths)
			calculatedSLIAmount = utils.RoundHalfUpFloat64(calculatedSLIAmount, 2)
			calculatedSLIVat = utils.RoundHalfUpFloat64(calculatedSLIAmount*sliEntity.VatRate/100, 2)
			invoiceLineCalculationsReady = true
		}
		if invoiceLineCalculationsReady {
			amount += calculatedSLIAmount
			vat += calculatedSLIVat
			invoiceLine := invoicepb.InvoiceLine{
				Name:                    sliEntity.Name,
				Price:                   utils.RoundHalfUpFloat64(calculatePriceForBilledType(sliEntity.Price, sliEntity.Billed, invoiceEntity.BillingCycleInMonths), 2),
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

func (h *invoiceService) FillOffCyclePrepaidInvoice(ctx context.Context, invoiceEntity *neo4jentity.InvoiceEntity, sliEntities neo4jentity.ServiceLineItemEntities) (*neo4jentity.InvoiceEntity, []*invoicepb.InvoiceLine, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceService.FillOffCyclePrepaidInvoice")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, common.GetTenantFromContext(ctx))
	span.SetTag(tracing.SpanTagEntityId, invoiceEntity.Id)

	// filter out not applicable SLIs
	referenceDate := invoiceEntity.PeriodStartDate
	filteredSliEntities := neo4jentity.ServiceLineItemEntities{}
	for _, sliEntity := range sliEntities {
		// process only monthly, quarterly and annually SLIs
		if sliEntity.Billed != neo4jenum.BilledTypeMonthly &&
			sliEntity.Billed != neo4jenum.BilledTypeQuarterly &&
			sliEntity.Billed != neo4jenum.BilledTypeAnnually &&
			sliEntity.Billed != neo4jenum.BilledTypeOnce {
			continue
		}
		// SLIs that started on or after reference date are not applicable
		if sliEntity.StartedAt.After(referenceDate) || sliEntity.StartedAt.Equal(referenceDate) {
			continue
		}
		// One time invoiced and cancelled SLIs are not applicable
		if sliEntity.Billed == neo4jenum.BilledTypeOnce {
			if sliEntity.Quantity <= 0 || sliEntity.Price <= 0 {
				continue
			}
			if sliEntity.Canceled {
				continue
			}
			ilDbNodeAndInvoiceId, err := h.services.Neo4jRepositories.InvoiceLineReadRepository.GetLatestInvoiceLineWithInvoiceIdByServiceLineItemParentId(ctx, common.GetTenantFromContext(ctx), sliEntity.ParentID)
			if err != nil {
				tracing.TraceErr(span, err)
				h.log.Errorf("Error getting latest invoice line for sli parent id {%s}: {%s}", sliEntity.ParentID, err.Error())
				return nil, nil, err
			}
			if ilDbNodeAndInvoiceId != nil {
				continue
			}
		}
		filteredSliEntities = append(filteredSliEntities, sliEntity)
	}
	// sort SLIs by startedAt
	sort.Slice(filteredSliEntities, func(i, j int) bool {
		return filteredSliEntities[i].StartedAt.Before(filteredSliEntities[j].StartedAt)
	})
	// group SLIs by parent id
	sliByParentID := map[string][]neo4jentity.ServiceLineItemEntity{}
	for _, sliEntity := range filteredSliEntities {
		sliByParentID[sliEntity.ParentID] = append(sliByParentID[sliEntity.ParentID], sliEntity)
	}

	span.LogFields(log.Int("result - amount of SLIs to process", len(filteredSliEntities)))

	amount, vat := float64(0), float64(0)
	var invoiceLines []*invoicepb.InvoiceLine

	proratedSliFound := false
	// iterate SLIs by parent id
	for parentId, slis := range sliByParentID {
		// get latest SLI that is active on reference date
		var sliEntityToInvoice *neo4jentity.ServiceLineItemEntity
		for _, sliEntity := range slis {
			if sliEntity.IsActiveAt(invoiceEntity.PeriodStartDate) {
				sliEntityToInvoice = &sliEntity
			}
		}
		// if no SLI is active on reference date, skip
		if sliEntityToInvoice == nil {
			span.LogFields(log.String("result - no active SLI for parent id", parentId))
			continue
		}
		// get invoice line for latest invoiced SLI per parent
		ilDbNodeAndInvoiceId, err := h.services.Neo4jRepositories.InvoiceLineReadRepository.GetLatestInvoiceLineWithInvoiceIdByServiceLineItemParentId(ctx, common.GetTenantFromContext(ctx), parentId)
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("Error getting latest invoice line for sli parent id {%s}: {%s}", parentId, err.Error())
			return nil, nil, err
		}
		finalSLIAmount, calculatedSLIVat := float64(0), float64(0)
		if sliEntityToInvoice.Billed == neo4jenum.BilledTypeOnce {
			quantity := sliEntityToInvoice.Quantity
			if quantity <= 0 {
				quantity = 1
			}
			finalSLIAmount = utils.RoundHalfUpFloat64(float64(quantity)*sliEntityToInvoice.Price, 2)
			if finalSLIAmount <= 0 {
				continue
			}
			calculatedSLIVat = utils.RoundHalfUpFloat64(finalSLIAmount*sliEntityToInvoice.VatRate/100, 2)
		} else {
			proratedInvoicedSLIAmount := float64(0)
			if ilDbNodeAndInvoiceId != nil {
				previousInvoiceDbNode, err := h.services.Neo4jRepositories.InvoiceReadRepository.GetInvoiceById(ctx, common.GetTenantFromContext(ctx), ilDbNodeAndInvoiceId.LinkedNodeId)
				if err != nil {
					tracing.TraceErr(span, err)
					h.log.Errorf("Error getting invoice {%s}: {%s}", ilDbNodeAndInvoiceId.LinkedNodeId, err.Error())
					return nil, nil, err
				}
				previousInvoiceEntity := mapper.MapDbNodeToInvoiceEntity(previousInvoiceDbNode)
				// if previous invoice is for different cycle, charge full amount
				if !previousInvoiceEntity.PeriodEndDate.Before(invoiceEntity.PeriodEndDate) {
					// calculate already invoiced amount, prorated for the period
					invoiceLineEntity := mapper.MapDbNodeToInvoiceLineEntity(ilDbNodeAndInvoiceId.Node)
					calculatedInvoicedSLIAmountFor1Year := calculateSLIAmountForCycleInvoicing(invoiceLineEntity.Quantity, invoiceLineEntity.Price, invoiceLineEntity.BilledType, 12)
					proratedInvoicedSLIAmount = prorateAnnualSLIAmount(sliEntityToInvoice.StartedAt, invoiceEntity.PeriodEndDate, calculatedInvoicedSLIAmountFor1Year)
					proratedInvoicedSLIAmount = utils.RoundHalfUpFloat64(proratedInvoicedSLIAmount, 2)
				}
			}

			calculatedSLIAmountFor1Year := calculateSLIAmountForCycleInvoicing(sliEntityToInvoice.Quantity, sliEntityToInvoice.Price, sliEntityToInvoice.Billed, 12)
			proratedSLIAmount := prorateAnnualSLIAmount(sliEntityToInvoice.StartedAt, invoiceEntity.PeriodEndDate, calculatedSLIAmountFor1Year)
			proratedSLIAmount = utils.RoundHalfUpFloat64(proratedSLIAmount, 2)
			finalSLIAmount = proratedSLIAmount - proratedInvoicedSLIAmount
			span.LogFields(log.Float64(fmt.Sprintf("result - final amount for SLI with parent id %s", parentId), finalSLIAmount))
			if finalSLIAmount <= 0 {
				continue
			}
			calculatedSLIVat = utils.RoundHalfUpFloat64(finalSLIAmount*sliEntityToInvoice.VatRate/100, 2)
			proratedSliFound = true
		}
		amount += finalSLIAmount
		vat += calculatedSLIVat
		invoiceLine := invoicepb.InvoiceLine{
			Name:                    sliEntityToInvoice.Name,
			Price:                   utils.RoundHalfUpFloat64(calculatePriceForBilledType(sliEntityToInvoice.Price, sliEntityToInvoice.Billed, invoiceEntity.BillingCycleInMonths), 2),
			Quantity:                sliEntityToInvoice.Quantity,
			Amount:                  finalSLIAmount,
			Total:                   finalSLIAmount + calculatedSLIVat,
			Vat:                     calculatedSLIVat,
			ServiceLineItemId:       sliEntityToInvoice.ID,
			ServiceLineItemParentId: sliEntityToInvoice.ParentID,
		}
		switch sliEntityToInvoice.Billed {
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
	}

	if !proratedSliFound && len(invoiceLines) > 0 {
		// if no prorated SLI found, then invoice contains only once billed SLIs
		// accept the invoice if today is monthly anniversary of the contract invoicing start date

		// UPDATE: The rule is on hold, invoice will be issued even if contains only one time SLIs

		//if !isMonthlyAnniversary(invoiceEntity.PeriodEndDate.AddDate(0, 0, 1)) {
		//	invoiceLines = []*invoicepb.InvoiceLine{}
		//}
	}

	invoiceEntity.Amount = amount
	invoiceEntity.Vat = vat
	invoiceEntity.TotalAmount = amount + vat

	return invoiceEntity, invoiceLines, nil
}

func calculateInvoiceCycleEnd(start time.Time, billingCycleInMonths int64) time.Time {
	end := start.AddDate(0, int(billingCycleInMonths), 0)
	previousDay := end.AddDate(0, 0, -1)
	return previousDay
}

func calculateSLIAmountForCycleInvoicing(quantity int64, price float64, billed neo4jenum.BilledType, billingCycleInMonths int64) float64 {
	if quantity == 0 || price == 0 {
		return 0
	}
	unitAmount := calculatePriceForBilledType(price, billed, billingCycleInMonths)
	unitAmount = utils.RoundHalfUpFloat64(unitAmount, 2)
	return float64(quantity) * unitAmount
}

func calculatePriceForBilledType(price float64, billed neo4jenum.BilledType, billingCycleInMonths int64) float64 {
	if billed == neo4jenum.BilledTypeOnce {
		return price
	}

	if billingCycleInMonths == 0 || billed.InMonths() == 0 {
		return 0
	}

	return price * float64(billingCycleInMonths) / float64(billed.InMonths())
}

func prorateAnnualSLIAmount(startDate, endDate time.Time, amount float64) float64 {
	start := utils.ToDate(startDate)
	end := utils.ToDate(endDate)
	days := end.Sub(start).Hours() / 24
	proratedAmount := amount * (days / 365)
	if proratedAmount <= 0 {
		return 0
	}
	return proratedAmount
}
