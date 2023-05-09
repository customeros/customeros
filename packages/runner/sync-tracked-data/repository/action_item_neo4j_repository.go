package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-tracked-data/entity"
)

type ActionRepository interface {
	CreatePageViewAction(ctx context.Context, contactId string, pv entity.PageView) error
}

type actionRepository struct {
	driver *neo4j.DriverWithContext
}

func NewActionRepository(driver *neo4j.DriverWithContext) ActionRepository {
	return &actionRepository{
		driver: driver,
	}
}

func (r *actionRepository) CreatePageViewAction(ctx context.Context, contactId string, pv entity.PageView) error {
	session := (*r.driver).NewSession(ctx,
		neo4j.SessionConfig{
			AccessMode: neo4j.AccessModeWrite,
			BoltLogger: neo4j.ConsoleBoltLogger()})
	defer session.Close(ctx)

	params := map[string]interface{}{
		"tenant":         pv.Tenant,
		"contactId":      contactId,
		"pvId":           pv.ID,
		"start":          pv.Start,
		"end":            pv.End,
		"application":    pv.Application,
		"sessionId":      pv.SessionID,
		"trackerName":    pv.TrackerName,
		"pageUrl":        pv.Url,
		"pageTitle":      pv.Title,
		"orderInSession": pv.OrderInSession,
		"engagedTime":    pv.EngagedTime,
		"source":         "openline",
		"sourceOfTruth":  "openline",
		"appSource":      pv.Application,
	}

	query := "MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(t:Tenant {name:$tenant}) " +
		" MERGE (c)-[:HAS_ACTION]->(a:TimelineEvent:PageView {id:$pvId, trackerName:$trackerName})" +
		" ON CREATE SET " +
		" 	a:%s, a:%s, " +
		" 	a.startedAt=$start, " +
		" 	a.endedAt=$end, " +
		" 	a.application=$application, " +
		" 	a.pageUrl=$pageUrl, " +
		" 	a.pageTitle=$pageTitle, " +
		" 	a.sessionId=$sessionId, " +
		" 	a.orderInSession=$orderInSession, " +
		" 	a.engagedTime=$engagedTime, " +
		" 	a.source=$source, " +
		" 	a.sourceOfTruth=$sourceOfTruth, " +
		" 	a.appSource=$appSource " +
		" ON MATCH SET " +
		" 	a.endedAt=$end, " +
		" 	a.engagedTime=$engagedTime, " +
		" 	a.orderInSession=$orderInSession "

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		_, err := tx.Run(ctx, fmt.Sprintf(query, "PageView_"+pv.Tenant, "TimelineEvent_"+pv.Tenant), params)
		return nil, err
	})

	return err
}
