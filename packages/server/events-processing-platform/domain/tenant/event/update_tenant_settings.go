package event

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/validator"
	tenantpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/tenant"
	"github.com/pkg/errors"
	"time"
)

type TenantSettingsUpdateEvent struct {
	Tenant            string    `json:"tenant" validate:"required"`
	UpdatedAt         time.Time `json:"updatedAt"`
	DefaultCurrency   string    `json:"defaultCurrency,omitempty"`
	InvoicingEnabled  bool      `json:"invoicingEnabled,omitempty"`
	InvoicingPostpaid bool      `json:"invoicingPostpaid,omitempty"`
	LogoUrl           string    `json:"logoUrl,omitempty"`
	FieldsMask        []string  `json:"fieldsMask,omitempty"`
}

func NewTenantSettingsUpdateEvent(aggregate eventstore.Aggregate, request *tenantpb.UpdateTenantSettingsRequest, updatedAt time.Time, fieldsMaks []string) (eventstore.Event, error) {
	eventData := TenantSettingsUpdateEvent{
		Tenant:            aggregate.GetTenant(),
		UpdatedAt:         updatedAt,
		DefaultCurrency:   request.DefaultCurrency,
		InvoicingEnabled:  request.InvoicingEnabled,
		InvoicingPostpaid: request.InvoicingPostpaid,
		LogoUrl:           request.LogoUrl,
		FieldsMask:        fieldsMaks,
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "failed to validate TenantSettingsUpdateEvent")
	}

	event := eventstore.NewBaseEvent(aggregate, TenantUpdateSettingsV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for TenantSettingsUpdateEvent")
	}

	return event, nil
}

func (e TenantSettingsUpdateEvent) UpdateLogoUrl() bool {
	return len(e.FieldsMask) == 0 || utils.Contains(e.FieldsMask, FieldMaskLogoUrl)
}

func (e TenantSettingsUpdateEvent) UpdateDefaultCurrency() bool {
	return len(e.FieldsMask) == 0 || utils.Contains(e.FieldsMask, FieldMaskDefaultCurrency)
}

func (e TenantSettingsUpdateEvent) UpdateInvoicingEnabled() bool {
	return len(e.FieldsMask) == 0 || utils.Contains(e.FieldsMask, FieldMaskInvoicingEnabled)
}

func (e TenantSettingsUpdateEvent) UpdateInvoicingPostpaid() bool {
	return len(e.FieldsMask) == 0 || utils.Contains(e.FieldsMask, FieldMaskInvoicingPostpaid)
}
