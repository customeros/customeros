# Test: Authorize me successfully
GET https://api.customeros.ai/me
X-CUSTOMER-OS-API-KEY:

HTTP 401
[Asserts]
header "Content-Type" contains "application/json"
jsonpath "$.errors[0].message" == "Api key is required"
