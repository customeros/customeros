package repository

import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils"
	"time"
)

type ActionRepository interface {
	GetContactActions(session neo4j.Session, tenant, contactId string, from, to time.Time, labels []string) ([]*dbtype.Node, error)
}

type actionRepository struct {
	driver *neo4j.Driver
}

func NewActionRepository(driver *neo4j.Driver) ActionRepository {
	return &actionRepository{
		driver: driver,
	}
}

func (r *actionRepository) GetContactActions(session neo4j.Session, tenant, contactId string, from, to time.Time, labels []string) ([]*dbtype.Node, error) {
	params := map[string]any{
		"tenant":    tenant,
		"contactId": contactId,
		"from":      from,
		"to":        to,
	}
	filterByTypeCypherFragment := ""
	if len(labels) > 0 {
		params["nodeLabels"] = labels
		filterByTypeCypherFragment = "AND size([label IN labels(a) WHERE label IN $nodeLabels | 1]) > 0"
	}
	query := fmt.Sprintf("MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(t:Tenant {name:$tenant}), "+
		" (c)-[:HAS_ACTION]->(a:Action) "+
		" WHERE a.startedAt >= datetime($from) AND a.startedAt <= datetime($to) "+
		" %s "+
		" RETURN a ORDER BY a.startedAt DESC", filterByTypeCypherFragment)

	records, err := session.ReadTransaction(func(tx neo4j.Transaction) (any, error) {
		queryResult, err := tx.Run(query, params)
		if err != nil {
			return nil, err
		}
		return queryResult.Collect()
	})
	if err != nil {
		return nil, err
	}
	actionDbNodes := []*dbtype.Node{}
	for _, v := range records.([]*neo4j.Record) {
		actionDbNodes = append(actionDbNodes, utils.NodePtr(v.Values[0].(dbtype.Node)))
	}
	return actionDbNodes, err
}
