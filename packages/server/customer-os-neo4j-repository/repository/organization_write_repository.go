package repository

import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/constants"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/model"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"golang.org/x/net/context"
	"strings"
	"time"
)

// DEPRECATED use save
type OrganizationCreateFields struct {
	AggregateVersion   int64                              `json:"aggregateVersion"`
	SourceFields       model.SourceFields                 `json:"sourceFields"`
	CreatedAt          time.Time                          `json:"createdAt"`
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
	Employees          int64                              `json:"employees"`
	Market             string                             `json:"market"`
	LastFundingRound   string                             `json:"lastFundingRound"`
	LastFundingAmount  string                             `json:"lastFundingAmount"`
	ReferenceId        string                             `json:"referenceId"`
	Note               string                             `json:"note"`
	LogoUrl            string                             `json:"logoUrl"`
	IconUrl            string                             `json:"iconUrl"`
	Headquarters       string                             `json:"headquarters"`
	YearFounded        *int64                             `json:"yearFounded"`
	EmployeeGrowthRate string                             `json:"employeeGrowthRate"`
	SlackChannelId     string                             `json:"slackChannelId"`
	Relationship       neo4jenum.OrganizationRelationship `json:"relationship"`
	Stage              neo4jenum.OrganizationStage        `json:"stage"`
	LeadSource         string                             `json:"leadSource"`
	IcpFit             bool                               `json:"icpFit"`
}

// DEPRECATED use save
type OrganizationUpdateFields struct {
	AggregateVersion         int64                              `json:"aggregateVersion"`
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
	Employees                int64                              `json:"employees"`
	Market                   string                             `json:"market"`
	LastFundingRound         string                             `json:"lastFundingRound"`
	LastFundingAmount        string                             `json:"lastFundingAmount"`
	ReferenceId              string                             `json:"referenceId"`
	Note                     string                             `json:"note"`
	LogoUrl                  string                             `json:"logoUrl"`
	IconUrl                  string                             `json:"iconUrl"`
	Headquarters             string                             `json:"headquarters"`
	YearFounded              *int64                             `json:"yearFounded"`
	EmployeeGrowthRate       string                             `json:"employeeGrowthRate"`
	SlackChannelId           string                             `json:"slackChannelId"`
	EnrichDomain             string                             `json:"enrichDomain"`
	EnrichSource             string                             `json:"enrichSource"`
	Source                   string                             `json:"source"`
	Relationship             neo4jenum.OrganizationRelationship `json:"relationship"`
	Stage                    neo4jenum.OrganizationStage        `json:"stage"`
	IcpFit                   bool                               `json:"icpFit"`
	UpdateName               bool                               `json:"updateName"`
	UpdateDescription        bool                               `json:"updateDescription"`
	UpdateHide               bool                               `json:"updateHide"`
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
	UpdateIconUrl            bool                               `json:"updateIconUrl"`
	UpdateEmployeeGrowthRate bool                               `json:"updateEmployeeGrowthRate"`
	UpdateSlackChannelId     bool                               `json:"updateSlackChannelId"`
	UpdateRelationship       bool                               `json:"updateRelationship"`
	UpdateStage              bool                               `json:"updateStage"`
	UpdateIcpFit             bool                               `json:"updateIcpFit"`
}

