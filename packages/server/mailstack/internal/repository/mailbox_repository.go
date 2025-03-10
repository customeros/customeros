package repository

import (
	"context"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/opentracing/opentracing-go"
	"gorm.io/gorm"

	"github.com/customeros/customeros/packages/server/mailstack/interfaces"
	"github.com/customeros/customeros/packages/server/mailstack/internal/models"
)

type mailboxRepository struct {
	db *gorm.DB
}

func NewMailboxRepository(db *gorm.DB) interfaces.MailboxRepository {
	return &mailboxRepository{db: db}
}

func (r *mailboxRepository) GetMailboxes(ctx context.Context) ([]*models.Mailbox, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "mailboxRepository.GetMailboxes")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	var mailboxes []*models.Mailbox
	result := r.db.Find(&mailboxes)
	if result.Error != nil {
		tracing.TraceErr(span, result.Error)
		return nil, result.Error
	}
	return mailboxes, nil
}

func (r *mailboxRepository) GetMailbox(ctx context.Context, id string) (*models.Mailbox, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "mailboxRepository.GetMailbox")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	var mailbox models.Mailbox
	err := r.db.First(&mailbox, "id = ?", id).Error
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	return &mailbox, nil
}

func (r *mailboxRepository) SaveMailbox(ctx context.Context, mailbox models.Mailbox) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "mailboxRepository.SaveMailbox")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	return r.db.Save(&mailbox).Error
}

func (r *mailboxRepository) DeleteMailbox(ctx context.Context, id string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "mailboxRepository.DeleteMailbox")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	return r.db.Delete(&models.Mailbox{}, "id = ?", id).Error
}
