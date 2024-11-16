package invoice

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/tenant/event"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	tenantpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/tenant"
	"github.com/openline-ai/openline-customer-os/packages/server/events/event/common"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
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

func (a *TenantAggregate) When(evt eventstore.Event) error {
	switch evt.GetEventType() {
	case event.TenantAddBillingProfileV1:
		return a.onAddBillingProfile(evt)
	case event.TenantUpdateBillingProfileV1:
		return a.onUpdateBillingProfile(evt)
	default:
		return nil
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
