package contract

import (
	"context"
	"fmt"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"github.com/opentracing/opentracing-go/log"
	"time"

	neo4jmodel "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/model"

	neo4jentity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"
	neo4jenum "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/enum"
	neo4jmapper "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/mapper"
	neoRepo "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/repository"

	"github.com/pkg/errors"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/constants"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/data_fields"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/dto"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/model"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/events"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
)

type contractService struct {
	log          logger.Logger
	neo4j        *neoRepo.Repositories
	events       *events.EventsService
	opportunity  interfaces.OpportunityService
	organization interfaces.OrganizationService
}

func NewContractService(log logger.Logger, neo4j *neoRepo.Repositories, events *events.EventsService, opportunity interfaces.OpportunityService, org interfaces.OrganizationService) interfaces.ContractService {
	return &contractService{
		log:          log,
		neo4j:        neo4j,
		events:       events,
		opportunity:  opportunity,
		organization: org,
	}
}

func (s *contractService) SetOpportunityService(opportunity interfaces.OpportunityService) {
	s.opportunity = opportunity
}

func (s *contractService) SetOrganizationService(org interfaces.OrganizationService) {
	s.organization = org
}

func (s *contractService) IsInitialized() bool {
	return utils.IsInitialized(s)
}

func (s *contractService) GetById(ctx context.Context, contractId string) (*neo4jentity.ContractEntity, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "ContractService.GetById")
	defer spans.Finish()

	spans.LogKV("contractId", contractId)

	if contractDbNode, err := s.neo4j.ContractReadRepository.GetContractById(ctx, common.GetContext(ctx).Tenant, contractId); err != nil {
		spans.TraceError(err)
		wrappedErr := errors.Wrap(err, fmt.Sprintf("Contract with id {%s} not found", contractId))
		return nil, wrappedErr
	} else {
		return neo4jmapper.MapDbNodeToContractEntity(contractDbNode), nil
	}
}

func (s *contractService) Save(ctx context.Context, txWithPostCommit *utils.TxWithPostCommit, id *string, dataFields data_fields.ContractSaveFields) (string, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "ContractService.Save")
	defer spans.Finish()

	spans.LogObjectAsJson("dataFields", dataFields)

	// validate tenant
	err := common.ValidateTenant(ctx)
	if err != nil {
		spans.TraceError(err)
		return "", err
	}
	tenant := common.GetTenantFromContext(ctx)

	createFlow := false
	contractId := ""

	if utils.IfNotNilString(id) == "" {
		createFlow = true
		spans.LogKV("flow", "create")
		contractId, err = s.neo4j.CommonReadRepository.GenerateId(ctx, tenant, model.NodeLabelContract)
		if err != nil {
			spans.TraceError(err)
			return "", err
		}
		// set default fields for create flow
		if dataFields.CreatedAt == nil {
			dataFields.CreatedAt = utils.NowPtr()
		} else {
			dataFields.CreatedAt = utils.TimePtr(utils.NowIfZero(*dataFields.CreatedAt))
		}
		if utils.IfNotNilString(dataFields.AppSource) == "" {
			dataFields.AppSource = utils.StringPtr(common.GetAppSourceFromContext(ctx))
		}
		if utils.IfNotNilString(dataFields.Source) == "" {
			dataFields.Source = utils.StringPtr(neo4jentity.DataSourceOpenline.String())
		}
		if utils.IfNotNilString(dataFields.Name) == "" {
			if utils.IfNotNilString(dataFields.OrganizationId) != "" {
				organizationEntity, err := s.organization.GetById(ctx, tenant, utils.IfNotNilString(dataFields.OrganizationId))
				if err != nil {
					spans.TraceError(errors.Wrap(err, "unable to get organization"))
					s.log.Errorf("unable to get organization: %s", err.Error())
					return "", err
				}
				dataFields.Name = utils.StringPtr(organizationEntity.Name)
			}
		}
	} else {
		spans.LogKV("flow", "update")
		contractId = *id

		// validate contract exists
		exists, err := s.neo4j.CommonReadRepository.ExistsById(ctx, tenant, contractId, model.NodeLabelContract)
		if err != nil || !exists {
			err = errors.New("contract not found")
			spans.TraceError(err)
			return "", err
		}
	}
	spans.TagEntity(contractId)

	var beforeUpdateContractEntity *neo4jentity.ContractEntity

	_, err = utils.ExecuteWriteInTransactionWithPostCommitActions(ctx, s.neo4j.Neo4jDriver, s.neo4j.Database, txWithPostCommit, func(txWithPostCommit *utils.TxWithPostCommit) (any, error) {

		if createFlow {
			err := s.neo4j.ContractWriteRepository.CreateForOrganization(ctx, txWithPostCommit.Tx, tenant, contractId, dataFields)
			if err != nil {
				spans.TraceError(err)
				return "", err
			}

			if dataFields.ExternalSystem != nil && dataFields.ExternalSystem.Available() {
				externalSystemData := neo4jmodel.ExternalSystem{
					ExternalSystemId: dataFields.ExternalSystem.ExternalSystemId,
					ExternalUrl:      dataFields.ExternalSystem.ExternalUrl,
					ExternalId:       dataFields.ExternalSystem.ExternalId,
					ExternalIdSecond: dataFields.ExternalSystem.ExternalIdSecond,
					ExternalSource:   dataFields.ExternalSystem.ExternalSource,
					SyncDate:         dataFields.ExternalSystem.SyncDate,
				}
				err = s.neo4j.ExternalSystemWriteRepository.LinkWithEntityInTx(ctx, txWithPostCommit.Tx, tenant, contractId, model.NodeLabelContract, externalSystemData)
				if err != nil {
					spans.TraceError(err)
					s.log.Errorf("Error while linking contract %s with external system %s: %s", contractId, dataFields.ExternalSystem.ExternalSystemId, err.Error())
					return "", err
				}
			}
		} else {
			contractDbNode, err := s.neo4j.ContractReadRepository.GetContractById(ctx, tenant, contractId)
			if err != nil {
				spans.TraceError(err)
				return contractId, err
			}
			beforeUpdateContractEntity = neo4jmapper.MapDbNodeToContractEntity(contractDbNode)

			err = s.neo4j.ContractWriteRepository.UpdateContract(ctx, txWithPostCommit.Tx, tenant, contractId, dataFields)
			if err != nil {
				spans.TraceError(err)
				return "", err
			}

			if dataFields.ExternalSystem != nil && dataFields.ExternalSystem.Available() {
				externalSystemData := neo4jmodel.ExternalSystem{
					ExternalSystemId: dataFields.ExternalSystem.ExternalSystemId,
					ExternalUrl:      dataFields.ExternalSystem.ExternalUrl,
					ExternalId:       dataFields.ExternalSystem.ExternalId,
					ExternalIdSecond: dataFields.ExternalSystem.ExternalIdSecond,
					ExternalSource:   dataFields.ExternalSystem.ExternalSource,
					SyncDate:         dataFields.ExternalSystem.SyncDate,
				}
				err = s.neo4j.ExternalSystemWriteRepository.LinkWithEntityInTx(ctx, txWithPostCommit.Tx, tenant, contractId, model.NodeLabelContract, externalSystemData)
				if err != nil {
					spans.TraceError(err)
					s.log.Errorf("Error while linking contract %s with external system %s: %s", contractId, dataFields.ExternalSystem.ExternalSystemId, err.Error())
					return "", err
				}
			}
		}

		txWithPostCommit.AddPostCommitAction(func(ctx context.Context) error {
			// send events
			if createFlow {
				err = s.events.Publisher.PublishFanoutEvent(ctx, contractId, model.CONTRACT, dto.CreateContract{ContractSaveFields: dataFields})
				if err != nil {
					spans.TraceError(errors.Wrap(err, "unable to publish message CreateContract"))
				}
				s.events.Publisher.PublishNotification(ctx, tenant, contractId, model.CONTRACT, utils.NewEventCompletedDetails().WithCreate())
			} else {
				err = s.events.Publisher.PublishFanoutEvent(ctx, contractId, model.CONTRACT, dto.UpdateContract{ContractSaveFields: dataFields})
				if err != nil {
					spans.TraceError(errors.Wrap(err, "unable to publish message UpdateContract"))
				}
				if dataFields.AppSource == nil || *dataFields.AppSource != constants.AppSourceCustomerOsApi {
					s.events.Publisher.PublishNotification(ctx, tenant, contractId, model.CONTRACT, utils.NewEventCompletedDetails().WithUpdate())
				}
			}

			// post save actions
			if createFlow {
				err = s.postCreateContract(ctx, tenant, contractId, dataFields)
				if err != nil {
					spans.TraceError(err)
					s.log.Errorf("Error while post create contract %s: %s", contractId, err.Error())
				}
			} else {
				err = s.postUpdateContract(ctx, tenant, contractId, beforeUpdateContractEntity)
				if err != nil {
					spans.TraceError(err)
					s.log.Errorf("Error while post create contract %s: %s", contractId, err.Error())
				}
				if dataFields.AppSource == nil || *dataFields.AppSource != constants.AppSourceCustomerOsApi {
					s.events.Publisher.PublishNotification(ctx, tenant, contractId, model.CONTRACT, utils.NewEventCompletedDetails().WithUpdate())
				}
			}

			return nil
		})
		return nil, nil
	})
	if err != nil {
		spans.TraceError(err)
		return "", err
	}

	return contractId, nil
}

