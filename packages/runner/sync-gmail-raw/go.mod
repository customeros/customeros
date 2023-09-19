module github.com/openline-ai/openline-customer-os/packages/runner/sync-gmail-raw

go 1.20

//replace github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module => ./../../server/customer-os-common-module
//replace github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-auth => ./../../server/customer-os-common-auth

require (
	github.com/caarlos0/env/v6 v6.10.1
	github.com/google/uuid v1.3.0
	github.com/joho/godotenv v1.5.1
	github.com/neo4j/neo4j-go-driver/v5 v5.12.0
	github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-auth v0.0.0-20230919071321-b1ec8acfb667
	github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module v0.0.0-20230905205324-bb882f6be687
	github.com/pkg/errors v0.9.1
	github.com/robfig/cron v1.2.0
	github.com/sirupsen/logrus v1.9.3
	golang.org/x/net v0.15.0
	golang.org/x/oauth2 v0.12.0
	google.golang.org/api v0.134.0
	gorm.io/driver/postgres v1.5.2
	gorm.io/gorm v1.25.4
)

require (
	cloud.google.com/go/compute v1.20.1 // indirect
	cloud.google.com/go/compute/metadata v0.2.3 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/google/s2a-go v0.1.4 // indirect
	github.com/googleapis/enterprise-certificate-proxy v0.2.5 // indirect
	github.com/googleapis/gax-go/v2 v2.12.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackc/pgx/v5 v5.3.1 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	go.opencensus.io v0.24.0 // indirect
	go.uber.org/multierr v1.10.0 // indirect
	go.uber.org/zap v1.25.0 // indirect
	golang.org/x/crypto v0.13.0 // indirect
	golang.org/x/sys v0.12.0 // indirect
	golang.org/x/text v0.13.0 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20230720185612-659f7aaaa771 // indirect
	google.golang.org/grpc v1.57.0 // indirect
	google.golang.org/protobuf v1.31.0 // indirect
)
