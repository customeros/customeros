package repository

import (
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/customer-os-api/utils"
)

type UserRepository interface {
	FindOwnerForContact(tx neo4j.Transaction, tenant, contactId string) (*dbtype.Node, error)
}

type userRepository struct {
	driver *neo4j.Driver
	repos  *RepositoryContainer
}

func NewUserRepository(driver *neo4j.Driver, repos *RepositoryContainer) UserRepository {
	return &userRepository{
		driver: driver,
		repos:  repos,
	}
}

func (r *userRepository) FindOwnerForContact(tx neo4j.Transaction, tenant, contactId string) (*dbtype.Node, error) {
	if queryResult, err := tx.Run(`
			MATCH (t:Tenant {name:$tenant})<-[:CONTACT_BELONGS_TO_TENANT]-(c:Contact {id:$contactId})<-[:OWNS]-(u:User)
			RETURN u`,
		map[string]any{
			"tenant":    tenant,
			"contactId": contactId,
		}); err != nil {
		return nil, err
	} else {
		dbRecords, err := queryResult.Collect()
		if err != nil {
			return nil, err
		} else if len(dbRecords) == 0 {
			return nil, nil
		} else {
			return utils.NodePtr(dbRecords[0].Values[0].(dbtype.Node)), nil
		}
	}

}
