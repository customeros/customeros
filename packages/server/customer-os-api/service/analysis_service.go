package service

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"golang.org/x/exp/slices"
)

type AnalysisService interface {
	GetAnalysisById(ctx context.Context, id string) (*entity.AnalysisEntity, error)
	GetDescribesForAnalysis(ctx context.Context, ids []string) (*entity.AnalysisDescribes, error)
	GetDescribedByForXX(ctx context.Context, ids []string, linkedWith repository.LinkedWith) (*entity.AnalysisEntities, error)

	Create(ctx context.Context, newAnalysis *AnalysisCreateData) (*entity.AnalysisEntity, error)

	convertDbNodesAnalysisDescribes(records []*utils.DbNodeAndId) entity.AnalysisDescribes
	mapDbNodeToAnalysisEntity(node dbtype.Node) *entity.AnalysisEntity
}

type AnalysisDescriptionData struct {
	InteractionEventId   *string
	InteractionSessionId *string
	MeetingId            *string
}

type AnalysisCreateData struct {
	AnalysisEntity *entity.AnalysisEntity
	Describes      []AnalysisDescriptionData
	Source         entity.DataSource
	SourceOfTruth  entity.DataSource
}

type analysisService struct {
	log          logger.Logger
	repositories *repository.Repositories
	services     *Services
}

func NewAnalysisService(log logger.Logger, repositories *repository.Repositories, services *Services) AnalysisService {
	return &analysisService{
		log:          log,
		repositories: repositories,
		services:     services,
	}
}

func (s *analysisService) GetDescribesForAnalysis(ctx context.Context, ids []string) (*entity.AnalysisDescribes, error) {
	records, err := s.repositories.AnalysisRepository.GetDescribesForAnalysis(ctx, common.GetTenantFromContext(ctx), ids)
	if err != nil {
		return nil, err
	}

	analysisDescribes := s.convertDbNodesAnalysisDescribes(records)

	return &analysisDescribes, nil
}

func (s *analysisService) GetDescribedByForXX(ctx context.Context, ids []string, linkedWith repository.LinkedWith) (*entity.AnalysisEntities, error) {
	records, err := s.repositories.AnalysisRepository.GetDescribedByForXX(ctx, common.GetTenantFromContext(ctx), ids, linkedWith)
	if err != nil {
		return nil, err
	}

	analysisDescribes := s.convertDbNodesToAnalysis(records)

	return &analysisDescribes, nil
}

func (s *analysisService) Create(ctx context.Context, newAnalysis *AnalysisCreateData) (*entity.AnalysisEntity, error) {
	session := utils.NewNeo4jWriteSession(ctx, s.getNeo4jDriver())
	defer session.Close(ctx)

	interactionEventDbNode, err := session.ExecuteWrite(ctx, s.createAnalysisInDBTxWork(ctx, newAnalysis))
	if err != nil {
		return nil, err
	}
	return s.mapDbNodeToAnalysisEntity(*interactionEventDbNode.(*dbtype.Node)), nil
}

func (s *analysisService) createAnalysisInDBTxWork(ctx context.Context, newAnalysis *AnalysisCreateData) func(tx neo4j.ManagedTransaction) (any, error) {
	return func(tx neo4j.ManagedTransaction) (any, error) {
		tenant := common.GetContext(ctx).Tenant
		analysisDbNode, err := s.repositories.AnalysisRepository.Create(ctx, tx, tenant, *newAnalysis.AnalysisEntity, newAnalysis.Source, newAnalysis.SourceOfTruth)
		if err != nil {
			return nil, err
		}
		var analysisId = utils.GetPropsFromNode(*analysisDbNode)["id"].(string)

		for _, describes := range newAnalysis.Describes {
			if describes.InteractionSessionId != nil {
				err := s.repositories.AnalysisRepository.LinkWithDescribesXXInTx(ctx, tx, tenant, repository.LINKED_WITH_INTERACTION_SESSION, *describes.InteractionSessionId, analysisId)
				if err != nil {
					return nil, err
				}
			}
			if describes.InteractionEventId != nil {
				err := s.repositories.AnalysisRepository.LinkWithDescribesXXInTx(ctx, tx, tenant, repository.LINKED_WITH_INTERACTION_EVENT, *describes.InteractionEventId, analysisId)
				if err != nil {
					return nil, err
				}
			}
			if describes.MeetingId != nil {
				err := s.repositories.AnalysisRepository.LinkWithDescribesXXInTx(ctx, tx, tenant, repository.LINKED_WITH_MEETING, *describes.MeetingId, analysisId)
				if err != nil {
					return nil, err
				}
			}
		}

		return analysisDbNode, nil
	}
}

