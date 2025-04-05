package postgres_repository

import (
	"context"
	"errors"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/coserrors"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	"gorm.io/gorm"

	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
)

type WebhooksRepository interface {
	FindAll(ctx context.Context) (*[]postgres_entity.Webhooks, error)
	Find(ctx context.Context, webhook postgres_entity.Webhooks) (*postgres_entity.Webhooks, error)
	FindLastRotationCount(ctx context.Context, integration enum.Source) (int, error)
	Create(ctx context.Context, webhook postgres_entity.Webhooks) (*postgres_entity.Webhooks, error)
	Deactivate(ctx context.Context, webhook postgres_entity.Webhooks) (ok bool, err error)
}

type webhooksRepository struct {
	gormDb *gorm.DB
}

func NewWebhooksRepository(gormDb *gorm.DB) WebhooksRepository {
	return &webhooksRepository{gormDb: gormDb}
}

func (r *webhooksRepository) Create(ctx context.Context, webhook postgres_entity.Webhooks) (*postgres_entity.Webhooks, error) {
	spans, ctx := telemetry.StartPostgresSpan(ctx, "WebhooksRepository.Create")
	defer spans.Finish()

	var created postgres_entity.Webhooks
	now := utils.Now()

	err := r.gormDb.Transaction(func(tx *gorm.DB) error {
		// First, deactivate any existing webhooks for this integration
		err := tx.Exec(`
            UPDATE webhooks 
            SET enabled = FALSE, updated_at = $1
            WHERE integration = $2 
            AND tenant = $3
            AND enabled = TRUE
        `, now, webhook.Integration, webhook.Tenant).Error
		if err != nil {
			spans.TraceError(err)
			return err
		}

		// Then create the new webhook
		return tx.Raw(`
            INSERT INTO webhooks (
                tenant,
                integration,
                webhook_path,
                enabled,
                created_at,
                updated_at,
                rotation_count,
                secret
            ) VALUES ($1, $2, $3, TRUE, $4, $5, $6, $7)
            RETURNING *
        `,
			webhook.Tenant,
			webhook.Integration,
			webhook.WebhookPath,
			now,
			now,
			webhook.RotationCount,
			webhook.Secret,
		).Scan(&created).Error
	})
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}

	return &created, nil
}

func (r *webhooksRepository) FindAll(ctx context.Context) (*[]postgres_entity.Webhooks, error) {
	spans, ctx := telemetry.StartPostgresSpan(ctx, "WebhooksRepository.FindAll")
	defer spans.Finish()

	tenant := common.GetTenantFromContext(ctx)
	if tenant == "" {
		err := errors.New("Tenant not set in context")
		spans.TraceError(err)
		return nil, err
	}

	var webhooks []postgres_entity.Webhooks
	err := r.gormDb.
		Where("enabled = ? AND tenant = ?", true, tenant).
		Order("created_at DESC").
		Find(&webhooks).Error
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}

	return &webhooks, nil
}

func (r *webhooksRepository) Find(ctx context.Context, webhook postgres_entity.Webhooks) (*postgres_entity.Webhooks, error) {
	spans, ctx := telemetry.StartPostgresSpan(ctx, "WebhooksRepository.Find")
	defer spans.Finish()

	if webhook.Tenant == "" {
		webhook.Tenant = common.GetTenantFromContext(ctx)
		if webhook.Tenant == "" {
			err := errors.New("Tenant not set in context")
			spans.TraceError(err)
			return nil, err
		}
	}

	var foundWebhook postgres_entity.Webhooks
	query := r.gormDb.Where("enabled = ?", true)

	// Add additional filters based on non-zero fields in webhook
	if webhook != (postgres_entity.Webhooks{}) {
		query = query.Where(&webhook)
	}

	err := query.First(&foundWebhook).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		spans.TraceError(err)
		return nil, err
	}

	return &foundWebhook, nil
}

func (r *webhooksRepository) FindLastRotationCount(ctx context.Context, integration enum.Source) (int, error) {
	spans, ctx := telemetry.StartPostgresSpan(ctx, "WebhooksRepository.FindLastRotationCount")
	defer spans.Finish()

	tenant := common.GetTenantFromContext(ctx)
	if tenant == "" {
		return 0, coserrors.ErrTenantNotSet
	}

	var webhook postgres_entity.Webhooks
	err := r.gormDb.
		Where("tenant = ?", tenant).
		Where("integration = ?", integration.String()).
		Order("rotation_count DESC").
		First(&webhook).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, nil
		}
		spans.TraceError(err)
		return 0, err
	}
	return webhook.RotationCount, nil
}

func (r *webhooksRepository) Deactivate(ctx context.Context, webhook postgres_entity.Webhooks) (bool, error) {
	spans, ctx := telemetry.StartPostgresSpan(ctx, "WebhooksRepository.Update")
	defer spans.Finish()

	if webhook.Tenant == "" {
		webhook.Tenant = common.GetTenantFromContext(ctx)
		if webhook.Tenant == "" {
			err := errors.New("Tenant not set in context")
			spans.TraceError(err)
			return false, err
		}
	}

	if webhook.WebhookPath == "" {
		err := errors.New("Webhook path is missing")
		spans.TraceError(err)
		return false, err
	}

	result := r.gormDb.Model(&postgres_entity.Webhooks{}).
		Where("webhook_path = ? AND tenant = ? AND enabled = true", webhook.WebhookPath, webhook.Tenant).
		Updates(map[string]interface{}{
			"enabled":    false,
			"updated_at": utils.NowPtr(),
		})

	if result.Error != nil {
		spans.TraceError(result.Error)
		return false, result.Error
	}

	// Check if any rows were affected
	if result.RowsAffected == 0 {
		return false, nil
	}

	return true, nil
}
