package config

type ConcurrencyConfig struct {
	UserSyncConcurrency             int `env:"USER_SYNC_CONCURRENCY" envDefault:"1" validate:"required,min=1,max=100"`
	ContactSyncConcurrency          int `env:"CONTACT_SYNC_CONCURRENCY" envDefault:"1" validate:"required,min=1,max=100"`
	OrganizationSyncConcurrency     int `env:"ORGANIZATION_SYNC_CONCURRENCY" envDefault:"1" validate:"required,min=1,max=100"`
	IssueSyncConcurrency            int `env:"ISSUE_SYNC_CONCURRENCY" envDefault:"1" validate:"required,min=1,max=100"`
	CommentSyncConcurrency          int `env:"COMMENT_SYNC_CONCURRENCY" envDefault:"1" validate:"required,min=1,max=100"`
	LogEntrySyncConcurrency         int `env:"LOG_ENTRY_SYNC_CONCURRENCY" envDefault:"1" validate:"required,min=1,max=100"`
	InteractionEventSyncConcurrency int `env:"INTERACTION_EVENT_SYNC_CONCURRENCY" envDefault:"1" validate:"required,min=1,max=100"`
}
