# Test: Download Bulk Email Verification Results - processing
POST {{cos_url}}/verify/v1/email/bulk
Content-Type: multipart/form-data
X-CUSTOMER-OS-API-KEY: {{api_key}}
[MultipartFormData]
file: file,download-bulk-email-verification-results-processing.csv;

HTTP 200
[Captures]
job_id: jsonpath "$.jobId"
[Asserts]
jsonpath "$.message" == "File uploaded successfully"

GET {{cos_url}}/verify/v1/email/bulk/results/{{job_id}}/download
Content-Type: application/json
X-CUSTOMER-OS-API-KEY: {{api_key}}
[Options]
#delay: 30000

HTTP 202
[Asserts]
header "Content-Type" contains "application/json"
jsonpath "$.requestId" exists
jsonpath "$.status" == "processing"
jsonpath "$.message" == "The bulk request is still being processed. Please try again later."
