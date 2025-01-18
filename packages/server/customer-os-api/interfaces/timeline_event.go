package cosapi_interfaces

import (
	"context"
	"time"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graphql/model"
)

type TimelineEventService interface {
	GetTimelineEventsForContact(ctx context.Context, contactId string, from *time.Time, size int, types []model.TimelineEventType) (*entity.TimelineEventEntities, error)
	GetTimelineEventsTotalCountForContact(ctx context.Context, contactId string, types []model.TimelineEventType) (int64, error)
	GetTimelineEventsForOrganization(ctx context.Context, organizationId string, from *time.Time, size int, types []model.TimelineEventType) (*entity.TimelineEventEntities, error)
	GetTimelineEventsTotalCountForOrganization(ctx context.Context, organizationId string, types []model.TimelineEventType) (int64, error)
	GetTimelineEventsWithIds(ctx context.Context, ids []string) (*entity.TimelineEventEntities, error)
	GetInboundCommsCountCountByOrganizations(ctx context.Context, organizationIds []string) (map[string]int64, error)
	GetOutboundCommsCountCountByOrganizations(ctx context.Context, organizationIds []string) (map[string]int64, error)
}
