package country

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events/events"
	"time"
)

type Country struct {
	ID           string        `json:"id"`
	CreatedAt    time.Time     `json:"createdAt"`
	UpdatedAt    time.Time     `json:"updatedAt"`
	SourceFields events.Source `json:"source"`

	Name      string `json:"name"`
	CodeA2    string `json:"codeA2"`
	CodeA3    string `json:"codeA3"`
	PhoneCode string `json:"phoneCode"`
}
