package resolver

import (
	"context"
	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	"testing"
	"time"

	"github.com/99designs/gqlgen/client"
	commonModel "github.com/customeros/customeros/packages/server/customer-os-common-module/model"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"
	"github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/enum"
	neo4jtest "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/customeros/customeros/packages/server/customer-os-api/graphql/model"
	"github.com/customeros/customeros/packages/server/customer-os-api/utils/decode"
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

	searchBy := postgres_entity.ColumnViewTypeOrganizationsName
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

	searchBy := postgres_entity.ColumnViewTypeOrganizationsWebsite
	searchTerm := "customeros"

	assertSearch(t, searchBy, searchTerm, commonModel.ComparisonOperatorIsEmpty, 5, 1)
	assertSearch(t, searchBy, searchTerm, commonModel.ComparisonOperatorIsNotEmpty, 5, 4)
	assertSearch(t, searchBy, searchTerm, commonModel.ComparisonOperatorContains, 5, 3)
	assertSearch(t, searchBy, searchTerm, commonModel.ComparisonOperatorNotContains, 5, 2)
}

func TestQueryResolver_UIOrganizationsSearch_FilterByPrimaryDomain(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)

	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{ID: "org1"})
	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{ID: "org2"})
	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{ID: "org3"})
	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{ID: "org4"})

	neo4jtest.LinkDomainToOrganization(ctx, driver, "org1", neo4jentity.DomainEntity{Domain: "domain1", IsPrimary: utils.BoolPtr(true)})
	neo4jtest.LinkDomainToOrganization(ctx, driver, "org1", neo4jentity.DomainEntity{Domain: "domain1secondary1"})
	neo4jtest.LinkDomainToOrganization(ctx, driver, "org1", neo4jentity.DomainEntity{Domain: "domain1secondary2"})
	neo4jtest.LinkDomainToOrganization(ctx, driver, "org2", neo4jentity.DomainEntity{Domain: "domain2", IsPrimary: utils.BoolPtr(true)})
	neo4jtest.LinkDomainToOrganization(ctx, driver, "org3", neo4jentity.DomainEntity{Domain: "domain3", IsPrimary: utils.BoolPtr(true)})

	require.Equal(t, 4, neo4jtest.GetCountOfNodes(ctx, driver, "Organization"))
	require.Equal(t, 5, neo4jtest.GetCountOfNodes(ctx, driver, "Domain"))
	require.Equal(t, 5, neo4jtest.GetCountOfRelationships(ctx, driver, "HAS_DOMAIN"))

	searchBy := postgres_entity.ColumnViewTypeOrganizationsPrimaryDomains

	assertSearch(t, searchBy, "", commonModel.ComparisonOperatorIsEmpty, 4, 1)
	assertSearch(t, searchBy, "", commonModel.ComparisonOperatorIsNotEmpty, 4, 3)
	assertSearch(t, searchBy, "domain", commonModel.ComparisonOperatorContains, 4, 3)
	assertSearch(t, searchBy, "domain1", commonModel.ComparisonOperatorContains, 4, 1)
	assertSearch(t, searchBy, "domain2", commonModel.ComparisonOperatorContains, 4, 1)
	assertSearch(t, searchBy, "domain3", commonModel.ComparisonOperatorContains, 4, 1)
	assertSearch(t, searchBy, "domain4", commonModel.ComparisonOperatorContains, 4, 0)
	assertSearch(t, searchBy, "domain", commonModel.ComparisonOperatorNotContains, 4, 0)
	assertSearch(t, searchBy, "domain1", commonModel.ComparisonOperatorNotContains, 4, 2)
	assertSearch(t, searchBy, "domain4", commonModel.ComparisonOperatorNotContains, 4, 3)
}

func TestQueryResolver_UIOrganizationsSearch_FilterByRelationship(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)

	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{Relationship: ""})
	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{})
	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{Relationship: enum.OrganizationRelationshipProspect})
	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{Relationship: enum.OrganizationRelationshipCustomer})

	require.Equal(t, 4, neo4jtest.GetCountOfNodes(ctx, driver, "Organization"))

	searchBy := postgres_entity.ColumnViewTypeOrganizationsRelationship

	assertSearch(t, searchBy, "", commonModel.ComparisonOperatorIsEmpty, 4, 2)
	assertSearch(t, searchBy, []string{enum.OrganizationRelationshipCustomer.String()}, commonModel.ComparisonOperatorIn, 4, 1)
	assertSearch(t, searchBy, []string{enum.OrganizationRelationshipCustomer.String()}, commonModel.ComparisonOperatorNotIn, 4, 3)
	assertSearch(t, searchBy, []string{enum.OrganizationRelationshipProspect.String(), enum.OrganizationRelationshipCustomer.String()}, commonModel.ComparisonOperatorIn, 4, 2)
	assertSearch(t, searchBy, []string{enum.OrganizationRelationshipProspect.String(), enum.OrganizationRelationshipCustomer.String()}, commonModel.ComparisonOperatorNotIn, 4, 2)
}

func TestQueryResolver_UIOrganizationsSearch_FilterByOnboardingStatus(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)

	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{OnboardingDetails: neo4jentity.OnboardingDetails{Status: ""}})
	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{OnboardingDetails: neo4jentity.OnboardingDetails{Status: string(enum.OnboardingStatusNotApplicable)}})
	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{OnboardingDetails: neo4jentity.OnboardingDetails{Status: string(enum.OnboardingStatusDone)}})

	require.Equal(t, 3, neo4jtest.GetCountOfNodes(ctx, driver, "Organization"))

	searchBy := postgres_entity.ColumnViewTypeOrganizationsOnboardingStatus

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

	searchBy := postgres_entity.ColumnViewTypeOrganizationsRenewalLikelihood

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

	searchBy := postgres_entity.ColumnViewTypeOrganizationsRenewalDate

	assertSearch(t, searchBy, "", commonModel.ComparisonOperatorIsEmpty, 3, 1)
	assertSearch(t, searchBy, []time.Time{firstOfDecember, midOfJanuary}, commonModel.ComparisonOperatorBetween, 3, 1)
	assertSearch(t, searchBy, []time.Time{firstOfJanuary, firstOfJanuary}, commonModel.ComparisonOperatorBetween, 3, 1)
	assertSearch(t, searchBy, []time.Time{firstOfJanuary, firstOfFebruary}, commonModel.ComparisonOperatorBetween, 3, 2)
	assertSearch(t, searchBy, []time.Time{firstOfFebruary, firstOfFebruary}, commonModel.ComparisonOperatorBetween, 3, 1)
	assertSearch(t, searchBy, []time.Time{midOfJanuary, firstOfMarch}, commonModel.ComparisonOperatorBetween, 3, 1)
	assertSearch(t, searchBy, []time.Time{midOfFebruary, firstOfMarch}, commonModel.ComparisonOperatorBetween, 3, 0)
}

