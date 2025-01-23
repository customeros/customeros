package invoice

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/config"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/data_fields"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/dto"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/events"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/postmark"
	webhook "github.com/customeros/customeros/packages/server/customer-os-common-module/webhook_temporal"
	postgresrepository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"
	"github.com/emersion/go-message/mail"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"time"

	neo4jentity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"
	neo4jenum "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/enum"
	neo4jmapper "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/mapper"
	neo4jrepository "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/repository"
	neoRepo "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/repository"
	commonpb "github.com/customeros/customeros/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/common"
	invoicepb "github.com/customeros/customeros/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/invoice"
	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/clients/grpc_client"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/model"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
)

type invoiceService struct {
	log                  logger.Logger
	grpc                 *grpc_client.Clients
	neo4j                *neoRepo.Repositories
	postgresRepositories *postgresrepository.Repositories
	events               *events.EventsService
	contract             interfaces.ContractService
	sli                  interfaces.ServiceLineItemService
	tenantSettings       interfaces.TenantSettingsService
	postmarkService      interfaces.PostmarkService
	fileService          interfaces.FileService
	cfg                  *config.ExternalServicesConfig
}

func NewInvoiceService(log logger.Logger,
	grpc *grpc_client.Clients,
	neo4j *neoRepo.Repositories,
	postgresRepositories *postgresrepository.Repositories,
	cfg *config.ExternalServicesConfig,
	events *events.EventsService,
	contract interfaces.ContractService,
	sli interfaces.ServiceLineItemService,
	tenantSettings interfaces.TenantSettingsService,
	postmarkService interfaces.PostmarkService,
	fileService interfaces.FileService) interfaces.InvoiceService {
	return &invoiceService{
		log:                  log,
		grpc:                 grpc,
		neo4j:                neo4j,
		postgresRepositories: postgresRepositories,
		cfg:                  cfg,
		events:               events,
		contract:             contract,
		sli:                  sli,
		tenantSettings:       tenantSettings,
		postmarkService:      postmarkService,
		fileService:          fileService,
	}
}

type InvoiceActionMetadata struct {
	Status        string  `json:"status"`
	Currency      string  `json:"currency"`
	Amount        float64 `json:"amount"`
	InvoiceNumber string  `json:"number"`
	InvoiceId     string  `json:"id"`
}

func (s *invoiceService) SetContractService(contract interfaces.ContractService) {
	s.contract = contract
}

func (s *invoiceService) SetServiceLineItemService(sli interfaces.ServiceLineItemService) {
	s.sli = sli
}

func (s *invoiceService) IsInitialized() bool {
	return utils.IsInitialized(s)
}

func (s *invoiceService) GenerateNewRandomInvoiceNumber() string {
	digits := "0123456789"
	consonants := "BCDFGHJKLMNPQRSTVWXYZ"
	invoiceNumber := utils.GenerateRandomStringFromCharset(3, consonants) + "-" + utils.GenerateRandomStringFromCharset(5, digits)
	return invoiceNumber
}

func (s *invoiceService) GetById(ctx context.Context, tx *neo4j.ManagedTransaction, invoiceId string) (*neo4jentity.InvoiceEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceService.GetById")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("invoiceId", invoiceId))

	if invoiceDbNode, err := s.neo4j.InvoiceReadRepository.GetInvoiceById(ctx, tx, common.GetTenantFromContext(ctx), invoiceId); err != nil {
		tracing.TraceErr(span, err)
		wrappedErr := errors.Wrap(err, fmt.Sprintf("Invoice with id {%s} not found", invoiceId))
		return nil, wrappedErr
	} else {
		return neo4jmapper.MapDbNodeToInvoiceEntity(invoiceDbNode), nil
	}
}

