package mapper

import (
	organization_grpc_service "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/organization"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/model"
)

func MapCustomFieldDataType(input organization_grpc_service.CustomFieldDataType) model.CustomFieldDataType {
	switch input {
	case organization_grpc_service.CustomFieldDataType_TEXT:
		return model.CustomFieldDataTypeText
	case organization_grpc_service.CustomFieldDataType_BOOL:
		return model.CustomFieldDataTypeBool
	case organization_grpc_service.CustomFieldDataType_DATETIME:
		return model.CustomFieldDataTypeDatetime
	case organization_grpc_service.CustomFieldDataType_INTEGER:
		return model.CustomFieldDataTypeInteger
	case organization_grpc_service.CustomFieldDataType_DECIMAL:
		return model.CustomFieldDataTypeDecimal
	default:
		return model.CustomFieldDataTypeText
	}
}
