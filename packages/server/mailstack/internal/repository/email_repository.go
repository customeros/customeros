package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/customeros/customeros/packages/server/mailstack/interfaces"
	"github.com/customeros/customeros/packages/server/mailstack/internal/models"
)

type GormEmailRepository struct {
	db *gorm.DB
}

func NewEmailRepository(db *gorm.DB) interfaces.EmailRepository {
	return &GormEmailRepository{
		db: db,
	}
}

// Create adds a new email to the database
func (r *GormEmailRepository) Create(ctx context.Context, email *models.Email) error {
	return r.db.WithContext(ctx).Create(email).Error
}

// GetByID retrieves an email by its ID
func (r *GormEmailRepository) GetByID(ctx context.Context, id string) (*models.Email, error) {
	var email models.Email
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&email).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &email, nil
}

// GetByUID retrieves an email by its UID within a specific mailbox and folder
func (r *GormEmailRepository) GetByUID(ctx context.Context, mailboxID, folder string, uid uint32) (*models.Email, error) {
	var email models.Email
	if err := r.db.WithContext(ctx).
		Where("mailbox_id = ? AND folder = ? AND uid = ?", mailboxID, folder, uid).
		First(&email).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &email, nil
}

// GetByMessageID retrieves an email by its Message-ID header
func (r *GormEmailRepository) GetByMessageID(ctx context.Context, messageID string) (*models.Email, error) {
	var email models.Email
	if err := r.db.WithContext(ctx).Where("message_id = ?", messageID).First(&email).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &email, nil
}

// ListByMailbox retrieves emails for a specific mailbox with pagination
func (r *GormEmailRepository) ListByMailbox(ctx context.Context, mailboxID string, limit, offset int) ([]*models.Email, int64, error) {
	var emails []*models.Email
	var count int64

	// Count total emails
	if err := r.db.WithContext(ctx).Model(&models.Email{}).
		Where("mailbox_id = ?", mailboxID).
		Count(&count).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated emails
	if err := r.db.WithContext(ctx).
		Where("mailbox_id = ?", mailboxID).
		Order("received_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&emails).Error; err != nil {
		return nil, 0, err
	}

	return emails, count, nil
}

// ListByFolder retrieves emails for a specific mailbox and folder with pagination
func (r *GormEmailRepository) ListByFolder(ctx context.Context, mailboxID, folder string, limit, offset int) ([]*models.Email, int64, error) {
	var emails []*models.Email
	var count int64

	// Count total emails in folder
	if err := r.db.WithContext(ctx).Model(&models.Email{}).
		Where("mailbox_id = ? AND folder = ?", mailboxID, folder).
		Count(&count).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated emails from folder
	if err := r.db.WithContext(ctx).
		Where("mailbox_id = ? AND folder = ?", mailboxID, folder).
		Order("received_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&emails).Error; err != nil {
		return nil, 0, err
	}

	return emails, count, nil
}

// ListByThread retrieves all emails in a conversation thread
func (r *GormEmailRepository) ListByThread(ctx context.Context, threadID string) ([]*models.Email, error) {
	var emails []*models.Email

	if err := r.db.WithContext(ctx).
		Where("thread_id = ?", threadID).
		Order("received_at ASC").
		Find(&emails).Error; err != nil {
		return nil, err
	}

	return emails, nil
}

// Search searches emails by query string
func (r *GormEmailRepository) Search(ctx context.Context, query string, limit, offset int) ([]*models.Email, int64, error) {
	var emails []*models.Email
	var count int64

	// Build search condition - search in subject, body, sender, recipients
	searchCondition := "subject ILIKE ? OR body_text ILIKE ? OR from_address ILIKE ? OR from_name ILIKE ?"
	searchParam := "%" + query + "%"

	// Count total matching emails
	if err := r.db.WithContext(ctx).Model(&models.Email{}).
		Where(searchCondition, searchParam, searchParam, searchParam, searchParam).
		Count(&count).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated search results
	if err := r.db.WithContext(ctx).
		Where(searchCondition, searchParam, searchParam, searchParam, searchParam).
		Order("received_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&emails).Error; err != nil {
		return nil, 0, err
	}

	return emails, count, nil
}

// Update updates an email record
func (r *GormEmailRepository) Update(ctx context.Context, email *models.Email) error {
	return r.db.WithContext(ctx).Save(email).Error
}
