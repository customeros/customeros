package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/model"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"time"
)

type PlayerFields struct {
	AuthId       string             `json:"authId"`
	Provider     string             `json:"provider"`
	IdentityId   string             `json:"identityId"`
	CreatedAt    time.Time          `json:"createdAt"`
	SourceFields model.SourceFields `json:"sourceFields"`
}

type PlayerWriteRepository interface {
	Merge(ctx context.Context, userId string, data entity.PlayerEntity) error

	SetDefaultUser(ctx context.Context, tenant, userId, playerId string, relation entity.PlayerRelation) error
	LinkWithUser(ctx context.Context, tenant, userId, playerId string, relation entity.PlayerRelation) error
	UnlinkUser(ctx context.Context, tenant, userId, playerId string, relation entity.PlayerRelation) error
}

type playerWriteRepository struct {
	driver   *neo4j.DriverWithContext
	database string
}

func NewPlayerWriteRepository(driver *neo4j.DriverWithContext, database string) PlayerWriteRepository {
	return &playerWriteRepository{
		driver:   driver,
		database: database,
	}
}

func (r *playerWriteRepository) Merge(c context.Context, userId string, data entity.PlayerEntity) error {
	span, ctx := opentracing.StartSpanFromContext(c, "PlayerWriteRepository.Merge")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	tenant := common.GetTenantFromContext(ctx)

	span.LogFields(log.String("userId", userId))
	tracing.LogObjectAsJson(span, "data", data)

	cypher := `MATCH (:Tenant)<-[:USER_BELONGS_TO_TENANT]-(u:User {id:$userId})
				MERGE (p:Player {authId:$authId, provider:$provider})
				ON CREATE SET p.id=randomUUID(),
							  p.identityId=$identityId,
							  p.createdAt=$createdAt,	
							  p.updatedAt=datetime(),
							  p.appSource=$appSource,	
							  p.source=$source,
							  p.sourceOfTruth=$sourceOfTruth
				ON MATCH SET p.updatedAt=datetime()
				MERGE (p)-[r:IDENTIFIES]->(u)
				SET r.default = CASE WHEN NOT EXISTS((p)-[:IDENTIFIES {default: true}]->(:User)) THEN true ELSE false END`
	params := map[string]any{
		"tenant":        tenant,
		"userId":        userId,
		"authId":        data.AuthId,
		"provider":      data.Provider,
		"identityId":    data.IdentityId,
		"createdAt":     utils.NowIfZero(data.CreatedAt),
		"appSource":     data.AppSource,
		"source":        utils.StringFirstNonEmpty(data.Source, constants.SourceOpenline),
		"sourceOfTruth": utils.StringFirstNonEmpty(data.SourceOfTruth, constants.SourceOpenline),
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}

func (r *playerWriteRepository) SetDefaultUser(c context.Context, tenant, userId, playerId string, relation entity.PlayerRelation) error {
	span, ctx := opentracing.StartSpanFromContext(c, "PlayerWriteRepository.SetDefaultUser")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	cypher := fmt.Sprintf(`
			MATCH (p:Player {id:$playerId})-[r:%s]->(u:User_%s)
			SET r.default=
				CASE u.id
					WHEN $userId THEN true
					ELSE false
				END
			RETURN DISTINCT(p)`, relation, tenant)
	params := map[string]any{
		"playerId": playerId,
		"userId":   userId,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}

func (r *playerWriteRepository) LinkWithUser(c context.Context, tenant, userId, playerId string, relation entity.PlayerRelation) error {
	span, ctx := opentracing.StartSpanFromContext(c, "PlayerWriteRepository.LinkWithUser")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	cypher := fmt.Sprintf(`
			MATCH (p:Player {id:$playerId}), (u:User {id:$userId})-[:USER_BELONGS_TO_TENANT]->(t:Tenant {name:$tenant})
			MERGE (p)-[r:%s]->(u)
			SET r.default= CASE
				WHEN NOT EXISTS((p)-[:%s {default: true}]->(:User)) THEN true
				ELSE false
			END
			RETURN p`, relation, tenant)
	params := map[string]any{
		"playerId": playerId,
		"userId":   userId,
		"tenant":   tenant,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}

func (r *playerWriteRepository) UnlinkUser(c context.Context, tenant, userId, playerId string, relation entity.PlayerRelation) error {
	span, ctx := opentracing.StartSpanFromContext(c, "PlayerWriteRepository.UnlinkUser")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	cypher := fmt.Sprintf(`
			MATCH (p:Player {id:$playerId}), (u:User_%s {id:$userId})
							MATCH (p)-[r:%s]->(u)
							DELETE r return p`, relation, tenant)
	params := map[string]any{
		"playerId": playerId,
		"userId":   userId,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}
