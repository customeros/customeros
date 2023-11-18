package contract

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/constants"
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contract/model"
	opportunitycmd "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/opportunity/command"
	opportunitycmdhandler "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/opportunity/command_handler"
	opportunitymodel "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/opportunity/model"
	servicelineitemmodel "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/service_line_item/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/graph_db"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/graph_db/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"math"
	"time"
)

type contractHandler struct {
	repositories        *repository.Repositories
	opportunityCommands *opportunitycmdhandler.CommandHandlers
	log                 logger.Logger
}

func NewContractHandler(log logger.Logger, repositories *repository.Repositories, opportunityCommands *opportunitycmdhandler.CommandHandlers) *contractHandler {
	return &contractHandler{
		repositories:        repositories,
		opportunityCommands: opportunityCommands,
		log:                 log,
	}
}

func (h *contractHandler) UpdateRenewalNextCycleDate(ctx context.Context, tenant, contractId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContractHandler.CalculateNextCycleDate")
	defer span.Finish()
	span.LogFields(log.String("contractId", contractId))

	contract, renewalOpportunity, done := h.assertContractAndRenewalOpportunity(ctx, span, tenant, contractId)
	if done {
		return nil
	}

	renewedAt := h.calculateNextCycleDate(contract.ServiceStartedAt, contract.RenewalCycle)
	err := h.opportunityCommands.UpdateRenewalOpportunityNextCycleDate.Handle(ctx, opportunitycmd.NewUpdateRenewalOpportunityNextCycleDateCommand(renewalOpportunity.Id, tenant, "", constants.AppSourceEventProcessingPlatform, nil, renewedAt))
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("UpdateRenewalOpportunityNextCycleDate command failed: %v", err.Error())
		return nil
	}

	return nil
}

func (h *contractHandler) calculateNextCycleDate(serviceStartedAt *time.Time, renewalCycle string) *time.Time {
	if serviceStartedAt == nil {
		return nil
	}

	renewalCycleNext := *serviceStartedAt
	for {
		switch renewalCycle {
		case string(model.MonthlyRenewalCycleString):
			renewalCycleNext = renewalCycleNext.AddDate(0, 1, 0)
		case string(model.AnnuallyRenewalCycleString):
			renewalCycleNext = renewalCycleNext.AddDate(1, 0, 0)
		default:
			return nil // invalid
		}

		if renewalCycleNext.After(utils.Now()) {
			break
		}
	}
	return &renewalCycleNext
}

func (h *contractHandler) UpdateRenewalArr(ctx context.Context, tenant, contractId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContractHandler.UpdateRenewalArr")
	defer span.Finish()
	span.LogFields(log.String("contractId", contractId))

	contract, renewalOpportunity, done := h.assertContractAndRenewalOpportunity(ctx, span, tenant, contractId)
	if done {
		return nil
	}

	// if contract already ended, return
	if contract.EndedAt != nil && contract.EndedAt.Before(utils.Now()) {
		return nil
	}

	maxArr, err := h.calculateMaxArr(ctx, tenant, contract, renewalOpportunity)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while calculating ARR for contract %s: %s", contractId, err.Error())
		return nil
	}
	// adjust with likelihood
	currentArr := h.calculateCurrentArrByLikelihood(maxArr, renewalOpportunity.RenewalDetails.RenewalLikelihood)

	err = h.opportunityCommands.UpdateOpportunity.Handle(ctx, opportunitycmd.NewUpdateOpportunityCommand(renewalOpportunity.Id, tenant, "",
		opportunitymodel.OpportunityDataFields{
			Amount:    currentArr,
			MaxAmount: maxArr,
		},
		commonmodel.Source{
			Source:    constants.SourceOpenline,
			AppSource: constants.AppSourceEventProcessingPlatform,
		},
		commonmodel.ExternalSystem{},
		nil,
		[]string{opportunitymodel.FieldMaskAmount, opportunitymodel.FieldMaskMaxAmount}))
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("UpdateOpportunity command failed: %v", err.Error())
		return nil
	}

	return nil
}

