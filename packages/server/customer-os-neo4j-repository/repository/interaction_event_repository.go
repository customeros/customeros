package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

// TODO delete me
type InteractionEventRepository interface {
	LinkInteractionEventToSession(ctx context.Context, tx neo4j.ManagedTransaction, tenant, interactionEventId, interactionSessionId string) error
	InteractionEventSentByEmail(ctx context.Context, tx neo4j.ManagedTransaction, tenant, interactionEventId, emailId string) error
	InteractionEventSentToEmails(ctx context.Context, tx neo4j.ManagedTransaction, tenant, interactionEventId, sentType string, emailsId []string) error
}

type interactionEventRepository struct {
	driver *neo4j.DriverWithContext
}

func NewInteractionEventRepository(driver *neo4j.DriverWithContext) InteractionEventRepository {
	return &interactionEventRepository{
		driver: driver,
	}
}

func (r *interactionEventRepository) LinkInteractionEventToSession(ctx context.Context, tx neo4j.ManagedTransaction, tenant, interactionEventId, interactionSessionId string) error {
	query := "MATCH (is:InteractionSession_%s {id:$interactionSessionId}) " +
		" MATCH (ie:InteractionEvent {id:$interactionEventId})" +
		" MERGE (ie)-[:PART_OF]->(is) "
	_, err := tx.Run(ctx, fmt.Sprintf(query, tenant),
		map[string]interface{}{
			"tenant":               tenant,
			"interactionSessionId": interactionSessionId,
			"interactionEventId":   interactionEventId,
		})
	return err
}

func (r *interactionEventRepository) InteractionEventSentByEmail(ctx context.Context, tx neo4j.ManagedTransaction, tenant, interactionEventId, emailId string) error {
	query := "MATCH (is:InteractionEvent_%s {id:$interactionEventId}) " +
		" MATCH (e:Email_%s {id: $emailId}) " +
		" MERGE (is)-[:SENT_BY]->(e) "
	_, err := tx.Run(ctx, fmt.Sprintf(query, tenant, tenant),
		map[string]interface{}{
			"tenant":             tenant,
			"interactionEventId": interactionEventId,
			"emailId":            emailId,
		})
	return err
}

func (r *interactionEventRepository) InteractionEventSentToEmails(ctx context.Context, tx neo4j.ManagedTransaction, tenant, interactionEventId, sentType string, emailsId []string) error {
	query := "MATCH (ie:InteractionEvent_%s {id:$interactionEventId}) " +
		" MATCH (e:Email_%s) WHERE e.id in $emailsId " +
		" MERGE (ie)-[:SENT_TO {type: $sentType}]->(e) "
	_, err := tx.Run(ctx, fmt.Sprintf(query, tenant, tenant),
		map[string]interface{}{
			"tenant":             tenant,
			"interactionEventId": interactionEventId,
			"sentType":           sentType,
			"emailsId":           emailsId,
		})
	return err
}
