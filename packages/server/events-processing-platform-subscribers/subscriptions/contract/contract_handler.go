package contract

import (
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/grpc_client"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/subscriptions"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/tracing"
	commonpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/common"
	opportunitypb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/opportunity"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"time"
)

type contractHandler struct {
	repositories *repository.Repositories
	log          logger.Logger
	grpcClients  *grpc_client.Clients
}

func NewContractHandler(log logger.Logger, repositories *repository.Repositories, grpcClients *grpc_client.Clients) *contractHandler {
	return &contractHandler{
		repositories: repositories,
		log:          log,
		grpcClients:  grpcClients,
	}
}

func (h *contractHandler) UpdateActiveRenewalOpportunityRenewDateAndArr(ctx context.Context, tenant, contractId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContractHandler.UpdateActiveRenewalOpportunityRenewDateAndArr")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	span.LogFields(log.String("contractId", contractId))

	contract, renewalOpportunity, done := h.assertContractAndRenewalOpportunity(ctx, tenant, contractId)
	if done {
		return nil
	}

	err := h.updateRenewalNextCycleDate(ctx, tenant, contract, renewalOpportunity, span)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil
	}
	err = h.updateRenewalArr(ctx, tenant, contract, renewalOpportunity, span)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil
	}
	return nil
}

func (h *contractHandler) UpdateActiveRenewalOpportunityNextCycleDate(ctx context.Context, tenant, contractId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContractHandler.CalculateNextCycleDate")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	span.LogFields(log.String("contractId", contractId))

	contract, renewalOpportunity, done := h.assertContractAndRenewalOpportunity(ctx, tenant, contractId)
	if done {
		return nil
	}

	return h.updateRenewalNextCycleDate(ctx, tenant, contract, renewalOpportunity, span)
}

func (h *contractHandler) UpdateActiveRenewalOpportunityArr(ctx context.Context, tenant, contractId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContractHandler.UpdateActiveRenewalOpportunityArr")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	span.LogFields(log.String("contractId", contractId))

	contract, renewalOpportunity, done := h.assertContractAndRenewalOpportunity(ctx, tenant, contractId)
	if done {
		return nil
	}

	return h.updateRenewalArr(ctx, tenant, contract, renewalOpportunity, span)
}

func (h *contractHandler) UpdateActiveRenewalOpportunityLikelihood(ctx context.Context, tenant, contractId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContractHandler.UpdateActiveRenewalOpportunityLikelihood")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	span.LogFields(log.String("contractId", contractId))

	opportunityDbNode, err := h.repositories.Neo4jRepositories.OpportunityReadRepository.GetActiveRenewalOpportunityForContract(ctx, tenant, contractId)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while getting renewal opportunity for contract %s: %s", contractId, err.Error())
		return err
	}
	if opportunityDbNode == nil {
		h.log.Infof("No open renewal opportunity found for contract %s", contractId)
		return nil
	}
	contractDbNode, err := h.repositories.Neo4jRepositories.ContractReadRepository.GetContractById(ctx, tenant, contractId)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while getting contract %s: %s", contractId, err.Error())
		return err
	}
	contractEntity := neo4jmapper.MapDbNodeToContractEntity(contractDbNode)
	opportunityEntity := neo4jmapper.MapDbNodeToOpportunityEntity(opportunityDbNode)

	var renewalLikelihood neo4jenum.RenewalLikelihood
	if contractEntity.EndedAt != nil &&
		opportunityEntity.RenewalDetails.RenewalLikelihood != neo4jenum.RenewalLikelihoodZero &&
		opportunityEntity.RenewalDetails.RenewedAt != nil &&
		contractEntity.EndedAt.Before(*opportunityEntity.RenewalDetails.RenewedAt) {
		// check if likelihood should be set to Zero
		renewalLikelihood = neo4jenum.RenewalLikelihoodZero
	} else if opportunityEntity.RenewalDetails.RenewalLikelihood == neo4jenum.RenewalLikelihoodZero &&
		opportunityEntity.RenewalDetails.RenewedAt != nil &&
		(contractEntity.EndedAt == nil || contractEntity.EndedAt.After(*opportunityEntity.RenewalDetails.RenewedAt)) {
		// check if likelihood should be set to Medium
		renewalLikelihood = neo4jenum.RenewalLikelihoodMedium
	}

	if renewalLikelihood != "" {
		ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
		_, err = subscriptions.CallEventsPlatformGRPCWithRetry[*opportunitypb.OpportunityIdGrpcResponse](func() (*opportunitypb.OpportunityIdGrpcResponse, error) {
			return h.grpcClients.OpportunityClient.UpdateRenewalOpportunity(ctx, &opportunitypb.UpdateRenewalOpportunityGrpcRequest{
				Tenant:            tenant,
				Id:                opportunityEntity.Id,
				RenewalLikelihood: renewalLikelihoodForGrpcRequest(renewalLikelihood),
				SourceFields: &commonpb.SourceFields{
					AppSource: constants.AppSourceEventProcessingPlatform,
				},
				FieldsMask: []opportunitypb.OpportunityMaskField{opportunitypb.OpportunityMaskField_OPPORTUNITY_PROPERTY_RENEWAL_LIKELIHOOD},
			})
		})
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("UpdateRenewalOpportunity failed: %s", err.Error())
			return errors.Wrap(err, "UpdateRenewalOpportunity")
		}
	}

	return nil
}

