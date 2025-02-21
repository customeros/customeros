package cosapi_interfaces

import (
	"context"
	neo4jentity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"
	neo4jenum "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/enum"
	"time"
)

type ServiceLineItemService interface {
	Create(ctx context.Context, serviceLineItemDetails ServiceLineItemCreateData) (string, error)
	Update(ctx context.Context, serviceLineItemDetails ServiceLineItemUpdateData) error
	Delete(ctx context.Context, serviceLineItemId string) (bool, error)
	Close(ctx context.Context, serviceLineItemId string, endedAt *time.Time) error
	NewVersion(ctx context.Context, data ServiceLineItemNewVersionData) (string, error)
}

type ServiceLineItemCreateData struct {
	ContractId        string                            `json:"contractId"`
	SliDescription    *string                           `json:"sliDescription"`
	SkuId             string                            `json:"skuId"`
	SliPrice          float64                           `json:"sliPrice"`
	SliQuantity       int64                             `json:"sliQuantity"`
	SliBilledType     neo4jenum.BilledType              `json:"sliBilledType"`
	ExternalReference *neo4jentity.ExternalSystemEntity `json:"externalReference"`
	Source            neo4jentity.DataSource            `json:"source"`
	AppSource         string                            `json:"appSource"`
	StartedAt         *time.Time                        `json:"startedAt"`
	EndedAt           *time.Time                        `json:"endedAt"`
	SliVatRate        float64                           `json:"sliVatRate"`
}

type ServiceLineItemNewVersionData struct {
	Id          string                 `json:"id"`
	SkuId       string                 `json:"skuId"`
	Description *string                `json:"sliDescription"`
	Price       float64                `json:"sliPrice"`
	Quantity    int64                  `json:"sliQuantity"`
	Comments    string                 `json:"sliComments"`
	Source      neo4jentity.DataSource `json:"source"`
	AppSource   string                 `json:"appSource"`
	VatRate     float64                `json:"sliVatRate"`
	StartedAt   *time.Time             `json:"startedAt"`
}

type ServiceLineItemUpdateData struct {
	Id                      string                 `json:"id"`
	IsRetroactiveCorrection bool                   `json:"isRetroactiveCorrection"`
	SkuId                   string                 `json:"skuId"`
	SliDescription          *string                `json:"sliDescription"`
	SliPrice                float64                `json:"sliPrice"`
	SliQuantity             int64                  `json:"sliQuantity"`
	SliBilledType           neo4jenum.BilledType   `json:"sliBilledType"`
	SliComments             string                 `json:"sliComments"`
	Source                  neo4jentity.DataSource `json:"source"`
	AppSource               string                 `json:"appSource"`
	SliVatRate              float64                `json:"sliVatRate"`
	StartedAt               *time.Time             `json:"startedAt"`
}
