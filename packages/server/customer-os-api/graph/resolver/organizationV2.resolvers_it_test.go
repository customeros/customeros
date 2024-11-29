package resolver

import (
	"context"
	"github.com/99designs/gqlgen/client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils/decode"
	commonModel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	neo4jtest "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/test"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestQueryResolver_UIOrganizationsSearch_FilterByName(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)

	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{Name: ""})
	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{Name: "A closed organization"})
	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{Name: "OPENLINE"})
	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{Name: "the openline"})
	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{Name: "some other open organization"})
	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{Name: "OpEnLiNe"})

	require.Equal(t, 6, neo4jtest.GetCountOfNodes(ctx, driver, "Organization"))

	searchBy := model.ColumnViewTypeOrganizationsName
	searchTerm := "open"

	assertSearch(t, searchBy, searchTerm, commonModel.ComparisonOperatorIsEmpty, 6, 1)
	assertSearch(t, searchBy, searchTerm, commonModel.ComparisonOperatorIsNotEmpty, 6, 5)
	assertSearch(t, searchBy, searchTerm, commonModel.ComparisonOperatorContains, 6, 4)
	assertSearch(t, searchBy, searchTerm, commonModel.ComparisonOperatorNotContains, 6, 2)
}

func TestQueryResolver_UIOrganizationsSearch_FilterByWebsite(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)

	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{Website: ""})
	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{Website: "https://www.customeros.ai"})
	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{Website: "customeros.ai"})
	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{Website: "www.customeros.ai"})
	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{Website: "www.google.com"})

	require.Equal(t, 5, neo4jtest.GetCountOfNodes(ctx, driver, "Organization"))

	searchBy := model.ColumnViewTypeOrganizationsWebsite
	searchTerm := "customeros"

	assertSearch(t, searchBy, searchTerm, commonModel.ComparisonOperatorIsEmpty, 5, 1)
	assertSearch(t, searchBy, searchTerm, commonModel.ComparisonOperatorIsNotEmpty, 5, 4)
	assertSearch(t, searchBy, searchTerm, commonModel.ComparisonOperatorContains, 5, 3)
	assertSearch(t, searchBy, searchTerm, commonModel.ComparisonOperatorNotContains, 5, 2)
}

func TestQueryResolver_UIOrganizationsSearch_FilterByRelationship(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)

	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{Relationship: ""})
	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{Relationship: enum.OrganizationRelationshipProspect})
	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{Relationship: enum.OrganizationRelationshipCustomer})

	require.Equal(t, 3, neo4jtest.GetCountOfNodes(ctx, driver, "Organization"))

	searchBy := model.ColumnViewTypeOrganizationsRelationship

	assertSearch(t, searchBy, "", commonModel.ComparisonOperatorIsEmpty, 3, 1)
	assertSearch(t, searchBy, []string{enum.OrganizationRelationshipCustomer.String()}, commonModel.ComparisonOperatorIn, 3, 1)
	assertSearch(t, searchBy, []string{enum.OrganizationRelationshipProspect.String(), enum.OrganizationRelationshipCustomer.String()}, commonModel.ComparisonOperatorIn, 3, 2)
}

func TestQueryResolver_UIOrganizationsSearch_FilterByOnboardingStatus(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)

	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{OnboardingDetails: neo4jentity.OnboardingDetails{Status: ""}})
	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{OnboardingDetails: neo4jentity.OnboardingDetails{Status: string(enum.OnboardingStatusNotApplicable)}})
	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{OnboardingDetails: neo4jentity.OnboardingDetails{Status: string(enum.OnboardingStatusDone)}})

	require.Equal(t, 3, neo4jtest.GetCountOfNodes(ctx, driver, "Organization"))

	searchBy := model.ColumnViewTypeOrganizationsOnboardingStatus

	assertSearch(t, searchBy, "", commonModel.ComparisonOperatorIsEmpty, 3, 1)
	assertSearch(t, searchBy, []string{string(enum.OnboardingStatusNotApplicable)}, commonModel.ComparisonOperatorIn, 3, 1)
	assertSearch(t, searchBy, []string{string(enum.OnboardingStatusNotApplicable), string(enum.OnboardingStatusDone)}, commonModel.ComparisonOperatorIn, 3, 2)
}

