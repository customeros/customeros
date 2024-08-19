package entity

import "time"

type CacheIpData struct {
	ID        string    `gorm:"primary_key;type:uuid;default:gen_random_uuid()" json:"id"`
	CreatedAt time.Time `gorm:"column:created_at;type:timestamp;DEFAULT:current_timestamp" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:timestamp;DEFAULT:current_timestamp" json:"updatedAt"`
	Ip        string    `gorm:"column:ip;type:varchar(255);NOT NULL;index:idx_cache_ip_data_ip,unique" json:"ip"`
	Data      string    `gorm:"column:data;type:text" json:"data"`
}

func (CacheIpData) TableName() string {
	return "cache_ip_data"
}

type IPDataResponseBody struct {
	StatusCode    int     `json:"status_code"`
	Message       string  `json:"message"`
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
		IsVpn           bool          `json:"is_vpn"`
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
