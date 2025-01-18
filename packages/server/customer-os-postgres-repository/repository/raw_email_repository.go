package repository

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"gorm.io/gorm"

	"github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
)

type RawEmailRepository interface {
	CountForUsername(ctx context.Context, externalSystem, tenant, username string) (int64, error)
	GetByMessageId(ctx context.Context, externalSystem, tenant, username, messageId string) (*entity.RawEmail, error)
	EmailExistsByMessageId(ctx context.Context, externalSystem, tenant, username, messageId string) (bool, error)
	Store(ctx context.Context, externalSystem, tenant, username, providerMessageId, messageId, rawEmail string, sentAt time.Time, state entity.EmailImportState) error
	GetEmailsIdsForSync(externalSystem, tenantName string) ([]entity.RawEmail, error)
	GetEmailsIdsForUserForSync(tenantName, userEmailAddress string) ([]entity.RawEmail, error)
	GetEmailForProcess(id uuid.UUID) (*entity.RawEmail, error)
	GetEmailForSyncByMessageId(tenant, usernameSource, messageId string) (*entity.RawEmail, error)
	UpdateRawEmailTable(id uuid.UUID, dbUpdateRecord entity.UpdateRawEmailTable) error
}

type rawEmailRepositoryImpl struct {
	gormDb *gorm.DB
}

func NewRawEmailRepository(gormDb *gorm.DB) RawEmailRepository {
	return &rawEmailRepositoryImpl{gormDb: gormDb}
}

func (repo *rawEmailRepositoryImpl) CountForUsername(ctx context.Context, externalSystem, tenant, username string) (int64, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "RawEmailRepository.CountForUsername")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	tracing.TagTenant(span, tenant)

	var result int64
	err := repo.gormDb.Model(entity.RawEmail{}).Where("external_system = ? AND tenant = ? AND username = ?", externalSystem, tenant, username).Count(&result).Error
	if err != nil {
		tracing.TraceErr(span, err)
		return 0, err
	}

	return result, nil
}

func (repo *rawEmailRepositoryImpl) GetByMessageId(ctx context.Context, externalSystem, tenant, username, messageId string) (*entity.RawEmail, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "RawEmailRepository.GetByMessageId")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	tracing.TagTenant(span, tenant)

	var result *entity.RawEmail
	err := repo.gormDb.Where("external_system = ? AND tenant = ? AND username = ? AND message_id = ?", externalSystem, tenant, username, messageId).First(&result).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return result, nil
}

func (repo *rawEmailRepositoryImpl) EmailExistsByMessageId(ctx context.Context, externalSystem, tenant, username, messageId string) (bool, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "RawEmailRepository.EmailExistsByMessageId")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	tracing.TagTenant(span, tenant)

	var result int64
	err := repo.gormDb.Model(entity.RawEmail{}).Where("external_system = ? AND tenant = ? AND username = ? AND message_id = ?", externalSystem, tenant, username, messageId).Count(&result).Error
	if err != nil {
		tracing.TraceErr(span, err)
		return false, err
	}

	return result > 0, nil
}

func (repo *rawEmailRepositoryImpl) Store(ctx context.Context, externalSystem, tenant, username, providerMessageId, messageId, rawEmail string, sentAt time.Time, state entity.EmailImportState) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "RawEmailRepository.Store")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	tracing.TagTenant(span, tenant)

	result := entity.RawEmail{}
	err := repo.gormDb.Find(&result, "external_system = ? AND tenant = ? AND username = ? AND message_id = ?", externalSystem, tenant, username, messageId).Error
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if result.Tenant != "" {
		err := errors.New("RawEmailRepository.Store - email already exists")
		tracing.TraceErr(span, err)
		return err
	}

	result.ProviderMessageId = providerMessageId
	result.MessageId = messageId

	result.CreatedAt = utils.Now()
	result.SentAt = sentAt
	result.State = state
	result.ExternalSystem = externalSystem
	result.Tenant = tenant
	result.Username = username
	result.Data = rawEmail
	result.Status = "PENDING"

	err = repo.gormDb.Save(&result).Error
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}

func (repo *rawEmailRepositoryImpl) GetEmailsIdsForSync(externalSystem, tenantName string) ([]entity.RawEmail, error) {
	result := []entity.RawEmail{}
	err := repo.gormDb.Order("sent_at desc").Select([]string{"id"}).Limit(25).Find(&result, "external_system = ? AND tenant = ? AND status = 'PENDING'", externalSystem, tenantName).Error
	if err != nil {
		logrus.Errorf("Failed getting rawEmails: %s; %s", externalSystem, tenantName)
		return nil, err
	}

	return result, nil
}

func (repo *rawEmailRepositoryImpl) GetEmailsIdsForUserForSync(tenantName, userSource string) ([]entity.RawEmail, error) {
	result := []entity.RawEmail{}
	err := repo.gormDb.Order("sent_at desc").Select([]string{"id", "external_system"}).Limit(100).Find(&result, "tenant = ? AND username = ? AND status = 'PENDING'", tenantName, userSource).Error
	if err != nil {
		logrus.Errorf("Failed getting rawEmails: %s; %s", tenantName, userSource)
		return nil, err
	}

	return result, nil
}

func (repo *rawEmailRepositoryImpl) GetEmailForProcess(id uuid.UUID) (*entity.RawEmail, error) {
	result := entity.RawEmail{}
	err := repo.gormDb.First(&result, id).Error
	if err != nil {
		logrus.Errorf("Failed getting rawEmail: %s", id)
		return nil, err
	}

	return &result, nil
}

func (repo *rawEmailRepositoryImpl) GetEmailForSyncByMessageId(tenant, usernameSource, messageId string) (*entity.RawEmail, error) {
	var result entity.RawEmail
	err := repo.gormDb.Where("tenant = ? AND username = ? AND message_id = ?", tenant, usernameSource, messageId).Find(&result).Error
	if err != nil {
		logrus.Errorf("GetEmailForSyncByMessageId - failed: %s; %s; %s", tenant, usernameSource, messageId)
		return nil, err
	}

	return &result, nil
}

func (repo *rawEmailRepositoryImpl) UpdateRawEmailTable(id uuid.UUID, dbUpdateRecord entity.UpdateRawEmailTable) error {
	tx := repo.gormDb.Model(&entity.RawEmail{}).Where("id = ?", id)

	tx.Update("status", dbUpdateRecord.EmailProcessingStatus)
	tx.Update("reason", dbUpdateRecord.Reason)
	if dbUpdateRecord.Error != nil {
		tx.Update("error", fmt.Sprintf("%v", dbUpdateRecord.Error))
	}
	if dbUpdateRecord.BouncedEmails != nil && len(*dbUpdateRecord.BouncedEmails) > 0 {
		jsonBytes, _ := json.Marshal(dbUpdateRecord.BouncedEmails)
		bounced := string(jsonBytes)
		tx.Update("bounced_emails", bounced)
	}

	err := tx.Error
	if err != nil {
		logrus.Errorf("Failed marking email as sent to event store: %v", id)
		return err
	}

	return nil
}
