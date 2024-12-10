# Test: Authorize me successfully
GET https://api.customeros.ai/me
X-CUSTOMER-OS-API-KEY: {{api_key}}

HTTP 200
[Asserts]
header "Content-Type" contains "application/json"
jsonpath "$.status" == "success"
jsonpath "$.requestId" exists
jsonpath "$.tenant" == "PureThePhantom"