func TestQueryResolver_UIOrganizationsSearch_FilterByForecastArr(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)

	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{RenewalSummary: neo4jentity.RenewalSummary{ArrForecast: nil}})
	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{RenewalSummary: neo4jentity.RenewalSummary{ArrForecast: utils.Float64Ptr(2000)}})
	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{RenewalSummary: neo4jentity.RenewalSummary{ArrForecast: utils.Float64Ptr(2005)}})

	require.Equal(t, 3, neo4jtest.GetCountOfNodes(ctx, driver, "Organization"))

	searchBy := postgres_entity.ColumnViewTypeOrganizationsForecastArr

	assertSearch(t, searchBy, "", commonModel.ComparisonOperatorIsEmpty, 3, 1)
	assertSearch(t, searchBy, "", commonModel.ComparisonOperatorIsNotEmpty, 3, 2)
	assertSearch(t, searchBy, 2005, commonModel.ComparisonOperatorGte, 3, 1)
	assertSearch(t, searchBy, 2005, commonModel.ComparisonOperatorGt, 3, 0)
	assertSearch(t, searchBy, 2004, commonModel.ComparisonOperatorLte, 3, 1)
	assertSearch(t, searchBy, 2005, commonModel.ComparisonOperatorLte, 3, 2)
	assertSearch(t, searchBy, 2005, commonModel.ComparisonOperatorLt, 3, 1)
	assertSearch(t, searchBy, 2005, commonModel.ComparisonOperatorEquals, 3, 1)
	assertSearch(t, searchBy, 2006, commonModel.ComparisonOperatorEquals, 3, 0)
	assertSearch(t, searchBy, 2005, commonModel.ComparisonOperatorNotEquals, 3, 2)
	assertSearch(t, searchBy, 2010, commonModel.ComparisonOperatorNotEquals, 3, 3)
}

func TestQueryResolver_UIOrganizationsSearch_FilterByOwner(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)

	neo4jtest.CreateUserWithId(ctx, driver, tenantName, "owner1")
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, "owner2")

	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{ID: "1"})
	neo4jtest.LinkNodes(ctx, driver, "owner1", "1", "OWNS")

	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{ID: "2"})
	neo4jtest.LinkNodes(ctx, driver, "owner2", "2", "OWNS")

	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{ID: "empty"})

	require.Equal(t, 3, neo4jtest.GetCountOfNodes(ctx, driver, "Organization"))
	require.Equal(t, 2, neo4jtest.GetCountOfNodes(ctx, driver, "User"))
	require.Equal(t, 2, neo4jtest.GetCountOfRelationships(ctx, driver, "OWNS"))

	require.Equal(t, 3, neo4jtest.GetCountOfNodes(ctx, driver, "Organization"))

	searchBy := postgres_entity.ColumnViewTypeOrganizationsOwner

	assertSearch(t, searchBy, []string{"owner1"}, commonModel.ComparisonOperatorIn, 3, 1)
	assertSearch(t, searchBy, []string{"owner2"}, commonModel.ComparisonOperatorIn, 3, 1)
	assertSearch(t, searchBy, []string{"owner1", "owner2"}, commonModel.ComparisonOperatorIn, 3, 2)
}

func TestQueryResolver_UIOrganizationsSearch_FilterByLastTouchpoint(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)

	touchpoint1 := "A"
	touchpoint2 := "B"

	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{LastTouchpointType: nil})
	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{LastTouchpointType: &touchpoint1})
	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{LastTouchpointType: &touchpoint2})

	require.Equal(t, 3, neo4jtest.GetCountOfNodes(ctx, driver, "Organization"))

	searchBy := postgres_entity.ColumnViewTypeOrganizationsLastTouchpoint

	assertSearch(t, searchBy, "", commonModel.ComparisonOperatorIsEmpty, 3, 1)
	assertSearch(t, searchBy, "", commonModel.ComparisonOperatorIsNotEmpty, 3, 2)
	assertSearch(t, searchBy, []string{"A"}, commonModel.ComparisonOperatorIn, 3, 1)
	assertSearch(t, searchBy, []string{"B"}, commonModel.ComparisonOperatorIn, 3, 1)
	assertSearch(t, searchBy, []string{"A", "B"}, commonModel.ComparisonOperatorIn, 3, 2)
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

	searchBy := postgres_entity.ColumnViewTypeOrganizationsLastTouchpointDate

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

	searchBy := postgres_entity.ColumnViewTypeOrganizationsStage

	assertSearch(t, searchBy, "", commonModel.ComparisonOperatorIsEmpty, 4, 1)
	assertSearch(t, searchBy, []string{enum.Trial.String()}, commonModel.ComparisonOperatorIn, 4, 1)
	assertSearch(t, searchBy, []string{enum.Lead.String()}, commonModel.ComparisonOperatorIn, 4, 2)
	assertSearch(t, searchBy, []string{enum.Trial.String(), enum.Lead.String()}, commonModel.ComparisonOperatorIn, 4, 3)
	assertSearch(t, searchBy, []string{enum.Trial.String()}, commonModel.ComparisonOperatorNotIn, 4, 3)
	assertSearch(t, searchBy, []string{enum.Lead.String()}, commonModel.ComparisonOperatorNotIn, 4, 2)
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

	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{ID: "org3"})

	require.Equal(t, 3, neo4jtest.GetCountOfNodes(ctx, driver, "Organization"))
	require.Equal(t, 2, neo4jtest.GetCountOfNodes(ctx, driver, "Social"))
	require.Equal(t, 2, neo4jtest.GetCountOfRelationships(ctx, driver, "HAS"))

	searchBy := postgres_entity.ColumnViewTypeOrganizationsSocials

	assertSearch(t, searchBy, "", commonModel.ComparisonOperatorIsEmpty, 3, 1)
	assertSearch(t, searchBy, "", commonModel.ComparisonOperatorIsNotEmpty, 3, 2)
	assertSearch(t, searchBy, "", commonModel.ComparisonOperatorContains, 3, 2)
	assertSearch(t, searchBy, "openline", commonModel.ComparisonOperatorContains, 3, 2)
	assertSearch(t, searchBy, "linkedin", commonModel.ComparisonOperatorContains, 3, 1)
	assertSearch(t, searchBy, "twitter", commonModel.ComparisonOperatorContains, 3, 1)
	assertSearch(t, searchBy, "google", commonModel.ComparisonOperatorContains, 3, 0)

	assertSearch(t, searchBy, "google", commonModel.ComparisonOperatorNotContains, 3, 2)
	assertSearch(t, searchBy, "linkedin", commonModel.ComparisonOperatorNotContains, 3, 1)
	assertSearch(t, searchBy, "twitter", commonModel.ComparisonOperatorNotContains, 3, 1)
}

