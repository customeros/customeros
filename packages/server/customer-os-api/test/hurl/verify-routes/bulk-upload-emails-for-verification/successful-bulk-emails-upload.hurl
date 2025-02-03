# Test: Successful bulk emails upload
POST {{cos_url}}/verify/v1/email/bulk
Content-Type: multipart/form-data
X-CUSTOMER-OS-API-KEY: {{api_key}}
[MultipartFormData]
file: file,successful-bulk-emails-upload.csv;

HTTP *
[Asserts]
status == 200
jsonpath "$.message" == "File uploaded successfully"
header "Content-Type" contains "application/json"
jsonpath "$.jobId" matches "^api_[a-z0-9]+$"
jsonpath "$.resultUrl" matches "^https?://.*\/verify\/v1\/email\/bulk\/results\/api_[a-z0-9]+$"
jsonpath "$.estimatedCompletionTs" isInteger
