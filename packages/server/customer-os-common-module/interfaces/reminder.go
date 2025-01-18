package interfaces

import (
	"context"
	"time"

	neo4j_entity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"
)

type ReminderService interface {
	CreateReminder(ctx context.Context, tenant, userId, orgId, content string, dueDate time.Time) (string, error)
	UpdateReminder(ctx context.Context, tenant, id string, content *string, dueDate *time.Time, dismissed, sent *bool) error
	GetReminderById(ctx context.Context, id string) (*neo4j_entity.ReminderEntity, error)
	RemindersForOrganization(ctx context.Context, organizationID string, dismissed *bool) ([]*neo4j_entity.ReminderEntity, error)

	SendNotification(ctx context.Context, reminderId, fronteraPublicPath string) error
}
