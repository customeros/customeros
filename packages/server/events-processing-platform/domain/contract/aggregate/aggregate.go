package aggregate

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/aggregate"
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contract/event"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contract/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/pkg/errors"
	"strings"
)

const (
	ContractAggregateType eventstore.AggregateType = "contract"
)

type ContractAggregate struct {
	*aggregate.CommonTenantIdAggregate
	Contract *model.Contract
}

func NewContractAggregateWithTenantAndID(tenant, id string) *ContractAggregate {
	contractAggregate := ContractAggregate{}
	contractAggregate.CommonTenantIdAggregate = aggregate.NewCommonAggregateWithTenantAndId(ContractAggregateType, tenant, id)
	contractAggregate.SetWhen(contractAggregate.When)
	contractAggregate.Contract = &model.Contract{}
	contractAggregate.Tenant = tenant

	return &contractAggregate
}

type ContractTempAggregate struct {
	*aggregate.CommonTenantIdTempAggregate
}

func NewContractTempAggregateWithTenantAndID(tenant, id string) *ContractTempAggregate {
	contractTempAggregate := ContractTempAggregate{}
	contractTempAggregate.CommonTenantIdTempAggregate = aggregate.NewCommonTempAggregateWithTenantAndId(ContractAggregateType, tenant, id)
	contractTempAggregate.Tenant = tenant

	return &contractTempAggregate
}

func (a *ContractAggregate) When(evt eventstore.Event) error {
	switch evt.GetEventType() {
	case event.ContractCreateV1:
		return a.onContractCreate(evt)
	case event.ContractUpdateV1:
		return a.onContractUpdate(evt)
	case event.ContractUpdateStatusV1:
		return a.onContractRefreshStatus(evt)
	case event.ContractRolloutRenewalOpportunityV1:
		return nil
	case event.ContractDeleteV1:
		return a.onContractDelete(evt)
	default:
		if strings.HasPrefix(evt.GetEventType(), constants.EsInternalStreamPrefix) {
			return nil
		}
		err := eventstore.ErrInvalidEventType
		err.EventType = evt.GetEventType()
		return err
	}
}

func (a *ContractAggregate) onContractCreate(evt eventstore.Event) error {
	var eventData event.ContractCreateEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}

	a.Contract.ID = a.ID
	a.Contract.Tenant = a.Tenant
	a.Contract.OrganizationId = eventData.OrganizationId
	a.Contract.Name = eventData.Name
	a.Contract.ContractUrl = eventData.ContractUrl
	a.Contract.CreatedByUserId = eventData.CreatedByUserId
	a.Contract.ServiceStartedAt = eventData.ServiceStartedAt
	a.Contract.SignedAt = eventData.SignedAt
	a.Contract.RenewalCycle = eventData.RenewalCycle
	a.Contract.Status = eventData.Status
	a.Contract.Currency = eventData.Currency
	a.Contract.BillingCycle = eventData.BillingCycle
	a.Contract.InvoicingStartDate = eventData.InvoicingStartDate
	a.Contract.CreatedAt = eventData.CreatedAt
	a.Contract.UpdatedAt = eventData.UpdatedAt
	a.Contract.Source = eventData.Source
	a.Contract.InvoicingEnabled = eventData.InvoicingEnabled
	if eventData.ExternalSystem.Available() {
		a.Contract.ExternalSystems = []commonmodel.ExternalSystem{eventData.ExternalSystem}
	}
	a.Contract.PayOnline = eventData.PayOnline
	a.Contract.PayAutomatically = eventData.PayAutomatically
	a.Contract.CanPayWithCard = eventData.CanPayWithCard
	a.Contract.CanPayWithDirectDebit = eventData.CanPayWithDirectDebit
	a.Contract.CanPayWithBankTransfer = eventData.CanPayWithBankTransfer
	return nil
}

