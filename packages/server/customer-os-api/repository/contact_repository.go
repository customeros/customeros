package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/db"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type ContactRepository interface {
	// Deprecated, use events-platform
	Delete(ctx context.Context, session neo4j.SessionWithContext, tenant, contactId string) error
	// Deprecated, use events-platform
	SetOwner(ctx context.Context, tx neo4j.ManagedTransaction, tenant, contactId, userId string) error
	// Deprecated, use events-platform
	RemoveOwner(ctx context.Context, tx neo4j.ManagedTransaction, tenant, contactId string) error
	GetPaginatedContacts(ctx context.Context, session neo4j.SessionWithContext, tenant string, skip, limit int, filter *utils.CypherFilter, sort *utils.CypherSort) (*utils.DbNodesWithTotalCount, error)
	GetPaginatedContactsForOrganization(ctx context.Context, session neo4j.SessionWithContext, tenant, organizationId string, skip, limit int, filter *utils.CypherFilter, sort *utils.CypherSort) (*utils.DbNodesWithTotalCount, error)
	GetAllForJobRoles(ctx context.Context, tenant string, jobRoleIds []string) ([]*utils.DbNodeAndId, error)
	GetContactsForPhoneNumber(ctx context.Context, tenant, phoneNumber string) ([]*dbtype.Node, error)
	// Deprecated, use events-platform
	MergeContactPropertiesInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant string, primaryContactId, mergedContactId string, sourceOfTruth neo4jentity.DataSource) error
	// Deprecated, use events-platform
	MergeContactRelationsInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant string, primaryContactId, mergedContactId string) error
	// Deprecated, use events-platform
	UpdateMergedContactLabelsInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant string, mergedContactId string) error
	GetAllForEmails(ctx context.Context, tenant string, emailIds []string) ([]*utils.DbNodeAndId, error)
	GetAllForPhoneNumbers(ctx context.Context, tenant string, phoneNumberIds []string) ([]*utils.DbNodeAndId, error)
	// Deprecated, use events-platform
	Archive(ctx context.Context, tenant, contactId string) error
	// Deprecated, use events-platform
	RestoreFromArchive(ctx context.Context, tenant, contactId string) error
	GetById(ctx context.Context, tenant, contactId string) (*dbtype.Node, error)

	GetBillableContactStats(ctx context.Context) (*neo4j.Record, error)
}

type contactRepository struct {
	driver   *neo4j.DriverWithContext
	database string
}