func (s *contractService) SoftDelete(ctx context.Context, contractId string) error {
	spans, ctx := telemetry.StartServiceSpan(ctx, "ContractService.SoftDelete")
	defer spans.Finish()

	spans.TagEntity(contractId)

	// validate tenant
	err := common.ValidateTenant(ctx)
	if err != nil {
		spans.TraceError(err)
		return err
	}
	tenant := common.GetTenantFromContext(ctx)

	// fetch organization of the contract
	organizationDbNode, err := s.neo4j.OrganizationReadRepository.GetOrganizationByContractId(ctx, tenant, contractId)
	if err != nil {
		spans.TraceError(err)
		s.log.Errorf("Error while getting organization for contract %s: %s", contractId, err.Error())
		return nil
	}
	if organizationDbNode == nil {
		s.log.Errorf("Organization not found for contract %s", contractId)
		return nil
	}
	organization := neo4jmapper.MapDbNodeToOrganizationEntity(organizationDbNode)

	err = s.neo4j.ContractWriteRepository.SoftDelete(ctx, tenant, contractId)
	if err != nil {
		spans.TraceError(err)
		s.log.Errorf("Error while deleting contract %s: %s", contractId, err.Error())
		return err
	}

	err = s.organization.UpdateRenewalSummary(ctx, organization.ID)
	if err != nil {
		spans.TraceError(err)
		s.log.Errorf("Error while updating renewal summary for organization %s: %s", organization.ID, err.Error())
	}

	err = s.neo4j.InvoiceWriteRepository.DeletePreviewCycleInvoices(ctx, tenant, contractId, "")
	if err != nil {
		spans.TraceError(err)
		s.log.Errorf("Error while deleting preview invoice for contract %s: %s", contractId, err.Error())
		return err
	}

	s.events.Publisher.PublishNotification(ctx, tenant, contractId, model.CONTRACT, utils.NewEventCompletedDetails().WithDelete())

	return nil
}

