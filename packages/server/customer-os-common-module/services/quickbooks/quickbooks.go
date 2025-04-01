package quickbooks

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/config"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	postgres_repository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
)

type quickbooksService struct {
	postgres *postgres_repository.Repositories
	qbConfig *config.QuickbooksConfig
}

func NewQuickbooksService(qbConfig *config.QuickbooksConfig, postgres *postgres_repository.Repositories) interfaces.QuickbooksService {
	return &quickbooksService{
		qbConfig: qbConfig,
		postgres: postgres,
	}
}

func (s *quickbooksService) GetAndStoreAccessToken(ctx context.Context, realmId string, requestData url.Values) (*postgres_entity.QuickbooksSettingsEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "QuickbooksService.getAuthToken")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	tenant := common.GetTenantFromContext(ctx)

	// Encode the form data
	requestBody := requestData.Encode()

	request, err := http.NewRequest("POST", "https://oauth.platform.intuit.com/oauth2/v1/tokens/bearer", nil)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	request.Body = ioutil.NopCloser(strings.NewReader(requestBody))

	toString := base64.StdEncoding.EncodeToString([]byte(s.qbConfig.ClientId + ":" + s.qbConfig.ClientSecret))
	request.Header.Set("Authorization", "Basic "+toString)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Perform the HTTP request
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	defer resp.Body.Close()

	// Read and print the response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	// convert body to OauthSlackResponse
	var quickbooksResponse interfaces.OauthQuickbooksResponse
	err = json.Unmarshal(body, &quickbooksResponse)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	if quickbooksResponse.Error == nil {
		now := utils.Now()

		qbSettings, err := s.postgres.QuickbooksSettingsRepository.Get(ctx, tenant)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}

		entity := postgres_entity.QuickbooksSettingsEntity{}
		if qbSettings != nil {
			entity = *qbSettings
		} else {
			entity.Tenant = tenant
			entity.RealmId = realmId
		}

		entity.AccessToken = quickbooksResponse.AccessToken
		entity.AccessTokenExpiresIn = quickbooksResponse.ExpiresIn
		entity.AccessTokenExpiresAt = now.Add(time.Duration(quickbooksResponse.ExpiresIn) * time.Second)
		entity.RefreshToken = quickbooksResponse.RefreshToken
		entity.RefreshTokenExpiresIn = quickbooksResponse.XRefreshTokenExpiresIn
		entity.RefreshTokenExpiresAt = now.Add(time.Duration(quickbooksResponse.XRefreshTokenExpiresIn) * time.Second)
		entity.RefreshTokenExpired = false

		stored, err := s.postgres.QuickbooksSettingsRepository.Save(ctx, entity)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}

		return stored, nil
	} else {
		span.LogFields(log.Object("error", *quickbooksResponse.Error))

		tracing.TraceErr(span, fmt.Errorf("error: %s", *quickbooksResponse.Error))
		return nil, fmt.Errorf("error: %s", *quickbooksResponse.Error)
	}
}

func (s *quickbooksService) RevokeAccess(ctx context.Context) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "QuickbooksService.RevokeAccess")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	tenant := common.GetTenantFromContext(ctx)

	// Retrieve QuickBooks settings for the tenant
	qbSettings, err := s.postgres.QuickbooksSettingsRepository.Get(ctx, tenant)
	if err != nil {
		tracing.TraceErr(span, err)
		return fmt.Errorf("failed to retrieve QuickBooks settings for tenant %s: %w", tenant, err)
	}
	if qbSettings == nil {
		err = errors.New("QuickBooks settings not found")
		tracing.TraceErr(span, err)
		return err
	}

	// Ensure that an access token exists
	token := qbSettings.AccessToken
	if token == "" {
		err := fmt.Errorf("no access token available for tenant %s", tenant)
		tracing.TraceErr(span, err)
		return err
	}

	// Prepare the revoke request
	revokeURL := "https://developer.api.intuit.com/v2/oauth2/tokens/revoke"
	formData := url.Values{}
	formData.Set("token", token)
	requestBody := formData.Encode()

	req, err := http.NewRequestWithContext(ctx, "POST", revokeURL, strings.NewReader(requestBody))
	if err != nil {
		tracing.TraceErr(span, err)
		return fmt.Errorf("failed to create revoke request: %w", err)
	}

	// Build the basic authorization header using client credentials
	credentials := s.qbConfig.ClientId + ":" + s.qbConfig.ClientSecret
	encodedCreds := base64.StdEncoding.EncodeToString([]byte(credentials))
	req.Header.Set("Authorization", "Basic "+encodedCreds)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Execute the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		tracing.TraceErr(span, err)
		return fmt.Errorf("failed to send revoke request: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		tracing.TraceErr(span, err)
		return fmt.Errorf("failed to read revoke response: %w", err)
	}

	span.LogFields(log.String("quickbooks.response", string(bodyBytes)))

	// Check for a successful response
	if resp.StatusCode != http.StatusOK {
		errMsg := fmt.Sprintf("revoke quickbooks request returned status %d", resp.StatusCode)
		tracing.TraceErr(span, fmt.Errorf(errMsg))
	}

	// http.StatusBadRequest is returned when revoking access from already removed app
	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusBadRequest {
		err = s.postgres.QuickbooksSettingsRepository.Delete(ctx, tenant)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}
	}

	return nil
}

