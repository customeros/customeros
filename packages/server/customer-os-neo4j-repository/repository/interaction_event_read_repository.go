package neo4j_repository

import (
	"context"
	"fmt"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
)

type InteractionEventReadRepository interface {
	GetAllForInteractionSessions(ctx context.Context, tenant string, ids []string, returnContent bool) ([]*utils.DbPropsAndId, error)
	GetAllForMeetings(ctx context.Context, tenant string, ids []string, returnContent bool) ([]*utils.DbPropsAndId, error)
	GetAllForIssues(ctx context.Context, tenant string, issueIds []string, returnContent bool) ([]*utils.DbPropsAndId, error)
	GetSentByFor(ctx context.Context, tenant string, ids []string) ([]*utils.DbNodeWithRelationAndId, error)
	GetSentToFor(ctx context.Context, tenant string, ids []string) ([]*utils.DbNodeWithRelationAndId, error)
	GetReplyToFor(ctx context.Context, tenant string, ids []string, returnContent bool) ([]*utils.DbPropsAndId, error)
	GetInteractionEventByCustomerOSIdentifier(ctx context.Context, customerOSInternalIdentifier string) (*dbtype.Node, error)
	InteractionEventSentByUser(ctx context.Context, tenant, interactionEventId string) (bool, error)
	GetInteractionEventIdByExternalId(ctx context.Context, tenant, externalSystemId, externalId string) (string, error)
}

type interactionEventReadRepository struct {
	driver   *neo4j.DriverWithContext
	database string
}

func NewInteractionEventReadRepository(driver *neo4j.DriverWithContext, database string) InteractionEventReadRepository {
	return &interactionEventReadRepository{
		driver:   driver,
		database: database,
	}
}

func (r *interactionEventReadRepository) prepareReadSession(ctx context.Context) neo4j.SessionWithContext {
	return utils.NewNeo4jReadSession(ctx, *r.driver, utils.WithDatabaseName(r.database))
}

func (r *interactionEventReadRepository) GetAllForInteractionSessions(ctx context.Context, tenant string, ids []string, returnContent bool) ([]*utils.DbPropsAndId, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "InteractionEventReadRepository.GetAllForInteractionSessions")
	defer spans.Finish()

	spans.LogKV("returnContent", returnContent)

	cypherReturnFragment := "e {.*}"
	if !returnContent {
		cypherReturnFragment = "e {.*, content: ''}"
	}

	cypher := fmt.Sprintf(`MATCH (s:InteractionSession)<-[:PART_OF]-(e:InteractionEvent_%s) 
		 WHERE s.id IN $ids AND s:InteractionSession_%s
		 RETURN %s, s.id ORDER BY e.createdAt ASC`, tenant, tenant, cypherReturnFragment)
	params := map[string]any{
		"ids": ids,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
			return nil, err
		} else {
			return utils.ExtractAllRecordsAsDbPropsAndId(ctx, queryResult, err)
		}
	})
	if err != nil {
		return nil, err
	}
	return result.([]*utils.DbPropsAndId), err
}

func (r *interactionEventReadRepository) GetAllForMeetings(ctx context.Context, tenant string, ids []string, returnContent bool) ([]*utils.DbPropsAndId, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "InteractionEventReadRepository.GetAllForMeetings")
	defer spans.Finish()

	spans.LogKV("returnContent", returnContent)

	cypherReturnFragment := "e {.*}"
	if !returnContent {
		cypherReturnFragment = "e {.*, content: ''}"
	}

	cypher := fmt.Sprintf(`MATCH (m:Meeting)<-[:PART_OF]-(e:InteractionEvent) 
		 WHERE m.id IN $ids AND m:Meeting_%s
		 RETURN %s, m.id ORDER BY e.createdAt ASC`, tenant, cypherReturnFragment)
	params := map[string]any{
		"ids": ids,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
			return nil, err
		} else {
			return utils.ExtractAllRecordsAsDbPropsAndId(ctx, queryResult, err)
		}
	})
	if err != nil {
		return nil, err
	}
	return result.([]*utils.DbPropsAndId), err
}

