# Test: Download Bulk Email Verification Results - job with valid id
GET {{cos_url}}/verify/v1/email/bulk/results/68c2c6e7-cf25-41e2-abee-882d05ae62f1/download
Content-Type: application/json
X-CUSTOMER-OS-API-KEY: {{api_key}}
[Options]

HTTP 404
[Asserts]
header "Content-Type" contains "application/json"
jsonpath "$.requestId" exists
jsonpath "$.status" == "error"
jsonpath "$.message" == "Unable to find request"
