package event

import "github.com/openline-ai/openline-customer-os/packages/server/events/event/common"

// Deprecated
type BaseRequest struct {
	ObjectID       string        `json:"objectID" validate:"required"`
	Tenant         string        `json:"tenant" validate:"required"`
	LoggedInUserId string        `json:"loggedInUserId"`
	AppSource      string        `json:"appSource"`
	SourceFields   common.Source `json:"sourceFields"`
}

func NewBaseRequest(objectID, tenant, loggedInUserId string, sourceFields common.Source) BaseRequest {
	return BaseRequest{
		ObjectID:       objectID,
		Tenant:         tenant,
		LoggedInUserId: loggedInUserId,
		SourceFields:   sourceFields,
	}
}
