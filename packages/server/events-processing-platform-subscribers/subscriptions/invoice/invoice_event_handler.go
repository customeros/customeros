package invoice

import (
	"encoding/base64"
	"fmt"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/clients/grpc_client"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/dto"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/model"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/postmark"
	commontracing "github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"
	neo4jmapper "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/mapper"
	neo4j_repository "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/repository"
	postgres_repository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"
	"github.com/customeros/customeros/packages/server/events-processing-platform/domain/invoice"
	"github.com/customeros/customeros/packages/server/events/eventstore"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"io/ioutil"
	"net/mail"
	"os"
	"path/filepath"

	"github.com/customeros/customeros/packages/server/events-processing-platform-subscribers/config"
	"github.com/customeros/customeros/packages/server/events-processing-platform-subscribers/constants"
	"github.com/customeros/customeros/packages/server/events-processing-platform-subscribers/logger"
	"github.com/customeros/customeros/packages/server/events-processing-platform-subscribers/tracing"
)

type eventMetadata struct {
	UserId string `json:"user-id"`
}

type InvoiceActionMetadata struct {
	Status        string  `json:"status"`
	Currency      string  `json:"currency"`
	Amount        float64 `json:"amount"`
	InvoiceNumber string  `json:"number"`
	InvoiceId     string  `json:"id"`
}

type InvoiceEventHandler struct {
	log            logger.Logger
	cfg            config.Config
	grpcClients    *grpc_client.Clients
	neo4j          *neo4j_repository.Repositories
	postgres       *postgres_repository.Repositories
	invoice        interfaces.InvoiceService
	fileStore      interfaces.FileService
	postmark       interfaces.PostmarkService
	eventPublisher interfaces.EventPublisher
}

func NewInvoiceEventHandler(
	log logger.Logger,
	cfg config.Config,
	grpcClients *grpc_client.Clients,
	neo4j *neo4j_repository.Repositories,
	postgres *postgres_repository.Repositories,
	invoice interfaces.InvoiceService,
	fileStore interfaces.FileService,
	postmark interfaces.PostmarkService,
	eventPublisher interfaces.EventPublisher,
) *InvoiceEventHandler {
	return &InvoiceEventHandler{
		log:            log,
		cfg:            cfg,
		grpcClients:    grpcClients,
		neo4j:          neo4j,
		postgres:       postgres,
		invoice:        invoice,
		fileStore:      fileStore,
		postmark:       postmark,
		eventPublisher: eventPublisher,
	}
}

