package service

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/data_fields"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/dto"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	commonModel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/constants"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"
	organizationpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/organization"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type ActionLikelihoodMetadata struct {
	Likelihood string `json:"likelihood"`
	Reason     string `json:"reason"`
}

type OpportunityService interface {
	GetById(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, opportunityId string) (*neo4jentity.OpportunityEntity, error)
	GetOpportunitiesForContracts(ctx context.Context, tenant string, contractIds []string) (*neo4jentity.OpportunityEntities, error)
	GetOpportunitiesForOrganizations(ctx context.Context, tenant string, organizationIds []string) (*neo4jentity.OpportunityEntities, error)
	GetPaginatedOrganizationOpportunities(ctx context.Context, tenant string, page int, limit int) (*utils.Pagination, error)

	Save(ctx context.Context, txWithPostCommit *utils.TxWithPostCommit, opportunityId *string, input *data_fields.OpportunityFields) (string, error)
	CreateRenewalOpportunity(ctx context.Context, txWithPostCommit *utils.TxWithPostCommit, input *data_fields.OpportunityFields) (string, error)
	CloseWon(ctx context.Context, txWithPostCommit *utils.TxWithPostCommit, tenant, opportunityId string) error
	CloseLost(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, opportunityId string) error
	Archive(ctx context.Context, tenant, opportunityId string) error

	RolloutRenewalOpportunity(ctx context.Context, contractId string) error
}

type opportunityService struct {
	log      logger.Logger
	services *Services
}

func NewOpportunityService(log logger.Logger, services *Services) OpportunityService {
	return &opportunityService{
		log:      log,
		services: services,
	}
}

func (s *opportunityService) GetById(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, opportunityId string) (*neo4jentity.OpportunityEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OpportunityService.GetById")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogKV("opportunityId", opportunityId)

	if common.GetTenantFromContext(ctx) == "" {
		tracing.TagTenant(span, tenant)
	}

	if opportunityDbNode, err := s.services.Neo4jRepositories.OpportunityReadRepository.GetOpportunityById(ctx, tx, tenant, opportunityId); err != nil {
		tracing.TraceErr(span, err)
		wrappedErr := errors.Wrap(err, fmt.Sprintf("opportunity with id {%s} not found", opportunityId))
		return nil, wrappedErr
	} else {
		return neo4jmapper.MapDbNodeToOpportunityEntity(opportunityDbNode), nil
	}
}

func (s *opportunityService) GetOpportunitiesForContracts(ctx context.Context, tenant string, contractIDs []string) (*neo4jentity.OpportunityEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OpportunityService.GetOpportunitiesForContracts")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("contractIDs", contractIDs))

	opportunities, err := s.services.Neo4jRepositories.OpportunityReadRepository.GetForContracts(ctx, tenant, contractIDs)
	if err != nil {
		return nil, err
	}
	opportunityEntities := make(neo4jentity.OpportunityEntities, 0, len(opportunities))
	for _, v := range opportunities {
		opportunityEntity := neo4jmapper.MapDbNodeToOpportunityEntity(v.Node)
		opportunityEntity.DataloaderKey = v.LinkedNodeId
		opportunityEntities = append(opportunityEntities, *opportunityEntity)
	}
	return &opportunityEntities, nil
}

func (s *opportunityService) GetOpportunitiesForOrganizations(ctx context.Context, tenant string, organizationIds []string) (*neo4jentity.OpportunityEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OpportunityService.GetOpportunitiesForOrganizations")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("organizationIds", organizationIds))

	opportunities, err := s.services.Neo4jRepositories.OpportunityReadRepository.GetForOrganizations(ctx, tenant, organizationIds)
	if err != nil {
		return nil, err
	}
	opportunityEntities := make(neo4jentity.OpportunityEntities, 0, len(opportunities))
	for _, v := range opportunities {
		opportunityEntity := neo4jmapper.MapDbNodeToOpportunityEntity(v.Node)
		opportunityEntity.DataloaderKey = v.LinkedNodeId
		opportunityEntities = append(opportunityEntities, *opportunityEntity)
	}
	return &opportunityEntities, nil
}

