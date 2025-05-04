package neo4j_repository

import (
	"fmt"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/data_fields"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	"github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/constants"
	neo4jenum "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/enum"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"

	"golang.org/x/net/context"
	"strings"
	"time"
)

type OrganizationWriteRepository interface {
	Save(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, organizationId string, data data_fields.OrganizationFields) error
	LinkWithDomain(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, organizationId, domain string) (bool, error)
	UnlinkDomain(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, organizationId, domain string) error
	ReplaceOwner(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, organizationId, userId string) error
	// Deprecated -> use Save with Hide property
	SetVisibility(ctx context.Context, tenant, organizationId string, hide bool) error
	UpdateLastTouchpoint(ctx context.Context, tenant, organizationId string, touchpointAt *time.Time, touchpointId, touchpointType string) error
	SetCustomerOsIdIfMissing(ctx context.Context, tenant, organizationId, customerOsId string) error
	LinkWithParentOrganization(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, subOrganizationId, parentOrganizationId, subOrganizationType string) error
	UnlinkParentOrganization(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, subOrganizationId, parentOrganizationId string) error
	UpdateArr(ctx context.Context, tenant, organizationId string) error
	UpdateRenewalSummary(ctx context.Context, tenant, organizationId string, likelihood *string, likelihoodOrder *int64, nextRenewalDate *time.Time) error
	WebScrapeRequested(ctx context.Context, tenant, organizationId, url string, attempt int64, requestedAt time.Time) error
	UpdateOnboardingStatus(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, organizationId string, data data_fields.OrganizationOnboardingStatusFields) error
	UpdateTimeProperty(ctx context.Context, tenant, organizationId, property string, value *time.Time) error
	UpdateFloatProperty(ctx context.Context, tenant, organizationId, property string, value float64) error
	UpdateStringProperty(ctx context.Context, tenant, organizationId, property string, value string) error
	Archive(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, organizationId string) error
	ResetEnrichAttempts(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, organizationId string) error
	RefreshContactCountByOrgId(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, organizationId string) error
	RefreshContactCountByContactId(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, contactId string) error
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

func (r *organizationWriteRepository) Save(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, organizationId string, data data_fields.OrganizationFields) error {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "OrganizationWriteRepository.Save")
	defer spans.Finish()

	spans.TagEntity(organizationId)

	spans.LogObjectAsJson("data", data)

	_, err := utils.ExecuteWriteInTransaction(ctx, r.driver, r.database, tx, func(tx neo4j.ManagedTransaction) (any, error) {

		//create if not exists
		cypherCreate := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant}) MERGE(t)<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization:Organization_%s {id:$organizationId})
				ON CREATE SET
					org.source = $source,
					org.appSource = $appSource,
					org.createdAt = datetime(),
					org.updatedAt = datetime(),
					org.onboardingStatus = $onboardingStatus,
					org.lastTouchpointAt = datetime(),
					org.lastTouchpointType = $lastTouchpointType,
					org.hide=false`, tenant)
		paramsCreate := map[string]any{
			"tenant":             tenant,
			"organizationId":     organizationId,
			"source":             utils.IfNotNilString(data.Source),
			"appSource":          utils.IfNotNilString(data.AppSource),
			"onboardingStatus":   string(neo4jenum.OnboardingStatusNotApplicable),
			"lastTouchpointType": neo4jenum.TouchpointTypeActionCreated.String(),
		}

		spans.LogKV("cypherCreate", cypherCreate)
		spans.LogObjectAsJson("paramsCreate", paramsCreate)

		_, err := tx.Run(ctx, cypherCreate, paramsCreate)
		if err != nil {
			spans.TraceError(err)
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
			cypherUpdate += `org.icpFitUpdatedAt = CASE WHEN org.icpFit IS NULL OR org.icpFit <> $icpFit THEN datetime() ELSE org.icpFitUpdatedAt END,`
			cypherUpdate += `org.icpCheckedAt = CASE 
			WHEN $icpFit = $icp_Fit OR $icpFit = $icp_NotFit THEN datetime() ELSE
				CASE WHEN $icpFit = $icp_NotSet THEN NULL ELSE org.icpCheckedAt END
			END,`
			paramsUpdate["icpFit"] = (*data.IcpFit).String()
			paramsUpdate["icp_Fit"] = enum.IcpIsFit.String()
			paramsUpdate["icp_NotFit"] = enum.IcpNotFit.String()
			paramsUpdate["icp_NotSet"] = enum.IcpNotSet.String()
		}
		if data.IcpFitReasons != nil {
			cypherUpdate += `org.icpFitReasons = $icpFitReasons,`
			paramsUpdate["icpFitReasons"] = *data.IcpFitReasons
		}
		if utils.IfNotNilString(data.EnrichDomain) != "" && utils.IfNotNilString(data.EnrichSource) != "" {
			cypherUpdate += `org.enrichDomain = $enrichDomain, org.enrichSource = $enrichSource, org.enrichedAt = $enrichedAt,`
			paramsUpdate["enrichDomain"] = *data.EnrichDomain
			paramsUpdate["enrichSource"] = *data.EnrichSource
			paramsUpdate["enrichedAt"] = utils.Now()
		}
		cypherUpdate += `org.updatedAt = datetime()`

		spans.LogKV("cypherUpdate", cypherUpdate)
		spans.LogObjectAsJson("paramsUpdate", paramsUpdate)

		_, err = tx.Run(ctx, cypherUpdate, paramsUpdate)
		if err != nil {
			spans.TraceError(err)
			return nil, err
		}

		return nil, nil
	})

	return err
}

func (r *organizationWriteRepository) LinkWithDomain(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, organizationId, domain string) (bool, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "OrganizationWriteRepository.LinkWithDomain")
	defer spans.Finish()

	spans.TagEntity(organizationId)

	cypher := `MATCH (d:Domain {domain: $domain}) 
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
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	result, err := utils.ExecuteWriteInTransaction(ctx, r.driver, r.database, tx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		resultWithContext, err := tx.Run(ctx, cypher, params)
		return utils.ExtractSingleRecordFirstValueAsType[bool](ctx, resultWithContext, err)
	})
	if err != nil {
		spans.TraceError(err)
		return false, err
	}
	spans.LogKV("result", result.(bool))
	return result.(bool), err
}

