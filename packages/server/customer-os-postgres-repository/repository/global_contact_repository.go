package postgres_repository

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	"gorm.io/gorm"
)

type GlobalContactRepository interface {
	Create(ctx context.Context, contact *postgres_entity.GlobalContact) (*postgres_entity.GlobalContact, error)
	Update(ctx context.Context, contact *postgres_entity.GlobalContact) (*postgres_entity.GlobalContact, error)
	GetById(ctx context.Context, id uint64) (*postgres_entity.GlobalContact, error)
	GetByLinkedInIdentifier(ctx context.Context, linkedInIdentifier string) ([]*postgres_entity.GlobalContact, error)
	GetByLinkedInAlias(ctx context.Context, linkedInAlias string) ([]*postgres_entity.GlobalContact, error)
	GetByLinkedInIdentifierAndDomain(ctx context.Context, linkedInIdentifier string, primaryDomain string) (*postgres_entity.GlobalContact, error)
	GetByWorkEmail(ctx context.Context, workEmail string) (*postgres_entity.GlobalContact, error)
	GetByWorkEmailAndDomain(ctx context.Context, workEmail string, primaryDomain string) (*postgres_entity.GlobalContact, error)
	GetByPersonalEmail(ctx context.Context, personalEmail string) (*postgres_entity.GlobalContact, error)
	GetByPersonalEmailAndDomain(ctx context.Context, personalEmail string, primaryDomain string) (*postgres_entity.GlobalContact, error)
	GetContactsToFetchPhoto(ctx context.Context, limit int) ([]*postgres_entity.GlobalContact, error)
	SetProfilePhoto(ctx context.Context, id uint64, photoPath string) error
	SetDownloadStatus(ctx context.Context, id uint64, status enum.DownloadStatus) error
	GetContactsToFindWorkEmailWithBetterContact(ctx context.Context, retryAfterDays int, limit int) ([]*postgres_entity.GlobalContact, error)
	MarkBetterContactRequested(ctx context.Context, id uint64, betterContactRequestId string) error
	GetContactsToSetWorkEmailFromBetterContactResponse(ctx context.Context, limit int) ([]*postgres_entity.GlobalContact, error)
	MarkBetterContactSet(ctx context.Context, id uint64) error
	GetGlobalContactsToSyncIntoTenantContacts(ctx context.Context, daysFromPreviousSync, forceSyncAfterDays, limit int) ([]*postgres_entity.GlobalContact, error)
	GetGlobalContactsByEmailAddresses(ctx context.Context, emailAddresses []string) ([]*postgres_entity.GlobalContact, error)
}

type globalContactRepository struct {
	db *gorm.DB
}

func NewGlobalContactRepository(gormDb *gorm.DB) GlobalContactRepository {
	return &globalContactRepository{db: gormDb}
}

func (r *globalContactRepository) addSortingClauses(query *gorm.DB) *gorm.DB {
	return query.
		Order("CASE WHEN job_ended_at IS NULL THEN 0 ELSE 1 END"). // Prioritize records with null job_ended_at
		Order("job_ended_at DESC").                                // Then most recent job_ended_at
		Order("job_started_at DESC")                               // Then most recent job_started_at
}

func (r *globalContactRepository) GetById(ctx context.Context, id uint64) (*postgres_entity.GlobalContact, error) {
	spans, _ := telemetry.StartPostgresSpan(ctx, "GlobalContactRepository.GetById")
	defer spans.Finish()

	spans.TagEntity(strconv.FormatUint(id, 10))

	var contact postgres_entity.GlobalContact
	result := r.addSortingClauses(
		r.db.WithContext(ctx).Where("id = ?", id),
	).First(&contact)

	if result.Error != nil {
		spans.TraceError(result.Error)
		return nil, result.Error
	}

	return &contact, nil
}

