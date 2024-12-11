# Test: Successful enrich person callback
GET {{cos_url}}/enrich/v1/person
Content-Type: application/json
X-CUSTOMER-OS-API-KEY: {{api_key}}
[QueryStringParams]
linkedinUrl: https://www.linkedin.com/in/mateocafe/
includeMobileNumber: true

HTTP/2 200
[Captures]
requestId: jsonpath "$.requestId"
[Asserts]
jsonpath "$.status" == "success"
jsonpath "$.requestId" exists

GET https://api.customeros.ai/enrich/v1/person/results/{{requestId}}
Content-Type: application/json
X-CUSTOMER-OS-API-KEY: {{api_key}}

HTTP 500
[Asserts]
body == ""
