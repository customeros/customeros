package models

import (
	comutils "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	grpccommon "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/common"
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
	s.SourceOfTruth = grpcSource.SourceOfTruth
	s.AppSource = grpcSource.AppSource
}