func TestQueryResolver_UIOrganizationsSearch_FilterByLeadSource(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)

	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{})
	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{LeadSource: ""})
	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{LeadSource: "A"})
	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{LeadSource: "B"})
	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{LeadSource: "AB"})

	require.Equal(t, 5, neo4jtest.GetCountOfNodes(ctx, driver, "Organization"))

	searchBy := postgres_entity.ColumnViewTypeOrganizationsLeadSource

	assertSearch(t, searchBy, "", commonModel.ComparisonOperatorIsEmpty, 5, 2)
	assertSearch(t, searchBy, "", commonModel.ComparisonOperatorIsNotEmpty, 5, 3)
	assertSearch(t, searchBy, []string{"A", "B"}, commonModel.ComparisonOperatorIn, 5, 2)
	assertSearch(t, searchBy, []string{"A"}, commonModel.ComparisonOperatorNotIn, 5, 4)
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

	searchBy := postgres_entity.ColumnViewTypeOrganizationsCreatedDate

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

	searchBy := postgres_entity.ColumnViewTypeOrganizationsEmployeeCount

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

func TestQueryResolver_UIOrganizationsSearch_FilterByContactCount(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)

	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{})
	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{DerivedData: neo4jentity.DerivedData{ContactCount: 2000}})
	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{DerivedData: neo4jentity.DerivedData{ContactCount: 2005}})

	require.Equal(t, 3, neo4jtest.GetCountOfNodes(ctx, driver, "Organization"))

	searchBy := postgres_entity.ColumnViewTypeOrganizationsContactCount

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

	searchBy := postgres_entity.ColumnViewTypeOrganizationsYearFounded

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

	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{})
	org2 := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{})
	org3 := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{})
	org4 := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{})
	org5 := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{})

	neo4jtest.LinkIndustryToOrganization(ctx, driver, org2, "11", "Industry_11")
	neo4jtest.LinkIndustryToOrganization(ctx, driver, org3, "111", "Industry_111")
	neo4jtest.LinkIndustryToOrganization(ctx, driver, org4, "1110", "Industry_1110")
	neo4jtest.LinkIndustryToOrganization(ctx, driver, org5, "1110", "Industry_1110")

	require.Equal(t, 5, neo4jtest.GetCountOfNodes(ctx, driver, "Organization"))
	require.Equal(t, 3, neo4jtest.GetCountOfNodes(ctx, driver, "Industry"))

	searchBy := postgres_entity.ColumnViewTypeOrganizationsIndustry

	assertSearch(t, searchBy, "", commonModel.ComparisonOperatorIsEmpty, 5, 1)
	assertSearch(t, searchBy, "", commonModel.ComparisonOperatorIsNotEmpty, 5, 4)
	assertSearch(t, searchBy, []string{"11", "1110"}, commonModel.ComparisonOperatorIn, 5, 3)
	assertSearch(t, searchBy, []string{"1110"}, commonModel.ComparisonOperatorNotIn, 5, 3)
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

	searchBy := postgres_entity.ColumnViewTypeOrganizationsChurnDate

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

	searchBy := postgres_entity.ColumnViewTypeOrganizationsLtv

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

func TestQueryResolver_UIOrganizationsSearch_FilterByCountry(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)

	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{ID: "org1"})
	neo4jtest.CreateLocation(ctx, driver, tenantName, neo4jentity.LocationEntity{Id: "l1", CountryCodeA2: "CA"})
	neo4jtest.LinkNodes(ctx, driver, "org1", "l1", "ASSOCIATED_WITH")

	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{ID: "org2"})
	neo4jtest.CreateLocation(ctx, driver, tenantName, neo4jentity.LocationEntity{Id: "l2", CountryCodeA2: "US"})
	neo4jtest.LinkNodes(ctx, driver, "org2", "l2", "ASSOCIATED_WITH")

	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{ID: "org3"})
	neo4jtest.CreateLocation(ctx, driver, tenantName, neo4jentity.LocationEntity{Id: "l3"})
	neo4jtest.LinkNodes(ctx, driver, "org3", "l3", "ASSOCIATED_WITH")

	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{ID: "org4"})

	require.Equal(t, 4, neo4jtest.GetCountOfNodes(ctx, driver, "Organization"))
	require.Equal(t, 3, neo4jtest.GetCountOfNodes(ctx, driver, "Location"))
	require.Equal(t, 3, neo4jtest.GetCountOfRelationships(ctx, driver, "ASSOCIATED_WITH"))

	searchBy := postgres_entity.ColumnViewTypeOrganizationsCountry

	assertSearch(t, searchBy, "", commonModel.ComparisonOperatorIsEmpty, 4, 2)
	assertSearch(t, searchBy, "", commonModel.ComparisonOperatorIsNotEmpty, 4, 2)
	assertSearch(t, searchBy, []string{"X"}, commonModel.ComparisonOperatorIn, 4, 0)
	assertSearch(t, searchBy, []string{"US"}, commonModel.ComparisonOperatorIn, 4, 1)
	assertSearch(t, searchBy, []string{"CA"}, commonModel.ComparisonOperatorIn, 4, 1)
	assertSearch(t, searchBy, []string{"US", "CA"}, commonModel.ComparisonOperatorIn, 4, 2)

	assertSearch(t, searchBy, []string{"X"}, commonModel.ComparisonOperatorNotIn, 4, 4)
	assertSearch(t, searchBy, []string{"CA"}, commonModel.ComparisonOperatorNotIn, 4, 3)
}

