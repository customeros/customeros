package invoice

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/mail"
	"os"
	"sort"
	"time"

	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/subscriptions"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/data"
	fsc "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/file_store_client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/neo4jutil"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/grpc_client"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/notifications"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/webhook"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/invoice"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	commonpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/common"
	invoicepb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/invoice"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

type RequestBodyInvoiceReady struct {
	Tenant                       string `json:"tenant"`
	Currency                     string `json:"currency"`
	AmountInSmallestCurrencyUnit int64  `json:"amountInSmallestCurrencyUnit"`
	StripeCustomerId             string `json:"stripeCustomerId"`
	InvoiceId                    string `json:"invoiceId"`
	InvoiceDescription           string `json:"invoiceDescription"`
	CustomerOsId                 string `json:"customerOsId"`
}

type InvoiceEventHandler struct {
	log              logger.Logger
	repositories     *repository.Repositories
	cfg              config.Config
	grpcClients      *grpc_client.Clients
	fsc              fsc.FileStoreApiService
	postmarkProvider *notifications.PostmarkProvider
}

func NewInvoiceEventHandler(log logger.Logger, repositories *repository.Repositories, cfg config.Config, grpcClients *grpc_client.Clients, fsc fsc.FileStoreApiService, postmarkProvider *notifications.PostmarkProvider) *InvoiceEventHandler {
	return &InvoiceEventHandler{
		log:              log,
		repositories:     repositories,
		cfg:              cfg,
		grpcClients:      grpcClients,
		fsc:              fsc,
		postmarkProvider: postmarkProvider,
	}
}

func (h *InvoiceEventHandler) onInvoiceFillRequestedV1(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceEventHandler.onInvoiceFillRequestedV1")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData invoice.InvoiceFillRequestedEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}
	tracing.LogObjectAsJson(span, "eventData", eventData)
	invoiceId := invoice.GetInvoiceObjectID(evt.GetAggregateID(), eventData.Tenant)
	span.SetTag(tracing.SpanTagEntityId, invoiceId)

	invoiceDbNode, err := h.repositories.Neo4jRepositories.InvoiceReadRepository.GetInvoiceById(ctx, eventData.Tenant, invoiceId)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("error getting invoice {%s}: {%s}", invoiceId, err.Error())
		return err
	}
	if invoiceDbNode == nil {
		err = errors.Errorf("invoice {%s} not found", invoiceId)
		tracing.TraceErr(span, err)
		h.log.Errorf("error getting invoice {%s}: {%s}", invoiceId, err.Error())
		return err
	}
	invoiceEntity := neo4jmapper.MapDbNodeToInvoiceEntity(invoiceDbNode)

	if invoiceEntity.OffCycle {
		return h.fillOffCyclePrepaidInvoice(ctx, eventData.Tenant, eventData.ContractId, *invoiceEntity)
	} else {
		return h.fillCycleInvoice(ctx, eventData.Tenant, eventData.ContractId, *invoiceEntity)
	}
}

func (h *InvoiceEventHandler) fillCycleInvoice(ctx context.Context, tenant, contractId string, invoiceEntity neo4jentity.InvoiceEntity) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceEventHandler.fillCycleInvoice")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	span.SetTag(tracing.SpanTagEntityId, invoiceEntity.Id)
	span.LogFields(log.String("contractId", contractId))

	sliDbNodes, err := h.repositories.Neo4jRepositories.ServiceLineItemReadRepository.GetAllForContract(ctx, tenant, contractId)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error getting service line items for contract %s: %s", contractId, err.Error())
		return err
	}

	var sliEntities neo4jentity.ServiceLineItemEntities
	for _, sliDbNode := range sliDbNodes {
		sliEntity := neo4jmapper.MapDbNodeToServiceLineItemEntity(sliDbNode)
		if sliEntity != nil {
			sliEntities = append(sliEntities, *neo4jmapper.MapDbNodeToServiceLineItemEntity(sliDbNode))
		}
	}

	amount, vat, totalAmount := float64(0), float64(0), float64(0)
	var invoiceLines []*invoicepb.InvoiceLine

	referenceTime := invoiceEntity.PeriodStartDate
	if invoiceEntity.Postpaid {
		referenceTime = utils.EndOfDayInUTC(invoiceEntity.PeriodEndDate)
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
		if !sliEntity.IsActiveAt(referenceTime) {
			continue
		}

		if sliEntity.Quantity <= 0 {
			continue
		}

		calculatedSLIAmount, calculatedSLIVat := float64(0), float64(0)
		invoiceLineCalculationsReady := false
		// process one time SLIs
		if sliEntity.Billed == neo4jenum.BilledTypeOnce {
			// Check any version of SLI not invoiced
			result, err := h.repositories.Neo4jRepositories.InvoiceLineReadRepository.GetLatestInvoiceLineWithInvoiceIdByServiceLineItemParentId(ctx, tenant, sliEntity.ParentID)
			if err != nil {
				tracing.TraceErr(span, err)
				h.log.Errorf("Error getting latest invoice line for sli parent id {%s}: {%s}", sliEntity.ParentID, err.Error())
			}
			if result != nil {
				// SLI already invoiced
				continue
			}
			quantity := sliEntity.Quantity
			if sliEntity.Quantity <= 0 {
				quantity = 1
			}
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
		err = errors.Errorf("Unprocessed SLI %s", sliEntity.ID)
		tracing.TraceErr(span, err)
		h.log.Errorf("Error processing SLI during invoicing %s: %s", sliEntity.ID, err.Error())
	}

	if len(invoiceLines) == 0 {
		return errors.Wrap(err, "No invoice lines to invoice")
	}

	totalAmount = amount + vat

	return h.prepareAndCallFillInvoice(ctx, tenant, contractId, invoiceEntity, amount, vat, totalAmount, invoiceLines, span)
}

