cd server
cd customer-os-postgres-repository
go mod tidy
cd ../customer-os-common-module
go mod tidy
cd ../events-subscribers
go mod tidy
cd ../customer-os-api
go mod tidy
cd ../customer-os-neo4j-repository
go mod tidy
cd ../customer-os-webhooks
go mod tidy
cd ../mailsherpa-api
go mod tidy
cd ../apigator
go mod tidy
cd ../core-crm
go mod tidy
cd ../ai
go mod tidy
cd ../..
cd runner
cd customer-os-data-upkeeper
go mod tidy
cd ../integrity-checker
go mod tidy
cd ../..
