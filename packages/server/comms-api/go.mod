module github.com/openline-ai/openline-customer-os/packages/server/comms-api

go 1.20

replace github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module => ./../customer-os-common-module

replace github.com/openline-ai/openline-customer-os/packages/server/customer-os-api => ./../customer-os-api

replace github.com/openline-ai/openline-customer-os/packages/server/events-processing-common => ./../events-processing-common

require (
	github.com/99designs/gqlgen v0.17.35
	github.com/DusanKasan/parsemail v1.2.0
	github.com/caarlos0/env/v6 v6.10.1
	github.com/emersion/go-message v0.16.0
	github.com/gin-contrib/cors v1.4.0
	github.com/gin-gonic/gin v1.9.1
	github.com/google/uuid v1.3.0
	github.com/gorilla/websocket v1.5.0
	github.com/joho/godotenv v1.5.1
	github.com/machinebox/graphql v0.2.2
	github.com/openline-ai/openline-customer-os/packages/server/customer-os-api v0.0.0-20230712112125-a28f2e7c190d
	github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module v0.0.0-20230624051924-41a31ae34f9c
	github.com/redis/go-redis/v9 v9.0.5
	github.com/sirupsen/logrus v1.9.3
	github.com/stretchr/testify v1.8.4
	github.com/vektah/gqlparser/v2 v2.5.8
	golang.org/x/oauth2 v0.10.0
	google.golang.org/api v0.131.0
	google.golang.org/grpc v1.56.2
	gorm.io/driver/postgres v1.5.2
	gorm.io/gorm v1.25.2
)

require (
	cloud.google.com/go/compute v1.20.1 // indirect
	cloud.google.com/go/compute/metadata v0.2.3 // indirect
	github.com/agnivade/levenshtein v1.1.1 // indirect
	github.com/bytedance/sonic v1.9.1 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/chenzhuoyu/base64x v0.0.0-20221115062448-fe3a3abad311 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/emersion/go-textwrapper v0.0.0-20200911093747-65d896831594 // indirect
	github.com/gabriel-vasile/mimetype v1.4.2 // indirect
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-playground/validator/v10 v10.14.1 // indirect
	github.com/goccy/go-json v0.10.2 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/google/s2a-go v0.1.4 // indirect
	github.com/googleapis/enterprise-certificate-proxy v0.2.5 // indirect
	github.com/googleapis/gax-go/v2 v2.12.0 // indirect
	github.com/hashicorp/golang-lru/v2 v2.0.3 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackc/pgx/v5 v5.3.1 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/cpuid/v2 v2.2.4 // indirect
	github.com/leodido/go-urn v1.2.4 // indirect
	github.com/matryer/is v1.4.1 // indirect
	github.com/mattn/go-isatty v0.0.19 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/neo4j/neo4j-go-driver/v5 v5.10.0 // indirect
	github.com/pelletier/go-toml/v2 v2.0.8 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	github.com/ugorji/go/codec v1.2.11 // indirect
	go.opencensus.io v0.24.0 // indirect
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/multierr v1.8.0 // indirect
	go.uber.org/zap v1.24.0 // indirect
	golang.org/x/arch v0.3.0 // indirect
	golang.org/x/crypto v0.11.0 // indirect
	golang.org/x/net v0.12.0 // indirect
	golang.org/x/sys v0.10.0 // indirect
	golang.org/x/text v0.11.0 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20230706204954-ccb25ca9f130 // indirect
	google.golang.org/protobuf v1.31.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
