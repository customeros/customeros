package neo4j_repository

import (
	"fmt"
	neo4jenum "github.com/customeros/customeros/packages/server/customer-os-common-module/model"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/db"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"

	"golang.org/x/net/context"
)

type TenantAndEmailId struct {
	Tenant  string
	EmailId string
}

type EmailReadRepository interface {
	GetEmailIdIfExists(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, email string) (string, error)
	GetEmailForUser(ctx context.Context, tenant string, userId string) (*dbtype.Node, error)
	GetById(ctx context.Context, tenant, emailId string) (*dbtype.Node, error)
	GetFirstByEmail(ctx context.Context, tenant, email string) (*dbtype.Node, error)
	GetAllEmailNodesForLinkedEntityIds(ctx context.Context, tenant string, entityType neo4jenum.EntityType, entityIds []string) ([]*utils.DbNodeWithRelationAndId, error)
	GetPrimaryEmailNodesForLinkedEntityIds(ctx context.Context, tenant string, entityType neo4jenum.EntityType, entityIds []string) ([]*utils.DbNodeWithRelationAndId, error)
	GetEmailsForValidation(ctx context.Context, delayFromLastUpdateInSeconds, delayFromLastValidationAttemptInMinutes, limit int) ([]TenantAndEmailId, error)
	IsLinkedToEntityByEmailAddress(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, email, entityId string, entityType neo4jenum.EntityType) (bool, error)
	GetOrphanEmailNodes(ctx context.Context, limit, hoursFromLastUpdate int) ([]TenantAndEmailId, error)
	IsOrphanEmail(ctx context.Context, tenant, emailId string) (bool, error)
	GetTestEmailForFlows(ctx context.Context, tenant string) (string, error)
}

type emailReadRepository struct {
	driver   *neo4j.DriverWithContext
	database string
}

func NewEmailReadRepository(driver *neo4j.DriverWithContext, database string) EmailReadRepository {
	return &emailReadRepository{
		driver:   driver,
		database: database,
	}
}

func (r *emailReadRepository) prepareReadSession(ctx context.Context) neo4j.SessionWithContext {
	return utils.NewNeo4jReadSession(ctx, *r.driver, utils.WithDatabaseName(r.database))
}

func (r *emailReadRepository) GetEmailIdIfExists(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, email string) (string, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "EmailReadRepository.GetEmailIdIfExists")
	defer spans.Finish()

	spans.LogKV("email", email)

	cypher := fmt.Sprintf(`MATCH (e:Email_%s) WHERE (lower(e.email) = lower($email) OR lower(e.rawEmail) = lower($email)) AND $email <> '' RETURN e.id ORDER BY e.createdAt LIMIT 1`, tenant)
	params := map[string]any{
		"email": email,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	result, err := utils.ExecuteReadInTransaction(ctx, r.driver, r.database, tx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
			return nil, err
		} else {
			return queryResult.Collect(ctx)
		}
	})
	if err != nil {
		spans.TraceError(err)
		return "", err
	}
	if len(result.([]*db.Record)) == 0 {
		spans.LogKV("result.found", false)
		return "", nil
	}
	spans.LogKV("result", result.([]*db.Record)[0].Values[0].(string))
	return result.([]*db.Record)[0].Values[0].(string), err
}

func (r *emailReadRepository) GetEmailForUser(ctx context.Context, tenant string, userId string) (*dbtype.Node, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "EmailReadRepository.GetEmailForUser")
	defer spans.Finish()

	spans.LogKV("userId", userId)

	cypher := fmt.Sprintf("MATCH (e:Email_%s)<-[:HAS]-(u:User {id:$userId})-[:USER_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}) WHERE u:User_%s return e", tenant, tenant)
	params := map[string]any{
		"userId": userId,
		"tenant": tenant,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
			return nil, err
		} else {
			return utils.ExtractFirstRecordFirstValueAsDbNodePtr(ctx, queryResult, err)
		}
	})
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}

	if result == nil {
		spans.LogKV("result.found", false)
		return nil, nil
	}

	spans.LogKV("result.found", true)
	return result.(*dbtype.Node), nil
}

