# Test: Unsuccessful webhook creation invalid integration
POST {{cos_url}}/webhooks/v1/hooks
Content-Type: application/json
X-CUSTOMER-OS-API-KEY: {{api_key}}
{
    "integration": "slack123"
}

HTTP *
[Asserts]
status == 400
header "Content-Type" == "application/json; charset=utf-8"
jsonpath "$.status" == "error"
jsonpath "$.requestId" exists
jsonpath "$.message" == "Invalid integration value"
