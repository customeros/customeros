package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
)

func MapEntityToAnalysis(entity *entity.AnalysisEntity) *model.Analysis {
	return &model.Analysis{
		ID:            entity.Id,
		CreatedAt:     *entity.CreatedAt,
		Content:       utils.StringPtrNillable(entity.Content),
		ContentType:   utils.StringPtrNillable(entity.ContentType),
		AnalysisType:  utils.StringPtrNillable(entity.AnalysisType),
		Source:        MapDataSourceToModel(entity.Source),
		SourceOfTruth: MapDataSourceToModel(entity.SourceOfTruth),
		AppSource:     entity.AppSource,
	}
}

func MapEntitiesToAnalysis(entities *entity.AnalysisEntities) []*model.Analysis {
	var analyses []*model.Analysis
	for _, interactionEventEntity := range *entities {
		analyses = append(analyses, MapEntityToAnalysis(&interactionEventEntity))
	}
	return analyses
}
func MapAnalysisInputToEntity(input *model.AnalysisInput) *entity.AnalysisEntity {
	return &entity.AnalysisEntity{
		AnalysisType: utils.IfNotNilString(input.AnalysisType),
		Content:      utils.IfNotNilString(input.Content),
		ContentType:  utils.IfNotNilString(input.ContentType),
		AppSource:    input.AppSource,
	}
}
