module github.com/openline-ai/openline-customer-os/packages/runner/customer-os-data-upkeeper

go 1.20

//replace github.com/openline-ai/openline-customer-os/packages/server/events-processing-common => ../../server/events-processing-common

require (
	github.com/caarlos0/env/v6 v6.10.1
	github.com/joho/godotenv v1.5.1
	github.com/neo4j/neo4j-go-driver/v5 v5.11.0
	github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module v0.0.0-20230828091946-278e377874e7
	github.com/openline-ai/openline-customer-os/packages/server/events-processing-common v0.0.0-20230828160335-6a557406abc8
	github.com/opentracing/opentracing-go v1.2.0
	github.com/robfig/cron v1.2.0
	github.com/stretchr/testify v1.8.4
	google.golang.org/grpc v1.57.0
	gorm.io/gorm v1.25.4
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackc/pgx/v5 v5.3.1 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/sirupsen/logrus v1.9.3 // indirect
	github.com/uber/jaeger-client-go v2.30.0+incompatible // indirect
	github.com/uber/jaeger-lib v2.4.1+incompatible // indirect
	go.uber.org/atomic v1.11.0 // indirect
	go.uber.org/multierr v1.10.0 // indirect
	go.uber.org/zap v1.25.0 // indirect
	golang.org/x/crypto v0.11.0 // indirect
	golang.org/x/net v0.12.0 // indirect
	golang.org/x/sys v0.11.0 // indirect
	golang.org/x/text v0.11.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20230525234030-28d5490b6b19 // indirect
	google.golang.org/protobuf v1.31.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	gorm.io/driver/postgres v1.5.2 // indirect
)
