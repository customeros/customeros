package data_fields

import (
	neo4jenum "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/enum"
	"time"
)

type SLIFields struct {
	AppSource   *string               `json:"appSource,omitempty"`
	Source      *string               `json:"source,omitempty"`
	CreatedAt   *time.Time            `json:"createdAt,omitempty"`
	BilledType  *neo4jenum.BilledType `json:"billedType,omitempty"`
	Quantity    *int64                `json:"quantity,omitempty"`
	Price       *float64              `json:"price,omitempty"`
	SkuId       *string               `json:"skuId,omitempty"`
	Description *string               `json:"description,omitempty"`
	ContractId  *string               `json:"contractId,omitempty"`
	StartedAt   *time.Time            `json:"startedAt,omitempty"`
	EndedAt     *time.Time            `json:"endedAt,omitempty"`
	TaxRate     *float64              `json:"taxRate,omitempty"`
	Comments    *string               `json:"comments,omitempty"`
	ParentId    *string               `json:"parentId,omitempty"`
	NewVersion  *bool                 `json:"newVersion,omitempty"`
}

func (s *SLIFields) IsNewVersion() bool {
	return s.NewVersion != nil && *s.NewVersion
}
