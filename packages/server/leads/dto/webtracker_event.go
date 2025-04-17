package dto

import (
	"time"

	"github.com/customeros/customeros/packages/server/leads/internal/enum"
)

type WebTrackerEvent struct {
	ID               string               `json:"id"`
	VisitorID        string               `json:"visitorId"`
	IP               string               `json:"ip" `
	EventType        enum.WebTrackerEvent `json:"eventType"`
	EventData        string               `json:"eventData"`
	Timestamp        time.Time            `json:"timestamp"`
	Href             string               `json:"href"`
	Referrer         string               `json:"referrer"`
	UserAgent        string               `json:"userAgent"`
	Language         string               `json:"language"`
	CookiesEnabled   bool                 `json:"cookiesEnabled"`
	ScreenResolution string               `json:"screenResolution"`
}