func (r *globalContactRepository) Create(ctx context.Context, contact *postgres_entity.GlobalContact) (*postgres_entity.GlobalContact, error) {
	spans, _ := telemetry.StartPostgresSpan(ctx, "GlobalContactRepository.Create")
	defer spans.Finish()

	result := r.db.WithContext(ctx).Create(contact)
	if result.Error != nil {
		spans.TraceError(result.Error)
		return nil, result.Error
	}

	return contact, nil
}

func (r *globalContactRepository) Update(ctx context.Context, contact *postgres_entity.GlobalContact) (*postgres_entity.GlobalContact, error) {
	spans, _ := telemetry.StartPostgresSpan(ctx, "GlobalContactRepository.Update")
	defer spans.Finish()

	spans.TagEntity(strconv.FormatUint(contact.ID, 10))

	result := r.db.WithContext(ctx).Save(contact)
	if result.Error != nil {
		spans.TraceError(result.Error)
		return nil, result.Error
	}

	return contact, nil
}

func (r *globalContactRepository) GetByLinkedInIdentifier(ctx context.Context, linkedInIdentifier string) ([]*postgres_entity.GlobalContact, error) {
	spans, _ := telemetry.StartPostgresSpan(ctx, "GlobalContactRepository.GetByLinkedInIdentifier")
	defer spans.Finish()

	spans.LogKV("linkedInIdentifier", linkedInIdentifier)

	if linkedInIdentifier == "" {
		spans.LogKV("result.count", 0)
		return nil, nil
	}

	contacts := make([]*postgres_entity.GlobalContact, 0)
	result := r.addSortingClauses(
		r.db.WithContext(ctx).Where("linked_in_identifier = ?", linkedInIdentifier),
	).Find(&contacts)

	if result.Error != nil {
		spans.TraceError(result.Error)
		return nil, result.Error
	}

	spans.LogKV("result.count", len(contacts))
	return contacts, nil
}

func (r *globalContactRepository) GetByLinkedInAlias(ctx context.Context, linkedInAlias string) ([]*postgres_entity.GlobalContact, error) {
	spans, _ := telemetry.StartPostgresSpan(ctx, "GlobalContactRepository.GetByLinkedInAlias")
	defer spans.Finish()

	spans.LogKV("linkedInAlias", linkedInAlias)

	if linkedInAlias == "" {
		spans.LogKV("result.count", 0)
		return nil, nil
	}

	contacts := make([]*postgres_entity.GlobalContact, 0)
	result := r.addSortingClauses(
		r.db.WithContext(ctx).Where("linked_in_alias = ?", linkedInAlias),
	).Find(&contacts)

	if result.Error != nil {
		spans.TraceError(result.Error)
		return nil, result.Error
	}

	spans.LogKV("result.count", len(contacts))
	return contacts, nil
}

func (r *globalContactRepository) GetByLinkedInIdentifierAndDomain(ctx context.Context, linkedInIdentifier string, primaryDomain string) (*postgres_entity.GlobalContact, error) {
	spans, ctx := telemetry.StartPostgresSpan(ctx, "GlobalContactRepository.GetByLinkedInIdentifierAndDomain")
	defer spans.Finish()

	spans.LogKV("linkedInIdentifier", linkedInIdentifier, "primaryDomain", primaryDomain)

	var contact postgres_entity.GlobalContact
	result := r.addSortingClauses(
		r.db.WithContext(ctx).Where("linked_in_identifier = ? AND primary_domain = ?", linkedInIdentifier, primaryDomain),
	).First(&contact)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		spans.TraceError(result.Error)
		return nil, result.Error
	}

	return &contact, nil
}

func (r *globalContactRepository) GetByWorkEmail(ctx context.Context, workEmail string) (*postgres_entity.GlobalContact, error) {
	spans, _ := telemetry.StartPostgresSpan(ctx, "GlobalContactRepository.GetByWorkEmail")
	defer spans.Finish()

	spans.LogKV("workEmail", workEmail)

	var contact postgres_entity.GlobalContact
	result := r.addSortingClauses(
		r.db.WithContext(ctx).Where("work_email = ?", workEmail),
	).First(&contact)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		spans.TraceError(result.Error)
		return nil, result.Error
	}

	return &contact, nil
}