func (h *InvoiceEventHandler) onInvoiceVoidV1(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceEventHandler.onInvoiceVoidV1")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData invoice.InvoiceVoidEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "evt.GetJsonData"))
		return errors.Wrap(err, "evt.GetJsonData")
	}
	tracing.LogObjectAsJson(span, "eventData", eventData)
	invoiceId := invoice.GetInvoiceObjectID(evt.GetAggregateID(), eventData.Tenant)
	span.SetTag(tracing.SpanTagEntityId, invoiceId)
	span.SetTag(tracing.SpanTagTenant, eventData.Tenant)

	invoiceEntity, err := h.invoice.GetById(ctx, nil, invoiceId)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "InvoiceService.GetById"))
		return err
	}
	if invoiceEntity.DryRun {
		span.LogFields(log.String("result", "dry run, skipping"))
		return nil
	}

	invoiceLines, err := h.invoice.GetInvoiceLinesForInvoices(ctx, []string{invoiceId})
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "GetInvoiceLinesForInvoices"))
		return err
	}

	allInvoiceLinesEmpty := true
	for _, invoiceLine := range *invoiceLines {
		if invoiceLine.Amount != float64(0) {
			allInvoiceLinesEmpty = false
			break
		}
	}

	if invoiceEntity.DryRun || allInvoiceLinesEmpty {
		return nil
	}

	err = h.eventPublisher.PublishFanoutEvent(ctx, invoiceId, model.INVOICE, dto.InvoiceVoided{})
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "PublishEvent"))
		h.log.Errorf("error publishing invoice voided event for invoice %s: %s", invoiceId, err.Error())
	}

	// void notification already sent, skip
	if invoiceEntity.InvoiceInternalFields.VoidInvoiceNotificationSentAt != nil {
		return nil
	}

	// load contract
	contractEntity := neo4jentity.ContractEntity{}
	contractNode, err := h.neo4j.ContractReadRepository.GetContractForInvoice(ctx, eventData.Tenant, invoiceEntity.Id)
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

	// load tenant billing profile from neo4j
	tenantBillingProfileEntity, err := h.loadTenantBillingProfile(ctx, eventData.Tenant, false)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "loadTenantBillingProfile"))
		return nil
	}

	cc := contractEntity.InvoiceEmailCC
	cc = utils.RemoveEmpties(cc)
	cc = utils.RemoveDuplicates(cc)

	bcc := utils.AddToListIfNotExists(contractEntity.InvoiceEmailBCC, tenantBillingProfileEntity.SendInvoicesBcc)
	bcc = utils.RemoveEmpties(bcc)
	bcc = utils.RemoveDuplicates(bcc)

	postmarkEmail := interfaces.PostmarkEmail{
		WorkflowId:    postmark.WorkflowInvoiceVoided,
		MessageStream: postmark.PostmarkMessageStreamInvoice,
		From:          invoiceEntity.Provider.Email,
		To:            invoiceEntity.Customer.Email,
		CC:            cc,
		BCC:           bcc,
		Subject:       fmt.Sprintf(postmark.WorkflowInvoiceVoidedSubject, invoiceEntity.Number), // "Voided invoice " + invoiceEntity.Number,
		TemplateData: map[string]string{
			"{{userFirstName}}":  invoiceEntity.Customer.Name,
			"{{invoiceNumber}}":  invoiceEntity.Number,
			"{{currencySymbol}}": invoiceEntity.Currency.Symbol(),
			"{{amtDue}}":         fmt.Sprintf("%.2f", invoiceEntity.TotalAmount),
			"{{issueDate}}":      invoiceEntity.CreatedAt.Format("02 Jan 2006"),
		},
		Attachments: []interfaces.PostmarkEmailAttachment{},
	}

	err = h.appendProviderLogoToEmail(ctx, eventData.Tenant, invoiceEntity.Provider.LogoRepositoryFileId, &postmarkEmail)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "appendProviderLogoToEmail"))
		h.log.Errorf("Error appending provider logo to email for invoice %s: %s", invoiceId, err.Error())
		return nil
	}

	err = h.appendCustomerOSLogoToEmail(ctx, &postmarkEmail)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "appendCustomerOSLogoToEmail"))
		h.log.Errorf("Error appending customeros logo to email for invoice %s: %s", invoiceId, err.Error())
		return nil
	}

	err = h.postmark.SendNotification(ctx, postmarkEmail, eventData.Tenant)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "SendNotification"))
		h.log.Errorf("Error sending invoice voided notification for invoice %s: %s", invoiceId, err.Error())
		return nil
	}

	// Request was successful
	err = h.neo4j.InvoiceWriteRepository.SetVoidInvoiceNotificationSentAt(ctx, eventData.Tenant, invoiceId)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "SetVoidInvoiceNotificationSentAt"))
		h.log.Errorf("Error setting invoice void notification sent at for invoice %s: %s", invoiceId, err.Error())
		return nil
	}

	return nil
}

