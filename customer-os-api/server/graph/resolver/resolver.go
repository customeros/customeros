package resolver

import (
	"github.com/openline-ai/openline-customer-os/customer-os-api/service"
)

//go:generate go run github.com/99designs/gqlgen
// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	ServiceContainer *service.ServiceContainer
}

func NewResolver(serviceContainer *service.ServiceContainer) *Resolver {
	return &Resolver{ServiceContainer: serviceContainer}
}