func (h *contractHandler) updateRenewalNextCycleDate(ctx context.Context, tenant string, contractEntity *neo4jentity.ContractEntity, renewalOpportunityEntity *neo4jentity.OpportunityEntity, span opentracing.Span) error {
	if contractEntity.IsEnded() && renewalOpportunityEntity != nil {
		ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
		_, err := subscriptions.CallEventsPlatformGRPCWithRetry[*opportunitypb.OpportunityIdGrpcResponse](func() (*opportunitypb.OpportunityIdGrpcResponse, error) {
			return h.grpcClients.OpportunityClient.CloseLooseOpportunity(ctx, &opportunitypb.CloseLooseOpportunityGrpcRequest{
				Tenant:    tenant,
				Id:        renewalOpportunityEntity.Id,
				AppSource: constants.AppSourceEventProcessingPlatform,
			})
		})
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("CloseLooseOpportunity failed: %s", err.Error())
			return errors.Wrap(err, "CloseLooseOpportunity")
		}
		return nil
	}

	renewedAt := h.calculateNextCycleDate(contractEntity.ServiceStartedAt, contractEntity.RenewalCycle, contractEntity.RenewalPeriods)
	if !utils.IsEqualTimePtr(renewedAt, renewalOpportunityEntity.RenewalDetails.RenewedAt) {
		ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
		_, err := subscriptions.CallEventsPlatformGRPCWithRetry[*opportunitypb.OpportunityIdGrpcResponse](func() (*opportunitypb.OpportunityIdGrpcResponse, error) {
			return h.grpcClients.OpportunityClient.UpdateRenewalOpportunityNextCycleDate(ctx, &opportunitypb.UpdateRenewalOpportunityNextCycleDateGrpcRequest{
				OpportunityId: renewalOpportunityEntity.Id,
				Tenant:        tenant,
				AppSource:     constants.AppSourceEventProcessingPlatform,
				RenewedAt:     utils.ConvertTimeToTimestampPtr(renewedAt),
			})
		})
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("UpdateRenewalOpportunityNextCycleDate failed: %s", err.Error())
			return errors.Wrap(err, "UpdateRenewalOpportunityNextCycleDate")
		}
	}

	return nil
}

func (h *contractHandler) calculateNextCycleDate(serviceStartedAt *time.Time, renewalCycle neo4jenum.RenewalCycle, renewalPeriods *int64) *time.Time {
	if serviceStartedAt == nil {
		return nil
	}

	renewalCycleNext := *serviceStartedAt
	for {
		switch renewalCycle {
		case neo4jenum.RenewalCycleMonthlyRenewal:
			renewalCycleNext = renewalCycleNext.AddDate(0, 1, 0)
		case neo4jenum.RenewalCycleQuarterlyRenewal:
			renewalCycleNext = renewalCycleNext.AddDate(0, 3, 0)
		case neo4jenum.RenewalCycleAnnualRenewal:
			renewalYears := 1
			if renewalPeriods != nil {
				renewalYears = int(*renewalPeriods)
			}
			renewalCycleNext = renewalCycleNext.AddDate(int(renewalYears), 0, 0)
		default:
			return nil // invalid
		}

		if renewalCycleNext.After(utils.Now()) {
			break
		}
	}
	return &renewalCycleNext
}

