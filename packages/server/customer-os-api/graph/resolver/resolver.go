package resolver

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/service/container"
)

//go:generate go get github.com/99designs/gqlgen@v0.17.20
//go:generate go run github.com/99designs/gqlgen
// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	ServiceContainer *container.ServiceContainer
}

func NewResolver(serviceContainer *container.ServiceContainer) *Resolver {
	return &Resolver{ServiceContainer: serviceContainer}
}
