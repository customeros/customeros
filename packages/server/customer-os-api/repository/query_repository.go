package repository

import (
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils"
)

type QueryRepository interface {
	GetOrganizationsAndContacts(session neo4j.Session, tenant string, skip int, limit int) (*utils.PairDbNodesWithTotalCount, error)
}

type queryRepository struct {
	driver *neo4j.Driver
}

func NewQueryRepository(driver *neo4j.Driver) QueryRepository {
	return &queryRepository{
		driver: driver,
	}
}

func (r *queryRepository) GetOrganizationsAndContacts(session neo4j.Session, tenant string, skip int, limit int) (*utils.PairDbNodesWithTotalCount, error) {
	result := new(utils.PairDbNodesWithTotalCount)
	dbRecords, err := session.ReadTransaction(func(tx neo4j.Transaction) (any, error) {
		queryResult, err := tx.Run(`
		CALL {
		  MATCH (t:Tenant {name:$tenant})--(o:Organization)
		  MATCH (t)--(c:Contact)
		  MATCH (o)--(c)
		  RETURN count(o) as o, count(c) as c
		  UNION
		  MATCH (t:Tenant {name:$tenant})--(o:Organization)
		  WHERE NOT (o)-[:LINKED]-(:Contact)
		  RETURN count(o) as o, 0 as c
		  UNION
		  MATCH (t:Tenant {name:$tenant})--(c:Contact)
		  WHERE NOT (c)-[:LINKED]-(:Organization)
		  RETURN 0 as o, count(c) as c
		}
		RETURN max(o), max (c)`,
			map[string]any{
				"tenant": tenant,
			})
		if err != nil {
			return nil, err
		}
		countRecord, err := queryResult.Single()
		if err != nil {
			return nil, err
		}
		result.Count = max(countRecord.Values[0].(int64), countRecord.Values[1].(int64))

		if queryResult, err := tx.Run(`
		CALL {
		  MATCH (t:Tenant {name:$tenant})--(o:Organization)
		  MATCH (t)--(c:Contact)
		  MATCH (o)--(c)
		  RETURN o, c
		  UNION
		  MATCH (t:Tenant {name:$tenant})--(o:Organization)
		  WHERE NOT (o)-[:LINKED]-(:Contact)
		  RETURN o, null as c
		  UNION
		  MATCH (t:Tenant {name:$tenant})--(c:Contact)
		  WHERE NOT (c)-[:LINKED]-(:Organization)
		  RETURN null as o, c
		}
		RETURN o, c
		SKIP $skip LIMIT $limit`,
			map[string]any{
				"tenant": tenant,
				"skip":   skip,
				"limit":  limit,
			}); err != nil {
			return nil, err
		} else {
			return queryResult.Collect()
		}
	})
	if err != nil {
		return nil, err
	}

	for _, v := range dbRecords.([]*neo4j.Record) {
		pair := new(utils.Pair[*dbtype.Node, *dbtype.Node])
		if v.Values[0] != nil {
			node := v.Values[0].(dbtype.Node)
			pair.First = &node
		}
		if v.Values[1] != nil {
			node := v.Values[1].(dbtype.Node)
			pair.Second = &node
		}

		result.Pairs = append(result.Pairs, pair)
	}
	return result, err
}

func max(a int64, b int64) int64 {
	if a < b {
		return b
	}
	return a
}
