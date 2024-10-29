package service

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	commonModel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/constants"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"
	commonpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/common"
	opportunitypb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/opportunity"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
)

type OpportunityService interface {
	GetById(ctx context.Context, tenant, opportunityId string) (*neo4jentity.OpportunityEntity, error)
	GetOpportunitiesForContracts(ctx context.Context, tenant string, contractIds []string) (*neo4jentity.OpportunityEntities, error)
	GetOpportunitiesForOrganizations(ctx context.Context, tenant string, organizationIds []string) (*neo4jentity.OpportunityEntities, error)
	GetPaginatedOrganizationOpportunities(ctx context.Context, tenant string, page int, limit int) (*utils.Pagination, error)

	Save(ctx context.Context, tx *neo4j.ManagedTransaction, tenant string, organizationId, opportunityId *string, input *repository.OpportunitySaveFields) (*string, error)
	CloseWon(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, opportunityId string) error
	CloseLost(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, opportunityId string) error
	Archive(ctx context.Context, tenant, opportunityId string) error
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

func (s *opportunityService) GetById(ctx context.Context, tenant, opportunityId string) (*neo4jentity.OpportunityEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OpportunityService.GetById")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("opportunityId", opportunityId))

	if opportunityDbNode, err := s.services.Neo4jRepositories.OpportunityReadRepository.GetOpportunityById(ctx, tenant, opportunityId); err != nil {
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

func (s *opportunityService) Save(ctx context.Context, tx *neo4j.ManagedTransaction, tenant string, organizationId, opportunityId *string, input *repository.OpportunitySaveFields) (*string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OpportunityService.Save")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("input", input))

	var err error
	var existing *neo4jentity.OpportunityEntity

	if organizationId == nil && opportunityId == nil {
		err := fmt.Errorf("(OpportunityService.Save) organizationId and opportunityId are nil")
		tracing.TraceErr(span, err)
		return nil, err
	}

	if organizationId != nil {
		existsById, err := s.services.Neo4jRepositories.CommonReadRepository.ExistsById(ctx, tenant, *organizationId, commonModel.NodeLabelOrganization)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}
		if !existsById {
			err := fmt.Errorf("(OpportunityService.Save) organization with id {%s} not found", *organizationId)
			tracing.TraceErr(span, err)
			return nil, err
		}
	}

	if opportunityId != nil {
		existing, err = s.GetById(ctx, tenant, *opportunityId)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}
		if existing == nil {
			err := fmt.Errorf("(OpportunityService.Save) opportunity with id {%s} not found", *opportunityId)
			tracing.TraceErr(span, err)
			return nil, err
		}
	}

	if opportunityId == nil {

		if input.InternalType == "" {
			input.InternalType = neo4jenum.OpportunityInternalTypeNBO.String()
			input.UpdateInternalType = true
		}

		if input.Currency == "" {
			tenantSettings, err := s.services.TenantService.GetTenantSettings(ctx)
			if err != nil {
				tracing.TraceErr(span, err)
				return nil, err
			}
			input.Currency = tenantSettings.BaseCurrency
			input.UpdateCurrency = true
		}

		generatedId, err := s.services.Neo4jRepositories.CommonReadRepository.GenerateId(ctx, tenant, commonModel.NodeLabelOpportunity)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}
		opportunityId = &generatedId
	}

	// Changing external stage should set internal stage back to OPEN
	if existing != nil && input.ExternalStage != "" && existing.ExternalStage != input.ExternalStage && existing.InternalStage != neo4jenum.OpportunityInternalStageOpen {
		input.InternalStage = neo4jenum.OpportunityInternalStageOpen.String()
		input.UpdateInternalStage = true
	}

	_, err = utils.ExecuteWriteInTransaction(ctx, s.services.Neo4jRepositories.Neo4jDriver, s.services.Neo4jRepositories.Database, tx, func(tx neo4j.ManagedTransaction) (any, error) {

		err = s.services.Neo4jRepositories.OpportunityWriteRepository.Save(ctx, &tx, tenant, *opportunityId, *input)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}

		if organizationId != nil {
			err = s.services.Neo4jRepositories.CommonWriteRepository.Link(ctx, &tx, tenant, repository.LinkDetails{
				FromEntityId:   *organizationId,
				FromEntityType: commonModel.ORGANIZATION,
				Relationship:   commonModel.HAS_OPPORTUNITY,
				ToEntityId:     *opportunityId,
				ToEntityType:   commonModel.OPPORTUNITY,
			})
			if err != nil {
				tracing.TraceErr(span, err)
				return nil, err
			}
		}

		if input.UpdateOwnerId {
			if input.OwnerId != "" {
				err = s.services.Neo4jRepositories.OpportunityWriteRepository.ReplaceOwner(ctx, &tx, tenant, *opportunityId, input.OwnerId)
				if err != nil {
					tracing.TraceErr(span, err)
					return nil, err
				}
			} else {
				if existing != nil {
					err = s.services.Neo4jRepositories.OpportunityWriteRepository.RemoveOwner(ctx, &tx, tenant, *opportunityId)
					if err != nil {
						tracing.TraceErr(span, err)
						return nil, err
					}
				}
			}
		}

		if input.UpdateInternalStage {
			if input.InternalStage == neo4jenum.OpportunityInternalStageClosedWon.String() {
				err := s.CloseWon(ctx, &tx, tenant, *opportunityId)
				if err != nil {
					tracing.TraceErr(span, err)
					return nil, err
				}
			} else if input.InternalStage == neo4jenum.OpportunityInternalStageClosedLost.String() {
				err := s.CloseLost(ctx, &tx, tenant, *opportunityId)
				if err != nil {
					tracing.TraceErr(span, err)
					return nil, err
				}
			}
		}

		//TODO when we migrate the renewal opportunities to the new model, we will need to uncomment this
		//if (input.UpdateAmount || input.UpdateMaxAmount) && existing.InternalType == neo4jenum.OpportunityInternalTypeRenewal {
		//	// if amount changed, recalculate organization combined ARR forecast
		//	organizationDbNode, err := s.services.Neo4jRepositories.OrganizationReadRepository.GetOrganizationByOpportunityId(ctx, tenant, *opportunityId)
		//	if err != nil {
		//		tracing.TraceErr(span, err)
		//		return nil, err
		//	}
		//	if organizationDbNode == nil {
		//		err := fmt.Errorf("organization not found")
		//		tracing.TraceErr(span, err)
		//		return nil, err
		//	}
		//	organization := neo4jmapper.MapDbNodeToOrganizationEntity(organizationDbNode)
		//
		//	ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
		//	_, err = utils.CallEventsPlatformGRPCWithRetry[*organizationpb.OrganizationIdGrpcResponse](func() (*organizationpb.OrganizationIdGrpcResponse, error) {
		//		return s.services.GrpcClients.OrganizationClient.RefreshArr(ctx, &organizationpb.OrganizationIdGrpcRequest{
		//			Tenant:         tenant,
		//			OrganizationId: organization.ID,
		//			AppSource:      input.AppSource,
		//		})
		//	})
		//	if err != nil {
		//		tracing.TraceErr(span, err)
		//		return nil, err
		//	}
		//}

		return nil, nil
	})

	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	//TODO put back after we integrate
	//if input.AppSource != constants.AppSourceCustomerOsApi {
	details := utils.NewEventCompletedDetails()

	if existing == nil {
		details.WithCreate()
	} else {
		details.WithUpdate()
	}

	s.services.RabbitMQService.PublishEventCompleted(ctx, tenant, *opportunityId, commonModel.OPPORTUNITY, details)
	//}

	return opportunityId, nil
}

