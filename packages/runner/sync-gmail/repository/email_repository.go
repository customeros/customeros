package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/db"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/pkg/errors"
	"time"
)

type EmailRepository interface {
	GetEmailId(ctx context.Context, tenant, email string) (string, error)
	GetEmailIdInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant, email string) (string, error)
	FindEmailByUserId(ctx context.Context, tenant string, userId string) (*dbtype.Node, error)
	CreateEmail(ctx context.Context, tx neo4j.ManagedTransaction, tenant, email, source, appSource string) (string, error)
	CreateContactWithEmailLinkedToOrganization(ctx context.Context, tx neo4j.ManagedTransaction, tenant, organizationId, email, firstName, lastName, source, appSource string) (string, error)
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

func (r *emailRepository) GetEmailIdInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant, email string) (string, error) {
	dbResult, err := tx.Run(ctx,
		"MATCH (e:Email {rawEmail: $email})-[:EMAIL_ADDRESS_BELONGS_TO_TENANT]->(t:Tenant {name:$tenant}) "+
			" RETURN e.id limit 1",
		map[string]interface{}{
			"tenant": tenant,
			"email":  email,
		})
	if err != nil {
		return "", err
	}
	records, err := dbResult.Collect(ctx)
	if err != nil {
		return "", err
	}
	if len(records) == 0 {
		return "", nil
	}
	return records[0].Values[0].(string), nil
}

func (r *emailRepository) FindEmailByUserId(ctx context.Context, tenant string, userId string) (*dbtype.Node, error) {
	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	result, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		queryResult, err := tx.Run(ctx, `
			MATCH (:Tenant {name:$tenant})<-[:USER_BELONGS_TO_TENANT]-(u:User{id:$userId})-[:HAS]->(e:Email) 
			RETURN DISTINCT e limit 1`,
			map[string]any{
				"tenant": tenant,
				"userId": userId,
			})
		return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
	})
	if err != nil {
		return nil, err
	}
	return result.(*dbtype.Node), nil
}

func (r *emailRepository) CreateEmail(ctx context.Context, tx neo4j.ManagedTransaction, tenant, email, source, appSource string) (string, error) {
	dbResult, err := tx.Run(ctx, fmt.Sprintf(
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
			" return e.id ", "Email_"+tenant),
		map[string]interface{}{
			"tenant":        tenant,
			"email":         email,
			"source":        source,
			"sourceOfTruth": "openline",
			"appSource":     appSource,
			"now":           time.Now().UTC(),
		})
	if err != nil {
		return "", err
	}
	records, err := dbResult.Collect(ctx)
	if err != nil {
		return "", err
	}
	if len(records) == 0 {
		return "", errors.New("no email created")
	}
	return records[0].Values[0].(string), nil
}

func (r *emailRepository) CreateContactWithEmailLinkedToOrganization(ctx context.Context, tx neo4j.ManagedTransaction, tenant, organizationId, email, firstName, lastName, source, appSource string) (string, error) {
	dbResult, err := tx.Run(ctx, fmt.Sprintf(
		" MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(o:Organization{id: $organizationId}) WITH t, o "+
			" MERGE (e:Email {rawEmail: $email})-[:EMAIL_ADDRESS_BELONGS_TO_TENANT]->(t) "+
			" ON CREATE SET "+
			"				e.id=randomUUID(), "+
			"				e.createdAt=$now, "+
			"				e.updatedAt=$now, "+
			"				e.source=$source, "+
			"				e.sourceOfTruth=$sourceOfTruth, "+
			"				e.appSource=$appSource, "+
			"				e:%s "+
			" WITH t, e, o "+
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
			" WITH t, e, o, c "+
			" MERGE (c)-[:WORKS_AS]->(j:JobRole)-[:ROLE_IN]->(o) "+
			" ON CREATE SET c.updatedAt=$now, "+
			" 	j.id=randomUUID(), "+
			"   j.source=$source, "+
			"   j.sourceOfTruth=$source, "+
			"   j.appSource=$appSource, "+
			"   j.createdAt=$now, "+
			"   j.updatedAt=$now, "+
			"   j:JobRole_%s "+
			" RETURN e.id limit 1", "Email_"+tenant, "Contact_"+tenant, tenant),
		map[string]interface{}{
			"tenant":         tenant,
			"organizationId": organizationId,
			"email":          email,
			"firstName":      firstName,
			"lastName":       lastName,
			"source":         source,
			"sourceOfTruth":  "openline",
			"appSource":      appSource,
			"now":            time.Now().UTC(),
		})
	if err != nil {
		return "", err
	}
	records, err := dbResult.Collect(ctx)
	if err != nil {
		return "", err
	}
	if len(records) == 0 {
		return "", errors.New("no email created")
	}
	return records[0].Values[0].(string), nil
}
