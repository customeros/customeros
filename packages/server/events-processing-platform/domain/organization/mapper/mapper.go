package mapper

import (
	organization_grpc_service "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/organization"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/models"
)

func MapCustomFieldDataType(input organization_grpc_service.CustomFieldDataType) models.CustomFieldDataType {
	switch input {
	case organization_grpc_service.CustomFieldDataType_TEXT:
		return models.CustomFieldDataTypeText
	case organization_grpc_service.CustomFieldDataType_BOOL:
		return models.CustomFieldDataTypeBool
	case organization_grpc_service.CustomFieldDataType_DATETIME:
		return models.CustomFieldDataTypeDatetime
	case organization_grpc_service.CustomFieldDataType_INTEGER:
		return models.CustomFieldDataTypeInteger
	case organization_grpc_service.CustomFieldDataType_DECIMAL:
		return models.CustomFieldDataTypeDecimal
	default:
		return models.CustomFieldDataTypeText
	}
}