func (s *opportunityService) GetPaginatedOrganizationOpportunities(ctx context.Context, tenant string, page int, limit int) (*utils.Pagination, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OpportunityService.GetPaginatedOrganizationOpportunities")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Int("page", page), log.Int("limit", limit))

	var paginatedResult = utils.Pagination{
		Limit: limit,
		Page:  page,
	}

	dbNodesWithTotalCount, err := s.services.Neo4jRepositories.OpportunityReadRepository.GetPaginatedOpportunitiesLinkedToAnOrganization(ctx, tenant, paginatedResult.GetSkip(), paginatedResult.GetLimit())
	if err != nil {
		return nil, err
	}
	paginatedResult.SetTotalRows(dbNodesWithTotalCount.Count)

	opportunities := neo4jentity.OpportunityEntities{}

	for _, v := range dbNodesWithTotalCount.Nodes {
		opportunities = append(opportunities, *neo4jmapper.MapDbNodeToOpportunityEntity(v))
	}
	paginatedResult.SetRows(&opportunities)
	return &paginatedResult, nil
}

func (s *opportunityService) CreateRenewalOpportunity(ctx context.Context, txWithPostCommit *utils.TxWithPostCommit, input *data_fields.OpportunityFields) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OpportunityService.CreateRenewalOpportunity")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	input.InternalType = utils.ToPtr(neo4jenum.OpportunityInternalTypeRenewal)

	return s.Save(ctx, txWithPostCommit, nil, input)
}