func (s *invoiceService) GetByIdAcrossAllTenants(ctx context.Context, invoiceId string) (*neo4jentity.InvoiceEntity, string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceService.GetByIdAcrossAllTenants")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogKV("invoiceId", invoiceId)

	invoiceDbNode, tenant, err := s.neo4j.InvoiceReadRepository.GetInvoiceByIdAcrossAllTenants(ctx, invoiceId)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "error getting invoice by id"))
		wrappedErr := errors.Wrap(err, fmt.Sprintf("Invoice with id {%s} not found", invoiceId))
		return nil, "", wrappedErr
	}
	if invoiceDbNode == nil {
		return nil, "", nil
	} else {
		return neo4jmapper.MapDbNodeToInvoiceEntity(invoiceDbNode), tenant, nil
	}
}

func (s *invoiceService) GetByNumber(ctx context.Context, number string) (*neo4jentity.InvoiceEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceService.GetByNumber")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("number", number))

	if invoiceDbNode, err := s.neo4j.InvoiceReadRepository.GetInvoiceByNumber(ctx, common.GetTenantFromContext(ctx), number); err != nil {
		tracing.TraceErr(span, err)
		wrappedErr := errors.Wrap(err, fmt.Sprintf("Invoice with number {%s} not found", number))
		return nil, wrappedErr
	} else {
		return neo4jmapper.MapDbNodeToInvoiceEntity(invoiceDbNode), nil
	}
}

