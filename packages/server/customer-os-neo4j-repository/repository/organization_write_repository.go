package repository

import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/constants"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"golang.org/x/net/context"
	"strings"
	"time"
)

type OrganizationCreateFields struct {
	SourceFields       model.Source                       `json:"sourceFields"`
	CreatedAt          time.Time                          `json:"createdAt"`
	UpdatedAt          time.Time                          `json:"updatedAt"`
	Name               string                             `json:"name"`
	Hide               bool                               `json:"hide"`
	Description        string                             `json:"description"`
	Website            string                             `json:"website"`
	Industry           string                             `json:"industry"`
	SubIndustry        string                             `json:"subIndustry"`
	IndustryGroup      string                             `json:"industryGroup"`
	TargetAudience     string                             `json:"targetAudience"`
	ValueProposition   string                             `json:"valueProposition"`
	IsPublic           bool                               `json:"isPublic"`
	IsCustomer         bool                               `json:"isCustomer"`
	Employees          int64                              `json:"employees"`
	Market             string                             `json:"market"`
	LastFundingRound   string                             `json:"lastFundingRound"`
	LastFundingAmount  string                             `json:"lastFundingAmount"`
	ReferenceId        string                             `json:"referenceId"`
	Note               string                             `json:"note"`
	LogoUrl            string                             `json:"logoUrl"`
	Headquarters       string                             `json:"headquarters"`
	YearFounded        *int64                             `json:"yearFounded"`
	EmployeeGrowthRate string                             `json:"employeeGrowthRate"`
	SlackChannelId     string                             `json:"slackChannelId"`
	Relationship       neo4jenum.OrganizationRelationship `json:"relationship"`
	Stage              neo4jenum.OrganizationStage        `json:"stage"`
	LeadSource         string                             `json:"leadSource"`
}

type OrganizationUpdateFields struct {
	Name                     string                             `json:"name"`
	Hide                     bool                               `json:"hide"`
	Description              string                             `json:"description"`
	Website                  string                             `json:"website"`
	Industry                 string                             `json:"industry"`
	SubIndustry              string                             `json:"subIndustry"`
	IndustryGroup            string                             `json:"industryGroup"`
	TargetAudience           string                             `json:"targetAudience"`
	ValueProposition         string                             `json:"valueProposition"`
	IsPublic                 bool                               `json:"isPublic"`
	IsCustomer               bool                               `json:"isCustomer"`
	Employees                int64                              `json:"employees"`
	Market                   string                             `json:"market"`
	LastFundingRound         string                             `json:"lastFundingRound"`
	LastFundingAmount        string                             `json:"lastFundingAmount"`
	ReferenceId              string                             `json:"referenceId"`
	Note                     string                             `json:"note"`
	LogoUrl                  string                             `json:"logoUrl"`
	Headquarters             string                             `json:"headquarters"`
	YearFounded              *int64                             `json:"yearFounded"`
	EmployeeGrowthRate       string                             `json:"employeeGrowthRate"`
	SlackChannelId           string                             `json:"slackChannelId"`
	WebScrapedUrl            string                             `json:"webScrapedUrl"`
	EnrichDomain             string                             `json:"enrichDomain"`
	EnrichSource             string                             `json:"enrichSource"`
	Source                   string                             `json:"source"`
	UpdatedAt                time.Time                          `json:"updatedAt"`
	Relationship             neo4jenum.OrganizationRelationship `json:"relationship"`
	Stage                    neo4jenum.OrganizationStage        `json:"stage"`
	UpdateName               bool                               `json:"updateName"`
	UpdateDescription        bool                               `json:"updateDescription"`
	UpdateHide               bool                               `json:"updateHide"`
	UpdateIsCustomer         bool                               `json:"updateIsCustomer"`
	UpdateWebsite            bool                               `json:"updateWebsite"`
	UpdateIndustry           bool                               `json:"updateIndustry"`
	UpdateSubIndustry        bool                               `json:"updateSubIndustry"`
	UpdateIndustryGroup      bool                               `json:"updateIndustryGroup"`
	UpdateTargetAudience     bool                               `json:"updateTargetAudience"`
	UpdateValueProposition   bool                               `json:"updateValueProposition"`
	UpdateLastFundingRound   bool                               `json:"updateLastFundingRound"`
	UpdateLastFundingAmount  bool                               `json:"updateLastFundingAmount"`
	UpdateReferenceId        bool                               `json:"updateReferenceId"`
	UpdateNote               bool                               `json:"updateNote"`
	UpdateIsPublic           bool                               `json:"updateIsPublic"`
	UpdateEmployees          bool                               `json:"updateEmployees"`
	UpdateMarket             bool                               `json:"updateMarket"`
	UpdateYearFounded        bool                               `json:"updateYearFounded"`
	UpdateHeadquarters       bool                               `json:"updateHeadquarters"`
	UpdateLogoUrl            bool                               `json:"updateLogoUrl"`
	UpdateEmployeeGrowthRate bool                               `json:"updateEmployeeGrowthRate"`
	UpdateSlackChannelId     bool                               `json:"updateSlackChannelId"`
	UpdateRelationship       bool                               `json:"updateRelationship"`
	UpdateStage              bool                               `json:"updateStage"`
}