func (r *globalContactRepository) GetByWorkEmailAndDomain(ctx context.Context, workEmail string, primaryDomain string) (*postgres_entity.GlobalContact, error) {
	spans, _ := telemetry.StartPostgresSpan(ctx, "GlobalContactRepository.GetByWorkEmailAndDomain")
	defer spans.Finish()

	spans.LogKV("workEmail", workEmail, "primaryDomain", primaryDomain)

	var contact postgres_entity.GlobalContact
	result := r.addSortingClauses(
		r.db.WithContext(ctx).Where("work_email = ? AND primary_domain = ?", workEmail, primaryDomain),
	).First(&contact)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		spans.TraceError(result.Error)
		return nil, result.Error
	}

	return &contact, nil
}

func (r *globalContactRepository) GetByPersonalEmail(ctx context.Context, personalEmail string) (*postgres_entity.GlobalContact, error) {
	spans, _ := telemetry.StartPostgresSpan(ctx, "GlobalContactRepository.GetByPersonalEmail")
	defer spans.Finish()

	spans.LogKV("personalEmail", personalEmail)

	var contact postgres_entity.GlobalContact
	result := r.addSortingClauses(
		r.db.WithContext(ctx).Where("personal_email = ?", personalEmail),
	).First(&contact)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		spans.TraceError(result.Error)
		return nil, result.Error
	}

	return &contact, nil
}

func (r *globalContactRepository) GetByPersonalEmailAndDomain(ctx context.Context, personalEmail string, primaryDomain string) (*postgres_entity.GlobalContact, error) {
	spans, _ := telemetry.StartPostgresSpan(ctx, "GlobalContactRepository.GetByPersonalEmailAndDomain")
	defer spans.Finish()

	spans.LogKV("personalEmail", personalEmail, "primaryDomain", primaryDomain)

	var contact postgres_entity.GlobalContact
	result := r.addSortingClauses(
		r.db.WithContext(ctx).Where("personal_email = ? AND primary_domain = ?", personalEmail, primaryDomain),
	).First(&contact)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		spans.TraceError(result.Error)
		return nil, result.Error
	}

	return &contact, nil
}

func (r *globalContactRepository) GetContactsToFetchPhoto(ctx context.Context, limit int) ([]*postgres_entity.GlobalContact, error) {
	spans, ctx := telemetry.StartPostgresSpan(ctx, "GlobalContactRepository.GetContactsToFetchPhoto")
	defer spans.Finish()

	contacts := make([]*postgres_entity.GlobalContact, 0)
	result := r.db.WithContext(ctx).
		Where("download_status = ? OR download_status IS NULL", enum.DownloadNotStarted.String()).
		Where("profile_photo_external_url IS NOT NULL AND profile_photo_external_url != ''").
		Where("((linked_in_identifier IS NOT NULL AND linked_in_identifier != '') OR (work_email IS NOT NULL AND work_email != '') OR (personal_email IS NOT NULL AND personal_email != ''))").
		Order("created_at DESC").
		Limit(limit).
		Find(&contacts)
	if result.Error != nil {
		spans.TraceError(result.Error)
		return nil, result.Error
	}
	spans.LogKV("result.count", len(contacts))
	return contacts, nil
}

func (r *globalContactRepository) SetProfilePhoto(ctx context.Context, id uint64, photoPath string) error {
	spans, ctx := telemetry.StartPostgresSpan(ctx, "GlobalContactRepository.SetProfilePhoto")
	defer spans.Finish()

	result := r.db.WithContext(ctx).Model(&postgres_entity.GlobalContact{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"profile_photo_path": photoPath,
		})
	if result.Error != nil {
		spans.TraceError(result.Error)
		return result.Error
	}

	if result.RowsAffected == 0 {
		err := errors.New("no contact found with provided ID")
		spans.TraceError(err)
		return err
	}

	return nil
}

