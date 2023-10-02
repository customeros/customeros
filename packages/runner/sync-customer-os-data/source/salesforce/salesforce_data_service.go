package salesforce

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/common"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/config"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/entity"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/logger"
	app_repo "github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/repository"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/source"
	source_entity "github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/source/entity"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/source/repository"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"golang.org/x/oauth2"
	"gorm.io/gorm"
	"net/http"
	"strings"
)

const (
	UserTableSuffix        = "user"
	AccountTableSuffix     = "account"
	ContactTableSuffix     = "contact"
	LeadTableSuffix        = "lead"
	FeeditemTableSuffix    = "feeditem"
	OpportunityTableSuffix = "opportunity"
	ContentnoteTableSuffix = "contentnote"
)

var sourceTableSuffixByDataType = map[string][]string{
	string(common.USERS):         {UserTableSuffix},
	string(common.ORGANIZATIONS): {AccountTableSuffix, LeadTableSuffix, OpportunityTableSuffix},
	string(common.CONTACTS):      {ContactTableSuffix, LeadTableSuffix},
	string(common.LOG_ENTRIES):   {FeeditemTableSuffix, ContentnoteTableSuffix},
}

type salesforceDataService struct {
	airbyteStoreDb *config.RawDataStoreDB
	repositories   *app_repo.Repositories
	tenant         string
	instance       string
	processingIds  map[string]source.ProcessingEntity
	dataFuncs      map[common.SyncedEntityType]func(context.Context, int, string) []any
	log            logger.Logger
}

func NewSalesforceDataService(airbyteStoreDb *config.RawDataStoreDB, repositories *app_repo.Repositories, tenant string, log logger.Logger) source.SourceDataService {
	dataService := salesforceDataService{
		airbyteStoreDb: airbyteStoreDb,
		repositories:   repositories,
		tenant:         tenant,
		processingIds:  map[string]source.ProcessingEntity{},
		log:            log,
	}
	dataService.dataFuncs = map[common.SyncedEntityType]func(context.Context, int, string) []any{}
	dataService.dataFuncs[common.USERS] = dataService.GetUsersForSync
	dataService.dataFuncs[common.ORGANIZATIONS] = dataService.GetOrganizationsForSync
	dataService.dataFuncs[common.CONTACTS] = dataService.GetContactsForSync
	dataService.dataFuncs[common.LOG_ENTRIES] = dataService.GetLogEntriesForSync
	return &dataService
}

func (s *salesforceDataService) GetDataForSync(ctx context.Context, dataType common.SyncedEntityType, batchSize int, runId string) []interface{} {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SalesforceDataService.GetDataForSync")
	defer span.Finish()
	tracing.SetDefaultSyncServiceSpanTags(ctx, span)
	span.LogFields(log.String("dataType", string(dataType)), log.Int("batchSize", batchSize))

	if ok := s.dataFuncs[dataType]; ok != nil {
		return s.dataFuncs[dataType](ctx, batchSize, runId)
	} else {
		s.log.Infof("No %s data function for %s", s.SourceId(), dataType)
		return nil
	}
}

func (s *salesforceDataService) Init() {
	err := s.getDb().AutoMigrate(&source_entity.SyncStatusForAirbyte{})
	if err != nil {
		s.log.Error(err)
	}
}

func (s *salesforceDataService) getDb() *gorm.DB {
	schemaName := s.SourceId()

	if len(s.instance) > 0 {
		schemaName = schemaName + "_" + s.instance
	}
	schemaName = schemaName + "_" + s.tenant
	return s.airbyteStoreDb.GetDBHandler(&config.Context{
		Schema: schemaName,
	})
}

func (s *salesforceDataService) SourceId() string {
	return string(entity.AirbyteSourceSalesforce)
}

func (s *salesforceDataService) Close() {
	s.processingIds = make(map[string]source.ProcessingEntity)
}

