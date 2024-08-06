package model

type IpLookupRequest struct {
	Ip string `json:"ip"`
}

type IpLookupData struct {
	Ip string `json:"ip"`
}

type IpLookupResponse struct {
	Status  string        `json:"status"`
	Message string        `json:"message,omitempty"`
	Data    *IpLookupData `json:"data,omitempty"`
}