func NewContactRepository(driver *neo4j.DriverWithContext, database string) ContactRepository {
	return &contactRepository{
		driver:   driver,
		database: database,
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
		queryResult, err := tx.Run(ctx, fmt.Sprintf("MATCH (:Tenant {name:$tenant})<-[:CONTACT_BELONGS_TO_TENANT]-(c:Contact) %s WITH * WHERE c.hide IS NULL OR c.hide = false RETURN count(c) as count", filterCypherStr),
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
				" WITH * WHERE c.hide IS NULL OR c.hide = false "+
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
			 WITH * WHERE c.hide IS NULL OR c.hide = false
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
				" WITH * WHERE c.hide IS NULL OR c.hide = false "+
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
			OPTIONAL MATCH (c)-[:WORKS_AS]->(j:JobRole)
			OPTIONAL MATCH (c)--(soc:Social)
			OPTIONAL MATCH (c)-[:HAS_ACTION]->(pv:PageView)
			OPTIONAL MATCH (c)--(alt:AlternateContact)
            DETACH DELETE alt, f, fs, j, pv, soc, c`,
			map[string]interface{}{
				"contactId": contactId,
				"tenant":    tenant,
			})
		return nil, err
	})
	return err
}

func (r *contactRepository) GetAllForJobRoles(ctx context.Context, tenant string, jobRoleIds []string) ([]*utils.DbNodeAndId, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationRepository.GetAllForJobRoles")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)
	span.LogFields(log.String("jobRoleIds", fmt.Sprintf("%v", jobRoleIds)))

	query := `MATCH (t:Tenant {name:$tenant})<-[:CONTACT_BELONGS_TO_TENANT]-(c:Contact)-[:WORKS_AS]->(j:JobRole)
				WHERE j.id IN $jobRoleIds
				RETURN c, j.id`

	span.LogFields(log.String("query", query))

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)
	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, query,
			map[string]any{
				"tenant":     tenant,
				"jobRoleIds": jobRoleIds,
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

func (r *contactRepository) MergeContactPropertiesInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant string, primaryContactId, mergedContactId string, sourceOfTruth neo4jentity.DataSource) error {
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
				primary.username = CASE WHEN primary.username is null OR primary.username = '' THEN merged.username ELSE primary.username END, 
				primary.prefix = CASE WHEN primary.prefix is null OR primary.prefix = '' THEN merged.prefix ELSE primary.prefix END, 
				primary.sourceOfTruth=$sourceOfTruth,
				primary.updatedAt = datetime()
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

	// merge job roles that not sharing same organization
	if _, err := tx.Run(ctx, matchQuery+
		` WITH primary, merged 
		MATCH (merged)-[rel:WORKS_AS]->(jb:JobRole)
		WHERE NOT EXISTS ((primary)-[:WORKS_AS]->(:JobRole)-[:ROLE_IN]->(:Organization)<-[:ROLE_IN]-(jb))
		MERGE (primary)-[newRel:WORKS_AS]->(jb) 
		ON CREATE SET 	newRel.mergedFrom = $mergedContactId, 
						newRel.createdAt = $now
					SET	rel.merged=true, 
						jb.primary=false`, params); err != nil {
		return err
	}

	// merge job roles that share same organization (happens when merging contacts within same organization)
	// merge SENT_TO relationships
	if _, err := tx.Run(ctx, matchQuery+
		` WITH primary, merged 
		MATCH (merged)-[:WORKS_AS]->(jb:JobRole)
		MATCH (primary)-[:WORKS_AS]->(primaryJb:JobRole)-[:ROLE_IN]->(:Organization)<-[:ROLE_IN]-(jb)
		MATCH (jb)<-[rel:SENT_TO]-(ie:InteractionEvent)
		MERGE (primaryJb)<-[newRel:SENT_TO]-(ie) 
		ON CREATE SET 	newRel.mergedFrom = $mergedContactId, 
						newRel.createdAt = $now,
						newRel.type = rel.type
					SET	rel.merged=true`, params); err != nil {
		return err
	}

	// merge job roles that share same organization (happens when merging contacts within same organization)
	// merge SENT_BY relationships
	if _, err := tx.Run(ctx, matchQuery+
		` WITH primary, merged 
		MATCH (merged)-[:WORKS_AS]->(jb:JobRole)
		MATCH (primary)-[:WORKS_AS]->(primaryJb:JobRole)-[:ROLE_IN]->(:Organization)<-[:ROLE_IN]-(jb)
		MATCH (jb)<-[rel:SENT_BY]-(ie:InteractionEvent)
		MERGE (primaryJb)<-[newRel:SENT_BY]-(ie) 
		ON CREATE SET 	newRel.mergedFrom = $mergedContactId, 
						newRel.createdAt = $now,
						newRel.type = rel.type
					SET	rel.merged=true`, params); err != nil {
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
			SET c.archived=true, c.updatedAt=datetime(), c:Contact_%s
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

func (r *contactRepository) GetBillableContactStats(ctx context.Context) (*neo4j.Record, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactRepository.GetBillableContactStats")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	query := `MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(o:Organization)
			OPTIONAL MATCH (o)<-[:ROLE_IN]-(j:JobRole)<-[:WORKS_AS]-(c:Contact) 
WITH o.hide as hide, count(distinct(o)) as orgs, count(DISTINCT(c)) as contacts
WITH
  CASE WHEN hide = true THEN orgs ELSE 0 END AS grey_orgs,
  CASE WHEN hide = true THEN 0 ELSE orgs END AS white_orgs,
  CASE WHEN hide = true THEN contacts ELSE 0 END AS grey_contacts,
  CASE WHEN hide = true THEN 0 ELSE contacts END AS white_contacts
RETURN sum(white_orgs) AS whiteOrgs, sum(white_contacts) AS whiteContacts, sum(grey_orgs) AS greyOrgs, sum(grey_contacts) AS greyContacts`
	params := map[string]any{
		"tenant": common.GetTenantFromContext(ctx),
	}
	span.LogFields(log.String("query", query), log.Object("params", params))

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)
	dbRecord, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, query, params)
		if err != nil {
			return nil, err
		}
		return queryResult.Single(ctx)
	})
	return dbRecord.(*db.Record), err
}

func (r *contactRepository) GetById(ctx context.Context, tenant, contactId string) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactRepository.GetById")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	query := `MATCH (:Tenant {name:$tenant})<-[:CONTACT_BELONGS_TO_TENANT]-(c:Contact {id:$contactId}) RETURN c`
	params := map[string]any{
		"tenant":    tenant,
		"contactId": contactId,
	}
	span.LogFields(log.String("query", query), log.Object("params", params))

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	dbRecord, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, query, params)
		return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
	})
	if err != nil {
		return nil, err
	}
	return dbRecord.(*dbtype.Node), err
}
