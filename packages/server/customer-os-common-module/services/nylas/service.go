package nylas

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/config"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
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

// getProviderAccessToken handles token refresh and returns a valid access token
func (s *nylasService) getProviderAccessToken(ctx context.Context, email string, provider enum.OAuthEmailProvider) (string, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "NylasService.getProviderAccessToken")
	defer spans.Finish()
	spans.LogKV("email", email, "provider", provider)

	tenant := common.GetTenantFromContext(ctx)

	accessToken := ""

	switch provider {
	case enum.ProviderGoogle:
		// Get Gmail service which will handle token refresh
		gmailService, err := s.googleService.GetGmailService(ctx, email, tenant)
		if err != nil {
			spans.TraceError(err)
			return "", fmt.Errorf("failed to get Gmail service: %v", err)
		}
		if gmailService == nil {
			return "", fmt.Errorf("failed to get Gmail service for email: %s", email)
		}

		// Get the access token from the Gmail service
		// The Gmail service will handle token refresh if needed
		accessToken, err = s.googleService.GetAccessToken(ctx, tenant, email)
		if err != nil {
			spans.TraceError(err)
			return "", fmt.Errorf("failed to get access token: %v", err)
		}
	default:
		return "", fmt.Errorf("unsupported provider: %s", provider)
	}

	return accessToken, nil
}

// ConnectAccount connects a new Nylas account
func (s *nylasService) ConnectAccount(ctx context.Context, email string, provider enum.OAuthEmailProvider) (*postgres_entity.NylasAccount, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "NylasService.ConnectAccount")
	defer spans.Finish()
	spans.LogKV("email", email, "provider", provider)

	// validate tenant
	err := common.ValidateTenant(ctx)
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}
	tenant := common.GetTenantFromContext(ctx)

	// Check if account already exists
	existingAccount, err := s.postgres.NylasAccountRepository.GetByTenantAndEmail(ctx, tenant, email)
	if err == nil && existingAccount != nil {
		return existingAccount, nil
	}

	// Connect new account
	accountID, err := s.prepareNylasAccountID(ctx, email, provider)
	if err != nil {
		spans.TraceError(err)
		return nil, fmt.Errorf("failed to connect account: %v", err)
	}

	// Save account to database
	account := &postgres_entity.NylasAccount{
		Tenant:         tenant,
		Email:          email,
		NylasAccountId: accountID,
		Provider:       provider.String(),
	}

	savedAccount, err := s.postgres.NylasAccountRepository.Save(ctx, account)
	if err != nil {
		spans.TraceError(err)
		return nil, fmt.Errorf("failed to save account: %v", err)
	}

	return savedAccount, nil
}

func (s *nylasService) DisconnectAccount(ctx context.Context, email string) error {
	spans, ctx := telemetry.StartServiceSpan(ctx, "NylasService.DisconnectAccount")
	defer spans.Finish()
	spans.LogKV("email", email)

	// validate tenant
	err := common.ValidateTenant(ctx)
	if err != nil {
		spans.TraceError(err)
		return err
	}
	tenant := common.GetTenantFromContext(ctx)

	account, err := s.postgres.NylasAccountRepository.GetByTenantAndEmail(ctx, tenant, email)
	if err != nil {
		spans.TraceError(err)
		return fmt.Errorf("failed to get Nylas account: %v", err)
	}
	if account == nil {
		return fmt.Errorf("Nylas account not found")
	}

	// First remove the account from Nylas
	req, err := http.NewRequestWithContext(ctx, "DELETE", fmt.Sprintf("%s/v3/connect/accounts/%s", s.config.APIUrl, account.NylasAccountId), nil)
	if err != nil {
		spans.TraceError(err)
		return fmt.Errorf("failed to create request to remove Nylas account: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+s.config.APIKey)
	req.Header.Set("Accept", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		spans.TraceError(err)
		return fmt.Errorf("failed to remove Nylas account: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		spans.TraceError(fmt.Errorf("unexpected status code: %d", resp.StatusCode))
		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("failed to remove Nylas account: unexpected status code %d", resp.StatusCode)
		}
	}

	// Delete from postgres database
	err = s.postgres.NylasAccountRepository.Delete(ctx, account)
	if err != nil {
		spans.TraceError(err)
		return fmt.Errorf("failed to delete Nylas account from database: %v", err)
	}

	return nil
}

// GetAccount retrieves a Nylas account for the given email
func (s *nylasService) GetAccount(ctx context.Context, email string) (*postgres_entity.NylasAccount, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "NylasService.GetAccount")
	defer spans.Finish()
	spans.LogKV("email", email)

	// validate tenant
	err := common.ValidateTenant(ctx)
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}
	tenant := common.GetTenantFromContext(ctx)

	account, err := s.postgres.NylasAccountRepository.GetByTenantAndEmail(ctx, tenant, email)
	if err != nil {
		spans.TraceError(err)
		return nil, fmt.Errorf("failed to get Nylas account: %v", err)
	}
	if account == nil {
		return nil, nil
	}

	return account, nil
}

