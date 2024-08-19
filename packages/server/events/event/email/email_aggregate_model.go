package email

import (
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/events/event/common"
	"time"
)

type Email struct {
	ID        string        `json:"id"`
	RawEmail  string        `json:"rawEmail"`
	Email     string        `json:"email"`
	Source    common.Source `json:"source"`
	CreatedAt time.Time     `json:"createdAt"`
	UpdatedAt time.Time     `json:"updatedAt"`
}

func (p *Email) String() string {
	return fmt.Sprintf("Email{ID: %s, RawEmail: %s, Email: %s, Source: %s, CreatedAt: %s, UpdatedAt: %s}", p.ID, p.RawEmail, p.Email, p.Source, p.CreatedAt, p.UpdatedAt)
}