func (h *InvoiceEventHandler) onInvoicePayNotificationV1(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceEventHandler.onInvoicePayNotificationV1")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData invoice.InvoicePayNotificationEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "evt.GetJsonData"))
		return errors.Wrap(err, "evt.GetJsonData")
	}
	tracing.LogObjectAsJson(span, "eventData", eventData)
	invoiceId := invoice.GetInvoiceObjectID(evt.GetAggregateID(), eventData.Tenant)
	span.SetTag(tracing.SpanTagEntityId, invoiceId)
	span.SetTag(tracing.SpanTagTenant, eventData.Tenant)

	var invoiceEntity neo4jentity.InvoiceEntity
	var contractEntity neo4jentity.ContractEntity

	// load invoice entity
	invoiceNode, err := h.neo4j.InvoiceReadRepository.GetInvoiceById(ctx, nil, eventData.Tenant, invoiceId)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "GetInvoice"))
		return nil
	}
	if invoiceNode != nil {
		invoiceEntity = *neo4jmapper.MapDbNodeToInvoiceEntity(invoiceNode)
	} else {
		tracing.TraceErr(span, errors.New("invoiceNode is nil"))
		return nil
	}

	// Do not send email if invoice is dry run or total amount is 0 or invoice is not due or overdue
	invoiceStatusAllowedForPayNotification := invoiceEntity.IsDue() || invoiceEntity.IsOverdue()
	if invoiceEntity.DryRun || invoiceEntity.TotalAmount == float64(0) || !invoiceStatusAllowedForPayNotification {
		span.LogFields(log.String("result", "skipped pay notification"))
		return nil
	}

	if invoiceEntity.Provider.Email == "" {
		h.log.Warnf("Provider email address is empty for invoice %s", invoiceId)
		return nil
	}

	// load contract entity
	contractNode, err := h.neo4j.ContractReadRepository.GetContractForInvoice(ctx, eventData.Tenant, invoiceEntity.Id)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "GetContractForInvoice"))
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

	// load tenant billing profile from neo4j
	tenantBillingProfileEntity, err := h.loadTenantBillingProfile(ctx, eventData.Tenant, false)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "loadTenantBillingProfile"))
		return err
	}

	workflowId := ""
	if contractEntity.PayOnline || contractEntity.PayAutomatically {
		workflowId = postmark.WorkflowInvoiceReadyWithPaymentLink
	} else {
		workflowId = postmark.WorkflowInvoiceReadyNoPaymentLink
	}

	cc := contractEntity.InvoiceEmailCC
	cc = utils.RemoveEmpties(cc)
	cc = utils.RemoveDuplicates(cc)

	bcc := utils.AddToListIfNotExists(contractEntity.InvoiceEmailBCC, tenantBillingProfileEntity.SendInvoicesBcc)
	bcc = utils.RemoveEmpties(bcc)
	bcc = utils.RemoveDuplicates(bcc)

	paymentLink := ""
	// prepare payment link for email only if invoice payment link was generated
	if contractEntity.PayOnline || contractEntity.PayAutomatically {
		paymentLink = h.cfg.CommonServices.Internal.CustomerOsApi.ApiUrl + "/invoice/" + invoiceEntity.Id + "/pay"
	}

	postmarkEmail := interfaces.PostmarkEmail{
		WorkflowId:    workflowId,
		MessageStream: postmark.PostmarkMessageStreamInvoice,
		From:          invoiceEntity.Provider.Email,
		To:            contractEntity.InvoiceEmail,
		CC:            cc,
		BCC:           bcc,
		Subject:       fmt.Sprintf(postmark.WorkflowInvoiceReadySubject, invoiceEntity.Number),
		TemplateData: map[string]string{
			"{{organizationName}}": invoiceEntity.Customer.Name,
			"{{invoiceNumber}}":    invoiceEntity.Number,
			"{{currencySymbol}}":   invoiceEntity.Currency.Symbol(),
			"{{amtDue}}":           fmt.Sprintf("%.2f", invoiceEntity.TotalAmount),
			"{{paymentLink}}":      paymentLink,
		},
		Attachments: []interfaces.PostmarkEmailAttachment{},
	}

	err = h.appendInvoiceFileToEmailAsAttachment(ctx, eventData.Tenant, invoiceEntity, &postmarkEmail)
	if err != nil {
		wrappedErr := errors.Wrap(err, "InvoiceSubscriber.onInvoicePayNotificationV1.AppendInvoiceFileToEmailAsAttachment")
		tracing.TraceErr(span, wrappedErr)
		h.log.Errorf("Error appending invoice file to email attachment for invoice %s: %s", invoiceId, err.Error())
		return wrappedErr
	}

	err = h.appendProviderLogoToEmail(ctx, eventData.Tenant, invoiceEntity.Provider.LogoRepositoryFileId, &postmarkEmail)
	if err != nil {
		wrappedErr := errors.Wrap(err, "InvoiceSubscriber.onInvoicePayNotificationV1.AppendProviderLogoToEmail")
		tracing.TraceErr(span, wrappedErr)
		h.log.Errorf("Error appending provider logo to email for invoice %s: %s", invoiceId, err.Error())
		return wrappedErr
	}

	err = h.appendCustomerOSLogoToEmail(ctx, &postmarkEmail)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "InvoiceSubscriber.onInvoicePayNotificationV1.AppendCustomerOSLogoToEmail"))
		h.log.Errorf("Error appending customeros logo to email for invoice %s: %s", invoiceId, err.Error())
		return nil
	}

	err = h.postmark.SendNotification(ctx, postmarkEmail, eventData.Tenant)
	if err != nil {
		wrappedErr := errors.Wrap(err, "InvoiceSubscriber.onInvoicePayNotificationV1.SendNotification")
		tracing.TraceErr(span, wrappedErr)
		h.log.Errorf("Error sending invoice pay request notification for invoice %s: %s", invoiceId, err.Error())
		return nil
	}

	h.createInvoiceAction(ctx, eventData.Tenant, invoiceEntity)

	// Request was successful
	err = h.neo4j.InvoiceWriteRepository.SetPayInvoiceNotificationSentAt(ctx, eventData.Tenant, invoiceId)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "SetPayInvoiceNotificationSentAt"))
		h.log.Errorf("Error setting invoice pay notification sent at for invoice %s: %s", invoiceId, err.Error())
		return err
	}

	return nil
}