func (s *opportunityService) Save(ctx context.Context, txWithPostCommit *utils.TxWithPostCommit, id *string, input *data_fields.OpportunityFields) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OpportunityService.Save")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("input", input))

	// validate tenant
	err := common.ValidateTenant(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return "", err
	}
	tenant := common.GetTenantFromContext(ctx)

	var existing *neo4jentity.OpportunityEntity
	createFlow := false
	opportunityId := ""
	likelihoodChanged, adjustedRateChanged, amountChanged := false, false, false

	if utils.IfNotNilString(id) == "" {
		if utils.IfNotNilString(input.OrganizationId) == "" && !input.IsRenewal() {
			err := fmt.Errorf("(OpportunityService.Save) organizationId and opportunityId and contractId are empty")
			tracing.TraceErr(span, err)
			return "", err
		}
		if utils.IfNotNilString(input.ContractId) == "" && input.IsRenewal() {
			err := fmt.Errorf("(OpportunityService.Save) contractId and opportunityId and contractId are empty")
			tracing.TraceErr(span, err)
			return "", err
		}
	}

	if input.OrganizationId != nil {
		existsById, err := s.services.Neo4jRepositories.CommonReadRepository.ExistsById(ctx, tenant, *input.OrganizationId, commonModel.NodeLabelOrganization)
		if err != nil {
			tracing.TraceErr(span, err)
			return "", err
		}
		if !existsById {
			err := fmt.Errorf("(OpportunityService.Save) organization with id {%s} not found", *input.OrganizationId)
			tracing.TraceErr(span, err)
			return "", err
		}
	}

	if utils.IfNotNilString(id) == "" {
		span.LogKV("flow", "create")

		// set default values if not provided
		if input.CreatedAt == nil || input.CreatedAt.IsZero() {
			input.CreatedAt = utils.NowPtr()
		}
		if utils.IfNotNilString(input.Source) == "" {
			input.Source = utils.StringPtr(neo4jentity.DataSourceOpenline.String())
		}
		if utils.IfNotNilString(input.AppSource) == "" {
			input.AppSource = utils.StringPtr(common.GetAppSourceFromContext(ctx))
		}
		if utils.IfNotNilString(input.InternalType) == "" {
			input.InternalType = utils.ToPtr(neo4jenum.OpportunityInternalTypeNBO)
		}
		if utils.IfNotNilString(input.InternalStage) == "" {
			input.InternalStage = utils.StringPtr(neo4jenum.OpportunityInternalStageOpen.String())
		}

		if input.Currency == nil || utils.IfNotNilString(input.Currency.String()) == "" {
			tenantSettings, err := s.services.TenantSettingsService.GetTenantSettings(ctx)
			if err != nil {
				tracing.TraceErr(span, err)
				return "", err
			}
			input.Currency = utils.ToPtr(tenantSettings.BaseCurrency)
		}

		// default values for renewal opportunity
		if input.IsRenewal() {
			if input.RenewalLikelihood == nil {
				input.RenewalLikelihood = utils.ToPtr(neo4jenum.RenewalLikelihoodHigh)
			}
			if *input.RenewalLikelihood == neo4jenum.RenewalLikelihoodHigh && utils.IfNotNilInt64(input.RenewalAdjustedRate) == 0 {
				input.RenewalAdjustedRate = utils.ToPtr(int64(100))
			}
			if utils.IfNotNilInt64(input.RenewalAdjustedRate) < 0 {
				input.RenewalAdjustedRate = utils.ToPtr(int64(0))
			} else if utils.IfNotNilInt64(input.RenewalAdjustedRate) > 100 {
				input.RenewalAdjustedRate = utils.ToPtr(int64(100))
			}
			if input.RenewalApproved == nil {
				input.RenewalApproved = utils.BoolPtr(false)
			}
		}

		// validate input
		if input.IsRenewal() {
			// check if active renewal opportunity already exists for this contract
			opportunityDbNode, err := s.services.Neo4jRepositories.OpportunityReadRepository.GetActiveRenewalOpportunityForContract(ctx, tenant, utils.IfNotNilString(input.ContractId))
			if err != nil {
				tracing.TraceErr(span, err)
				return "", err
			}
			if opportunityDbNode != nil {
				opportunity := neo4jmapper.MapDbNodeToOpportunityEntity(opportunityDbNode)
				if opportunity.RenewalDetails.RenewedAt != nil && opportunity.RenewalDetails.RenewedAt.After(utils.Now()) {
					span.LogFields(log.String("result", "active renewal opportunity already exists, skip creation"))
					s.log.Infof("active renewal opportunity already exists for contract %s", utils.IfNotNilString(input.ContractId))
					return "", nil
				}
			}
		}

		generatedId, err := s.services.Neo4jRepositories.CommonReadRepository.GenerateId(ctx, tenant, commonModel.NodeLabelOpportunity)
		if err != nil {
			tracing.TraceErr(span, err)
			return "", err
		}
		opportunityId = generatedId

	} else {
		span.LogKV("flow", "update")
		opportunityId = utils.IfNotNilString(id)

		existing, err = s.GetById(ctx, nil, tenant, opportunityId)
		if err != nil {
			tracing.TraceErr(span, err)
			return "", err
		}
		if existing == nil {
			err := fmt.Errorf("(OpportunityService.Save) opportunity with id {%s} not found", opportunityId)
			tracing.TraceErr(span, err)
			return "", err
		}

		// Changing external stage should set internal stage back to OPEN
		if utils.IfNotNilString(input.ExternalStage) != "" && existing.ExternalStage != utils.IfNotNilString(input.ExternalStage) && existing.InternalStage != neo4jenum.OpportunityInternalStageOpen {
			input.InternalStage = utils.StringPtr(neo4jenum.OpportunityInternalStageOpen.String())
		}

		if existing.IsRenewal() {
			if input.RenewalLikelihood != nil && input.RenewalLikelihood.ToV2() == "" {
				input.RenewalLikelihood = utils.ToPtr(neo4jenum.RenewalLikelihoodHigh)
				input.RenewalAdjustedRate = utils.ToPtr(int64(100))
			}
			if utils.IfNotNilInt64(input.RenewalAdjustedRate) < 0 {
				input.RenewalAdjustedRate = utils.ToPtr(int64(0))
			} else if utils.IfNotNilInt64(input.RenewalAdjustedRate) > 100 {
				input.RenewalAdjustedRate = utils.ToPtr(int64(100))
			}
		}

		likelihoodChanged = input.RenewalLikelihood != nil && existing.RenewalDetails.RenewalLikelihood.String() != utils.IfNotNilString(input.RenewalLikelihood)
		adjustedRateChanged = input.RenewalAdjustedRate != nil && existing.RenewalDetails.RenewalAdjustedRate != utils.IfNotNilInt64(input.RenewalAdjustedRate)
		amountChanged = input.Amount != nil && existing.Amount != utils.IfNotNilFloat64(input.Amount)
	}

	tracing.TagEntity(span, opportunityId)

	_, err = utils.ExecuteWriteInTransactionWithPostCommitActions(ctx, s.services.Neo4jRepositories.Neo4jDriver, s.services.Neo4jRepositories.Database, txWithPostCommit, func(txWithPostCommit *utils.TxWithPostCommit) (any, error) {
		newRenewalOpportunityCreated := false
		if createFlow && input.IsRenewal() {
			newRenewalOpportunityCreated, err = s.services.Neo4jRepositories.OpportunityWriteRepository.CreateRenewal(ctx, txWithPostCommit.Tx, tenant, opportunityId, *input)
			if err != nil {
				tracing.TraceErr(span, err)
				s.log.Errorf("error while saving renewal opportunity %s: %s", opportunityId, err.Error())
				return nil, err
			}
		} else if !createFlow && existing.IsRenewal() {
			err = s.services.Neo4jRepositories.OpportunityWriteRepository.UpdateRenewal(ctx, txWithPostCommit.Tx, tenant, opportunityId, *input)
			if err != nil {
				tracing.TraceErr(span, err)
				s.log.Errorf("error while updating renewal opportunity %s: %s", opportunityId, err.Error())
				return nil, err
			}
		} else {
			err = s.services.Neo4jRepositories.OpportunityWriteRepository.Save(ctx, txWithPostCommit.Tx, tenant, opportunityId, *input)
			if err != nil {
				tracing.TraceErr(span, err)
				return nil, err
			}
		}

		if utils.IfNotNilString(input.OrganizationId) != "" {
			err = s.services.Neo4jRepositories.CommonWriteRepository.Link(ctx, txWithPostCommit.Tx, tenant, repository.LinkDetails{
				FromEntityId:   *input.OrganizationId,
				FromEntityType: commonModel.ORGANIZATION,
				Relationship:   commonModel.HAS_OPPORTUNITY,
				ToEntityId:     opportunityId,
				ToEntityType:   commonModel.OPPORTUNITY,
			})
			if err != nil {
				tracing.TraceErr(span, err)
				return nil, err
			}
		}

		if input.OwnerId != nil {
			if utils.IfNotNilString(input.OwnerId) != "" {
				err = s.services.Neo4jRepositories.OpportunityWriteRepository.ReplaceOwner(ctx, txWithPostCommit.Tx, tenant, opportunityId, *input.OwnerId)
				if err != nil {
					tracing.TraceErr(span, err)
					return nil, err
				}
			} else {
				if existing != nil {
					err = s.services.Neo4jRepositories.OpportunityWriteRepository.RemoveOwner(ctx, txWithPostCommit.Tx, tenant, opportunityId)
					if err != nil {
						tracing.TraceErr(span, err)
						return nil, err
					}
				}
			}
		}

		if input.InternalStage != nil {
			if utils.IfNotNilString(input.InternalStage) == neo4jenum.OpportunityInternalStageClosedWon.String() {
				err := s.CloseWon(ctx, txWithPostCommit, tenant, opportunityId)
				if err != nil {
					tracing.TraceErr(span, err)
					return nil, err
				}
			} else if utils.IfNotNilString(input.InternalStage) == neo4jenum.OpportunityInternalStageClosedLost.String() {
				err := s.CloseLost(ctx, txWithPostCommit.Tx, tenant, opportunityId)
				if err != nil {
					tracing.TraceErr(span, err)
					return nil, err
				}
			}
		}

		// post create renewal actions
		txWithPostCommit.AddPostCommitAction(func(ctx context.Context) error {
			if newRenewalOpportunityCreated {
				err = s.services.ContractService.UpdateActiveRenewalOpportunityRenewDateAndArr(ctx, tenant, utils.IfNotNilString(input.ContractId))
				if err != nil {
					tracing.TraceErr(span, err)
					s.log.Errorf("error while updating renewal opportunity %s: %s", opportunityId, err.Error())
					return nil
				}
			}
			return nil
		})

		// post update renewal actions
		txWithPostCommit.AddPostCommitAction(func(ctx context.Context) error {
			if !createFlow && existing.IsRenewal() {
				// refresh organization renewal summary
				if likelihoodChanged {
					organizationDbNode, err := s.services.Neo4jRepositories.OrganizationReadRepository.GetOrganizationByOpportunityId(ctx, tenant, opportunityId)
					if err != nil {
						tracing.TraceErr(span, err)
						s.log.Errorf("error while getting organization for opportunity %s: %s", opportunityId, err.Error())
						return nil
					}
					if organizationDbNode == nil {
						return nil
					}
					organization := neo4jmapper.MapDbNodeToOrganizationEntity(organizationDbNode)

					ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
					_, err = utils.CallEventsPlatformGRPCWithRetry[*organizationpb.OrganizationIdGrpcResponse](func() (*organizationpb.OrganizationIdGrpcResponse, error) {
						return s.services.GrpcClients.OrganizationClient.RefreshRenewalSummary(ctx, &organizationpb.RefreshRenewalSummaryGrpcRequest{
							Tenant:         tenant,
							OrganizationId: organization.ID,
							AppSource:      common.GetAppSourceFromContext(ctx),
						})
					})
					if err != nil {
						tracing.TraceErr(span, err)
						s.log.Errorf("RefreshRenewalSummary failed: %v", err.Error())
					}
				}
				// likelihood change action
				if likelihoodChanged {
					contractDbNode, err := s.services.Neo4jRepositories.ContractReadRepository.GetContractByOpportunityId(ctx, tenant, opportunityId)
					if err != nil {
						tracing.TraceErr(span, err)
						s.log.Errorf("error while getting contract for opportunity %s: %s", opportunityId, err.Error())
					}
					if contractDbNode == nil {
						return nil
					}
					contractEntity := neo4jmapper.MapDbNodeToContractEntity(contractDbNode)

					metadata, err := utils.ToJson(ActionLikelihoodMetadata{
						Reason:     utils.IfNotNilString(input.Comments),
						Likelihood: input.RenewalLikelihood.String(),
					})
					userName := ""
					userDbNode, err := s.services.Neo4jRepositories.UserReadRepository.GetUserById(ctx, tenant, common.GetUserIdFromContext(ctx))
					if err != nil {
						tracing.TraceErr(span, err)
					}
					if userDbNode != nil {
						userEntity := neo4jmapper.MapDbNodeToUserEntity(userDbNode)
						userName = userEntity.GetFullName()
					}
					message := fmt.Sprintf("Renewal likelihood set to %s", cases.Title(language.English).String(input.RenewalLikelihood.String()))
					if userName != "" {
						message += fmt.Sprintf(" by %s", userName)
					}

					extraActionProperties := map[string]interface{}{
						"comments": utils.IfNotNilString(input.Comments),
					}
					_, err = s.services.Neo4jRepositories.ActionWriteRepository.CreateWithProperties(ctx, tenant, contractEntity.Id, commonModel.CONTRACT, neo4jenum.ActionRenewalLikelihoodUpdated, message, metadata, utils.Now(), common.GetAppSourceFromContext(ctx), extraActionProperties)
					if err != nil {
						tracing.TraceErr(span, err)
						s.log.Errorf("error while creating action for opportunity %s: %s", opportunityId, err.Error())
					}
				}
				// adjusted rate change action
				if (likelihoodChanged || adjustedRateChanged) && !amountChanged {
					contractDbNode, err := s.services.Neo4jRepositories.ContractReadRepository.GetContractByOpportunityId(ctx, tenant, opportunityId)
					if err != nil {
						tracing.TraceErr(span, err)
						s.log.Errorf("error while getting contract for opportunity %s: %s", opportunityId, err.Error())
						return nil
					}
					if contractDbNode == nil {
						return nil
					}
					contract := neo4jmapper.MapDbNodeToContractEntity(contractDbNode)
					err = s.services.ContractService.UpdateActiveRenewalOpportunityArr(ctx, contract.Id)
					if err != nil {
						tracing.TraceErr(span, err)
						s.log.Errorf("error while updating renewal opportunity %s: %s", opportunityId, err.Error())
						return nil
					}
				} else if amountChanged {
					s.sendEventToUpdateOrganizationArr(ctx, tenant, opportunityId, span)
				}
			}
			return nil
		})

		txWithPostCommit.AddPostCommitAction(func(ctx context.Context) error {
			if createFlow {
				err = s.services.RabbitMQService.PublishEvent(ctx, opportunityId, commonModel.OPPORTUNITY, dto.CreateOpportunity{*input})
				if err != nil {
					tracing.TraceErr(span, errors.Wrap(err, "unable to publish message CreateOpportunity"))
				}
				s.services.RabbitMQService.PublishEventCompleted(ctx, tenant, opportunityId, commonModel.OPPORTUNITY, utils.NewEventCompletedDetails().WithCreate())
			} else {
				err = s.services.RabbitMQService.PublishEvent(ctx, opportunityId, commonModel.OPPORTUNITY, dto.UpdateOpportunity{*input})
				if err != nil {
					tracing.TraceErr(span, errors.Wrap(err, "unable to publish message UpdateOpportunity"))
				}
				if common.GetAppSourceFromContext(ctx) != constants.AppSourceCustomerOsApi {
					s.services.RabbitMQService.PublishEventCompleted(ctx, tenant, opportunityId, commonModel.OPPORTUNITY, utils.NewEventCompletedDetails().WithUpdate())
				}
			}

			return nil
		})

		return nil, nil
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return "", err
	}

	return opportunityId, nil
}

