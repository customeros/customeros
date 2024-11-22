package repository

import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/data_fields"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/constants"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"golang.org/x/net/context"
	"strings"
	"time"
)

type OrganizationWriteRepository interface {
	Save(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, organizationId string, data data_fields.OrganizationFields) error
	LinkWithDomain(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, organizationId, domain string) (bool, error)
	UnlinkFromDomain(ctx context.Context, tenant, organizationId, domain string) error
	ReplaceOwner(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, organizationId, userId string) error
	// Deprecated -> use Save with Hide property
	SetVisibility(ctx context.Context, tenant, organizationId string, hide bool) error
	UpdateLastTouchpoint(ctx context.Context, tenant, organizationId string, touchpointAt *time.Time, touchpointId, touchpointType string) error
	SetCustomerOsIdIfMissing(ctx context.Context, tenant, organizationId, customerOsId string) error
	LinkWithParentOrganization(ctx context.Context, tenant, organizationId, parentOrganizationId, subOrganizationType string) error
	UnlinkParentOrganization(ctx context.Context, tenant, organizationId, parentOrganizationId string) error
	UpdateArr(ctx context.Context, tenant, organizationId string) error
	UpdateRenewalSummary(ctx context.Context, tenant, organizationId string, likelihood *string, likelihoodOrder *int64, nextRenewalDate *time.Time) error
	WebScrapeRequested(ctx context.Context, tenant, organizationId, url string, attempt int64, requestedAt time.Time) error
	UpdateOnboardingStatus(ctx context.Context, tenant, organizationId, status, comments string, statusOrder *int64, updatedAt time.Time) error
	UpdateTimeProperty(ctx context.Context, tenant, organizationId, property string, value *time.Time) error
	UpdateFloatProperty(ctx context.Context, tenant, organizationId, property string, value float64) error
	UpdateStringProperty(ctx context.Context, tenant, organizationId, property string, value string) error
	Archive(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, organizationId string) error
	ResetEnrichAttempts(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, organizationId string) error
}

type organizationWriteRepository struct {
	driver   *neo4j.DriverWithContext
	database string
}

func NewOrganizationWriteRepository(driver *neo4j.DriverWithContext, database string) OrganizationWriteRepository {
	return &organizationWriteRepository{
		driver:   driver,
		database: database,
	}
}

func (r *organizationWriteRepository) prepareWriteSession(ctx context.Context) neo4j.SessionWithContext {
	return utils.NewNeo4jWriteSession(ctx, *r.driver, utils.WithDatabaseName(r.database))
}

