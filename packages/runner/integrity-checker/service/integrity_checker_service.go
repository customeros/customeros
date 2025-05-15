package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/customeros/customeros/packages/runner/integrity-checker/caches"
	"github.com/customeros/customeros/packages/runner/integrity-checker/config"
	"github.com/customeros/customeros/packages/runner/integrity-checker/logger"
	"github.com/customeros/customeros/packages/runner/integrity-checker/model"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	neo4jRepository "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/repository"
	postgresRepository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"
	"github.com/pkg/errors"
)

type DbType string

const (
	Neo4j    DbType = "Neo4j"
	Postgres DbType = "Postgres"
)

type IntegrityCheckerService interface {
	RunNeo4jIntegrityCheckerQueries()
	RunOpenlinePostgresIntegrityCheckerQueries()
}

type integrityCheckerService struct {
	cfg      *config.Config
	log      logger.Logger
	cache    *caches.Cache
	neo4j    *neo4jRepository.Repositories
	postgres *postgresRepository.Repositories
}

type integrityCheckerResult struct {
	Name                 string `json:"name"`
	Success              bool   `json:"success"`
	CountOfDataWithIssue int64  `json:"countOfDataWithIssue"`
	TechError            string `json:"techError"`
}

func (i integrityCheckerResult) String() string {
	return fmt.Sprintf("(name: %s, success: %t, countOfDataWithIssue: %d, techError: %s)",
		i.Name, i.Success, i.CountOfDataWithIssue, i.TechError)
}

func NewIntegrityCheckerService(cfg *config.Config, log logger.Logger, neo4j *neo4jRepository.Repositories, postgres *postgresRepository.Repositories, cache *caches.Cache) IntegrityCheckerService {
	return &integrityCheckerService{
		cfg:      cfg,
		log:      log,
		neo4j:    neo4j,
		postgres: postgres,
		cache:    cache,
	}
}

func (s *integrityCheckerService) RunNeo4jIntegrityCheckerQueries() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Cancel context on exit

	spans, ctx := telemetry.StartSpan(ctx, "IntegrityCheckerService.RunNeo4jIntegrityCheckerQueries")
	defer spans.Finish()

	integrityCheckerQueries, err := s.getQueriesFromS3(ctx, "neo4j-integrity-checker-queries.json")
	if err != nil {
		spans.TraceError(err)
		s.log.Errorf("Error getting queries from S3: %v", err)
	}
	result := s.executeNeo4jQueries(ctx, integrityCheckerQueries)
	spans.LogObjectAsJson("integrityCheckerResult", result)
	s.log.Infof("Neo4j integrity checker result: %v", result)

	_ = s.alertInSlack(ctx, result, Neo4j)
}

func (s *integrityCheckerService) RunOpenlinePostgresIntegrityCheckerQueries() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Cancel context on exit

	spans, ctx := telemetry.StartSpan(ctx, "IntegrityCheckerService.RunOpenlinePostgresIntegrityCheckerQueries")
	defer spans.Finish()

	integrityCheckerQueries, err := s.getQueriesFromS3(ctx, "postgres-openline-integrity-checker-queries.json")
	if err != nil {
		spans.TraceError(err)
		s.log.Errorf("Error getting queries from S3: %v", err)
	}
	result := s.executePostgresQueries(ctx, integrityCheckerQueries)
	spans.LogObjectAsJson("integrityCheckerResult", result)
	s.log.Infof("Postgres (openline DB) integrity checker result: %v", result)

	_ = s.alertInSlack(ctx, result, Postgres)
}

func (s *integrityCheckerService) getQueriesFromS3(ctx context.Context, filename string) (model.IntegrityCheckQueries, error) {
	spans, ctx := telemetry.StartSpan(ctx, "IntegrityCheckerService.getQueriesFromS3")
	defer spans.Finish()

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
		Config: aws.Config{
			Region: aws.String(s.cfg.AWS.Region),
		},
	}))
	downloader := s3manager.NewDownloader(sess)

	buffer := &aws.WriteAtBuffer{}
	_, err := downloader.Download(buffer,
		&s3.GetObjectInput{
			Bucket: aws.String(s.cfg.AWS.Bucket),
			Key:    aws.String(filename),
		})
	if err != nil {
		spans.TraceError(err)
		s.log.Errorf("Error downloading queries from S3: %v", err)
		return model.IntegrityCheckQueries{}, err
	}

	var queries model.IntegrityCheckQueries
	if err := json.Unmarshal(buffer.Bytes(), &queries); err != nil {
		spans.TraceError(err)
		s.log.Errorf("Error unmarshalling queries: %v", err)
		return model.IntegrityCheckQueries{}, err
	}

	return queries, nil
}

