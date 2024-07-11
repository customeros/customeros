package events

// Deprecated
type BaseRequest struct {
	ObjectID       string `json:"objectID" validate:"required"`
	Tenant         string `json:"tenant" validate:"required"`
	LoggedInUserId string `json:"loggedInUserId"`
	AppSource      string `json:"appSource"`
	SourceFields   Source `json:"sourceFields"`
}

func NewBaseRequest(objectID, tenant, loggedInUserId string, sourceFields Source) BaseRequest {
	return BaseRequest{
		ObjectID:       objectID,
		Tenant:         tenant,
		LoggedInUserId: loggedInUserId,
		SourceFields:   sourceFields,
	}
}
