package model

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"time"
)
import grpccommon "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/common"

type ExternalSystem struct {
	ExternalSystemId string     `json:"externalSystemId,omitempty"`
	ExternalUrl      string     `json:"externalUrl,omitempty"`
	ExternalId       string     `json:"externalId,omitempty"`
	ExternalIdSecond string     `json:"externalIdSecond,omitempty"`
	ExternalSource   string     `json:"externalSource,omitempty"`
	SyncDate         *time.Time `json:"syncDate,omitempty"`
}

func (e *ExternalSystem) String() string {
	output, _ := utils.ToJson(e)
	return output
}

func (e *ExternalSystem) FromGrpc(grpcExternalSystem *grpccommon.ExternalSystemFields) {
	if grpcExternalSystem == nil {
		return
	}
	e.ExternalSystemId = grpcExternalSystem.ExternalSystemId
	e.ExternalUrl = grpcExternalSystem.ExternalUrl
	e.ExternalId = grpcExternalSystem.ExternalId
	e.ExternalIdSecond = grpcExternalSystem.ExternalIdSecond
	e.ExternalSource = grpcExternalSystem.ExternalSource
	e.SyncDate = utils.TimestampProtoToTimePtr(grpcExternalSystem.SyncDate)
}

func (e *ExternalSystem) Available() bool {
	return e.ExternalSystemId != "" && e.ExternalId != ""
}
