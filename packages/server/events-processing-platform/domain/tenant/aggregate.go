package invoice

import (
	"context"
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/aggregate"
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/tenant/event"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	tenantpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/tenant"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"strings"
)

const (
	TenantAggregateType eventstore.AggregateType = "tenant"
)

type TenantAggregate struct {
	*aggregate.CommonIdAggregate
	TenantDetails *Tenant
}

func GetTenantName(aggregateID string) string {
	return strings.ReplaceAll(aggregateID, string(TenantAggregateType)+"-", "")
}

func LoadTenantAggregate(ctx context.Context, eventStore eventstore.AggregateStore, tenant string, options eventstore.LoadAggregateOptions) (*TenantAggregate, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "LoadTenantAggregate")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)

	tenantAggregate := NewTenantAggregate(tenant)

	err := aggregate.LoadAggregate(ctx, eventStore, tenantAggregate, options)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return tenantAggregate, nil
}

func NewTenantAggregate(tenant string) *TenantAggregate {
	tenantAggregate := TenantAggregate{}
	tenantAggregate.CommonIdAggregate = aggregate.NewCommonAggregateWithId(TenantAggregateType, tenant)
	tenantAggregate.SetWhen(tenantAggregate.When)
	tenantAggregate.TenantDetails = &Tenant{}
	tenantAggregate.Tenant = tenant

	return &tenantAggregate
}

func (a *TenantAggregate) HandleRequest(ctx context.Context, request any) (any, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TenantAggregate.HandleRequest")
	defer span.Finish()

	switch r := request.(type) {
	case *tenantpb.AddBillingProfileRequest:
		return a.AddBillingProfile(ctx, r)
	case *tenantpb.UpdateBillingProfileRequest:
		return r.Id, a.UpdateBillingProfile(ctx, r)
	case *tenantpb.UpdateTenantSettingsRequest:
		return nil, a.UpdateTenantSettings(ctx, r)
	default:
		tracing.TraceErr(span, eventstore.ErrInvalidRequestType)
		return nil, eventstore.ErrInvalidRequestType
	}
}