func (a *ContractAggregate) onContractUpdate(evt eventstore.Event) error {
	var eventData event.ContractUpdateEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}

	// Update only if the source of truth is 'openline' or the new source matches the source of truth
	if eventData.Source == constants.SourceOpenline {
		a.Contract.Source.SourceOfTruth = eventData.Source
	}

	if eventData.Source != a.Contract.Source.SourceOfTruth && a.Contract.Source.SourceOfTruth == constants.SourceOpenline {
		// Update fields only if they are empty
		if a.Contract.Name == "" && eventData.UpdateName() {
			a.Contract.Name = eventData.Name
		}
		if a.Contract.ContractUrl == "" && eventData.UpdateContractUrl() {
			a.Contract.ContractUrl = eventData.ContractUrl
		}
	} else {
		// Update fields unconditionally
		if eventData.UpdateName() {
			a.Contract.Name = eventData.Name
		}
		if eventData.UpdateContractUrl() {
			a.Contract.ContractUrl = eventData.ContractUrl
		}
	}

	a.Contract.UpdatedAt = eventData.UpdatedAt
	if eventData.UpdateRenewalCycle() {
		a.Contract.RenewalCycle = eventData.RenewalCycle
	}
	if eventData.UpdateStatus() {
		a.Contract.Status = eventData.Status
	}
	if eventData.UpdateServiceStartedAt() {
		a.Contract.ServiceStartedAt = eventData.ServiceStartedAt
	}
	if eventData.UpdateSignedAt() {
		a.Contract.SignedAt = eventData.SignedAt
	}
	if eventData.UpdateEndedAt() {
		a.Contract.EndedAt = eventData.EndedAt
	}
	if eventData.UpdateCurrency() {
		a.Contract.Currency = eventData.Currency
	}
	if eventData.UpdateBillingCycle() {
		a.Contract.BillingCycle = eventData.BillingCycle
	}
	if eventData.UpdateInvoicingStartDate() {
		a.Contract.InvoicingStartDate = eventData.InvoicingStartDate
	}
	if eventData.UpdateAddressLine1() {
		a.Contract.AddressLine1 = eventData.AddressLine1
	}
	if eventData.UpdateAddressLine2() {
		a.Contract.AddressLine2 = eventData.AddressLine2
	}
	if eventData.UpdateLocality() {
		a.Contract.Locality = eventData.Locality
	}
	if eventData.UpdateCountry() {
		a.Contract.Country = eventData.Country
	}
	if eventData.UpdateZip() {
		a.Contract.Zip = eventData.Zip
	}
	if eventData.UpdateOrganizationLegalName() {
		a.Contract.OrganizationLegalName = eventData.OrganizationLegalName
	}
	if eventData.UpdateInvoiceEmail() {
		a.Contract.InvoiceEmail = eventData.InvoiceEmail
	}
	if eventData.UpdateInvoiceNote() {
		a.Contract.InvoiceNote = eventData.InvoiceNote
	}
	if eventData.UpdateNextInvoiceDate() {
		a.Contract.NextInvoiceDate = eventData.NextInvoiceDate
	}
	if eventData.UpdateCanPayWithCard() {
		a.Contract.CanPayWithCard = eventData.CanPayWithCard
	}
	if eventData.UpdateCanPayWithDirectDebit() {
		a.Contract.CanPayWithDirectDebit = eventData.CanPayWithDirectDebit
	}
	if eventData.UpdateCanPayWithBankTransfer() {
		a.Contract.CanPayWithBankTransfer = eventData.CanPayWithBankTransfer
	}
	if eventData.UpdateInvoicingEnabled() {
		a.Contract.InvoicingEnabled = eventData.InvoicingEnabled
	}
	if eventData.UpdatePayOnline() {
		a.Contract.PayOnline = eventData.PayOnline
	}
	if eventData.UpdatePayAutomatically() {
		a.Contract.PayAutomatically = eventData.PayAutomatically
	}

	if eventData.ExternalSystem.Available() {
		found := false
		for _, externalSystem := range a.Contract.ExternalSystems {
			if externalSystem.ExternalSystemId == eventData.ExternalSystem.ExternalSystemId && externalSystem.ExternalId == eventData.ExternalSystem.ExternalId {
				found = true
				externalSystem.ExternalUrl = eventData.ExternalSystem.ExternalUrl
				externalSystem.ExternalSource = eventData.ExternalSystem.ExternalSource
				externalSystem.SyncDate = eventData.ExternalSystem.SyncDate
				if eventData.ExternalSystem.ExternalIdSecond != "" {
					externalSystem.ExternalIdSecond = eventData.ExternalSystem.ExternalIdSecond
				}
			}
		}
		if !found {
			a.Contract.ExternalSystems = append(a.Contract.ExternalSystems, eventData.ExternalSystem)
		}
	}

	return nil
}

func (a *ContractAggregate) onContractRefreshStatus(evt eventstore.Event) error {
	var eventData event.ContractUpdateStatusEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}

	a.Contract.Status = eventData.Status
	return nil
}

func (a *ContractAggregate) onContractDelete(evt eventstore.Event) error {
	var eventData event.ContractDeleteEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}

	a.Contract.Removed = true
	return nil
}
