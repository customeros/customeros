package resolver

import (
	"github.com.openline-ai.customer-os-analytics-api/repository"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	RepositoryHandler *repository.RepositoryHandler
}

func NewResolver(repositoryHandler *repository.RepositoryHandler) *Resolver {
	return &Resolver{RepositoryHandler: repositoryHandler}
}