func (r *organizationWriteRepository) Save(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, organizationId string, data data_fields.OrganizationFields) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationWriteRepository.Save")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)

	span.SetTag(tracing.SpanTagEntityId, organizationId)

	tracing.LogObjectAsJson(span, "data", data)

	_, err := utils.ExecuteWriteInTransaction(ctx, r.driver, r.database, tx, func(tx neo4j.ManagedTransaction) (any, error) {

		//create if not exists
		cypherCreate := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant}) MERGE(t)<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization:Organization_%s {id:$organizationId})
				ON CREATE SET
					org.source = $source,
					org.appSource = $appSource,
					org.createdAt = datetime(),
					org.updatedAt = datetime(),
					org.onboardingStatus = $onboardingStatus,
					org.hide=false`, tenant)
		paramsCreate := map[string]any{
			"tenant":           tenant,
			"organizationId":   organizationId,
			"source":           utils.IfNotNilString(data.Source),
			"appSource":        utils.IfNotNilString(data.AppSource),
			"onboardingStatus": string(neo4jenum.OnboardingStatusNotApplicable),
		}

		span.LogFields(log.String("cypherCreate", cypherCreate))
		tracing.LogObjectAsJson(span, "paramsCreate", paramsCreate)

		_, err := tx.Run(ctx, cypherCreate, paramsCreate)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}

		paramsUpdate := map[string]any{
			"tenant":         tenant,
			"organizationId": organizationId,
			"now":            utils.Now(),
		}

		cypherUpdate := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization:Organization_%s {id:$organizationId}) SET `, tenant)

		if data.Name != nil {
			cypherUpdate += `org.name = $name,`
			paramsUpdate["name"] = *data.Name
		}
		if data.Description != nil {
			cypherUpdate += `org.description = $description,`
			paramsUpdate["description"] = *data.Description
		}
		if data.Hide != nil {
			cypherUpdate += `org.hide = $hide,`
			cypherUpdate += `org.hiddenAt = CASE WHEN $hide = true THEN datetime() ELSE null END,`
			paramsUpdate["hide"] = *data.Hide
		}
		if data.Website != nil {
			cypherUpdate += `org.website = $website,`
			paramsUpdate["website"] = *data.Website
		}
		if data.Industry != nil {
			cypherUpdate += `org.industry = $industry,`
			paramsUpdate["industry"] = *data.Industry
		}
		if data.SubIndustry != nil {
			cypherUpdate += `org.subIndustry = $subIndustry,`
			paramsUpdate["subIndustry"] = *data.SubIndustry
		}
		if data.IndustryGroup != nil {
			cypherUpdate += `org.industryGroup = $industryGroup,`
			paramsUpdate["industryGroup"] = *data.IndustryGroup
		}
		if data.TargetAudience != nil {
			cypherUpdate += `org.targetAudience = $targetAudience,`
			paramsUpdate["targetAudience"] = *data.TargetAudience
		}
		if data.ValueProposition != nil {
			cypherUpdate += `org.valueProposition = $valueProposition,`
			paramsUpdate["valueProposition"] = data.ValueProposition
		}
		if data.LastFundingRound != nil {
			cypherUpdate += `org.lastFundingRound = $lastFundingRound,`
			paramsUpdate["lastFundingRound"] = *data.LastFundingRound
		}
		if data.LastFundingAmount != nil {
			cypherUpdate += `org.lastFundingAmount = $lastFundingAmount,`
			paramsUpdate["lastFundingAmount"] = *data.LastFundingAmount
		}
		if data.CustomerOsId != nil {
			cypherUpdate += `org.customerOsId = $customerOsId,`
			paramsUpdate["customerOsId"] = *data.CustomerOsId
		}
		if data.ReferenceId != nil {
			cypherUpdate += `org.referenceId = $referenceId,`
			paramsUpdate["referenceId"] = *data.ReferenceId
		}
		if data.Note != nil {
			cypherUpdate += `org.note = $note,`
			paramsUpdate["note"] = *data.Note
		}
		if data.IsPublic != nil {
			cypherUpdate += `org.isPublic = $isPublic,`
			paramsUpdate["isPublic"] = *data.IsPublic
		}
		if data.Employees != nil {
			cypherUpdate += `org.employees = $employees,`
			paramsUpdate["employees"] = *data.Employees
		}
		if data.Market != nil {
			cypherUpdate += `org.market = $market,`
			paramsUpdate["market"] = *data.Market
		}
		if data.YearFounded != nil {
			cypherUpdate += `org.yearFounded = $yearFounded,`
			paramsUpdate["yearFounded"] = *data.YearFounded
		}
		if data.Headquarters != nil {
			cypherUpdate += `org.headquarters = $headquarters,`
			paramsUpdate["headquarters"] = *data.Headquarters
		}
		if data.LogoUrl != nil {
			cypherUpdate += `org.logoUrl = $logoUrl,`
			paramsUpdate["logoUrl"] = *data.LogoUrl
		}
		if data.IconUrl != nil {
			cypherUpdate += `org.iconUrl = $iconUrl,`
			paramsUpdate["iconUrl"] = *data.IconUrl
		}
		if data.EmployeeGrowthRate != nil {
			cypherUpdate += `org.employeeGrowthRate = $employeeGrowthRate,`
			paramsUpdate["employeeGrowthRate"] = *data.EmployeeGrowthRate
		}
		if data.SlackChannelId != nil {
			cypherUpdate += `org.slackChannelId = $slackChannelId,`
			paramsUpdate["slackChannelId"] = *data.SlackChannelId
		}
		if data.Relationship != nil {
			cypherUpdate += `org.relationship = $relationship,`
			paramsUpdate["relationship"] = data.Relationship.String()
		}
		if data.Stage != nil {
			cypherUpdate += `org.stage = $stage,`
			cypherUpdate += `org.stageUpdatedAt = CASE WHEN (org.stage is null OR org.stage = '') AND (org.stage is null OR org.stage <> $stage) THEN $now ELSE org.stageUpdatedAt END,`
			paramsUpdate["stage"] = data.Stage.String()
		}
		if data.LeadSource != nil {
			cypherUpdate += `org.leadSource = $leadSource,`
			paramsUpdate["leadSource"] = *data.LeadSource
		}
		if data.IcpFit != nil {
			cypherUpdate += `org.icpFit = $icpFit,`
			paramsUpdate["icpFit"] = *data.IcpFit
		}
		if utils.IfNotNilString(data.EnrichDomain) != "" && utils.IfNotNilString(data.EnrichSource) != "" {
			cypherUpdate += `org.enrichDomain = $enrichDomain, org.enrichSource = $enrichSource, org.enrichedAt = $enrichedAt,`
			paramsUpdate["enrichDomain"] = *data.EnrichDomain
			paramsUpdate["enrichSource"] = *data.EnrichSource
			paramsUpdate["enrichedAt"] = utils.Now()
		}
		cypherUpdate += `org.updatedAt = datetime()`

		span.LogFields(log.String("cypherUpdate", cypherUpdate))
		tracing.LogObjectAsJson(span, "paramsUpdate", paramsUpdate)

		_, err = tx.Run(ctx, cypherUpdate, paramsUpdate)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}

		return nil, nil
	})

	return err
}

