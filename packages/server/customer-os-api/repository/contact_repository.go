package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
	"time"
)

type ContactRepository interface {
	Create(ctx context.Context, tx neo4j.ManagedTransaction, tenant string, newContact entity.ContactEntity) (*dbtype.Node, error)
	Update(ctx context.Context, tx neo4j.ManagedTransaction, tenant, contactId string, contactDtls *entity.ContactEntity) (*dbtype.Node, error)
	Delete(ctx context.Context, session neo4j.SessionWithContext, tenant, contactId string) error
	SetOwner(ctx context.Context, tx neo4j.ManagedTransaction, tenant, contactId, userId string) error
	RemoveOwner(ctx context.Context, tx neo4j.ManagedTransaction, tenant, contactId string) error
	LinkWithEntityTemplateInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant string, obj *model.CustomFieldEntityType, entityTemplateId string) error
	GetPaginatedContacts(ctx context.Context, session neo4j.SessionWithContext, tenant string, skip, limit int, filter *utils.CypherFilter, sort *utils.CypherSort) (*utils.DbNodesWithTotalCount, error)
	GetPaginatedContactsForOrganization(ctx context.Context, session neo4j.SessionWithContext, tenant, organizationId string, skip, limit int, filter *utils.CypherFilter, sort *utils.CypherSort) (*utils.DbNodesWithTotalCount, error)
	GetAllForConversation(ctx context.Context, session neo4j.SessionWithContext, tenant, conversationId string) ([]*dbtype.Node, error)
	GetContactForRole(ctx context.Context, session neo4j.SessionWithContext, tenant, roleId string) (*dbtype.Node, error)
	GetContactsForEmail(ctx context.Context, tenant, email string) ([]*dbtype.Node, error)
	GetContactsForPhoneNumber(ctx context.Context, tenant, phoneNumber string) ([]*dbtype.Node, error)
	AddTag(ctx context.Context, tenant, contactId, tagId string) (*dbtype.Node, error)
	RemoveTag(ctx context.Context, tenant, contactId, tagId string) (*dbtype.Node, error)
	AddOrganization(ctx context.Context, tenant, contactId, organizationId, source, appSource string) (*dbtype.Node, error)
	RemoveOrganization(ctx context.Context, tenant, contactId, organizationId string) (*dbtype.Node, error)
	MergeContactPropertiesInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant string, primaryContactId, mergedContactId string, sourceOfTruth entity.DataSource) error
	MergeContactRelationsInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant string, primaryContactId, mergedContactId string) error
	UpdateMergedContactLabelsInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant string, mergedContactId string) error
	GetAllForEmails(ctx context.Context, tenant string, emailIds []string) ([]*utils.DbNodeAndId, error)
	GetAllForPhoneNumbers(ctx context.Context, tenant string, phoneNumberIds []string) ([]*utils.DbNodeAndId, error)
	Archive(ctx context.Context, tenant, contactId string) error
	RestoreFromArchive(ctx context.Context, tenant, contactId string) error

	GetAllContactPhoneNumberRelationships(ctx context.Context, size int) ([]*neo4j.Record, error)
	GetAllContactEmailRelationships(ctx context.Context, size int) ([]*neo4j.Record, error)
}

type contactRepository struct {
	driver *neo4j.DriverWithContext
}

func NewContactRepository(driver *neo4j.DriverWithContext) ContactRepository {
	return &contactRepository{
		driver: driver,
	}
}

func (r *contactRepository) SetOwner(ctx context.Context, tx neo4j.ManagedTransaction, tenant, contactId, userId string) error {
	queryResult, err := tx.Run(ctx, `
			MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(t:Tenant {name:$tenant}),
					(u:User {id:$userId})-[:USER_BELONGS_TO_TENANT]->(t)
			MERGE (u)-[r:OWNS]->(c)
			RETURN r`,
		map[string]interface{}{
			"tenant":    tenant,
			"contactId": contactId,
			"userId":    userId,
		})
	if err != nil {
		return err
	}
	_, err = queryResult.Single(ctx)
	return err
}

func (r *contactRepository) RemoveOwner(ctx context.Context, tx neo4j.ManagedTransaction, tenant, contactId string) error {
	_, err := tx.Run(ctx, `
			MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(t:Tenant {name:$tenant}),
					(c)<-[r:OWNS]-()
			DELETE r`,
		map[string]interface{}{
			"tenant":    tenant,
			"contactId": contactId,
		})
	return err
}

