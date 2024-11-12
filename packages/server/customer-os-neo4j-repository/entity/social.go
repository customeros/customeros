package entity

import (
	"strings"
	"time"
)

type SocialProperty string

const (
	SocialPropertyId             SocialProperty = "id"
	SocialPropertyCreatedAt      SocialProperty = "createdAt"
	SocialPropertyUpdatedAt      SocialProperty = "updatedAt"
	SocialPropertySource         SocialProperty = "source"
	SocialPropertyAppSource      SocialProperty = "appSource"
	SocialPropertyUrl            SocialProperty = "url"
	SocialPropertyAlias          SocialProperty = "alias"
	SocialPropertyFollowersCount SocialProperty = "followersCount"
	SocialPropertyExternalId     SocialProperty = "externalId"
)

type SocialEntity struct {
	DataLoaderKey
	Id             string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	Source         DataSource
	AppSource      string
	Url            string
	Alias          string
	FollowersCount int64
	ExternalId     string
}

type SocialEntities []SocialEntity

func (s SocialEntity) IsLinkedin() bool {
	return strings.Contains(s.Url, "linkedin.com")
}

func (s SocialEntity) ExtractLinkedinCompanyIdentifierFromUrl() string {
	if !s.IsLinkedin() {
		return ""
	}

	identifier := s.Url

	// trim trailing / from url
	identifier = strings.TrimSuffix(s.Url, "/")

	// remove all chars before linkedin.com/company
	if i := strings.Index(identifier, "linkedin.com/company"); i != -1 {
		identifier = identifier[i:]
	}

	if strings.HasPrefix(identifier, "linkedin.com/company") {
		// get last part of url
		parts := strings.Split(identifier, "/")
		identifier = parts[len(parts)-1]
		if identifier == "company" {
			identifier = ""
		}
	}
	return identifier
}

func (s SocialEntity) ExtractLinkedinPersonIdentifierFromUrl() string {
	if !s.IsLinkedin() {
		return ""
	}

	identifier := s.Url

	// trim trailing / from url
	identifier = strings.TrimSuffix(s.Url, "/")

	// remove all chars before linkedin.com/company
	if i := strings.Index(identifier, "linkedin.com/in"); i != -1 {
		identifier = identifier[i:]
	}

	if strings.HasPrefix(identifier, "linkedin.com/in") {
		// get last part of url
		parts := strings.Split(identifier, "/")
		identifier = parts[len(parts)-1]
		if identifier == "in" {
			identifier = ""
		}
	}
	return identifier
}