func (r *organizationWriteRepository) LinkWithDomain(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, organizationId, domain string) (bool, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationWriteRepository.LinkWithDomain")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, organizationId)

	cypher := `MERGE (d:Domain {domain: $domain}) 
  				ON CREATE SET 	d.createdAt = datetime(), 
                				d.updatedAt = datetime()
				WITH d
				MATCH (t:Tenant {name: $tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization {id: $organizationId})
				OPTIONAL MATCH (d)<-[:HAS_DOMAIN]-(otherOrg:Organization)-[:ORGANIZATION_BELONGS_TO_TENANT]->(t)
				WHERE org <> otherOrg
				WITH d, org, COUNT(otherOrg) AS existingOrgCount
				FOREACH (_ IN CASE WHEN existingOrgCount = 0 THEN [1] ELSE [] END | 
 					MERGE (org)-[rel:HAS_DOMAIN]->(d)
  					SET org.updatedAt = datetime()
				)
				RETURN existingOrgCount = 0 AS linked`
	params := map[string]any{
		"tenant":         tenant,
		"organizationId": organizationId,
		"domain":         strings.ToLower(domain),
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	result, err := utils.ExecuteWriteInTransaction(ctx, r.driver, r.database, tx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		resultWithContext, err := tx.Run(ctx, cypher, params)
		return utils.ExtractSingleRecordFirstValueAsType[bool](ctx, resultWithContext, err)
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return false, err
	}
	span.LogFields(log.Bool("result", result.(bool)))
	return result.(bool), err
}

func (r *organizationWriteRepository) UnlinkFromDomain(ctx context.Context, tenant, organizationId, domain string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationWriteRepository.UnlinkFromDomain")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, organizationId)

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization {id:$organizationId})
		 MATCH (org)-[rel:HAS_DOMAIN]->(d:Domain {domain:$domain})
		 SET org.updatedAt = datetime()
		 DELETE rel`
	params := map[string]any{
		"tenant":         tenant,
		"organizationId": organizationId,
		"domain":         strings.ToLower(domain),
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}

func (r *organizationWriteRepository) ReplaceOwner(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, organizationId, userId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationWriteRepository.ReplaceOwner")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)

	span.LogFields(log.String("organizationId", organizationId), log.String("userId", userId))

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization {id:$organizationId})
			OPTIONAL MATCH (:User)-[rel:OWNS]->(org)
			DELETE rel
			WITH org, t
			MATCH (t)<-[:USER_BELONGS_TO_TENANT]-(u:User {id:$userId})
			WHERE (u.internal=false OR u.internal is null) AND (u.bot=false OR u.bot is null) AND (u.test=false OR u.test is null)
			MERGE (u)-[:OWNS]->(org)
			SET org.updatedAt=datetime(), org.sourceOfTruth=$source`

	params := map[string]any{
		"tenant":         tenant,
		"organizationId": organizationId,
		"userId":         userId,
		"source":         constants.SourceOpenline,
		"now":            utils.Now(),
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	_, err := utils.ExecuteWriteInTransaction(ctx, r.driver, r.database, tx, func(tx neo4j.ManagedTransaction) (any, error) {

		_, err := tx.Run(ctx, cypher, params)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}

		return nil, nil
	})

	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}

