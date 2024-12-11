# Test: Authorize me with invalid key
GET https://api.customeros.ai/me
X-CUSTOMER-OS-API-KEY: 131a1111-1111-1111-1111-11113b4757b1

HTTP 401
[Asserts]
header "Content-Type" contains "application/json"
jsonpath "$.errors[0].message" == "Invalid api key"