func TestQueryResolver_UIOrganizationsSearch_FilterByCity(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)

	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{ID: "org1"})
	neo4jtest.CreateLocation(ctx, driver, tenantName, neo4jentity.LocationEntity{Id: "l1", Locality: "C1"})
	neo4jtest.LinkNodes(ctx, driver, "org1", "l1", "ASSOCIATED_WITH")

	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{ID: "org2"})
	neo4jtest.CreateLocation(ctx, driver, tenantName, neo4jentity.LocationEntity{Id: "l2", Locality: "C2"})
	neo4jtest.LinkNodes(ctx, driver, "org2", "l2", "ASSOCIATED_WITH")

	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{ID: "org3"})

	require.Equal(t, 3, neo4jtest.GetCountOfNodes(ctx, driver, "Organization"))
	require.Equal(t, 2, neo4jtest.GetCountOfNodes(ctx, driver, "Location"))
	require.Equal(t, 2, neo4jtest.GetCountOfRelationships(ctx, driver, "ASSOCIATED_WITH"))

	searchBy := postgres_entity.ColumnViewTypeOrganizationsCity

	assertSearch(t, searchBy, "", commonModel.ComparisonOperatorIsEmpty, 3, 1)
	assertSearch(t, searchBy, "", commonModel.ComparisonOperatorIsNotEmpty, 3, 2)
	assertSearch(t, searchBy, "", commonModel.ComparisonOperatorContains, 3, 2)
	assertSearch(t, searchBy, "c", commonModel.ComparisonOperatorContains, 3, 2)
	assertSearch(t, searchBy, "c1", commonModel.ComparisonOperatorContains, 3, 1)
	assertSearch(t, searchBy, "c2", commonModel.ComparisonOperatorContains, 3, 1)
	assertSearch(t, searchBy, "c3", commonModel.ComparisonOperatorContains, 3, 0)

	assertSearch(t, searchBy, "c1", commonModel.ComparisonOperatorNotContains, 3, 1)
	assertSearch(t, searchBy, "c2", commonModel.ComparisonOperatorNotContains, 3, 1)
	assertSearch(t, searchBy, "c3", commonModel.ComparisonOperatorNotContains, 3, 2)
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

	searchBy := postgres_entity.ColumnViewTypeOrganizationsHeadquarters

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

	searchBy := postgres_entity.ColumnViewTypeOrganizationsIsPublic

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
	neo4jtest.CreateTag(ctx, driver, tenantName, neo4jentity.TagEntity{Id: "tag3", Name: "A"})
	neo4jtest.CreateTag(ctx, driver, tenantName, neo4jentity.TagEntity{Id: "tag4", Name: "c"})
	neo4jtest.LinkNodes(ctx, driver, "org2", "tag3", "TAGGED")
	neo4jtest.LinkNodes(ctx, driver, "org2", "tag4", "TAGGED")

	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{ID: "org3"})

	require.Equal(t, 3, neo4jtest.GetCountOfNodes(ctx, driver, "Organization"))
	require.Equal(t, 4, neo4jtest.GetCountOfNodes(ctx, driver, "Tag"))
	require.Equal(t, 4, neo4jtest.GetCountOfRelationships(ctx, driver, "TAGGED"))

	searchBy := postgres_entity.ColumnViewTypeOrganizationsTags

	assertSearch(t, searchBy, "", commonModel.ComparisonOperatorIsEmpty, 3, 1)
	assertSearch(t, searchBy, "", commonModel.ComparisonOperatorIsNotEmpty, 3, 2)
	assertSearch(t, searchBy, []string{"tag3", "tag2", "no-tag"}, commonModel.ComparisonOperatorIn, 3, 2)
	assertSearch(t, searchBy, []string{"tag4"}, commonModel.ComparisonOperatorIn, 3, 1)
	assertSearch(t, searchBy, []string{"notag1", "notag2", "notag3"}, commonModel.ComparisonOperatorIn, 3, 0)

	assertSearch(t, searchBy, []string{"notag1"}, commonModel.ComparisonOperatorNotIn, 3, 3)
	assertSearch(t, searchBy, []string{"tag2", "notag2"}, commonModel.ComparisonOperatorNotIn, 3, 2)
}

func TestQueryResolver_UIOrganizationsSearch_FilterByParentOrganization(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)

	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{ID: "subsidiary1"})
	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{ID: "parent1", Name: "Parent1"})
	neo4jtest.LinkNodes(ctx, driver, "subsidiary1", "parent1", "SUBSIDIARY_OF")

	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{ID: "subsidiary2"})
	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{ID: "parent2", Name: "Parent2"})
	neo4jtest.LinkNodes(ctx, driver, "subsidiary2", "parent2", "SUBSIDIARY_OF")

	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{ID: "org3"})

	require.Equal(t, 5, neo4jtest.GetCountOfNodes(ctx, driver, "Organization"))
	require.Equal(t, 2, neo4jtest.GetCountOfRelationships(ctx, driver, "SUBSIDIARY_OF"))

	searchBy := postgres_entity.ColumnViewTypeOrganizationsParentOrganization

	assertSearch(t, searchBy, "", commonModel.ComparisonOperatorIsEmpty, 5, 3)
	assertSearch(t, searchBy, "", commonModel.ComparisonOperatorIsNotEmpty, 5, 2)
	assertSearch(t, searchBy, "", commonModel.ComparisonOperatorContains, 5, 2)
	assertSearch(t, searchBy, "parent", commonModel.ComparisonOperatorContains, 5, 2)
	assertSearch(t, searchBy, "parent1", commonModel.ComparisonOperatorContains, 5, 1)
	assertSearch(t, searchBy, "parent2", commonModel.ComparisonOperatorContains, 5, 1)
	assertSearch(t, searchBy, "parent3", commonModel.ComparisonOperatorContains, 5, 0)

	assertSearch(t, searchBy, "parent1", commonModel.ComparisonOperatorNotContains, 5, 1)
	assertSearch(t, searchBy, "parent2", commonModel.ComparisonOperatorNotContains, 5, 1)
	assertSearch(t, searchBy, "parent3", commonModel.ComparisonOperatorNotContains, 5, 2)
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

	searchBy := postgres_entity.ColumnViewTypeOrganizationsUpdatedDate

	assertSearch(t, searchBy, "", commonModel.ComparisonOperatorIsEmpty, 3, 0)
	assertSearch(t, searchBy, []time.Time{firstOfDecember, midOfJanuary}, commonModel.ComparisonOperatorBetween, 3, 1)
	assertSearch(t, searchBy, []time.Time{firstOfJanuary, firstOfJanuary}, commonModel.ComparisonOperatorBetween, 3, 1)
	assertSearch(t, searchBy, []time.Time{firstOfJanuary, firstOfFebruary}, commonModel.ComparisonOperatorBetween, 3, 2)
	assertSearch(t, searchBy, []time.Time{firstOfFebruary, firstOfFebruary}, commonModel.ComparisonOperatorBetween, 3, 1)
	assertSearch(t, searchBy, []time.Time{midOfJanuary, firstOfMarch}, commonModel.ComparisonOperatorBetween, 3, 1)
	assertSearch(t, searchBy, []time.Time{midOfFebruary, firstOfMarch}, commonModel.ComparisonOperatorBetween, 3, 0)
	assertSearch(t, searchBy, firstOfJanuary, commonModel.ComparisonOperatorGte, 3, 3)
	assertSearch(t, searchBy, firstOfJanuary, commonModel.ComparisonOperatorGt, 3, 2)
	assertSearch(t, searchBy, firstOfFebruary, commonModel.ComparisonOperatorGte, 3, 2)
	assertSearch(t, searchBy, firstOfFebruary, commonModel.ComparisonOperatorGt, 3, 1)
	assertSearch(t, searchBy, firstOfJanuary, commonModel.ComparisonOperatorLt, 3, 0)
	assertSearch(t, searchBy, firstOfJanuary, commonModel.ComparisonOperatorLte, 3, 1)
	assertSearch(t, searchBy, firstOfFebruary, commonModel.ComparisonOperatorLt, 3, 1)
	assertSearch(t, searchBy, firstOfFebruary, commonModel.ComparisonOperatorLte, 3, 2)
}

