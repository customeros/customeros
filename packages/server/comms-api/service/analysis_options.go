package service

import "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"

type AnalysisOption func(*AnalysisOptions)

type AnalysisOptions struct {
	analysisType *string
	content      *string
	contentType  *string
	appSource    *string
	tenant       *string
	username     *string
	describes    *model.AnalysisDescriptionInput
}

func WithAnalysisType(value *string) AnalysisOption {
	return func(options *AnalysisOptions) {
		options.analysisType = value
	}
}

func WithAnalysisContent(value *string) AnalysisOption {
	return func(options *AnalysisOptions) {
		options.content = value
	}
}

func WithAnalysisContentType(value *string) AnalysisOption {
	return func(options *AnalysisOptions) {
		options.contentType = value
	}
}

func WithAnalysisAppSource(value *string) AnalysisOption {
	return func(options *AnalysisOptions) {
		options.appSource = value
	}
}

func WithAnalysisTenant(value *string) AnalysisOption {
	return func(options *AnalysisOptions) {
		options.tenant = value
	}
}

func WithAnalysisUsername(value *string) AnalysisOption {
	return func(options *AnalysisOptions) {
		options.username = value
	}
}

func WithAnalysisDescribes(value *model.AnalysisDescriptionInput) AnalysisOption {
	return func(options *AnalysisOptions) {
		options.describes = value
	}
}
