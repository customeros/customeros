package invoice

import (
	"context"
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/tenant/event"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	tenantpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/tenant"
	"github.com/openline-ai/openline-customer-os/packages/server/events/events/common"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	events2 "github.com/openline-ai/openline-customer-os/packages/server/events/utils"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"strings"
)

const (
	TenantAggregateType eventstore.AggregateType = "tenant"
)

type TenantAggregate struct {
	*eventstore.CommonIdAggregate
	TenantDetails *Tenant
}

func GetTenantName(aggregateID string) string {
	return strings.ReplaceAll(aggregateID, string(TenantAggregateType)+"-", "")
}

func NewTenantAggregate(tenant string) *TenantAggregate {
	tenantAggregate := TenantAggregate{}
	tenantAggregate.CommonIdAggregate = eventstore.NewCommonAggregateWithId(TenantAggregateType, tenant)
	tenantAggregate.SetWhen(tenantAggregate.When)
	tenantAggregate.TenantDetails = &Tenant{}
	tenantAggregate.Tenant = tenant

	return &tenantAggregate
}

func (a *TenantAggregate) HandleGRPCRequest(ctx context.Context, request any, params map[string]any) (any, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TenantAggregate.HandleGRPCRequest")
	defer span.Finish()

	switch r := request.(type) {
	case *tenantpb.AddBillingProfileRequest:
		return a.AddBillingProfile(ctx, r)
	case *tenantpb.UpdateBillingProfileRequest:
		return r.Id, a.UpdateBillingProfile(ctx, r)
	case *tenantpb.UpdateTenantSettingsRequest:
		return nil, a.UpdateTenantSettings(ctx, r)
	case *tenantpb.AddBankAccountGrpcRequest:
		return a.AddBankAccount(ctx, r)
	case *tenantpb.UpdateBankAccountGrpcRequest:
		return r.Id, a.UpdateBankAccount(ctx, r)
	case *tenantpb.DeleteBankAccountGrpcRequest:
		return nil, a.DeleteBankAccount(ctx, r)
	default:
		return nil, nil
	}
}

func (a *TenantAggregate) AddBillingProfile(ctx context.Context, request *tenantpb.AddBillingProfileRequest) (string, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "TenantAggregate.AddBillingProfile")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.GetTenant())
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("AggregateVersion", a.GetVersion()))

	sourceFields := common.Source{}
	sourceFields.FromGrpc(request.SourceFields)

	createdAtNotNil := utils.IfNotNilTimeWithDefault(utils.TimestampProtoToTimePtr(request.CreatedAt), utils.Now())

	billingProfileId := uuid.New().String()

	addBillingProfileEvent, err := event.NewTenantBillingProfileCreateEvent(a, sourceFields, billingProfileId, request, createdAtNotNil)
	if err != nil {
		tracing.TraceErr(span, err)
		return "", errors.Wrap(err, "TenantBillingProfileCreateEvent")
	}
	eventstore.EnrichEventWithMetadataExtended(&addBillingProfileEvent, span, eventstore.EventMetadata{
		Tenant: request.Tenant,
		UserId: request.LoggedInUserId,
		App:    sourceFields.AppSource,
	})

	return billingProfileId, a.Apply(addBillingProfileEvent)
}

func (a *TenantAggregate) UpdateBillingProfile(ctx context.Context, r *tenantpb.UpdateBillingProfileRequest) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "TenantAggregate.UpdateBillingProfile")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.GetTenant())
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("AggregateVersion", a.GetVersion()))

	updatedAtNotNil := utils.IfNotNilTimeWithDefault(utils.TimestampProtoToTimePtr(r.UpdatedAt), utils.Now())
	fieldsMaks := extractTenantBillingProfileFieldsMask(r.FieldsMask)

	updateBillingProfileEvent, err := event.NewTenantBillingProfileUpdateEvent(a, r.Id, r, updatedAtNotNil, fieldsMaks)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "TenantBillingProfileUpdateEvent")
	}
	eventstore.EnrichEventWithMetadataExtended(&updateBillingProfileEvent, span, eventstore.EventMetadata{
		Tenant: r.Tenant,
		UserId: r.LoggedInUserId,
		App:    r.AppSource,
	})

	return a.Apply(updateBillingProfileEvent)
}

