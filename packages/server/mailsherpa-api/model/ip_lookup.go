package model

import (
	postgresentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
)

type IpLookupRequest struct {
	Ip string `json:"ip"`
}

type IpLookupResponse struct {
	Status  string                             `json:"status"`
	Message string                             `json:"message,omitempty"`
	IpData  *postgresentity.IPDataResponseBody `json:"ipdata,omitempty"`
}
