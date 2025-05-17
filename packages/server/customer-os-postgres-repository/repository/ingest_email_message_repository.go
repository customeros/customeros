package postgres_repository

import (
	"errors"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"golang.org/x/net/context"
	"gorm.io/gorm"

	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
)

type IngestEmailMessageRepository interface {
	Store(ctx context.Context, tenant, username, provider, messageId string, ingestEmailMessage *postgres_entity.IngestEmailMessage) error
	UpdateState(ctx context.Context, id string, state postgres_entity.IngestEmailMessageState) error

	CountForUsername(ctx context.Context, tenant, username, provider string) (int64, error)
	GetByMessageId(ctx context.Context, externalSystem, tenant, username, messageId string) (*postgres_entity.IngestEmailMessage, error)
	EmailExistsByMessageId(ctx context.Context, tenant, username, provider, messageId string) (bool, error)
	GetEmail(ctx context.Context, id string) (*postgres_entity.IngestEmailMessage, error)
	GetDistinctUsersForPendingMessages(ctx context.Context) ([]postgres_entity.IngestEmailMessage, error)
	GetEmailsForUserForSync(ctx context.Context, tenant, username string) ([]postgres_entity.IngestEmailMessage, error)
}

type ingestEmailMessageRepositoryImpl struct {
	gormDb *gorm.DB
}

func NewIngestEmailMessageRepository(gormDb *gorm.DB) IngestEmailMessageRepository {
	return &ingestEmailMessageRepositoryImpl{gormDb: gormDb}
}

func (repo *ingestEmailMessageRepositoryImpl) Store(ctx context.Context, tenant, username, provider, messageId string, ingestEmailMessage *postgres_entity.IngestEmailMessage) error {
	spans, ctx := telemetry.StartPostgresSpan(ctx, "IngestEmailMessageRepository.Store")
	defer spans.Finish()

	spans.TagTenant(tenant)

	result := postgres_entity.IngestEmailMessage{}
	err := repo.gormDb.Find(&result, "tenant = ? AND username = ? AND provider = ? AND message_id = ?", tenant, username, provider, messageId).Error
	if err != nil {
		spans.TraceError(err)
		return err
	}

	if result.Tenant != "" {
		err := errors.New("IngestEmailMessageRepository.Store - email already exists")
		spans.TraceError(err)
		return err
	}

	err = repo.gormDb.Save(&ingestEmailMessage).Error
	if err != nil {
		spans.TraceError(err)
		return err
	}

	return nil
}

func (repo *ingestEmailMessageRepositoryImpl) UpdateState(ctx context.Context, id string, state postgres_entity.IngestEmailMessageState) error {
	spans, ctx := telemetry.StartPostgresSpan(ctx, "IngestEmailMessageRepository.UpdateState")
	defer spans.Finish()

	err := repo.gormDb.Model(&postgres_entity.IngestEmailMessage{}).
		Where("id = ?", id).
		Update("state", state).
		Error
	if err != nil {
		spans.TraceError(err)
		return err
	}

	return nil
}

func (repo *ingestEmailMessageRepositoryImpl) CountForUsername(ctx context.Context, tenant, username, provider string) (int64, error) {
	spans, _ := telemetry.StartPostgresSpan(ctx, "IngestEmailMessageRepository.CountForUsername")
	defer spans.Finish()

	var result int64
	err := repo.gormDb.Model(postgres_entity.IngestEmailMessage{}).Where("provider = ? AND tenant = ? AND username = ?", provider, tenant, username).Count(&result).Error
	if err != nil {
		spans.TraceError(err)
		return 0, err
	}

	return result, nil
}

func (repo *ingestEmailMessageRepositoryImpl) GetByMessageId(ctx context.Context, externalSystem, tenant, username, messageId string) (*postgres_entity.IngestEmailMessage, error) {
	spans, _ := telemetry.StartPostgresSpan(ctx, "IngestEmailMessageRepository.GetByMessageId")
	defer spans.Finish()

	var result *postgres_entity.IngestEmailMessage
	err := repo.gormDb.Where("provider = ? AND tenant = ? AND username = ? AND message_id = ?", externalSystem, tenant, username, messageId).First(&result).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return result, nil
}

func (repo *ingestEmailMessageRepositoryImpl) EmailExistsByMessageId(ctx context.Context, tenant, username, provider, messageId string) (bool, error) {
	spans, _ := telemetry.StartPostgresSpan(ctx, "IngestEmailMessageRepository.EmailExistsByMessageId")
	defer spans.Finish()

	var result int64
	err := repo.gormDb.Model(postgres_entity.IngestEmailMessage{}).Where("provider = ? AND tenant = ? AND username = ? AND message_id = ?", provider, tenant, username, messageId).Count(&result).Error
	if err != nil {
		spans.TraceError(err)
		return false, err
	}

	return result > 0, nil
}

func (repo *ingestEmailMessageRepositoryImpl) GetEmail(ctx context.Context, id string) (*postgres_entity.IngestEmailMessage, error) {
	spans, ctx := telemetry.StartPostgresSpan(ctx, "IngestEmailMessageRepository.GetEmail")
	defer spans.Finish()

	result := postgres_entity.IngestEmailMessage{}
	err := repo.gormDb.First(&result, "id = ? ", id).Error
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}

	return &result, nil
}

func (repo *ingestEmailMessageRepositoryImpl) GetDistinctUsersForPendingMessages(ctx context.Context) ([]postgres_entity.IngestEmailMessage, error) {
	spans, ctx := telemetry.StartPostgresSpan(ctx, "IngestEmailMessageRepository.GetDistinctUsersForPendingMessages")
	defer spans.Finish()

	var results []postgres_entity.IngestEmailMessage

	err := repo.gormDb.Select("DISTINCT tenant, username").Where("state = ?", postgres_entity.IngestEmailMessageStatePending).Find(&results).Error
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}

	return results, nil
}

func (repo *ingestEmailMessageRepositoryImpl) GetEmailsForUserForSync(ctx context.Context, tenant, username string) ([]postgres_entity.IngestEmailMessage, error) {
	spans, ctx := telemetry.StartPostgresSpan(ctx, "IngestEmailMessageRepository.GetEmailsForUserForSync")
	defer spans.Finish()

	results := []postgres_entity.IngestEmailMessage{}

	err := repo.gormDb.Order("sent_at desc").Limit(50).Find(&results, "tenant = ? AND username = ? AND state = ?", tenant, username, postgres_entity.IngestEmailMessageStatePending).Error
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}

	return results, nil
}