func (h *InvoiceEventHandler) onInvoiceRemindNotificationV1(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceEventHandler.onInvoiceRemindNotificationV1")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData invoice.InvoiceRemindNotificationEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "evt.GetJsonData"))
		return errors.Wrap(err, "evt.GetJsonData")
	}
	tracing.LogObjectAsJson(span, "eventData", eventData)
	invoiceId := invoice.GetInvoiceObjectID(evt.GetAggregateID(), eventData.Tenant)
	span.SetTag(tracing.SpanTagEntityId, invoiceId)
	span.SetTag(tracing.SpanTagTenant, eventData.Tenant)

	var invoiceEntity neo4jentity.InvoiceEntity
	var contractEntity neo4jentity.ContractEntity

	// load invoice
	invoiceNode, err := h.neo4j.InvoiceReadRepository.GetInvoiceById(ctx, nil, eventData.Tenant, invoiceId)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "GetInvoice"))
		return nil
	}
	if invoiceNode != nil {
		invoiceEntity = *neo4jmapper.MapDbNodeToInvoiceEntity(invoiceNode)
	} else {
		tracing.TraceErr(span, errors.New("invoiceNode is nil"))
		return nil
	}

	if invoiceEntity.DryRun || invoiceEntity.TotalAmount == float64(0) || !invoiceEntity.IsOverdue() {
		tracing.TraceErr(span, errors.New("remind invoice notification requested for not applicable invoice"))
		return nil
	}

	if invoiceEntity.Provider.Email == "" {
		h.log.Warnf("Provider email address is empty for invoice %s", invoiceId)
		return nil
	}

	// load contract
	contractNode, err := h.neo4j.ContractReadRepository.GetContractForInvoice(ctx, eventData.Tenant, invoiceEntity.Id)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "GetContractForInvoice"))
		return errors.Wrap(err, "InvoiceSubscriber.onInvoiceRemindNotificationV1.GetContractForInvoice")
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

	// load tenant billing profile from neo4j
	tenantBillingProfileEntity, err := h.loadTenantBillingProfile(ctx, eventData.Tenant, false)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "loadTenantBillingProfile"))
		return err
	}
	tenantSettingsDbNode, err := h.neo4j.TenantReadRepository.GetTenantSettings(ctx, eventData.Tenant)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "GetTenantSettings"))
		return err
	}
	tenantSettingsEntity := neo4jmapper.MapDbNodeToTenantSettingsEntity(tenantSettingsDbNode)

	workflowId := ""
	if invoiceEntity.PaymentDetails.PaymentLink == "" {
		workflowId = postmark.WorkflowInvoiceRemindNoPaymentLink
	} else {
		workflowId = postmark.WorkflowInvoiceRemindWithPaymentLink
	}

	cc := contractEntity.InvoiceEmailCC
	cc = utils.RemoveEmpties(cc)
	cc = utils.RemoveDuplicates(cc)

	bcc := utils.AddToListIfNotExists(contractEntity.InvoiceEmailBCC, tenantBillingProfileEntity.SendInvoicesBcc)
	bcc = utils.RemoveEmpties(bcc)
	bcc = utils.RemoveDuplicates(bcc)

	paymentLink := ""
	// prepare payment link for email only if invoice payment link was generated
	if invoiceEntity.PaymentDetails.PaymentLink != "" {
		paymentLink = h.cfg.CommonServices.Internal.CustomerOsApi.ApiUrl + "/invoice/" + invoiceEntity.Id + "/pay"
	}

	postmarkEmail := interfaces.PostmarkEmail{
		WorkflowId:    workflowId,
		MessageStream: postmark.PostmarkMessageStreamInvoice,
		From:          invoiceEntity.Provider.Email,
		To:            contractEntity.InvoiceEmail,
		CC:            cc,
		BCC:           bcc,
		Subject:       fmt.Sprintf(postmark.WorkflowInvoiceRemindSubject, invoiceEntity.Number),
		TemplateData: map[string]string{
			"{{organizationName}}": invoiceEntity.Customer.Name,
			"{{invoiceNumber}}":    invoiceEntity.Number,
			"{{currencySymbol}}":   invoiceEntity.Currency.Symbol(),
			"{{amtDue}}":           fmt.Sprintf("%.2f", invoiceEntity.TotalAmount),
			"{{paymentLink}}":      paymentLink,
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

	err = h.appendInvoiceFileToEmailAsAttachment(ctx, eventData.Tenant, invoiceEntity, &postmarkEmail)
	if err != nil {
		wrappedErr := errors.Wrap(err, "InvoiceSubscriber.onInvoiceRemindNotificationV1.AppendInvoiceFileToEmailAsAttachment")
		tracing.TraceErr(span, wrappedErr)
		h.log.Errorf("Error appending invoice file to email attachment for invoice %s: %s", invoiceId, err.Error())
		return wrappedErr
	}

	err = h.appendProviderLogoToEmail(ctx, eventData.Tenant, invoiceEntity.Provider.LogoRepositoryFileId, &postmarkEmail)
	if err != nil {
		wrappedErr := errors.Wrap(err, "InvoiceSubscriber.onInvoiceRemindNotificationV1.AppendProviderLogoToEmail")
		tracing.TraceErr(span, wrappedErr)
		h.log.Errorf("Error appending provider logo to email for invoice %s: %s", invoiceId, err.Error())
		return wrappedErr
	}

	err = h.appendCustomerOSLogoToEmail(ctx, &postmarkEmail)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "InvoiceSubscriber.onInvoiceRemindNotificationV1.AppendCustomerOSLogoToEmail"))
		h.log.Errorf("Error appending customeros logo to email for invoice %s: %s", invoiceId, err.Error())
		return nil
	}

	err = h.postmark.SendNotification(ctx, postmarkEmail, eventData.Tenant)
	if err != nil {
		wrappedErr := errors.Wrap(err, "InvoiceSubscriber.onInvoiceRemindNotificationV1.SendNotification")
		tracing.TraceErr(span, wrappedErr)
		h.log.Errorf("Error sending invoice remind request notification for invoice %s: %s", invoiceId, err.Error())
		return nil
	}

	// Request was successful
	err = h.neo4j.CommonWriteRepository.UpdateTimeProperty(ctx, eventData.Tenant, model.NodeLabelInvoice, invoiceEntity.Id, string(neo4jentity.InvoicePropertyLastRemindInvoiceNotificationSentAt), utils.NowPtr())
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "UpdateTimeProperty"))
		h.log.Errorf("Error setting invoice remind notification sent at for invoice %s: %s", invoiceId, err.Error())
		return err
	}

	return nil
}

