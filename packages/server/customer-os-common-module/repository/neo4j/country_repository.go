package repository

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/db"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/repository/neo4j/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
)

type CountryRepository interface {
	Create(ctx context.Context, entity entity.CountryEntity) (*dbtype.Node, error)
	Update(ctx context.Context, entity entity.CountryEntity) (*dbtype.Node, error)
	GetCountryByCodeA3(ctx context.Context, codeA3 string) (*dbtype.Node, error)
	GetCountriesPaginated(ctx context.Context, skip, limit int) (*utils.DbNodesWithTotalCount, error)
	GetCountries(ctx context.Context) ([]*dbtype.Node, error)
}

type countryRepository struct {
	driver *neo4j.DriverWithContext
}

func NewCountryRepository(driver *neo4j.DriverWithContext) CountryRepository {
	return &countryRepository{
		driver: driver,
	}
}

func (r *countryRepository) Create(ctx context.Context, entity entity.CountryEntity) (*dbtype.Node, error) {
	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := "MATCH (t:Tenant {name:$tenant}) " +
		" MATCH (c:Country {id: randomUUID()})" +
		" ON CREATE SET c.name=$name, " +
		"				c.codeA2=$codeA2, " +
		"				c.codeA3=$codeA3, " +
		"				c.phoneCode=$phoneCode, " +
		" 				c.createdAt=datetime({timezone: 'UTC'}), " +
		" 				c.updatedAt=datetime({timezone: 'UTC'}), " +
		" RETURN c"

	result, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		queryResult, err := tx.Run(ctx, query,
			map[string]any{
				"name":      entity.Name,
				"codeA2":    entity.CodeA2,
				"codeA3":    entity.CodeA3,
				"phoneCode": entity.PhoneCode,
			})
		return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
	})
	if err != nil {
		return nil, err
	}
	return result.(*dbtype.Node), nil
}

func (r *countryRepository) Update(ctx context.Context, entity entity.CountryEntity) (*dbtype.Node, error) {
	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	dbNode, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, `
			MATCH (c:Country {id:$countryId})
			SET c.name=$name, 
				c.codeA2=$codeA2,
				c.codeA3=$codeA3,
				c.phoneCode=$phoneCode,
				c.updatedAt=datetime({timezone: 'UTC'}),
			RETURN c`,
			map[string]any{
				"countryId": entity.Id,
				"codeA2":    entity.CodeA2,
				"codeA3":    entity.CodeA3,
				"phoneCode": entity.PhoneCode,
			})
		return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
	})
	if err != nil {
		return nil, err
	}
	return dbNode.(*dbtype.Node), nil
}

func (r *countryRepository) GetCountryByCodeA3(ctx context.Context, codeA3 string) (*dbtype.Node, error) {
	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	dbRecord, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, "MATCH (c:Country {codeA3:$codeA3} ) RETURN c",
			map[string]any{
				"codeA3": codeA3,
			})
		if err != nil {
			return nil, err
		}
		return queryResult.Single(ctx)
	})
	return utils.NodePtr(dbRecord.(*db.Record).Values[0].(dbtype.Node)), err
}

func (r *countryRepository) GetCountriesPaginated(ctx context.Context, skip, limit int) (*utils.DbNodesWithTotalCount, error) {
	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	dbNodesWithTotalCount := new(utils.DbNodesWithTotalCount)

	dbRecords, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, "MATCH (c:Country) RETURN count(c) as count", map[string]any{})
		if err != nil {
			return nil, err
		}
		count, _ := queryResult.Single(ctx)
		dbNodesWithTotalCount.Count = count.Values[0].(int64)

		params := map[string]any{
			"skip":  skip,
			"limit": limit,
		}

		queryResult, err = tx.Run(ctx, "MATCH (c:Country) RETURN c ORDER BY c.name SKIP $skip LIMIT $limit", params)
		return queryResult.Collect(ctx)
	})
	if err != nil {
		return nil, err
	}
	for _, v := range dbRecords.([]*neo4j.Record) {
		dbNodesWithTotalCount.Nodes = append(dbNodesWithTotalCount.Nodes, utils.NodePtr(v.Values[0].(neo4j.Node)))
	}
	return dbNodesWithTotalCount, nil
}

func (r *countryRepository) GetCountries(ctx context.Context) ([]*dbtype.Node, error) {
	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	dbRecords, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, "MATCH (c:Country) RETURN c ORDER BY c.name", map[string]any{})
		if err != nil {
			return nil, err
		}
		return queryResult.Collect(ctx)
	})

	if err != nil {
		return nil, err
	}
	dbNodes := []*dbtype.Node{}
	for _, v := range dbRecords.([]*neo4j.Record) {
		if v.Values[0] != nil {
			dbNodes = append(dbNodes, utils.NodePtr(v.Values[0].(dbtype.Node)))
		}
	}
	return dbNodes, err
}
