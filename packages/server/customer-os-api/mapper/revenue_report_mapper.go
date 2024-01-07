package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
)

func MapFileGeneratorResponseData(data *entity.FileGeneratorResponseData) *model.FileGeneratorResponse {
	if data == nil {
		return nil
	}
	return &model.FileGeneratorResponse{
		Success: data.Success,
		Message: data.Message,
	}
}
