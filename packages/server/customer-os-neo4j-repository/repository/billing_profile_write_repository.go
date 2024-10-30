package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/model"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"time"
)

type BillingProfileCreateFields struct {
	OrganizationId string             `json:"organizationId"`
	LegalName      string             `json:"legalName"`
	TaxId          string             `json:"taxId"`
	CreatedAt      time.Time          `json:"createdAt"`
	SourceFields   model.SourceFields `json:"sourceFields"`
}

type BillingProfileUpdateFields struct {
	OrganizationId  string `json:"organizationId"`
	LegalName       string `json:"legalName"`
	TaxId           string `json:"taxId"`
	UpdateLegalName bool   `json:"updateLegalName"`
	UpdateTaxId     bool   `json:"updateTaxId"`
}

type BillingProfileWriteRepository interface {
	Create(ctx context.Context, tenant, billingProfileId string, data BillingProfileCreateFields) error
	Update(ctx context.Context, tenant, billingProfileId string, data BillingProfileUpdateFields) error
	LinkEmailToBillingProfile(ctx context.Context, tenant, organizationId, billingProfileId, emailId string, primary bool) error
	UnlinkEmailFromBillingProfile(ctx context.Context, tenant, organizationId, billingProfileId, emailId string) error
	LinkLocationToBillingProfile(ctx context.Context, tenant, organizationId, billingProfileId, locationId string) error
	UnlinkLocationFromBillingProfile(ctx context.Context, tenant, organizationId, billingProfileId, locationId string) error
}

type billingProfileWriteRepository struct {
	driver   *neo4j.DriverWithContext
	database string
}

func NewBillingProfileWriteRepository(driver *neo4j.DriverWithContext, database string) BillingProfileWriteRepository {
	return &billingProfileWriteRepository{
		driver:   driver,
		database: database,
	}
}

func (r *billingProfileWriteRepository) Create(ctx context.Context, tenant, billingProfileId string, data BillingProfileCreateFields) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "BillingProfileWriteRepository.Create")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)

	span.SetTag(tracing.SpanTagEntityId, billingProfileId)
	tracing.LogObjectAsJson(span, "data", data)

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization {id:$orgId})
							MERGE (bp:BillingProfile {id:$billingProfileId})<-[:HAS_BILLING_PROFILE]-(org)
							ON CREATE SET 
								org.updatedAt=datetime(),
								bp:BillingProfile_%s,
								bp.createdAt=$createdAt,
								bp.updatedAt=datetime(),
								bp.source=$source,
								bp.sourceOfTruth=$sourceOfTruth,
								bp.appSource=$appSource,
								bp.legalName=$legalName,
								bp.taxId=$taxId`, tenant)
	params := map[string]any{
		"tenant":           tenant,
		"billingProfileId": billingProfileId,
		"orgId":            data.OrganizationId,
		"createdAt":        data.CreatedAt,
		"source":           data.SourceFields.Source,
		"sourceOfTruth":    data.SourceFields.Source,
		"appSource":        data.SourceFields.AppSource,
		"legalName":        data.LegalName,
		"taxId":            data.TaxId,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}

func (r *billingProfileWriteRepository) Update(ctx context.Context, tenant, billingProfileId string, data BillingProfileUpdateFields) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "BillingProfileWriteRepository.Create")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, billingProfileId)
	tracing.LogObjectAsJson(span, "data", data)

	cypher := `MATCH (:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization {id:$orgId})-[:HAS_BILLING_PROFILE]->(bp:BillingProfile {id:$billingProfileId})
							SET bp.updatedAt=datetime(), org.updatedAt=datetime()
								`
	params := map[string]any{
		"tenant":           tenant,
		"billingProfileId": billingProfileId,
		"orgId":            data.OrganizationId,
	}
	if data.UpdateLegalName {
		cypher += ", bp.legalName=$legalName"
		params["legalName"] = data.LegalName
	}
	if data.UpdateTaxId {
		cypher += ", bp.taxId=$taxId"
		params["taxId"] = data.TaxId
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}

func (r *billingProfileWriteRepository) LinkEmailToBillingProfile(ctx context.Context, tenant, organizationId, billingProfileId, emailId string, primary bool) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "BillingProfileWriteRepository.LinkEmailToBillingProfile")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, billingProfileId)

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization {id:$orgId})-[:HAS_BILLING_PROFILE]->(bp:BillingProfile {id:$billingProfileId}),
					(t)<-[:EMAIL_ADDRESS_BELONGS_TO_TENANT]->(e:Email {id:$emailId})
				MERGE (e)<-[rel:HAS]-(bp)
				SET
					bp.updatedAt=datetime(),
					rel.primary=$primary
				WITH bp
				OPTIONAL MATCH (bp)-[rel2:HAS]->(oe:Email)
				WHERE oe.id <> $emailId AND $primary = true
				SET rel2.primary=false`
	params := map[string]any{
		"tenant":           tenant,
		"billingProfileId": billingProfileId,
		"orgId":            organizationId,
		"emailId":          emailId,
		"primary":          primary,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}

func (r *billingProfileWriteRepository) UnlinkEmailFromBillingProfile(ctx context.Context, tenant, organizationId, billingProfileId, emailId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "BillingProfileWriteRepository.UnlinkEmailFromBillingProfile")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, billingProfileId)

	cypher := `MATCH (:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(:Organization {id:$orgId})-[:HAS_BILLING_PROFILE]->(bp:BillingProfile {id:$billingProfileId})-[rel:HAS]->(e:Email {id:$emailId})
				SET bp.updatedAt=datetime()
				DELETE rel`
	params := map[string]any{
		"tenant":           tenant,
		"billingProfileId": billingProfileId,
		"orgId":            organizationId,
		"emailId":          emailId,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}

func (r *billingProfileWriteRepository) LinkLocationToBillingProfile(ctx context.Context, tenant, organizationId, billingProfileId, locationId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "BillingProfileWriteRepository.LinkLocationToBillingProfile")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, billingProfileId)

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization {id:$orgId})-[:HAS_BILLING_PROFILE]->(bp:BillingProfile {id:$billingProfileId}),
					(t)<-[:LOCATION_BELONGS_TO_TENANT]->(loc:Location {id:$locationId})
				MERGE (loc)<-[rel:HAS]-(bp)
				SET bp.updatedAt=datetime()`
	params := map[string]any{
		"tenant":           tenant,
		"billingProfileId": billingProfileId,
		"orgId":            organizationId,
		"locationId":       locationId,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}

func (r *billingProfileWriteRepository) UnlinkLocationFromBillingProfile(ctx context.Context, tenant, organizationId, billingProfileId, locationId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "BillingProfileWriteRepository.UnlinkLocationFromBillingProfile")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, billingProfileId)

	cypher := `MATCH (:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(:Organization {id:$orgId})-[:HAS_BILLING_PROFILE]->(bp:BillingProfile {id:$billingProfileId})-[rel:HAS]->(loc:Location {id:$locationId})
				SET bp.updatedAt=datetime()
				DELETE rel`
	params := map[string]any{
		"tenant":           tenant,
		"billingProfileId": billingProfileId,
		"orgId":            organizationId,
		"locationId":       locationId,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}