func (s *opportunityService) CloseWon(ctx context.Context, txWithPostCommit *utils.TxWithPostCommit, tenant, opportunityId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OpportunityService.CloseWon")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.SetTag(tracing.SpanTagEntityId, opportunityId)

	_, err := utils.ExecuteWriteInTransactionWithPostCommitActions(ctx, s.services.Neo4jRepositories.Neo4jDriver, s.services.Neo4jRepositories.Database, txWithPostCommit, func(txWithPostCommit *utils.TxWithPostCommit) (any, error) {
		opportunityEntity, err := s.GetById(ctx, txWithPostCommit.Tx, tenant, opportunityId)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}
		if opportunityEntity == nil {
			err = fmt.Errorf("opportunity not found")
			tracing.TraceErr(span, err)
			return nil, err
		}

		// check opportunity is not already closed won
		if opportunityEntity.InternalStage == neo4jenum.OpportunityInternalStageClosedWon {
			return nil, nil
		}

		// Set opportunity as closed won
		err = s.services.Neo4jRepositories.OpportunityWriteRepository.CloseWon(ctx, txWithPostCommit.Tx, tenant, opportunityId, utils.Now())
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}

		// clean external stage
		if opportunityEntity.IsNBO() && opportunityEntity.ExternalStage != "" {
			err = s.services.Neo4jRepositories.OpportunityWriteRepository.Save(ctx, txWithPostCommit.Tx, tenant, opportunityId, data_fields.OpportunityFields{
				ExternalStage: utils.StringPtr(""),
			})
			if err != nil {
				tracing.TraceErr(span, err)
				return nil, err
			}
		}

		// For NBO set organization as customer
		if opportunityEntity.IsNBO() {
			organizationDbNode, err := s.services.Neo4jRepositories.OrganizationReadRepository.GetOrganizationByOpportunityId(ctx, tenant, opportunityId)
			if err != nil {
				tracing.TraceErr(span, err)
			}
			if organizationDbNode != nil {
				organizationEntity := neo4jmapper.MapDbNodeToOrganizationEntity(organizationDbNode)
				// Make organization customer if it's not already
				if organizationEntity.Relationship != neo4jenum.OrganizationRelationshipCustomer && organizationEntity.Stage != neo4jenum.Trial {
					_, err := s.services.OrganizationService.Save(ctx, txWithPostCommit, &organizationEntity.ID, data_fields.OrganizationFields{
						Relationship: utils.ToPtr(neo4jenum.OrganizationRelationshipCustomer),
						Stage:        utils.ToPtr(neo4jenum.OrganizationRelationshipCustomer.DefaultStage()),
					})
					if err != nil {
						tracing.TraceErr(span, err)
						return nil, err
					}
				}
			}
		}

		txWithPostCommit.AddPostCommitAction(func(ctx context.Context) error {
			// create new renewal opportunity
			if opportunityEntity.IsRenewal() {
				// get contract id for opportunity
				contractDbNode, err := s.services.Neo4jRepositories.ContractReadRepository.GetContractByOpportunityId(ctx, tenant, opportunityId)
				if err != nil {
					tracing.TraceErr(span, err)
					return err
				}
				contractEntity := neo4jmapper.MapDbNodeToContractEntity(contractDbNode)
				// create new renewal opportunity
				_, err = s.services.OpportunityService.CreateRenewalOpportunity(ctx, nil, &data_fields.OpportunityFields{
					ContractId: utils.StringPtr(contractEntity.Id),
				})
				if err != nil {
					tracing.TraceErr(span, err)
					return err
				}
			}

			return nil
		})

		txWithPostCommit.AddPostCommitAction(func(ctx context.Context) error {
			s.services.RabbitMQService.PublishEventCompleted(ctx, tenant, opportunityId, commonModel.OPPORTUNITY, utils.NewEventCompletedDetails().WithUpdate())
			return nil
		})

		return nil, nil
	})

	return err
}