func (h *InvoiceEventHandler) fillOffCyclePrepaidInvoice(ctx context.Context, tenant, contractId string, invoiceEntity neo4jentity.InvoiceEntity) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceEventHandler.fillOffCyclePrepaidInvoice")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	span.SetTag(tracing.SpanTagEntityId, invoiceEntity.Id)
	span.LogFields(log.String("contractId", contractId))

	sliDbNodes, err := h.repositories.Neo4jRepositories.ServiceLineItemReadRepository.GetAllForContract(ctx, tenant, contractId)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error getting service line items for contract %s: %s", contractId, err.Error())
		return err
	}

	var sliEntities neo4jentity.ServiceLineItemEntities
	for _, sliDbNode := range sliDbNodes {
		sliEntity := neo4jmapper.MapDbNodeToServiceLineItemEntity(sliDbNode)
		if sliEntity != nil {
			sliEntities = append(sliEntities, *neo4jmapper.MapDbNodeToServiceLineItemEntity(sliDbNode))
		}
	}
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
		// One time invoiced SLIs are not applicable
		if sliEntity.Billed == neo4jenum.BilledTypeOnce {
			ilDbNodeAndInvoiceId, err := h.repositories.Neo4jRepositories.InvoiceLineReadRepository.GetLatestInvoiceLineWithInvoiceIdByServiceLineItemParentId(ctx, tenant, sliEntity.ParentID)
			if err != nil {
				tracing.TraceErr(span, err)
				h.log.Errorf("Error getting latest invoice line for sli parent id {%s}: {%s}", sliEntity.ParentID, err.Error())
				return err
			}
			if ilDbNodeAndInvoiceId != nil {
				continue
			}
		}

		if sliEntity.Quantity <= 0 {
			continue
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

	amount, vat, totalAmount := float64(0), float64(0), float64(0)
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
		ilDbNodeAndInvoiceId, err := h.repositories.Neo4jRepositories.InvoiceLineReadRepository.GetLatestInvoiceLineWithInvoiceIdByServiceLineItemParentId(ctx, tenant, parentId)
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("Error getting latest invoice line for sli parent id {%s}: {%s}", parentId, err.Error())
			return err
		}
		finalSLIAmount, calculatedSLIVat := float64(0), float64(0)
		if sliEntityToInvoice.Billed == neo4jenum.BilledTypeOnce {
			quantity := sliEntityToInvoice.Quantity
			if quantity <= 0 {
				quantity = 1
			}
			finalSLIAmount = utils.TruncateFloat64(float64(quantity)*sliEntityToInvoice.Price, 2)
			if finalSLIAmount <= 0 {
				continue
			}
			calculatedSLIVat = utils.TruncateFloat64(finalSLIAmount*sliEntityToInvoice.VatRate/100, 2)
		} else {
			proratedInvoicedSLIAmount := float64(0)
			if ilDbNodeAndInvoiceId != nil {
				previousInvoiceDbNode, err := h.repositories.Neo4jRepositories.InvoiceReadRepository.GetInvoiceById(ctx, tenant, ilDbNodeAndInvoiceId.LinkedNodeId)
				if err != nil {
					tracing.TraceErr(span, err)
					h.log.Errorf("Error getting invoice {%s}: {%s}", ilDbNodeAndInvoiceId.LinkedNodeId, err.Error())
					return err
				}
				previousInvoiceEntity := neo4jmapper.MapDbNodeToInvoiceEntity(previousInvoiceDbNode)
				// if previous invoice is for different cycle, charge full amount
				if !previousInvoiceEntity.PeriodEndDate.Before(invoiceEntity.PeriodEndDate) {
					// calculate already invoiced amount, prorated for the period
					invoiceLineEntity := neo4jmapper.MapDbNodeToInvoiceLineEntity(ilDbNodeAndInvoiceId.Node)
					calculatedInvoicedSLIAmountFor1Year := calculateSLIAmountForCycleInvoicing(invoiceLineEntity.Quantity, invoiceLineEntity.Price, invoiceLineEntity.BilledType, neo4jenum.BillingCycleAnnuallyBilling)
					proratedInvoicedSLIAmount = prorateAnnualSLIAmount(sliEntityToInvoice.StartedAt, invoiceEntity.PeriodEndDate, calculatedInvoicedSLIAmountFor1Year)
					proratedInvoicedSLIAmount = utils.TruncateFloat64(proratedInvoicedSLIAmount, 2)
				}
			}

			calculatedSLIAmountFor1Year := calculateSLIAmountForCycleInvoicing(sliEntityToInvoice.Quantity, sliEntityToInvoice.Price, sliEntityToInvoice.Billed, neo4jenum.BillingCycleAnnuallyBilling)
			proratedSLIAmount := prorateAnnualSLIAmount(sliEntityToInvoice.StartedAt, invoiceEntity.PeriodEndDate, calculatedSLIAmountFor1Year)
			proratedSLIAmount = utils.TruncateFloat64(proratedSLIAmount, 2)
			finalSLIAmount = proratedSLIAmount - proratedInvoicedSLIAmount
			span.LogFields(log.Float64(fmt.Sprintf("result - final amount for SLI with parent id %s", parentId), finalSLIAmount))
			if finalSLIAmount <= 0 {
				continue
			}
			calculatedSLIVat = utils.TruncateFloat64(finalSLIAmount*sliEntityToInvoice.VatRate/100, 2)
			proratedSliFound = true
		}
		amount += finalSLIAmount
		vat += calculatedSLIVat
		invoiceLine := invoicepb.InvoiceLine{
			Name:                    sliEntityToInvoice.Name,
			Price:                   utils.TruncateFloat64(calculatePriceForBilledType(sliEntityToInvoice.Price, sliEntityToInvoice.Billed, invoiceEntity.BillingCycle), 2),
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
	totalAmount = amount + vat

	if !proratedSliFound && len(invoiceLines) > 0 {
		// if no prorated SLI found, then invoice contains only once billed SLIs
		// accept the invoice if today is monthly anniversary of the contract invoicing start date

		// UPDATE: The rule is on hold, invoice will be issued even if contains only one time SLIs

		//if !isMonthlyAnniversary(invoiceEntity.PeriodEndDate.AddDate(0, 0, 1)) {
		//	invoiceLines = []*invoicepb.InvoiceLine{}
		//}
	}

	if totalAmount == 0 || len(invoiceLines) == 0 {
		_, err = subscriptions.CallEventsPlatformGRPCWithRetry[*invoicepb.InvoiceIdResponse](func() (*invoicepb.InvoiceIdResponse, error) {
			return h.grpcClients.InvoiceClient.PermanentlyDeleteDraftInvoice(ctx, &invoicepb.PermanentlyDeleteDraftInvoiceRequest{
				Tenant:    tenant,
				InvoiceId: invoiceEntity.Id,
				AppSource: constants.AppSourceEventProcessingPlatform,
			})
		})
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("Error permanently deleting draft invoice {%s}: {%s}", invoiceEntity.Id, err.Error())
		}
		return err
	} else {
		return h.prepareAndCallFillInvoice(ctx, tenant, contractId, invoiceEntity, amount, vat, totalAmount, invoiceLines, span)
	}
}