func (s *salesforceDataService) GetUsersForSync(ctx context.Context, batchSize int, runId string) []any {
	s.processingIds = make(map[string]source.ProcessingEntity)
	currentEntity := string(common.USERS)
	var users []any
	for _, sourceTableSuffix := range sourceTableSuffixByDataType[currentEntity] {
		airbyteRecords, err := repository.GetAirbyteUnprocessedRawRecords(ctx, s.getDb(), batchSize, runId, currentEntity, sourceTableSuffix)
		if err != nil {
			s.log.Error(err)
			return nil
		}
		for _, v := range airbyteRecords {
			if len(users) >= batchSize {
				break
			}
			outputJSON, err := MapUser(v.AirbyteData)
			user, err := source.MapJsonToUser(outputJSON, v.AirbyteAbId, s.SourceId())
			if err != nil {
				user = entity.UserData{
					BaseData: entity.BaseData{
						SyncId: v.AirbyteAbId,
					},
				}
			}
			s.processingIds[v.AirbyteAbId] = source.ProcessingEntity{
				ExternalId:  user.ExternalId,
				Entity:      currentEntity,
				TableSuffix: sourceTableSuffix,
			}
			users = append(users, user)
		}
	}
	return users
}

func (s *salesforceDataService) GetOrganizationsForSync(ctx context.Context, batchSize int, runId string) []any {
	s.processingIds = make(map[string]source.ProcessingEntity)
	currentEntity := string(common.ORGANIZATIONS)

	var organizations []any
	for _, sourceTableSuffix := range sourceTableSuffixByDataType[currentEntity] {
		airbyteRecords, err := repository.GetAirbyteUnprocessedRawRecords(ctx, s.getDb(), batchSize, runId, currentEntity, sourceTableSuffix)
		if err != nil {
			s.log.Error(err)
			return nil
		}
		for _, v := range airbyteRecords {
			if len(organizations) >= batchSize {
				break
			}
			outputJSON, err := MapOrganization(v.AirbyteData)
			organization, err := source.MapJsonToOrganization(outputJSON, v.AirbyteAbId, s.SourceId())
			if err != nil {
				organization = entity.OrganizationData{
					BaseData: entity.BaseData{
						SyncId: v.AirbyteAbId,
					},
				}
			}

			s.processingIds[v.AirbyteAbId] = source.ProcessingEntity{
				ExternalId:  organization.ExternalId,
				Entity:      currentEntity,
				TableSuffix: sourceTableSuffix,
			}
			organizations = append(organizations, organization)
		}
	}
	return organizations
}

func (s *salesforceDataService) GetContactsForSync(ctx context.Context, batchSize int, runId string) []any {
	s.processingIds = make(map[string]source.ProcessingEntity)
	currentEntity := string(common.CONTACTS)

	var contacts []any
	for _, sourceTableSuffix := range sourceTableSuffixByDataType[currentEntity] {
		airbyteRecords, err := repository.GetAirbyteUnprocessedRawRecords(ctx, s.getDb(), batchSize, runId, currentEntity, sourceTableSuffix)
		if err != nil {
			s.log.Error(err)
			return nil
		}
		for _, v := range airbyteRecords {
			if len(contacts) >= batchSize {
				break
			}
			outputJSON, err := MapContact(v.AirbyteData)
			contact, err := source.MapJsonToContact(outputJSON, v.AirbyteAbId, s.SourceId())
			if err != nil {
				contact = entity.ContactData{
					BaseData: entity.BaseData{
						SyncId: v.AirbyteAbId,
					},
				}
			}

			s.processingIds[v.AirbyteAbId] = source.ProcessingEntity{
				ExternalId:  contact.ExternalId,
				Entity:      currentEntity,
				TableSuffix: sourceTableSuffix,
			}
			contacts = append(contacts, contact)
		}
	}
	return contacts
}