func (s *opportunityService) CloseLost(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, opportunityId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OpportunityService.CloseLost")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.SetTag(tracing.SpanTagEntityId, opportunityId)

	opportunity, err := s.GetById(ctx, nil, tenant, opportunityId)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	if opportunity == nil {
		err = fmt.Errorf("opportunity not found")
		tracing.TraceErr(span, err)
		return err
	}

	// check opportunity is not already closed lost
	if opportunity.InternalStage == neo4jenum.OpportunityInternalStageClosedLost {
		return nil
	}

	err = s.services.Neo4jRepositories.OpportunityWriteRepository.CloseLost(ctx, tx, tenant, opportunityId, utils.Now())
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	s.services.RabbitMQService.PublishEventCompleted(ctx, tenant, opportunityId, commonModel.OPPORTUNITY, utils.NewEventCompletedDetails().WithUpdate())

	return nil
}

func (s *opportunityService) Archive(ctx context.Context, tenant, opportunityId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OpportunityService.Archive")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.SetTag(tracing.SpanTagEntityId, opportunityId)

	opportunity, err := s.GetById(ctx, nil, tenant, opportunityId)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	if opportunity == nil {
		err = fmt.Errorf("opportunity not found")
		tracing.TraceErr(span, err)
		return err
	}

	if opportunity.InternalType == neo4jenum.OpportunityInternalTypeRenewal {
		err = errors.New("Renewal opportunity cannot be archived")
		tracing.TraceErr(span, err)
		return err
	}

	err = s.services.Neo4jRepositories.OpportunityWriteRepository.Archive(ctx, tenant, opportunityId)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	s.services.RabbitMQService.PublishEventCompleted(ctx, tenant, opportunityId, commonModel.OPPORTUNITY, utils.NewEventCompletedDetails().WithDelete())

	return nil
}

