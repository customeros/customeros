package repository

import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"golang.org/x/net/context"
	"time"
)

type DescribesType string

const (
	DESCRIBES_TYPE_INTERACTION_SESSION DescribesType = "InteractionSession"
	DESCRIBES_TYPE_INTERACTION_EVENT   DescribesType = "InteractionEvent"
	DESCRIBES_TYPE_MEETING             DescribesType = "Meeting"
)

type AnalysisRepository interface {
	LinkWithDescribesXXInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant string, describesType DescribesType, analysisId, describedId string) error
	GetDescribesForAnalysis(ctx context.Context, tenant string, ids []string) ([]*utils.DbNodeAndId, error)
	Create(ctx context.Context, tx neo4j.ManagedTransaction, tenant string, newAnalysis entity.AnalysisEntity, source, sourceOfTruth entity.DataSource) (*dbtype.Node, error)
}

type analysisRepository struct {
	driver *neo4j.DriverWithContext
}

func NewAnalysisRepository(driver *neo4j.DriverWithContext) AnalysisRepository {
	return &analysisRepository{
		driver: driver,
	}
}

func (r *analysisRepository) LinkWithDescribesXXInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant string, describesType DescribesType, analysisId, describedId string) error {

	query := fmt.Sprintf(`MATCH (d:%s_%s {id:$describedId}) `, describesType, tenant)
	query += fmt.Sprintf(`MATCH (a:Analysis_%s {id:$analysisId}) `, tenant)
	query += `MERGE (a)-[r:DESCRIBES]->(d) `
	query += `return r `

	queryResult, err := tx.Run(ctx, query,
		map[string]any{
			"describedId": describedId,
			"analysisId":  analysisId,
		})
	if err != nil {
		return err
	}
	_, err = queryResult.Single(ctx)
	return err
}

func (r *analysisRepository) GetDescribesForAnalysis(ctx context.Context, tenant string, ids []string) ([]*utils.DbNodeAndId, error) {
	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := "MATCH (a:Analysis_%s)-[DESCRIBES]->(d) " +
		" WHERE a.id IN $ids " +
		" RETURN d, a.id"

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, fmt.Sprintf(query, tenant),
			map[string]any{
				"tenant": tenant,
				"ids":    ids,
			}); err != nil {
			return nil, err
		} else {
			return utils.ExtractAllRecordsAsDbNodeAndId(ctx, queryResult, err)
		}
	})
	if err != nil {
		return nil, err
	}
	return result.([]*utils.DbNodeAndId), err
}

func (r *analysisRepository) Create(ctx context.Context, tx neo4j.ManagedTransaction, tenant string, newAnalysis entity.AnalysisEntity, source, sourceOfTruth entity.DataSource) (*dbtype.Node, error) {
	var createdAt time.Time
	createdAt = utils.Now()
	if newAnalysis.CreatedAt != nil {
		createdAt = *newAnalysis.CreatedAt
	}

	query := "MERGE (a:Analysis_%s {id:randomUUID()}) ON CREATE SET " +
		" a:Analysis, " +
		" a:TimelineEvent, " +
		" a:TimelineEvent_%s, " +
		" a.source=$source, " +
		" a.createdAt=$createdAt, " +
		" a.analysisType=$analysisType, " +
		" a.content=$content, " +
		" a.contentType=$contentType, " +
		" a.sourceOfTruth=$sourceOfTruth, " +
		" a.appSource=$appSource " +
		" RETURN a"

	if queryResult, err := tx.Run(ctx, fmt.Sprintf(query, tenant, tenant),
		map[string]interface{}{
			"tenant":        tenant,
			"source":        source,
			"createdAt":     createdAt,
			"analysisType":  newAnalysis.AnalysisType,
			"content":       newAnalysis.Content,
			"contentType":   newAnalysis.ContentType,
			"sourceOfTruth": sourceOfTruth,
			"appSource":     newAnalysis.AppSource,
		}); err != nil {
		return nil, err
	} else {
		return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
	}
}