func (r *emailReadRepository) GetById(ctx context.Context, tenant, emailId string) (*dbtype.Node, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "EmailReadRepository.GetById")
	defer spans.Finish()

	spans.LogKV("emailId", emailId)

	cypher := `MATCH (:Tenant {name:$tenant})<-[:EMAIL_ADDRESS_BELONGS_TO_TENANT]-(e:Email {id:$emailId}) return e`
	params := map[string]any{
		"tenant":  tenant,
		"emailId": emailId,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	dbRecord, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
		return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
	})
	if err != nil {
		return nil, err
	}
	spans.LogKV("result.found", dbRecord != nil)
	return dbRecord.(*dbtype.Node), err
}

func (r *emailReadRepository) GetFirstByEmail(ctx context.Context, tenant, email string) (*dbtype.Node, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "EmailRepository.GetFirstByEmail")
	defer spans.Finish()

	spans.LogKV("email", email)

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:EMAIL_ADDRESS_BELONGS_TO_TENANT]-(e:Email) 
		       WHERE e.rawEmail = $email OR e.email = $email RETURN e ORDER BY e.createdAt LIMIT 1`
	params := map[string]any{
		"tenant": tenant,
		"email":  email,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	dbRecord, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
		return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
	})

	if err != nil && err.Error() == "Result contains no more records" {
		spans.LogKV("result.found", false)
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	if dbRecord == nil {
		spans.LogKV("result.found", false)
		return nil, nil
	}

	spans.LogKV("result.found", true)
	return dbRecord.(*dbtype.Node), err
}

func (r *emailReadRepository) GetAllEmailNodesForLinkedEntityIds(ctx context.Context, tenant string, entityType neo4jenum.EntityType, entityIds []string) ([]*utils.DbNodeWithRelationAndId, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "EmailReadRepository.GetAllEmailNodesForLinkedEntityIds")
	defer spans.Finish()

	cypher := ""
	switch entityType {
	case neo4jenum.CONTACT:
		cypher = `MATCH (t:Tenant {name:$tenant})<-[:CONTACT_BELONGS_TO_TENANT]-(entity:Contact)`
	case neo4jenum.USER:
		cypher = `MATCH (t:Tenant {name:$tenant})<-[:USER_BELONGS_TO_TENANT]-(entity:User)`
	case neo4jenum.ORGANIZATION:
		cypher = `MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(entity:Organization)`
	}
	cypher = cypher + `, (entity)-[rel:HAS]->(e:Email)-[:EMAIL_ADDRESS_BELONGS_TO_TENANT]->(t)
					WHERE entity.id IN $entityIds
					RETURN e, rel, entity.id ORDER BY e.email, e.rawEmail`
	params := map[string]any{
		"tenant":    tenant,
		"entityIds": entityIds,
	}

	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
			return nil, err
		} else {
			return utils.ExtractAllRecordsAsDbNodeWithRelationAndId(ctx, queryResult, err)
		}
	})
	if err != nil {
		return nil, err
	}
	return result.([]*utils.DbNodeWithRelationAndId), err
}

func (r *emailReadRepository) GetPrimaryEmailNodesForLinkedEntityIds(ctx context.Context, tenant string, entityType neo4jenum.EntityType, entityIds []string) ([]*utils.DbNodeWithRelationAndId, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "EmailReadRepository.GetPrimaryEmailNodesForLinkedEntityIds")
	defer spans.Finish()

	cypher := ""
	switch entityType {
	case neo4jenum.CONTACT:
		cypher = `MATCH (t:Tenant {name:$tenant})<-[:CONTACT_BELONGS_TO_TENANT]-(entity:Contact)`
	case neo4jenum.USER:
		cypher = `MATCH (t:Tenant {name:$tenant})<-[:USER_BELONGS_TO_TENANT]-(entity:User)`
	case neo4jenum.ORGANIZATION:
		cypher = `MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(entity:Organization)`
	}
	cypher = cypher + `, (entity)-[rel:HAS]->(e:Email)-[:EMAIL_ADDRESS_BELONGS_TO_TENANT]->(t)
					WHERE entity.id IN $entityIds AND rel.primary = true
					RETURN e, rel, entity.id ORDER BY e.email, e.rawEmail`
	params := map[string]any{
		"tenant":    tenant,
		"entityIds": entityIds,
	}

	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
			return nil, err
		} else {
			return utils.ExtractAllRecordsAsDbNodeWithRelationAndId(ctx, queryResult, err)
		}
	})
	if err != nil {
		return nil, err
	}
	return result.([]*utils.DbNodeWithRelationAndId), err
}

func (r *emailReadRepository) GetEmailsForValidation(ctx context.Context, delayFromLastUpdateInSeconds, delayFromLastValidationAttemptInMinutes, limit int) ([]TenantAndEmailId, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "EmailReadRepository.GetEmailsForValidation")
	defer spans.Finish()

	spans.LogKV("delayFromLastUpdateInSeconds", delayFromLastUpdateInSeconds)
	spans.LogKV("delayFromLastValidationAttemptInMinutes", delayFromLastValidationAttemptInMinutes)
	spans.LogKV("limit", limit)

	cypher := `MATCH (t:Tenant {active:true})<-[:EMAIL_ADDRESS_BELONGS_TO_TENANT]-(e:Email)
				WHERE
					e.rawEmail <> '' AND					
					(e.updatedAt < datetime() - duration({seconds: $delayFromLastUpdateInSeconds})) AND					
					(
						(
							e.techValidatedAt IS NULL AND 
							(e.techValidationRequestedAt IS NULL OR e.techValidationRequestedAt < datetime() - duration({minutes: $delayFromLastValidationAttemptInMinutes}))
						)
						OR
						(
							e.retryValidation = true AND
							(e.techValidationRequestedAt IS NULL OR e.techValidationRequestedAt < datetime() - duration({days: $delayForRetryValidationInDays}))
						)
					)
				WITH t.name as tenant, e.id as emailId
				ORDER BY 
    			CASE 
        			WHEN e.techValidationRequestedAt IS NULL THEN 0 
        			ELSE 1 
    			END ASC,
				e.techValidationRequestedAt ASC
				LIMIT $limit
				RETURN DISTINCT tenant, emailId`
	params := map[string]any{
		"delayFromLastUpdateInSeconds":            delayFromLastUpdateInSeconds,
		"delayFromLastValidationAttemptInMinutes": delayFromLastValidationAttemptInMinutes,
		"delayForRetryValidationInDays":           7,
		"limit":                                   limit,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	records, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
		if err != nil {
			return nil, err
		}
		return queryResult.Collect(ctx)
	})
	if err != nil {
		return nil, err
	}
	output := make([]TenantAndEmailId, 0)
	for _, v := range records.([]*neo4j.Record) {
		output = append(output,
			TenantAndEmailId{
				Tenant:  v.Values[0].(string),
				EmailId: v.Values[1].(string),
			})
	}
	spans.LogKV("result.count", len(output))
	return output, nil
}

func (r *emailReadRepository) IsLinkedToEntityByEmailAddress(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, email, entityId string, entityType neo4jenum.EntityType) (bool, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "EmailReadRepository.IsLinkedToEntityByEmailAddress")
	defer spans.Finish()

	spans.LogKV("email", email)
	spans.LogKV("entityId", entityId)
	spans.LogKV("entityType", entityType.String())

	cypher := ""
	switch entityType {
	case neo4jenum.CONTACT:
		cypher = `MATCH (t:Tenant {name:$tenant})<-[:CONTACT_BELONGS_TO_TENANT]-(c:Contact {id:$entityId})-[:HAS]->(e:Email)-[:EMAIL_ADDRESS_BELONGS_TO_TENANT]->(t)
					WHERE e.rawEmail = $email OR e.email = $email
					RETURN e`
	case neo4jenum.USER:
		cypher = `MATCH (t:Tenant {name:$tenant})<-[:USER_BELONGS_TO_TENANT]-(u:User {id:$entityId})-[:HAS]->(e:Email)-[:EMAIL_ADDRESS_BELONGS_TO_TENANT]->(t)
					WHERE e.rawEmail = $email OR e.email = $email
					RETURN e`
	case neo4jenum.ORGANIZATION:
		cypher = `MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(o:Organization {id:$entityId})-[:HAS]->(e:Email)-[:EMAIL_ADDRESS_BELONGS_TO_TENANT]->(t)
					WHERE e.rawEmail = $email OR e.email = $email
					RETURN e`
	}
	params := map[string]any{
		"tenant":   tenant,
		"email":    email,
		"entityId": entityId,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	result, err := utils.ExecuteReadInTransaction(ctx, r.driver, r.database, tx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
			return nil, err
		} else {
			return queryResult.Collect(ctx)
		}
	})
	if err != nil {
		return false, err
	}
	return len(result.([]*db.Record)) > 0, nil
}

func (r *emailReadRepository) GetOrphanEmailNodes(ctx context.Context, limit, hoursFromLastUpdate int) ([]TenantAndEmailId, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "EmailReadRepository.GetOrphanEmailNodes")
	defer spans.Finish()

	spans.LogKV("limit", limit)
	spans.LogKV("hoursFromLastUpdate", hoursFromLastUpdate)

	cypher := `MATCH (t:Tenant)<-[:EMAIL_ADDRESS_BELONGS_TO_TENANT]-(e:Email)
				WHERE e.updatedAt < datetime() - duration({hours: $hoursFromLastUpdate})
				OPTIONAL MATCH (e)--(other) 
				WHERE NOT (other:Tenant OR other:Domain)
				WITH t, e, other
				WHERE other IS NULL     
				RETURN DISTINCT t.name AS tenant, e.id AS emailId
				LIMIT $limit`
	params := map[string]any{
		"hoursFromLastUpdate": hoursFromLastUpdate,
		"limit":               limit,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	records, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
		if err != nil {
			return nil, err
		}
		return queryResult.Collect(ctx)
	})
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}
	output := make([]TenantAndEmailId, 0)
	for _, v := range records.([]*neo4j.Record) {
		output = append(output,
			TenantAndEmailId{
				Tenant:  v.Values[0].(string),
				EmailId: v.Values[1].(string),
			})
	}
	spans.LogKV("result.count", len(output))
	return output, nil
}

func (r *emailReadRepository) IsOrphanEmail(ctx context.Context, tenant, emailId string) (bool, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "EmailReadRepository.IsOrphanEmail")
	defer spans.Finish()

	spans.TagEntity(emailId)

	cypher := `MATCH (:Tenant)<-[:EMAIL_ADDRESS_BELONGS_TO_TENANT]-(e:Email {id:$emailId})
				OPTIONAL MATCH (e)--(other) 
				WHERE NOT (other:Tenant OR other:Domain)
				WITH e, other
				WHERE other IS NOT NULL
				RETURN count(other) = 0`
	params := map[string]any{
		"emailId": emailId,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
		return utils.ExtractSingleRecordFirstValueAsType[bool](ctx, queryResult, err)
	})
	if err != nil {
		spans.TraceError(err)
		return false, err
	}
	return result.(bool), err
}

func (r *emailReadRepository) GetTestEmailForFlows(ctx context.Context, tenant string) (string, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "EmailReadRepository.GetTestEmailForFlows")
	defer spans.Finish()

	cypher := fmt.Sprintf(`match (t:Tenant{name:$tenant})--(e:Email_%s) where e.rawEmail =~ ".*testcustomeros.com" return e.rawEmail`, tenant)
	params := map[string]any{
		"tenant": tenant,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
			return nil, err
		} else {
			return queryResult.Collect(ctx)
		}
	})
	if err != nil {
		spans.TraceError(err)
		return "", err
	}
	if len(result.([]*db.Record)) == 0 {
		spans.LogKV("result.found", false)
		return "", nil
	}
	spans.LogKV("result", result.([]*db.Record)[0].Values[0].(string))
	return result.([]*db.Record)[0].Values[0].(string), err
}