func (h *contractHandler) calculateMaxArr(ctx context.Context, tenant string, contract *entity.ContractEntity, renewalOpportunity *entity.OpportunityEntity) (float64, error) {
	var arr float64

	// Fetch service line items for the contract from the database
	sliDbNodes, err := h.repositories.ServiceLineItemRepository.GetAllForContract(ctx, tenant, contract.Id)
	if err != nil {
		return 0, err
	}
	serviceLineItems := entity.ServiceLineItemEntities{}
	for _, sliDbNode := range sliDbNodes {
		serviceLineItems = append(serviceLineItems, *graph_db.MapDbNodeToServiceLineItemEntity(*sliDbNode))
	}

	for _, sli := range serviceLineItems {
		annualPrice := float64(0)

		if sli.Billed == string(servicelineitemmodel.OnceBilledString) {
			annualPrice = sli.Price
		} else {
			annualPrice = float64(sli.Price) * float64(sli.Quantity)
			if sli.Billed == string(servicelineitemmodel.MonthlyBilledString) {
				annualPrice *= 12
			}
		}
		// Add to total ARR
		arr += annualPrice
	}

	// Adjust with end date
	currentDate := utils.Now()
	if contract.EndedAt != nil {
		// if end date before next renewal cycle, return set current arr to 0
		if renewalOpportunity.RenewalDetails.RenewedAt != nil && contract.EndedAt.Before(*renewalOpportunity.RenewalDetails.RenewedAt) {
			arr = 0
		}
		arr = prorateArr(arr, monthsUntilContractEnd(currentDate, *contract.EndedAt))
	}

	return arr, nil
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
	return monthlyRate * float64(monthsRemaining)
}

func (h *contractHandler) calculateCurrentArrByLikelihood(amount float64, likelihood string) float64 {
	var likelihoodFactor float64
	switch opportunitymodel.RenewalLikelihoodString(likelihood) {
	case opportunitymodel.RenewalLikelihoodStringHigh:
		likelihoodFactor = 1
	case opportunitymodel.RenewalLikelihoodStringMedium:
		likelihoodFactor = 0.5
	case opportunitymodel.RenewalLikelihoodStringLow:
		likelihoodFactor = 0.25
	case opportunitymodel.RenewalLikelihoodStringZero:
		likelihoodFactor = 0
	default:
		likelihoodFactor = 1
	}

	return math.Trunc(amount*likelihoodFactor*100) / 100
}

func (h *contractHandler) assertContractAndRenewalOpportunity(ctx context.Context, span opentracing.Span, tenant, contractId string) (*entity.ContractEntity, *entity.OpportunityEntity, bool) {
	if h.opportunityCommands == nil {
		tracing.TraceErr(span, errors.New("OpportunityCommands is nil"))
		h.log.Errorf("OpportunityCommands is nil")
		return nil, nil, true
	}

	contractDbNode, err := h.repositories.ContractRepository.GetContractById(ctx, tenant, contractId)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while getting contract %s: %s", contractId, err.Error())
		return nil, nil, true
	}
	contract := graph_db.MapDbNodeToContractEntity(*contractDbNode)

	// if contract is not frequency based, return
	if !model.IsFrequencyBasedRenewalCycle(contract.RenewalCycle) {
		return nil, nil, true
	}

	currentRenewalOpportunityDbNode, err := h.repositories.OpportunityRepository.GetOpenRenewalOpportunityForContract(ctx, tenant, contractId)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while getting renewal opportunity for contract %s: %s", contractId, err.Error())
		return nil, nil, true
	}
	// if there is no renewal opportunity, create one
	if currentRenewalOpportunityDbNode == nil {
		err = h.opportunityCommands.CreateRenewalOpportunity.Handle(ctx, opportunitycmd.NewCreateRenewalOpportunityCommand("", tenant, "", contractId, "", commonmodel.Source{}, nil, nil))
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("CreateRenewalOpportunity command failed: %v", err.Error())
			return nil, nil, true
		}
		return nil, nil, true
	}

	currentRenewalOpportunity := graph_db.MapDbNodeToOpportunityEntity(*currentRenewalOpportunityDbNode)

	return contract, currentRenewalOpportunity, false
}