func (s *quickbooksService) SaveProduct(ctx context.Context, id, productName string, archived bool, price float64) (*interfaces.QuickbooksSaveProductResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "QuickbooksService.SaveProduct")
	defer span.Finish()
	tracing.TagComponentService(span)
	tracing.TagTenant(span, common.GetTenantFromContext(ctx))
	span.LogFields(log.String("productName", productName), log.Bool("archived", archived), log.Float64("price", price), log.String("id", id))

	tenant := common.GetTenantFromContext(ctx)

	qbSettings, err := s.postgres.QuickbooksSettingsRepository.Get(ctx, tenant)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	if qbSettings == nil {
		err = errors.New("Quickbooks settings not found")
		tracing.TraceErr(span, err)
		return nil, err
	}

	if qbSettings.SalesAccountId == "" {
		salesAccountId, err := s.GetAccountIdByName(ctx, "CustomerOS Sales")
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}

		if salesAccountId != "" {
			qbSettings.SalesAccountId = salesAccountId
			_, err = s.postgres.QuickbooksSettingsRepository.Save(ctx, *qbSettings)
			if err != nil {
				tracing.TraceErr(span, err)
				return nil, err
			}
		} else {
			salesAccountUrl := fmt.Sprintf(s.qbConfig.Url+"/v3/company/%s/account", qbSettings.RealmId)
			salesAccountRequest := map[string]interface{}{
				"Name":        "CustomerOS Sales",
				"AccountType": "Income",
			}

			qbAccountResponse, err := s.performRequest(ctx, qbSettings, salesAccountUrl, "POST", salesAccountRequest, true)
			if err != nil {
				tracing.TraceErr(span, err)
				return nil, err
			}

			var qbAccount interfaces.QuickbooksSaveAccountResponse
			err = json.Unmarshal(qbAccountResponse, &qbAccount)
			if err != nil {
				tracing.TraceErr(span, err)
				return nil, err
			}

			qbSettings.SalesAccountId = qbAccount.Account.Id
			_, err = s.postgres.QuickbooksSettingsRepository.Save(ctx, *qbSettings)
			if err != nil {
				tracing.TraceErr(span, err)
				return nil, err
			}
		}
	}

	var qbSaveProductRequest map[string]interface{}
	var qbProduct interfaces.QuickbooksGetProductResponse

	if id != "" {
		productByIdUrl := fmt.Sprintf(s.qbConfig.Url+"/v3/company/%s/item/%s", qbSettings.RealmId, id)
		qbProductResponse, err := s.performRequest(ctx, qbSettings, productByIdUrl, "GET", nil, true)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}

		err = json.Unmarshal(qbProductResponse, &qbProduct)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}

		qbProduct.Item.Type = "Service"
		qbProduct.Item.Name = productName
		qbProduct.Item.UnitPrice = price
		qbProduct.Item.Active = !archived

		// nullify PurchaseCost if ExpenseAccountRef is not set
		if qbProduct.Item.ExpenseAccountRef == nil {
			qbProduct.Item.PurchaseCost = nil
		}

		jsonBytes, _ := json.Marshal(qbProduct.Item)
		err = json.Unmarshal(jsonBytes, &qbSaveProductRequest)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}
	} else {
		qbSaveProductRequest = map[string]interface{}{
			"Name":      productName,
			"Active":    !archived,
			"UnitPrice": price,
			"Type":      "Service",
			"IncomeAccountRef": map[string]interface{}{
				"value": qbSettings.SalesAccountId,
			},
		}
	}

	requestUrl := fmt.Sprintf(s.qbConfig.Url+"/v3/company/%s/item", qbSettings.RealmId)
	qbResponse, err := s.performRequest(ctx, qbSettings, requestUrl, "POST", qbSaveProductRequest, true)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	var qbProductResponse interfaces.QuickbooksSaveProductResponse
	err = json.Unmarshal(qbResponse, &qbProductResponse)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return &qbProductResponse, nil
}