func assertSearch(t *testing.T, filterName postgres_entity.ColumnViewType, searchValue any, operator commonModel.ComparisonOperator, totalAvailable int64, totalElements int64) {
	rawResponse, err := c.RawPost(getQuery("organization/ui_organizations_search"),
		client.Var("limit", 10),
		client.Var("filterName", filterName),
		client.Var("filterValue", searchValue),
		client.Var("filterOperation", operator),
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

func TestQueryResolver_UIOrganizationsSearch_SortByName(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)

	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{ID: "empty", Name: ""})
	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{ID: "a", Name: "a"})
	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{ID: "b", Name: "b"})
	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{ID: "ab", Name: "ab"})
	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{ID: "c", Name: "c"})
	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{ID: "aa", Name: "aa"})
	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{ID: "abc", Name: "abc"})
	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{ID: "ba", Name: "ba"})
	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{ID: "z", Name: "z"})

	require.Equal(t, 9, neo4jtest.GetCountOfNodes(ctx, driver, "Organization"))

	expectedAsc := []string{"a", "aa", "ab", "abc", "b", "ba", "c", "z", "empty"}
	expectedDesc := []string{"z", "c", "ba", "b", "abc", "ab", "aa", "a", "empty"}

	verifySortOrder(t, postgres_entity.ColumnViewTypeOrganizationsName, commonModel.SortingDirectionAsc, expectedAsc)
	verifySortOrder(t, postgres_entity.ColumnViewTypeOrganizationsName, commonModel.SortingDirectionDesc, expectedDesc)
}

func TestQueryResolver_UIOrganizationsSearch_SortByWebsite(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)

	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{ID: "empty", Website: ""})
	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{ID: "a", Website: "https://www.a"})
	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{ID: "b", Website: "https://www.b"})
	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{ID: "ab", Website: "https://www.ab"})
	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{ID: "c", Website: "https://www.c"})
	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{ID: "aa", Website: "https://www.aa"})
	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{ID: "abc", Website: "https://www.abc"})
	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{ID: "ba", Website: "https://www.ba"})
	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{ID: "z", Website: "https://www.z"})

	require.Equal(t, 9, neo4jtest.GetCountOfNodes(ctx, driver, "Organization"))

	expectedAsc := []string{"a", "aa", "ab", "abc", "b", "ba", "c", "z", "empty"}
	expectedDesc := []string{"z", "c", "ba", "b", "abc", "ab", "aa", "a", "empty"}

	verifySortOrder(t, postgres_entity.ColumnViewTypeOrganizationsWebsite, commonModel.SortingDirectionAsc, expectedAsc)
	verifySortOrder(t, postgres_entity.ColumnViewTypeOrganizationsWebsite, commonModel.SortingDirectionDesc, expectedDesc)
}

func TestQueryResolver_UIOrganizationsSearch_SortByPrimaryDomain(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)

	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{ID: "org1"})
	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{ID: "org2"})
	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{ID: "org3"})
	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{ID: "org4"})

	neo4jtest.LinkDomainToOrganization(ctx, driver, "org1", neo4jentity.DomainEntity{Domain: "domain1", IsPrimary: utils.BoolPtr(true)})
	neo4jtest.LinkDomainToOrganization(ctx, driver, "org1", neo4jentity.DomainEntity{Domain: "domain1secondary1"})
	neo4jtest.LinkDomainToOrganization(ctx, driver, "org1", neo4jentity.DomainEntity{Domain: "domain1secondary2"})
	neo4jtest.LinkDomainToOrganization(ctx, driver, "org2", neo4jentity.DomainEntity{Domain: "domain2", IsPrimary: utils.BoolPtr(true)})
	neo4jtest.LinkDomainToOrganization(ctx, driver, "org3", neo4jentity.DomainEntity{Domain: "domain3", IsPrimary: utils.BoolPtr(true)})

	require.Equal(t, 4, neo4jtest.GetCountOfNodes(ctx, driver, "Organization"))
	require.Equal(t, 5, neo4jtest.GetCountOfNodes(ctx, driver, "Domain"))
	require.Equal(t, 5, neo4jtest.GetCountOfRelationships(ctx, driver, "HAS_DOMAIN"))

	expectedAsc := []string{"org1", "org2", "org3", "org4"}
	expectedDesc := []string{"org3", "org2", "org1", "org4"}

	verifySortOrder(t, postgres_entity.ColumnViewTypeOrganizationsPrimaryDomains, commonModel.SortingDirectionAsc, expectedAsc)
	verifySortOrder(t, postgres_entity.ColumnViewTypeOrganizationsPrimaryDomains, commonModel.SortingDirectionDesc, expectedDesc)
}

func TestQueryResolver_UIOrganizationsSearch_SortByRelationship(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)

	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{ID: "empty", Relationship: ""})
	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{ID: "a", Relationship: "a"})
	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{ID: "b", Relationship: "b"})
	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{ID: "ab", Relationship: "ab"})
	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{ID: "c", Relationship: "c"})
	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{ID: "aa", Relationship: "aa"})
	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{ID: "abc", Relationship: "abc"})
	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{ID: "ba", Relationship: "ba"})
	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{ID: "z", Relationship: "z"})

	require.Equal(t, 9, neo4jtest.GetCountOfNodes(ctx, driver, "Organization"))

	expectedAsc := []string{"a", "aa", "ab", "abc", "b", "ba", "c", "z", "empty"}
	expectedDesc := []string{"z", "c", "ba", "b", "abc", "ab", "aa", "a", "empty"}

	verifySortOrder(t, postgres_entity.ColumnViewTypeOrganizationsRelationship, commonModel.SortingDirectionAsc, expectedAsc)
	verifySortOrder(t, postgres_entity.ColumnViewTypeOrganizationsRelationship, commonModel.SortingDirectionDesc, expectedDesc)
}

func TestQueryResolver_UIOrganizationsSearch_SortByOnboardingStatus(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)

	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{ID: "empty", OnboardingDetails: neo4jentity.OnboardingDetails{SortingOrder: nil}})
	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{ID: "1", OnboardingDetails: neo4jentity.OnboardingDetails{SortingOrder: utils.Int64Ptr(1)}})
	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{ID: "2", OnboardingDetails: neo4jentity.OnboardingDetails{SortingOrder: utils.Int64Ptr(2)}})

	require.Equal(t, 3, neo4jtest.GetCountOfNodes(ctx, driver, "Organization"))

	expectedAsc := []string{"1", "2", "empty"}
	expectedDesc := []string{"2", "1", "empty"}

	verifySortOrder(t, postgres_entity.ColumnViewTypeOrganizationsOnboardingStatus, commonModel.SortingDirectionAsc, expectedAsc)
	verifySortOrder(t, postgres_entity.ColumnViewTypeOrganizationsOnboardingStatus, commonModel.SortingDirectionDesc, expectedDesc)
}

func TestQueryResolver_UIOrganizationsSearch_SortByRenewalLikelihood(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)

	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{ID: "empty", RenewalSummary: neo4jentity.RenewalSummary{RenewalLikelihoodOrder: nil}})
	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{ID: "1", RenewalSummary: neo4jentity.RenewalSummary{RenewalLikelihoodOrder: utils.Int64Ptr(1)}})
	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{ID: "2", RenewalSummary: neo4jentity.RenewalSummary{RenewalLikelihoodOrder: utils.Int64Ptr(2)}})

	require.Equal(t, 3, neo4jtest.GetCountOfNodes(ctx, driver, "Organization"))

	expectedAsc := []string{"1", "2", "empty"}
	expectedDesc := []string{"2", "1", "empty"}

	verifySortOrder(t, postgres_entity.ColumnViewTypeOrganizationsRenewalLikelihood, commonModel.SortingDirectionAsc, expectedAsc)
	verifySortOrder(t, postgres_entity.ColumnViewTypeOrganizationsRenewalLikelihood, commonModel.SortingDirectionDesc, expectedDesc)
}

