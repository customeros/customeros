package event

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/validator"
	tenantpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/tenant"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/pkg/errors"
	"time"
)

type TenantSettingsUpdateEvent struct {
	Tenant               string    `json:"tenant" validate:"required"`
	UpdatedAt            time.Time `json:"updatedAt"`
	BaseCurrency         string    `json:"baseCurrency,omitempty"`
	InvoicingEnabled     bool      `json:"invoicingEnabled,omitempty"`
	InvoicingPostpaid    bool      `json:"invoicingPostpaid,omitempty"`
	LogoRepositoryFileId string    `json:"logoRepositoryFileId,omitempty"`
	FieldsMask           []string  `json:"fieldsMask,omitempty"`
}

func NewTenantSettingsUpdateEvent(aggregate eventstore.Aggregate, request *tenantpb.UpdateTenantSettingsRequest, updatedAt time.Time, fieldsMaks []string) (eventstore.Event, error) {
	eventData := TenantSettingsUpdateEvent{
		Tenant:               aggregate.GetTenant(),
		UpdatedAt:            updatedAt,
		BaseCurrency:         request.BaseCurrency,
		InvoicingEnabled:     request.InvoicingEnabled,
		InvoicingPostpaid:    request.InvoicingPostpaid,
		LogoRepositoryFileId: request.LogoRepositoryFileId,
		FieldsMask:           fieldsMaks,
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

func (e TenantSettingsUpdateEvent) UpdateInvoicingEnabled() bool {
	return len(e.FieldsMask) == 0 || utils.Contains(e.FieldsMask, FieldMaskInvoicingEnabled)
}

func (e TenantSettingsUpdateEvent) UpdateInvoicingPostpaid() bool {
	return len(e.FieldsMask) == 0 || utils.Contains(e.FieldsMask, FieldMaskInvoicingPostpaid)
}

func (e TenantSettingsUpdateEvent) UpdateLogoRepositoryFileId() bool {
	return len(e.FieldsMask) == 0 || utils.Contains(e.FieldsMask, FieldMaskLogoRepositoryFileId)
}

func (e TenantSettingsUpdateEvent) UpdateBaseCurrency() bool {
	return len(e.FieldsMask) == 0 || utils.Contains(e.FieldsMask, FieldMaskBaseCurrency)
}
