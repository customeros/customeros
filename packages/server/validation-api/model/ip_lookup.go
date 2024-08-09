package model

import "time"

type IpLookupRequest struct {
	Ip string `json:"ip"`
}

type IpLookupData struct {
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

type IpLookupResponse struct {
	Status  string        `json:"status"`
	Message string        `json:"message,omitempty"`
	IpData  *IpLookupData `json:"ipdata,omitempty"`
}