func (s *quickbooksService) GetProduct(ctx context.Context, id string) (*interfaces.QuickbooksGetProductResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "QuickbooksService.GetProduct")
	defer span.Finish()
	tracing.TagComponentService(span)
	tracing.TagTenant(span, common.GetTenantFromContext(ctx))
	span.LogFields(log.String("id", id))

	tenant := common.GetTenantFromContext(ctx)

	qbSettings, err := s.postgres.QuickbooksSettingsRepository.Get(ctx, tenant)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	if qbSettings == nil {
		err = errors.New("Quickbooks settings not found")
		tracing.TraceErr(span, err)
		return nil, err
	}

	requestUrl := fmt.Sprintf(s.qbConfig.Url+"/v3/company/%s/item/%s", qbSettings.RealmId, id)
	resp, err := s.performRequest(ctx, qbSettings, requestUrl, "GET", nil, true)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	var quickbooksResponse interfaces.QuickbooksGetProductResponse
	err = json.Unmarshal(resp, &quickbooksResponse)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return &quickbooksResponse, nil
}

func (s *quickbooksService) SaveCustomer(ctx context.Context, id, customerName string) (*interfaces.QuickbooksSaveCustomerResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "QuickbooksService.SaveCustomer")
	defer span.Finish()

	tenant := common.GetTenantFromContext(ctx)

	qbSettings, err := s.postgres.QuickbooksSettingsRepository.Get(ctx, tenant)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	if qbSettings == nil {
		err = errors.New("Quickbooks settings not found")
		tracing.TraceErr(span, err)
		return nil, err
	}

	request := map[string]interface{}{
		"GivenName": customerName,
	}

	if id != "" {
		request["Id"] = id
	}

	resp, err := s.performRequest(ctx, qbSettings, s.qbConfig.Url+"/v3/company/"+qbSettings.RealmId+"/customer", "POST", request, true)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	// convert body to OauthSlackResponse
	var quickbooksResponse interfaces.QuickbooksSaveCustomerResponse
	err = json.Unmarshal(resp, &quickbooksResponse)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	if quickbooksResponse.Fault != nil {
		span.LogFields(log.Object("error", quickbooksResponse.Fault))
		return nil, fmt.Errorf("error: %s", quickbooksResponse.Fault.Error[0].Message)
	}

	return &quickbooksResponse, nil
}

