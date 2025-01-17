package interfaces

import (
	"context"
	"time"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
)

type ReminderService interface {
	CreateReminder(ctx context.Context, tenant, userId, orgId, content string, dueDate time.Time) (string, error)
	UpdateReminder(ctx context.Context, tenant, id string, content *string, dueDate *time.Time, dismissed, sent *bool) error
	GetReminderById(ctx context.Context, id string) (*entity.ReminderEntity, error)
	RemindersForOrganization(ctx context.Context, organizationID string, dismissed *bool) ([]*entity.ReminderEntity, error)

	SendNotification(ctx context.Context, reminderId, fronteraPublicPath string) error
}