func (h *contractHandler) updateRenewalArr(ctx context.Context, tenant string, contract *neo4jentity.ContractEntity, renewalOpportunity *neo4jentity.OpportunityEntity, span opentracing.Span) error {
	// if contract already ended, return
	if contract.IsEnded() {
		span.LogFields(log.Bool("contract ended", true))
		return nil
	}

	maxArr, err := h.calculateMaxArr(ctx, tenant, contract, renewalOpportunity, span)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while calculating ARR for contract %s: %s", contract.Id, err.Error())
		return nil
	}
	// adjust with likelihood
	currentArr := h.calculateCurrentArrByLikelihood(maxArr, renewalOpportunity.RenewalDetails.RenewalLikelihood.String())

	ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
	_, err = subscriptions.CallEventsPlatformGRPCWithRetry[*opportunitypb.OpportunityIdGrpcResponse](func() (*opportunitypb.OpportunityIdGrpcResponse, error) {
		return h.grpcClients.OpportunityClient.UpdateOpportunity(ctx, &opportunitypb.UpdateOpportunityGrpcRequest{
			Tenant:    tenant,
			Id:        renewalOpportunity.Id,
			Amount:    currentArr,
			MaxAmount: maxArr,
			SourceFields: &commonpb.SourceFields{
				AppSource: constants.AppSourceEventProcessingPlatform,
				Source:    constants.SourceOpenline,
			},
			FieldsMask: []opportunitypb.OpportunityMaskField{
				opportunitypb.OpportunityMaskField_OPPORTUNITY_PROPERTY_AMOUNT,
				opportunitypb.OpportunityMaskField_OPPORTUNITY_PROPERTY_MAX_AMOUNT,
			},
		})
	})
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("UpdateOpportunity failed: %s", err.Error())
	}

	return nil
}

func (h *contractHandler) calculateMaxArr(ctx context.Context, tenant string, contract *neo4jentity.ContractEntity, renewalOpportunity *neo4jentity.OpportunityEntity, span opentracing.Span) (float64, error) {
	var arr float64

	// Fetch service line items for the contract from the database
	sliDbNodes, err := h.repositories.Neo4jRepositories.ServiceLineItemReadRepository.GetAllForContract(ctx, tenant, contract.Id)
	if err != nil {
		tracing.TraceErr(span, err)
		return 0, err
	}
	serviceLineItems := neo4jentity.ServiceLineItemEntities{}
	for _, sliDbNode := range sliDbNodes {
		sli := neo4jmapper.MapDbNodeToServiceLineItemEntity(sliDbNode)
		serviceLineItems = append(serviceLineItems, *sli)
	}

	span.LogFields(log.Int("service line items count", len(serviceLineItems)))
	for _, sli := range serviceLineItems {
		if sli.IsEnded() {
			span.LogFields(log.Bool(fmt.Sprintf("service line item {%s} ended", sli.ID), true))
			continue
		}
		span.LogFields(log.Object(fmt.Sprintf("service line item {%s}:", sli.ID), sli))
		annualPrice := float64(0)
		if sli.Billed == neo4jenum.BilledTypeAnnually {
			annualPrice = float64(sli.Price) * float64(sli.Quantity)
		} else if sli.Billed == neo4jenum.BilledTypeMonthly {
			annualPrice = float64(sli.Price) * float64(sli.Quantity)
			annualPrice *= 12
		} else if sli.Billed == neo4jenum.BilledTypeQuarterly {
			annualPrice = float64(sli.Price) * float64(sli.Quantity)
			annualPrice *= 4
		}
		span.LogFields(log.Float64(fmt.Sprintf("service line item {%s} added ARR value:", sli.ID), annualPrice))
		// Add to total ARR
		arr += annualPrice
	}

	// Adjust with end date
	if contract.EndedAt != nil {
		span.LogFields(log.Bool("ARR prorated with contract end date", true))
		arr = prorateArr(arr, monthsUntilContractEnd(utils.Now(), *contract.EndedAt))
	}

	return utils.TruncateFloat64(arr, 2), nil
}

