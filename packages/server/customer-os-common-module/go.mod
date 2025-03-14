module github.com/customeros/customeros/packages/server/customer-os-common-module

go 1.23.0

toolchain go1.23.1

replace github.com/customeros/customeros/packages/server/customer-os-postgres-repository => ../customer-os-postgres-repository

replace github.com/customeros/customeros/packages/server/customer-os-neo4j-repository => ../customer-os-neo4j-repository

require (
	github.com/BurntSushi/toml v1.4.0
	github.com/PuerkitoBio/goquery v1.10.2
	github.com/araddon/dateparse v0.0.0-20210429162001-6b43995a97de
	github.com/biter777/countries v1.7.5
	github.com/cenkalti/backoff/v4 v4.3.0
	github.com/cloudflare/cloudflare-go v0.115.0
	github.com/coocood/freecache v1.2.4
	github.com/customeros/customeros/packages/server/customer-os-neo4j-repository v0.0.0-20241202190357-68898897475c
	github.com/customeros/customeros/packages/server/customer-os-postgres-repository v0.0.0-20240920114849-ff5ab459a427
	github.com/customeros/mailsherpa v0.3.9
	github.com/customeros/mailwatcher v0.1.6
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/docker/docker v28.0.1+incompatible
	github.com/dustin/go-humanize v1.0.1
	github.com/emersion/go-message v0.18.2
	github.com/forPelevin/gomoji v1.3.0
	github.com/gin-gonic/gin v1.10.0
	github.com/google/generative-ai-go v0.19.0
	github.com/google/uuid v1.6.0
	github.com/h2non/filetype v1.1.3
	github.com/matoous/go-nanoid/v2 v2.1.0
	github.com/mitchellh/mapstructure v1.5.0
	github.com/mrz1836/postmark v1.7.0
	github.com/neo4j/neo4j-go-driver/v5 v5.28.0
	github.com/nyaruka/phonenumbers v1.5.0
	github.com/opensearch-project/opensearch-go/v2 v2.3.0
	github.com/opentracing/opentracing-go v1.2.0
	github.com/pkg/errors v0.9.1
	github.com/rabbitmq/amqp091-go v1.10.0
	github.com/sirupsen/logrus v1.9.3
	github.com/smartystreets/smartystreets-go-sdk v1.21.1
	github.com/stretchr/testify v1.10.0
	github.com/stripe/stripe-go/v81 v81.4.0
	github.com/testcontainers/testcontainers-go v0.35.0
	github.com/uber/jaeger-client-go v2.30.0+incompatible
	go.uber.org/zap v1.27.0
	golang.org/x/exp v0.0.0-20241204233417-43b7b7cde48d
	golang.org/x/net v0.37.0
	golang.org/x/oauth2 v0.28.0
	google.golang.org/api v0.226.0
	google.golang.org/grpc v1.71.0
	google.golang.org/protobuf v1.36.5
	gorm.io/datatypes v1.2.5
	gorm.io/driver/postgres v1.5.11
	gorm.io/gorm v1.25.12
)