func (s *integrityCheckerService) executeNeo4jQueries(ctx context.Context, queries model.IntegrityCheckQueries) []integrityCheckerResult {
	spans, ctx := telemetry.StartSpan(ctx, "IntegrityCheckerService.executeNeo4jQueries")
	defer spans.Finish()

	var output []integrityCheckerResult
	var queriesToExecute []model.Query
	for _, query := range queries.Queries {
		queriesToExecute = append(queriesToExecute, query)
	}
	for _, group := range queries.Groups {
		for _, query := range group.Queries {
			queriesToExecute = append(queriesToExecute, query)
		}
	}

	for _, query := range queriesToExecute {
		select {
		case <-ctx.Done():
			return output
		default:
			count, err := s.neo4j.CommonReadRepository.ExecuteIntegrityCheckerQuery(ctx, query.Name, query.Query)
			checkerResult := integrityCheckerResult{
				Name:                 query.Name,
				Success:              err == nil && count == int64(0),
				CountOfDataWithIssue: count,
			}
			if err != nil {
				checkerResult.TechError = err.Error()
			}
			output = append(output, checkerResult)
		}
	}

	return output
}

func (s *integrityCheckerService) executePostgresQueries(ctx context.Context, queries model.IntegrityCheckQueries) []integrityCheckerResult {
	spans, ctx := telemetry.StartSpan(ctx, "IntegrityCheckerService.executePostgresQueries")
	defer spans.Finish()

	var output []integrityCheckerResult
	var queriesToExecute []model.Query
	for _, query := range queries.Queries {
		queriesToExecute = append(queriesToExecute, query)
	}
	for _, group := range queries.Groups {
		for _, query := range group.Queries {
			queriesToExecute = append(queriesToExecute, query)
		}
	}

	for _, query := range queriesToExecute {
		select {
		case <-ctx.Done():
			return output
		default:
			count, err := s.postgres.CommonRepository.GetCountFromPlainQuery(ctx, query.Query)
			checkerResult := integrityCheckerResult{
				Name:                 query.Name,
				Success:              err == nil && count == int64(0),
				CountOfDataWithIssue: count,
			}
			if err != nil {
				checkerResult.TechError = err.Error()
			}
			output = append(output, checkerResult)
		}
	}

	return output
}

func (s *integrityCheckerService) alertInSlack(ctx context.Context, results []integrityCheckerResult, dbType DbType) error {
	spans, ctx := telemetry.StartSpan(ctx, "IntegrityCheckerService.alertInSlack")
	defer spans.Finish()

	// if no webhook is configured, return early
	if s.cfg.SlackConfig.DataAlertsRegisteredWebhook == "" {
		spans.TraceError(errors.New("no slack webhook configured"))
		return nil
	}

	spans.LogKV("slackUrl", s.cfg.SlackConfig.DataAlertsRegisteredWebhook)

	var alertMessages []string
	hasAlert := false
	var issues []struct {
		message string
		count   int
	}

	for _, result := range results {
		if result.Success && result.CountOfDataWithIssue == 0 && result.TechError == "" {
			continue
		}

		if result.CountOfDataWithIssue > 0 {
			issues = append(issues, struct {
				message string
				count   int
			}{message: result.Name, count: int(result.CountOfDataWithIssue)})
			hasAlert = true
		}
	}

	// sort issues by count descending
	sort.Slice(issues, func(i, j int) bool {
		return issues[i].count > issues[j].count
	})

	for _, issue := range issues {
		alertMessages = append(alertMessages, fmt.Sprintf("%s: %d", issue.message, issue.count))
	}

	// do not send messages to slack if no changes from previous run
	var previousAlertMessages []string
	var err error
	if dbType == Neo4j {
		previousAlertMessages, err = s.cache.GetPreviousNeo4jAlertMessages()
	} else {
		previousAlertMessages, err = s.cache.GetPreviousPostgresAlertMessages()
	}
	if err != nil {
		spans.TraceError(errors.Wrap(err, "error getting previous alert messages"))
	}
	if utils.StringSlicesEqualIgnoreOrder(previousAlertMessages, alertMessages) {
		spans.LogKV("result", "no changes from previous run")
		return nil
	}

	if dbType == Neo4j {
		err = s.cache.SetPreviousNeo4jAlertMessages(alertMessages)
	} else {
		err = s.cache.SetPreviousPostgresAlertMessages(alertMessages)
	}
	if err != nil {
		spans.TraceError(errors.Wrap(err, "error setting previous alert messages"))
	}

	// If no alerts, return early
	if !hasAlert {
		spans.LogKV("result", "no alerts to send")
		return nil
	}

	// Create the message text
	messageText := fmt.Sprintf("%s Data Integrity Issues Summary:\n%s", dbType, strings.Join(alertMessages, "\n"))

	// Create a struct to hold the JSON data
	type SlackMessage struct {
		Text string `json:"text"`
	}
	message := SlackMessage{Text: messageText}

	// Convert struct to JSON
	jsonData, err := json.Marshal(message)
	if err != nil {
		spans.TraceError(err)
		return fmt.Errorf("error encoding JSON: %w", err)
	}

	spans.LogObjectAsJson("request", string(jsonData))

	// Send POST request
	resp, err := http.Post(s.cfg.SlackConfig.DataAlertsRegisteredWebhook, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		spans.TraceError(err)
		return fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	spans.LogKV("result.status", resp.Status)

	// Check response status
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		spans.TraceError(errors.New("unexpected status code: " + resp.Status))
		spans.LogKV("response.body", string(body))
		return fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	spans.LogKV("result", "alert sent successfully")
	return nil
}