func (s *invoiceService) GetInvoiceLinesForInvoices(ctx context.Context, invoiceIds []string) (*neo4jentity.InvoiceLineEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceService.GetInvoiceLinesForInvoices")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("invoiceIds", invoiceIds))

	invoiceLines, err := s.neo4j.InvoiceLineReadRepository.GetAllForInvoices(ctx, common.GetTenantFromContext(ctx), invoiceIds)
	if err != nil {
		return nil, err
	}
	invoiceLineEntities := make(neo4jentity.InvoiceLineEntities, 0, len(invoiceLines))
	for _, v := range invoiceLines {
		invoiceLineEntity := neo4jmapper.MapDbNodeToInvoiceLineEntity(v.Node)
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

	invoices, err := s.neo4j.InvoiceReadRepository.GetAllForContracts(ctx, common.GetTenantFromContext(ctx), contractIds)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	invoiceEntities := make(neo4jentity.InvoiceEntities, 0, len(invoices))
	for _, v := range invoices {
		invoiceEntity := neo4jmapper.MapDbNodeToInvoiceEntity(v.Node)
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

func (s *invoiceService) SimulateInvoice(ctx context.Context, simulateInvoicesWithChanges *interfaces.SimulateInvoiceRequestData) ([]*interfaces.SimulateInvoiceResponseData, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceService.SimulateInvoice")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("simulateInvoicesWithChanges", simulateInvoicesWithChanges))

	if len(simulateInvoicesWithChanges.ServiceLines) == 0 {
		err := fmt.Errorf("no invoice lines to simulate")
		tracing.TraceErr(span, err)
		return nil, err
	}

	var response []*interfaces.SimulateInvoiceResponseData

	tenantSettings, err := s.tenantSettings.GetTenantSettings(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	// fetch existing contract and set the next invoice date
	contract, err := s.contract.GetById(ctx, simulateInvoicesWithChanges.ContractId)
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

	existingSlis, err := s.sli.GetServiceLineItemsForContract(ctx, contract.Id)
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

			// identify first SLI change in the simualtion and replace it in the sliEntities
			for _, sliData := range simulateInvoicesWithChanges.ServiceLines {

				prorationNeeded := false

				if sliData.ServiceStarted.Before(invoicePeriodStartGeneration.Add(1)) {
					continue
				}

				if sliData.ServiceStarted.After(calculateInvoiceCycleEnd(invoicePeriodStartGeneration, contract.BillingCycleInMonths)) {
					continue
				}

				if sliData.ServiceLineItemID == "" {
					// new sli item - adding it to the sli entities and trigger proration
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
					// existing sli item - to check if there is any change in the sli item to decide if proration is needed
					existingSli, err := s.sli.GetById(ctx, sliData.ServiceLineItemID)
					if err != nil {
						tracing.TraceErr(span, err)
						return nil, err
					}

					if sliData.ServiceStarted.After(utils.Today()) {
						// new version
					} else {
						// update existing version

						if (existingSli.Billed != sliData.BillingCycle || existingSli.Price != sliData.Price || existingSli.Quantity != sliData.Quantity || existingSli.StartedAt != sliData.ServiceStarted || existingSli.VatRate != utils.IfNotNilFloat64(sliData.TaxRate)) && (existingSli.Price*float64(existingSli.Quantity) < sliData.Price*float64(sliData.Quantity)) {
							// existing sli item - new version - proration needed
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

	// no proration needed - only on cycle invoice
	if len(response) == 0 {

		contract.NextInvoiceDate = &invoiceDate

		// build sli entities to reflect the changes for the period
		onCycleSliEntities := neo4jentity.ServiceLineItemEntities{}

		// prepare simulation invoice lines grouped by parent id
		invoiceLinesGroupedByParentId := map[string][]interfaces.SimulateInvoiceRequestServiceLineData{}
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

func (s *invoiceService) SimulateOnCycleInvoice(ctx context.Context, contract *neo4jentity.ContractEntity, sliEntities *neo4jentity.ServiceLineItemEntities, span opentracing.Span) (*interfaces.SimulateInvoiceResponseData, error) {
	invoiceEntity := &neo4jentity.InvoiceEntity{}
	invoiceLines := []*invoicepb.InvoiceLine{}

	tenantSettings, err := s.tenantSettings.GetTenantSettings(ctx)
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

	onCycleInvoice := &interfaces.SimulateInvoiceResponseData{
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

func (s *invoiceService) SimulateOffCycleInvoice(ctx context.Context, contract *neo4jentity.ContractEntity, sliEntities *neo4jentity.ServiceLineItemEntities, span opentracing.Span) (*interfaces.SimulateInvoiceResponseData, error) {
	invoiceEntity := &neo4jentity.InvoiceEntity{}
	invoiceLines := []*invoicepb.InvoiceLine{}

	tenantSettings, err := s.tenantSettings.GetTenantSettings(ctx)
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

	onCycleInvoice := &interfaces.SimulateInvoiceResponseData{
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

	contract, err := s.contract.GetById(ctx, contractId)
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

	tenantSettings, err := s.tenantSettings.GetTenantSettings(ctx)
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
		return s.grpc.InvoiceClient.NewInvoiceForContract(ctx, &dryRunInvoiceRequest)
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

	err := s.UpdateInvoice(ctx, nil, invoiceId, neo4jrepository.InvoiceUpdateFields{
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
		return s.grpc.InvoiceClient.VoidInvoice(ctx, &invoicepb.VoidInvoiceRequest{
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
			result, err := h.neo4j.InvoiceLineReadRepository.GetLatestInvoiceLineWithInvoiceIdByServiceLineItemParentId(ctx, common.GetTenantFromContext(ctx), sliEntity.ParentID)
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

// Deprecated
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
			ilDbNodeAndInvoiceId, err := h.neo4j.InvoiceLineReadRepository.GetLatestInvoiceLineWithInvoiceIdByServiceLineItemParentId(ctx, common.GetTenantFromContext(ctx), sliEntity.ParentID)
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
		ilDbNodeAndInvoiceId, err := h.neo4j.InvoiceLineReadRepository.GetLatestInvoiceLineWithInvoiceIdByServiceLineItemParentId(ctx, common.GetTenantFromContext(ctx), parentId)
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
				previousInvoiceDbNode, err := h.neo4j.InvoiceReadRepository.GetInvoiceById(ctx, nil, common.GetTenantFromContext(ctx), ilDbNodeAndInvoiceId.LinkedNodeId)
				if err != nil {
					tracing.TraceErr(span, err)
					h.log.Errorf("Error getting invoice {%s}: {%s}", ilDbNodeAndInvoiceId.LinkedNodeId, err.Error())
					return nil, nil, err
				}
				previousInvoiceEntity := neo4jmapper.MapDbNodeToInvoiceEntity(previousInvoiceDbNode)
				// if previous invoice is for different cycle, charge full amount
				if !previousInvoiceEntity.PeriodEndDate.Before(invoiceEntity.PeriodEndDate) {
					// calculate already invoiced amount, prorated for the period
					invoiceLineEntity := neo4jmapper.MapDbNodeToInvoiceLineEntity(ilDbNodeAndInvoiceId.Node)
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

	invoiceDbNodes, err := s.neo4j.InvoiceReadRepository.GetNonDryRunInvoicesForOrganization(ctx, tenant, organizationId)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	invoiceEntities := make(neo4jentity.InvoiceEntities, 0, len(invoiceDbNodes))
	for _, v := range invoiceDbNodes {
		invoiceEntity := neo4jmapper.MapDbNodeToInvoiceEntity(v)
		invoiceEntities = append(invoiceEntities, *invoiceEntity)
	}
	return &invoiceEntities, nil
}

func (s *invoiceService) UpdateInvoice(ctx context.Context, txWithPostCommit *utils.TxWithPostCommit, invoiceId string, data neo4jrepository.InvoiceUpdateFields) error {
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

	invoiceEntityBeforeUpdate, err := s.GetById(ctx, nil, invoiceId)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error while getting invoice %s: %s", invoiceId, err.Error())
		return err
	}

	_, err = utils.ExecuteWriteInTransactionWithPostCommitActions(ctx, s.neo4j.Neo4jDriver, s.neo4j.Database, txWithPostCommit, func(txWithPostCommit *utils.TxWithPostCommit) (any, error) {

		err = s.neo4j.InvoiceWriteRepository.UpdateInvoice(ctx, txWithPostCommit.Tx, tenant, invoiceId, data)
		if err != nil {
			return nil, err
		}

		invoiceEntityAfterUpdate, err := s.GetById(ctx, txWithPostCommit.Tx, invoiceId)
		if err != nil {
			s.log.Errorf("Error while getting invoice %s: %s", invoiceId, err.Error())
			return nil, err
		}

		s.createInvoiceAction(ctx, txWithPostCommit.Tx, tenant, invoiceEntityBeforeUpdate.Status, *invoiceEntityAfterUpdate)

		txWithPostCommit.AddPostCommitAction(func(ctx context.Context) error {
			// status changed
			if invoiceEntityBeforeUpdate.Status != invoiceEntityAfterUpdate.Status {
				if invoiceEntityAfterUpdate.Status == neo4jenum.InvoiceStatusVoid {
					err = s.VoidInvoice(ctx, invoiceId, common.GetAppSourceFromContext(ctx))
					if err != nil {
						tracing.TraceErr(span, errors.Wrap(err, "void invoice"))
						s.log.Errorf("Error while voiding invoice %s: %s", invoiceId, err.Error())
					}
				} else if invoiceEntityAfterUpdate.Status == neo4jenum.InvoiceStatusPaid {
					err = s.sendPaidInvoiceNotification(ctx, invoiceId)
					if err != nil {
						tracing.TraceErr(span, errors.Wrap(err, "send paid invoice notification"))
					}
				}
			}

			return nil
		})

		txWithPostCommit.AddPostCommitAction(func(ctx context.Context) error {
			// status changed
			if invoiceEntityBeforeUpdate.Status != invoiceEntityAfterUpdate.Status {
				if invoiceEntityAfterUpdate.Status == neo4jenum.InvoiceStatusPaid {
					// do not dispatch invoice paid event if it was already dispatched
					if invoiceEntityAfterUpdate.InvoiceInternalFields.InvoicePaidWebhookProcessedAt == nil {
						// dispatch invoice paid event
						err = s.dispatchInvoicePaidEvent(ctx, tenant, *invoiceEntityAfterUpdate)
						if err != nil {
							tracing.TraceErr(span, errors.Wrap(err, "dispatchInvoicePaidEvent"))
							s.log.Errorf("Error dispatching invoice paid event for invoice %s: %s", invoiceId, err.Error())
						} else {
							err = s.neo4j.CommonWriteRepository.UpdateTimeProperty(ctx, tenant, model.NodeLabelInvoice, invoiceEntityAfterUpdate.Id, string(neo4jentity.InvoicePropertyPaidWebhookProcessedAt), utils.NowPtr())
							if err != nil {
								tracing.TraceErr(span, errors.Wrap(err, "UpdateTimeProperty"))
								s.log.Errorf("Error setting invoice paid webhook processed for invoice %s: %s", invoiceEntityAfterUpdate.Id, err.Error())
							}
						}
					}
				}
			}

			return nil
		})

		txWithPostCommit.AddPostCommitAction(func(ctx context.Context) error {
			updateDto := dto.UpdateInvoice{}
			if data.UpdateStatus {
				updateDto.Status = utils.StringPtr(data.Status.String())
			}
			if data.UpdatePaymentLink {
				updateDto.PaymentLink = utils.StringPtr(data.PaymentLink)
				updateDto.PaymentLinkValidUntil = data.PaymentLinkValidUntil
			}

			err = s.events.Publisher.PublishEvent(ctx, invoiceId, model.INVOICE, updateDto)
			if err != nil {
				tracing.TraceErr(span, errors.Wrap(err, "unable to publish message UpdateInvoice"))
			}

			return nil
		})

		return nil, nil
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}

func (s *invoiceService) createInvoiceAction(ctx context.Context, tx *neo4j.ManagedTransaction, tenant string, previousStatus neo4jenum.InvoiceStatus, invoiceEntity neo4jentity.InvoiceEntity) {
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

	actionType := enum.ActionNA
	message := ""
	switch invoiceEntity.Status {
	case neo4jenum.InvoiceStatusDue:
		message = "Invoice N° " + invoiceEntity.Number + " issued with an amount of " + invoiceEntity.Currency.Symbol() + utils.FormatAmount(invoiceEntity.TotalAmount, 2)
		actionType = enum.ActionInvoiceIssued
	case neo4jenum.InvoiceStatusPaid:
		message = "Invoice N° " + invoiceEntity.Number + " paid in full: " + invoiceEntity.Currency.Symbol() + utils.FormatAmount(invoiceEntity.TotalAmount, 2)
		actionType = enum.ActionInvoicePaid
	case neo4jenum.InvoiceStatusVoid:
		message = "Invoice N° " + invoiceEntity.Number + " voided"
		actionType = enum.ActionInvoiceVoided
	case neo4jenum.InvoiceStatusOverdue:
		message = "Invoice N° " + invoiceEntity.Number + " overdue"
		actionType = enum.ActionInvoiceOverdue
	}
	if actionType == enum.ActionNA {
		span.LogFields(log.String("result", "status not supported"))
		return
	}
	if invoiceEntity.Status == neo4jenum.InvoiceStatusDue {
		_, err = s.neo4j.ActionWriteRepository.MergeByActionType(ctx, tx, tenant, invoiceEntity.Id, model.INVOICE, actionType, message, metadata, utils.Now(), common.GetAppSourceFromContext(ctx))
	} else {
		actionId, err := s.neo4j.CommonReadRepository.GenerateId(ctx, tenant, model.NodeLabelAction)
		if err != nil {
			tracing.TraceErr(span, err)
			return
		}
		actionFields := data_fields.ActionFields{
			ActionType: utils.ToPtr(actionType),
			AppSource:  utils.StringPtr(common.GetAppSourceFromContext(ctx)),
			Content:    utils.StringPtr(message),
			Metadata:   utils.StringPtr(metadata),
		}
		err = s.neo4j.ActionWriteRepository.CreateV2(ctx, tx, tenant, actionId, invoiceEntity.Id, model.INVOICE, actionFields)
	}
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Failed creating invoice action for invoice %s: %s", invoiceEntity.Id, err.Error())
	}
}

func (s *invoiceService) sendPaidInvoiceNotification(ctx context.Context, invoiceId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceService.sendPaidInvoiceNotification")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	tracing.TagEntity(span, invoiceId)

	tenant := common.GetTenantFromContext(ctx)
	eventTriggeredByUser := common.GetUserIdFromContext(ctx) != ""

	var invoiceEntity *neo4jentity.InvoiceEntity
	var contractEntity neo4jentity.ContractEntity

	//load invoice details
	invoiceEntity, err := s.GetById(ctx, nil, invoiceId)
	if err != nil {
		return nil
	}

	if invoiceEntity.DryRun || invoiceEntity.TotalAmount == float64(0) {
		return nil
	}

	// paid notification already sent, skip
	if invoiceEntity.InvoiceInternalFields.PaidInvoiceNotificationSentAt != nil {
		return nil
	}

	// load contract
	contractNode, err := s.neo4j.ContractReadRepository.GetContractForInvoice(ctx, tenant, invoiceEntity.Id)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "GetContractForInvoice"))
		return errors.Wrap(err, "InvoiceSubscriber.onInvoicePaidV1.GetContractForInvoice")
	}
	if contractNode != nil {
		contractEntity = *neo4jmapper.MapDbNodeToContractEntity(contractNode)
	} else {
		tracing.TraceErr(span, errors.New("contractNode is nil"))
		return errors.New("contractNode is nil")
	}

	if contractEntity.InvoiceEmail == "" || !isValidEmailSyntax(contractEntity.InvoiceEmail) {
		tracing.TraceErr(span, errors.New("contractEntity.InvoiceEmail is empty or invalid"))
		return nil
	}

	//load tenant billing profile from neo4j
	tenantBillingProfileEntity, err := s.loadTenantBillingProfile(ctx, tenant, false)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "loadTenantBillingProfile"))
		return nil
	}
	tenantSettingsDbNode, err := s.neo4j.TenantReadRepository.GetTenantSettings(ctx, tenant)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "GetTenantSettings"))
		return err
	}
	tenantSettingsEntity := neo4jmapper.MapDbNodeToTenantSettingsEntity(tenantSettingsDbNode)

	cc := contractEntity.InvoiceEmailCC
	cc = utils.RemoveEmpties(cc)
	cc = utils.RemoveDuplicates(cc)

	bcc := utils.AddToListIfNotExists(contractEntity.InvoiceEmailBCC, tenantBillingProfileEntity.SendInvoicesBcc)
	bcc = utils.RemoveEmpties(bcc)
	bcc = utils.RemoveDuplicates(bcc)

	postmarkEmail := interfaces.PostmarkEmail{
		MessageStream: postmark.PostmarkMessageStreamInvoice,
		From:          invoiceEntity.Provider.Email,
		To:            contractEntity.InvoiceEmail,
		CC:            cc,
		BCC:           bcc,
		TemplateData: map[string]string{
			"{{userFirstName}}":  invoiceEntity.Customer.Name,
			"{{invoiceNumber}}":  invoiceEntity.Number,
			"{{currencySymbol}}": invoiceEntity.Currency.Symbol(),
			"{{amtDue}}":         fmt.Sprintf("%.2f", invoiceEntity.TotalAmount),
			"{{paymentDate}}":    utils.Now().Format("02 Jan 2006"),
		},
		Attachments: []interfaces.PostmarkEmailAttachment{},
	}
	if tenantSettingsEntity.StripeCustomerPortalLink != "" {
		postmarkEmail.TemplateData["{{stripeFooterHtml}}"] = fmt.Sprintf(`PS: If you pay by card you can manage your billing details <a href="%s">here</a>.`, tenantSettingsEntity.StripeCustomerPortalLink)
		postmarkEmail.TemplateData["{{stripeFooterTxt}}"] = `PS: If you pay by card you can manage your billing details here.`
		postmarkEmail.TemplateData["{{stripeFooterLink}}"] = tenantSettingsEntity.StripeCustomerPortalLink
	} else {
		postmarkEmail.TemplateData["{{stripeFooterHtml}}"] = ""
		postmarkEmail.TemplateData["{{stripeFooterTxt}}"] = ""
		postmarkEmail.TemplateData["{{stripeFooterLink}}"] = ""
	}

	if eventTriggeredByUser {
		postmarkEmail.WorkflowId = postmark.WorkflowInvoicePaymentReceived
		postmarkEmail.Subject = fmt.Sprintf(postmark.WorkflowInvoicePaymentReceivedSubject, invoiceEntity.Number, invoiceEntity.Provider.Name)
	} else {
		postmarkEmail.WorkflowId = postmark.WorkflowInvoicePaid
		postmarkEmail.Subject = fmt.Sprintf(postmark.WorkflowInvoicePaidSubject, invoiceEntity.Number, invoiceEntity.Provider.Name)
	}

	err = s.appendInvoiceFileToEmailAsAttachment(ctx, tenant, *invoiceEntity, &postmarkEmail)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "appendInvoiceFileToEmailAsAttachment"))
		s.log.Errorf("Error appending invoice file to email attachment for invoice %s: %s", invoiceId, err.Error())
		return nil
	}

	err = s.appendProviderLogoToEmail(ctx, tenant, invoiceEntity.Provider.LogoRepositoryFileId, &postmarkEmail)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "appendProviderLogoToEmail"))
		s.log.Errorf("Error appending provider logo to email for invoice %s: %s", invoiceId, err.Error())
		return nil
	}

	err = s.appendCustomerOSLogoToEmail(ctx, &postmarkEmail)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "appendCustomerOSLogoToEmail"))
		s.log.Errorf("Error appending customeros logo to email for invoice %s: %s", invoiceId, err.Error())
		return nil
	}

	err = s.postmarkService.SendNotification(ctx, postmarkEmail, tenant)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "SendNotification"))
		s.log.Errorf("Error sending invoice paid notification for invoice %s: %s", invoiceId, err.Error())
		return nil
	}

	// Request was successful
	err = s.neo4j.InvoiceWriteRepository.SetPaidInvoiceNotificationSentAt(ctx, tenant, invoiceId)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "SetPaidInvoiceNotificationSentAt"))
		s.log.Errorf("Error setting invoice paid notification sent at for invoice %s: %s", invoiceId, err.Error())
		return nil
	}

	return nil
}