func (s *analysisService) GetAnalysisById(ctx context.Context, id string) (*entity.AnalysisEntity, error) {
	session := utils.NewNeo4jReadSession(ctx, s.getNeo4jDriver())
	defer session.Close(ctx)

	queryResult, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		result, err := tx.Run(ctx, fmt.Sprintf(`
			MATCH (a:Analysis_%s {id:$id}) RETURN a`,
			common.GetTenantFromContext(ctx)),
			map[string]interface{}{
				"id": id,
			})
		record, err := result.Single(ctx)
		if err != nil {
			return nil, err
		}
		return record.Values[0], nil
	})
	if err != nil {
		return nil, err
	}

	return s.mapDbNodeToAnalysisEntity(queryResult.(dbtype.Node)), nil
}

func (s *analysisService) mapDbNodeToAnalysisEntity(node dbtype.Node) *entity.AnalysisEntity {
	props := utils.GetPropsFromNode(node)
	createdAt := utils.GetTimePropOrEpochStart(props, "createdAt")
	analysisEntity := entity.AnalysisEntity{
		Id:            utils.GetStringPropOrEmpty(props, "id"),
		CreatedAt:     &createdAt,
		AnalysisType:  utils.GetStringPropOrEmpty(props, "analysisType"),
		Content:       utils.GetStringPropOrEmpty(props, "content"),
		ContentType:   utils.GetStringPropOrEmpty(props, "contentType"),
		AppSource:     utils.GetStringPropOrEmpty(props, "appSource"),
		Source:        entity.GetDataSource(utils.GetStringPropOrEmpty(props, "source")),
		SourceOfTruth: entity.GetDataSource(utils.GetStringPropOrEmpty(props, "sourceOfTruth")),
	}
	return &analysisEntity
}

func (s *analysisService) convertDbNodesToAnalysis(records []*utils.DbNodeAndId) entity.AnalysisEntities {
	analysises := entity.AnalysisEntities{}
	for _, v := range records {
		analysis := s.mapDbNodeToAnalysisEntity(*v.Node)
		analysis.DataloaderKey = v.LinkedNodeId
		analysises = append(analysises, *analysis)

	}
	return analysises
}

func (s *analysisService) convertDbNodesAnalysisDescribes(records []*utils.DbNodeAndId) entity.AnalysisDescribes {
	analysisDescribes := entity.AnalysisDescribes{}
	for _, v := range records {
		if slices.Contains(v.Node.Labels, entity.NodeLabel_InteractionSession) {
			sessionEntity := s.services.InteractionSessionService.mapDbNodeToInteractionSessionEntity(*v.Node)
			sessionEntity.DataloaderKey = v.LinkedNodeId
			analysisDescribes = append(analysisDescribes, sessionEntity)
		} else if slices.Contains(v.Node.Labels, entity.NodeLabel_InteractionEvent) {
			eventEntity := s.services.InteractionEventService.mapDbNodeToInteractionEventEntity(*v.Node)
			eventEntity.DataloaderKey = v.LinkedNodeId
			analysisDescribes = append(analysisDescribes, eventEntity)
		} else if slices.Contains(v.Node.Labels, entity.NodeLabel_Meeting) {
			meetingEntity := s.services.MeetingService.mapDbNodeToMeetingEntity(*v.Node)
			meetingEntity.DataloaderKey = v.LinkedNodeId
			analysisDescribes = append(analysisDescribes, meetingEntity)
		}
	}
	return analysisDescribes
}

func MapAnalysisDescriptionInputToDescriptionData(input []*model.AnalysisDescriptionInput) []AnalysisDescriptionData {
	var inputData []AnalysisDescriptionData
	for _, analysisDescriptionInput := range input {
		inputData = append(inputData, AnalysisDescriptionData{
			InteractionEventId:   analysisDescriptionInput.InteractionEventID,
			InteractionSessionId: analysisDescriptionInput.InteractionSessionID,
			MeetingId:            analysisDescriptionInput.MeetingID,
		})
	}
	return inputData
}
func (s *analysisService) getNeo4jDriver() neo4j.DriverWithContext {
	return *s.repositories.Drivers.Neo4jDriver
}
