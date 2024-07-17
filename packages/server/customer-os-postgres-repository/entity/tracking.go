package entity

import "time"

type TrackingIdentificationState string

const (
	TrackingIdentificationStateError               TrackingIdentificationState = "ERROR"                // tracking record processing error
	TrackingIdentificationStateNew                 TrackingIdentificationState = "NEW"                  // New tracking record
	TrackingIdentificationStatePrefilteredPass     TrackingIdentificationState = "PREFILTER_PASS"       // tracking record passed the IPData prefilter
	TrackingIdentificationStatePrefilteredFail     TrackingIdentificationState = "PREFILTER_FAIL"       // tracking record failed the IPData prefilter
	TrackingIdentificationStateIdentified          TrackingIdentificationState = "IDENTIFIED"           // tracking record identified with scraping
	TrackingIdentificationStateNotIdentified       TrackingIdentificationState = "NOT_IDENTIFIED"       // tracking record not identified with scraping
	TrackingIdentificationStateOrganizationCreated TrackingIdentificationState = "ORGANIZATION_CREATED" // organization created for tracking record
	TrackingIdentificationStateOrganizationExists  TrackingIdentificationState = "ORGANIZATION_EXISTS"  // organization already exists for tracking record
)

type Tracking struct {
	ID        string    `gorm:"primary_key;type:uuid;default:gen_random_uuid()" json:"id"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`

	Tenant string `gorm:"column:tenant;type:varchar(255);" json:"tenant"`

	UserId    string `gorm:"column:user_id;type:varchar(255);NOT NULL;" json:"userId"`
	IP        string `gorm:"column:ip;type:varchar(255);" json:"ip" `
	EventType string `gorm:"column:event_type;type:varchar(255);" json:"eventType"`
	EventData string `gorm:"column:event_data;type:text;" json:"eventData"`
	Timestamp int    `gorm:"column:timestamp;type:bigint;" json:"timestamp"`

	Href             string `gorm:"column:href;type:varchar(255);" json:"href"`
	Origin           string `gorm:"column:origin;type:varchar(255);" json:"origin"`
	Search           string `gorm:"column:search;type:varchar(255);" json:"search"`
	Hostname         string `gorm:"column:hostname;type:varchar(255);" json:"hostname"`
	Pathname         string `gorm:"column:pathname;type:varchar(255);" json:"pathname"`
	Referrer         string `gorm:"column:referrer;type:varchar(255);" json:"referrer"`
	UserAgent        string `gorm:"column:user_agent;type:text;" json:"userAgent"`
	Language         string `gorm:"column:language;type:varchar(255);" json:"language"`
	CookiesEnabled   bool   `gorm:"column:cookies_enabled;type:boolean;" json:"cookiesEnabled"`
	ScreenResolution string `gorm:"column:screen_resolution;type:varchar(255);" json:"screenResolution"`

	State            TrackingIdentificationState `gorm:"column:state;type:varchar(50);" json:"state"`
	OrganizationId   *string                     `gorm:"column:organization_id;type:varchar(255);" json:"organizationId"`
	OrganizationName *string                     `gorm:"column:organization_name;type:varchar(255);" json:"organizationName"`
	Notified         bool                        `gorm:"column:notified;type:boolean;default:false" json:"notified"`
}

func (Tracking) TableName() string {
	return "tracking"
}