func (s *opportunityService) CloseWon(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, opportunityId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OpportunityService.CloseWon")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.SetTag(tracing.SpanTagEntityId, opportunityId)

	opportunity, err := s.GetById(ctx, tenant, opportunityId)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	if opportunity == nil {
		err = fmt.Errorf("opportunity not found")
		tracing.TraceErr(span, err)
		return err
	}

	// check opportunity is not already closed won
	if opportunity.InternalStage == neo4jenum.OpportunityInternalStageClosedWon {
		return nil
	}

	_, err = utils.ExecuteWriteInTransaction(ctx, s.services.Neo4jRepositories.Neo4jDriver, s.services.Neo4jRepositories.Database, tx, func(tx neo4j.ManagedTransaction) (any, error) {

		err = s.services.Neo4jRepositories.OpportunityWriteRepository.CloseWon(ctx, &tx, tenant, opportunityId, utils.Now())
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}

		// clean external stage
		if opportunity.InternalType == neo4jenum.OpportunityInternalTypeNBO && opportunity.ExternalStage != "" {
			err = s.services.Neo4jRepositories.OpportunityWriteRepository.Save(ctx, &tx, tenant, opportunityId, repository.OpportunitySaveFields{
				ExternalStage:       "",
				UpdateExternalStage: true,
			})
			if err != nil {
				tracing.TraceErr(span, err)
				return nil, err
			}
		}

		// set organization as customer
		if opportunity.InternalType == neo4jenum.OpportunityInternalTypeNBO {
			organizationDbNode, err := s.services.Neo4jRepositories.OrganizationReadRepository.GetOrganizationByOpportunityId(ctx, tenant, opportunityId)
			if err != nil {
				tracing.TraceErr(span, err)
			}
			if organizationDbNode != nil {
				organizationEntity := neo4jmapper.MapDbNodeToOrganizationEntity(organizationDbNode)
				// Make organization customer if it's not already
				if organizationEntity.Relationship != neo4jenum.Customer && organizationEntity.Stage != neo4jenum.Trial {

					//TODO use TX
					err := s.services.Neo4jRepositories.OrganizationWriteRepository.UpdateOrganization(ctx, tenant, organizationEntity.ID, repository.OrganizationUpdateFields{
						Relationship:       neo4jenum.Customer,
						Stage:              neo4jenum.Customer.DefaultStage(),
						UpdateRelationship: true,
						UpdateStage:        true,
					})
					if err != nil {
						tracing.TraceErr(span, err)
						return nil, err
					}
				}
			}
		}

		// create new renewal opportunity
		if opportunity.InternalType == neo4jenum.OpportunityInternalTypeRenewal {
			// get contract id for opportunity
			contractDbNode, err := s.services.Neo4jRepositories.ContractReadRepository.GetContractByOpportunityId(ctx, tenant, opportunityId)
			if err != nil {
				tracing.TraceErr(span, err)
				return nil, err
			}
			contractEntity := neo4jmapper.MapDbNodeToContractEntity(contractDbNode)
			// create new renewal opportunity
			_, err = utils.CallEventsPlatformGRPCWithRetry[*opportunitypb.OpportunityIdGrpcResponse](func() (*opportunitypb.OpportunityIdGrpcResponse, error) {
				return s.services.GrpcClients.OpportunityClient.CreateRenewalOpportunity(ctx, &opportunitypb.CreateRenewalOpportunityGrpcRequest{
					Tenant:     tenant,
					ContractId: contractEntity.Id,
					SourceFields: &commonpb.SourceFields{
						Source:    constants.SourceOpenline,
						AppSource: common.GetAppSourceFromContext(ctx),
					},
				})
			})
			if err != nil {
				tracing.TraceErr(span, err)
				return nil, err
			}
		}

		return nil, nil
	})

	utils.EventCompleted(ctx, tenant, commonModel.OPPORTUNITY.String(), opportunityId, s.services.GrpcClients, utils.NewEventCompletedDetails().WithUpdate())
	s.services.RabbitMQService.PublishEventCompleted(ctx, tenant, opportunityId, commonModel.OPPORTUNITY, utils.NewEventCompletedDetails().WithUpdate())

	return nil
}