func (s *invoiceService) appendCustomerOSLogoToEmail(ctx context.Context, postmarkEmail *interfaces.PostmarkEmail) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceService.appendCustomerOSLogoToEmail")
	defer span.Finish()

	currentDir, err := os.Getwd()
	if err != nil {
		return errors.Wrap(err, "os.Getwd")
	}

	file, err := utils.GetFileByName(filepath.Join(currentDir, "/static", "customer-os.png"))

	b, err := ioutil.ReadAll(file)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "ioutil.ReadAll"))
		return err
	}

	postmarkEmail.Attachments = append(postmarkEmail.Attachments, interfaces.PostmarkEmailAttachment{
		Filename:       "customer-os-encoded",
		ContentEncoded: base64.StdEncoding.EncodeToString(b),
		ContentType:    "image/png",
		ContentID:      "cid:customer-os-encoded",
	})

	return nil
}

func (s *invoiceService) appendProviderLogoToEmail(ctx context.Context, tenant, logoFileId string, postmarkEmail *interfaces.PostmarkEmail) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceService.appendProviderLogoToEmail")
	defer span.Finish()

	if logoFileId == "" {
		return nil
	}

	fileInfo, err := s.fileService.GetById(ctx, logoFileId)
	if err != nil {
		return err
	}
	if fileInfo == nil {
		return nil
	}

	fileBytes, err := s.fileService.GetFileBytes(ctx, logoFileId)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	postmarkEmail.Attachments = append(postmarkEmail.Attachments, interfaces.PostmarkEmailAttachment{
		Filename:       "provider-logo-file-encoded",
		ContentEncoded: base64.StdEncoding.EncodeToString(*fileBytes),
		ContentType:    fileInfo.MimeType,
		ContentID:      "cid:provider-logo-file-encoded",
	})

	return nil
}

