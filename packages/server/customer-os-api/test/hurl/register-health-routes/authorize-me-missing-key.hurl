# Test: Authorize me with missing key
GET {{cos_url}}/me
X-CUSTOMER-OS-API-KEY:

HTTP 401
[Asserts]
header "Content-Type" contains "application/json"
jsonpath "$.errors[0].message" == "Api key is required"
