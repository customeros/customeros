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
	req.Header.Set("Authorization", "Bearer eyJlbmMiOiJBMTI4Q0JDLUhTMjU2IiwiYWxnIjoiZGlyIn0..Zotbwr8fTUtS9jhALhsoDQ.2Ytm6ENIA1Gh7sxAPxeLaQqd82zlr-34hXHEUwsfpK1w5j2pudzVO5nPXBssr3Or7RskNuKuxyUzH1Tot7acQYhoN_w5D6Xu2aQ05gxP14FrBEDJxW7s3bjMdTZO9TL1qNeCsMEVZEArEeNWPImsilR-BfngMarhrFAjcliA2-GVkI3XdeGwnWtEGqCJOyL7hOi0bfP8mg5gJe8pOZhRpmf_OsZMXvSTy6-N5yaOJ0PB6Y8flakFHNIQWXVvmmLxzFe1zl8mADTnkv3s4h5zaePHvprk7X2tZSWAChwFWB2G9HBW6OXMWg_TIShtOH7uIcXxaw_epJ6CAV9fRaiaJst-4Z6e6KNBgz2BwrHsGzgs-EPrQxvnPvIhTtcrt7okMSdIzKbzbh0x8ZRkIG08xUPpHCcBw0jqCeKVE7u1Lp1jKM1JuFUzQiKDA36iUxE6R4GMMcbe8Qb3x5b5huAnMNlzOUX3ZSZbvvk-Ad7Sp2h0hqm0v0CjFqPRxjHHzyLnwzYzQZHb4Vmd1yApmUsvFWZA3tgZjFFOZB6PEAAXiLDXBjxE7T5Wa9ULHL2QAsyWfwgaqL5H1angciPsPxuaMBP0LbiB_0D0MDyC3m2tZkvI8njvRhrpFRoOSafvw7c0Mi6M_lbLGZQO3c9daPK2sVTq0ZnVRmx23AzHApTDh3bObfuUDj0f-Z1ytXiYFD2mQiE6d8ARVnyMWZDXBLSMJl092DIaBRkLXD6oCXU9WB4b1y198ny91bLOTB5m_Xey.AFausD6vgvRHnr4Xo4RoAQ")

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
