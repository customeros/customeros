# Test: Enrich Organization Missing Required Params
GET {{cos_url}}/enrich/v1/organization
Content-Type: application/json
X-CUSTOMER-OS-API-KEY: {{api_key}}

HTTP *
[Asserts]
status == 400
header "Content-Type" contains "application/json; charset=utf-8"
jsonpath "$.status" == "error"
jsonpath "$.requestId" exists
jsonpath "$.error" == "Bad Request"
jsonpath "$.message" == "Missing linkedinUrl or domain"
