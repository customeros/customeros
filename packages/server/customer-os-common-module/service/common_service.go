package service

import "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"

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

type commonService struct {
	services *Services
}

type CommonService interface {
}

func NewCommonService(services *Services) CommonService {
	return &commonService{
		services: services,
	}
}
