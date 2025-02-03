# Test: Authorize me with invalid key
GET {{cos_url}}/me
X-CUSTOMER-OS-API-KEY: 131a1111-1111-1111-1111-11113b4757b1

HTTP *
[Asserts]
status == 401
header "Content-Type" contains "application/json"
jsonpath "$.status" == "error"
jsonpath "$.requestId" exists
jsonpath "$.message" == "invalid API key"
