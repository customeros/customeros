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
	"time"
)

type OrganizationCreateFields struct {
	SourceFields       model.Source `json:"sourceFields"`
	CreatedAt          time.Time    `json:"createdAt"`
	UpdatedAt          time.Time    `json:"updatedAt"`
	Name               string       `json:"name"`
	Hide               bool         `json:"hide"`
	Description        string       `json:"description"`
	Website            string       `json:"website"`
	Industry           string       `json:"industry"`
	SubIndustry        string       `json:"subIndustry"`
	IndustryGroup      string       `json:"industryGroup"`
	TargetAudience     string       `json:"targetAudience"`
	ValueProposition   string       `json:"valueProposition"`
	IsPublic           bool         `json:"isPublic"`
	IsCustomer         bool         `json:"isCustomer"`
	Employees          int64        `json:"employees"`
	Market             string       `json:"market"`
	LastFundingRound   string       `json:"lastFundingRound"`
	LastFundingAmount  string       `json:"lastFundingAmount"`
	ReferenceId        string       `json:"referenceId"`
	Note               string       `json:"note"`
	LogoUrl            string       `json:"logoUrl"`
	Headquarters       string       `json:"headquarters"`
	YearFounded        *int64       `json:"yearFounded"`
	EmployeeGrowthRate string       `json:"employeeGrowthRate"`
}

type OrganizationUpdateFields struct {
	Name                     string    `json:"name"`
	Hide                     bool      `json:"hide"`
	Description              string    `json:"description"`
	Website                  string    `json:"website"`
	Industry                 string    `json:"industry"`
	SubIndustry              string    `json:"subIndustry"`
	IndustryGroup            string    `json:"industryGroup"`
	TargetAudience           string    `json:"targetAudience"`
	ValueProposition         string    `json:"valueProposition"`
	IsPublic                 bool      `json:"isPublic"`
	IsCustomer               bool      `json:"isCustomer"`
	Employees                int64     `json:"employees"`
	Market                   string    `json:"market"`
	LastFundingRound         string    `json:"lastFundingRound"`
	LastFundingAmount        string    `json:"lastFundingAmount"`
	ReferenceId              string    `json:"referenceId"`
	Note                     string    `json:"note"`
	LogoUrl                  string    `json:"logoUrl"`
	Headquarters             string    `json:"headquarters"`
	YearFounded              *int64    `json:"yearFounded"`
	EmployeeGrowthRate       string    `json:"employeeGrowthRate"`
	WebScrapedUrl            string    `json:"webScrapedUrl"`
	Source                   string    `json:"source"`
	UpdatedAt                time.Time `json:"updatedAt"`
	UpdateName               bool      `json:"updateName"`
	UpdateDescription        bool      `json:"updateDescription"`
	UpdateHide               bool      `json:"updateHide"`
	UpdateIsCustomer         bool      `json:"updateIsCustomer"`
	UpdateWebsite            bool      `json:"updateWebsite"`
	UpdateIndustry           bool      `json:"updateIndustry"`
	UpdateSubIndustry        bool      `json:"updateSubIndustry"`
	UpdateIndustryGroup      bool      `json:"updateIndustryGroup"`
	UpdateTargetAudience     bool      `json:"updateTargetAudience"`
	UpdateValueProposition   bool      `json:"updateValueProposition"`
	UpdateLastFundingRound   bool      `json:"updateLastFundingRound"`
	UpdateLastFundingAmount  bool      `json:"updateLastFundingAmount"`
	UpdateReferenceId        bool      `json:"updateReferenceId"`
	UpdateNote               bool      `json:"updateNote"`
	UpdateIsPublic           bool      `json:"updateIsPublic"`
	UpdateEmployees          bool      `json:"updateEmployees"`
	UpdateMarket             bool      `json:"updateMarket"`
	UpdateYearFounded        bool      `json:"updateYearFounded"`
	UpdateHeadquarters       bool      `json:"updateHeadquarters"`
	UpdateLogoUrl            bool      `json:"updateLogoUrl"`
	UpdateEmployeeGrowthRate bool      `json:"updateEmployeeGrowthRate"`
}

