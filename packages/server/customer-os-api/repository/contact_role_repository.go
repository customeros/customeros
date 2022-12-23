package repository

import (
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils"
)

type ContactRoleRepository interface {
	GetRolesForContact(session neo4j.Session, tenant, contactId string) ([]*dbtype.Node, error)
}

type contactRoleRepository struct {
	driver *neo4j.Driver
}

func NewContactRoleRepository(driver *neo4j.Driver) ContactRoleRepository {
	return &contactRoleRepository{
		driver: driver,
	}
}

func (r *contactRoleRepository) GetRolesForContact(session neo4j.Session, tenant, contactId string) ([]*dbtype.Node, error) {
	records, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		queryResult, err := tx.Run(`
				MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}),
              			(c)-[:HAS_ROLE]->(r:Role) 
				RETURN r`,
			map[string]interface{}{
				"contactId": contactId,
				"tenant":    tenant,
			})
		if err != nil {
			return nil, err
		}
		return queryResult.Collect()
	})
	if err != nil {
		return nil, err
	}
	dbNodes := []*dbtype.Node{}
	for _, v := range records.([]*neo4j.Record) {
		if v.Values[0] != nil {
			dbNodes = append(dbNodes, utils.NodePtr(v.Values[0].(dbtype.Node)))
		}
	}
	return dbNodes, err
}
