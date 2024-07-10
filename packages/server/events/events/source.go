package events

import (
	comutils "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	grpccommon "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/common"
)

type Source struct {
	Source        string `json:"source"`
	SourceOfTruth string `json:"sourceOfTruth"`
	AppSource     string `json:"appSource"`
}

func (s *Source) Available() bool {
	return s.Source != "" || s.SourceOfTruth != "" || s.AppSource != ""
}

func (s *Source) String() string {
	output, _ := comutils.ToJson(s)
	return output
}

func (s *Source) FromGrpc(grpcSource *grpccommon.SourceFields) {
	if grpcSource == nil {
		s.SetDefaultValues()
		return
	}
	s.Source = grpcSource.Source
	s.SourceOfTruth = comutils.StringFirstNonEmpty(grpcSource.SourceOfTruth, grpcSource.Source)
	s.AppSource = grpcSource.AppSource
}

func SourceFromGrpc(grpcSource *grpccommon.SourceFields) Source {
	s := Source{}
	s.Source = grpcSource.Source
	s.SourceOfTruth = comutils.StringFirstNonEmpty(grpcSource.SourceOfTruth, grpcSource.Source)
	s.AppSource = grpcSource.AppSource
	s.SetDefaultValues()
	return s
}

func (s *Source) SetDefaultValues() {
	if s.Source == "" {
		s.Source = "openline" //todo constants
	}
	if s.SourceOfTruth == "" {
		s.SourceOfTruth = s.Source
	}
	if s.AppSource == "" {
		s.AppSource = "event-processing-platform" //todo constants
	}
}
