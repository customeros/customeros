# Test: Unsuccessful get bulk emails upload - unexisting jobId with valid id
GET {{cos_url}}/verify/v1/email/bulk/results/02d64272-9067-4875-9717-3a42d7077541
X-CUSTOMER-OS-API-KEY: {{api_key}}

HTTP 404
[Asserts]
header "Content-Type" contains "application/json"
jsonpath "$.requestId" exists
jsonpath "$.status" == "error"
jsonpath "$.message" == "Resource not found"
