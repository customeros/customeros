package neo4j

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/runner/customer-os-data-upkeeper/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"time"
)

type OrganizationForWebScrape struct {
	Tenant         string
	OrganizationId string
	Url            string
}

type OrganizationRepository interface {
	GetOrganizationsForWebScrape(ctx context.Context, limit int) ([]OrganizationForWebScrape, error)
}

type organizationRepository struct {
	driver *neo4j.DriverWithContext
}

func NewOrganizationRepository(driver *neo4j.DriverWithContext) OrganizationRepository {
	return &organizationRepository{
		driver: driver,
	}
}

func (r *organizationRepository) GetOrganizationsForWebScrape(ctx context.Context, limit int) ([]OrganizationForWebScrape, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationRepository.GetOrganizationsForWebScrape")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(span)
	span.LogFields(log.Int("limit", limit))
	if limit <= 0 {
		limit = 1
	}

	aDayAgo := utils.Now().Add(-24 * time.Hour)
	cypher := `MATCH (org:Organization)-[:ORGANIZATION_BELONGS_TO_TENANT]->(t:Tenant)
				WHERE (org.webScrapeLastRequestedAt IS NULL OR org.webScrapeLastRequestedAt < $aDayAgo)
					AND (org.webScrapeAttempts < 4 or org.webScrapeAttempts IS NULL)
				WITH org, t OPTIONAL MATCH (org)-[:HAS_DOMAIN]->(domain:Domain) WHERE domain.domain <> "" AND domain.domain IS NOT NULL
				WITH t, org, COLLECT(DISTINCT domain.domain) as domains
				WITH t, org, CASE WHEN org.website IS NOT NULL AND org.website <> "" THEN (domains + org.website) ELSE domains END AS urls
				WHERE size(urls) > 0 AND (org.webScrapedUrl IS NULL OR NOT org.webScrapedUrl in urls)
				RETURN t.name AS tenant_name, org.id AS org_id, HEAD(urls) AS url ORDER BY org.createdAt DESC LIMIT $limit`
	params := map[string]any{
		"limit":   limit,
		"aDayAgo": aDayAgo,
	}
	span.LogFields(log.String("query", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
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
	output := make([]OrganizationForWebScrape, 0)
	for _, v := range records.([]*neo4j.Record) {
		output = append(output,
			OrganizationForWebScrape{
				Tenant:         v.Values[0].(string),
				OrganizationId: v.Values[1].(string),
				Url:            v.Values[2].(string),
			})
	}
	span.LogFields(log.Int("output - length", len(output)))
	return output, nil
}