func (a *TenantAggregate) UpdateTenantSettings(ctx context.Context, r *tenantpb.UpdateTenantSettingsRequest) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "TenantAggregate.UpdateTenantSettings")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.GetTenant())
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("AggregateVersion", a.GetVersion()))

	updatedAtNotNil := utils.IfNotNilTimeWithDefault(utils.TimestampProtoToTimePtr(r.UpdatedAt), utils.Now())
	fieldsMaks := extractTenantSettingsFieldsMask(r.FieldsMask)

	updateSettingsEvent, err := event.NewTenantSettingsUpdateEvent(a, r, updatedAtNotNil, fieldsMaks)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "TenantSettingsUpdateEvent")
	}
	eventstore.EnrichEventWithMetadataExtended(&updateSettingsEvent, span, eventstore.EventMetadata{
		Tenant: r.Tenant,
		UserId: r.LoggedInUserId,
		App:    r.AppSource,
	})

	return a.Apply(updateSettingsEvent)
}

func (a *TenantAggregate) AddBankAccount(ctx context.Context, r *tenantpb.AddBankAccountGrpcRequest) (string, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "TenantAggregate.AddBankAccount")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.GetTenant())
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("AggregateVersion", a.GetVersion()))

	sourceFields := common.Source{}
	sourceFields.FromGrpc(r.SourceFields)

	createdAtNotNil := utils.IfNotNilTimeWithDefault(utils.TimestampProtoToTimePtr(r.CreatedAt), utils.Now())

	bankAccountId := uuid.New().String()

	addBankAccountEvent, err := event.NewTenantBankAccountCreateEvent(a, sourceFields, bankAccountId, r, createdAtNotNil)
	if err != nil {
		tracing.TraceErr(span, err)
		return "", errors.Wrap(err, "TenantBankAccountCreateEvent")
	}
	eventstore.EnrichEventWithMetadataExtended(&addBankAccountEvent, span, eventstore.EventMetadata{
		Tenant: r.Tenant,
		UserId: r.LoggedInUserId,
		App:    sourceFields.AppSource,
	})

	return bankAccountId, a.Apply(addBankAccountEvent)
}

func (a *TenantAggregate) UpdateBankAccount(ctx context.Context, r *tenantpb.UpdateBankAccountGrpcRequest) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "TenantAggregate.UpdateBankAccount")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.GetTenant())
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("AggregateVersion", a.GetVersion()))

	updatedAtNotNil := utils.IfNotNilTimeWithDefault(utils.TimestampProtoToTimePtr(r.UpdatedAt), utils.Now())
	fieldsMaks := extractTenantBankAccountFieldsMask(r.FieldsMask)

	updateBankAccountEvent, err := event.NewTenantBankAccountUpdateEvent(a, r.Id, r, updatedAtNotNil, fieldsMaks)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "TenantBankAccountUpdateEvent")
	}
	eventstore.EnrichEventWithMetadataExtended(&updateBankAccountEvent, span, eventstore.EventMetadata{
		Tenant: r.Tenant,
		UserId: r.LoggedInUserId,
		App:    r.AppSource,
	})

	return a.Apply(updateBankAccountEvent)
}

func (a *TenantAggregate) DeleteBankAccount(ctx context.Context, r *tenantpb.DeleteBankAccountGrpcRequest) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "TenantAggregate.DeleteBankAccount")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.GetTenant())
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("AggregateVersion", a.GetVersion()))

	deleteBankAccountEvent, err := event.NewTenantBankAccountDeleteEvent(a, r.Id, utils.Now())
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "TenantBankAccountDeleteEvent")
	}
	eventstore.EnrichEventWithMetadataExtended(&deleteBankAccountEvent, span, eventstore.EventMetadata{
		Tenant: r.Tenant,
		UserId: r.LoggedInUserId,
		App:    r.AppSource,
	})

	return a.Apply(deleteBankAccountEvent)
}