func (r *organizationWriteRepository) SetVisibility(ctx context.Context, tenant, organizationId string, hide bool) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationWriteRepository.SetVisibility")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, organizationId)
	span.LogFields(log.Bool("hide", hide))

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization {id:$id})
			WHERE org:Organization_%s
		 SET	org.hide = $hide,
				org.hiddenAt = CASE WHEN $hide = true THEN datetime() ELSE org.hiddenAt END,
				org.updatedAt = datetime()`, tenant)
	params := map[string]any{
		"id":     organizationId,
		"tenant": tenant,
		"hide":   hide,
		"now":    utils.Now(),
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}

func (r *organizationWriteRepository) UpdateLastTouchpoint(ctx context.Context, tenant, organizationId string, touchpointAt *time.Time, touchpointId, touchpointType string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationWriteRepository.UpdateLastTouchpoint")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	span.LogFields(log.String("organizationId", organizationId), log.String("touchpointId", touchpointId), log.Object("touchpointAt", touchpointAt))

	cypher := `MATCH (:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization {id:$organizationId})
		 SET 	org.updatedAt = CASE WHEN org.lastTouchpointId <> $touchpointId THEN datetime() ELSE org.updatedAt END,
				org.lastTouchpointAt=$touchpointAt, 
				org.lastTouchpointId=$touchpointId, 
				org.lastTouchpointType=$touchpointType`
	params := map[string]any{
		"tenant":         tenant,
		"organizationId": organizationId,
		"touchpointAt":   utils.TimePtrAsAny(touchpointAt),
		"touchpointId":   touchpointId,
		"touchpointType": touchpointType,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}

func (r *organizationWriteRepository) SetCustomerOsIdIfMissing(ctx context.Context, tenant, organizationId, customerOsId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationWriteRepository.SetCustomerOsIdIfMissing")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, organizationId)
	span.LogFields(log.String("customerOsId", customerOsId))

	cypher := `MATCH (:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization {id:$organizationId})
		 SET org.customerOsId = CASE WHEN (org.customerOsId IS NULL OR org.customerOsId = '') AND $customerOsId <> '' THEN $customerOsId ELSE org.customerOsId END,
			org.updatedAt = datetime()`
	params := map[string]any{
		"tenant":         tenant,
		"organizationId": organizationId,
		"customerOsId":   customerOsId,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}

func (r *organizationWriteRepository) LinkWithParentOrganization(ctx context.Context, tenant, organizationId, parentOrganizationId, subOrganizationType string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationWriteRepository.LinkWithParentOrganization")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, organizationId)
	span.LogFields(log.String("parentOrganizationId", parentOrganizationId), log.String("subOrganizationType", subOrganizationType))

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(parent:Organization {id:$parentOrganizationId}),
		 			(t)<-[:ORGANIZATION_BELONGS_TO_TENANT]-(sub:Organization {id:$subOrganizationId}) 
		 	MERGE (sub)-[rel:SUBSIDIARY_OF]->(parent) 
		 		ON CREATE SET rel.type=$type 
		 		ON MATCH SET rel.type=$type
				SET sub.updatedAt = datetime(),
					parent.updatedAt = datetime()`
	params := map[string]any{
		"tenant":               tenant,
		"subOrganizationId":    organizationId,
		"parentOrganizationId": parentOrganizationId,
		"type":                 subOrganizationType,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}

