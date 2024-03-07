package eventstore

import (
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
)

// Deprecated
type BaseRequest struct {
	ObjectID       string             `json:"objectID" validate:"required"`
	Tenant         string             `json:"tenant" validate:"required"`
	LoggedInUserId string             `json:"loggedInUserId"`
	AppSource      string             `json:"appSource"`
	SourceFields   commonmodel.Source `json:"sourceFields"`
}

func NewBaseRequest(objectID, tenant, loggedInUserId string, sourceFields commonmodel.Source) BaseRequest {
	return BaseRequest{
		ObjectID:       objectID,
		Tenant:         tenant,
		LoggedInUserId: loggedInUserId,
		SourceFields:   sourceFields,
	}
}