func (a *TenantAggregate) When(evt eventstore.Event) error {
	switch evt.GetEventType() {
	case event.TenantAddBillingProfileV1:
		return a.onAddBillingProfile(evt)
	case event.TenantUpdateBillingProfileV1:
		return a.onUpdateBillingProfile(evt)
	case event.TenantUpdateSettingsV1:
		return a.onUpdateTenantSettings(evt)
	case event.TenantAddBankAccountV1:
		return a.onAddBankAccount(evt)
	case event.TenantUpdateBankAccountV1:
		return a.onUpdateBankAccount(evt)
	case event.TenantDeleteBankAccountV1:
		return a.onDeleteBankAccount(evt)
	default:
		if strings.HasPrefix(evt.GetEventType(), events2.EsInternalStreamPrefix) {
			return nil
		}
		err := eventstore.ErrInvalidEventType
		err.EventType = evt.GetEventType()
		return err
	}
}

func (a *TenantAggregate) onAddBillingProfile(evt eventstore.Event) error {
	var eventData event.TenantBillingProfileCreateEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}

	if a.TenantDetails.HasBillingProfile(eventData.Id) {
		return nil
	}
	tenantBillingProfile := TenantBillingProfile{
		Id:                     eventData.Id,
		CreatedAt:              eventData.CreatedAt,
		Phone:                  eventData.Phone,
		AddressLine1:           eventData.AddressLine1,
		AddressLine2:           eventData.AddressLine2,
		AddressLine3:           eventData.AddressLine3,
		Locality:               eventData.Locality,
		Country:                eventData.Country,
		Region:                 eventData.Region,
		Zip:                    eventData.Zip,
		LegalName:              eventData.LegalName,
		VatNumber:              eventData.VatNumber,
		SendInvoicesFrom:       eventData.SendInvoicesFrom,
		SendInvoicesBcc:        eventData.SendInvoicesBcc,
		CanPayWithPigeon:       eventData.CanPayWithPigeon,
		CanPayWithBankTransfer: eventData.CanPayWithBankTransfer,
		SourceFields:           eventData.SourceFields,
		Check:                  eventData.Check,
	}
	a.TenantDetails.BillingProfiles = append(a.TenantDetails.BillingProfiles, tenantBillingProfile)

	return nil
}

func (a *TenantAggregate) onUpdateBillingProfile(evt eventstore.Event) error {
	var eventData event.TenantBillingProfileUpdateEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}

	if !a.TenantDetails.HasBillingProfile(eventData.Id) {
		tenantBillingProfile := TenantBillingProfile{
			Id: eventData.Id,
		}
		a.TenantDetails.BillingProfiles = append(a.TenantDetails.BillingProfiles, tenantBillingProfile)
	}

	tenantBillingProfile := a.TenantDetails.GetBillingProfile(eventData.Id)
	if eventData.UpdatePhone() {
		tenantBillingProfile.Phone = eventData.Phone
	}
	if eventData.UpdateAddressLine1() {
		tenantBillingProfile.AddressLine1 = eventData.AddressLine1
	}
	if eventData.UpdateAddressLine2() {
		tenantBillingProfile.AddressLine2 = eventData.AddressLine2
	}
	if eventData.UpdateAddressLine3() {
		tenantBillingProfile.AddressLine3 = eventData.AddressLine3
	}
	if eventData.UpdateLocality() {
		tenantBillingProfile.Locality = eventData.Locality
	}
	if eventData.UpdateCountry() {
		tenantBillingProfile.Country = eventData.Country
	}
	if eventData.UpdateRegion() {
		tenantBillingProfile.Region = eventData.Region
	}
	if eventData.UpdateZip() {
		tenantBillingProfile.Zip = eventData.Zip
	}
	if eventData.UpdateLegalName() {
		tenantBillingProfile.LegalName = eventData.LegalName
	}
	if eventData.UpdateVatNumber() {
		tenantBillingProfile.VatNumber = eventData.VatNumber
	}
	if eventData.UpdateSendInvoicesFrom() {
		tenantBillingProfile.SendInvoicesFrom = eventData.SendInvoicesFrom
	}
	if eventData.UpdateSendInvoicesBcc() {
		tenantBillingProfile.SendInvoicesBcc = eventData.SendInvoicesBcc
	}
	if eventData.UpdateCanPayWithPigeon() {
		tenantBillingProfile.CanPayWithPigeon = eventData.CanPayWithPigeon
	}
	if eventData.UpdateCanPayWithBankTransfer() {
		tenantBillingProfile.CanPayWithBankTransfer = eventData.CanPayWithBankTransfer
	}
	if eventData.UpdateCheck() {
		tenantBillingProfile.Check = eventData.Check
	}
	return nil
}

