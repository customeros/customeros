package cosapi_interfaces

import (
	"context"
	"github.com/customeros/customeros/packages/server/customer-os-api/entity"
)

type CalendarService interface {
	GetAllForUsers(ctx context.Context, userIds []string) (*entity.CalendarEntities, error)
}
