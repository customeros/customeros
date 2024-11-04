package service

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	neo4jrepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"
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
	GetByIdAcrossAllTenants(ctx context.Context, invoiceId string) (*neo4jentity.InvoiceEntity, string, error)
	GetByNumber(ctx context.Context, number string) (*neo4jentity.InvoiceEntity, error)
	GetInvoiceLinesForInvoices(ctx context.Context, invoiceIds []string) (*neo4jentity.InvoiceLineEntities, error)
	GetInvoicesForContracts(ctx context.Context, contractIds []string) (*neo4jentity.InvoiceEntities, error)
	GetNonDryRunInvoicesForOrganization(ctx context.Context, tenant, organizationId string) (*neo4jentity.InvoiceEntities, error)
	SimulateInvoice(ctx context.Context, invoiceData *SimulateInvoiceRequestData) ([]*SimulateInvoiceResponseData, error)
	NextInvoiceDryRun(ctx context.Context, contractId, appSource string) (string, error)
	PayInvoice(ctx context.Context, invoiceId, appSource string) error
	VoidInvoice(ctx context.Context, invoiceId, appSource string) error
	UpdateInvoice(ctx context.Context, invoiceId string, data neo4jrepository.InvoiceUpdateFields) error

	FillCycleInvoice(ctx context.Context, invoiceEntity *neo4jentity.InvoiceEntity, sliEntities neo4jentity.ServiceLineItemEntities) (*neo4jentity.InvoiceEntity, []*invoicepb.InvoiceLine, error)
	// Deprecated: Method should be re-worked. DO NOT ENABLE IN PROD
	FillOffCyclePrepaidInvoice(ctx context.Context, invoiceEntity *neo4jentity.InvoiceEntity, sliEntities neo4jentity.ServiceLineItemEntities) (*neo4jentity.InvoiceEntity, []*invoicepb.InvoiceLine, error)
}
type invoiceService struct {
	log      logger.Logger
	services *Services
}

func NewInvoiceService(services *Services) InvoiceService {
	return &invoiceService{
		services: services,
	}
}

type InvoiceActionMetadata struct {
	Status        string  `json:"status"`
	Currency      string  `json:"currency"`
	Amount        float64 `json:"amount"`
	InvoiceNumber string  `json:"number"`
	InvoiceId     string  `json:"id"`
}

type SimulateInvoiceRequestData struct {
	ContractId   string
	ServiceLines []SimulateInvoiceRequestServiceLineData
}
type SimulateInvoiceRequestServiceLineData struct {
	Key               string
	ServiceLineItemID string
	ParentID          string
	Description       string
	Comments          string
	BillingCycle      enum.BilledType
	Price             float64
	Quantity          int64
	ServiceStarted    time.Time
	ServiceEnded      *time.Time
	TaxRate           *float64
	Canceled          bool
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

func (s *invoiceService) GetByIdAcrossAllTenants(ctx context.Context, invoiceId string) (*neo4jentity.InvoiceEntity, string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceService.GetByIdAcrossAllTenants")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogKV("invoiceId", invoiceId)

	invoiceDbNode, tenant, err := s.services.Neo4jRepositories.InvoiceReadRepository.GetInvoiceByIdAcrossAllTenants(ctx, invoiceId)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "error getting invoice by id"))
		wrappedErr := errors.Wrap(err, fmt.Sprintf("Invoice with id {%s} not found", invoiceId))
		return nil, "", wrappedErr
	}
	if invoiceDbNode == nil {
		return nil, "", nil
	} else {
		return mapper.MapDbNodeToInvoiceEntity(invoiceDbNode), tenant, nil
	}
}

