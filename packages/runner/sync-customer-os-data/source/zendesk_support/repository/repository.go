package repository

import (
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/source/zendesk_support/entity"
	"gorm.io/gorm"
	"time"
)

func GetUsers(db *gorm.DB, limit int, runId string) (entity.Users, error) {
	var owners entity.Users

	cte := `
		WITH UpToDateData AS (
   		SELECT row_number() OVER (PARTITION BY id ORDER BY updated_at DESC) AS row_num, *
   		FROM users
		)`
	err := db.
		Raw(cte+" SELECT u.* FROM UpToDateData u left join openline_sync_status_users s "+
			" on u.id = s.id and u._airbyte_ab_id = s._airbyte_ab_id and u._airbyte_users_hashid = s._airbyte_users_hashid "+
			" WHERE u.row_num = ? "+
			" and (u.role <> ?) "+
			" and (s.synced_to_customer_os is null or s.synced_to_customer_os = ?) "+
			" and (s.synced_to_customer_os_attempt is null or s.synced_to_customer_os_attempt < ?) "+
			" and (s.run_id is null or s.run_id <> ?) "+
			" limit ?", 1, "end-user", false, 10, runId, limit).
		Find(&owners).Error

	if err != nil {
		return nil, err
	}
	return owners, nil
}

func MarkUserProcessed(db *gorm.DB, user entity.User, synced bool, runId string) error {
	syncStatusUser := entity.SyncStatusUser{
		Id:                 user.Id,
		AirbyteAbId:        user.AirbyteAbId,
		AirbyteUsersHashid: user.AirbyteUsersHashid,
	}
	db.FirstOrCreate(&syncStatusUser, syncStatusUser)

	return db.Model(&syncStatusUser).
		Where(&entity.SyncStatusUser{Id: user.Id, AirbyteAbId: user.AirbyteAbId, AirbyteUsersHashid: user.AirbyteUsersHashid}).
		Updates(entity.SyncStatusUser{
			SyncedToCustomerOs: synced,
			SyncedAt:           time.Now(),
			SyncAttempt:        syncStatusUser.SyncAttempt + 1,
			RunId:              runId,
		}).Error
}
