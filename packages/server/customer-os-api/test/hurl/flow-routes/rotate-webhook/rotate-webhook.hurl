# Test: Rotate Webhook [BUG COS-6332]
POST {{cos_url}}/flows/v1/hooks
Content-Type: application/json
X-CUSTOMER-OS-API-KEY: {{api_key}}
{
    "integration": "slack"
}

HTTP 201
[Captures]
hook_url: jsonpath "$.hook.url"
integration: jsonpath "$.hook.integration"
secret: jsonpath "$.hook.secret"
tenant_id: regex "/flows/v1/([^/]+)"
integration_id: regex "/i/([a-zA-Z0-9]+)"
[Asserts]
header "Content-Type" == "application/json; charset=utf-8"
jsonpath "$.status" == "success"
jsonpath "$.requestId" exists
jsonpath "$.hook.url" matches "/flows/v1/[a-zA-Z0-9-]+/i/[a-zA-Z0-9-]+"
jsonpath "$.hook.integration" == "slack"
jsonpath "$.hook.secret" matches "^[a-zA-Z0-9_-]{32,}$"

GET {{cos_url}}/flows/v1/hooks
Content-Type: application/json
X-CUSTOMER-OS-API-KEY: {{api_key}}

HTTP 200
[Asserts]
jsonpath "$.hooks[*].url" includes "{{hook_url}}"
jsonpath "$.hooks[?(@.url == '{{hook_url}}')].active" includes true


POST {{cos_url}}/flows/v1/{{tenant_id}}/i/{{integration_id}}/rotate
Content-Type: application/json
X-CUSTOMER-OS-API-KEY: {{api_key}}

HTTP 201
[Captures]
new_hook_url: jsonpath "$.hook.url"
[Asserts]
jsonpath "$.requestId" exists
jsonpath "$.status" == "success"
jsonpath "$.hook.url" exists
jsonpath "$.hook.url" != "{{hook_url}}"
jsonpath "$.hook.integration" == "{{integration}}"
jsonpath "$.hook.secret" exists
jsonpath "$.hook.secret" != "{{secret}}"

GET {{cos_url}}/flows/v1/hooks
Content-Type: application/json
X-CUSTOMER-OS-API-KEY: {{api_key}}

HTTP 200
[Asserts]
jsonpath "$.hooks[*].url" not includes "{{hook_url}}"
jsonpath "$.hooks[*].url" includes "{{new_hook_url}}"
jsonpath "$.hooks[?(@.url == '{{new_hook_url}}')].active" includes true
