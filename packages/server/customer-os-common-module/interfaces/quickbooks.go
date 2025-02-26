package interfaces

import (
	"context"
	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	"net/url"
	"time"
)

type QuickbooksService interface {
	QuickbooksConnected(ctx context.Context) (bool, error)
	GetAndStoreAccessToken(ctx context.Context, realmId string, requestData url.Values) (*postgres_entity.QuickbooksSettingsEntity, error)
	RevokeAccess(ctx context.Context) error
	SaveProduct(ctx context.Context, id, productName string, archived bool, price float64) (*QuickbooksSaveProductResponse, error)
	SaveCustomer(ctx context.Context, id, customerName string) (*QuickbooksSaveCustomerResponse, error)
	SaveInvoice(ctx context.Context, customerId, invoiceNumber string, invoiceDate time.Time, lines []QuickbooksInvoiceLine) (*QuickbooksSaveInvoiceResponse, error)
	PayInvoice(ctx context.Context, customerId, invoiceId string, totalAmount float64) (*QuickbooksSavePaymentResponse, error)
	VoidInvoice(ctx context.Context, invoiceId string) (*QuickbooksSaveInvoiceResponse, error)
}

type QuickbooksInvoiceLine struct {
	DetailType          string  `json:"DetailType"`
	Amount              float64 `json:"Amount"`
	SalesItemLineDetail struct {
		ItemRef struct {
			Value string `json:"value"`
		} `json:"ItemRef"`
		ServiceDate string `json:"ServiceDate"`
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

type QuickbooksSaveInvoiceResponse struct {
	QuickbooksCheckFaultResponse
	Invoice *struct {
		Id string `json:"Id"`
	} `json:"Invoice"`
}

type QuickbooksSavePaymentResponse struct {
	QuickbooksCheckFaultResponse
	Payment *struct {
		Id string `json:"Id"`
	} `json:"Payment"`
}

type QuickbooksGetInvoiceResponse struct {
	QuickbooksCheckFaultResponse
	Invoice *struct {
		Id        string `json:"Id"`
		SyncToken string `json:"SyncToken"`
	} `json:"Invoice"`
}

type QuickbooksGetProductResponse struct {
	Item struct {
		Type             string  `json:"Type"`
		Name             string  `json:"Name"`
		Active           bool    `json:"Active"`
		UnitPrice        float64 `json:"UnitPrice"`
		Sku              *string `json:"Sku"`
		IncomeAccountRef struct {
			Value string `json:"value"`
			Name  string `json:"name"`
		} `json:"IncomeAccountRef"`
		Sparse    bool   `json:"sparse"`
		Id        string `json:"Id"`
		SyncToken string `json:"SyncToken"`
		MetaData  struct {
			CreateTime      time.Time `json:"CreateTime"`
			LastUpdatedTime time.Time `json:"LastUpdatedTime"`
		} `json:"MetaData"`

		Description *string `json:"Description"`
		SubItem     *bool   `json:"SubItem"`
		ParentRef   *struct {
			Value string `json:"value"`
			Name  string `json:"name"`
		} `json:"ParentRef"`
		Level              *int     `json:"Level"`
		FullyQualifiedName *string  `json:"FullyQualifiedName"`
		Taxable            *bool    `json:"Taxable"`
		PurchaseDesc       *string  `json:"PurchaseDesc"`
		PurchaseCost       *float64 `json:"PurchaseCost"`
		ExpenseAccountRef  *struct {
			Value string `json:"value"`
			Name  string `json:"name"`
		} `json:"ExpenseAccountRef"`
		PrefVendorRef *struct {
			Value string `json:"value"`
			Name  string `json:"name"`
		} `json:"PrefVendorRef"`
		TrackQtyOnHand       *bool `json:"TrackQtyOnHand"`
		TaxClassificationRef *struct {
			Value string `json:"value"`
			Name  string `json:"name"`
		} `json:"TaxClassificationRef"`
		DeferredRevenue *bool   `json:"DeferredRevenue"`
		Domain          *string `json:"domain"`
	} `json:"Item"`
	Time time.Time `json:"time"`
}

type QuickbooksSearchAccountResponse struct {
	QuickbooksCheckFaultResponse
	QueryResponse struct {
		Account []struct {
			Sparse bool   `json:"sparse"`
			Id     string `json:"Id"`
		} `json:"Account"`
		StartPosition int `json:"startPosition"`
		MaxResults    int `json:"maxResults"`
	} `json:"QueryResponse"`
	Time time.Time `json:"time"`
}

type QuickbooksSaveAccountResponse struct {
	QuickbooksCheckFaultResponse
	Account *struct {
		Id string `json:"Id"`
	} `json:"Account"`
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