func (r *interactionEventReadRepository) GetAllForIssues(ctx context.Context, tenant string, issueIds []string, returnContent bool) ([]*utils.DbPropsAndId, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "InteractionEventReadRepository.GetAllForIssues")
	defer spans.Finish()

	spans.LogKV("returnContent", returnContent)

	cypherReturnFragment := "e {.*}"
	if !returnContent {
		cypherReturnFragment = "e {.*, content: ''}"
	}

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})<-[:ISSUE_BELONGS_TO_TENANT]-(i:Issue)<-[:PART_OF]-(e:InteractionEvent) 
				WHERE i.id IN $issueIds
				RETURN %s, i.id ORDER BY e.createdAt ASC`, cypherReturnFragment)
	params := map[string]any{
		"tenant":   tenant,
		"issueIds": issueIds,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	session := utils.NewNeo4jReadSession(ctx, *r.driver, utils.WithDatabaseName(r.database))
	defer session.Close(ctx)
	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
			return nil, err
		} else {
			return utils.ExtractAllRecordsAsDbPropsAndId(ctx, queryResult, err)
		}
	})
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}
	return result.([]*utils.DbPropsAndId), nil
}

func (r *interactionEventReadRepository) GetSentByFor(ctx context.Context, tenant string, ids []string) ([]*utils.DbNodeWithRelationAndId, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "InteractionEventReadRepository.GetSentByFor")
	defer spans.Finish()

	cypher := fmt.Sprintf(`MATCH (ie:InteractionEvent)-[rel:SENT_BY]->(p:Email|PhoneNumber|User|Contact|Organization|JobRole) 
		WHERE ie.id IN $ids AND ie:InteractionEvent_%s
		RETURN p, rel, ie.id`, tenant)
	params := map[string]any{
		"ids": ids,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
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

func (r *interactionEventReadRepository) GetSentToFor(ctx context.Context, tenant string, ids []string) ([]*utils.DbNodeWithRelationAndId, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "InteractionEventReadRepository.GetSentToFor")
	defer spans.Finish()

	cypher := fmt.Sprintf(`MATCH (ie:InteractionEvent)-[rel:SENT_TO]->(p:Email|PhoneNumber|User|Contact|Organization|JobRole) 
		 WHERE ie.id IN $ids AND ie:InteractionEvent_%s RETURN p, rel, ie.id`, tenant)
	params := map[string]any{
		"ids": ids,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
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

func (r *interactionEventReadRepository) GetReplyToFor(ctx context.Context, tenant string, ids []string, returnContent bool) ([]*utils.DbPropsAndId, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "InteractionEventReadRepository.GetReplyToFor")
	defer spans.Finish()

	spans.LogKV("returnContent", returnContent)

	cypherReturnFragment := "rie {.*}"
	if !returnContent {
		cypherReturnFragment = "rie {.*, content: ''}"
	}

	cypher := fmt.Sprintf(`MATCH (ie:InteractionEvent)-[rel:REPLIES_TO]->(rie:InteractionEvent_%s) 
		 	WHERE ie.id IN $ids AND ie:InteractionEvent_%s 
			RETURN %s, ie.id`, tenant, tenant, cypherReturnFragment)
	params := map[string]any{
		"ids": ids,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
			return nil, err
		} else {
			return utils.ExtractAllRecordsAsDbPropsAndId(ctx, queryResult, err)
		}
	})
	if err != nil {
		return nil, err
	}
	return result.([]*utils.DbPropsAndId), err
}

func (r *interactionEventReadRepository) GetInteractionEventByCustomerOSIdentifier(ctx context.Context, customerOSInternalIdentifier string) (*dbtype.Node, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "InteractionEventReadRepository.GetInteractionEventByCustomerOSIdentifier")
	defer spans.Finish()
	spans.LogKV("customerOSInternalIdentifier", customerOSInternalIdentifier)

	cypher := `MATCH (i:InteractionEvent {customerOSInternalIdentifier:$customerOSInternalIdentifier}) WHERE i:InteractionEvent RETURN i`
	params := map[string]any{
		"customerOSInternalIdentifier": customerOSInternalIdentifier,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
			return nil, err
		} else {
			return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
		}
	})
	if err != nil && err.Error() == "Result contains no more records" {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return result.(*dbtype.Node), nil
}

func (r *interactionEventReadRepository) InteractionEventSentByUser(ctx context.Context, tenant, interactionEventId string) (bool, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "InteractionEventReadRepository.InteractionEventSentByUser")
	defer spans.Finish()

	spans.LogKV("interactionEventId", interactionEventId)

	cypher := fmt.Sprintf(`MATCH (i:InteractionEvent {id:$id}) WHERE i:InteractionEvent_%s AND (i)-[:SENT_BY]->(:User) OR (i)-[:SENT_BY]->(:Email|PhoneNumber)--(:User) return count(i) > 0`, tenant)
	params := map[string]any{
		"id": interactionEventId,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
			return nil, err
		} else {
			return utils.ExtractSingleRecordFirstValueAsType[bool](ctx, queryResult, err)
		}
	})
	if err != nil {
		spans.TraceError(err)
		return false, err
	}
	spans.LogKV("result", result.(bool))
	return result.(bool), nil
}

func (r *interactionEventReadRepository) GetInteractionEventIdByExternalId(ctx context.Context, tenant, externalSystemId, externalId string) (string, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "InteractionEventReadRepository.GetInteractionEventIdByExternalId")
	defer spans.Finish()

	cypher := fmt.Sprintf(`MATCH (ie:InteractionEvent_%s)-[IS_LINKED_WITH {externalId:$externalId}]-(e:ExternalSystem{id:$externalSystemId}) RETURN ie.id`, tenant)
	params := map[string]any{
		"externalId":       externalId,
		"externalSystemId": externalSystemId,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
			return nil, err
		} else {
			return utils.ExtractSingleRecordFirstValueAsString(ctx, queryResult, err)
		}
	})
	if err != nil && err.Error() == "Result contains no more records" {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return result.(string), nil
}
