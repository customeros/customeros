package repository

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/user/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/helper"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type PlayerRepository interface {
	Merge(ctx context.Context, tenant, userId string, event events.UserAddPlayerInfoEvent) error
}

type playerRepository struct {
	driver *neo4j.DriverWithContext
}

func NewPlayerRepository(driver *neo4j.DriverWithContext) PlayerRepository {
	return &playerRepository{
		driver: driver,
	}
}

func (r *playerRepository) Merge(ctx context.Context, tenant, userId string, event events.UserAddPlayerInfoEvent) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PlayerRepository.Merge")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(ctx, span, tenant)
	span.LogFields(log.String("userId", userId), log.Object("event", event))

	query := `MATCH (:Tenant)<-[:USER_BELONGS_TO_TENANT]-(u:User {id:$userId})
				MERGE (p:Player {authId:$authId, provider:$provider})
				ON CREATE SET p.id=randomUUID(),
							  p.identityId=$identityId,
							  p.createdAt=$createdAt,	
							  p.updatedAt=$updatedAt,
							  p.appSource=$appSource,	
							  p.source=$source,
							  p.sourceOfTruth=$sourceOfTruth
				ON MATCH SET p.updatedAt=$updatedAt
				MERGE (p)-[r:IDENTIFIES]->(u)
				SET r.default = CASE WHEN NOT EXISTS((p)-[:IDENTIFIES {default: true}]->(:User)) THEN true ELSE false END
				`

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	return r.executeQuery(ctx, query, map[string]any{
		"tenant":        tenant,
		"userId":        userId,
		"authId":        event.AuthId,
		"provider":      event.Provider,
		"identityId":    event.IdentityId,
		"createdAt":     event.CreatedAt,
		"updatedAt":     utils.Now(),
		"appSource":     helper.GetAppSource(event.SourceFields.AppSource),
		"source":        helper.GetSource(event.SourceFields.Source),
		"sourceOfTruth": helper.GetSourceOfTruth(event.SourceFields.SourceOfTruth),
	})
}

func (r *playerRepository) executeQuery(ctx context.Context, query string, params map[string]any) error {
	return utils.ExecuteWriteQuery(ctx, *r.driver, query, params)
}