func (s *quickbooksService) SaveInvoice(ctx context.Context, quickbooksCustomerId, invoiceNumber string, invoiceDate, dueDate time.Time, invoiceEmail string,
	lines []interfaces.QuickbooksInvoiceLine) (*interfaces.QuickbooksSaveInvoiceResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "QuickbooksService.SaveInvoice")
	defer span.Finish()
	tenant := common.GetTenantFromContext(ctx)
	tracing.TagTenant(span, tenant)
	span.LogFields(
		log.String("quickbooksCustomerId", quickbooksCustomerId),
		log.String("invoiceNumber", invoiceNumber),
		log.Object("invoiceDate", invoiceDate),
		log.Object("dueDate", dueDate),
		log.String("invoiceEmail", invoiceEmail))

	qbSettings, err := s.postgres.QuickbooksSettingsRepository.Get(ctx, tenant)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	if qbSettings == nil {
		err = errors.New("Quickbooks settings not found")
		tracing.TraceErr(span, err)
		return nil, err
	}

	request := map[string]interface{}{
		"CustomerRef": map[string]interface{}{
			"value": quickbooksCustomerId,
		},
		"TxnDate":   invoiceDate.Format("2006/01/02"),
		"DueDate":   dueDate.Format("2006/01/02"),
		"Line":      lines,
		"DocNumber": invoiceNumber,
		"BillEmail": map[string]interface{}{
			"Address": invoiceEmail,
		},
	}

	resp, err := s.performRequest(ctx, qbSettings, s.qbConfig.Url+"/v3/company/"+qbSettings.RealmId+"/invoice", "POST", request, true)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	// convert body to OauthSlackResponse
	var quickbooksResponse interfaces.QuickbooksSaveInvoiceResponse
	err = json.Unmarshal(resp, &quickbooksResponse)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	if quickbooksResponse.Fault != nil {
		span.LogFields(log.Object("error", quickbooksResponse.Fault))
		return nil, fmt.Errorf("error: %s", quickbooksResponse.Fault.Error[0].Message)
	}

	return &quickbooksResponse, nil
}

func (s *quickbooksService) VoidInvoice(ctx context.Context, invoiceId string) (*interfaces.QuickbooksSaveInvoiceResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "QuickbooksService.VoidInvoice")
	defer span.Finish()

	tenant := common.GetTenantFromContext(ctx)

	qbSettings, err := s.postgres.QuickbooksSettingsRepository.Get(ctx, tenant)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	if qbSettings == nil {
		err = errors.New("Quickbooks settings not found")
		tracing.TraceErr(span, err)
		return nil, err
	}

	request := map[string]interface{}{
		"Id": invoiceId,
	}

	if invoiceId != "" {
		invoiceByIdUrl := fmt.Sprintf(s.qbConfig.Url+"/v3/company/%s/invoice/%s", qbSettings.RealmId, invoiceId)
		qbInvoiceResponse, err := s.performRequest(ctx, qbSettings, invoiceByIdUrl, "GET", request, true)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}

		var qbInvoice interfaces.QuickbooksGetInvoiceResponse
		err = json.Unmarshal(qbInvoiceResponse, &qbInvoice)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}

		request["SyncToken"] = qbInvoice.Invoice.SyncToken
	}

	resp, err := s.performRequest(ctx, qbSettings, s.qbConfig.Url+"/v3/company/"+qbSettings.RealmId+"/invoice?operation=void", "POST", request, true)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	// convert body to OauthSlackResponse
	var quickbooksResponse interfaces.QuickbooksSaveInvoiceResponse
	err = json.Unmarshal(resp, &quickbooksResponse)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	if quickbooksResponse.Fault != nil {
		span.LogFields(log.Object("error", quickbooksResponse.Fault))
		return nil, fmt.Errorf("error: %s", quickbooksResponse.Fault.Error[0].Message)
	}

	return &quickbooksResponse, nil
}

