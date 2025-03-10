package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	common_interfaces "github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/opentracing/opentracing-go"
	"gorm.io/gorm"

	"github.com/customeros/customeros/packages/server/mailstack/interfaces"
	"github.com/customeros/customeros/packages/server/mailstack/internal/models"
)

type emailAttachmentRepository struct {
	db      *gorm.DB
	storage common_interfaces.StorageService
}

func NewEmailAttachmentRepository(db *gorm.DB, storageService common_interfaces.StorageService) interfaces.EmailAttachmentRepository {
	return &emailAttachmentRepository{
		db:      db,
		storage: storageService,
	}
}

// Create adds a new attachment to the database
func (r *emailAttachmentRepository) Create(ctx context.Context, attachment *models.EmailAttachment) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "emailAttachmentRepository.Create")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	return r.db.WithContext(ctx).Create(attachment).Error
}

// GetByID retrieves an attachment by its ID
func (r *emailAttachmentRepository) GetByID(ctx context.Context, id string) (*models.EmailAttachment, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "emailAttachmentRepository.GetByID")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	var attachment models.EmailAttachment
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&attachment).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		tracing.TraceErr(span, err)
		return nil, err
	}
	return &attachment, nil
}

// ListByEmail retrieves all attachments for a specific email
func (r *emailAttachmentRepository) ListByEmail(ctx context.Context, emailID string) ([]*models.EmailAttachment, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "emailAttachmentRepository.ListByEmail")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	var attachments []*models.EmailAttachment
	err := r.db.WithContext(ctx).
		Where("email_id = ?", emailID).
		Find(&attachments).Error
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	return attachments, nil
}

// Store saves attachment data to the configured storage service
func (r *emailAttachmentRepository) Store(ctx context.Context, attachment *models.EmailAttachment, data []byte) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "emailAttachmentRepository.Store")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	// Generate a storage key if one doesn't exist
	if attachment.StorageKey == "" {
		attachment.StorageKey = fmt.Sprintf("attachments/%s/%s", attachment.EmailID, attachment.ID)
	}

	// Store the file in the storage service
	if err := r.storage.Upload(ctx, attachment.StorageKey, data, attachment.ContentType); err != nil {
		tracing.TraceErr(span, err)
		return fmt.Errorf("failed to upload attachment: %w", err)
	}

	// Update the attachment record with the storage key and size
	attachment.Size = len(data)
	attachment.UpdatedAt = time.Now()

	return r.db.WithContext(ctx).Save(attachment).Error
}

// GetData retrieves the attachment data from storage
func (r *emailAttachmentRepository) GetData(ctx context.Context, id string) ([]byte, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "emailAttachmentRepository.GetData")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	// Get the attachment metadata
	attachment, err := r.GetByID(ctx, id)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	if attachment == nil {
		tracing.TraceErr(span, err)
		return nil, errors.New("attachment not found")
	}

	// Retrieve the file from storage
	data, err := r.storage.Download(ctx, attachment.StorageKey)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, fmt.Errorf("failed to download attachment: %w", err)
	}

	return data, nil
}

// Delete removes an attachment from both database and storage
func (r *emailAttachmentRepository) Delete(ctx context.Context, id string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "emailAttachmentRepository.Delete")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	// Get the attachment first
	attachment, err := r.GetByID(ctx, id)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	if attachment == nil {
		return nil // Already deleted
	}

	// Delete from storage
	if attachment.StorageKey != "" {
		if err := r.storage.Delete(ctx, attachment.StorageKey); err != nil {
			// Log the error but continue with DB deletion
			fmt.Printf("Failed to delete attachment from storage: %v", err)
		}
	}

	// Delete from database
	return r.db.WithContext(ctx).Delete(&models.EmailAttachment{}, "id = ?", id).Error
}
