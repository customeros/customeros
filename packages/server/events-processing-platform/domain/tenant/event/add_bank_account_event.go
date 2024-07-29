package event

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/validator"
	tenantpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/tenant"
	"github.com/openline-ai/openline-customer-os/packages/server/events/event/common"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/pkg/errors"
	"time"
)

type TenantBankAccountCreateEvent struct {
	Tenant              string        `json:"tenant" validate:"required"`
	Id                  string        `json:"id" validate:"required"`
	CreatedAt           time.Time     `json:"createdAt"`
	SourceFields        common.Source `json:"sourceFields"`
	BankName            string        `json:"bankName,omitempty"`
	BankTransferEnabled bool          `json:"bankTransferEnabled"`
	AllowInternational  bool          `json:"allowInternational"`
	Currency            string        `json:"currency"`
	Iban                string        `json:"iban,omitempty"`
	Bic                 string        `json:"bic,omitempty"`
	SortCode            string        `json:"sortCode,omitempty"`
	AccountNumber       string        `json:"accountNumber,omitempty"`
	RoutingNumber       string        `json:"routingNumber,omitempty"`
	OtherDetails        string        `json:"otherDetails,omitempty"`
}

func NewTenantBankAccountCreateEvent(aggregate eventstore.Aggregate, sourceFields common.Source, id string, request *tenantpb.AddBankAccountGrpcRequest, createdAt time.Time) (eventstore.Event, error) {
	eventData := TenantBankAccountCreateEvent{
		Tenant:              aggregate.GetTenant(),
		Id:                  id,
		CreatedAt:           createdAt,
		SourceFields:        sourceFields,
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
		return eventstore.Event{}, errors.Wrap(err, "failed to validate TenantBankAccountCreateEvent")
	}

	event := eventstore.NewBaseEvent(aggregate, TenantAddBankAccountV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for TenantBankAccountCreateEvent")
	}

	return event, nil
}