func (s *contractService) postCreateContract(ctx context.Context, tenant, contractId string, dataFields data_fields.ContractSaveFields) error {
	spans, ctx := telemetry.StartServiceSpan(ctx, "ContractService.postCreateContract")
	defer spans.Finish()
	spans.TagEntity(contractId)

	_, _, err := s.updateStatus(ctx, tenant, contractId)
	if err != nil {
		spans.TraceError(err)
		s.log.Errorf("Error while updating contract %s status: %s", contractId, err.Error())
	}

	if dataFields.LengthInMonths != nil && *dataFields.LengthInMonths > 0 {
		_, err = s.opportunity.CreateRenewalOpportunity(ctx, nil, &data_fields.OpportunityFields{
			ContractId: &contractId,
			Source:     dataFields.Source,
			AppSource:  dataFields.AppSource,
		})
		if err != nil {
			spans.TraceError(err)
			s.log.Errorf("CreateRenewalOpportunity failed: %s", err.Error())
		}
	}

	return nil
}

func (s *contractService) postUpdateContract(ctx context.Context, tenant string, contractId string, beforeUpdateContractEntity *neo4jentity.ContractEntity) error {
	spans, ctx := telemetry.StartServiceSpan(ctx, "ContractService.postCreateContract")
	defer spans.Finish()
	spans.TagEntity(contractId)

	_, statusChanged, err := s.updateStatus(ctx, tenant, contractId)
	if err != nil {
		spans.TraceError(err)
		s.log.Errorf("Error while updating contract %s status: %s", contractId, err.Error())
	}

	if statusChanged {
		err = s.updateOrganizationRelationship(ctx, tenant, contractId, statusChanged)
		if err != nil {
			spans.TraceError(err)
			s.log.Errorf("Error while updating organization relationship for contract %s: %s", contractId, err.Error())
		}
	}
	contractDbNode, err := s.neo4j.ContractReadRepository.GetContractById(ctx, tenant, contractId)
	if err != nil {
		spans.TraceError(err)
		return err
	}
	afterUpdateContractEntity := neo4jmapper.MapDbNodeToContractEntity(contractDbNode)

	if beforeUpdateContractEntity.LengthInMonths > 0 && afterUpdateContractEntity.LengthInMonths == 0 {
		err = s.neo4j.ContractWriteRepository.SuspendActiveRenewalOpportunity(ctx, tenant, contractId)
		if err != nil {
			spans.TraceError(err)
			s.log.Errorf("Error while suspending renewal opportunity for contract %s: %s", contractId, err.Error())
		}
		organizationDbNode, err := s.neo4j.OrganizationReadRepository.GetOrganizationByContractId(ctx, tenant, contractId)
		if err != nil {
			spans.TraceError(err)
			s.log.Errorf("Error while getting organization for contract %s: %s", contractId, err.Error())
			return nil
		}
		if organizationDbNode == nil {
			s.log.Errorf("Organization not found for contract %s", contractId)
			return nil
		}
		organization := neo4jmapper.MapDbNodeToOrganizationEntity(organizationDbNode)

		err = s.organization.UpdateRenewalSummary(ctx, organization.ID)
		if err != nil {
			spans.TraceError(err)
			s.log.Errorf("Error while updating renewal summary for organization %s: %s", organization.ID, err.Error())
		}

		err = s.neo4j.OrganizationWriteRepository.UpdateArr(ctx, tenant, organization.ID)
		if err != nil {
			spans.TraceError(err)
			s.log.Errorf("Error while updating ARR for organization %s: %s", organization.ID, err.Error())
		}

	} else {
		if beforeUpdateContractEntity.LengthInMonths == 0 && afterUpdateContractEntity.LengthInMonths > 0 {
			err = s.neo4j.ContractWriteRepository.ActivateSuspendedRenewalOpportunity(ctx, tenant, contractId)
			if err != nil {
				spans.TraceError(err)
				s.log.Errorf("Error while activating renewal opportunity for contract %s: %s", contractId, err.Error())
			}
		}
		err = s.UpdateActiveRenewalOpportunityRenewDateAndArr(ctx, tenant, contractId)
		if err != nil {
			spans.TraceError(err)
			s.log.Errorf("error while updating renewal opportunity for contract %s: %s", contractId, err.Error())
		}
	}

	if beforeUpdateContractEntity.ContractStatus != afterUpdateContractEntity.ContractStatus {
		s.createActionForStatusChange(ctx, tenant, contractId, string(afterUpdateContractEntity.ContractStatus), afterUpdateContractEntity.Name)
	}

	err = s.UpdateActiveRenewalOpportunityLikelihood(ctx, tenant, contractId)
	if err != nil {
		spans.TraceError(err)
		s.log.Errorf("error while updating renewal opportunity for contract %s: %s", contractId, err.Error())
	}
	return nil
}

func (s *contractService) updateStatus(ctx context.Context, tenant, contractId string) (string, bool, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "ContractService.updateStatus")
	defer spans.Finish()
	spans.TagEntity(contractId)

	contractDbNode, err := s.neo4j.ContractReadRepository.GetContractById(ctx, tenant, contractId)
	if err != nil {
		spans.TraceError(err)
		s.log.Errorf("Error while getting contract %s: %s", contractId, err.Error())
		return "", false, err
	}
	contractEntity := neo4jmapper.MapDbNodeToContractEntity(contractDbNode)

	status, err := s.deriveContractStatus(ctx, tenant, *contractEntity)
	if err != nil {
		spans.TraceError(err)
		s.log.Errorf("Error while deriving contract %s status: %s", contractId, err.Error())
		return "", false, err
	}
	statusChanged := contractEntity.ContractStatus.String() != status

	if statusChanged {
		err = s.neo4j.ContractWriteRepository.UpdateStatus(ctx, tenant, contractId, status)
		if err != nil {
			spans.TraceError(err)
			s.log.Errorf("Error while updating contract %s status: %s", contractId, err.Error())
			return "", false, err
		}

		s.events.Publisher.PublishNotification(ctx, tenant, contractId, model.CONTRACT, utils.NewEventCompletedDetails().WithUpdate())

		err = s.events.Publisher.PublishFanoutEvent(ctx, contractId, model.CONTRACT, dto.ChangeStatusForContract{Status: status})
		if err != nil {
			spans.TraceError(errors.Wrap(err, "unable to publish message ChangeStatusForContract"))
		}
	}

	return status, statusChanged, nil
}

