package entity

import neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"

// Deprecated
type SourceFields struct {
	Source        neo4jentity.DataSource `json:"source"`
	SourceOfTruth neo4jentity.DataSource `json:"sourceOfTruth"`
	AppSource     string                 `json:"appSource"`
}