func TestQueryResolver_UIOrganizationsSearch_FilterByRenewalLikelihood(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)

	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{RenewalSummary: neo4jentity.RenewalSummary{RenewalLikelihood: ""}})
	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{RenewalSummary: neo4jentity.RenewalSummary{RenewalLikelihood: string(enum.RenewalLikelihoodZero)}})
	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{RenewalSummary: neo4jentity.RenewalSummary{RenewalLikelihood: string(enum.RenewalLikelihoodHigh)}})

	require.Equal(t, 3, neo4jtest.GetCountOfNodes(ctx, driver, "Organization"))

	searchBy := model.ColumnViewTypeOrganizationsRenewalLikelihood

	assertSearch(t, searchBy, "", commonModel.ComparisonOperatorIsEmpty, 3, 1)
	assertSearch(t, searchBy, []string{string(enum.RenewalLikelihoodZero)}, commonModel.ComparisonOperatorIn, 3, 1)
	assertSearch(t, searchBy, []string{string(enum.RenewalLikelihoodZero), string(enum.RenewalLikelihoodHigh)}, commonModel.ComparisonOperatorIn, 3, 2)
}

func TestQueryResolver_UIOrganizationsSearch_FilterByRenewalDate(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)

	firstOfDecember := utils.FirstTimeOfMonth(2023, 12)
	firstOfJanuary := utils.FirstTimeOfMonth(2024, 1)
	midOfJanuary := utils.FirstTimeOfMonth(2024, 1).Add(time.Hour * 24 * 15)
	firstOfFebruary := utils.FirstTimeOfMonth(2024, 2)
	midOfFebruary := utils.FirstTimeOfMonth(2024, 2).Add(time.Hour * 24 * 15)
	firstOfMarch := utils.FirstTimeOfMonth(2024, 3)

	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{RenewalSummary: neo4jentity.RenewalSummary{NextRenewalAt: nil}})
	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{RenewalSummary: neo4jentity.RenewalSummary{NextRenewalAt: &firstOfJanuary}})
	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{RenewalSummary: neo4jentity.RenewalSummary{NextRenewalAt: &firstOfFebruary}})

	require.Equal(t, 3, neo4jtest.GetCountOfNodes(ctx, driver, "Organization"))

	searchBy := model.ColumnViewTypeOrganizationsRenewalDate

	assertSearch(t, searchBy, "", commonModel.ComparisonOperatorIsEmpty, 3, 1)
	assertSearch(t, searchBy, []time.Time{firstOfDecember, midOfJanuary}, commonModel.ComparisonOperatorBetween, 3, 1)
	assertSearch(t, searchBy, []time.Time{firstOfJanuary, firstOfJanuary}, commonModel.ComparisonOperatorBetween, 3, 1)
	assertSearch(t, searchBy, []time.Time{firstOfJanuary, firstOfFebruary}, commonModel.ComparisonOperatorBetween, 3, 2)
	assertSearch(t, searchBy, []time.Time{firstOfFebruary, firstOfFebruary}, commonModel.ComparisonOperatorBetween, 3, 1)
	assertSearch(t, searchBy, []time.Time{midOfJanuary, firstOfMarch}, commonModel.ComparisonOperatorBetween, 3, 1)
	assertSearch(t, searchBy, []time.Time{midOfFebruary, firstOfMarch}, commonModel.ComparisonOperatorBetween, 3, 0)
}

func TestQueryResolver_UIOrganizationsSearch_FilterByLastTouchpointAt(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)

	firstOfDecember := utils.FirstTimeOfMonth(2023, 12)
	firstOfJanuary := utils.FirstTimeOfMonth(2024, 1)
	midOfJanuary := utils.FirstTimeOfMonth(2024, 1).Add(time.Hour * 24 * 15)
	firstOfFebruary := utils.FirstTimeOfMonth(2024, 2)
	midOfFebruary := utils.FirstTimeOfMonth(2024, 2).Add(time.Hour * 24 * 15)
	firstOfMarch := utils.FirstTimeOfMonth(2024, 3)

	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{LastTouchpointAt: nil})
	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{LastTouchpointAt: &firstOfJanuary})
	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{LastTouchpointAt: &firstOfFebruary})

	require.Equal(t, 3, neo4jtest.GetCountOfNodes(ctx, driver, "Organization"))

	searchBy := model.ColumnViewTypeOrganizationsLastTouchpointDate

	assertSearch(t, searchBy, "", commonModel.ComparisonOperatorIsEmpty, 3, 1)
	assertSearch(t, searchBy, []time.Time{firstOfDecember, midOfJanuary}, commonModel.ComparisonOperatorBetween, 3, 1)
	assertSearch(t, searchBy, []time.Time{firstOfJanuary, firstOfJanuary}, commonModel.ComparisonOperatorBetween, 3, 1)
	assertSearch(t, searchBy, []time.Time{firstOfJanuary, firstOfFebruary}, commonModel.ComparisonOperatorBetween, 3, 2)
	assertSearch(t, searchBy, []time.Time{firstOfFebruary, firstOfFebruary}, commonModel.ComparisonOperatorBetween, 3, 1)
	assertSearch(t, searchBy, []time.Time{midOfJanuary, firstOfMarch}, commonModel.ComparisonOperatorBetween, 3, 1)
	assertSearch(t, searchBy, []time.Time{midOfFebruary, firstOfMarch}, commonModel.ComparisonOperatorBetween, 3, 0)
}

