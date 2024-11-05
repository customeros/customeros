package aggregate

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contract/event"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contract/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	contractpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/contract"
	events2 "github.com/openline-ai/openline-customer-os/packages/server/events/constants"
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/events/event/common"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"time"
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
	case *contractpb.UpdateContractGrpcRequest:
		return nil, a.updateContract(ctx, r)
	case *contractpb.SoftDeleteContractGrpcRequest:
		return nil, a.softDeleteContract(ctx, r)
	case *contractpb.RolloutRenewalOpportunityOnExpirationGrpcRequest:
		return nil, a.rolloutRenewalOpportunityOnExpiration(ctx, r)
	default:
		tracing.TraceErr(span, eventstore.ErrInvalidRequestType)
		return nil, eventstore.ErrInvalidRequestType
	}
}

func (a *ContractAggregate) updateContract(ctx context.Context, request *contractpb.UpdateContractGrpcRequest) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "ContractAggregate.updateContract")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.Tenant)
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("aggregateVersion", a.GetVersion()))
	tracing.LogObjectAsJson(span, "request", request)

	externalSystem := commonmodel.ExternalSystem{}
	externalSystem.FromGrpc(request.ExternalSystemFields)

	sourceFields := commonmodel.Source{}
	sourceFields.FromGrpc(request.SourceFields)
	source := utils.StringFirstNonEmpty(sourceFields.Source, a.Contract.Source.SourceOfTruth)

	signedAt := utils.TimestampProtoToTimePtr(request.SignedAt)
	if signedAt != nil && signedAt.Equal(time.Time{}) {
		signedAt = nil
	}
	endedAt := utils.TimestampProtoToTimePtr(request.EndedAt)
	if endedAt != nil && endedAt.Equal(time.Time{}) {
		endedAt = nil
	}
	serviceStartedAt := utils.TimestampProtoToTimePtr(request.ServiceStartedAt)
	if serviceStartedAt != nil && serviceStartedAt.Equal(time.Time{}) {
		serviceStartedAt = nil
	}
	invoicingStartDate := utils.TimestampProtoToTimePtr(request.InvoicingStartDate)
	if invoicingStartDate != nil && invoicingStartDate.Equal(time.Time{}) {
		invoicingStartDate = nil
	}
	updatedAtNotNil := utils.IfNotNilTimeWithDefault(utils.TimestampProtoToTimePtr(request.UpdatedAt), utils.Now())

	dataFields := model.ContractDataFields{
		Name:                   request.Name,
		ServiceStartedAt:       serviceStartedAt,
		SignedAt:               signedAt,
		EndedAt:                endedAt,
		InvoicingStartDate:     invoicingStartDate,
		ContractUrl:            request.ContractUrl,
		Currency:               request.Currency,
		BillingCycleInMonths:   request.BillingCycleInMonths,
		AddressLine1:           request.AddressLine1,
		AddressLine2:           request.AddressLine2,
		Locality:               request.Locality,
		Country:                request.Country,
		Region:                 request.Region,
		Zip:                    request.Zip,
		OrganizationLegalName:  request.OrganizationLegalName,
		InvoiceEmail:           request.InvoiceEmailTo,
		InvoiceEmailCC:         request.InvoiceEmailCc,
		InvoiceEmailBCC:        request.InvoiceEmailBcc,
		InvoiceNote:            request.InvoiceNote,
		NextInvoiceDate:        utils.TimestampProtoToTimePtr(request.NextInvoiceDate),
		CanPayWithCard:         request.CanPayWithCard,
		CanPayWithDirectDebit:  request.CanPayWithDirectDebit,
		CanPayWithBankTransfer: request.CanPayWithBankTransfer,
		PayOnline:              request.PayOnline,
		PayAutomatically:       request.PayAutomatically,
		InvoicingEnabled:       request.InvoicingEnabled,
		AutoRenew:              request.AutoRenew,
		Check:                  request.Check,
		DueDays:                request.DueDays,
		LengthInMonths:         request.LengthInMonths,
		Approved:               request.Approved,
	}

	fieldsMask := extractFieldsMask(request.FieldsMask)

	// Set the approved field to true if the contract is already approved
	if a.Contract.Approved && utils.Contains(fieldsMask, event.FieldMaskApproved) {
		dataFields.Approved = true
	}

	// Validate the dates
	if isUpdated(event.FieldMaskEndedAt, fieldsMask) && dataFields.EndedAt != nil && (dataFields.SignedAt != nil && dataFields.EndedAt.Before(*dataFields.SignedAt) ||
		dataFields.ServiceStartedAt != nil && dataFields.EndedAt.Before(*dataFields.ServiceStartedAt)) {
		return errors.New(events2.FieldValidation + ": endedAt date must be after both signedAt and serviceStartedAt dates")
	}

	// Determine contract status based start and end dates
	if !isUpdated(event.FieldMaskServiceStartedAt, fieldsMask) {
		dataFields.ServiceStartedAt = a.Contract.ServiceStartedAt
	}
	if !isUpdated(event.FieldMaskEndedAt, fieldsMask) {
		dataFields.EndedAt = a.Contract.EndedAt
	}

	// Set renewal periods
	if !isUpdated(event.FieldMaskLengthInMonths, fieldsMask) {
		dataFields.LengthInMonths = a.Contract.LengthInMonths
	}

	updateEvent, err := event.NewContractUpdateEvent(
		a,
		dataFields,
		externalSystem,
		source,
		sourceFields.AppSource,
		updatedAtNotNil,
		fieldsMask,
	)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewContractUpdateEvent")
	}
	eventstore.EnrichEventWithMetadataExtended(&updateEvent, span, eventstore.EventMetadata{
		Tenant: a.Tenant,
		UserId: request.LoggedInUserId,
		App:    sourceFields.AppSource,
	})

	return a.Apply(updateEvent)
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

