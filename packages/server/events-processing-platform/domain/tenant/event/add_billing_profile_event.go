package event

import (
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/validator"
	tenantpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/tenant"
	"github.com/pkg/errors"
	"time"
)

type TenantBillingProfileCreateEvent struct {
	Tenant                            string             `json:"tenant" validate:"required"`
	Id                                string             `json:"id" validate:"required"`
	CreatedAt                         time.Time          `json:"createdAt"`
	SourceFields                      commonmodel.Source `json:"sourceFields"`
	Phone                             string             `json:"phone"`
	AddressLine1                      string             `json:"addressLine1"`
	AddressLine2                      string             `json:"addressLine2"`
	AddressLine3                      string             `json:"addressLine3"`
	Locality                          string             `json:"locality"`
	Country                           string             `json:"country"`
	Zip                               string             `json:"zip"`
	LegalName                         string             `json:"legalName"`
	DomesticPaymentsBankInfo          string             `json:"domesticPaymentsBankInfo"`
	DomesticPaymentsBankName          string             `json:"domesticPaymentsBankName"`
	DomesticPaymentsAccountNumber     string             `json:"domesticPaymentsAccountNumber"`
	DomesticPaymentsSortCode          string             `json:"domesticPaymentsSortCode"`
	InternationalPaymentsBankInfo     string             `json:"internationalPaymentsBankInfo"`
	InternationalPaymentsSwiftBic     string             `json:"internationalPaymentsSwiftBic"`
	InternationalPaymentsBankName     string             `json:"internationalPaymentsBankName"`
	InternationalPaymentsBankAddress  string             `json:"internationalPaymentsBankAddress"`
	InternationalPaymentsInstructions string             `json:"internationalPaymentsInstructions"`
	VatNumber                         string             `json:"vatNumber"`
	SendInvoicesFrom                  string             `json:"sendInvoicesFrom"`
	SendInvoicesBcc                   string             `json:"sendInvoicesBcc"`
	CanPayWithCard                    bool               `json:"canPayWithCard"`
	CanPayWithDirectDebitSEPA         bool               `json:"canPayWithDirectDebitSEPA"`
	CanPayWithDirectDebitACH          bool               `json:"canPayWithDirectDebitACH"`
	CanPayWithDirectDebitBacs         bool               `json:"canPayWithDirectDebitBacs"`
	CanPayWithPigeon                  bool               `json:"canPayWithPigeon"`
	CanPayWithBankTransfer            bool               `json:"canPayWithBankTransfer"`
}

func NewTenantBillingProfileCreateEvent(aggregate eventstore.Aggregate, sourceFields commonmodel.Source, id string, request *tenantpb.AddBillingProfileRequest, createdAt time.Time) (eventstore.Event, error) {
	eventData := TenantBillingProfileCreateEvent{
		Tenant:                            aggregate.GetTenant(),
		Id:                                id,
		CreatedAt:                         createdAt,
		SourceFields:                      sourceFields,
		Phone:                             request.Phone,
		AddressLine1:                      request.AddressLine1,
		AddressLine2:                      request.AddressLine2,
		AddressLine3:                      request.AddressLine3,
		Locality:                          request.Locality,
		Country:                           request.Country,
		Zip:                               request.Zip,
		LegalName:                         request.LegalName,
		DomesticPaymentsBankInfo:          request.DomesticPaymentsBankInfo,
		DomesticPaymentsBankName:          request.DomesticPaymentsBankName,
		DomesticPaymentsAccountNumber:     request.DomesticPaymentsAccountNumber,
		DomesticPaymentsSortCode:          request.DomesticPaymentsSortCode,
		InternationalPaymentsBankInfo:     request.InternationalPaymentsBankInfo,
		InternationalPaymentsSwiftBic:     request.InternationalPaymentsSwiftBic,
		InternationalPaymentsBankName:     request.InternationalPaymentsBankName,
		InternationalPaymentsBankAddress:  request.InternationalPaymentsBankAddress,
		InternationalPaymentsInstructions: request.InternationalPaymentsInstructions,
		VatNumber:                         request.VatNumber,
		SendInvoicesFrom:                  request.SendInvoicesFrom,
		SendInvoicesBcc:                   request.SendInvoicesBcc,
		CanPayWithCard:                    request.CanPayWithCard,
		CanPayWithDirectDebitSEPA:         request.CanPayWithDirectDebitSEPA,
		CanPayWithDirectDebitACH:          request.CanPayWithDirectDebitACH,
		CanPayWithDirectDebitBacs:         request.CanPayWithDirectDebitBacs,
		CanPayWithPigeon:                  request.CanPayWithPigeon,
		CanPayWithBankTransfer:            request.CanPayWithBankTransfer,
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