func (a *TenantAggregate) onUpdateTenantSettings(evt eventstore.Event) error {
	var eventData event.TenantSettingsUpdateEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}

	if eventData.UpdateBaseCurrency() {
		a.TenantDetails.TenantSettings.BaseCurrency = eventData.BaseCurrency
	}
	if eventData.UpdateInvoicingEnabled() {
		a.TenantDetails.TenantSettings.InvoicingEnabled = eventData.InvoicingEnabled
	}
	if eventData.UpdateLogoRepositoryFileId() {
		a.TenantDetails.TenantSettings.LogoRepositoryFileId = eventData.LogoRepositoryFileId
	}
	return nil
}

func (a *TenantAggregate) onAddBankAccount(evt eventstore.Event) error {
	var eventData event.TenantBankAccountCreateEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}

	if a.TenantDetails.HasBankAccount(eventData.Id) {
		return nil
	}
	bankAccount := BankAccount{
		Id:                  eventData.Id,
		CreatedAt:           eventData.CreatedAt,
		BankName:            eventData.BankName,
		BankTransferEnabled: eventData.BankTransferEnabled,
		AllowInternational:  eventData.AllowInternational,
		Currency:            eventData.Currency,
		Iban:                eventData.Iban,
		Bic:                 eventData.Bic,
		SortCode:            eventData.SortCode,
		AccountNumber:       eventData.AccountNumber,
		RoutingNumber:       eventData.RoutingNumber,
		OtherDetails:        eventData.OtherDetails,
		SourceFields:        eventData.SourceFields,
	}
	a.TenantDetails.BankAccounts = append(a.TenantDetails.BankAccounts, bankAccount)

	return nil
}

func (a *TenantAggregate) onUpdateBankAccount(evt eventstore.Event) error {
	var eventData event.TenantBankAccountUpdateEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}

	if !a.TenantDetails.HasBankAccount(eventData.Id) {
		bankAccount := BankAccount{
			Id: eventData.Id,
		}
		a.TenantDetails.BankAccounts = append(a.TenantDetails.BankAccounts, bankAccount)
	}

	bankAccount := a.TenantDetails.GetBankAccount(eventData.Id)
	if eventData.UpdateBankName() {
		bankAccount.BankName = eventData.BankName
	}
	if eventData.UpdateBankTransferEnabled() {
		bankAccount.BankTransferEnabled = eventData.BankTransferEnabled
	}
	if eventData.UpdateAllowInternational() {
		bankAccount.AllowInternational = eventData.AllowInternational
	}
	if eventData.UpdateCurrency() {
		bankAccount.Currency = eventData.Currency
	}
	if eventData.UpdateIban() {
		bankAccount.Iban = eventData.Iban
	}
	if eventData.UpdateBic() {
		bankAccount.Bic = eventData.Bic
	}
	if eventData.UpdateSortCode() {
		bankAccount.SortCode = eventData.SortCode
	}
	if eventData.UpdateAccountNumber() {
		bankAccount.AccountNumber = eventData.AccountNumber
	}
	if eventData.UpdateRoutingNumber() {
		bankAccount.RoutingNumber = eventData.RoutingNumber
	}
	if eventData.UpdateOtherDetails() {
		bankAccount.OtherDetails = eventData.OtherDetails
	}
	return nil
}