func (s *quickbooksService) PayInvoice(ctx context.Context, customerId, invoiceId, invoiceNumber string, totalAmount float64, paymentIncomeAccountName string) (*interfaces.QuickbooksSavePaymentResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "QuickbooksService.PayInvoice")
	defer span.Finish()
	tracing.TagTenant(span, common.GetTenantFromContext(ctx))
	tracing.TagEntity(span, invoiceId)
	span.LogFields(log.String("paymentIncomeAccountName", paymentIncomeAccountName))
	span.LogFields(log.String("totalAmount", fmt.Sprintf("%f", totalAmount)))
	span.LogFields(log.String("invoiceNumber", invoiceNumber))
	tenant := common.GetTenantFromContext(ctx)

	qbSettings, err := s.postgres.QuickbooksSettingsRepository.Get(ctx, tenant)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	if qbSettings == nil {
		err = errors.New("Quickbooks settings not found")
		tracing.TraceErr(span, err)
		return nil, err
	}

	request := map[string]interface{}{
		"CustomerRef": map[string]interface{}{
			"value": customerId,
		},
		"TotalAmt":      totalAmount,
		"PaymentRefNum": invoiceNumber,
		"Line": []map[string]interface{}{
			{
				"Amount": totalAmount,
				"LinkedTxn": []map[string]interface{}{
					{
						"TxnId":   invoiceId,
						"TxnType": "Invoice",
					},
				},
			},
		},
	}

	// If paymentIncomeAccountName is provided, try to find the account ID
	if paymentIncomeAccountName != "" {
		accountId, err := s.GetAccountIdByName(ctx, paymentIncomeAccountName)
		if err != nil {
			tracing.TraceErr(span, err)
		} else if accountId != "" {
			request["DepositToAccountRef"] = map[string]interface{}{
				"value": accountId,
			}
		}
	}

	resp, err := s.performRequest(ctx, qbSettings, s.qbConfig.Url+"/v3/company/"+qbSettings.RealmId+"/payment", "POST", request, true)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	// convert body to OauthSlackResponse
	var quickbooksResponse interfaces.QuickbooksSavePaymentResponse
	err = json.Unmarshal(resp, &quickbooksResponse)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	if quickbooksResponse.Fault != nil {
		span.LogFields(log.Object("error", quickbooksResponse.Fault))
		return nil, fmt.Errorf("error: %s", quickbooksResponse.Fault.Error[0].Message)
	}

	return &quickbooksResponse, nil
}

func (s *quickbooksService) performRequest(ctx context.Context, qbSettings *postgres_entity.QuickbooksSettingsEntity, requestUrl string, requestMethod string, requestBody map[string]interface{}, rerunOnTokenRefresh bool) ([]byte, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "QuickbooksService.performRequest")
	defer span.Finish()
	span.LogFields(log.String("requestUrl", requestUrl))
	span.LogFields(log.String("requestMethod", requestMethod))
	tracing.LogObjectAsJson(span, "requestBody", requestBody)
	span.LogFields(log.Bool("rerunOnTokenRefresh", rerunOnTokenRefresh))

	payload, err := json.Marshal(requestBody)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	// Create a new HTTP request
	req, err := http.NewRequest(requestMethod, requestUrl, bytes.NewBuffer(payload))
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	// Set headers
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+qbSettings.AccessToken)

	// Perform the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	defer resp.Body.Close()

	var bodyBytes []byte

	// Read response body
	bodyBytes, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	span.LogFields(log.String("result.quickbooksBody", string(bodyBytes)))

	// convert body to OauthSlackResponse
	var quickbooksCheckFaultResponse interfaces.QuickbooksCheckFaultResponse
	err = json.Unmarshal(bodyBytes, &quickbooksCheckFaultResponse)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to unmarshal response"))
		return nil, err
	}

	fault := quickbooksCheckFaultResponse.Fault
	if fault != nil {
		if (*fault).Type == "AUTHENTICATION" {
			//refresh token
			requestData := url.Values{}
			requestData.Set("grant_type", "refresh_token")
			requestData.Set("refresh_token", qbSettings.RefreshToken)

			qbSettings, err = s.GetAndStoreAccessToken(ctx, qbSettings.RealmId, requestData)
			if err != nil {
				tracing.TraceErr(span, err)
				return nil, err
			}

			// re-run original request
			if rerunOnTokenRefresh {
				return s.performRequest(ctx, qbSettings, requestUrl, requestMethod, requestBody, false)
			}
		} else {
			err = fmt.Errorf("error: %s", fault.Error[0].Message)
			tracing.TraceErr(span, err)
			return nil, err
		}
	}

	return bodyBytes, nil
}

func (s *quickbooksService) QuickbooksConnected(ctx context.Context) (bool, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "QuickbooksService.QuickbooksConnected")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	tenant := common.GetTenantFromContext(ctx)

	qbSettings, err := s.postgres.QuickbooksSettingsRepository.Get(ctx, tenant)
	if err != nil {
		tracing.TraceErr(span, err)
		return false, err
	}

	return qbSettings != nil, nil
}