func TestQueryResolver_UIOrganizationsSearch_SortByRenewalDate(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)

	firstOfDecember := utils.FirstTimeOfMonth(2023, 12)
	firstOfJanuary := utils.FirstTimeOfMonth(2024, 1)

	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{ID: "empty", RenewalSummary: neo4jentity.RenewalSummary{NextRenewalAt: nil}})
	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{ID: "1", RenewalSummary: neo4jentity.RenewalSummary{NextRenewalAt: &firstOfDecember}})
	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{ID: "2", RenewalSummary: neo4jentity.RenewalSummary{NextRenewalAt: &firstOfJanuary}})

	require.Equal(t, 3, neo4jtest.GetCountOfNodes(ctx, driver, "Organization"))

	expectedAsc := []string{"1", "2", "empty"}
	expectedDesc := []string{"2", "1", "empty"}

	verifySortOrder(t, postgres_entity.ColumnViewTypeOrganizationsRenewalDate, commonModel.SortingDirectionAsc, expectedAsc)
	verifySortOrder(t, postgres_entity.ColumnViewTypeOrganizationsRenewalDate, commonModel.SortingDirectionDesc, expectedDesc)
}

func TestQueryResolver_UIOrganizationsSearch_SortByForecastArr(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)

	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{ID: "empty", RenewalSummary: neo4jentity.RenewalSummary{ArrForecast: nil}})
	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{ID: "1", RenewalSummary: neo4jentity.RenewalSummary{ArrForecast: utils.Float64Ptr(1)}})
	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{ID: "2", RenewalSummary: neo4jentity.RenewalSummary{ArrForecast: utils.Float64Ptr(2)}})

	require.Equal(t, 3, neo4jtest.GetCountOfNodes(ctx, driver, "Organization"))

	expectedAsc := []string{"1", "2", "empty"}
	expectedDesc := []string{"2", "1", "empty"}

	verifySortOrder(t, postgres_entity.ColumnViewTypeOrganizationsForecastArr, commonModel.SortingDirectionAsc, expectedAsc)
	verifySortOrder(t, postgres_entity.ColumnViewTypeOrganizationsForecastArr, commonModel.SortingDirectionDesc, expectedDesc)
}

func TestQueryResolver_UIOrganizationsSearch_SortByOwner(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)

	neo4jtest.CreateUser(ctx, driver, tenantName, neo4jentity.UserEntity{Id: "owner1", FirstName: "owner1"})
	neo4jtest.CreateUser(ctx, driver, tenantName, neo4jentity.UserEntity{Id: "owner2", FirstName: "owner2"})

	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{ID: "1"})
	neo4jtest.LinkNodes(ctx, driver, "owner1", "1", "OWNS")

	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{ID: "2"})
	neo4jtest.LinkNodes(ctx, driver, "owner2", "2", "OWNS")

	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{ID: "empty"})

	require.Equal(t, 3, neo4jtest.GetCountOfNodes(ctx, driver, "Organization"))
	require.Equal(t, 2, neo4jtest.GetCountOfNodes(ctx, driver, "User"))
	require.Equal(t, 2, neo4jtest.GetCountOfRelationships(ctx, driver, "OWNS"))

	require.Equal(t, 3, neo4jtest.GetCountOfNodes(ctx, driver, "Organization"))

	expectedAsc := []string{"1", "2", "empty"}
	expectedDesc := []string{"2", "1", "empty"}

	verifySortOrder(t, postgres_entity.ColumnViewTypeOrganizationsOwner, commonModel.SortingDirectionAsc, expectedAsc)
	verifySortOrder(t, postgres_entity.ColumnViewTypeOrganizationsOwner, commonModel.SortingDirectionDesc, expectedDesc)
}

func TestQueryResolver_UIOrganizationsSearch_SortByLastTouchpoint(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)

	firstOfDecember := utils.FirstTimeOfMonth(2023, 12)
	firstOfJanuary := utils.FirstTimeOfMonth(2024, 1)

	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{ID: "empty", LastTouchpointAt: nil})
	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{ID: "1", LastTouchpointAt: &firstOfDecember})
	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{ID: "2", LastTouchpointAt: &firstOfJanuary})

	require.Equal(t, 3, neo4jtest.GetCountOfNodes(ctx, driver, "Organization"))

	expectedAsc := []string{"1", "2", "empty"}
	expectedDesc := []string{"2", "1", "empty"}

	verifySortOrder(t, postgres_entity.ColumnViewTypeOrganizationsLastTouchpoint, commonModel.SortingDirectionAsc, expectedAsc)
	verifySortOrder(t, postgres_entity.ColumnViewTypeOrganizationsLastTouchpoint, commonModel.SortingDirectionDesc, expectedDesc)
}

func TestQueryResolver_UIOrganizationsSearch_SortByLastTouchpointDate(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)

	firstOfDecember := utils.FirstTimeOfMonth(2023, 12)
	firstOfJanuary := utils.FirstTimeOfMonth(2024, 1)

	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{ID: "empty", LastTouchpointAt: nil})
	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{ID: "1", LastTouchpointAt: &firstOfDecember})
	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{ID: "2", LastTouchpointAt: &firstOfJanuary})

	require.Equal(t, 3, neo4jtest.GetCountOfNodes(ctx, driver, "Organization"))

	expectedAsc := []string{"1", "2", "empty"}
	expectedDesc := []string{"2", "1", "empty"}

	verifySortOrder(t, postgres_entity.ColumnViewTypeOrganizationsLastTouchpointDate, commonModel.SortingDirectionAsc, expectedAsc)
	verifySortOrder(t, postgres_entity.ColumnViewTypeOrganizationsLastTouchpointDate, commonModel.SortingDirectionDesc, expectedDesc)
}

func TestQueryResolver_UIOrganizationsSearch_SortByStage(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)

	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{ID: "empty", Stage: ""})
	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{ID: "a", Stage: "a"})
	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{ID: "b", Stage: "b"})

	require.Equal(t, 3, neo4jtest.GetCountOfNodes(ctx, driver, "Organization"))

	expectedAsc := []string{"a", "b", "empty"}
	expectedDesc := []string{"b", "a", "empty"}

	verifySortOrder(t, postgres_entity.ColumnViewTypeOrganizationsStage, commonModel.SortingDirectionAsc, expectedAsc)
	verifySortOrder(t, postgres_entity.ColumnViewTypeOrganizationsStage, commonModel.SortingDirectionDesc, expectedDesc)
}

