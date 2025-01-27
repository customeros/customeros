package interfaces

import (
	"context"
	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	"net/url"
)

type QuickbooksService interface {
	GetAndStoreAccessToken(ctx context.Context, realmId string, requestData url.Values) (*postgres_entity.QuickbooksSettingsEntity, error)
	SaveProduct(ctx context.Context, id, productName string, archived bool) (*QuickbooksSaveProductResponse, error)
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
	QuickbooksCheckFaultResponse
	Customer *struct {
		Id string `json:"Id"`
	} `json:"Customer"`
}

type QuickbooksGetProductResponse struct {
	QuickbooksCheckFaultResponse
	Product *struct {
		Id        string `json:"Id"`
		SyncToken string `json:"SyncToken"`
	} `json:"Item"`
}

type QuickbooksSaveProductResponse struct {
	QuickbooksCheckFaultResponse
	Product *struct {
		Id string `json:"Id"`
	} `json:"Item"`
}

type QuickbooksCheckFaultResponse struct {
	Fault *struct {
		Error []struct {
			Message string      `json:"message"`
			Detail  string      `json:"detail"`
			Code    string      `json:"code"`
			Element interface{} `json:"element"`
		} `json:"error"`
		Type string `json:"type"`
	} `json:"fault"`
}
