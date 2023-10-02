package models

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"time"
)
import grpccommon "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/common"

type ExternalSystem struct {
	ExternalSystemId string     `json:"externalSystemId"`
	ExternalUrl      string     `json:"externalUrl"`
	ExternalId       string     `json:"externalId"`
	ExternalIdSecond string     `json:"externalIdSecond"`
	ExternalSource   string     `json:"externalSource"`
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
	e.SyncDate = utils.TimestampProtoToTime(grpcExternalSystem.SyncDate)
}

func (e *ExternalSystem) Available() bool {
	return e.ExternalSystemId != "" && e.ExternalId != ""
}
