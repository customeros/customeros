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

type OrganizationWriteRepository interface {
	CreateOrganization(ctx context.Context, tenant, organizationId string, data OrganizationCreateFields) error
	CreateOrganizationInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant, organizationId string, data OrganizationCreateFields) error
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
