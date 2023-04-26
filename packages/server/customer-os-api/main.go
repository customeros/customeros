package main

import (
	"context"
	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/caarlos0/env/v6"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/config"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/dataloader"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/generated"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/resolver"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/grpc_client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/service"
	commonService "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	"github.com/sirupsen/logrus"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"google.golang.org/grpc"
)

const customerOSApiPort = "10000"

func InitDB(cfg *config.Config) (db *config.StorageDB, err error) {
	if db, err = config.NewDBConn(cfg); err != nil {
		logrus.Fatalf("Coud not open db connection: %s", err.Error())
	}
	return
}

func graphqlHandler(cfg *config.Config, driver neo4j.DriverWithContext,
	commonServices *commonService.Services, gRPCconn *grpc.ClientConn) gin.HandlerFunc {

	grpcContainer := grpc_client.InitClients(gRPCconn)
	serviceContainer := service.InitServices(&driver, commonServices, grpcContainer)
	// instantiate graph resolver
	graphResolver := resolver.NewResolver(serviceContainer, grpcContainer)
	// make a data loader
	loader := dataloader.NewDataLoader(serviceContainer)
	schemaConfig := generated.Config{Resolvers: graphResolver}

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(schemaConfig))
	srv.SetRecoverFunc(func(ctx context.Context, err interface{}) error {
		return gqlerror.Errorf("Internal server error!")
	})
	srv.SetErrorPresenter(func(ctx context.Context, e error) *gqlerror.Error {
		err := graphql.DefaultErrorPresenter(ctx, e)
		// Error hook place, Returned error can be customized. Check https://gqlgen.com/reference/errors/
		return err
	})
	srv.Use(extension.FixedComplexityLimit(cfg.GraphQL.FixedComplexityLimit))

	return func(c *gin.Context) {
		customCtx := &common.CustomContext{
			Tenant: c.Keys["TenantName"].(string),
		}
		if c.Keys["UserId"] != nil {
			customCtx.UserId = c.Keys["UserId"].(string)
		}

		dataloaderSrv := dataloader.Middleware(loader, srv)
		h := common.WithContext(customCtx, dataloaderSrv)

		h.ServeHTTP(c.Writer, c.Request)
	}
}

// Defining the Playground handler
func playgroundHandler() gin.HandlerFunc {
	h := playground.Handler("GraphQL", "/query")

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func init() {
}

func main() {
	cfg := loadConfiguration()
	config.InitLogger(cfg)

	db, _ := InitDB(cfg)
	defer db.SqlDB.Close()

	// Setting up Neo4j
	neo4jDriver, err := config.NewDriver(cfg)
	if err != nil {
		logrus.Fatalf("Could not establish connection with neo4j at: %v, error: %v", cfg.Neo4j.Target, err.Error())
	}
	ctx := context.Background()
	defer neo4jDriver.Close(ctx)

	// Setting up Postgres repositories
	commonServices := commonService.InitServices(db.GormDB, &neo4jDriver)

	// Setting up gRPC client
	var gRPCconn *grpc.ClientConn
	if cfg.Service.EventsProcessingPlatformEnabled {
		df := grpc_client.NewDialFactory(cfg)
		gRPCconn, err = df.GetEventsProcessingPlatformConn()
		if err != nil {
			logrus.Fatalf("Failed to connect: %v", err)
		}
		defer df.Close(gRPCconn)
	}

	// Setting up Gin
	r := gin.Default()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"*"}
	r.Use(cors.New(corsConfig))

	r.POST("/query",
		commonService.TenantUserContextEnhancer(ctx, commonService.USERNAME_OR_TENANT, commonServices.CommonRepositories),
		commonService.ApiKeyCheckerHTTP(commonServices.CommonRepositories.AppKeyRepository, commonService.CUSTOMER_OS_API),
		graphqlHandler(cfg, neo4jDriver, commonServices, gRPCconn))
	if cfg.GraphQL.PlaygroundEnabled {
		r.GET("/", playgroundHandler())
	}
	r.GET("/health", healthCheckHandler)
	r.GET("/readiness", healthCheckHandler)

	port := cfg.ApiPort
	if port == "" {
		port = customerOSApiPort
	}

	r.Run(":" + port)
}

func loadConfiguration() *config.Config {
	if err := godotenv.Load(); err != nil {
		logrus.Warn("Error loading .env file")
	}

	cfg := config.Config{}
	if err := env.Parse(&cfg); err != nil {
		logrus.Errorf("%+v\n", err)
	}

	return &cfg
}

func healthCheckHandler(c *gin.Context) {

	c.JSON(200, gin.H{"status": "OK"})
}