func isMonthlyAnniversary(date time.Time) bool {
	now := utils.Now()
	if now.Day() == date.Day() {
		return true
	}
	return false
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

func calculateSLIAmountForCycleInvoicing(quantity int64, price float64, billed neo4jenum.BilledType, cycle neo4jenum.BillingCycle) float64 {
	sliAmount := float64(quantity) * price
	if sliAmount == 0 {
		return sliAmount
	}

	return calculatePriceForBilledType(sliAmount, billed, cycle)
}

func prorateAnnualSLIAmount(startDate, endDate time.Time, amount float64) float64 {
	start := utils.StartOfDayInUTC(startDate)
	end := utils.StartOfDayInUTC(endDate)
	days := end.Sub(start).Hours() / 24
	proratedAmount := amount * (days / 365)
	if proratedAmount <= 0 {
		return 0
	}
	return proratedAmount
}

func (h *InvoiceEventHandler) prepareAndCallFillInvoice(ctx context.Context, tenant string, contractId string, invoiceEntity neo4jentity.InvoiceEntity, amount, vat, totalAmount float64, invoiceLines []*invoicepb.InvoiceLine, span opentracing.Span) error {
	var contractEntity neo4jentity.ContractEntity
	var tenantSettingsEntity *neo4jentity.TenantSettingsEntity

	//load contract from neo4j
	contract, err := h.repositories.Neo4jRepositories.ContractReadRepository.GetContractById(ctx, tenant, contractId)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "InvoiceSubscriber.onInvoiceFillV1.GetContractById")
	}
	if contract != nil {
		contractEntity = *neo4jmapper.MapDbNodeToContractEntity(contract)
	} else {
		return errors.New("contract is nil")
	}

	//load tenant settings from neo4j
	tenantSettings, err := h.repositories.Neo4jRepositories.TenantReadRepository.GetTenantSettings(ctx, tenant)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	if tenantSettings != nil {
		tenantSettingsEntity = neo4jmapper.MapDbNodeToTenantSettingsEntity(tenantSettings)
	} else {
		tracing.TraceErr(span, errors.New("tenantSettings is nil"))
		return errors.New("tenantSettings is nil")
	}

	//load tenant billing profile from neo4j
	tenantBillingProfileEntity, err := h.loadTenantBillingProfile(ctx, tenant, true)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	contractCountry := contractEntity.Country
	countryDbNode, _ := h.repositories.Neo4jRepositories.CountryReadRepository.GetCountryByCodeIfExists(ctx, contractCountry)
	if countryDbNode != nil {
		countryEntity := neo4jmapper.MapDbNodeToCountryEntity(countryDbNode)
		contractCountry = countryEntity.Name
	}
	tenantBillingProfileCountry := tenantBillingProfileEntity.Country
	countryDbNode, _ = h.repositories.Neo4jRepositories.CountryReadRepository.GetCountryByCodeIfExists(ctx, tenantBillingProfileCountry)
	if countryDbNode != nil {
		countryEntity := neo4jmapper.MapDbNodeToCountryEntity(countryDbNode)
		tenantBillingProfileCountry = countryEntity.Name
	}

	err = h.callFillInvoice(ctx,
		tenant,
		invoiceEntity.Id,
		tenantBillingProfileEntity.DomesticPaymentsBankInfo,
		tenantBillingProfileEntity.InternationalPaymentsBankInfo,
		contractEntity.OrganizationLegalName,
		contractEntity.InvoiceEmail,
		contractEntity.AddressLine1, contractEntity.AddressLine2, contractEntity.Zip, contractEntity.Locality, contractCountry,
		tenantSettingsEntity.LogoRepositoryFileId,
		tenantBillingProfileEntity.LegalName,
		tenantBillingProfileEntity.SendInvoicesFrom,
		tenantBillingProfileEntity.AddressLine1, tenantBillingProfileEntity.AddressLine2, tenantBillingProfileEntity.Zip, tenantBillingProfileEntity.Locality, tenantBillingProfileCountry,
		contractEntity.InvoiceNote,
		amount,
		vat,
		totalAmount,
		invoiceLines,
		span)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	return nil
}

