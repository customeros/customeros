package repository

import (
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/source/zendesk_support/entity"
	"gorm.io/gorm"
	"time"
)

func GetUsers(db *gorm.DB, limit int, runId string) (entity.Users, error) {
	var users entity.Users

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
		Find(&users).Error

	if err != nil {
		return nil, err
	}
	return users, nil
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

func GetOrganizations(db *gorm.DB, limit int, runId string) (entity.Organizations, error) {
	var organizations entity.Organizations

	cte := `
		WITH UpToDateData AS (
   		SELECT row_number() OVER (PARTITION BY id ORDER BY updated_at DESC) AS row_num, *
   		FROM organizations
		)`
	err := db.
		Raw(cte+" SELECT u.* FROM UpToDateData u left join openline_sync_status_organizations s "+
			" on u.id = s.id and u._airbyte_ab_id = s._airbyte_ab_id and u._airbyte_organizations_hashid = s._airbyte_organizations_hashid "+
			" WHERE u.row_num = ? "+
			" and (s.synced_to_customer_os is null or s.synced_to_customer_os = ?) "+
			" and (s.synced_to_customer_os_attempt is null or s.synced_to_customer_os_attempt < ?) "+
			" and (s.run_id is null or s.run_id <> ?) "+
			" limit ?", 1, false, 10, runId, limit).
		Find(&organizations).Error

	if err != nil {
		return nil, err
	}
	return organizations, nil
}

func MarkOrganizationProcessed(db *gorm.DB, organization entity.Organization, synced bool, runId string) error {
	syncStatusOrganization := entity.SyncStatusOrganization{
		Id:                         organization.Id,
		AirbyteAbId:                organization.AirbyteAbId,
		AirbyteOrganizationsHashid: organization.AirbyteOrganizationsHashid,
	}
	db.FirstOrCreate(&syncStatusOrganization, syncStatusOrganization)

	return db.Model(&syncStatusOrganization).
		Where(&entity.SyncStatusOrganization{Id: organization.Id, AirbyteAbId: organization.AirbyteAbId, AirbyteOrganizationsHashid: organization.AirbyteOrganizationsHashid}).
		Updates(entity.SyncStatusOrganization{
			SyncedToCustomerOs: synced,
			SyncedAt:           time.Now(),
			SyncAttempt:        syncStatusOrganization.SyncAttempt + 1,
			RunId:              runId,
		}).Error
}

func GetUsersAsOrganizations(db *gorm.DB, limit int, runId string) (entity.UsersAsOrganizations, error) {
	var usersAsOrganizations entity.UsersAsOrganizations

	cte := `
		WITH UpToDateData AS (
   		SELECT row_number() OVER (PARTITION BY id ORDER BY updated_at DESC) AS row_num, *
   		FROM users
		)`
	err := db.
		Raw(cte+" SELECT u.* FROM UpToDateData u left join openline_sync_status_users_as_organizations s "+
			" on u.id = s.id and u._airbyte_ab_id = s._airbyte_ab_id and u._airbyte_users_hashid = s._airbyte_users_hashid "+
			" WHERE u.row_num = ? "+
			" and (u.role = ?) "+
			" and (s.synced_to_customer_os is null or s.synced_to_customer_os = ?) "+
			" and (s.synced_to_customer_os_attempt is null or s.synced_to_customer_os_attempt < ?) "+
			" and (s.run_id is null or s.run_id <> ?) "+
			" limit ?", 1, "end-user", false, 10, runId, limit).
		Find(&usersAsOrganizations).Error

	if err != nil {
		return nil, err
	}
	return usersAsOrganizations, nil
}

func MarkUserAsOrganizationProcessed(db *gorm.DB, userAsOrganization entity.UserAsOrganization, synced bool, runId string) error {
	syncStatusUserAsOrganization := entity.SyncStatusUserAsOrganization{
		Id:                 userAsOrganization.Id,
		AirbyteAbId:        userAsOrganization.AirbyteAbId,
		AirbyteUsersHashid: userAsOrganization.AirbyteUsersHashid,
	}
	db.FirstOrCreate(&syncStatusUserAsOrganization, syncStatusUserAsOrganization)

	return db.Model(&syncStatusUserAsOrganization).
		Where(&entity.SyncStatusUserAsOrganization{Id: userAsOrganization.Id, AirbyteAbId: userAsOrganization.AirbyteAbId, AirbyteUsersHashid: userAsOrganization.AirbyteUsersHashid}).
		Updates(entity.SyncStatusUserAsOrganization{
			SyncedToCustomerOs: synced,
			SyncedAt:           time.Now(),
			SyncAttempt:        syncStatusUserAsOrganization.SyncAttempt + 1,
			RunId:              runId,
		}).Error
}

func GetTickets(db *gorm.DB, limit int, runId string) (entity.Tickets, error) {
	var tickets entity.Tickets

	cte := `
		WITH UpToDateData AS (
   		SELECT row_number() OVER (PARTITION BY id ORDER BY updated_at DESC) AS row_num, *
   		FROM tickets
		)`
	err := db.
		Raw(cte+" SELECT u.* FROM UpToDateData u left join openline_sync_status_tickets s "+
			" on u.id = s.id and u._airbyte_ab_id = s._airbyte_ab_id and u._airbyte_tickets_hashid = s._airbyte_tickets_hashid "+
			" WHERE u.row_num = ? "+
			" and (s.synced_to_customer_os is null or s.synced_to_customer_os = ?) "+
			" and (s.synced_to_customer_os_attempt is null or s.synced_to_customer_os_attempt < ?) "+
			" and (s.run_id is null or s.run_id <> ?) "+
			" limit ?", 1, false, 10, runId, limit).
		Find(&tickets).Error

	if err != nil {
		return nil, err
	}
	return tickets, nil
}

func MarkTicketProcessed(db *gorm.DB, ticket entity.Ticket, synced bool, runId string) error {
	syncStatusTicket := entity.SyncStatusTicket{
		Id:                   ticket.Id,
		AirbyteAbId:          ticket.AirbyteAbId,
		AirbyteTicketsHashid: ticket.AirbyteTicketsHashid,
	}
	db.FirstOrCreate(&syncStatusTicket, syncStatusTicket)

	return db.Model(&syncStatusTicket).
		Where(&entity.SyncStatusTicket{Id: ticket.Id, AirbyteAbId: ticket.AirbyteAbId, AirbyteTicketsHashid: ticket.AirbyteTicketsHashid}).
		Updates(entity.SyncStatusTicket{
			SyncedToCustomerOs: synced,
			SyncedAt:           time.Now(),
			SyncAttempt:        syncStatusTicket.SyncAttempt + 1,
			RunId:              runId,
		}).Error
}

func GetTicketComments(db *gorm.DB, limit int, runId string) (entity.TicketComments, error) {
	var ticketComments entity.TicketComments

	cte := `
		WITH UpToDateData AS (
   		SELECT row_number() OVER (PARTITION BY id ORDER BY created_at DESC) AS row_num, *
   		FROM ticket_comments
		)`
	err := db.
		Raw(cte+" SELECT u.* FROM UpToDateData u left join openline_sync_status_ticket_comments s "+
			" on u.id = s.id and u._airbyte_ab_id = s._airbyte_ab_id and u._airbyte_ticket_comments_hashid = s._airbyte_ticket_comments_hashid "+
			" WHERE u.row_num = ? "+
			" and (s.synced_to_customer_os is null or s.synced_to_customer_os = ?) "+
			" and (s.synced_to_customer_os_attempt is null or s.synced_to_customer_os_attempt < ?) "+
			" and (s.run_id is null or s.run_id <> ?) "+
			" limit ?", 1, false, 10, runId, limit).
		Find(&ticketComments).Error

	if err != nil {
		return nil, err
	}
	return ticketComments, nil
}

func MarkTicketCommentProcessed(db *gorm.DB, ticketComment entity.TicketComment, synced bool, runId string) error {
	syncStatusTicketComment := entity.SyncStatusTicketComment{
		Id:                          ticketComment.Id,
		AirbyteAbId:                 ticketComment.AirbyteAbId,
		AirbyteTicketCommentsHashid: ticketComment.AirbyteTicketCommentsHashid,
	}
	db.FirstOrCreate(&syncStatusTicketComment, syncStatusTicketComment)

	return db.Model(&syncStatusTicketComment).
		Where(&entity.SyncStatusTicketComment{Id: ticketComment.Id, AirbyteAbId: ticketComment.AirbyteAbId, AirbyteTicketCommentsHashid: ticketComment.AirbyteTicketCommentsHashid}).
		Updates(entity.SyncStatusTicketComment{
			SyncedToCustomerOs: synced,
			SyncedAt:           time.Now(),
			SyncAttempt:        syncStatusTicketComment.SyncAttempt + 1,
			RunId:              runId,
		}).Error
}
