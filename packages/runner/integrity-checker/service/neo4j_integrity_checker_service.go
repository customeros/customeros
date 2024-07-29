package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/openline-ai/openline-customer-os/packages/runner/integrity-checker/caches"
	"github.com/openline-ai/openline-customer-os/packages/runner/integrity-checker/config"
	"github.com/openline-ai/openline-customer-os/packages/runner/integrity-checker/logger"
	"github.com/openline-ai/openline-customer-os/packages/runner/integrity-checker/model"
	"github.com/openline-ai/openline-customer-os/packages/runner/integrity-checker/repository"
	"github.com/openline-ai/openline-customer-os/packages/runner/integrity-checker/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"io"
	"net/http"
	"sort"
	"strings"
)

type Neo4jIntegrityCheckerService interface {
	RunIntegrityCheckerQueries()
}

type neo4jIntegrityCheckerService struct {
	cfg          *config.Config
	log          logger.Logger
	repositories *repository.Repositories
	cache        *caches.Cache
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

func NewNeo4jIntegrityCheckerService(cfg *config.Config, log logger.Logger, repositories *repository.Repositories) Neo4jIntegrityCheckerService {
	return &neo4jIntegrityCheckerService{
		cfg:          cfg,
		log:          log,
		repositories: repositories,
		cache:        caches.NewCache(),
	}
}

func (s *neo4jIntegrityCheckerService) RunIntegrityCheckerQueries() {
	ctx, cancel := context.WithCancel(context.Background())

	span, ctx := tracing.StartTracerSpan(ctx, "Neo4jIntegrityCheckerService.RunIntegrityCheckerQueries")
	defer span.Finish()

	defer cancel() // Cancel context on exit

	integrityCheckerQueries, err := s.getQueriesFromS3(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error getting queries from S3: %v", err)
	}
	result := s.executeQueries(ctx, integrityCheckerQueries)
	tracing.LogObjectAsJson(span, "integrityCheckerResult", result)
	s.log.Infof("Integrity checker result: %v", result)

	s.sendMetrics(ctx, result)
	_ = s.alertInSlack(ctx, result)
}

func (s *neo4jIntegrityCheckerService) getQueriesFromS3(ctx context.Context) (model.IntegrityCheckQueries, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Neo4jIntegrityCheckerService.getQueriesFromS3")
	defer span.Finish()

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
			Key:    aws.String("neo4j-integrity-checker-queries.json"),
		})
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error downloading queries from S3: %v", err)
		return model.IntegrityCheckQueries{}, err
	}

	var queries model.IntegrityCheckQueries
	if err := json.Unmarshal(buffer.Bytes(), &queries); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error unmarshalling queries: %v", err)
		return model.IntegrityCheckQueries{}, err
	}

	return queries, nil
}

func (s *neo4jIntegrityCheckerService) executeQueries(ctx context.Context, queries model.IntegrityCheckQueries) []integrityCheckerResult {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Neo4jIntegrityCheckerService.RunIntegrityCheckerQueries")
	defer span.Finish()

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
		// Check if context is cancelled
		select {
		case <-ctx.Done():
			return output
		default:
			// Continue fetching organizations
		}

		count, err := s.repositories.Neo4jRepositories.CommonReadRepository.ExecuteIntegrityCheckerQuery(ctx, query.Name, query.Query)
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

	return output
}

func (s *neo4jIntegrityCheckerService) sendMetrics(ctx context.Context, results []integrityCheckerResult) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Neo4jIntegrityCheckerService.sendMetrics")
	defer span.Finish()

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
		Config: aws.Config{
			Region: aws.String(s.cfg.AWS.Region),
		},
	}))

	svc := cloudwatch.New(sess)

	var metrics []*cloudwatch.MetricDatum
	totalProblematicNodes := int64(0)
	totalFailedQueries := 0

	dimensions := []*cloudwatch.Dimension{{
		Name:  aws.String("Environment"),
		Value: aws.String(s.cfg.AWS.MetricsDimensionEnvironment),
	}, {
		Name:  aws.String("Service"),
		Value: aws.String(s.cfg.AWS.MetricsDimensionNeo4jIntegrityChecks),
	}}

	executionTime := utils.Now()

	for _, result := range results {
		if result.Success && result.CountOfDataWithIssue == 0 && result.TechError == "" {
			continue
		}
		metrics = append(metrics, &cloudwatch.MetricDatum{
			MetricName: aws.String(strings.ReplaceAll(strings.ToLower(result.Name), " ", "_")),
			Value:      aws.Float64(float64(result.CountOfDataWithIssue)),
			Unit:       aws.String("Count"),
			Timestamp:  utils.TimePtr(executionTime),
			Dimensions: dimensions,
		})
		if result.TechError != "" {
			totalFailedQueries++
		}
		totalProblematicNodes += result.CountOfDataWithIssue
	}

	metrics = append(metrics, &cloudwatch.MetricDatum{
		MetricName: aws.String("failed_neo4j_queries"),
		Value:      aws.Float64(float64(totalFailedQueries)),
		Unit:       aws.String("Count"),
		Timestamp:  utils.TimePtr(utils.Now()),
		Dimensions: dimensions,
	})

	metrics = append(metrics, &cloudwatch.MetricDatum{
		MetricName: aws.String("neo4j_integrity_checker_data_issues"),
		Value:      aws.Float64(float64(totalProblematicNodes)),
		Unit:       aws.String("Count"),
		Timestamp:  utils.TimePtr(utils.Now()),
		Dimensions: dimensions,
	})

	_, err := svc.PutMetricData(&cloudwatch.PutMetricDataInput{
		Namespace:  aws.String(s.cfg.AWS.CloudWatchNamespace),
		MetricData: metrics,
	})

	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error reporting metrics: %v", err)
		return
	}
}

func (h *neo4jIntegrityCheckerService) alertInSlack(ctx context.Context, results []integrityCheckerResult) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Neo4jIntegrityCheckerService.alertInSlack")
	defer span.Finish()

	// if no webhook is configured, return early
	if h.cfg.SlackConfig.DataAlertsRegisteredWebhook == "" {
		return nil
	}

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
	previousAlertMessages, err := h.cache.GetPreviousAlertMessages()
	if err != nil {
		tracing.TraceErr(span, err)
	}
	if utils.StringSlicesEqualIgnoreOrder(previousAlertMessages, alertMessages) {
		return nil
	}

	err = h.cache.SetPreviousAlertMessages(alertMessages)
	if err != nil {
		tracing.TraceErr(span, err)
	}

	// If no alerts, return early
	if !hasAlert {
		return nil
	}

	// Create the message text
	messageText := "Data Integrity Issues Summary:\n" + strings.Join(alertMessages, "\n")

	// Create a struct to hold the JSON data
	type SlackMessage struct {
		Text string `json:"text"`
	}
	message := SlackMessage{Text: messageText}

	// Convert struct to JSON
	jsonData, err := json.Marshal(message)
	if err != nil {
		span.LogFields(log.Error(err))
		return fmt.Errorf("error encoding JSON: %w", err)
	}

	// Send POST request
	resp, err := http.Post(h.cfg.SlackConfig.DataAlertsRegisteredWebhook, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		span.LogFields(log.Error(err))
		return fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	span.LogFields(log.String("result.status", resp.Status))

	// Check response status
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	return nil
}