func (s *quickbooksService) SaveJournalEntry(ctx context.Context, txnDate time.Time, journalLineItems []interfaces.QuickbooksJournalEntryLine) (*interfaces.QuickbooksJournalEntryResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "QuickbooksService.SaveJournalEntry")
	defer span.Finish()
	tenant := common.GetTenantFromContext(ctx)
	tracing.TagTenant(span, tenant)
	tracing.LogObjectAsJson(span, "journalLineItems", journalLineItems)

	txnDateStr := txnDate.Format("2006/01/02")
	span.LogFields(log.String("txnDate", txnDateStr))

	// Retrieve QuickBooks settings for the tenant.
	qbSettingsEntity, err := s.postgres.QuickbooksSettingsRepository.Get(ctx, tenant)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	if qbSettingsEntity == nil {
		span.LogFields(log.String("error", "QuickBooks settings not found"))
		return nil, nil
	}

	// Build the request payload.
	request := map[string]interface{}{
		"TxnDate": txnDateStr,
		"Line":    journalLineItems,
	}

	// Construct the URL for creating a journal entry.
	requestUrl := fmt.Sprintf("%s/v3/company/%s/journalentry", s.qbConfig.Url, qbSettingsEntity.RealmId)

	// Perform the request.
	resp, err := s.performRequest(ctx, qbSettingsEntity, requestUrl, "POST", request, true)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	// Unmarshal the response into a QuickbooksJournalEntry.
	var qbJournalEntryResponse interfaces.QuickbooksJournalEntryResponse
	err = json.Unmarshal(resp, &qbJournalEntryResponse)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	if qbJournalEntryResponse.Fault != nil {
		span.LogFields(log.Object("error", qbJournalEntryResponse.Fault))
		return nil, fmt.Errorf("error: %s", qbJournalEntryResponse.Fault.Error[0].Message)
	}

	return &qbJournalEntryResponse, nil
}

func (s *quickbooksService) ZeroJournalEntry(ctx context.Context, journalEntryId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "QuickbooksService.ZeroJournalEntry")
	defer span.Finish()
	tenant := common.GetTenantFromContext(ctx)
	tracing.TagTenant(span, tenant)
	span.LogFields(log.String("journalEntryId", journalEntryId))

	// Retrieve QuickBooks settings for the tenant.
	qbSettings, err := s.postgres.QuickbooksSettingsRepository.Get(ctx, tenant)
	if err != nil {
		tracing.TraceErr(span, err)
		return fmt.Errorf("failed to retrieve QuickBooks settings for tenant %s: %w", tenant, err)
	}
	if qbSettings == nil {
		err = errors.New("Quickbooks settings not found")
		tracing.TraceErr(span, err)
		return err
	}

	// Construct URL to fetch the journal entry by ID.
	getURL := fmt.Sprintf("%s/v3/company/%s/journalentry/%s", s.qbConfig.Url, qbSettings.RealmId, journalEntryId)
	// Fetch the existing journal entry.
	resp, err := s.performRequest(ctx, qbSettings, getURL, "GET", nil, true)
	if err != nil {
		tracing.TraceErr(span, err)
		return fmt.Errorf("failed to fetch journal entry: %w", err)
	}

	var journalEntryResp interfaces.QuickbooksJournalEntryResponse
	err = json.Unmarshal(resp, &journalEntryResp)
	if err != nil {
		tracing.TraceErr(span, err)
		return fmt.Errorf("failed to unmarshal journal entry response: %w", err)
	}

	// Set all line item amounts to zero.
	for i := range journalEntryResp.JournalEntry.Line {
		journalEntryResp.JournalEntry.Line[i].Amount = 0
	}

	// Build the update payload.
	updatePayload := map[string]interface{}{
		"TxnDate": journalEntryResp.JournalEntry.TxnDate, // Preserve the original transaction date
		"Line":    journalEntryResp.JournalEntry.Line,    // Updated lines with zero amounts
	}

	// Construct URL for updating the journal entry (QuickBooks updates via POST to the same endpoint).
	updateURL := fmt.Sprintf("%s/v3/company/%s/journalentry", s.qbConfig.Url, qbSettings.RealmId)
	updateResp, err := s.performRequest(ctx, qbSettings, updateURL, "POST", updatePayload, true)
	if err != nil {
		tracing.TraceErr(span, err)
		return fmt.Errorf("failed to update journal entry: %w", err)
	}

	var updateJournalEntryResp interfaces.QuickbooksJournalEntryResponse
	err = json.Unmarshal(updateResp, &updateJournalEntryResp)
	if err != nil {
		tracing.TraceErr(span, err)
		return fmt.Errorf("failed to unmarshal updated journal entry response: %w", err)
	}

	if updateJournalEntryResp.Fault != nil {
		span.LogFields(log.Object("error", updateJournalEntryResp.Fault))
		return fmt.Errorf("error updating journal entry: %s", updateJournalEntryResp.Fault.Error[0].Message)
	}

	return nil
}

