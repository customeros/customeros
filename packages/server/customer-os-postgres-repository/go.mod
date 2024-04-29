module github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository

go 1.21

replace github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module => ../customer-os-common-module

require (
	github.com/google/uuid v1.6.0
	github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module v0.0.0-20231227140027-08b87eac2360
	github.com/opentracing/opentracing-go v1.2.0
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.9.3
	golang.org/x/net v0.24.0
	gorm.io/gorm v1.25.10
)

require (
	github.com/cenkalti/backoff/v4 v4.3.0 // indirect
	github.com/h2non/filetype v1.1.3 // indirect
	github.com/neo4j/neo4j-go-driver/v5 v5.20.0 // indirect
	github.com/uber/jaeger-client-go v2.30.0+incompatible // indirect
	go.uber.org/zap v1.27.0 // indirect
	google.golang.org/protobuf v1.33.0 // indirect
)

require (
	github.com/HdrHistogram/hdrhistogram-go v1.1.2 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/uber/jaeger-lib v2.4.1+incompatible // indirect
	go.uber.org/atomic v1.11.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/sys v0.19.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240227224415-6ceb2ff114de // indirect
	google.golang.org/grpc v1.63.2 // indirect
)
