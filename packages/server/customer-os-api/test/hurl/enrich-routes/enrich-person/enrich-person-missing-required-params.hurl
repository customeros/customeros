# Test: Enrich Person Missing Required Params
GET {{cos_url}}/enrich/v1/person
Content-Type: application/json
X-CUSTOMER-OS-API-KEY: {{api_key}}
[QueryStringParams]
includeMobileNumber: true

HTTP *
[Asserts]
status == 400
header "Content-Type" contains "application/json; charset=utf-8"
jsonpath "$.status" == "error"
jsonpath "$.requestId" exists
jsonpath "$.message" == "Missing linkedin or email"
