package entity

import "time"

type EnrichDetailsBrandfetch struct {
	ID        uint64    `gorm:"primary_key;autoIncrement:true" json:"id"`
	Domain    string    `gorm:"column:domain;type:varchar(255);DEFAULT:'';NOT NULL" json:"domain"`
	Data      string    `gorm:"column:data;type:text;DEFAULT:'';NOT NULL" json:"data"`
	CreatedAt time.Time `gorm:"column:created_at;type:timestamp;DEFAULT:current_timestamp" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:timestamp;;DEFAULT:current_timestamp" json:"updatedAt"`
	Success   bool      `gorm:"column:success;type:boolean;DEFAULT:false" json:"success"`
}

func (EnrichDetailsBrandfetch) TableName() string {
	return "enrich_details_brandfetch"
}

type BrandfetchResponseBody struct {
	Message         string            `json:"message,omitempty"`
	Id              string            `json:"id,omitempty"`
	Name            string            `json:"name,omitempty"`
	Domain          string            `json:"domain,omitempty"`
	Claimed         bool              `json:"claimed"`
	Description     string            `json:"description,omitempty"`
	LongDescription string            `json:"longDescription,omitempty"`
	Links           []BrandfetchLink  `json:"links,omitempty"`
	Logos           []BranfetchLogo   `json:"logos,omitempty"`
	QualityScore    float64           `json:"qualityScore,omitempty"`
	Company         BrandfetchCompany `json:"company,omitempty"`
	IsNsfw          bool              `json:"isNsfw"`
}

type BrandfetchLink struct {
	Name string `json:"name,omitempty"`
	Url  string `json:"url,omitempty"`
}

type BranfetchLogo struct {
	Theme   string                `json:"theme,omitempty"`
	Type    string                `json:"type,omitempty"`
	Formats []BranfetchLogoFormat `json:"formats,omitempty"`
}

type BranfetchLogoFormat struct {
	Src        string `json:"src,omitempty"`
	Background string `json:"background,omitempty"`
	Format     string `json:"format,omitempty"`
	Height     int64  `json:"height,omitempty"`
	Width      int64  `json:"width,omitempty"`
	Size       int64  `json:"size,omitempty"`
}

type BrandfetchCompany struct {
	Employees   any                  `json:"employees,omitempty"`
	FoundedYear int64                `json:"foundedYear,omitempty"`
	Industries  []BrandfetchIndustry `json:"industries,omitempty"`
	Kind        string               `json:"kind,omitempty"`
	Location    struct {
		City          string `json:"city,omitempty"`
		Country       string `json:"country,omitempty"`
		CountryCodeA2 string `json:"countryCode,omitempty"`
		Region        string `json:"region,omitempty"`
		State         string `json:"state,omitempty"`
		SubRegion     string `json:"subRegion,omitempty"`
	} `json:"location,omitempty"`
}

func (b *BrandfetchCompany) LocationIsEmpty() bool {
	return b.Location.City == "" && b.Location.Country == "" && b.Location.State == "" && b.Location.CountryCodeA2 == ""
}

type BrandfetchIndustry struct {
	Score  float64 `json:"score,omitempty"`
	Name   string  `json:"name,omitempty"`
	Emoji  string  `json:"emoji,omitempty"`
	Slug   string  `json:"slug,omitempty"`
	Parent struct {
		Emoji string `json:"emoji,omitempty"`
		Name  string `json:"name,omitempty"`
		Slug  string `json:"slug,omitempty"`
	} `json:"parent,omitempty"`
}
