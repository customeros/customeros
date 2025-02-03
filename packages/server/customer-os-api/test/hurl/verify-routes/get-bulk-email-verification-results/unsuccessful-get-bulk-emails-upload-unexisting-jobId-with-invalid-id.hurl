# Test: Unsuccessful get bulk emails upload - unexisting jobId with invalid id
GET {{cos_url}}/verify/v1/email/bulk/results/02d64272
X-CUSTOMER-OS-API-KEY: {{api_key}}

HTTP *
[Asserts]
status == 404
header "Content-Type" contains "application/json"
jsonpath "$.requestId" exists
jsonpath "$.status" == "error"
jsonpath "$.error" == "Not Found"
