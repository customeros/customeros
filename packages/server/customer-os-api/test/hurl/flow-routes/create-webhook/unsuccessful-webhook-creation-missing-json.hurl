# Test: Unsuccessful webhook creation missing json
POST {{cos_url}}/flows/v1/hooks
Content-Type: application/json
X-CUSTOMER-OS-API-KEY: {{api_key}}

HTTP 400
#[Captures]
#hook_url: jsonpath "$.hook.url"
#hook_secret: jsonpath "$.hook.secret"

[Asserts]
header "Content-Type" == "application/json; charset=utf-8"
jsonpath "$.status" == "error"
jsonpath "$.requestId" exists
jsonpath "$.message" == "Missing parameter: integration"
