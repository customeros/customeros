package invoice

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/data"
	fsc "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/file_store_client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/neo4jutil"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/invoice"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/grpc_client"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	commonpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/common"
	invoicepb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/invoice"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"net/http"
	"os"
	"time"
)

type RequestBodyInvoiceReady struct {
	Tenant                       string `json:"tenant"`
	Currency                     string `json:"currency"`
	AmountInSmallestCurrencyUnit int64  `json:"amountInSmallestCurrencyUnit"`
	StripeCustomerId             string `json:"stripeCustomerId"`
	InvoiceId                    string `json:"invoiceId"`
	InvoiceDescription           string `json:"invoiceDescription"`
}

type InvoiceEventHandler struct {
	log          logger.Logger
	repositories *repository.Repositories
	cfg          config.Config
	grpcClients  *grpc_client.Clients
	fsc          fsc.FileStoreApiService
}

func NewInvoiceEventHandler(log logger.Logger, repositories *repository.Repositories, cfg config.Config, grpcClients *grpc_client.Clients, fsc fsc.FileStoreApiService) *InvoiceEventHandler {
	return &InvoiceEventHandler{
		log:          log,
		repositories: repositories,
		cfg:          cfg,
		grpcClients:  grpcClients,
		fsc:          fsc,
	}
}

