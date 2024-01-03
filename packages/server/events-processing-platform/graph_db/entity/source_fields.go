package entity

import neo4jentity "github.com/openline-ai/customer-os-neo4j-repository/entity"

type SourceFields struct {
	Source        neo4jentity.DataSource `json:"source"`
	SourceOfTruth neo4jentity.DataSource `json:"sourceOfTruth"`
	AppSource     string                 `json:"appSource"`
}
