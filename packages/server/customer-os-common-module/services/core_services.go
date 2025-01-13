package service

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/interfaces"
)

type CoreServices struct {
	EmailService            interfaces.EmailService
	JobRoleService          interfaces.JobRoleService
	SocialService           interfaces.SocialService
	ContactService          interfaces.ContactService
	OrganizationService     interfaces.OrganizationService
	InteractionEventService interfaces.InteractionEventService
	MailService             interfaces.MailService
}