func (s *quickbooksService) GetAccountIdByName(ctx context.Context, accountName string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "QuickbooksService.GetAccountIdByName")
	defer span.Finish()
	tenant := common.GetTenantFromContext(ctx)
	tracing.TagTenant(span, tenant)
	span.LogFields(log.String("accountName", accountName))

	// Retrieve QuickBooks settings for the tenant.
	qbSettings, err := s.postgres.QuickbooksSettingsRepository.Get(ctx, tenant)
	if err != nil {
		tracing.TraceErr(span, err)
		return "", fmt.Errorf("failed to retrieve QuickBooks settings for tenant %s: %w", tenant, err)
	}
	if qbSettings == nil {
		err = errors.New("QuickBooks settings not found")
		tracing.TraceErr(span, err)
		return "", err
	}

	// Construct the URL for querying the account by name.
	queryURL := fmt.Sprintf("%s/v3/company/%s/query?query=select+Id+from+Account+where+Name='%s'", s.qbConfig.Url, qbSettings.RealmId, accountName)
	// Perform the request.
	resp, err := s.performRequest(ctx, qbSettings, queryURL, "GET", nil, true)
	if err != nil {
		tracing.TraceErr(span, err)
		return "", fmt.Errorf("failed to query account: %w", err)
	}

	var searchAccountResp interfaces.QuickbooksSearchAccountResponse
	err = json.Unmarshal(resp, &searchAccountResp)
	if err != nil {
		tracing.TraceErr(span, err)
		return "", fmt.Errorf("failed to unmarshal account response: %w", err)
	}

	if len(searchAccountResp.QueryResponse.Account) > 0 {
		return searchAccountResp.QueryResponse.Account[0].Id, nil
	}
	return "", nil
}

func (s *quickbooksService) SavePaymentLinkingJournalEntryToInvoice(ctx context.Context, quickbooksCustomerId string,
	quickbooksInvoiceId string, quickbooksJournalEntryId string, txnDate time.Time, totalAmount float64) (*interfaces.QuickbooksSavePaymentResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "QuickbooksService.SavePaymentLinkingJournalEntryToInvoice")
	defer span.Finish()
	tenant := common.GetTenantFromContext(ctx)
	tracing.TagTenant(span, tenant)
	span.LogFields(
		log.String("quickbooksCustomerId", quickbooksCustomerId),
		log.String("quickbooksInvoiceId", quickbooksInvoiceId),
		log.String("quickbooksJournalEntryId", quickbooksJournalEntryId),
		log.Float64("totalAmount", totalAmount),
		log.Object("txnDate", txnDate))

	// Retrieve QuickBooks settings for the tenant.
	qbSettings, err := s.postgres.QuickbooksSettingsRepository.Get(ctx, tenant)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, fmt.Errorf("failed to retrieve QuickBooks settings for tenant %s: %w", tenant, err)
	}
	if qbSettings == nil {
		err = errors.New("QuickBooks settings not found")
		tracing.TraceErr(span, err)
		return nil, err
	}

	// Construct the payment payload.
	paymentData := map[string]interface{}{
		"TxnDate":  txnDate.Format("2006/01/02"),
		"TotalAmt": totalAmount,
		"CustomerRef": map[string]interface{}{
			"value": quickbooksCustomerId,
		},
		"Line": []map[string]interface{}{
			{
				"Amount": totalAmount,
				"LinkedTxn": []map[string]interface{}{
					{
						"TxnId":   quickbooksJournalEntryId,
						"TxnType": "JournalEntry",
					},
				},
			},
		},
		"LinkedTxn": []map[string]interface{}{
			{
				"TxnId":   quickbooksInvoiceId,
				"TxnType": "Invoice",
			},
		},
	}

	// Construct the URL for creating the payment.
	requestUrl := fmt.Sprintf("%s/v3/company/%s/payment", s.qbConfig.Url, qbSettings.RealmId)
	// Perform the request.
	quickbooksResponse, err := s.performRequest(ctx, qbSettings, requestUrl, "POST", paymentData, true)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, fmt.Errorf("failed to save payment: %w", err)
	}

	// Unmarshal the response into QuickbooksSavePaymentResponse.
	var qbPaymentResponse interfaces.QuickbooksSavePaymentResponse
	err = json.Unmarshal(quickbooksResponse, &qbPaymentResponse)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	if qbPaymentResponse.Fault != nil {
		span.LogFields(log.Object("error", qbPaymentResponse.Fault))
		return nil, fmt.Errorf("error: %s", qbPaymentResponse.Fault.Error[0].Message)
	}

	return &qbPaymentResponse, nil
}