func (s *contractService) deriveContractStatus(ctx context.Context, tenant string, contractEntity neo4jentity.ContractEntity) (string, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "ContractService.deriveContractStatus")
	defer spans.Finish()

	now := utils.Now()

	// If endedAt is not nil and is in the past, the contract is considered Ended.
	if contractEntity.IsEnded() {
		spans.LogKV("result.status", neo4jenum.ContractStatusEnded.String())
		return neo4jenum.ContractStatusEnded.String(), nil
	}

	// check if contract is draft
	if !contractEntity.Approved {
		spans.LogKV("result.status", neo4jenum.ContractStatusDraft.String())
		return neo4jenum.ContractStatusDraft.String(), nil
	}

	// Check contract is scheduled
	if contractEntity.ServiceStartedAt == nil || contractEntity.ServiceStartedAt.After(now) {
		spans.LogKV("result.status", neo4jenum.ContractStatusScheduled.String())
		return neo4jenum.ContractStatusScheduled.String(), nil
	}

	// Check if contract is out of contract
	if !contractEntity.AutoRenew {
		// fetch active renewal opportunity for the contract
		opportunityDbNode, err := s.neo4j.OpportunityReadRepository.GetActiveRenewalOpportunityForContract(ctx, tenant, contractEntity.Id)
		if err != nil {
			spans.TraceError(err)
			return "", err
		}
		if opportunityDbNode != nil {
			opportunityEntity := neo4jmapper.MapDbNodeToOpportunityEntity(opportunityDbNode)
			if opportunityEntity.RenewalDetails.RenewedAt != nil && opportunityEntity.RenewalDetails.RenewedAt.Before(now) {
				spans.LogKV("result.status", neo4jenum.ContractStatusLive.String())
				return neo4jenum.ContractStatusOutOfContract.String(), nil
			}
		}
	}

	// Otherwise, the contract is considered Live.
	spans.LogKV("result.status", neo4jenum.ContractStatusLive.String())
	return neo4jenum.ContractStatusLive.String(), nil
}

func (s *contractService) updateOrganizationRelationship(ctx context.Context, tenant, contractId string, statusChanged bool) error {
	spans, ctx := telemetry.StartServiceSpan(ctx, "ContractService.UpdateOrganizationRelationship")
	defer spans.Finish()
	spans.LogKV("contractId", contractId, "statusChanged", statusChanged)

	if !statusChanged {
		return nil
	}

	// get contract
	contractDbNode, err := s.neo4j.ContractReadRepository.GetContractById(ctx, tenant, contractId)
	if err != nil {
		spans.TraceError(err)
		s.log.Errorf("Error while getting contract %s: %s", contractId, err.Error())
		return err
	}
	contractEntity := neo4jmapper.MapDbNodeToContractEntity(contractDbNode)

	// get organization for contract
	organizationDbNode, err := s.neo4j.OrganizationReadRepository.GetOrganizationByContractId(ctx, tenant, contractId)
	if err != nil {
		spans.TraceError(err)
		s.log.Errorf("Error while getting organization for contract %s: %s", contractId, err.Error())
		return err
	}
	orgEntity := neo4jmapper.MapDbNodeToOrganizationEntity(organizationDbNode)

	// get all contracts for organization
	orgContracts, err := s.neo4j.ContractReadRepository.GetContractsForOrganizations(ctx, tenant, []string{orgEntity.ID})
	if err != nil {
		spans.TraceError(err)
		s.log.Errorf("Error while getting contracts for organization %s: %s", orgEntity.ID, err.Error())
		return err
	}
	orgContractEntities := []neo4jentity.ContractEntity{}
	for _, orgContract := range orgContracts {
		orgContractEntities = append(orgContractEntities, *neo4jmapper.MapDbNodeToContractEntity(orgContract.Node))
	}

	if contractEntity.ContractStatus == neo4jenum.ContractStatusEnded {
		// check no other contract is active
		activeContractFound := false
		for _, orgContract := range orgContractEntities {
			if orgContract.ContractStatus != neo4jenum.ContractStatusDraft && orgContract.ContractStatus != neo4jenum.ContractStatusEnded {
				activeContractFound = true
				break
			}
		}

		if !activeContractFound {
			_, err = s.organization.Save(ctx, nil, &orgEntity.ID, data_fields.OrganizationFields{
				Relationship: utils.ToPtr(neo4jenum.OrganizationRelationshipFormerCustomer),
				Stage:        utils.ToPtr(neo4jenum.Target),
			})
			if err != nil {
				spans.TraceError(err)
				s.log.Errorf("UpdateOrganization failed: %s", err.Error())
				return errors.Wrap(err, "UpdateOrganization")
			}
			err = s.organization.UpdateDerivedData(ctx, orgEntity.ID)
			if err != nil {
				spans.TraceError(err)
				s.log.Errorf("UpdateDerivedData failed: %s", err.Error())
			}
		}
	}

	return nil
}