type OrganizationSaveFields struct {
	SourceFields model.SourceFields `json:"sourceFields"`

	//not stored directly in node
	Domains        []string             `json:"domains"`
	ExternalSystem model.ExternalSystem `json:"externalSystem"`

	OwnerId       string `json:"ownerId"`
	UpdateOwnerId bool   `json:"updateOwnerId"`

	Hide               bool                               `json:"hide"`
	Name               string                             `json:"name"`
	Description        string                             `json:"description"`
	Website            string                             `json:"website"`
	Industry           string                             `json:"industry"`
	SubIndustry        string                             `json:"subIndustry"`
	IndustryGroup      string                             `json:"industryGroup"`
	TargetAudience     string                             `json:"targetAudience"`
	ValueProposition   string                             `json:"valueProposition"`
	IsPublic           bool                               `json:"isPublic"`
	Employees          int64                              `json:"employees"`
	Market             string                             `json:"market"`
	LastFundingRound   string                             `json:"lastFundingRound"`
	LastFundingAmount  string                             `json:"lastFundingAmount"`
	CustomerOsId       string                             `json:"customerOsId"`
	ReferenceId        string                             `json:"referenceId"`
	Note               string                             `json:"note"`
	LogoUrl            string                             `json:"logoUrl"`
	IconUrl            string                             `json:"iconUrl"`
	Headquarters       string                             `json:"headquarters"`
	YearFounded        int64                              `json:"yearFounded"`
	EmployeeGrowthRate string                             `json:"employeeGrowthRate"`
	SlackChannelId     string                             `json:"slackChannelId"`
	EnrichDomain       string                             `json:"enrichDomain"`
	EnrichSource       string                             `json:"enrichSource"`
	LeadSource         string                             `json:"leadSource"`
	Relationship       neo4jenum.OrganizationRelationship `json:"relationship"`
	Stage              neo4jenum.OrganizationStage        `json:"stage"`
	IcpFit             bool                               `json:"icpFit"`

	UpdateHide               bool `json:"updateHide"`
	UpdateName               bool `json:"updateName"`
	UpdateDescription        bool `json:"updateDescription"`
	UpdateWebsite            bool `json:"updateWebsite"`
	UpdateIndustry           bool `json:"updateIndustry"`
	UpdateSubIndustry        bool `json:"updateSubIndustry"`
	UpdateIndustryGroup      bool `json:"updateIndustryGroup"`
	UpdateTargetAudience     bool `json:"updateTargetAudience"`
	UpdateValueProposition   bool `json:"updateValueProposition"`
	UpdateLastFundingRound   bool `json:"updateLastFundingRound"`
	UpdateLastFundingAmount  bool `json:"updateLastFundingAmount"`
	UpdateCustomerOsId       bool `json:"updateCustomerOsId"`
	UpdateReferenceId        bool `json:"updateReferenceId"`
	UpdateNote               bool `json:"updateNote"`
	UpdateIsPublic           bool `json:"updateIsPublic"`
	UpdateEmployees          bool `json:"updateEmployees"`
	UpdateMarket             bool `json:"updateMarket"`
	UpdateYearFounded        bool `json:"updateYearFounded"`
	UpdateHeadquarters       bool `json:"updateHeadquarters"`
	UpdateLogoUrl            bool `json:"updateLogoUrl"`
	UpdateIconUrl            bool `json:"updateIconUrl"`
	UpdateEmployeeGrowthRate bool `json:"updateEmployeeGrowthRate"`
	UpdateSlackChannelId     bool `json:"updateSlackChannelId"`
	UpdateLeadSource         bool `json:"updateLeadSource"`
	UpdateRelationship       bool `json:"updateRelationship"`
	UpdateStage              bool `json:"updateStage"`
	UpdateIcpFit             bool `json:"updateIcpFit"`
}

type OrganizationWriteRepository interface {
	Save(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, organizationId string, data OrganizationSaveFields) error

	//Deprecated
	ReserveOrganizationId(ctx context.Context, tenant, organizationId string) (string, error)
	//Deprecated
	CreateOrganization(ctx context.Context, tenant, organizationId string, data OrganizationCreateFields) error
	//Deprecated
	CreateOrganizationInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant, organizationId string, data OrganizationCreateFields) error
	//Deprecated
	UpdateOrganization(ctx context.Context, tenant, organizationId string, data OrganizationUpdateFields) error
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

