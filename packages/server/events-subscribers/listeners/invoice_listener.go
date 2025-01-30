package listeners

import (
	"fmt"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/dto"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	model2 "github.com/customeros/customeros/packages/server/customer-os-common-module/model"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	neo4j_entity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"
	"github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/mapper"
	"github.com/customeros/customeros/packages/server/events-subscribers/model"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

// check if organization exists in quickbooks, if not create it
// check if all line items have skuId and the skuId has been pushed to quickbooks
// save invoice to quickbooks
func OnInvoiceFinalized(ctx context.Context, dependencies *model.DependencyContainer, input any) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Listeners.OnInvoiceFinalized")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)

	message := input.(*dto.Event)
	invoiceId := message.Event.EntityId
	span.SetTag(tracing.SpanTagEntityId, invoiceId)

	tenant := common.GetTenantFromContext(ctx)
	if tenant == "" {
		err := errors.New("Missing tenant in context")
		tracing.TraceErr(span, err)
		return err
	}

	quickbooksSettingsEntity, err := dependencies.CommonServices.PostgresRepositories.QuickbooksSettingsRepository.Get(ctx, tenant)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if quickbooksSettingsEntity == nil {
		span.LogFields(log.String("error", "Quickbooks settings not found"))
		return nil
	}

	invoice, err := dependencies.CommonServices.InvoiceService.GetById(ctx, nil, invoiceId)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if invoice == nil {
		err := errors.New("Invoice not found")
		tracing.TraceErr(span, err)
		return err
	}

	invoiceLines, err := dependencies.CommonServices.InvoiceService.GetInvoiceLinesForInvoices(ctx, []string{invoiceId})
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if invoiceLines == nil || len(*invoiceLines) == 0 {
		span.LogFields(log.String("skip", "Invoice lines not found"))
		return nil
	}

	//validate all invoice line have skuId and the skuId has been pushed to quickbooks
	for _, invoiceLine := range *invoiceLines {
		if invoiceLine.SkuId == "" {
			span.LogFields(log.String("skip", fmt.Sprintf("SkuId not found for invoice line %s", invoiceLine.Id)))
			return nil
		}

		sku, err := dependencies.CommonServices.PostgresRepositories.SkuRepository.Get(ctx, tenant, invoiceLine.SkuId)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}

		if sku == nil {
			err := errors.New(fmt.Sprintf("Sku not found for invoice line %s", invoiceLine.Id))
			tracing.TraceErr(span, err)
			return err
		}

		if sku.QuickbooksId == "" {
			span.LogFields(log.String("skip", fmt.Sprintf("QuickbooksId not found for sku %s", sku.ID)))
			return nil
		}
	}

	//save invoice to quickbooks
	quickbooksInvoiceLines := make([]interfaces.QuickbooksInvoiceLine, 0)
	for _, invoiceLine := range *invoiceLines {
		sku, err := dependencies.CommonServices.PostgresRepositories.SkuRepository.Get(ctx, tenant, invoiceLine.SkuId)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}

		quickbooksInvoiceLine := interfaces.QuickbooksInvoiceLine{
			DetailType: "SalesItemLineDetail",
			Amount:     invoiceLine.Amount,
		}
		quickbooksInvoiceLine.SalesItemLineDetail.ItemRef.Value = sku.QuickbooksId

		quickbooksInvoiceLines = append(quickbooksInvoiceLines, quickbooksInvoiceLine)
	}

	organizationNode, err := dependencies.CommonServices.Neo4jRepositories.OrganizationReadRepository.GetOrganizationByInvoiceId(ctx, tenant, invoiceId)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if organizationNode == nil {
		err := errors.New("Organization not found for invoice")
		tracing.TraceErr(span, err)
		return err
	}

	organization := mapper.MapDbNodeToOrganizationEntity(organizationNode)

	if organization.QuickbooksCustomerId == "" {
		quickbooksSaveCustomerResponse, err := dependencies.CommonServices.QuickbooksService.SaveCustomer(ctx, "", organization.Name)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}

		if quickbooksSaveCustomerResponse.Customer == nil {
			err := errors.New("Quickbooks customer not found")
			tracing.TraceErr(span, err)
			return err
		}

		organization.QuickbooksCustomerId = quickbooksSaveCustomerResponse.Customer.Id

		err = dependencies.Neo4jRepositories.CommonWriteRepository.UpdateStringProperty(ctx, nil, tenant, model2.NodeLabelOrganization, organization.ID, string(neo4j_entity.OrganizationPropertyQuickbooksCustomerId), organization.QuickbooksCustomerId)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}
	}

	savedInvoiced, err := dependencies.CommonServices.QuickbooksService.SaveInvoice(ctx, organization.QuickbooksCustomerId, invoice.IssuedDate, quickbooksInvoiceLines)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if savedInvoiced == nil {
		err := errors.New("Invoice not saved in quickbooks")
		tracing.TraceErr(span, err)
		return err
	}

	err = dependencies.Neo4jRepositories.CommonWriteRepository.UpdateStringProperty(ctx, nil, tenant, model2.NodeLabelInvoice, invoice.Id, string(neo4j_entity.InvoicePropertyQuickbooksInvoiceId), savedInvoiced.Invoice.Id)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}

