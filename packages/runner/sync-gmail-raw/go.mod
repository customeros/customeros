module github.com/openline-ai/openline-customer-os/packages/runner/sync-gmail-raw

go 1.21

replace github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto => ../../server/events-processing-proto

replace github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module => ./../../server/customer-os-common-module

replace github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository => ./../../server/customer-os-neo4j-repository

replace github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository => ./../../server/customer-os-postgres-repository

replace github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-auth => ./../../server/customer-os-common-auth

require (
	github.com/araddon/dateparse v0.0.0-20210429162001-6b43995a97de
	github.com/caarlos0/env/v6 v6.10.1
	github.com/google/uuid v1.6.0
	github.com/joho/godotenv v1.5.1
	github.com/neo4j/neo4j-go-driver/v5 v5.20.0
	github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-auth v0.0.0-20240219173641-99951e96f1dc
	github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module v0.0.0-20240219173641-99951e96f1dc
	github.com/pkg/errors v0.9.1
	github.com/robfig/cron v1.2.0
	github.com/sirupsen/logrus v1.9.3
	google.golang.org/api v0.165.0
	gorm.io/driver/postgres v1.5.7
	gorm.io/gorm v1.25.10
)

require (
	cloud.google.com/go/compute v1.24.0 // indirect
	cloud.google.com/go/compute/metadata v0.2.3 // indirect
	github.com/cenkalti/backoff/v4 v4.3.0 // indirect
	github.com/felixge/httpsnoop v1.0.4 // indirect
	github.com/go-logr/logr v1.4.1 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/google/s2a-go v0.1.7 // indirect
	github.com/googleapis/enterprise-certificate-proxy v0.3.2 // indirect
	github.com/googleapis/gax-go/v2 v2.12.1 // indirect
	github.com/h2non/filetype v1.1.3 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20231201235250-de7065d80cb9 // indirect
	github.com/jackc/pgx/v5 v5.5.4 // indirect
	github.com/jackc/puddle/v2 v2.2.1 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository v0.0.0-20240410144729-44cbe53c019c // indirect
	github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository v0.0.0-20240410144729-44cbe53c019c // indirect
	github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto v0.0.0-20240413132139-bfffc416fdeb // indirect
	github.com/opentracing/opentracing-go v1.2.0 // indirect
	github.com/uber/jaeger-client-go v2.30.0+incompatible // indirect
	github.com/uber/jaeger-lib v2.4.1+incompatible // indirect
	go.opencensus.io v0.24.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.49.0 // indirect
	go.opentelemetry.io/otel v1.24.0 // indirect
	go.opentelemetry.io/otel/metric v1.24.0 // indirect
	go.opentelemetry.io/otel/trace v1.24.0 // indirect
	go.uber.org/atomic v1.11.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.27.0 // indirect
	golang.org/x/crypto v0.22.0 // indirect
	golang.org/x/exp v0.0.0-20240205201215-2c58cdc269a3 // indirect
	golang.org/x/net v0.24.0 // indirect
	golang.org/x/oauth2 v0.17.0 // indirect
	golang.org/x/sync v0.6.0 // indirect
	golang.org/x/sys v0.19.0 // indirect
	golang.org/x/text v0.15.0 // indirect
	google.golang.org/appengine v1.6.8 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240227224415-6ceb2ff114de // indirect
	google.golang.org/grpc v1.63.2 // indirect
	google.golang.org/protobuf v1.34.0 // indirect
)