func TestQueryResolver_UIOrganizationsSearch_FilterByStage(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)

	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{Stage: ""})
	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{Stage: enum.Lead})
	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{Stage: enum.Trial})
	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{Stage: enum.Lead})

	require.Equal(t, 4, neo4jtest.GetCountOfNodes(ctx, driver, "Organization"))

	searchBy := model.ColumnViewTypeOrganizationsStage

	assertSearch(t, searchBy, "", commonModel.ComparisonOperatorIsEmpty, 4, 1)
	assertSearch(t, searchBy, []string{enum.Trial.String()}, commonModel.ComparisonOperatorIn, 4, 1)
	assertSearch(t, searchBy, []string{enum.Lead.String()}, commonModel.ComparisonOperatorIn, 4, 2)
	assertSearch(t, searchBy, []string{enum.Trial.String(), enum.Lead.String()}, commonModel.ComparisonOperatorIn, 4, 3)
}

func TestQueryResolver_UIOrganizationsSearch_FilterBySocials(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)

	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{ID: "org1"})
	neo4jtest.CreateSocial(ctx, driver, tenantName, neo4jentity.SocialEntity{Id: "social1", Url: "https://www.linkedin.com/company/openline-ai"})
	neo4jtest.LinkNodes(ctx, driver, "org1", "social1", "HAS")

	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{ID: "org2"})
	neo4jtest.CreateSocial(ctx, driver, tenantName, neo4jentity.SocialEntity{Id: "social2", Url: "https://www.twitter.com/company/openline-ai"})
	neo4jtest.LinkNodes(ctx, driver, "org2", "social2", "HAS")

	require.Equal(t, 2, neo4jtest.GetCountOfNodes(ctx, driver, "Organization"))
	require.Equal(t, 2, neo4jtest.GetCountOfNodes(ctx, driver, "Social"))
	require.Equal(t, 2, neo4jtest.GetCountOfRelationships(ctx, driver, "HAS"))

	searchBy := model.ColumnViewTypeOrganizationsSocials

	assertSearch(t, searchBy, "", commonModel.ComparisonOperatorContains, 2, 2)
	assertSearch(t, searchBy, "openline", commonModel.ComparisonOperatorContains, 2, 2)
	assertSearch(t, searchBy, "linkedin", commonModel.ComparisonOperatorContains, 2, 1)
	assertSearch(t, searchBy, "twitter", commonModel.ComparisonOperatorContains, 2, 1)
}

func TestQueryResolver_UIOrganizationsSearch_FilterByLeadSource(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)

	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{LeadSource: ""})
	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{LeadSource: "A"})
	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{LeadSource: "B"})
	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{LeadSource: "AB"})
	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{LeadSource: "C"})

	require.Equal(t, 5, neo4jtest.GetCountOfNodes(ctx, driver, "Organization"))

	searchBy := model.ColumnViewTypeOrganizationsLeadSource

	assertSearch(t, searchBy, "", commonModel.ComparisonOperatorIsEmpty, 5, 1)
	assertSearch(t, searchBy, "", commonModel.ComparisonOperatorIsNotEmpty, 5, 4)
	assertSearch(t, searchBy, "C", commonModel.ComparisonOperatorContains, 5, 1)
	assertSearch(t, searchBy, "A", commonModel.ComparisonOperatorContains, 5, 2)
	assertSearch(t, searchBy, "A", commonModel.ComparisonOperatorNotContains, 5, 3)
}

