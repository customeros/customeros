package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/openline-ai/openline-customer-os/packages/runner/integrity-checker/config"
	"github.com/openline-ai/openline-customer-os/packages/runner/integrity-checker/logger"
	"github.com/openline-ai/openline-customer-os/packages/runner/integrity-checker/model"
	"github.com/openline-ai/openline-customer-os/packages/runner/integrity-checker/repository"
	"github.com/openline-ai/openline-customer-os/packages/runner/integrity-checker/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
	"strings"
)

type Neo4jIntegrityCheckerService interface {
	RunIntegrityCheckerQueries()
}

type neo4jIntegrityCheckerService struct {
	cfg          *config.Config
	log          logger.Logger
	repositories *repository.Repositories
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