func (a *TenantAggregate) onDeleteBankAccount(evt eventstore.Event) error {
	var eventData event.TenantBankAccountDeleteEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}

	for i, bankAccount := range a.TenantDetails.BankAccounts {
		if bankAccount.Id == eventData.Id {
			a.TenantDetails.BankAccounts = append(a.TenantDetails.BankAccounts[:i], a.TenantDetails.BankAccounts[i+1:]...)
			break
		}
	}
	return nil
}

func extractTenantBillingProfileFieldsMask(requestFieldsMask []tenantpb.TenantBillingProfileFieldMask) []string {
	var fieldsMask []string
	for _, requestFieldMask := range requestFieldsMask {
		switch requestFieldMask {
		case tenantpb.TenantBillingProfileFieldMask_TENANT_BILLING_PROFILE_FIELD_PHONE:
			fieldsMask = append(fieldsMask, event.FieldMaskPhone)
		case tenantpb.TenantBillingProfileFieldMask_TENANT_BILLING_PROFILE_FIELD_ADDRESS_LINE_1:
			fieldsMask = append(fieldsMask, event.FieldMaskAddressLine1)
		case tenantpb.TenantBillingProfileFieldMask_TENANT_BILLING_PROFILE_FIELD_ADDRESS_LINE_2:
			fieldsMask = append(fieldsMask, event.FieldMaskAddressLine2)
		case tenantpb.TenantBillingProfileFieldMask_TENANT_BILLING_PROFILE_FIELD_ADDRESS_LINE_3:
			fieldsMask = append(fieldsMask, event.FieldMaskAddressLine3)
		case tenantpb.TenantBillingProfileFieldMask_TENANT_BILLING_PROFILE_FIELD_LOCALITY:
			fieldsMask = append(fieldsMask, event.FieldMaskLocality)
		case tenantpb.TenantBillingProfileFieldMask_TENANT_BILLING_PROFILE_FIELD_COUNTRY:
			fieldsMask = append(fieldsMask, event.FieldMaskCountry)
		case tenantpb.TenantBillingProfileFieldMask_TENANT_BILLING_PROFILE_FIELD_REGION:
			fieldsMask = append(fieldsMask, event.FieldMaskRegion)
		case tenantpb.TenantBillingProfileFieldMask_TENANT_BILLING_PROFILE_FIELD_ZIP:
			fieldsMask = append(fieldsMask, event.FieldMaskZip)
		case tenantpb.TenantBillingProfileFieldMask_TENANT_BILLING_PROFILE_FIELD_LEGAL_NAME:
			fieldsMask = append(fieldsMask, event.FieldMaskLegalName)
		case tenantpb.TenantBillingProfileFieldMask_TENANT_BILLING_PROFILE_FIELD_VAT_NUMBER:
			fieldsMask = append(fieldsMask, event.FieldMaskVatNumber)
		case tenantpb.TenantBillingProfileFieldMask_TENANT_BILLING_PROFILE_FIELD_SEND_INVOICES_FROM:
			fieldsMask = append(fieldsMask, event.FieldMaskSendInvoicesFrom)
		case tenantpb.TenantBillingProfileFieldMask_TENANT_BILLING_PROFILE_FIELD_SEND_INVOICES_BCC:
			fieldsMask = append(fieldsMask, event.FieldMaskSendInvoicesBcc)
		case tenantpb.TenantBillingProfileFieldMask_TENANT_BILLING_PROFILE_FIELD_CAN_PAY_WITH_PIGEON:
			fieldsMask = append(fieldsMask, event.FieldMaskCanPayWithPigeon)
		case tenantpb.TenantBillingProfileFieldMask_TENANT_BILLING_PROFILE_FIELD_CAN_PAY_WITH_BANK_TRANSFER:
			fieldsMask = append(fieldsMask, event.FieldMaskCanPayWithBankTransfer)
		case tenantpb.TenantBillingProfileFieldMask_TENANT_BILLING_PROFILE_FIELD_CHECK:
			fieldsMask = append(fieldsMask, event.FieldMaskCheck)
		}
	}
	fieldsMask = utils.RemoveDuplicates(fieldsMask)
	return fieldsMask
}

