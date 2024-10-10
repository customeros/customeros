package service

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	model2 "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"
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
	ReplaceOwner(ctx context.Context, tenant string, opportunityId, userId string) error
	RemoveOwner(ctx context.Context, tenant string, opportunityId string) error
	CloseWon(ctx context.Context, tenant, opportunityId string) error
	CloseLost(ctx context.Context, tenant, opportunityId string) error
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
	var existingOpportunity *neo4jentity.OpportunityEntity

	if organizationId == nil && opportunityId == nil {
		err := fmt.Errorf("(OpportunityService.Save) organizationId and opportunityId are nil")
		tracing.TraceErr(span, err)
		return nil, err
	}

	if organizationId != nil {
		existsById, err := s.services.Neo4jRepositories.CommonReadRepository.ExistsById(ctx, tenant, *organizationId, model2.NodeLabelOrganization)
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
		existingOpportunity, err = s.GetById(ctx, tenant, *opportunityId)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}
		if existingOpportunity == nil {
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

		generatedId, err := s.services.Neo4jRepositories.CommonReadRepository.GenerateId(ctx, tenant, model2.NodeLabelOpportunity)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}
		opportunityId = &generatedId
	}

	// Changing external stage should set internal stage back to OPEN
	if existingOpportunity != nil && input.ExternalStage != "" && existingOpportunity.ExternalStage != input.ExternalStage && existingOpportunity.InternalStage != neo4jenum.OpportunityInternalStageOpen {
		input.InternalStage = neo4jenum.OpportunityInternalStageOpen.String()
		input.UpdateInternalStage = true
	}

	err = s.services.Neo4jRepositories.OpportunityWriteRepository.Save(ctx, tx, tenant, *opportunityId, *input)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	if organizationId != nil {
		err = s.services.Neo4jRepositories.CommonWriteRepository.Link(ctx, tx, tenant, repository.LinkDetails{
			FromEntityId:   *organizationId,
			FromEntityType: model2.ORGANIZATION,
			Relationship:   model2.HAS_OPPORTUNITY,
			ToEntityId:     *opportunityId,
			ToEntityType:   model2.OPPORTUNITY,
		})
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}
	}

	//TODO use TX
	if input.OwnerUserId != "" {
		err = s.services.Neo4jRepositories.OpportunityWriteRepository.ReplaceOwner(ctx, tenant, *opportunityId, input.OwnerUserId)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}
	}
	if input.UpdateOwnerUserId {
		if input.OwnerUserId != "" {
			err = s.services.Neo4jRepositories.OpportunityWriteRepository.ReplaceOwner(ctx, tenant, *opportunityId, input.OwnerUserId)
			if err != nil {
				tracing.TraceErr(span, err)
				return nil, err
			}
		} else {
			err = s.services.Neo4jRepositories.OpportunityWriteRepository.RemoveOwner(ctx, tenant, *opportunityId)
			if err != nil {
				tracing.TraceErr(span, err)
				return nil, err
			}
		}
	}

	//TODO FIX
	//if (input.UpdateAmount || input.UpdateMaxAmount) && existingOpportunity.InternalType == neo4jenum.OpportunityInternalTypeRenewal {
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

	return opportunityId, nil
}

func (s *opportunityService) ReplaceOwner(ctx context.Context, tenant, opportunityId, userId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OpportunityService.ReplaceOwner")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.SetTag(tracing.SpanTagEntityId, opportunityId)
	span.LogFields(log.String("userId", userId))

	opportunityExists, _ := s.services.Neo4jRepositories.CommonReadRepository.ExistsById(ctx, tenant, opportunityId, model2.NodeLabelOpportunity)
	if !opportunityExists {
		err := fmt.Errorf("(OpportunityService.ReplaceOwner) opportunity with id {%s} not found", opportunityId)
		s.log.Error(err.Error())
		tracing.TraceErr(span, err)
		return err
	}

	userExists, _ := s.services.Neo4jRepositories.CommonReadRepository.ExistsById(ctx, tenant, userId, model2.NodeLabelUser)
	if !userExists {
		err := fmt.Errorf("(OpportunityService.ReplaceOwner) user with id {%s} not found", userId)
		s.log.Error(err.Error())
		tracing.TraceErr(span, err)
		return err
	}

	err := s.services.Neo4jRepositories.OpportunityWriteRepository.ReplaceOwner(ctx, tenant, opportunityId, userId)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}

func (s *opportunityService) RemoveOwner(ctx context.Context, tenant, opportunityId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OpportunityService.RemoveOwner")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.SetTag(tracing.SpanTagEntityId, opportunityId)

	opportunityExists, _ := s.services.Neo4jRepositories.CommonReadRepository.ExistsById(ctx, tenant, opportunityId, model2.NodeLabelOpportunity)
	if !opportunityExists {
		err := fmt.Errorf("(OpportunityService.ReplaceOwner) opportunity with id {%s} not found", opportunityId)
		s.log.Error(err.Error())
		tracing.TraceErr(span, err)
		return err
	}

	err := s.services.Neo4jRepositories.OpportunityWriteRepository.RemoveOwner(ctx, tenant, opportunityId)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}

func (s *opportunityService) CloseWon(ctx context.Context, tenant, opportunityId string) error {
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
		err = fmt.Errorf("opportunity already closed won")
		tracing.TraceErr(span, err)
		return err
	}

	_, err = utils.ExecuteWriteInTransaction(ctx, s.services.Neo4jRepositories.Neo4jDriver, s.services.Neo4jRepositories.Database, nil, func(tx neo4j.ManagedTransaction) (any, error) {

		//todo use TX
		err = s.services.Neo4jRepositories.OpportunityWriteRepository.CloseWon(ctx, tenant, opportunityId, utils.Now())
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
			//TODO implement
			// get contract id for opportunity
			//contractDbNode, err := s.services.Neo4jRepositories.ContractReadRepository.GetContractByOpportunityId(ctx, tenant, opportunityId)
			//if err != nil {
			//	tracing.TraceErr(span, err)
			//	return nil, err
			//}
			//contractEntity := neo4jmapper.MapDbNodeToContractEntity(contractDbNode)

			// create new renewal opportunity
			//_, err = subscriptions.CallEventsPlatformGRPCWithRetry[*opportunitypb.OpportunityIdGrpcResponse](func() (*opportunitypb.OpportunityIdGrpcResponse, error) {
			//	return h.grpcClients.OpportunityClient.CreateRenewalOpportunity(ctx, &opportunitypb.CreateRenewalOpportunityGrpcRequest{
			//		Tenant:     eventData.Tenant,
			//		ContractId: contractEntity.Id,
			//		SourceFields: &commonpb.SourceFields{
			//			Source:    constants.SourceOpenline,
			//			AppSource: constants.AppSourceEventProcessingPlatformSubscribers,
			//		},
			//	})
			//})
			//if err != nil {
			//	tracing.TraceErr(span, err)
			//	h.log.Errorf("CreateRenewalOpportunity failed: %s", err.Error())
			//}
		}

		return nil, nil
	})

	return nil
}

func (s *opportunityService) CloseLost(ctx context.Context, tenant, opportunityId string) error {
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
		err = fmt.Errorf("opportunity already closed lost")
		tracing.TraceErr(span, err)
		return err
	}

	err = s.services.Neo4jRepositories.OpportunityWriteRepository.CloseLost(ctx, tenant, opportunityId, utils.Now())
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}