func monthsUntilContractEnd(start, end time.Time) int {
	yearDiff := end.Year() - start.Year()
	monthDiff := int(end.Month()) - int(start.Month())

	// Total difference in months
	totalMonths := yearDiff*12 + monthDiff

	// If the end day is before the start day in the month, subtract a month
	if end.Day() < start.Day() {
		totalMonths--
	}

	if totalMonths < 0 {
		totalMonths = 0
	}

	return totalMonths
}

func prorateArr(arr float64, monthsRemaining int) float64 {
	if monthsRemaining > 12 {
		return arr
	}
	monthlyRate := arr / 12
	return utils.TruncateFloat64(monthlyRate*float64(monthsRemaining), 2)
}

func (h *contractHandler) calculateCurrentArrByLikelihood(amount float64, likelihood string) float64 {
	var likelihoodFactor float64
	switch neo4jenum.RenewalLikelihood(likelihood) {
	case neo4jenum.RenewalLikelihoodHigh:
		likelihoodFactor = 1
	case neo4jenum.RenewalLikelihoodMedium:
		likelihoodFactor = 0.5
	case neo4jenum.RenewalLikelihoodLow:
		likelihoodFactor = 0.25
	case neo4jenum.RenewalLikelihoodZero:
		likelihoodFactor = 0
	default:
		likelihoodFactor = 1
	}

	return utils.TruncateFloat64(amount*likelihoodFactor, 2)
}

func (h *contractHandler) assertContractAndRenewalOpportunity(ctx context.Context, tenant, contractId string) (*neo4jentity.ContractEntity, *neo4jentity.OpportunityEntity, bool) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContractHandler.assertContractAndRenewalOpportunity")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	span.LogFields(log.String("contractId", contractId))

	contractDbNode, err := h.repositories.Neo4jRepositories.ContractReadRepository.GetContractById(ctx, tenant, contractId)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while getting contract %s: %s", contractId, err.Error())
		return nil, nil, true
	}
	contract := neo4jmapper.MapDbNodeToContractEntity(contractDbNode)

	// if contract is not frequency based, return
	if !neo4jenum.IsFrequencyBasedRenewalCycle(contract.RenewalCycle) {
		return nil, nil, true
	}

	currentRenewalOpportunityDbNode, err := h.repositories.Neo4jRepositories.OpportunityReadRepository.GetActiveRenewalOpportunityForContract(ctx, tenant, contractId)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while getting renewal opportunity for contract %s: %s", contractId, err.Error())
		return nil, nil, true
	}

	// if there is no renewal opportunity, create one
	if currentRenewalOpportunityDbNode == nil {
		if !contract.IsEnded() {
			ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
			_, err = subscriptions.CallEventsPlatformGRPCWithRetry[*opportunitypb.OpportunityIdGrpcResponse](func() (*opportunitypb.OpportunityIdGrpcResponse, error) {
				return h.grpcClients.OpportunityClient.CreateRenewalOpportunity(ctx, &opportunitypb.CreateRenewalOpportunityGrpcRequest{
					Tenant:     tenant,
					ContractId: contractId,
					SourceFields: &commonpb.SourceFields{
						AppSource: constants.AppSourceEventProcessingPlatform,
					},
				})
			})
			if err != nil {
				tracing.TraceErr(span, err)
				h.log.Errorf("CreateRenewalOpportunity command failed: %v", err.Error())
				return nil, nil, true
			}
			span.LogFields(log.Bool("renewal opportunity create requested", true))
		}
		return nil, nil, true
	}

	currentRenewalOpportunity := neo4jmapper.MapDbNodeToOpportunityEntity(currentRenewalOpportunityDbNode)

	return contract, currentRenewalOpportunity, false
}

func renewalLikelihoodForGrpcRequest(renewalLikelihood neo4jenum.RenewalLikelihood) opportunitypb.RenewalLikelihood {
	switch renewalLikelihood {
	case neo4jenum.RenewalLikelihoodHigh:
		return opportunitypb.RenewalLikelihood_HIGH_RENEWAL
	case neo4jenum.RenewalLikelihoodMedium:
		return opportunitypb.RenewalLikelihood_MEDIUM_RENEWAL
	case neo4jenum.RenewalLikelihoodLow:
		return opportunitypb.RenewalLikelihood_LOW_RENEWAL
	case neo4jenum.RenewalLikelihoodZero:
		return opportunitypb.RenewalLikelihood_ZERO_RENEWAL
	default:
		return opportunitypb.RenewalLikelihood_HIGH_RENEWAL
	}
}
