# Test: Unsuccessful email verification - missing address
GET {{cos_url}}/verify/v1/email
Content-Type: application/json
X-CUSTOMER-OS-API-KEY: {{api_key}}
[QueryStringParams]

HTTP 400
[Asserts]
header "Content-Type" contains "application/json; charset=utf-8"
jsonpath "$.status" == "error"
jsonpath "$.requestId" exists
jsonpath "$.message" == "Missing parameter: address"
