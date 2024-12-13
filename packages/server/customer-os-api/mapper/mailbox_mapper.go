package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	postgresEntity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
)

func MapEntityToMailbox(entity *postgresEntity.TenantSettingsMailbox) *model.Mailbox {
	return &model.Mailbox{
		Domain:          entity.Domain,
		Mailbox:         entity.MailboxUsername,
		Created:         entity.CreatedAt,
		UserID:          &entity.UserId,
		RampUpCurrent:   entity.RampUpCurrent,
		RampUpMax:       entity.RampUpMax,
		RampUpRate:      entity.RampUpRate,
		CurrentFlowIds:  []string{},
		ScheduledEmails: 0,
	}
}