func (h *InvoiceEventHandler) callFillInvoice(ctx context.Context, tenant, invoiceId, domesticPaymentsBankInfo, internationalPaymentsBankInfo,
	customerName, customerEmail, customerAddressLine1, customerAddressLine2, customerAddressZip, customerAddressLocality, customerAddressCountry,
	providerLogoRepositoryFileId, providerName, providerEmail, providerAddressLine1, providerAddressLine2, providerAddressZip, providerAddressLocality, providerAddressCountry,
	note string, amount, vat, total float64, invoiceLines []*invoicepb.InvoiceLine, span opentracing.Span) error {
	ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
	now := time.Now()
	_, err := subscriptions.CallEventsPlatformGRPCWithRetry[*invoicepb.InvoiceIdResponse](func() (*invoicepb.InvoiceIdResponse, error) {
		return h.grpcClients.InvoiceClient.FillInvoice(ctx, &invoicepb.FillInvoiceRequest{
			Tenant:                        tenant,
			InvoiceId:                     invoiceId,
			Note:                          note,
			DomesticPaymentsBankInfo:      domesticPaymentsBankInfo,
			InternationalPaymentsBankInfo: internationalPaymentsBankInfo,
			Customer: &invoicepb.FillInvoiceCustomer{
				Name:         customerName,
				Email:        customerEmail,
				AddressLine1: customerAddressLine1,
				AddressLine2: customerAddressLine2,
				Zip:          customerAddressZip,
				Locality:     customerAddressLocality,
				Country:      customerAddressCountry,
			},
			Provider: &invoicepb.FillInvoiceProvider{
				LogoRepositoryFileId: providerLogoRepositoryFileId,
				Name:                 providerName,
				Email:                providerEmail,
				AddressLine1:         providerAddressLine1,
				AddressLine2:         providerAddressLine2,
				Zip:                  providerAddressZip,
				Locality:             providerAddressLocality,
				Country:              providerAddressCountry,
			},
			Amount:       amount,
			Vat:          vat,
			Total:        total,
			InvoiceLines: invoiceLines,
			UpdatedAt:    utils.ConvertTimeToTimestampPtr(&now),
			AppSource:    constants.AppSourceEventProcessingPlatform,
			Status:       invoicepb.InvoiceStatus_INVOICE_STATUS_DUE,
		})
	})
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error sending the fill invoice request for invoice %s: %s", invoiceId, err.Error())
		return err
	}
	return nil
}

func (h *InvoiceEventHandler) onInvoicePdfGeneratedV1(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceEventHandler.onInvoicePdfGeneratedV1")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData invoice.InvoicePdfGeneratedEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}
	tracing.LogObjectAsJson(span, "eventData", eventData)
	invoiceId := invoice.GetInvoiceObjectID(evt.GetAggregateID(), eventData.Tenant)
	span.SetTag(tracing.SpanTagEntityId, invoiceId)

	invoiceDbNode, err := h.repositories.Neo4jRepositories.InvoiceReadRepository.GetInvoiceById(ctx, eventData.Tenant, invoiceId)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error getting invoice %s: %s", invoiceId, err.Error())
		return err
	}
	if invoiceDbNode == nil {
		err = errors.Errorf("Invoice %s not found", invoiceId)
		tracing.TraceErr(span, err)
		h.log.Errorf("Error getting invoice %s: %s", invoiceId, err.Error())
		return err
	}
	invoiceEntity := neo4jmapper.MapDbNodeToInvoiceEntity(invoiceDbNode)

	if invoiceEntity.DryRun {
		return nil
	}
	// do not invoke invoice finalized webhook if it was already invoked
	if invoiceEntity.InvoiceInternalFields.PaymentRequestedAt == nil {
		err = h.integrationAppInvoiceFinalizedWebhook(ctx, eventData.Tenant, *invoiceEntity)
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("Error invoking invoice ready webhook for invoice %s: %s", invoiceId, err.Error())
		}
	}

	// dispatch invoice finalized event
	err = h.dispatchInvoiceFinalizedEvent(ctx, eventData.Tenant, *invoiceEntity)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error dispatching invoice finalized event for invoice %s: %s", invoiceId, err.Error())
		// TODO: must implement retry mechanism for dispatching invoice finalized event
	}

	return nil
}

func (h *InvoiceEventHandler) integrationAppInvoiceFinalizedWebhook(ctx context.Context, tenant string, invoice neo4jentity.InvoiceEntity) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceEventHandler.integrationAppInvoiceFinalizedWebhook")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	tracing.LogObjectAsJson(span, "invoice", invoice)

	if h.cfg.EventNotifications.EndPoints.InvoiceReady == "" {
		return nil
	}

	// get organization linked to invoice
	organizationDbNode, err := h.repositories.Neo4jRepositories.OrganizationReadRepository.GetOrganizationByInvoiceId(ctx, tenant, invoice.Id)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error getting organization for invoice %s: %s", invoice.Id, err.Error())
		return err
	}
	organizationEntity := neo4jentity.OrganizationEntity{}
	if organizationDbNode != nil {
		organizationEntity = *neo4jmapper.MapDbNodeToOrganizationEntity(organizationDbNode)
	}

	// get stripe customer id for organization
	stripeCustomerIds, err := h.repositories.Neo4jRepositories.ExternalSystemReadRepository.GetAllExternalIdsForLinkedEntity(ctx, tenant, neo4jenum.Stripe.String(), organizationEntity.ID, neo4jutil.NodeLabelOrganization)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error getting stripe customer id for organization %s: %s", organizationEntity.ID, err.Error())
		return err
	}
	identifiedStripeCustomerId := ""
	if len(stripeCustomerIds) == 1 {
		identifiedStripeCustomerId = stripeCustomerIds[0]
	}

	// convert amount to the smallest currency unit
	amountInSmallestCurrencyUnit, err := data.InSmallestCurrencyUnit(invoice.Currency.String(), invoice.TotalAmount)
	if err != nil {
		return fmt.Errorf("error converting amount to smallest currency unit: %v", err.Error())
	}

	requestBody := RequestBodyInvoiceReady{
		Tenant:                       tenant,
		Currency:                     invoice.Currency.String(),
		AmountInSmallestCurrencyUnit: amountInSmallestCurrencyUnit,
		StripeCustomerId:             identifiedStripeCustomerId,
		InvoiceId:                    invoice.Id,
		InvoiceDescription:           fmt.Sprintf("Invoice %s", invoice.Number),
		CustomerOsId:                 organizationEntity.CustomerOsId,
	}

	// Convert the request body to JSON
	requestBodyJSON, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("error encoding JSON: %v", err)
	}

	// Create an HTTP client
	client := &http.Client{}

	// Create a POST request with headers and body
	req, err := http.NewRequest("POST", h.cfg.EventNotifications.EndPoints.InvoiceReady, bytes.NewBuffer(requestBodyJSON))
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

	// Request was successful
	err = h.repositories.Neo4jRepositories.InvoiceWriteRepository.SetInvoicePaymentRequested(ctx, tenant, invoice.Id)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error setting invoice payment requested for invoice %s: %s", invoice.Id, err.Error())
	}

	return nil
}

