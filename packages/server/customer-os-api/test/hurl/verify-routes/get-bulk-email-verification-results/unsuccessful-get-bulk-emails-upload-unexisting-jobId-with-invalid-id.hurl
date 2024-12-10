# Test: Unsuccessful get bulk emails upload - unexisting jobId with valid id
GET https://api.customeros.ai/verify/v1/email/bulk/results/02d64272
X-CUSTOMER-OS-API-KEY: {{api_key}}

HTTP 500
[Asserts]
header "Content-Type" contains "application/json"
jsonpath "$.requestId" exists
jsonpath "$.status" == "error"
jsonpath "$.message" == "Failed to retrieve results"
