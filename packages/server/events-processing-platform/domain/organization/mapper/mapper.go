package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/model"
	organizationpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/organization"
)

func MapCustomFieldDataType(input organizationpb.CustomFieldDataType) model.CustomFieldDataType {
	switch input {
	case organizationpb.CustomFieldDataType_TEXT:
		return model.CustomFieldDataTypeText
	case organizationpb.CustomFieldDataType_BOOL:
		return model.CustomFieldDataTypeBool
	case organizationpb.CustomFieldDataType_DATETIME:
		return model.CustomFieldDataTypeDatetime
	case organizationpb.CustomFieldDataType_INTEGER:
		return model.CustomFieldDataTypeInteger
	case organizationpb.CustomFieldDataType_DECIMAL:
		return model.CustomFieldDataTypeDecimal
	default:
		return model.CustomFieldDataTypeText
	}
}
