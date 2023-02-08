package repository

import (
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils"
)

type LocationRepository interface {
	GetAllForContact(tenant, contactId string) ([]*dbtype.Node, error)
	GetAllForOrganization(tenant, organizationId string) ([]*dbtype.Node, error)
}

type locationRepository struct {
	driver *neo4j.Driver
}

func NewLocationRepository(driver *neo4j.Driver) LocationRepository {
	return &locationRepository{
		driver: driver,
	}
}

func (r *locationRepository) GetAllForContact(tenant, contactId string) ([]*dbtype.Node, error) {
	session := utils.NewNeo4jWriteSession(*r.driver)
	defer session.Close()

	result, err := session.ReadTransaction(func(tx neo4j.Transaction) (any, error) {
		if queryResult, err := tx.Run(`
			MATCH (:Tenant {name:$tenant})<-[:CONTACT_BELONGS_TO_TENANT]-(:Contact {id:$contactId})-[:ASSOCIATED_WITH]->(loc:Location)
			RETURN loc`,
			map[string]any{
				"tenant":    tenant,
				"contactId": contactId,
			}); err != nil {
			return nil, err
		} else {
			return utils.ExtractAllRecordsFirstValueAsNodePtrs(queryResult, err)
		}
	})
	return result.([]*dbtype.Node), err
}

func (r *locationRepository) GetAllForOrganization(tenant, organizationId string) ([]*dbtype.Node, error) {
	session := utils.NewNeo4jWriteSession(*r.driver)
	defer session.Close()

	result, err := session.ReadTransaction(func(tx neo4j.Transaction) (any, error) {
		if queryResult, err := tx.Run(`
			MATCH (:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(:Organization {id:$organizationId})-[:ASSOCIATED_WITH]->(loc:Location)
			RETURN loc`,
			map[string]any{
				"tenant":         tenant,
				"organizationId": organizationId,
			}); err != nil {
			return nil, err
		} else {
			return utils.ExtractAllRecordsFirstValueAsNodePtrs(queryResult, err)
		}
	})
	return result.([]*dbtype.Node), err
}