func (r *organizationWriteRepository) UnlinkParentOrganization(ctx context.Context, tenant, organizationId, parentOrganizationId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationWriteRepository.UnlinkParentOrganization")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, organizationId)
	span.LogFields(log.String("parentOrganizationId", parentOrganizationId))

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(parent:Organization {id:$parentOrganizationId})<-[rel:SUBSIDIARY_OF]-(sub:Organization {id:$subOrganizationId})-[:ORGANIZATION_BELONGS_TO_TENANT]->(t)
		 		DELETE rel
				SET sub.updatedAt = datetime(),
					parent.updatedAt = datetime()`
	params := map[string]any{
		"tenant":               tenant,
		"subOrganizationId":    organizationId,
		"parentOrganizationId": parentOrganizationId,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}

func (r *organizationWriteRepository) UpdateArr(ctx context.Context, tenant, organizationId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationWriteRepository.UpdateArr")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, organizationId)

	cypher := `MATCH (t:Tenant {name: $tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization {id: $organizationId})
				OPTIONAL MATCH (org)-[:HAS_CONTRACT]->(c:Contract) WHERE c.status <> $statusDraft
				WITH *
				OPTIONAL MATCH (c)-[:ACTIVE_RENEWAL]->(op:Opportunity)
				WITH org, COALESCE(sum(op.amount), 0) as arr, COALESCE(sum(op.maxAmount), 0) as maxArr
				SET org.renewalForecastArr = arr, org.renewalForecastMaxArr = maxArr, org.updatedAt = datetime()`
	params := map[string]any{
		"tenant":         tenant,
		"organizationId": organizationId,
		"statusDraft":    neo4jenum.ContractStatusDraft.String(),
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}

func (r *organizationWriteRepository) UpdateRenewalSummary(ctx context.Context, tenant, organizationId string, likelihood *string, likelihoodOrder *int64, nextRenewalDate *time.Time) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationWriteRepository.UpdateRenewalSummary")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, organizationId)
	span.LogFields(log.Object("likelihood", likelihood), log.Object("likelihoodOrder", likelihoodOrder), log.Object("nextRenewalDate", nextRenewalDate))

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization {id:$organizationId})
				SET org.derivedRenewalLikelihood = $derivedRenewalLikelihood,
					org.derivedRenewalLikelihoodOrder = $derivedRenewalLikelihoodOrder,
					org.derivedNextRenewalAt = $derivedNextRenewalAt,
					org.updatedAt = datetime()`
	params := map[string]any{
		"tenant":                        tenant,
		"organizationId":                organizationId,
		"derivedRenewalLikelihood":      likelihood,
		"derivedRenewalLikelihoodOrder": likelihoodOrder,
		"derivedNextRenewalAt":          utils.TimePtrAsAny(nextRenewalDate),
		"now":                           utils.Now(),
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}

func (r *organizationWriteRepository) WebScrapeRequested(ctx context.Context, tenant, organizationId, url string, attempt int64, requestedAt time.Time) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationWriteRepository.WebScrapeRequested")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, organizationId)
	span.LogFields(log.String("url", url))

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization {id:$organizationId})
		 	SET org.webScrapeLastRequestedAt=$requestedAt, 
				org.webScrapeLastRequestedUrl=$url, 
				org.webScrapeAttempts=$attempt`
	params := map[string]any{
		"tenant":         tenant,
		"organizationId": organizationId,
		"url":            url,
		"attempt":        attempt,
		"requestedAt":    requestedAt,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}

func (r *organizationWriteRepository) UpdateOnboardingStatus(ctx context.Context, tenant, organizationId, status, comments string, statusOrder *int64, updatedAt time.Time) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationWriteRepository.UpdateOnboardingStatus")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, organizationId)

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization {id:$organizationId})
				SET org.onboardingUpdatedAt = CASE WHEN org.onboardingStatus IS NULL OR org.onboardingStatus <> $status THEN $updatedAt ELSE org.onboardingUpdatedAt END,
					org.onboardingStatus=$status,
					org.onboardingStatusOrder=$statusOrder,
					org.onboardingComments=$comments,
					org.onboardingUpdatedAt=$updatedAt,
					org.updatedAt=datetime()`
	params := map[string]any{
		"tenant":         tenant,
		"organizationId": organizationId,
		"status":         status,
		"statusOrder":    statusOrder,
		"comments":       comments,
		"updatedAt":      updatedAt,
		"now":            utils.Now(),
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}

func (r *organizationWriteRepository) UpdateTimeProperty(ctx context.Context, tenant, organizationId, property string, value *time.Time) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationWriteRepository.UpdateTimeProperty")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, organizationId)
	span.LogFields(log.String("property", property), log.Object("value", value))

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name: $tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization {id: $organizationId})
			SET org.%s = $value`, property)
	params := map[string]any{
		"tenant":         tenant,
		"organizationId": organizationId,
		"property":       property,
		"value":          utils.TimePtrAsAny(value),
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}

func (r *organizationWriteRepository) UpdateFloatProperty(ctx context.Context, tenant, organizationId, property string, value float64) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationWriteRepository.UpdateFloatProperty")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, organizationId)
	span.LogFields(log.String("property", property), log.Float64("value", value))

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name: $tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization {id: $organizationId})
			SET org.%s = $value`, property)
	params := map[string]any{
		"tenant":         tenant,
		"organizationId": organizationId,
		"property":       property,
		"value":          value,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}

