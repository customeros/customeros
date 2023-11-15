package contract

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/constants"
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contract/model"
	opportunitycmd "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/opportunity/command"
	opportunitycmdhandler "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/opportunity/command_handler"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/graph_db"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
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

	if h.opportunityCommands == nil {
		tracing.TraceErr(span, errors.New("OpportunityCommands is nil"))
		h.log.Errorf("OpportunityCommands is nil")
		return nil
	}

	contractDbNode, err := h.repositories.ContractRepository.GetContract(ctx, tenant, contractId)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while getting contract %s: %s", contractId, err.Error())
		return nil
	}
	contract := graph_db.MapDbNodeToContractEntity(*contractDbNode)

	// if contract is not frequency based, return
	if !model.IsFrequencyBasedRenewalCycle(contract.RenewalCycle) {
		return nil
	}

	currentRenewalOpportunityDbNode, err := h.repositories.OpportunityRepository.GetOpenRenewalOpportunityForContract(ctx, tenant, contractId)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while getting renewal opportunity for contract %s: %s", contractId, err.Error())
		return nil
	}
	// if there is no renewal opportunity, create one
	if currentRenewalOpportunityDbNode == nil {
		err = h.opportunityCommands.CreateRenewalOpportunity.Handle(ctx, opportunitycmd.NewCreateRenewalOpportunityCommand("", tenant, "", contractId, commonmodel.Source{}, nil, nil))
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("CreateRenewalOpportunity command failed: %v", err.Error())
			return nil
		}
		return nil
	}

	currentRenewalOpportunity := graph_db.MapDbNodeToOpportunityEntity(*currentRenewalOpportunityDbNode)

	// renewal opportunity exists, calculate next cycle date
	renewedAt := h.calculateNextCycleDate(contract.ServiceStartedAt, contract.RenewalCycle)
	if contract.ServiceStartedAt == nil {
		err = h.opportunityCommands.UpdateRenewalOpportunityNextCycleDate.Handle(ctx, opportunitycmd.NewUpdateRenewalOpportunityNextCycleDateCommand(currentRenewalOpportunity.Id, tenant, "", constants.AppSourceEventProcessingPlatform, nil, renewedAt))
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("UpdateRenewalOpportunityNextCycleDate command failed: %v", err.Error())
			return nil
		}
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
