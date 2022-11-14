package repository

import (
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/customer-os-api/utils"
)

type ContactTypeRepository interface {
	Create(tenant string, contactType *entity.ContactTypeEntity) (*dbtype.Node, error)
}

type contactTypeRepository struct {
	driver *neo4j.Driver
	repos  *RepositoryContainer
}

func NewContactTypeRepository(driver *neo4j.Driver, repos *RepositoryContainer) ContactTypeRepository {
	return &contactTypeRepository{
		driver: driver,
		repos:  repos,
	}
}

func (r *contactTypeRepository) Create(tenant string, contactType *entity.ContactTypeEntity) (*dbtype.Node, error) {
	session := (*r.driver).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	if result, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		queryResult, err := tx.Run(`
			MATCH (t:Tenant {name:$tenant})
			MERGE (t)-[:USES_CONTACT_TYPE]->(c:ContactType {id:randomUUID()})
			ON CREATE SET c.name=$name
			RETURN c`,
			map[string]any{
				"tenant": tenant,
				"name":   contactType.Name,
			})
		return utils.ExtractSingleRecordFirstValueAsNode(queryResult, err)
	}); err != nil {
		return nil, err
	} else {
		return result.(*dbtype.Node), nil
	}
}
