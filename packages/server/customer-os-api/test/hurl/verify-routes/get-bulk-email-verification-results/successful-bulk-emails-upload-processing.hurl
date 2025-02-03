# Test: Successful get bulk emails upload processing
POST {{cos_url}}/verify/v1/email/bulk
Content-Type: multipart/form-data
X-CUSTOMER-OS-API-KEY: {{api_key}}
[MultipartFormData]
file: file,successful-bulk-emails-upload-processing.csv;

HTTP *
[Captures]
job_id: jsonpath "$.jobId"
[Asserts]
status == 200
jsonpath "$.message" == "File uploaded successfully"

GET {{cos_url}}/verify/v1/email/bulk/results/{{job_id}}
Content-Type: application/json
X-CUSTOMER-OS-API-KEY: {{api_key}}

HTTP *
[Asserts]
status == 200
header "Content-Type" contains "application/json"
jsonpath "$.jobId" == "{{job_id}}"
jsonpath "$.status" == "processing"
jsonpath "$.fileName" == "successful-bulk-emails-upload-processing.csv"
jsonpath "$.message" == "Completed 0 of 2 emails"
jsonpath "$.estimatedCompletionTs" exists