func TestQueryResolver_UIOrganizationsSearch_SortByLeadsource(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)

	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{ID: "empty", LeadSource: ""})
	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{ID: "a", LeadSource: "a"})
	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{ID: "b", LeadSource: "b"})

	require.Equal(t, 3, neo4jtest.GetCountOfNodes(ctx, driver, "Organization"))

	expectedAsc := []string{"a", "b", "empty"}
	expectedDesc := []string{"b", "a", "empty"}

	verifySortOrder(t, postgres_entity.ColumnViewTypeOrganizationsLeadSource, commonModel.SortingDirectionAsc, expectedAsc)
	verifySortOrder(t, postgres_entity.ColumnViewTypeOrganizationsLeadSource, commonModel.SortingDirectionDesc, expectedDesc)
}

func TestQueryResolver_UIOrganizationsSearch_SortByCreatedDate(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)

	firstOfDecember := utils.FirstTimeOfMonth(2023, 12)
	firstOfJanuary := utils.FirstTimeOfMonth(2024, 1)

	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{ID: "1", CreatedAt: firstOfDecember})
	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{ID: "2", CreatedAt: firstOfJanuary})

	require.Equal(t, 2, neo4jtest.GetCountOfNodes(ctx, driver, "Organization"))

	expectedAsc := []string{"1", "2"}
	expectedDesc := []string{"2", "1"}

	verifySortOrder(t, postgres_entity.ColumnViewTypeOrganizationsCreatedDate, commonModel.SortingDirectionAsc, expectedAsc)
	verifySortOrder(t, postgres_entity.ColumnViewTypeOrganizationsCreatedDate, commonModel.SortingDirectionDesc, expectedDesc)
}

func TestQueryResolver_UIOrganizationsSearch_SortByEmployeeCount(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)

	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{ID: "1", Employees: 1})
	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{ID: "2", Employees: 2})

	require.Equal(t, 2, neo4jtest.GetCountOfNodes(ctx, driver, "Organization"))

	expectedAsc := []string{"1", "2"}
	expectedDesc := []string{"2", "1"}

	verifySortOrder(t, postgres_entity.ColumnViewTypeOrganizationsEmployeeCount, commonModel.SortingDirectionAsc, expectedAsc)
	verifySortOrder(t, postgres_entity.ColumnViewTypeOrganizationsEmployeeCount, commonModel.SortingDirectionDesc, expectedDesc)
}

func TestQueryResolver_UIOrganizationsSearch_SortByContactCount(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)

	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{ID: "1", DerivedData: neo4jentity.DerivedData{ContactCount: 1}})
	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{ID: "2", DerivedData: neo4jentity.DerivedData{ContactCount: 2}})

	require.Equal(t, 2, neo4jtest.GetCountOfNodes(ctx, driver, "Organization"))

	expectedAsc := []string{"1", "2"}
	expectedDesc := []string{"2", "1"}

	verifySortOrder(t, postgres_entity.ColumnViewTypeOrganizationsContactCount, commonModel.SortingDirectionAsc, expectedAsc)
	verifySortOrder(t, postgres_entity.ColumnViewTypeOrganizationsContactCount, commonModel.SortingDirectionDesc, expectedDesc)
}

func TestQueryResolver_UIOrganizationsSearch_SortByYearFounded(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)

	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{ID: "empty", YearFounded: nil})
	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{ID: "1", YearFounded: utils.Int64Ptr(1)})
	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{ID: "2", YearFounded: utils.Int64Ptr(2)})

	require.Equal(t, 3, neo4jtest.GetCountOfNodes(ctx, driver, "Organization"))

	expectedAsc := []string{"1", "2", "empty"}
	expectedDesc := []string{"2", "1", "empty"}

	verifySortOrder(t, postgres_entity.ColumnViewTypeOrganizationsYearFounded, commonModel.SortingDirectionAsc, expectedAsc)
	verifySortOrder(t, postgres_entity.ColumnViewTypeOrganizationsYearFounded, commonModel.SortingDirectionDesc, expectedDesc)
}

func TestQueryResolver_UIOrganizationsSearch_SortByIndustry(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)

	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{ID: "empty"})
	org2 := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{ID: "a"})
	org3 := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{ID: "b"})

	neo4jtest.LinkIndustryToOrganization(ctx, driver, org2, "1", "A")
	neo4jtest.LinkIndustryToOrganization(ctx, driver, org3, "2", "B")

	require.Equal(t, 3, neo4jtest.GetCountOfNodes(ctx, driver, commonModel.NodeLabelOrganization))
	require.Equal(t, 2, neo4jtest.GetCountOfNodes(ctx, driver, commonModel.NodeLabelIndustry))

	expectedAsc := []string{"a", "b", "empty"}
	expectedDesc := []string{"b", "a", "empty"}

	verifySortOrder(t, postgres_entity.ColumnViewTypeOrganizationsIndustry, commonModel.SortingDirectionAsc, expectedAsc)
	verifySortOrder(t, postgres_entity.ColumnViewTypeOrganizationsIndustry, commonModel.SortingDirectionDesc, expectedDesc)
}

func TestQueryResolver_UIOrganizationsSearch_SortByChurnDate(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)

	firstOfDecember := utils.FirstTimeOfMonth(2023, 12)
	firstOfJanuary := utils.FirstTimeOfMonth(2024, 1)

	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{ID: "empty", DerivedData: neo4jentity.DerivedData{ChurnedAt: nil}})
	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{ID: "1", DerivedData: neo4jentity.DerivedData{ChurnedAt: &firstOfDecember}})
	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{ID: "2", DerivedData: neo4jentity.DerivedData{ChurnedAt: &firstOfJanuary}})

	require.Equal(t, 3, neo4jtest.GetCountOfNodes(ctx, driver, "Organization"))

	expectedAsc := []string{"1", "2", "empty"}
	expectedDesc := []string{"2", "1", "empty"}

	verifySortOrder(t, postgres_entity.ColumnViewTypeOrganizationsChurnDate, commonModel.SortingDirectionAsc, expectedAsc)
	verifySortOrder(t, postgres_entity.ColumnViewTypeOrganizationsChurnDate, commonModel.SortingDirectionDesc, expectedDesc)
}

func TestQueryResolver_UIOrganizationsSearch_SortLtv(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)

	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{ID: "a", DerivedData: neo4jentity.DerivedData{Ltv: float64(1)}})
	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{ID: "b", DerivedData: neo4jentity.DerivedData{Ltv: float64(2)}})

	require.Equal(t, 2, neo4jtest.GetCountOfNodes(ctx, driver, "Organization"))

	expectedAsc := []string{"a", "b"}
	expectedDesc := []string{"b", "a"}

	verifySortOrder(t, postgres_entity.ColumnViewTypeOrganizationsLtv, commonModel.SortingDirectionAsc, expectedAsc)
	verifySortOrder(t, postgres_entity.ColumnViewTypeOrganizationsLtv, commonModel.SortingDirectionDesc, expectedDesc)
}

