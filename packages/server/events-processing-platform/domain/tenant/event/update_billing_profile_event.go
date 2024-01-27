package event

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/validator"
	tenantpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/tenant"
	"github.com/pkg/errors"
	"time"
)

type TenantBillingProfileUpdateEvent struct {
	Tenant                        string    `json:"tenant" validate:"required"`
	Id                            string    `json:"id" validate:"required"`
	UpdatedAt                     time.Time `json:"updatedAt"`
	Email                         string    `json:"email,omitempty"`
	Phone                         string    `json:"phone,omitempty"`
	AddressLine1                  string    `json:"addressLine1,omitempty"`
	AddressLine2                  string    `json:"addressLine2,omitempty"`
	AddressLine3                  string    `json:"addressLine3,omitempty"`
	Locality                      string    `json:"locality,omitempty"`
	Country                       string    `json:"country,omitempty"`
	Zip                           string    `json:"zip,omitempty"`
	LegalName                     string    `json:"legalName,omitempty"`
	DomesticPaymentsBankInfo      string    `json:"domesticPaymentsBankInfo,omitempty"`
	InternationalPaymentsBankInfo string    `json:"internationalPaymentsBankInfo,omitempty"`
	FieldsMask                    []string  `json:"fieldsMask,omitempty"`
}

func NewTenantBillingProfileUpdateEvent(aggregate eventstore.Aggregate, id string, request *tenantpb.UpdateBillingProfileRequest, updatedAt time.Time, fieldsMaks []string) (eventstore.Event, error) {
	eventData := TenantBillingProfileUpdateEvent{
		Tenant:                        aggregate.GetTenant(),
		Id:                            id,
		UpdatedAt:                     updatedAt,
		Email:                         request.Email,
		Phone:                         request.Phone,
		AddressLine1:                  request.AddressLine1,
		AddressLine2:                  request.AddressLine2,
		AddressLine3:                  request.AddressLine3,
		Locality:                      request.Locality,
		Country:                       request.Country,
		Zip:                           request.Zip,
		LegalName:                     request.LegalName,
		DomesticPaymentsBankInfo:      request.DomesticPaymentsBankInfo,
		InternationalPaymentsBankInfo: request.InternationalPaymentsBankInfo,
		FieldsMask:                    fieldsMaks,
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "failed to validate TenantBillingProfileUpdateEvent")
	}

	event := eventstore.NewBaseEvent(aggregate, TenantUpdateBillingProfileV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for TenantBillingProfileUpdateEvent")
	}

	return event, nil
}

func (e TenantBillingProfileUpdateEvent) UpdateEmail() bool {
	return len(e.FieldsMask) == 0 || utils.Contains(e.FieldsMask, FieldMaskEmail)
}

func (e TenantBillingProfileUpdateEvent) UpdatePhone() bool {
	return len(e.FieldsMask) == 0 || utils.Contains(e.FieldsMask, FieldMaskPhone)
}

func (e TenantBillingProfileUpdateEvent) UpdateAddressLine1() bool {
	return len(e.FieldsMask) == 0 || utils.Contains(e.FieldsMask, FieldMaskAddressLine1)
}

func (e TenantBillingProfileUpdateEvent) UpdateAddressLine2() bool {
	return len(e.FieldsMask) == 0 || utils.Contains(e.FieldsMask, FieldMaskAddressLine2)
}

func (e TenantBillingProfileUpdateEvent) UpdateAddressLine3() bool {
	return len(e.FieldsMask) == 0 || utils.Contains(e.FieldsMask, FieldMaskAddressLine3)
}

func (e TenantBillingProfileUpdateEvent) UpdateLocality() bool {
	return len(e.FieldsMask) == 0 || utils.Contains(e.FieldsMask, FieldMaskLocality)
}

func (e TenantBillingProfileUpdateEvent) UpdateCountry() bool {
	return len(e.FieldsMask) == 0 || utils.Contains(e.FieldsMask, FieldMaskCountry)
}

func (e TenantBillingProfileUpdateEvent) UpdateZip() bool {
	return len(e.FieldsMask) == 0 || utils.Contains(e.FieldsMask, FieldMaskZip)
}

func (e TenantBillingProfileUpdateEvent) UpdateLegalName() bool {
	return len(e.FieldsMask) == 0 || utils.Contains(e.FieldsMask, FieldMaskLegalName)
}

func (e TenantBillingProfileUpdateEvent) UpdateDomesticPaymentsBankInfo() bool {
	return len(e.FieldsMask) == 0 || utils.Contains(e.FieldsMask, FieldMaskDomesticPaymentsBankInfo)
}

func (e TenantBillingProfileUpdateEvent) UpdateInternationalPaymentsBankInfo() bool {
	return len(e.FieldsMask) == 0 || utils.Contains(e.FieldsMask, FieldMaskInternationalPaymentsBankInfo)
}
