package server

import (
	"context"
	"fmt"
	"github.com/AleksK1NG/es-microservice/config"
	"github.com/AleksK1NG/es-microservice/pkg/constants"
	"github.com/AleksK1NG/es-microservice/pkg/elasticsearch"
	serviceErrors "github.com/AleksK1NG/es-microservice/pkg/service_errors"
	"github.com/AleksK1NG/es-microservice/pkg/utils"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strings"
	"time"
)

const (
	waitShotDownDuration = 3 * time.Second
)

func (server *server) initMongoDBCollections(ctx context.Context) {
	err := server.mongoClient.Database(server.cfg.Mongo.Db).CreateCollection(ctx, server.cfg.MongoCollections.Orders)
	if err != nil {
		if !utils.CheckErrMessages(err, serviceErrors.ErrMsgMongoCollectionAlreadyExists) {
			server.log.Warnf("(CreateCollection) err: {%v}", err)
		}
	}

	indexOptions := options.Index().SetSparse(true).SetUnique(true)
	index, err := server.mongoClient.Database(server.cfg.Mongo.Db).Collection(server.cfg.MongoCollections.Orders).Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: constants.OrderIdIndex, Value: 1}},
		Options: indexOptions,
	})
	if err != nil && !utils.CheckErrMessages(err, serviceErrors.ErrMsgAlreadyExists) {
		server.log.Warnf("(CreateOne) err: {%v}", err)
	}
	server.log.Infof("(CreatedIndex) index: {%server}", index)

	list, err := server.mongoClient.Database(server.cfg.Mongo.Db).Collection(server.cfg.MongoCollections.Orders).Indexes().List(ctx)
	if err != nil {
		server.log.Warnf("(initMongoDBCollections) [List] err: {%v}", err)
	}

	if list != nil {
		var results []bson.M
		if err := list.All(ctx, &results); err != nil {
			server.log.Warnf("(All) err: {%v}", err)
		}
		server.log.Infof("(indexes) results: {%#v}", results)
	}

	collections, err := server.mongoClient.Database(server.cfg.Mongo.Db).ListCollectionNames(ctx, bson.M{})
	if err != nil {
		server.log.Warnf("(ListCollections) err: {%v}", err)
	}
	server.log.Infof("(Collections) created collections: {%v}", collections)
}

func (server *server) initElasticClient(ctx context.Context) error {
	elasticClient, err := elasticsearch.NewElasticClient(server.cfg.Elastic)
	if err != nil {
		return err
	}
	server.elasticClient = elasticClient

	info, code, err := server.elasticClient.Ping(server.cfg.Elastic.URL).Do(ctx)
	if err != nil {
		return errors.Wrap(err, "client.Ping")
	}
	server.log.Infof("Elasticsearch returned with code {%d} and version {%server}", code, info.Version.Number)

	esVersion, err := server.elasticClient.ElasticsearchVersion(server.cfg.Elastic.URL)
	if err != nil {
		return errors.Wrap(err, "client.ElasticsearchVersion")
	}
	server.log.Infof("Elasticsearch version {%server}", esVersion)

	return nil
}

func (server *server) runMetrics(cancel context.CancelFunc) {
	metricsServer := echo.New()
	go func() {
		metricsServer.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
			StackSize:         stackSize,
			DisablePrintStack: true,
			DisableStackAll:   true,
		}))
		metricsServer.GET(server.cfg.Probes.PrometheusPath, echo.WrapHandler(promhttp.Handler()))
		server.log.Infof("Metrics server is running on port: {%server}", server.cfg.Probes.PrometheusPort)
		if err := metricsServer.Start(server.cfg.Probes.PrometheusPort); err != nil {
			server.log.Errorf("metricsServer.Start: {%v}", err)
			cancel()
		}
	}()
}

func (server *server) getHttpMetricsCb() func(err error) {
	return func(err error) {
		if err != nil {
			server.metrics.ErrorHttpRequests.Inc()
		} else {
			server.metrics.SuccessHttpRequests.Inc()
		}
	}
}

func (server *server) getGrpcMetricsCb() func(err error) {
	return func(err error) {
		if err != nil {
			server.metrics.ErrorGrpcRequests.Inc()
		} else {
			server.metrics.SuccessGrpcRequests.Inc()
		}
	}
}

func (server *server) waitShootDown(duration time.Duration) {
	go func() {
		time.Sleep(duration)
		server.doneCh <- struct{}{}
	}()
}

func GetMicroserviceName(cfg *config.Config) string {
	return fmt.Sprintf("(%s)", strings.ToUpper(cfg.ServiceName))
}