func (h *InvoiceEventHandler) onInvoiceCreateForContractV1(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceEventHandler.onInvoiceCreateForContractV1")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData invoice.InvoiceForContractCreateEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}
	tracing.LogObjectAsJson(span, "eventData", eventData)
	invoiceId := invoice.GetInvoiceObjectID(evt.GetAggregateID(), eventData.Tenant)
	span.SetTag(tracing.SpanTagEntityId, invoiceId)

	sliDbNodes, err := h.repositories.Neo4jRepositories.ServiceLineItemReadRepository.GetAllForContract(ctx, eventData.Tenant, eventData.ContractId)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error getting service line items for contract %s: %s", eventData.ContractId, err.Error())
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
	invoiceLines := []*invoicepb.InvoiceLine{}

	referenceTime := eventData.PeriodStartDate
	for _, sliEntity := range sliEntities {
		// skip for now one time and usage SLIs
		if sliEntity.Billed == neo4jenum.BilledTypeOnce || sliEntity.Billed == neo4jenum.BilledTypeUsage {
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
		// process monthly, quarterly and annually SLIs
		if sliEntity.Billed == neo4jenum.BilledTypeMonthly || sliEntity.Billed == neo4jenum.BilledTypeQuarterly || sliEntity.Billed == neo4jenum.BilledTypeAnnually {
			calculatedSLIAmount := calculateSLIAmountForCycleInvoicing(sliEntity.Quantity, sliEntity.Price, sliEntity.Billed, neo4jenum.DecodeBillingCycle(eventData.BillingCycle))
			calculatedSLIAmount = utils.TruncateFloat64(calculatedSLIAmount, 2)
			amount += calculatedSLIAmount
			vat += float64(0)
			invoiceLine := invoicepb.InvoiceLine{
				Name:                    sliEntity.Name,
				Price:                   sliEntity.Price,
				Quantity:                sliEntity.Quantity,
				Amount:                  calculatedSLIAmount,
				Total:                   calculatedSLIAmount,
				Vat:                     float64(0),
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

	var contractEntity *neo4jentity.ContractEntity
	var tenantSettingsEntity *neo4jentity.TenantSettingsEntity
	var tenantBillingProfileEntity *neo4jentity.TenantBillingProfileEntity

	//load contract from neo4j
	contract, err := h.repositories.Neo4jRepositories.ContractReadRepository.GetContractById(ctx, eventData.Tenant, eventData.ContractId)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "InvoiceSubscriber.onInvoiceFillV1.GetContractById")
	}
	if contract != nil {
		contractEntity = neo4jmapper.MapDbNodeToContractEntity(contract)
	} else {
		return errors.New("contract is nil")
	}

	//load tenant settings from neo4j
	tenantSettings, err := h.repositories.Neo4jRepositories.TenantReadRepository.GetTenantSettings(ctx, eventData.Tenant)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "InvoiceSubscriber.onInvoiceFillV1.GetTenantSettings")
	}
	if tenantSettings != nil {
		tenantSettingsEntity = neo4jmapper.MapDbNodeToTenantSettingsEntity(tenantSettings)
	} else {
		return errors.New("tenantSettings is nil")
	}

	//load tenant billing profile from neo4j
	tenantBillingProfiles, err := h.repositories.Neo4jRepositories.TenantReadRepository.GetTenantBillingProfiles(ctx, eventData.Tenant)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "InvoiceSubscriber.onInvoiceFillV1.GetTenantSettings")
	}
	if tenantBillingProfiles == nil || len(tenantBillingProfiles) > 0 {
		tenantBillingProfileEntity = neo4jmapper.MapDbNodeToTenantBillingProfileEntity(tenantBillingProfiles[0])
	} else {
		return errors.New("tenantBillingProfiles is nil or empty")
	}

	err = h.callFillInvoice(ctx,
		eventData.Tenant,
		invoiceId,
		tenantBillingProfileEntity.DomesticPaymentsBankInfo,
		tenantBillingProfileEntity.InternationalPaymentsBankInfo,
		contractEntity.OrganizationLegalName,
		contractEntity.InvoiceEmail,
		contractEntity.AddressLine1, contractEntity.AddressLine2, contractEntity.Zip, contractEntity.Locality, contractEntity.Country,
		tenantSettingsEntity.LogoUrl,
		tenantBillingProfileEntity.LegalName,
		tenantBillingProfileEntity.AddressLine1, tenantBillingProfileEntity.AddressLine2, tenantBillingProfileEntity.Zip, tenantBillingProfileEntity.Locality, tenantBillingProfileEntity.Country,
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

func calculateSLIAmountForCycleInvoicing(quantity int64, price float64, billed neo4jenum.BilledType, cycle neo4jenum.BillingCycle) float64 {
	sliAmount := float64(quantity) * price
	if sliAmount == 0 {
		return sliAmount
	}
	switch cycle {
	case neo4jenum.BillingCycleMonthlyBilling:
		switch billed {
		case neo4jenum.BilledTypeMonthly:
			return sliAmount
		case neo4jenum.BilledTypeQuarterly:
			return sliAmount / 3
		case neo4jenum.BilledTypeAnnually:
			return sliAmount / 12
		}
	case neo4jenum.BillingCycleQuarterlyBilling:
		switch billed {
		case neo4jenum.BilledTypeMonthly:
			return sliAmount * 3
		case neo4jenum.BilledTypeQuarterly:
			return sliAmount
		case neo4jenum.BilledTypeAnnually:
			return sliAmount / 4
		}
	case neo4jenum.BillingCycleAnnuallyBilling:
		switch billed {
		case neo4jenum.BilledTypeMonthly:
			return sliAmount * 12
		case neo4jenum.BilledTypeQuarterly:
			return sliAmount * 4
		case neo4jenum.BilledTypeAnnually:
			return sliAmount
		}
	}
	return float64(0)
}

func (h *InvoiceEventHandler) callFillInvoice(ctx context.Context, tenant, invoiceId, domesticPaymentsBankInfo, internationalPaymentsBankInfo,
	customerName, customerEmail, customerAddressLine1, customerAddressLine2, customerAddressZip, customerAddressLocality, customerAddressCountry,
	providerLogoUrl, providerName, providerAddressLine1, providerAddressLine2, providerAddressZip, providerAddressLocality, providerAddressCountry,
	note string, amount, vat, total float64, invoiceLines []*invoicepb.InvoiceLine, span opentracing.Span) error {
	ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
	now := time.Now()
	_, err := h.grpcClients.InvoiceClient.FillInvoice(ctx, &invoicepb.FillInvoiceRequest{
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
			LogoUrl:      providerLogoUrl,
			Name:         providerName,
			AddressLine1: providerAddressLine1,
			AddressLine2: providerAddressLine2,
			Zip:          providerAddressZip,
			Locality:     providerAddressLocality,
			Country:      providerAddressCountry,
		},
		Amount:       amount,
		Vat:          vat,
		Total:        total,
		InvoiceLines: invoiceLines,
		UpdatedAt:    utils.ConvertTimeToTimestampPtr(&now),
		AppSource:    constants.AppSourceEventProcessingPlatform,
		Status:       invoicepb.InvoiceStatus_INVOICE_STATUS_DUE,
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

	// do not invoke invoice ready webhook if it was already invoked
	if invoiceEntity.InvoiceInternalFields.PaymentRequestedAt == nil {
		err = h.invokeInvoiceReadyWebhook(ctx, eventData.Tenant, *invoiceEntity)
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("Error invoking invoice ready webhook for invoice %s: %s", invoiceId, err.Error())
			return err
		}
	}

	return nil
}

func (h *InvoiceEventHandler) invokeInvoiceReadyWebhook(ctx context.Context, tenant string, invoice neo4jentity.InvoiceEntity) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceEventHandler.invokeInvoiceReadyWebhook")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	tracing.LogObjectAsJson(span, "invoice", invoice)

	if invoice.DryRun {
		return nil
	}
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
	stripeCustomerId, err := h.repositories.Neo4jRepositories.ExternalSystemReadRepository.GetFirstExternalIdForLinkedEntity(ctx, tenant, neo4jenum.Stripe.String(), organizationEntity.ID, neo4jutil.NodeLabelOrganization)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error getting stripe customer id for organization %s: %s", organizationEntity.ID, err.Error())
		return err
	}

	// convert amount to the smallest currency unit
	amountInSmallestCurrencyUnit, err := data.InSmallestCurrencyUnit(invoice.Currency.String(), invoice.Amount)
	if err != nil {
		return fmt.Errorf("error converting amount to smallest currency unit: %v", err.Error())
	}

	requestBody := RequestBodyInvoiceReady{
		Tenant:                       tenant,
		Currency:                     invoice.Currency.String(),
		AmountInSmallestCurrencyUnit: amountInSmallestCurrencyUnit,
		StripeCustomerId:             stripeCustomerId,
		InvoiceId:                    invoice.Id,
		InvoiceDescription:           fmt.Sprintf("Invoice %s", invoice.Number),
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
	var invoiceLineEntities []*neo4jentity.InvoiceLineEntity = []*neo4jentity.InvoiceLineEntity{}

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

	data := map[string]interface{}{
		"CustomerName":                  invoiceEntity.Customer.Name,
		"CustomerEmail":                 invoiceEntity.Customer.Email,
		"CustomerAddressLine1":          invoiceEntity.Customer.AddressLine1,
		"CustomerAddressLine2":          invoiceEntity.Customer.AddressLine2,
		"CustomerAddressLine3":          utils.JoinNonEmpty(", ", invoiceEntity.Customer.Zip, invoiceEntity.Customer.Locality, invoiceEntity.Customer.Country),
		"ProviderLogoUrl":               invoiceEntity.Provider.LogoUrl,
		"ProviderLogoExtension":         GetFileExtensionFromUrl(invoiceEntity.Provider.LogoUrl),
		"ProviderName":                  invoiceEntity.Provider.Name,
		"ProviderAddressLine1":          invoiceEntity.Provider.AddressLine1,
		"ProviderAddressLine2":          invoiceEntity.Provider.AddressLine2,
		"ProviderAddressLine3":          utils.JoinNonEmpty(", ", invoiceEntity.Provider.Zip, invoiceEntity.Provider.Locality, invoiceEntity.Provider.Country),
		"InvoiceNumber":                 invoiceEntity.Number,
		"InvoiceIssueDate":              invoiceEntity.CreatedAt.Format("02 Jan 2006"),
		"InvoiceDueDate":                invoiceEntity.DueDate.Format("02 Jan 2006"),
		"InvoiceCurrency":               invoiceEntity.Currency.String() + "" + invoiceEntity.Currency.Symbol(),
		"InvoiceSubtotal":               fmt.Sprintf("%.2f", invoiceEntity.TotalAmount),
		"InvoiceTax":                    "0.00",
		"InvoiceTotal":                  fmt.Sprintf("%.2f", invoiceEntity.TotalAmount),
		"InvoiceAmountDue":              fmt.Sprintf("%.2f", invoiceEntity.TotalAmount),
		"InvoiceLineItems":              []map[string]string{},
		"Note":                          invoiceEntity.Note,
		"DomesticPaymentsBankInfo":      invoiceEntity.DomesticPaymentsBankInfo,
		"InternationalPaymentsBankInfo": invoiceEntity.InternationalPaymentsBankInfo,
	}

	for _, line := range invoiceLineEntities {
		data["InvoiceLineItems"] = append(data["InvoiceLineItems"].([]map[string]string), map[string]string{
			"Name":        line.Name,
			"Description": "",
			"Quantity":    fmt.Sprintf("%d", line.Quantity),
			"UnitPrice":   fmt.Sprintf("%.2f", line.Price),
			"Amount":      fmt.Sprintf("%.2f", line.Amount),
		})
	}

	//prepare the temp html file
	tmpInvoiceFile, err := os.CreateTemp("", "invoice_*.html")
	if err != nil {
		return errors.Wrap(err, "ioutil.TempFile")
	}
	defer os.Remove(tmpInvoiceFile.Name()) // Delete the temporary HTML file when done
	defer tmpInvoiceFile.Close()

	//fill the template with data and store it in temp
	err = FillInvoiceHtmlTemplate(tmpInvoiceFile, data)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "FillInvoiceHtmlTemplate")
	}

	//convert the temp to pdf
	pdfBytes, err := ConvertInvoiceHtmlToPdf(h.cfg.Subscriptions.InvoiceSubscription.PdfConverterUrl, tmpInvoiceFile, data)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "ConvertInvoiceHtmlToPdf")
	}

	// Save the PDF file to disk
	//err = ioutil.WriteFile("output.pdf", *pdfBytes, 0644)
	//if err != nil {
	//	return errors.Wrap(err, "ioutil.WriteFile")
	//}

	fileDTO, err := h.fsc.UploadSingleFileBytes(eventData.Tenant, *pdfBytes)
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
	_, err := s.grpcClients.InvoiceClient.PdfGeneratedInvoice(ctx, &invoicepb.PdfGeneratedInvoiceRequest{
		Tenant:           tenant,
		InvoiceId:        invoiceId,
		RepositoryFileId: repositoryFileId,
		AppSource:        constants.AppSourceEventProcessingPlatform,
	})
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error sending the pdf generated request for invoice %s: %s", invoiceId, err.Error())
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

	// TODO Implement

	return nil
}