func TestQueryResolver_UIOrganizationsSearch_FilterByCreatedAt(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)

	firstOfDecember := utils.FirstTimeOfMonth(2023, 12)
	firstOfJanuary := utils.FirstTimeOfMonth(2024, 1)
	midOfJanuary := utils.FirstTimeOfMonth(2024, 1).Add(time.Hour * 24 * 15)
	firstOfFebruary := utils.FirstTimeOfMonth(2024, 2)
	midOfFebruary := utils.FirstTimeOfMonth(2024, 2).Add(time.Hour * 24 * 15)
	firstOfMarch := utils.FirstTimeOfMonth(2024, 3)

	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{CreatedAt: firstOfJanuary})
	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{CreatedAt: firstOfFebruary})

	require.Equal(t, 2, neo4jtest.GetCountOfNodes(ctx, driver, "Organization"))

	searchBy := model.ColumnViewTypeOrganizationsCreatedDate

	assertSearch(t, searchBy, "", commonModel.ComparisonOperatorIsEmpty, 2, 0)
	assertSearch(t, searchBy, []time.Time{firstOfDecember, midOfJanuary}, commonModel.ComparisonOperatorBetween, 2, 1)
	assertSearch(t, searchBy, []time.Time{firstOfJanuary, firstOfJanuary}, commonModel.ComparisonOperatorBetween, 2, 1)
	assertSearch(t, searchBy, []time.Time{firstOfJanuary, firstOfFebruary}, commonModel.ComparisonOperatorBetween, 2, 2)
	assertSearch(t, searchBy, []time.Time{firstOfFebruary, firstOfFebruary}, commonModel.ComparisonOperatorBetween, 2, 1)
	assertSearch(t, searchBy, []time.Time{midOfJanuary, firstOfMarch}, commonModel.ComparisonOperatorBetween, 2, 1)
	assertSearch(t, searchBy, []time.Time{midOfFebruary, firstOfMarch}, commonModel.ComparisonOperatorBetween, 2, 0)
}

func TestQueryResolver_UIOrganizationsSearch_FilterByEmployeeCount(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)

	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{})
	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{Employees: 2000})
	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{Employees: 2005})

	require.Equal(t, 3, neo4jtest.GetCountOfNodes(ctx, driver, "Organization"))

	searchBy := model.ColumnViewTypeOrganizationsEmployeeCount

	assertSearch(t, searchBy, "", commonModel.ComparisonOperatorIsEmpty, 3, 0)
	assertSearch(t, searchBy, "", commonModel.ComparisonOperatorIsNotEmpty, 3, 3)
	assertSearch(t, searchBy, 2005, commonModel.ComparisonOperatorGte, 3, 1)
	assertSearch(t, searchBy, 2005, commonModel.ComparisonOperatorGt, 3, 0)
	assertSearch(t, searchBy, 2004, commonModel.ComparisonOperatorLte, 3, 2)
	assertSearch(t, searchBy, 2005, commonModel.ComparisonOperatorLte, 3, 3)
	assertSearch(t, searchBy, 2005, commonModel.ComparisonOperatorLt, 3, 2)
	assertSearch(t, searchBy, 2005, commonModel.ComparisonOperatorEquals, 3, 1)
	assertSearch(t, searchBy, 2006, commonModel.ComparisonOperatorEquals, 3, 0)
	assertSearch(t, searchBy, 2005, commonModel.ComparisonOperatorNotEquals, 3, 2)
	assertSearch(t, searchBy, 2010, commonModel.ComparisonOperatorNotEquals, 3, 3)
}

func TestQueryResolver_UIOrganizationsSearch_FilterByYearFounded(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)

	founded1 := int64(2000)
	founded2 := int64(2005)

	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{YearFounded: nil})
	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{YearFounded: &founded1})
	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{YearFounded: &founded2})

	require.Equal(t, 3, neo4jtest.GetCountOfNodes(ctx, driver, "Organization"))

	searchBy := model.ColumnViewTypeOrganizationsYearFounded

	assertSearch(t, searchBy, "", commonModel.ComparisonOperatorIsEmpty, 3, 1)
	assertSearch(t, searchBy, "", commonModel.ComparisonOperatorIsNotEmpty, 3, 2)
	assertSearch(t, searchBy, 2005, commonModel.ComparisonOperatorGte, 3, 1)
	assertSearch(t, searchBy, 2005, commonModel.ComparisonOperatorGt, 3, 0)
	assertSearch(t, searchBy, 2005, commonModel.ComparisonOperatorLte, 3, 2)
	assertSearch(t, searchBy, 2005, commonModel.ComparisonOperatorLt, 3, 1)
	assertSearch(t, searchBy, 2005, commonModel.ComparisonOperatorEquals, 3, 1)
	assertSearch(t, searchBy, 2006, commonModel.ComparisonOperatorEquals, 3, 0)
	assertSearch(t, searchBy, 2005, commonModel.ComparisonOperatorNotEquals, 3, 2)
	assertSearch(t, searchBy, 2010, commonModel.ComparisonOperatorNotEquals, 3, 3)
}