func (s *salesforceDataService) GetLogEntriesForSync(ctx context.Context, batchSize int, runId string) []any {
	s.processingIds = make(map[string]source.ProcessingEntity)
	currentEntity := string(common.LOG_ENTRIES)

	tenantSettings, err := s.repositories.TenantSettingsRepository.GetTenantSettings(ctx, s.tenant)
	if err != nil {
		s.log.Errorf("Failed to get tenant settings for tenant %s: %v", s.tenant, err.Error())
	}

	var logEntries []any
	for _, sourceTableSuffix := range sourceTableSuffixByDataType[currentEntity] {
		airbyteRecords, err := repository.GetAirbyteUnprocessedRawRecords(ctx, s.getDb(), batchSize, runId, currentEntity, sourceTableSuffix)
		if err != nil {
			s.log.Error(err)
			return nil
		}
		for _, v := range airbyteRecords {
			if len(logEntries) >= batchSize {
				break
			}
			outputJSON, err := MapLogEntry(v.AirbyteData)
			logEntry, err := source.MapJsonToLogEntry(outputJSON, v.AirbyteAbId, s.SourceId())
			if err != nil {
				logEntry = entity.LogEntryData{
					BaseData: entity.BaseData{
						SyncId: v.AirbyteAbId,
					},
				}
			}
			if sourceTableSuffix == ContentnoteTableSuffix && logEntry.ExternalSystem != "" && tenantSettings.SalesforceClientId != "" {
				linkedEntityAccountOrLeadId, err := fetchLinkedEntityForContentNote(ctx, s.tenant, tenantSettings.SalesforceClientId, tenantSettings.SalesforceClientSecret, tenantSettings.SalesforceRefreshToken, logEntry.ExternalId)
				if err != nil {
					s.log.Errorf("Failed to fetch linked entity for content note %s: %v", logEntry.ExternalId, err.Error())
					logEntry = entity.LogEntryData{
						BaseData: entity.BaseData{
							SyncId: v.AirbyteAbId,
						},
					}
				} else {
					s.log.Infof("Linked entity for content note %s: %v", logEntry.ExternalId, linkedEntityAccountOrLeadId)
					logEntry.LoggedOrganization = entity.ReferencedOrganization{
						ExternalId: linkedEntityAccountOrLeadId,
					}
				}
			}
			s.processingIds[v.AirbyteAbId] = source.ProcessingEntity{
				ExternalId:  logEntry.ExternalId,
				Entity:      currentEntity,
				TableSuffix: sourceTableSuffix,
			}
			logEntries = append(logEntries, logEntry)
		}
	}
	return logEntries
}

func (s *salesforceDataService) MarkProcessed(ctx context.Context, syncId, runId string, synced, skipped bool, reason string) error {
	v, ok := s.processingIds[syncId]
	if ok {
		err := repository.MarkAirbyteRawRecordProcessed(ctx, s.getDb(), v.Entity, v.TableSuffix, syncId, synced, skipped, runId, v.ExternalId, reason)
		if err != nil {
			s.log.Errorf("error while marking %s with external reference %s as synced for %s", v.Entity, v.ExternalId, s.SourceId())
		}
		return err
	}
	return nil
}

const (
	SalesforceOAuthURL   = "https://login.salesforce.com/services/oauth2/token"
	ServicesDataQueryURL = "/services/data/v57.0/query/?q=%s"
)

// Define a struct to hold token information
type TokenInfo struct {
	Token       *oauth2.Token
	InstanceURL string
}

// Define a map to store tokens with tenant as the key
var tokens = make(map[string]*TokenInfo)

