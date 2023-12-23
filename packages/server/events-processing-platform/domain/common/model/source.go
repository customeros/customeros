package model

import (
	comutils "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/constants"
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
		return
	}
	s.Source = grpcSource.Source
	s.SourceOfTruth = comutils.StringFirstNonEmpty(grpcSource.SourceOfTruth, grpcSource.Source)
	s.AppSource = grpcSource.AppSource
}

func (s *Source) SetDefaultValues() {
	if s.Source == "" {
		s.Source = constants.SourceOpenline
	}
	if s.SourceOfTruth == "" {
		s.SourceOfTruth = s.Source
	}
	if s.AppSource == "" {
		s.AppSource = constants.AppSourceEventProcessingPlatform
	}
}
