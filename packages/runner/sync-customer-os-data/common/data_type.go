package common

type SyncedEntityType string

const (
	USERS              SyncedEntityType = "users"
	CONTACTS           SyncedEntityType = "contacts"
	ORGANIZATIONS      SyncedEntityType = "organizations"
	NOTES              SyncedEntityType = "notes"
	EMAIL_MESSAGES     SyncedEntityType = "email_messages"
	ISSUES             SyncedEntityType = "issues"
	LOG_ENTRIES        SyncedEntityType = "log_entries"
	MEETINGS           SyncedEntityType = "meetings"
	INTERACTION_EVENTS SyncedEntityType = "interaction_events"
)
