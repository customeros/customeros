package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/db"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/constants"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/entity"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
	"time"
)

type NoteRepository interface {
	GetMatchedNoteId(ctx context.Context, tenant string, note entity.NoteData) (string, error)
	MergeNote(ctx context.Context, tenant string, syncDate time.Time, note entity.NoteData) error
	NoteLinkWithContactByExternalId(ctx context.Context, tenant, noteId, contactExternalId, externalSystem string) error
	NoteLinkWithOrganizationByExternalId(ctx context.Context, tenant, noteId, organizationExternalId, externalSystem string) error
	NoteLinkWithIssueReporterContactOrOrganization(ctx context.Context, tenant, noteId, issueId, externalSystem string) error
	NoteMentionedTag(ctx context.Context, tenant, noteId, tagName, externalSystem string) error
	NoteLinkWithCreatorUserByExternalId(ctx context.Context, tenant, noteId, userExternalId, externalSystem string) error
	NoteLinkWithCreatorUserByExternalOwnerId(ctx context.Context, tenant, noteId, userExternalOwnerId, externalSystem string) error
	NoteLinkWithCreatorByExternalId(ctx context.Context, tenant, noteId, creatorExternalId, externalSystem string) error
}

type noteRepository struct {
	driver *neo4j.DriverWithContext
}

func NewNoteRepository(driver *neo4j.DriverWithContext) NoteRepository {
	return &noteRepository{
		driver: driver,
	}
}

func (r *noteRepository) GetMatchedNoteId(ctx context.Context, tenant string, note entity.NoteData) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "NoteRepository.GetMatchedNoteId")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := `MATCH (t:Tenant {name:$tenant})<-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]-(e:ExternalSystem {id:$externalSystem})
				OPTIONAL MATCH (e)<-[:IS_LINKED_WITH {externalId:$noteExternalId}]-(n:Note)
				WITH n WHERE n is not null
				return n.id limit 1`

	dbRecords, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, query,
			map[string]interface{}{
				"tenant":         tenant,
				"externalSystem": note.ExternalSystem,
				"noteExternalId": note.ExternalId,
			})
		if err != nil {
			return nil, err
		}
		return queryResult.Collect(ctx)
	})
	if err != nil {
		return "", err
	}
	noteIDs := dbRecords.([]*db.Record)
	if len(noteIDs) > 0 {
		return noteIDs[0].Values[0].(string), nil
	}
	return "", nil
}