func (h *InvoiceEventHandler) dispatchInvoiceFinalizedEvent(ctx context.Context, tenant string, invoice neo4jentity.InvoiceEntity) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "InvoiceEventHandler.dispatchInvoiceFinalizedEvent")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	tracing.LogObjectAsJson(span, "invoice", invoice)

	// get organization linked to invoice to build payload for webhook
	organizationDbNode, err := h.repositories.Neo4jRepositories.OrganizationReadRepository.GetOrganizationByInvoiceId(ctx, tenant, invoice.Id)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error getting organization for invoice %s: %s", invoice.Id, err.Error())
		return err
	}
	organizationEntity := neo4jentity.OrganizationEntity{}
	if organizationDbNode != nil {
		organizationEntity = *neo4jmapper.MapDbNodeToOrganizationEntity(organizationDbNode)
	}

	// get contract linked to invoice to build payload for webhook
	contractDbNode, err := h.repositories.Neo4jRepositories.ContractReadRepository.GetContractsForOrganizations(ctx, tenant, []string{organizationEntity.ID})
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error getting contract for invoice %s: %s", invoice.Id, err.Error())
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
	invoiceLineDbNodes, err := h.repositories.Neo4jRepositories.InvoiceLineReadRepository.GetAllForInvoice(ctx, tenant, invoice.Id)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error getting invoice line items for invoice %s: %s", invoice.Id, err.Error())
		return err
	}

	ilEntities := []*neo4jentity.InvoiceLineEntity{}
	for _, ilDbNode := range invoiceLineDbNodes {
		ilEntity := neo4jmapper.MapDbNodeToInvoiceLineEntity(ilDbNode)
		ilEntities = append(ilEntities, ilEntity)
	}

	webhookPayload := webhook.PopulateInvoiceFinalizedPayload(&invoice, &organizationEntity, &contractEntity, ilEntities)

	// dispatch the event
	err = webhook.DispatchWebhook(tenant, webhook.WebhookEventInvoiceFinalized, webhookPayload, h.repositories, h.cfg)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error dispatching invoice finalized event for invoice %s: %s", invoice.Id, err.Error())
		return err
	}

	return nil
}

