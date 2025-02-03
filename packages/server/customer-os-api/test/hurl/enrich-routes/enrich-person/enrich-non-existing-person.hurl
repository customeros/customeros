# Test: Enrich non-existing Person [BUG COS-6663]
GET {{cos_url}}/enrich/v1/person
Content-Type: application/json
X-CUSTOMER-OS-API-KEY: {{api_key}}
[QueryStringParams]
linkedinUrl: https://www.linkedin.com/in/asdqwezxc11aqz1/
includeMobileNumber: true

HTTP *
[Asserts]
status == 404
header "Content-Type" contains "application/json; charset=utf-8"
jsonpath "$.status" == "warning"
jsonpath "$.requestId" exists
jsonpath "$.message" == "Person not found"
body == ""