func (r *noteRepository) MergeNote(ctx context.Context, tenant string, syncDate time.Time, note entity.NoteData) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "NoteRepository.MergeNote")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	// Create new Note if it does not exist
	// If Note exists, and sourceOfTruth is acceptable then update Note.
	//   otherwise create/update AlternateNote for incoming source, with a new relationship 'ALTERNATE'
	query := "MATCH (t:Tenant {name:$tenant})<-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]-(ext:ExternalSystem {id:$externalSystem}) " +
		" MERGE (n:Note {id:$noteId}) " +
		" ON CREATE SET " +
		"				n.createdAt=$createdAt, " +
		"				n.updatedAt=$createdAt, " +
		"              	n.source=$source, " +
		"				n.sourceOfTruth=$sourceOfTruth, " +
		"				n.appSource=$appSource, " +
		"              	n.html=$content, " +
		"              	n.content=$content, " +
		"              	n.contentType=$contentType, " +
		"              	n.text=$text, " +
		"				n:Note_%s, " +
		" 				n:TimelineEvent, " +
		"				n:TimelineEvent_%s " +
		" ON MATCH SET 	n.html = CASE WHEN n.sourceOfTruth=$sourceOfTruth OR n.html is null or n.html = '' THEN $content ELSE n.html END, " +
		"             	n.content = CASE WHEN n.sourceOfTruth=$sourceOfTruth OR n.content is null or n.content = '' THEN $content ELSE n.content END, " +
		"             	n.contentType = CASE WHEN n.sourceOfTruth=$sourceOfTruth OR n.contentType is null or n.contentType = '' THEN $contentType ELSE n.contentType END, " +
		"             	n.text = CASE WHEN n.sourceOfTruth=$sourceOfTruth OR n.text is null or n.text = '' THEN $text ELSE n.text END, " +
		"				n.updatedAt = $now " +
		" WITH n, ext " +
		" MERGE (n)-[r:IS_LINKED_WITH {externalId:$externalId}]->(ext) " +
		" ON CREATE SET r.syncDate=$syncDate " +
		" ON MATCH SET r.syncDate=$syncDate " +
		" WITH n " +
		" FOREACH (x in CASE WHEN n.sourceOfTruth <> $sourceOfTruth THEN [n] ELSE [] END | " +
		"  MERGE (x)-[:ALTERNATE]->(alt:AlternateNote {source:$source, id:x.id}) " +
		"    SET alt.updatedAt=$now, alt.appSource=$appSource, alt.html=$content, alt.content=$content, alt.contentType=$contentType, alt.text=$text " +
		" ) " +
		" RETURN n.id"

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, fmt.Sprintf(query, tenant, tenant),
			map[string]interface{}{
				"tenant":         tenant,
				"noteId":         note.Id,
				"source":         note.ExternalSystem,
				"sourceOfTruth":  note.ExternalSystem,
				"appSource":      constants.AppSourceSyncCustomerOsData,
				"externalSystem": note.ExternalSystem,
				"externalId":     note.ExternalId,
				"syncDate":       syncDate,
				"content":        note.Content,
				"contentType":    note.ContentType,
				"text":           note.Text,
				"createdAt":      utils.TimePtrFirstNonNilNillableAsAny(note.CreatedAt),
				"now":            time.Now().UTC(),
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

func (r *noteRepository) NoteLinkWithContactByExternalId(ctx context.Context, tenant, noteId, contactExternalId, externalSystem string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "NoteRepository.NoteLinkWithContactByExternalId")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := `MATCH (t:Tenant {name:$tenant})<-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]-(e:ExternalSystem {id:$externalSystem})<-[:IS_LINKED_WITH {externalId:$contactExternalId}]-(c:Contact)
				MATCH (n:Note {id:$noteId})-[:IS_LINKED_WITH]->(e)
				MERGE (c)-[:NOTED]->(n)`

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, query,
			map[string]interface{}{
				"tenant":            tenant,
				"externalSystem":    externalSystem,
				"noteId":            noteId,
				"contactExternalId": contactExternalId,
			})
		return nil, err
	})
	return err
}

func (r *noteRepository) NoteLinkWithOrganizationByExternalId(ctx context.Context, tenant, noteId, organizationExternalId, externalSystem string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "NoteRepository.NoteLinkWithOrganizationByExternalId")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := `MATCH (t:Tenant {name:$tenant})<-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]-(e:ExternalSystem {id:$externalSystem})<-[:IS_LINKED_WITH {externalId:$organizationExternalId}]-(org:Organization)
				MATCH (n:Note {id:$noteId})-[:IS_LINKED_WITH]->(e)
				MERGE (org)-[:NOTED]->(n)`

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, query,
			map[string]interface{}{
				"tenant":                 tenant,
				"externalSystem":         externalSystem,
				"noteId":                 noteId,
				"organizationExternalId": organizationExternalId,
			})
		return nil, err
	})
	return err
}

func (r *noteRepository) NoteLinkWithIssueReporterContactOrOrganization(ctx context.Context, tenant, noteId, issueId, externalSystem string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "NoteRepository.NoteLinkWithIssueReporterContactOrOrganization")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := `MATCH (t:Tenant {name:$tenant}),
				(t)<-[:ISSUE_BELONGS_TO_TENANT]-(i:Issue {id:$issueId}),
				(t)<-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]-(e:ExternalSystem {id:$externalSystem})<-[:IS_LINKED_WITH]-(n:Note {id:$noteId}),
				(i)-[:REPORTED_BY]->(reporter)
				WHERE "Contact" in labels(reporter) OR "Organization" in labels(reporter)
				MERGE (reporter)-[:NOTED]->(n)`

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, query,
			map[string]interface{}{
				"tenant":         tenant,
				"issueId":        issueId,
				"noteId":         noteId,
				"externalSystem": externalSystem,
			})
		return nil, err
	})
	return err
}

