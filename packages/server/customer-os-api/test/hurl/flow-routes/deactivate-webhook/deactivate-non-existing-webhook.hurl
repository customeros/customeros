# Test: Deactivate non-existing webhooks [Bug COS-6248]
POST {{cos_url}}/webhooks/v1/hooks
Content-Type: application/json
X-CUSTOMER-OS-API-KEY: {{api_key}}
{
    "integration": "slack"
}

HTTP *
[Captures]
hook_url: jsonpath "$.hook.url"
tenant_id: regex "/webhooks/v1/([^/]+)"
integration_id: regex "/i/([a-zA-Z0-9]+)"
[Asserts]
status == 201
header "Content-Type" == "application/json; charset=utf-8"
jsonpath "$.status" == "success"
jsonpath "$.requestId" exists
jsonpath "$.hook.url" matches "/webhooks/v1/[a-zA-Z0-9-]+/i/[a-zA-Z0-9-]+"
jsonpath "$.hook.integration" == "slack"
jsonpath "$.hook.secret" matches "^[a-zA-Z0-9_-]{32,}$"

GET {{cos_url}}/webhooks/v1/hooks
Content-Type: application/json
X-CUSTOMER-OS-API-KEY: {{api_key}}

HTTP *
[Asserts]
status == 200
jsonpath "$.hooks[*].url" includes "{{hook_url}}"
jsonpath "$.hooks[?(@.url == '{{hook_url}}')].active" includes true


DELETE {{cos_url}}/webhooks/v1/{{tenant_id}}/i/does-not-exist
Content-Type: application/json
X-CUSTOMER-OS-API-KEY: {{api_key}}

HTTP *
[Asserts]
status == 404
jsonpath "$.requestId" exists
jsonpath "$.status" == "error"
jsonpath "$.message" == "Webhook successfully deactivated"
body == ""