func (h *InvoiceEventHandler) generateInvoicePDFV1(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceSubscriber.generateInvoicePDFV1")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", evt.GetAggregateID()))

	var eventData invoice.InvoiceFillEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	invoiceId := invoice.GetInvoiceObjectID(evt.GetAggregateID(), eventData.Tenant)

	var invoiceEntity *neo4jentity.InvoiceEntity
	var invoiceLineEntities = []*neo4jentity.InvoiceLineEntity{}

	//load invoice
	invoiceNode, err := h.repositories.Neo4jRepositories.InvoiceReadRepository.GetInvoiceById(ctx, eventData.Tenant, invoiceId)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "InvoiceSubscriber.onInvoiceFillV1.GetInvoice")
	}
	if invoiceNode != nil {
		invoiceEntity = neo4jmapper.MapDbNodeToInvoiceEntity(invoiceNode)
	} else {
		return errors.New("invoiceNode is nil")
	}

	invoiceLinesNodes, err := h.repositories.Neo4jRepositories.InvoiceLineReadRepository.GetAllForInvoice(ctx, eventData.Tenant, invoiceId)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "InvoiceSubscriber.onInvoiceFillV1.GetAllForInvoice")
	}
	if invoiceLinesNodes != nil {
		for _, invoiceLineNode := range invoiceLinesNodes {
			invoiceLineEntities = append(invoiceLineEntities, neo4jmapper.MapDbNodeToInvoiceLineEntity(invoiceLineNode))
		}
	} else {
		return errors.New("invoiceLinesNodes is nil")
	}

	invoiceHasVat := false

	if invoiceEntity.Vat > 0 {
		invoiceHasVat = true
	}

	data := map[string]interface{}{
		"Tenant":                        eventData.Tenant,
		"CustomerName":                  invoiceEntity.Customer.Name,
		"CustomerEmail":                 invoiceEntity.Customer.Email,
		"CustomerAddressLine1":          invoiceEntity.Customer.AddressLine1,
		"CustomerAddressLine2":          invoiceEntity.Customer.AddressLine2,
		"CustomerAddressLine3":          utils.JoinNonEmpty(", ", invoiceEntity.Customer.Locality, invoiceEntity.Customer.Zip),
		"CustomerCountry":               invoiceEntity.Customer.Country,
		"ProviderLogoExtension":         "",
		"ProviderLogoRepositoryFileId":  invoiceEntity.Provider.LogoRepositoryFileId,
		"ProviderName":                  invoiceEntity.Provider.Name,
		"ProviderEmail":                 invoiceEntity.Provider.Email,
		"ProviderAddressLine1":          invoiceEntity.Provider.AddressLine1,
		"ProviderAddressLine2":          invoiceEntity.Provider.AddressLine2,
		"ProviderAddressLine3":          utils.JoinNonEmpty(", ", invoiceEntity.Provider.Locality, invoiceEntity.Provider.Zip),
		"ProviderCountry":               invoiceEntity.Provider.Country,
		"InvoiceNumber":                 invoiceEntity.Number,
		"InvoiceIssueDate":              invoiceEntity.CreatedAt.Format("02 Jan 2006"),
		"InvoiceDueDate":                invoiceEntity.DueDate.Format("02 Jan 2006"),
		"InvoiceCurrency":               invoiceEntity.Currency.String() + "" + invoiceEntity.Currency.Symbol(),
		"InvoiceSubtotal":               utils.FormatAmount(invoiceEntity.Amount, 2),
		"InvoiceTotal":                  utils.FormatAmount(invoiceEntity.TotalAmount, 2),
		"InvoiceAmountDue":              utils.FormatAmount(invoiceEntity.TotalAmount, 2),
		"InvoiceLineItems":              []map[string]string{},
		"Note":                          invoiceEntity.Note,
		"DomesticPaymentsBankInfo":      invoiceEntity.DomesticPaymentsBankInfo,
		"InternationalPaymentsBankInfo": invoiceEntity.InternationalPaymentsBankInfo,
	}

	if invoiceHasVat {
		data["InvoiceVat"] = fmt.Sprintf("%.2f", invoiceEntity.Vat)
	}

	for _, line := range invoiceLineEntities {
		invoiceLineItem := map[string]string{
			"Name":      line.Name,
			"Quantity":  fmt.Sprintf("%d", line.Quantity),
			"UnitPrice": invoiceEntity.Currency.Symbol() + utils.FormatAmount(line.Price, 2),
			"Amount":    invoiceEntity.Currency.Symbol() + utils.FormatAmount(line.Amount, 2),
			"Vat":       invoiceEntity.Currency.Symbol() + utils.FormatAmount(line.Vat, 2),
		}
		if line.BilledType == neo4jenum.BilledTypeOnce {
			sliDbNode, err := h.repositories.Neo4jRepositories.ServiceLineItemReadRepository.GetServiceLineItemById(ctx, eventData.Tenant, line.ServiceLineItemId)
			if err == nil && sliDbNode != nil {
				sliEntity := neo4jmapper.MapDbNodeToServiceLineItemEntity(sliDbNode)
				invoiceLineItem["InvoiceLineSubtitle"] = sliEntity.StartedAt.Format("02 Jan 2006")
			}
		}
		if _, ok := invoiceLineItem["InvoiceLineSubtitle"]; !ok {
			invoiceLineItem["InvoiceLineSubtitle"] = fmt.Sprintf("%s - %s", invoiceEntity.PeriodStartDate.Format("02 Jan 2006"), invoiceEntity.PeriodEndDate.Format("02 Jan 2006"))
		}

		if invoiceHasVat {
			invoiceLineItem["InvoiceHasVat"] = "true"
		}

		data["InvoiceLineItems"] = append(data["InvoiceLineItems"].([]map[string]string), invoiceLineItem)
	}

	//prepare the temp html file
	tmpInvoiceFile, err := os.CreateTemp("", "invoice_*.html")
	if err != nil {
		return errors.Wrap(err, "os.TempFile")
	}
	defer os.Remove(tmpInvoiceFile.Name()) // Delete the temporary HTML file when done
	defer tmpInvoiceFile.Close()

	if invoiceEntity.Provider.LogoRepositoryFileId != "" {
		fileMetadata, err := h.fsc.GetFileMetadata(eventData.Tenant, invoiceEntity.Provider.LogoRepositoryFileId, span)
		if err != nil {
			tracing.TraceErr(span, err)
			return errors.Wrap(err, "InvoiceSubscriber.onInvoiceFillV1.GetFileMetadata")
		}

		data["ProviderLogoExtension"] = GetFileExtensionFromMetadata(fileMetadata)
	}

	//fill the template with data and store it in temp
	err = FillInvoiceHtmlTemplate(tmpInvoiceFile, data)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "FillInvoiceHtmlTemplate")
	}

	//convert the temp to pdf
	pdfBytes, err := ConvertInvoiceHtmlToPdf(h.fsc, h.cfg.Subscriptions.InvoiceSubscription.PdfConverterUrl, tmpInvoiceFile, data, span)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "ConvertInvoiceHtmlToPdf")
	}

	if pdfBytes == nil {
		return errors.New("pdfBytes is nil")
	}

	//TODO remove this at some point when we are sure that the pdf is generated correctly
	// Save the PDF file to disk
	os.WriteFile("output.pdf", *pdfBytes, 0644)

	basePath := fmt.Sprintf("/INVOICE/%d/%s", invoiceEntity.CreatedAt.Year(), invoiceEntity.CreatedAt.Format("01"))

	if invoiceEntity.DryRun {
		basePath = basePath + "/DRY_RUN"
	}

	fileDTO, err := h.fsc.UploadSingleFileBytes(eventData.Tenant, basePath, invoiceEntity.Id, "Invoice - "+invoiceEntity.Number+".pdf", *pdfBytes, span)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "InvoiceSubscriber.onInvoiceFillV1.UploadSingleFileBytes")
	}

	if fileDTO.Id == "" {
		return errors.New("fileDTO.Id is empty")
	}

	err = h.callPdfGeneratedInvoice(ctx, eventData.Tenant, invoiceId, fileDTO.Id, span)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "InvoiceSubscriber.onInvoiceFillV1.CallPdfGeneratedInvoice")
	}

	return nil
}

func (s *InvoiceEventHandler) callPdfGeneratedInvoice(ctx context.Context, tenant, invoiceId, repositoryFileId string, span opentracing.Span) error {
	ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
	_, err := subscriptions.CallEventsPlatformGRPCWithRetry[*invoicepb.InvoiceIdResponse](func() (*invoicepb.InvoiceIdResponse, error) {
		return s.grpcClients.InvoiceClient.PdfGeneratedInvoice(ctx, &invoicepb.PdfGeneratedInvoiceRequest{
			Tenant:           tenant,
			InvoiceId:        invoiceId,
			RepositoryFileId: repositoryFileId,
			AppSource:        constants.AppSourceEventProcessingPlatform,
		})
	})
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error sending the pdf generated request for invoice %s: %s", invoiceId, err.Error())
		return err
	}
	return nil
}

