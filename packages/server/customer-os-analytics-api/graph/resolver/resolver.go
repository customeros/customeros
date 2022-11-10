package resolver

import (
	"github.com.openline-ai.customer-os-analytics-api/repository"
)

//go:generate go run github.com/99designs/gqlgen
// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	RepositoryContainer *repository.RepositoryContainer
}

func NewResolver(repositoryContainer *repository.RepositoryContainer) *Resolver {
	return &Resolver{RepositoryContainer: repositoryContainer}
}