func (a *TenantAggregate) AddBillingProfile(ctx context.Context, request *tenantpb.AddBillingProfileRequest) (string, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "TenantAggregate.AddBillingProfile")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.GetTenant())
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("AggregateVersion", a.GetVersion()))

	sourceFields := commonmodel.Source{}
	sourceFields.FromGrpc(request.SourceFields)

	createdAtNotNil := utils.IfNotNilTimeWithDefault(utils.TimestampProtoToTimePtr(request.CreatedAt), utils.Now())

	billingProfileId := uuid.New().String()

	addBillingProfileEvent, err := event.NewTenantBillingProfileCreateEvent(a, sourceFields, billingProfileId, request, createdAtNotNil)
	if err != nil {
		tracing.TraceErr(span, err)
		return "", errors.Wrap(err, "TenantBillingProfileCreateEvent")
	}
	aggregate.EnrichEventWithMetadataExtended(&addBillingProfileEvent, span, aggregate.EventMetadata{
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
	aggregate.EnrichEventWithMetadataExtended(&updateBillingProfileEvent, span, aggregate.EventMetadata{
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
	aggregate.EnrichEventWithMetadataExtended(&updateSettingsEvent, span, aggregate.EventMetadata{
		Tenant: r.Tenant,
		UserId: r.LoggedInUserId,
		App:    r.AppSource,
	})

	return a.Apply(updateSettingsEvent)
}

func (a *TenantAggregate) When(evt eventstore.Event) error {
	switch evt.GetEventType() {
	case event.TenantAddBillingProfileV1:
		return a.onAddBillingProfile(evt)
	case event.TenantUpdateBillingProfileV1:
		return a.onUpdateBillingProfile(evt)
	case event.TenantUpdateSettingsV1:
		return a.onUpdateTenantSettings(evt)
	default:
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
		Id:                                eventData.Id,
		CreatedAt:                         eventData.CreatedAt,
		Email:                             eventData.Email,
		Phone:                             eventData.Phone,
		AddressLine1:                      eventData.AddressLine1,
		AddressLine2:                      eventData.AddressLine2,
		AddressLine3:                      eventData.AddressLine3,
		Locality:                          eventData.Locality,
		Country:                           eventData.Country,
		Zip:                               eventData.Zip,
		LegalName:                         eventData.LegalName,
		DomesticPaymentsBankInfo:          eventData.DomesticPaymentsBankInfo,
		DomesticPaymentsBankName:          eventData.DomesticPaymentsBankName,
		DomesticPaymentsAccountNumber:     eventData.DomesticPaymentsAccountNumber,
		DomesticPaymentsSortCode:          eventData.DomesticPaymentsSortCode,
		InternationalPaymentsBankInfo:     eventData.InternationalPaymentsBankInfo,
		InternationalPaymentsSwiftBic:     eventData.InternationalPaymentsSwiftBic,
		InternationalPaymentsBankName:     eventData.InternationalPaymentsBankName,
		InternationalPaymentsBankAddress:  eventData.InternationalPaymentsBankAddress,
		InternationalPaymentsInstructions: eventData.InternationalPaymentsInstructions,
		SourceFields:                      eventData.SourceFields,
	}
	a.TenantDetails.BillingProfiles = append(a.TenantDetails.BillingProfiles, &tenantBillingProfile)

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
		a.TenantDetails.BillingProfiles = append(a.TenantDetails.BillingProfiles, &tenantBillingProfile)
	}

	tenantBillingProfile := a.TenantDetails.GetBillingProfile(eventData.Id)
	if eventData.UpdateEmail() {
		tenantBillingProfile.Email = eventData.Email
	}
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
	if eventData.UpdateZip() {
		tenantBillingProfile.Zip = eventData.Zip
	}
	if eventData.UpdateLegalName() {
		tenantBillingProfile.LegalName = eventData.LegalName
	}
	if eventData.UpdateDomesticPaymentsBankInfo() {
		tenantBillingProfile.DomesticPaymentsBankInfo = eventData.DomesticPaymentsBankInfo
	}
	if eventData.UpdateInternationalPaymentsBankInfo() {
		tenantBillingProfile.InternationalPaymentsBankInfo = eventData.InternationalPaymentsBankInfo
	}
	return nil
}

func (a *TenantAggregate) onUpdateTenantSettings(evt eventstore.Event) error {
	var eventData event.TenantSettingsUpdateEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}

	if eventData.UpdateDefaultCurrency() {
		a.TenantDetails.TenantSettings.DefaultCurrency = eventData.DefaultCurrency
	}
	if eventData.UpdateInvoicingEnabled() {
		a.TenantDetails.TenantSettings.InvoicingEnabled = eventData.InvoicingEnabled
	}
	if eventData.UpdateLogoUrl() {
		a.TenantDetails.TenantSettings.LogoUrl = eventData.LogoUrl
	}
	return nil
}

func extractTenantBillingProfileFieldsMask(requestFieldsMask []tenantpb.TenantBillingProfileFieldMask) []string {
	var fieldsMask []string
	for _, requestFieldMask := range requestFieldsMask {
		switch requestFieldMask {
		case tenantpb.TenantBillingProfileFieldMask_TENANT_BILLING_PROFILE_FIELD_EMAIL:
			fieldsMask = append(fieldsMask, event.FieldMaskEmail)
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
		case tenantpb.TenantBillingProfileFieldMask_TENANT_BILLING_PROFILE_FIELD_ZIP:
			fieldsMask = append(fieldsMask, event.FieldMaskZip)
		case tenantpb.TenantBillingProfileFieldMask_TENANT_BILLING_PROFILE_FIELD_LEGAL_NAME:
			fieldsMask = append(fieldsMask, event.FieldMaskLegalName)
		case tenantpb.TenantBillingProfileFieldMask_TENANT_BILLING_PROFILE_FIELD_DOMESTIC_PAYMENTS_BANK_INFO:
			fieldsMask = append(fieldsMask, event.FieldMaskDomesticPaymentsBankInfo)
		case tenantpb.TenantBillingProfileFieldMask_TENANT_BILLING_PROFILE_FIELD_INTERNATIONAL_PAYMENTS_BANK_INFO:
			fieldsMask = append(fieldsMask, event.FieldMaskInternationalPaymentsBankInfo)
		}
	}
	fieldsMask = utils.RemoveDuplicates(fieldsMask)
	return fieldsMask
}

func extractTenantSettingsFieldsMask(inputFieldsMask []tenantpb.TenantSettingsFieldMask) []string {
	var fieldsMask []string
	for _, requestFieldMask := range inputFieldsMask {
		switch requestFieldMask {
		case tenantpb.TenantSettingsFieldMask_TENANT_SETTINGS_FIELD_LOGO_URL:
			fieldsMask = append(fieldsMask, event.FieldMaskLogoUrl)
		case tenantpb.TenantSettingsFieldMask_TENANT_SETTINGS_FIELD_DEFAULT_CURRENCY:
			fieldsMask = append(fieldsMask, event.FieldMaskDefaultCurrency)
		case tenantpb.TenantSettingsFieldMask_TENANT_SETTINGS_FIELD_INVOICING_ENABLED:
			fieldsMask = append(fieldsMask, event.FieldMaskInvoicingEnabled)
		}
	}
	fieldsMask = utils.RemoveDuplicates(fieldsMask)
	return fieldsMask
}
