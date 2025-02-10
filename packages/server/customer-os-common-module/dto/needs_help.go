package dto

import (
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
)

type CompanyNeedsHelp struct{}

func (e CompanyNeedsHelp) Name() enum.AgentListenerEvent {
	return enum.EventCompanyNeedsHelp
}
