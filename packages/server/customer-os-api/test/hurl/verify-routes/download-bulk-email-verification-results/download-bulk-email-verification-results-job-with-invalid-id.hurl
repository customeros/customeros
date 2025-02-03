# Test: Download Bulk Email Verification Results - job with invalid id
GET {{cos_url}}/verify/v1/email/bulk/results/abc/download
Content-Type: application/json
X-CUSTOMER-OS-API-KEY: {{api_key}}
[Options]

HTTP *
[Asserts]
status == 404
header "Content-Type" contains "application/json"
jsonpath "$.requestId" exists
jsonpath "$.status" == "error"
jsonpath "$.message" == "Unable to find request"
