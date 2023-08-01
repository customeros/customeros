package repository

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
)

type CalendarRepository interface {
	GetAllForUsers(ctx context.Context, tenant string, userIds []string) ([]*utils.DbNodeAndId, error)
}

type calendarRepository struct {
	driver *neo4j.DriverWithContext
}

func NewCalendarRepository(driver *neo4j.DriverWithContext) CalendarRepository {
	return &calendarRepository{
		driver: driver,
	}
}

func (r *calendarRepository) GetAllForUsers(ctx context.Context, tenant string, userIds []string) ([]*utils.DbNodeAndId, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CalendarRepository.GetAllForUsers")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, `
			MATCH (:Tenant {name:$tenant})<-[:USER_BELONGS_TO_TENANT]-(u:User)-[:HAS_CALENDAR]->(cal:Calendar)
			WHERE u.id IN $userIds
			RETURN cal, u.id as userId`,
			map[string]any{
				"tenant":  tenant,
				"userIds": userIds,
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
