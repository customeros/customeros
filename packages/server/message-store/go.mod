module github.com/openline-ai/openline-customer-os/packages/server/message-store

go 1.19

require (
	entgo.io/ent v0.11.4-0.20221001062602-1029a2d3ba2a
	github.com/99designs/gqlgen v0.17.20
	github.com/caarlos0/env/v6 v6.10.1
	github.com/lib/pq v1.10.7
	github.com/machinebox/graphql v0.2.2
	github.com/openline-ai/openline-customer-os/customer-os-api v0.0.0-00010101000000-000000000000
	github.com/vektah/gqlparser/v2 v2.5.1
	google.golang.org/grpc v1.50.1
	google.golang.org/protobuf v1.28.1
)

require (
	ariga.io/atlas v0.8.1 // indirect
	github.com/agext/levenshtein v1.2.3 // indirect
	github.com/agnivade/levenshtein v1.1.1 // indirect
	github.com/apparentlymart/go-textseg/v13 v13.0.0 // indirect
	github.com/go-openapi/inflect v0.19.0 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/go-cmp v0.5.9 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/hashicorp/hcl/v2 v2.14.1 // indirect
	github.com/matryer/is v1.4.0 // indirect
	github.com/mitchellh/go-wordwrap v1.0.1 // indirect
	github.com/neo4j/neo4j-go-driver/v4 v4.4.4 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/zclconf/go-cty v1.11.1 // indirect
	golang.org/x/mod v0.7.0 // indirect
	golang.org/x/net v0.2.0 // indirect
	golang.org/x/sys v0.2.0 // indirect
	golang.org/x/text v0.4.0 // indirect
	google.golang.org/genproto v0.0.0-20220617124728-180714bec0ad // indirect
)

replace github.com/openline-ai/openline-customer-os/customer-os-api => ./../customer-os-api
