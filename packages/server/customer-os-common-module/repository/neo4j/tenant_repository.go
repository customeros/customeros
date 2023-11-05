package repository

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type tenantRepository struct {
	driver *neo4j.DriverWithContext
}

type TenantRepository interface {
	TenantExists(ctx context.Context, name string) (bool, error)
}

func NewTenantRepository(driver *neo4j.DriverWithContext) TenantRepository {
	return &tenantRepository{
		driver: driver,
	}
}

func (u *tenantRepository) TenantExists(ctx context.Context, tenantName string) (bool, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TenantRepository.TenantExists")
	defer span.Finish()
	span.SetTag(tracing.SpanTagComponent, "neo4jRepository")
	span.LogFields(log.String("tenantName", tenantName))

	session := (*u.driver).NewSession(
		ctx,
		neo4j.SessionConfig{
			AccessMode: neo4j.AccessModeRead,
			BoltLogger: neo4j.ConsoleBoltLogger(),
		},
	)
	defer session.Close(ctx)

	records, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, `
			MATCH (t:Tenant {name:$name})
			RETURN t.id`,
			map[string]interface{}{
				"name": tenantName,
			})
		if err != nil {
			return nil, err
		}
		return queryResult.Collect(ctx)
	})
	if err != nil {
		return false, err
	}
	if len(records.([]*neo4j.Record)) > 0 {
		return true, nil
	} else {
		return false, nil
	}
}
