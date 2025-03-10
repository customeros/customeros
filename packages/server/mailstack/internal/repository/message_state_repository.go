package repository

import (
	"context"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/opentracing/opentracing-go"
	"gorm.io/gorm"

	"github.com/customeros/customeros/packages/server/mailstack/interfaces"
	"github.com/customeros/customeros/packages/server/mailstack/internal/models"
)

type messageStateRepository struct {
	db *gorm.DB
}

func NewMessageStateRepository(db *gorm.DB) interfaces.MessageStateRepository {
	return &messageStateRepository{db: db}
}

func (r *messageStateRepository) GetLastSeenUID(ctx context.Context, mailboxID, folderName string) (uint32, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "messageStateRepository.GetLastSeenUID")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	var state models.MessageState
	err := r.db.First(&state, "mailbox_id = ? AND folder_name = ?", mailboxID, folderName).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return 0, nil // Return 0 for new folders
		}
		tracing.TraceErr(span, err)
		return 0, err
	}

	return state.LastSeenUID, nil
}

func (r *messageStateRepository) UpdateLastSeenUID(ctx context.Context, mailboxID, folderName string, uid uint32) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "messageStateRepository.GetLastSeenUID")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	return r.db.Exec(
		`INSERT INTO message_states (mailbox_id, folder_name, last_seen_uid, last_checked_at) 
		VALUES (?, ?, ?, NOW()) 
		ON CONFLICT (mailbox_id, folder_name) 
		DO UPDATE SET last_seen_uid = ?, last_checked_at = NOW()`,
		mailboxID, folderName, uid, uid,
	).Error
}
