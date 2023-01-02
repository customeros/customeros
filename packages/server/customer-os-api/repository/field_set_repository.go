package repository

import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils"
)

type FieldSetRepository interface {
	LinkWithFieldSetTemplateInTx(tx neo4j.Transaction, tenant, fieldSetId, templateId string) error
	MergeFieldSetToContactInTx(tx neo4j.Transaction, tenant, contactId string, entity entity.FieldSetEntity) (*dbtype.Node, *dbtype.Relationship, error)

	UpdateForContactInTx(tx neo4j.Transaction, tenant, contactId string, entity entity.FieldSetEntity) (*dbtype.Node, *dbtype.Relationship, error)
}

type fieldSetRepository struct {
	driver *neo4j.Driver
}

func NewFieldSetRepository(driver *neo4j.Driver) FieldSetRepository {
	return &fieldSetRepository{
		driver: driver,
	}
}

func (r *fieldSetRepository) LinkWithFieldSetTemplateInTx(tx neo4j.Transaction, tenant, fieldSetId, templateId string) error {
	txResult, err := tx.Run(`
			MATCH (f:FieldSet {id:$fieldSetId})<-[:HAS_COMPLEX_PROPERTY]-(c:Contact)-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}),
					(c)-[:IS_DEFINED_BY]->(e:EntityTemplate)-[:CONTAINS]->(d:FieldSetTemplate {id:$templateId})
			MERGE (f)-[r:IS_DEFINED_BY]->(d)
			RETURN r`,
		map[string]any{
			"templateId": templateId,
			"fieldSetId": fieldSetId,
			"tenant":     tenant,
		})
	if err != nil {
		return err
	}
	_, err = txResult.Single()
	return err
}

func (r *fieldSetRepository) MergeFieldSetToContactInTx(tx neo4j.Transaction, tenant, contactId string, entity entity.FieldSetEntity) (*dbtype.Node, *dbtype.Relationship, error) {
	query := "MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}) " +
		" MERGE (f:FieldSet {name: $name})<-[r:HAS_COMPLEX_PROPERTY]-(c) " +
		" ON CREATE SET f.id=randomUUID(), r.added=datetime({timezone: 'UTC'}), f:%s " +
		" RETURN f, r"
	queryResult, err := tx.Run(fmt.Sprintf(query, "FieldSet_"+tenant),
		map[string]interface{}{
			"tenant":    tenant,
			"contactId": contactId,
			"name":      entity.Name,
		})
	return utils.ExtractSingleRecordNodeAndRelationship(queryResult, err)
}

func (r *fieldSetRepository) UpdateForContactInTx(tx neo4j.Transaction, tenant, contactId string, entity entity.FieldSetEntity) (*dbtype.Node, *dbtype.Relationship, error) {
	queryResult, err := tx.Run(`
			MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}),
					(c)-[r:HAS_COMPLEX_PROPERTY]->(s:FieldSet {id:$fieldSetId})
            SET s.name=$name
			RETURN s, r`,
		map[string]interface{}{
			"tenant":     tenant,
			"contactId":  contactId,
			"fieldSetId": entity.Id,
			"name":       entity.Name,
		})
	return utils.ExtractSingleRecordNodeAndRelationship(queryResult, err)
}
