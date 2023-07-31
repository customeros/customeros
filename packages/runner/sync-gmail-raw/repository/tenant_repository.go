package repository

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"golang.org/x/net/context"
)

type TenantRepository interface {
	GetAll(ctx context.Context) ([]*dbtype.Node, error)
}

type tenantRepository struct {
	driver *neo4j.DriverWithContext
}

func NewTenantRepository(driver *neo4j.DriverWithContext) TenantRepository {
	return &tenantRepository{
		driver: driver,
	}
}

func (r *tenantRepository) GetAll(ctx context.Context) ([]*dbtype.Node, error) {
	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, `MATCH (t:Tenant) return t`, map[string]any{})
		return utils.ExtractAllRecordsFirstValueAsDbNodePtrs(ctx, queryResult, err)
	})
	if err != nil {
		return nil, err
	}

	return result.([]*dbtype.Node), nil
}