func isValidEmailSyntax(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func (s *invoiceService) loadTenantBillingProfile(ctx context.Context, tenant string, failIfNotFound bool) (neo4jentity.TenantBillingProfileEntity, error) {
	tenantBillingProfiles, err := s.neo4j.TenantReadRepository.GetTenantBillingProfiles(ctx, tenant)
	if err != nil {
		return neo4jentity.TenantBillingProfileEntity{}, err
	}
	if len(tenantBillingProfiles) == 0 {
		if failIfNotFound {
			return neo4jentity.TenantBillingProfileEntity{}, errors.New("tenantBillingProfiles not available")
		} else {
			return neo4jentity.TenantBillingProfileEntity{}, nil
		}
	}
	tenantBillingProfileEntity := neo4jmapper.MapDbNodeToTenantBillingProfileEntity(tenantBillingProfiles[0])
	return *tenantBillingProfileEntity, nil
}

func (s *invoiceService) appendInvoiceFileToEmailAsAttachment(ctx context.Context, tenant string, invoice neo4jentity.InvoiceEntity, postmarkEmail *interfaces.PostmarkEmail) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceService.appendInvoiceFileToEmailAsAttachment")
	defer span.Finish()
	tracing.TagTenant(span, tenant)

	fileInfo, err := s.fileService.GetById(ctx, invoice.RepositoryFileId)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	if fileInfo == nil {
		return nil
	}

	invoiceFileBytes, err := s.fileService.GetFileBytes(ctx, invoice.RepositoryFileId)
	if err != nil {
		return err
	}

	postmarkEmail.Attachments = append(postmarkEmail.Attachments, interfaces.PostmarkEmailAttachment{
		Filename:       "Invoice " + invoice.Number + ".pdf",
		ContentEncoded: base64.StdEncoding.EncodeToString(*invoiceFileBytes),
		ContentType:    "application/pdf",
	})

	return nil
}