func (r *organizationWriteRepository) UpdateStringProperty(ctx context.Context, tenant, organizationId, property string, value string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationWriteRepository.UpdateFloatProperty")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, organizationId)
	span.LogFields(log.String("property", property), log.String("value", value))

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name: $tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization {id: $organizationId})
			SET org.%s = $value`, property)
	params := map[string]any{
		"tenant":         tenant,
		"organizationId": organizationId,
		"property":       property,
		"value":          value,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}

func (r *organizationWriteRepository) Archive(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, organizationId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationRepository.Delete")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	cypher := fmt.Sprintf(`MATCH (org:Organization {id:$organizationId})-[currentRel:ORGANIZATION_BELONGS_TO_TENANT]->(t:Tenant {name:$tenant})
			MERGE (org)-[newRel:ARCHIVED]->(t)
			SET org.archived=true, org.archivedAt=$now, org.updatedAt=datetime(), org:ArchivedOrganization_%s
            DELETE currentRel
			REMOVE org:Organization_%s`, tenant, tenant)
	params := map[string]any{
		"tenant":         tenant,
		"organizationId": organizationId,
		"now":            utils.Now(),
	}

	span.LogFields(log.String("cypher", cypher))
	span.LogFields(log.Object("params", params))

	_, err := utils.ExecuteWriteInTransaction(ctx, r.driver, r.database, tx, func(tx neo4j.ManagedTransaction) (any, error) {

		_, err := tx.Run(ctx, cypher, params)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}

		return nil, nil
	})

	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}

func (r *organizationWriteRepository) ResetEnrichAttempts(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, organizationId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationWriteRepository.ResetEnrichAttempts")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	tracing.TagEntity(span, organizationId)

	cypher := `MATCH (t:Tenant {name: $tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(o:Organization {id: $organizationId})
	WHERE o.enrichedAt IS NULL
	REMOVE o.techEnrichAttempts, o.techEnrichRequestedAt`
	params := map[string]any{
		"tenant":         tenant,
		"organizationId": organizationId,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	_, err := utils.ExecuteWriteInTransaction(ctx, r.driver, r.database, tx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, cypher, params)
		if err != nil {
			return nil, err
		}
		return nil, nil
	})
	if err != nil {
		tracing.TraceErr(span, err)
	}

	return err
}