func (s *contractService) startOnboardingIfEligible(ctx context.Context, tenant, contractId string, spans telemetry.Spans) {
	// TODO temporary not eligible for all contracts
	return

	contractDbNode, err := s.neo4j.ContractReadRepository.GetContractById(ctx, tenant, contractId)
	if err != nil {
		spans.TraceError(err)
		return
	}
	if contractDbNode == nil {
		return
	}
	contractEntity := neo4jmapper.MapDbNodeToContractEntity(contractDbNode)

	if contractEntity.IsEligibleToStartOnboarding() {
		organizationDbNode, err := s.neo4j.OrganizationReadRepository.GetOrganizationByContractId(ctx, tenant, contractEntity.Id)
		if err != nil {
			spans.TraceError(err)
			s.log.Errorf("Error while getting organization for contract %s: %s", contractEntity.Id, err.Error())
			return
		}
		if organizationDbNode == nil {
			return
		}
		organization := neo4jmapper.MapDbNodeToOrganizationEntity(organizationDbNode)
		err = s.organization.UpdateOnboardingStatus(ctx, nil, organization.ID, data_fields.OrganizationOnboardingStatusFields{
			CausedByContractId: &contractEntity.Id,
			Status:             utils.ToPtr(neo4jenum.OnboardingStatusNotStarted),
		})
		if err != nil {
			spans.TraceError(err)
			s.log.Errorf("UpdateOnboardingStatus failed: %v", err.Error())
		}
	}
}

func (s *contractService) UpdateActiveRenewalOpportunityRenewDateAndArr(ctx context.Context, tenant, contractId string) error {
	spans, ctx := telemetry.StartServiceSpan(ctx, "ContractService.updateActiveRenewalOpportunityRenewDateAndArr")
	defer spans.Finish()
	spans.LogKV("contractId", contractId)

	contract, renewalOpportunity, done := s.assertContractAndRenewalOpportunity(ctx, tenant, contractId)
	if done {
		return nil
	}

	err := s.updateRenewalOpportunityRenewedAt(ctx, tenant, contract, renewalOpportunity)
	if err != nil {
		spans.TraceError(err)
		return nil
	}
	err = s.updateRenewalArr(ctx, tenant, contract, renewalOpportunity, *spans)
	if err != nil {
		spans.TraceError(err)
		return nil
	}
	return nil
}

func (s *contractService) UpdateActiveRenewalOpportunityArr(ctx context.Context, contractId string) error {
	spans, ctx := telemetry.StartServiceSpan(ctx, "ContractService.UpdateActiveRenewalOpportunityArr")
	defer spans.Finish()

	spans.LogKV("contractId", contractId)

	tenant := common.GetTenantFromContext(ctx)

	contract, renewalOpportunity, done := s.assertContractAndRenewalOpportunity(ctx, tenant, contractId)
	if done {
		return nil
	}
	err := s.updateRenewalArr(ctx, tenant, contract, renewalOpportunity, *spans)
	if err != nil {
		spans.TraceError(err)
		return nil
	}
	return nil
}

func (s *contractService) assertContractAndRenewalOpportunity(ctx context.Context, tenant, contractId string) (*neo4jentity.ContractEntity, *neo4jentity.OpportunityEntity, bool) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "ContractService.assertContractAndRenewalOpportunity")
	defer spans.Finish()
	spans.LogKV("contractId", contractId)

	contractDbNode, err := s.neo4j.ContractReadRepository.GetContractById(ctx, tenant, contractId)
	if err != nil {
		spans.TraceError(err)
		s.log.Errorf("Error while getting contract %s: %s", contractId, err.Error())
		return nil, nil, true
	}
	contract := neo4jmapper.MapDbNodeToContractEntity(contractDbNode)

	// if contract is not frequency based, return
	if contract.LengthInMonths == 0 {
		return nil, nil, true
	}

	currentRenewalOpportunityDbNode, err := s.neo4j.OpportunityReadRepository.GetActiveRenewalOpportunityForContract(ctx, tenant, contractId)
	if err != nil {
		spans.TraceError(err)
		s.log.Errorf("Error while getting renewal opportunity for contract %s: %s", contractId, err.Error())
		return nil, nil, true
	}

	// if there is no renewal opportunity, create one
	if currentRenewalOpportunityDbNode == nil {
		if !contract.IsEnded() {
			_, err = s.opportunity.CreateRenewalOpportunity(ctx, nil, &data_fields.OpportunityFields{
				ContractId: &contractId,
			})
			if err != nil {
				spans.TraceError(err)
				s.log.Errorf("CreateRenewalOpportunity command failed: %v", err.Error())
				return nil, nil, true
			}
			spans.LogFields(log.Bool("renewal opportunity create requested", true))
		}
		return nil, nil, true
	}

	currentRenewalOpportunity := neo4jmapper.MapDbNodeToOpportunityEntity(currentRenewalOpportunityDbNode)

	return contract, currentRenewalOpportunity, false
}

