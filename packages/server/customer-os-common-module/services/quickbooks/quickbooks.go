package quickbooks

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
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
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
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

		quickbooksSettingsEntity, err := s.postgres.QuickbooksSettingsRepository.Get(ctx, tenant)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}

		entity := postgres_entity.QuickbooksSettingsEntity{}
		if quickbooksSettingsEntity != nil {
			entity = *quickbooksSettingsEntity
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

	span.LogFields(log.String("revokeResponse", string(bodyBytes)))

	// Check for a successful response
	if resp.StatusCode != http.StatusOK {
		errMsg := fmt.Sprintf("revoke quickbooks request returned status %d", resp.StatusCode)
		tracing.TraceErr(span, fmt.Errorf(errMsg))
		span.LogFields(log.String("quickbooks.response", string(bodyBytes)))
		return fmt.Errorf(errMsg)
	}

	err = s.postgres.QuickbooksSettingsRepository.Delete(ctx, tenant)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
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

	quickbooksSettingsEntity, err := s.postgres.QuickbooksSettingsRepository.Get(ctx, tenant)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	if quickbooksSettingsEntity == nil {
		span.LogFields(log.String("result.error", "Quickbooks settings not found"))
		return nil, nil
	}

	if quickbooksSettingsEntity.SalesAccountId == "" {
		//search for the sales account
		searchSalesAccountUrl := s.qbConfig.Url + fmt.Sprintf("/v3/company/%s/query?query=select+Id+from+Account+where+Name='CustomerOS+Sales'", quickbooksSettingsEntity.RealmId)
		searchSalesAccountResponse, err := s.performRequest(ctx, quickbooksSettingsEntity, searchSalesAccountUrl, "POST", nil, true)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}

		var searchSalesAccount interfaces.QuickbooksSearchAccountResponse
		err = json.Unmarshal(searchSalesAccountResponse, &searchSalesAccount)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}

		if len(searchSalesAccount.QueryResponse.Account) > 0 {
			quickbooksSettingsEntity.SalesAccountId = searchSalesAccount.QueryResponse.Account[0].Id
			_, err = s.postgres.QuickbooksSettingsRepository.Save(ctx, *quickbooksSettingsEntity)
			if err != nil {
				tracing.TraceErr(span, err)
				return nil, err
			}
		} else {
			salesAccountUrl := fmt.Sprintf(s.qbConfig.Url+"/v3/company/%s/account", quickbooksSettingsEntity.RealmId)
			salesAccountRequest := map[string]interface{}{
				"Name":        "CustomerOS Sales",
				"AccountType": "Income",
			}

			qbAccountResponse, err := s.performRequest(ctx, quickbooksSettingsEntity, salesAccountUrl, "POST", salesAccountRequest, true)
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

			quickbooksSettingsEntity.SalesAccountId = qbAccount.Account.Id
			_, err = s.postgres.QuickbooksSettingsRepository.Save(ctx, *quickbooksSettingsEntity)
			if err != nil {
				tracing.TraceErr(span, err)
				return nil, err
			}
		}
	}

	var qbSaveProductRequest map[string]interface{}
	var qbProduct interfaces.QuickbooksGetProductResponse

	if id != "" {
		productByIdUrl := fmt.Sprintf(s.qbConfig.Url+"/v3/company/%s/item/%s", quickbooksSettingsEntity.RealmId, id)
		qbProductResponse, err := s.performRequest(ctx, quickbooksSettingsEntity, productByIdUrl, "GET", nil, true)
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
				"value": quickbooksSettingsEntity.SalesAccountId,
			},
		}
	}

	requestUrl := fmt.Sprintf(s.qbConfig.Url+"/v3/company/%s/item", quickbooksSettingsEntity.RealmId)
	qbResponse, err := s.performRequest(ctx, quickbooksSettingsEntity, requestUrl, "POST", qbSaveProductRequest, true)
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

func (s *quickbooksService) SaveCustomer(ctx context.Context, id, customerName string) (*interfaces.QuickbooksSaveCustomerResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "QuickbooksService.SaveCustomer")
	defer span.Finish()

	tenant := common.GetTenantFromContext(ctx)

	quickbooksSettingsEntity, err := s.postgres.QuickbooksSettingsRepository.Get(ctx, tenant)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	//TODO check how we return in case of missing integration
	if quickbooksSettingsEntity == nil {
		span.LogFields(log.String("error", "Quickbooks settings not found"))
		return nil, nil
	}

	request := map[string]interface{}{
		"GivenName": customerName,
	}

	if id != "" {
		request["Id"] = id
	}

	resp, err := s.performRequest(ctx, quickbooksSettingsEntity, s.qbConfig.Url+"/v3/company/"+quickbooksSettingsEntity.RealmId+"/customer", "POST", request, true)
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

