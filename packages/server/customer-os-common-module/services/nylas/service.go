package nylas

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/config"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	postgresEntity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	postgresRepository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"
)

type nylasService struct {
	config        *config.NylasConfig
	googleService interfaces.GoogleService
	httpClient    *http.Client
	postgres      *postgresRepository.Repositories
}

func NewNylasService(config *config.NylasConfig, googleService interfaces.GoogleService, postgres *postgresRepository.Repositories) interfaces.NylasService {
	return &nylasService{
		config:        config,
		googleService: googleService,
		httpClient:    &http.Client{},
		postgres:      postgres,
	}
}

// getProviderAccessToken handles token refresh and returns a valid refresh token
func (s *nylasService) getProviderRefreshToken(ctx context.Context, email string, provider enum.OAuthEmailProvider) (string, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "NylasService.getProviderRefreshToken")
	defer spans.Finish()
	spans.LogKV("email", email, "provider", provider)

	tenant := common.GetTenantFromContext(ctx)

	refreshToken := ""

	switch provider {
	case enum.ProviderGoogle:
		// get gmail service
		gmailService, err := s.googleService.GetGmailService(ctx, tenant, email)
		if err != nil {
			spans.TraceError(err)
			return "", fmt.Errorf("failed to get Gmail service: %v", err)
		}
		if gmailService == nil {
			spans.TraceError(fmt.Errorf("Gmail service is nil"))
			return "", fmt.Errorf("Gmail service is nil")
		}

		requiredScopes := []string{"https://www.googleapis.com/auth/calendar"}
		refreshToken, err = s.googleService.GetRefreshToken(ctx, tenant, email, provider, requiredScopes)
		if err != nil {
			spans.TraceError(err)
			return "", fmt.Errorf("failed to get refresh token: %v", err)
		}
	default:
		return "", fmt.Errorf("unsupported provider: %s", provider)
	}

	return refreshToken, nil
}