func (s *contractService) updateRenewalOpportunityRenewedAt(ctx context.Context, tenant string, contractEntity *neo4jentity.ContractEntity, renewalOpportunityEntity *neo4jentity.OpportunityEntity) error {
	spans, ctx := telemetry.StartServiceSpan(ctx, "ContractService.updateRenewalOpportunityRenewedAt")
	defer spans.Finish()

	if renewalOpportunityEntity == nil {
		err := fmt.Errorf("renewalOpportunityEntity is nil")
		spans.TraceError(err)
		return nil
	}

	// IF contract already ended, close the renewal opportunity
	if contractEntity.IsEnded() {
		err := s.opportunity.CloseLost(ctx, nil, tenant, renewalOpportunityEntity.Id)
		if err != nil {
			spans.TraceError(err)
			s.log.Errorf("CloseLooseOpportunity failed: %s", err.Error())
			return errors.Wrap(err, "CloseLooseOpportunity")
		}
		return nil
	}

	// Choose starting date for renewal calculation
	if renewalOpportunityEntity.RenewalDetails.RenewedAt != nil {
		return nil
	}
	startRenewalDateCalculation := contractEntity.ServiceStartedAt
	previousClosedWonRenewalDbNode, err := s.neo4j.OpportunityReadRepository.GetPreviousClosedWonRenewalOpportunityForContract(ctx, tenant, contractEntity.Id)
	if err != nil {
		spans.TraceError(err)
		return nil
	}
	if previousClosedWonRenewalDbNode != nil {
		previousRenewalOpportunityEntity := neo4jmapper.MapDbNodeToOpportunityEntity(previousClosedWonRenewalDbNode)
		if previousRenewalOpportunityEntity.RenewalDetails.RenewedAt != nil {
			startRenewalDateCalculation = previousRenewalOpportunityEntity.RenewalDetails.RenewedAt
		}
	}
	spans.LogFields(log.Object("startRenewalDateCalculation", startRenewalDateCalculation))

	// Calculate until first future date if auto-renew is enabled or renewal is approved
	calculateUntilFirstFutureDate := contractEntity.AutoRenew
	spans.LogFields(log.Bool("calculateUntilFirstFutureDate", calculateUntilFirstFutureDate))

	renewedAt := calculateNextCycleDate(startRenewalDateCalculation, contractEntity.LengthInMonths, calculateUntilFirstFutureDate)
	spans.LogFields(log.Object("result.renewedAt", renewedAt))
	if !utils.IsEqualTimePtr(renewedAt, renewalOpportunityEntity.RenewalDetails.RenewedAt) {
		_, err = s.opportunity.Save(ctx, nil, &renewalOpportunityEntity.Id, &data_fields.OpportunityFields{
			RenewedAt: renewedAt,
		})
		if err != nil {
			spans.TraceError(err)
			s.log.Errorf("UpdateRenewalOpportunityNextCycleDate failed: %s", err.Error())
			return err
		}
	}

	return nil
}

func (s *contractService) updateRenewalArr(ctx context.Context, tenant string, contract *neo4jentity.ContractEntity, renewalOpportunity *neo4jentity.OpportunityEntity, spans telemetry.Spans) error {
	// if contract already ended, return
	if contract.IsEnded() {
		spans.LogFields(log.Bool("contract ended", true))
		return nil
	}

	maxArr, err := s.calculateMaxArr(ctx, tenant, contract, renewalOpportunity, spans)
	if err != nil {
		spans.TraceError(err)
		s.log.Errorf("Error while calculating ARR for contract %s: %s", contract.Id, err.Error())
		return nil
	}
	// adjust with likelihood
	currentArr := calculateCurrentArrByAdjustedRate(maxArr, renewalOpportunity.RenewalDetails.RenewalAdjustedRate)

	_, err = s.opportunity.Save(ctx, nil, &renewalOpportunity.Id, &data_fields.OpportunityFields{
		Amount:    &currentArr,
		MaxAmount: &maxArr,
	})
	if err != nil {
		spans.TraceError(err)
		s.log.Errorf("UpdateOpportunity failed: %s", err.Error())
		return err
	}

	return nil
}

func (s *contractService) calculateMaxArr(ctx context.Context, tenant string, contract *neo4jentity.ContractEntity, renewalOpportunity *neo4jentity.OpportunityEntity, spans telemetry.Spans) (float64, error) {
	var arr float64

	// Fetch service line items for the contract from the database
	sliDbNodes, err := s.neo4j.ServiceLineItemReadRepository.GetServiceLineItemsForContract(ctx, tenant, contract.Id)
	if err != nil {
		spans.TraceError(err)
		return 0, err
	}
	serviceLineItems := neo4jentity.ServiceLineItemEntities{}
	for _, sliDbNode := range sliDbNodes {
		sli := neo4jmapper.MapDbNodeToServiceLineItemEntity(sliDbNode)
		serviceLineItems = append(serviceLineItems, *sli)
	}

	spans.LogKV("service line items count", len(serviceLineItems))
	for _, sli := range serviceLineItems {
		if sli.IsEnded() {
			spans.LogFields(log.Bool(fmt.Sprintf("service line item {%s} ended", sli.ID), true))
			continue
		}
		spans.LogFields(log.Object(fmt.Sprintf("service line item {%s}:", sli.ID), sli))
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
		spans.LogKV(fmt.Sprintf("service line item {%s} added ARR value:", sli.ID), annualPrice)
		// Add to total ARR
		arr += annualPrice
	}

	// Adjust with end date
	if contract.EndedAt != nil {
		spans.LogFields(log.Bool("ARR prorated with contract end date", true))
		arr = prorateArr(arr, monthsUntilContractEnd(utils.Now(), *contract.EndedAt))
	}

	return utils.RoundHalfUpFloat64(arr, 2), nil
}

func calculateCurrentArrByAdjustedRate(maxAmount float64, rate int64) float64 {
	if rate == 0 {
		return 0
	} else if rate == 100 {
		return maxAmount
	}
	return utils.RoundHalfUpFloat64(maxAmount*float64(rate)/100, 2)
}

type ActionStatusMetadata struct {
	Status       string `json:"status"`
	ContractName string `json:"contract-name"`
	Comment      string `json:"comment"`
}

func (s *contractService) createActionForStatusChange(ctx context.Context, tenant, contractId, status, contractName string) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "ContractService.createActionForStatusChange")
	defer spans.Finish()
	var name string
	spans.LogKV("contractId", contractId, "status", status, "contractName", contractName)

	// if status is not one of the predefined statuses, return
	if status != string(neo4jenum.ContractStatusLive) &&
		status != string(neo4jenum.ContractStatusEnded) &&
		status != string(neo4jenum.ContractStatusOutOfContract) {
		return
	}

	if contractName != "" {
		name = contractName
	} else {
		name = "Unnamed contract"
	}
	actionStatusMetadata := ActionStatusMetadata{
		Status:       status,
		ContractName: name,
		Comment:      name + " is now " + status,
	}

	message := ""
	switch status {
	case string(neo4jenum.ContractStatusLive):
		message = contractName + " is now live"
		actionStatusMetadata.Comment = contractName + " is now live"
	case string(neo4jenum.ContractStatusEnded):
		message = contractName + " has ended"
		actionStatusMetadata.Comment = contractName + " has ended"
	case string(neo4jenum.ContractStatusOutOfContract):
		message = contractName + " is now out of contract"
		actionStatusMetadata.Comment = contractName + " is now out of contract"
	}
	metadata, err := utils.ToJson(actionStatusMetadata)
	if err != nil {
		spans.TraceError(err)
		s.log.Errorf("Failed creating status update action for contract %s: %s", contractId, err.Error())
		return
	}
	_, err = s.neo4j.ActionWriteRepository.Create(ctx, tenant, contractId, model.CONTRACT, enum.ActionContractStatusUpdated, message, metadata, utils.Now(), common.GetAppSourceFromContext(ctx))
	if err != nil {
		spans.TraceError(err)
		s.log.Errorf("Failed creating status update action for contract %s: %s", contractId, err.Error())
	}
}

