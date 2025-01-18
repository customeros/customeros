package mapper

import (
	postgresEntity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"

	"github.com/customeros/customeros/packages/server/customer-os-api/graphql/model"
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