func (r *noteRepository) NoteMentionedTag(ctx context.Context, tenant, noteId, tagName, externalSystem string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "NoteRepository.NoteMentionedTag")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := "MATCH (t:Tenant {name:$tenant}) " +
		" MERGE (tag:Tag {name:$tagName})-[:TAG_BELONGS_TO_TENANT]->(t) " +
		" ON CREATE SET tag.id=randomUUID(), " +
		"				tag.createdAt=$now, " +
		"				tag.updatedAt=$now, " +
		"				tag.source=$source," +
		"				tag.sourceOfTruth=$sourceOfTruth," +
		"				tag.appSource=$appSource," +
		"				tag:Tag_%s  " +
		" WITH DISTINCT tag " +
		" MATCH (n:Note_%s {id:$noteId}) " +
		" MERGE (n)-[r:MENTIONED]->(tag) " +
		" return r"

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, fmt.Sprintf(query, tenant, tenant),
			map[string]interface{}{
				"tenant":        tenant,
				"noteId":        noteId,
				"tagName":       tagName,
				"source":        externalSystem,
				"sourceOfTruth": externalSystem,
				"appSource":     constants.AppSourceSyncCustomerOsData,
				"now":           time.Now().UTC(),
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

func (r *noteRepository) NoteLinkWithCreatorUserByExternalId(ctx context.Context, tenant, noteId, userExternalId, externalSystem string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "NoteRepository.NoteLinkWithCreatorUserByExternalId")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, `
				MATCH (t:Tenant {name:$tenant})<-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]-(e:ExternalSystem {id:$externalSystem})<-[:IS_LINKED_WITH {externalId:$userExternalId}]-(u:User)
				MATCH (n:Note {id:$noteId})-[:IS_LINKED_WITH]->(e)
				MERGE (u)-[:CREATED]->(n)
				`,
			map[string]interface{}{
				"tenant":         tenant,
				"externalSystem": externalSystem,
				"noteId":         noteId,
				"userExternalId": userExternalId,
			})
		return nil, err
	})
	return err
}

func (r *noteRepository) NoteLinkWithCreatorUserByExternalOwnerId(ctx context.Context, tenant, noteId, userExternalOwnerId, externalSystem string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "NoteRepository.NoteLinkWithCreatorUserByExternalOwnerId")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, `
				MATCH (t:Tenant {name:$tenant})<-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]-(e:ExternalSystem {id:$externalSystem})<-[:IS_LINKED_WITH {externalOwnerId:$userExternalOwnerId}]-(u:User)
				MATCH (n:Note {id:$noteId})-[:IS_LINKED_WITH]->(e)
				MERGE (u)-[:CREATED]->(n)
				`,
			map[string]interface{}{
				"tenant":              tenant,
				"externalSystem":      externalSystem,
				"noteId":              noteId,
				"userExternalOwnerId": userExternalOwnerId,
			})
		return nil, err
	})
	return err
}

func (r *noteRepository) NoteLinkWithCreatorByExternalId(ctx context.Context, tenant, noteId, creatorExternalId, externalSystem string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "NoteRepository.NoteLinkWithCreatorByExternalId")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, `
				MATCH (t:Tenant {name:$tenant})<-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]-(e:ExternalSystem {id:$externalSystem})<-[:IS_LINKED_WITH {externalId:$creatorExternalId}]-(c)
				WHERE c:User OR c:Contact OR c:Organization
				MATCH (n:Note {id:$noteId})-[:IS_LINKED_WITH]->(e)
				MERGE (c)-[:CREATED]->(n) `,
			map[string]interface{}{
				"tenant":            tenant,
				"externalSystem":    externalSystem,
				"noteId":            noteId,
				"creatorExternalId": creatorExternalId,
			})
		return nil, err
	})
	return err
}