func (h *InvoiceEventHandler) createInvoiceAction(ctx context.Context, tenant string, invoiceEntity neo4jentity.InvoiceEntity) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceEventHandler.createInvoiceAction")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	span.LogFields(log.String("InvoiceId", invoiceEntity.Id))

	if invoiceEntity.DryRun || invoiceEntity.TotalAmount == float64(0) {
		return
	}

	metadata, err := utils.ToJson(InvoiceActionMetadata{
		Status:        invoiceEntity.Status.String(),
		Currency:      invoiceEntity.Currency.String(),
		Amount:        invoiceEntity.TotalAmount,
		InvoiceNumber: invoiceEntity.Number,
		InvoiceId:     invoiceEntity.Id,
	})

	actionType := enum.ActionInvoiceSent
	message := "Sent invoice N° " + invoiceEntity.Number + " with an amount of " + invoiceEntity.Currency.Symbol() + utils.FormatAmount(invoiceEntity.TotalAmount, 2)

	_, err = h.neo4j.ActionWriteRepository.MergeByActionType(ctx, nil, tenant, invoiceEntity.Id, model.INVOICE, actionType, message, metadata, utils.Now(), constants.AppSourceEventProcessingPlatformSubscribers)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "ActionWriteRepository.MergeByActionType"))
		h.log.Errorf("Failed creating invoice action for invoice %s: %s", invoiceEntity.Id, err.Error())
	}
}