func (s *opportunityService) CloseLost(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, opportunityId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OpportunityService.CloseLost")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.SetTag(tracing.SpanTagEntityId, opportunityId)

	opportunity, err := s.GetById(ctx, tenant, opportunityId)
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

	utils.EventCompleted(ctx, tenant, commonModel.OPPORTUNITY.String(), opportunityId, s.services.GrpcClients, utils.NewEventCompletedDetails().WithUpdate())
	s.services.RabbitMQService.PublishEventCompleted(ctx, tenant, opportunityId, commonModel.OPPORTUNITY, utils.NewEventCompletedDetails().WithUpdate())

	return nil
}

func (s *opportunityService) Archive(ctx context.Context, tenant, opportunityId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OpportunityService.Archive")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.SetTag(tracing.SpanTagEntityId, opportunityId)

	opportunity, err := s.GetById(ctx, tenant, opportunityId)
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

	utils.EventCompleted(ctx, tenant, commonModel.OPPORTUNITY.String(), opportunityId, s.services.GrpcClients, utils.NewEventCompletedDetails().WithDelete())
	s.services.RabbitMQService.PublishEventCompleted(ctx, tenant, opportunityId, commonModel.OPPORTUNITY, utils.NewEventCompletedDetails().WithDelete())

	return nil
}