func (r *contactRepository) Create(ctx context.Context, tx neo4j.ManagedTransaction, tenant string, newContact entity.ContactEntity) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactRepository.Create")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	var createdAt time.Time
	createdAt = utils.Now()
	if newContact.CreatedAt != nil {
		createdAt = *newContact.CreatedAt
	}

	query := `MATCH (t:Tenant {name:$tenant}) 
			MERGE (c:Contact {id:randomUUID()})-[:CONTACT_BELONGS_TO_TENANT]->(t) 
			ON CREATE SET 
		 		c.prefix=$prefix, 
				c.name=$name,
		 		c.firstName=$firstName, 
		 		c.lastName=$lastName, 
				c.description=$description, 
				c.timezone=$timezone, 
		 		c.createdAt=$createdAt, 
		 		c.updatedAt=$createdAt, 
		 		c.source=$source, 
		 		c.appSource=$appSource, 
		 		c.sourceOfTruth=$sourceOfTruth, 
		 		c:Contact_%s 
			RETURN c`

	if queryResult, err := tx.Run(ctx, fmt.Sprintf(query, tenant),
		map[string]interface{}{
			"tenant":        tenant,
			"prefix":        newContact.Prefix,
			"name":          newContact.Name,
			"firstName":     newContact.FirstName,
			"lastName":      newContact.LastName,
			"source":        newContact.Source,
			"sourceOfTruth": newContact.SourceOfTruth,
			"appSource":     newContact.AppSource,
			"description":   newContact.Description,
			"timezone":      newContact.Timezone,
			"createdAt":     createdAt,
		}); err != nil {
		return nil, err
	} else {
		return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
	}
}

