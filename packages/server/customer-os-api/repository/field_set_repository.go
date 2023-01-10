package repository

import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils"
	"time"
)

type FieldSetRepository interface {
	LinkWithFieldSetTemplateInTx(tx neo4j.Transaction, tenant, fieldSetId, templateId string) error
	MergeFieldSetToContactInTx(tx neo4j.Transaction, tenant, contactId string, entity entity.FieldSetEntity) (*dbtype.Node, error)

	UpdateFieldSetForContactInTx(tx neo4j.Transaction, tenant, contactId string, entity entity.FieldSetEntity) (*dbtype.Node, error)
	DeleteByIdFromContact(session neo4j.Session, tenant, contactId, fieldSetId string) error
	FindAllForContact(session neo4j.Session, tenant, contactId string) ([]*neo4j.Record, error)
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

func (r *fieldSetRepository) MergeFieldSetToContactInTx(tx neo4j.Transaction, tenant, contactId string, entity entity.FieldSetEntity) (*dbtype.Node, error) {
	query := "MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}) " +
		" MERGE (f:FieldSet {name: $name})<-[r:HAS_COMPLEX_PROPERTY]-(c) " +
		" ON CREATE SET f.id=randomUUID(), f.createdAt=$createdAt, f.source=$source, f.sourceOfTruth=$sourceOfTruth, f:%s " +
		" RETURN f"
	queryResult, err := tx.Run(fmt.Sprintf(query, "FieldSet_"+tenant),
		map[string]interface{}{
			"tenant":        tenant,
			"contactId":     contactId,
			"name":          entity.Name,
			"source":        entity.Source,
			"sourceOfTruth": entity.SourceOfTruth,
			"createdAt":     time.Now().UTC(),
		})
	return utils.ExtractSingleRecordFirstValueAsNode(queryResult, err)
}

func (r *fieldSetRepository) UpdateFieldSetForContactInTx(tx neo4j.Transaction, tenant, contactId string, entity entity.FieldSetEntity) (*dbtype.Node, error) {
	queryResult, err := tx.Run(`
			MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}),
					(c)-[r:HAS_COMPLEX_PROPERTY]->(s:FieldSet {id:$fieldSetId})
            SET s.name=$name, s.sourceOfTruth=$sourceOfTruth
			RETURN s`,
		map[string]interface{}{
			"tenant":        tenant,
			"contactId":     contactId,
			"fieldSetId":    entity.Id,
			"name":          entity.Name,
			"sourceOfTruth": entity.SourceOfTruth,
		})
	return utils.ExtractSingleRecordFirstValueAsNode(queryResult, err)
}

func (r *fieldSetRepository) DeleteByIdFromContact(session neo4j.Session, tenant, contactId, fieldSetId string) error {
	_, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		_, err := tx.Run(`
			MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}),
                  (c)-[:HAS_COMPLEX_PROPERTY]->(s:FieldSet {id:$fieldSetId}),
				  (s)-[:HAS_PROPERTY]->(f:CustomField)
            DETACH DELETE f, s`,
			map[string]any{
				"contactId":  contactId,
				"fieldSetId": fieldSetId,
				"tenant":     tenant,
			})
		return nil, err
	})
	return err
}

func (r *fieldSetRepository) FindAllForContact(session neo4j.Session, tenant, contactId string) ([]*neo4j.Record, error) {
	records, err := session.ReadTransaction(func(tx neo4j.Transaction) (any, error) {
		queryResult, err := tx.Run(`
				MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}),
              			(c)-[:HAS_COMPLEX_PROPERTY]->(s:FieldSet) 
				RETURN s ORDER BY s.name`,
			map[string]any{
				"contactId": contactId,
				"tenant":    tenant})
		if err != nil {
			return nil, err
		}
		return queryResult.Collect()
	})
	return records.([]*neo4j.Record), err
}
