package enum

import "github.com/customeros/customeros/packages/server/customer-os-common-module/utils"

type BaseResponse struct {
	RequestID string `json:"requestId" example:"api_1234567890abcdef"`
	// Status indicates the result of the operation ("success" or "error")
	Status string `json:"status" example:"success"`
}

func BuildBaseResponse(status Status) BaseResponse {
	return BaseResponse{
		RequestID: utils.GenerateNanoIdWithPrefix("api", 16),
		Status:    string(status),
	}
}
