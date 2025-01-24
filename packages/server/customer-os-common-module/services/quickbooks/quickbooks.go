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

		entity := postgres_entity.QuickbooksSettingsEntity{
			Tenant:                tenant,
			RealmId:               realmId,
			AccessToken:           quickbooksResponse.AccessToken,
			AccessTokenExpiresIn:  quickbooksResponse.ExpiresIn,
			AccessTokenExpiresAt:  now.Add(time.Duration(quickbooksResponse.ExpiresIn) * time.Second),
			RefreshToken:          quickbooksResponse.RefreshToken,
			RefreshTokenExpiresIn: quickbooksResponse.XRefreshTokenExpiresIn,
			RefreshTokenExpiresAt: now.Add(time.Duration(quickbooksResponse.XRefreshTokenExpiresIn) * time.Second),
			RefreshTokenExpired:   false,
		}

		if quickbooksSettingsEntity != nil {
			entity.Id = quickbooksSettingsEntity.Id
		}

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

//TODO create a CustomerOS invoice Account in QB and link services to it

func (s *quickbooksService) SaveProduct(ctx context.Context, id, productName string) (*interfaces.QuickbooksSaveProductResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "QuickbooksService.SaveProduct")
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
		"Name": productName,
		"Type": "Service",
		"IncomeAccountRef": map[string]interface{}{
			"value": "1", // TODO HOW DO WE GET THIS ID??
		},
	}

	if id != "" {
		request["Id"] = id
	}

	payload, err := json.Marshal(request)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	// Create a new HTTP request
	req, err := http.NewRequest("POST", fmt.Sprintf("https://sandbox-quickbooks.api.intuit.com/v3/company/%s/item?minorversion=73", quickbooksSettingsEntity.RealmId), bytes.NewBuffer(payload))
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

	// Read response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	// convert body to OauthSlackResponse
	var quickbooksResponse interfaces.QuickbooksSaveProductResponse
	err = json.Unmarshal(body, &quickbooksResponse)
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

	payload, err := json.Marshal(request)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	// Create a new HTTP request
	req, err := http.NewRequest("POST", fmt.Sprintf("https://sandbox-quickbooks.api.intuit.com/v3/company/%s/customer?minorversion=73", quickbooksSettingsEntity.RealmId), bytes.NewBuffer(payload))
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

	// Read response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	// convert body to OauthSlackResponse
	var quickbooksResponse interfaces.QuickbooksSaveCustomerResponse
	err = json.Unmarshal(body, &quickbooksResponse)
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

func (s *quickbooksService) SaveInvoice(ctx context.Context, customerId string, invoiceDate time.Time, lines []interfaces.QuickbooksInvoiceLine) (*interfaces.QuickbooksSaveCustomerResponse, error) {
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

	payload, err := json.Marshal(request)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	// Create a new HTTP request
	req, err := http.NewRequest("POST", fmt.Sprintf("https://sandbox-quickbooks.api.intuit.com/v3/company/%s/customer?minorversion=73", quickbooksSettingsEntity.RealmId), bytes.NewBuffer(payload))
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

	// Read response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	// convert body to OauthSlackResponse
	var quickbooksResponse interfaces.QuickbooksSaveCustomerResponse
	err = json.Unmarshal(body, &quickbooksResponse)
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