func (s *invoiceService) GetByNumber(ctx context.Context, number string) (*neo4jentity.InvoiceEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceService.GetByNumber")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("number", number))

	if invoiceDbNode, err := s.services.Neo4jRepositories.InvoiceReadRepository.GetInvoiceByNumber(ctx, common.GetTenantFromContext(ctx), number); err != nil {
		tracing.TraceErr(span, err)
		wrappedErr := errors.Wrap(err, fmt.Sprintf("Invoice with number {%s} not found", number))
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
		tracing.TraceErr(span, err)
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

	if len(simulateInvoicesWithChanges.ServiceLines) == 0 {
		err := fmt.Errorf("no invoice lines to simulate")
		tracing.TraceErr(span, err)
		return nil, err
	}

	var response []*SimulateInvoiceResponseData

	tenantSettings, err := s.services.TenantService.GetTenantSettings(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	//fetch existing contract and set the next invoice date
	contract, err := s.services.ContractService.GetById(ctx, simulateInvoicesWithChanges.ContractId)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	invoiceDate := utils.Today()
	if contract.NextInvoiceDate != nil {
		invoiceDate = *contract.NextInvoiceDate
	} else if contract.InvoicingStartDate != nil {
		invoiceDate = *contract.InvoicingStartDate
	}

	existingSlis, err := s.services.ServiceLineItemService.GetServiceLineItemsForContract(ctx, contract.Id)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	sliEntities := *existingSlis

	// TODO temporary disable prorated off-cycle simulated invoices
	allowProratedOffCycleInvoices := false
	if !tenantSettings.InvoicingPostpaid && allowProratedOffCycleInvoices {
		//determine the interval to compute invoices
		//[ invoicing starts, max service start date ]
		invoicePeriodStartGeneration := invoiceDate
		if contract.InvoicingStartDate != nil {
			invoicePeriodStartGeneration = *contract.InvoicingStartDate
		}
		invoicePeriodEndGeneration := time.Time{}

		for _, sliData := range simulateInvoicesWithChanges.ServiceLines {
			if sliData.ServiceStarted.After(invoicePeriodEndGeneration) {
				invoicePeriodEndGeneration = sliData.ServiceStarted
			}
		}
		invoicePeriodEndGeneration = calculateInvoiceCycleEnd(invoicePeriodEndGeneration, contract.BillingCycleInMonths)

		for true {
			if invoicePeriodStartGeneration.After(invoicePeriodEndGeneration) {
				break
			}

			//identify first SLI change in the simualtion and replace it in the sliEntities
			for _, sliData := range simulateInvoicesWithChanges.ServiceLines {

				prorationNeeded := false

				if sliData.ServiceStarted.Before(invoicePeriodStartGeneration.Add(1)) {
					continue
				}

				if sliData.ServiceStarted.After(calculateInvoiceCycleEnd(invoicePeriodStartGeneration, contract.BillingCycleInMonths)) {
					continue
				}

				if sliData.ServiceLineItemID == "" {
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
					existingSli, err := s.services.ServiceLineItemService.GetById(ctx, sliData.ServiceLineItemID)
					if err != nil {
						tracing.TraceErr(span, err)
						return nil, err
					}

					if sliData.ServiceStarted.After(utils.Today()) {
						//new version
					} else {
						//update existing version

						if (existingSli.Billed != sliData.BillingCycle || existingSli.Price != sliData.Price || existingSli.Quantity != sliData.Quantity || existingSli.StartedAt != sliData.ServiceStarted || existingSli.VatRate != utils.IfNotNilFloat64(sliData.TaxRate)) && (existingSli.Price*float64(existingSli.Quantity) < sliData.Price*float64(sliData.Quantity)) {
							//existing sli item - new version - proration needed
							for i, sliEntity := range sliEntities {
								if sliEntity.ID == sliData.ServiceLineItemID {
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

				if prorationNeeded {

					proratedInvoice, err := s.SimulateOffCycleInvoice(ctx, contract, &sliEntities, span)
					if err != nil {
						tracing.TraceErr(span, err)
						return nil, err
					}

					response = append(response, proratedInvoice)

					contract.NextInvoiceDate = utils.Ptr(proratedInvoice.Invoice.PeriodEndDate.AddDate(0, 0, 1))
					onCycleInvoice, err := s.SimulateOnCycleInvoice(ctx, contract, &sliEntities, span)
					if err != nil {
						tracing.TraceErr(span, err)
						return nil, err
					}

					response = append(response, onCycleInvoice)
				}

			}

			nextInvoiceDate := invoicePeriodStartGeneration.AddDate(0, int(contract.BillingCycleInMonths), 0)
			invoicePeriodStartGeneration = nextInvoiceDate

			nextOnCycleDate := invoiceDate.AddDate(0, int(contract.BillingCycleInMonths), 0)
			contract.NextInvoiceDate = &nextOnCycleDate
		}
	}

	//no proration needed - only on cycle invoice
	if len(response) == 0 {

		contract.NextInvoiceDate = &invoiceDate

		//build sli entities to reflect the changes for the period
		onCycleSliEntities := neo4jentity.ServiceLineItemEntities{}

		// prepare simulation invoice lines grouped by parent id
		invoiceLinesGroupedByParentId := map[string][]SimulateInvoiceRequestServiceLineData{}
		for i := range simulateInvoicesWithChanges.ServiceLines {
			sliData := &simulateInvoicesWithChanges.ServiceLines[i]
			if sliData.ServiceLineItemID == "" {
				sliData.ServiceLineItemID = sliData.Key
			}
			if sliData.ServiceLineItemID == "" {
				sliData.ServiceLineItemID = uuid.New().String()
			}
			if sliData.ParentID == "" {
				sliData.ParentID = sliData.ServiceLineItemID
			}
			invoiceLinesGroupedByParentId[sliData.ParentID] = append(invoiceLinesGroupedByParentId[sliData.ParentID], *sliData)
		}
		// per parent group sort by service started date and set end date for each sli as previous sli start date
		for _, sliDataGroup := range invoiceLinesGroupedByParentId {
			sort.Slice(sliDataGroup, func(i, j int) bool {
				return sliDataGroup[i].ServiceStarted.Before(sliDataGroup[j].ServiceStarted)
			})
			for i, sliData := range sliDataGroup {
				if i > 0 {
					sliDataGroup[i-1].ServiceEnded = &sliData.ServiceStarted
				}
			}
		}

		for _, sliData := range simulateInvoicesWithChanges.ServiceLines {
			sliEntity := neo4jentity.ServiceLineItemEntity{
				ID:        utils.IfNotNilString(sliData.ServiceLineItemID),
				ParentID:  utils.IfNotNilString(sliData.ParentID),
				Name:      sliData.Description,
				Comments:  sliData.Comments,
				Billed:    sliData.BillingCycle,
				Price:     sliData.Price,
				Quantity:  sliData.Quantity,
				StartedAt: sliData.ServiceStarted,
				EndedAt:   sliData.ServiceEnded,
				VatRate:   utils.IfNotNilFloat64(sliData.TaxRate),
				Canceled:  sliData.Canceled,
			}

			onCycleSliEntities = append(onCycleSliEntities, sliEntity)
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

func (s *invoiceService) SimulateOffCycleInvoice(ctx context.Context, contract *neo4jentity.ContractEntity, sliEntities *neo4jentity.ServiceLineItemEntities, span opentracing.Span) (*SimulateInvoiceResponseData, error) {
	invoiceEntity := &neo4jentity.InvoiceEntity{}
	invoiceLines := []*invoicepb.InvoiceLine{}

	tenantSettings, err := s.services.TenantService.GetTenantSettings(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	invoicePeriodEnd := contract.InvoicingStartDate
	invoicePeriodStart := utils.Ptr(utils.FirstTimeOfMonth(9999, 12))
	for _, sliData := range *sliEntities {
		if sliData.StartedAt.Before(*invoicePeriodStart) {
			invoicePeriodStart = utils.Ptr(sliData.StartedAt.Add(time.Hour * 24))
		}
	}

	for true {
		invoicePeriodEnd = utils.Ptr(calculateInvoiceCycleEnd(*invoicePeriodEnd, contract.BillingCycleInMonths))
		if invoicePeriodEnd.After(*invoicePeriodStart) {
			break
		} else {
			invoicePeriodEnd = utils.Ptr(invoicePeriodEnd.AddDate(0, 0, 1))
		}
	}

	invoiceEntity.Number = s.GenerateNewRandomInvoiceNumber()
	invoiceEntity.OffCycle = true
	invoiceEntity.Postpaid = false
	invoiceEntity.BillingCycleInMonths = contract.BillingCycleInMonths
	invoiceEntity.PeriodStartDate = *invoicePeriodStart
	invoiceEntity.PeriodEndDate = *invoicePeriodEnd
	invoiceEntity.IssuedDate = *invoicePeriodStart
	invoiceEntity.DueDate = invoicePeriodEnd.AddDate(0, 0, int(contract.DueDays))

	if contract.Currency != "" {
		invoiceEntity.Currency = contract.Currency
	} else {
		invoiceEntity.Currency = tenantSettings.BaseCurrency
	}

	sliEntitiesForProration := neo4jentity.ServiceLineItemEntities{}
	for _, sliData := range *sliEntities {
		sliEntity := neo4jentity.ServiceLineItemEntity{
			Name:      sliData.Name,
			Comments:  sliData.Comments,
			Billed:    sliData.Billed,
			Price:     sliData.Price,
			Quantity:  sliData.Quantity,
			StartedAt: sliData.StartedAt,
			EndedAt:   nil,
			VatRate:   sliData.VatRate,
		}
		sliEntitiesForProration = append(sliEntitiesForProration, sliEntity)
	}

	invoiceEntity, invoiceLines, err = s.FillOffCyclePrepaidInvoice(ctx, invoiceEntity, sliEntitiesForProration)
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

func (s *invoiceService) PayInvoice(ctx context.Context, invoiceId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceService.PayInvoice")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("invoiceId", invoiceId))

	err := s.UpdateInvoice(ctx, invoiceId, neo4jrepository.InvoiceUpdateFields{
		Status:       neo4jenum.InvoiceStatusPaid,
		UpdateStatus: true,
	})
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error from events processing: %s", err.Error())
		return err
	}
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

	cancelledSliParentIds := []string{}
	for _, sliEntity := range sliEntities {
		if sliEntity.Canceled && sliEntity.ParentID != "" {
			cancelledSliParentIds = append(cancelledSliParentIds, sliEntity.ParentID)
		}
	}

	reasonForSliExcludedFromInvoicing := map[string]string{}

	for _, sliEntity := range sliEntities {
		// skip paused subscriptions
		if sliEntity.Paused {
			reasonForSliExcludedFromInvoicing[sliEntity.ID] = "SLI is paused"
			continue
		}
		// skip for now usage SLIs
		if sliEntity.Billed == neo4jenum.BilledTypeUsage {
			reasonForSliExcludedFromInvoicing[sliEntity.ID] = "Billed type is Usage"
			continue
		}
		// skip SLI if of None type
		if sliEntity.Billed == neo4jenum.BilledTypeNone {
			reasonForSliExcludedFromInvoicing[sliEntity.ID] = "Billed type is None"
			continue
		}
		// skip SLI if ended on the reference time
		if sliEntity.EndedAt != nil && sliEntity.EndedAt.Before(referenceTime) {
			reasonForSliExcludedFromInvoicing[sliEntity.ID] = "SLI ended before reference time"
			continue
		}
		// skip SLI if not active on the reference time
		if sliEntity.IsRecurrent() && !sliEntity.IsActiveAt(referenceTime) {
			reasonForSliExcludedFromInvoicing[sliEntity.ID] = "SLI is not active at reference time"
			continue
		}
		// skip ONE TIME SLI if started after the end period
		if sliEntity.IsOneTime() && sliEntity.StartedAt.After(periodEndTime) {
			reasonForSliExcludedFromInvoicing[sliEntity.ID] = "One time SLI started after the end period"
			continue
		}

		// Any Cancelled SLI should not be invoiced
		if sliEntity.Canceled {
			reasonForSliExcludedFromInvoicing[sliEntity.ID] = "SLI is cancelled"
			continue
		}

		// Any SLI that has any version of same group cancelled should not be invoiced
		if sliEntity.ParentID != "" && utils.Contains(cancelledSliParentIds, sliEntity.ParentID) {
			reasonForSliExcludedFromInvoicing[sliEntity.ID] = "SLI is part of a cancelled group"
			continue
		}

		// skip SLI if quantity or price is negative
		if sliEntity.Quantity < 0 {
			reasonForSliExcludedFromInvoicing[sliEntity.ID] = "Quantity is negative"
			continue
		}

		// skip SLI if price is negative for non one times
		if sliEntity.Price < 0 && !sliEntity.IsOneTime() {
			reasonForSliExcludedFromInvoicing[sliEntity.ID] = "Price is negative"
			continue
		}

		// skip SLI if quantity is zero
		if sliEntity.Quantity == 0 {
			reasonForSliExcludedFromInvoicing[sliEntity.ID] = "Quantity is zero"
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
				reasonForSliExcludedFromInvoicing[sliEntity.ID] = "SLI already invoiced"
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
				Total:                   utils.RoundHalfUpFloat64(calculatedSLIAmount+calculatedSLIVat, 2),
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

	span.LogFields(log.Object("result.ignored_SLIs", reasonForSliExcludedFromInvoicing))

	invoiceEntity.Amount = utils.RoundHalfUpFloat64(amount, 2)
	invoiceEntity.Vat = utils.RoundHalfUpFloat64(vat, 2)
	invoiceEntity.TotalAmount = utils.RoundHalfUpFloat64(amount+vat, 2)

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
			if sliEntity.Quantity <= 0 || sliEntity.Price == 0 {
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
				continue
			}
			finalSLIAmount = utils.RoundHalfUpFloat64(float64(quantity)*sliEntityToInvoice.Price, 2)
			if finalSLIAmount == 0 {
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
			finalSLIAmount = utils.RoundHalfUpFloat64(proratedSLIAmount-proratedInvoicedSLIAmount, 2)
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
			Total:                   utils.RoundHalfUpFloat64(finalSLIAmount+calculatedSLIVat, 2),
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

	invoiceEntity.Amount = utils.RoundHalfUpFloat64(amount, 2)
	invoiceEntity.Vat = utils.RoundHalfUpFloat64(vat, 2)
	invoiceEntity.TotalAmount = utils.RoundHalfUpFloat64(amount+vat, 2)

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

func (s *invoiceService) GetNonDryRunInvoicesForOrganization(ctx context.Context, tenant, organizationId string) (*neo4jentity.InvoiceEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceService.GetNonDryRunInvoicesForOrganization")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	invoiceDbNodes, err := s.services.Neo4jRepositories.InvoiceReadRepository.GetNonDryRunInvoicesForOrganization(ctx, tenant, organizationId)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	invoiceEntities := make(neo4jentity.InvoiceEntities, 0, len(invoiceDbNodes))
	for _, v := range invoiceDbNodes {
		invoiceEntity := mapper.MapDbNodeToInvoiceEntity(v)
		invoiceEntities = append(invoiceEntities, *invoiceEntity)
	}
	return &invoiceEntities, nil
}

func (s *invoiceService) UpdateInvoice(ctx context.Context, invoiceId string, data neo4jrepository.InvoiceUpdateFields) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceService.UpdateInvoice")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	tracing.TagEntity(span, invoiceId)
	span.LogFields(log.String("invoiceId", invoiceId))

	// validate tenant
	err := common.ValidateTenant(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	tenant := common.GetTenantFromContext(ctx)

	invoiceEntityBeforeUpdate, err := s.GetById(ctx, invoiceId)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error while getting invoice %s: %s", invoiceId, err.Error())
		return err
	}

	err = s.services.Neo4jRepositories.InvoiceWriteRepository.UpdateInvoice(ctx, tenant, invoiceId, data)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	invoiceEntityAfterUpdate, err := s.GetById(ctx, invoiceId)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error while getting invoice %s: %s", invoiceId, err.Error())
		return err
	}

	s.createInvoiceAction(ctx, tenant, invoiceEntityBeforeUpdate.Status, *invoiceEntityAfterUpdate)

	// status changed
	if invoiceEntityBeforeUpdate.Status != invoiceEntityAfterUpdate.Status {
		if invoiceEntityAfterUpdate.Status == neo4jenum.InvoiceStatusVoid {
			err = s.VoidInvoice(ctx, invoiceId, common.GetAppSourceFromContext(ctx))
			if err != nil {
				tracing.TraceErr(span, errors.Wrap(err, "void invoice"))
				s.log.Errorf("Error while voiding invoice %s: %s", invoiceId, err.Error())
			}
		}
	}

	return nil
}

func (s *invoiceService) createInvoiceAction(ctx context.Context, tenant string, previousStatus neo4jenum.InvoiceStatus, invoiceEntity neo4jentity.InvoiceEntity) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceEventHandler.createInvoiceAction")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	span.LogFields(log.String("invoiceId", invoiceEntity.Id))
	span.LogFields(log.String("previousStatus", previousStatus.String()))
	span.LogFields(log.String("newStatus", invoiceEntity.Status.String()))
	span.LogFields(log.Bool("dryRun", invoiceEntity.DryRun))
	span.LogFields(log.Float64("totalAmount", invoiceEntity.TotalAmount))

	if previousStatus == invoiceEntity.Status {
		return
	}
	if invoiceEntity.DryRun || invoiceEntity.TotalAmount == float64(0) {
		span.LogFields(log.String("result", "dry run or total amount is 0"))
		return
	}

	metadata, err := utils.ToJson(InvoiceActionMetadata{
		Status:        invoiceEntity.Status.String(),
		Currency:      invoiceEntity.Currency.String(),
		Amount:        invoiceEntity.TotalAmount,
		InvoiceNumber: invoiceEntity.Number,
		InvoiceId:     invoiceEntity.Id,
	})

	actionType := neo4jenum.ActionNA
	message := ""
	switch invoiceEntity.Status {
	case neo4jenum.InvoiceStatusDue:
		message = "Invoice N° " + invoiceEntity.Number + " issued with an amount of " + invoiceEntity.Currency.Symbol() + utils.FormatAmount(invoiceEntity.TotalAmount, 2)
		actionType = neo4jenum.ActionInvoiceIssued
	case neo4jenum.InvoiceStatusPaid:
		message = "Invoice N° " + invoiceEntity.Number + " paid in full: " + invoiceEntity.Currency.Symbol() + utils.FormatAmount(invoiceEntity.TotalAmount, 2)
		actionType = neo4jenum.ActionInvoicePaid
	case neo4jenum.InvoiceStatusVoid:
		message = "Invoice N° " + invoiceEntity.Number + " voided"
		actionType = neo4jenum.ActionInvoiceVoided
	case neo4jenum.InvoiceStatusOverdue:
		message = "Invoice N° " + invoiceEntity.Number + " overdue"
		actionType = neo4jenum.ActionInvoiceOverdue
	}
	if actionType == neo4jenum.ActionNA {
		span.LogFields(log.String("result", "status not supported"))
		return
	}
	if invoiceEntity.Status == neo4jenum.InvoiceStatusDue {
		_, err = s.services.Neo4jRepositories.ActionWriteRepository.MergeByActionType(ctx, nil, tenant, invoiceEntity.Id, model.INVOICE, actionType, message, metadata, utils.Now(), common.GetAppSourceFromContext(ctx))
	} else {
		_, err = s.services.Neo4jRepositories.ActionWriteRepository.Create(ctx, tenant, invoiceEntity.Id, model.INVOICE, actionType, message, metadata, utils.Now(), common.GetAppSourceFromContext(ctx))
	}
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Failed creating invoice action for invoice %s: %s", invoiceEntity.Id, err.Error())
	}
}
