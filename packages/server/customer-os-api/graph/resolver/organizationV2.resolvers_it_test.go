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
}
