package customer_base

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var apiKey = os.Getenv("TESTIFY_TENANT_API_KEY")

func init() {
	if apiKey == "" {
		panic("TESTIFY_TENANT_API_KEY environment variable is required")
	}
}

type OrganizationPayload struct {
	CustomID     string `json:"customId"`
	ICPFit       bool   `json:"icpFit"`
	LeadSource   string `json:"leadSource"`
	LinkedinURL  string `json:"linkedinUrl"`
	Name         string `json:"name"`
	Relationship string `json:"relationship"`
	Website      string `json:"website"`
}

type Organization struct {
	Status       string `json:"status"`
	Organization struct {
		ID            string   `json:"id"`
		CustomId      string   `json:"customId"`
		CosId         string   `json:"cosId"`
		Name          string   `json:"name"`
		Website       string   `json:"website"`
		LeadSource    string   `json:"leadSource"`
		LinkedinURL   string   `json:"linkedinURL"`
		Relationship  string   `json:"relationship"`
		Stage         string   `json:"stage"`
		IcpFit        bool     `json:"icpFit"`
		Domains       []string `json:"domains"`
		ExternalLinks []struct {
			Name    string `json:"name"`
			ID      string `json:"id"`
			Primary bool   `json:"primary"`
		} `json:"organization"`
	}
}

type ExternalSystemResponse struct {
	Status       string `json:"status"`
	Organization struct {
		OrganizationId string `json:"organizationId"`
		ExternalSystem string `json:"externalSystem"`
		ExternalId     string `json:"externalId"`
		Primary        bool   `json:"primary"`
	} `json:"organization"`
}

// Helper functions come before test functions
func NewTestOrganizationPayload() OrganizationPayload {
	uniqueID := uuid.New().String()
	return OrganizationPayload{
		CustomID:     strings.ReplaceAll(uniqueID, "-", ""),
		ICPFit:       true,
		LeadSource:   "Web Search",
		LinkedinURL:  "https://linkedin.com/in/" + uniqueID,
		Name:         "Testify-" + uniqueID,
		Relationship: "CUSTOMER",
		Website:      "https://testify" + uniqueID + ".ai",
	}
}

func createOrganization(t *testing.T) (struct {
	Organization struct {
		ID string `json:"id"`
	} `json:"organization"`
	Status string `json:"status"`
}, OrganizationPayload) {
	payload := NewTestOrganizationPayload()

	jsonData, err := json.Marshal(payload)
	assert.NoError(t, err)

	url := "https://api.customeros.ai/customerbase/v1/organizations"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	assert.NoError(t, err)

	req.Header.Set("X-CUSTOMER-OS-API-KEY", apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var result struct {
		Organization struct {
			ID string `json:"id"`
		} `json:"organization"`
		Status string `json:"status"`
	}

	err = json.NewDecoder(resp.Body).Decode(&result)
	assert.NoError(t, err)
	return result, payload
}

func getOrganization(t *testing.T, organizationID string) Organization {
	url := fmt.Sprintf("https://api.customeros.ai/customerbase/v1/organizations/%s", organizationID)
	req, err := http.NewRequest("GET", url, nil)
	assert.NoError(t, err)

	req.Header.Set("X-CUSTOMER-OS-API-KEY", apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var org Organization
	err = json.NewDecoder(resp.Body).Decode(&org)
	assert.NoError(t, err)
	return org
}

func setPrimaryLink(t *testing.T, org Organization, payload map[string]string) ExternalSystemResponse {
	url := fmt.Sprintf("https://api.customeros.ai/customerbase/v1/organizations/%s/links/stripe/primary", org.Organization.ID)
	jsonData, err := json.Marshal(payload)
	assert.NoError(t, err)
	setPrimaryLinkReq, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
	assert.NoError(t, err)
	setPrimaryLinkReq.Header.Set("X-CUSTOMER-OS-API-KEY", apiKey)
	setPrimaryLinkReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	setPrimaryLinkResponse, err := client.Do(setPrimaryLinkReq)
	assert.NoError(t, err)
	defer setPrimaryLinkResponse.Body.Close()

	assert.Equal(t, http.StatusOK, setPrimaryLinkResponse.StatusCode)

	var setPrimaryLinkResult ExternalSystemResponse
	err = json.NewDecoder(setPrimaryLinkResponse.Body).Decode(&setPrimaryLinkResult)
	assert.NoError(t, err)
	return setPrimaryLinkResult
}

func TestCreateOrganization(t *testing.T) {
	orgId, _ := createOrganization(t)
	assert.Equal(t, "success", orgId.Status)
	assert.NotEmpty(t, orgId.Organization.ID)
}

func TestGetOrganization(t *testing.T) {
	orgId, expectedOrg := createOrganization(t)
	assert.Equal(t, "success", orgId.Status)
	assert.NotEmpty(t, orgId.Organization.ID)

	org := getOrganization(t, orgId.Organization.ID)
	assert.Equal(t, "success", org.Status)
	assert.Equal(t, orgId.Organization.ID, org.Organization.ID)
	assert.Equal(t, expectedOrg.ICPFit, org.Organization.IcpFit)
	assert.Equal(t, expectedOrg.LeadSource, org.Organization.LeadSource)
	assert.Equal(t, expectedOrg.LinkedinURL, org.Organization.LinkedinURL)
	assert.Equal(t, expectedOrg.Name, org.Organization.Name)
	assert.Equal(t, expectedOrg.Relationship, org.Organization.Relationship)
	assert.Equal(t, "ONBOARDING", org.Organization.Stage)
	assert.Equal(t, expectedOrg.Website, org.Organization.Website)
	assert.NotEmpty(t, org.Organization.CosId)
	assert.Nil(t, org.Organization.Domains)
	assert.Nil(t, org.Organization.ExternalLinks)
}

func TestSetPrimaryLink(t *testing.T) {
	orgId, _ := createOrganization(t)
	assert.Equal(t, "success", orgId.Status)
	assert.NotEmpty(t, orgId.Organization.ID)

	org := getOrganization(t, orgId.Organization.ID)
	firstPayload := map[string]string{"externalId": "stripe-1234"}
	initialPrimaryLink := setPrimaryLink(t, org, firstPayload)
	assert.Equal(t, "success", initialPrimaryLink.Status)
	assert.Equal(t, "stripe", initialPrimaryLink.Organization.ExternalSystem)
	assert.Equal(t, "stripe-1234", initialPrimaryLink.Organization.ExternalId)
	assert.NotEmpty(t, initialPrimaryLink.Organization.Primary)
	assert.NotEmpty(t, initialPrimaryLink.Organization.OrganizationId)

	secondPayload := map[string]string{"externalId": "stripe-5678"}
	updatedPrimaryLink := setPrimaryLink(t, org, secondPayload)
	assert.Equal(t, "success", updatedPrimaryLink.Status)
	assert.Equal(t, "stripe", updatedPrimaryLink.Organization.ExternalSystem)
	assert.Equal(t, "stripe-5678", updatedPrimaryLink.Organization.ExternalId)
	assert.NotEmpty(t, updatedPrimaryLink.Organization.Primary)
	assert.NotEmpty(t, updatedPrimaryLink.Organization.OrganizationId)
}
