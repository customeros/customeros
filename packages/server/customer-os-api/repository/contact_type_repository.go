package repository

import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils"
)

type ContactTypeRepository interface {
	Create(tenant string, contactType *entity.ContactTypeEntity) (*dbtype.Node, error)
	Update(tenant string, contactType *entity.ContactTypeEntity) (*dbtype.Node, error)
	Delete(tenant string, id string) error
	FindAll(tenant string) ([]*dbtype.Node, error)
	FindForContact(tenant, contactId string) (*dbtype.Node, error)
}

type contactTypeRepository struct {
	driver *neo4j.Driver
}

func NewContactTypeRepository(driver *neo4j.Driver) ContactTypeRepository {
	return &contactTypeRepository{
		driver: driver,
	}
}

func (r *contactTypeRepository) Create(tenant string, contactType *entity.ContactTypeEntity) (*dbtype.Node, error) {
	session := utils.NewNeo4jWriteSession(*r.driver)
	defer session.Close()

	query := "MATCH (t:Tenant {name:$tenant})" +
		" MERGE (t)<-[:CONTACT_TYPE_BELONGS_TO_TENANT]-(c:ContactType {id:randomUUID()})" +
		" ON CREATE SET c.name=$name, c:%s" +
		" RETURN c"

	if result, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		queryResult, err := tx.Run(fmt.Sprintf(query, "ContactType_"+tenant),
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

func (r *contactTypeRepository) Update(tenant string, contactType *entity.ContactTypeEntity) (*dbtype.Node, error) {
	session := utils.NewNeo4jWriteSession(*r.driver)
	defer session.Close()

	if result, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		queryResult, err := tx.Run(`
			MATCH (t:Tenant {name:$tenant})<-[:CONTACT_TYPE_BELONGS_TO_TENANT]-(c:ContactType {id:$id})
			SET c.name=$name
			RETURN c`,
			map[string]any{
				"tenant": tenant,
				"id":     contactType.Id,
				"name":   contactType.Name,
			})
		return utils.ExtractSingleRecordFirstValueAsNode(queryResult, err)
	}); err != nil {
		return nil, err
	} else {
		return result.(*dbtype.Node), nil
	}
}

func (r *contactTypeRepository) Delete(tenant string, id string) error {
	session := utils.NewNeo4jWriteSession(*r.driver)
	defer session.Close()

	if _, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		_, err := tx.Run(`
			MATCH (t:Tenant {name:$tenant})<-[r:CONTACT_TYPE_BELONGS_TO_TENANT]-(c:ContactType {id:$id})
			DELETE r, c`,
			map[string]any{
				"tenant": tenant,
				"id":     id,
			})
		return nil, err
	}); err != nil {
		return err
	} else {
		return nil
	}
}

func (r *contactTypeRepository) FindAll(tenant string) ([]*dbtype.Node, error) {
	session := utils.NewNeo4jReadSession(*r.driver)
	defer session.Close()

	records, err := session.ReadTransaction(func(tx neo4j.Transaction) (any, error) {
		if queryResult, err := tx.Run(`
			MATCH (t:Tenant {name:$tenant})<-[:CONTACT_TYPE_BELONGS_TO_TENANT]-(c:ContactType)
			RETURN c ORDER BY c.name`,
			map[string]any{
				"tenant": tenant,
			}); err != nil {
			return nil, err
		} else {
			return queryResult.Collect()
		}
	})
	contactTypeDbNodes := []*dbtype.Node{}
	for _, v := range records.([]*neo4j.Record) {
		contactTypeDbNodes = append(contactTypeDbNodes, utils.NodePtr(v.Values[0].(dbtype.Node)))
	}
	return contactTypeDbNodes, err
}

func (r *contactTypeRepository) FindForContact(tenant, contactId string) (*dbtype.Node, error) {
	session := utils.NewNeo4jReadSession(*r.driver)
	defer session.Close()

	dbRecords, err := session.ReadTransaction(func(tx neo4j.Transaction) (any, error) {
		if queryResult, err := tx.Run(`
			MATCH (t:Tenant {name:$tenant})<-[:CONTACT_BELONGS_TO_TENANT]-(c:Contact {id:$contactId})-[:IS_OF_TYPE]->(o:ContactType)
			RETURN o`,
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
	} else if len(dbRecords.([]*neo4j.Record)) == 0 {
		return nil, nil
	} else {
		return utils.NodePtr(dbRecords.([]*neo4j.Record)[0].Values[0].(dbtype.Node)), nil
	}
}