func (r *globalContactRepository) SetDownloadStatus(ctx context.Context, id uint64, status enum.DownloadStatus) error {
	spans, ctx := telemetry.StartPostgresSpan(ctx, "GlobalContactRepository.SetDownloadStatus")
	defer spans.Finish()

	spans.TagEntity(strconv.FormatUint(id, 10))
	spans.LogKV("status", status.String())

	result := r.db.WithContext(ctx).Model(&postgres_entity.GlobalContact{}).
		Where("id = ?", id).
		UpdateColumn("download_status", status.String())

	if result.Error != nil {
		spans.TraceError(result.Error)
		return result.Error
	}

	if result.RowsAffected == 0 {
		err := errors.New("no contact found with provided ID")
		spans.TraceError(err)
		return err
	}

	return nil
}

func (r *globalContactRepository) GetContactsToFindWorkEmailWithBetterContact(ctx context.Context, retryAfterDays, limit int) ([]*postgres_entity.GlobalContact, error) {
	spans, ctx := telemetry.StartPostgresSpan(ctx, "GlobalContactRepository.GetContactsToFindWorkEmailWithBetterContact")
	defer spans.Finish()

	spans.LogKV("retryAfterDays", retryAfterDays, "limit", limit)

	contacts := make([]*postgres_entity.GlobalContact, 0)
	result := r.db.WithContext(ctx).
		Where("work_email IS NULL OR work_email = ''").
		Where("primary_domain IS NOT NULL AND primary_domain != ''").
		Where("linked_in_identifier IS NOT NULL AND linked_in_identifier != ''").
		Where("job_ended_at IS NULL").
		Where("bettercontact_requested_at IS NULL OR bettercontact_requested_at < ?", utils.Now().Add(-time.Duration(retryAfterDays)*24*time.Hour)).
		Order("bettercontact_requested_at IS NULL DESC, COALESCE(bettercontact_requested_at, created_at) ASC").
		Limit(limit).
		Find(&contacts)
	if result.Error != nil {
		spans.TraceError(result.Error)
		return nil, result.Error
	}
	spans.LogKV("result.count", len(contacts))
	return contacts, nil
}

func (r *globalContactRepository) MarkBetterContactRequested(ctx context.Context, id uint64, betterContactRequestId string) error {
	spans, ctx := telemetry.StartPostgresSpan(ctx, "GlobalContactRepository.MarkBetterContactRequested")
	defer spans.Finish()

	spans.TagEntity(strconv.FormatUint(id, 10))
	spans.LogKV("betterContactRequestId", betterContactRequestId)

	result := r.db.WithContext(ctx).Model(&postgres_entity.GlobalContact{}).
		Where("id = ?", id).
		UpdateColumn("bettercontact_requested_at", utils.Now()).
		UpdateColumn("bettercontact_request_id", betterContactRequestId).
		UpdateColumn("bettercontact_set_at", nil).
		UpdateColumn("bettercontact_check_response_at", nil)

	if result.Error != nil {
		spans.TraceError(result.Error)
		return result.Error
	}

	if result.RowsAffected == 0 {
		err := errors.New("no contact found with provided ID")
		spans.TraceError(err)
		return err
	}

	return nil
}

