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
	Services *container.Services
}

func NewResolver(serviceContainer *container.Services) *Resolver {
	return &Resolver{Services: serviceContainer}
}
