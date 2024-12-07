# Test: Enrich non-existing Person
GET https://api.customeros.ai/enrich/v1/person
Content-Type: application/json
X-CUSTOMER-OS-API-KEY: {{api_key}}
[QueryStringParams]
includeMobileNumber: true

HTTP 400
[Asserts]
header "Content-Type" contains "application/json; charset=utf-8"
jsonpath "$.status" == "error"
jsonpath "$.requestId" exists
jsonpath "$.message" == "Missing linkedinUrl or email"