func (r *contactRepository) Update(ctx context.Context, tx neo4j.ManagedTransaction, tenant, contactId string, contactDtls *entity.ContactEntity) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactRepository.Update")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	if queryResult, err := tx.Run(ctx, `
			MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant})
			SET c.firstName=$firstName,
				c.lastName=$lastName,
				c.name=$name,
				c.description=$description,
				c.timezone=$timezone,
				c.prefix=$prefix,
				c.updatedAt=$now,
				c.sourceOfTruth=$sourceOfTruth
			RETURN c`,
		map[string]interface{}{
			"tenant":        tenant,
			"contactId":     contactId,
			"name":          contactDtls.Name,
			"firstName":     contactDtls.FirstName,
			"lastName":      contactDtls.LastName,
			"description":   contactDtls.Description,
			"timezone":      contactDtls.Timezone,
			"prefix":        contactDtls.Prefix,
			"sourceOfTruth": string(contactDtls.SourceOfTruth),
			"now":           utils.Now(),
		}); err != nil {
		return nil, err
	} else {
		return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
	}
}

func (r *contactRepository) LinkWithEntityTemplateInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant string, obj *model.CustomFieldEntityType, entityTemplateId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactRepository.LinkWithEntityTemplateInTx")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	var rel string
	var extends model.EntityTemplateExtension
	if obj.EntityType == model.EntityTypeContact {
		rel = "CONTACT_BELONGS_TO_TENANT"
		extends = model.EntityTemplateExtensionContact
	} else {
		rel = "ORGANIZATION_BELONGS_TO_TENANT"
		extends = model.EntityTemplateExtensionOrganization
	}

	queryResult, err := tx.Run(ctx, fmt.Sprintf(`
			MATCH (c:%s {id:$Id})-[:%s]->(:Tenant {name:$tenant})<-[:ENTITY_TEMPLATE_BELONGS_TO_TENANT]-(e:EntityTemplate {id:$entityTemplateId})
			WHERE e.extends=$extends
			MERGE (c)-[r:IS_DEFINED_BY]->(e)
			RETURN r`, obj.EntityType, rel),
		map[string]any{
			"entityTemplateId": entityTemplateId,
			"Id":               obj.ID,
			"tenant":           tenant,
			"extends":          extends,
		})
	if err != nil {
		return err
	}
	_, err = queryResult.Single(ctx)
	return err
}

func (r *contactRepository) GetPaginatedContacts(ctx context.Context, session neo4j.SessionWithContext, tenant string, skip, limit int, filter *utils.CypherFilter, sort *utils.CypherSort) (*utils.DbNodesWithTotalCount, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactRepository.GetPaginatedContacts")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	dbNodesWithTotalCount := new(utils.DbNodesWithTotalCount)

	dbRecords, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		filterCypherStr, filterParams := filter.CypherFilterFragment("c")
		countParams := map[string]any{
			"tenant": tenant,
		}
		utils.MergeMapToMap(filterParams, countParams)
		queryResult, err := tx.Run(ctx, fmt.Sprintf("MATCH (:Tenant {name:$tenant})<-[:CONTACT_BELONGS_TO_TENANT]-(c:Contact) %s RETURN count(c) as count", filterCypherStr),
			countParams)
		if err != nil {
			return nil, err
		}
		count, _ := queryResult.Single(ctx)
		dbNodesWithTotalCount.Count = count.Values[0].(int64)

		params := map[string]any{
			"tenant": tenant,
			"skip":   skip,
			"limit":  limit,
		}
		utils.MergeMapToMap(filterParams, params)

		queryResult, err = tx.Run(ctx, fmt.Sprintf(
			"MATCH (:Tenant {name:$tenant})<-[:CONTACT_BELONGS_TO_TENANT]-(c:Contact) "+
				" %s "+
				" RETURN c "+
				" %s "+
				" SKIP $skip LIMIT $limit", filterCypherStr, sort.SortingCypherFragment("c")),
			params)
		return queryResult.Collect(ctx)
	})
	if err != nil {
		return nil, err
	}
	for _, v := range dbRecords.([]*neo4j.Record) {
		dbNodesWithTotalCount.Nodes = append(dbNodesWithTotalCount.Nodes, utils.NodePtr(v.Values[0].(neo4j.Node)))
	}
	return dbNodesWithTotalCount, nil
}

func (r *contactRepository) GetPaginatedContactsForOrganization(ctx context.Context, session neo4j.SessionWithContext, tenant, organizationId string, skip, limit int, filter *utils.CypherFilter, sort *utils.CypherSort) (*utils.DbNodesWithTotalCount, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactRepository.GetPaginatedContactsForOrganization")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	dbNodesWithTotalCount := new(utils.DbNodesWithTotalCount)

	dbRecords, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		filterCypherStr, filterParams := filter.CypherFilterFragment("c")
		countParams := map[string]any{
			"tenant":         tenant,
			"organizationId": organizationId,
		}
		utils.MergeMapToMap(filterParams, countParams)

		query := `MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization {id:$organizationId})<-[:ROLE_IN]-(:JobRole)<-[:WORKS_AS]-(c:Contact) 
			 %s 
			 RETURN count(distinct(c)) as count`

		queryResult, err := tx.Run(ctx, fmt.Sprintf(query, filterCypherStr),
			countParams)
		if err != nil {
			return nil, err
		}
		count, _ := queryResult.Single(ctx)
		dbNodesWithTotalCount.Count = count.Values[0].(int64)

		params := map[string]any{
			"tenant":         tenant,
			"skip":           skip,
			"limit":          limit,
			"organizationId": organizationId,
		}
		utils.MergeMapToMap(filterParams, params)

		queryResult, err = tx.Run(ctx, fmt.Sprintf(
			"MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization {id:$organizationId})<-[:ROLE_IN]-(:JobRole)<-[:WORKS_AS]-(c:Contact) "+
				" %s "+
				" RETURN distinct(c) "+
				" %s "+
				" SKIP $skip LIMIT $limit", filterCypherStr, sort.SortingCypherFragment("c")),
			params)
		return queryResult.Collect(ctx)
	})
	if err != nil {
		return nil, err
	}
	for _, v := range dbRecords.([]*neo4j.Record) {
		dbNodesWithTotalCount.Nodes = append(dbNodesWithTotalCount.Nodes, utils.NodePtr(v.Values[0].(neo4j.Node)))
	}
	return dbNodesWithTotalCount, nil
}

func (r *contactRepository) Delete(ctx context.Context, session neo4j.SessionWithContext, tenant, contactId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactRepository.Delete")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, `
			MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant})
			OPTIONAL MATCH (c)-[:HAS_PROPERTY]->(f:CustomField)
			OPTIONAL MATCH (c)-[:HAS_COMPLEX_PROPERTY]->(fs:FieldSet)
			OPTIONAL MATCH (c)-[:WORKS_AS]->(j:JobRole)
			OPTIONAL MATCH (c)--(alt:AlternateContact)
            DETACH DELETE alt, f, fs, j, c`,
			map[string]interface{}{
				"contactId": contactId,
				"tenant":    tenant,
			})
		return nil, err
	})
	return err
}

func (r *contactRepository) GetAllForConversation(ctx context.Context, session neo4j.SessionWithContext, tenant, conversationId string) ([]*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactRepository.GetAllForConversation")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	dbRecords, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, `
			MATCH (t:Tenant {name:$tenant})<-[:CONTACT_BELONGS_TO_TENANT]-(c:Contact)-[:PARTICIPATES]->(o:Conversation {id:$conversationId})
			RETURN c`,
			map[string]any{
				"tenant":         tenant,
				"conversationId": conversationId,
			}); err != nil {
			return nil, err
		} else {
			return queryResult.Collect(ctx)
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

func (r *contactRepository) GetContactForRole(ctx context.Context, session neo4j.SessionWithContext, tenant, roleId string) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactRepository.GetContactForRole")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, `
			MATCH (:JobRole {id:$roleId})<-[:WORKS_AS]-(c:Contact)-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant})
			RETURN c`,
			map[string]any{
				"tenant": tenant,
				"roleId": roleId,
			}); err != nil {
			return nil, err
		} else {
			return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
		}
	})
	if err != nil {
		return nil, err
	}
	return result.(*dbtype.Node), nil
}

func (r *contactRepository) AddTag(ctx context.Context, tenant, contactId, tagId string) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactRepository.AddTag")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := "MATCH (t:Tenant {name:$tenant}), " +
		" (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(t), " +
		" (tag:Tag {id:$tagId})-[:TAG_BELONGS_TO_TENANT]->(t) " +
		" MERGE (c)-[rel:TAGGED]->(tag) " +
		" ON CREATE SET rel.taggedAt=$now, c.updatedAt=$now " +
		" RETURN c"

	if result, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, query,
			map[string]any{
				"tenant":    tenant,
				"contactId": contactId,
				"tagId":     tagId,
				"now":       utils.Now(),
			})
		return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
	}); err != nil {
		return nil, err
	} else {
		return result.(*dbtype.Node), nil
	}
}

func (r *contactRepository) RemoveTag(ctx context.Context, tenant, contactId, tagId string) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactRepository.RemoveTag")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := "MATCH (t:Tenant {name:$tenant}), " +
		" (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(t), " +
		" (tag:Tag {id:$tagId})-[:TAG_BELONGS_TO_TENANT]->(t) " +
		" OPTIONAL MATCH (c)-[rel:TAGGED]->(tag) " +
		" SET c.updatedAt = CASE WHEN rel is not null THEN $now ELSE c.updatedAt END " +
		" DELETE rel " +
		" RETURN c"

	if result, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, query,
			map[string]any{
				"tenant":    tenant,
				"contactId": contactId,
				"tagId":     tagId,
				"now":       utils.Now(),
			})
		return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
	}); err != nil {
		return nil, err
	} else {
		return result.(*dbtype.Node), nil
	}
}

func (r *contactRepository) AddOrganization(ctx context.Context, tenant, contactId, organizationId, source, appSource string) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactRepository.AddOrganization")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := `MATCH (t:Tenant {name:$tenant}), 
		 (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(t), 
		 (org:Organization {id:$organizationId})-[:ORGANIZATION_BELONGS_TO_TENANT]->(t) 
		 MERGE (c)-[:WORKS_AS]->(j:JobRole)-[:ROLE_IN]->(org)  
		 ON CREATE SET c.updatedAt=$now,
		 				j.id=randomUUID(), 
						j.source=$source, 
						j.sourceOfTruth=$source, 
						j.appSource=$appSource, 
						j.createdAt=$now, 
						j.updatedAt=$now,
						j.startedAt=$now,
						j:JobRole_%s
		 RETURN c`

	if result, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, fmt.Sprintf(query, tenant),
			map[string]any{
				"tenant":         tenant,
				"contactId":      contactId,
				"organizationId": organizationId,
				"now":            utils.Now(),
				"source":         source,
				"appSource":      appSource,
			})
		return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
	}); err != nil {
		return nil, err
	} else {
		return result.(*dbtype.Node), nil
	}
}

func (r *contactRepository) RemoveOrganization(ctx context.Context, tenant, contactId, organizationId string) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactRepository.RemoveOrganization")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := "MATCH (t:Tenant {name:$tenant}), " +
		" (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(t), " +
		" (org:Organization {id:$organizationId})-[:ORGANIZATION_BELONGS_TO_TENANT]->(t) " +
		" OPTIONAL MATCH (c)-[rel:WORKS_AS]->(j:JobRole)-[:ROLE_IN]->(org) " +
		" SET c.updatedAt = CASE WHEN rel is not null THEN $now ELSE c.updatedAt END " +
		" DETACH DELETE j " +
		" RETURN c"

	if result, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, query,
			map[string]any{
				"tenant":         tenant,
				"contactId":      contactId,
				"organizationId": organizationId,
				"now":            utils.Now(),
			})
		return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
	}); err != nil {
		return nil, err
	} else {
		return result.(*dbtype.Node), nil
	}
}

func (r *contactRepository) GetContactsForEmail(ctx context.Context, tenant, email string) ([]*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactRepository.GetContactsForEmail")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		if queryResult, err := tx.Run(ctx, `
			MATCH (:Tenant {name:$tenant})<-[:CONTACT_BELONGS_TO_TENANT]-(c:Contact)-[:HAS]->(e:Email) 
			WHERE e.email=$email OR e.rawEmail=$email
			RETURN DISTINCT c`,
			map[string]interface{}{
				"email":  email,
				"tenant": tenant,
			}); err != nil {
			return nil, err
		} else {
			return utils.ExtractAllRecordsFirstValueAsDbNodePtrs(ctx, queryResult, err)
		}
	})
	if err != nil {
		return nil, err
	}
	return result.([]*dbtype.Node), err
}

func (r *contactRepository) GetContactsForPhoneNumber(ctx context.Context, tenant, phoneNumber string) ([]*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactRepository.GetContactsForPhoneNumber")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		if queryResult, err := tx.Run(ctx, `
			MATCH (:Tenant {name:$tenant})<-[:CONTACT_BELONGS_TO_TENANT]-(c:Contact)-[:HAS]->(p:PhoneNumber) 
			WHERE p.e164=$phoneNumber OR p.rawPhoneNumber=$phoneNumber
			RETURN DISTINCT c`,
			map[string]interface{}{
				"phoneNumber": phoneNumber,
				"tenant":      tenant,
			}); err != nil {
			return nil, err
		} else {
			return utils.ExtractAllRecordsFirstValueAsDbNodePtrs(ctx, queryResult, err)
		}
	})
	if err != nil {
		return nil, err
	}
	return result.([]*dbtype.Node), err
}

func (r *contactRepository) MergeContactPropertiesInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant string, primaryContactId, mergedContactId string, sourceOfTruth entity.DataSource) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactRepository.MergeContactPropertiesInTx")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	_, err := tx.Run(ctx, `
			MATCH (t:Tenant {name:$tenant})<-[:CONTACT_BELONGS_TO_TENANT]-(primary:Contact {id:$primaryContactId}),
			(t)<-[:CONTACT_BELONGS_TO_TENANT]-(merged:Contact {id:$mergedContactId})
			SET primary.firstName = CASE WHEN primary.firstName is null OR primary.firstName = '' THEN merged.firstName ELSE primary.firstName END, 
				primary.lastName = CASE WHEN primary.lastName is null OR primary.lastName = '' THEN merged.lastName ELSE primary.lastName END, 
				primary.name = CASE WHEN primary.name is null OR primary.name = '' THEN merged.name ELSE primary.name END, 
				primary.description = CASE WHEN primary.description is null OR primary.description = '' THEN merged.description ELSE primary.description END, 
				primary.timezone = CASE WHEN primary.timezone is null OR primary.timezone = '' THEN merged.timezone ELSE primary.timezone END, 
				primary.profilePhotoUrl = CASE WHEN primary.profilePhotoUrl is null OR primary.profilePhotoUrl = '' THEN merged.profilePhotoUrl ELSE primary.profilePhotoUrl END, 
				primary.prefix = CASE WHEN primary.prefix is null OR primary.prefix = '' THEN merged.prefix ELSE primary.prefix END, 
				primary.sourceOfTruth=$sourceOfTruth,
				primary.updatedAt = $now
			`,
		map[string]any{
			"tenant":           tenant,
			"primaryContactId": primaryContactId,
			"mergedContactId":  mergedContactId,
			"sourceOfTruth":    string(sourceOfTruth),
			"now":              utils.Now(),
		})
	return err
}

func (r *contactRepository) MergeContactRelationsInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant string, primaryContactId, mergedContactId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactRepository.MergeContactRelationsInTx")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	matchQuery := "MATCH (t:Tenant {name:$tenant})<-[:CONTACT_BELONGS_TO_TENANT]-(primary:Contact {id:$primaryContactId}), " +
		"(t)<-[:CONTACT_BELONGS_TO_TENANT]-(merged:Contact {id:$mergedContactId})"

	params := map[string]any{
		"tenant":           tenant,
		"primaryContactId": primaryContactId,
		"mergedContactId":  mergedContactId,
		"now":              utils.Now(),
	}

	if _, err := tx.Run(ctx, matchQuery+" "+
		" WITH primary, merged "+
		" MATCH (merged)-[rel:PARTICIPATES]->(c:Conversation) "+
		" MERGE (primary)-[newRel:PARTICIPATES]->(c)"+
		" ON CREATE SET newRel.mergedFrom = $mergedContactId, "+
		"				newRel.createdAt = $now "+
		"			SET	rel.merged=true", params); err != nil {
		return err
	}

	if _, err := tx.Run(ctx, matchQuery+" "+
		" WITH primary, merged "+
		" MATCH (merged)-[rel:INITIATED]->(c:Conversation) "+
		" MERGE (primary)-[newRel:INITIATED]->(c)"+
		" ON CREATE SET newRel.mergedFrom = $mergedContactId, "+
		"				newRel.createdAt = $now "+
		"			SET	rel.merged=true", params); err != nil {
		return err
	}

	if _, err := tx.Run(ctx, matchQuery+" "+
		" WITH primary, merged "+
		" MATCH (merged)-[rel:HAS_ACTION]->(p:PageView) "+
		" MERGE (primary)-[newRel:HAS_ACTION]->(p)"+
		" ON CREATE SET newRel.mergedFrom = $mergedContactId, "+
		"				newRel.createdAt = $now "+
		"			SET	rel.merged=true", params); err != nil {
		return err
	}

	if _, err := tx.Run(ctx, matchQuery+
		" WITH primary, merged "+
		" MATCH (merged)-[rel:TAGGED]->(t:Tag) "+
		" MERGE (primary)-[newRel:TAGGED]->(t) "+
		" ON CREATE SET newRel.taggedAt=rel.taggedAt, "+
		"				newRel.mergedFrom = $mergedContactId, "+
		"				newRel.createdAt = $now "+
		"			SET	rel.merged=true", params); err != nil {
		return err
	}

	if _, err := tx.Run(ctx, matchQuery+
		" WITH primary, merged "+
		" MATCH (merged)-[rel:ASSOCIATED_WITH]->(loc:Location) "+
		" MERGE (primary)-[newRel:ASSOCIATED_WITH]->(loc) "+
		" ON CREATE SET newRel.mergedFrom = $mergedContactId, "+
		"				newRel.createdAt = $now "+
		"			SET	rel.merged=true", params); err != nil {
		return err
	}

	if _, err := tx.Run(ctx, matchQuery+
		" WITH primary, merged "+
		" MATCH (merged)-[rel:HAS_PROPERTY]->(c:CustomField) "+
		" MERGE (primary)-[newRel:HAS_PROPERTY]->(c) "+
		" ON CREATE SET newRel.mergedFrom = $mergedContactId, "+
		"				newRel.createdAt = $now "+
		"			SET	rel.merged=true", params); err != nil {
		return err
	}

	if _, err := tx.Run(ctx, matchQuery+
		" WITH primary, merged "+
		" MATCH (merged)<-[rel:OWNS]-(u:User) "+
		" MERGE (primary)<-[newRel:OWNS]-(u) "+
		" ON CREATE SET newRel.mergedFrom = $mergedContactId, "+
		"				newRel.createdAt = $now "+
		"			SET	rel.merged=true", params); err != nil {
		return err
	}

	if _, err := tx.Run(ctx, matchQuery+
		" WITH primary, merged "+
		" MATCH (merged)-[rel:WORKS_AS]->(jb:JobRole) "+
		" MERGE (primary)-[newRel:WORKS_AS]->(jb) "+
		" ON CREATE SET newRel.mergedFrom = $mergedContactId, "+
		"				newRel.createdAt = $now "+
		"			SET	rel.merged=true, "+
		"				jb.primary=false", params); err != nil {
		return err
	}

	if _, err := tx.Run(ctx, matchQuery+
		" WITH primary, merged "+
		" MATCH (merged)-[rel:HAS]->(e:Email) "+
		" MERGE (primary)-[newRel:HAS]->(e) "+
		" ON CREATE SET newRel.primary=false, "+
		"				newRel.label=rel.label, "+
		"				newRel.mergedFrom = $mergedContactId, "+
		"				newRel.createdAt = $now "+
		"			SET	rel.merged=true", params); err != nil {
		return err
	}

	if _, err := tx.Run(ctx, matchQuery+
		" WITH primary, merged "+
		" MATCH (merged)-[rel:HAS]->(p:PhoneNumber) "+
		" MERGE (primary)-[newRel:HAS]->(p) "+
		" ON CREATE SET newRel.primary=false, "+
		"				newRel.label=rel.label, "+
		"               newRel.mergedFrom = $mergedContactId, "+
		"				newRel.createdAt = $now "+
		"			SET	rel.merged=true", params); err != nil {
		return err
	}

	if _, err := tx.Run(ctx, matchQuery+
		" WITH primary, merged "+
		" MATCH (merged)-[rel:NOTED]->(n:Note) "+
		" MERGE (primary)-[newRel:NOTED]->(n) "+
		" ON CREATE SET newRel.mergedFrom = $mergedContactId, "+
		"				newRel.createdAt = $now "+
		"			SET	rel.merged=true", params); err != nil {
		return err
	}

	if _, err := tx.Run(ctx, matchQuery+
		" WITH primary, merged "+
		" MATCH (merged)-[rel:CREATED]->(n:Note) "+
		" MERGE (primary)-[newRel:CREATED]->(n) "+
		" ON CREATE SET newRel.mergedFrom = $mergedContactId, "+
		"				newRel.createdAt = $now "+
		"			SET	rel.merged=true", params); err != nil {
		return err
	}

	if _, err := tx.Run(ctx, matchQuery+
		" WITH primary, merged "+
		" MATCH (merged)<-[rel:SENT_BY]-(i:InteractionEvent) "+
		" MERGE (primary)<-[newRel:SENT_BY]-(i) "+
		" ON CREATE SET newRel.mergedFrom = $mergedContactId, "+
		"				newRel.createdAt = $now, "+
		"				newRel.type = rel.type "+
		"			SET	rel.merged=true", params); err != nil {
		return err
	}

	if _, err := tx.Run(ctx, matchQuery+
		" WITH primary, merged "+
		" MATCH (merged)<-[rel:SENT_TO]-(i:InteractionEvent) "+
		" MERGE (primary)<-[newRel:SENT_TO]-(i) "+
		" ON CREATE SET newRel.mergedFrom = $mergedContactId, "+
		"				newRel.createdAt = $now, "+
		"				newRel.type = rel.type "+
		"			SET	rel.merged=true", params); err != nil {
		return err
	}

	if _, err := tx.Run(ctx, matchQuery+
		" WITH primary, merged "+
		" MATCH (merged)<-[rel:ATTENDED_BY]-(s:InteractionSession) "+
		" MERGE (primary)<-[newRel:ATTENDED_BY]-(s) "+
		" ON CREATE SET newRel.mergedFrom = $mergedContactId, "+
		"				newRel.createdAt = $now, "+
		"				newRel.type = rel.type "+
		"			SET	rel.merged=true", params); err != nil {
		return err
	}

	if _, err := tx.Run(ctx, matchQuery+
		" WITH primary, merged "+
		" MATCH (merged)<-[rel:CREATED_BY]-(m:Meeting) "+
		" MERGE (primary)<-[newRel:CREATED_BY]-(m) "+
		" ON CREATE SET newRel.mergedFrom = $mergedContactId, "+
		"				newRel.createdAt = $now "+
		"			SET	rel.merged=true", params); err != nil {
		return err
	}

	if _, err := tx.Run(ctx, matchQuery+
		" WITH primary, merged "+
		" MATCH (merged)<-[rel:ATTENDED_BY]-(m:Meeting) "+
		" MERGE (primary)<-[newRel:ATTENDED_BY]-(m) "+
		" ON CREATE SET newRel.mergedFrom = $mergedContactId, "+
		"				newRel.createdAt = $now "+
		"			SET	rel.merged=true", params); err != nil {
		return err
	}

	if _, err := tx.Run(ctx, matchQuery+
		" WITH primary, merged "+
		" MATCH (merged)-[rel:IS_LINKED_WITH]->(ext:ExternalSystem) "+
		" MERGE (primary)-[newRel:IS_LINKED_WITH {externalId:rel.externalId}]->(ext) "+
		" ON CREATE SET newRel.syncDate=rel.syncDate, "+
		"				newRel.externalUrl=rel.externalUrl, "+
		"				newRel.mergedFrom = $mergedContactId, "+
		"				newRel.createdAt = $now "+
		"			SET	rel.merged=true", params); err != nil {
		return err
	}

	if _, err := tx.Run(ctx, matchQuery+
		" WITH primary, merged "+
		" MERGE (merged)-[rel:IS_MERGED_INTO]->(primary)"+
		" ON CREATE SET rel.mergedAt=$now ", params); err != nil {
		return err
	}

	return nil
}

func (r *contactRepository) UpdateMergedContactLabelsInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant string, mergedContactId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactRepository.UpdateMergedContactLabelsInTx")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	query := "MATCH (t:Tenant {name:$tenant})<-[:CONTACT_BELONGS_TO_TENANT]-(c:Contact {id:$contactId}) " +
		" SET c:MergedContact:%s " +
		" REMOVE c:Contact:%s"

	_, err := tx.Run(ctx, fmt.Sprintf(query, "MergedContact_"+tenant, "Contact_"+tenant),
		map[string]any{
			"tenant":    tenant,
			"contactId": mergedContactId,
		})
	return err
}

func (r *contactRepository) GetAllForEmails(ctx context.Context, tenant string, emailIds []string) ([]*utils.DbNodeAndId, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactRepository.GetAllForEmails")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, `
			MATCH (t:Tenant {name:$tenant})<-[:CONTACT_BELONGS_TO_TENANT]-(c:Contact)-[:HAS]->(e:Email)-[:EMAIL_ADDRESS_BELONGS_TO_TENANT]->(t)
			WHERE e.id IN $emailIds
			RETURN c, e.id as emailId ORDER BY c.firstName, c.lastName`,
			map[string]any{
				"tenant":   tenant,
				"emailIds": emailIds,
			}); err != nil {
			return nil, err
		} else {
			return utils.ExtractAllRecordsAsDbNodeAndId(ctx, queryResult, err)
		}
	})
	if err != nil {
		return nil, err
	}
	return result.([]*utils.DbNodeAndId), err
}

func (r *contactRepository) GetAllForPhoneNumbers(ctx context.Context, tenant string, phoneNumberIds []string) ([]*utils.DbNodeAndId, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactRepository.GetAllForPhoneNumbers")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, `
			MATCH (t:Tenant {name:$tenant})<-[:CONTACT_BELONGS_TO_TENANT]-(c:Contact)-[:HAS]->(p:PhoneNumber)-[:PHONE_NUMBER_BELONGS_TO_TENANT]->(t)
			WHERE p.id IN $phoneNumberIds
			RETURN c, p.id as phoneNumberId ORDER BY c.firstName, c.lastName`,
			map[string]any{
				"tenant":         tenant,
				"phoneNumberIds": phoneNumberIds,
			}); err != nil {
			return nil, err
		} else {
			return utils.ExtractAllRecordsAsDbNodeAndId(ctx, queryResult, err)
		}
	})
	if err != nil {
		return nil, err
	}
	return result.([]*utils.DbNodeAndId), err
}

func (r *contactRepository) GetAllContactPhoneNumberRelationships(ctx context.Context, size int) ([]*neo4j.Record, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactRepository.GetAllContactPhoneNumberRelationships")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	dbRecords, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, `
			MATCH (t:Tenant)<-[:CONTACT_BELONGS_TO_TENANT]-(c:Contact)-[rel:HAS]->(p:PhoneNumber)
 			WHERE (rel.syncedWithEventStore is null or rel.syncedWithEventStore=false)
			RETURN rel, c.id, p.id, t.name limit $size`,
			map[string]any{
				"size": size,
			}); err != nil {
			return nil, err
		} else {
			return queryResult.Collect(ctx)
		}
	})
	if err != nil {
		return nil, err
	}
	return dbRecords.([]*neo4j.Record), err
}

func (r *contactRepository) GetAllContactEmailRelationships(ctx context.Context, size int) ([]*neo4j.Record, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactRepository.GetAllContactEmailRelationships")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	dbRecords, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, `
			MATCH (t:Tenant)<-[:CONTACT_BELONGS_TO_TENANT]-(c:Contact)-[rel:HAS]->(e:Email)
 			WHERE (rel.syncedWithEventStore is null or rel.syncedWithEventStore=false)
			RETURN rel, c.id, e.id, t.name limit $size`,
			map[string]any{
				"size": size,
			}); err != nil {
			return nil, err
		} else {
			return queryResult.Collect(ctx)
		}
	})
	if err != nil {
		return nil, err
	}
	return dbRecords.([]*neo4j.Record), err
}

func (r *contactRepository) Archive(ctx context.Context, tenant, contactId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactRepository.Archive")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := `
			MATCH (c:Contact {id:$contactId})-[r:CONTACT_BELONGS_TO_TENANT]->(t:Tenant {name:$tenant})
			MERGE (c)-[newRel:ARCHIVED]->(t)
			SET c.archived=true, newRel.archivedAt=$now, c:ArchivedContact_%s
            DELETE r
			REMOVE c:Contact_%s
			`

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		_, err := tx.Run(ctx, fmt.Sprintf(query, tenant, tenant),
			map[string]interface{}{
				"contactId": contactId,
				"tenant":    tenant,
				"now":       utils.Now(),
			})

		return nil, err
	})
	return err
}

func (r *contactRepository) RestoreFromArchive(ctx context.Context, tenant, contactId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactRepository.RestoreFromArchive")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := `
			MATCH (c:Contact {id:$contactId})-[r:ARCHIVED]->(t:Tenant {name:$tenant})
			MERGE (c)-[newRel:CONTACT_BELONGS_TO_TENANT]->(t)
			SET c.archived=true, c.updatedAt=$now, c:Contact_%s
            DELETE r
			REMOVE c.archived
			REMOVE c:ArchivedContact_%s
			`

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		_, err := tx.Run(ctx, fmt.Sprintf(query, tenant, tenant),
			map[string]interface{}{
				"contactId": contactId,
				"tenant":    tenant,
				"now":       utils.Now(),
			})

		return nil, err
	})
	return err
}