func (s *quickbooksService) SaveInvoice(ctx context.Context, customerId string, invoiceDate time.Time, lines []interfaces.QuickbooksInvoiceLine) (*interfaces.QuickbooksSaveInvoiceResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "QuickbooksService.SaveInvoice")
	defer span.Finish()

	tenant := common.GetTenantFromContext(ctx)

	quickbooksSettingsEntity, err := s.postgres.QuickbooksSettingsRepository.Get(ctx, tenant)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	//TODO check how we return in case of missing integration
	if quickbooksSettingsEntity == nil {
		span.LogFields(log.String("error", "Quickbooks settings not found"))
		return nil, nil
	}

	request := map[string]interface{}{
		"CustomerRef": map[string]interface{}{
			"value": customerId,
		},
		"TxnDate": invoiceDate.Format("2006/01/25"),
		"Line":    lines,
	}

	resp, err := s.performRequest(ctx, quickbooksSettingsEntity, s.qbConfig.Url+"/v3/company/"+quickbooksSettingsEntity.RealmId+"/invoice", "POST", request, true)
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

	quickbooksSettingsEntity, err := s.postgres.QuickbooksSettingsRepository.Get(ctx, tenant)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	//TODO check how we return in case of missing integration
	if quickbooksSettingsEntity == nil {
		span.LogFields(log.String("error", "Quickbooks settings not found"))
		return nil, nil
	}

	request := map[string]interface{}{
		"Id": invoiceId,
	}

	if invoiceId != "" {
		invoiceByIdUrl := fmt.Sprintf(s.qbConfig.Url+"/v3/company/%s/invoice/%s", quickbooksSettingsEntity.RealmId, invoiceId)
		qbInvoiceResponse, err := s.performRequest(ctx, quickbooksSettingsEntity, invoiceByIdUrl, "GET", request, true)
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

	resp, err := s.performRequest(ctx, quickbooksSettingsEntity, s.qbConfig.Url+"/v3/company/"+quickbooksSettingsEntity.RealmId+"/invoice?operation=void", "POST", request, true)
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

func (s *quickbooksService) PayInvoice(ctx context.Context, customerId, invoiceId string, totalAmount float64) (*interfaces.QuickbooksSavePaymentResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "QuickbooksService.PayInvoice")
	defer span.Finish()

	tenant := common.GetTenantFromContext(ctx)

	quickbooksSettingsEntity, err := s.postgres.QuickbooksSettingsRepository.Get(ctx, tenant)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	//TODO check how we return in case of missing integration
	if quickbooksSettingsEntity == nil {
		span.LogFields(log.String("error", "Quickbooks settings not found"))
		return nil, nil
	}

	request := map[string]interface{}{
		"CustomerRef": map[string]interface{}{
			"value": customerId,
		},
		"TotalAmt": totalAmount,
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

	resp, err := s.performRequest(ctx, quickbooksSettingsEntity, s.qbConfig.Url+"/v3/company/"+quickbooksSettingsEntity.RealmId+"/payment", "POST", request, true)
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

func (s *quickbooksService) performRequest(ctx context.Context, quickbooksSettingsEntity *postgres_entity.QuickbooksSettingsEntity, requestUrl string, requestMethod string, requestBody map[string]interface{}, rerunOnTokenRefresh bool) ([]byte, error) {
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
	req.Header.Set("Authorization", "Bearer "+quickbooksSettingsEntity.AccessToken)

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
		tracing.TraceErr(span, err)
		return nil, err
	}

	fault := quickbooksCheckFaultResponse.Fault
	if fault != nil {
		if (*fault).Type == "AUTHENTICATION" {
			//refresh token
			requestData := url.Values{}
			requestData.Set("grant_type", "refresh_token")
			requestData.Set("refresh_token", quickbooksSettingsEntity.RefreshToken)

			quickbooksSettingsEntity, err = s.GetAndStoreAccessToken(ctx, quickbooksSettingsEntity.RealmId, requestData)
			if err != nil {
				tracing.TraceErr(span, err)
				return nil, err
			}

			// retry the HTTP request
			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				tracing.TraceErr(span, err)
				return nil, err
			}
			defer resp.Body.Close()

			// Read response body
			bodyBytes, err = ioutil.ReadAll(resp.Body)
			if err != nil {
				tracing.TraceErr(span, err)
				return nil, err
			}

			err = json.Unmarshal(bodyBytes, &quickbooksCheckFaultResponse)
			if err != nil {
				tracing.TraceErr(span, err)
				return nil, err
			}
			// re-run original request
			if rerunOnTokenRefresh {
				return s.performRequest(ctx, quickbooksSettingsEntity, requestUrl, requestMethod, requestBody, false)
			}

		} else {
			span.LogFields(log.Object("error", fault))

			return nil, fmt.Errorf("error: %s", fault.Error[0].Message)
		}
	}

	return bodyBytes, nil
}

func (s *quickbooksService) QuickbooksConnected(ctx context.Context) (bool, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "QuickbooksService.QuickbooksConnected")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	tenant := common.GetTenantFromContext(ctx)

	quickbooksSettingsEntity, err := s.postgres.QuickbooksSettingsRepository.Get(ctx, tenant)
	if err != nil {
		tracing.TraceErr(span, err)
		return false, err
	}

	return quickbooksSettingsEntity != nil, nil
}