func TestQueryResolver_UIOrganizationsSearch_SortByCountry(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)

	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{ID: "1"})
	neo4jtest.CreateLocation(ctx, driver, tenantName, neo4jentity.LocationEntity{Id: "l1", Country: "C1"})
	neo4jtest.LinkNodes(ctx, driver, "1", "l1", "ASSOCIATED_WITH")

	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{ID: "2"})
	neo4jtest.CreateLocation(ctx, driver, tenantName, neo4jentity.LocationEntity{Id: "l2", Country: "C2"})
	neo4jtest.LinkNodes(ctx, driver, "2", "l2", "ASSOCIATED_WITH")

	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{ID: "empty"})

	require.Equal(t, 3, neo4jtest.GetCountOfNodes(ctx, driver, "Organization"))
	require.Equal(t, 2, neo4jtest.GetCountOfNodes(ctx, driver, "Location"))
	require.Equal(t, 2, neo4jtest.GetCountOfRelationships(ctx, driver, "ASSOCIATED_WITH"))

	require.Equal(t, 3, neo4jtest.GetCountOfNodes(ctx, driver, "Organization"))

	expectedAsc := []string{"1", "2", "empty"}
	expectedDesc := []string{"2", "1", "empty"}

	verifySortOrder(t, postgres_entity.ColumnViewTypeOrganizationsCountry, commonModel.SortingDirectionAsc, expectedAsc)
	verifySortOrder(t, postgres_entity.ColumnViewTypeOrganizationsCountry, commonModel.SortingDirectionDesc, expectedDesc)
}

func TestQueryResolver_UIOrganizationsSearch_SortByCity(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)

	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{ID: "1"})
	neo4jtest.CreateLocation(ctx, driver, tenantName, neo4jentity.LocationEntity{Id: "l1", Locality: "C1"})
	neo4jtest.LinkNodes(ctx, driver, "1", "l1", "ASSOCIATED_WITH")

	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{ID: "2"})
	neo4jtest.CreateLocation(ctx, driver, tenantName, neo4jentity.LocationEntity{Id: "l2", Locality: "C2"})
	neo4jtest.LinkNodes(ctx, driver, "2", "l2", "ASSOCIATED_WITH")

	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{ID: "empty"})

	require.Equal(t, 3, neo4jtest.GetCountOfNodes(ctx, driver, "Organization"))
	require.Equal(t, 2, neo4jtest.GetCountOfNodes(ctx, driver, "Location"))
	require.Equal(t, 2, neo4jtest.GetCountOfRelationships(ctx, driver, "ASSOCIATED_WITH"))

	require.Equal(t, 3, neo4jtest.GetCountOfNodes(ctx, driver, "Organization"))

	expectedAsc := []string{"1", "2", "empty"}
	expectedDesc := []string{"2", "1", "empty"}

	verifySortOrder(t, postgres_entity.ColumnViewTypeOrganizationsCity, commonModel.SortingDirectionAsc, expectedAsc)
	verifySortOrder(t, postgres_entity.ColumnViewTypeOrganizationsCity, commonModel.SortingDirectionDesc, expectedDesc)
}

func TestQueryResolver_UIOrganizationsSearch_SortByIsPublic(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)

	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{ID: "empty"}) // false by default
	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{ID: "1", IsPublic: true})
	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{ID: "2", IsPublic: false})

	require.Equal(t, 3, neo4jtest.GetCountOfNodes(ctx, driver, "Organization"))

	expectedAsc := []string{"1", "2", "empty"}
	expectedDesc := []string{"2", "empty", "1"}

	verifySortOrder(t, postgres_entity.ColumnViewTypeOrganizationsIsPublic, commonModel.SortingDirectionAsc, expectedAsc)
	verifySortOrder(t, postgres_entity.ColumnViewTypeOrganizationsIsPublic, commonModel.SortingDirectionDesc, expectedDesc)
}

func TestQueryResolver_UIOrganizationsSearch_SortByParentOrganization(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)

	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{ID: "subsidiary1"})
	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{ID: "parent1", Name: "Parent1"})
	neo4jtest.LinkNodes(ctx, driver, "subsidiary1", "parent1", "SUBSIDIARY_OF")

	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{ID: "subsidiary2"})
	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{ID: "parent2", Name: "Parent2"})
	neo4jtest.LinkNodes(ctx, driver, "subsidiary2", "parent2", "SUBSIDIARY_OF")

	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{ID: "empty"})

	require.Equal(t, 5, neo4jtest.GetCountOfNodes(ctx, driver, "Organization"))

	sortAsc := assertSort(t, postgres_entity.ColumnViewTypeOrganizationsParentOrganization, commonModel.SortingDirectionAsc)
	assert.Equal(t, "subsidiary1", sortAsc[0])
	assert.Equal(t, "subsidiary2", sortAsc[1])

	sortDesc := assertSort(t, postgres_entity.ColumnViewTypeOrganizationsParentOrganization, commonModel.SortingDirectionDesc)
	assert.Equal(t, "subsidiary2", sortDesc[0])
	assert.Equal(t, "subsidiary1", sortDesc[1])
}

func TestQueryResolver_UIOrganizationsSearch_SortByUpdatedDate(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)

	firstOfDecember := utils.FirstTimeOfMonth(2023, 12)
	firstOfJanuary := utils.FirstTimeOfMonth(2024, 1)

	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{ID: "empty"}) // default to now
	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{ID: "1", UpdatedAt: firstOfDecember})
	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{ID: "2", UpdatedAt: firstOfJanuary})

	require.Equal(t, 3, neo4jtest.GetCountOfNodes(ctx, driver, "Organization"))

	expectedAsc := []string{"1", "2", "empty"}
	expectedDesc := []string{"empty", "2", "1"}

	verifySortOrder(t, postgres_entity.ColumnViewTypeOrganizationsUpdatedDate, commonModel.SortingDirectionAsc, expectedAsc)
	verifySortOrder(t, postgres_entity.ColumnViewTypeOrganizationsUpdatedDate, commonModel.SortingDirectionDesc, expectedDesc)
}

func verifySortOrder(t *testing.T, sortBy postgres_entity.ColumnViewType, direction commonModel.SortingDirection, expectedOrder []string) {
	sortedResult := assertSort(t, sortBy, direction)
	assert.Equal(t, len(expectedOrder), len(sortedResult), "Mismatch in result length")
	for i, expected := range expectedOrder {
		assert.Equal(t, expected, sortedResult[i], "Mismatch at index %d", i)
	}
}

func assertSort(t *testing.T, sortBy postgres_entity.ColumnViewType, sortDirection commonModel.SortingDirection) []string {
	rawResponse, err := c.RawPost(getQuery("organization/ui_organizations_sort"),
		client.Var("limit", 10),
		client.Var("sortByField", string(sortBy)),
		client.Var("sortByDirection", sortDirection),
	)
	assertRawResponseSuccess(t, rawResponse, err)

	var organizations struct {
		Ui_Organizations_Search model.OrganizationSearchResult
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &organizations)
	require.Nil(t, err)
	require.NotNil(t, organizations)

	return organizations.Ui_Organizations_Search.Ids
}
