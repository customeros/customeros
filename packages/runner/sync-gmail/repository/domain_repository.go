package repository

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"time"
)

type DomainRepository interface {
	GetDomain(ctx context.Context, domain string) (*dbtype.Node, error)
	GetDomainInTx(ctx context.Context, tx neo4j.ManagedTransaction, domain string) (*dbtype.Node, error)
	CreateDomain(ctx context.Context, tx neo4j.ManagedTransaction, domain, source, appSource string, now time.Time) (*dbtype.Node, error)
}

type domainRepository struct {
	driver *neo4j.DriverWithContext
}

func NewDomainRepository(driver *neo4j.DriverWithContext) DomainRepository {
	return &domainRepository{
		driver: driver,
	}
}

func (r *domainRepository) GetDomain(ctx context.Context, domain string) (*dbtype.Node, error) {
	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := `MATCH (d:Domain{domain:$domain}) RETURN d`

	if result, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, query,
			map[string]any{
				"domain": domain,
			})
		return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
	}); err != nil && err.Error() == "Result contains no more records" {
		return nil, nil
	} else if err != nil {
		return nil, nil
	} else {
		return result.(*dbtype.Node), nil
	}
}

func (r *domainRepository) GetDomainInTx(ctx context.Context, tx neo4j.ManagedTransaction, domain string) (*dbtype.Node, error) {
	query := `MATCH (d:Domain{domain:$domain}) RETURN d`

	queryResult, err := tx.Run(ctx, query,
		map[string]any{
			"domain": domain,
		})
	if err != nil {
		return nil, err
	}
	result, err := utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
	if err != nil && err.Error() == "Result contains no more records" {
		return nil, nil
	} else if err != nil {
		return nil, nil
	} else {
		return result, nil
	}
}

func (r *domainRepository) CreateDomain(ctx context.Context, tx neo4j.ManagedTransaction, domain, source, appSource string, now time.Time) (*dbtype.Node, error) {
	query := "MERGE (d:Domain {domain:$domain}) " +
		" ON CREATE SET " +
		"  d.createdAt=$now, " +
		"  d.updatedAt=$now, " +
		"  d.source=$source, " +
		"  d.appSource=$appSource " +
		" RETURN d"

	queryResult, err := tx.Run(ctx, query,
		map[string]any{
			"domain":    domain,
			"source":    source,
			"appSource": appSource,
			"now":       now,
		})
	if err != nil {
		return nil, err
	}
	result, err := utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
	if err != nil {
		return nil, err
	} else {
		return result, nil
	}
}
