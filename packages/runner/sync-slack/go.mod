module github.com/openline-ai/openline-customer-os/packages/runner/sync-slack

go 1.23

toolchain go1.23.1

replace github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto => ../../server/events-processing-proto

replace github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module => ./../../server/customer-os-common-module

replace github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository => ./../../server/customer-os-neo4j-repository

replace github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository => ./../../server/customer-os-postgres-repository

require (
	github.com/caarlos0/env/v6 v6.10.1
	github.com/coocood/freecache v1.2.4
	github.com/google/uuid v1.6.0
	github.com/joho/godotenv v1.5.1
	github.com/machinebox/graphql v0.2.2
	github.com/neo4j/neo4j-go-driver/v5 v5.27.0
	github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module v0.0.0-20240925071556-618778a7336e
	github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository v0.0.0-20240920114849-ff5ab459a427
	github.com/opentracing/opentracing-go v1.2.0
	github.com/pkg/errors v0.9.1
	github.com/robfig/cron v1.2.0
	github.com/sirupsen/logrus v1.9.3
	github.com/slack-go/slack v0.15.0
	github.com/stretchr/testify v1.10.0
	gorm.io/driver/postgres v1.5.11
	gorm.io/gorm v1.25.12
	zgo.at/zcache v1.2.0
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/bytedance/sonic v1.12.6 // indirect
	github.com/bytedance/sonic/loader v0.2.1 // indirect
	github.com/cenkalti/backoff/v4 v4.3.0 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/cloudwego/base64x v0.1.4 // indirect
	github.com/cloudwego/iasm v0.2.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/forPelevin/gomoji v1.2.0 // indirect
	github.com/gabriel-vasile/mimetype v1.4.8 // indirect
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/gin-gonic/gin v1.10.0 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-playground/validator/v10 v10.24.0 // indirect
	github.com/go-sql-driver/mysql v1.8.1 // indirect
	github.com/goccy/go-json v0.10.4 // indirect
	github.com/gorilla/websocket v1.4.2 // indirect
	github.com/h2non/filetype v1.1.3 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20231201235250-de7065d80cb9 // indirect
	github.com/jackc/pgx/v5 v5.5.5 // indirect
	github.com/jackc/puddle/v2 v2.2.1 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/cpuid/v2 v2.2.9 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	github.com/matoous/go-nanoid/v2 v2.1.0 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/pelletier/go-toml/v2 v2.2.3 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rivo/uniseg v0.4.7 // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	github.com/uber/jaeger-client-go v2.30.0+incompatible // indirect
	github.com/uber/jaeger-lib v2.4.1+incompatible // indirect
	github.com/ugorji/go/codec v1.2.12 // indirect
	go.uber.org/atomic v1.11.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.27.0 // indirect
	golang.org/x/arch v0.12.0 // indirect
	golang.org/x/crypto v0.32.0 // indirect
	golang.org/x/net v0.34.0 // indirect
	golang.org/x/sync v0.10.0 // indirect
	golang.org/x/sys v0.29.0 // indirect
	golang.org/x/text v0.21.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20241118233622-e639e219e697 // indirect
	google.golang.org/grpc v1.69.4 // indirect
	google.golang.org/protobuf v1.36.2 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	gorm.io/datatypes v1.2.5 // indirect
	gorm.io/driver/mysql v1.5.6 // indirect
)
