package models

import (
	"time"

	"gorm.io/gorm"

	"github.com/customeros/customeros/packages/server/leads/internal/enum"
)

type WebTrackerEvent struct {
	ID           string            `gorm:"column:id;type:varchar(50);primaryKey;not null" json:"id"`
	Event        enum.Events       `gorm:"column:event;type:varchar(50);index;not null" json:"event"`
	Publisher    enum.LeadsService `gorm:"column:publisher;type:varchar(50);index;not null" json:"publisher"`
	Timestamp    time.Time         `gorm:"not null;index"`
	Tenant       string            `gorm:"column:tenant;type:varchar(50);index;not null" json:"tenant"`
	TrackerID    string            `gorm:"column:tracker_id;type:varchar(50);index;not null" json:"trackerId"`
	SessionID    string            `gorm:"column:session_id;type:varchar(50);index" json:"sessionId"`
	Payload      []byte            `gorm:"column:payload;type:bytea" json:"-"`
	HasError     bool              `gorm:"column:has_error;type:boolean" json:"hasError"`
	ErrorMessage string            `gorm:"column:error_message;type:varchar(255)" json:"errorMessage"`
}

// TableName overrides the table name
func (WebTrackerEvent) TableName() string {
	return "webtracker_events"
}

// BeforeCreate hook to ensure the timestamp is set
func (e *WebTrackerEvent) BeforeCreate(tx *gorm.DB) error {
	if e.Timestamp.IsZero() {
		e.Timestamp = time.Now()
	}
	return nil
}

// SetupTimescaleDB initializes the TimescaleDB specifics for this model
func SetupTimescaleDB(db *gorm.DB) error {
	// Migrate the schema
	if err := db.AutoMigrate(&WebTrackerEvent{}); err != nil {
		return err
	}

	// Convert to hypertable - this only needs to be done once
	if err := db.Exec(`SELECT create_hypertable('web_events_source', 'timestamp', 
		chunk_time_interval => INTERVAL '1 day',
		if_not_exists => TRUE)`).Error; err != nil {
		return err
	}

	// Create continuous aggregate for hourly web statistics
	if err := db.Exec(`
		CREATE MATERIALIZED VIEW IF NOT EXISTS hourly_web_stats
		WITH (timescaledb.continuous) AS
		SELECT
			time_bucket('1 hour', timestamp) AS hour,
			tenant,
			publisher,
			event,
			session_id,
			has_error,
			count(*) AS event_count
		FROM web_events_source
		GROUP BY hour, tenant, publisher, event, session_id, has_error
	`).Error; err != nil {
		return err
	}

	// Add refresh policy for continuous aggregate
	if err := db.Exec(`
		SELECT add_continuous_aggregate_policy('hourly_web_stats',
			start_offset => INTERVAL '1 day',
			end_offset => INTERVAL '1 hour',
			schedule_interval => INTERVAL '1 hour',
			if_not_exists => TRUE)
	`).Error; err != nil {
		return err
	}

	// Create additional indexes for common query patterns
	if err := db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_web_events_has_error_timestamp ON web_events_source (has_error, timestamp DESC);
		CREATE INDEX IF NOT EXISTS idx_web_events_publisher_event ON web_events_source (publisher, event);
		CREATE INDEX IF NOT EXISTS idx_web_events_session_timestamp ON web_events_source (session_id, timestamp DESC);
	`).Error; err != nil {
		return err
	}
	return nil
}
