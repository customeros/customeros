# Test: Successful get bulk emails upload completed
POST {{cos_url}}/verify/v1/email/bulk
Content-Type: multipart/form-data
X-CUSTOMER-OS-API-KEY: {{api_key}}
[MultipartFormData]
file: file,successful-bulk-emails-upload-completed.csv;

HTTP 200
[Captures]
job_id: jsonpath "$.jobId"
[Asserts]
jsonpath "$.message" == "File uploaded successfully"

GET {{cos_url}}/verify/v1/email/bulk/results/{{job_id}}
Content-Type: application/json
X-CUSTOMER-OS-API-KEY: {{api_key}}
[Options]
delay: 30000

HTTP 200
[Asserts]
header "Content-Type" contains "application/json"
jsonpath "$.jobId" == "{{job_id}}"
jsonpath "$.status" == "completed"
jsonpath "$.fileName" == "successful-bulk-emails-upload-completed.csv"
jsonpath "$.message" == ""
jsonpath "$.results.totalEmails" == 2
jsonpath "$.results.deliverable" == 0
jsonpath "$.results.undeliverable" == 2
jsonpath "$.results.downloadUrl" == "{{cos_url}}/verify/v1/email/bulk/results/{{job_id}}/download"
jsonpath "$.estimatedCompletionTs" == 0
