package interfaces

import (
	"context"
	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	"net/url"
	"time"
)

type QuickbooksService interface {
	GetAndStoreAccessToken(ctx context.Context, realmId string, requestData url.Values) (*postgres_entity.QuickbooksSettingsEntity, error)
	SaveCustomer(ctx context.Context, id, customerName string) (*QuickbooksSaveCustomerResponse, error)
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
		Taxable         bool `json:"Taxable"`
		Job             bool `json:"Job"`
		BillWithParent  bool `json:"BillWithParent"`
		Balance         int  `json:"Balance"`
		BalanceWithJobs int  `json:"BalanceWithJobs"`
		CurrencyRef     struct {
			Value string `json:"value"`
			Name  string `json:"name"`
		} `json:"CurrencyRef"`
		PreferredDeliveryMethod string `json:"PreferredDeliveryMethod"`
		IsProject               bool   `json:"IsProject"`
		Domain                  string `json:"domain"`
		Sparse                  bool   `json:"sparse"`
		Id                      string `json:"Id"`
		SyncToken               string `json:"SyncToken"`
		MetaData                struct {
			CreateTime      time.Time `json:"CreateTime"`
			LastUpdatedTime time.Time `json:"LastUpdatedTime"`
		} `json:"MetaData"`
		GivenName          string `json:"GivenName"`
		FullyQualifiedName string `json:"FullyQualifiedName"`
		DisplayName        string `json:"DisplayName"`
		PrintOnCheckName   string `json:"PrintOnCheckName"`
		Active             bool   `json:"Active"`
		DefaultTaxCodeRef  struct {
			Value string `json:"value"`
		} `json:"DefaultTaxCodeRef"`
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