func TestQueryResolver_UIOrganizationsSearch_FilterByIndustry(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)

	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{Industry: ""})
	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{Industry: "A"})
	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{Industry: "B"})
	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{Industry: "AB"})
	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{Industry: "C"})

	require.Equal(t, 5, neo4jtest.GetCountOfNodes(ctx, driver, "Organization"))

	searchBy := model.ColumnViewTypeOrganizationsIndustry

	assertSearch(t, searchBy, "", commonModel.ComparisonOperatorIsEmpty, 5, 1)
	assertSearch(t, searchBy, "", commonModel.ComparisonOperatorIsNotEmpty, 5, 4)
	assertSearch(t, searchBy, "C", commonModel.ComparisonOperatorContains, 5, 1)
	assertSearch(t, searchBy, "A", commonModel.ComparisonOperatorContains, 5, 2)
	assertSearch(t, searchBy, "A", commonModel.ComparisonOperatorNotContains, 5, 3)
}

func TestQueryResolver_UIOrganizationsSearch_FilterByChurnedAt(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)

	firstOfDecember := utils.FirstTimeOfMonth(2023, 12)
	firstOfJanuary := utils.FirstTimeOfMonth(2024, 1)
	midOfJanuary := utils.FirstTimeOfMonth(2024, 1).Add(time.Hour * 24 * 15)
	firstOfFebruary := utils.FirstTimeOfMonth(2024, 2)
	midOfFebruary := utils.FirstTimeOfMonth(2024, 2).Add(time.Hour * 24 * 15)
	firstOfMarch := utils.FirstTimeOfMonth(2024, 3)

	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{DerivedData: neo4jentity.DerivedData{ChurnedAt: nil}})
	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{DerivedData: neo4jentity.DerivedData{ChurnedAt: &firstOfJanuary}})
	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{DerivedData: neo4jentity.DerivedData{ChurnedAt: &firstOfFebruary}})

	require.Equal(t, 3, neo4jtest.GetCountOfNodes(ctx, driver, "Organization"))

	searchBy := model.ColumnViewTypeOrganizationsChurnDate

	assertSearch(t, searchBy, "", commonModel.ComparisonOperatorIsEmpty, 3, 1)
	assertSearch(t, searchBy, []time.Time{firstOfDecember, midOfJanuary}, commonModel.ComparisonOperatorBetween, 3, 1)
	assertSearch(t, searchBy, []time.Time{firstOfJanuary, firstOfJanuary}, commonModel.ComparisonOperatorBetween, 3, 1)
	assertSearch(t, searchBy, []time.Time{firstOfJanuary, firstOfFebruary}, commonModel.ComparisonOperatorBetween, 3, 2)
	assertSearch(t, searchBy, []time.Time{firstOfFebruary, firstOfFebruary}, commonModel.ComparisonOperatorBetween, 3, 1)
	assertSearch(t, searchBy, []time.Time{midOfJanuary, firstOfMarch}, commonModel.ComparisonOperatorBetween, 3, 1)
	assertSearch(t, searchBy, []time.Time{midOfFebruary, firstOfMarch}, commonModel.ComparisonOperatorBetween, 3, 0)
}