func (r *globalContactRepository) GetContactsToSetWorkEmailFromBetterContactResponse(ctx context.Context, limit int) ([]*postgres_entity.GlobalContact, error) {
	spans, ctx := telemetry.StartPostgresSpan(ctx, "GlobalContactRepository.GetContactsToSetWorkEmailFromBetterContactResponse")
	defer spans.Finish()

	contacts := make([]*postgres_entity.GlobalContact, 0)
	result := r.db.WithContext(ctx).
		Where("bettercontact_requested_at IS NOT NULL").
		Where("bettercontact_requested_at < ?", utils.Now().Add(-90*time.Second)).
		Where("bettercontact_set_at IS NULL").
		Where("bettercontact_request_id IS NOT NULL AND bettercontact_request_id != ''").
		Order("CASE WHEN bettercontact_check_response_at IS NULL THEN 0 ELSE 1 END ASC").
		Order("bettercontact_check_response_at ASC").
		Limit(limit).
		Find(&contacts)
	if result.Error != nil {
		spans.TraceError(result.Error)
		return nil, result.Error
	}

	// Update bettercontact_check_response_at for found contacts
	if len(contacts) > 0 {
		var ids []uint64
		for _, contact := range contacts {
			ids = append(ids, contact.ID)
		}

		result := r.db.WithContext(ctx).Model(&postgres_entity.GlobalContact{}).
			Where("id IN ?", ids).
			UpdateColumn("bettercontact_check_response_at", utils.Now())

		if result.Error != nil {
			spans.TraceError(result.Error)
			return contacts, result.Error
		}
	}

	spans.LogKV("result.count", len(contacts))
	return contacts, nil
}

func (r *globalContactRepository) MarkBetterContactSet(ctx context.Context, id uint64) error {
	spans, ctx := telemetry.StartPostgresSpan(ctx, "GlobalContactRepository.MarkBetterContactSet")
	defer spans.Finish()

	spans.TagEntity(strconv.FormatUint(id, 10))

	result := r.db.WithContext(ctx).Model(&postgres_entity.GlobalContact{}).
		Where("id = ?", id).
		UpdateColumn("bettercontact_set_at", utils.Now())

	if result.Error != nil {
		spans.TraceError(result.Error)
		return result.Error
	}

	return nil
}

func (r *globalContactRepository) GetGlobalContactsToSyncIntoTenantContacts(ctx context.Context, daysFromPreviousSync, forceSyncAfterDays, limit int) ([]*postgres_entity.GlobalContact, error) {
	spans, _ := telemetry.StartPostgresSpan(ctx, "GlobalContactRepository.GetGlobalContactsToSyncIntoTenantContacts")
	defer spans.Finish()

	spans.LogKV("daysFromPreviousSync", daysFromPreviousSync, "forceSyncAfterDays", forceSyncAfterDays, "limit", limit)

	contacts := make([]*postgres_entity.GlobalContact, 0)
	result := r.db.WithContext(ctx).
		Where("(synced_to_neo_at IS NULL OR (synced_to_neo_at < ? AND updated_at > synced_to_neo_at) OR (synced_to_neo_at < ?))", utils.Now().Add(-24*time.Hour*time.Duration(daysFromPreviousSync)), utils.Now().Add(-24*time.Hour*time.Duration(forceSyncAfterDays))).
		Order("synced_to_neo_at IS NULL DESC, COALESCE(synced_to_neo_at, created_at) ASC").
		Limit(limit).
		Find(&contacts)
	if result.Error != nil {
		spans.TraceError(result.Error)
		return nil, result.Error
	}
	spans.LogKV("result.count", len(contacts))
	return contacts, nil
}

func (r *globalContactRepository) GetGlobalContactsByEmailAddresses(ctx context.Context, emailAddresses []string) ([]*postgres_entity.GlobalContact, error) {
	spans, ctx := telemetry.StartPostgresSpan(ctx, "GlobalContactRepository.GetGlobalContactsByEmailAddresses")
	defer spans.Finish()

	spans.LogKV("emailAddresses", emailAddresses)

	contacts := make([]*postgres_entity.GlobalContact, 0)
	result := r.db.WithContext(ctx).
		Where("personal_email IN ? OR work_email IN ?", emailAddresses, emailAddresses).
		Find(&contacts)
	if result.Error != nil {
		spans.TraceError(result.Error)
		return nil, result.Error
	}
	spans.LogKV("result.count", len(contacts))
	return contacts, nil
}
