module github.com/openline-ai/openline-customer-os/packages/runner/customer-os-dedup

go 1.21

replace github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module => ./../../server/customer-os-common-module

replace github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository => ./../../server/customer-os-neo4j-repository

replace github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository => ./../../server/customer-os-postgres-repository

require (
	github.com/caarlos0/env/v6 v6.10.1
	github.com/joho/godotenv v1.5.1
	github.com/machinebox/graphql v0.2.2
	github.com/neo4j/neo4j-go-driver/v5 v5.19.0
	github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module v0.0.0-20240223192554-76c272b433c9
	github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository v0.0.0-20240223192554-76c272b433c9
	github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository v0.0.0-20240410144729-44cbe53c019c
	github.com/opentracing/opentracing-go v1.2.0
	github.com/pkg/errors v0.9.1
	github.com/robfig/cron v1.2.0
	github.com/stretchr/testify v1.9.0
	gorm.io/gorm v1.25.9
)

require (
	github.com/cenkalti/backoff/v4 v4.3.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/h2non/filetype v1.1.3 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackc/pgx/v5 v5.5.4 // indirect
	github.com/jackc/puddle/v2 v2.2.1 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/matryer/is v1.4.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/sirupsen/logrus v1.9.3 // indirect
	github.com/uber/jaeger-client-go v2.30.0+incompatible // indirect
	github.com/uber/jaeger-lib v2.4.1+incompatible // indirect
	go.uber.org/atomic v1.11.0 // indirect
	go.uber.org/multierr v1.10.0 // indirect
	go.uber.org/zap v1.27.0 // indirect
	golang.org/x/crypto v0.22.0 // indirect
	golang.org/x/net v0.24.0 // indirect
	golang.org/x/sync v0.6.0 // indirect
	golang.org/x/sys v0.19.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
	google.golang.org/protobuf v1.33.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	gorm.io/driver/postgres v1.5.7 // indirect
)