type OrganizationWriteRepository interface {
	CreateOrganization(ctx context.Context, tenant, organizationId string, data OrganizationCreateFields) error
	CreateOrganizationInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant, organizationId string, data OrganizationCreateFields) error
	UpdateOrganization(ctx context.Context, tenant, organizationId string, data OrganizationUpdateFields) error
	LinkWithDomain(ctx context.Context, tenant, organizationId, domain string) error
	UnlinkFromDomain(ctx context.Context, tenant, organizationId, domain string) error
	ReplaceOwner(ctx context.Context, tenant, organizationId, userId string) error
	SetVisibility(ctx context.Context, tenant, organizationId string, hide bool) error
	UpdateLastTouchpoint(ctx context.Context, tenant, organizationId string, touchpointAt time.Time, touchpointId, touchpointType string) error
	SetCustomerOsIdIfMissing(ctx context.Context, tenant, organizationId, customerOsId string) error
	LinkWithParentOrganization(ctx context.Context, tenant, organizationId, parentOrganizationId, subOrganizationType string) error
	UnlinkParentOrganization(ctx context.Context, tenant, organizationId, parentOrganizationId string) error
	UpdateArr(ctx context.Context, tenant, organizationId string) error
	UpdateRenewalSummary(ctx context.Context, tenant, organizationId string, likelihood *string, likelihoodOrder *int64, nextRenewalDate *time.Time) error
	WebScrapeRequested(ctx context.Context, tenant, organizationId, url string, attempt int64, requestedAt time.Time) error
	UpdateOnboardingStatus(ctx context.Context, tenant, organizationId, status, comments string, statusOrder *int64, updatedAt time.Time) error
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

func (r *organizationWriteRepository) CreateOrganization(ctx context.Context, tenant, organizationId string, data OrganizationCreateFields) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationWriteRepository.CreateOrganization")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, organizationId)

	session := utils.NewNeo4jWriteSession(ctx, *r.driver, utils.WithDatabaseName(r.database))
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		return nil, r.CreateOrganizationInTx(ctx, tx, tenant, organizationId, data)
	})
	return err
}