func (r *organizationWriteRepository) UnlinkDomain(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, organizationId, domain string) error {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "OrganizationWriteRepository.UnlinkDomain")
	defer spans.Finish()

	spans.TagEntity(organizationId)

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization {id:$organizationId})
		 MATCH (org)-[rel:HAS_DOMAIN]->(d:Domain {domain:$domain})
		 SET org.updatedAt = datetime()
		 DELETE rel`
	params := map[string]any{
		"tenant":         tenant,
		"organizationId": organizationId,
		"domain":         strings.ToLower(domain),
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	_, err := utils.ExecuteWriteInTransaction(ctx, r.driver, r.database, tx, func(tx neo4j.ManagedTransaction) (any, error) {
		return tx.Run(ctx, cypher, params)
	})
	if err != nil {
		spans.TraceError(err)
		return err
	}
	return nil
}

func (r *organizationWriteRepository) ReplaceOwner(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, organizationId, userId string) error {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "OrganizationWriteRepository.ReplaceOwner")
	defer spans.Finish()

	spans.LogKV("organizationId", organizationId)
	spans.LogKV("userId", userId)

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization {id:$organizationId})
			OPTIONAL MATCH (:User)-[rel:OWNS]->(org)
			DELETE rel
			WITH org, t
			MATCH (t)<-[:USER_BELONGS_TO_TENANT]-(u:User {id:$userId})
			WHERE (u.internal=false OR u.internal is null) AND (u.bot=false OR u.bot is null) AND (u.test=false OR u.test is null)
			MERGE (u)-[:OWNS]->(org)
			SET org.updatedAt=datetime()`

	params := map[string]any{
		"tenant":         tenant,
		"organizationId": organizationId,
		"userId":         userId,
		"source":         constants.SourceOpenline,
		"now":            utils.Now(),
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	_, err := utils.ExecuteWriteInTransaction(ctx, r.driver, r.database, tx, func(tx neo4j.ManagedTransaction) (any, error) {

		_, err := tx.Run(ctx, cypher, params)
		if err != nil {
			spans.TraceError(err)
			return nil, err
		}

		return nil, nil
	})

	if err != nil {
		spans.TraceError(err)
		return err
	}

	return nil
}

