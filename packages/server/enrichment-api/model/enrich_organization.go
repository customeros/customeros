package model

import (
	postgresEntity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"strings"
)

type EnrichOrganizationRequest struct {
	Domain      string `json:"domain"`
	LinkedinUrl string `json:"linkedinUrl"`
}

func (e *EnrichOrganizationRequest) Normalize() {
	e.LinkedinUrl = strings.TrimSpace(e.LinkedinUrl)
	e.Domain = strings.TrimSpace(e.Domain)
}

type EnrichOrganizationResponse struct {
	Status              string                         `json:"status"`
	Message             string                         `json:"message,omitempty"`
	Success             bool                           `json:"success"`
	PrimaryEnrichSource string                         `json:"primaryEnrichSource"`
	Data                EnrichOrganizationResponseData `json:"data"`
}

type EnrichOrganizationResponseData struct {
	Name      string `json:"name"`
	Domain    string `json:"domain"`
	Employees int64  `json:"employees"`
}

type EnrichOrganizationScrapinResponse struct {
	Status            string                              `json:"status"`
	Message           string                              `json:"message,omitempty"`
	RecordId          uint64                              `json:"recordId,omitempty"`
	OrganizationFound bool                                `json:"organizationFound"`
	Data              *postgresEntity.ScrapInResponseBody `json:"data,omitempty"`
}

//
//
//	if brandfetch.Company.FoundedYear > 0 {
//		organizationFieldsMask = append(organizationFieldsMask, organizationpb.OrganizationMaskField_ORGANIZATION_PROPERTY_YEAR_FOUNDED)
//		updateGrpcRequest.YearFounded = &brandfetch.Company.FoundedYear
//	}
//	if brandfetch.Description != "" {
//		organizationFieldsMask = append(organizationFieldsMask, organizationpb.OrganizationMaskField_ORGANIZATION_PROPERTY_VALUE_PROPOSITION)
//		updateGrpcRequest.ValueProposition = brandfetch.Description
//	}
//	if brandfetch.LongDescription != "" {
//		organizationFieldsMask = append(organizationFieldsMask, organizationpb.OrganizationMaskField_ORGANIZATION_PROPERTY_DESCRIPTION)
//		updateGrpcRequest.Description = brandfetch.LongDescription
//	}
//
//	// Set public indicator
//	if brandfetch.Company.Kind != "" {
//		organizationFieldsMask = append(organizationFieldsMask, organizationpb.OrganizationMaskField_ORGANIZATION_PROPERTY_IS_PUBLIC)
//		if brandfetch.Company.Kind == "PUBLIC_COMPANY" {
//			updateGrpcRequest.IsPublic = true
//		} else {
//			updateGrpcRequest.IsPublic = false
//		}
//	}
//
//	// Set company name
//	if brandfetch.Name != "" {
//		organizationFieldsMask = append(organizationFieldsMask, organizationpb.OrganizationMaskField_ORGANIZATION_PROPERTY_NAME)
//		updateGrpcRequest.Name = brandfetch.Name
//	} else if brandfetch.Domain != "" {
//		organizationFieldsMask = append(organizationFieldsMask, organizationpb.OrganizationMaskField_ORGANIZATION_PROPERTY_NAME)
//		updateGrpcRequest.Name = brandfetch.Domain
//	}
//
//	if brandfetch.Domain != "" && organizationEntity.Website == "" {
//		organizationFieldsMask = append(organizationFieldsMask, organizationpb.OrganizationMaskField_ORGANIZATION_PROPERTY_WEBSITE)
//		updateGrpcRequest.Website = brandfetch.Domain
//	}
//
//	// Set company logo and icon urls
//	logoUrl := ""
//	iconUrl := ""
//	if len(brandfetch.Logos) > 0 {
//		for _, logo := range brandfetch.Logos {
//			if logo.Type == "icon" {
//				iconUrl = logo.Formats[0].Src
//			} else if logo.Type == "symbol" && iconUrl == "" {
//				iconUrl = logo.Formats[0].Src
//			} else if logo.Type == "logo" {
//				logoUrl = logo.Formats[0].Src
//			} else if logo.Type == "other" && logoUrl == "" {
//				logoUrl = logo.Formats[0].Src
//			}
//		}
//	}
//	if logoUrl != "" {
//		organizationFieldsMask = append(organizationFieldsMask, organizationpb.OrganizationMaskField_ORGANIZATION_PROPERTY_LOGO_URL)
//		updateGrpcRequest.LogoUrl = logoUrl
//	}
//	if iconUrl != "" {
//		organizationFieldsMask = append(organizationFieldsMask, organizationpb.OrganizationMaskField_ORGANIZATION_PROPERTY_ICON_URL)
//		updateGrpcRequest.IconUrl = iconUrl
//	}
//
//	// set industry
//	industryName := ""
//	industryMaxScore := float64(0)
//	if len(brandfetch.Company.Industries) > 0 {
//		for _, industry := range brandfetch.Company.Industries {
//			if industry.Name != "" && industry.Score > industryMaxScore {
//				industryName = industry.Name
//				industryMaxScore = industry.Score
//			}
//		}
//	}
//	if industryName != "" {
//		organizationFieldsMask = append(organizationFieldsMask, organizationpb.OrganizationMaskField_ORGANIZATION_PROPERTY_INDUSTRY)
//		updateGrpcRequest.Industry = industryName
//	}
//
//	updateGrpcRequest.FieldsMask = organizationFieldsMask
//	tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
//	_, err := subscriptions.CallEventsPlatformGRPCWithRetry[*organizationpb.OrganizationIdGrpcResponse](func() (*organizationpb.OrganizationIdGrpcResponse, error) {
//		return h.grpcClients.OrganizationClient.UpdateOrganization(ctx, &updateGrpcRequest)
//	})
//	if err != nil {
//		tracing.TraceErr(span, err)
//		h.log.Errorf("Error updating organization: %s", err.Error())
//	}
//
//	//add location
//	if !brandfetch.Company.LocationIsEmpty() {
//		_, err := subscriptions.CallEventsPlatformGRPCWithRetry[*locationpb.LocationIdGrpcResponse](func() (*locationpb.LocationIdGrpcResponse, error) {
//			return h.grpcClients.OrganizationClient.AddLocation(ctx, &organizationpb.OrganizationAddLocationGrpcRequest{
//				Tenant:         tenant,
//				OrganizationId: organizationEntity.ID,
//				LocationDetails: &locationpb.LocationDetails{
//					Country:       brandfetch.Company.Location.Country,
//					CountryCodeA2: brandfetch.Company.Location.CountryCodeA2,
//					Locality:      brandfetch.Company.Location.City,
//					Region:        brandfetch.Company.Location.State,
//				},
//				SourceFields: &commonpb.SourceFields{
//					AppSource: constants.AppBrandfetch,
//					Source:    constants.SourceOpenline,
//				},
//			})
//		})
//		if err != nil {
//			tracing.TraceErr(span, err)
//		}
//	}
//
//	//add socials
//	for _, link := range brandfetch.Links {
//		h.addSocial(ctx, organizationEntity.ID, tenant, link.Url, constants.AppBrandfetch)
//	}
