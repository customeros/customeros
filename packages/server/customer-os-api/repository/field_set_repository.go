package repository

import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"golang.org/x/net/context"
)

type FieldSetRepository interface {
	LinkWithFieldSetTemplateInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant, fieldSetId, templateId string) error
	MergeFieldSetToContactInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant, contactId string, entity entity.FieldSetEntity) (*dbtype.Node, error)

	UpdateFieldSetForContactInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant, contactId string, entity entity.FieldSetEntity) (*dbtype.Node, error)
	DeleteByIdFromContact(ctx context.Context, session neo4j.SessionWithContext, tenant, contactId, fieldSetId string) error
	FindAllForContact(ctx context.Context, session neo4j.SessionWithContext, tenant, contactId string) ([]*neo4j.Record, error)
}

type fieldSetRepository struct {
	driver *neo4j.DriverWithContext
}

func NewFieldSetRepository(driver *neo4j.DriverWithContext) FieldSetRepository {
	return &fieldSetRepository{
		driver: driver,
	}
}

func (r *fieldSetRepository) LinkWithFieldSetTemplateInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant, fieldSetId, templateId string) error {
	txResult, err := tx.Run(ctx, `
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
	_, err = txResult.Single(ctx)
	return err
}

func (r *fieldSetRepository) MergeFieldSetToContactInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant, contactId string, entity entity.FieldSetEntity) (*dbtype.Node, error) {
	query := "MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}) " +
		" MERGE (f:FieldSet {name: $name})<-[r:HAS_COMPLEX_PROPERTY]-(c) " +
		" ON CREATE SET f.id=randomUUID(), " +
		"				f.createdAt=$now, " +
		"				f.createdAt=$now, " +
		"				f.source=$source, " +
		"				f.sourceOfTruth=$sourceOfTruth, " +
		"				f:%s " +
		" RETURN f"
	queryResult, err := tx.Run(ctx, fmt.Sprintf(query, "FieldSet_"+tenant),
		map[string]interface{}{
			"tenant":        tenant,
			"contactId":     contactId,
			"name":          entity.Name,
			"source":        entity.Source,
			"sourceOfTruth": entity.SourceOfTruth,
			"now":           utils.Now(),
		})
	return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
}

func (r *fieldSetRepository) UpdateFieldSetForContactInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant, contactId string, entity entity.FieldSetEntity) (*dbtype.Node, error) {
	queryResult, err := tx.Run(ctx, `
			MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}),
					(c)-[r:HAS_COMPLEX_PROPERTY]->(s:FieldSet {id:$fieldSetId})
            SET s.name=$name, s.sourceOfTruth=$sourceOfTruth, s.updatedAt=$now
			RETURN s`,
		map[string]interface{}{
			"tenant":        tenant,
			"contactId":     contactId,
			"fieldSetId":    entity.Id,
			"name":          entity.Name,
			"sourceOfTruth": entity.SourceOfTruth,
			"now":           utils.Now(),
		})
	return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
}

func (r *fieldSetRepository) DeleteByIdFromContact(ctx context.Context, session neo4j.SessionWithContext, tenant, contactId, fieldSetId string) error {
	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, `
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

func (r *fieldSetRepository) FindAllForContact(ctx context.Context, session neo4j.SessionWithContext, tenant, contactId string) ([]*neo4j.Record, error) {
	records, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, `
				MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}),
              			(c)-[:HAS_COMPLEX_PROPERTY]->(s:FieldSet) 
				RETURN s ORDER BY s.name`,
			map[string]any{
				"contactId": contactId,
				"tenant":    tenant})
		if err != nil {
			return nil, err
		}
		return queryResult.Collect(ctx)
	})
	return records.([]*neo4j.Record), err
}