func (r *organizationWriteRepository) SetVisibility(ctx context.Context, tenant, organizationId string, hide bool) error {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "OrganizationWriteRepository.SetVisibility")
	defer spans.Finish()

	spans.TagEntity(organizationId)
	spans.LogKV("hide", hide)

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
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		spans.TraceError(err)
	}
	return err
}

func (r *organizationWriteRepository) UpdateLastTouchpoint(ctx context.Context, tenant, organizationId string, touchpointAt *time.Time, touchpointId, touchpointType string) error {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "OrganizationWriteRepository.UpdateLastTouchpoint")
	defer spans.Finish()

	spans.LogKV("organizationId", organizationId)
	spans.LogKV("touchpointId", touchpointId)
	spans.LogKV("touchpointAt", touchpointAt)

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
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		spans.TraceError(err)
	}
	return err
}

func (r *organizationWriteRepository) SetCustomerOsIdIfMissing(ctx context.Context, tenant, organizationId, customerOsId string) error {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "OrganizationWriteRepository.SetCustomerOsIdIfMissing")
	defer spans.Finish()

	spans.TagEntity(organizationId)
	spans.LogKV("customerOsId", customerOsId)

	cypher := `MATCH (:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization {id:$organizationId})
		 SET org.customerOsId = CASE WHEN (org.customerOsId IS NULL OR org.customerOsId = '') AND $customerOsId <> '' THEN $customerOsId ELSE org.customerOsId END,
			org.updatedAt = datetime()`
	params := map[string]any{
		"tenant":         tenant,
		"organizationId": organizationId,
		"customerOsId":   customerOsId,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		spans.TraceError(err)
	}
	return err
}

func (r *organizationWriteRepository) LinkWithParentOrganization(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, subOrganizationId, parentOrganizationId, subOrganizationType string) error {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "OrganizationWriteRepository.LinkWithParentOrganization")
	defer spans.Finish()

	spans.TagEntity(subOrganizationId)
	spans.LogKV("parentOrganizationId", parentOrganizationId)
	spans.LogKV("subOrganizationType", subOrganizationType)

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(parent:Organization {id:$parentOrganizationId}),
		 			(t)<-[:ORGANIZATION_BELONGS_TO_TENANT]-(sub:Organization {id:$subOrganizationId}) 
		 	MERGE (sub)-[rel:SUBSIDIARY_OF]->(parent) 
		 		ON CREATE SET rel.type=$type 
		 		ON MATCH SET rel.type=$type
				SET sub.updatedAt = datetime(),
					parent.updatedAt = datetime()`
	params := map[string]any{
		"tenant":               tenant,
		"subOrganizationId":    subOrganizationId,
		"parentOrganizationId": parentOrganizationId,
		"type":                 subOrganizationType,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	_, err := utils.ExecuteWriteInTransaction(ctx, r.driver, r.database, tx, func(tx neo4j.ManagedTransaction) (any, error) {
		return tx.Run(ctx, cypher, params)
	})
	if err != nil {
		spans.TraceError(err)
		return err
	}
	return nil
}

func (r *organizationWriteRepository) UnlinkParentOrganization(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, subOrganizationId, parentOrganizationId string) error {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "OrganizationWriteRepository.UnlinkParentOrganization")
	defer spans.Finish()

	spans.TagEntity(subOrganizationId)
	spans.LogKV("parentOrganizationId", parentOrganizationId)

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(parent:Organization {id:$parentOrganizationId})<-[rel:SUBSIDIARY_OF]-(sub:Organization {id:$subOrganizationId})-[:ORGANIZATION_BELONGS_TO_TENANT]->(t)
		 		DELETE rel
				SET sub.updatedAt = datetime(),
					parent.updatedAt = datetime()`
	params := map[string]any{
		"tenant":               tenant,
		"subOrganizationId":    subOrganizationId,
		"parentOrganizationId": parentOrganizationId,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	_, err := utils.ExecuteWriteInTransaction(ctx, r.driver, r.database, tx, func(tx neo4j.ManagedTransaction) (any, error) {
		return tx.Run(ctx, cypher, params)
	})
	if err != nil {
		spans.TraceError(err)
		return err
	}
	return nil
}

