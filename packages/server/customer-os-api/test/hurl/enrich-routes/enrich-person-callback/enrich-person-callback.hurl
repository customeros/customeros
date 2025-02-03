# Test: Successful enrich person callback [BUG COS-6173]
GET {{cos_url}}/enrich/v1/person
Content-Type: application/json
X-CUSTOMER-OS-API-KEY: {{api_key}}
[QueryStringParams]
linkedinUrl: https://www.linkedin.com/in/mateocafe/
includeMobileNumber: true

HTTP/2 *
[Captures]
requestId: jsonpath "$.requestId"
[Asserts]
status == 200
jsonpath "$.status" == "success"
jsonpath "$.requestId" exists

GET {{cos_url}}/enrich/v1/person/results/{{requestId}}
Content-Type: application/json
X-CUSTOMER-OS-API-KEY: {{api_key}}

HTTP *
[Asserts]
status == 500
body == ""
