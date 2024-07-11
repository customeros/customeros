package event

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/validator"
	tenantpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/tenant"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/pkg/errors"
	"time"
)

type TenantBankAccountUpdateEvent struct {
	Tenant              string    `json:"tenant" validate:"required"`
	Id                  string    `json:"id" validate:"required"`
	UpdatedAt           time.Time `json:"updatedAt"`
	FieldsMask          []string  `json:"fieldsMask,omitempty"`
	BankName            string    `json:"bankName,omitempty"`
	BankTransferEnabled bool      `json:"bankTransferEnabled,omitempty"`
	AllowInternational  bool      `json:"allowInternational,omitempty"`
	Currency            string    `json:"currency,omitempty"`
	Iban                string    `json:"iban,omitempty"`
	Bic                 string    `json:"bic,omitempty"`
	SortCode            string    `json:"sortCode,omitempty"`
	AccountNumber       string    `json:"accountNumber,omitempty"`
	RoutingNumber       string    `json:"routingNumber,omitempty"`
	OtherDetails        string    `json:"otherDetails,omitempty"`
}

func NewTenantBankAccountUpdateEvent(aggregate eventstore.Aggregate, id string, request *tenantpb.UpdateBankAccountGrpcRequest, updatedAt time.Time, fieldsMaks []string) (eventstore.Event, error) {
	eventData := TenantBankAccountUpdateEvent{
		Tenant:              aggregate.GetTenant(),
		Id:                  id,
		UpdatedAt:           updatedAt,
		FieldsMask:          fieldsMaks,
		BankName:            request.BankName,
		BankTransferEnabled: request.BankTransferEnabled,
		AllowInternational:  request.AllowInternational,
		Currency:            request.Currency,
		Iban:                request.Iban,
		Bic:                 request.Bic,
		SortCode:            request.SortCode,
		AccountNumber:       request.AccountNumber,
		RoutingNumber:       request.RoutingNumber,
		OtherDetails:        request.OtherDetails,
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "failed to validate TenantBankAccountUpdateEvent")
	}

	event := eventstore.NewBaseEvent(aggregate, TenantUpdateBankAccountV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for TenantBankAccountUpdateEvent")
	}

	return event, nil
}

func (e TenantBankAccountUpdateEvent) UpdateBankName() bool {
	return utils.Contains(e.FieldsMask, FieldMaskBankAccountBankName)
}

func (e TenantBankAccountUpdateEvent) UpdateBankTransferEnabled() bool {
	return utils.Contains(e.FieldsMask, FieldMaskBankAccountBankTransferEnabled)
}

func (e TenantBankAccountUpdateEvent) UpdateAllowInternational() bool {
	return utils.Contains(e.FieldsMask, FieldMaskBankAccountAllowInternational)
}

func (e TenantBankAccountUpdateEvent) UpdateCurrency() bool {
	return utils.Contains(e.FieldsMask, FieldMaskBankAccountCurrency)
}

func (e TenantBankAccountUpdateEvent) UpdateIban() bool {
	return utils.Contains(e.FieldsMask, FieldMaskBankAccountIban)
}

func (e TenantBankAccountUpdateEvent) UpdateBic() bool {
	return utils.Contains(e.FieldsMask, FieldMaskBankAccountBic)
}

func (e TenantBankAccountUpdateEvent) UpdateSortCode() bool {
	return utils.Contains(e.FieldsMask, FieldMaskBankAccountSortCode)
}

func (e TenantBankAccountUpdateEvent) UpdateAccountNumber() bool {
	return utils.Contains(e.FieldsMask, FieldMaskBankAccountAccountNumber)
}

func (e TenantBankAccountUpdateEvent) UpdateRoutingNumber() bool {
	return utils.Contains(e.FieldsMask, FieldMaskBankAccountRoutingNumber)
}

func (e TenantBankAccountUpdateEvent) UpdateOtherDetails() bool {
	return utils.Contains(e.FieldsMask, FieldMaskBankAccountOtherDetails)
}
