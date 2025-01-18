package common_srv

import "github.com/customeros/customeros/packages/server/customer-os-common-module/model"

type LinkWith struct {
	Type         model.EntityType `json:"type"`
	Id           string           `json:"id"`
	Relationship string           `json:"relationship"`
}

func (l LinkWith) IsContact() bool {
	return l.Type == model.CONTACT
}

func (l LinkWith) IsOrganization() bool {
	return l.Type == model.ORGANIZATION
}

type ServiceOptions struct {
	SkipCompletedEvents bool
}

func PublishCompletedEvents(options ...ServiceOptions) bool {
	if options == nil {
		return true
	}
	if len(options) == 0 {
		return true
	}
	return !options[0].SkipCompletedEvents
}