func OnInvoicePaid(ctx context.Context, dependencies *model.DependencyContainer, input any) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Listeners.OnInvoicePaid")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)

	message := input.(*dto.Event)
	invoiceId := message.Event.EntityId
	span.SetTag(tracing.SpanTagEntityId, invoiceId)

	tenant := common.GetTenantFromContext(ctx)
	if tenant == "" {
		err := errors.New("Missing tenant in context")
		tracing.TraceErr(span, err)
		return err
	}

	quickbooksSettingsEntity, err := dependencies.CommonServices.PostgresRepositories.QuickbooksSettingsRepository.Get(ctx, tenant)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if quickbooksSettingsEntity == nil {
		span.LogFields(log.String("error", "Quickbooks settings not found"))
		return nil
	}

	invoice, err := dependencies.CommonServices.InvoiceService.GetById(ctx, nil, invoiceId)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if invoice == nil {
		err := errors.New("Invoice not found")
		tracing.TraceErr(span, err)
		return err
	}

	organizationNode, err := dependencies.CommonServices.Neo4jRepositories.OrganizationReadRepository.GetOrganizationByInvoiceId(ctx, tenant, invoiceId)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if organizationNode == nil {
		err := errors.New("Organization not found for invoice")
		tracing.TraceErr(span, err)
		return err
	}

	organization := mapper.MapDbNodeToOrganizationEntity(organizationNode)

	paymentResponse, err := dependencies.CommonServices.QuickbooksService.PayInvoice(ctx, organization.QuickbooksCustomerId, invoice.QuickbooksInvoiceId, invoice.TotalAmount)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if paymentResponse == nil {
		err := errors.New("Invoice not paid in quickbooks")
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}

func OnInvoiceVoided(ctx context.Context, dependencies *model.DependencyContainer, input any) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Listeners.OnInvoiceVoided")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)

	message := input.(*dto.Event)
	invoiceId := message.Event.EntityId
	span.SetTag(tracing.SpanTagEntityId, invoiceId)

	tenant := common.GetTenantFromContext(ctx)
	if tenant == "" {
		err := errors.New("Missing tenant in context")
		tracing.TraceErr(span, err)
		return err
	}

	quickbooksSettingsEntity, err := dependencies.CommonServices.PostgresRepositories.QuickbooksSettingsRepository.Get(ctx, tenant)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if quickbooksSettingsEntity == nil {
		span.LogFields(log.String("error", "Quickbooks settings not found"))
		return nil
	}

	invoice, err := dependencies.CommonServices.InvoiceService.GetById(ctx, nil, invoiceId)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if invoice == nil {
		err := errors.New("Invoice not found")
		tracing.TraceErr(span, err)
		return err
	}

	voidedResponse, err := dependencies.CommonServices.QuickbooksService.VoidInvoice(ctx, invoice.QuickbooksInvoiceId)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if voidedResponse == nil {
		err := errors.New("Invoice not voided in quickbooks")
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}
