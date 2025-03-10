package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/opentracing/opentracing-go"
	"gorm.io/gorm"

	"github.com/customeros/customeros/packages/server/mailstack/interfaces"
	"github.com/customeros/customeros/packages/server/mailstack/internal/models"
)

type mailboxSyncRepository struct {
	db *gorm.DB
}

func NewMailboxSyncRepository(db *gorm.DB) interfaces.MailboxSyncRepository {
	return &mailboxSyncRepository{db: db}
}

// GetSyncState retrieves the sync state for a specific mailbox and folder
func (r *mailboxSyncRepository) GetSyncState(ctx context.Context, mailboxID, folderName string) (*models.MailboxSyncState, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "mailboxSyncRepository.GetSyncState")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	var state models.MailboxSyncState
	result := r.db.WithContext(ctx).
		Where("mailbox_id = ? AND folder_name = ?", mailboxID, folderName).
		First(&state)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil // No sync state yet
		}
		tracing.TraceErr(span, result.Error)
		return nil, fmt.Errorf("failed to get sync state: %w", result.Error)
	}

	return &state, nil
}

// SaveSyncState saves the sync state for a mailbox folder
func (r *mailboxSyncRepository) SaveSyncState(ctx context.Context, state *models.MailboxSyncState) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "mailboxSyncRepository.SaveSyncState")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	// Set the last sync time
	state.LastSync = time.Now()

	// Try to update first
	result := r.db.WithContext(ctx).
		Model(&models.MailboxSyncState{}).
		Where("mailbox_id = ? AND folder_name = ?", state.MailboxID, state.FolderName).
		Updates(map[string]interface{}{
			"last_uid":   state.LastUID,
			"last_sync":  state.LastSync,
			"updated_at": time.Now(),
		})

	// If no record was updated, create a new one
	if result.RowsAffected == 0 {
		result = r.db.WithContext(ctx).Create(state)
	}

	if result.Error != nil {
		tracing.TraceErr(span, result.Error)
		return fmt.Errorf("failed to save sync state: %w", result.Error)
	}

	return nil
}

// DeleteSyncState deletes the sync state for a mailbox folder
func (r *mailboxSyncRepository) DeleteSyncState(ctx context.Context, mailboxID, folderName string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "mailboxSyncRepository.DeleteSyncState")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	result := r.db.WithContext(ctx).
		Where("mailbox_id = ? AND folder_name = ?", mailboxID, folderName).
		Delete(&models.MailboxSyncState{})

	if result.Error != nil {
		tracing.TraceErr(span, result.Error)
		return fmt.Errorf("failed to delete sync state: %w", result.Error)
	}

	return nil
}

// DeleteMailboxSyncStates deletes all sync states for a mailbox
func (r *mailboxSyncRepository) DeleteMailboxSyncStates(ctx context.Context, mailboxID string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "mailboxSyncRepository.DeleteMailboxSyncStates")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	result := r.db.WithContext(ctx).
		Where("mailbox_id = ?", mailboxID).
		Delete(&models.MailboxSyncState{})

	if result.Error != nil {
		tracing.TraceErr(span, result.Error)
		return fmt.Errorf("failed to delete mailbox sync states: %w", result.Error)
	}

	return nil
}

// GetAllSyncStates gets all sync states
func (r *mailboxSyncRepository) GetAllSyncStates(ctx context.Context) (map[string]map[string]uint32, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "mailboxSyncRepository.GetAllSyncStates")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	var states []models.MailboxSyncState
	if err := r.db.WithContext(ctx).Find(&states).Error; err != nil {
		tracing.TraceErr(span, err)
		return nil, fmt.Errorf("failed to get all sync states: %w", err)
	}

	result := make(map[string]map[string]uint32)
	for _, state := range states {
		if _, ok := result[state.MailboxID]; !ok {
			result[state.MailboxID] = make(map[string]uint32)
		}
		result[state.MailboxID][state.FolderName] = state.LastUID
	}

	return result, nil
}

// GetMailboxSyncStates gets all sync states for a mailbox
func (r *mailboxSyncRepository) GetMailboxSyncStates(ctx context.Context, mailboxID string) (map[string]uint32, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "mailboxSyncRepository.GetMailboxSyncStates")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	var states []models.MailboxSyncState
	if err := r.db.WithContext(ctx).Where("mailbox_id = ?", mailboxID).Find(&states).Error; err != nil {
		tracing.TraceErr(span, err)
		return nil, fmt.Errorf("failed to get mailbox sync states: %w", err)
	}

	result := make(map[string]uint32)
	for _, state := range states {
		result[state.FolderName] = state.LastUID
	}

	return result, nil
}