func extractTenantSettingsFieldsMask(inputFieldsMask []tenantpb.TenantSettingsFieldMask) []string {
	var fieldsMask []string
	for _, requestFieldMask := range inputFieldsMask {
		switch requestFieldMask {
		case tenantpb.TenantSettingsFieldMask_TENANT_SETTINGS_FIELD_LOGO_REPOSITORY_FILE_ID:
			fieldsMask = append(fieldsMask, event.FieldMaskLogoRepositoryFileId)
		case tenantpb.TenantSettingsFieldMask_TENANT_SETTINGS_FIELD_BASE_CURRENCY:
			fieldsMask = append(fieldsMask, event.FieldMaskBaseCurrency)
		case tenantpb.TenantSettingsFieldMask_TENANT_SETTINGS_FIELD_INVOICING_ENABLED:
			fieldsMask = append(fieldsMask, event.FieldMaskInvoicingEnabled)
		case tenantpb.TenantSettingsFieldMask_TENANT_SETTINGS_FIELD_INVOICING_POSTPAID:
			fieldsMask = append(fieldsMask, event.FieldMaskInvoicingPostpaid)
		}
	}
	fieldsMask = utils.RemoveDuplicates(fieldsMask)
	return fieldsMask
}

func extractTenantBankAccountFieldsMask(inputFieldsMask []tenantpb.BankAccountFieldMask) []string {
	var fieldsMask []string
	for _, requestFieldMask := range inputFieldsMask {
		switch requestFieldMask {
		case tenantpb.BankAccountFieldMask_BANK_ACCOUNT_FIELD_BANK_NAME:
			fieldsMask = append(fieldsMask, event.FieldMaskBankAccountBankName)
		case tenantpb.BankAccountFieldMask_BANK_ACCOUNT_FIELD_BANK_TRANSFER_ENABLED:
			fieldsMask = append(fieldsMask, event.FieldMaskBankAccountBankTransferEnabled)
		case tenantpb.BankAccountFieldMask_BANK_ACCOUNT_FIELD_ALLOW_INTERNATIONAL:
			fieldsMask = append(fieldsMask, event.FieldMaskBankAccountAllowInternational)
		case tenantpb.BankAccountFieldMask_BANK_ACCOUNT_FIELD_CURRENCY:
			fieldsMask = append(fieldsMask, event.FieldMaskBankAccountCurrency)
		case tenantpb.BankAccountFieldMask_BANK_ACCOUNT_FIELD_IBAN:
			fieldsMask = append(fieldsMask, event.FieldMaskBankAccountIban)
		case tenantpb.BankAccountFieldMask_BANK_ACCOUNT_FIELD_BIC:
			fieldsMask = append(fieldsMask, event.FieldMaskBankAccountBic)
		case tenantpb.BankAccountFieldMask_BANK_ACCOUNT_FIELD_SORT_CODE:
			fieldsMask = append(fieldsMask, event.FieldMaskBankAccountSortCode)
		case tenantpb.BankAccountFieldMask_BANK_ACCOUNT_FIELD_ACCOUNT_NUMBER:
			fieldsMask = append(fieldsMask, event.FieldMaskBankAccountAccountNumber)
		case tenantpb.BankAccountFieldMask_BANK_ACCOUNT_FIELD_ROUTING_NUMBER:
			fieldsMask = append(fieldsMask, event.FieldMaskBankAccountRoutingNumber)
		case tenantpb.BankAccountFieldMask_BANK_ACCOUNT_FIELD_OTHER_DETAILS:
			fieldsMask = append(fieldsMask, event.FieldMaskBankAccountOtherDetails)
		}
	}
	fieldsMask = utils.RemoveDuplicates(fieldsMask)
	return fieldsMask
}