type OrganizationWriteRepository interface {
	CreateOrganization(ctx context.Context, tenant, organizationId string, data OrganizationCreateFields) error
	CreateOrganizationInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant, organizationId string, data OrganizationCreateFields) error
	UpdateOrganization(ctx context.Context, tenant, organizationId string, data OrganizationUpdateFields) error
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
						org.syncedWithEventStore = true 
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
		"source":             data.SourceFields.Source,
		"sourceOfTruth":      data.SourceFields.SourceOfTruth,
		"appSource":          data.SourceFields.AppSource,
		"createdAt":          data.CreatedAt,
		"updatedAt":          data.UpdatedAt,
		"onboardingStatus":   string(neo4jenum.OnboardingStatusNotApplicable),
		"overwrite":          data.SourceFields.Source == constants.SourceOpenline,
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
		"id":                 organizationId,
		"tenant":             tenant,
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
		"source":             data.Source,
		"updatedAt":          data.UpdatedAt,
		"overwrite":          data.Source == constants.SourceOpenline || data.Source == constants.SourceWebscrape,
	}
	cypher := `MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization {id:$id}) SET `
	if data.UpdateName {
		cypher += `org.name = CASE WHEN org.sourceOfTruth=$source OR $overwrite=true OR org.name = '' THEN $name ELSE org.name END,`
	}
	if data.UpdateDescription {
		cypher += `org.description = CASE WHEN org.sourceOfTruth=$source OR $overwrite=true OR org.description = '' THEN $description ELSE org.description END,`
	}
	if data.UpdateHide {
		cypher += `org.hide = CASE WHEN $overwrite=true OR (org.sourceOfTruth=$source AND $hide = false) THEN $hide ELSE org.hide END,`
	}
	if data.UpdateIsCustomer {
		cypher += `org.isCustomer = CASE WHEN $overwrite=true OR (org.sourceOfTruth=$source AND $isCustomer = true) THEN $isCustomer ELSE org.isCustomer END,`
	}
	if data.UpdateWebsite {
		cypher += `org.website = CASE WHEN org.sourceOfTruth=$source OR $overwrite=true OR org.website is null OR org.website = '' THEN $website ELSE org.website END,`
	}
	if data.UpdateIndustry {
		cypher += `org.industry = CASE WHEN org.sourceOfTruth=$source OR $overwrite=true OR org.industry is null OR org.industry = '' THEN $industry ELSE org.industry END,`
	}
	if data.UpdateSubIndustry {
		cypher += `org.subIndustry = CASE WHEN org.sourceOfTruth=$source OR $overwrite=true OR org.subIndustry is null OR org.subIndustry = '' THEN $subIndustry ELSE org.subIndustry END,`
	}
	if data.UpdateIndustryGroup {
		cypher += `org.industryGroup = CASE WHEN org.sourceOfTruth=$source OR $overwrite=true OR org.industryGroup is null OR org.industryGroup = '' THEN $industryGroup ELSE org.industryGroup END,`
	}
	if data.UpdateTargetAudience {
		cypher += `org.targetAudience = CASE WHEN org.sourceOfTruth=$source OR $overwrite=true OR org.targetAudience is null OR org.targetAudience = '' THEN $targetAudience ELSE org.targetAudience END,`
	}
	if data.UpdateValueProposition {
		cypher += `org.valueProposition = CASE WHEN org.sourceOfTruth=$source OR $overwrite=true OR org.valueProposition is null OR org.valueProposition = '' THEN $valueProposition ELSE org.valueProposition END,`
	}
	if data.UpdateLastFundingRound {
		cypher += `org.lastFundingRound = CASE WHEN org.sourceOfTruth=$source OR $overwrite=true OR org.lastFundingRound is null OR org.lastFundingRound = '' THEN $lastFundingRound ELSE org.lastFundingRound END,`
	}
	if data.UpdateLastFundingAmount {
		cypher += `org.lastFundingAmount = CASE WHEN org.sourceOfTruth=$source OR $overwrite=true OR org.lastFundingAmount is null OR org.lastFundingAmount = '' THEN $lastFundingAmount ELSE org.lastFundingAmount END,`
	}
	if data.UpdateReferenceId {
		cypher += `org.referenceId = CASE WHEN org.sourceOfTruth=$source OR $overwrite=true OR org.referenceId is null OR org.referenceId = '' THEN $referenceId ELSE org.referenceId END,`
	}
	if data.UpdateNote {
		cypher += `org.note = CASE WHEN org.sourceOfTruth=$source OR $overwrite=true OR org.note is null OR org.note = '' THEN $note ELSE org.note END,`
	}
	if data.UpdateIsPublic {
		cypher += `org.isPublic = CASE WHEN org.sourceOfTruth=$source OR $overwrite=true OR org.isPublic is null THEN $isPublic ELSE org.isPublic END,`
	}
	if data.UpdateEmployees {
		cypher += `org.employees = CASE WHEN org.sourceOfTruth=$source OR $overwrite=true OR org.employees is null THEN $employees ELSE org.employees END,`
	}
	if data.UpdateMarket {
		cypher += `org.market = CASE WHEN org.sourceOfTruth=$source OR $overwrite=true OR org.market is null OR org.market = '' THEN $market ELSE org.market END,`
	}
	if data.UpdateYearFounded {
		cypher += `org.yearFounded = CASE WHEN org.sourceOfTruth=$source OR $overwrite=true OR org.yearFounded is null OR org.yearFounded = 0 THEN $yearFounded ELSE org.yearFounded END,`
	}
	if data.UpdateHeadquarters {
		cypher += `org.headquarters = CASE WHEN org.sourceOfTruth=$source OR $overwrite=true OR org.headquarters is null OR org.headquarters = '' THEN $headquarters ELSE org.headquarters END,`
	}
	if data.UpdateLogoUrl {
		cypher += `org.logoUrl = CASE WHEN org.sourceOfTruth=$source OR $overwrite=true OR org.logoUrl is null OR org.logoUrl = '' THEN $logoUrl ELSE org.logoUrl END,`
	}
	if data.UpdateEmployeeGrowthRate {
		cypher += `org.employeeGrowthRate = CASE WHEN org.sourceOfTruth=$source OR $overwrite=true OR org.employeeGrowthRate is null OR org.employeeGrowthRate = '' THEN $employeeGrowthRate ELSE org.employeeGrowthRate END,`
	}
	if data.WebScrapedUrl != "" {
		params["webScrapedUrl"] = data.WebScrapedUrl
		params["webScrapedAt"] = utils.Now()
		cypher += `org.webScrapedUrl = $webScrapedUrl, org.webScrapedAt = $webScrapedAt,`
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
