package repository

import (
	"errors"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"golang.org/x/net/context"
)

type TenantRepository interface {
	Merge(ctx context.Context, tenant entity.TenantEntity) (*dbtype.Node, error)
	GetForDomain(ctx context.Context, tenant string, domain string) (*dbtype.Node, error)
	LinkWithDomain(ctx context.Context, tenant, domain string) error
}

type tenantRepository struct {
	driver *neo4j.DriverWithContext
}

func NewTenantRepository(driver *neo4j.DriverWithContext) TenantRepository {
	return &tenantRepository{
		driver: driver,
	}
}

func (r *tenantRepository) LinkWithDomain(ctx context.Context, tenant, domain string) error {
	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)
	query := `
			MATCH (t:Tenant {name:$tenant}),
			(d:Domain {domain: $domain})
			MERGE (t)-[:HAS_DOMAIN]->(d)`
	if _, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, query,
			map[string]any{
				"tenant": tenant,
				"domain": domain,
			})
		return nil, err
	}); err != nil {
		return err
	}
	return nil
}

func (r *tenantRepository) Merge(ctx context.Context, tenant entity.TenantEntity) (*dbtype.Node, error) {
	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := "MERGE (d:Tenant {name:$name}) " +
		" ON CREATE SET " +
		"  d.id=randomUUID(), " +
		"  d.createdAt=$now, " +
		"  d.updatedAt=$now, " +
		"  d.source=$source, " +
		"  d.appSource=$appSource " +
		" RETURN d"

	if result, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, query,
			map[string]any{
				"name":      tenant.Name,
				"source":    tenant.Source,
				"appSource": tenant.AppSource,
				"now":       utils.Now(),
			})
		return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
	}); err != nil {
		return nil, err
	} else {
		return result.(*dbtype.Node), nil
	}
}

func (r *tenantRepository) GetForDomain(ctx context.Context, tenant string, domain string) (*dbtype.Node, error) {
	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, `
			MATCH (t:Tenant)-[:HAS_DOMAIN]->(d:Domain)
			WHERE d.domain=$domainName
			RETURN DISTINCT t LIMIT 1`,
			map[string]any{
				"tenant":     tenant,
				"domainName": domain,
			}); err != nil {
			return nil, err
		} else {
			return utils.ExtractAllRecordsFirstValueAsNodePtrs(ctx, queryResult, err)
		}
	})
	if err != nil {
		return nil, err
	}
	convertedResult, isOk := result.([]*dbtype.Node)
	if !isOk {
		return nil, errors.New("GetForDomain: cannot convert result")
	}
	if len(convertedResult) == 0 {
		return nil, nil
	}
	return convertedResult[0], err
}
