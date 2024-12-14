# Test: Unsuccessful webhook creation invalid integration
POST {{cos_url}}/flows/v1/hooks
Content-Type: application/json
X-CUSTOMER-OS-API-KEY: {{api_key}}
{
    "integration": "slack123"
}

HTTP 400
#[Captures]
#hook_url: jsonpath "$.hook.url"
#hook_secret: jsonpath "$.hook.secret"

[Asserts]
header "Content-Type" == "application/json; charset=utf-8"
jsonpath "$.status" == "error"
jsonpath "$.requestId" exists
jsonpath "$.message" == "please provide a valid integration value"
