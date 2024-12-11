# Test: Enrich non-existing Organization
GET {{cos_url}}/enrich/v1/organization
Content-Type: application/json
X-CUSTOMER-OS-API-KEY: {{api_key}}
[QueryStringParams]
linkedinUrl: https://www.linkedin.com/company/customer-osasd/
includeMobileNumber: true

HTTP 404
[Asserts]
header "Content-Type" contains "application/json; charset=utf-8"
jsonpath "$.status" == "warning"
jsonpath "$.requestId" exists
jsonpath "$.message" == "Organization not found"