func (h *InvoiceEventHandler) appendInvoiceFileToEmailAsAttachment(ctx context.Context, tenant string, invoice neo4jentity.InvoiceEntity, postmarkEmail *interfaces.PostmarkEmail) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceEventHandler.appendInvoiceFileToEmailAsAttachment")
	defer span.Finish()
	commontracing.TagTenant(span, tenant)

	fileInfo, err := h.fileStore.GetById(ctx, invoice.RepositoryFileId)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	if fileInfo == nil {
		return nil
	}

	invoiceFileBytes, err := h.fileStore.GetFileBytes(ctx, invoice.RepositoryFileId)
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

func (h *InvoiceEventHandler) appendProviderLogoToEmail(ctx context.Context, tenant, logoFileId string, postmarkEmail *interfaces.PostmarkEmail) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceEventHandler.appendProviderLogoToEmail")
	defer span.Finish()

	if logoFileId == "" {
		return nil
	}

	fileInfo, err := h.fileStore.GetById(ctx, logoFileId)
	if err != nil {
		return err
	}
	if fileInfo == nil {
		return nil
	}

	fileBytes, err := h.fileStore.GetFileBytes(ctx, logoFileId)
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

func (h *InvoiceEventHandler) appendCustomerOSLogoToEmail(ctx context.Context, postmarkEmail *interfaces.PostmarkEmail) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceEventHandler.appendCustomerOSLogoToEmail")
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

func isValidEmailSyntax(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func (h *InvoiceEventHandler) loadTenantBillingProfile(ctx context.Context, tenant string, failIfNotFound bool) (neo4jentity.TenantBillingProfileEntity, error) {
	tenantBillingProfiles, err := h.neo4j.TenantReadRepository.GetTenantBillingProfiles(ctx, tenant)
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
