package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/entity"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
)

type ActionRepository interface {
	OrganizationCreatedAction(ctx context.Context, tenant, organizationId, source, appSource string) error
}

type actionRepository struct {
	driver *neo4j.DriverWithContext
}

func NewActionRepository(driver *neo4j.DriverWithContext) ActionRepository {
	return &actionRepository{
		driver: driver,
	}
}

func (r *actionRepository) OrganizationCreatedAction(ctx context.Context, tenant, organizationId, source, appSource string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ActionRepository.OrganizationCreatedAction")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := fmt.Sprintf("MATCH (p:Organization_%s {id:$organizationId}) "+
		"		MERGE (p)<-[:ACTION_ON]-(a:Action {id:randomUUID()}) "+
		"		ON CREATE SET 	a.type=$type, "+
		"						a.createdAt=$createdAt, "+
		"						a.updatedAt=$createdAt, "+
		"						a.source=$source, "+
		"						a.appSource=$appSource, "+
		"						a:Action_%s, "+
		"						a:TimelineEvent, "+
		"						a:TimelineEvent_%s return a", tenant, tenant, tenant)

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, fmt.Sprintf(query),
			map[string]interface{}{
				"tenant":         tenant,
				"organizationId": organizationId,
				"createdAt":      utils.Now(),
				"type":           entity.ActionCreated,
				"source":         source,
				"appSource":      appSource,
			})
		if err != nil {
			return nil, err
		}
		_, err = queryResult.Single(ctx)
		if err != nil {
			return nil, err
		}
		return nil, nil
	})
	return err
}
