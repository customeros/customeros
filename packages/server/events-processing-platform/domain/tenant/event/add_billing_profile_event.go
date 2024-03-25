package event

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/validator"
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	tenantpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/tenant"
	"github.com/pkg/errors"
	"time"
)

type TenantBillingProfileCreateEvent struct {
	Tenant                 string             `json:"tenant" validate:"required"`
	Id                     string             `json:"id" validate:"required"`
	CreatedAt              time.Time          `json:"createdAt"`
	SourceFields           commonmodel.Source `json:"sourceFields"`
	Phone                  string             `json:"phone"`
	AddressLine1           string             `json:"addressLine1"`
	AddressLine2           string             `json:"addressLine2"`
	AddressLine3           string             `json:"addressLine3"`
	Locality               string             `json:"locality"`
	Country                string             `json:"country"`
	Region                 string             `json:"region"`
	Zip                    string             `json:"zip"`
	LegalName              string             `json:"legalName"`
	VatNumber              string             `json:"vatNumber"`
	SendInvoicesFrom       string             `json:"sendInvoicesFrom"`
	SendInvoicesBcc        string             `json:"sendInvoicesBcc"`
	CanPayWithPigeon       bool               `json:"canPayWithPigeon"`
	CanPayWithBankTransfer bool               `json:"canPayWithBankTransfer"`
	InvoiceNote            string             `json:"invoiceNote"`
}

func NewTenantBillingProfileCreateEvent(aggregate eventstore.Aggregate, sourceFields commonmodel.Source, id string, request *tenantpb.AddBillingProfileRequest, createdAt time.Time) (eventstore.Event, error) {
	eventData := TenantBillingProfileCreateEvent{
		Tenant:                 aggregate.GetTenant(),
		Id:                     id,
		CreatedAt:              createdAt,
		SourceFields:           sourceFields,
		Phone:                  request.Phone,
		AddressLine1:           request.AddressLine1,
		AddressLine2:           request.AddressLine2,
		AddressLine3:           request.AddressLine3,
		Locality:               request.Locality,
		Country:                request.Country,
		Region:                 request.Region,
		Zip:                    request.Zip,
		LegalName:              request.LegalName,
		VatNumber:              request.VatNumber,
		SendInvoicesFrom:       request.SendInvoicesFrom,
		SendInvoicesBcc:        request.SendInvoicesBcc,
		CanPayWithPigeon:       request.CanPayWithPigeon,
		CanPayWithBankTransfer: request.CanPayWithBankTransfer,
		InvoiceNote:            request.InvoiceNote,
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "failed to validate TenantBillingProfileCreateEvent")
	}

	event := eventstore.NewBaseEvent(aggregate, TenantAddBillingProfileV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for TenantBillingProfileCreateEvent")
	}

	return event, nil
}
