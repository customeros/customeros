package service

import "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"

type LinkWith struct {
	Type         model.EntityType
	Id           string
	Relationship string
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
