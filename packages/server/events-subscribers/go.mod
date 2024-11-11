module github.com/openline-ai/openline-customer-os/packages/server/events-subscribers

go 1.22.7

toolchain go1.23.1

replace github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module => ../customer-os-common-module

replace github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository => ../customer-os-postgres-repository

replace github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository => ../customer-os-neo4j-repository

replace github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto => ../events-processing-proto

replace github.com/openline-ai/openline-customer-os/packages/server/enrichment-api => ../enrichment-api

replace github.com/openline-ai/openline-customer-os/packages/server/validation-api => ../validation-api

require (
	github.com/caarlos0/env/v6 v6.10.1
	github.com/joho/godotenv v1.5.1
	github.com/neo4j/neo4j-go-driver/v5 v5.26.0
	github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module v0.0.0-20240925071556-618778a7336e
	github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository v0.0.0-20240410144729-44cbe53c019c
	github.com/openline-ai/openline-customer-os/packages/server/enrichment-api v0.0.0-00010101000000-000000000000
	github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto v0.0.0-20241004122044-3a0040d9c64c
	github.com/opentracing/opentracing-go v1.2.0
	github.com/pkg/errors v0.9.1
	golang.org/x/net v0.31.0
	gorm.io/driver/postgres v1.5.9
	gorm.io/gorm v1.25.12
)

require (
	cloud.google.com/go/auth v0.10.1 // indirect
	cloud.google.com/go/auth/oauth2adapt v0.2.5 // indirect
	cloud.google.com/go/compute/metadata v0.5.2 // indirect
	github.com/Boostport/mjml-go v0.15.0 // indirect
	github.com/PuerkitoBio/goquery v1.9.2 // indirect
	github.com/andybalholm/brotli v1.1.0 // indirect
	github.com/andybalholm/cascadia v1.3.2 // indirect
	github.com/araddon/dateparse v0.0.0-20210429162001-6b43995a97de // indirect
	github.com/aws/aws-sdk-go v1.55.3 // indirect
	github.com/bytedance/sonic v1.12.1 // indirect
	github.com/bytedance/sonic/loader v0.2.0 // indirect
	github.com/cenkalti/backoff/v4 v4.3.0 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/cloudwego/base64x v0.1.4 // indirect
	github.com/cloudwego/iasm v0.2.0 // indirect
	github.com/coocood/freecache v1.2.4 // indirect
	github.com/customeros/mailsherpa v0.2.18 // indirect
	github.com/emersion/go-message v0.18.1 // indirect
	github.com/felixge/httpsnoop v1.0.4 // indirect
	github.com/forPelevin/gomoji v1.2.0 // indirect
	github.com/gabriel-vasile/mimetype v1.4.5 // indirect
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/gin-gonic/gin v1.10.0 // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-playground/validator/v10 v10.22.1 // indirect
	github.com/goccy/go-json v0.10.3 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/google/go-querystring v1.1.0 // indirect
	github.com/google/s2a-go v0.1.8 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/googleapis/enterprise-certificate-proxy v0.3.4 // indirect
	github.com/googleapis/gax-go/v2 v2.13.0 // indirect
	github.com/h2non/filetype v1.1.3 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20231201235250-de7065d80cb9 // indirect
	github.com/jackc/pgx/v5 v5.5.5 // indirect
	github.com/jackc/puddle/v2 v2.2.1 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/cpuid/v2 v2.2.8 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	github.com/machinebox/graphql v0.2.2 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/mrz1836/postmark v1.6.5 // indirect
	github.com/novuhq/go-novu v0.1.2 // indirect
	github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-ai v0.0.0-20241028063002-b2d4aaf8096f // indirect
	github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository v0.0.0-20240920114849-ff5ab459a427 // indirect
	github.com/pelletier/go-toml/v2 v2.2.2 // indirect
	github.com/rabbitmq/amqp091-go v1.10.0 // indirect
	github.com/rivo/uniseg v0.4.7 // indirect
	github.com/sirupsen/logrus v1.9.3 // indirect
	github.com/tetratelabs/wazero v1.8.0 // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	github.com/uber/jaeger-client-go v2.30.0+incompatible // indirect
	github.com/uber/jaeger-lib v2.4.1+incompatible // indirect
	github.com/ugorji/go/codec v1.2.12 // indirect
	go.opencensus.io v0.24.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.54.0 // indirect
	go.opentelemetry.io/otel v1.30.0 // indirect
	go.opentelemetry.io/otel/metric v1.30.0 // indirect
	go.opentelemetry.io/otel/trace v1.30.0 // indirect
	go.uber.org/atomic v1.11.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.27.0 // indirect
	golang.org/x/arch v0.9.0 // indirect
	golang.org/x/crypto v0.29.0 // indirect
	golang.org/x/exp v0.0.0-20240719175910-8a7402abbf56 // indirect
	golang.org/x/oauth2 v0.24.0 // indirect
	golang.org/x/sync v0.9.0 // indirect
	golang.org/x/sys v0.27.0 // indirect
	golang.org/x/text v0.20.0 // indirect
	google.golang.org/api v0.205.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20241021214115-324edc3d5d38 // indirect
	google.golang.org/grpc v1.68.0 // indirect
	google.golang.org/protobuf v1.35.1 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