func (r *organizationWriteRepository) CreateOrganizationInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant, organizationId string, data OrganizationCreateFields) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationWriteRepository.CreateOrganizationInTx")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, organizationId)
	tracing.LogObjectAsJson(span, "data", data)

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant}) 
		 MERGE (t)<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization:Organization_%s {id:$id}) 
		 ON CREATE SET 	org.name = $name,
						org.description = $description,
						org.hide = $hide,
						org.website = $website,
						org.industry = $industry,
						org.subIndustry = $subIndustry,
						org.industryGroup = $industryGroup,
						org.targetAudience = $targetAudience,
						org.valueProposition = $valueProposition,
						org.lastFundingRound = $lastFundingRound,
						org.lastFundingAmount = $lastFundingAmount,
						org.referenceId = $referenceId,
						org.note = $note,
						org.isPublic = $isPublic,
						org.isCustomer = $isCustomer,
						org.source = $source,
						org.sourceOfTruth = $sourceOfTruth,
						org.employees = $employees,
						org.market = $market,
						org.logoUrl = $logoUrl,
						org.headquarters = $headquarters,
						org.yearFounded = $yearFounded,
						org.employeeGrowthRate = $employeeGrowthRate,
						org.appSource = $appSource,
						org.createdAt = $createdAt,
						org.updatedAt = $updatedAt,
						org.onboardingStatus = $onboardingStatus,
						org.relationship = $relationship,
						org.stage = $stage,
						org.stageUpdatedAt = $now,
						org.slackChannelId = $slackChannelId,
						org.syncedWithEventStore = true,
						org.leadSource = $leadSource
		 ON MATCH SET 	org.name = CASE WHEN org.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR org.name is null OR org.name = '' THEN $name ELSE org.name END,
						org.description = CASE WHEN org.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR org.description is null OR org.description = '' THEN $description ELSE org.description END,
						org.hide = CASE WHEN $overwrite=true OR (org.sourceOfTruth=$sourceOfTruth AND $hide = false) THEN $hide ELSE org.hide END,
						org.isCustomer = CASE WHEN $overwrite=true OR (org.sourceOfTruth=$sourceOfTruth AND $isCustomer = true) THEN $isCustomer ELSE org.isCustomer END,
						org.website = CASE WHEN org.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR org.website is null OR org.website = '' THEN $website ELSE org.website END,
						org.industry = CASE WHEN org.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR org.industry is null OR org.industry = '' THEN $industry ELSE org.industry END,
						org.subIndustry = CASE WHEN org.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR org.subIndustry is null OR org.subIndustry = '' THEN $subIndustry ELSE org.subIndustry END,
						org.industryGroup = CASE WHEN org.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR org.industryGroup is null OR org.industryGroup = '' THEN $industryGroup ELSE org.industryGroup END,
						org.targetAudience = CASE WHEN org.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR org.targetAudience is null OR org.targetAudience = '' THEN $targetAudience ELSE org.targetAudience END,
						org.valueProposition = CASE WHEN org.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR org.valueProposition is null OR org.valueProposition = '' THEN $valueProposition ELSE org.valueProposition END,
						org.lastFundingRound = CASE WHEN org.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR org.lastFundingRound is null OR org.lastFundingRound = '' THEN $lastFundingRound ELSE org.lastFundingRound END,
						org.lastFundingAmount = CASE WHEN org.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR org.lastFundingAmount is null OR org.lastFundingAmount = '' THEN $lastFundingAmount ELSE org.lastFundingAmount END,
						org.referenceId = CASE WHEN org.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR org.referenceId is null OR org.referenceId = '' THEN $referenceId ELSE org.referenceId END,
						org.logoUrl = CASE WHEN org.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR org.logoUrl is null OR org.logoUrl = '' THEN $logoUrl ELSE org.logoUrl END,
						org.headquarters = CASE WHEN org.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR org.headquarters is null OR org.headquarters = '' THEN $headquarters ELSE org.headquarters END,
						org.employeeGrowthRate = CASE WHEN org.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR org.employeeGrowthRate is null OR org.employeeGrowthRate = '' THEN $employeeGrowthRate ELSE org.employeeGrowthRate END,
						org.yearFounded = CASE WHEN org.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR org.yearFounded is null OR org.yearFounded = 0 THEN $yearFounded ELSE org.yearFounded END,
						org.note = CASE WHEN org.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR org.note is null OR org.note = '' THEN $note ELSE org.note END,
						org.isPublic = CASE WHEN org.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR org.isPublic is null THEN $isPublic ELSE org.isPublic END,
						org.employees = CASE WHEN org.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR org.employees is null THEN $employees ELSE org.employees END,
						org.market = CASE WHEN org.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR org.market is null OR org.market = '' THEN $market ELSE org.market END,
						org.slackChannelId = CASE WHEN org.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR org.slackChannelId is null OR org.slackChannelId = '' THEN $slackChannelId ELSE org.slackChannelId END,
						org.relationship = CASE WHEN org.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR org.relationship is null OR org.relationship = '' THEN $relationship ELSE org.relationship END,
						org.stage = CASE WHEN org.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR org.stage is null OR org.stage = '' THEN $stage ELSE org.stage END,
						org.stageUpdatedAt = CASE WHEN (org.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR org.stage is null OR org.stage = '') AND (org.stage is null OR org.stage <> $stage) THEN $now ELSE org.stageUpdatedAt END,
						org.updatedAt=$updatedAt,
						org.syncedWithEventStore = true`, tenant)
	params := map[string]any{
		"id":                 organizationId,
		"name":               data.Name,
		"hide":               data.Hide,
		"description":        data.Description,
		"website":            data.Website,
		"industry":           data.Industry,
		"subIndustry":        data.SubIndustry,
		"industryGroup":      data.IndustryGroup,
		"targetAudience":     data.TargetAudience,
		"valueProposition":   data.ValueProposition,
		"isPublic":           data.IsPublic,
		"isCustomer":         data.IsCustomer,
		"tenant":             tenant,
		"employees":          data.Employees,
		"market":             data.Market,
		"lastFundingRound":   data.LastFundingRound,
		"lastFundingAmount":  data.LastFundingAmount,
		"referenceId":        data.ReferenceId,
		"note":               data.Note,
		"logoUrl":            data.LogoUrl,
		"headquarters":       data.Headquarters,
		"yearFounded":        data.YearFounded,
		"employeeGrowthRate": data.EmployeeGrowthRate,
		"slackChannelId":     data.SlackChannelId,
		"source":             data.SourceFields.Source,
		"sourceOfTruth":      data.SourceFields.SourceOfTruth,
		"appSource":          data.SourceFields.AppSource,
		"createdAt":          data.CreatedAt,
		"updatedAt":          data.UpdatedAt,
		"onboardingStatus":   string(neo4jenum.OnboardingStatusNotApplicable),
		"overwrite":          data.SourceFields.Source == constants.SourceOpenline,
		"relationship":       data.Relationship.String(),
		"stage":              data.Stage.String(),
		"leadSource":         data.LeadSource,
		"now":                utils.Now(),
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}

func (r *organizationWriteRepository) UpdateOrganization(ctx context.Context, tenant, organizationId string, data OrganizationUpdateFields) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationWriteRepository.UpdateOrganization")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, organizationId)
	tracing.LogObjectAsJson(span, "data", data)

	params := map[string]any{
		"id":        organizationId,
		"tenant":    tenant,
		"source":    data.Source,
		"updatedAt": data.UpdatedAt,
		"overwrite": data.Source == constants.SourceOpenline || data.Source == constants.SourceWebscrape,
		"now":       utils.Now(),
	}
	cypher := `MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization {id:$id}) SET `
	if data.UpdateName {
		cypher += `org.name = CASE WHEN org.sourceOfTruth=$source OR $overwrite=true OR org.name = '' THEN $name ELSE org.name END,`
		params["name"] = data.Name
	}
	if data.UpdateDescription {
		cypher += `org.description = CASE WHEN org.sourceOfTruth=$source OR $overwrite=true OR org.description = '' THEN $description ELSE org.description END,`
		params["description"] = data.Description
	}
	if data.UpdateHide {
		cypher += `org.hide = CASE WHEN $overwrite=true OR $hide = false THEN $hide ELSE org.hide END,`
		params["hide"] = data.Hide
	}
	if data.UpdateIsCustomer {
		cypher += `org.isCustomer = CASE WHEN $overwrite=true OR (org.sourceOfTruth=$source AND $isCustomer = true) THEN $isCustomer ELSE org.isCustomer END,`
		params["isCustomer"] = data.IsCustomer
	}
	if data.UpdateWebsite {
		cypher += `org.website = CASE WHEN org.sourceOfTruth=$source OR $overwrite=true OR org.website is null OR org.website = '' THEN $website ELSE org.website END,`
		params["website"] = data.Website
	}
	if data.UpdateIndustry {
		cypher += `org.industry = CASE WHEN org.sourceOfTruth=$source OR $overwrite=true OR org.industry is null OR org.industry = '' THEN $industry ELSE org.industry END,`
		params["industry"] = data.Industry
	}
	if data.UpdateSubIndustry {
		cypher += `org.subIndustry = CASE WHEN org.sourceOfTruth=$source OR $overwrite=true OR org.subIndustry is null OR org.subIndustry = '' THEN $subIndustry ELSE org.subIndustry END,`
		params["subIndustry"] = data.SubIndustry
	}
	if data.UpdateIndustryGroup {
		cypher += `org.industryGroup = CASE WHEN org.sourceOfTruth=$source OR $overwrite=true OR org.industryGroup is null OR org.industryGroup = '' THEN $industryGroup ELSE org.industryGroup END,`
		params["industryGroup"] = data.IndustryGroup
	}
	if data.UpdateTargetAudience {
		cypher += `org.targetAudience = CASE WHEN org.sourceOfTruth=$source OR $overwrite=true OR org.targetAudience is null OR org.targetAudience = '' THEN $targetAudience ELSE org.targetAudience END,`
		params["targetAudience"] = data.TargetAudience
	}
	if data.UpdateValueProposition {
		cypher += `org.valueProposition = CASE WHEN org.sourceOfTruth=$source OR $overwrite=true OR org.valueProposition is null OR org.valueProposition = '' THEN $valueProposition ELSE org.valueProposition END,`
		params["valueProposition"] = data.ValueProposition
	}
	if data.UpdateLastFundingRound {
		cypher += `org.lastFundingRound = CASE WHEN org.sourceOfTruth=$source OR $overwrite=true OR org.lastFundingRound is null OR org.lastFundingRound = '' THEN $lastFundingRound ELSE org.lastFundingRound END,`
		params["lastFundingRound"] = data.LastFundingRound
	}
	if data.UpdateLastFundingAmount {
		cypher += `org.lastFundingAmount = CASE WHEN org.sourceOfTruth=$source OR $overwrite=true OR org.lastFundingAmount is null OR org.lastFundingAmount = '' THEN $lastFundingAmount ELSE org.lastFundingAmount END,`
		params["lastFundingAmount"] = data.LastFundingAmount
	}
	if data.UpdateReferenceId {
		cypher += `org.referenceId = CASE WHEN org.sourceOfTruth=$source OR $overwrite=true OR org.referenceId is null OR org.referenceId = '' THEN $referenceId ELSE org.referenceId END,`
		params["referenceId"] = data.ReferenceId
	}
	if data.UpdateNote {
		cypher += `org.note = CASE WHEN org.sourceOfTruth=$source OR $overwrite=true OR org.note is null OR org.note = '' THEN $note ELSE org.note END,`
		params["note"] = data.Note
	}
	if data.UpdateIsPublic {
		cypher += `org.isPublic = CASE WHEN org.sourceOfTruth=$source OR $overwrite=true OR org.isPublic is null THEN $isPublic ELSE org.isPublic END,`
		params["isPublic"] = data.IsPublic
	}
	if data.UpdateEmployees {
		cypher += `org.employees = CASE WHEN org.sourceOfTruth=$source OR $overwrite=true OR org.employees is null THEN $employees ELSE org.employees END,`
		params["employees"] = data.Employees
	}
	if data.UpdateMarket {
		cypher += `org.market = CASE WHEN org.sourceOfTruth=$source OR $overwrite=true OR org.market is null OR org.market = '' THEN $market ELSE org.market END,`
		params["market"] = data.Market
	}
	if data.UpdateYearFounded {
		cypher += `org.yearFounded = CASE WHEN org.sourceOfTruth=$source OR $overwrite=true OR org.yearFounded is null OR org.yearFounded = 0 THEN $yearFounded ELSE org.yearFounded END,`
		params["yearFounded"] = data.YearFounded
	}
	if data.UpdateHeadquarters {
		cypher += `org.headquarters = CASE WHEN org.sourceOfTruth=$source OR $overwrite=true OR org.headquarters is null OR org.headquarters = '' THEN $headquarters ELSE org.headquarters END,`
		params["headquarters"] = data.Headquarters
	}
	if data.UpdateLogoUrl {
		cypher += `org.logoUrl = CASE WHEN org.sourceOfTruth=$source OR $overwrite=true OR org.logoUrl is null OR org.logoUrl = '' THEN $logoUrl ELSE org.logoUrl END,`
		params["logoUrl"] = data.LogoUrl
	}
	if data.UpdateEmployeeGrowthRate {
		cypher += `org.employeeGrowthRate = CASE WHEN org.sourceOfTruth=$source OR $overwrite=true OR org.employeeGrowthRate is null OR org.employeeGrowthRate = '' THEN $employeeGrowthRate ELSE org.employeeGrowthRate END,`
		params["employeeGrowthRate"] = data.EmployeeGrowthRate
	}
	if data.UpdateSlackChannelId {
		cypher += `org.slackChannelId = CASE WHEN org.sourceOfTruth=$source OR $overwrite=true OR org.slackChannelId is null OR org.slackChannelId = '' THEN $slackChannelId ELSE org.slackChannelId END,`
		params["slackChannelId"] = data.SlackChannelId
	}
	if data.UpdateRelationship {
		cypher += `org.relationship = CASE WHEN org.sourceOfTruth=$source OR $overwrite=true OR org.relationship is null OR org.relationship = '' THEN $relationship ELSE org.relationship END,`
		params["relationship"] = data.Relationship.String()
	}
	if data.UpdateStage {
		cypher += `org.stage = CASE WHEN org.sourceOfTruth=$source OR $overwrite=true OR org.stage is null OR org.stage = '' THEN $stage ELSE org.stage END,`
		cypher += `org.stageUpdatedAt = CASE WHEN (org.sourceOfTruth=$source OR $overwrite=true OR org.stage is null OR org.stage = '') AND (org.stage is null OR org.stage <> $stage) THEN $now ELSE org.stageUpdatedAt END,`
		params["stage"] = data.Stage.String()
	}
	if data.WebScrapedUrl != "" {
		params["webScrapedUrl"] = data.WebScrapedUrl
		params["webScrapedAt"] = utils.Now()
		cypher += `org.webScrapedUrl = $webScrapedUrl, org.webScrapedAt = $webScrapedAt,`
	}
	if data.EnrichDomain != "" && data.EnrichSource != "" {
		params["enrichDomain"] = data.EnrichDomain
		params["enrichSource"] = data.EnrichSource
		params["enrichedAt"] = utils.Now()
		cypher += `org.enrichDomain = $enrichDomain, org.enrichSource = $enrichSource, org.enrichedAt = $enrichedAt,`
	}
	cypher += ` org.sourceOfTruth = case WHEN $overwrite=true THEN $source ELSE org.sourceOfTruth END,
				org.updatedAt = $updatedAt,
				org.syncedWithEventStore = true`

	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}

func (r *organizationWriteRepository) LinkWithDomain(ctx context.Context, tenant, organizationId, domain string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationWriteRepository.MergeOrganizationDomain")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, organizationId)
	span.SetTag(tracing.SpanTagEntityId, organizationId)

	cypher := `MERGE (d:Domain {domain:$domain}) 
				ON CREATE SET 	d.id=randomUUID(), 
								d.createdAt=$now, 
								d.updatedAt=$now,
								d.appSource=$appSource
				WITH d
				MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization {id:$organizationId})
				MERGE (org)-[rel:HAS_DOMAIN]->(d)
				SET rel.syncedWithEventStore = true`
	params := map[string]any{
		"tenant":         tenant,
		"organizationId": organizationId,
		"domain":         strings.ToLower(domain),
		"appSource":      constants.AppSourceEventProcessingPlatform,
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

func (r *organizationWriteRepository) UnlinkFromDomain(ctx context.Context, tenant, organizationId, domain string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationWriteRepository.UnlinkFromDomain")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, organizationId)

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization {id:$organizationId})
		 MATCH (org)-[rel:HAS_DOMAIN]->(d:Domain {domain:$domain})
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

func (r *organizationWriteRepository) ReplaceOwner(ctx context.Context, tenant, organizationId, userId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationWriteRepository.ReplaceOwner")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.LogFields(log.String("organizationId", organizationId), log.String("userId", userId))

	query := `MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization {id:$organizationId})
			OPTIONAL MATCH (:User)-[rel:OWNS]->(org)
			DELETE rel
			WITH org, t
			MATCH (t)<-[:USER_BELONGS_TO_TENANT]-(u:User {id:$userId})
			WHERE (u.internal=false OR u.internal is null) AND (u.bot=false OR u.bot is null)
			MERGE (u)-[:OWNS]->(org)
			SET org.updatedAt=$now, org.sourceOfTruth=$source`

	session := utils.NewNeo4jWriteSession(ctx, *r.driver, utils.WithDatabaseName(r.database))
	defer session.Close(ctx)

	return utils.ExecuteWriteQuery(ctx, *r.driver, query, map[string]any{
		"tenant":         tenant,
		"organizationId": organizationId,
		"userId":         userId,
		"source":         constants.SourceOpenline,
		"now":            utils.Now(),
	})
}

func (r *organizationWriteRepository) SetVisibility(ctx context.Context, tenant, organizationId string, hide bool) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationWriteRepository.SetVisibility")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, organizationId)
	span.LogFields(log.Bool("hide", hide))

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization {id:$id})
			WHERE org:Organization_%s
		 SET	org.hide = $hide,
				org.updatedAt = $now`, tenant)
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