func extractFieldsMask(requestFieldsMask []contractpb.ContractFieldMask) []string {
	var fieldsMask []string
	for _, requestFieldMask := range requestFieldsMask {
		switch requestFieldMask {
		case contractpb.ContractFieldMask_CONTRACT_FIELD_NAME:
			fieldsMask = append(fieldsMask, event.FieldMaskName)
		case contractpb.ContractFieldMask_CONTRACT_FIELD_SERVICE_STARTED_AT:
			fieldsMask = append(fieldsMask, event.FieldMaskServiceStartedAt)
		case contractpb.ContractFieldMask_CONTRACT_FIELD_SIGNED_AT:
			fieldsMask = append(fieldsMask, event.FieldMaskSignedAt)
		case contractpb.ContractFieldMask_CONTRACT_FIELD_ENDED_AT:
			fieldsMask = append(fieldsMask, event.FieldMaskEndedAt)
		case contractpb.ContractFieldMask_CONTRACT_FIELD_CONTRACT_URL:
			fieldsMask = append(fieldsMask, event.FieldMaskContractURL)
		case contractpb.ContractFieldMask_CONTRACT_FIELD_BILLING_CYCLE:
			fieldsMask = append(fieldsMask, event.FieldMaskBillingCycle)
		case contractpb.ContractFieldMask_CONTRACT_FIELD_INVOICING_START_DATE:
			fieldsMask = append(fieldsMask, event.FieldMaskInvoicingStartDate)
		case contractpb.ContractFieldMask_CONTRACT_FIELD_CURRENCY:
			fieldsMask = append(fieldsMask, event.FieldMaskCurrency)
		case contractpb.ContractFieldMask_CONTRACT_FIELD_ADDRESS_LINE_1:
			fieldsMask = append(fieldsMask, event.FieldMaskAddressLine1)
		case contractpb.ContractFieldMask_CONTRACT_FIELD_ADDRESS_LINE_2:
			fieldsMask = append(fieldsMask, event.FieldMaskAddressLine2)
		case contractpb.ContractFieldMask_CONTRACT_FIELD_LOCALITY:
			fieldsMask = append(fieldsMask, event.FieldMaskLocality)
		case contractpb.ContractFieldMask_CONTRACT_FIELD_COUNTRY:
			fieldsMask = append(fieldsMask, event.FieldMaskCountry)
		case contractpb.ContractFieldMask_CONTRACT_FIELD_REGION:
			fieldsMask = append(fieldsMask, event.FieldMaskRegion)
		case contractpb.ContractFieldMask_CONTRACT_FIELD_ZIP:
			fieldsMask = append(fieldsMask, event.FieldMaskZip)
		case contractpb.ContractFieldMask_CONTRACT_FIELD_ORGANIZATION_LEGAL_NAME:
			fieldsMask = append(fieldsMask, event.FieldMaskOrganizationLegalName)
		case contractpb.ContractFieldMask_CONTRACT_FIELD_INVOICE_EMAIL_TO:
			fieldsMask = append(fieldsMask, event.FieldMaskInvoiceEmail)
		case contractpb.ContractFieldMask_CONTRACT_FIELD_INVOICE_EMAIL_CC:
			fieldsMask = append(fieldsMask, event.FieldMaskInvoiceEmailCC)
		case contractpb.ContractFieldMask_CONTRACT_FIELD_INVOICE_EMAIL_BCC:
			fieldsMask = append(fieldsMask, event.FieldMaskInvoiceEmailBCC)
		case contractpb.ContractFieldMask_CONTRACT_FIELD_INVOICE_NOTE:
			fieldsMask = append(fieldsMask, event.FieldMaskInvoiceNote)
		case contractpb.ContractFieldMask_CONTRACT_FIELD_NEXT_INVOICE_DATE:
			fieldsMask = append(fieldsMask, event.FieldMaskNextInvoiceDate)
		case contractpb.ContractFieldMask_CONTRACT_FIELD_INVOICING_ENABLED:
			fieldsMask = append(fieldsMask, event.FieldMaskInvoicingEnabled)
		case contractpb.ContractFieldMask_CONTRACT_FIELD_CAN_PAY_WITH_CARD:
			fieldsMask = append(fieldsMask, event.FieldMaskCanPayWithCard)
		case contractpb.ContractFieldMask_CONTRACT_FIELD_CAN_PAY_WITH_DIRECT_DEBIT:
			fieldsMask = append(fieldsMask, event.FieldMaskCanPayWithDirectDebit)
		case contractpb.ContractFieldMask_CONTRACT_FIELD_CAN_PAY_WITH_BANK_TRANSFER:
			fieldsMask = append(fieldsMask, event.FieldMaskCanPayWithBankTransfer)
		case contractpb.ContractFieldMask_CONTRACT_FIELD_PAY_ONLINE:
			fieldsMask = append(fieldsMask, event.FieldMaskPayOnline)
		case contractpb.ContractFieldMask_CONTRACT_FIELD_PAY_AUTOMATICALLY:
			fieldsMask = append(fieldsMask, event.FieldMaskPayAutomatically)
		case contractpb.ContractFieldMask_CONTRACT_FIELD_AUTO_RENEW:
			fieldsMask = append(fieldsMask, event.FieldMaskAutoRenew)
		case contractpb.ContractFieldMask_CONTRACT_FIELD_CHECK:
			fieldsMask = append(fieldsMask, event.FieldMaskCheck)
		case contractpb.ContractFieldMask_CONTRACT_FIELD_DUE_DAYS:
			fieldsMask = append(fieldsMask, event.FieldMaskDueDays)
		case contractpb.ContractFieldMask_CONTRACT_FIELD_LENGTH_IN_MONTHS:
			fieldsMask = append(fieldsMask, event.FieldMaskLengthInMonths)
		case contractpb.ContractFieldMask_CONTRACT_FIELD_APPROVED:
			fieldsMask = append(fieldsMask, event.FieldMaskApproved)
		case contractpb.ContractFieldMask_CONTRACT_FIELD_BILLING_CYCLE_IN_MONTHS:
			fieldsMask = append(fieldsMask, event.FieldMaskBillingCycleInMonths)
		}
	}
	fieldsMask = utils.RemoveDuplicates(fieldsMask)
	return fieldsMask
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
	case event.ContractUpdateV1:
		return a.onContractUpdate(evt)
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
