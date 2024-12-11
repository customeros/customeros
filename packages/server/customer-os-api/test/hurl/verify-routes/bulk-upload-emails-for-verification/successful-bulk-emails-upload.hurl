# Test: Successful bulk emails upload
POST {{cos_url}}/verify/v1/email/bulk
Content-Type: multipart/form-data
X-CUSTOMER-OS-API-KEY: {{api_key}}
[MultipartFormData]
file: file,successful-bulk-emails-upload.csv;

HTTP 200
[Asserts]
jsonpath "$.message" == "File uploaded successfully"
header "Content-Type" contains "application/json"
jsonpath "$.jobId" matches "^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"
jsonpath "$.resultUrl" matches "^https?://.*\/verify\/v1\/email\/bulk\/results\/[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"
jsonpath "$.estimatedCompletionTs" isInteger
