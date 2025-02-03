# Test: Unsuccessful email verification - missing address X
GET {{cos_url}}/verify/v1/email
Content-Type: application/json
X-CUSTOMER-OS-API-KEY: {{api_key}}
[QueryStringParams]

HTTP *
[Asserts]
status == 400
header "Content-Type" contains "application/json; charset=utf-8"
jsonpath "$.status" == "error"
jsonpath "$.requestId" exists
jsonpath "$.message" == "Missing parameter: address"