// prepareNylasAccountID prepares the Nylas account ID for a user
func (s *nylasService) prepareNylasAccountID(ctx context.Context, email string, provider enum.OAuthEmailProvider) (string, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "NylasService.prepareNylasAccountID")
	defer spans.Finish()
	spans.LogKV("email", email, "provider", provider)

	// Get tenant from context
	tenant := common.GetTenantFromContext(ctx)
	if tenant == "" {
		spans.TraceError(fmt.Errorf("tenant not found in context"))
		return "", fmt.Errorf("tenant not found in context")
	}

	// Check if account exists in database
	account, err := s.postgres.NylasAccountRepository.GetByTenantAndEmail(ctx, tenant, email)
	if err != nil {
		spans.TraceError(err)
		return "", fmt.Errorf("failed to get Nylas account: %v", err)
	}

	if account != nil {
		return account.NylasAccountId, nil
	}

	// Connect new account
	accountID, err := s.ConnectAccount(ctx, email, provider)
	if err != nil {
		spans.TraceError(err)
		return "", fmt.Errorf("failed to connect account: %v", err)
	}

	return accountID.NylasAccountId, nil
}

// ListCalendars lists all calendars for a user
func (s *nylasService) ListCalendars(ctx context.Context, email string, provider enum.OAuthEmailProvider) ([]*interfaces.Calendar, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "NylasService.ListCalendars")
	defer spans.Finish()
	spans.LogKV("email", email, "provider", provider)

	// Get Nylas account ID
	accountID, err := s.prepareNylasAccountID(ctx, email, provider)
	if err != nil {
		spans.TraceError(err)
		return nil, fmt.Errorf("failed to get Nylas account ID: %v", err)
	}

	// Create request to Nylas v3 API
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/v3/calendars?account_id=%s", s.config.APIUrl, accountID), nil)
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

	if resp.StatusCode != http.StatusOK {
		spans.TraceError(fmt.Errorf("unexpected status code: %d", resp.StatusCode))
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Parse response
	var calendars []*interfaces.Calendar
	if err := json.NewDecoder(resp.Body).Decode(&calendars); err != nil {
		spans.TraceError(err)
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return calendars, nil
}

// Update ListEvents to use prepareNylasAccountID
func (s *nylasService) ListEvents(ctx context.Context, calendarID string, startTime, endTime time.Time) ([]*interfaces.CalendarEvent, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "NylasService.ListEvents")
	defer spans.Finish()
	spans.LogObjectAsJson("request", map[string]interface{}{
		"calendarID": calendarID,
		"startTime":  startTime,
		"endTime":    endTime,
	})

	// Get email from context
	email := ctx.Value("email").(string)
	if email == "" {
		return nil, fmt.Errorf("no email found in context")
	}

	// Get provider from context
	provider := ctx.Value("provider").(enum.OAuthEmailProvider)
	if provider == "" {
		return nil, fmt.Errorf("no provider found in context")
	}

	// Get Nylas account ID for this email
	accountID, err := s.prepareNylasAccountID(ctx, email, provider)
	if err != nil {
		spans.TraceError(err)
		return nil, fmt.Errorf("failed to get Nylas account ID: %v", err)
	}

	// Create request to Nylas v3 API
	url := fmt.Sprintf("%s/v3/events?account_id=%s&calendar_id=%s&start=%d&end=%d",
		s.config.APIUrl,
		accountID,
		calendarID,
		startTime.Unix(),
		endTime.Unix(),
	)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		spans.TraceError(err)
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	// Set headers - use Nylas API key for authentication
	req.Header.Set("Authorization", "Bearer "+s.config.APIKey)
	req.Header.Set("Accept", "application/json")

	// Make request
	resp, err := s.httpClient.Do(req)
	if err != nil {
		spans.TraceError(err)
		return nil, fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		spans.TraceError(fmt.Errorf("unexpected status code: %d", resp.StatusCode))
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Parse response
	var response struct {
		Data []struct {
			ID          string `json:"id"`
			Title       string `json:"title"`
			Description string `json:"description"`
			Start       int64  `json:"start"`
			End         int64  `json:"end"`
			Location    string `json:"location"`
			Status      string `json:"status"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		spans.TraceError(err)
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	// Convert to our interface type
	result := make([]*interfaces.CalendarEvent, len(response.Data))
	for i, event := range response.Data {
		result[i] = &interfaces.CalendarEvent{
			ID:          event.ID,
			Title:       event.Title,
			Description: event.Description,
			StartTime:   time.Unix(event.Start, 0),
			EndTime:     time.Unix(event.End, 0),
			Location:    event.Location,
			Status:      event.Status,
		}
	}

	return result, nil
}

func (s *nylasService) CreateEvent(ctx context.Context, calendarID string, event *interfaces.CalendarEvent) (*interfaces.CalendarEvent, error) {
	// TODO: Implement Nylas API call to create event
	return nil, nil
}

func (s *nylasService) UpdateEvent(ctx context.Context, calendarID string, eventID string, event *interfaces.CalendarEvent) (*interfaces.CalendarEvent, error) {
	// TODO: Implement Nylas API call to update event
	return nil, nil
}

func (s *nylasService) DeleteEvent(ctx context.Context, calendarID string, eventID string) error {
	// TODO: Implement Nylas API call to delete event
	return nil
}

func (s *nylasService) GetEvent(ctx context.Context, calendarID string, eventID string) (*interfaces.CalendarEvent, error) {
	// TODO: Implement Nylas API call to get event
	return nil, nil
}

func (s *nylasService) GetCalendar(ctx context.Context, calendarID string) (*interfaces.Calendar, error) {
	// TODO: Implement Nylas API call to get calendar
	return nil, nil
}

// getNylasProvider maps our provider to Nylas provider
func (s *nylasService) getNylasProvider(provider enum.OAuthEmailProvider) (string, error) {
	switch provider {
	case enum.ProviderGoogle:
		return string(interfaces.NylasProviderGoogle), nil
	default:
		return "", fmt.Errorf("unsupported provider for Nylas: %s", provider)
	}
}
