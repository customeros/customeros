module github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data

go 1.20

require (
	github.com/caarlos0/env/v6 v6.10.1
	github.com/google/uuid v1.3.0
	github.com/jackc/pgtype v1.14.0
	github.com/joho/godotenv v1.5.1
	github.com/neo4j/neo4j-go-driver/v5 v5.8.0
	github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module v0.0.0-20230504131725-b9ae747edbd4
	github.com/openline-ai/openline-customer-os/packages/server/events-processing-common v0.0.0-20230501083326-bbe2c8ef1fbc
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.9.0
	golang.org/x/net v0.9.0
	google.golang.org/grpc v1.55.0
	gorm.io/driver/postgres v1.5.0
	gorm.io/gorm v1.25.0
	zgo.at/zcache v1.2.0
)

require (
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/jackc/pgio v1.0.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackc/pgx/v4 v4.17.2 // indirect
	github.com/jackc/pgx/v5 v5.3.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/lib/pq v1.10.7 // indirect
	golang.org/x/crypto v0.6.0 // indirect
	golang.org/x/sys v0.7.0 // indirect
	golang.org/x/text v0.9.0 // indirect
	google.golang.org/genproto v0.0.0-20230306155012-7f2fa6fef1f4 // indirect
	google.golang.org/protobuf v1.30.0 // indirect
)
