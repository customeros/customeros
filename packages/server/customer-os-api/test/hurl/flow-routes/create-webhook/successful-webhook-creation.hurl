# Test: Successful webhook creation
POST {{cos_url}}/flows/v1/hooks
Content-Type: application/json
X-CUSTOMER-OS-API-KEY: {{api_key}}
{
    "integration": "slack"
}

HTTP 201
[Asserts]
header "Content-Type" == "application/json; charset=utf-8"
jsonpath "$.status" == "success"
jsonpath "$.requestId" exists
jsonpath "$.hook.url" matches "/flows/v1/[a-zA-Z0-9-]+/i/[a-zA-Z0-9-]+"
jsonpath "$.hook.integration" == "slack"
jsonpath "$.hook.secret" matches "^[a-zA-Z0-9_-]{32,}$"
