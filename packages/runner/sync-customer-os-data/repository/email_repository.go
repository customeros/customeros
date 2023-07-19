package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/db"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"time"
)

type EmailRepository interface {
	GetEmailId(ctx context.Context, tenant, email string) (string, error)
	GetEmailIdOrCreateContactByEmail(ctx context.Context, tenant, email, firstName, lastName, externalSystemId string) (string, error)
	GetEmailIdOrCreateUserByEmail(ctx context.Context, tenant, email, firstName, lastName, externalSystemId string) (string, error)
	GetAllCrossTenantsWithRawEmail(ctx context.Context, size int) ([]*utils.DbNodeAndId, error)
}

type emailRepository struct {
	driver *neo4j.DriverWithContext
}

func NewEmailRepository(driver *neo4j.DriverWithContext) EmailRepository {
	return &emailRepository{
		driver: driver,
	}
}

func (r *emailRepository) GetEmailId(ctx context.Context, tenant, email string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailRepository.GetEmailId")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	records, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		queryResult, err := tx.Run(ctx,
			"MATCH (e:Email {rawEmail: $email})-[:EMAIL_ADDRESS_BELONGS_TO_TENANT]->(t:Tenant {name:$tenant}) "+
				" RETURN e.id limit 1",
			map[string]interface{}{
				"tenant": tenant,
				"email":  email,
			})
		if err != nil {
			return nil, err
		}
		return queryResult.Collect(ctx)
	})
	if err != nil {
		return "", err
	}
	if len(records.([]*db.Record)) == 0 {
		return "", nil
	}
	return records.([]*db.Record)[0].Values[0].(string), nil
}

func (r *emailRepository) GetEmailIdOrCreateContactByEmail(ctx context.Context, tenant, email, firstName, lastName, externalSystemId string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailRepository.GetEmailIdOrCreateContactByEmail")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	records, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		queryResult, err := tx.Run(ctx, fmt.Sprintf(
			" MATCH (t:Tenant {name:$tenant}) "+
				" MERGE (e:Email {rawEmail: $email})-[:EMAIL_ADDRESS_BELONGS_TO_TENANT]->(t) "+
				" ON CREATE SET "+
				"				e.id=randomUUID(), "+
				"				e.createdAt=$now, "+
				"				e.updatedAt=$now, "+
				"				e.source=$source, "+
				"				e.sourceOfTruth=$sourceOfTruth, "+
				"				e.appSource=$appSource, "+
				"				e:%s "+
				" WITH DISTINCT t, e "+
				" MERGE (e)<-[rel:HAS]-(c:Contact)-[:CONTACT_BELONGS_TO_TENANT]->(t) "+
				" ON CREATE SET rel.primary=true, "+
				"				c.id=randomUUID(), "+
				"				c.firstName=$firstName, "+
				"				c.lastName=$lastName, "+
				"				c.createdAt=$now, "+
				"				c.updatedAt=$now, "+
				"				c.source=$source, "+
				"				c.sourceOfTruth=$sourceOfTruth, "+
				"				c.appSource=$appSource, "+
				"               c:%s"+
				" RETURN e.id limit 1", "Email_"+tenant, "Contact_"+tenant),
			map[string]interface{}{
				"tenant":        tenant,
				"email":         email,
				"firstName":     firstName,
				"lastName":      lastName,
				"source":        externalSystemId,
				"sourceOfTruth": externalSystemId,
				"appSource":     externalSystemId,
				"now":           time.Now().UTC(),
			})
		if err != nil {
			return nil, err
		}
		return queryResult.Collect(ctx)
	})
	if err != nil {
		return "", err
	}
	if len(records.([]*db.Record)) == 0 {
		return "", errors.New("no contact created")
	}
	return records.([]*db.Record)[0].Values[0].(string), nil
}

func (r *emailRepository) GetEmailIdOrCreateUserByEmail(ctx context.Context, tenant, email, firstName, lastName, externalSystemId string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailRepository.GetEmailIdOrCreateUserByEmail")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	records, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		queryResult, err := tx.Run(ctx, fmt.Sprintf(
			" MATCH (t:Tenant {name:$tenant}) "+
				" MERGE (e:Email {rawEmail: $email})-[:EMAIL_ADDRESS_BELONGS_TO_TENANT]->(t) "+
				" ON CREATE SET "+
				"				e.id=randomUUID(), "+
				"				e.createdAt=$now, "+
				"				e.updatedAt=$now, "+
				"				e.source=$source, "+
				"				e.sourceOfTruth=$sourceOfTruth, "+
				"				e.appSource=$appSource, "+
				"				e:%s "+
				" WITH DISTINCT t, e "+
				" MERGE (e)<-[rel:HAS]-(u:User)-[:USER_BELONGS_TO_TENANT]->(t) "+
				" ON CREATE SET rel.primary=true, "+
				"				u.id=randomUUID(), "+
				"				u.firstName=$firstName, "+
				"				u.lastName=$lastName, "+
				"				u.createdAt=$now, "+
				"				u.updatedAt=$now, "+
				"				u.source=$source, "+
				"				u.sourceOfTruth=$sourceOfTruth, "+
				"				u.appSource=$appSource, "+
				"               u:%s"+
				" RETURN e.id limit 1", "Email_"+tenant, "User_"+tenant),
			map[string]interface{}{
				"tenant":        tenant,
				"email":         email,
				"firstName":     firstName,
				"lastName":      lastName,
				"source":        externalSystemId,
				"sourceOfTruth": externalSystemId,
				"appSource":     externalSystemId,
				"now":           time.Now().UTC(),
			})
		if err != nil {
			return nil, err
		}
		return queryResult.Collect(ctx)
	})
	if err != nil {
		return "", err
	}
	if len(records.([]*db.Record)) == 0 {
		return "", errors.New("no contact created")
	}
	return records.([]*db.Record)[0].Values[0].(string), nil
}

func (r *emailRepository) GetAllCrossTenantsWithRawEmail(ctx context.Context, size int) ([]*utils.DbNodeAndId, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailRepository.GetAllCrossTenantsWithRawEmail")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, `
			MATCH (e:Email)--(t:Tenant)
 			WHERE (e.syncedWithEventStore is null or e.syncedWithEventStore=false)
			and (e.rawEmail is not null and e.rawEmail <> '')
			RETURN e, t.name limit $size`,
			map[string]any{
				"size": size,
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