func TestQueryResolver_UIOrganizationsSearch_FilterByLtv(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)

	ltv1 := float64(2000)
	ltv2 := float64(2005)

	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{})
	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{DerivedData: neo4jentity.DerivedData{Ltv: ltv1}})
	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{DerivedData: neo4jentity.DerivedData{Ltv: ltv2}})

	require.Equal(t, 3, neo4jtest.GetCountOfNodes(ctx, driver, "Organization"))

	searchBy := model.ColumnViewTypeOrganizationsLtv

	assertSearch(t, searchBy, "", commonModel.ComparisonOperatorIsEmpty, 3, 0)
	assertSearch(t, searchBy, "", commonModel.ComparisonOperatorIsNotEmpty, 3, 3)
	assertSearch(t, searchBy, 2005, commonModel.ComparisonOperatorGte, 3, 1)
	assertSearch(t, searchBy, 2005.0, commonModel.ComparisonOperatorGte, 3, 1)
	assertSearch(t, searchBy, 2005, commonModel.ComparisonOperatorGt, 3, 0)
	assertSearch(t, searchBy, 2005, commonModel.ComparisonOperatorLte, 3, 3)
	assertSearch(t, searchBy, 2005, commonModel.ComparisonOperatorLt, 3, 2)
	assertSearch(t, searchBy, 0, commonModel.ComparisonOperatorEquals, 3, 1)
	assertSearch(t, searchBy, 2005, commonModel.ComparisonOperatorEquals, 3, 1)
	assertSearch(t, searchBy, 2006, commonModel.ComparisonOperatorEquals, 3, 0)
	assertSearch(t, searchBy, 2005, commonModel.ComparisonOperatorNotEquals, 3, 2)
	assertSearch(t, searchBy, 2010, commonModel.ComparisonOperatorNotEquals, 3, 3)
}

func TestQueryResolver_UIOrganizationsSearch_FilterByHeadquarters(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)

	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{Headquarters: ""})
	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{Headquarters: "A"})
	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{Headquarters: "B"})
	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{Headquarters: "AB"})
	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{Headquarters: "C"})

	require.Equal(t, 5, neo4jtest.GetCountOfNodes(ctx, driver, "Organization"))

	searchBy := model.ColumnViewTypeOrganizationsHeadquarters

	assertSearch(t, searchBy, "", commonModel.ComparisonOperatorIsEmpty, 5, 1)
	assertSearch(t, searchBy, "", commonModel.ComparisonOperatorIsNotEmpty, 5, 4)
	assertSearch(t, searchBy, "C", commonModel.ComparisonOperatorContains, 5, 1)
	assertSearch(t, searchBy, "A", commonModel.ComparisonOperatorContains, 5, 2)
	assertSearch(t, searchBy, "A", commonModel.ComparisonOperatorNotContains, 5, 3)
}

func TestQueryResolver_UIOrganizationsSearch_FilterByIsPublic(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)

	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{})
	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{IsPublic: true})
	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{IsPublic: false})

	require.Equal(t, 3, neo4jtest.GetCountOfNodes(ctx, driver, "Organization"))

	searchBy := model.ColumnViewTypeOrganizationsIsPublic

	assertSearch(t, searchBy, true, commonModel.ComparisonOperatorEquals, 3, 1)
	assertSearch(t, searchBy, false, commonModel.ComparisonOperatorEquals, 3, 2)
	assertSearch(t, searchBy, true, commonModel.ComparisonOperatorNotEquals, 3, 2)
	assertSearch(t, searchBy, false, commonModel.ComparisonOperatorNotEquals, 3, 1)
}

func TestQueryResolver_UIOrganizationsSearch_FilterByTags(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)

	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{ID: "org1"})
	neo4jtest.CreateTag(ctx, driver, tenantName, neo4jentity.TagEntity{Id: "tag1", Name: "A"})
	neo4jtest.CreateTag(ctx, driver, tenantName, neo4jentity.TagEntity{Id: "tag2", Name: "B"})
	neo4jtest.LinkNodes(ctx, driver, "org1", "tag1", "TAGGED")
	neo4jtest.LinkNodes(ctx, driver, "org1", "tag2", "TAGGED")

	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{ID: "org2"})
	neo4jtest.CreateTag(ctx, driver, tenantName, neo4jentity.TagEntity{Id: "tag3", Name: "a"})
	neo4jtest.CreateTag(ctx, driver, tenantName, neo4jentity.TagEntity{Id: "tag4", Name: "c"})
	neo4jtest.LinkNodes(ctx, driver, "org2", "tag3", "TAGGED")
	neo4jtest.LinkNodes(ctx, driver, "org2", "tag4", "TAGGED")

	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{ID: "org3"})

	require.Equal(t, 3, neo4jtest.GetCountOfNodes(ctx, driver, "Organization"))
	require.Equal(t, 4, neo4jtest.GetCountOfNodes(ctx, driver, "Tag"))
	require.Equal(t, 4, neo4jtest.GetCountOfRelationships(ctx, driver, "TAGGED"))

	searchBy := model.ColumnViewTypeOrganizationsTags

	assertSearch(t, searchBy, "", commonModel.ComparisonOperatorIsEmpty, 3, 1)
	assertSearch(t, searchBy, "", commonModel.ComparisonOperatorIsNotEmpty, 3, 2)
	assertSearch(t, searchBy, "", commonModel.ComparisonOperatorContains, 3, 2)
	assertSearch(t, searchBy, "A", commonModel.ComparisonOperatorContains, 3, 2)
	assertSearch(t, searchBy, "B", commonModel.ComparisonOperatorContains, 3, 1)
	assertSearch(t, searchBy, "C", commonModel.ComparisonOperatorContains, 3, 1)
	assertSearch(t, searchBy, "D", commonModel.ComparisonOperatorContains, 3, 0)

	assertSearch(t, searchBy, "D", commonModel.ComparisonOperatorNotContains, 3, 3)
	assertSearch(t, searchBy, "B", commonModel.ComparisonOperatorNotContains, 3, 2)
	assertSearch(t, searchBy, "A", commonModel.ComparisonOperatorNotContains, 3, 1)
}

