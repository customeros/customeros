cd server
cd customer-os-postgres-repository
go mod tidy
cd ../customer-os-common-module
go mod tidy
cd ../anthropic-api
go mod tidy
cd ../comms-api
go mod tidy
cd ../customer-os-analytics-api
go mod tidy
cd ../customer-os-api
go mod tidy
cd ../customer-os-common-auth
go mod tidy
cd ../customer-os-neo4j-repository
go mod tidy
cd ../customer-os-platform-admin-api
go mod tidy
cd ../customer-os-webhooks
go mod tidy
cd ../events-processing-platform
go mod tidy
cd ../events-processing-platform-subscribers
go mod tidy
cd ../file-store-api
go mod tidy
cd ../openai-api
go mod tidy
cd ../settings-api
go mod tidy
cd ../user-admin-api
go mod tidy
cd ../validation-api
go mod tidy
cd ../..
cd runner
cd customer-os-data-upkeeper
go mod tidy
cd ../sync-gmail-raw
go mod tidy
cd ../sync-gmail
go mod tidy
cd ../sync-customer-os-data
go mod tidy
cd ../customer-os-dedup
go mod tidy
cd ../..
