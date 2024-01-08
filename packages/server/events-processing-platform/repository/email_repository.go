package repository

import (
	"context"
	"time"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type EmailRepository interface {
	LinkWithContact(ctx context.Context, tenant, contactId, emailId, label string, primary bool, updatedAt time.Time) error
	LinkWithOrganization(ctx context.Context, tenant, organizationId, emailId, label string, primary bool, updatedAt time.Time) error
	LinkWithUser(ctx context.Context, tenant, userId, emailId, label string, primary bool, updatedAt time.Time) error
}

type emailRepository struct {
	driver *neo4j.DriverWithContext
}

func NewEmailRepository(driver *neo4j.DriverWithContext) EmailRepository {
	return &emailRepository{
		driver: driver,
	}
}

func (r *emailRepository) LinkWithContact(ctx context.Context, tenant, contactId, emailId, label string, primary bool, updatedAt time.Time) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailRepository.LinkWithContact")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(ctx, span, tenant)
	span.LogFields(log.String("emailId", emailId), log.String("contactId", contactId))

	query := `MATCH (t:Tenant {name:$tenant})<-[:CONTACT_BELONGS_TO_TENANT]-(c:Contact {id:$contactId}),
				(t)<-[:EMAIL_ADDRESS_BELONGS_TO_TENANT]-(e:Email {id:$emailId})
		MERGE (c)-[rel:HAS]->(e)
		SET	rel.primary = $primary,
			rel.label = $label,	
			c.updatedAt = $updatedAt,
			rel.syncedWithEventStore = true`
	params := map[string]any{
		"tenant":    tenant,
		"contactId": contactId,
		"emailId":   emailId,
		"label":     label,
		"primary":   primary,
		"updatedAt": updatedAt,
	}
	span.LogFields(log.String("query", query), log.Object("params", params))

	return r.executeQuery(ctx, query, params)
}

func (r *emailRepository) LinkWithOrganization(ctx context.Context, tenant, organizationId, emailId, label string, primary bool, updatedAt time.Time) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailRepository.LinkWithOrganization")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(ctx, span, tenant)
	span.LogFields(log.String("emailId", emailId), log.String("organizationId", organizationId))

	query := `
		MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization {id:$organizationId}),
				(t)<-[:EMAIL_ADDRESS_BELONGS_TO_TENANT]-(e:Email {id:$emailId})
		MERGE (org)-[rel:HAS]->(e)
		SET	rel.primary = $primary,
			rel.label = $label,	
			org.updatedAt = $updatedAt,
			rel.syncedWithEventStore = true`
	params := map[string]any{
		"tenant":         tenant,
		"organizationId": organizationId,
		"emailId":        emailId,
		"label":          label,
		"primary":        primary,
		"updatedAt":      updatedAt,
	}
	span.LogFields(log.String("query", query), log.Object("params", params))

	return r.executeQuery(ctx, query, params)
}

func (r *emailRepository) LinkWithUser(ctx context.Context, tenant, userId, emailId, label string, primary bool, updatedAt time.Time) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailRepository.LinkWithUser")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(ctx, span, tenant)
	span.LogFields(log.String("emailId", emailId), log.String("userId", userId))

	query := `
		MATCH (t:Tenant {name:$tenant})<-[:USER_BELONGS_TO_TENANT]-(u:User {id:$userId}),
				(t)<-[:EMAIL_ADDRESS_BELONGS_TO_TENANT]-(e:Email {id:$emailId})
		MERGE (u)-[rel:HAS]->(e)
		SET	rel.primary = $primary,
			rel.label = $label,	
			u.updatedAt = $updatedAt,
			rel.syncedWithEventStore = true`
	params := map[string]any{
		"tenant":    tenant,
		"userId":    userId,
		"emailId":   emailId,
		"label":     label,
		"primary":   primary,
		"updatedAt": updatedAt,
	}
	span.LogFields(log.String("query", query), log.Object("params", params))
	return r.executeQuery(ctx, query, params)
}

func (r *emailRepository) executeQuery(ctx context.Context, query string, params map[string]any) error {
	return utils.ExecuteWriteQuery(ctx, *r.driver, query, params)
}
