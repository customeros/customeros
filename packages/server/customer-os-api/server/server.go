package server

import (
	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/config"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/dataloader"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/generated"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/resolver"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/grpc_client"
	cosHandler "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/handler"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/mapper"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/rest"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/service"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/validator"
	commonService "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
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

	if err := validator.GetValidator().Struct(server.cfg); err != nil {
		return errors.Wrap(err, "cfg validate")
	}

	config.InitLogrusLogger(server.cfg)

	// Setting up tracing
	if server.cfg.Jaeger.Enabled {
		tracer, closer, err := tracing.NewJaegerTracer(&server.cfg.Jaeger, server.log)
		if err != nil {
			server.log.Fatalf("Could not initialize jaeger tracer: %s", err.Error())
		}
		defer closer.Close()
		opentracing.SetGlobalTracer(tracer)
	}

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
	grpcContainer := grpc_client.InitClients(gRPCconn)

	serviceContainer := service.InitServices(server.log, &neo4jDriver, commonServices, grpcContainer)
	r.Use(cors.New(corsConfig))

	r.POST("/query",
		cosHandler.TracingEnhancer(ctx, "/query"),
		commonService.TenantUserContextEnhancer(ctx, commonService.USERNAME_OR_TENANT, commonServices.CommonRepositories),
		commonService.ApiKeyCheckerHTTP(commonServices.CommonRepositories.AppKeyRepository, commonService.CUSTOMER_OS_API),
		server.graphqlHandler(grpcContainer, serviceContainer))
	if server.cfg.GraphQL.PlaygroundEnabled {
		r.GET("/",
			cosHandler.TracingEnhancer(ctx, "/"),
			playgroundHandler())
	}
	r.GET("/whoami",
		cosHandler.TracingEnhancer(ctx, "/whoami"),
		commonService.ApiKeyCheckerHTTP(commonServices.CommonRepositories.AppKeyRepository, commonService.CUSTOMER_OS_API),
		rest.WhoamiHandler(serviceContainer))
	r.POST("/admin/query",
		cosHandler.TracingEnhancer(ctx, "/admin/query"),
		adminApiHandler.GetAdminApiHandlerEnhancer(),
		server.graphqlHandler(grpcContainer, serviceContainer))

	if server.cfg.GraphQL.PlaygroundEnabled {
		r.GET("/admin/",
			cosHandler.TracingEnhancer(ctx, "/admin"),
			playgroundAdminHandler())
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

func (server *server) graphqlHandler(grpcContainer *grpc_client.Clients, serviceContainer *service.Services) gin.HandlerFunc {
	// instantiate graph resolver
	graphResolver := resolver.NewResolver(server.log, serviceContainer, grpcContainer)
	// make a data loader
	loader := dataloader.NewDataLoader(serviceContainer)
	schemaConfig := generated.Config{Resolvers: graphResolver}
	schemaConfig.Directives.HasRole = cosHandler.GetRoleChecker()
	schemaConfig.Directives.HasTenant = cosHandler.GetTenantChecker()
	schemaConfig.Directives.HasIdentityId = cosHandler.GetIdentityIdChecker()

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
		if c.Keys[commonService.KEY_TENANT_NAME] != nil {
			customCtx.Tenant = c.Keys[commonService.KEY_TENANT_NAME].(string)
		}
		if c.Keys[commonService.KEY_USER_ROLES] != nil {
			customCtx.Roles = mapper.MapRolesToModel(c.Keys[commonService.KEY_USER_ROLES].([]string))
		}
		if c.Keys[commonService.KEY_USER_ID] != nil {
			customCtx.UserId = c.Keys[commonService.KEY_USER_ID].(string)
		}
		if c.Keys[commonService.KEY_IDENTITY_ID] != nil {
			customCtx.IdentityId = c.Keys[commonService.KEY_IDENTITY_ID].(string)
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
