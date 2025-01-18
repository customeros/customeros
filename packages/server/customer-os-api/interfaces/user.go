package cosapi_interfaces

import (
	"context"

	commonModel "github.com/customeros/customeros/packages/server/customer-os-common-module/model"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"

	"github.com/customeros/customeros/packages/server/customer-os-api/graphql/model"
)

type UserService interface {
	GetAll(ctx context.Context, page, limit int, filter *model.Filter, sortBy []*commonModel.SortBy) (*utils.Pagination, error)
	ContainsRole(parentCtx context.Context, allowedRoles []model.Role) bool
}
