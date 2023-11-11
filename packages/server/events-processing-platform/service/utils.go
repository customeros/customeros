package service

import (
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"time"
)

func convertCreateAndUpdateProtoTimestampsToTime(createdAtProto, updatedAtProto *timestamp.Timestamp) (*time.Time, *time.Time) {
	createdAt := utils.TimestampProtoToTimePtr(createdAtProto)
	updatedAt := utils.TimestampProtoToTimePtr(updatedAtProto)
	return createdAt, updatedAt
}