func (r *organizationWriteRepository) UpdateArr(ctx context.Context, tenant, organizationId string) error {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "OrganizationWriteRepository.UpdateArr")
	defer spans.Finish()

	spans.TagEntity(organizationId)

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
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		spans.TraceError(err)
	}
	return err
}

func (r *organizationWriteRepository) UpdateRenewalSummary(ctx context.Context, tenant, organizationId string, likelihood *string, likelihoodOrder *int64, nextRenewalDate *time.Time) error {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "OrganizationWriteRepository.UpdateRenewalSummary")
	defer spans.Finish()

	spans.TagEntity(organizationId)
	spans.LogKV("likelihood", utils.IfNotNilString(likelihood))
	spans.LogKV("likelihoodOrder", likelihoodOrder)
	spans.LogKV("nextRenewalDate", nextRenewalDate)

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
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		spans.TraceError(err)
	}
	return err
}

func (r *organizationWriteRepository) WebScrapeRequested(ctx context.Context, tenant, organizationId, url string, attempt int64, requestedAt time.Time) error {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "OrganizationWriteRepository.WebScrapeRequested")
	defer spans.Finish()

	spans.TagEntity(organizationId)
	spans.LogKV("url", url)

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
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		spans.TraceError(err)
	}
	return err
}

func (r *organizationWriteRepository) UpdateOnboardingStatus(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, organizationId string, data data_fields.OrganizationOnboardingStatusFields) error {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "OrganizationWriteRepository.UpdateOnboardingStatus")
	defer spans.Finish()

	spans.TagEntity(organizationId)

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization {id:$organizationId})
				SET org.onboardingUpdatedAt = CASE WHEN org.onboardingStatus IS NULL OR (org.onboardingStatus <> $status AND $status IS NULL) THEN datetime() ELSE org.onboardingUpdatedAt END,
					org.onboardingStatus = CASE WHEN $status IS NULL THEN org.onboardingStatus ELSE $status END,
					org.onboardingStatusOrder = CASE WHEN $status IS NULL THEN org.onboardingStatusOrder ELSE $statusOrder END,
					org.onboardingComments = CASE WHEN $comments IS NULL THEN org.onboardingComments ELSE $comments END,
					org.updatedAt=datetime()`
	params := map[string]any{
		"tenant":         tenant,
		"organizationId": organizationId,
		"comments":       data.Comments,
	}
	if data.Status != nil {
		params["status"] = data.Status.String()
		params["statusOrder"] = data.Status.GetOrder()
	} else {
		params["status"] = nil
		params["statusOrder"] = nil
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	_, err := utils.ExecuteWriteInTransaction(ctx, r.driver, r.database, tx, func(tx neo4j.ManagedTransaction) (any, error) {
		return tx.Run(ctx, cypher, params)
	})
	if err != nil {
		spans.TraceError(err)
		return err
	}
	return nil
}

func (r *organizationWriteRepository) UpdateTimeProperty(ctx context.Context, tenant, organizationId, property string, value *time.Time) error {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "OrganizationWriteRepository.UpdateTimeProperty")
	defer spans.Finish()

	spans.TagEntity(organizationId)
	spans.LogKV("property", property)
	spans.LogObjectAsJson("value", value)

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name: $tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization {id: $organizationId})
			SET org.%s = $value`, property)
	params := map[string]any{
		"tenant":         tenant,
		"organizationId": organizationId,
		"property":       property,
		"value":          utils.TimePtrAsAny(value),
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		spans.TraceError(err)
	}
	return err
}

