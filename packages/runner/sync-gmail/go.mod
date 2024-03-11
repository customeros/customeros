module github.com/openline-ai/openline-customer-os/packages/runner/sync-gmail

go 1.21

//replace github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module => ./../../server/customer-os-common-module

require (
	github.com/caarlos0/env/v6 v6.10.1
	github.com/google/uuid v1.6.0
	github.com/joho/godotenv v1.5.1
	github.com/neo4j/neo4j-go-driver/v5 v5.18.0
	github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module v0.0.0-20240220064825-20118a9bac6a
	github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto v0.0.0-20240220064825-20118a9bac6a
	github.com/opentracing/opentracing-go v1.2.0
	github.com/pkg/errors v0.9.1
	github.com/robfig/cron v1.2.0
	github.com/sirupsen/logrus v1.9.3
	google.golang.org/grpc v1.62.1
	gorm.io/driver/postgres v1.5.7
	gorm.io/gorm v1.25.7-0.20240204074919-46816ad31dde
)

require (
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/h2non/filetype v1.1.3 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackc/pgx/v5 v5.5.0 // indirect
	github.com/jackc/puddle/v2 v2.2.1 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/stretchr/objx v0.5.0 // indirect
	github.com/uber/jaeger-client-go v2.30.0+incompatible // indirect
	github.com/uber/jaeger-lib v2.4.1+incompatible // indirect
	go.uber.org/atomic v1.11.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.26.0 // indirect
	golang.org/x/crypto v0.18.0 // indirect
	golang.org/x/net v0.20.0 // indirect
	golang.org/x/sync v0.6.0 // indirect
	golang.org/x/sys v0.16.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240123012728-ef4313101c80 // indirect
	google.golang.org/protobuf v1.32.0 // indirect
)
