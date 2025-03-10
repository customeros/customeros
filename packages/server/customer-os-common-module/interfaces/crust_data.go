package interfaces

import (
	"context"
)

// CrustDataResponse represents the response structure from the Crust Data API
type CrustDataResponse struct {
	Profiles          []CrustDataProfile `json:"profiles"`
	TotalDisplayCount string             `json:"total_display_count"`
}

// CrustDataProfile represents a single person's profile in the Crust Data response
type CrustDataProfile struct {
	Name                             string                `json:"name"`
	Location                         string                `json:"location"`
	LinkedinProfileURL               string                `json:"linkedin_profile_url"`
	LinkedinProfileURN               string                `json:"linkedin_profile_urn"`
	DefaultPositionTitle             string                `json:"default_position_title"`
	DefaultPositionCompanyLinkedinID string                `json:"default_position_company_linkedin_id"`
	DefaultPositionIsDecisionMaker   bool                  `json:"default_position_is_decision_maker"`
	FlagshipProfileURL               string                `json:"flagship_profile_url"`
	ProfilePictureURL                string                `json:"profile_picture_url"`
	Headline                         string                `json:"headline"`
	Summary                          *string               `json:"summary"`
	NumOfConnections                 int                   `json:"num_of_connections"`
	RelatedColleagueCompanyID        *int64                `json:"related_colleague_company_id"`
	Skills                           []string              `json:"skills"`
	Employer                         []CrustDataEmployment `json:"employer"`
	EducationBackground              []CrustDataEducation  `json:"education_background"`
	Emails                           []string              `json:"emails"`
	Websites                         []string              `json:"websites"`
	TwitterHandle                    *string               `json:"twitter_handle"`
	Languages                        []string              `json:"languages"`
	Pronoun                          *string               `json:"pronoun"`
	QueryPersonLinkedinURN           string                `json:"query_person_linkedin_urn"`
	LinkedinSlugOrURNs               []string              `json:"linkedin_slug_or_urns"`
	CurrentTitle                     string                `json:"current_title"`
}

// CrustDataEmployment represents employment information in a profile
type CrustDataEmployment struct {
	Title             string               `json:"title"`
	CompanyName       string               `json:"company_name"`
	CompanyLinkedinID *string              `json:"company_linkedin_id"`
	CompanyLogoURL    *string              `json:"company_logo_url"`
	StartDate         *string              `json:"start_date"`
	EndDate           *string              `json:"end_date"`
	PositionID        int64                `json:"position_id"`
	Description       *string              `json:"description"`
	Location          *string              `json:"location"`
	RichMedia         []CrustDataRichMedia `json:"rich_media"`
}

// CrustDataEducation represents education information in a profile
type CrustDataEducation struct {
	DegreeName           *string `json:"degree_name"`
	InstituteName        string  `json:"institute_name"`
	FieldOfStudy         string  `json:"field_of_study"`
	StartDate            *string `json:"start_date"`
	EndDate              *string `json:"end_date"`
	InstituteLinkedinID  string  `json:"institute_linkedin_id"`
	InstituteLinkedinURL string  `json:"institute_linkedin_url"`
	InstituteLogoURL     string  `json:"institute_logo_url"`
}

// CrustDataRichMedia represents media content in employment history
type CrustDataRichMedia struct {
	Title        string `json:"title"`
	Description  string `json:"description"`
	URL          string `json:"url"`
	ThumbnailURL string `json:"thumbnail_url"`
}

// CrustDataService defines the interface for interacting with Crust Data
type CrustDataService interface {
	SearchPeople(ctx context.Context, companyDomain string, jobTitles []string) (*CrustDataResponse, error)
}