func (r *organizationWriteRepository) UpdateFloatProperty(ctx context.Context, tenant, organizationId, property string, value float64) error {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "OrganizationWriteRepository.UpdateFloatProperty")
	defer spans.Finish()

	spans.TagEntity(organizationId)
	spans.LogKV("property", property)
	spans.LogKV("value", value)

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name: $tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization {id: $organizationId})
			SET org.%s = $value`, property)
	params := map[string]any{
		"tenant":         tenant,
		"organizationId": organizationId,
		"property":       property,
		"value":          value,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		spans.TraceError(err)
	}
	return err
}

func (r *organizationWriteRepository) UpdateStringProperty(ctx context.Context, tenant, organizationId, property string, value string) error {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "OrganizationWriteRepository.UpdateFloatProperty")
	defer spans.Finish()

	spans.TagEntity(organizationId)
	spans.LogKV("property", property)
	spans.LogKV("value", value)

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name: $tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization {id: $organizationId})
			SET org.%s = $value`, property)
	params := map[string]any{
		"tenant":         tenant,
		"organizationId": organizationId,
		"property":       property,
		"value":          value,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		spans.TraceError(err)
	}
	return err
}

func (r *organizationWriteRepository) Archive(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, organizationId string) error {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "OrganizationRepository.Delete")
	defer spans.Finish()

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

	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	_, err := utils.ExecuteWriteInTransaction(ctx, r.driver, r.database, tx, func(tx neo4j.ManagedTransaction) (any, error) {

		_, err := tx.Run(ctx, cypher, params)
		if err != nil {
			spans.TraceError(err)
			return nil, err
		}

		return nil, nil
	})

	if err != nil {
		spans.TraceError(err)
		return err
	}

	return nil
}

func (r *organizationWriteRepository) ResetEnrichAttempts(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, organizationId string) error {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "OrganizationWriteRepository.ResetEnrichAttempts")
	defer spans.Finish()

	spans.TagEntity(organizationId)

	cypher := `MATCH (t:Tenant {name: $tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(o:Organization {id: $organizationId})
	WHERE o.enrichedAt IS NULL
	REMOVE o.techEnrichAttempts, o.techEnrichRequestedAt`
	params := map[string]any{
		"tenant":         tenant,
		"organizationId": organizationId,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	_, err := utils.ExecuteWriteInTransaction(ctx, r.driver, r.database, tx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, cypher, params)
		if err != nil {
			return nil, err
		}
		return nil, nil
	})
	if err != nil {
		spans.TraceError(err)
	}

	return err
}

func (r *organizationWriteRepository) RefreshContactCountByOrgId(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, organizationId string) error {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "OrganizationWriteRepository.RefreshContactCountByOrgId")
	defer spans.Finish()

	cypher := `MATCH (:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization {id:$organizationId})
			WITH org
				OPTIONAL MATCH (org)--(:JobRole)--(c:Contact)
				WHERE c.hide=false OR c.hide IS NULL
			WITH org, COUNT(DISTINCT c) AS contactCount
				SET org.derivedContactCount = contactCount`
	params := map[string]any{
		"tenant":         tenant,
		"organizationId": organizationId,
	}

	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	_, err := utils.ExecuteWriteInTransaction(ctx, r.driver, r.database, tx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, cypher, params)
		if err != nil {
			return nil, err
		}
		return nil, nil
	})
	if err != nil {
		spans.TraceError(err)
	}

	return err
}

func (r *organizationWriteRepository) RefreshContactCountByContactId(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, contactId string) error {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "OrganizationWriteRepository.RefreshContactCountByContactId")
	defer spans.Finish()

	cypher := `MATCH (:Tenant {name:$tenant})<-[:CONTACT_BELONGS_TO_TENANT]-(:Contact {id:$contactId})--(:JobRole)--(org:Organization)
			WITH org
				OPTIONAL MATCH (org)--(:JobRole)--(c:Contact)
				WHERE c.hide=false OR c.hide IS NULL
			WITH org, COUNT(DISTINCT c) AS contactCount
				SET org.derivedContactCount = contactCount, org.updatedAt = datetime()`
	params := map[string]any{
		"tenant":    tenant,
		"contactId": contactId,
	}

	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	_, err := utils.ExecuteWriteInTransaction(ctx, r.driver, r.database, tx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, cypher, params)
		if err != nil {
			return nil, err
		}
		return nil, nil
	})
	if err != nil {
		spans.TraceError(err)
	}

	return err
}
