package repository

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"time"
)

type PlayerFields struct {
	AuthId       string       `json:"authId"`
	Provider     string       `json:"provider"`
	IdentityId   string       `json:"identityId"`
	CreatedAt    time.Time    `json:"createdAt"`
	UpdatedAt    time.Time    `json:"updatedAt"`
	SourceFields model.Source `json:"sourceFields"`
}

type PlayerWriteRepository interface {
	Merge(ctx context.Context, tenant, userId string, data PlayerFields) error
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

func (r *playerWriteRepository) Merge(ctx context.Context, tenant, userId string, data PlayerFields) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PlayerWriteRepository.Merge")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.LogFields(log.String("userId", userId))
	tracing.LogObjectAsJson(span, "data", data)

	cypher := `MATCH (:Tenant)<-[:USER_BELONGS_TO_TENANT]-(u:User {id:$userId})
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
				SET r.default = CASE WHEN NOT EXISTS((p)-[:IDENTIFIES {default: true}]->(:User)) THEN true ELSE false END`
	params := map[string]any{
		"tenant":        tenant,
		"userId":        userId,
		"authId":        data.AuthId,
		"provider":      data.Provider,
		"identityId":    data.IdentityId,
		"createdAt":     data.CreatedAt,
		"updatedAt":     data.UpdatedAt,
		"appSource":     data.SourceFields.AppSource,
		"source":        data.SourceFields.Source,
		"sourceOfTruth": data.SourceFields.SourceOfTruth,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}