func (s *quickbooksService) ZeroPaymentLinkingJournalEntryToInvoice(ctx context.Context, quickbooksPaymentId, quickbooksCustomerId, quickbooksJournalEntryId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "QuickbooksService.ZeroPaymentLinkingJournalEntryToInvoice")
	defer span.Finish()
	tenant := common.GetTenantFromContext(ctx)
	tracing.TagTenant(span, tenant)
	span.LogFields(log.String("quickbooksPaymentId", quickbooksPaymentId),
		log.String("quickbooksCustomerId", quickbooksCustomerId),
		log.String("quickbooksJournalEntryId", quickbooksJournalEntryId))

	// Retrieve QuickBooks settings for the tenant.
	qbSettings, err := s.postgres.QuickbooksSettingsRepository.Get(ctx, tenant)
	if err != nil {
		tracing.TraceErr(span, err)
		return fmt.Errorf("failed to retrieve QuickBooks settings for tenant %s: %w", tenant, err)
	}
	if qbSettings == nil {
		err = errors.New("QuickBooks settings not found")
		tracing.TraceErr(span, err)
		return err
	}

	// Step 1: Fetch the existing payment to get its SyncToken.
	getURL := fmt.Sprintf("%s/v3/company/%s/payment/%s", s.qbConfig.Url, qbSettings.RealmId, quickbooksPaymentId)
	resp, err := s.performRequest(ctx, qbSettings, getURL, "GET", nil, true)
	if err != nil {
		tracing.TraceErr(span, err)
		return fmt.Errorf("failed to fetch existing payment: %w", err)
	}

	// Unmarshal the GET response into a Quickbooks payment response structure.
	var qbGetPaymentResp interfaces.QuickbooksGetPaymentResponse
	err = json.Unmarshal(resp, &qbGetPaymentResp)
	if err != nil {
		tracing.TraceErr(span, err)
		return fmt.Errorf("failed to unmarshal payment response: %w", err)
	}

	// Extract the SyncToken from the fetched payment.
	syncToken := qbGetPaymentResp.Payment.SyncToken
	if syncToken == "" {
		err = errors.New("no SyncToken found in payment record")
		tracing.TraceErr(span, err)
		return err
	}

	// Construct the payment payload.
	paymentData := map[string]interface{}{
		"Id":        quickbooksPaymentId,
		"SyncToken": syncToken,
		"TotalAmt":  0,
		"CustomerRef": map[string]interface{}{
			"value": quickbooksCustomerId,
		},
		"Line": []map[string]interface{}{
			{
				"Amount": 0,
				"LinkedTxn": []map[string]interface{}{
					{
						"TxnId":   quickbooksJournalEntryId,
						"TxnType": "JournalEntry",
					},
				},
			},
		},
	}

	// Construct the URL for creating the payment.
	requestUrl := fmt.Sprintf("%s/v3/company/%s/payment", s.qbConfig.Url, qbSettings.RealmId)
	// Perform the request.
	_, err = s.performRequest(ctx, qbSettings, requestUrl, "POST", paymentData, true)
	if err != nil {
		tracing.TraceErr(span, err)
		return fmt.Errorf("failed to save payment: %w", err)
	}

	return nil
}