func (s *contractService) UpdateActiveRenewalOpportunityLikelihood(ctx context.Context, tenant, contractId string) error {
	spans, ctx := telemetry.StartServiceSpan(ctx, "ContractService.UpdateActiveRenewalOpportunityLikelihood")
	defer spans.Finish()
	spans.LogKV("contractId", contractId)

	opportunityDbNode, err := s.neo4j.OpportunityReadRepository.GetActiveRenewalOpportunityForContract(ctx, tenant, contractId)
	if err != nil {
		spans.TraceError(err)
		s.log.Errorf("Error while getting renewal opportunity for contract %s: %s", contractId, err.Error())
		return err
	}
	if opportunityDbNode == nil {
		s.log.Infof("No open renewal opportunity found for contract %s", contractId)
		return nil
	}
	contractDbNode, err := s.neo4j.ContractReadRepository.GetContractById(ctx, tenant, contractId)
	if err != nil {
		spans.TraceError(err)
		s.log.Errorf("Error while getting contract %s: %s", contractId, err.Error())
		return err
	}
	contractEntity := neo4jmapper.MapDbNodeToContractEntity(contractDbNode)
	opportunityEntity := neo4jmapper.MapDbNodeToOpportunityEntity(opportunityDbNode)

	var renewalLikelihood neo4jenum.RenewalLikelihood
	renewalAdjustedRate := opportunityEntity.RenewalDetails.RenewalAdjustedRate
	if contractEntity.EndedAt != nil &&
		opportunityEntity.RenewalDetails.RenewalLikelihood != neo4jenum.RenewalLikelihoodZero &&
		opportunityEntity.RenewalDetails.RenewedAt != nil &&
		contractEntity.EndedAt.Before(*opportunityEntity.RenewalDetails.RenewedAt) {
		// check if likelihood should be set to Zero
		renewalLikelihood = neo4jenum.RenewalLikelihoodZero
		renewalAdjustedRate = int64(0)
	} else if opportunityEntity.RenewalDetails.RenewalLikelihood == neo4jenum.RenewalLikelihoodZero &&
		opportunityEntity.RenewalDetails.RenewedAt != nil &&
		(contractEntity.EndedAt == nil || contractEntity.EndedAt.After(*opportunityEntity.RenewalDetails.RenewedAt)) {
		// check if likelihood should be set to Medium
		renewalLikelihood = neo4jenum.RenewalLikelihoodMedium
		renewalAdjustedRate = int64(50)
	}

	if renewalLikelihood != "" {
		_, err = s.opportunity.Save(ctx, nil, &opportunityEntity.Id, &data_fields.OpportunityFields{
			RenewalLikelihood:   &renewalLikelihood,
			RenewalAdjustedRate: &renewalAdjustedRate,
		})
		if err != nil {
			spans.TraceError(err)
			s.log.Errorf("UpdateRenewalOpportunity failed: %s", err.Error())
			return errors.Wrap(err, "UpdateRenewalOpportunity")
		}
	}

	return nil
}

func calculateNextCycleDate(from *time.Time, lengthInMonths int64, calculateUntilFirstFutureDate bool) *time.Time {
	if from == nil || lengthInMonths <= 0 {
		return nil
	}

	renewalCycleNext := *from
	for {
		renewalCycleNext = renewalCycleNext.AddDate(0, int(lengthInMonths), 0)
		// Break the loop either when the next cycle date is in the future
		// or if we are not calculating until the first future date.
		if renewalCycleNext.After(utils.Now()) || !calculateUntilFirstFutureDate {
			break
		}
	}
	return &renewalCycleNext
}