func (h *InvoiceEventHandler) onInvoiceVoidV1(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceEventHandler.onInvoiceVoidV1")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData invoice.InvoiceVoidEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}
	tracing.LogObjectAsJson(span, "eventData", eventData)
	invoiceId := invoice.GetInvoiceObjectID(evt.GetAggregateID(), eventData.Tenant)
	span.SetTag(tracing.SpanTagEntityId, invoiceId)
	span.SetTag(tracing.SpanTagTenant, eventData.Tenant)

	var invoiceEntity *neo4jentity.InvoiceEntity

	//load invoice
	invoiceNode, err := h.repositories.Neo4jRepositories.InvoiceReadRepository.GetInvoiceById(ctx, eventData.Tenant, invoiceId)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "InvoiceSubscriber.onInvoiceFillV1.GetInvoice")
	}
	if invoiceNode != nil {
		invoiceEntity = neo4jmapper.MapDbNodeToInvoiceEntity(invoiceNode)
	} else {
		return errors.New("invoiceNode is nil")
	}

	if invoiceEntity.DryRun {
		return nil
	}

	//load tenant billing profile from neo4j
	tenantBillingProfileEntity, err := h.loadTenantBillingProfile(ctx, eventData.Tenant, false)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	postmarkEmail := notifications.PostmarkEmail{
		WorkflowId:    notifications.WorkflowInvoiceVoided,
		MessageStream: notifications.PostmarkMessageStreamInvoice,
		From:          invoiceEntity.Provider.Email,
		To:            invoiceEntity.Customer.Email,
		BCC:           tenantBillingProfileEntity.SendInvoicesBcc,
		Subject:       "Voided invoice " + invoiceEntity.Number,
		TemplateData: map[string]string{
			"{{userFirstName}}":  invoiceEntity.Customer.Name,
			"{{invoiceNumber}}":  invoiceEntity.Number,
			"{{currencySymbol}}": invoiceEntity.Currency.Symbol(),
			"{{amtDue}}":         fmt.Sprintf("%.2f", invoiceEntity.TotalAmount),
			"{{issueDate}}":      invoiceEntity.CreatedAt.Format("02 Jan 2006"),
		},
		Attachments: []notifications.PostmarkEmailAttachment{},
	}

	err = h.AppendProviderLogoToEmail(eventData.Tenant, invoiceEntity.Provider.LogoRepositoryFileId, &postmarkEmail, span)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error appending provider logo to email for invoice %s: %s", invoiceId, err.Error())
		return err
	}

	err = h.postmarkProvider.SendNotification(ctx, postmarkEmail, span, eventData.Tenant)

	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error sending invoice voided notification for invoice %s: %s", invoiceId, err.Error())
		return err
	}

	// Request was successful
	err = h.repositories.Neo4jRepositories.InvoiceWriteRepository.SetPaidInvoiceNotificationSentAt(ctx, eventData.Tenant, invoiceId)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error setting invoice voided notification sent at for invoice %s: %s", invoiceId, err.Error())
		return err
	}

	return nil
}

func (h *InvoiceEventHandler) onInvoicePaidV1(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceEventHandler.onInvoicePaidV1")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData invoice.InvoicePaidEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}
	tracing.LogObjectAsJson(span, "eventData", eventData)
	invoiceId := invoice.GetInvoiceObjectID(evt.GetAggregateID(), eventData.Tenant)
	span.SetTag(tracing.SpanTagEntityId, invoiceId)
	span.SetTag(tracing.SpanTagTenant, eventData.Tenant)

	var invoiceEntity neo4jentity.InvoiceEntity
	var contractEntity neo4jentity.ContractEntity

	//load invoice
	invoiceNode, err := h.repositories.Neo4jRepositories.InvoiceReadRepository.GetInvoiceById(ctx, eventData.Tenant, invoiceId)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "InvoiceSubscriber.onInvoicePaidV1.GetInvoice")
	}
	if invoiceNode != nil {
		invoiceEntity = *neo4jmapper.MapDbNodeToInvoiceEntity(invoiceNode)
	} else {
		return errors.New("invoiceNode is nil")
	}

	// load contract
	contractNode, err := h.repositories.Neo4jRepositories.ContractReadRepository.GetContractForInvoice(ctx, eventData.Tenant, invoiceEntity.Id)
	if err != nil {
		tracing.TraceErr(span, err)
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
		return errors.New("contractEntity.InvoiceEmail is empty or invalid")
	}

	//load tenant billing profile from neo4j
	tenantBillingProfileEntity, err := h.loadTenantBillingProfile(ctx, eventData.Tenant, false)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	postmarkEmail := notifications.PostmarkEmail{
		WorkflowId:    notifications.WorkflowInvoicePaid,
		MessageStream: notifications.PostmarkMessageStreamInvoice,
		From:          invoiceEntity.Provider.Email,
		To:            contractEntity.InvoiceEmail,
		BCC:           tenantBillingProfileEntity.SendInvoicesBcc,
		Subject:       "Paid Invoice " + invoiceEntity.Number + " from " + invoiceEntity.Provider.Name,
		TemplateData: map[string]string{
			"{{userFirstName}}":  invoiceEntity.Customer.Name,
			"{{invoiceNumber}}":  invoiceEntity.Number,
			"{{currencySymbol}}": invoiceEntity.Currency.Symbol(),
			"{{amtDue}}":         fmt.Sprintf("%.2f", invoiceEntity.TotalAmount),
			"{{paymentDate}}":    invoiceEntity.DueDate.Format("02 Jan 2006"),
		},
		Attachments: []notifications.PostmarkEmailAttachment{},
	}

	err = h.AppendInvoiceFileToEmailAsAttachment(eventData.Tenant, invoiceEntity, &postmarkEmail, span)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error appending invoice file to email attachment for invoice %s: %s", invoiceId, err.Error())
		return err
	}

	err = h.AppendProviderLogoToEmail(eventData.Tenant, invoiceEntity.Provider.LogoRepositoryFileId, &postmarkEmail, span)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error appending provider logo to email for invoice %s: %s", invoiceId, err.Error())
		return err
	}

	err = h.postmarkProvider.SendNotification(ctx, postmarkEmail, span, eventData.Tenant)

	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error sending invoice paid notification for invoice %s: %s", invoiceId, err.Error())
		return err
	}

	// Request was successful
	err = h.repositories.Neo4jRepositories.InvoiceWriteRepository.SetPaidInvoiceNotificationSentAt(ctx, eventData.Tenant, invoiceId)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error setting invoice paid notification sent at for invoice %s: %s", invoiceId, err.Error())
		return err
	}

	return nil
}

