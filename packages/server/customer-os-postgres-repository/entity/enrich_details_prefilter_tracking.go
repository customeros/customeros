package entity

import (
	"time"
)

type EnrichDetailsPreFilterTracking struct {
	ID             string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	CreatedAt      time.Time `gorm:"column:created_at;type:timestamp;DEFAULT:current_timestamp" json:"createdAt"`
	UpdatedAt      time.Time `gorm:"column:updated_at;type:timestamp;" json:"updatedAt"`
	IP             string    `gorm:"column:ip;uniqueIndex:ip_unique;type:varchar(255);" json:"ip"`
	ShouldIdentify *bool     `gorm:"column:should_identify;type:boolean;" json:"shouldIdentify"`
	Response       *string   `gorm:"column:response;type:text;" json:"response"`
}

func (EnrichDetailsPreFilterTracking) TableName() string {
	return "enrich_details_prefilter_tracking"
}

type IPDataResponseBody struct {
	Ip            string  `json:"ip"`
	City          string  `json:"city"`
	Region        string  `json:"region"`
	RegionCode    string  `json:"region_code"`
	RegionType    string  `json:"region_type"`
	CountryName   string  `json:"country_name"`
	CountryCode   string  `json:"country_code"`
	ContinentName string  `json:"continent_name"`
	ContinentCode string  `json:"continent_code"`
	Latitude      float64 `json:"latitude"`
	Longitude     float64 `json:"longitude"`
	Asn           struct {
		Asn    string `json:"asn"`
		Name   string `json:"name"`
		Domain string `json:"domain"`
		Route  string `json:"route"`
		Type   string `json:"type"`
	} `json:"asn"`
	Carrier *struct {
		Name string `json:"name"`
		Mcc  string `json:"mcc"`
		Mnc  string `json:"mnc"`
	} `json:"carrier"`
	TimeZone struct {
		Name        string    `json:"name"`
		Abbr        string    `json:"abbr"`
		Offset      string    `json:"offset"`
		IsDst       bool      `json:"is_dst"`
		CurrentTime time.Time `json:"current_time"`
	} `json:"time_zone"`
	Threat struct {
		IsTor           bool          `json:"is_tor"`
		IsIcloudRelay   bool          `json:"is_icloud_relay"`
		IsProxy         bool          `json:"is_proxy"`
		IsDatacenter    bool          `json:"is_datacenter"`
		IsAnonymous     bool          `json:"is_anonymous"`
		IsKnownAttacker bool          `json:"is_known_attacker"`
		IsKnownAbuser   bool          `json:"is_known_abuser"`
		IsThreat        bool          `json:"is_threat"`
		IsBogon         bool          `json:"is_bogon"`
		Blocklists      []interface{} `json:"blocklists"`
	} `json:"threat"`
	Count string `json:"count"`
}