// Authenticate into Nylas and save the grant
func (s *nylasService) GrantAccess(ctx context.Context, email, userId, refreshToken string, nylasProvider interfaces.NylasProvider) (*postgresEntity.NylasGrant, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "NylasService.GrantAccess")
	defer spans.Finish()
	spans.LogKV("email", email, "userId", userId, "nylasProvider", nylasProvider)

	// validate tenant
	err := common.ValidateTenant(ctx)
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}
	tenant := common.GetTenantFromContext(ctx)

	// Check if account already exists
	nylasGrantEntity, err := s.postgres.NylasGrantRepository.GetByTenantAndEmail(ctx, tenant, email)
	if err != nil {
		spans.TraceError(err)
		return nil, fmt.Errorf("failed to check if Nylas account exists: %v", err)
	}
	if nylasGrantEntity == nil {
		nylasGrantEntity = &postgresEntity.NylasGrant{}
	}

	// Prepare request body
	requestBody := struct {
		Provider string `json:"provider"`
		Settings struct {
			RefreshToken string `json:"refresh_token"`
		} `json:"settings"`
		Email string `json:"email,omitempty"`
		State string `json:"state,omitempty"`
	}{
		Provider: string(nylasProvider),
		Settings: struct {
			RefreshToken string `json:"refresh_token"`
		}{
			RefreshToken: refreshToken,
		},
		Email: email,
		State: tenant, // Using tenant as state to track the origin
	}

	// Log request body (without sensitive data)
	spans.LogObjectAsJson("request", map[string]interface{}{
		"provider": nylasProvider,
		"email":    email,
		"state":    tenant,
	})

	// Marshal request body
	bodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		spans.TraceError(err)
		return nil, fmt.Errorf("failed to marshal request body: %v", err)
	}

	// Create request to Nylas v3 API
	req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("%s/v3/connect/custom", s.config.APIUrl), bytes.NewBuffer(bodyBytes))
	if err != nil {
		spans.TraceError(err)
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	// Set headers
	req.Header.Set("Authorization", "Bearer "+s.config.APIKey)
	req.Header.Set("Accept", "application/json, application/gzip")
	req.Header.Set("Content-Type", "application/json")

	// Make request
	resp, err := s.httpClient.Do(req)
	if err != nil {
		spans.TraceError(err)
		return nil, fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	bodyBytes, err = io.ReadAll(resp.Body)
	if err != nil {
		spans.TraceError(err)
		return nil, fmt.Errorf("failed to read error response body: %v", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		spans.LogKV("response", string(bodyBytes))
		spans.TraceError(fmt.Errorf("unexpected status code: %d", resp.StatusCode))
		return nil, fmt.Errorf("unexpected status code: %d, response: %s", resp.StatusCode, string(bodyBytes))
	}

	// Parse response
	var responseStruct postgresEntity.NylasGrantResponse

	if err := json.NewDecoder(bytes.NewReader(bodyBytes)).Decode(&responseStruct); err != nil {
		spans.TraceError(err)
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	// Save grant to database
	nylasGrantEntity.Tenant = tenant
	nylasGrantEntity.Email = email
	nylasGrantEntity.UserId = userId
	nylasGrantEntity.NylasGrantId = responseStruct.Data.ID
	nylasGrantEntity.NylasProvider = string(nylasProvider)
	nylasGrantEntity.NylasConnectResponse = string(bodyBytes)

	savedGrant, err := s.postgres.NylasGrantRepository.Save(ctx, nylasGrantEntity)
	if err != nil {
		spans.TraceError(err)
		return nil, fmt.Errorf("failed to save grant: %v", err)
	}

	return savedGrant, nil
}

func (s *nylasService) RevokeAccess(ctx context.Context, email string) error {
	spans, ctx := telemetry.StartServiceSpan(ctx, "NylasService.RevokeAccess")
	defer spans.Finish()
	spans.LogKV("email", email)

	// validate tenant
	err := common.ValidateTenant(ctx)
	if err != nil {
		spans.TraceError(err)
		return err
	}
	tenant := common.GetTenantFromContext(ctx)

	grant, err := s.postgres.NylasGrantRepository.GetByTenantAndEmail(ctx, tenant, email)
	if err != nil {
		spans.TraceError(err)
		return fmt.Errorf("failed to get Nylas grant: %v", err)
	}
	if grant == nil {
		return fmt.Errorf("Nylas grant not found")
	}

	// First remove the grant from Nylas
	req, err := http.NewRequestWithContext(ctx, "DELETE", fmt.Sprintf("%s/v3/grants/%s", s.config.APIUrl, grant.NylasGrantId), nil)
	if err != nil {
		spans.TraceError(err)
		return fmt.Errorf("failed to create request to revoke Nylas grant: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+s.config.APIKey)
	req.Header.Set("Accept", "application/json, application/gzip")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		spans.TraceError(err)
		return fmt.Errorf("failed to revoke Nylas grant: %v", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		spans.TraceError(err)
		return fmt.Errorf("failed to read response body: %v", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		spans.LogKV("response", string(bodyBytes))
		spans.TraceError(fmt.Errorf("unexpected status code: %d", resp.StatusCode))
		return fmt.Errorf("failed to revoke Nylas grant: status %d, response: %s", resp.StatusCode, string(bodyBytes))
	}

	// Delete from postgres database
	err = s.postgres.NylasGrantRepository.DeleteByGrantId(ctx, grant.NylasGrantId)
	if err != nil {
		spans.TraceError(err)
		return fmt.Errorf("failed to delete Nylas grant from database: %v", err)
	}

	return nil
}

// GetAccount retrieves a Nylas account for the given email
func (s *nylasService) GetGrant(ctx context.Context, email string) (*postgresEntity.NylasGrant, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "NylasService.GetGrant")
	defer spans.Finish()
	spans.LogKV("email", email)

	// validate tenant
	err := common.ValidateTenant(ctx)
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}
	tenant := common.GetTenantFromContext(ctx)
	userId := common.GetUserIdFromContext(ctx)

	var grant *postgresEntity.NylasGrant
	if email != "" {
		grant, err = s.postgres.NylasGrantRepository.GetByTenantAndEmail(ctx, tenant, email)
		if err != nil {
			spans.TraceError(err)
			return nil, fmt.Errorf("failed to get Nylas grant: %v", err)
		}
		if grant == nil {
			return nil, nil
		}
	} else if userId != "" {
		grant, err = s.postgres.NylasGrantRepository.GetByTenantAndUserId(ctx, tenant, userId)
		if err != nil {
			spans.TraceError(err)
			return nil, fmt.Errorf("failed to get Nylas grant: %v", err)
		}
		if grant == nil {
			return nil, nil
		}
	}

	return grant, nil
}

// ListCalendars lists all calendars for a user
func (s *nylasService) ListCalendars(ctx context.Context, email string) ([]*interfaces.NylasCalendar, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "NylasService.ListCalendars")
	defer spans.Finish()
	spans.LogKV("email", email)

	// Get Nylas grant ID
	grant, err := s.postgres.NylasGrantRepository.GetByTenantAndEmail(ctx, common.GetTenantFromContext(ctx), email)
	if err != nil {
		spans.TraceError(err)
		return nil, fmt.Errorf("failed to get Nylas grant ID: %v", err)
	}
	if grant == nil {
		return nil, fmt.Errorf("Nylas grant not found")
	}

	// Create request to Nylas v3 API
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/v3/grants/%s/calendars?limit=200", s.config.APIUrl, grant.NylasGrantId), nil)
	if err != nil {
		spans.TraceError(err)
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	// Set headers
	req.Header.Set("Authorization", "Bearer "+s.config.APIKey)
	req.Header.Set("Accept", "application/json")

	// Make request
	resp, err := s.httpClient.Do(req)
	if err != nil {
		spans.TraceError(err)
		return nil, fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		spans.TraceError(err)
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		spans.LogKV("response", string(bodyBytes))
		spans.TraceError(fmt.Errorf("unexpected status code: %d", resp.StatusCode))
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Parse response
	var response interfaces.NylasCalendarsResponse
	if err := json.NewDecoder(bytes.NewReader(bodyBytes)).Decode(&response); err != nil {
		spans.TraceError(err)
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return response.Data, nil
}

func (s *nylasService) GetDefaultCalendar(ctx context.Context, email string) (*interfaces.NylasCalendar, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "NylasService.GetDefaultCalendar")
	defer spans.Finish()
	spans.LogKV("email", email)

	// Get all calendars for the user
	calendars, err := s.ListCalendars(ctx, email)
	if err != nil {
		spans.TraceError(err)
		return nil, fmt.Errorf("failed to list calendars: %v", err)
	}

	if len(calendars) == 0 {
		return nil, nil
	}

	// First try to find a calendar marked as primary
	for _, calendar := range calendars {
		if calendar.IsPrimary {
			spans.LogKV("result.primaryFound", true)
			spans.LogObjectAsJson("result.calendar", calendar)
			return calendar, nil
		}
	}

	// If no primary calendar found, return the first one
	spans.LogKV("result.primaryFound", false)
	spans.LogObjectAsJson("result.calendar", calendars[0])
	return calendars[0], nil
}

func (s *nylasService) GetCalendarAvailability(ctx context.Context, email, calendarID string, startTime, endTime time.Time, bufferBefore, bufferAfter int) (*interfaces.NylasAvailabilityResponse, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "NylasService.GetCalendarAvailability")
	defer spans.Finish()
	spans.LogObjectAsJson("request", map[string]interface{}{
		"email":        email,
		"calendarID":   calendarID,
		"startTime":    startTime,
		"endTime":      endTime,
		"bufferBefore": bufferBefore,
		"bufferAfter":  bufferAfter,
	})

	if email == "" {
		return nil, fmt.Errorf("grantID is required")
	}
	if calendarID == "" {
		return nil, fmt.Errorf("calendarID is required")
	}

	// Create availability request
	request := interfaces.NylasAvailabilityRequest{
		Participants: []interfaces.NylasAvailabilityParticipant{
			{
				Email:       email,
				CalendarIds: []string{calendarID},
			},
		},
		StartTime:       startTime.UTC().Unix(),
		EndTime:         endTime.UTC().Unix(),
		IntervalMinutes: 5, // Match our time slot generation interval
		DurationMinutes: 5, // Default duration
		RoundTo:         5, // Match our time slot generation interval
		AvailabilityRules: interfaces.NylasAvailabilityRules{
			AvailabilityMethod: "max-availability",
			Buffer: struct {
				Before int `json:"before"`
				After  int `json:"after"`
			}{
				Before: bufferBefore,
				After:  bufferAfter,
			},
			TentativeAsBusy: true,
		},
	}

	// Marshal request body
	bodyBytes, err := json.Marshal(request)
	if err != nil {
		spans.TraceError(err)
		return nil, fmt.Errorf("failed to marshal request body: %v", err)
	}

	// Create request to Nylas v3 API
	req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("%s/v3/calendars/availability", s.config.APIUrl), bytes.NewBuffer(bodyBytes))
	if err != nil {
		spans.TraceError(err)
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	// Set headers
	req.Header.Set("Authorization", "Bearer "+s.config.APIKey)
	req.Header.Set("Accept", "application/json, application/gzip")
	req.Header.Set("Content-Type", "application/json")

	// Make request
	resp, err := s.httpClient.Do(req)
	if err != nil {
		spans.TraceError(err)
		return nil, fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	respBodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		spans.TraceError(err)
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		spans.LogKV("response", string(respBodyBytes))
		spans.TraceError(fmt.Errorf("unexpected status code: %d", resp.StatusCode))
		return nil, fmt.Errorf("unexpected status code: %d, response: %s", resp.StatusCode, string(respBodyBytes))
	}

	// Parse response
	var response interfaces.NylasAvailabilityResponse
	if err := json.NewDecoder(bytes.NewReader(respBodyBytes)).Decode(&response); err != nil {
		spans.TraceError(err)
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return &response, nil
}