func (h *InvoiceEventHandler) onInvoicePayNotificationV1(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceEventHandler.onInvoicePayNotificationV1")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData invoice.InvoicePayNotificationEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}
	tracing.LogObjectAsJson(span, "eventData", eventData)
	invoiceId := invoice.GetInvoiceObjectID(evt.GetAggregateID(), eventData.Tenant)
	span.SetTag(tracing.SpanTagEntityId, invoiceId)
	span.SetTag(tracing.SpanTagTenant, eventData.Tenant)

	var invoiceEntity neo4jentity.InvoiceEntity
	var contractEntity neo4jentity.ContractEntity

	//load invoice
	invoiceNode, err := h.repositories.Neo4jRepositories.InvoiceReadRepository.GetInvoiceById(ctx, eventData.Tenant, invoiceId)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "InvoiceSubscriber.onInvoicePayNotificationV1.GetInvoice")
	}
	if invoiceNode != nil {
		invoiceEntity = *neo4jmapper.MapDbNodeToInvoiceEntity(invoiceNode)
	} else {
		tracing.TraceErr(span, errors.New("invoiceNode is nil"))
		return errors.New("invoiceNode is nil")
	}

	contractNode, err := h.repositories.Neo4jRepositories.ContractReadRepository.GetContractForInvoice(ctx, eventData.Tenant, invoiceEntity.Id)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "InvoiceSubscriber.onInvoicePayNotificationV1.GetContractForInvoice")
	}
	if contractNode != nil {
		contractEntity = *neo4jmapper.MapDbNodeToContractEntity(contractNode)
	} else {
		tracing.TraceErr(span, errors.New("contractNode is nil"))
		return errors.New("contractNode is nil")
	}

	if contractEntity.InvoiceEmail == "" || !isValidEmailSyntax(contractEntity.InvoiceEmail) {
		tracing.TraceErr(span, errors.New("contractEntity.InvoiceEmail is empty or invalid"))
		return errors.New("contractEntity.InvoiceEmail is empty or invalid")
	}

	if invoiceEntity.PaymentDetails.PaymentLink == "" {
		tracing.TraceErr(span, errors.New("invoiceEntity.PaymentDetails.PaymentLink is empty"))
		return errors.New("invoiceEntity.PaymentDetails.PaymentLink is empty")
	}

	//load tenant billing profile from neo4j
	tenantBillingProfileEntity, err := h.loadTenantBillingProfile(ctx, eventData.Tenant, false)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	postmarkEmail := notifications.PostmarkEmail{
		WorkflowId:    notifications.WorkflowInvoiceReady,
		MessageStream: notifications.PostmarkMessageStreamInvoice,
		From:          invoiceEntity.Provider.Email,
		To:            contractEntity.InvoiceEmail,
		BCC:           tenantBillingProfileEntity.SendInvoicesBcc,
		Subject:       "New invoice " + invoiceEntity.Number,
		TemplateData: map[string]string{
			"{{organizationName}}": invoiceEntity.Customer.Name,
			"{{invoiceNumber}}":    invoiceEntity.Number,
			"{{currencySymbol}}":   invoiceEntity.Currency.Symbol(),
			"{{amtDue}}":           fmt.Sprintf("%.2f", invoiceEntity.TotalAmount),
			"{{paymentLink}}":      invoiceEntity.PaymentDetails.PaymentLink,
		},
		Attachments: []notifications.PostmarkEmailAttachment{},
	}

	err = h.AppendInvoiceFileToEmailAsAttachment(eventData.Tenant, invoiceEntity, &postmarkEmail, span)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error appending invoice file to email attachment for invoice %s: %s", invoiceId, err.Error())
		return err
	}

	err = h.AppendProviderLogoToEmail(eventData.Tenant, invoiceEntity.Provider.LogoRepositoryFileId, &postmarkEmail, span)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error appending provider logo to email for invoice %s: %s", invoiceId, err.Error())
		return err
	}

	err = h.postmarkProvider.SendNotification(ctx, postmarkEmail, span, eventData.Tenant)

	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error sending invoice pay request notification for invoice %s: %s", invoiceId, err.Error())
		return err
	}

	// Request was successful
	err = h.repositories.Neo4jRepositories.InvoiceWriteRepository.SetPayInvoiceNotificationSentAt(ctx, eventData.Tenant, invoiceId)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error setting invoice pay notification sent at for invoice %s: %s", invoiceId, err.Error())
		return err
	}

	return nil
}

func (h *InvoiceEventHandler) AppendInvoiceFileToEmailAsAttachment(tenant string, invoice neo4jentity.InvoiceEntity, postmarkEmail *notifications.PostmarkEmail, span opentracing.Span) error {
	invoiceFileBytes, err := h.fsc.GetFileBytes(tenant, invoice.RepositoryFileId, span)
	if err != nil {
		return err
	}

	postmarkEmail.Attachments = append(postmarkEmail.Attachments, notifications.PostmarkEmailAttachment{
		Filename:       "Invoice " + invoice.Number + ".pdf",
		ContentEncoded: base64.StdEncoding.EncodeToString(*invoiceFileBytes),
		ContentType:    "application/pdf",
	})

	return nil
}

func (h *InvoiceEventHandler) AppendProviderLogoToEmail(tenant, logoFileId string, postmarkEmail *notifications.PostmarkEmail, span opentracing.Span) error {
	if logoFileId == "" {
		return nil
	}

	metadata, fileBytes, err := h.fsc.GetFile(tenant, logoFileId, span)
	if err != nil {
		return err
	}

	postmarkEmail.Attachments = append(postmarkEmail.Attachments, notifications.PostmarkEmailAttachment{
		Filename:       "provider-logo-file-encoded",
		ContentEncoded: base64.StdEncoding.EncodeToString(*fileBytes),
		ContentType:    metadata.MimeType,
		ContentID:      "cid:provider-logo-file-encoded",
	})

	return nil
}

func isValidEmailSyntax(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func (h *InvoiceEventHandler) loadTenantBillingProfile(ctx context.Context, tenant string, failIfNotFound bool) (neo4jentity.TenantBillingProfileEntity, error) {
	tenantBillingProfiles, err := h.repositories.Neo4jRepositories.TenantReadRepository.GetTenantBillingProfiles(ctx, tenant)
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
