package aggregate

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contract/event"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contract/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	contractpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/contract"
	events2 "github.com/openline-ai/openline-customer-os/packages/server/events/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

const (
	ContractAggregateType eventstore.AggregateType = "contract"
)

func GetContractObjectID(aggregateID, tenant string) string {
	return eventstore.GetAggregateObjectID(aggregateID, tenant, ContractAggregateType)
}

type ContractAggregate struct {
	*eventstore.CommonTenantIdAggregate
	Contract *model.Contract
}

func NewContractAggregateWithTenantAndID(tenant, id string) *ContractAggregate {
	contractAggregate := ContractAggregate{}
	contractAggregate.CommonTenantIdAggregate = eventstore.NewCommonAggregateWithTenantAndId(ContractAggregateType, tenant, id)
	contractAggregate.SetWhen(contractAggregate.When)
	contractAggregate.Contract = &model.Contract{}
	contractAggregate.Tenant = tenant

	return &contractAggregate
}

func (a *ContractAggregate) HandleGRPCRequest(ctx context.Context, request any, params map[string]any) (any, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContractAggregate.HandleGRPCRequest")
	defer span.Finish()

	switch r := request.(type) {
	case *contractpb.SoftDeleteContractGrpcRequest:
		return nil, a.softDeleteContract(ctx, r)
	case *contractpb.RolloutRenewalOpportunityOnExpirationGrpcRequest:
		return nil, a.rolloutRenewalOpportunityOnExpiration(ctx, r)
	default:
		tracing.TraceErr(span, eventstore.ErrInvalidRequestType)
		return nil, eventstore.ErrInvalidRequestType
	}
}

func (a *ContractAggregate) rolloutRenewalOpportunityOnExpiration(ctx context.Context, request *contractpb.RolloutRenewalOpportunityOnExpirationGrpcRequest) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "ContractAggregate.rolloutRenewalOpportunityOnExpiration")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.Tenant)
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("aggregateVersion", a.GetVersion()))
	tracing.LogObjectAsJson(span, "request", request)

	updateEvent, err := event.NewRolloutRenewalOpportunityEvent(a)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewRolloutRenewalOpportunityEvent")
	}
	eventstore.EnrichEventWithMetadataExtended(&updateEvent, span, eventstore.EventMetadata{
		Tenant: a.Tenant,
		UserId: request.LoggedInUserId,
		App:    request.GetAppSource(),
	})

	return a.Apply(updateEvent)
}

func isUpdated(field string, fieldsMask []string) bool {
	return len(fieldsMask) == 0 || utils.Contains(fieldsMask, field)
}

func (a *ContractAggregate) softDeleteContract(ctx context.Context, r *contractpb.SoftDeleteContractGrpcRequest) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "ContractAggregate.softDeleteContract")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.Tenant)
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("aggregateVersion", a.GetVersion()))
	tracing.LogObjectAsJson(span, "request", r)

	deleteEvent, err := event.NewContractDeleteEvent(a, utils.Now())
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewContractDeleteEvent")
	}
	eventstore.EnrichEventWithMetadataExtended(&deleteEvent, span, eventstore.EventMetadata{
		Tenant: a.Tenant,
		UserId: r.LoggedInUserId,
		App:    r.AppSource,
	})

	return a.Apply(deleteEvent)
}

func (a *ContractAggregate) When(evt eventstore.Event) error {
	switch evt.GetEventType() {
	case event.ContractUpdateStatusV1:
		return a.onContractRefreshStatus(evt)
	case event.ContractRolloutRenewalOpportunityV1:
		return nil
	case event.ContractDeleteV1:
		return a.onContractDelete(evt)
	default:
		return nil
	}
}

func (a *ContractAggregate) onContractUpdate(evt eventstore.Event) error {
	var eventData event.ContractUpdateEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}

	// Update only if the source of truth is 'openline' or the new source matches the source of truth
	if eventData.Source == events2.SourceOpenline {
		a.Contract.Source.SourceOfTruth = eventData.Source
	}

	if eventData.Source != a.Contract.Source.SourceOfTruth && a.Contract.Source.SourceOfTruth == events2.SourceOpenline {
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
	if eventData.UpdateLengthInMonths() {
		a.Contract.LengthInMonths = eventData.LengthInMonths
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
	if eventData.UpdateBillingCycleInMonths() {
		a.Contract.BillingCycleInMonths = eventData.BillingCycleInMonths
	} else if eventData.UpdateBillingCycle() {
		switch eventData.BillingCycle {
		case "MONTHLY":
			a.Contract.BillingCycleInMonths = 1
		case "QUARTERLY":
			a.Contract.BillingCycleInMonths = 3
		case "ANNUALLY":
			a.Contract.BillingCycleInMonths = 12
		case "":
			a.Contract.BillingCycleInMonths = 0
		}
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
	if eventData.UpdateCountry() {
		a.Contract.Region = eventData.Region
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
	if eventData.UpdateInvoiceEmailCC() {
		a.Contract.InvoiceEmailCC = eventData.InvoiceEmailCC
	}
	if eventData.UpdateInvoiceEmailBCC() {
		a.Contract.InvoiceEmailBCC = eventData.InvoiceEmailBCC
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
	if eventData.UpdateAutoRenew() {
		a.Contract.AutoRenew = eventData.AutoRenew
	}
	if eventData.UpdateCheck() {
		a.Contract.Check = eventData.Check
	}
	if eventData.UpdateDueDays() {
		a.Contract.DueDays = eventData.DueDays
	}
	if eventData.UpdateApproved() {
		a.Contract.Approved = eventData.Approved
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
