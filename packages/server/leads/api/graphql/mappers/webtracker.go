package mappers

import (
	"github.com/customeros/customeros/packages/server/leads/api/graphql/graphql_model"
	"github.com/customeros/customeros/packages/server/leads/internal/models"
)

func ToWebtrackerGraphQL(w *models.WebTracker) *graphql_model.Webtracker {
	return &graphql_model.Webtracker{
		ID:                w.ID,
		Domain:            w.Domain,
		CnameHost:         w.CNAMEHost,
		CnameTarget:       w.CNAMETarget,
		IsCnameConfigured: w.IsCNAMEConfigured,
		IsProxyActive:     w.IsProxyActive,
		IsArchived:        w.IsArchived,
		LastEventAt:       *w.LastEventAt,
		CreatedAt:         w.CreatedAt,
		UpdatedAt:         w.UpdatedAt,
	}
}

func ToWebtrackerDBModel(dto *graphql_model.Webtracker) *models.WebTracker {
	return &models.WebTracker{
		ID:                dto.ID,
		Domain:            dto.Domain,
		CNAMEHost:         dto.CnameHost,
		CNAMETarget:       dto.CnameTarget,
		IsCNAMEConfigured: dto.IsCnameConfigured,
		IsProxyActive:     dto.IsProxyActive,
		IsArchived:        dto.IsArchived,
		LastEventAt:       &dto.LastEventAt,
		CreatedAt:         dto.CreatedAt,
		UpdatedAt:         dto.UpdatedAt,
	}
}

func ToWebtrackerGraphQLSlice(models []models.WebTracker) []*graphql_model.Webtracker {
	dtos := make([]*graphql_model.Webtracker, len(models))
	for i, model := range models {
		dtos[i] = ToWebtrackerGraphQL(&model)
	}
	return dtos
}

func ToWebtrackerDBModelSlice(dtos []*graphql_model.Webtracker) []models.WebTracker {
	models := make([]models.WebTracker, len(dtos))
	for i, dto := range dtos {
		models[i] = *ToWebtrackerDBModel(dto)
	}
	return models
}
