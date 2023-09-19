package models

import (
	comutils "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/utils"
	"time"
)
import grpccommon "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/common"

type ExternalSystem struct {
	ExternalSystemId string     `json:"externalSystemId"`
	ExternalUrl      string     `json:"externalUrl"`
	ExternalId       string     `json:"externalId"`
	ExternalSource   string     `json:"externalSource"`
	SyncDate         *time.Time `json:"syncDate,omitempty"`
}

func (e *ExternalSystem) String() string {
	output, _ := comutils.ToJson(e)
	return output
}

func (e *ExternalSystem) FromGrpc(grpcExternalSystem *grpccommon.ExternalSystemFields) {
	if grpcExternalSystem == nil {
		return
	}
	e.ExternalSystemId = grpcExternalSystem.ExternalSystemId
	e.ExternalUrl = grpcExternalSystem.ExternalUrl
	e.ExternalId = grpcExternalSystem.ExternalId
	e.ExternalSource = grpcExternalSystem.ExternalSource
	e.SyncDate = utils.TimestampProtoToTime(grpcExternalSystem.SyncDate)
}

func (e *ExternalSystem) Available() bool {
	return e.ExternalSystemId != "" && e.ExternalId != ""
}
