package event

import (
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/validator"
	tenantpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/tenant"
	"github.com/pkg/errors"
	"time"
)

type CreateTenantBillingProfileEvent struct {
	Tenant                            string             `json:"tenant" validate:"required"`
	Id                                string             `json:"id" validate:"required"`
	CreatedAt                         time.Time          `json:"createdAt"`
	SourceFields                      commonmodel.Source `json:"sourceFields"`
	Email                             string             `json:"email"`
	Phone                             string             `json:"phone"`
	AddressLine1                      string             `json:"addressLine1"`
	AddressLine2                      string             `json:"addressLine2"`
	AddressLine3                      string             `json:"addressLine3"`
	LegalName                         string             `json:"legalName"`
	DomesticPaymentsBankName          string             `json:"domesticPaymentsBankName"`
	DomesticPaymentsAccountNumber     string             `json:"domesticPaymentsAccountNumber"`
	DomesticPaymentsSortCode          string             `json:"domesticPaymentsSortCode"`
	InternationalPaymentsSwiftBic     string             `json:"internationalPaymentsSwiftBic"`
	InternationalPaymentsBankName     string             `json:"internationalPaymentsBankName"`
	InternationalPaymentsBankAddress  string             `json:"internationalPaymentsBankAddress"`
	InternationalPaymentsInstructions string             `json:"internationalPaymentsInstructions"`
}

func NewCreateTenantBillingProfileEvent(aggregate eventstore.Aggregate, sourceFields commonmodel.Source, id string, request *tenantpb.AddBillingProfileRequest, createdAt time.Time) (eventstore.Event, error) {
	eventData := CreateTenantBillingProfileEvent{
		Tenant:                            aggregate.GetTenant(),
		Id:                                id,
		CreatedAt:                         createdAt,
		SourceFields:                      sourceFields,
		Email:                             request.Email,
		Phone:                             request.Phone,
		AddressLine1:                      request.AddressLine1,
		AddressLine2:                      request.AddressLine2,
		AddressLine3:                      request.AddressLine3,
		LegalName:                         request.LegalName,
		DomesticPaymentsBankName:          request.DomesticPaymentsBankName,
		DomesticPaymentsAccountNumber:     request.DomesticPaymentsAccountNumber,
		DomesticPaymentsSortCode:          request.DomesticPaymentsSortCode,
		InternationalPaymentsSwiftBic:     request.InternationalPaymentsSwiftBic,
		InternationalPaymentsBankName:     request.InternationalPaymentsBankName,
		InternationalPaymentsBankAddress:  request.InternationalPaymentsBankAddress,
		InternationalPaymentsInstructions: request.InternationalPaymentsInstructions,
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "failed to validate CreateTenantBillingProfileEvent")
	}

	event := eventstore.NewBaseEvent(aggregate, TenantAddBillingProfileV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for CreateTenantBillingProfileEvent")
	}

	return event, nil
}