func (r *organizationWriteRepository) ReserveOrganizationId(ctx context.Context, tenant, inputId string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationWriteRepository.ReserveOrganizationId")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	span.LogFields(log.String("organizationId", inputId))

	orgId := utils.NewUUIDIfEmpty(inputId)

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant}) 
							MERGE (t)<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization:Organization_%s {id:$id})
							SET org.updatedAt = datetime()`, tenant)
	params := map[string]any{
		"id":     orgId,
		"tenant": tenant,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return orgId, err
}

func (r *organizationWriteRepository) CreateOrganization(ctx context.Context, tenant, organizationId string, data OrganizationCreateFields) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationWriteRepository.CreateOrganization")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
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
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
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
						org.source = $source,
						org.sourceOfTruth = $sourceOfTruth,
						org.employees = $employees,
						org.market = $market,
						org.logoUrl = $logoUrl,
						org.iconUrl = $iconUrl,
						org.headquarters = $headquarters,
						org.yearFounded = $yearFounded,
						org.employeeGrowthRate = $employeeGrowthRate,
						org.appSource = $appSource,
						org.createdAt = $createdAt,
						org.updatedAt = datetime(),
						org.onboardingStatus = $onboardingStatus,
						org.relationship = $relationship,
						org.stage = $stage,
						org.stageUpdatedAt = datetime(),
						org.slackChannelId = $slackChannelId,
						org.leadSource = $leadSource,
						org.icpFit = $icpFit,
						org.aggregateVersion = $aggregateVersion,
						org.lastTouchpointAt = datetime()
		 ON MATCH SET 	org.source = CASE WHEN org.source IS NULL THEN $source ELSE org.source END,
						org.sourceOfTruth = CASE WHEN org.sourceOfTruth IS NULL THEN $sourceOfTruth ELSE org.sourceOfTruth END,
						org.appSource = CASE WHEN org.appSource IS NULL THEN $appSource ELSE org.appSource END,
						org.createdAt = CASE WHEN org.createdAt IS NULL THEN $createdAt ELSE org.createdAt END,
						org.leadSource = CASE WHEN org.leadSource IS NULL THEN $leadSource ELSE org.leadSource END,
						org.lastTouchpointAt = CASE WHEN org.lastTouchpointAt IS NULL THEN datetime() ELSE org.lastTouchpointAt END,
						org.onboardingStatus = CASE WHEN org.onboardingStatus IS NULL THEN $onboardingStatus ELSE org.onboardingStatus END,
						org.name = CASE WHEN org.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR org.name IS NULL OR org.name = '' THEN $name ELSE org.name END,
						org.description = CASE WHEN org.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR org.description IS NULL OR org.description = '' THEN $description ELSE org.description END,
						org.hide = CASE WHEN $overwrite=true OR org.hide IS NULL OR (org.sourceOfTruth=$sourceOfTruth AND $hide = false) THEN $hide ELSE org.hide END,
						org.hiddenAt = CASE WHEN $hide = true THEN datetime() ELSE org.hiddenAt END,
						org.website = CASE WHEN org.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR org.website IS NULL OR org.website = '' THEN $website ELSE org.website END,
						org.industry = CASE WHEN org.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR org.industry IS NULL OR org.industry = '' THEN $industry ELSE org.industry END,
						org.subIndustry = CASE WHEN org.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR org.subIndustry IS NULL OR org.subIndustry = '' THEN $subIndustry ELSE org.subIndustry END,
						org.industryGroup = CASE WHEN org.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR org.industryGroup IS NULL OR org.industryGroup = '' THEN $industryGroup ELSE org.industryGroup END,
						org.targetAudience = CASE WHEN org.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR org.targetAudience IS NULL OR org.targetAudience = '' THEN $targetAudience ELSE org.targetAudience END,
						org.valueProposition = CASE WHEN org.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR org.valueProposition IS NULL OR org.valueProposition = '' THEN $valueProposition ELSE org.valueProposition END,
						org.lastFundingRound = CASE WHEN org.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR org.lastFundingRound IS NULL OR org.lastFundingRound = '' THEN $lastFundingRound ELSE org.lastFundingRound END,
						org.lastFundingAmount = CASE WHEN org.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR org.lastFundingAmount IS NULL OR org.lastFundingAmount = '' THEN $lastFundingAmount ELSE org.lastFundingAmount END,
						org.referenceId = CASE WHEN org.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR org.referenceId IS NULL OR org.referenceId = '' THEN $referenceId ELSE org.referenceId END,
						org.logoUrl = CASE WHEN org.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR org.logoUrl is null OR org.logoUrl = '' THEN $logoUrl ELSE org.logoUrl END,
						org.iconUrl = CASE WHEN org.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR org.iconUrl is null OR org.iconUrl = '' THEN $iconUrl ELSE org.iconUrl END,
						org.headquarters = CASE WHEN org.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR org.headquarters is null OR org.headquarters = '' THEN $headquarters ELSE org.headquarters END,
						org.employeeGrowthRate = CASE WHEN org.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR org.employeeGrowthRate is null OR org.employeeGrowthRate = '' THEN $employeeGrowthRate ELSE org.employeeGrowthRate END,
						org.yearFounded = CASE WHEN org.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR org.yearFounded is null OR org.yearFounded = 0 THEN $yearFounded ELSE org.yearFounded END,
						org.note = CASE WHEN org.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR org.note IS NULL OR org.note = '' THEN $note ELSE org.note END,
						org.isPublic = CASE WHEN org.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR org.isPublic is null THEN $isPublic ELSE org.isPublic END,
						org.employees = CASE WHEN org.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR org.employees is null THEN $employees ELSE org.employees END,
						org.market = CASE WHEN org.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR org.market IS NULL OR org.market = '' THEN $market ELSE org.market END,
						org.slackChannelId = CASE WHEN org.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR org.slackChannelId is null OR org.slackChannelId = '' THEN $slackChannelId ELSE org.slackChannelId END,
						org.relationship = CASE WHEN org.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR org.relationship is null OR org.relationship = '' THEN $relationship ELSE org.relationship END,
						org.stage = CASE WHEN org.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR org.stage is null OR org.stage = '' THEN $stage ELSE org.stage END,
						org.stageUpdatedAt = CASE WHEN (org.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR org.stage is null OR org.stage = '') AND (org.stage is null OR org.stage <> $stage) THEN datetime() ELSE org.stageUpdatedAt END,
						org.icpFit = CASE WHEN org.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR org.icpFit IS NULL THEN $icpFit ELSE org.icpFit END,
						org.aggregateVersion = $aggregateVersion,
						org.updatedAt=datetime()`, tenant)
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
		"tenant":             tenant,
		"employees":          data.Employees,
		"market":             data.Market,
		"lastFundingRound":   data.LastFundingRound,
		"lastFundingAmount":  data.LastFundingAmount,
		"referenceId":        data.ReferenceId,
		"note":               data.Note,
		"logoUrl":            data.LogoUrl,
		"iconUrl":            data.IconUrl,
		"headquarters":       data.Headquarters,
		"yearFounded":        data.YearFounded,
		"employeeGrowthRate": data.EmployeeGrowthRate,
		"slackChannelId":     data.SlackChannelId,
		"source":             data.SourceFields.Source,
		"sourceOfTruth":      data.SourceFields.SourceOfTruth,
		"appSource":          data.SourceFields.AppSource,
		"createdAt":          data.CreatedAt,
		"onboardingStatus":   string(neo4jenum.OnboardingStatusNotApplicable),
		"overwrite":          data.SourceFields.Source == constants.SourceOpenline,
		"relationship":       data.Relationship.String(),
		"stage":              data.Stage.String(),
		"icpFit":             data.IcpFit,
		"leadSource":         data.LeadSource,
		"aggregateVersion":   data.AggregateVersion,
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
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, organizationId)
	tracing.LogObjectAsJson(span, "data", data)

	params := map[string]any{
		"id":               organizationId,
		"tenant":           tenant,
		"source":           data.Source,
		"overwrite":        data.Source == constants.SourceOpenline || data.Source == constants.SourceWebscrape,
		"now":              utils.Now(),
		"aggregateVersion": data.AggregateVersion,
	}
	cypher := `MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization {id:$id})
				WHERE org.aggregateVersion IS NULL OR org.aggregateVersion < $aggregateVersion
				SET `
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
		cypher += `org.hiddenAt = CASE WHEN $hide = true THEN datetime() ELSE org.hiddenAt END,`
		params["hide"] = data.Hide
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
	if data.UpdateIconUrl {
		cypher += `org.iconUrl = CASE WHEN org.sourceOfTruth=$source OR $overwrite=true OR org.iconUrl is null OR org.iconUrl = '' THEN $iconUrl ELSE org.iconUrl END,`
		params["iconUrl"] = data.IconUrl
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
	if data.UpdateIcpFit {
		cypher += `org.icpFit = CASE WHEN org.sourceOfTruth=$source OR $overwrite=true THEN $icpFit ELSE org.icpFit END,`
		params["icpFit"] = data.IcpFit
	}
	if data.EnrichDomain != "" && data.EnrichSource != "" {
		params["enrichDomain"] = data.EnrichDomain
		params["enrichSource"] = data.EnrichSource
		params["enrichedAt"] = utils.Now()
		cypher += `org.enrichDomain = $enrichDomain, org.enrichSource = $enrichSource, org.enrichedAt = $enrichedAt,`
	}
	cypher += ` org.sourceOfTruth = case WHEN $overwrite=true THEN $source ELSE org.sourceOfTruth END,
				org.updatedAt = datetime(),
				org.aggregateVersion = $aggregateVersion`

	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}

