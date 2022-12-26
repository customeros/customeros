package repository

import (
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils"
)

type AddressRepository interface {
	FindAllForContact(session neo4j.Session, tenant, contactId string) ([]*dbtype.Node, error)
	FindAllForCompany(session neo4j.Session, tenant, companyId string) ([]*dbtype.Node, error)
}

type addressRepository struct {
	driver *neo4j.Driver
}

func NewAddressRepository(driver *neo4j.Driver) AddressRepository {
	return &addressRepository{
		driver: driver,
	}
}

func (r *addressRepository) FindAllForContact(session neo4j.Session, tenant, contactId string) ([]*dbtype.Node, error) {
	dbRecords, err := session.ReadTransaction(func(tx neo4j.Transaction) (any, error) {
		if queryResult, err := tx.Run(`
			MATCH (t:Tenant {name:$tenant})<-[:CONTACT_BELONGS_TO_TENANT]-(c:Contact {id:$contactId})-[:LOCATED_AT]->(a:Address)
			RETURN a`,
			map[string]any{
				"tenant":    tenant,
				"contactId": contactId,
			}); err != nil {
			return nil, err
		} else {
			return queryResult.Collect()
		}
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

func (r *addressRepository) FindAllForCompany(session neo4j.Session, tenant, companyId string) ([]*dbtype.Node, error) {
	dbRecords, err := session.ReadTransaction(func(tx neo4j.Transaction) (any, error) {
		if queryResult, err := tx.Run(`
			MATCH (t:Tenant {name:$tenant})<-[:COMPANY_BELONGS_TO_TENANT]-(c:Company {id:$companyId})-[:LOCATED_AT]->(a:Address)
			RETURN a`,
			map[string]any{
				"tenant":    tenant,
				"companyId": companyId,
			}); err != nil {
			return nil, err
		} else {
			return queryResult.Collect()
		}
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
