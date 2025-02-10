package dto

import (
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
)

var _ interfaces.AgentEvents = CompanyNeedsHelp{}

type CompanyNeedsHelp struct{}

func (e CompanyNeedsHelp) Name() enum.AgentListenerEvent {
	return enum.EventCompanyNeedsHelp
}