require (
	cloud.google.com/go v0.116.0 // indirect
	cloud.google.com/go/ai v0.8.0 // indirect
	cloud.google.com/go/auth v0.15.0 // indirect
	cloud.google.com/go/auth/oauth2adapt v0.2.7 // indirect
	cloud.google.com/go/compute/metadata v0.6.0 // indirect
	cloud.google.com/go/longrunning v0.5.7 // indirect
	dario.cat/mergo v1.0.0 // indirect
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/Azure/go-ansiterm v0.0.0-20230124172434-306776ec8161 // indirect
	github.com/Microsoft/go-winio v0.6.2 // indirect
	github.com/andybalholm/brotli v1.1.0 // indirect
	github.com/andybalholm/cascadia v1.3.3 // indirect
	github.com/bytedance/sonic/loader v0.2.1 // indirect
	github.com/cloudwego/base64x v0.1.4 // indirect
	github.com/cloudwego/iasm v0.2.0 // indirect
	github.com/containerd/log v0.1.0 // indirect
	github.com/containerd/platforms v0.2.1 // indirect
	github.com/cpuguy83/dockercfg v0.3.2 // indirect
	github.com/distribution/reference v0.6.0 // indirect
	github.com/docker/go-connections v0.5.0 // indirect
	github.com/docker/go-units v0.5.0 // indirect
	github.com/facebookgo/clock v0.0.0-20150410010913-600d898af40a // indirect
	github.com/felixge/httpsnoop v1.0.4 // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-ole/go-ole v1.2.6 // indirect
	github.com/go-sql-driver/mysql v1.8.1 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/google/go-querystring v1.1.0 // indirect
	github.com/google/s2a-go v0.1.9 // indirect
	github.com/googleapis/enterprise-certificate-proxy v0.3.5 // indirect
	github.com/googleapis/gax-go/v2 v2.14.1 // indirect
	github.com/grpc-ecosystem/go-grpc-middleware v1.4.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.25.1 // indirect
	github.com/jackc/puddle/v2 v2.2.1 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/jpillora/backoff v1.0.0 // indirect
	github.com/klauspost/compress v1.17.9 // indirect
	github.com/lib/pq v1.10.9 // indirect
	github.com/likexian/gokit v0.25.15 // indirect
	github.com/likexian/whois v1.15.5 // indirect
	github.com/likexian/whois-parser v1.24.20 // indirect
	github.com/lucasepe/codename v0.2.0 // indirect
	github.com/lufia/plan9stats v0.0.0-20211012122336-39d0f177ccd0 // indirect
	github.com/magiconair/properties v1.8.7 // indirect
	github.com/miekg/dns v1.1.63 // indirect
	github.com/moby/docker-image-spec v1.3.1 // indirect
	github.com/moby/patternmatcher v0.6.0 // indirect
	github.com/moby/sys/sequential v0.5.0 // indirect
	github.com/moby/sys/user v0.1.0 // indirect
	github.com/moby/sys/userns v0.1.0 // indirect
	github.com/moby/term v0.5.0 // indirect
	github.com/morikuni/aec v1.0.0 // indirect
	github.com/nexus-rpc/sdk-go v0.3.0 // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.1.0 // indirect
	github.com/pborman/uuid v1.2.1 // indirect
	github.com/power-devops/perfstat v0.0.0-20210106213030-5aafc221ea8c // indirect
	github.com/rdegges/go-ipify v0.0.0-20150526035502-2d94a6a86c40 // indirect
	github.com/rivo/uniseg v0.4.7 // indirect
	github.com/robfig/cron v1.2.0 // indirect
	github.com/shirou/gopsutil/v3 v3.23.12 // indirect
	github.com/shoenig/go-m1cpu v0.1.6 // indirect
	github.com/stretchr/objx v0.5.2 // indirect
	github.com/tetratelabs/wazero v1.8.0 // indirect
	github.com/tklauser/go-sysconf v0.3.12 // indirect
	github.com/tklauser/numcpus v0.6.1 // indirect
	github.com/yusufpapurcu/wmi v1.2.3 // indirect
	go.opentelemetry.io/auto/sdk v1.1.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.59.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.59.0 // indirect
	go.opentelemetry.io/otel v1.34.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.34.0 // indirect
	go.opentelemetry.io/otel/metric v1.34.0 // indirect
	go.opentelemetry.io/otel/trace v1.34.0 // indirect
	go.opentelemetry.io/proto/otlp v1.5.0 // indirect
	go.temporal.io/api v1.44.1 // indirect
	golang.org/x/mod v0.22.0 // indirect
	golang.org/x/sync v0.12.0 // indirect
	golang.org/x/time v0.11.0 // indirect
	golang.org/x/tools v0.28.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20250115164207-1a7da9e5054f // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gorm.io/driver/mysql v1.5.6 // indirect
)

require (
	github.com/Boostport/mjml-go v0.15.0
	github.com/aws/aws-sdk-go v1.55.3
	github.com/bytedance/sonic v1.12.6 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/gabriel-vasile/mimetype v1.4.8 // indirect
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-playground/validator/v10 v10.25.0
	github.com/goccy/go-json v0.10.5 // indirect
	github.com/golang/mock v1.6.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20231201235250-de7065d80cb9 // indirect
	github.com/jackc/pgx/v5 v5.5.5 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/cpuid/v2 v2.2.9 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/novuhq/go-novu v0.1.2
	github.com/pelletier/go-toml/v2 v2.2.3 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	github.com/uber/jaeger-lib v2.4.1+incompatible // indirect
	github.com/ugorji/go/codec v1.2.12 // indirect
	go.temporal.io/sdk v1.33.0
	go.uber.org/atomic v1.11.0 // indirect
	go.uber.org/multierr v1.11.0
	golang.org/x/arch v0.12.0 // indirect
	golang.org/x/crypto v0.36.0 // indirect
	golang.org/x/sys v0.31.0 // indirect
	golang.org/x/text v0.23.0
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250303144028-a0af3efb3deb // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