func (r *organizationWriteRepository) UpdateLastTouchpoint(ctx context.Context, tenant, organizationId string, touchpointAt time.Time, touchpointId, touchpointType string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationWriteRepository.UpdateLastTouchpoint")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.LogFields(log.String("organizationId", organizationId), log.String("touchpointId", touchpointId), log.Object("touchpointAt", touchpointAt))

	cypher := `MATCH (:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization {id:$organizationId})
		 SET org.lastTouchpointAt=$touchpointAt, org.lastTouchpointId=$touchpointId, org.lastTouchpointType=$touchpointType`
	params := map[string]any{
		"tenant":         tenant,
		"organizationId": organizationId,
		"touchpointAt":   touchpointAt,
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
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, organizationId)
	span.LogFields(log.String("customerOsId", customerOsId))

	cypher := `MATCH (:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization {id:$organizationId})
		 SET org.customerOsId = CASE WHEN (org.customerOsId IS NULL OR org.customerOsId = '') AND $customerOsId <> '' THEN $customerOsId ELSE org.customerOsId END`
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
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, organizationId)
	span.LogFields(log.String("parentOrganizationId", parentOrganizationId), log.String("subOrganizationType", subOrganizationType))

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(parent:Organization {id:$parentOrganizationId}),
		 			(t)<-[:ORGANIZATION_BELONGS_TO_TENANT]-(sub:Organization {id:$subOrganizationId}) 
		 	MERGE (sub)-[rel:SUBSIDIARY_OF]->(parent) 
		 		ON CREATE SET rel.type=$type 
		 		ON MATCH SET rel.type=$type`
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
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, organizationId)
	span.LogFields(log.String("parentOrganizationId", parentOrganizationId))

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(parent:Organization {id:$parentOrganizationId})<-[rel:SUBSIDIARY_OF]-(sub:Organization {id:$subOrganizationId})-[:ORGANIZATION_BELONGS_TO_TENANT]->(t)
		 		DELETE rel`
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
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, organizationId)

	cypher := `MATCH (t:Tenant {name: $tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization {id: $organizationId})
				OPTIONAL MATCH (org)-[:HAS_CONTRACT]->(c:Contract)
				OPTIONAL MATCH (c)-[:ACTIVE_RENEWAL]->(op:Opportunity)
				WITH org, COALESCE(sum(op.amount), 0) as arr, COALESCE(sum(op.maxAmount), 0) as maxArr
				SET org.renewalForecastArr = arr, org.renewalForecastMaxArr = maxArr, org.updatedAt = $now`
	params := map[string]any{
		"tenant":         tenant,
		"organizationId": organizationId,
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

func (r *organizationWriteRepository) UpdateRenewalSummary(ctx context.Context, tenant, organizationId string, likelihood *string, likelihoodOrder *int64, nextRenewalDate *time.Time) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationWriteRepository.UpdateRenewalSummary")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, organizationId)
	span.LogFields(log.Object("likelihood", likelihood), log.Object("likelihoodOrder", likelihoodOrder), log.Object("nextRenewalDate", nextRenewalDate))

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization {id:$organizationId})
				SET org.derivedRenewalLikelihood = $derivedRenewalLikelihood,
					org.derivedRenewalLikelihoodOrder = $derivedRenewalLikelihoodOrder,
					org.derivedNextRenewalAt = $derivedNextRenewalAt,
					org.updatedAt = $now`
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
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
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
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, organizationId)

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization {id:$organizationId})
				SET org.onboardingUpdatedAt = CASE WHEN org.onboardingStatus IS NULL OR org.onboardingStatus <> $status THEN $updatedAt ELSE org.onboardingUpdatedAt END,
					org.onboardingStatus=$status,
					org.onboardingStatusOrder=$statusOrder,
					org.onboardingComments=$comments,
					org.onboardingUpdatedAt=$updatedAt,
					org.updatedAt=$now`
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