func (s *invoiceService) dispatchInvoicePaidEvent(ctx context.Context, tenant string, invoice neo4jentity.InvoiceEntity) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "InvoiceService.dispatchInvoicePaidEvent")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	tracing.LogObjectAsJson(span, "invoice", invoice)

	// get organization linked to invoice to build payload for webhook
	organizationDbNode, err := s.neo4j.OrganizationReadRepository.GetOrganizationByInvoiceId(ctx, tenant, invoice.Id)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "GetOrganizationByInvoiceId"))
		s.log.Errorf("Error getting organization for invoice %s: %s", invoice.Id, err.Error())
		return err
	}
	organizationEntity := neo4jentity.OrganizationEntity{}
	if organizationDbNode != nil {
		organizationEntity = *neo4jmapper.MapDbNodeToOrganizationEntity(organizationDbNode)
	}

	// get contract linked to invoice to build payload for webhook
	contractDbNode, err := s.neo4j.ContractReadRepository.GetContractsForOrganizations(ctx, tenant, []string{organizationEntity.ID})
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "GetContractsForOrganizations"))
		s.log.Errorf("Error getting contract for invoice %s: %s", invoice.Id, err.Error())
		return err
	}

	contractEntity := neo4jentity.ContractEntity{}
	if len(contractDbNode) > 0 && contractDbNode[0] != nil {
		node := contractDbNode[0].Node
		if node != nil {
			contractEntity = *neo4jmapper.MapDbNodeToContractEntity(node)
		}
	}

	// get invoice line items linked to invoice to build payload for webhook
	invoiceLineDbNodes, err := s.neo4j.InvoiceLineReadRepository.GetAllForInvoice(ctx, tenant, invoice.Id)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "GetAllForInvoice"))
		s.log.Errorf("Error getting invoice line items for invoice %s: %s", invoice.Id, err.Error())
		return err
	}

	ilEntities := []*neo4jentity.InvoiceLineEntity{}
	for _, ilDbNode := range invoiceLineDbNodes {
		ilEntity := neo4jmapper.MapDbNodeToInvoiceLineEntity(ilDbNode)
		ilEntities = append(ilEntities, ilEntity)
	}

	webhookPayload := webhook.PopulateInvoicePayload(&invoice, &organizationEntity, &contractEntity, ilEntities)
	// dispatch the event
	err = webhook.DispatchWebhook(
		ctx,
		tenant,
		webhook.WebhookEventInvoiceStatusPaid,
		webhookPayload,
		s.postgresRepositories,
		*s.cfg,
	)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "DispatchWebhook"))
		s.log.Errorf("Error dispatching invoice paid event for invoice %s: %s", invoice.Id, err.Error())
		return err
	}

	return nil
}
