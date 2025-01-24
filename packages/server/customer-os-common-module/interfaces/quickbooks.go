package interfaces

import (
	"context"
	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	"net/url"
	"time"
)

type QuickbooksService interface {
	GetAndStoreAccessToken(ctx context.Context, realmId string, requestData url.Values) (*postgres_entity.QuickbooksSettingsEntity, error)
	SaveProduct(ctx context.Context, id, productName string) (*QuickbooksSaveProductResponse, error)
	SaveCustomer(ctx context.Context, id, customerName string) (*QuickbooksSaveCustomerResponse, error)
}

type QuickbooksInvoiceLine struct {
	DetailType          string  `json:"DetailType"`
	Amount              float64 `json:"Amount"`
	SalesItemLineDetail struct {
		ItemRef struct {
			Value string `json:"value"`
		} `json:"ItemRef"`
	} `json:"SalesItemLineDetail"`
}

type OauthQuickbooksResponse struct {
	ExpiresIn              int     `json:"expires_in"`
	TokenType              string  `json:"token_type"`
	XRefreshTokenExpiresIn int     `json:"x_refresh_token_expires_in"`
	RefreshToken           string  `json:"refresh_token"`
	AccessToken            string  `json:"access_token"`
	Error                  *string `json:"error"`
}

type QuickbooksSaveCustomerResponse struct {
	Customer *struct {
		Id string `json:"Id"`
	} `json:"Customer"`
	Fault *struct {
		Error []struct {
			Message string `json:"Message"`
			Detail  string `json:"Detail"`
			Code    string `json:"code"`
		} `json:"Error"`
		Type string `json:"type"`
	} `json:"Fault"`
	Time time.Time `json:"time"`
}

type QuickbooksSaveProductResponse struct {
	Customer *struct {
		Id string `json:"Id"`
	} `json:"Customer"`
	Fault *struct {
		Error []struct {
			Message string `json:"Message"`
			Detail  string `json:"Detail"`
			Code    string `json:"code"`
		} `json:"Error"`
		Type string `json:"type"`
	} `json:"Fault"`
	Time time.Time `json:"time"`
}