func (r *organizationWriteRepository) Save(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, organizationId string, data OrganizationSaveFields) error {
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
					org.createdAt = datetime(),
					org.updatedAt = datetime(),
					org.hide=false`, tenant)
		paramsCreate := map[string]any{
			"tenant":         tenant,
			"organizationId": organizationId,
			"source":         utils.FirstNotEmptyString(data.SourceFields.Source, neo4jentity.DataSourceOpenline.String()),
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

		if data.UpdateName {
			cypherUpdate += `org.name = $name,`
			paramsUpdate["name"] = data.Name
		}
		if data.UpdateDescription {
			cypherUpdate += `org.description = $description,`
			paramsUpdate["description"] = data.Description
		}
		if data.UpdateHide {
			cypherUpdate += `org.hide = $hide,`
			cypherUpdate += `org.hiddenAt = CASE WHEN $hide = true THEN datetime() ELSE null END,`
			paramsUpdate["hide"] = data.Hide
		}
		if data.UpdateWebsite {
			cypherUpdate += `org.website = $website,`
			paramsUpdate["website"] = data.Website
		}
		if data.UpdateIndustry {
			cypherUpdate += `org.industry = $industry,`
			paramsUpdate["industry"] = data.Industry
		}
		if data.UpdateSubIndustry {
			cypherUpdate += `org.subIndustry = $subIndustry,`
			paramsUpdate["subIndustry"] = data.SubIndustry
		}
		if data.UpdateIndustryGroup {
			cypherUpdate += `org.industryGroup = $industryGroup,`
			paramsUpdate["industryGroup"] = data.IndustryGroup
		}
		if data.UpdateTargetAudience {
			cypherUpdate += `org.targetAudience = $targetAudience,`
			paramsUpdate["targetAudience"] = data.TargetAudience
		}
		if data.UpdateValueProposition {
			cypherUpdate += `org.valueProposition = $valueProposition,`
			paramsUpdate["valueProposition"] = data.ValueProposition
		}
		if data.UpdateLastFundingRound {
			cypherUpdate += `org.lastFundingRound = $lastFundingRound,`
			paramsUpdate["lastFundingRound"] = data.LastFundingRound
		}
		if data.UpdateLastFundingAmount {
			cypherUpdate += `org.lastFundingAmount = $lastFundingAmount,`
			paramsUpdate["lastFundingAmount"] = data.LastFundingAmount
		}
		if data.UpdateCustomerOsId {
			cypherUpdate += `org.customerOsId = $customerOsId,`
			paramsUpdate["customerOsId"] = data.CustomerOsId
		}
		if data.UpdateReferenceId {
			cypherUpdate += `org.referenceId = $referenceId,`
			paramsUpdate["referenceId"] = data.ReferenceId
		}
		if data.UpdateNote {
			cypherUpdate += `org.note = $note,`
			paramsUpdate["note"] = data.Note
		}
		if data.UpdateIsPublic {
			cypherUpdate += `org.isPublic = $isPublic,`
			paramsUpdate["isPublic"] = data.IsPublic
		}
		if data.UpdateEmployees {
			cypherUpdate += `org.employees = $employees,`
			paramsUpdate["employees"] = data.Employees
		}
		if data.UpdateMarket {
			cypherUpdate += `org.market = $market,`
			paramsUpdate["market"] = data.Market
		}
		if data.UpdateYearFounded {
			cypherUpdate += `org.yearFounded = $yearFounded,`
			paramsUpdate["yearFounded"] = data.YearFounded
		}
		if data.UpdateHeadquarters {
			cypherUpdate += `org.headquarters = $headquarters,`
			paramsUpdate["headquarters"] = data.Headquarters
		}
		if data.UpdateLogoUrl {
			cypherUpdate += `org.logoUrl = $logoUrl,`
			paramsUpdate["logoUrl"] = data.LogoUrl
		}
		if data.UpdateIconUrl {
			cypherUpdate += `org.iconUrl = $iconUrl,`
			paramsUpdate["iconUrl"] = data.IconUrl
		}
		if data.UpdateEmployeeGrowthRate {
			cypherUpdate += `org.employeeGrowthRate = $employeeGrowthRate,`
			paramsUpdate["employeeGrowthRate"] = data.EmployeeGrowthRate
		}
		if data.UpdateSlackChannelId {
			cypherUpdate += `org.slackChannelId = $slackChannelId,`
			paramsUpdate["slackChannelId"] = data.SlackChannelId
		}
		if data.UpdateRelationship {
			cypherUpdate += `org.relationship = $relationship,`
			paramsUpdate["relationship"] = data.Relationship.String()
		}
		if data.UpdateStage {
			cypherUpdate += `org.stage = $stage,`
			cypherUpdate += `org.stageUpdatedAt = CASE WHEN (org.stage is null OR org.stage = '') AND (org.stage is null OR org.stage <> $stage) THEN $now ELSE org.stageUpdatedAt END,`
			paramsUpdate["stage"] = data.Stage.String()
		}
		if data.UpdateLeadSource {
			cypherUpdate += `org.leadSource = $leadSource,`
			paramsUpdate["leadSource"] = data.LeadSource
		}
		if data.UpdateIcpFit {
			cypherUpdate += `org.icpFit = $icpFit,`
			paramsUpdate["icpFit"] = data.IcpFit
		}
		if data.EnrichDomain != "" && data.EnrichSource != "" {
			cypherUpdate += `org.enrichDomain = $enrichDomain, org.enrichSource = $enrichSource, org.enrichedAt = $enrichedAt,`
			paramsUpdate["enrichDomain"] = data.EnrichDomain
			paramsUpdate["enrichSource"] = data.EnrichSource
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
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationWriteRepository.MergeOrganizationDomain")
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
				WITH d, org, COUNT(otherOrg) AS existingOrgCount
				WHERE existingOrgCount = 0
				MERGE (org)-[rel:HAS_DOMAIN]->(d)
				SET org.updatedAt = datetime()
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
			WHERE (u.internal=false OR u.internal is null) AND (u.bot=false OR u.bot is null)
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
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
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
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
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
			SET org.archived=true, org.archivedAt=$now, org:ArchivedOrganization_%s
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