func fetchLinkedEntityForContentNote(ctx context.Context, tenant, salesforceClientID, salesforceClientSecret, salesforceRefreshToken, contentNoteID string) (string, error) {

	// Step 1: Obtain an access token using the provided credentials
	token, err := getAccessTokenIfNeeded(tenant, salesforceClientID, salesforceClientSecret, salesforceRefreshToken)
	if err != nil {
		return "", err
	}

	// Retrieve the instance URL from the map
	instanceURL := tokens[tenant].InstanceURL

	// Build the Salesforce API URL for ContentDocumentLink
	query := fmt.Sprintf("SELECT+FIELDS(ALL)+from+ContentDocumentLink+where+contentDocumentId+=+'%s'+limit+200", contentNoteID)
	contentDocumentLinkURL := instanceURL + fmt.Sprintf(ServicesDataQueryURL, query)

	// Create a new HTTP request
	req, err := http.NewRequest("GET", contentDocumentLinkURL, nil)
	if err != nil {
		return "", err
	}
	// Set the authorization header with the access token
	req.Header.Set("Authorization", "Bearer "+token.AccessToken)

	// Send the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Failed to fetch linked entity: %s", resp.Status)
	}

	// Parse the response to extract the linked entity ID
	var sfContentDocumentLinkResponse struct {
		Records []struct {
			LinkedEntityId string `json:"LinkedEntityId"`
			IsDeleted      bool   `json:"IsDeleted"`
		} `json:"records"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&sfContentDocumentLinkResponse); err != nil {
		return "", err
	}

	for _, record := range sfContentDocumentLinkResponse.Records {
		if strings.HasPrefix(record.LinkedEntityId, sfAccountIdPrefix) || strings.HasPrefix(record.LinkedEntityId, sfLeadIdPrefix) {
			return record.LinkedEntityId, nil
		}
	}

	return "", nil
}

func getAccessTokenIfNeeded(tenant, ClientID, ClientSecret, RefreshToken string) (*oauth2.Token, error) {
	// Check if the token for the given tenant exists and is not expired
	if tokenInfo, ok := tokens[tenant]; ok {
		if !tokenInfo.Token.Valid() {
			// Token is expired, need to refresh
			return refreshAccessToken(tenant, ClientID, ClientSecret, RefreshToken)
		}
		return tokenInfo.Token, nil
	}

	// Token is missing, need to obtain a new one
	return obtainNewAccessToken(tenant, ClientID, ClientSecret, RefreshToken)
}

func refreshAccessToken(tenant, ClientID, ClientSecret, RefreshToken string) (*oauth2.Token, error) {
	// Configure the OAuth2 client credentials
	conf := &oauth2.Config{
		ClientID:     ClientID,
		ClientSecret: ClientSecret,
		Endpoint: oauth2.Endpoint{
			TokenURL: SalesforceOAuthURL,
		},
	}

	// Create a context
	ctx := context.Background()

	// Exchange the refresh token for a new access token
	token := &oauth2.Token{
		RefreshToken: RefreshToken,
	}

	newToken, err := conf.TokenSource(ctx, token).Token()
	if err != nil {
		return nil, err
	}

	// Update the token information in the map
	tokens[tenant] = &TokenInfo{
		Token:       newToken,
		InstanceURL: newToken.Extra("instance_url").(string),
	}

	return newToken, nil
}

func obtainNewAccessToken(tenant, ClientID, ClientSecret, RefreshToken string) (*oauth2.Token, error) {
	// Configure the OAuth2 client credentials
	conf := &oauth2.Config{
		ClientID:     ClientID,
		ClientSecret: ClientSecret,
		Endpoint: oauth2.Endpoint{
			TokenURL: SalesforceOAuthURL,
		},
	}

	// Create a context
	ctx := context.Background()

	// Exchange the refresh token for an access token
	token := &oauth2.Token{
		RefreshToken: RefreshToken,
	}

	newToken, err := conf.TokenSource(ctx, token).Token()
	if err != nil {
		return nil, err
	}

	// Update the token information in the map
	tokens[tenant] = &TokenInfo{
		Token:       newToken,
		InstanceURL: newToken.Extra("instance_url").(string),
	}

	return newToken, nil
}
