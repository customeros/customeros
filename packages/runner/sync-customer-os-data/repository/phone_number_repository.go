package repository

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
)

type PhoneNumberRepository interface {
	GetAllCrossTenantsWithRawPhoneNumber(ctx context.Context, size int) ([]*utils.DbNodeAndId, error)
}

type phoneNumberRepository struct {
	driver *neo4j.DriverWithContext
}

func NewPhoneNumberRepository(driver *neo4j.DriverWithContext) PhoneNumberRepository {
	return &phoneNumberRepository{
		driver: driver,
	}
}

func (r *phoneNumberRepository) GetAllCrossTenantsWithRawPhoneNumber(ctx context.Context, size int) ([]*utils.DbNodeAndId, error) {
	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, `
			MATCH (p:PhoneNumber)--(t:Tenant)
 			WHERE (p.syncedWithEventStore is null or p.syncedWithEventStore=false)
			and (p.rawPhoneNumber is not null and p.rawPhoneNumber <> '')
			RETURN p, t.name limit $size`,
			map[string]any{
				"size": size,
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
