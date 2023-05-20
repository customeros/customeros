package server

import (
	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/config"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/dataloader"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/generated"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/resolver"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/grpc_client"
	cosHandler "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/handler"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/service"
	commonService "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	"github.com/sirupsen/logrus"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"os"
	"os/signal"
	"syscall"
)

type server struct {
	cfg    *config.Config
	log    logger.Logger
	doneCh chan struct{}
}

func NewServer(cfg *config.Config, log logger.Logger) *server {
	return &server{cfg: cfg, log: log, doneCh: make(chan struct{})}
}

func (server *server) Run(parentCtx context.Context) error {
	ctx, cancel := signal.NotifyContext(parentCtx, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	config.InitLogger(server.cfg)

	// Initialize postgres db
	db, _ := InitDB(server.cfg)
	defer db.SqlDB.Close()

	// Setting up Neo4j
	neo4jDriver, err := config.NewDriver(server.cfg)
	if err != nil {
		server.log.Fatalf("Could not establish connection with neo4j at: %v, error: %v", server.cfg.Neo4j.Target, err.Error())
	}
	defer neo4jDriver.Close(ctx)

	// Setting up Postgres repositories
	commonServices := commonService.InitServices(db.GormDB, &neo4jDriver)

	// Setting up gRPC client
	var gRPCconn *grpc.ClientConn
	if server.cfg.Service.EventsProcessingPlatformEnabled {
		df := grpc_client.NewDialFactory(server.cfg)
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
	adminApiHandler := cosHandler.NewAdminApiHandler(server.cfg, commonServices)
	r.Use(cors.New(corsConfig))

	r.POST("/query",
		commonService.TenantUserContextEnhancer(ctx, commonService.USERNAME_OR_TENANT, commonServices.CommonRepositories),
		commonService.ApiKeyCheckerHTTP(commonServices.CommonRepositories.AppKeyRepository, commonService.CUSTOMER_OS_API),
		cosHandler.GetUserRoleHandlerEnhancer(),
		server.graphqlHandler(neo4jDriver, commonServices, gRPCconn))
	if server.cfg.GraphQL.PlaygroundEnabled {
		r.GET("/", playgroundHandler())
	}
	r.POST("/admin/query",
		adminApiHandler.GetAdminApiHandlerEnhancer(),
		server.graphqlHandler(neo4jDriver, commonServices, gRPCconn))

	if server.cfg.GraphQL.PlaygroundEnabled {
		r.GET("/admin/", playgroundAdminHandler())
	}

	r.GET("/health", healthCheckHandler)
	r.GET("/readiness", healthCheckHandler)

	r.Run(":" + server.cfg.ApiPort)

	<-server.doneCh
	server.log.Infof("Application %s exited properly", constants.ServiceName)
	return nil
}

func InitDB(cfg *config.Config) (db *config.StorageDB, err error) {
	if db, err = config.NewDBConn(cfg); err != nil {
		logrus.Fatalf("Coud not open db connection: %s", err.Error())
	}
	return
}

func (server *server) graphqlHandler(driver neo4j.DriverWithContext, commonServices *commonService.Services, gRPCconn *grpc.ClientConn) gin.HandlerFunc {
	grpcContainer := grpc_client.InitClients(gRPCconn)
	serviceContainer := service.InitServices(server.log, &driver, commonServices, grpcContainer)
	// instantiate graph resolver
	graphResolver := resolver.NewResolver(server.log, serviceContainer, grpcContainer)
	// make a data loader
	loader := dataloader.NewDataLoader(serviceContainer)
	schemaConfig := generated.Config{Resolvers: graphResolver}
	schemaConfig.Directives.HasRole = cosHandler.GetRoleChecker()
	schemaConfig.Directives.HasTenant = cosHandler.GetTenantChecker()

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(schemaConfig))
	srv.SetRecoverFunc(func(ctx context.Context, err interface{}) error {
		return gqlerror.Errorf("Internal server error!")
	})
	srv.SetErrorPresenter(func(ctx context.Context, e error) *gqlerror.Error {
		err := graphql.DefaultErrorPresenter(ctx, e)
		// Error hook place, Returned error can be customized. Check https://gqlgen.com/reference/errors/
		return err
	})
	srv.Use(extension.FixedComplexityLimit(server.cfg.GraphQL.FixedComplexityLimit))

	return func(c *gin.Context) {
		customCtx := &common.CustomContext{}
		if c.Keys["TenantName"] != nil {
			customCtx.Tenant = c.Keys["TenantName"].(string)
		}
		if c.Keys["Role"] != nil {
			customCtx.Role = c.Keys["Role"].(model.Role)
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

func playgroundAdminHandler() gin.HandlerFunc {
	h := playground.Handler("GraphQL", "/admin/query")

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func healthCheckHandler(c *gin.Context) {
	c.JSON(200, gin.H{"status": "OK"})
}