func (s *opportunityService) RolloutRenewalOpportunity(ctx context.Context, contractId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OpportunityService.RolloutRenewalOpportunity")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.SetTag(tracing.SpanTagEntityId, contractId)

	// validate tenant
	err := common.ValidateTenant(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	tenant := common.GetTenantFromContext(ctx)

	contractDbNode, err := s.services.Neo4jRepositories.ContractReadRepository.GetContractById(ctx, tenant, contractId)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	contractEntity := neo4jmapper.MapDbNodeToContractEntity(contractDbNode)

	if contractEntity.LengthInMonths <= 0 {
		return nil
	}

	currentRenewalOpportunityDbNode, err := s.services.Neo4jRepositories.OpportunityReadRepository.GetActiveRenewalOpportunityForContract(ctx, tenant, contractId)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "Failed getting renewal opportunity for contract"))
		s.log.Errorf("Error while getting renewal opportunity for contract %s: %s", contractId, err.Error())
	}

	ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
	if currentRenewalOpportunityDbNode != nil {
		currentOpportunity := neo4jmapper.MapDbNodeToOpportunityEntity(currentRenewalOpportunityDbNode)

		err = s.services.OpportunityService.CloseWon(ctx, nil, tenant, currentOpportunity.Id)
		if err != nil {
			tracing.TraceErr(span, err)
			s.log.Errorf("CloseWinOpportunity failed: %s", err.Error())
			return err
		}
	}

	err = s.services.ContractService.RecalculateContractLtv(ctx, contractId)
	if err != nil {
		tracing.TraceErr(span, err)
	}

	// Add action in timeline
	status := "Renewed"
	metadata, err := utils.ToJson(ActionStatusMetadata{
		Status: status,
	})
	message := contractEntity.Name + " renewed"

	_, err = s.services.Neo4jRepositories.ActionWriteRepository.Create(ctx, tenant, contractId, commonModel.CONTRACT, neo4jenum.ActionContractRenewed, message, metadata, utils.Now(), common.GetAppSourceFromContext(ctx))
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Failed creating renewed action for contract %s: %s", contractId, err.Error())
	}

	s.services.RabbitMQService.PublishEventCompleted(ctx, tenant, contractId, commonModel.CONTACT, utils.NewEventCompletedDetails().WithUpdate())

	return nil
}

func (s *opportunityService) sendEventToUpdateOrganizationArr(ctx context.Context, tenant, opportunityId string, span opentracing.Span) {
	// if amount changed, recalculate organization combined ARR forecast
	organizationDbNode, err := s.services.Neo4jRepositories.OrganizationReadRepository.GetOrganizationByOpportunityId(ctx, tenant, opportunityId)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("error while getting organization for opportunity %s: %s", opportunityId, err.Error())
		return
	}
	if organizationDbNode == nil {
		return
	}
	organization := neo4jmapper.MapDbNodeToOrganizationEntity(organizationDbNode)

	ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
	_, err = utils.CallEventsPlatformGRPCWithRetry[*organizationpb.OrganizationIdGrpcResponse](func() (*organizationpb.OrganizationIdGrpcResponse, error) {
		return s.services.GrpcClients.OrganizationClient.RefreshArr(ctx, &organizationpb.OrganizationIdGrpcRequest{
			Tenant:         tenant,
			OrganizationId: organization.ID,
			AppSource:      common.GetAppSourceFromContext(ctx),
		})
	})
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("RefreshArr failed: %v", err.Error())
	}
}
