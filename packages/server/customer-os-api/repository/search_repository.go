package repository

import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/db"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils"
)

type SearchRepository interface {
	SearchBasic(tenant, keyword string) ([]*db.Record, error)
}

type searchRepository struct {
	driver *neo4j.Driver
}

func NewSearchRepository(driver *neo4j.Driver) SearchRepository {
	return &searchRepository{
		driver: driver,
	}
}

func (r *searchRepository) SearchBasic(tenant, keyword string) ([]*db.Record, error) {
	session := utils.NewNeo4jReadSession(*r.driver)
	defer session.Close()

	params := map[string]any{
		"tenant":        tenant,
		"keyword":       fmt.Sprintf("%s~", keyword),
		"indexStandard": fmt.Sprintf("basicSearchStandard_%s", tenant),
		"indexSimple":   fmt.Sprintf("basicSearchSimple_%s", tenant),
		"limit":         50,
	}
	query := "CALL { " +
		" CALL db.index.fulltext.queryNodes($indexStandard, $keyword) YIELD node, score RETURN score, node, labels(node) as labels limit $limit " +
		" union" +
		" CALL db.index.fulltext.queryNodes($indexSimple, $keyword) YIELD node, score RETURN score, node, labels(node) as labels limit $limit " +
		"} " +
		" with labels, node, score order by score desc " +
		" with labels, node, collect(score) as scores " +
		" return labels, head(scores) as score, node order by score desc limit $limit"

	records, err := session.ReadTransaction(func(tx neo4j.Transaction) (any, error) {
		queryResult, err := tx.Run(query, params)
		if err != nil {
			return nil, err
		}
		return queryResult.Collect()
	})
	return records.([]*db.Record), err
}