func TestQueryResolver_UIOrganizationsSearch_FilterByUpdatedAt(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)

	firstOfDecember := utils.FirstTimeOfMonth(2023, 12)
	firstOfJanuary := utils.FirstTimeOfMonth(2024, 1)
	midOfJanuary := utils.FirstTimeOfMonth(2024, 1).Add(time.Hour * 24 * 15)
	firstOfFebruary := utils.FirstTimeOfMonth(2024, 2)
	midOfFebruary := utils.FirstTimeOfMonth(2024, 2).Add(time.Hour * 24 * 15)
	firstOfMarch := utils.FirstTimeOfMonth(2024, 3)

	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{})
	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{UpdatedAt: firstOfJanuary})
	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{UpdatedAt: firstOfFebruary})

	require.Equal(t, 3, neo4jtest.GetCountOfNodes(ctx, driver, "Organization"))

	searchBy := model.ColumnViewTypeOrganizationsUpdatedDate

	assertSearch(t, searchBy, "", commonModel.ComparisonOperatorIsEmpty, 3, 0)
	assertSearch(t, searchBy, []time.Time{firstOfDecember, midOfJanuary}, commonModel.ComparisonOperatorBetween, 3, 1)
	assertSearch(t, searchBy, []time.Time{firstOfJanuary, firstOfJanuary}, commonModel.ComparisonOperatorBetween, 3, 1)
	assertSearch(t, searchBy, []time.Time{firstOfJanuary, firstOfFebruary}, commonModel.ComparisonOperatorBetween, 3, 2)
	assertSearch(t, searchBy, []time.Time{firstOfFebruary, firstOfFebruary}, commonModel.ComparisonOperatorBetween, 3, 1)
	assertSearch(t, searchBy, []time.Time{midOfJanuary, firstOfMarch}, commonModel.ComparisonOperatorBetween, 3, 1)
	assertSearch(t, searchBy, []time.Time{midOfFebruary, firstOfMarch}, commonModel.ComparisonOperatorBetween, 3, 0)
}

func assertSearch(t *testing.T, filterName model.ColumnViewType, searchValue any, operator commonModel.ComparisonOperator, totalAvailable int64, totalElements int64) {
	rawResponse, err := c.RawPost(getQuery("organization/ui_organizations_search"),
		client.Var("limit", 10),
		client.Var("filterName", filterName),
		client.Var("filterValue", searchValue),
		client.Var("filterOperation", operator),
		client.Var("sortByField", model.ColumnViewTypeOrganizationsName),
		client.Var("sortByDirection", commonModel.SortingDirectionAsc),
	)
	assertRawResponseSuccess(t, rawResponse, err)

	var organizations struct {
		Ui_Organizations_Search model.OrganizationSearchResult
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &organizations)
	require.Nil(t, err)
	require.NotNil(t, organizations)

	searchResult := organizations.Ui_Organizations_Search
	require.Equal(t, totalAvailable, searchResult.TotalAvailable)
	require.Equal(t, totalElements, searchResult.TotalElements)
	require.Equal(t, int(totalElements), len(searchResult.Ids))
}