func prorateArr(arr float64, monthsRemaining int) float64 {
	if monthsRemaining > 12 {
		return arr
	}
	monthlyRate := arr / 12
	return utils.RoundHalfUpFloat64(monthlyRate*float64(monthsRemaining), 2)
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

func (s *contractService) RefreshContractStatus(ctx context.Context, contractId string) error {
	spans, ctx := telemetry.StartServiceSpan(ctx, "ContractService.RefreshContractStatus")
	defer spans.Finish()

	spans.TagEntity(contractId)

	// validate tenant
	err := common.ValidateTenant(ctx)
	if err != nil {
		spans.TraceError(err)
		return err
	}
	tenant := common.GetTenantFromContext(ctx)

	status, statusChanged, err := s.updateStatus(ctx, tenant, contractId)
	if err != nil {
		spans.TraceError(err)
		s.log.Errorf("Error while updating contract %s status: %s", contractId, err.Error())
		return err
	}
	spans.LogKV("result.status", status)
	spans.LogFields(log.Bool("result.statusChanged", statusChanged))

	if statusChanged {
		err = s.updateOrganizationRelationship(ctx, tenant, contractId, statusChanged)
		if err != nil {
			spans.TraceError(err)
			s.log.Errorf("Error while updating organization relationship for contract %s: %s", contractId, err.Error())
		}

		contractDbNode, err := s.neo4j.ContractReadRepository.GetContractById(ctx, tenant, contractId)
		if err != nil {
			spans.TraceError(err)
			return err
		}
		contractEntity := neo4jmapper.MapDbNodeToContractEntity(contractDbNode)
		s.createActionForStatusChange(ctx, tenant, contractId, status, contractEntity.Name)

		s.startOnboardingIfEligible(ctx, tenant, contractId, *spans)
		s.events.Publisher.PublishNotification(ctx, tenant, contractId, model.CONTRACT, utils.NewEventCompletedDetails().WithUpdate())
	}

	if status == neo4jenum.ContractStatusEnded.String() {
		err = s.UpdateActiveRenewalOpportunityRenewDateAndArr(ctx, tenant, contractId)
		if err != nil {
			spans.TraceError(err)
			s.log.Errorf("error while updating renewal opportunity for contract %s: %s", contractId, err.Error())
		}

		err = s.neo4j.InvoiceWriteRepository.DeletePreviewCycleInvoices(ctx, tenant, contractId, "")
		if err != nil {
			spans.TraceError(err)
			s.log.Errorf("Error while deleting preview invoice for contract %s: %s", contractId, err.Error())
		}
	}

	return nil
}

func (s *contractService) RecalculateContractLtv(ctx context.Context, contractId string) error {
	spans, ctx := telemetry.StartServiceSpan(ctx, "ContractService.RecalculateContractLtv")
	defer spans.Finish()

	spans.TagEntity(contractId)

	// validate tenant
	err := common.ValidateTenant(ctx)
	if err != nil {
		spans.TraceError(err)
		return err
	}
	tenant := common.GetTenantFromContext(ctx)

	// Get all invoices for the contract
	invoiceDbNodes, err := s.neo4j.InvoiceReadRepository.GetAllForContracts(ctx, tenant, []string{contractId})
	if err != nil {
		spans.TraceError(err)
		return err
	}

	contractDbNode, err := s.neo4j.ContractReadRepository.GetContractById(ctx, tenant, contractId)
	if err != nil {
		spans.TraceError(err)
		return err
	}
	contractEntity := neo4jmapper.MapDbNodeToContractEntity(contractDbNode)

	// Calculate total LTV from actual invoices
	ltv := 0.0
	for _, invoiceDbNode := range invoiceDbNodes {
		invoiceEntity := neo4jmapper.MapDbNodeToInvoiceEntity(invoiceDbNode.Node)

		// Skip dry run invoices and voided invoices
		if invoiceEntity.DryRun || invoiceEntity.Status == neo4jenum.InvoiceStatusVoid {
			continue
		}

		// Add invoice amount to total LTV
		ltv += invoiceEntity.Amount
	}

	truncatedLtv := utils.TruncateFloat64(ltv, 2)

	if contractEntity.Ltv != truncatedLtv {
		err = s.neo4j.ContractWriteRepository.SetLtv(ctx, tenant, contractId, truncatedLtv)
		if err != nil {
			spans.TraceError(err)
			s.log.Errorf("Error while updating contract %s ltv: %s", contractId, err.Error())
			return err
		}

		// get organization for contract
		organizationDbNode, err := s.neo4j.OrganizationReadRepository.GetOrganizationByContractId(ctx, tenant, contractId)
		if err != nil {
			spans.TraceError(err)
			s.log.Errorf("Error while getting organization for contract %s: %s", contractId, err.Error())
			return nil
		}
		organizationEntity := neo4jmapper.MapDbNodeToOrganizationEntity(organizationDbNode)

		// request organization ltv refresh
		if organizationEntity.ID != "" {
			err = s.organization.UpdateDerivedData(ctx, organizationEntity.ID)
			if err != nil {
				spans.TraceError(err)
				s.log.Errorf("Error while updating organization %s ltv: %s", organizationEntity.ID, err.Error())
			}
		}

		err = s.events.Publisher.PublishFanoutEvent(ctx, contractId, model.CONTRACT, dto.UpdateContract{Ltv: &truncatedLtv})
		if err != nil {
			spans.TraceError(errors.Wrap(err, "unable to publish message UpdateContract"))
		}

		s.events.Publisher.PublishNotification(ctx, tenant, contractId, model.CONTRACT, utils.NewEventCompletedDetails().WithUpdate())
	}
	return nil
}

func (s *contractService) GetContractsForOrganizations(ctx context.Context, organizationIDs []string) (*neo4jentity.ContractEntities, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "ContractService.GetContractsForOrganizations")
	defer spans.Finish()

	spans.LogFields(log.Object("organizationIDs", organizationIDs))

	contracts, err := s.neo4j.ContractReadRepository.GetContractsForOrganizations(ctx, common.GetTenantFromContext(ctx), organizationIDs)
	if err != nil {
		return nil, err
	}
	contractEntities := make(neo4jentity.ContractEntities, 0, len(contracts))
	for _, v := range contracts {
		contractEntity := neo4jmapper.MapDbNodeToContractEntity(v.Node)
		contractEntity.DataloaderKey = v.LinkedNodeId
		contractEntities = append(contractEntities, *contractEntity)
	}
	return &contractEntities, nil
}

func (s *contractService) GetContractForInvoice(ctx context.Context, invoiceId string) (*neo4jentity.ContractEntity, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "ContractService.GetContractForInvoice")
	defer spans.Finish()

	spans.LogKV("invoiceId", invoiceId)

	contractDbNode, err := s.neo4j.ContractReadRepository.GetContractForInvoice(ctx, common.GetTenantFromContext(ctx), invoiceId)
	if err != nil {
		return nil, err
	}
	if contractDbNode == nil {
		err = errors.New("contract not found")
		spans.TraceError(err)
		return nil, err
	}
	return neo4jmapper.MapDbNodeToContractEntity(contractDbNode), nil
}
