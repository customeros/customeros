package repository

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type CountryRepository interface {
	GetAllForPhoneNumbers(ctx context.Context, ids []string) ([]*utils.DbNodeAndId, error)
}

type countryRepository struct {
	driver *neo4j.DriverWithContext
}

func NewCountryRepository(driver *neo4j.DriverWithContext) CountryRepository {
	return &countryRepository{
		driver: driver,
	}
}

func (r *countryRepository) GetAllForPhoneNumbers(ctx context.Context, ids []string) ([]*utils.DbNodeAndId, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CountryRepository.GetAllForPhoneNumbers")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)
	span.LogFields(log.Object("ids", ids))

	query := `MATCH (:Tenant {name:$tenant})<-[:PHONE_NUMBER_BELONGS_TO_TENANT]-(p:PhoneNumber)-[:LINKED_TO]->(c:Country)
		 		WHERE p.id IN $ids 
		 		RETURN c, p.id`
	span.LogFields(log.String("query", query))

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, query,
			map[string]any{
				"tenant": common.GetTenantFromContext(ctx),
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
