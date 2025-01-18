package cosapi_interfaces

import (
	"context"
	neo4jentity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"
)

// TODO deprecate and remove
type CustomFieldTemplateService interface {
	FindLinkedWithCustomField(ctx context.Context, customFieldId string) (*neo4jentity.CustomFieldTemplateEntity, error)
}
