package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/db"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/entity"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
	"time"
)

type MeetingRepository interface {
	GetMatchedMeetingId(ctx context.Context, tenant string, meeting entity.MeetingData) (string, error)
	MergeMeeting(ctx context.Context, tenant string, syncDate time.Time, meeting entity.MeetingData) error
	MergeMeetingLocation(ctx context.Context, tenant string, meeting entity.MeetingData) error
	MeetingLinkWithCreatorUserByExternalId(ctx context.Context, tenant, meetingId, userExternalId, externalSystem string) error
	MeetingLinkWithAttendedByUserByExternalId(ctx context.Context, tenant, meetingId, userExternalId, externalSystem string) error
	MeetingLinkWithAttendedByContactByExternalId(ctx context.Context, tenant, meetingId, contactExternalId, externalSystem string) error
}

type meetingRepository struct {
	driver *neo4j.DriverWithContext
}

func NewMeetingRepository(driver *neo4j.DriverWithContext) MeetingRepository {
	return &meetingRepository{
		driver: driver,
	}
}

func (r *meetingRepository) GetMatchedMeetingId(ctx context.Context, tenant string, meeting entity.MeetingData) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MeetingRepository.GetMatchedMeetingId")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := `MATCH (t:Tenant {name:$tenant})<-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]-(e:ExternalSystem {id:$externalSystem})
				OPTIONAL MATCH (e)<-[:IS_LINKED_WITH {externalId:$meetingExternalId}]-(m1:Meeting)
				OPTIONAL MATCH (m2:Meeting_%s)
					WHERE m2.meetingExternalUrl=$meetingExternalUrl
				WITH coalesce(m1, m2) as meeting
				WHERE meeting is not null
				return meeting.id limit 1`

	dbRecords, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, fmt.Sprintf(query, tenant),
			map[string]interface{}{
				"tenant":             tenant,
				"externalSystem":     meeting.ExternalSystem,
				"meetingExternalId":  meeting.ExternalId,
				"meetingExternalUrl": meeting.MeetingUrl,
			})
		if err != nil {
			return nil, err
		}
		return queryResult.Collect(ctx)
	})
	if err != nil {
		return "", err
	}
	meetingIDs := dbRecords.([]*db.Record)
	if len(meetingIDs) > 0 {
		return meetingIDs[0].Values[0].(string), nil
	}
	return "", nil
}

func (r *meetingRepository) MergeMeeting(ctx context.Context, tenant string, syncDate time.Time, meeting entity.MeetingData) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MeetingRepository.MergeMeeting")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	// Create new Meeting if it does not exist
	// If Note exists, and sourceOfTruth is acceptable then update Meeting.
	query := "MATCH (t:Tenant {name:$tenant})<-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]-(ext:ExternalSystem {id:$externalSystem}) " +
		" MERGE (m:Meeting {id:$meetingId}) " +
		" ON CREATE SET " +
		"				m.createdAt=$createdAt, " +
		"				m.updatedAt=$updatedAt, " +
		"				m.startedAt=$startedAt, " +
		"				m.endedAt=$endedAt, " +
		"              	m.source=$source, " +
		"				m.sourceOfTruth=$sourceOfTruth, " +
		"				m.appSource=$appSource, " +
		"              	m.name=$name, " +
		"              	m.agenda=$agenda, " +
		"              	m.agendaContentType=$agendaContentType, " +
		"              	m.conferenceUrl=$conferenceUrl, " +
		"              	m.meetingExternalUrl=$meetingExternalUrl, " +
		"				m:Meeting_%s, " +
		" 				m:TimelineEvent, " +
		"				m:TimelineEvent_%s " +
		" ON MATCH SET 	m.agenda = CASE WHEN m.sourceOfTruth=$sourceOfTruth OR m.agenda is null or m.agenda = '' THEN $agenda ELSE m.agenda END, " +
		"             	m.agendaContentType = CASE WHEN m.sourceOfTruth=$sourceOfTruth OR m.agendaContentType is null or m.agendaContentType = '' THEN $agendaContentType ELSE m.agendaContentType END, " +
		"             	m.name = CASE WHEN m.sourceOfTruth=$sourceOfTruth OR m.name is null or m.name = '' THEN $name ELSE m.name END, " +
		"             	m.conferenceUrl = CASE WHEN m.sourceOfTruth=$sourceOfTruth OR m.conferenceUrl is null or m.conferenceUrl = '' THEN $conferenceUrl ELSE m.conferenceUrl END, " +
		"             	m.meetingExternalUrl = CASE WHEN m.sourceOfTruth=$sourceOfTruth OR m.meetingExternalUrl is null or m.meetingExternalUrl = '' THEN $meetingExternalUrl ELSE m.meetingExternalUrl END, " +
		"				m.updatedAt = $now " +
		" WITH m, ext " +
		" MERGE (m)-[r:IS_LINKED_WITH {externalId:$externalId}]->(ext) " +
		" ON CREATE SET r.syncDate=$syncDate " +
		" ON MATCH SET r.syncDate=$syncDate " +
		" RETURN m.id"

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, fmt.Sprintf(query, tenant, tenant),
			map[string]interface{}{
				"tenant":             tenant,
				"meetingId":          meeting.Id,
				"source":             meeting.ExternalSystem,
				"sourceOfTruth":      meeting.ExternalSystem,
				"appSource":          meeting.ExternalSystem,
				"externalSystem":     meeting.ExternalSystem,
				"externalId":         meeting.ExternalId,
				"syncDate":           syncDate,
				"name":               meeting.Name,
				"meetingExternalUrl": meeting.MeetingUrl,
				"conferenceUrl":      meeting.ConferenceUrl,
				"agenda":             meeting.Agenda,
				"agendaContentType":  meeting.ContentType,
				"createdAt":          utils.TimePtrFirstNonNilNillableAsAny(meeting.CreatedAt),
				"updatedAt":          utils.TimePtrFirstNonNilNillableAsAny(meeting.UpdatedAt),
				"startedAt":          utils.TimePtrFirstNonNilNillableAsAny(meeting.StartedAt),
				"endedAt":            utils.TimePtrFirstNonNilNillableAsAny(meeting.EndedAt),
				"now":                time.Now().UTC(),
			})
		if err != nil {
			return nil, err
		}
		_, err = queryResult.Single(ctx)
		if err != nil {
			return nil, err
		}
		return nil, nil
	})
	return err
}

func (r *meetingRepository) MergeMeetingLocation(ctx context.Context, tenant string, meeting entity.MeetingData) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MeetingRepository.MergeMeetingLocation")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	// Create new Location if it does not exist with given source and raw address
	query := "MATCH (m:Meeting_%s {id:$meetingId}) " +
		" MERGE (m)-[:ASSOCIATED_WITH]->(loc:Location {source:$source}) " +
		" ON CREATE SET " +
		"	loc.id=randomUUID(), " +
		"	loc.rawAddress=$rawAddress, " +
		"	loc.appSource=$appSource, " +
		"	loc.sourceOfTruth=$sourceOfTruth, " +
		"	loc.createdAt=$createdAt, " +
		"	loc.updatedAt=$updatedAt, " +
		"	loc:Location_%s " +
		" ON MATCH SET 	" +
		"	loc.rawAddress = CASE WHEN loc.sourceOfTruth=$sourceOfTruth OR loc.rawAddress is null or loc.rawAddress = '' THEN $rawAddress ELSE loc.rawAddress END, " +
		"   loc.updatedAt = CASE WHEN loc.sourceOfTruth=$sourceOfTruth THEN $now ELSE loc.updatedAt END " +
		" WITH loc " +
		" MATCH (t:Tenant {name:$tenant}) " +
		" MERGE (loc)-[:LOCATION_BELONGS_TO_TENANT]->(t) "

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, fmt.Sprintf(query, tenant, tenant),
			map[string]interface{}{
				"tenant":        tenant,
				"meetingId":     meeting.Id,
				"rawAddress":    meeting.Location,
				"source":        meeting.ExternalSystem,
				"sourceOfTruth": meeting.ExternalSystem,
				"appSource":     meeting.ExternalSystem,
				"createdAt":     utils.TimePtrFirstNonNilNillableAsAny(meeting.CreatedAt),
				"updatedAt":     utils.TimePtrFirstNonNilNillableAsAny(meeting.UpdatedAt),
				"now":           time.Now().UTC(),
			})
		return nil, err
	})
	return err
}

func (r *meetingRepository) MeetingLinkWithCreatorUserByExternalId(ctx context.Context, tenant, meetingId, userExternalId, externalSystem string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MeetingRepository.MeetingLinkWithCreatorUserByExternalId")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, `
				MATCH (t:Tenant {name:$tenant})<-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]-(ext:ExternalSystem {id:$externalSystem})<-[:IS_LINKED_WITH {externalId:$userExternalId}]-(u:User)
				MATCH (m:Meeting {id:$meetingId})-[:IS_LINKED_WITH]->(ext)
				MERGE (m)-[:CREATED_BY]->(u)
				`,
			map[string]interface{}{
				"tenant":         tenant,
				"externalSystem": externalSystem,
				"meetingId":      meetingId,
				"userExternalId": userExternalId,
			})
		return nil, err
	})
	return err
}

func (r *meetingRepository) MeetingLinkWithAttendedByUserByExternalId(ctx context.Context, tenant, meetingId, userExternalId, externalSystem string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MeetingRepository.MeetingLinkWithAttendedByUserByExternalId")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, `
				MATCH (t:Tenant {name:$tenant})<-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]-(ext:ExternalSystem {id:$externalSystem})<-[:IS_LINKED_WITH {externalId:$userExternalId}]-(u:User)
				MATCH (m:Meeting {id:$meetingId})-[:IS_LINKED_WITH]->(ext)
				MERGE (m)-[:ATTENDED_BY]->(u)
				`,
			map[string]interface{}{
				"tenant":         tenant,
				"externalSystem": externalSystem,
				"meetingId":      meetingId,
				"userExternalId": userExternalId,
			})
		return nil, err
	})
	return err
}

func (r *meetingRepository) MeetingLinkWithAttendedByContactByExternalId(ctx context.Context, tenant, meetingId, contactExternalId, externalSystem string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MeetingRepository.MeetingLinkWithAttendedByContactByExternalId")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, `
				MATCH (t:Tenant {name:$tenant})<-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]-(ext:ExternalSystem {id:$externalSystem})<-[:IS_LINKED_WITH {externalId:$contactExternalId}]-(c:Contact)
				MATCH (m:Meeting {id:$meetingId})-[:IS_LINKED_WITH]->(ext)
				MERGE (m)-[:ATTENDED_BY]->(c)
				`,
			map[string]interface{}{
				"tenant":            tenant,
				"externalSystem":    externalSystem,
				"meetingId":         meetingId,
				"contactExternalId": contactExternalId,
			})
		return nil, err
	})
	return err
}
