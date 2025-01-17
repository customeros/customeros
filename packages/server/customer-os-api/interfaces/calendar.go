package cosapi_interfaces

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
)

type CalendarService interface {
	GetAllForUsers(ctx context.Context, userIds []string) (*entity.CalendarEntities, error)
}
